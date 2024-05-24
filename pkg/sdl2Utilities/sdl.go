package sdl2utilities

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type SDL struct {
	WinWidth  int
	WinHeight int
	window    *sdl.Window
	renderer  *sdl.Renderer
	tex       *sdl.Texture
	Screen    []uint32
}

type Pos struct {
	X, Y int
}

func (s *SDL) Init_Sdl(winWidth, winHeight int) {
	s.WinHeight = winHeight
	s.WinWidth = winWidth
	s.Screen = make([]uint32, s.WinHeight*s.WinWidth)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	window, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	s.window = window
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	s.renderer = renderer
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		panic(err)
	}
	s.tex = tex
}

func (s *SDL) DrawScreen() {
	var frameStart time.Time
	var elapsedTime float32

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		s.tex.UpdateRGBA(nil, s.Screen, s.WinWidth)
		s.renderer.Copy(s.tex, nil, nil)
		s.renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())

		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

	}
}

func (s *SDL) SetPixel(x, y int, c Color) {
	index := (y*s.WinWidth + x)
	if index < len(s.Screen) && index >= 0 {
		s.Screen[index] = c.toUint32()
	}
}
func (s *SDL) SetPixelUint(x, y int, color uint32) {
	index := (y*s.WinWidth + x)
	if index < len(s.Screen) && index >= 0 {
		s.Screen[index] = color
	}
}
func (s *SDL) Clearscreen() {
	for i := range s.Screen {
		s.Screen[i] = 0
	}
}

type Texture struct {
	Pos
	Pixels []uint32
	W, H   int
}
type Color struct {
	r, g, b, a byte
}

func (c *Color) toUint32() uint32 {
	ui := uint32(0)
	ui = uint32((uint32(c.a) << 24) + (uint32(c.b) << 16) + (uint32(c.g) << 8) + uint32(c.r))
	return ui
}

func (t *Texture) Draw(s *SDL) {
	for y := 0; y < t.H; y++ {
		for x := 0; x < t.W; x++ {
			index := (y*t.W + x)
			s.SetPixelUint(t.X+x, t.Y+y, t.Pixels[index])
		}
	}
}
