# Procedural noise

## Purpose

Deterministic scalar fields for future content (terrain hints, density maps, scatter, variation). Lives in pure `sim` — no raylib.

## Package

`internal/sim/noise`

| Sampler | Constructor | Notes |
|---------|-------------|-------|
| Perlin | `NewPerlin` | Improved Perlin (2002), 2D + 3D |
| Simplex | `NewSimplex` | Classic simplex lattice, 2D + 3D |
| OpenSimplex | `NewOpenSimplex` | OpenSimplex-family alternative, 2D + 3D |
| Worley | `NewWorley` | Cellular / Voronoi distances |

All implement `noise.Sampler` (`Sample2` / `Sample3`). Defaults: `DefaultPerlinConfig`, `DefaultSimplexConfig`, `DefaultOpenSimplexConfig`, `DefaultWorleyConfig`.

## Configuration layers (shared)

Each algorithm config embeds / includes:

| Layer | Controls |
|-------|----------|
| `Seeded` | `Seed` |
| `Domain` | `Frequency`, `Offset*`, `RotationDeg`, `ScaleX/Y/Z` |
| `Fractal` | `Type` (FBM / Billow / Ridged / PingPong), `Octaves`, `Lacunarity`, `Persistence`, `WeightedStrength`, `PingPongStrength` |
| `Output` | `Amplitude`, `Bias`, `Absolute`, `Invert`, `To01`, `Clamp01`, `Power` |
| `Warp` | Domain warp: `Enabled`, `Strength`, `Frequency`, `Fractal`, `SeedOffset` |

## Worley-only

| Field | Options |
|-------|---------|
| `Distance` | Euclidean, EuclideanSq, Manhattan, Chebyshev, Minkowski (`MinkowskiP`), Hybrid (`HybridBlend`) |
| `Return` | F1, F2, F2−F1, F1+F2, F1×F2, CellValue, F1/F2 |
| `Jitter` | 0 (grid centers) … 1 (full cell) |
| `DistanceScale` / `DisableRemap` | map distances into ~[-1, 1] or leave raw |

## Example

```go
cfg := noise.DefaultSimplexConfig()
cfg.Seed = 42
cfg.Frequency = 0.02
cfg.Octaves = 5
cfg.Type = noise.FractalRidged
cfg.Warp = noise.Warp{Enabled: true, Strength: 12, Frequency: 0.01}
cfg.Output = noise.Output{To01: true, Clamp01: true}
n := noise.NewSimplex(cfg)
v := n.Sample2(x, y) // ~[0,1]
```

## Non-goals (for now)

- GPU shaders / texture baking UI
- Wiring into map generation (next consumer)

See also: [architecture.md](architecture.md).
