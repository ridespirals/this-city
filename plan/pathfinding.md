# Pathfinding & path decisions

## Purpose

Choose how agents leave junctions and compute routes across the street **Network**, including when the map or endpoints change.

## PathDecision (runtime)

Component on agents (`World.Decisions`) in `internal/sim`:

| Mode | Behavior |
|------|----------|
| `DecideRandom` (default) | Uniform pick among incident edges; avoid U-turn when alternatives exist |
| `DecideRoute` | Consume `Route []EdgeID` in order |

API: `sim.ChooseNext`. Fill routes with `PathDecision.SetRoute(edges)` after a search.

## Graph

`sim.Network`:

- **Nodes:** junctions  
- **Edges:** Bézier segments, weight ≈ `Poly.Length` × optional cost multiplier  
- Travel is **bidirectional**  
- `Version()` bumps on topology / cost / block changes  
- `SetEdgeBlocked` / `SetEdgeCostMul` — temporary closures and congestion  

## Package `internal/sim/pathfind`

Same idea as `internal/sim/noise`: algorithms live in a focused subpackage.

| File | Contents |
|------|----------|
| `pathfind.go` | `Algo`, `Query`, `Result`, `Find`, A*/Dijkstra/BFS/DFS, bidirectional |
| `dstar.go` | `DStarLite` incremental replanner |
| `dynamic.go` | `DynamicRoute` (cache + full replan on change) |
| `graph.go` | Network adjacency / heuristic helpers |
| `heap.go` | Priority queues |

### One-shot search

| Algo | Helper |
|------|--------|
| `AlgoAStar` | `AStar` |
| `AlgoDijkstra` | `Dijkstra` |
| `AlgoBFS` | `BFS` |
| `AlgoDFS` | `DFS` |
| `AlgoBidirectionalBFS` | `BidirectionalBFS` |
| `AlgoBidirectionalAStar` | `BidirectionalAStar` |
| `AlgoBidirectionalDijkstra` | via `Find` |

```go
import "github.com/ridespirals/this-city/internal/sim/pathfind"

r := pathfind.AStar(world.Network, from, to)
if r.Found {
    dec.SetRoute(r.Edges)
}
```

### Bidirectional / “meet in the middle”

**Bidirectional search**: grow from start and goal until frontiers meet, then stitch.

### Dynamic pathing

| Tool | When |
|------|------|
| `DynamicRoute` | Infrequent changes; recomputes when endpoints or `Network.Version` change |
| `DStarLite` | Frequent blocks/cost changes and/or moving start toward a stable goal |

**D\*** is the classic name; **D\* Lite** is what we implement.

```go
d := pathfind.NewDStarLite(net, start, goal)
r := d.Replan()
net.SetEdgeBlocked(edgeID, true)
r = d.Replan()
d.SetStart(currentNode)
r = d.Replan()
```

## Behaviors (planned)

| Behavior | Decision use |
|----------|----------------|
| `walk` (now) | Random at junctions |
| `travel` | `pathfind` → `DecideRoute` |
| `patrol` / `wander` / `flee` | See [agents.md](agents.md) |

## Testing

`go test ./internal/sim/pathfind/...` — line/branch graphs, blocks, cost muls, figure-8 tip-to-tip.

See also: [paths.md](paths.md), [map-format.md](map-format.md), [noise.md](noise.md).
