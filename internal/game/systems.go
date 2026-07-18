package game

import "github.com/ridespirals/this-city/internal/sim"

// Machines maps AgentBrain.Machine keys to FSM definitions.
type Machines map[string]*sim.Definition

// TickBrains runs OnUpdate/transitions for every AgentBrain.
func TickBrains(w *sim.World, machines Machines, dt float32) {
	if w == nil || machines == nil {
		return
	}
	type pair struct {
		e     sim.Entity
		brain sim.AgentBrain
	}
	var batch []pair
	w.Brains.ForEach(func(e sim.Entity, brain sim.AgentBrain) {
		batch = append(batch, pair{e, brain})
	})
	for _, p := range batch {
		def := machines[p.brain.Machine]
		if def == nil {
			continue
		}
		brain := p.brain
		def.Tick(&brain, dt, nil)
		w.Brains.Set(p.e, brain)
	}
}
