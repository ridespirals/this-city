package sim

// PathID identifies a path in a PathSet.
type PathID uint32

// NilPath is an unset path id.
const NilPath PathID = 0

// DefaultPathSamples is the per-segment subdivision count for runtime polylines.
const DefaultPathSamples = 16

// Path is an ordered chain of cubic Bézier segments plus a sampled polyline.
type Path struct {
	ID       PathID
	Segments []CubicBezier
	Poly     Polyline
}

// Rebuild refreshes the polyline approximation from Segments.
func (p *Path) Rebuild(stepsPerSeg int) {
	if p == nil {
		return
	}
	if stepsPerSeg < 1 {
		stepsPerSeg = DefaultPathSamples
	}
	p.Poly = BuildPolyline(p.Segments, stepsPerSeg)
}

// PathSet owns authored paths (not ECS entities).
type PathSet struct {
	next  PathID
	paths map[PathID]*Path
}

func newPathSet() *PathSet {
	return &PathSet{
		next:  1,
		paths: make(map[PathID]*Path),
	}
}

// Add creates a path from segments, rebuilds its polyline, and returns its id.
func (ps *PathSet) Add(segments []CubicBezier) PathID {
	if ps == nil {
		return NilPath
	}
	id := ps.next
	ps.next++
	p := &Path{ID: id, Segments: append([]CubicBezier(nil), segments...)}
	p.Rebuild(DefaultPathSamples)
	ps.paths[id] = p
	return id
}

// SetSegments replaces a path's geometry and rebuilds its polyline.
func (ps *PathSet) SetSegments(id PathID, segments []CubicBezier) bool {
	p, ok := ps.Get(id)
	if !ok {
		return false
	}
	p.Segments = append([]CubicBezier(nil), segments...)
	p.Rebuild(DefaultPathSamples)
	return true
}

// Remove deletes a path. Followers referencing it are not auto-cleared.
func (ps *PathSet) Remove(id PathID) bool {
	if ps == nil || id == NilPath {
		return false
	}
	if _, ok := ps.paths[id]; !ok {
		return false
	}
	delete(ps.paths, id)
	return true
}

// Get returns a path by id.
func (ps *PathSet) Get(id PathID) (*Path, bool) {
	if ps == nil || id == NilPath {
		return nil, false
	}
	p, ok := ps.paths[id]
	return p, ok
}

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

// ForEach invokes fn for each path. Order is unspecified.
func (ps *PathSet) ForEach(fn func(p *Path)) {
	if ps == nil {
		return
	}
	for _, p := range ps.paths {
		fn(p)
	}
}

// Len returns the number of paths.
func (ps *PathSet) Len() int {
	if ps == nil {
		return 0
	}
	return len(ps.paths)
}

// PathFollower advances an entity along a path at constant speed.
type PathFollower struct {
	Path     PathID
	Distance float32
	Speed    float32 // world units per second
	Forward  bool    // false = travel toward start
	PingPong bool    // reverse direction at ends (vs clamp)
}
