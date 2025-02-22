package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Player *Player
	Mob    *Sprite
	Box    *Sprite
	Cam    *Camera
}

func NewGame() *Game {

	// v := image.Rect(100, 100, 200, 200).(*ebiten.Image)
	img := ebiten.NewImage(50, 30)
	fmt.Printf("%+v", img)
	vector.StrokeRect(
		img,
		float32(img.Bounds().Min.X), float32(img.Bounds().Min.Y),
		float32(img.Bounds().Dx()), float32(img.Bounds().Dy()),
		1.0,
		color.RGBA{255, 0, 0, 255},
		true)

	return &Game{
		Player: NewPlayer(),
		Mob:    NewSprite(skelly),
		Box: &Sprite{
			Img:    img,
			Width:  50,
			Height: 30,
			X:      320 / 2, Y: 240 / 2,
			Dx: 0, Dy: 0,
			TargetX: 320 / 2, TargetY: 240 / 2,
			speed: 2.0,
		},
		Cam: NewCamera(0, 0),
	}
}

func (g *Game) Update() error {

	g.Player.Update(g)
	// update cam
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

	opts.GeoM.Reset()

	opts.GeoM.Translate(float64(g.Box.X), float64(g.Box.Y))
	opts.GeoM.Translate(g.Cam.X, g.Cam.Y)

	// g.Box.Img.ColorModel().Convert(color.RGBA{255, 0, 0, 255})
	screen.DrawImage(
		g.Box.Img,
		&opts,
	)

	x, y := ebiten.CursorPosition()

	msg := fmt.Sprintf("fps :%0.2f\n", ebiten.ActualFPS())
	msg += fmt.Sprintf("mouse x y: %d %d\n", x, y)
	msg += fmt.Sprintf("player x y:  %f %f\n", g.Player.X, g.Player.Y)
	msg += fmt.Sprintf("target x y:  %f %f\n", g.Player.TargetX, g.Player.TargetY)
	msg += fmt.Sprintf("speed x y:  %f %f\n", g.Player.Dx, g.Player.Dy)
	msg += fmt.Sprintf("player frame: %d\n", g.Player.Frame)
	ebitenutil.DebugPrint(screen, msg)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
