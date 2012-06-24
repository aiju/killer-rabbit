package main

import (
	"net"
	"bytes"
	"encoding/binary"
	"time"
)

const (
	MaxMessage = 65536
	UpdateInterval = 100 * time.Millisecond
)

const (
	TJoin = iota
	TUpdate
	RAck
)

type Client struct {
	Id [4]byte
	Addr *net.UDPAddr
}

type MessageHeader struct {
	Id [4]byte
	Seq [4]byte
	Type uint8
}

type Message struct {
	MessageHeader
	Aux []byte
	Addr *net.UDPAddr
}

type Server struct {
	Conn *net.UDPConn
	Clients map[[4]byte] *Client
	St *State
}

func (s *Server) Recv() (m *Message) {
	for {
		buf := make([]byte, MaxMessage)
		m = new(Message)
		n, addr, err := s.Conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		bufb := bytes.NewBuffer(buf[:n])
		binary.Read(bufb, binary.LittleEndian, &m.MessageHeader)
		m.Aux = buf[bufb.Len():n]
		m.Addr = addr
		return
	}
	return nil
}

func (s *Server) Send(m *Message){
	buf := make([]byte, 0, MaxMessage)
	bufb := bytes.NewBuffer(buf)
	binary.Write(bufb, binary.LittleEndian, &m.MessageHeader)
	bufb.Write(m.Aux)
	s.Conn.WriteToUDP(buf, m.Addr)
}

func (s *Server) Ack(m *Message) {
	s.Send(&Message{MessageHeader{Id: m.Id, Seq: m.Seq, Type: RAck}, nil, m.Addr})
}

func (s *Server) Handle(m *Message) {
	switch m.Type {
	case TJoin:
		s.Ack(m)
	}
}

func (s *Server) Update(c *Client) {
	
}

func StartServer(port int) error {
	var err error

	s := new(Server)
	s.St = new(State)
	s.Conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		return err
	}
	s.Clients = make(map[[4]byte] *Client)
	in := make(chan *Message)
	go func() {
		for {
			in <- s.Recv()
		}
	} ()
	ticks := time.Tick(UpdateInterval)
	for {
		select {
		case m := <-in:
			s.Handle(m)
		case <-ticks:
			s.St.Advance()
		}
	}
	return nil
}
