package levels

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	name       string
	background string
	difficulty int
	levelFile  string
}

func NewLevel(name, background string, difficulty int, levelFile string) *Level {
	return &Level{
		name:       name,
		background: background,
		difficulty: difficulty,
		levelFile:  levelFile,
	}
}

func (l *Level) Load() {
	// Load resources specific to the Earth level
	// For example, load background image, obstacle data, etc.
	fmt.Printf("Loading Earth level: %s\n", l.name)
	// Load background image from l.background
	// Load obstacle data from l.levelFile
}

func (l *Level) Update() error {
	// Update game logic specific to the Earth level
	// For example, update obstacle movement, collision detection, etc.
	fmt.Printf("Updating Earth level: %s\n", l.name)
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	// Draw the Earth level specific elements
	// For example, draw the background, obstacles, etc.
	fmt.Printf("Drawing Earth level: %s\n", l.name)
	// Draw the background image onto the screen
	// Draw the obstacles and other level elements
}

func (l *Level) IsComplete() bool {
	// Check if the Earth level is completed based on your game's criteria
	// For example, check if a goal is reached or challenges are overcome
	return false // Replace with your condition
}
