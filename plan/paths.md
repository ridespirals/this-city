# Paths (Bezier Streets)

## Purpose

Author and run along streets and walkways as cubic Bezier curves, forming a navigable **network** for agents.

## Implementation

| Piece | Location |
|-------|----------|
| `Network`, `Node`, `Edge` | `internal/sim/network.go` |
| `PathFollower`, junction advance | `internal/sim/path_follow.go` |
| `PathDecision` | `internal/sim/path_decision.go` |
| Map JSON | `internal/sim/mapfile.go`, [`maps/`](../maps/) |
| Demo Figure 8 map | `sim.FigureEightMap`, `game.LoadDemoMap`, `maps/figure-8.json` |
| Draw | `render.DrawPaths` uses `Edge.Poly` only |

## Authoring model

- **Node:** junction point (intersection).
- **Edge:** cubic Bézier between two nodes (`P0`/`P1` snapped to node positions).
- **Group:** editor chain id for linear draw-path tools.
- Map files: [map-format.md](map-format.md).

## Runtime

Each edge holds a **polyline approximation** (chord-length) for constant-speed following.

## Path follower

- `Edge`, `Distance` (0 at From → Length at To), `Speed`, `Forward`.
- On reaching a node, **`PathDecision`** chooses the next edge ([pathfinding.md](pathfinding.md)).
- Dead ends: U-turn on the same edge.

## Drawing

- Stroke `Edge.Poly`; control handles only in edit tools.
- Sim owns geometry; render never re-samples Béziers.

See also: [pathfinding.md](pathfinding.md), [map-format.md](map-format.md), [editor.md](editor.md).
