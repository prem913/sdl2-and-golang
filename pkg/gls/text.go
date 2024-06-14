package gls

var CharMap = NumberToTextureByteArray()

type TextBox struct {
	charLen  uint32
	charSize uint32
	charGap  uint32
	text     string
	W, H     uint32
	Tex      *Texture
}

func NewTextBox(charLen, charSize uint32, charGap uint32) *TextBox {
	return &TextBox{
		charLen:  charLen,
		charSize: charSize,
		charGap:  charGap,
		Tex:      NewTexture(0, 0, 0, make([]byte, 0)),
		text:     "",
	}
}

func (t *TextBox) UpdateText(text string) {
	t.text = text
	t.W = t.GetWidth()
	t.H = t.GetHeight()
	// fmt.Println("Textbox : ", t.W, t.H)
}

func (t *TextBox) GetWidth() uint32 {
	return t.charLen * (t.charSize + t.charGap)
}

func (t *TextBox) GetHeight() uint32 {
	numNewLines := uint32(1)
	numChars := 0
	for _, i := range t.text {
		if i == '#' {
			numNewLines++
			numChars = 0
			continue
		}
		numChars++
		if numChars > int(t.charLen) {
			numNewLines++
			numChars = 0
		}
	}
	// fmt.Println("New Lines : ", numNewLines)
	return (t.charSize + t.charGap) * numNewLines
}

func (t *TextBox) UpdateTexture() {
	pix := make([]byte, int(t.W*t.H*4))
	t.Tex = &Texture{
		Pixels: pix,
		W:      int(t.W),
		H:      int(t.H),
		Pitch:  int(t.W) * 4,
	}
	// border
	// for y := 0; y < int(t.H); y++ {
	// 	for x := 0; x < int(t.W); x++ {
	// 		index := (y * int(t.W) * 4) + x*4
	// 		if x == 0 || x == int(t.W-1) || y == 0 || int(t.H-1) == y {
	// 			pix[index] = 255
	// 			pix[index+1] = 255
	// 			pix[index+2] = 255
	// 			pix[index+3] = 255
	// 		}
	// 	}
	// }

	numchar := 0
	curX := 0
	curY := 0

	for _, c := range t.text {
		if c == '#' {
			curY += int(t.charSize) + int(t.charGap)
			numchar = 0
			curX = 0
			continue
		}
		if c == ' ' {
			curX += int(t.charSize) + int(t.charGap)
			continue
		}
		if numchar > int(t.charLen) {
			curY += int(t.charGap) + int(t.charSize)
			numchar = 0
			curX = 0
		}
		tex := CharMap[c]
		for y := 0; y < tex.H; y++ {
			for x := 0; x < tex.W; x++ {
				Y := y + curY
				X := x + curX

				Idx := (Y * int(t.W) * 4) + X*4
				idx := (y * tex.Pitch) + x*4

				pix[Idx] = tex.Pixels[idx]
				pix[Idx+1] = tex.Pixels[idx+1]
				pix[Idx+2] = tex.Pixels[idx+2]
				pix[Idx+3] = tex.Pixels[idx+3]
			}
		}
		curX += int(t.charSize) + int(t.charGap)
		numchar++
	}

}
