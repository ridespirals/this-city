// Package render adapts raylib for drawing and raw input.
// It must not run agent AI or own path math beyond display.
package render

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/config"
	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/sim"
)

// Window wraps raylib window lifecycle.
type Window struct {
	width  int32
	height int32
}

// Open creates the raylib window from config.C.Window.
func Open() *Window {
	cfg := config.C.Window
	if cfg.Width <= 0 {
		cfg.Width = 1280
	}
	if cfg.Height <= 0 {
		cfg.Height = 720
	}
	if cfg.Title == "" {
		cfg.Title = "This City"
	}
	if cfg.TargetFPS <= 0 {
		cfg.TargetFPS = 60
	}
	rl.InitWindow(cfg.Width, cfg.Height, cfg.Title)
	rl.SetTargetFPS(cfg.TargetFPS)
	return &Window{width: cfg.Width, height: cfg.Height}
}

// Close shuts down the raylib window.
func (w *Window) Close() {
	if w == nil {
		return
	}
	rl.CloseWindow()
}

// ShouldClose reports whether the user requested quit.
func (w *Window) ShouldClose() bool {
	return rl.WindowShouldClose()
}

// FrameDT returns the last frame time in seconds, clamped for sim stability.
func (w *Window) FrameDT() float32 {
	dt := rl.GetFrameTime()
	max := config.C.Sim.MaxDT
	if max <= 0 {
		max = 0.1
	}
	if dt > max {
		return max
	}
	if dt < 0 {
		return 0
	}
	return dt
}

// BeginFrame starts a draw frame with a clear color.
func BeginFrame(clear rl.Color) {
	rl.BeginDrawing()
	rl.ClearBackground(clear)
}

// EndFrame finishes a draw frame.
func EndFrame() {
	rl.EndDrawing()
}

// FrameInfo is HUD data for the current frame.
type FrameInfo struct {
	Paused  bool
	SimTime float64
	Phase   string
}

// DrawHUD paints status text in screen space (top-right-ish).
func DrawHUD(info FrameInfo) {
	phase := info.Phase
	if phase == "" {
		phase = "dev"
	}
	width := config.C.Window.Width
	if width <= 0 {
		width = 1280
	}
	x := width - 360
	TextTitle(x, 40, "This City", rl.RayWhite)
	TextBody(x, 80, phase, rl.LightGray)
	status := "running"
	if info.Paused {
		status = "paused (Space)"
	}
	TextBody(x, 110, status, rl.LightGray)
	TextBody(x, 140, fmt.Sprintf("sim time: %.1fs", info.SimTime), rl.LightGray)
}

// DrawPaths strokes network edge polylines from sim geometry (read-only).
func DrawPaths(w *sim.World, selected sim.EdgeID) {
	if w == nil || w.Network == nil {
		return
	}
	w.Network.ForEachEdge(func(e *sim.Edge) {
		stroke := rl.NewColor(90, 110, 140, 255)
		width := float32(4)
		if e.ID == selected {
			stroke = rl.NewColor(140, 190, 230, 255)
			width = 6
		}
		pts := e.Poly.Points
		for i := 1; i < len(pts); i++ {
			a, b := pts[i-1], pts[i]
			rl.DrawLineEx(rl.NewVector2(a.X, a.Y), rl.NewVector2(b.X, b.Y), width, stroke)
		}
	})
}

// DrawPathHandles draws Bézier control points for the selected edge.
func DrawPathHandles(w *sim.World, edgeID sim.EdgeID) {
	if w == nil || edgeID == sim.NilEdge {
		return
	}
	e, ok := w.Network.GetEdge(edgeID)
	if !ok {
		return
	}
	seg := e.Curve
	rl.DrawLineEx(rl.NewVector2(seg.P0.X, seg.P0.Y), rl.NewVector2(seg.C0.X, seg.C0.Y), 1, rl.DarkGray)
	rl.DrawLineEx(rl.NewVector2(seg.P1.X, seg.P1.Y), rl.NewVector2(seg.C1.X, seg.C1.Y), 1, rl.DarkGray)
	rl.DrawCircle(int32(seg.P0.X), int32(seg.P0.Y), 6, rl.RayWhite)
	rl.DrawCircle(int32(seg.P1.X), int32(seg.P1.Y), 6, rl.RayWhite)
	rl.DrawCircle(int32(seg.C0.X), int32(seg.C0.Y), 5, rl.NewColor(220, 180, 80, 255))
	rl.DrawCircle(int32(seg.C1.X), int32(seg.C1.Y), 5, rl.NewColor(220, 180, 80, 255))
}

// DrawWorld draws agents and events from sim state (read-only).
func DrawWorld(w *sim.World, selected sim.Entity) {
	if w == nil {
		return
	}
	w.Transforms.ForEach(func(e sim.Entity, xf sim.Transform2D) {
		if ev, ok := w.Events.Get(e); ok {
			drawEvent(xf, ev, e == selected)
			return
		}
		color := rl.NewColor(120, 180, 220, 255)
		label := ""
		if role, ok := w.Roles.Get(e); ok {
			switch role.Role {
			case sim.RoleCivilian:
				color = rl.NewColor(100, 180, 220, 255)
				label = "civ"
			case sim.RolePolice:
				color = rl.NewColor(80, 120, 220, 255)
				label = "cop"
			case sim.RoleDebug:
				color = rl.NewColor(80, 200, 140, 255)
			}
		}
		if brain, ok := w.Brains.Get(e); ok {
			label = string(brain.State)
		}
		radius := float32(14)
		if xf.Scale > 0 {
			radius *= xf.Scale
		}
		rl.DrawCircle(int32(xf.X), int32(xf.Y), radius, color)
		if e == selected {
			rl.DrawCircleLines(int32(xf.X), int32(xf.Y), radius+4, rl.RayWhite)
		}
		dx := float32(math.Cos(float64(xf.Rotation))) * (radius + 6)
		dy := float32(math.Sin(float64(xf.Rotation))) * (radius + 6)
		rl.DrawLineEx(
			rl.NewVector2(xf.X, xf.Y),
			rl.NewVector2(xf.X+dx, xf.Y+dy),
			2,
			rl.RayWhite,
		)
		if label != "" {
			TextLabel(int32(xf.X)-12, int32(xf.Y)-34, label, rl.RayWhite)
		}
	})
}

func drawEvent(xf sim.Transform2D, ev sim.EventSource, selected bool) {
	color := rl.NewColor(200, 80, 80, 255)
	switch ev.Kind {
	case sim.EventDistress:
		color = rl.NewColor(220, 160, 60, 255)
	case sim.EventAttraction:
		color = rl.NewColor(200, 100, 200, 255)
	case sim.EventBench:
		color = rl.NewColor(140, 110, 70, 255)
	}
	rl.DrawRectangle(int32(xf.X)-10, int32(xf.Y)-10, 20, 20, color)
	if selected {
		rl.DrawRectangleLines(int32(xf.X)-14, int32(xf.Y)-14, 28, 28, rl.RayWhite)
	}
	TextCaption(int32(xf.X)-20, int32(xf.Y)-28, sim.EventKindName(ev.Kind), rl.LightGray)
}

// DrawGhost draws a placement preview at the cursor.
func DrawGhost(ed *editor.Editor, worldPos sim.Vec2) {
	if ed == nil {
		return
	}
	switch ed.ActiveTool {
	case editor.ToolPlaceCivilian:
		rl.DrawCircleLines(int32(worldPos.X), int32(worldPos.Y), 14, rl.NewColor(100, 180, 220, 180))
	case editor.ToolPlacePolice:
		rl.DrawCircleLines(int32(worldPos.X), int32(worldPos.Y), 14, rl.NewColor(80, 120, 220, 180))
	case editor.ToolPlaceEvent:
		rl.DrawRectangleLines(int32(worldPos.X)-10, int32(worldPos.Y)-10, 20, 20, rl.NewColor(200, 200, 200, 180))
	case editor.ToolDrawPath:
		rl.DrawCircle(int32(worldPos.X), int32(worldPos.Y), 4, rl.NewColor(200, 200, 200, 180))
		if n := len(ed.DraftAnchors); n > 0 {
			last := ed.DraftAnchors[n-1]
			rl.DrawLineEx(
				rl.NewVector2(last.X, last.Y),
				rl.NewVector2(worldPos.X, worldPos.Y),
				2,
				rl.NewColor(180, 180, 180, 120),
			)
		}
	}
}

// CollectEditorInput reads raylib into an editor.FrameInput and updates the camera.
func CollectEditorInput(cam *Camera, ed *editor.Editor) editor.FrameInput {
	mouse := rl.GetMousePosition()
	in := editor.FrameInput{
		CursorScreen:  sim.Vec2{X: mouse.X, Y: mouse.Y},
		CursorWorld:   cam.ScreenToWorld(mouse.X, mouse.Y),
		LeftPressed:   rl.IsMouseButtonPressed(rl.MouseButtonLeft),
		LeftDown:      rl.IsMouseButtonDown(rl.MouseButtonLeft),
		DeletePressed: rl.IsKeyPressed(rl.KeyDelete) || rl.IsKeyPressed(rl.KeyBackspace),
		CycleEvent:    rl.IsKeyPressed(rl.KeyE),
	}

	for i, key := range []int32{rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour, rl.KeyFive, rl.KeySix} {
		if rl.IsKeyPressed(key) {
			in.HasToolHotkey = true
			in.ToolHotkey = editor.Tool(i)
			break
		}
	}

	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		d := rl.GetMouseDelta()
		cam.Pan(d.X, d.Y)
	}
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		factor := float32(1.1)
		if wheel < 0 {
			factor = 1 / factor
		}
		cam.ZoomAt(mouse.X, mouse.Y, factor)
		in.CursorWorld = cam.ScreenToWorld(mouse.X, mouse.Y)
	}

	_ = ed
	return in
}
