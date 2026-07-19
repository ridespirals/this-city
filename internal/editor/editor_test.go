package editor

import (
	"testing"

	"github.com/ridespirals/this-city/internal/config"
	"github.com/ridespirals/this-city/internal/game"
	"github.com/ridespirals/this-city/internal/sim"
)

func TestPlaceCivilianViaClick(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	_ = game.LoadDemoMap(s.World)
	ed := New()
	ed.SetTool(ToolPlaceCivilian)
	ed.Update(s, FrameInput{
		CursorWorld:  sim.Vec2{X: 640, Y: 250},
		CursorScreen: sim.Vec2{X: 400, Y: 400},
		LeftPressed:  true,
	})
	if s.World.Roles.Len() != 1 {
		t.Fatalf("roles=%d", s.World.Roles.Len())
	}
	var brain sim.AgentBrain
	s.World.Brains.ForEach(func(_ sim.Entity, b sim.AgentBrain) { brain = b })
	if brain.Machine != game.MachineWalk || brain.State != game.StateWalk {
		t.Fatalf("brain=%+v", brain)
	}
}

func TestToolbarHitSelectsTool(t *testing.T) {
	s := game.NewSession(sim.NewWorld())
	ed := New()
	ed.SetTool(ToolPlacePolice)
	l := config.C.UI.ToolbarLayout()
	ed.Update(s, FrameInput{
		CursorScreen: sim.Vec2{X: l.X + l.Pad + 10, Y: l.Y + l.Pad + 10},
		LeftPressed:  true,
	})
	if ed.ActiveTool != ToolSelect {
		t.Fatalf("tool=%v want Select", ed.ActiveTool)
	}
}

func TestDrawPathCreatesEdges(t *testing.T) {
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
	if s.World.Network.EdgeCount() != 2 {
		t.Fatalf("edges=%d", s.World.Network.EdgeCount())
	}
	if ed.DraftGroup == 0 {
		t.Fatal("expected draft group")
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
	ed.SetTool(ToolSelect)
	ed.Update(s, FrameInput{DeletePressed: true})
	if s.World.Events.Len() != 0 {
		t.Fatal("expected delete")
	}
}

func TestWorldPickRadiusAccountsForZoom(t *testing.T) {
	if WorldPickRadius(36, 1) != 36 {
		t.Fatalf("zoom1=%v", WorldPickRadius(36, 1))
	}
	if WorldPickRadius(36, 0.5) != 72 {
		t.Fatalf("zoom0.5=%v", WorldPickRadius(36, 0.5))
	}
}

func TestToolbarContainsPanel(t *testing.T) {
	l := config.C.UI.ToolbarLayout()
	inside := sim.Vec2{X: l.X + l.Pad + 2, Y: l.Y + l.Pad + 2}
	if !ToolbarContains(inside) {
		t.Fatal("expected inside panel")
	}
	outside := sim.Vec2{X: l.X + l.Width() + 20, Y: l.Y}
	if ToolbarContains(outside) {
		t.Fatal("expected outside panel")
	}
}
