package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

var dir = [16]int16{
	-1, 90, 180, 135, 270, -1, 225, 180, 0, 45, -1, 90, 315, 0, 270, -1,
}

func ProcessInput(move chan MoveMsg, quit chan bool, squit chan bool) {
	var movem MoveMsg
	var moving, oldmoving int

	for {
		select {
		case ev := <-sdl.Events:
			if _, ok := ev.(sdl.QuitEvent); ok {
				quit <- true
				squit <- true
				return
			}
			if k, ok := ev.(sdl.KeyboardEvent); ok {
				switch k.Type {
				case sdl.KEYDOWN:
					switch k.Keysym.Sym {
					case 's':
						moving &= ^4
						moving |= 1
					case 'a':
						moving &= ^8
						moving |= 2
					case 'w':
						moving &= ^1
						moving |= 4
					case 'd':
						moving &= ^2
						moving |= 8
					}
				case sdl.KEYUP:
					switch k.Keysym.Sym {
					case 's':
						moving &= ^1
					case 'a':
						moving &= ^2
					case 'w':
						moving &= ^4
					case 'd':
						moving &= ^8
					}
				}
				if moving != oldmoving {
					movem.Moving = dir[moving] >= 0
					if dir[moving] >= 0 {
						movem.D = dir[moving]
					}
					move <- movem
					oldmoving = moving
				}
			}
			if k, ok := ev.(sdl.MouseMotionEvent); ok {
				movem.FX = int16(k.X)
				movem.FY = int16(k.Y)
				move <- movem
			}
			if k, ok := ev.(sdl.MouseButtonEvent); ok {
				if k.Button == sdl.BUTTON_LEFT && k.State == sdl.PRESSED {
					movem.FireStart = true
					move <- movem
					movem.FireStart = false
				}
			}

		}
	}
}

