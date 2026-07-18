// Package game owns agent/event/schedule systems and ticks the simulation.
// It may import sim but must not import raylib.
package game

import "github.com/ridespirals/this-city/internal/sim"

// Session is the runnable game layer over a sim world.
type Session struct {
	World  *sim.World
	Paused bool
	Time   float64 // accumulated sim time in seconds
}

// NewSession wraps a world for ticking.
func NewSession(world *sim.World) *Session {
	return &Session{World: world}
}

// Tick advances simulation systems by dt seconds when not paused.
func (s *Session) Tick(dt float32) {
	if s == nil || s.Paused {
		return
	}
	s.Time += float64(dt)
	// Phase 3+: run ECS / FSM / path systems here.
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
