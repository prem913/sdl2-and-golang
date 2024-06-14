package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	gls "github.com/prem913/gl_go/pkg/gls"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindWidth = 1600
	WinHeight = 800
)

type GameState int

const (
	INIT = iota
	PAUSE
	START
	GAMEOVER
)

var (
	gameState GameState = INIT
)

var texs = gls.LoadSprite("./cmd/spaceinvaders/assets/shipsall.png", map[string][4]int{
	"bosslvl1": {4, 4, 58, 58},
	"bosslvl2": {130, 130, 56, 56},
	"player":   {5, 68, 26, 22},
	"bullet":   {133, 68, 19, 14},
})

func RectCollision(x1, y1, w1, h1, x2, y2, w2, h2 float32) bool {
	l1 := x1 - w1/2
	r1 := x1 + w1/2
	t1 := y1 - h1/2
	b1 := y1 + h1/2

	l2 := x2 - w2/2
	r2 := x2 + w2/2
	t2 := y2 - h2/2
	b2 := y2 + h2/2

	return !(r1 < l2 || l1 > r2 || t1 > b2 || b1 < t2)
}

func HandleEnemyCollision(pb *Bullets, enemies *enemies) {
	for _, bullet := range pb.list {
		for _, enemy := range enemies.list {

			collided := RectCollision(bullet.pos.X, bullet.pos.Y, float32(bullet.tex.W), float32(bullet.tex.H), enemy.pos.X, enemy.pos.Y, float32(enemy.tex.W), float32(enemy.tex.H))

			if enemy.invincibility <= 0 && !bullet.destroyed && collided {
				enemy.health -= bullet.damage
				enemy.invincibility = enemy.inviCooldown
				bullet.destroyed = true
			}
		}
	}
}

func HandlePlayerCollision(eb *Bullets, player *player) {
	for _, bullet := range eb.list {
		collided := RectCollision(bullet.pos.X, bullet.pos.Y, float32(bullet.tex.W), float32(bullet.tex.H), player.pos.X, player.pos.Y, float32(player.tex.W), float32(player.tex.H))
		if player.invincibility <= 0 && collided {
			player.health -= bullet.damage
			bullet.destroyed = true
		}
	}
}

func HandleEnemyBullets(enemies *enemies, tex *gls.Texture) {
	for _, enemy := range enemies.list {
		if enemy.firecountdown <= 0 {
			enemies.bullets.Add(NewBullet(tex, enemy.pos.X, enemy.pos.Y, 100, enemy.speed/2, 900))
			enemy.ResetFireCooldown()
		}
	}
}

func HandlePlayerBullets(player *player, keyState []uint8, tex *gls.Texture) {
	if keyState[sdl.SCANCODE_SPACE] != 0 && player.firecountdown <= 0 {
		player.bullets.Add(NewBullet(tex, player.pos.X, player.pos.Y, 100, 0, -900))
		player.ResetFireCooldown()
	}
}

func SpawnEnemeis(enemies *enemies) {
	for i := 0; i < 10; i++ {
		enemies.addEnemy(NewEnemy("bosslvl2", float32(rand.Int32N(WindWidth)), 100, 0.5, 100+float32(rand.Int32N(3))*100, 100, 0.2, 0.3+float32(rand.IntN(7))*0.1))
	}
}

//	func returnRandomEnemy(){
//	 list :=
//	}
type Framerate struct {
	frames     float32
	start      time.Time
	framecount float32
}

func NewFramerate() *Framerate {
	return &Framerate{
		frames:     0,
		start:      time.Now(),
		framecount: 0,
	}
}

func (f *Framerate) run() {
	f.framecount++
	if time.Since(f.start).Milliseconds() >= 1000 {
		f.frames = f.framecount
		f.framecount = 0
		f.start = time.Now()
	}
	fmt.Printf("fps : %.1f \r", f.frames)
}

func main() {
	var s gls.SDL
	s.Init_Sdl(WindWidth, WinHeight,"SPI")

	keyState := sdl.GetKeyboardState()

	texs["bosslvl1"].FlipX()
	texs["bosslvl2"].FlipX()

	var texEnemyBullet = gls.CopyTexture(texs["bullet"])
	texEnemyBullet.FlipX()

	player := NewPlayer(100, Lerp(100, float32(WinHeight), 0.9), 800, 800, 0.2, 0.4)
	enemies := NewEnemies()

	SpawnEnemeis(enemies)

	fps := NewFramerate()

	s.DrawScreen(func(delta float32) {
		HandleEnemyCollision(&player.bullets, enemies)
		HandlePlayerCollision(enemies.bullets, player)
		HandleEnemyBullets(enemies, texEnemyBullet)
		HandlePlayerBullets(player, keyState, texs["bullet"])
		player.update(delta, keyState)
		enemies.update(delta, &s)
		enemies.bullets.update(delta, &s)
		player.bullets.update(delta, &s)

	}, func() {
		var wg sync.WaitGroup
		wg.Add(4)
			s.Clearscreen()
		go func() {
			player.tex.DrawAlpha(*player.pos, &s)
			wg.Done()
		}()
		go func() {
			enemies.draw(&s)
			wg.Done()
		}()
		go func() {
			enemies.bullets.draw(&s)
			wg.Done()
		}()
		go func() {
			player.bullets.draw(&s)
			wg.Done()
		}()

		wg.Wait()
		fps.run()
	})
}

func Lerp(l, r, p float32) float32 {
	return (r-l)*p + l
}
