package main

import (
	"bytes"
	"encoding/gob"
	"time"
	"math"
)

type Rect struct {
	cx, cy, r, d float64
}

type Ent interface {
	Move(t time.Duration)
	Draw()
	Hit() bool
	Bounds() *Rect
}

type Player struct {
	X, Y, V, D float64
	Weapon
	Id     uint32
	FX, FY int16
}

type Bullet struct {
	X, Y, V, D float64
}

type State struct {
	Ents    []Ent
	Bullets []*Bullet
	Players []*Player
}

type MoveMsg struct {
	Moving, FireStart bool
	D                 int16
	FX, FY            int16
}

func BasicMove(t time.Duration, X, Y, V, D float64) (float64, float64) {
	s := float64(t)/1e9 * V
	return X + math.Cos(D * deg) * s , Y + math.Sin(D * deg) * s
}

func (s *State) Advance(t time.Duration) {
	for _, e := range s.Ents {
		e.Move(t)
	}
	for i := 0; i < len(s.Bullets); i++ {
		e := s.Bullets[i]
		e.Move(t)
		for k := 0; k < len(s.Ents); k++ {
			f := s.Ents[k]
			r := f.Bounds()
			if r != nil && e.Hit(*r) {
				if f.Hit() {
					s.Ents = append(s.Ents[:k], s.Ents[k+1:]...)
					s.Bullets = append(s.Bullets[:i], s.Bullets[i+1:]...)
					i--
					goto next
				}
			}
		}
	next:
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
	p.X, p.Y = 200, 200
	p.Weapon = new(Pistol)
	p.Id = id
	s.Players = append(s.Players, p)
	s.Ents = append(s.Ents, &Monster{400, 200, 50, 0})
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
		p.D = float64(e.D)
		p.FX = e.FX
		p.FY = e.FY
}

func (e MoveMsg) Encode() []byte {
	bufb := new(bytes.Buffer)
	enc := gob.NewEncoder(bufb)
	vital(enc.Encode(e))
	return bufb.Bytes()
}

func MoveMsgDecode(buf []byte) (m MoveMsg) {
	bufb := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(bufb)
	vital(dec.Decode(&m))
	return
}

func (p *Player) Move(t time.Duration) {
	p.X, p.Y = BasicMove(t, p.X, p.Y, p.V, p.D)
}

func (b *Bullet) Move(t time.Duration) {
	b.X, b.Y = BasicMove(t, b.X, b.Y, b.V, b.D)
}

func (p *Player) R() float64 {
	return math.Atan2(float64(p.FY)-p.Y, float64(p.FX)-p.X) * 180 / math.Pi
}

func (b *Bullet) Hit(r Rect) bool {
	x := b.X - r.cx
	y := b.Y - r.cy
	R := r.r
	return x < R && x > -R && y < R && y > -R
}
