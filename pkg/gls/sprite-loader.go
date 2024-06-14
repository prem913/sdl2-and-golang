package gls

// takes Image Path and a map with sprite position data in the order [x,y,w,h]
// return map with same names along with textures
func LoadSprite(path string, spritesdata map[string][4]int) map[string]*Texture {
	imgTexture, err := LoadImageTexture(path)
	if err != nil {
		panic(err)
	}

	res := make(map[string]*Texture)
	for name, data := range spritesdata {
		X, Y, W, H := data[0], data[1], data[2], data[3]
		pix := make([]byte, W*H*4)
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				texIndex := (Y+y)*imgTexture.Pitch + (x+X)*4
				pixIndex := y*W*4 + x*4

				pix[pixIndex] = imgTexture.Pixels[texIndex]
				pix[pixIndex+1] = imgTexture.Pixels[texIndex+1]
				pix[pixIndex+2] = imgTexture.Pixels[texIndex+2]
				pix[pixIndex+3] = imgTexture.Pixels[texIndex+3]
			}
		}
		res[name] = NewTexture(W, H, W*4, pix)
	}

	return res
}
