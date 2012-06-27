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
	Id     uint32
	FX, FY int16
}

type State struct {
	Ents    []*Ent
	Players []*Player
}

type MoveMsg struct {
	Moving, FireStart bool
	D                 int16
	FX, FY            int16
}

func (e *Ent) Move(t time.Duration) {
	e.X += int16(int64(e.V) * (int64(t) / 1000) * int64(Cos[e.D]) / 1e9)
	e.Y += int16(int64(e.V) * (int64(t) / 1000) * int64(Sin[e.D]) / 1e9)
}

func (s *State) Advance(t time.Duration) {
	for _, e := range s.Ents {
		e.Move(t)
	}
	for _, e := range s.Players {
		e.Move(t)
	}
}

func (s State) Client(id uint32) *Player {
	for _, e := range s.Players {
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
	s.Players = append(s.Players, p)
	return p
}

func (s *State) RemovePlayer(id uint32) {
	for i, e := range s.Players {
		if e.Id == id {
			s.Players = append(s.Players[:i], s.Players[i+1:]...)
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

func (s *State) Copy() *State {
	return StateDecode(s.Serialize())
}

func (e MoveMsg) Process(s *State, p *Player) {
		if e.Moving {
			p.V = 100
		} else {
			p.V = 0
		}
		if e.FireStart {
			p.Weapon.Fire(s, p)
		}
		p.D = e.D
		p.FX = e.FX
		p.FY = e.FY
}
