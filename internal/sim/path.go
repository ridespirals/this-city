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

// Get returns a path by id.
func (ps *PathSet) Get(id PathID) (*Path, bool) {
	if ps == nil || id == NilPath {
		return nil, false
	}
	p, ok := ps.paths[id]
	return p, ok
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
