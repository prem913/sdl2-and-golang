package main

import (
	gls "github.com/prem913/gl_go/pkg/gls"
)

// Direction is an array of length 4 type bool representing (right,left,top,botton) directions

type Bulleter interface {
	GetDamage() float32

	// Returns x,y,w,h
	GetBounds() (float32, float32, float32, float32)
	GetSpeed() (float32, float32)
}

type Bullet struct {
	tex       *gls.Texture
	pos       *gls.Pos
	damage    float32
	xv        float32
	yv        float32
	destroyed bool
}

func NewBullet(tex *gls.Texture, x, y, damage, speedX, speedY float32) *Bullet {
	return &Bullet{
		tex:       tex,
		pos:       gls.NewPos(x, y),
		damage:    damage,
		xv:        speedX,
		yv:        speedY,
		destroyed: false,
	}
}

func (b *Bullet) GetDamage() float32 {
	return b.damage
}

func (b *Bullet) GetBounds() (float32, float32, float32, float32) {
	return b.pos.X, b.pos.Y, float32(b.tex.W), float32(b.tex.H)
}

func (b *Bullet) GetDirection() (float32, float32) {
	return b.xv, b.yv
}

// TODO: paralalize
type Bullets struct {
	list []*Bullet
}

func NewBullets(buffSize uint32) *Bullets {
	b := Bullets{
		list: make([]*Bullet, 0, buffSize),
	}
	return &b
}

func (b *Bullets) Add(bullet *Bullet) {
	b.list = append(b.list, bullet)
}

func (b *Bullets) draw(s *gls.SDL) {
	for _, p := range b.list {
		p.tex.DrawAlpha(*p.pos, s)
	}
}

func (b *Bullets) update(delta float32, s *gls.SDL) {
  // fmt.Printf("Bullets Now : %v \r",len(b.list))
	for _, bullet := range b.list {
		// move bullet
		bullet.pos.X += (delta * bullet.xv)
		bullet.pos.Y += (delta * bullet.yv)
	}

	// keeping it simple for now
	// remove bullets if bullet goes out of window or destroyed
	newlist := make([]*Bullet, 0, 100)
	for _, bullet := range b.list {
		if bullet.pos.X+float32(bullet.tex.W)/2 < 0 || bullet.pos.X-float32(bullet.tex.W)/2 > float32(s.WinWidth) || bullet.pos.Y+float32(bullet.tex.H)/2 < 0 || bullet.pos.Y+float32(bullet.tex.H)/2 > float32(s.WinHeight) || bullet.destroyed {
			continue
		}
		newlist = append(newlist, bullet)
	}
	b.list = newlist
}

// func (b *bullets) update(delta float32, keyState []uint8, player *player) {
//
// 	newpos := make([]gls.Pos, 0, 100)
// 	b.countdown -= delta
// 	if keyState[sdl.SCANCODE_SPACE] != 0 {
// 		if b.countdown < 0 {
// 			newpos = append(newpos, *player.pos)
// 			b.countdown = b.cooldown
// 		}
// 	}
// 	for _, p := range b.pos {
// 		p.Y -= b.speed * delta
// 		if p.Y > 0 {
// 			newpos = append(newpos, p)
// 		}
// 	}
// 	b.pos = newpos
// 	// fmt.Printf("Total Bullets : %d\r", len(newpos))
// }
