package game

import (
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestSpawnAndPick(t *testing.T) {
	w := sim.NewWorld()
	a := SpawnAgent(w, sim.RolePolice, sim.Vec2{X: 10, Y: 10})
	_ = SpawnAgent(w, sim.RoleCivilian, sim.Vec2{X: 100, Y: 100})
	got := PickEntity(w, sim.Vec2{X: 12, Y: 11}, 20)
	if got != a {
		t.Fatalf("pick=%v want %v", got, a)
	}
	if !DeleteEntity(w, a) || w.Alive(a) {
		t.Fatal("delete failed")
	}
}

func TestSetPathFromAnchors(t *testing.T) {
	w := sim.NewWorld()
	id := SetPathFromAnchors(w, sim.NilPath, []sim.Vec2{{0, 0}, {50, 0}, {50, 50}})
	p, ok := w.Paths.Get(id)
	if !ok || len(p.Segments) != 2 {
		t.Fatalf("ok=%v segs=%d", ok, len(p.Segments))
	}
	id2 := SetPathFromAnchors(w, id, []sim.Vec2{{0, 0}, {10, 0}})
	if id2 != id {
		t.Fatal("expected update in place")
	}
	p, _ = w.Paths.Get(id)
	if len(p.Segments) != 1 {
		t.Fatalf("segs=%d", len(p.Segments))
	}
}
