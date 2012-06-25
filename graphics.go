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
	gl.Vertex2f(gl.Float(x+5), gl.Float(y+5))
	gl.Vertex2f(gl.Float(x-5), gl.Float(y+5))
	gl.Vertex2f(gl.Float(x-5), gl.Float(y-5))
	gl.Vertex2f(gl.Float(x+5), gl.Float(y-5))
	gl.End()
}

func RenderScene(s *State, dt time.Duration) {
	for _, e := range s.P {
		x := float64(e.X)/10 + 1e-9*float64(dt)*float64(e.V)*math.Cos(float64(e.D)*math.Pi/180)
		y := float64(e.Y)/10 + 1e-9*float64(dt)*float64(e.V)*math.Sin(float64(e.D)*math.Pi/180)
		Quad(x, y)
	}
}

func StartGraphics(update chan *State, move chan Ent) error {
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
	for {
		select {
		case ev := <-sdl.Events:
			if _, ok := ev.(sdl.QuitEvent); ok {
				return nil
			}
			if k, ok := ev.(sdl.KeyboardEvent); ok {
				if k.Keysym.Sym == 'a' {
					if k.Type == sdl.KEYDOWN {
						move <- Ent{0, 0, 100, 180}
					}
					if k.Type == sdl.KEYUP {
						move <- Ent{0, 0, 0, 180}
					}
				}
				if k.Keysym.Sym == 's' {
					if k.Type == sdl.KEYDOWN {
						move <- Ent{0, 0, 100, 90}
					}
					if k.Type == sdl.KEYUP {
						move <- Ent{0, 0, 0, 90}
					}
				}
				if k.Keysym.Sym == 'd' {
					if k.Type == sdl.KEYDOWN {
						move <- Ent{0, 0, 100, 0}
					}
					if k.Type == sdl.KEYUP {
						move <- Ent{0, 0, 0, 0}
					}
				}
				if k.Keysym.Sym == 'w' {
					if k.Type == sdl.KEYDOWN {
						move <- Ent{0, 0, 100, 270}
					}
					if k.Type == sdl.KEYUP {
						move <- Ent{0, 0, 0, 270}
					}
				}
			}
		case s = <-update:
			ls = time.Now()
		case <-tick:
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.LoadIdentity()
			gl.Translatef(0, 0, -1)
			RenderScene(s, time.Now().Sub(ls))
			sdl.GL_SwapBuffers()
		}
	}
	return nil
}
