package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var playerImg *ebiten.Image

func init() {
	fmt.Println("Loading player data")
	pi, _, err := ebitenutil.NewImageFromFile("assets/NinjaSpriteSheet.png")
	playerImg = pi // Why do I need to redeclare the var
	if err != nil {
		log.Fatal(err)
	}
}

type Sprite struct {
	Img    *ebiten.Image
	X, Y   float32 // why a float?
	Dx, Dy float32 // change in x and y, ie velocity, Do i need this if I have target
}

type PlayerState uint8

const (
	// This also feels bad
	// Better if in player go file / package
	Down PlayerState = iota
	Up
	Left
	Right
	DownJump
	UpJump
	LeftJump
	RightJump
)

type Player struct {
	*Sprite
	Frame            int
	TargetX, TargetY float32
	PlayerFace       PlayerState
	Animations       map[PlayerState]*Animation
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
			Down:      NewAnimation(4, 12, 4, 20.0),
			Up:        NewAnimation(5, 13, 4, 20.0),
			Left:      NewAnimation(6, 14, 4, 20.0),
			Right:     NewAnimation(7, 15, 4, 20.0),
			DownJump:  NewAnimation(20, 20, 0, 0.0),
			UpJump:    NewAnimation(21, 21, 0, 0.0),
			LeftJump:  NewAnimation(22, 21, 0, 0.0),
			RightJump: NewAnimation(23, 21, 0, 0.0),
		},
	}
}

func (p *Player) IsWalking() bool {
	return (p.X != p.TargetX || p.Y != p.TargetY)
}

// Bad idea to pass jumped, maybe create an input object or just check it here?
func (p *Player) AnimationState(dx, dy float64, jumped bool) *Animation {
	// dy and dx -> The diffrance between the target and the current position

	// player is standing still
	if dx == 0 && dy == 0 {
		if jumped {
			switch p.PlayerFace {
			case Up:
				return p.Animations[UpJump]
			case Down:
				return p.Animations[DownJump]
			case Left:
				return p.Animations[LeftJump]
			case Right:
				return p.Animations[RightJump]

			}
		}
	}

	//check the longest leg of the path to make player face they way when walking by mouse click
	// Is there better code then this nasty if checker
	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			if jumped {
				p.PlayerFace = RightJump
			} else {
				p.PlayerFace = Right
			}
			return p.Animations[p.PlayerFace]
		}
		if dx < 0 {
			p.PlayerFace = Left
			return p.Animations[Left]
		}
	}
	if dy > 0 {
		p.PlayerFace = Down
		return p.Animations[Down]
	}
	if dy < 0 {
		p.PlayerFace = Up
		return p.Animations[Up]
	}

	return nil // why return nil? Why not always return an animations?
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
