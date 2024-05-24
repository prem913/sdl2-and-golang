package sdl2utilities

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func LoadImageFromPath(path string) image.Image {
	// Read image from file that already exists
	existingImageFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	imageData, err := png.Decode(existingImageFile)
	if err != nil {
		panic(err)
	}

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	return imageData
	// loadedImage, err := png.Decode(existingImageFile)
	// if err != nil {
	// 	panic(err)
	//   }
	// fmt.Println(loadedImage)
}

func LoadImageTexture(imgPath string) *Texture {
	img := LoadImageFromPath(imgPath)

	var texture Texture
	imgBounds := img.Bounds()
	texture.W = imgBounds.Max.X - imgBounds.Min.X
	texture.H = imgBounds.Max.Y - imgBounds.Min.Y
	texture.Pixels = make([]uint32, texture.W*texture.W)
	fmt.Printf("W : %d H : %d min : %d max : %d", texture.W, texture.H, imgBounds.Min.X, imgBounds.Max.X)
	for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
		for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
			index := y*texture.W + x
			r, g, b, a := img.At(x, y).RGBA()
			color := Color{byte(r / 256), byte(g / 256), byte(b / 256), byte(a / 256)}
			texture.Pixels[index] = color.toUint32()
		}
	}
	return &texture
}
