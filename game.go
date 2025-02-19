package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Player *Player
}

func NewGame() *Game {

	return &Game{
		Player: NewPlayer(),
	}

}

func (g *Game) Update() error {

	// reset player velocity, it it had any. Target is the goal the player much reach
	g.Player.Dx = 0
	g.Player.Dy = 0
	// each tick we need to figure out the next frame we want to show the user.
	g.Player.Frame = 0

	// TODO: replace jumped bool with better system
	jumped := false
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.TargetX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.TargetX += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.TargetY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.TargetY += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		jumped = true
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Fires once per click
		x, y := ebiten.CursorPosition()
		fmt.Printf("x%d y%d\n", x, y)
		g.Player.TargetX = float32(x)
		g.Player.TargetY = float32(y)
	}

	// Compute the vector from the player's current position to the target.
	dx := float64(g.Player.TargetX - g.Player.X)
	dy := float64(g.Player.TargetY - g.Player.Y)

	// Handle the way the player is facing
	// This is only movement
	// dx and dy are not velocities, they are distance, rename?

	// Need to check if the player is moving
	// if g.Player.IsWalking() {
	// var a *Animation
	// if jumped {
	// 	a = g.Player.Animations[Jump]
	// } else {
	a := g.Player.AnimationState(dx, dy, jumped)
	// }

	if a != nil {
		a.Update()
		// Should Player.Frame be public?
		// should the frame be set by AnimationState?
		g.Player.Frame = a.frame
	} else {
		g.Player.Frame = int(g.Player.PlayerFace)
	}

	// Calculate the distance to the target.
	dist := math.Hypot(dx, dy)

	// If the player is close enough to the target, snap to it.
	speed := 2.5 // move into sprite or player object
	if dist < speed {
		g.Player.X = g.Player.TargetX
		g.Player.Y = g.Player.TargetY
	} else {
		// Normalize the vector and move the player by speed
		g.Player.X += float32((dx / dist) * speed)
		g.Player.Y += float32((dy / dist) * speed)
	}
	return nil
}

func getSpriteByIndex(index int, img *ebiten.Image) *ebiten.Image {
	Tilesize := 16
	WidthInTiles := 4
	// formula to index sprites like an array.
	// idx 0 -> 4  left to right
	x := (index % WidthInTiles) * Tilesize
	y := (index / WidthInTiles) * Tilesize

	var rect image.Rectangle = image.Rect(x, y, x+Tilesize, y+Tilesize)

	return img.SubImage(rect).(*ebiten.Image)
}

func (g *Game) Draw(screen *ebiten.Image) {

	// fill the screen with a nice sky color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	// set the translation of our drawImageOptions to the player's position
	opts.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))

	screen.DrawImage(
		// grab a subimage of the spritesheet
		getSpriteByIndex(g.Player.Frame, g.Player.Img), &opts,
	)

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
