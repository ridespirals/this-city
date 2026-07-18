package sim

// Entity is an opaque ID. Index selects a slot; Generation detects stale references.
type Entity struct {
	Index      uint32
	Generation uint32
}

// NilEntity is the zero entity (never alive).
var NilEntity Entity

// IsNil reports whether e is the zero value.
func (e Entity) IsNil() bool {
	return e == NilEntity
}
