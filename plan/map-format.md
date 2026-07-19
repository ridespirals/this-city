# Map storage format

## Purpose

Persist and share city street networks (nodes + Bézier edges) independently of entity/agent state. Scene save/load (Phase 8) will compose map data with entities/events.

## Current format (`*.json`)

Located under [`maps/`](../maps/). Schema:

```json
{
  "name": "figure-8",
  "nodes": [{"id": 1, "x": 640, "y": 360}],
  "edges": [{
    "id": 1,
    "from": 1,
    "to": 2,
    "c0": {"x": 560, "y": 220},
    "c1": {"x": 420, "y": 220}
  }]
}
```

| Field | Meaning |
|-------|---------|
| `nodes[].id` | Stable id within the file (remapped on load) |
| `nodes[].x/y` | Junction position |
| `edges[].from/to` | Endpoint node ids |
| `edges[].c0/c1` | Cubic Bézier controls; `P0`/`P1` come from node positions |

## Code

| API | Package |
|-----|---------|
| `MapFile`, `LoadNetworkJSON`, `LoadNetworkFile`, `ApplyMapFile`, `ExportMapFile` | `internal/sim/mapfile.go` |
| Built-in Figure 8 map | `sim.FigureEightMap()` + [`maps/figure-8.json`](../maps/figure-8.json) |
| Demo loader | `game.LoadDemoMap` (file if present, else built-in) |

## Sample map: figure-8

Classic horizontal ∞:

- **Center** — degree-4 crossing
- **Left / right tips** — lobe endpoints
- **Four arcs** — upper/lower left and upper/lower right

Agents use [`PathDecision`](pathfinding.md) at the center junction.

## Future

- Embed maps with `go:embed`
- Version field + migrations
- Full scene files: map + entities + events + camera
- Editor “Save map…” / “Load map…”

See also: [paths.md](paths.md), [pathfinding.md](pathfinding.md).
