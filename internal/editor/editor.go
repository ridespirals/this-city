// Package editor owns toolbar state and placement tools.
// Mutations go through game command APIs; raylib stays out of this package.
package editor

import (
	"github.com/ridespirals/this-city/internal/game"
	"github.com/ridespirals/this-city/internal/sim"
)

// Tool is the active editor mode.
type Tool int

const (
	ToolSelect Tool = iota
	ToolPlaceCivilian
	ToolPlacePolice
	ToolPlaceEvent
	ToolDrawPath
	ToolEditPath
	ToolCount
)

// ToolName returns a short toolbar label.
func ToolName(t Tool) string {
	switch t {
	case ToolSelect:
		return "Select"
	case ToolPlaceCivilian:
		return "Civilian"
	case ToolPlacePolice:
		return "Police"
	case ToolPlaceEvent:
		return "Event"
	case ToolDrawPath:
		return "Draw path"
	case ToolEditPath:
		return "Edit path"
	default:
		return "?"
	}
}

// Toolbar layout (screen space) — keep in sync with render.DrawToolbar.
const (
	ToolbarX    float32 = 16
	ToolbarY    float32 = 16
	ToolbarBtnW float32 = 110
	ToolbarBtnH float32 = 32
	ToolbarGap  float32 = 6
	ToolbarPad  float32 = 8
)

// FrameInput is a raylib-free snapshot of one frame of editor input.
type FrameInput struct {
	CursorWorld   sim.Vec2
	CursorScreen  sim.Vec2
	LeftPressed   bool
	LeftDown      bool
	DeletePressed bool
	CycleEvent    bool
	ToolHotkey    Tool
	HasToolHotkey bool
}

// Editor holds UI/tool state that is not part of the sim world.
type Editor struct {
	ActiveTool    Tool
	EventKind     sim.EventKind
	Selected      sim.Entity
	SelectedEdge  sim.EdgeID
	SelectedGroup uint32

	DraftAnchors []sim.Vec2
	DraftGroup   uint32

	Dragging  bool
	DragEdge  sim.EdgeID
	DragWhich int // 0=P0/node From, 1=C0, 2=C1, 3=P1/node To
}

// New returns an editor with the select tool active.
func New() *Editor {
	return &Editor{
		ActiveTool:   ToolSelect,
		EventKind:    sim.EventCrime,
		Selected:     sim.NilEntity,
		SelectedEdge: sim.NilEdge,
	}
}

// SetTool changes the active toolbar tool and clears in-progress drags.
func (e *Editor) SetTool(tool Tool) {
	if e == nil || tool < 0 || tool >= ToolCount {
		return
	}
	e.ActiveTool = tool
	e.Dragging = false
	if tool != ToolDrawPath {
		e.clearDraft()
	}
}

func (e *Editor) clearDraft() {
	e.DraftAnchors = nil
	e.DraftGroup = 0
}

// ToolbarHit returns the tool under screen position, or false.
func ToolbarHit(screen sim.Vec2) (Tool, bool) {
	x := ToolbarX + ToolbarPad
	y := ToolbarY + ToolbarPad
	for t := Tool(0); t < ToolCount; t++ {
		if screen.X >= x && screen.X < x+ToolbarBtnW &&
			screen.Y >= y && screen.Y < y+ToolbarBtnH {
			return t, true
		}
		y += ToolbarBtnH + ToolbarGap
	}
	return 0, false
}

// ToolbarHeight is the total screen height of the toolbar panel.
func ToolbarHeight() float32 {
	return ToolbarPad*2 + float32(ToolCount)*ToolbarBtnH + float32(ToolCount-1)*ToolbarGap
}

// Update applies one frame of input against the session via game commands.
func (e *Editor) Update(s *game.Session, in FrameInput) {
	if e == nil || s == nil || s.World == nil {
		return
	}
	w := s.World

	if in.HasToolHotkey {
		e.SetTool(in.ToolHotkey)
	}
	if in.CycleEvent {
		e.EventKind = sim.EventKind((int(e.EventKind) + 1) % sim.EventKindCount)
	}

	if in.LeftPressed {
		if tool, ok := ToolbarHit(in.CursorScreen); ok {
			e.SetTool(tool)
			return
		}
		e.onWorldClick(s, in.CursorWorld)
	}

	if e.Dragging && in.LeftDown {
		e.onDrag(w, in.CursorWorld)
	}
	if e.Dragging && !in.LeftDown {
		e.Dragging = false
	}

	if in.DeletePressed {
		e.onDelete(w)
	}
}

func (e *Editor) onWorldClick(s *game.Session, pos sim.Vec2) {
	w := s.World
	switch e.ActiveTool {
	case ToolSelect:
		e.Selected = game.PickEntity(w, pos, 24)
		e.SelectedEdge = sim.NilEdge
		if !e.Selected.IsNil() {
			e.Dragging = true
		}
	case ToolPlaceCivilian:
		e.Selected = game.SpawnWalkerOnNearestEdge(w, sim.RoleCivilian, pos, 90)
	case ToolPlacePolice:
		e.Selected = game.SpawnWalkerOnNearestEdge(w, sim.RolePolice, pos, 100)
	case ToolPlaceEvent:
		e.Selected = game.SpawnEvent(w, e.EventKind, pos)
	case ToolDrawPath:
		e.DraftAnchors = append(e.DraftAnchors, pos)
		if len(e.DraftAnchors) >= 2 {
			e.DraftGroup = game.SetPathFromAnchors(w, e.DraftGroup, e.DraftAnchors)
			e.SelectedGroup = e.DraftGroup
		}
	case ToolEditPath:
		e.Selected = sim.NilEntity
		e.SelectedEdge = game.PickEdge(w, pos, 20)
		if e.SelectedEdge != sim.NilEdge {
			if edge, ok := w.Network.GetEdge(e.SelectedEdge); ok {
				e.SelectedGroup = edge.Group
			}
			if which, ok := pickHandle(w, e.SelectedEdge, pos, 14); ok {
				e.Dragging = true
				e.DragEdge = e.SelectedEdge
				e.DragWhich = which
			}
		}
	}
}

func (e *Editor) onDrag(w *sim.World, pos sim.Vec2) {
	switch e.ActiveTool {
	case ToolSelect:
		if !e.Selected.IsNil() {
			if _, ok := w.Followers.Get(e.Selected); ok {
				return
			}
			game.MoveEntity(w, e.Selected, pos)
		}
	case ToolEditPath:
		if e.DragEdge == sim.NilEdge {
			return
		}
		edge, ok := w.Network.GetEdge(e.DragEdge)
		if !ok {
			return
		}
		switch e.DragWhich {
		case 0:
			w.Network.MoveNode(edge.From, pos)
		case 3:
			w.Network.MoveNode(edge.To, pos)
		case 1:
			w.Network.SetEdgeCurve(e.DragEdge, pos, edge.Curve.C1)
		case 2:
			w.Network.SetEdgeCurve(e.DragEdge, edge.Curve.C0, pos)
		}
	}
}

func (e *Editor) onDelete(w *sim.World) {
	if e.ActiveTool == ToolEditPath || e.ActiveTool == ToolDrawPath {
		if e.SelectedGroup != 0 {
			game.DeletePathGroup(w, e.SelectedGroup)
			if e.DraftGroup == e.SelectedGroup {
				e.clearDraft()
			}
			e.SelectedGroup = 0
			e.SelectedEdge = sim.NilEdge
			return
		}
		if e.SelectedEdge != sim.NilEdge {
			game.DeleteEdge(w, e.SelectedEdge)
			e.SelectedEdge = sim.NilEdge
			return
		}
	}
	if !e.Selected.IsNil() {
		game.DeleteEntity(w, e.Selected)
		e.Selected = sim.NilEntity
	}
}

func pickHandle(w *sim.World, id sim.EdgeID, pos sim.Vec2, maxDist float32) (which int, ok bool) {
	e, found := w.Network.GetEdge(id)
	if !found {
		return -1, false
	}
	bestD := maxDist * maxDist
	which = -1
	pts := [4]sim.Vec2{e.Curve.P0, e.Curve.C0, e.Curve.C1, e.Curve.P1}
	for j, pt := range pts {
		d := pt.Sub(pos)
		dd := d.Dot(d)
		if dd <= bestD {
			bestD = dd
			which = j
			ok = true
		}
	}
	return which, ok
}
