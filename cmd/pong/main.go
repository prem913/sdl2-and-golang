package main

import (
	"fmt"
	"math"
	"time"

	sdl2utilities "github.com/prem913/gl_go/pkg/sdl2Utilities"
	"github.com/veandco/go-sdl2/sdl"
)

// Types
type GameStateType int

// Global Constants
const (
	START = iota
	PAUSE
	GAMEOVER
	INIT
)

const (
	WindWidth int = 1024
	WinHeight int = 720
)

// Game State
var GameState GameStateType = INIT
var Score uint32 = 0
var paddleSpeed float32 = 800
var ballSpeedX, ballSpeedY float32 = 800, 900
var highScore int = 32
var curScore float32 = 0

// Game Objects
var numTextures = sdl2utilities.NumberToTextureByteArray()

func DrawText(str string, pos sdl2utilities.Pos, s *sdl2utilities.SDL, digitSpace float32) {
	curPos := pos
	for _, i := range str {
		if i == '#' {
			curPos.X = pos.X
			curPos.Y += 15 + digitSpace
			continue
		}
		curPos.X += 15 + digitSpace
		tex, exists := numTextures[i]
		if !exists {
			continue
		}
		tex.DrawAlpha(curPos, s)

	}
}

type Paddle struct {
	sdl2utilities.Pos
	tex   *sdl2utilities.Texture
	speed float32
}

func NewPaddle(x, y, speed float32, w, h int, color *sdl2utilities.Color) *Paddle {
	pix := make([]byte, w*h*4)
	pitch := w * 4
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := (y*pitch + x*4)
			pix[index], pix[index+1], pix[index+2], pix[index+3] = color.RGBA()
		}
	}
	tex := sdl2utilities.NewTexture(w, h, pitch, pix)
	return &Paddle{
		Pos:   sdl2utilities.Pos{X: x, Y: y},
		speed: speed,
		tex:   tex,
	}
}

func (p *Paddle) update(keyState []uint8, delta float32) {
	if keyState[sdl.SCANCODE_UP] != 0 || keyState[sdl.SCANCODE_W] != 0 && p.Y-float32(p.tex.H)/2 > 0 {
		p.Y -= p.speed * float32(delta)
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 || keyState[sdl.SCANCODE_S] != 0 && p.Y+float32(p.tex.H)/2 < float32(WinHeight) {
		p.Y += p.speed * float32(delta)
	}
}
func (p *Paddle) aiUpdate(b *Ball) {
	p.Y = b.Y
}

// func (p *Paddle) draw(pixels []uint32) {
// 	// startX and startY becomes mid points of the rectangle
// 	startX := p.X - p.W/2
// 	startY := p.Y - p.H/2

// 	for y := 0; y < int(p.h); y++ {
// 		for x := 0; x < int(p.w); x++ {
// 			setPixel(int(startX)+x, int(startY)+y, p.color, pixels)
// 		}
// 	}
// }

type Ball struct {
	sdl2utilities.Pos
	radius float32
	xv     float32
	yv     float32
	tex    *sdl2utilities.Texture
}

func NewBall(X, Y, xv, yv, radius float32) *Ball {
	h := 2 * radius
	w := h
	pix := make([]byte, int(4*h*w))
	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			if x*x+y*y < radius*radius {
				index := int(((y+radius)*w + (x + radius))) * 4
				pix[index] = 255
				pix[index+1] = 255
				pix[index+2] = 255
				pix[index+3] = 255
			}

		}
	}
	tex := sdl2utilities.NewTexture(int(w), int(h), int(w*4), pix)
	return &Ball{
		Pos:    sdl2utilities.Pos{X: X, Y: Y},
		xv:     xv,
		yv:     yv,
		tex:    tex,
		radius: radius,
	}
}

// func (b *ball) draw(pixels []uint32) {

// 	for y := -b.radius; y < b.radius; y++ {
// 		for x := -b.radius; x < b.radius; x++ {
// 			if int(x*x+y*y) < int(b.radius*b.radius) {
// 				setPixel(int(b.x+x), int(b.y+y), b.color, pixels)
// 			}
// 		}
// 	}
// }

func (b *Ball) update(leftPaddle, rightPaddle *Paddle, elapsedTime float32) {
	b.X += b.xv * elapsedTime
	b.Y += b.yv * elapsedTime

	// ball hit top and bottom
	if b.Y-b.radius < 0 {
		b.yv = -b.yv
		b.Y = b.radius

	}
	if b.Y+b.radius > float32(WinHeight) {
		b.yv = -b.yv
		b.Y = float32(WinHeight) - b.radius
	}

	if b.X-b.radius < leftPaddle.X+float32(leftPaddle.tex.W/2) {
		if int(b.Y) > int(leftPaddle.Y)-leftPaddle.tex.H/2 && int(b.Y) < int(leftPaddle.Y)+leftPaddle.tex.H/2 {
			b.xv = -b.xv
			b.X = leftPaddle.X + float32(leftPaddle.tex.W)/2 + b.radius
		}
	}

	if b.X+b.radius > rightPaddle.X-float32(rightPaddle.tex.W/2) {
		if b.Y > rightPaddle.Y-float32(rightPaddle.tex.H/2) && b.Y < rightPaddle.Y+float32(rightPaddle.tex.H/2) {
			b.xv = -b.xv
			b.X = rightPaddle.X - float32(rightPaddle.tex.W)/2 - b.radius
		}
	}

	// close the game for now
	if b.X < 0 || b.X > float32(WindWidth) {
		b.X, b.Y = float32(WindWidth/2), float32(WinHeight/2)
		GameState = GAMEOVER
	}
}

func Lerp(l, r, p float32) float32 {
	return (r-l)*p + l
}

func main() {
	player1 := NewPaddle(Lerp(0, float32(WindWidth), 0.05), Lerp(0, float32(WinHeight), 0.5), paddleSpeed, 20, 180, sdl2utilities.NewColor(255, 255, 255, 255))
	aiplayer := NewPaddle(Lerp(0, float32(WindWidth), 0.95), Lerp(0, float32(WinHeight), 0.5), paddleSpeed, 20, 180, sdl2utilities.NewColor(255, 255, 255, 255))
	ball := NewBall(100, 100, ballSpeedX, ballSpeedY, 10)

	welcomeTextBox := sdl2utilities.NewTextBox(20, 15, 10)
	welcomeTextBox.UpdateText("       WELCOME      #         TO         #        PONG        #  SPACE TO CONTINUE ")
	welcomeTextBox.UpdateTexture()
	pauseTextBox := sdl2utilities.NewTextBox(20, 15, 5)
	pauseTextBox.UpdateText("PAUSE")
	pauseTextBox.UpdateTexture()
	var s sdl2utilities.SDL
	s.Init_Sdl(WindWidth, WinHeight)
	fps := 0
	frames := 0
	start := time.Now()

	// go func() {
	// 	for {
	// 		fmt.Printf("FPS : %d \r", fps)
	// 		time.Sleep(1000)
	// 	}
	// }()
	countdown := 4

	s.DrawScreen(func(delta float32) {
		keyState := sdl.GetKeyboardState()
		if keyState[sdl.SCANCODE_SPACE] != 0 {
			switch GameState {
			case INIT:
				GameState = START
			case START:
				GameState = PAUSE
			case PAUSE:
				GameState = START
			case GAMEOVER:
				GameState = INIT
			}
			sdl.Delay(150)

		}
		switch GameState {
		case START:
			player1.update(keyState, delta)
			if countdown == 0 {
				aiplayer.aiUpdate(ball)
				ball.update(player1, aiplayer, delta)
				curScore += 0.1
			}
		case PAUSE:
			countdown = 4
		case GAMEOVER:
			countdown = 4
		}

	}, func() {
		frames++
		s.Clearscreen()
		switch GameState {
		case INIT:
			// DrawText(fmt.Sprintf("  WELCOME #    TO    #   PONG   #HIGH SCORE %d", highScore), getCenter(WindWidth-200, WinHeight), &s, 3)
			welcomeTextBox.Tex.Draw(getCenter(WindWidth, WinHeight), &s)

		case START:
			ball.tex.DrawAlpha(ball.Pos, &s)
			player1.tex.DrawAlpha(player1.Pos, &s)
			aiplayer.tex.DrawAlpha(aiplayer.Pos, &s)
			DrawText(fmt.Sprintf("FPS %d#SCORE %d", fps, int(curScore)), sdl2utilities.Pos{X: 15, Y: 15}, &s, 2)
			if time.Since(start).Milliseconds() > 1000 {
				if countdown > 0 && GameState == START {
					countdown--
				}
				fps = frames
				frames = 0
				start = time.Now()
			}
			if countdown > 0 {
				DrawText(fmt.Sprint(countdown), getCenter(WindWidth, WinHeight), &s, 0)
			}
		case GAMEOVER:
			highScore = int(math.Max(float64(curScore), float64(highScore)))
			DrawText(fmt.Sprintf("GAME OVER#SCORE %d", int(curScore)), getCenter(WindWidth-200, WinHeight), &s, 3)
			curScore = 0
		case PAUSE:
			DrawText(fmt.Sprintf("PAUSE#SCORE %d", int(curScore)), getCenter(WindWidth-200, WinHeight), &s, 3)
			// pauseTextBox.UpdateText(fmt.Sprintf("       PAUSE        #HIGH SCORE : %d", int(curScore)))
			// pauseTextBox.UpdateTexture()
			// pauseTextBox.Tex.Draw(getCenter(WindWidth, WinHeight), &s)
		}

	})

}

func GetDigits(number int) []byte {
	digits := make([]byte, 0, 2)
	for number != 0 {
		digits = append(digits, byte(number%10))
		number /= 10
	}
	n := len(digits) - 1
	i := 0
	for i < n {
		t := digits[i]
		digits[i] = digits[n]
		digits[n] = t
		i++
		n--
	}
	return digits
}

func getCenter(winwidth, winheight int) sdl2utilities.Pos {
	return sdl2utilities.Pos{
		X: float32(winwidth / 2),
		Y: float32(winheight / 2),
	}
}
