package sim

// DefaultPathSamples is the per-segment subdivision count for runtime polylines.
const DefaultPathSamples = 16

// AnchorsToSegments builds cubic segments between consecutive anchors with
// straight-line control points at 1/3 and 2/3.
func AnchorsToSegments(anchors []Vec2) []CubicBezier {
	if len(anchors) < 2 {
		return nil
	}
	out := make([]CubicBezier, 0, len(anchors)-1)
	for i := 0; i < len(anchors)-1; i++ {
		a, b := anchors[i], anchors[i+1]
		d := b.Sub(a)
		out = append(out, CubicBezier{
			P0: a,
			C0: a.Add(d.Scale(1.0 / 3)),
			C1: a.Add(d.Scale(2.0 / 3)),
			P1: b,
		})
	}
	return out
}
