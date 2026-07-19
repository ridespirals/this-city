package game

import "github.com/ridespirals/this-city/internal/sim"

// SpawnAgent places a role-tagged agent with a Walk brain at pos.
func SpawnAgent(w *sim.World, role sim.Role, pos sim.Vec2) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Transforms.Set(e, sim.Transform2D{X: pos.X, Y: pos.Y, Scale: 1})
	w.Roles.Set(e, sim.RoleTag{Role: role})
	AttachWalkBrain(w, e)
	w.Decisions.Set(e, sim.DefaultPathDecision())
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

// SetPathFromAnchors creates or updates an editor edge-group from anchor points.
// group 0 means allocate a new group. Returns the group id.
func SetPathFromAnchors(w *sim.World, group uint32, anchors []sim.Vec2) uint32 {
	if w == nil || w.Network == nil {
		return 0
	}
	g, _ := w.Network.AnchorsToChain(anchors, group)
	return g
}

// DeletePathGroup removes an editor-authored edge group.
func DeletePathGroup(w *sim.World, group uint32) {
	if w == nil || w.Network == nil {
		return
	}
	w.Network.RemoveGroup(group)
}

// DeleteEdge removes a single network edge.
func DeleteEdge(w *sim.World, id sim.EdgeID) bool {
	if w == nil || w.Network == nil {
		return false
	}
	return w.Network.RemoveEdge(id)
}

// StampPathPiece places a reusable path piece at world position `at`.
// Returns the editor group id (0 if nothing was stamped).
func StampPathPiece(w *sim.World, piece sim.PathPiece, at sim.Vec2, mergeEps float32) uint32 {
	if w == nil || w.Network == nil || len(piece.Curves) == 0 {
		return 0
	}
	g, _ := sim.StampPiece(w.Network, piece, at, mergeEps)
	return g
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

// PickEdge returns the nearest network edge within maxDist.
func PickEdge(w *sim.World, pos sim.Vec2, maxDist float32) sim.EdgeID {
	if w == nil || w.Network == nil {
		return sim.NilEdge
	}
	id, _, ok := w.Network.NearestEdge(pos, maxDist)
	if !ok {
		return sim.NilEdge
	}
	return id
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
