package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	skelly *ebiten.Image
)

func init() {
	sp, _, err := ebitenutil.NewImageFromFile("assets/SkeletonSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	skelly = sp
}

type Sprite struct {
	Img              *ebiten.Image // sprite sheet
	Width, Height    int           // sprite sheet might have diffrent images or might just be 1 large one
	X, Y             float32       // why a float?
	Dx, Dy           float32       // change in x and y, ie velocity, Do i need this if I have target
	TargetX, TargetY float32
	speed            float32 // rename it
}

func NewSprite(img *ebiten.Image) *Sprite {
	return &Sprite{
		Img:    img,
		Width:  16, // x
		Height: 16, // y
		X:      100, Y: 100,
		Dx: 0, Dy: 0,
		TargetX: 0, TargetY: 40,
		speed: 2.0,
	}
}

func (s *Sprite) GetHitBox() image.Rectangle {
	return image.Rect(
		int(s.X),
		int(s.Y),
		int(s.X)+s.Width,
		int(s.Y)+s.Height,
	)
}

type Animation struct {
	First, Last  int
	Step         int
	SpeedInTps   float32 // Tps -> ticks per second
	frameCounter float32 // why a float, how fast is a tick?
	frame        int     // current frame
}

func NewAnimation(first, last, step int, speedInTps float32) *Animation {
	return &Animation{
		first,
		last,
		step,
		speedInTps,
		0,
		first,
	}
}

func (a *Animation) Update() {

	a.frameCounter += 1.0 // again why a float? to allow for faster animations?
	if a.frameCounter >= a.SpeedInTps {
		a.frameCounter = 0 // reset counter
		a.frame += a.Step  // next frame
		if a.frame > a.Last {
			// loop back to the beginning
			a.frame = a.First
		}
	}
}

// Why keep frame private here? Attach to player?
func (a *Animation) Frame() int {
	return a.frame
}
