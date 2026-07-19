// Command this-city is the desktop entrypoint: window, loop, and layer wiring.
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/config"
	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/game"
	"github.com/ridespirals/this-city/internal/render"
	"github.com/ridespirals/this-city/internal/sim"
)

func main() {
	// Optional: tweak before Open, e.g. config.C.UI.Scale = 1.25

	world := sim.NewWorld()
	session := game.NewSession(world)
	session.SpawnDemo()
	ed := editor.New()
	cam := render.NewCamera()
	input := &render.InputTracker{}

	win := render.Open()
	defer win.Close()

	uiFonts := render.LoadFonts()
	defer uiFonts.Unload()

	bg := rl.NewColor(28, 32, 40, 255)

	for !win.ShouldClose() {
		// Poll once per frame (edge-detect), then update, then draw.
		in := input.Poll(cam)
		if input.SpacePressed {
			session.TogglePause()
		}
		if input.EscPressed {
			break
		}

		ed.Update(session, in)

		dt := sim.ClampDT(win.FrameDT(), config.C.Sim.MaxDT)
		session.Tick(dt)

		render.BeginFrame(bg)
		cam.Begin()
		render.DrawPaths(session.World, ed.SelectedEdge)
		if ed.ActiveTool == editor.ToolEditPath || ed.ActiveTool == editor.ToolDrawPath {
			render.DrawPathHandles(session.World, ed.SelectedEdge)
		}
		render.DrawWorld(session.World, ed.Selected)
		render.DrawGhost(ed, in.CursorWorld)
		cam.End()

		render.DrawToolbar(ed)
		render.DrawHUD(render.FrameInfo{
			Paused:  session.Paused,
			SimTime: session.Time,
			Phase:   "Figure 8 · walk + path decisions",
		})
		render.EndFrame()
	}
}
