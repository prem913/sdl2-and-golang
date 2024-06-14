package main

import (
	"fmt"
	"math/rand"
	"time"

	gls "github.com/prem913/gl_go/pkg/gls"
	"github.com/veandco/go-sdl2/sdl"
	// "github.com/veandco/go-sdl2/sdl"
)

const (
	WindWidth = 510
	WinHeight = 920
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
	gravity   float32   = 600
	score     float32   = 0
)

type Bird struct {
	pos       gls.Pos
	texs      []gls.Texture
	curtex    int
	animtime  float32
	animspeed float32
	speed     float32
	jumpcd    float32
	jumptime  float32
	boost     float32
	maxspeed  float32
	flying    bool
}

func NewBird(pos gls.Pos, animspeed, jumpcd, boost, maxspeed float32) *Bird {
	var err error
	var tex *gls.Texture
	texs := make([]gls.Texture, 0, 3)
	tex, err = gls.LoadImageTexture("./cmd/flappy_birb/assets/bluebird-upflap.png")
	texs = append(texs, *tex)
	tex, err = gls.LoadImageTexture("./cmd/flappy_birb/assets/bluebird-midflap.png")
	texs = append(texs, *tex)
	tex, err = gls.LoadImageTexture("./cmd/flappy_birb/assets/bluebird-downflap.png")
	texs = append(texs, *tex)

	if err != nil {
		panic("couldn't load textures")
	}

	return &Bird{
		pos:       pos,
		texs:      texs,
		curtex:    0,
		animtime:  0,
		animspeed: animspeed,
		speed:     0,
		jumpcd:    jumpcd,
		jumptime:  0,
		boost:     boost,
		maxspeed:  maxspeed,
	}
}

func (b *Bird) update(keyState []uint8, delta, gravity float32) {
	// animate
	b.animtime -= delta
	if b.animtime <= 0 {
		b.animtime = b.animspeed
		b.curtex = (b.curtex + 1) % len(b.texs)
	}

	b.jumptime -= delta
	if b.jumptime <= 0 && keyState[sdl.SCANCODE_SPACE] != 0 {
		b.speed = -b.boost
		if b.speed > b.maxspeed {
			b.speed = b.maxspeed
		}
		b.jumptime = b.jumpcd
	}

	b.speed += gravity * delta
	b.pos.Y += b.speed * delta

}

func (b *Bird) draw(s *gls.SDL) {
	if b.speed < -100 {
		b.texs[b.curtex].DrawAlphaRotate(b.pos, float64(Lerp(0, 45, -b.speed/b.maxspeed)), s)
	} else if b.speed > 100 {
		b.texs[b.curtex].DrawAlphaRotate(b.pos, float64(-Lerp(0, 45, b.speed/b.maxspeed)), s)
	} else {
		b.texs[b.curtex].DrawAlpha(b.pos, s)
	}
}

type Base struct {
	tex   gls.Texture
	pos   gls.Pos
	speed float32
}

func NewBase(speed float32) *Base {
	tex, err := gls.LoadImageTexture("./cmd/flappy_birb/assets/base.png")
	if err != nil {
		panic("failed loading assets")
	}

	return &Base{
		tex:   *tex,
		pos:   *gls.NewPos(WindWidth/2, float32(WinHeight-tex.H/2)),
		speed: speed,
	}
}

func (b *Base) update(delta float32) {
	b.pos.X -= delta * b.speed
	if b.pos.X < 0 {
		b.pos.X = WindWidth / 2
	}
}

func (b *Base) draw(s *gls.SDL) {
	b.tex.DrawAlpha(b.pos, s)
	p := b.pos
	q := b.pos
	p.X += float32(b.tex.W)
	q.X -= float32(b.tex.W)
	b.tex.DrawAlpha(p, s)
	b.tex.DrawAlpha(q, s)
}

type Pipes struct {
	upTex   gls.Texture
	downTex gls.Texture
	pipes   []gls.Pos
	speed   float32
	count   int
}

func NewPipes(speed float32, count int) *Pipes {
	uptex, _ := gls.LoadImageTexture("./cmd/flappy_birb/assets/pipeNorth.png")
	downtex, _ := gls.LoadImageTexture("./cmd/flappy_birb/assets/pipeSouth.png")
	pipes := make([]gls.Pos, count)
	initpos := *gls.NewPos(0, 0)
	for i := 0; i < count; i++ {
		initpos.X += WindWidth / 2
		initpos.Y = WinHeight / 2
		pipes[i] = initpos
	}

	return &Pipes{
		upTex:   *uptex,
		downTex: *downtex,
		pipes:   pipes,
		speed:   speed,
		count:   count,
	}
}

func (p *Pipes) update(delta float32) {
	for i := 0; i < p.count; i++ {
		p.pipes[i].X -= delta * p.speed
		if p.pipes[i].X < 0 {
			var maxX float32 = 0
			for j := 0; j < p.count; j++ {
				if i == j {
					continue
				}
				if maxX < p.pipes[j].X {
					maxX = p.pipes[j].X
				}
			}
			p.pipes[i].X += WindWidth/2 + maxX
			p.pipes[i].Y = Lerp(WinHeight/2-100, WinHeight/2+100, rand.Float32())
		}
	}
}

func (p *Pipes) draw(s *gls.SDL) {
	for _, pipe := range p.pipes {
		p.downTex.DrawAlpha(*gls.NewPos(pipe.X, pipe.Y-float32(p.downTex.H)), s)
		p.upTex.DrawAlpha(*gls.NewPos(pipe.X, pipe.Y-float32(p.upTex.H)*2), s)
		p.upTex.DrawAlpha(*gls.NewPos(pipe.X, pipe.Y+float32(p.upTex.H)), s)
		p.downTex.DrawAlpha(*gls.NewPos(pipe.X, pipe.Y+float32(p.downTex.H)*2), s)
	}
}

func Collisions(p *Pipes, b *Bird) {
	for _, pipe := range p.pipes {
		// top pipes
		pw, ph, bw, bh := float32(p.upTex.W), float32(p.upTex.H), float32(b.texs[0].W), float32(b.texs[0].H)
		if RectCollision(pipe.X, pipe.Y-ph, pw, ph, b.pos.X, b.pos.Y, bw, bh) || RectCollision(pipe.X, pipe.Y+ph, pw, ph, b.pos.X, b.pos.Y, bw, bh) {
			gameState = GAMEOVER
			fmt.Println("Score : ", score)
		}

	}
}

func resetGame(p *Pipes, b *Bird) {
	score = 0
	initpos := *gls.NewPos(0, 0)
	for i := 0; i < p.count; i++ {
		initpos.X += WindWidth / 2
		initpos.Y = WinHeight / 2
		p.pipes[i] = initpos
	}
	b.pos = getCenter()
}

type Background struct {
	texs  []gls.Texture
	pos   gls.Pos
	speed float32
	day   int
}

func NewBackground(speed float32) *Background {
	daytex, _ := gls.LoadImageTexture("./cmd/flappy_birb/assets/background-day.png")
	nighttex, _ := gls.LoadImageTexture("./cmd/flappy_birb/assets/background-night.png")
	daytex.Scale(int(float32(daytex.W)*1.8), int(float32(daytex.H)*1.8))
	nighttex.Scale(int(float32(nighttex.W)*1.8), int(float32(nighttex.H)*1.8))
	texs := []gls.Texture{*daytex, *nighttex}
	return &Background{
		texs:  texs,
		pos:   *gls.NewPos(WindWidth/2, WinHeight-float32(daytex.H/2)),
		speed: speed,
	}
}
func (b *Background) update(delta float32) {
	b.pos.X -= delta * b.speed
	if b.pos.X < -float32(b.texs[0].W/2) {
		b.pos.X = WindWidth / 2
	}
	if score > 10 {
		b.day = 1
	}
}
func (b *Background) draw(s *gls.SDL) {
	tex := b.texs[b.day]
	tex.Draw(b.pos, s)
	pp := b.pos
	pp.X = b.pos.X + WindWidth
	tex.Draw(pp, s)
	pp.X = b.pos.X - WindWidth
	tex.Draw(pp, s)
}

func main() {
	var s gls.SDL
	s.Init_Sdl(WindWidth, WinHeight, "Flappy bird")

	keyState := sdl.GetKeyboardState()
	fps := NewFramerate()

	bird := NewBird(getCenter(), 0.1, 0.1, 300, 300)
	pipes := NewPipes(300, 400)
	background := NewBackground(30)

	texmap := loadAssets()

	bgTex := texmap["background"]
	bgTex.Scale(int(float32(bgTex.W)*1.8), int(float32(bgTex.H)*1.8))
	bgpos := *gls.NewPos(WindWidth/2, WinHeight-float32(bgTex.H/2))
	fmt.Println(bgTex.W, bgTex.H)

	base := NewBase(300)

	clickcd := float32(0)

	s.DrawScreen(func(delta float32) {
		if clickcd <= 0 {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				switch gameState {
				case INIT:
					gameState = START
					break
				case GAMEOVER:
					gameState = START
					resetGame(pipes, bird)
					break
				case PAUSE:
					gameState = START
					break
				}
			}
			clickcd -= 0.2
		} else {
			clickcd -= delta
		}

		switch gameState {
		case START:
			score += delta
			go func() {
				Collisions(pipes, bird)
			}()

			bird.update(keyState, delta, gravity)
			pipes.update(delta)
			base.update(delta)
			background.update(delta)
			break
		}
	}, func() {
		fps.run()
		switch gameState {
		case INIT:
			bgTex.Draw(bgpos, &s)
			bird.draw(&s)
			break
		case START:
			background.draw(&s)
				bird.draw(&s)
				pipes.draw(&s)
				base.draw(&s)
			break
		case GAMEOVER:
			texmap["gameover"].DrawAlpha(getCenter(), &s)
			break
		}
	})
}

func loadAssets() map[string]*gls.Texture {
	var err error
	texmap := make(map[string]*gls.Texture)
	texmap["background"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/background-day.png")
	// texmap["background"].Upscale()
	texmap["bird"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/bird.png")
	texmap["gameover"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/game-over.png")
	texmap["pipenorth"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/pipeNorth.png")
	texmap["pipesouth"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/pipeSouth.png")
	texmap["pipes"], err = gls.LoadImageTexture("./cmd/flappy_birb/assets/pipes.png")

	if err != nil {
		panic("Unable to load assets")
	}

	return texmap
}

func Lerp(l, r, p float32) float32 {
	return (r-l)*p + l
}
func getCenter() gls.Pos {
	return *gls.NewPos(WindWidth/2, WinHeight/2)
}

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
