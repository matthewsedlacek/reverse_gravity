package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	charWidth       = 32
	charHeight      = 32
	platformWidth   = 128
	platformHeight  = 16
	platformSpacing = 200
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
	distance         float64
	altitude         float64
	reachedSpace     bool
	gameOver         bool
	fontFace         font.Face
}

type Platform struct {
	x float64
	y float64
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Reverse Gravity")

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

	// Load the font
	fontData, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	fontFace := truetype.NewFace(fontData, &truetype.Options{
		Size: 12,
	})

	game := &Game{
		charX:            screenWidth / 2,
		charY:            screenHeight - charHeight - platformHeight - 50,
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
		fontFace:        fontFace,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	spacePressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyLeft)
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyRight)
	shiftPressed := inpututil.IsKeyJustPressed(ebiten.KeyShift)

	// Apply gravity
	gravity := 0.1 // Initial gravity value
	if shiftPressed {
		gravity = -1.0 // Reverse gravity if shift key is pressed
	}
	g.charYSpeed += gravity

	// Handle jumping
	if spacePressed && g.onGround && !g.prevSpacePressed {
		g.charYSpeed = -5 // Adjust the initial jump velocity to your preference
		g.onGround = false // Reset onGround status when jumping
	}

	g.prevSpacePressed = spacePressed

	// Apply horizontal movement
	if leftPressed {
		g.charX -= 3 // Adjust the horizontal movement speed to your preference
	}
	if rightPressed {
		g.charX += 3 // Adjust the horizontal movement speed to your preference
	}

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

	// Generate new platform if needed
	if g.charY < screenHeight/2 && len(g.platforms) < 50 {
		g.generateNewPlatform()
	}

	// Game over condition
	if g.charY > screenHeight {
		fmt.Println("Game Over")
		return ebiten.Termination
	}

	// Calculate distance traveled in pixels
	g.distance += math.Abs(g.charYSpeed)

	// Update altitude
	g.altitude = screenHeight - g.charY

	// Check if player reached space
	if !g.reachedSpace && g.altitude >= 10000 {
		g.reachedSpace = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var platformColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green color

	screenX := int(g.charX) - screenWidth/2
	screenY := int(g.charY) - screenHeight/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(screenX), -float64(screenY))

	// Render background images
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			bgX := (screenX / screenWidth) + i
			bgY := (screenY / screenHeight) + j
			op.GeoM.Translate(float64(bgX*screenWidth), float64(bgY*screenHeight))
			screen.DrawImage(g.backgroundImage, op)
			op.GeoM.Translate(-float64(bgX*screenWidth), -float64(bgY*screenHeight))
		}
	}
        
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

	altitudeStr := fmt.Sprintf("Altitude: %.1f", g.altitude)
	ebitenutil.DebugPrint(screen, altitudeStr)

	if g.reachedSpace {
		congratsStr := "Congratulations!\nYou made it to space!"
		// Display game over message
		gameOverWidth := measureTextWidth(congratsStr, g.fontFace)
		ebitenutil.DebugPrintAt(screen, congratsStr, (screenWidth-gameOverWidth)/2, screenHeight/2)
	}

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

func measureTextWidth(text string, face font.Face) int {
	totalWidth := 0

	for _, runeValue := range text {
		advance, ok := face.GlyphAdvance(runeValue)
		if !ok {
			// Handle the case where advance width is not available for the rune
			continue
		}

		width := int(advance.Ceil())
		totalWidth += width
	}

	return totalWidth
}
