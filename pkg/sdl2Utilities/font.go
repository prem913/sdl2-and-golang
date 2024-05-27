package sdl2utilities

// Function to create the 15x15 scaled bitmap from the 5x3 bitmap
func createTextureBitmap(numberBitmap []string) [15][15]int {
	var scaledBitmap [15][15]int
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			if numberBitmap[i][j] == '1' {
				for m := 0; m < 3; m++ {
					for n := 0; n < 5; n++ {
						scaledBitmap[i*3+m][j*5+n] = 1
					}
				}
			}
		}
	}
	return scaledBitmap
}

func createTextureByteArray(scaledBitmap [15][15]int) *Texture {
	byteArray := make([]byte, 0, 15*15*4)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if scaledBitmap[i][j] == 1 {
				byteArray = append(byteArray, 255, 255, 255, 255)
			} else {
				byteArray = append(byteArray, 0, 0, 0, 0)
			}
		}
	}
	tex := NewTexture(15, 15, 15*4, byteArray)
	return tex
}

// Function to convert a number to its 15-bit bitmap and then to a texture byte array
func NumberToTextureByteArray() map[rune]*Texture {
	bitmaps := map[rune]string{
		'0': "111101101101111",
		'1': "010110010010111",
		'2': "111001111100111",
		'3': "111001111001111",
		'4': "101101111001001",
		'5': "111100111001111",
		'6': "111100111101111",
		'7': "111001001001001",
		'8': "111101111101111",
		'9': "111101111001111",
		'A': "111101111111101",
		'B': "110101110101110",
		'C': "111100100100111",
		'D': "110101101101110",
		'E': "111100111100111",
		'F': "111100111100100",
		'G': "111100101101111",
		'H': "101101111101101",
		'I': "111010010010111",
		'J': "001001001101111",
		'K': "101101110101101",
		'L': "100100100100111",
		'M': "101111111101101",
		'N': "101111111111101",
		'O': "111101101101111",
		'P': "111101111100100",
		'Q': "111101101111011",
		'R': "111101110101101",
		'S': "111100111001111",
		'T': "111010010010010",
		'U': "101101101101111",
		'V': "101101101101010",
		'W': "101101111111101",
		'X': "101101010101101",
		'Y': "101101010010010",
		'Z': "111001010100111",
	}
	numTextures := make(map[rune]*Texture, len(bitmaps))
	for x := range bitmaps {
		numberBitmapStr := bitmaps[x]
		numberBitmap := make([]string, 5)
		for i := 0; i < 5; i++ {
			numberBitmap[i] = numberBitmapStr[i*3 : (i+1)*3]
		}
		scaledBitmap := createTextureBitmap(numberBitmap)
		byteArray := createTextureByteArray(scaledBitmap)
		numTextures[x] = byteArray
	}

	return numTextures
}
