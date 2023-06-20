package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// This file handles the game initialization, updates, and rendering. It orchestrates the different game entities.

type Game struct {
	// Game state and variables
	currentLevel int
}

func NewGame() *Game {
	g := &Game{
		currentLevel: 0,
	}
	return g
}

func (g *Game) Update() error {
	// Update game logic
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Set the dimensions of the game window
	screenWidth = 320
	screenHeight = 240
	return screenWidth, screenHeight
}
