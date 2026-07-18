package sim

// CubicBezier is a cubic Bézier segment (p0 → p1 with controls c0, c1).
type CubicBezier struct {
	P0, C0, C1, P1 Vec2
}

// Point evaluates the curve at t in [0,1].
func (b CubicBezier) Point(t float32) Vec2 {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	u := 1 - t
	uu := u * u
	uuu := uu * u
	tt := t * t
	ttt := tt * t
	// (1-t)^3 P0 + 3(1-t)^2 t C0 + 3(1-t) t^2 C1 + t^3 P1
	return b.P0.Scale(uuu).
		Add(b.C0.Scale(3 * uu * t)).
		Add(b.C1.Scale(3 * u * tt)).
		Add(b.P1.Scale(ttt))
}

// Sample returns steps+1 points from t=0..1 (inclusive). steps must be >= 1.
func (b CubicBezier) Sample(steps int) []Vec2 {
	if steps < 1 {
		steps = 1
	}
	out := make([]Vec2, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float32(i) / float32(steps)
		out = append(out, b.Point(t))
	}
	return out
}

// Polyline is a chord-length approximation of one or more Bézier segments.
type Polyline struct {
	Points []Vec2
	CumLen []float32 // CumLen[i] = distance from start to Points[i]
	Length float32
}

// BuildPolyline samples each segment with stepsPerSeg subdivisions and joins them.
// Adjacent shared endpoints are not duplicated.
func BuildPolyline(segments []CubicBezier, stepsPerSeg int) Polyline {
	if stepsPerSeg < 1 {
		stepsPerSeg = 1
	}
	var pts []Vec2
	for i, seg := range segments {
		samples := seg.Sample(stepsPerSeg)
		if i == 0 {
			pts = append(pts, samples...)
			continue
		}
		if len(samples) > 0 {
			pts = append(pts, samples[1:]...)
		}
	}
	return polylineFromPoints(pts)
}

func polylineFromPoints(pts []Vec2) Polyline {
	pl := Polyline{Points: pts, CumLen: make([]float32, len(pts))}
	if len(pts) == 0 {
		return pl
	}
	for i := 1; i < len(pts); i++ {
		pl.CumLen[i] = pl.CumLen[i-1] + pts[i].Sub(pts[i-1]).Len()
	}
	pl.Length = pl.CumLen[len(pts)-1]
	return pl
}

// SampleAt returns position and unit tangent at distance along the polyline.
// Distance is clamped to [0, Length].
func (p Polyline) SampleAt(distance float32) (pos Vec2, tangent Vec2) {
	if len(p.Points) == 0 {
		return Vec2{}, Vec2{X: 1}
	}
	if len(p.Points) == 1 || p.Length <= 0 {
		return p.Points[0], Vec2{X: 1}
	}
	if distance <= 0 {
		t := p.Points[1].Sub(p.Points[0]).Normalize()
		if t == (Vec2{}) {
			t = Vec2{X: 1}
		}
		return p.Points[0], t
	}
	if distance >= p.Length {
		n := len(p.Points)
		t := p.Points[n-1].Sub(p.Points[n-2]).Normalize()
		if t == (Vec2{}) {
			t = Vec2{X: 1}
		}
		return p.Points[n-1], t
	}
	// Binary search CumLen for segment containing distance.
	lo, hi := 0, len(p.CumLen)-1
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if p.CumLen[mid] <= distance {
			lo = mid
		} else {
			hi = mid
		}
	}
	segLen := p.CumLen[hi] - p.CumLen[lo]
	var u float32
	if segLen > 1e-6 {
		u = (distance - p.CumLen[lo]) / segLen
	}
	a, b := p.Points[lo], p.Points[hi]
	pos = a.Add(b.Sub(a).Scale(u))
	tangent = b.Sub(a).Normalize()
	if tangent == (Vec2{}) {
		tangent = Vec2{X: 1}
	}
	return pos, tangent
}
