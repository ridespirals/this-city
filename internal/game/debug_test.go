package game

import (
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestDebugAgentTogglesState(t *testing.T) {
	s := NewSession(sim.NewWorld())
	e := SpawnDebugAgent(s.World, 0, 0)
	brain, ok := s.World.Brains.Get(e)
	if !ok || brain.State != stateAlpha {
		t.Fatalf("spawn brain: ok=%v state=%q", ok, brain.State)
	}

	if !tickUntilState(s, e, stateBeta, 200) {
		brain, _ = s.World.Brains.Get(e)
		t.Fatalf("timed out waiting for beta; state=%q timer=%v", brain.State, brain.BB.Timer)
	}
	if !tickUntilState(s, e, stateAlpha, 200) {
		brain, _ = s.World.Brains.Get(e)
		t.Fatalf("timed out waiting for alpha; state=%q timer=%v", brain.State, brain.BB.Timer)
	}
}

func tickUntilState(s *Session, e sim.Entity, want sim.StateID, maxTicks int) bool {
	for i := 0; i < maxTicks; i++ {
		s.Tick(1.0 / 60)
		brain, ok := s.World.Brains.Get(e)
		if ok && brain.State == want {
			return true
		}
	}
	return false
}
