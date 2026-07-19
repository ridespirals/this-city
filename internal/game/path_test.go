package game

import (
	"math"
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestDemoFollowerMovesAlongNetwork(t *testing.T) {
	s := NewSession(sim.NewWorld())
	s.SpawnDemo()
	if s.World.Network.EdgeCount() < 1 || s.World.Followers.Len() < 1 {
		t.Fatalf("edges=%d followers=%d", s.World.Network.EdgeCount(), s.World.Followers.Len())
	}

	var e sim.Entity
	var start sim.Transform2D
	s.World.Followers.ForEach(func(ent sim.Entity, _ sim.PathFollower) {
		if e.IsNil() {
			e = ent
			start, _ = s.World.Transforms.Get(ent)
		}
	})

	for i := 0; i < 60; i++ {
		s.Tick(1.0 / 60)
	}
	end, _ := s.World.Transforms.Get(e)
	dist := math.Hypot(float64(end.X-start.X), float64(end.Y-start.Y))
	if dist < 20 {
		t.Fatalf("expected movement after 1s, moved %v", dist)
	}
}
