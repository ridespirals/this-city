# AGENTS.md — Guidance for AI coding agents

This file is the contract for automated assistants working on **This City**. Prefer it over inventing structure. Keep it updated when architecture, workflow, or conventions change.

## Project in one paragraph

Go + raylib city sim. ECS + FSM-driven agents (civilians, police), Bezier street network, path following/pathfinding, BSP spatial queries. Layers: pure `sim` → `game` → `render`/`editor`. Module path: `github.com/ridespirals/this-city`. Phase 2 scaffolding is in place; next is Phase 3 (ECS + FSM).

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

## Package layout

```
cmd/this-city/main.go     # window + loop wiring
internal/sim/             # pure sim (no raylib)
internal/game/            # tick/session (no raylib)
internal/render/          # raylib window + draw helpers
internal/editor/          # tool state (tools in Phase 5)
plan/                     # design docs (source of truth for intent)
```

## Current phase

**Phase 2 — Module + loop (complete).** Runnable `go run ./cmd/this-city`. Next: Phase 3 — ECS + FSM core (`plan/ecs.md`, `plan/state-machines.md`). Do not skip ahead to editor/agents unless asked.

## Suggested commit style

```
Short imperative summary (what/why).

Optional body: motivation, layer touched, doc updates.
```

## CI

- [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — on push to `main` and on PRs: `go test ./...` + build.
- [`.github/workflows/release.yml`](.github/workflows/release.yml) — on tags `v*`: multi-OS binaries + GitHub Release.
- Do not commit secrets; release uses `GITHUB_TOKEN`.
- When changing Go version or native deps, update both workflows (Linux apt packages for raylib).

## Release suggestions (human-owned)

- Tag **`v0.1.0`** for Phase 2 (`git tag v0.1.0 && git push origin v0.1.0`) once CI is green on `main`.
- Further tags at roadmap phase boundaries; see [plan/roadmap.md](plan/roadmap.md).
- Prefer annotated tags: `git tag -a v0.1.0 -m "v0.1.0"`.
