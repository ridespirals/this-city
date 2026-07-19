package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/sim"
)

// Camera is an overhead 2D camera (pan / zoom).
type Camera struct {
	Inner rl.Camera2D
}

// NewCamera returns a camera centered on the default window.
func NewCamera() *Camera {
	return &Camera{
		Inner: rl.Camera2D{
			Offset:   rl.NewVector2(DefaultWidth/2, DefaultHeight/2),
			Target:   rl.NewVector2(DefaultWidth/2, DefaultHeight/2),
			Rotation: 0,
			Zoom:     1,
		},
	}
}

// ScreenToWorld converts screen pixels to world coordinates.
func (c *Camera) ScreenToWorld(screenX, screenY float32) sim.Vec2 {
	if c == nil {
		return sim.Vec2{X: screenX, Y: screenY}
	}
	w := rl.GetScreenToWorld2D(rl.NewVector2(screenX, screenY), c.Inner)
	return sim.Vec2{X: w.X, Y: w.Y}
}

// Pan shifts the camera target by screen-space delta (divided by zoom).
func (c *Camera) Pan(dx, dy float32) {
	if c == nil {
		return
	}
	z := c.Inner.Zoom
	if z < 0.01 {
		z = 0.01
	}
	c.Inner.Target.X -= dx / z
	c.Inner.Target.Y -= dy / z
}

// ZoomAt multiplies zoom, clamped, keeping screen point stable.
func (c *Camera) ZoomAt(screenX, screenY, factor float32) {
	if c == nil || factor == 0 {
		return
	}
	before := c.ScreenToWorld(screenX, screenY)
	c.Inner.Zoom *= factor
	if c.Inner.Zoom < 0.25 {
		c.Inner.Zoom = 0.25
	}
	if c.Inner.Zoom > 4 {
		c.Inner.Zoom = 4
	}
	after := c.ScreenToWorld(screenX, screenY)
	c.Inner.Target.X += before.X - after.X
	c.Inner.Target.Y += before.Y - after.Y
}

// Begin applies the camera for world-space drawing.
func (c *Camera) Begin() {
	if c == nil {
		return
	}
	rl.BeginMode2D(c.Inner)
}

// End stops camera-space drawing.
func (c *Camera) End() {
	rl.EndMode2D()
}
