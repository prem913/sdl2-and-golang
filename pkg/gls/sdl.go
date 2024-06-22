package gls

import (
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type SDLOptions struct {
	WinH, WinW, WinX, WinY int32
	WinName                string
	onMouseMove            func(event *sdl.MouseMotionEvent)
}

type SDL struct {
	SDLOptions
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Tex      *sdl.Texture

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

func Init_Sdl(options SDLOptions) *SDL {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	var wposX int32 = sdl.WINDOWPOS_UNDEFINED
	var wposY int32 = sdl.WINDOWPOS_UNDEFINED
	if options.WinX != 0 && options.WinY != 0 {
		wposX = options.WinX
		wposY = options.WinY
	}
	Window, err := sdl.CreateWindow(options.WinName, wposX, wposY, options.WinW, options.WinH, sdl.WINDOW_OPENGL)

	if err != nil {
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, options.WinW, options.WinH)
	if err != nil {
		panic(err)
	}
	if options.onMouseMove == nil {
		options.onMouseMove = func(t *sdl.MouseMotionEvent) {
			// fmt.Printf("[%d ms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
		}
	}
	Screen := make([]byte, options.WinH*options.WinW*4)
	return &SDL{
		SDLOptions: options,
		Window:     Window,
		Tex:        tex,
		Renderer:   renderer,
		Screen:     Screen,
	}
}

type updatefunctype func(delta float32)
type drawfunctype func()

func (s *SDL) DrawScreen(updatefunc updatefunctype, drawfunc drawfunctype) {
	stop := false

	frames := make(chan []byte)

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
				switch t := event.(type) {
				case *sdl.QuitEvent:
					return
				case *sdl.MouseMotionEvent:
					s.onMouseMove(t)
        case *sdl.MouseButtonEvent:
          fmt.Println(t)
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
		s.Tex.Update(nil, unsafe.Pointer(&frame[0]), int(s.WinW)*4)
		s.Renderer.Copy(s.Tex, nil, nil)
		s.Renderer.Present()

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
func (s *SDL) ClearScreen() {
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
			ScreenY := y + int(p.Y) - tex.H/2
			ScreenX := x + int(p.X) - tex.W/2

			if ScreenX >= 0 && ScreenX < int(s.WinW) && ScreenY >= 0 && ScreenY < int(s.WinH) {
				texIndex := y*tex.Pitch + x*4
				ScreenIndex := ScreenY*int(s.WinW)*4 + ScreenX*4
				s.Screen[ScreenIndex] = tex.Pixels[texIndex]
				s.Screen[ScreenIndex+1] = tex.Pixels[texIndex+1]
				s.Screen[ScreenIndex+2] = tex.Pixels[texIndex+2]
				s.Screen[ScreenIndex+3] = tex.Pixels[texIndex+3]
			}
		}
	}
}

func (tex *Texture) DrawAlpha(p Pos, s *SDL) {
	for y := 0; y < tex.H; y++ {
		for x := 0; x < tex.W; x++ {
			ScreenY := y + int(p.Y) - tex.H/2
			ScreenX := x + int(p.X) - tex.W/2

			if ScreenX >= 0 && ScreenX < int(s.WinW) && ScreenY >= 0 && ScreenY < int(s.WinH) {
				texIndex := y*tex.Pitch + x*4
				ScreenIndex := ScreenY*int(s.WinW)*4 + ScreenX*4

				srcR, srcG, srcB, srcA := tex.Pixels[texIndex], tex.Pixels[texIndex+1], tex.Pixels[texIndex+2], tex.Pixels[texIndex+3]

				dstR, dstG, dstB, _ := s.Screen[ScreenIndex], s.Screen[ScreenIndex+1], s.Screen[ScreenIndex+2], s.Screen[ScreenIndex+3]

				rstR := (int(srcR)*255 + int(dstR)*(255-int(srcA))) / 255
				rstG := (int(srcG)*255 + int(dstG)*(255-int(srcA))) / 255
				rstB := (int(srcB)*255 + int(dstB)*(255-int(srcA))) / 255

				s.Screen[ScreenIndex] = byte(rstR)
				s.Screen[ScreenIndex+1] = byte(rstG)
				s.Screen[ScreenIndex+2] = byte(rstB)
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

			ScreenY := ry + int(py) - tex.H/2
			ScreenX := rx + int(px) - tex.W/2

			if ScreenX >= 0 && ScreenX < int(s.WinW) && ScreenY >= 0 && ScreenY < int(s.WinH) {
				texIndex := y*tex.Pitch + x*4
				ScreenIndex := ScreenY*int(s.WinW)*4 + ScreenX*4

				srcR, srcG, srcB, srcA := tex.Pixels[texIndex], tex.Pixels[texIndex+1], tex.Pixels[texIndex+2], tex.Pixels[texIndex+3]

				dstR, dstG, dstB, _ := s.Screen[ScreenIndex], s.Screen[ScreenIndex+1], s.Screen[ScreenIndex+2], s.Screen[ScreenIndex+3]

				rstR := (int(srcR)*255 + int(dstR)*(255-int(srcA))) / 255
				rstG := (int(srcG)*255 + int(dstG)*(255-int(srcA))) / 255
				rstB := (int(srcB)*255 + int(dstB)*(255-int(srcA))) / 255

				s.Screen[ScreenIndex] = byte(rstR)
				s.Screen[ScreenIndex+1] = byte(rstG)
				s.Screen[ScreenIndex+2] = byte(rstB)
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
			ScreenY := ny + int(p.Y) - nh/2
			ScreenX := nx + int(p.X) - nw/2

			if ScreenX >= 0 && ScreenX < int(s.WinH) && ScreenY >= 0 && ScreenY < int(s.WinH) {
				texIndex := oy*tex.Pitch + ox*4
				ScreenIndex := ScreenY*int(s.WinW)*4 + ScreenX*4

				srcR, srcG, srcB, srcA := tex.Pixels[texIndex], tex.Pixels[texIndex+1], tex.Pixels[texIndex+2], tex.Pixels[texIndex+3]

				dstR, dstG, dstB, _ := s.Screen[ScreenIndex], s.Screen[ScreenIndex+1], s.Screen[ScreenIndex+2], s.Screen[ScreenIndex+3]

				rstR := (int(srcR)*255 + int(dstR)*(255-int(srcA))) / 255
				rstG := (int(srcG)*255 + int(dstG)*(255-int(srcA))) / 255
				rstB := (int(srcB)*255 + int(dstB)*(255-int(srcA))) / 255

				s.Screen[ScreenIndex] = byte(rstR)
				s.Screen[ScreenIndex+1] = byte(rstG)
				s.Screen[ScreenIndex+2] = byte(rstB)
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
