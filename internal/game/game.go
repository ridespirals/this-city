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
	return &Session{
		World: world,
		Machines: Machines{
			MachineDebug: DebugMachine(),
		},
	}
}

// SpawnDemo adds a Phase-4 path and a ping-pong path follower.
func (s *Session) SpawnDemo() {
	if s == nil || s.World == nil {
		return
	}
	id := s.World.Paths.Add(DemoPath())
	SpawnPathFollower(s.World, id, 140)
}

// Tick advances simulation systems by dt seconds when not paused.
func (s *Session) Tick(dt float32) {
	if s == nil || s.Paused {
		return
	}
	s.Time += float64(dt)
	TickBrains(s.World, s.Machines, dt)
	sim.TickPathFollowers(s.World, dt)
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
