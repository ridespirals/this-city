# This City — Plan Index

Planning docs for the city simulation. Code follows these contracts; when scope or design changes, update the relevant files here plus the root [README.md](../README.md) and [AGENTS.md](../AGENTS.md).

## Status

**Current phase:** 5+ — Editor + ⌘ network map, Walk FSM, PathDecision. Next: Phase 6 — richer civilian/police behaviors. See [roadmap.md](roadmap.md).

## Locked decisions

| Decision | Choice |
|----------|--------|
| Language | Go |
| Graphics | raylib via [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go) |
| Architecture | ECS (not OOP inheritance) |
| Layers | `sim` (pure) → `game` → `render` / `editor` |
| Git | Humans own all git ops; agents only suggest |

## Document map

| Doc | Topic |
|-----|--------|
| [architecture.md](architecture.md) | Packages, dependency rules, layering |
| [ecs.md](ecs.md) | Entity / component / system design |
| [state-machines.md](state-machines.md) | Reusable FSM engine for agents |
| [paths.md](paths.md) | Bezier streets and path following |
| [spatial.md](spatial.md) | BSP trees and spatial queries |
| [pathfinding.md](pathfinding.md) | PathDecision, graph, A*/Dijkstra plans |
| [map-format.md](map-format.md) | JSON map storage (`maps/`) |
| [agents.md](agents.md) | Walk default + civilian/police plans |
| [events.md](events.md) | Placeable events / objects and responses |
| [editor.md](editor.md) | 2D overhead toolbar and placement |
| [rendering.md](rendering.md) | 2D now, 3D later; render isolation |
| [testing.md](testing.md) | How each layer is tested |
| [roadmap.md](roadmap.md) | Phased development and milestones |

## Development order (summary)

1. Docs foundation — done
2. Module + game loop — done
3. ECS + FSM core — done
4. Bezier paths + follower — done
5. Editor toolbar — **done**
6. Civilian + police FSMs ← **next**
7. BSP + pathfinding
8. Save/load stub
9. 3D track (later)

Details and exit criteria: [roadmap.md](roadmap.md).

## Package layout

```
cmd/this-city/main.go
internal/sim/       # pure: world, ECS, time, spatial, math — no raylib
internal/game/      # systems: agents, events, schedules — no raylib
internal/render/    # raylib draw/input adapters only
internal/editor/    # toolbar, placement (uses game + render)
```
