# Pathfinding

## Purpose

Find routes across the Bezier street network so agents travel between locations (home, park bench, patrol waypoints) without steering freely through buildings (v1 assumes travel stays on the network).

## Graph construction

Derived from path geometry ([paths.md](paths.md)):

- **Nodes:** junctions + notable stops (bench anchors, spawn points, patrol posts).
- **Edges:** Bezier segments (or chains) with weight = approximate length.
- Rebuild graph when the editor commits geometry changes.

## Algorithm (v1)

- **A\*** on the path graph with Euclidean (or octile) heuristic between node positions.
- Output: ordered list of node ids / edge ids.
- Agent `PathFollower` + brain consume the route; replan if blocked or destination changes.

## Off-network movement (later)

Short steers to sit on a bench or approach an event may leave the polyline briefly; rejoin nearest point on network afterward. Not required for first follower demo.

## Integration with BSP

- BSP does **not** replace A\*; it accelerates “who is near me” and picking.
- Optional: use spatial index to find nearest path node to an agent when starting a journey.

## Performance

- City graphs stay small initially; naive A\* is fine.
- Cache last route per agent; invalidate on graph version bump.

## Testing

- Hand-built graphs: shortest path correctness, unreachable goals, zero-length edges.
- Fuzz weights non-negative.

See also: [paths.md](paths.md), [spatial.md](spatial.md), [agents.md](agents.md).
