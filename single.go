package main

import "time"

func StartSingle(id uint32, update chan *State, move chan MoveMsg, quit chan bool) error {
	s := new(State)
	p := s.Spawn(id)
	tick := time.Tick(UpdateInterval)
	for {
		select {
		case <-tick:
			update <- s.Copy()
			s.Advance(UpdateInterval)
		case m := <-move:
			m.Process(s, p)
		case <- quit:
			return nil
		}
	}
	return nil
}
