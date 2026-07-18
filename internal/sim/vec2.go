package sim

import "math"

// Vec2 is a 2D vector / point in world space.
type Vec2 struct {
	X, Y float32
}

// Add returns a + b.
func (a Vec2) Add(b Vec2) Vec2 { return Vec2{a.X + b.X, a.Y + b.Y} }

// Sub returns a - b.
func (a Vec2) Sub(b Vec2) Vec2 { return Vec2{a.X - b.X, a.Y - b.Y} }

// Scale returns a * s.
func (a Vec2) Scale(s float32) Vec2 { return Vec2{a.X * s, a.Y * s} }

// Dot returns the dot product.
func (a Vec2) Dot(b Vec2) float32 { return a.X*b.X + a.Y*b.Y }

// Len returns the Euclidean length.
func (a Vec2) Len() float32 {
	return float32(math.Hypot(float64(a.X), float64(a.Y)))
}

// Normalize returns a unit vector, or zero if a is near-zero.
func (a Vec2) Normalize() Vec2 {
	l := a.Len()
	if l < 1e-6 {
		return Vec2{}
	}
	return a.Scale(1 / l)
}

// Angle returns the heading in radians (atan2).
func (a Vec2) Angle() float32 {
	return float32(math.Atan2(float64(a.Y), float64(a.X)))
}
