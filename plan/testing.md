# Testing

## Principles

- **`sim` and `game` are the test fortress** — pure Go tests, no window.
- Deterministic: fixed `dt`, seeded RNG, explicit world fixtures.
- **`render` / `editor` stay thin** — manual smoke; virtual input tests where cheap.

## What to test by layer

| Layer | Examples |
|-------|----------|
| `sim` ECS | create/destroy, generation reuse, component get/set |
| `sim` FSM | transition tables, guard/action order, OnEnter/Exit |
| `sim` paths | Bezier sample, arc-length advance, end-of-segment |
| `sim` BSP | radius query correctness, empty/degenerate |
| `sim` pathfinding | A\* on fixtures, unreachable |
| `game` agents | civilian hang-out; police investigate on crime |
| `game` events | lifetime expiry; priority pick |
| `editor` | tool → command mapping with fake input |
| `render` | compile/run smoke only |

## Patterns

- **World fixtures:** helpers that spawn a tiny city (two nodes, one edge, one bench).
- **Tick loops:** `for i := 0; i < N; i++ { game.Tick(1.0/60) }` then assert state.
- **Fake perception:** either real BSP with placed entities or inject blackboard targets for FSM-only tests.

## CI (when code exists)

- `go test ./...` on every PR.
- No GPU requirement for unit tests.
- Race detector optional on sim/game packages.

## Non-goals

- Full visual regression suite in v1.
- Flaky wall-clock timing tests.

See also: [architecture.md](architecture.md), [roadmap.md](roadmap.md).
