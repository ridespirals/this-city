package editor

import (
	"testing"

	"github.com/ridespirals/this-city/internal/game"
	"github.com/ridespirals/this-city/internal/sim"
)

func TestPlaceCivilianViaClick(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	ed := New()
	ed.SetTool(ToolPlaceCivilian)
	ed.Update(s, FrameInput{
		CursorWorld:  sim.Vec2{X: 100, Y: 200},
		CursorScreen: sim.Vec2{X: 400, Y: 400}, // outside toolbar
		LeftPressed:  true,
	})
	if s.World.Roles.Len() != 1 {
		t.Fatalf("roles=%d", s.World.Roles.Len())
	}
	var role sim.Role
	s.World.Roles.ForEach(func(_ sim.Entity, tag sim.RoleTag) { role = tag.Role })
	if role != sim.RoleCivilian {
		t.Fatalf("role=%v", role)
	}
}

func TestToolbarHitSelectsTool(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	ed := New()
	// First button is Select at toolbar pad offset.
	ed.SetTool(ToolPlacePolice)
	ed.Update(s, FrameInput{
		CursorScreen: sim.Vec2{X: ToolbarX + ToolbarPad + 10, Y: ToolbarY + ToolbarPad + 10},
		LeftPressed:  true,
	})
	if ed.ActiveTool != ToolSelect {
		t.Fatalf("tool=%v want Select", ed.ActiveTool)
	}
	if s.World.Len() != 0 {
		t.Fatal("toolbar click should not spawn")
	}
}

func TestDrawPathCreatesSegments(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	ed := New()
	ed.SetTool(ToolDrawPath)
	click := func(x, y float32) {
		ed.Update(s, FrameInput{
			CursorWorld:  sim.Vec2{X: x, Y: y},
			CursorScreen: sim.Vec2{X: 400, Y: 400},
			LeftPressed:  true,
		})
	}
	click(0, 0)
	click(100, 0)
	click(100, 100)
	if s.World.Paths.Len() != 1 {
		t.Fatalf("paths=%d", s.World.Paths.Len())
	}
	p, ok := s.World.Paths.Get(ed.DraftPathID)
	if !ok || len(p.Segments) != 2 {
		t.Fatalf("segments=%v ok=%v", len(p.Segments), ok)
	}
}

func TestDeleteSelectedEntity(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	ed := New()
	ed.SetTool(ToolPlaceEvent)
	ed.Update(s, FrameInput{
		CursorWorld:  sim.Vec2{X: 50, Y: 50},
		CursorScreen: sim.Vec2{X: 400, Y: 400},
		LeftPressed:  true,
	})
	if s.World.Events.Len() != 1 {
		t.Fatal("expected event")
	}
	ed.SetTool(ToolSelect)
	ed.Update(s, FrameInput{DeletePressed: true})
	if s.World.Events.Len() != 0 {
		t.Fatal("expected delete")
	}
}

func TestHotkeyAndCycleEvent(t *testing.T) {
	ed := New()
	ed.Update(nil, FrameInput{}) // nil session no panic
	s := game.NewSession(sim.NewWorld())
	ed.Update(s, FrameInput{HasToolHotkey: true, ToolHotkey: ToolPlaceEvent})
	if ed.ActiveTool != ToolPlaceEvent {
		t.Fatal("hotkey")
	}
	ed.Update(s, FrameInput{CycleEvent: true})
	if ed.EventKind != sim.EventDistress {
		t.Fatalf("kind=%v", ed.EventKind)
	}
}
