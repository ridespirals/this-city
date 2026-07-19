// Package game owns agent/event/schedule systems and ticks the simulation.
// It may import sim but must not import raylib.
package game

import "github.com/ridespirals/this-city/internal/sim"

// Session is the runnable game layer over a sim world.
type Session struct {
	World    *sim.World
	Machines Machines
	Paused   bool
	Time     float64 // accumulated sim time in seconds
}

// NewSession wraps a world for ticking and registers built-in machines.
func NewSession(world *sim.World) *Session {
	sim.EnsureRNG(world, 1)
	return &Session{
		World: world,
		Machines: Machines{
			MachineDebug: DebugMachine(),
			MachineWalk:  WalkMachine(),
		},
	}
}

// SpawnDemo loads the figure-8 sample map and a few random-walking agents.
func (s *Session) SpawnDemo() {
	if s == nil || s.World == nil {
		return
	}
	_ = LoadDemoMap(s.World)
	var first sim.EdgeID
	s.World.Network.ForEachEdge(func(e *sim.Edge) {
		if first == sim.NilEdge {
			first = e.ID
		}
	})
	if first != sim.NilEdge {
		SpawnPathFollower(s.World, first, 120)
		// Second agent on another edge, opposite direction bias via RNG choices.
		var second sim.EdgeID
		s.World.Network.ForEachEdge(func(e *sim.Edge) {
			if second == sim.NilEdge && e.ID != first {
				second = e.ID
			}
		})
		if second != sim.NilEdge {
			SpawnPathFollower(s.World, second, 100)
		}
	}
}

// Tick advances simulation systems by dt seconds when not paused.
func (s *Session) Tick(dt float32) {
	if s == nil || s.Paused {
		return
	}
	s.Time += float64(dt)
	TickBrains(s.World, s.Machines, dt)
	sim.TickPathFollowers(s.World, dt)
	TickEvents(s.World, dt)
}

// SetPaused freezes or resumes simulation ticks.
func (s *Session) SetPaused(paused bool) {
	if s == nil {
		return
	}
	s.Paused = paused
}

// TogglePause flips the pause flag and returns the new value.
func (s *Session) TogglePause() bool {
	if s == nil {
		return false
	}
	s.Paused = !s.Paused
	return s.Paused
}
