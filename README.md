# This City

A city simulation sandbox: agents living on Bezier streets, driven by state machines, with an overhead editor that grows toward a simple 3D game.

Civilians travel between places, hang out in the park, and react to things they notice. Police work shifts, patrol, take breaks, and investigate crimes or people who need help. Under the hood: an ECS, a reusable FSM engine, path following, pathfinding, and BSP-accelerated spatial queries—kept separate from rendering so the sim stays testable.

## Status

**Editor + ⌘ network map.** Sample streets are two stacked figure-8s (Mac ⌘). Agents use a default **Walk** brain and **PathDecision** (random at intersections; ready for A* routes). Maps load from [`maps/command-key.json`](maps/command-key.json). Plans in [`plan/`](plan/README.md).

Next up: Phase 6 — richer civilian/police behaviors (travel, wander, flee, patrol). See the [roadmap](plan/roadmap.md).

## Stack

| Piece | Choice |
|-------|--------|
| Language | Go |
| Graphics | [raylib](https://www.raylib.com/) via [raylib-go](https://github.com/gen2brain/raylib-go) |
| Structure | Entity Component System (not OOP hierarchies) |
| Layers | Pure **sim** → **game** logic → **render** / **editor** |

## Vision

| Now (near-term) | Later |
|-----------------|--------|
| Overhead 2D view | Simple 3D presentation |
| Toolbar: place agents, events, draw streets | Richer content and scenarios |
| Civilians + police FSMs | More roles and deeper investigations |
| Bezier paths, followers, A*, BSP | LOS, save/load, polish |

## Architecture (short)

```
cmd/this-city          # window + main loop
internal/editor        # toolbar, placement, path authoring
internal/render        # raylib only
internal/game          # agent/event systems (no raylib)
internal/sim           # ECS, FSM, paths, BSP, math (no raylib)
```

`sim` never imports raylib or UI code. Details: [plan/architecture.md](plan/architecture.md).

## Planning docs

Start at **[plan/README.md](plan/README.md)** — index of all system plans:

- [Architecture](plan/architecture.md) · [ECS](plan/ecs.md) · [State machines](plan/state-machines.md)
- [Paths](plan/paths.md) · [Spatial / BSP](plan/spatial.md) · [Pathfinding](plan/pathfinding.md)
- [Agents](plan/agents.md) · [Events](plan/events.md) · [Editor](plan/editor.md)
- [Rendering](plan/rendering.md) · [Testing](plan/testing.md) · [Roadmap](plan/roadmap.md)

AI contributors: see [AGENTS.md](AGENTS.md).

## Requirements

- Go 1.26+ (module targets the toolchain in `go.mod`)
- A desktop environment for the raylib window (`go test ./...` does not open a window)

## Running

```bash
go run ./cmd/this-city
```

- **1–6** / toolbar — select, civilian, police, event, draw path, edit path  
- **E** — cycle event kind (crime / distress / attraction / bench)  
- **Del** / **Backspace** — delete selected entity or path  
- **RMB drag** — pan · **wheel** — zoom  
- **Space** — pause / resume sim clock  
- **Esc** or close the window — quit  

```bash
go test ./...
```

## CI & releases

| Workflow | Trigger | What it does |
|----------|---------|----------------|
| [CI](.github/workflows/ci.yml) | Pushes to `main`, all pull requests | `go test ./...` and build `cmd/this-city` on Ubuntu |
| [Release](.github/workflows/release.yml) | Tags matching `v*` (e.g. `v0.1.0`) | Test, build Linux/macOS/Windows binaries, publish a GitHub Release |

Tag a release (after pushing the commit you want tagged):

```bash
git tag v0.1.0
git push origin v0.1.0
```

## License

[MIT](LICENSE) © 2026 John Varga.
