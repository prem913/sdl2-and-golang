package main

import (
	"fmt"
	"math/rand/v2"

	sdl2utilities "github.com/prem913/gl_go/pkg/sdl2Utilities"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindWidth = 800
	WinHeight = 800
)

var texs = sdl2utilities.LoadSprite("./cmd/spaceinvaders/assets/shipsall.png", map[string][4]int{
	"bosslvl1": {4, 4, 58, 58},
	"bosslvl2": {130, 130, 56, 56},
	"player":   {5, 68, 26, 22},
	"bullet":   {133, 68, 19, 14},
})

type player struct {
	pos   *sdl2utilities.Pos
	tex   *sdl2utilities.Texture
	speed float32
}


func NewPlayer(x, y, speed float32) *player {
	return &player{
		pos:   &sdl2utilities.Pos{X: x, Y: y},
		tex:   texs["player"],
		speed: speed,
	}
}
func (p *player) update(delta float32, keyState []uint8) {
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

type bullets struct {
	tex       *sdl2utilities.Texture
	pos       []sdl2utilities.Pos
	speed     float32
	cooldown  float32
	countdown float32
	damage    float32
}

func NewBullets(damage, speed, cooldown float32) *bullets {
	return &bullets{
		tex: texs["bullet"],
		pos: make([]sdl2utilities.Pos, 0, 100), speed: speed,
		cooldown:  cooldown,
		countdown: cooldown,
		damage:    damage,
	}
}

func (b *bullets) update(delta float32, keyState []uint8, player *player) {
	newpos := make([]sdl2utilities.Pos, 0, 100)
	b.countdown -= delta
	if keyState[sdl.SCANCODE_SPACE] != 0 {
		if b.countdown < 0 {
			newpos = append(newpos, *player.pos)
			b.countdown = b.cooldown
		}
	}
	for _, p := range b.pos {
		p.Y -= b.speed * delta
		if p.Y > 0 {
			newpos = append(newpos, p)
		}
	}
	b.pos = newpos
	// fmt.Printf("Total Bullets : %d\r", len(newpos))
}

func (b *bullets) draw(s *sdl2utilities.SDL) {
	for _, p := range b.pos {
		b.tex.DrawAlpha(p, s)
	}
}

type enemy struct {
	tex           *sdl2utilities.Texture
	pos           *sdl2utilities.Pos
	cooldown      float32
	countdown     float32
	speed         float32
	health        float32
	inviCooldown  float32
	invincibility float32
}

type enemyBullets struct {
	bullets []sdl2utilities.Pos
	damage  float32
}

func NewEnemyBullets(damage float32) *enemyBullets {
	return &enemyBullets{
		bullets: make([]sdl2utilities.Pos, 0, 100),
		damage:  damage,
	}
}

func NewEnemy(enemyName string, x, y, cooldown, speed, health, inviCooldown float32) *enemy {
	return &enemy{
		tex:          texs[enemyName],
		pos:          &sdl2utilities.Pos{X: x, Y: y},
		cooldown:     cooldown,
		countdown:    cooldown,
		speed:        speed,
		health:       health,
		inviCooldown: inviCooldown,
	}
}

func (e *enemy) update(delta float32, bullets *bullets, player *player) {
	e.pos.Y += delta * e.speed
	if player.pos.X > e.pos.X {
		e.pos.X += delta * e.speed
	} else if player.pos.X < e.pos.X {
		e.pos.X -= e.speed * delta
	}
	if e.invincibility > 0 {
		e.invincibility -= delta
	} else {
		e.invincibility = 0
	}
	for _, p := range bullets.pos {
		// b <- e , b -> e , b ^ e , b bottom e
		if e.invincibility == 0 && !(p.X+float32(bullets.tex.W)/2 < e.pos.X || p.X > e.pos.X+float32(e.tex.W)/2 || p.Y+float32(bullets.tex.H)/2 < e.pos.Y || p.Y-float32(bullets.tex.H)/2 > e.pos.Y+float32(e.tex.H)/2) {
			e.health -= bullets.damage
			e.invincibility = e.inviCooldown
		}
	}
}

type enemies struct {
	enemies []enemy
}

func NewEnemies() *enemies {
	return &enemies{
		enemies: make([]enemy, 0, 100),
	}
}

func (e *enemies) update(delta float32, bullets *bullets, player *player) {
	newenemies := make([]enemy, 0, 100)

	for _, e := range e.enemies {
		e.update(delta, bullets, player)
		if e.health > 0 {
			newenemies = append(newenemies, e)
		}
	}
	e.enemies = newenemies
}
func (e *enemies) addEnemy(enemy *enemy) {
	e.enemies = append(e.enemies, *enemy)
}
func (e *enemies) draw(s *sdl2utilities.SDL) {
	for _, enemy := range e.enemies {
		enemy.tex.DrawAlpha(*enemy.pos, s)
	}
}

func main() {
	var s sdl2utilities.SDL
	s.Init_Sdl(WindWidth, WinHeight)
	keyState := sdl.GetKeyboardState()
	player := NewPlayer(100, Lerp(100, float32(WinHeight), 0.9), 800)
	bullets := NewBullets(40, 900, 0.2)
	texs["bosslvl1"].FlipX()
	enemies := NewEnemies()
	bossSpawnCD := float32(1)

	s.DrawScreen(func(delta float32) {
		bossSpawnCD -= delta
		fmt.Println(len(enemies.enemies))
		if bossSpawnCD < 0 && len(enemies.enemies) < 6 {
			enemies.addEnemy(NewEnemy("bosslvl1", float32(rand.Int32N(WindWidth)), 0, 0.2, 100, 100, 0.2))
			bossSpawnCD = 1
		}
		player.update(delta, keyState)
		bullets.update(delta, keyState, player)
		enemies.update(delta, bullets, player)

	}, func() {
		s.Clearscreen()
		player.tex.DrawAlpha(*player.pos, &s)
		bullets.draw(&s)
		enemies.draw(&s)
	})

}

func Lerp(l, r, p float32) float32 {
	return (r-l)*p + l
}
