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
	brain, ok := w.Brains.Get(a)
	if !ok || brain.Machine != MachineWalk {
		t.Fatalf("expected walk brain, got %+v ok=%v", brain, ok)
	}
	if !DeleteEntity(w, a) || w.Alive(a) {
		t.Fatal("delete failed")
	}
}

func TestSetPathFromAnchors(t *testing.T) {
	w := sim.NewWorld()
	g := SetPathFromAnchors(w, 0, []sim.Vec2{{0, 0}, {50, 0}, {50, 50}})
	if g == 0 || w.Network.EdgeCount() != 2 {
		t.Fatalf("group=%d edges=%d", g, w.Network.EdgeCount())
	}
	g2 := SetPathFromAnchors(w, g, []sim.Vec2{{0, 0}, {10, 0}})
	if g2 != g || w.Network.EdgeCount() != 1 {
		t.Fatalf("group=%d edges=%d", g2, w.Network.EdgeCount())
	}
}
