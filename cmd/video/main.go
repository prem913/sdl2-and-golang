package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"

	"github.com/prem913/gl_go/pkg/gls"
)

func main() {
	// Path to the video file
	args := os.Args

	videoPath := args[1]

	// Execute ffmpeg command to extract frames
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", "-f", "image2pipe", "-vcodec", "png", "-")
	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing ffmpeg...")
	err := cmd.Run()
	fmt.Println("executing ffmpeg done")
	if err != nil {
		fmt.Println("Error running ffmpeg:", err)
		return
	}

	var s gls.SDL
	s.Init_Sdl(800, 800,"test")

	// Read the frames from the output buffer
	// for {


	s.DrawScreen(func(delta float32) {
		img, _ := png.Decode(&out)
		rgbaArray := ImageToRGBAArray(img)

		fmt.Println("Extracted frame RGBA array:", len(rgbaArray.Pixels))

	}, func() {
	})

	// Do something with rgbaArray, e.g., store it in a larger array for all frames
	//	}
}

// ImageToRGBAArray converts an image.Image to an RGBA array
func ImageToRGBAArray(img image.Image) *gls.Texture {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgbaArray := make([]byte, width*height*4)

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgbaArray[index] = byte(r >> 8)
			rgbaArray[index+1] = byte(g >> 8)
			rgbaArray[index+2] = byte(b >> 8)
			rgbaArray[index+3] = byte(a >> 8)
			index += 4
		}
	}

	return gls.NewTexture(width, height, width*4, rgbaArray)
}
