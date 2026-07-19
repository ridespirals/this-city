package sim

import (
	"math"
	"math/rand"
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

func TestPathFollowerDeadEndReverses(t *testing.T) {
	w := NewWorld()
	w.RNG = rand.New(rand.NewSource(1))
	n0 := w.Network.AddNode(Vec2{0, 0})
	n1 := w.Network.AddNode(Vec2{100, 0})
	e := w.Network.AddEdge(n0, n1, Vec2{33, 0}, Vec2{66, 0})
	ent := w.Create()
	w.Transforms.Set(ent, Transform2D{Scale: 1})
	w.Decisions.Set(ent, DefaultPathDecision())
	PlaceOnEdge(w, ent, e, 0, true, 100)
	// Travel past the end — only choice is U-turn.
	f, _ := w.Followers.Get(ent)
	f = AdvancePathFollower(w, ent, f, 1.2)
	if f.Forward {
		t.Fatalf("expected reverse after dead-end, dist=%v edge=%v", f.Distance, f.Edge)
	}
}
