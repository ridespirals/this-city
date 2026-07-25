package noise

import "math"

// DistanceMetric selects how Worley measures distance to feature points.
type DistanceMetric int

const (
	DistEuclidean DistanceMetric = iota
	DistEuclideanSq
	DistManhattan
	DistChebyshev
	DistMinkowski // uses MinkowskiP
	DistHybrid    // mix of Euclidean and Manhattan via HybridBlend
)

// WorleyReturn selects which cellular statistic becomes the sample value.
type WorleyReturn int

const (
	// ReturnF1 is distance to nearest feature (classic Worley).
	ReturnF1 WorleyReturn = iota
	// ReturnF2 is distance to second-nearest feature.
	ReturnF2
	// ReturnF2MinusF1 is F2 - F1 (cell edges / cracks).
	ReturnF2MinusF1
	// ReturnF1PlusF2 is F1 + F2.
	ReturnF1PlusF2
	// ReturnF1TimesF2 is F1 * F2.
	ReturnF1TimesF2
	// ReturnCellValue is a stable random value for the nearest cell (voronoi flat).
	ReturnCellValue
	// ReturnDistance2Div is F1 / max(F2, eps) (normalized-ish cracks).
	ReturnDistance2Div
)

// WorleyConfig configures Worley / cellular noise.
type WorleyConfig struct {
	Seeded
	Domain
	Fractal
	Output
	Warp Warp // named to avoid Frequency clash with Domain

	// Distance defaults to DistEuclidean.
	Distance DistanceMetric
	// Return defaults to ReturnF1.
	Return WorleyReturn
	// Jitter places feature points randomly inside each cell (0 = grid centers, 1 = full).
	// Default 1.
	Jitter float64
	// MinkowskiP is used when Distance == DistMinkowski. Default 3.
	MinkowskiP float64
	// HybridBlend mixes Euclidean (0) and Manhattan (1) when Distance == DistHybrid.
	HybridBlend float64
	// DistanceScale multiplies distance before remap. Default 1.
	DistanceScale float64
	// DisableRemap keeps non-negative distance (or cell value) without mapping to [-1,1].
	// Output still applies afterward.
	DisableRemap bool
}

// Worley is a Sampler for cellular noise.
type Worley struct {
	cfg  WorleyConfig
	warp *Worley
}

// DefaultWorleyConfig returns sensible cellular defaults (full jitter, F1 Euclidean).
func DefaultWorleyConfig() WorleyConfig {
	return WorleyConfig{
		Domain:        Domain{Frequency: 1},
		Fractal:       Fractal{Octaves: 1, Lacunarity: 2, Persistence: 0.5},
		Output:        Output{Amplitude: 1},
		Distance:      DistEuclidean,
		Return:        ReturnF1,
		Jitter:        1,
		MinkowskiP:    3,
		DistanceScale: 1,
	}
}

// NewWorley builds a Worley sampler.
func NewWorley(cfg WorleyConfig) *Worley {
	cfg.NormalizeDomain()
	cfg.NormalizeFractal()
	cfg.NormalizeOutput()
	cfg.Warp.NormalizeWarp()
	if cfg.Jitter < 0 {
		cfg.Jitter = 0
	}
	if cfg.Jitter > 1 {
		cfg.Jitter = 1
	}
	if cfg.MinkowskiP == 0 {
		cfg.MinkowskiP = 3
	}
	if cfg.DistanceScale == 0 {
		cfg.DistanceScale = 1
	}
	w := &Worley{cfg: cfg}
	if cfg.Warp.Enabled {
		wcfg := WorleyConfig{
			Seeded:   Seeded{Seed: cfg.Seed + cfg.Warp.SeedOffset + 7919},
			Domain:   Domain{Frequency: cfg.Warp.Frequency},
			Fractal:  cfg.Warp.Fractal,
			Output:   Output{Amplitude: 1},
			Distance: DistEuclidean,
			Return:   ReturnF1,
			Jitter:   1,
		}
		w.warp = NewWorley(wcfg)
	}
	return w
}

// Sample2 implements Sampler.
func (w *Worley) Sample2(x, y float64) float64 {
	if w == nil {
		return 0
	}
	if w.warp != nil {
		wx := w.warp.Sample2(x, y)
		wy := w.warp.Sample2(x+23.1, y+41.7)
		x += wx * w.cfg.Warp.Strength
		y += wy * w.cfg.Warp.Strength
	}
	x, y = w.cfg.map2(x, y)
	v := applyFractal2(w.cfg.Fractal, w.raw2, x, y)
	return w.cfg.apply(v)
}

// Sample3 implements Sampler.
func (w *Worley) Sample3(x, y, z float64) float64 {
	if w == nil {
		return 0
	}
	if w.warp != nil {
		wx := w.warp.Sample3(x, y, z)
		wy := w.warp.Sample3(x+23.1, y+41.7, z+7.9)
		wz := w.warp.Sample3(x+11.3, y+29.5, z+53.1)
		s := w.cfg.Warp.Strength
		x += wx * s
		y += wy * s
		z += wz * s
	}
	x, y, z = w.cfg.map3(x, y, z)
	v := applyFractal3(w.cfg.Fractal, w.raw3, x, y, z)
	return w.cfg.apply(v)
}

func (w *Worley) raw2(x, y float64) float64 {
	xi := fastFloor(x)
	yi := fastFloor(y)
	xf := x - float64(xi)
	yf := y - float64(yi)

	f1, f2 := math.MaxFloat64, math.MaxFloat64
	var nearestHash uint64

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			cx := xi + dx
			cy := yi + dy
			h := cellSeed(w.cfg.Seed, cx, cy, 0)
			rx, h := rand01(h)
			ry, h := rand01(h)
			px := float64(dx) + 0.5 + (rx-0.5)*w.cfg.Jitter
			py := float64(dy) + 0.5 + (ry-0.5)*w.cfg.Jitter
			d := w.distance(xf-px, yf-py, 0)
			if d < f1 {
				f2 = f1
				f1 = d
				nearestHash = h
			} else if d < f2 {
				f2 = d
			}
		}
	}
	return w.combine(f1, f2, nearestHash)
}

func (w *Worley) raw3(x, y, z float64) float64 {
	xi := fastFloor(x)
	yi := fastFloor(y)
	zi := fastFloor(z)
	xf := x - float64(xi)
	yf := y - float64(yi)
	zf := z - float64(zi)

	f1, f2 := math.MaxFloat64, math.MaxFloat64
	var nearestHash uint64

	for dz := -1; dz <= 1; dz++ {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				cx := xi + dx
				cy := yi + dy
				cz := zi + dz
				h := cellSeed(w.cfg.Seed, cx, cy, cz)
				rx, h := rand01(h)
				ry, h := rand01(h)
				rz, h := rand01(h)
				px := float64(dx) + 0.5 + (rx-0.5)*w.cfg.Jitter
				py := float64(dy) + 0.5 + (ry-0.5)*w.cfg.Jitter
				pz := float64(dz) + 0.5 + (rz-0.5)*w.cfg.Jitter
				d := w.distance(xf-px, yf-py, zf-pz)
				if d < f1 {
					f2 = f1
					f1 = d
					nearestHash = h
				} else if d < f2 {
					f2 = d
				}
			}
		}
	}
	return w.combine(f1, f2, nearestHash)
}

func (w *Worley) distance(dx, dy, dz float64) float64 {
	switch w.cfg.Distance {
	case DistEuclideanSq:
		return dx*dx + dy*dy + dz*dz
	case DistManhattan:
		return math.Abs(dx) + math.Abs(dy) + math.Abs(dz)
	case DistChebyshev:
		return math.Max(math.Abs(dx), math.Max(math.Abs(dy), math.Abs(dz)))
	case DistMinkowski:
		p := w.cfg.MinkowskiP
		return math.Pow(math.Pow(math.Abs(dx), p)+math.Pow(math.Abs(dy), p)+math.Pow(math.Abs(dz), p), 1/p)
	case DistHybrid:
		e := math.Sqrt(dx*dx + dy*dy + dz*dz)
		m := math.Abs(dx) + math.Abs(dy) + math.Abs(dz)
		t := clamp01(w.cfg.HybridBlend)
		return lerp(e, m, t)
	default: // DistEuclidean
		return math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
}

func (w *Worley) combine(f1, f2 float64, nearestHash uint64) float64 {
	var v float64
	switch w.cfg.Return {
	case ReturnF2:
		v = f2
	case ReturnF2MinusF1:
		v = f2 - f1
	case ReturnF1PlusF2:
		v = f1 + f2
	case ReturnF1TimesF2:
		v = f1 * f2
	case ReturnCellValue:
		// Map hash to [-1,1]
		r, _ := rand01(nearestHash ^ 0xA24BAED4963EE407)
		return r*2 - 1
	case ReturnDistance2Div:
		if f2 < 1e-12 {
			v = 0
		} else {
			v = f1 / f2
		}
	default: // ReturnF1
		v = f1
	}

	if w.cfg.DisableRemap {
		return v
	}
	// Remap typical distance into approx [-1, 1]: close → +1, far → -1.
	t := clamp01(v * w.cfg.DistanceScale)
	return 1 - 2*t
}
