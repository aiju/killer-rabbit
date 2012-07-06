package main

import (
	"time"
	"math/rand"
	"encoding/gob"
)

type Monster struct {
	X, Y, V, D float64
}

func (m *Monster) Draw() {
	Quad(m.X, m.Y, MonsterSize, m.D, "tux")
}

func (m *Monster) Move(t time.Duration) {
	m.X, m.Y = BasicMove(t, m.X, m.Y, m.V, m.D)
	m.D += rand.Float64() * float64(t) / 1e8
}

func (m *Monster) Bounds() *Rect {
	return &Rect{m.X, m.Y, MonsterSize, m.D}
}

func (m *Monster) Hit() bool {
	return true
}

func init() {
	gob.Register(&Monster{})
}
