package game

// This file defines the platforms in the game, including their positioning, collisions, and interaction with the player.

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Platform struct {
	image  *ebiten.Image
	x, y   float64
	width  int
	height int
}

func NewPlatform(image *ebiten.Image, x, y float64) *Player {
	width, height := image.Size()
	return &Player{
		image:  image,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (p *Platform) Update() {
	// Update Platform logic goes here
}

func (p *Platform) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.image, op)
}

func (p *Platform) SetPosition(x, y float64) {
	p.x = x
	p.y = y
}

func (p *Platform) GetPosition() (float64, float64) {
	return p.x, p.y
}

func (p *Platform) GetSize() (int, int) {
	return p.width, p.height
}
