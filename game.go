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
	Mob    *Sprite
	Cam    *Camera
}

func NewGame() *Game {

	return &Game{
		Player: NewPlayer(),
		Mob:    NewSprite(),
		Cam:    NewCamera(0, 0),
	}
}

func (g *Game) Update() error {

	// reset player velocity, if player had any.
	g.Player.Dx = 0 // both are unused remove?
	g.Player.Dy = 0

	// each tick we need to figure out the next frame we want to show the user.
	g.Player.Frame = 0

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

		// Ha this kinda works.
		// Somehow move into player?
		if g.Player.JumpTPS == 0 {
			// 20 seconds is how long is
			// 20 seconds for the cooldown
			g.Player.JumpTPS = 40
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Fires once per click
		x, y := ebiten.CursorPosition()
		fmt.Printf("x%d y%d\n", x, y)
		x -= int(g.Cam.X) + 8 // half of tile size
		y -= int(g.Cam.Y) + 8
		g.Player.TargetX = float32(x)
		g.Player.TargetY = float32(y)
	}

	// Compute the vector from the player's current position to the target.
	// dx and dy are not velocities, they are distance, rename?
	dx := float64(g.Player.TargetX - g.Player.X)
	dy := float64(g.Player.TargetY - g.Player.Y)

	a := g.Player.AnimationState(dx, dy)
	// We should always an AnimationState the player is in
	a.Update()
	g.Player.Frame = a.frame

	// Calculate the distance to the target.
	dist := math.Hypot(dx, dy)

	// If the player is close enough to the target, snap to it.
	speed := 2.5 //todo: move into sprite or player object
	if dist < speed {
		g.Player.X = g.Player.TargetX
		g.Player.Y = g.Player.TargetY
	} else {
		// Normalize the vector and move the player by speed
		g.Player.X += float32((dx / dist) * speed)
		g.Player.Y += float32((dy / dist) * speed)
	}

	g.Cam.FollowTarget(float64(g.Player.X)+8, float64(g.Player.Y)+8, 320, 240)

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
	opts.GeoM.Translate(g.Cam.X, g.Cam.Y) // move the player by cam offset

	screen.DrawImage(
		// grab a subimage of the spritesheet
		getSpriteByIndex(g.Player.Frame, g.Player.Img), &opts,
	)

	opts.GeoM.Reset()

	opts.GeoM.Translate(float64(g.Mob.X), float64(g.Mob.Y))
	// The cam x,y cords are the offset for all other entites
	// relitive to the target
	opts.GeoM.Translate(g.Cam.X, g.Cam.Y)

	screen.DrawImage(
		getSpriteByIndex(0, g.Mob.Img),
		&opts,
	)

	x, y := ebiten.CursorPosition()

	msg := fmt.Sprintf("fps :%0.2f\n", ebiten.ActualFPS())
	msg += fmt.Sprintf("mouse x y: %d %d\n", x, y)
	msg += fmt.Sprintf("player x y:  %f %f\n", g.Player.X, g.Player.Y)
	msg += fmt.Sprintf("player frame: %d\n", g.Player.Frame)
	ebitenutil.DebugPrint(screen, msg)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
