package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

const (
	TJoin = iota
	TQuit
	TMove
	TUpdate
	TPing
	RPong
	RAck
)

type SClient struct {
	Id         uint32
	Addr       *net.UDPAddr
	SSeq, CSeq uint32
	*Player
	LastPing time.Time
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
	*State
	LastUp time.Time
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
		c.LastPing = time.Now()
		c.Player = s.State.Spawn(c.Id)
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
	c.LastPing = time.Now()
	auxb := bytes.NewBuffer(m.Aux)
	switch m.Type {
	case TQuit:
		s.State.RemovePlayer(m.Id)
		delete(s.Clients, m.Id)
	case TMove:
		var e MoveMsg
		binary.Read(auxb, binary.LittleEndian, &e)
		e.Process(s.State, c.Player)
	}
}

func StartServer(port int) error {
	var err error

	s := new(Server)
	s.State = new(State)
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
		case s.LastUp = <-ticks:
			for _, c := range s.Clients {
				if c.LastPing.Before(time.Now().Add(-PingTimeout)) {
					s.Send(&Message{MessageHeader{c.Id, c.SSeq, TQuit}, nil, c.Addr})
					s.State.RemovePlayer(c.Id)
					delete(s.Clients, c.Id)
				}
			}
			p := s.State.Serialize()
			for _, c := range s.Clients {
				if c.Addr != nil {
					s.Send(&Message{MessageHeader{c.Id, c.SSeq, TUpdate}, p, c.Addr})
					c.SSeq++
				}
			}
			s.State.Advance(UpdateInterval)

		}
	}
	return nil
}
