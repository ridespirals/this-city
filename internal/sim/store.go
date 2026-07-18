package sim

// ComponentStore is a simple map-backed component table keyed by Entity.
type ComponentStore[T any] struct {
	data map[Entity]T
}

func newStore[T any]() *ComponentStore[T] {
	return &ComponentStore[T]{data: make(map[Entity]T)}
}

func (s *ComponentStore[T]) removeEntity(e Entity) {
	if s == nil {
		return
	}
	delete(s.data, e)
}

// Has reports whether e has this component.
func (s *ComponentStore[T]) Has(e Entity) bool {
	if s == nil {
		return false
	}
	_, ok := s.data[e]
	return ok
}

// Get returns the component and whether it exists.
func (s *ComponentStore[T]) Get(e Entity) (T, bool) {
	var zero T
	if s == nil {
		return zero, false
	}
	v, ok := s.data[e]
	return v, ok
}

// Set writes the component for e.
func (s *ComponentStore[T]) Set(e Entity, v T) {
	if s == nil {
		return
	}
	s.data[e] = v
}

// Remove deletes the component for e.
func (s *ComponentStore[T]) Remove(e Entity) {
	s.removeEntity(e)
}

// Len returns the number of component instances.
func (s *ComponentStore[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

// ForEach invokes fn for each entity/component pair. Iteration order is unspecified.
func (s *ComponentStore[T]) ForEach(fn func(e Entity, v T)) {
	if s == nil {
		return
	}
	for e, v := range s.data {
		fn(e, v)
	}
}
