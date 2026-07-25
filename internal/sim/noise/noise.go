// Package noise provides deterministic procedural noise for simulation content.
// No raylib; float64 samples in roughly [-1, 1] unless remapped by Output.
package noise

import "math"

// Sampler evaluates noise at a point. Implementations are safe for concurrent
// read-only use after construction (immutable tables).
type Sampler interface {
	// Sample2 returns noise at (x, y).
	Sample2(x, y float64) float64
	// Sample3 returns noise at (x, y, z).
	Sample3(x, y, z float64) float64
}

// FractalType selects how octaves are combined.
type FractalType int

const (
	// FractalFBM is standard fractional Brownian motion (sum of octaves).
	FractalFBM FractalType = iota
	// FractalBillow uses absolute value per octave (billowy / cloud-like).
	FractalBillow
	// FractalRidged is ridged multi-fractal (1 - |n|)^2 style ridges.
	FractalRidged
	// FractalPingPong folds the signal for a ping-pong / terrace look.
	FractalPingPong
)

// Fractal is multi-octave layering applied on top of a base noise function.
type Fractal struct {
	// Type defaults to FractalFBM.
	Type FractalType
	// Octaves is the number of layers (clamped to >= 1). Default 1 = no fractal.
	Octaves int
	// Lacunarity multiplies frequency each octave. Default 2.
	Lacunarity float64
	// Persistence (gain) multiplies amplitude each octave. Default 0.5.
	Persistence float64
	// WeightedStrength blends toward higher-octave weighting (0 = off, 1 = strong).
	// Used mainly with ridged/billow; 0 keeps classic behavior.
	WeightedStrength float64
	// PingPongStrength controls fold frequency for FractalPingPong. Default 2.
	PingPongStrength float64
}

// NormalizeFractal fills defaults and clamps.
func (f *Fractal) NormalizeFractal() {
	if f.Octaves < 1 {
		f.Octaves = 1
	}
	if f.Lacunarity == 0 {
		f.Lacunarity = 2
	}
	if f.Persistence == 0 {
		f.Persistence = 0.5
	}
	if f.PingPongStrength == 0 {
		f.PingPongStrength = 2
	}
}

// Domain is spatial transform applied before sampling.
type Domain struct {
	// Frequency scales input coordinates. Default 1.
	Frequency float64
	// Offset shifts the sample domain.
	OffsetX, OffsetY, OffsetZ float64
	// Rotation degrees about Z (2D) applied after offset, before frequency.
	RotationDeg float64
	// ScaleX/ScaleY/ScaleZ stretch axes (1 = uniform). Useful for elongated features.
	ScaleX, ScaleY, ScaleZ float64
}

// NormalizeDomain fills defaults.
func (d *Domain) NormalizeDomain() {
	if d.Frequency == 0 {
		d.Frequency = 1
	}
	if d.ScaleX == 0 {
		d.ScaleX = 1
	}
	if d.ScaleY == 0 {
		d.ScaleY = 1
	}
	if d.ScaleZ == 0 {
		d.ScaleZ = 1
	}
}

func (d Domain) map2(x, y float64) (float64, float64) {
	x += d.OffsetX
	y += d.OffsetY
	if d.RotationDeg != 0 {
		rad := d.RotationDeg * math.Pi / 180
		c, s := math.Cos(rad), math.Sin(rad)
		x, y = x*c-y*s, x*s+y*c
	}
	x *= d.ScaleX * d.Frequency
	y *= d.ScaleY * d.Frequency
	return x, y
}

func (d Domain) map3(x, y, z float64) (float64, float64, float64) {
	x += d.OffsetX
	y += d.OffsetY
	z += d.OffsetZ
	if d.RotationDeg != 0 {
		rad := d.RotationDeg * math.Pi / 180
		c, s := math.Cos(rad), math.Sin(rad)
		x, y = x*c-y*s, x*s+y*c
	}
	x *= d.ScaleX * d.Frequency
	y *= d.ScaleY * d.Frequency
	z *= d.ScaleZ * d.Frequency
	return x, y, z
}

// Output remaps the raw noise value after fractal combination.
type Output struct {
	// Amplitude multiplies the signal. Default 1.
	Amplitude float64
	// Bias is added after amplitude. Default 0.
	Bias float64
	// Absolute takes |n| before amplitude.
	Absolute bool
	// Invert negates after absolute (before amplitude).
	Invert bool
	// Clamp01 clamps final value to [0, 1] after remap.
	Clamp01 bool
	// To01 remaps assumed [-1,1] input to [0,1] before amplitude/bias.
	To01 bool
	// Power applies math.Pow(abs(n), Power)*sign after fractal (0/1 = off).
	Power float64
}

// NormalizeOutput fills defaults.
func (o *Output) NormalizeOutput() {
	if o.Amplitude == 0 {
		o.Amplitude = 1
	}
}

func (o Output) apply(v float64) float64 {
	if o.Absolute {
		v = math.Abs(v)
	}
	if o.Invert {
		v = -v
	}
	if o.Power != 0 && o.Power != 1 {
		sign := 1.0
		if v < 0 {
			sign = -1
			v = -v
		}
		v = sign * math.Pow(v, o.Power)
	}
	if o.To01 {
		v = v*0.5 + 0.5
	}
	v = v*o.Amplitude + o.Bias
	if o.Clamp01 {
		if v < 0 {
			v = 0
		} else if v > 1 {
			v = 1
		}
	}
	return v
}

// Warp is a domain-warp layer: sample warp noise to displace coordinates,
// then sample the source.
type Warp struct {
	// Enabled turns domain warp on.
	Enabled bool
	// Strength is displacement scale in world units (before source frequency).
	Strength float64
	// Fractal octaves for the warp field itself.
	Fractal Fractal
	// Frequency of the warp field (independent of source Domain.Frequency).
	Frequency float64
	// SeedOffset is added to the source seed for the warp RNG (decorrelate).
	SeedOffset int64
}

// NormalizeWarp fills defaults.
func (w *Warp) NormalizeWarp() {
	if w.Frequency == 0 {
		w.Frequency = 1
	}
	if w.Strength == 0 && w.Enabled {
		w.Strength = 1
	}
	w.Fractal.NormalizeFractal()
}

// Seeded holds a 64-bit seed used to build permutation / hash state.
type Seeded struct {
	Seed int64
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func fade(t float64) float64 { return t * t * t * (t*(t*6-15) + 10) }

func grad2(h int, x, y float64) float64 {
	switch h & 3 {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	default:
		return -x - y
	}
}

func grad3(h int, x, y, z float64) float64 {
	switch h & 15 {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	case 3:
		return -x - y
	case 4:
		return x + z
	case 5:
		return -x + z
	case 6:
		return x - z
	case 7:
		return -x - z
	case 8:
		return y + z
	case 9:
		return -y + z
	case 10:
		return y - z
	case 11:
		return -y - z
	case 12:
		return x + y
	case 13:
		return -x + y
	case 14:
		return -y + z
	default:
		return -y - z
	}
}
