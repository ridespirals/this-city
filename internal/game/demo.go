package game

import (
	"os"
	"path/filepath"

	"github.com/ridespirals/this-city/internal/sim"
)

// LoadDemoMap loads maps/figure-8.json when present, else the built-in figure-8 map.
func LoadDemoMap(w *sim.World) error {
	if w == nil {
		return nil
	}
	candidates := []string{
		"maps/figure-8.json",
		filepath.Join("..", "..", "maps", "figure-8.json"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return sim.LoadNetworkFile(p, w.Network)
		}
	}
	return sim.ApplyMapFile(sim.FigureEightMap(), w.Network)
}

// SpawnPathFollower creates a walking agent on edgeID.
func SpawnPathFollower(w *sim.World, edgeID sim.EdgeID, speed float32) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Roles.Set(e, sim.RoleTag{Role: sim.RoleDebug})
	AttachWalkBrain(w, e)
	w.Decisions.Set(e, sim.DefaultPathDecision())
	w.Transforms.Set(e, sim.Transform2D{Scale: 1})
	sim.PlaceOnEdge(w, e, edgeID, 0, true, speed)
	return e
}

// SpawnWalkerOnNearestEdge places a role agent on the nearest network edge to pos.
func SpawnWalkerOnNearestEdge(w *sim.World, role sim.Role, pos sim.Vec2, speed float32) sim.Entity {
	e := SpawnAgent(w, role, pos)
	if e.IsNil() || w.Network == nil {
		return e
	}
	edge, dist, ok := w.Network.NearestEdge(pos, 1e9)
	if !ok {
		return e
	}
	w.Decisions.Set(e, sim.DefaultPathDecision())
	sim.PlaceOnEdge(w, e, edge, dist, true, speed)
	return e
}
