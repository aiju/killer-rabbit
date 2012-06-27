package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	gl "github.com/chsc/gogl/gl21"
	"math"
	"time"
)

func SetupGL() {
	gl.ClearColor(1, 1, 1, 1)
	gl.Viewport(0, 0, 800, 600)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, 800, 600, 0, 0.01, 100)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func Quad(x, y float64) {
	gl.Begin(gl.QUADS)
	gl.Color4f(1, 0, 0, 1)
	gl.Vertex2f(gl.Float(x+2), gl.Float(y+2))
	gl.Vertex2f(gl.Float(x-2), gl.Float(y+2))
	gl.Vertex2f(gl.Float(x-2), gl.Float(y-2))
	gl.Vertex2f(gl.Float(x+2), gl.Float(y-2))
	gl.End()
}

func Tri(x, y float64, r float64) {
	gl.PushMatrix()
	gl.Translatef(gl.Float(x), gl.Float(y), 0)
	gl.Rotatef(gl.Float(r), 0, 0, 1)
	gl.Begin(gl.TRIANGLES)
	gl.Color4f(1, 0, 0, 1)
	gl.Vertex2f(8, 0)
	gl.Vertex2f(-8, -5)
	gl.Vertex2f(-8, 5)
	gl.End()
	gl.PopMatrix()
}

func RenderScene(s *State, dt time.Duration) {
	for _, e := range s.Ents {
		x := float64(e.X)/10 + 1e-9*float64(dt)*float64(e.V)*math.Cos(float64(e.D)*math.Pi/180)
		y := float64(e.Y)/10 + 1e-9*float64(dt)*float64(e.V)*math.Sin(float64(e.D)*math.Pi/180)
		Quad(x, y)
	}
	for _, e := range s.Players {
		x := float64(e.X)/10 + 1e-9*float64(dt)*float64(e.V)*math.Cos(float64(e.D)*math.Pi/180)
		y := float64(e.Y)/10 + 1e-9*float64(dt)*float64(e.V)*math.Sin(float64(e.D)*math.Pi/180)
		r := math.Atan2(float64(e.FY)-y, float64(e.FX)-x) * 180 / math.Pi
		Tri(x, y, r)
	}
}

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

func StartGraphics(id uint32, update chan *State, move chan MoveMsg, quit chan bool) error {
	rc := sdl.Init(sdl.INIT_VIDEO)
	if rc == -1 {
		return errors.New(sdl.GetError())
	}
	disp := sdl.SetVideoMode(800, 600, 32, sdl.OPENGL)
	if disp == nil {
		return errors.New(sdl.GetError())
	}
	err := gl.Init()
	if err != nil {
		return err
	}
	s := new(State)
	ls := time.Now()

	SetupGL()

	tick := time.Tick(time.Second / 50)
	squit := make(chan bool)
	go ProcessInput(move, quit, squit)
	for {
		select {
		case s = <-update:
			ls = time.Now()
		case <-tick:
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.LoadIdentity()
			gl.Translatef(0, 0, -1)
			RenderScene(s, time.Now().Sub(ls))
			sdl.GL_SwapBuffers()
		case <-squit:
			return nil
		}
	}
	return nil
}
