package sim

// StateID names a state within a machine definition.
type StateID string

// Blackboard is per-agent scratch data for guards and actions.
type Blackboard struct {
	Timer  float32
	Target Entity
	Tag    string
	Flags  uint32
}

// Context is the pure view passed to FSM hooks. Guards must not mutate;
// actions / OnEnter / OnExit / OnUpdate may mutate the blackboard.
type Context struct {
	BB   *Blackboard
	DT   float32
	View WorldView
}

// WorldView is an optional narrow facade for guards/actions. Nil is allowed.
// Game can supply a richer implementation later (nearby queries, transforms).
type WorldView interface{}

// Guard is a side-effect-free predicate. Nil means always true.
type Guard func(ctx *Context) bool

// Action runs on a taken transition (may mutate the blackboard).
type Action func(ctx *Context)

// Transition is an ordered edge from the owning state.
type Transition struct {
	To     StateID
	Guard  Guard
	Action Action
}

// StateHooks are optional lifecycle callbacks for a state.
type StateHooks struct {
	OnEnter  func(ctx *Context)
	OnExit   func(ctx *Context)
	OnUpdate func(ctx *Context)
}

// StateDef describes one state: hooks plus outbound transitions (priority order).
type StateDef struct {
	Hooks       StateHooks
	Transitions []Transition
}

// Definition is a static FSM table. Instances live on AgentBrain.
type Definition struct {
	Initial StateID
	States  map[StateID]StateDef
}

// Tick advances one brain instance for dt seconds using def.
// Order: OnUpdate(current) → first matching transition → OnExit → Action → OnEnter.
func (def *Definition) Tick(brain *AgentBrain, dt float32, view WorldView) {
	if def == nil || brain == nil {
		return
	}
	if brain.State == "" {
		brain.State = def.Initial
		enter(def, brain, view, 0)
	}
	ctx := &Context{BB: &brain.BB, DT: dt, View: view}
	cur, ok := def.States[brain.State]
	if !ok {
		return
	}
	if cur.Hooks.OnUpdate != nil {
		cur.Hooks.OnUpdate(ctx)
	}
	for _, tr := range cur.Transitions {
		if tr.Guard != nil && !tr.Guard(ctx) {
			continue
		}
		if cur.Hooks.OnExit != nil {
			cur.Hooks.OnExit(ctx)
		}
		if tr.Action != nil {
			tr.Action(ctx)
		}
		brain.State = tr.To
		enter(def, brain, view, dt)
		return
	}
}

func enter(def *Definition, brain *AgentBrain, view WorldView, dt float32) {
	next, ok := def.States[brain.State]
	if !ok || next.Hooks.OnEnter == nil {
		return
	}
	ctx := &Context{BB: &brain.BB, DT: dt, View: view}
	next.Hooks.OnEnter(ctx)
}
