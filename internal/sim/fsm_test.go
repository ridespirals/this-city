package sim

import "testing"

func TestFSMTimerToggle(t *testing.T) {
	const (
		alpha StateID = "alpha"
		beta  StateID = "beta"
	)
	def := &Definition{
		Initial: alpha,
		States: map[StateID]StateDef{
			alpha: {
				Hooks: StateHooks{
					OnEnter: func(ctx *Context) {
						ctx.BB.Timer = 0
						ctx.BB.Tag = "alpha"
					},
					OnUpdate: func(ctx *Context) {
						ctx.BB.Timer += ctx.DT
					},
				},
				Transitions: []Transition{{
					To: beta,
					Guard: func(ctx *Context) bool {
						return ctx.BB.Timer >= 1
					},
				}},
			},
			beta: {
				Hooks: StateHooks{
					OnEnter: func(ctx *Context) {
						ctx.BB.Timer = 0
						ctx.BB.Tag = "beta"
					},
					OnUpdate: func(ctx *Context) {
						ctx.BB.Timer += ctx.DT
					},
				},
				Transitions: []Transition{{
					To: alpha,
					Guard: func(ctx *Context) bool {
						return ctx.BB.Timer >= 1
					},
				}},
			},
		},
	}

	brain := &AgentBrain{}
	def.Tick(brain, 0, nil)
	if brain.State != alpha || brain.BB.Tag != "alpha" {
		t.Fatalf("initial: state=%q tag=%q", brain.State, brain.BB.Tag)
	}

	def.Tick(brain, 0.6, nil)
	if brain.State != alpha {
		t.Fatalf("after 0.6s: state=%q", brain.State)
	}
	def.Tick(brain, 0.6, nil)
	if brain.State != beta || brain.BB.Tag != "beta" {
		t.Fatalf("after 1.2s: state=%q tag=%q", brain.State, brain.BB.Tag)
	}

	def.Tick(brain, 1.0, nil)
	if brain.State != alpha {
		t.Fatalf("after another 1s: state=%q", brain.State)
	}
}

func TestFSMTransitionAction(t *testing.T) {
	const (
		a StateID = "a"
		b StateID = "b"
	)
	var actionFired bool
	def := &Definition{
		Initial: a,
		States: map[StateID]StateDef{
			a: {
				Transitions: []Transition{{
					To: b,
					Action: func(ctx *Context) {
						actionFired = true
						ctx.BB.Flags = 7
					},
				}},
			},
			b: {},
		},
	}
	brain := &AgentBrain{State: a}
	def.Tick(brain, 0.016, nil)
	if !actionFired || brain.State != b || brain.BB.Flags != 7 {
		t.Fatalf("actionFired=%v state=%q flags=%d", actionFired, brain.State, brain.BB.Flags)
	}
}
