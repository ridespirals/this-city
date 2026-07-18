// Package sim is the pure simulation core: ECS, math, paths, spatial indexes, and FSM.
// It must not import raylib or UI packages.
package sim

// World holds simulation state. Phase 2 is a placeholder; ECS arrives in Phase 3.
type World struct{}

// NewWorld returns an empty simulation world.
func NewWorld() *World {
	return &World{}
}
