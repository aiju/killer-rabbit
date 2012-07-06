package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/aiju/gl"
	"time"
	"os"
	"log"
	"image"
	_ "image/png"
)

var quadData = []float64{0, 0, 0, 1, 1, 0, 1, 1}
var quadBuf *gl.Buffer
var prog *gl.Program
var winscale gl.Mat4
var texs = make(map[string] gl.Texture)

func AddTex(tex, file string) {
	f, err := os.Open("data/" + file)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	t := gl.NewTexture2D(img, 0)
	t.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	t.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	texs[tex] = t
}

func Quad(cx, cy, d, r float64, tex string) {
	t, ok := texs[tex]
	if !ok {
		panic("no such texture " + tex)
	}
	prog.Use()
	prog.EnableAttrib("pos", quadBuf, 0, 2, 0, false)
	t.Enable(0, gl.TEXTURE_2D)
	prog.SetUniform("matrix", gl.Mul4(winscale, gl.Translate(cx-d/2, cy-d/2, 0), gl.Scale(d, d, 1), gl.Translate(0.5, 0.5, 0), gl.RotZ(r), gl.Translate(-0.5, -0.5, 0)))
	prog.SetUniform("tex", 0)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	prog.DisableAttrib("pos")
	prog.Unuse()
}

func (p *Player) Draw() {
	Quad(p.X, p.Y, PlayerSize, p.R(), "glenda")
}

func (p *Bullet) Draw() {
	Quad(p.X, p.Y, 4, p.D, "black")
}

func RenderScene(s *State) {
	for _, e := range s.Ents {
		e.Draw()
	}
	for _, e := range s.Players {
		e.Draw()
	}
	for _, e := range s.Bullets {
		e.Draw()
	}
}

func StartGraphics(id uint32, update chan *State, move chan MoveMsg, quit chan bool) error {
	var err error

	rc := sdl.Init(sdl.INIT_VIDEO)
	if rc == -1 {
		return errors.New(sdl.GetError())
	}
	disp := sdl.SetVideoMode(800, 600, 32, sdl.OPENGL)
	if disp == nil {
		return errors.New(sdl.GetError())
	}
	winscale = gl.Mul4(gl.Translate(-1, 1, 0), gl.Scale(1./400, -1./300, 1))
	gl.Init()
	gl.ClearColor(1, 1, 1, 1)
	s := new(State)
	ls := time.Now()
	quadBuf = gl.NewBuffer(gl.ARRAY_BUFFER, quadData, gl.STATIC_DRAW)
	prog, err = gl.MakeProgram([]string{vertexShader}, []string{fragmentShader})
	if err != nil {
		log.Fatal(err)
	}
	AddTex("glenda", "glenda.png")
	AddTex("tux", "tux.png")
	i := image.NewRGBA(image.Rect(0, 0, 800, 600))
	texs["black"] = gl.NewTexture2D(i, 0)

	tick := time.Tick(time.Second / 50)
	squit := make(chan bool)
	go ProcessInput(move, quit, squit)
	for {
		select {
		case s = <-update:
			ls = time.Now()
		case <-tick:
			s.Advance(time.Now().Sub(ls))
			ls = time.Now()
			gl.Clear(gl.COLOR_BUFFER_BIT)
			RenderScene(s)
			sdl.GL_SwapBuffers()
		case <-squit:
			return nil
		}
	}
	return nil
}

var vertexShader = `
#version 110

attribute vec2 pos;
uniform mat4 matrix;
varying vec2 texcoord;

void main() {
	gl_Position = matrix * vec4(pos, 0, 1);
	texcoord = pos;
}`

var fragmentShader = `
#version 110

varying vec2 texcoord;
uniform sampler2D tex;

void main() {
	gl_FragColor = texture2D(tex, texcoord);
}`
