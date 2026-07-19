package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/internal/editor"
	"github.com/ridespirals/this-city/internal/sim"
)

// InputTracker edge-detects buttons/keys from held state so short clicks
// are not lost if raylib's one-frame Is*Pressed flags are missed.
type InputTracker struct {
	leftDown  bool
	rightDown bool
	keysDown  [6]bool // 1..6 tool hotkeys
	spaceDown bool
	escDown   bool
	delDown   bool
	eDown     bool

	// App keys sampled in the same Poll as mouse (avoid double PollInputEvents).
	SpacePressed bool
	EscPressed   bool
}

// Poll gathers a FrameInput and updates camera controls.
// Call once per frame, before simulation/UI updates.
func (t *InputTracker) Poll(cam *Camera) editor.FrameInput {
	if t == nil {
		t = &InputTracker{}
	}
	// Pull OS events before sampling (EndDrawing also polls; this keeps
	// pre-draw UI updates on fresh input).
	rl.PollInputEvents()

	mouse := rl.GetMousePosition()
	zoom := float32(1)
	if cam != nil && cam.Inner.Zoom > 0 {
		zoom = cam.Inner.Zoom
	}

	left := rl.IsMouseButtonDown(rl.MouseButtonLeft)
	right := rl.IsMouseButtonDown(rl.MouseButtonRight)

	in := editor.FrameInput{
		CursorScreen: sim.Vec2{X: mouse.X, Y: mouse.Y},
		CursorWorld:  cam.ScreenToWorld(mouse.X, mouse.Y),
		Zoom:         zoom,
		LeftPressed:  left && !t.leftDown,
		LeftReleased: !left && t.leftDown,
		LeftDown:     left,
	}
	t.leftDown = left

	keyCodes := []int32{rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour, rl.KeyFive, rl.KeySix}
	for i, key := range keyCodes {
		down := rl.IsKeyDown(key)
		if down && !t.keysDown[i] {
			in.HasToolHotkey = true
			in.ToolHotkey = editor.Tool(i)
		}
		t.keysDown[i] = down
	}

	del := rl.IsKeyDown(rl.KeyDelete) || rl.IsKeyDown(rl.KeyBackspace)
	in.DeletePressed = del && !t.delDown
	t.delDown = del

	eKey := rl.IsKeyDown(rl.KeyE)
	in.CycleEvent = eKey && !t.eDown
	t.eDown = eKey

	sp := rl.IsKeyDown(rl.KeySpace)
	t.SpacePressed = sp && !t.spaceDown
	t.spaceDown = sp

	es := rl.IsKeyDown(rl.KeyEscape)
	t.EscPressed = es && !t.escDown
	t.escDown = es

	if right {
		d := rl.GetMouseDelta()
		cam.Pan(d.X, d.Y)
	}
	t.rightDown = right

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		factor := float32(1.1)
		if wheel < 0 {
			factor = 1 / factor
		}
		cam.ZoomAt(mouse.X, mouse.Y, factor)
		in.Zoom = cam.Inner.Zoom
		in.CursorWorld = cam.ScreenToWorld(mouse.X, mouse.Y)
	}

	return in
}
