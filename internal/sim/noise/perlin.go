package noise

// PerlinConfig configures classic improved Perlin noise (Ken Perlin, 2002).
type PerlinConfig struct {
	Seeded
	Domain
	Fractal
	Output
	Warp Warp // named to avoid Frequency clash with Domain
}

// Perlin is a Sampler for improved Perlin noise.
type Perlin struct {
	cfg  PerlinConfig
	perm permTable
	warp *Perlin // optional warp field (nil if Warp disabled)
}

// NewPerlin builds a Perlin sampler from cfg.
func NewPerlin(cfg PerlinConfig) *Perlin {
	cfg.NormalizeDomain()
	cfg.NormalizeFractal()
	cfg.NormalizeOutput()
	cfg.Warp.NormalizeWarp()
	p := &Perlin{cfg: cfg, perm: newPerm(cfg.Seed)}
	if cfg.Warp.Enabled {
		wcfg := PerlinConfig{
			Seeded:  Seeded{Seed: cfg.Seed + cfg.Warp.SeedOffset + 1013},
			Domain:  Domain{Frequency: cfg.Warp.Frequency},
			Fractal: cfg.Warp.Fractal,
			Output:  Output{Amplitude: 1},
		}
		p.warp = NewPerlin(wcfg)
	}
	return p
}

// Sample2 implements Sampler.
func (p *Perlin) Sample2(x, y float64) float64 {
	if p == nil {
		return 0
	}
	if p.warp != nil {
		wx := p.warp.Sample2(x, y)
		wy := p.warp.Sample2(x+19.1, y+67.3)
		x += wx * p.cfg.Warp.Strength
		y += wy * p.cfg.Warp.Strength
	}
	x, y = p.cfg.map2(x, y)
	v := applyFractal2(p.cfg.Fractal, p.raw2, x, y)
	return p.cfg.apply(v)
}

// Sample3 implements Sampler.
func (p *Perlin) Sample3(x, y, z float64) float64 {
	if p == nil {
		return 0
	}
	if p.warp != nil {
		wx := p.warp.Sample3(x, y, z)
		wy := p.warp.Sample3(x+19.1, y+67.3, z+11.7)
		wz := p.warp.Sample3(x+41.2, y+13.9, z+53.1)
		s := p.cfg.Warp.Strength
		x += wx * s
		y += wy * s
		z += wz * s
	}
	x, y, z = p.cfg.map3(x, y, z)
	v := applyFractal3(p.cfg.Fractal, p.raw3, x, y, z)
	return p.cfg.apply(v)
}

func (p *Perlin) raw2(x, y float64) float64 {
	x0 := fastFloor(x)
	y0 := fastFloor(y)
	xf := x - float64(x0)
	yf := y - float64(y0)
	u := fade(xf)
	v := fade(yf)

	aa := p.perm.hash2(x0, y0)
	ab := p.perm.hash2(x0, y0+1)
	ba := p.perm.hash2(x0+1, y0)
	bb := p.perm.hash2(x0+1, y0+1)

	x1 := lerp(grad2(aa, xf, yf), grad2(ba, xf-1, yf), u)
	x2 := lerp(grad2(ab, xf, yf-1), grad2(bb, xf-1, yf-1), u)
	// Improved Perlin 2D is roughly [-1,1]; slight scale keeps peaks near that range.
	return lerp(x1, x2, v)
}

func (p *Perlin) raw3(x, y, z float64) float64 {
	x0 := fastFloor(x)
	y0 := fastFloor(y)
	z0 := fastFloor(z)
	xf := x - float64(x0)
	yf := y - float64(y0)
	zf := z - float64(z0)
	u := fade(xf)
	v := fade(yf)
	w := fade(zf)

	n000 := grad3(p.perm.hash3(x0, y0, z0), xf, yf, zf)
	n100 := grad3(p.perm.hash3(x0+1, y0, z0), xf-1, yf, zf)
	n010 := grad3(p.perm.hash3(x0, y0+1, z0), xf, yf-1, zf)
	n110 := grad3(p.perm.hash3(x0+1, y0+1, z0), xf-1, yf-1, zf)
	n001 := grad3(p.perm.hash3(x0, y0, z0+1), xf, yf, zf-1)
	n101 := grad3(p.perm.hash3(x0+1, y0, z0+1), xf-1, yf, zf-1)
	n011 := grad3(p.perm.hash3(x0, y0+1, z0+1), xf, yf-1, zf-1)
	n111 := grad3(p.perm.hash3(x0+1, y0+1, z0+1), xf-1, yf-1, zf-1)

	nx00 := lerp(n000, n100, u)
	nx10 := lerp(n010, n110, u)
	nx01 := lerp(n001, n101, u)
	nx11 := lerp(n011, n111, u)
	nxy0 := lerp(nx00, nx10, v)
	nxy1 := lerp(nx01, nx11, v)
	return lerp(nxy0, nxy1, w)
}

func fastFloor(x float64) int {
	i := int(x)
	if float64(i) > x {
		return i - 1
	}
	return i
}
