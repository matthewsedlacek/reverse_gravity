package game

// This file defines the player character, including its movement, collisions, and gravity-related behavior.
import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	image  *ebiten.Image
	x, y   float64
	width  int
	height int
}

func NewPlayer(image *ebiten.Image, x, y float64) *Player {
	width, height := image.Size()
	return &Player{
		image:  image,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (p *Player) Update() {
	// Update player logic goes here
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.image, op)
}

func (p *Player) SetPosition(x, y float64) {
	p.x = x
	p.y = y
}

func (p *Player) GetPosition() (float64, float64) {
	return p.x, p.y
}

func (p *Player) GetSize() (int, int) {
	return p.width, p.height
}
