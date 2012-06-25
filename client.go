package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	Id         uint32
	CSeq, SSeq uint32
	Conn       *net.UDPConn
}

func (c *Client) Recv() (m *Message, err error) {
	buf := make([]byte, MaxMessage)
	m = new(Message)
	n, err := c.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	bufb := bytes.NewBuffer(buf[:n])
	binary.Read(bufb, binary.LittleEndian, &m.MessageHeader)
	m.Aux = bufb.Bytes()
	return m, nil
}

func (c *Client) Send(m *Message) {
	bufb := new(bytes.Buffer)
	binary.Write(bufb, binary.LittleEndian, &m.MessageHeader)
	bufb.Write(m.Aux)
	c.Conn.Write(bufb.Bytes())
}

func StartClient(addr *net.UDPAddr, update chan<- *State) error {
	var err error
	var m *Message
	var i int

	c := new(Client)
	c.Id = rand.Uint32()
	c.Conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	for i = 0; i < JoinRetry; i++ {
		c.Send(&Message{MessageHeader{c.Id, c.CSeq, TJoin}, nil, nil})
		c.Conn.SetReadDeadline(time.Now().Add(AckTimeout))
		m, err = c.Recv()
		if err != nil {
			if err.(net.Error).Timeout() {
				continue
			}
			return err
		}
		if m.Type == RAck && m.Id == c.Id && m.Seq == c.CSeq {
			break
		}
	}
	if i == JoinRetry {
		return EConnect
	}
	c.Conn.SetReadDeadline(time.Time{})
	for {
		m, err = c.Recv()
		if err != nil {
			return err
		}
		switch m.Type {
		case TUpdate:
			update <- StateDecode(m.Aux)
		}
	}
	return nil
}
