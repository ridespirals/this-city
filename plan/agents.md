# Agents

## Overview

Agents are ECS entities with `Role`, `AgentBrain` (FSM), `PathFollower`, `PathDecision`, and (later) `Perception`. Behavior differences come from FSM definitions and which systems run—not subclasses.

## Default brain: Walk (now)

Every spawned agent gets `MachineWalk` with a single `walk` state. Movement is continuous along the network; **`PathDecision`** (default: random) picks edges at junctions. No destination yet.

```
walk  →  (future transitions into travel / wander / flee / patrol)
```

## Behavior sets (planned)

| Behavior | Intent | PathDecision |
|----------|--------|--------------|
| `walk` | Keep moving on the network | Random at intersections |
| `travel` | Go from A to B | A*/Dijkstra → `DecideRoute` |
| `patrol` | Loop waypoints / beat | Fixed or cyclic route |
| `wander` | Walk, stop at benches, watch events | Random + perception interrupts |
| `flee` | Escape a threat entity/category | Prefer increasing separation |

Role FSMs (civilian / police) will compose these behaviors rather than replace the network follower.

## Civilian

**Intent:** Move between locations, sometimes hang out, notice and react to nearby events/objects.

### States (v1 target)

| State | Behavior |
|-------|----------|
| `Idle` | Pick next destination on a timer or schedule |
| `Travel` | Pathfind + follow path to destination |
| `HangOut` | Stay at POI (e.g. park bench); play idle timer |
| `React` | Respond to perceived event: approach, watch, or flee |

### Reactions

Driven by event kind + simple personality weights on the blackboard (optional in v1: fixed per reaction table).

- Attraction / street performance → approach or watch.
- Crime / distress → watch from distance or flee (civilians do not “investigate” like police).

### Schedule (light)

Optional destinations list or time-of-day slots later; v1 can use random POI picks from a set of anchors.

## Police

**Intent:** Work a shift, patrol, take breaks, detect crimes/people needing help, investigate.

### States (v1 target)

| State | Behavior |
|-------|----------|
| `StartShift` | Enter duty; set patrol route / shift end |
| `Patrol` | Follow patrol waypoints or loop edges |
| `OnBreak` | Pause patrol at break spot for a duration |
| `Investigate` | Move to event; remain until resolved, expired, or lost |

### Duty cycle

- Shift start → patrol.
- Periodic break transitions from `Patrol` (timer / schedule).
- On perceive `Crime` or `Distress` (and maybe `HelpRequest`): transition to `Investigate`.
- After investigate: resume `Patrol` or `OnBreak` if break was deferred.

### Authority

v1: investigation is presence + timer that “clears” or acknowledges the event. Arrest / chase combat is out of scope initially.

## Perception

- Radius query via BSP ([spatial.md](spatial.md)).
- Filter by event/agent relevance for the role.
- LOS stub: none at first (all in radius visible).

## Spawning

Editor toolbar places civilians/police with default components and the correct FSM definition ([editor.md](editor.md)).

## Testing

- Script perception + time; assert civilian enters `HangOut` at bench.
- Assert police enter `Investigate` when a `Crime` is within radius during `Patrol`.

See also: [state-machines.md](state-machines.md), [events.md](events.md), [ecs.md](ecs.md).
