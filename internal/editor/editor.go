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
	LeftPressed   bool // edge: just clicked
	LeftDown      bool
	DeletePressed bool
	CycleEvent    bool
	ToolHotkey    Tool // ToolCount means none
	HasToolHotkey bool
}

// Editor holds UI/tool state that is not part of the sim world.
type Editor struct {
	ActiveTool   Tool
	EventKind    sim.EventKind
	Selected     sim.Entity
	SelectedPath sim.PathID

	// Path drawing draft (anchor polyline).
	DraftAnchors []sim.Vec2
	DraftPathID  sim.PathID

	// Dragging for select / edit path.
	Dragging   bool
	DragPathID sim.PathID
	DragSeg    int
	DragWhich  int // 0=P0,1=C0,2=C1,3=P1
}

// New returns an editor with the select tool active.
func New() *Editor {
	return &Editor{
		ActiveTool:   ToolSelect,
		EventKind:    sim.EventCrime,
		Selected:     sim.NilEntity,
		SelectedPath: sim.NilPath,
		DraftPathID:  sim.NilPath,
		DragSeg:      -1,
	}
}

// SetTool changes the active toolbar tool and clears in-progress drags.
func (e *Editor) SetTool(tool Tool) {
	if e == nil || tool < 0 || tool >= ToolCount {
		return
	}
	e.ActiveTool = tool
	e.Dragging = false
	e.DragSeg = -1
	if tool != ToolDrawPath {
		// Keep draft if switching away briefly? Clear for predictability.
		e.clearDraft()
	}
}

func (e *Editor) clearDraft() {
	e.DraftAnchors = nil
	e.DraftPathID = sim.NilPath
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
		e.DragSeg = -1
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
		e.SelectedPath = sim.NilPath
		if !e.Selected.IsNil() {
			e.Dragging = true
		}
	case ToolPlaceCivilian:
		e.Selected = game.SpawnAgent(w, sim.RoleCivilian, pos)
	case ToolPlacePolice:
		e.Selected = game.SpawnAgent(w, sim.RolePolice, pos)
	case ToolPlaceEvent:
		e.Selected = game.SpawnEvent(w, e.EventKind, pos)
	case ToolDrawPath:
		e.DraftAnchors = append(e.DraftAnchors, pos)
		if len(e.DraftAnchors) >= 2 {
			e.DraftPathID = game.SetPathFromAnchors(w, e.DraftPathID, e.DraftAnchors)
			e.SelectedPath = e.DraftPathID
		}
	case ToolEditPath:
		e.Selected = sim.NilEntity
		// Prefer grabbing a control handle on the selected path; else pick a path.
		if e.SelectedPath != sim.NilPath {
			if seg, which, ok := pickHandle(w, e.SelectedPath, pos, 14); ok {
				e.Dragging = true
				e.DragPathID = e.SelectedPath
				e.DragSeg = seg
				e.DragWhich = which
				return
			}
		}
		e.SelectedPath = game.PickPath(w, pos, 20)
		if e.SelectedPath != sim.NilPath {
			if seg, which, ok := pickHandle(w, e.SelectedPath, pos, 14); ok {
				e.Dragging = true
				e.DragPathID = e.SelectedPath
				e.DragSeg = seg
				e.DragWhich = which
			}
		}
	}
}

func (e *Editor) onDrag(w *sim.World, pos sim.Vec2) {
	switch e.ActiveTool {
	case ToolSelect:
		if !e.Selected.IsNil() {
			// Don't fight path followers.
			if _, ok := w.Followers.Get(e.Selected); ok {
				return
			}
			game.MoveEntity(w, e.Selected, pos)
		}
	case ToolEditPath:
		if e.DragPathID == sim.NilPath || e.DragSeg < 0 {
			return
		}
		p, ok := w.Paths.Get(e.DragPathID)
		if !ok || e.DragSeg >= len(p.Segments) {
			return
		}
		segs := append([]sim.CubicBezier(nil), p.Segments...)
		seg := segs[e.DragSeg]
		switch e.DragWhich {
		case 0:
			seg.P0 = pos
		case 1:
			seg.C0 = pos
		case 2:
			seg.C1 = pos
		case 3:
			seg.P1 = pos
		}
		segs[e.DragSeg] = seg
		// Keep chain continuity: moving P1 updates next P0 and vice versa.
		if e.DragWhich == 3 && e.DragSeg+1 < len(segs) {
			n := segs[e.DragSeg+1]
			n.P0 = pos
			segs[e.DragSeg+1] = n
		}
		if e.DragWhich == 0 && e.DragSeg > 0 {
			prev := segs[e.DragSeg-1]
			prev.P1 = pos
			segs[e.DragSeg-1] = prev
		}
		w.Paths.SetSegments(e.DragPathID, segs)
	}
}

func (e *Editor) onDelete(w *sim.World) {
	if e.ActiveTool == ToolEditPath || e.ActiveTool == ToolDrawPath {
		if e.SelectedPath != sim.NilPath {
			game.DeletePath(w, e.SelectedPath)
			if e.DraftPathID == e.SelectedPath {
				e.clearDraft()
			}
			e.SelectedPath = sim.NilPath
			return
		}
	}
	if !e.Selected.IsNil() {
		game.DeleteEntity(w, e.Selected)
		e.Selected = sim.NilEntity
	}
}

func pickHandle(w *sim.World, id sim.PathID, pos sim.Vec2, maxDist float32) (seg, which int, ok bool) {
	p, found := w.Paths.Get(id)
	if !found {
		return -1, -1, false
	}
	bestD := maxDist * maxDist
	seg, which = -1, -1
	for i, s := range p.Segments {
		pts := [4]sim.Vec2{s.P0, s.C0, s.C1, s.P1}
		for j, pt := range pts {
			dx := pt.X - pos.X
			dy := pt.Y - pos.Y
			d := dx*dx + dy*dy
			if d <= bestD {
				bestD = d
				seg = i
				which = j
				ok = true
			}
		}
	}
	return seg, which, ok
}
