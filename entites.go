package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	playerImg *ebiten.Image
	skelly    *ebiten.Image
)

func init() {
	// TODO: Make a list of sprites and maybe expected types
	fmt.Println("Loading entity data")
	pi, _, err := ebitenutil.NewImageFromFile("assets/NinjaSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	playerImg = pi // Why do I need to redeclare the var

	sp, _, err := ebitenutil.NewImageFromFile("assets/SkeletonSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	skelly = sp
}

type Sprite struct {
	Img    *ebiten.Image
	X, Y   float32 // why a float?
	Dx, Dy float32 // change in x and y, ie velocity, Do i need this if I have target
}

func NewSprite() *Sprite {
	return &Sprite{
		Img: skelly,
		X:   50, Y: 50,
		Dx: 0, Dy: 50,
	}
}

type PlayerState uint8
type PlayerDirection uint8

const (
	Down PlayerDirection = iota
	Up
	Left
	Right
)

const (
	// This also feels bad
	// Better if in player go file / package
	DownStand PlayerState = iota
	UpStand
	LeftStand
	RightStand

	DownWalk
	UpWalk
	LeftWalk
	RightWalk

	DownJump
	UpJump
	LeftJump
	RightJump
)

type Player struct {
	*Sprite
	Frame            int
	TargetX, TargetY float32
	PlayerDirection  PlayerDirection
	PlayerFace       PlayerState
	Animations       map[PlayerState]*Animation
	JumpTPS          uint
	// store the last state? pushdown automata
	// This way I can get the rest frame
}

func NewPlayer() *Player {
	return &Player{
		Sprite: &Sprite{
			Img: playerImg,
			X:   10,
			Y:   10,
		},
		TargetX: 10,
		TargetY: 10,

		// Is there a better approch to handling all the animations that exist?
		Animations: map[PlayerState]*Animation{
			// The player is always animating
			DownStand:  NewAnimation(0, 0, 0, 0.0),
			UpStand:    NewAnimation(1, 0, 0, 0.0),
			LeftStand:  NewAnimation(2, 0, 0, 0.0),
			RightStand: NewAnimation(3, 0, 0, 0.0),

			DownWalk:  NewAnimation(4, 12, 4, 20.0),
			UpWalk:    NewAnimation(5, 13, 4, 20.0),
			LeftWalk:  NewAnimation(6, 14, 4, 20.0),
			RightWalk: NewAnimation(7, 15, 4, 20.0),

			DownJump:  NewAnimation(20, 20, 0, 0.0),
			UpJump:    NewAnimation(21, 21, 0, 0.0),
			LeftJump:  NewAnimation(22, 22, 0, 0.0),
			RightJump: NewAnimation(23, 23, 0, 0.0),
		},
	}
}

func (p *Player) IsWalking() bool {
	return (p.X != p.TargetX || p.Y != p.TargetY)
}

func (p *Player) AnimationState(dx, dy float64) *Animation {
	// dy and dx -> The difference between the target and the current position

	if dx == 0 && dy == 0 {
		switch p.PlayerDirection {
		case Up:
			p.PlayerFace = UpStand
		case Down:
			p.PlayerFace = DownStand
		case Left:
			p.PlayerFace = LeftStand
		case Right:
			p.PlayerFace = RightStand
		}
		// check the longest leg to set PlayerDirection, requried for mouse clicks
		// Is there better code then this nasty if checker
	} else if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			p.PlayerDirection = Right
			p.PlayerFace = RightWalk
		}
		if dx < 0 {
			p.PlayerDirection = Left
			p.PlayerFace = LeftWalk
		}
	} else {
		if dy > 0 {
			p.PlayerDirection = Down
			p.PlayerFace = DownWalk
		}
		if dy < 0 {
			p.PlayerDirection = Up
			p.PlayerFace = UpWalk
		}
	}
	// Player can jump at any point for now
	if p.JumpTPS > 20 {
		switch p.PlayerDirection {
		case Up:
			p.PlayerFace = UpJump
		case Down:
			p.PlayerFace = DownJump
		case Left:
			p.PlayerFace = LeftJump
		case Right:
			p.PlayerFace = RightJump
		}
	}

	if p.JumpTPS > 0 {
		p.JumpTPS -= 1
	}

	return p.Animations[p.PlayerFace] // why return nil? Why not always return an animations?
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
