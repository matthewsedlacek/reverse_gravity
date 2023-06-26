package game

import (
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matthewsedlacek/reverse_gravity/levels"
)

type Game struct {
	currentLevel int
	levels       []*levels.Level
	loading      bool
	loadingImage *ebiten.Image
}

func NewGame() *Game {
	g := &Game{
		currentLevel: 0,
		levels: []*levels.Level{
			levels.NewLevel("Space", "assets/images/space_background.png", 1.0, "space/level1.txt"),
			levels.NewLevel("Earth", "assets/images/earth_background.png", 2.0, "earth/level1.txt"),
			// Add more levels as needed
		},
		loading:      true,
		loadingImage: loadImage("assets/images/loading_screen.png"),
	}
	return g
}

func (g *Game) Update() error {
	if g.loading {
		// Perform loading operations here
		// Once loading is complete, set g.loading to false
		// and transition to the first level
		g.levels[g.currentLevel].Load()
		g.loading = false
		return nil
	}

	level := g.levels[g.currentLevel]
	if level.IsComplete() {
		g.currentLevel++
		if g.currentLevel >= len(g.levels) {
			// Game completed, show game over screen or handle end condition
			return ebiten.Termination
		}
		// Load the next level
		g.loading = true
		return nil
	}

	// Update level-specific logic
	if err := level.Update(); err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.loading {
		// Draw loading screen
		screen.DrawImage(g.loadingImage, nil)
		return
	}

	level := g.levels[g.currentLevel]
	// Draw level-specific elements
	level.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Specify the size of the game's logical screen
	return 640, 480
}

// loadImage loads an image file and returns an ebiten.Image.
func loadImage(path string) *ebiten.Image {
	_, image, err := ebitenutil.NewImageFromFile(path, 0)
	if err != nil {
		panic(err)
	}

	ebitenImg := ebiten.NewImageFromImage(image)

	return ebitenImg
}
