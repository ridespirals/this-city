package noise

import "math"

// OpenSimplexConfig configures OpenSimplex-style noise (lattice-based, patent-unencumbered family).
// This is a compact 2D/3D OpenSimplex2-inspired implementation suitable for gameplay content.
type OpenSimplexConfig struct {
	Seeded
	Domain
	Fractal
	Output
	Warp Warp // named to avoid Frequency clash with Domain
}

// OpenSimplex is a Sampler for OpenSimplex-family noise.
type OpenSimplex struct {
	cfg  OpenSimplexConfig
	perm permTable
	warp *OpenSimplex
}

// NewOpenSimplex builds an OpenSimplex sampler.
func NewOpenSimplex(cfg OpenSimplexConfig) *OpenSimplex {
	cfg.NormalizeDomain()
	cfg.NormalizeFractal()
	cfg.NormalizeOutput()
	cfg.Warp.NormalizeWarp()
	o := &OpenSimplex{cfg: cfg, perm: newPerm(cfg.Seed ^ 0x6C62272E07BB0142)}
	if cfg.Warp.Enabled {
		wcfg := OpenSimplexConfig{
			Seeded:  Seeded{Seed: cfg.Seed + cfg.Warp.SeedOffset + 4099},
			Domain:  Domain{Frequency: cfg.Warp.Frequency},
			Fractal: cfg.Warp.Fractal,
			Output:  Output{Amplitude: 1},
		}
		o.warp = NewOpenSimplex(wcfg)
	}
	return o
}

// Sample2 implements Sampler.
func (o *OpenSimplex) Sample2(x, y float64) float64 {
	if o == nil {
		return 0
	}
	if o.warp != nil {
		wx := o.warp.Sample2(x, y)
		wy := o.warp.Sample2(x+31.7, y+47.9)
		x += wx * o.cfg.Warp.Strength
		y += wy * o.cfg.Warp.Strength
	}
	x, y = o.cfg.map2(x, y)
	v := applyFractal2(o.cfg.Fractal, o.raw2, x, y)
	return o.cfg.apply(v)
}

// Sample3 implements Sampler.
func (o *OpenSimplex) Sample3(x, y, z float64) float64 {
	if o == nil {
		return 0
	}
	if o.warp != nil {
		wx := o.warp.Sample3(x, y, z)
		wy := o.warp.Sample3(x+31.7, y+47.9, z+12.3)
		wz := o.warp.Sample3(x+17.1, y+29.5, z+43.7)
		s := o.cfg.Warp.Strength
		x += wx * s
		y += wy * s
		z += wz * s
	}
	x, y, z = o.cfg.map3(x, y, z)
	v := applyFractal3(o.cfg.Fractal, o.raw3, x, y, z)
	return o.cfg.apply(v)
}

// OpenSimplex2-style 2D: evaluate contributions on a triangular lattice with
// stretched square grid (similar spirit to KS OpenSimplex2 "fast" 2D).
var (
	os2Stretch = (1 / math.Sqrt(3)) - 1
	os2Squish  = (math.Sqrt(3) - 1) / 2
	os2Norm    = 47.0
)

func (o *OpenSimplex) raw2(x, y float64) float64 {
	// Stretch
	s := (x + y) * os2Stretch
	xs := x + s
	ys := y + s
	xsb := fastFloor(xs)
	ysb := fastFloor(ys)

	// Squish
	sq := float64(xsb+ysb) * os2Squish
	xb := float64(xsb) - sq
	yb := float64(ysb) - sq
	dx0 := x - xb
	dy0 := y - yb

	xins := xs - float64(xsb)
	yins := ys - float64(ysb)
	inSum := xins + yins

	var value float64
	value += o.osContribute2(xsb, ysb, dx0, dy0)
	value += o.osContribute2(xsb+1, ysb+1, dx0-1-2*os2Squish, dy0-1-2*os2Squish)

	value += o.osContribute2(xsb+1, ysb, dx0-1-os2Squish, dy0-os2Squish)
	value += o.osContribute2(xsb, ysb+1, dx0-os2Squish, dy0-1-os2Squish)
	_ = inSum
	return value * os2Norm
}

func (o *OpenSimplex) osContribute2(ix, iy int, dx, dy float64) float64 {
	attn := 2.0/3.0 - dx*dx - dy*dy
	if attn <= 0 {
		return 0
	}
	attn *= attn
	return attn * attn * grad2(o.perm.hash2(ix, iy), dx, dy)
}

func (o *OpenSimplex) raw3(x, y, z float64) float64 {
	// Rotate to improve lattice orientation, then use simplex-like evaluation
	// with OpenSimplex-ish constants (practical compromise for 3D content).
	r := (2.0 / 3.0)
	xr := r*(x+y+z) - x
	yr := r*(x+y+z) - y
	zr := r*(x+y+z) - z

	// Reuse simplex 3D structure with separate permutation bias.
	s0 := (xr + yr + zr) * f3
	i := fastFloor(xr + s0)
	j := fastFloor(yr + s0)
	k := fastFloor(zr + s0)
	t0 := float64(i+j+k) * g3
	x0 := xr - (float64(i) - t0)
	y0 := yr - (float64(j) - t0)
	z0 := zr - (float64(k) - t0)

	var i1, j1, k1, i2, j2, k2 int
	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 0, 1
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 1, 0, 1
		}
	} else if y0 < z0 {
		i1, j1, k1, i2, j2, k2 = 0, 0, 1, 0, 1, 1
	} else if x0 < z0 {
		i1, j1, k1, i2, j2, k2 = 0, 1, 0, 0, 1, 1
	} else {
		i1, j1, k1, i2, j2, k2 = 0, 1, 0, 1, 1, 0
	}

	x1 := x0 - float64(i1) + g3
	y1 := y0 - float64(j1) + g3
	z1 := z0 - float64(k1) + g3
	x2 := x0 - float64(i2) + 2*g3
	y2 := y0 - float64(j2) + 2*g3
	z2 := z0 - float64(k2) + 2*g3
	x3 := x0 - 1 + 3*g3
	y3 := y0 - 1 + 3*g3
	z3 := z0 - 1 + 3*g3

	n0 := corner3(o.perm, i, j, k, x0, y0, z0)
	n1 := corner3(o.perm, i+i1, j+j1, k+k1, x1, y1, z1)
	n2 := corner3(o.perm, i+i2, j+j2, k+k2, x2, y2, z2)
	n3 := corner3(o.perm, i+1, j+1, k+1, x3, y3, z3)
	return 32 * (n0 + n1 + n2 + n3)
}
