package gls

import (
	"fmt"
	"math"
	"sync"
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

	mu     sync.RWMutex
	Screen []byte
}

type Pos struct {
	X, Y float32
}

func NewPos(x, y float32) *Pos {
	return &Pos{
		X: x,
		Y: y,
	}
}

func (s *SDL) Init_Sdl(winWidth, winHeight int,windowname string) {
	s.WinHeight = winHeight
	s.WinWidth = winWidth
	s.Screen = make([]byte, s.WinHeight*s.WinWidth*4)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	window, err := sdl.CreateWindow(windowname, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

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
	stop := false

	frames := make(chan []byte,50)

	var frameStart time.Time
	var elapsedTime float32

	for i := range s.Screen {
		if s.Screen[i] == 0 {
			continue
		}
	}

	go func() {
		defer func() {
			stop = true
			fmt.Println("loop exit")
		}()
		for {
			frameStart = time.Now()
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					return

				}
			}
			updatefunc(elapsedTime)
			elapsedTime = float32(time.Since(frameStart).Seconds())
			// if elapsedTime < 1 {
			// 	sdl.Delay(uint32(1 - elapsedTime / 1000.0))
			// 	elapsedTime = float32(time.Since(frameStart).Seconds())
			// }
		}
	}()
	go func() {
		for {
			drawfunc()
			copied := make([]byte, len(s.Screen))
			copy(copied, s.Screen)
			frames <- copied
		}
	}()

	for {
		if stop {
			return
		}
		frame := <-frames
		s.tex.Update(nil, unsafe.Pointer(&frame[0]), s.WinWidth*4)
		s.renderer.Copy(s.tex, nil, nil)
		s.renderer.Present()

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
func CopyTexture(tex *Texture) *Texture {
	newpix := make([]byte, len(tex.Pixels))
	copy(newpix, tex.Pixels)
	return NewTexture(tex.W, tex.H, tex.Pitch, newpix)
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

func (tex *Texture) DrawAlphaRotate(p Pos, deg float64, s *SDL) {
	rad := float64(deg * 0.0174533)
	for y := 0; y < tex.H; y++ {
		for x := 0; x < tex.W; x++ {
			px, py := p.X, p.Y
			ox, oy := float64(x), float64(y)
			rx := int(math.Round(ox*math.Cos(rad) + oy*math.Sin(rad)))
			ry := int(math.Round(oy*math.Cos(rad) - ox*math.Sin(rad)))

			screenY := ry + int(py) - tex.H/2
			screenX := rx + int(px) - tex.W/2

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

func (tex *Texture) DrawAlphaScaled(p Pos, nw, nh int, s *SDL) {
	ow, oh := tex.W, tex.H
	for ny := 0; ny < nh; ny++ {
		for nx := 0; nx < nw; nx++ {
			ox := fit(nx, ow, nw)
			oy := fit(ny, oh, nh)
			screenY := ny + int(p.Y) - nh/2
			screenX := nx + int(p.X) - nw/2

			if screenX >= 0 && screenX < s.WinWidth && screenY >= 0 && screenY < s.WinHeight {
				texIndex := oy*tex.Pitch + ox*4
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
func (tex *Texture) Upscale() {
	tex.Pixels = upscale(tex.Pixels, tex.W, tex.H)
	tex.W *= 2
	tex.H *= 2
	tex.Pitch *= 2
}
func upscale(pix []byte, w, h int) []byte {
	newpix := make([]byte, w*h*16)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idxold := (y * w * 4) + x*4
			idxnew := idxold * 2
			for i := 0; i < 4; i++ {
				idxnew := idxnew + (i * 4)
				newpix[idxnew] = pix[idxold]
				newpix[idxnew+1] = pix[idxold+1]
				newpix[idxnew+2] = pix[idxold+2]
				newpix[idxnew+3] = pix[idxold+3]
			}
		}
	}
	return newpix
}

func (tex *Texture) Scale(w, h int) {
	ow, oh := tex.W, tex.H

	l := w * h * 4
	// ol := ow * oh * 4

	rp := make([]byte, l)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := y*w*4 + x*4
			ox := fit(x, ow, w)
			oy := fit(y, oh, h)
			oindex := oy*ow*4 + ox*4
			if ox >= ow || oy >= oh {
				continue
			}
			rp[index] = tex.Pixels[oindex]
			rp[index+1] = tex.Pixels[oindex+1]
			rp[index+2] = tex.Pixels[oindex+2]
			rp[index+3] = tex.Pixels[oindex+3]
		}
	}

	tex.W = w
	tex.H = h
	tex.Pixels = rp
	tex.Pitch = w * 4
}

func fit(ov, nr, or int) int {
	p := float32(ov) / float32(or)
	return int(p * float32(nr))
}
