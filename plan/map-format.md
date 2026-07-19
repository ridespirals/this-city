# Map storage format

## Purpose

Persist and share city street networks (nodes + Bézier edges) independently of entity/agent state. Scene save/load (Phase 8) will compose map data with entities/events.

## Current format (`*.json`)

Located under [`maps/`](../maps/). Schema:

```json
{
  "name": "command-key",
  "nodes": [{"id": 1, "x": 640, "y": 250}],
  "edges": [{
    "id": 1,
    "from": 1,
    "to": 5,
    "c0": {"x": 520, "y": 180},
    "c1": {"x": 420, "y": 180}
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
| Built-in ⌘ map | `sim.CommandKeyMap()` + [`maps/command-key.json`](../maps/command-key.json) |
| Demo loader | `game.LoadDemoMap` (file if present, else built-in) |

## Sample map: command-key (⌘)

Two horizontal figure-8s stacked so they share left/right junctions:

- **Mt / Mb** — crossings of the top and bottom 8s (degree 4)
- **JL / JR** — stack intersections (degree 4)
- **Lt, Rt, Lb, Rb** — outer lobe tips

Agents use [`PathDecision`](pathfinding.md) at every junction.

## Future

- Embed maps with `go:embed`
- Version field + migrations
- Full scene files: map + entities + events + camera
- Editor “Save map…” / “Load map…”

See also: [paths.md](paths.md), [pathfinding.md](pathfinding.md).
