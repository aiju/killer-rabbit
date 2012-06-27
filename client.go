package main

import (
	"bytes"
	"encoding/binary"
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

func StartClient(addr *net.UDPAddr, id uint32, update chan<- *State, move <-chan MoveMsg, quit chan bool) error {
	var err error
	var m *Message
	var i int

	c := new(Client)
	c.Id = id
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
	c.CSeq++
	c.Conn.SetReadDeadline(time.Time{})
	end := 0
	tick := time.Tick(PingInterval)
	go func() {
		for {
			select {
			case e := <-move:
				c.Send(&Message{MessageHeader{c.Id, c.CSeq, TMove}, e.Encode(), nil})
				c.CSeq++
			case <-quit:
				c.Send(&Message{MessageHeader{c.Id, c.CSeq, TQuit}, nil, nil})
				c.Conn.Close()
				end = 1
			case <-tick:
				c.Send(&Message{MessageHeader{c.Id, c.CSeq, TPing}, nil, nil})
				c.CSeq++
			}
		}
	}()
	for {
		m, err = c.Recv()
		if end == 1 {
			return nil
		}
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
