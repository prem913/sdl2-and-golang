package sdl2utilities

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type SDL struct {
	WinWidth  int
	WinHeight int
	window    *sdl.Window
	renderer  *sdl.Renderer
	tex       *sdl.Texture
	Screen    []byte
}

type Pos struct {
	X, Y float32
}

func (s *SDL) Init_Sdl(winWidth, winHeight int) {
	s.WinHeight = winHeight
	s.WinWidth = winWidth
	s.Screen = make([]byte, s.WinHeight*s.WinWidth*4)
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

type updatefunctype func(delta float32)
type drawfunctype func()

func (s *SDL) DrawScreen(updatefunc updatefunctype, drawfunc drawfunctype) {
	var frameStart time.Time
	var elapsedTime float32

	for i := range s.Screen {
		if s.Screen[i] == 0 {
			continue
		}
	}

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return

			}
		}

		updatefunc(elapsedTime)
		drawfunc()

		s.tex.Update(nil, unsafe.Pointer(&s.Screen[0]), s.WinWidth*4)
		s.renderer.Copy(s.tex, nil, nil)
		s.renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())

		if elapsedTime < 0.005 {
			// sdl.Delay(uint32(elapsedTime / 1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

	}
}

func SetPixel(x, y, w, h int, c Color, pix []byte) {
	index := (y*w + x) * 4
	pix[index] = c.r
	pix[index+1] = c.g
	pix[index+2] = c.b
	pix[index+3] = c.a
}

func NewTexture(W, H, Pitch int, Pixles []byte) *Texture {
	return &Texture{
		Pixles,
		W,
		H,
		Pitch,
	}
}

//	func (s *SDL) SetPixelUint(x, y int, color uint32) {
//		index := (y*s.WinWidth + x)
//		if index < len(s.Screen) && index >= 0 {
//			s.Screen[index] = color
//		}
//	}
func (s *SDL) Clearscreen() {
	for i := range s.Screen {
		s.Screen[i] = 0
	}
}

type Texture struct {
	Pixels      []byte
	W, H, Pitch int
}
type Color struct {
	r, g, b, a byte
}

func NewColor(r, g, b, a byte) *Color {
	return &Color{r, g, b, a}
}

func (c *Color) RGBA() (byte, byte, byte, byte) {
	return c.r, c.g, c.b, c.a
}

// func (c *Color) toUint32() uint32 {
// 	ui := uint32(0)
// 	ui = uint32((uint32(c.a) << 24) + (uint32(c.b) << 16) + (uint32(c.g) << 8) + uint32(c.r))
// 	return ui
// }

// func toARGB(c uint32) (uint32, uint32, uint32, uint32) {
// 	return (c >> 24) & 0xFF, (c >> 16) & 0xFF, (c >> 8) & 0xFF, c & 0xFF
// }

func (tex *Texture) Draw(p Pos, s *SDL) {
	for y := 0; y < tex.H; y++ {
		for x := 0; x < tex.W; x++ {
			// make origin center
			screenY := y + int(p.Y) - tex.H/2
			screenX := x + int(p.X) - tex.W/2

			if screenX >= 0 && screenX < s.WinWidth && screenY >= 0 && screenY < s.WinHeight {
				texIndex := y*tex.Pitch + x*4
				screenIndex := screenY*s.WinWidth*4 + screenX*4
				s.Screen[screenIndex] = tex.Pixels[texIndex]
				s.Screen[screenIndex+1] = tex.Pixels[texIndex+1]
				s.Screen[screenIndex+2] = tex.Pixels[texIndex+2]
				s.Screen[screenIndex+3] = tex.Pixels[texIndex+3]
			}
		}
	}
}

func (tex *Texture) DrawAlpha(p Pos, s *SDL) {
	for y := 0; y < tex.H; y++ {
		for x := 0; x < tex.W; x++ {
			screenY := y + int(p.Y) - tex.H/2
			screenX := x + int(p.X) - tex.W/2

			if screenX >= 0 && screenX < s.WinWidth && screenY >= 0 && screenY < s.WinHeight {
				texIndex := y*tex.Pitch + x*4
				screenIndex := screenY*s.WinWidth*4 + screenX*4

				srcR, srcG, srcB, srcA := tex.Pixels[texIndex], tex.Pixels[texIndex+1], tex.Pixels[texIndex+2], tex.Pixels[texIndex+3]

				dstR, dstG, dstB, _ := s.Screen[screenIndex], s.Screen[screenIndex+1], s.Screen[screenIndex+2], s.Screen[screenIndex+3]

				rstR := (int(srcR)*255 + int(dstR)*(255-int(srcA))) / 255
				rstG := (int(srcG)*255 + int(dstG)*(255-int(srcA))) / 255
				rstB := (int(srcB)*255 + int(dstB)*(255-int(srcA))) / 255

				s.Screen[screenIndex] = byte(rstR)
				s.Screen[screenIndex+1] = byte(rstG)
				s.Screen[screenIndex+2] = byte(rstB)
			}

		}
	}
}
