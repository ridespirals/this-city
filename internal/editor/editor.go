// Package editor owns toolbar state and placement tools.
// Mutations go through game/sim command APIs (tools arrive in Phase 5).
package editor

// Tool is the active editor mode.
type Tool int

const (
	ToolSelect Tool = iota
	ToolPlaceCivilian
	ToolPlacePolice
	ToolPlaceEvent
	ToolDrawPath
	ToolEditPath
)

// Editor holds UI/tool state that is not part of the sim world.
type Editor struct {
	ActiveTool Tool
}

// New returns an editor with the select tool active.
func New() *Editor {
	return &Editor{ActiveTool: ToolSelect}
}

// SetTool changes the active toolbar tool.
func (e *Editor) SetTool(tool Tool) {
	if e == nil {
		return
	}
	e.ActiveTool = tool
}
