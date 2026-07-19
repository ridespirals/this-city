# Map storage format

## Purpose

Persist and share city street networks (nodes + Bézier edges) independently of entity/agent state. Scene save/load (Phase 8) will compose map data with entities/events.

## Current format (`*.json`)

Located under [`maps/`](../maps/). Schema:

```json
{
  "name": "dev-map",
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

## SVG import

Authoring source paths can live as SVG under [`assets/maps/`](../assets/maps/). Each `<path d="...">` becomes a **PathPiece** (sequence of cubic Béziers).

| API | Package |
|-----|---------|
| `ParseSVGPathData`, `ParseSVG`, `MapFileFromSVG` | `internal/sim` |
| `PathPiece`, `StampPiece`, `Recentered` | `internal/sim` |
| Embedded SVGs + `LoadStampPieces` | `assets/maps` (`mapsvg`) |
| CLI converter | `go run ./cmd/svg2map -o maps/out.json assets/maps/in.svg` |

Supported path commands: `M/L/H/V/C/S/Q/T/Z` (absolute and relative). Nodes within `DefaultSVGMergeEps` (1.5) snap together on import/stamp.

**Stamp pieces** are recentered (bbox center at origin) so the editor can place them repeatedly at the cursor. Full maps keep absolute SVG coordinates (typically matching the 1280×720 viewBox).

## Code

| API | Package |
|-----|---------|
| `MapFile`, `LoadNetworkJSON`, `LoadNetworkFile`, `ApplyMapFile`, `ExportMapFile` | `internal/sim/mapfile.go` |
| Default demo map | `maps/dev-map.json` (+ embedded `assets/maps/dev-map.svg`) |
| Figure-8 fallback | `sim.FigureEightMap()` + [`maps/figure-8.json`](../maps/figure-8.json) |
| Demo loader | `game.LoadDemoMap` (JSON → embedded SVG → figure-8) |

## Sample maps

### dev-map (default)

Imported from [`assets/maps/dev-map.svg`](../assets/maps/dev-map.svg) — one closed cubic loop (8 nodes / 8 edges) in the 1280×720 viewBox.

### figure-8

Classic horizontal ∞ kept as a regression / alternate sample:

- **Center** — degree-4 crossing
- **Left / right tips** — lobe endpoints
- **Four arcs** — upper/lower left and upper/lower right

## Future

- Version field + migrations
- Full scene files: map + entities + events + camera
- Editor “Save map…” / “Load map…”
- More SVG: groups, transforms, multi-layer labels

See also: [paths.md](paths.md), [pathfinding.md](pathfinding.md), [editor.md](editor.md).
