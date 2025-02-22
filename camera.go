package main

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	// The camera is a lie
	// the camera is just an offset for the player.
	c.X = (screenWidth / 2.0) - targetX
	c.Y = (screenHeight / 2.0) - targetY

	c.Y = 0
	c.X = 0
}
