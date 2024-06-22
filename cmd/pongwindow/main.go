package main

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/prem913/gl_go/pkg/gls"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	s := gls.Init_Sdl(gls.SDLOptions{
		WinW:    800,
		WinH:    800,
		WinName: "Window dayo",
		WinX:    1,
		WinY:    1,
	})
	// sdl.CreateWindow("2",100,100,40,40,sdl.WINDOW_SHOWN)
	context, _ := s.Window.GLCreateContext()
	defer sdl.GLDeleteContext(context)
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Viewport(0, 0, int32(800), int32(800))

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
			}
		}
		drawgl()
		s.Window.GLSwap()
	}
}
func drawgl() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	gl.Begin(gl.TRIANGLES)
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Vertex2f(0.5, 0.0)
	gl.Color3f(0.0, 1.0, 0.0)
	gl.Vertex2f(-0.5, -0.5)
	gl.Color3f(0.0, 0.0, 1.0)
	gl.Vertex2f(-0.5, 0.5)
	gl.Color3f(0.2, 0.0, 1.0)
	gl.Vertex2f(-0.5, 0.5)
	gl.End()
}
