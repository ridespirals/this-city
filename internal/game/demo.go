package game

import "github.com/ridespirals/this-city/internal/sim"

// DemoPath builds a hard-coded S-curve street for Phase 4.
func DemoPath() []sim.CubicBezier {
	return []sim.CubicBezier{
		{
			P0: sim.Vec2{X: 160, Y: 480},
			C0: sim.Vec2{X: 320, Y: 480},
			C1: sim.Vec2{X: 320, Y: 200},
			P1: sim.Vec2{X: 560, Y: 200},
		},
		{
			P0: sim.Vec2{X: 560, Y: 200},
			C0: sim.Vec2{X: 800, Y: 200},
			C1: sim.Vec2{X: 800, Y: 520},
			P1: sim.Vec2{X: 1120, Y: 520},
		},
	}
}

// SpawnPathFollower creates an agent that follows pathID at the given speed.
func SpawnPathFollower(w *sim.World, pathID sim.PathID, speed float32) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Roles.Set(e, sim.RoleTag{Role: sim.RoleDebug})
	w.Followers.Set(e, sim.PathFollower{
		Path:     pathID,
		Distance: 0,
		Speed:    speed,
		Forward:  true,
		PingPong: true,
	})
	// Snap transform to path start.
	f, _ := w.Followers.Get(e)
	w.Transforms.Set(e, sim.Transform2D{Scale: 1})
	w.Followers.Set(e, sim.AdvancePathFollower(w, e, f, 0))
	return e
}
