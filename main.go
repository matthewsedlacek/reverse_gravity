package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	charWidth       = 32
	charHeight      = 32
	platformWidth   = 128
	platformHeight  = 16
	platformSpacing = 200
	finishLineMiles = 100
)

type Game struct {
	charX            float64
	charY            float64
	charYSpeed       float64
	onGround         bool
	prevSpacePressed bool
	platforms        []*Platform
	backgroundImage  *ebiten.Image
	charImage        *ebiten.Image
	audioContext     *audio.Context
	jumpPlayer       *audio.Player
	finishLineMiles  float64
}

type Platform struct {
	x float64
	y float64
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Reverse Gravity")

	jumpPlayer, err := loadAudioPlayer("assets/sounds/slime_jump.wav")
	if err != nil {
		log.Fatal(err)
	}

	// Load character image
	_, backImage, err := ebitenutil.NewImageFromFile("assets/images/space_background.png")
	if err != nil {
		panic(err)
	}

	// Load character image
	_, characterImage, err := ebitenutil.NewImageFromFile("assets/images/scientist.png")
	if err != nil {
		panic(err)
	}

	backgroundImage := ebiten.NewImageFromImage(backImage)
	charImage := ebiten.NewImageFromImage(characterImage)

	rand.Seed(time.Now().UnixNano())

	game := &Game{
		charX:            screenWidth / 2,
		charY:            screenHeight - charHeight,
		charYSpeed:       0,
		onGround:         true,
		prevSpacePressed: false,
		platforms: []*Platform{
			{
				x: screenWidth/2 - platformWidth/2,
				y: screenHeight - platformHeight - charHeight,
			},
			{
				x: screenWidth/2 - platformWidth/2 - 100,
				y: screenHeight - platformHeight - charHeight - 100,
			},
			{
				x: screenWidth/2 - platformWidth/2 + 100,
				y: screenHeight - platformHeight - charHeight - 200,
			},
			{
				x: screenWidth/2 - platformWidth/2 - 100,
				y: screenHeight - platformHeight - charHeight - 300,
			},
		},
		backgroundImage: backgroundImage,
		charImage:       charImage,
		jumpPlayer:      jumpPlayer,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	spacePressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)

	// Apply gravity
	if !g.onGround {
		g.charYSpeed += 0.5
	}

	// Handle jumping
	if spacePressed && g.onGround && !g.prevSpacePressed {
		g.charYSpeed = -12
		g.jumpPlayer.Rewind()
		g.jumpPlayer.Play()
	}

	g.prevSpacePressed = spacePressed

	// Apply vertical speed
	g.charY += g.charYSpeed

	// Check collisions with platforms
	for _, platform := range g.platforms {
		if g.charX+charWidth >= platform.x && g.charX <= platform.x+platformWidth &&
			g.charY+charHeight >= platform.y && g.charY+charHeight <= platform.y+platformHeight &&
			g.charYSpeed > 0 {
			g.charY = platform.y - charHeight
			g.charYSpeed = 0
			g.onGround = true
		}
	}

	// Update onGround status
	if g.charY+charHeight < screenHeight {
		g.onGround = false
	}

	// Generate new platform if needed
	if g.charY < screenHeight/2 && len(g.platforms) < 8 {
		g.generateNewPlatform()
	}

	// Game over condition
	if g.charY > screenHeight {
		fmt.Println("Game Over")
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var platformColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green color

	screenX := int(g.charX) - screenWidth/2
	screenY := int(g.charY) - screenHeight/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(screenX), -float64(screenY))
	screen.DrawImage(g.backgroundImage, op)

	for _, platform := range g.platforms {
		platformX := platform.x - float64(screenX)
		platformY := platform.y - float64(screenY)

		// Draw colored square for the platform
		platformRect := image.Rect(int(platformX), int(platformY), int(platformX+platformWidth), int(platformY+platformHeight))
		ebitenutil.DrawRect(screen, float64(platformRect.Min.X), float64(platformRect.Min.Y), float64(platformRect.Dx()), float64(platformRect.Dy()), platformColor)
	}

	charOp := &ebiten.DrawImageOptions{}
	charOp.GeoM.Translate(screenWidth/2-charWidth/2, screenHeight/2-charHeight/2)
	screen.DrawImage(g.charImage, charOp)

	progressStr := fmt.Sprintf("Progress: %.1f miles / %.1f miles", g.charY/10, finishLineMiles)
	ebitenutil.DebugPrint(screen, progressStr)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) generateNewPlatform() {
	highestPlatform := g.platforms[0]
	for _, platform := range g.platforms {
		if platform.y < highestPlatform.y {
			highestPlatform = platform
		}
	}

	newPlatform := &Platform{
		x: g.randomPlatformX(),
		y: highestPlatform.y - platformSpacing,
	}

	g.platforms = append(g.platforms, newPlatform)
}

func (g *Game) randomPlatformX() float64 {
	return rand.Float64()*(screenWidth-platformWidth) + platformWidth/2
}

func loadAudioPlayer(filepath string) (*audio.Player, error) {

	// Create audio context
	audioContext, err := audio.NewContext(48000)
	if err != nil {
		panic(err)
	}

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

	return jumpPlayer, nil
}
