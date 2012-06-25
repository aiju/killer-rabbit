package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

const (
	TJoin = iota
	TUpdate
	RAck
)

type SClient struct {
	Id         uint32
	Addr       *net.UDPAddr
	SSeq, CSeq uint32
}

type MessageHeader struct {
	Id   uint32
	Seq  uint32
	Type uint8
}

type Message struct {
	MessageHeader
	Aux  []byte
	Addr *net.UDPAddr
}

type Server struct {
	Conn    *net.UDPConn
	Clients map[uint32]*SClient
	St      *State
}

func (s *Server) Recv() (m *Message, err error) {
	buf := make([]byte, MaxMessage)
	m = new(Message)
	n, addr, err := s.Conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}
	bufb := bytes.NewBuffer(buf[:n])
	binary.Read(bufb, binary.LittleEndian, &m.MessageHeader)
	m.Aux = bufb.Bytes()
	m.Addr = addr
	return m, nil
}

func (s *Server) Send(m *Message) {
	bufb := new(bytes.Buffer)
	binary.Write(bufb, binary.LittleEndian, &m.MessageHeader)
	bufb.Write(m.Aux)
	s.Conn.WriteToUDP(bufb.Bytes(), m.Addr)
}

func (s *Server) Ack(m *Message) {
	s.Send(&Message{MessageHeader{Id: m.Id, Seq: m.Seq, Type: RAck}, nil, m.Addr})
}

func (s *Server) Handle(m *Message) {
	if m.Type == TJoin {
		c := new(SClient)
		c.Id = m.Id
		c.CSeq = m.Seq
		c.Addr = m.Addr
		s.St.E = append(s.St.E, Ent{300, 300, 100, 0})
		s.Clients[m.Id] = c
		s.Ack(m)
		return
	}
	c := s.Clients[m.Id]
	if c == nil {
		return
	}
	c.CSeq = m.Seq
	c.Addr = m.Addr
}

func StartServer(port int) error {
	var err error

	s := new(Server)
	s.St = new(State)
	s.Conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		return err
	}
	s.Clients = make(map[uint32]*SClient)
	in := make(chan *Message)
	errch := make(chan error)
	go func() {
		for {
			m, err := s.Recv()
			if err != nil {
				errch <- err
				return
			}
			in <- m
		}
	}()
	ticks := time.Tick(UpdateInterval)
	for {
		select {
		case err := <-errch:
			return err
		case m := <-in:
			s.Handle(m)
		case <-ticks:
			s.St.Advance()
			p := s.St.Serialize()
			for _, c := range s.Clients {
				if c.Addr != nil {
					s.Send(&Message{MessageHeader{c.Id, c.SSeq, TUpdate}, p, c.Addr})
					c.SSeq++
				}
			}
		}
	}
	return nil
}
