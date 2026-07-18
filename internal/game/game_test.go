package game

import (
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestSessionTickRespectsPause(t *testing.T) {
	s := NewSession(sim.NewWorld())
	s.Tick(0.5)
	if s.Time != 0.5 {
		t.Fatalf("Time = %v, want 0.5", s.Time)
	}
	s.SetPaused(true)
	s.Tick(1)
	if s.Time != 0.5 {
		t.Fatalf("paused Time = %v, want 0.5", s.Time)
	}
	s.TogglePause()
	s.Tick(0.25)
	if s.Time != 0.75 {
		t.Fatalf("resumed Time = %v, want 0.75", s.Time)
	}
}
