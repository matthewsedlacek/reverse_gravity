package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	charWidth       = 32
	charHeight      = 32
	charMoveSpeed   = 4.0
	charJumpSpeed   = -8.0
	gravity         = 0.1
	platformWidth   = 200
	platformHeight  = 50
	platformSpacing = 150
)

type Game struct {
	charImage        *ebiten.Image
	platformImage    *ebiten.Image
	jumpSound        audio.Player
	charX, charY     float64
	charXSpeed       float64
	charYSpeed       float64
	platforms        []*Platform
	prevSpacePressed bool
	onGround         bool
}
type Platform struct {
	x, y float64
}

func (g *Game) Update() error {
	// Quit the game if the Escape key is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Move left and right
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.charXSpeed = -charMoveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.charXSpeed = charMoveSpeed
	} else {
		g.charXSpeed = 0
	}

	// Jump logic
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)

	if spacePressed && !g.prevSpacePressed && g.onGround {
		g.charYSpeed = charJumpSpeed
		// Play jump sound
		g.jumpSound.Rewind()
		g.jumpSound.Play()
	}

	// Apply gravity
	g.charYSpeed += gravity

	// Update character position
	g.charX += g.charXSpeed
	g.charY += g.charYSpeed

	// Collision detection with platforms
	g.onGround = false

	for _, platform := range g.platforms {
		if g.charY+charHeight >= platform.y && g.charY <= platform.y+platformHeight &&
			g.charX+charWidth >= platform.x && g.charX <= platform.x+platformWidth {
			// Character collided with platform
			g.onGround = true
			g.charYSpeed = 0
			g.charY = platform.y - charHeight
		}
	}

	g.prevSpacePressed = spacePressed

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Draw the character
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.charX, g.charY)
	screen.DrawImage(g.charImage, op)

	// Draw the platforms
	for _, platform := range g.platforms {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(platform.x, platform.y)
		screen.DrawImage(g.platformImage, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Reverse Gravity")

	// Create audio context
	audioContext, err := audio.NewContext(48000)
	if err != nil {
		panic(err)
	}

	// Load character and platform images
	_, characterImage, err := ebitenutil.NewImageFromFile("assets/images/scientist.png", 0)
	if err != nil {
		panic(err)
	}
	_, platformsImage, err := ebitenutil.NewImageFromFile("assets/images/platform.png", 0)
	if err != nil {
		panic(err)
	}

	charImage := ebiten.NewImageFromImage(characterImage)
	platformImage := ebiten.NewImageFromImage(platformsImage)

	// Load jump sound
	jumpFile, err := os.Open("assets/sounds/slime_jump.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer jumpFile.Close()

	jumpDecoded, err := wav.Decode(audioContext, jumpFile)
	if err != nil {
		log.Fatal(err)
	}

	jumpPlayer, err := audio.NewPlayer(audioContext, jumpDecoded)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		charImage:     charImage,
		platformImage: platformImage,
		jumpSound:     *jumpPlayer,
		charX:         100,
		charY:         screenHeight - platformHeight - 50,
		platforms: []*Platform{
			{100, screenHeight - platformHeight - 50},
			{300, screenHeight - platformHeight - 100},
			{200, screenHeight - platformHeight - 200},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
