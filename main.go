package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/matthewsedlacek/reverse_gravity/game"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	g *game.Game
)

func update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Handle game exit when the Escape key is pressed
		return ebiten.Termination
	}

	if err := g.Update(); err != nil {
		return err
	}

	g.Draw(screen)

	return nil
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Reverse Gravity")

	g = game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
