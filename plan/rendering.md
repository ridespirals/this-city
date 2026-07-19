# Rendering

## Purpose

Present the simulation and editor chrome via raylib, without owning game rules. Isolate so a future 3D renderer can replace or extend this layer.

## Stack

- Go + [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go) v0.60
- Package: `internal/render` (+ UI bits in `internal/editor`)
- Phase 2: `Window` lifecycle, clear/draw frame helpers, placeholder splash text

## 2D (current target)

- Clear background, draw path strokes, agents as simple sprites/shapes, events as icons/markers.
- Camera2D: pan/zoom.
- Toolbar drawn in screen space (not world space).

## Isolation rules

- **Read** sim/game state for drawing.
- **Write** only by emitting commands (or calling documented command functions).
- No FSM transitions inside draw code.
- No Bezier subdivision ownership—use sim-provided samples or shared pure functions in `sim`.

## Frame loop (typical)

```
process input → commands
if !paused: game.Tick(dt)
render.Draw(world, editorState)
```

## 3D (later)

- Add perspective camera and meshes; map agents to ground plane (X/Z) from sim `Transform2D` or promote to `Transform3D`.
- Keep `game` FSMs and path network in plane space initially.
- Swap draw implementations behind a small `Renderer` interface if both 2D and 3D coexist.

## Assets

- UI font: Space Mono (Regular + Bold), embedded from `assets/fonts/`, drawn via `render.Text` / `TextBold`.
- v1 world: primitives and colored shapes. Texture/atlas pipeline deferred.

## Testing

- Prefer not to unit-test pixels.
- Keep render thin; logic bugs should be catchable in `sim`/`game` tests.
- Optional screenshot/manual checklist in roadmap milestones.

See also: [architecture.md](architecture.md), [editor.md](editor.md).
