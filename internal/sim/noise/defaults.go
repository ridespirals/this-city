package noise

// DefaultPerlinConfig returns FBM-ready Perlin defaults.
func DefaultPerlinConfig() PerlinConfig {
	return PerlinConfig{
		Domain:  Domain{Frequency: 1},
		Fractal: Fractal{Type: FractalFBM, Octaves: 4, Lacunarity: 2, Persistence: 0.5},
		Output:  Output{Amplitude: 1},
	}
}

// DefaultSimplexConfig returns FBM-ready Simplex defaults.
func DefaultSimplexConfig() SimplexConfig {
	return SimplexConfig{
		Domain:  Domain{Frequency: 1},
		Fractal: Fractal{Type: FractalFBM, Octaves: 4, Lacunarity: 2, Persistence: 0.5},
		Output:  Output{Amplitude: 1},
	}
}

// DefaultOpenSimplexConfig returns FBM-ready OpenSimplex defaults.
func DefaultOpenSimplexConfig() OpenSimplexConfig {
	return OpenSimplexConfig{
		Domain:  Domain{Frequency: 1},
		Fractal: Fractal{Type: FractalFBM, Octaves: 4, Lacunarity: 2, Persistence: 0.5},
		Output:  Output{Amplitude: 1},
	}
}
