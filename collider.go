package main

import (
	"fmt"
	"image"
	"math"
)

type Colider interface {
	GetHitBox() image.Rectangle // rename?
}

func CheckCollisionHorizontal(sprite *Player, colliders []Colider) {
	// Check needs to happen before update of passed in sprite x cord
	/*
		More or less this prevents collisions for any passed sprite
		I can do 1 loop as the sprite has the "future" state within it self
		maybe a good way to clean that up would be to have functions for the math that requires
		The negtive is that the code is a bit more complex

		The steps are
		create the future x positon of the player
		Get it the larget value in absolute terms ie -101 is larger then -100
			(this prevents rubberbanding effect when walking against a wall)
		create the hitbox rect with only the future x value
		if the overlaps, update player state to prevent movement in the x direction

		do all the for the y after
	*/
	for _, collider := range colliders {
		colliderRect := collider.GetHitBox()

		// hitbox only get the current loc of the player
		// need to know the next one
		futureX := sprite.X + sprite.Dx

		// The hitbox we create uses ints but when the player can travel in floats
		// This means the players future box "could" be 101.5
		// the int cast will drop the .5, this would mean there is no overlap
		// on the next call it would hit bouncing back the cords causing a stutter effect
		// The 2nd issue is that Ceil will target to the GREATEST int,
		// this means we want to floor the value when the speed is negateve
		if sprite.Dx < 0 {
			futureX = float32(math.Floor(float64(futureX)))
		} else {
			futureX = float32(math.Ceil(float64(futureX)))
		}

		// is it better to get the og hitbox or just make a rect here?
		hitbox := sprite.GetHitBox()
		hitbox.Min.X = int(futureX)
		hitbox.Max.X = int(futureX) + sprite.Width

		if colliderRect.Overlaps(hitbox) {
			fmt.Printf("collided on x axis Speed %f  x%f  future x%f\n", sprite.Dx, sprite.X, futureX)
			if sprite.Dx < 0 { // moves right
				sprite.TargetX = float32(colliderRect.Max.X)
			} else if sprite.Dx > 0 { // moves left
				sprite.TargetX = float32(colliderRect.Min.X) - float32(sprite.Width)
			}
			sprite.Dx = 0 // prevent move
			sprite.X = sprite.TargetX
		}

		futureY := sprite.Y + sprite.Dy
		if sprite.Dy < 0 {
			futureY = float32(math.Floor(float64(futureY)))
		} else {
			futureY = float32(math.Ceil(float64(futureY)))
		}
		// like this?
		hitbox = image.Rect(
			int(sprite.X),
			int(futureY),
			int(sprite.X)+sprite.Width,
			int(futureY)+sprite.Height,
		)
		if colliderRect.Overlaps(hitbox) {
			fmt.Printf("collided on Y axis Speed %f\n", sprite.Dy)
			if sprite.Dy < 0 { // moves up
				sprite.TargetY = float32(colliderRect.Max.Y)
			} else if sprite.Dy > 0 { // moves down
				sprite.TargetY = float32(colliderRect.Min.Y) - float32(sprite.Height)
			}
			sprite.Dy = 0 // prevent move
			sprite.Y = sprite.TargetY
		}

	}
}
