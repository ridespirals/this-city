# Entity Component System (ECS)

## Why ECS

Agents share cross-cutting data (transform, path follower, perception) with different behaviors. Composition via components avoids deep OOP hierarchies and keeps systems testable in isolation.

## Core concepts

| Concept | Meaning |
|---------|---------|
| **Entity** | Opaque ID (index + generation to reuse slots safely) |
| **Component** | Plain data struct; no behavior methods that touch the world |
| **System** | Function that queries entities with a component set and updates them for `dt` |

## Storage (v1)

Implemented in `internal/sim`:

- Entity allocator: free-list + generation counters (generations start at 1; `{0,0}` is `NilEntity`).
- Components: `ComponentStore[T]` (`map[Entity]T`) on `World` — currently `Transforms`, `Brains`, `Roles`.
- Destroy clears all registered stores and bumps generation.
- No archetype graph until we need cache-friendly iteration at scale.

## Early components

| Component | Purpose |
|-----------|---------|
| `Transform2D` | Position, rotation (radians), optional scale |
| `Velocity` | Linear velocity for free movement / steering |
| `PathFollower` | Network edge id, distance, speed, forward |
| `PathDecision` | Random vs routed junction choice (`Route` for A*) |
| `AgentBrain` | FSM id/type, current state, blackboard ref |
| `Perception` | Radius, last sensed entity/event ids |
| `Role` | `Civilian`, `Police`, … |
| `Schedule` | Shift windows, break timers, destination intents |
| `EventSource` | Event kind, priority, lifetime, active flag |
| `Selectable` | Editor selection / pick id |

Add components when a system needs new data; avoid “god” components.

## Systems (ownership)

Systems live primarily in `internal/game`. Shared mechanical systems (integrate velocity, advance path follower) may live in `sim` if they are role-agnostic.

Suggested tick order (adjust as needed):

1. Apply pending commands (spawn/despawn, edit paths).
2. Schedule / timer updates.
3. Perception (BSP queries → fill `Perception`).
4. FSM / brain transitions.
5. Path follow / movement integration.
6. Event lifetime / cleanup.
7. Spatial index refresh if transforms moved (or deferred rebuild).

## Queries

- Explicit “entities with components A,B not C” helpers on the world.
- Systems must not assume iteration order unless documented (tests should not rely on map order).

## Destruction

- Destroying an entity removes all components and bumps generation.
- Held entity IDs in blackboards must be validated (generation check) before use.

## Relationship to OOP

- No class hierarchy for agent types.
- Role differences are data (`Role`) + which FSM definition is attached + which systems care.
- Shared behavior = shared systems/components.

See also: [state-machines.md](state-machines.md), [agents.md](agents.md), [architecture.md](architecture.md).
