# Events and Placeable Objects

## Purpose

Author-time and runtime objects that agents can perceive and respond to—crimes, distress calls, attractions, furniture anchors (benches), etc.

## Model

Events/objects are entities with:

| Component / field | Role |
|-------------------|------|
| `Transform2D` | Where it is |
| `EventSource` | Kind, priority, lifetime, active |
| Optional tags | `POI`, `Cover`, … |

Kinds (v1 set; extend as needed):

- `Crime` — police investigate; civilians may flee/watch.
- `Distress` / `HelpRequest` — police respond; civilians may watch.
- `Attraction` — civilians approach/hang out.
- `Bench` / `POI` — hang-out / break anchors (may be “objects” more than timed events).

Timed events expire via a lifetime system; POIs persist until deleted in the editor.

## Lifecycle

1. Editor or script **spawns** event entity.
2. Perception systems expose it to nearby agents.
3. FSM guards match kind + role → transition.
4. Resolution: police investigate completes → deactivate/despawn; or lifetime elapses; or editor deletes.

## Priority

If multiple events are in range, pick by: role relevance → higher priority → closer distance. Store chosen target on the blackboard.

## Commands

Prefer spawning through game commands (`SpawnEvent`, `DespawnEntity`) so undo/editor and tests share one path.

## Testing

- Spawn crime near police; after N ticks, police state is `Investigate` and event eventually clears.
- Expired events disappear from perception.

See also: [agents.md](agents.md), [editor.md](editor.md), [spatial.md](spatial.md).
