package sim

import (
	"math"
	"testing"
)

func TestCubicBezierEndpoints(t *testing.T) {
	b := CubicBezier{
		P0: Vec2{0, 0},
		C0: Vec2{1, 2},
		C1: Vec2{2, 2},
		P1: Vec2{3, 0},
	}
	p0 := b.Point(0)
	p1 := b.Point(1)
	if p0 != b.P0 || p1 != b.P1 {
		t.Fatalf("endpoints: got %v %v", p0, p1)
	}
}

func TestPolylineSampleAtConstantSpeedLength(t *testing.T) {
	// Straight-ish cubic that is nearly a line from (0,0) to (100,0).
	segs := []CubicBezier{{
		P0: Vec2{0, 0},
		C0: Vec2{33, 0},
		C1: Vec2{66, 0},
		P1: Vec2{100, 0},
	}}
	pl := BuildPolyline(segs, 20)
	if pl.Length < 99 || pl.Length > 101 {
		t.Fatalf("length = %v, want ~100", pl.Length)
	}
	pos, tan := pl.SampleAt(pl.Length * 0.5)
	if math.Abs(float64(pos.X-50)) > 2 || math.Abs(float64(pos.Y)) > 1 {
		t.Fatalf("mid pos = %v", pos)
	}
	if tan.X <= 0 {
		t.Fatalf("tangent should point +X, got %v", tan)
	}
}

func TestPathFollowerPingPong(t *testing.T) {
	w := NewWorld()
	id := w.Paths.Add([]CubicBezier{{
		P0: Vec2{0, 0},
		C0: Vec2{50, 0},
		C1: Vec2{50, 0},
		P1: Vec2{100, 0},
	}})
	path, _ := w.Paths.Get(id)
	e := w.Create()
	w.Transforms.Set(e, Transform2D{Scale: 1})
	f := PathFollower{Path: id, Speed: 100, Forward: true, PingPong: true}
	// Travel past the end in one big step.
	f = AdvancePathFollower(w, e, f, 1.2)
	if f.Forward {
		t.Fatalf("expected reverse after overshoot, distance=%v len=%v", f.Distance, path.Poly.Length)
	}
	xf, _ := w.Transforms.Get(e)
	if xf.X < 0 || xf.X > 100 {
		t.Fatalf("position out of range: %v", xf)
	}
}
