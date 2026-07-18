package game

import "github.com/ridespirals/this-city/internal/sim"

const (
	MachineDebug = "debug"

	stateAlpha sim.StateID = "alpha"
	stateBeta  sim.StateID = "beta"
)

// DebugMachine returns a timer-toggling FSM used to prove the engine in Phase 3.
func DebugMachine() *sim.Definition {
	return &sim.Definition{
		Initial: stateAlpha,
		States: map[sim.StateID]sim.StateDef{
			stateAlpha: {
				Hooks: sim.StateHooks{
					OnEnter: func(ctx *sim.Context) {
						ctx.BB.Timer = 0
						ctx.BB.Tag = "alpha"
					},
					OnUpdate: func(ctx *sim.Context) {
						ctx.BB.Timer += ctx.DT
					},
				},
				Transitions: []sim.Transition{{
					To: stateBeta,
					Guard: func(ctx *sim.Context) bool {
						return ctx.BB.Timer >= 1.5
					},
				}},
			},
			stateBeta: {
				Hooks: sim.StateHooks{
					OnEnter: func(ctx *sim.Context) {
						ctx.BB.Timer = 0
						ctx.BB.Tag = "beta"
					},
					OnUpdate: func(ctx *sim.Context) {
						ctx.BB.Timer += ctx.DT
					},
				},
				Transitions: []sim.Transition{{
					To: stateAlpha,
					Guard: func(ctx *sim.Context) bool {
						return ctx.BB.Timer >= 1.5
					},
				}},
			},
		},
	}
}

// SpawnDebugAgent creates a visible agent with the debug brain at (x, y).
func SpawnDebugAgent(w *sim.World, x, y float32) sim.Entity {
	if w == nil {
		return sim.NilEntity
	}
	e := w.Create()
	w.Transforms.Set(e, sim.Transform2D{X: x, Y: y, Scale: 1})
	w.Roles.Set(e, sim.RoleTag{Role: sim.RoleDebug})
	w.Brains.Set(e, sim.AgentBrain{
		Machine: MachineDebug,
		State:   stateAlpha,
		BB:      sim.Blackboard{Tag: "alpha"},
	})
	return e
}
