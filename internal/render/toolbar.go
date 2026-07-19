package render

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/sim"
)

// DrawToolbar paints tool buttons and a short hint line (screen space).
func DrawToolbar(ed *editor.Editor) {
	if ed == nil {
		return
	}
	panelH := editor.ToolbarHeight()
	panelW := editor.ToolbarPad*2 + editor.ToolbarBtnW
	rl.DrawRectangle(
		int32(editor.ToolbarX),
		int32(editor.ToolbarY),
		int32(panelW),
		int32(panelH),
		rl.NewColor(20, 24, 30, 220),
	)

	x := editor.ToolbarX + editor.ToolbarPad
	y := editor.ToolbarY + editor.ToolbarPad
	for t := editor.Tool(0); t < editor.ToolCount; t++ {
		bg := rl.NewColor(50, 58, 70, 255)
		fg := rl.LightGray
		if t == ed.ActiveTool {
			bg = rl.NewColor(70, 120, 160, 255)
			fg = rl.RayWhite
		}
		rl.DrawRectangle(int32(x), int32(y), int32(editor.ToolbarBtnW), int32(editor.ToolbarBtnH), bg)
		label := fmt.Sprintf("%d %s", int(t)+1, editor.ToolName(t))
		rl.DrawText(label, int32(x+8), int32(y+8), 16, fg)
		y += editor.ToolbarBtnH + editor.ToolbarGap
	}

	hintY := int32(editor.ToolbarY + panelH + 8)
	hint := "Space pause · RMB pan · wheel zoom · Del delete"
	if ed.ActiveTool == editor.ToolPlaceEvent {
		hint = fmt.Sprintf("Event: %s (E cycle) · %s", sim.EventKindName(ed.EventKind), hint)
	}
	if ed.ActiveTool == editor.ToolDrawPath {
		hint = "Click anchors to draw · Del removes path · " + hint
	}
	rl.DrawText(hint, int32(editor.ToolbarX), hintY, 14, rl.Gray)
}
