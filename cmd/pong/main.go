package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 800

type Pixel uint32

type color struct {
	r, g, b, a uint32
}

func (c *color) toUint32() uint32 {
	ui := uint32(0)
	ui = uint32((uint32(c.a) << 24) + (uint32(c.b) << 16) + (uint32(c.g) << 8) + uint32(c.r))
	return ui
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (b *ball) draw(pixels []uint32) {

	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			if int(x*x+y*y) < int(b.radius*b.radius) {
				setPixel(int(b.x+x), int(b.y+y), b.color, pixels)
			}
		}
	}
}

func (b *ball) update(leftPaddle, rightPaddle *paddle, elapsedTime float32) {
	b.x += b.xv * elapsedTime
	b.y += b.yv * elapsedTime

	// ball hit top
	if b.y-b.radius < 0 || b.y+b.radius > float32(winHeight) {
		b.yv = -b.yv
	}

	if b.x < leftPaddle.x+leftPaddle.w/2 {
		if b.y > leftPaddle.y-leftPaddle.h/2 && b.y < leftPaddle.y+leftPaddle.h/2 {
			b.xv = -b.xv
		}
	}

	if b.x > rightPaddle.x-rightPaddle.w/2 {
		if b.y > rightPaddle.y-rightPaddle.h/2 && b.y < rightPaddle.y+rightPaddle.h/2 {
			b.xv = -b.xv
		}
	}

	// close the game for now
	if b.x < 0 || b.x > float32(winWidth) {
		b.x, b.y = getCenter()
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
}

func (p *paddle) draw(pixels []uint32) {
	// startX and startY becomes mid points of the rectangle
	startX := p.x - p.w/2
	startY := p.y - p.h/2

	for y := 0; y < int(p.h); y++ {
		for x := 0; x < int(p.w); x++ {
			setPixel(int(startX)+x, int(startY)+y, p.color, pixels)
		}
	}
}

func (p *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= p.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += p.speed * elapsedTime
	}
}

func (p *paddle) aiUpdate(b *ball) {
	p.y = b.y
}

func setPixel(x, y int, c color, pixels []uint32) {
	index := (y*winWidth + x)
	if index < len(pixels) && index >= 0 {
		pixels[index] = c.toUint32()
	}
}
func clearPixels(pixels []uint32) {
	for i := range pixels {
		pixels[i] = 0
	}
}
func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	defer window.Destroy()

	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()
	if err != nil {
		panic(err)
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	defer tex.Destroy()
	if err != nil {
		panic(err)
	}
	pixels := make([]uint32, winHeight*winWidth)
	playerh, playerw := float32(100), float32(20)
	playerInitialPosition := pos{playerw / 2, playerh / 2}

	player1 := paddle{playerInitialPosition, playerw, playerh, 800, color{255, 255, 255, 255}}
	player2 := paddle{pos{float32(winWidth - 10), 0}, 20, 100, 300, color{255, 255, 255, 255}}
	ball1 := ball{pos{200, 200}, 10, 400, 500, color{255, 255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32
	frames := 0
	var minStart time.Time
	minStart = time.Now()

	for {
		frames++
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clearPixels(pixels)

		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball1)
		ball1.update(&player1, &player2, elapsedTime)

		player1.draw(pixels)
		player2.draw(pixels)
		ball1.draw(pixels)

		tex.UpdateRGBA(nil, pixels, winWidth)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())

		if time.Since(minStart).Milliseconds() > 1000 {
			minStart = time.Now()
			fmt.Printf("FPS : %d \r", frames)
			frames = 0
		}
		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

	}
}

func getCenter() (float32, float32) {
	return float32(winWidth) / 2, float32(winHeight) / 2
}
