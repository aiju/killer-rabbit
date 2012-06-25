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

type State struct {
	E []Ent
}

func (e Ent) Move(t time.Duration) (f Ent) {
	f = e
	f.X += int16(int64(e.V) * (int64(t) / 1000) * int64(Cos[e.D]) / 1e10)
	f.Y += int16(int64(e.V) * (int64(t) / 1000) * int64(Sin[e.D]) / 1e10)
	return
}

func (s *State) Advance() {
	for i := range s.E {
		s.E[i] = s.E[i].Move(UpdateInterval)
		s.E[i].D = (s.E[i].D + 10) % 360
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
