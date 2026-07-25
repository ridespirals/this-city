package noise

import "math"

// SimplexConfig configures classic simplex noise (skewed grid gradient noise).
// Note: the original simplex patent expired in 2022; OpenSimplex is also provided.
type SimplexConfig struct {
	Seeded
	Domain
	Fractal
	Output
	Warp Warp // named to avoid Frequency clash with Domain
}

// Simplex is a Sampler for classic 2D/3D simplex noise.
type Simplex struct {
	cfg  SimplexConfig
	perm permTable
	warp *Simplex
}

// NewSimplex builds a Simplex sampler.
func NewSimplex(cfg SimplexConfig) *Simplex {
	cfg.NormalizeDomain()
	cfg.NormalizeFractal()
	cfg.NormalizeOutput()
	cfg.Warp.NormalizeWarp()
	s := &Simplex{cfg: cfg, perm: newPerm(cfg.Seed)}
	if cfg.Warp.Enabled {
		wcfg := SimplexConfig{
			Seeded:  Seeded{Seed: cfg.Seed + cfg.Warp.SeedOffset + 2029},
			Domain:  Domain{Frequency: cfg.Warp.Frequency},
			Fractal: cfg.Warp.Fractal,
			Output:  Output{Amplitude: 1},
		}
		s.warp = NewSimplex(wcfg)
	}
	return s
}

// Sample2 implements Sampler.
func (s *Simplex) Sample2(x, y float64) float64 {
	if s == nil {
		return 0
	}
	if s.warp != nil {
		wx := s.warp.Sample2(x, y)
		wy := s.warp.Sample2(x+19.1, y+67.3)
		x += wx * s.cfg.Warp.Strength
		y += wy * s.cfg.Warp.Strength
	}
	x, y = s.cfg.map2(x, y)
	v := applyFractal2(s.cfg.Fractal, s.raw2, x, y)
	return s.cfg.apply(v)
}

// Sample3 implements Sampler.
func (s *Simplex) Sample3(x, y, z float64) float64 {
	if s == nil {
		return 0
	}
	if s.warp != nil {
		wx := s.warp.Sample3(x, y, z)
		wy := s.warp.Sample3(x+19.1, y+67.3, z+11.7)
		wz := s.warp.Sample3(x+41.2, y+13.9, z+53.1)
		st := s.cfg.Warp.Strength
		x += wx * st
		y += wy * st
		z += wz * st
	}
	x, y, z = s.cfg.map3(x, y, z)
	v := applyFractal3(s.cfg.Fractal, s.raw3, x, y, z)
	return s.cfg.apply(v)
}

const (
	f2 = 0.5 * (math.Sqrt2 - 1) // skew 2D
	g2 = (3 - math.Sqrt2) / 6   // unskew 2D
	f3 = 1.0 / 3.0
	g3 = 1.0 / 6.0
)

func (s *Simplex) raw2(x, y float64) float64 {
	s0 := (x + y) * f2
	i := fastFloor(x + s0)
	j := fastFloor(y + s0)
	t0 := float64(i+j) * g2
	x0 := x - (float64(i) - t0)
	y0 := y - (float64(j) - t0)

	var i1, j1 int
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}
	x1 := x0 - float64(i1) + g2
	y1 := y0 - float64(j1) + g2
	x2 := x0 - 1 + 2*g2
	y2 := y0 - 1 + 2*g2

	n0 := corner2(s.perm, i, j, x0, y0)
	n1 := corner2(s.perm, i+i1, j+j1, x1, y1)
	n2 := corner2(s.perm, i+1, j+1, x2, y2)
	// Scale to approx [-1, 1]
	return 70 * (n0 + n1 + n2)
}

func corner2(p permTable, ix, iy int, x, y float64) float64 {
	t := 0.5 - x*x - y*y
	if t < 0 {
		return 0
	}
	t *= t
	return t * t * grad2(p.hash2(ix, iy), x, y)
}

func (s *Simplex) raw3(x, y, z float64) float64 {
	s0 := (x + y + z) * f3
	i := fastFloor(x + s0)
	j := fastFloor(y + s0)
	k := fastFloor(z + s0)
	t0 := float64(i+j+k) * g3
	x0 := x - (float64(i) - t0)
	y0 := y - (float64(j) - t0)
	z0 := z - (float64(k) - t0)

	var i1, j1, k1, i2, j2, k2 int
	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1, i2, j2, k2 = 1, 0, 0, 1, 0, 1
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 1, 0, 1
		}
	} else {
		if y0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 0, 1, 0, 1, 1
		} else if x0 < z0 {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 0, 1, 1
		} else {
			i1, j1, k1, i2, j2, k2 = 0, 1, 0, 1, 1, 0
		}
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

	n0 := corner3(s.perm, i, j, k, x0, y0, z0)
	n1 := corner3(s.perm, i+i1, j+j1, k+k1, x1, y1, z1)
	n2 := corner3(s.perm, i+i2, j+j2, k+k2, x2, y2, z2)
	n3 := corner3(s.perm, i+1, j+1, k+1, x3, y3, z3)
	return 32 * (n0 + n1 + n2 + n3)
}

func corner3(p permTable, ix, iy, iz int, x, y, z float64) float64 {
	t := 0.6 - x*x - y*y - z*z
	if t < 0 {
		return 0
	}
	t *= t
	return t * t * grad3(p.hash3(ix, iy, iz), x, y, z)
}
