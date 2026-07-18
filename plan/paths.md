# Paths (Bezier Streets)

## Purpose

Author and run along streets and walkways as cubic Bezier curves, forming a navigable network for agents.

## Implementation (Phase 4)

| Piece | Location |
|-------|----------|
| `Vec2`, `CubicBezier`, `Polyline` | `internal/sim/vec2.go`, `bezier.go` |
| `Path`, `PathSet`, `PathFollower` | `internal/sim/path.go` |
| Advance / tick | `internal/sim/path_follow.go` |
| Demo path + spawn | `internal/game/demo.go` |
| Draw | `render.DrawPaths` uses `Path.Poly` only |

## Authoring model

- **Segment:** cubic Bezier (`p0`, `c0`, `c1`, `p1`) in world 2D.
- **Path / chain:** ordered segments; endpoints may share junctions.
- **Junction (node):** point where segments meet; used by the path graph ([pathfinding.md](pathfinding.md)).
- Editor tools create/edit control points; see [editor.md](editor.md).

## Runtime representation

- Store exact Bezier control points for editing and high-quality drawing.
- Build a **polyline approximation** (adaptive subdivision or fixed samples) for:
  - Arc-length / distance-along-path queries.
  - Constant-speed following.
  - Debug draw and hit-testing.

Recompute samples when a segment is edited.

## Path follower

Component fields (conceptual):

- `PathID` / segment index or edge id in the path graph.
- `Distance` along current edge (meters/units).
- `Speed`, `Forward` (bool or sign).
- Optional `GoalDistance` or node id for stop conditions.

Each tick:

1. Advance distance by `speed * dt`.
2. If past end of edge, either stop, reverse, or hand off to pathfinding for the next edge.
3. Sample position (and tangent → facing) from the polyline/Bezier.

## Constant speed

Prefer distance along the sampled polyline (chord-length cumulative). Good enough for v1; refine with true arc-length tables if wobble appears on sharp curves.

## Drawing

- `render` strokes Bezier or polyline; control handles only in editor mode.
- Sim owns geometry; render never invents path math.

## Interaction with pathfinding

- Bezier network → undirected (or directed) **graph** of nodes + edges.
- Pathfinding returns a sequence of edges; the follower consumes that sequence.

See also: [pathfinding.md](pathfinding.md), [spatial.md](spatial.md), [editor.md](editor.md).
