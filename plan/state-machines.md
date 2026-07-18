# State Machine Engine

## Purpose

Drive agent behavior with an explicit, testable FSM living in `internal/sim`, configured and interpreted by `internal/game`.

## Design

### Definition (static)

- **States:** named ids (`Idle`, `Travel`, `HangOut`, `Patrol`, `OnBreak`, `Investigate`, …).
- **Transitions:** `from` → `to` with:
  - **Guard:** pure predicate over blackboard + world view + `dt`/time.
  - **Action (optional):** run on transition (set target, clear investigation, start timer).
- **OnEnter / OnExit / OnUpdate:** optional hooks per state.

Definitions are data (tables or builder APIs), not subclasses.

### Instance (per agent)

Stored in `AgentBrain` (or adjacent component):

- Current state id.
- **Blackboard:** small key-value or typed struct (target entity, path destination, shift end time, investigation target, mood/reaction).
- Optional timers local to the state.

### Tick

```
for each agent with AgentBrain:
  run OnUpdate(current)
  evaluate transitions in priority order
  if guard passes: OnExit → action → change state → OnEnter
```

Guards must be side-effect free. Side effects belong in actions / OnEnter / OnUpdate.

## World access

FSM hooks receive a narrow **world view** interface (query nearby, get transform, get event kind) implemented by `game` over `sim`. The FSM engine itself stays in `sim` and does not import role logic.

## Hierarchy (later)

v1 is flat FSMs. Nested FSMs (e.g. Investigate sub-states) can wrap a child machine later without changing the ECS shape.

## Testing

- Unit-test definitions with a fake world view and fixed clock.
- Assert transition sequences given scripted perception and timers.
- No raylib in FSM tests.

## Example sketches

**Civilian (simplified):**

`Idle` → `Travel` (has destination) → `HangOut` (arrived at bench) → `Idle`  
Any → `React` when perceived event matches interest → back to prior intent when done.

**Police (simplified):**

`StartShift` → `Patrol` → `OnBreak` (timer) → `Patrol`  
`Patrol`/`OnBreak` → `Investigate` (crime/distress perceived) → `Patrol` when resolved or lost.

Full behavior notes: [agents.md](agents.md).

See also: [ecs.md](ecs.md), [testing.md](testing.md).
