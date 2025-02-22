package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayerState uint8
type PlayerDirection uint8

const (
	Down PlayerDirection = iota
	Up
	Left
	Right
)

const (
	// Feels like this all should be in a package
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
	Frame           int
	PlayerDirection PlayerDirection
	PlayerFace      PlayerState
	Animations      map[PlayerState]*Animation
	JumpTPS         uint
	// store the last state? pushdown automata
	// This way I can get the rest frame
}

var playerImg *ebiten.Image

func init() {
	// TODO: Make a list of sprites and maybe expected types
	fmt.Println("Loading Player data")
	pi, _, err := ebitenutil.NewImageFromFile("assets/NinjaSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	playerImg = pi // Why cant I assign directly
}

func NewPlayer() *Player {

	fmt.Printf("%+v", playerImg)
	return &Player{
		Sprite: &Sprite{
			Img:     playerImg,
			Width:   16,
			Height:  16,
			X:       10,
			Y:       10,
			TargetX: 10,
			TargetY: 10,
			speed:   1.5,
		},
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

// I need some kind of input mapper or maybe an event queue
func (p *Player) Update(g *Game) {

	// what do I need here, The keys pressed? The cam?
	// reset player velocity, if player had any.
	p.Dx = 0 // both are unused remove?
	p.Dy = 0

	// each tick we need to figure out the next frame we want to show the user.
	p.Frame = 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.TargetX -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.TargetX += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.TargetY -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.TargetY += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {

		// Ha this kinda works.
		// Somehow move into player?
		if p.JumpTPS == 0 {
			// 20 seconds is how long is
			// 20 seconds for the cooldown
			p.JumpTPS = 40
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Fires once per click
		x, y := ebiten.CursorPosition()
		fmt.Printf("x%d y%d\n", x, y)
		//TODO: I need to off set the where I clicked by the camera cords
		// Should player get this data passed? Is the cam globle state?
		x -= int(g.Cam.X) + 8 // half of tile size
		y -= int(g.Cam.Y) + 8
		p.TargetX = float32(x)
		p.TargetY = float32(y)
	}

	// Compute the vector from the player's current position to the target.
	// dx and dy are not velocities, they are distance, rename?
	dx := float64(p.TargetX - p.X)
	dy := float64(p.TargetY - p.Y)

	a := p.AnimationState(dx, dy)
	// We should always an AnimationState the player is in
	a.Update()
	p.Frame = a.frame

	// Calculate the distance to the target.
	dist := math.Hypot(dx, dy)

	// If the player is close enough to the target, snap to it.
	if dist < float64(p.speed) {
		p.Dx = float32(dx)
		p.Dy = float32(dy)
	} else {
		// Normalize the vector and move the player by speed
		m := func(dd float64) float32 {
			// Was debuging/messing around to figure out stutter. Not sure this func is needed anymore
			val := float32(dd/dist) * p.speed
			// val = float32(math.Ceil(float64(val)))
			return val
		}
		p.Dx = m(dx)
		p.Dy = m(dy)
	}

	colliders := []Colider{g.Mob, g.Box}
	// Needs the state before the player moves
	CheckCollisionHorizontal(p, colliders)

	// update player?
	p.X += p.Dx
	p.Y += p.Dy
}
