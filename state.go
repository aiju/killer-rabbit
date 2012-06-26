package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Ent struct {
	X, Y int16 // coordinates (1/10 pixels)
	V    int16 // magnitude of velocity in pixels per second
	D    int16 // direction in degrees, 0 <= D < 360
}

type Player struct {
	Ent
	Weapon
	Id             uint32
	MV, MD, FX, FY int16
}

type State struct {
	E []*Ent
	P []*Player
}

func (e *Ent) Move(t time.Duration) {
	e.X += int16(int64(e.V) * (int64(t) / 1000) * int64(Cos[e.D]) / 1e9)
	e.Y += int16(int64(e.V) * (int64(t) / 1000) * int64(Sin[e.D]) / 1e9)
}

func (s *State) Advance(t time.Duration) {
	for _, e := range s.E {
		e.Move(t)
	}
	for _, e := range s.P {
		e.Move(t)
	}
}

func (s State) Client(id uint32) *Player {
	for _, e := range s.P {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func (s *State) Spawn(id uint32) *Player {
	p := new(Player)
	p.Ent = Ent{3000, 3000, 0, 0}
	p.Weapon = new(Pistol)
	p.Id = id
	s.P = append(s.P, p)
	return p
}

func (s *State) RemovePlayer(id uint32) {
	for i, e := range s.P {
		if e.Id == id {
			s.P = append(s.P[:i], s.P[i+1:]...)
			return
		}
	}
}

func (s *State) Serialize() []byte {
	bufb := new(bytes.Buffer)
	enc := gob.NewEncoder(bufb)
	vital(enc.Encode(s))
	return bufb.Bytes()
}

func StateDecode(buf []byte) *State {
	s := new(State)
	bufb := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(bufb)
	vital(dec.Decode(s))
	return s
}
