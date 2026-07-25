package noise

import (
	"math"
	"testing"
)

func TestPerlinDeterministicAndBounded(t *testing.T) {
	cfg := DefaultPerlinConfig()
	cfg.Seed = 42
	cfg.Frequency = 0.05
	p := NewPerlin(cfg)
	a := p.Sample2(10, 20)
	b := p.Sample2(10, 20)
	if a != b {
		t.Fatalf("not deterministic: %v vs %v", a, b)
	}
	if math.Abs(a) > 1.5 {
		t.Fatalf("unexpected magnitude %v", a)
	}
	// Different seed → different value (almost surely).
	cfg.Seed = 43
	q := NewPerlin(cfg)
	if q.Sample2(10, 20) == a {
		t.Fatal("expected seed to change sample")
	}
}

func TestSimplexAndOpenSimplex(t *testing.T) {
	sc := DefaultSimplexConfig()
	sc.Seed = 7
	sc.Octaves = 3
	s := NewSimplex(sc)
	oc := DefaultOpenSimplexConfig()
	oc.Seed = 7
	oc.Octaves = 3
	o := NewOpenSimplex(oc)
	sv := s.Sample2(1.25, -3.5)
	ov := o.Sample2(1.25, -3.5)
	if sv == ov {
		t.Fatal("simplex and opensimplex should differ")
	}
	if math.IsNaN(sv) || math.IsNaN(ov) {
		t.Fatal("NaN")
	}
	s3 := s.Sample3(1, 2, 3)
	o3 := o.Sample3(1, 2, 3)
	if math.IsNaN(s3) || math.IsNaN(o3) {
		t.Fatal("NaN 3d")
	}
}

func TestWorleyModes(t *testing.T) {
	base := DefaultWorleyConfig()
	base.Seed = 99
	base.Frequency = 0.1
	base.Octaves = 1

	checks := []WorleyConfig{base}
	for _, dist := range []DistanceMetric{DistEuclidean, DistManhattan, DistChebyshev, DistEuclideanSq, DistMinkowski, DistHybrid} {
		c := base
		c.Distance = dist
		c.HybridBlend = 0.5
		checks = append(checks, c)
	}
	for _, ret := range []WorleyReturn{ReturnF1, ReturnF2, ReturnF2MinusF1, ReturnF1PlusF2, ReturnCellValue, ReturnDistance2Div} {
		c := base
		c.Return = ret
		checks = append(checks, c)
	}
	for i, c := range checks {
		w := NewWorley(c)
		v := w.Sample2(12.3, 45.6)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("case %d: bad value %v", i, v)
		}
		_ = w.Sample3(1, 2, 3)
	}
}

func TestFractalTypes(t *testing.T) {
	for _, ft := range []FractalType{FractalFBM, FractalBillow, FractalRidged, FractalPingPong} {
		cfg := DefaultPerlinConfig()
		cfg.Seed = 1
		cfg.Type = ft
		cfg.Octaves = 4
		cfg.WeightedStrength = 0.8
		p := NewPerlin(cfg)
		v := p.Sample2(0.5, 0.5)
		if math.IsNaN(v) {
			t.Fatalf("fractal %v NaN", ft)
		}
	}
}

func TestDomainWarpAndOutput(t *testing.T) {
	cfg := DefaultSimplexConfig()
	cfg.Seed = 5
	cfg.Octaves = 2
	cfg.Warp = Warp{Enabled: true, Strength: 8, Frequency: 0.02, Fractal: Fractal{Octaves: 2}}
	cfg.Output = Output{Amplitude: 0.5, Bias: 0.5, To01: true, Clamp01: true}
	s := NewSimplex(cfg)
	v := s.Sample2(100, 100)
	if v < 0 || v > 1 {
		t.Fatalf("expected clamped [0,1], got %v", v)
	}
}

func TestJitterZeroWorley(t *testing.T) {
	c := DefaultWorleyConfig()
	c.Jitter = 0
	c.Seed = 1
	w := NewWorley(c)
	// Grid-centered features should still be finite.
	if math.IsNaN(w.Sample2(0.25, 0.25)) {
		t.Fatal("NaN")
	}
}

func TestSamplerInterface(t *testing.T) {
	var _ Sampler = NewPerlin(DefaultPerlinConfig())
	var _ Sampler = NewSimplex(DefaultSimplexConfig())
	var _ Sampler = NewOpenSimplex(DefaultOpenSimplexConfig())
	var _ Sampler = NewWorley(DefaultWorleyConfig())
}
