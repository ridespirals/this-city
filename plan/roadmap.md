# Roadmap

Phased delivery from docs to a 3D-ready layered sim. Exit criteria are practical “done enough to build on.”

## Phase 1 — Docs foundation

**Status:** complete (docs in tree; human commits when ready).

- [x] `plan/` tree with per-system docs
- [x] Root `README.md` (human)
- [x] `AGENTS.md` (AI agents)

**Exit:** New contributors (human or AI) can explain layers, ECS, and next phase without guessing.

## Phase 2 — Module + loop

**Status:** complete.

- [x] Go module (`github.com/ridespirals/this-city`), raylib-go v0.60
- [x] Package skeletons: `internal/sim`, `game`, `render`, `editor`
- [x] `cmd/this-city`: window, clear color, frame clock, pause (Space), quit (Esc)
- Suggested tag: `v0.1.0`

**Exit:** `go run ./cmd/this-city` opens a window.

## Phase 3 — ECS + FSM core

- Entity allocator, component stores, basic systems runner.
- Generic FSM engine + unit tests.
- One debug agent that toggles states on a timer (no art required).

**Exit:** Tests prove ECS + FSM; debug agent visible or logged.

## Phase 4 — Bezier paths + follower

- Path data structures, sampling, path follower system.
- Render paths and a moving agent on a hard-coded or loaded path.

**Exit:** Agent follows a curve at roughly constant speed.

## Phase 5 — Editor toolbar

- Tools: select, place civilian/police, place events, draw/edit path.
- Pause toggle; commands mutate world.

**Exit:** Author a tiny scene without recompiling data.

## Phase 6 — Civilian + police FSMs

- Full v1 state sets from [agents.md](agents.md).
- Hang out at benches; patrol/breaks; investigate crimes.

**Exit:** Observable sandbox loop with both roles reacting to placed events.

## Phase 7 — BSP + pathfinding

- BSP radius queries power perception.
- Path graph + A\* between POIs/patrol nodes.

**Exit:** Agents path across a multi-segment network; perception scales past naive O(n²) for demo sizes.

## Phase 8 — Polish / save-load stub

- Serialize world (paths, entities, events) to JSON or similar.
- Load on startup; basic error handling.

**Exit:** Restart app and resume an authored scene.

## Phase 9 — 3D track (later)

- 3D camera/meshes in `render`; sim API remains stable.
- Ground-plane mapping from 2D transforms.

**Exit:** Same game systems visible in a simple 3D view.

## Release suggestions

| Milestone | Tag |
|-----------|-----|
| Window runs | `v0.1.0` |
| Paths + follower | `v0.2.0` |
| Editor + agents sandbox | `v0.3.0` |
| BSP + pathfinding | `v0.4.0` |
| Save/load | `v0.5.0` |

Humans cut tags; agents only suggest.

## Explicitly out of scope (for now)

- Multiplayer, combat depth, dialogue trees, full building interiors, shipping asset pipeline.
