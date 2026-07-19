package game

import "github.com/ridespirals/this-city/internal/sim"

const (
	MachineWalk = "walk"
	StateWalk   = sim.StateID("walk")
)

// WalkMachine is the default agent brain: stay in Walk until richer behaviors exist.
// Future machines: patrol, travel (A→B), wander, flee — see plan/agents.md.
func WalkMachine() *sim.Definition {
	return &sim.Definition{
		Initial: StateWalk,
		States: map[sim.StateID]sim.StateDef{
			StateWalk: {
				Hooks: sim.StateHooks{
					OnEnter: func(ctx *sim.Context) {
						ctx.BB.Tag = "walk"
					},
				},
			},
		},
	}
}

// AttachWalkBrain gives e the default walk state machine.
func AttachWalkBrain(w *sim.World, e sim.Entity) {
	if w == nil || !w.Alive(e) {
		return
	}
	w.Brains.Set(e, sim.AgentBrain{
		Machine: MachineWalk,
		State:   StateWalk,
		BB:      sim.Blackboard{Tag: "walk"},
	})
}
