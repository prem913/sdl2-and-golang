package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/prem913/gl_go/pkg/gls"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindWidth int = 600
	WinHeight int = 600
)

var (
	generation uint32 = 0
	LiveCells  uint32 = 0
	DeadCells  uint32 = 0
)

func setPixel(x, y, w int, pix []byte, value byte) {
	index := (y*w + x) * 4
	pix[index] = value
	pix[index+1] = value
	pix[index+2] = value
	pix[index+3] = value
}

func randomBorn(s *gls.Texture) {
	for i := 0; i < 1000; i++ {
		setPixel(rand.IntN(s.W), rand.IntN(s.H), s.W, s.Pixels, 255)
	}

}

func universe(w, h int) *gls.Texture {
	pix := make([]byte, w*h*4)

	for y := 0; y < w; y++ {
		for x := 0; x < h; x++ {
			index := (y*w + x) * 4
			live := byte(0)
			if rand.IntN(90) > 50 {
				live = 255
			}

			pix[index] = live
			pix[index+1] = live
			pix[index+2] = live
			pix[index+3] = live
		}
	}
	return gls.NewTexture(w, h, w*4, pix)
}

var dirX [8]int = [8]int{0, 0, 1, -1, -1, -1, 1, 1}
var dirY [8]int = [8]int{-1, 1, 0, 0, -1, 1, -1, 1}

func tick(tex *gls.Texture) {
	generation++
	LiveCells = 0
	DeadCells = 0
	pix := make([]byte, tex.W*tex.H*4)
	for y := 0; y < tex.W; y++ {
		for x := 0; x < tex.H; x++ {
			liveCells := 0
			deadCells := 0

			// for k := 0; k < 4; k++ {
			// 	for i := 0; i < 3; i++ {
			// 		dx := x + dirX[k] + i
			// 		dy := y + dirY[k] + i
			// 		if dx >= 0 && dy >= 0 && dx < WindWidth && dy < WinHeight {
			// 			idx := (dy*tex.W + dx) * 4

			// 			if tex.Pixels[idx] == 0 {
			// 				deadCells++
			// 			} else {
			// 				liveCells++
			// 			}
			// 		}
			// 	}
			// }
			for k := 0; k < 8; k++ {
				dx := x + dirX[k]
				dy := y + dirY[k]
				if dx >= 0 && dy >= 0 && dx < tex.W && dy < tex.H {
					idx := (dy*tex.W + dx) * 4

					if tex.Pixels[idx] == 0 {
						deadCells++
					} else {
						liveCells++
					}
				}
			}
			index := (y*tex.W + x) * 4
			live := tex.Pixels[index]

			if live == 0 {
				LiveCells++
				if liveCells == 3 {
					live = 255
				}
			} else {
				DeadCells++
				if liveCells < 2 {
					live = 0
				} else if liveCells > 3 {
					live = 0
				}
			}

			pix[index] = live
			pix[index+1] = live
			pix[index+2] = live
			pix[index+3] = live
		}
	}
	tex.Pixels = pix

}

func main() {
	s := gls.Init_Sdl(gls.SDLOptions{
		WinW:    int32(WindWidth),
		WinH:    int32(WinHeight),
		WinName: "Game of life",
	})
	uni := universe(400, 400)
	keyState := sdl.GetKeyboardState()

	s.DrawScreen(func(delta float32) {
		tick(uni)
		if keyState[sdl.SCANCODE_SPACE] != 0 {
			randomBorn(uni)
		}

	}, func() {
		uni.Draw(getCenter(), s)
		fmt.Printf("Generation : %d LiveCells : %d DeadCells : %d\r", generation, LiveCells, DeadCells)
	})

}

func getCenter() gls.Pos {
	return gls.Pos{X: float32(WindWidth) / 2, Y: float32(WinHeight) / 2}
}
