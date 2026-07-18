# Spatial Indexing (BSP)

## Purpose

Accelerate:

- Nearby-entity queries for perception.
- Editor click-picking.
- Optional locality for path segments and events.

## Approach (v1)

Axis-aligned **binary space partitioning** over the 2D map:

- Internal nodes: split on X or Y (alternate or by longest axis / object count).
- Leaves: small lists of entity ids (and optionally static segment ids).
- Rebuild after batch edits; support incremental insert/remove when agents move if rebuild cost shows up in profiles.

Start with **periodic rebuild** (e.g. each tick or every N ticks for dynamic entities) plus **full rebuild** when path geometry changes. Optimize later.

## Queries

| Query | Use |
|-------|-----|
| Point | Click pick, snap |
| Radius / circle | Perception, event notice |
| AABB | Editor marquee (later) |

Return candidates; caller applies exact distance / filters (role, event kind, LOS stub).

## What goes in the tree

- Dynamic: agents, active events (by `Transform2D`).
- Static (separate or tagged): path segment AABBs, furniture anchors (benches).

Keep static and dynamic separation if it simplifies rebuilds.

## LOS (stub)

v1 perception = radius query only. Line-of-sight (occluders) is a later filter on BSP candidates.

## Testing

- Insert known points; assert radius query membership.
- Degenerate cases: empty world, all points coincident, very large radius.

## Alternatives considered

- Uniform grid / spatial hash: simpler; revisit if BSP complexity outweighs benefit.
- Quadtree: acceptable substitute if implementation is clearer; docs may say “BSP” as the family of hierarchical spatial indexes.

See also: [pathfinding.md](pathfinding.md), [agents.md](agents.md), [events.md](events.md).
