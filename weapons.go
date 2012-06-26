package main

import (
	"encoding/gob"
	"math"
	"time"
)

type Weapon interface {
	Fire(s *State, p *Player)
}

type Pistol struct {
	time.Time
}

func (w *Pistol) Fire(s *State, p *Player) {
	if !w.Time.IsZero() && w.Time.After(time.Now().Add(-PistolCooldown)) {
		return
	}
	r := math.Atan2(float64(p.FY)-float64(p.Y)/10, float64(p.FX)-float64(p.X)/10) * 180 / math.Pi
	e := Ent{p.X, p.Y, 1000, int16(r+720) % 360}
	s.E = append(s.E, &e)
	w.Time = time.Now()
}

func init() {
	gob.Register(new(Pistol))
}
