# AGENTS.md — Guidance for AI coding agents

This file is the contract for automated assistants working on **This City**. Prefer it over inventing structure. Keep it updated when architecture, workflow, or conventions change.

## Project in one paragraph

Go + raylib city sim. ECS + FSM-driven agents (civilians, police), Bezier street network, path following/pathfinding, BSP spatial queries. Layers: pure `sim` → `game` → `render`/`editor`. Docs-first; code scaffolding starts at roadmap Phase 2.

## Read first

1. [README.md](README.md) — human status and scope  
2. [plan/README.md](plan/README.md) — plan index  
3. [plan/architecture.md](plan/architecture.md) — dependency rules  
4. Relevant `plan/*.md` for the feature you touch  
5. [plan/roadmap.md](plan/roadmap.md) — what phase we are in  

## Hard rules

### Layering

- `internal/sim`: **no** raylib, **no** `game`/`render`/`editor` imports.
- `internal/game`: may import `sim` only (among project packages); **no** raylib.
- `internal/render` / `internal/editor`: raylib and UI; mutate world via **commands**, not ad-hoc draws-side writes.
- Prefer ECS composition over OOP agent hierarchies.

### Documentation (every prompt)

When a change affects design, scope, status, or conventions, update as needed:

| Doc | Audience | Update when |
|-----|----------|-------------|
| `plan/<feature>.md` | Both | Feature/system design changes |
| `plan/README.md` / `plan/roadmap.md` | Both | Phase/status/order changes |
| `README.md` | Humans | Vision, status, how to run, feature summary |
| `AGENTS.md` | Agents | Workflow, layer rules, git policy, conventions |

Do not leave docs contradicting the code or the agreed plan.

### Git

- **Do not** run git writes: no commit, push, tag, branch create/checkout, rebase, amend, or config changes unless the human explicitly asks in that conversation.
- **Do** suggest: commit message summaries, when to branch, when to cut releases/tags.
- After code changes, always offer a concise **suggested commit message** (subject + short body focusing on why).

### Code quality

- Match existing style once code exists; no drive-by refactors.
- Tests for `sim` and `game` must not require a window.
- Deterministic ticks: fixed `dt`, seeded RNG in tests.
- No third-party ECS library unless the human decides to adopt one; hand-roll minimal first.

## Planned package layout

```
cmd/this-city/main.go
internal/sim/
internal/game/
internal/render/
internal/editor/
plan/                 # design docs (source of truth for intent)
```

## Current phase

**Phase 1 — Docs foundation.** No application module yet. Do not invent large code trees unless the human asks to start Phase 2.

## Suggested commit style

```
Short imperative summary (what/why).

Optional body: motivation, layer touched, doc updates.
```

Example after this docs bootstrap:

```
Document project architecture and phased roadmap.

Establish plan/, README, and AGENTS.md so sim, game, and render
layers have a shared contract before code scaffolding.
```

## Release suggestions (human-owned)

- `v0.1.0` — after Phase 2 (runnable window).
- Further tags at roadmap phase boundaries; see [plan/roadmap.md](plan/roadmap.md).
