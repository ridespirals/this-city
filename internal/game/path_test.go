package game

import (
	"math"
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestDemoFollowerMovesAlongPath(t *testing.T) {
	s := NewSession(sim.NewWorld())
	s.SpawnDemo()
	if s.World.Paths.Len() != 1 || s.World.Followers.Len() != 1 {
		t.Fatalf("demo paths=%d followers=%d", s.World.Paths.Len(), s.World.Followers.Len())
	}

	var e sim.Entity
	var start sim.Transform2D
	s.World.Followers.ForEach(func(ent sim.Entity, _ sim.PathFollower) {
		e = ent
		start, _ = s.World.Transforms.Get(ent)
	})

	for i := 0; i < 60; i++ {
		s.Tick(1.0 / 60)
	}
	end, _ := s.World.Transforms.Get(e)
	dx := float64(end.X - start.X)
	dy := float64(end.Y - start.Y)
	dist := math.Hypot(dx, dy)
	if dist < 40 {
		t.Fatalf("expected movement after 1s, moved %v (start=%v end=%v)", dist, start, end)
	}
}
