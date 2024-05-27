package sdl2utilities

import (
	"image"
	"image/png"
	"os"
)

func LoadImageFromPath(path string) (image.Image, error) {
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
	return imageData, nil
	// loadedImage, err := png.Decode(existingImageFile)
	// if err != nil {
	// 	panic(err)
	//   }
	// fmt.Println(loadedImage)
}

func LoadImageTexture(imgPath string) (*Texture, error) {
	img, err := LoadImageFromPath(imgPath)
	if err != nil {
		panic(err)
	}

	imgBounds := img.Bounds()
	w := imgBounds.Max.X - imgBounds.Min.X
	h := imgBounds.Max.Y - imgBounds.Min.Y
	pixels := make([]byte, w*h*4)
	index := 0
	for y := imgBounds.Min.Y; y < h; y++ {
		for x := imgBounds.Min.X; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rr, gg, bb, aa := byte(r/256), byte(g/256), byte(b/256), byte(a/256)
			pixels[index] = rr
			index++
			pixels[index] = gg
			index++
			pixels[index] = bb
			index++
			pixels[index] = aa
			index++
		}
	}
	return &Texture{pixels, w, h, w * 4}, nil
}
