package main

import (
	sdl2utilities "github.com/prem913/gl_go/pkg/sdl2Utilities"
)

func main() {
	tex := sdl2utilities.LoadImageTexture("./img2.png")
	var s sdl2utilities.SDL
	s.Init_Sdl(1000, 1000)
	s.Clearscreen()
	tex.Draw(&s)
	s.DrawScreen()
}
