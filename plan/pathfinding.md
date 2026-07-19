# Pathfinding & path decisions

## Purpose

Choose how agents leave junctions and (later) compute routes across the network with Dijkstra/A*.

## PathDecision (implemented)

Component on agents (`World.Decisions`):

| Mode | Behavior |
|------|----------|
| `DecideRandom` (default) | Uniform pick among incident edges; avoid U-turn when alternatives exist |
| `DecideRoute` | Consume `Route []EdgeID` in order (filled by A*/Dijkstra or scripts) |

API: `sim.ChooseNext(network, decision, arrival, rng)`.

If no route is assigned, agents **always** use random intersection choice while walking.

## Graph

Built as `sim.Network`:

- **Nodes:** junctions
- **Edges:** Bézier segments, weight ≈ `Poly.Length`
- Bidirectional travel (follower `Forward` flag)

## Algorithms (next)

- **A\*** / Dijkstra on the node/edge graph → fill `PathDecision.Route` and set `Mode = DecideRoute`.
- Used by future **travel** behavior (go A→B) and optionally **patrol** loops.
- Invalidate routes when the network version changes (editor edits).

## Behaviors (planned interaction)

| Behavior | Decision use |
|----------|----------------|
| `walk` (now) | Random at junctions |
| `travel` | A* route, then walk edges |
| `patrol` | Fixed or cyclic route list |
| `wander` | Random + occasional POI stops |
| `flee` | Prefer edges increasing distance from threat (or reverse A*) |

See [agents.md](agents.md).

## Testing

- Command-key map degree checks; follower visits multiple edges.
- `DecideRoute` overrides random when a route is queued.

See also: [paths.md](paths.md), [map-format.md](map-format.md).
