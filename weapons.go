package main

import (
	"encoding/gob"
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
	s.Bullets = append(s.Bullets, &Bullet{p.X, p.Y, 300, p.R()})
	w.Time = time.Now()
}

func init() {
	gob.Register(new(Pistol))
}
