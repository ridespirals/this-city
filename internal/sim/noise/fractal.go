package noise

import "math"

// base2/base3 are single-octave evaluators in "noise space" (already domain-mapped).
type base2 func(x, y float64) float64
type base3 func(x, y, z float64) float64

func applyFractal2(f Fractal, fn base2, x, y float64) float64 {
	f.NormalizeFractal()
	switch f.Type {
	case FractalBillow:
		return fractalBillow2(f, fn, x, y)
	case FractalRidged:
		return fractalRidged2(f, fn, x, y)
	case FractalPingPong:
		return fractalPingPong2(f, fn, x, y)
	default:
		return fractalFBM2(f, fn, x, y)
	}
}

func applyFractal3(f Fractal, fn base3, x, y, z float64) float64 {
	f.NormalizeFractal()
	switch f.Type {
	case FractalBillow:
		return fractalBillow3(f, fn, x, y, z)
	case FractalRidged:
		return fractalRidged3(f, fn, x, y, z)
	case FractalPingPong:
		return fractalPingPong3(f, fn, x, y, z)
	default:
		return fractalFBM3(f, fn, x, y, z)
	}
}

func fractalFBM2(f Fractal, fn base2, x, y float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := fn(x*freq, y*freq)
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}

func fractalFBM3(f Fractal, fn base3, x, y, z float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := fn(x*freq, y*freq, z*freq)
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}

func fractalBillow2(f Fractal, fn base2, x, y float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := math.Abs(fn(x*freq, y*freq))*2 - 1
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}

func fractalBillow3(f Fractal, fn base3, x, y, z float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := math.Abs(fn(x*freq, y*freq, z*freq))*2 - 1
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}

func fractalRidged2(f Fractal, fn base2, x, y float64) float64 {
	var sum, amp, freq, weight float64 = 0, 1, 1, 1
	var max float64
	for i := 0; i < f.Octaves; i++ {
		n := fn(x*freq, y*freq)
		n = 1 - math.Abs(n)
		n *= n
		n *= weight
		weight = clamp01(n * f.WeightedStrength)
		if f.WeightedStrength == 0 {
			weight = 1
		}
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	// Ridged is typically [0,1]; remap to [-1,1] for Output consistency.
	return sum*2 - 1
}

func fractalRidged3(f Fractal, fn base3, x, y, z float64) float64 {
	var sum, amp, freq, weight float64 = 0, 1, 1, 1
	var max float64
	for i := 0; i < f.Octaves; i++ {
		n := fn(x*freq, y*freq, z*freq)
		n = 1 - math.Abs(n)
		n *= n
		n *= weight
		weight = clamp01(n * f.WeightedStrength)
		if f.WeightedStrength == 0 {
			weight = 1
		}
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum*2 - 1
}

func pingPong(t, s float64) float64 {
	t = t*s + 1
	// triangle fold into [0,1] then to [-1,1]
	t = t - 2*math.Floor(t*0.5)
	if t < 0 {
		t = -t
	}
	if t > 1 {
		t = 2 - t
	}
	return t*2 - 1
}

func fractalPingPong2(f Fractal, fn base2, x, y float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := pingPong(fn(x*freq, y*freq), f.PingPongStrength)
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}

func fractalPingPong3(f Fractal, fn base3, x, y, z float64) float64 {
	var sum, amp, freq, max float64 = 0, 1, 1, 0
	for i := 0; i < f.Octaves; i++ {
		n := pingPong(fn(x*freq, y*freq, z*freq), f.PingPongStrength)
		sum += n * amp
		max += amp
		freq *= f.Lacunarity
		amp *= f.Persistence
	}
	if max > 0 {
		sum /= max
	}
	return sum
}
