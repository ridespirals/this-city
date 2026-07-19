package game

import "github.com/ridespirals/this-city/internal/sim"

// SpawnAgent places a role-tagged agent at pos (no FSM yet — Phase 6).
func SpawnAgent(w *sim.World, role sim.Role, pos sim.Vec2) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Transforms.Set(e, sim.Transform2D{X: pos.X, Y: pos.Y, Scale: 1})
	w.Roles.Set(e, sim.RoleTag{Role: role})
	return e
}

// SpawnEvent places an active event/POI at pos.
func SpawnEvent(w *sim.World, kind sim.EventKind, pos sim.Vec2) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Transforms.Set(e, sim.Transform2D{X: pos.X, Y: pos.Y, Scale: 1})
	lifetime := float32(0)
	priority := 1
	switch kind {
	case sim.EventCrime:
		lifetime = 60
		priority = 10
	case sim.EventDistress:
		lifetime = 45
		priority = 8
	case sim.EventAttraction:
		lifetime = 30
		priority = 3
	case sim.EventBench:
		lifetime = 0
		priority = 1
	}
	w.Events.Set(e, sim.EventSource{
		Kind:     kind,
		Priority: priority,
		Lifetime: lifetime,
		Active:   true,
	})
	return e
}

// DeleteEntity destroys an entity if it is alive.
func DeleteEntity(w *sim.World, e sim.Entity) bool {
	if w == nil {
		return false
	}
	return w.Destroy(e)
}

// SetPathFromAnchors creates or updates a path from anchor points.
// If id is NilPath, a new path is created. Returns the path id.
func SetPathFromAnchors(w *sim.World, id sim.PathID, anchors []sim.Vec2) sim.PathID {
	if w == nil {
		return sim.NilPath
	}
	segs := sim.AnchorsToSegments(anchors)
	if len(segs) == 0 {
		return id
	}
	if id == sim.NilPath {
		return w.Paths.Add(segs)
	}
	if !w.Paths.SetSegments(id, segs) {
		return w.Paths.Add(segs)
	}
	return id
}

// DeletePath removes a path from the world.
func DeletePath(w *sim.World, id sim.PathID) bool {
	if w == nil {
		return false
	}
	return w.Paths.Remove(id)
}

// PickEntity returns the nearest entity with a transform within maxDist, or NilEntity.
func PickEntity(w *sim.World, pos sim.Vec2, maxDist float32) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	best := sim.NilEntity
	bestD := maxDist * maxDist
	w.Transforms.ForEach(func(e sim.Entity, xf sim.Transform2D) {
		dx := xf.X - pos.X
		dy := xf.Y - pos.Y
		d := dx*dx + dy*dy
		if d <= bestD {
			bestD = d
			best = e
		}
	})
	return best
}

// PickPath returns the nearest path id within maxDist of pos (distance to polyline samples).
func PickPath(w *sim.World, pos sim.Vec2, maxDist float32) sim.PathID {
	if w == nil || w.Paths == nil {
		return sim.NilPath
	}
	best := sim.NilPath
	bestD := maxDist * maxDist
	w.Paths.ForEach(func(p *sim.Path) {
		for _, pt := range p.Poly.Points {
			dx := pt.X - pos.X
			dy := pt.Y - pos.Y
			d := dx*dx + dy*dy
			if d <= bestD {
				bestD = d
				best = p.ID
			}
		}
	})
	return best
}

// MoveEntity sets an entity's transform position (keeps rotation/scale).
func MoveEntity(w *sim.World, e sim.Entity, pos sim.Vec2) bool {
	if w == nil || !w.Alive(e) {
		return false
	}
	xf, ok := w.Transforms.Get(e)
	if !ok {
		return false
	}
	xf.X, xf.Y = pos.X, pos.Y
	w.Transforms.Set(e, xf)
	return true
}
