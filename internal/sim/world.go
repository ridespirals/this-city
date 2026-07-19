// Package sim is the pure simulation core: ECS, math, paths, spatial indexes, and FSM.
// It must not import raylib or UI packages.
package sim

type remover interface {
	removeEntity(Entity)
}

// World holds ECS entity slots and component stores.
type World struct {
	generations []uint32
	free        []uint32
	removers    []remover

	Paths *PathSet

	Transforms *ComponentStore[Transform2D]
	Brains     *ComponentStore[AgentBrain]
	Roles      *ComponentStore[RoleTag]
	Followers  *ComponentStore[PathFollower]
	Events     *ComponentStore[EventSource]
}

// NewWorld returns an empty simulation world with standard component stores.
func NewWorld() *World {
	w := &World{
		Paths:      newPathSet(),
		Transforms: newStore[Transform2D](),
		Brains:     newStore[AgentBrain](),
		Roles:      newStore[RoleTag](),
		Followers:  newStore[PathFollower](),
		Events:     newStore[EventSource](),
	}
	w.removers = []remover{w.Transforms, w.Brains, w.Roles, w.Followers, w.Events}
	return w
}

// Create allocates a new live entity.
func (w *World) Create() Entity {
	if w == nil {
		return NilEntity
	}
	var idx uint32
	if n := len(w.free); n > 0 {
		idx = w.free[n-1]
		w.free = w.free[:n-1]
	} else {
		idx = uint32(len(w.generations))
		// Generations start at 1 so Entity{0,0} (NilEntity) is never alive.
		w.generations = append(w.generations, 1)
	}
	return Entity{Index: idx, Generation: w.generations[idx]}
}

// Alive reports whether e refers to a currently live entity.
func (w *World) Alive(e Entity) bool {
	if w == nil || e.IsNil() {
		return false
	}
	if int(e.Index) >= len(w.generations) {
		return false
	}
	return w.generations[e.Index] == e.Generation
}

// Destroy removes all components for e and invalidates the ID. Returns false if e was not alive.
func (w *World) Destroy(e Entity) bool {
	if !w.Alive(e) {
		return false
	}
	for _, r := range w.removers {
		r.removeEntity(e)
	}
	w.generations[e.Index]++
	w.free = append(w.free, e.Index)
	return true
}

// Len returns the number of live entities (allocated slots minus free list).
func (w *World) Len() int {
	if w == nil {
		return 0
	}
	return len(w.generations) - len(w.free)
}
