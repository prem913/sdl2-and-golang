package main

import (
	gls "github.com/prem913/gl_go/pkg/gls"
	"github.com/veandco/go-sdl2/sdl"
)

type player struct {
	pos           *gls.Pos
	tex           *gls.Texture
	speed         float32
	invincibility float32
	inviCooldown  float32
	health        float32
	bullets       Bullets
	firecooldown  float32
	firecountdown float32
}

func NewPlayer(x, y, speed, health, inviCooldown, firecooldown float32) *player {
	return &player{
		pos:           &gls.Pos{X: x, Y: y},
		tex:           texs["player"],
		speed:         speed,
		inviCooldown:  inviCooldown,
		health:        health,
		bullets:       *NewBullets(100),
		firecooldown:  firecooldown,
		firecountdown: firecooldown,
	}
}
func (p *player) update(delta float32, keyState []uint8) {
	if p.invincibility > 0 {
		p.invincibility -= delta
	}

	if p.firecountdown > 0 {
		p.firecountdown -= delta
	}

	if keyState[sdl.SCANCODE_D] != 0 {
		if p.pos.X < WindWidth-float32(p.tex.W/2) {
			p.pos.X += p.speed * delta
		}
	}
	if keyState[sdl.SCANCODE_A] != 0 {
		if p.pos.X > float32(p.tex.W/2) {
			p.pos.X -= p.speed * delta
		}
	}
	if keyState[sdl.SCANCODE_W] != 0 {
		if p.pos.Y > float32(p.tex.H/2) {
			p.pos.Y -= p.speed * delta
		}
	}
	if keyState[sdl.SCANCODE_S] != 0 {
		if p.pos.Y < WinHeight-float32(p.tex.H/2) {
			p.pos.Y += p.speed * delta
		}
	}
}

func (p *player) ResetFireCooldown() {
	p.firecountdown = p.firecooldown
}

func (p *player) ResetCountdown() {
	p.invincibility = p.inviCooldown
}
