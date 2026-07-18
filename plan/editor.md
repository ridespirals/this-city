# Editor (2D Overhead Tools)

## Purpose

Author city geometry and populate the sim: place agents, events/objects, and draw Bezier paths from an overhead 2D view.

## View

- Orthographic top-down camera; pan/zoom.
- Optional grid and snap.
- Sim ticks while editing; **pause toggle** freezes game systems (editor still interactive).

## Toolbar tools (v1)

| Tool | Action |
|------|--------|
| Select | Click pick (BSP/point query); move/delete selection |
| Place civilian | Spawn civilian at cursor |
| Place police | Spawn police at cursor |
| Place event | Sub-menu or cycle: crime, distress, attraction, bench |
| Draw path | Click to place Bezier anchors / control points |
| Edit path | Adjust control points of selected segment |

Keep the toolbar visually simple: one active tool, clear cursor affordance.

## Commands

All mutations go through command APIs (not ad-hoc component writes from draw code):

- `SpawnAgent(role, pos)`
- `SpawnEvent(kind, pos)`
- `DeleteEntity(id)`
- `PathAddPoint` / `PathUpdateControl` / `PathCommitSegment`
- `SetPaused(bool)`

## Feedback

- Ghost preview for placement.
- Selected entity outline.
- Path handles visible only with path tools.

## Non-goals (v1)

- Full undo stack (nice-to-have soon after).
- Multi-select transform gizmos.
- Terrain painting / building footprints (later content).

## Testing

- Logic of tools (which command fires) can be tested without raylib by feeding virtual clicks.
- Manual smoke: place agents, draw a path, press play.

See also: [rendering.md](rendering.md), [paths.md](paths.md), [events.md](events.md).
