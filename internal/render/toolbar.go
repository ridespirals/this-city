package render

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/config"
	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/sim"
)

// DrawToolbar paints tool buttons and a short hint line (screen space).
func DrawToolbar(ed *editor.Editor) {
	if ed == nil {
		return
	}
	l := config.C.UI.ToolbarLayout()
	panelH := l.Height(int(editor.ToolCount))
	panelW := l.Width()
	rl.DrawRectangle(
		int32(l.X),
		int32(l.Y),
		int32(panelW),
		int32(panelH),
		rl.NewColor(20, 24, 30, 220),
	)

	labelSize := config.C.UI.SizeLabel()
	x := l.X + l.Pad
	y := l.Y + l.Pad
	for t := editor.Tool(0); t < editor.ToolCount; t++ {
		bg := rl.NewColor(50, 58, 70, 255)
		fg := rl.LightGray
		if t == ed.ActiveTool {
			bg = rl.NewColor(70, 120, 160, 255)
			fg = rl.RayWhite
		}
		rl.DrawRectangle(int32(x), int32(y), int32(l.BtnW), int32(l.BtnH), bg)
		label := fmt.Sprintf("%d %s", int(t)+1, editor.ToolName(t))
		textY := y + (l.BtnH-labelSize)/2
		if textY < y {
			textY = y
		}
		TextLabel(int32(x+l.TextPad), int32(textY), label, fg)
		y += l.BtnH + l.Gap
	}

	hintY := int32(l.Y + panelH + l.Pad)
	hint := "Space pause · RMB pan · wheel zoom · Del delete"
	if ed.ActiveTool == editor.ToolPlaceEvent {
		hint = fmt.Sprintf("Event: %s (E cycle) · %s", sim.EventKindName(ed.EventKind), hint)
	}
	if ed.ActiveTool == editor.ToolDrawPath {
		hint = "Click anchors to draw · Del removes path · " + hint
	}
	if ed.ActiveTool == editor.ToolStampPiece {
		name := "(no SVG pieces)"
		if p, ok := ed.CurrentPiece(); ok {
			name = p.Name
		}
		hint = fmt.Sprintf("Stamp: %s (P cycle) · click to place · %s", name, hint)
	}
	TextCaption(int32(l.X), hintY, hint, rl.Gray)
}
