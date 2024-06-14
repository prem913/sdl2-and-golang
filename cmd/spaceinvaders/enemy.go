package main

import gls "github.com/prem913/gl_go/pkg/gls"

type enemy struct {
	tex           *gls.Texture
	pos           *gls.Pos
	speed         float32
	health        float32
	inviCooldown  float32
	invincibility float32
	firecooldown  float32
	firecountdown float32
}

func NewEnemy(enemyName string, x, y, cooldown, speed, health, inviCooldown, firecooldown float32) *enemy {
	return &enemy{
		tex:           texs[enemyName],
		pos:           &gls.Pos{X: x, Y: y},
		speed:         speed,
		health:        health,
		inviCooldown:  inviCooldown,
		firecooldown:  firecooldown,
		firecountdown: firecooldown,
	}
}

func (e *enemy) update(delta float32,s *gls.SDL) {

	if e.invincibility > 0 {
		e.invincibility -= delta
	}

	if e.firecountdown > 0 {
		e.firecountdown -= delta
	}

	// e.pos.Y += delta * e.speed
  if e.pos.X - float32(e.tex.W)/2 < 0{
    e.pos.X = float32(e.tex.W) /2
    e.speed = -e.speed
  }
  if e.pos.X + float32(e.tex.W) / 2 > float32(s.WinWidth){
    e.pos.X = float32(s.WinWidth - (e.tex.W) / 2)
    e.speed = -e.speed
  }
  e.pos.X += (delta*e.speed)

	// enemy movement logic
}

type enemies struct {
	list    []*enemy
	bullets *Bullets
}

func NewEnemies() *enemies {
	return &enemies{
		list:    make([]*enemy, 0, 100),
		bullets: NewBullets(100),
	}
}

func (e *enemies) update(delta float32,s *gls.SDL) {
	newenemies := make([]*enemy, 0, 100)

	for _, e := range e.list {
		if e.health >= 0 {
			newenemies = append(newenemies, e)
			e.update(delta,s)
		}
	}
	e.list = newenemies
}
func (e *enemies) addEnemy(enemy *enemy) {
	e.list = append(e.list, enemy)
}
func (e *enemies) draw(s *gls.SDL) {
	for _, enemy := range e.list {
		enemy.tex.DrawAlpha(*enemy.pos, s)
	}
}
func (e *enemy) ResetFireCooldown() {
	e.firecountdown = e.firecooldown
}
