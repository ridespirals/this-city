# This City

A city simulation sandbox: agents living on Bezier streets, driven by state machines, with an overhead editor that grows toward a simple 3D game.

Civilians travel between places, hang out in the park, and react to things they notice. Police work shifts, patrol, take breaks, and investigate crimes or people who need help. Under the hood: an ECS, a reusable FSM engine, path following, pathfinding, and BSP-accelerated spatial queries—kept separate from rendering so the sim stays testable.

## Status

**Phase 1 — Docs foundation.** Architecture and feature plans are in [`plan/`](plan/README.md). Application code is not started yet.

Next up: Go module, raylib window, and package skeletons (`sim` / `game` / `render` / `editor`). See the [roadmap](plan/roadmap.md).

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

## Running

Not applicable until Phase 2. Expected entrypoint:

```bash
go run ./cmd/this-city
```

## License

[MIT](LICENSE) © 2026 John Varga.
