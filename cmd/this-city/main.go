// Command this-city is the desktop entrypoint: window, loop, and layer wiring.
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/game"
	"github.com/ridespirals/this-city/internal/render"
	"github.com/ridespirals/this-city/internal/sim"
)

func main() {
	world := sim.NewWorld()
	session := game.NewSession(world)
	session.SpawnDemo()
	ed := editor.New()
	cam := render.NewCamera()

	win := render.Open(render.DefaultConfig())
	defer win.Close()

	bg := rl.NewColor(28, 32, 40, 255)

	for !win.ShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			session.TogglePause()
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			break
		}

		in := render.CollectEditorInput(cam, ed)
		ed.Update(session, in)

		dt := win.FrameDT()
		dt = sim.ClampDT(dt, render.MaxDT)
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
			Phase:   "⌘ map · walk + path decisions",
		})
		render.EndFrame()
	}
}
