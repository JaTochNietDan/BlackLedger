# Current goal and workstream ownership

## User amendment — detail, interiors and event direction (2026-09-13)

The user requires substantially more detailed, high-quality 1950s noir gangster architecture and characters; current blocky assets do not meet visual acceptance. Model detailed building interiors for display on entry, preserving the room's actual public occupants and available gameplay actions. For committed simulation events, automatically open the city and frame the complete action before presenting its animation. Keep the Mac awake during active work. Validate visual quality and performance in the browser as production work proceeds; local high FPS alone does not establish production readiness.

## User revision — browser 3D city (2026-09-13)

The current user request supersedes the earlier 2D visual-production restriction below. Codex is to build one browser-rendered, rotatable, pannable, zoomable 3D 1950s mafia city, using locally authored Blender models and textures. Buildings must be selectable destinations; public simulation journeys must show pedestrians or their appropriate vehicles; committed violence must drive animated effects. Preserve Go authority and existing command handling. Use isolated saves for QA, verify footprints and browser behavior, measure performance toward 60 FPS, and record unfinished acceptance honestly. Historical handoff instructions remain below for context; no separate visual agent is currently running on this request.

## User revision — visual production handoff (2026-09-07)

The user has reassigned visual production to Claude or another separate agent. **Codex's ongoing goal is gameplay and player experience:** build and playtest the single-player mafia vertical slice, deepen action-driven progression and consequences, improve AI stories/NPC/faction behavior, maintain reliable saves and optional voices, run headless simulations, and complete a playable 20–30-minute rise-and-consequence campaign. Keep the Mac awake while active work is running.

Codex prepares the visual handoff, maintains the public API, fixes functional UX/gameplay bugs, playtests the whole product and integrates reviewed visual commits. Codex should record visual defects for the visual owner rather than spending subsequent autonomous turns generating art, reskinning screens or tuning rendering. Visual production is a parallel workstream, not a prerequisite that blocks simulation development.

The visual owner handles art assets, illustrated-city rendering, animation, lighting, UI skin and visual layout under `docs/VISUAL_HANDOFF.md`. Functional UI behavior remains part of Codex's remit. Shared-file/API changes need coordination; neither owner may silently break the other's interface.

Handoff deliverables: written brief and examples, acceptance checklist, ready-to-paste agent prompt, isolated `codex/visual-handoff` worktree and safe preview launcher. Preparing the handoff does not mean another agent has started. Whole-game completion remains unproven until integrated behavior has been tested; no requirement is satisfied merely by assigning it to someone else.

The task goal was recreated with this amended scope on 2026-09-07 after the controller reported no existing goal. It is active and matches this ownership split. No unfinished goal was marked complete to replace its text.

Build and playtest the full Black Ledger single-player mafia vertical slice described in DESIGN.md: a playable 20–30-minute rise from rented housing through contacts, crew and business ownership, into consequential rival incidents; permanent death and new people in a persistent city; validated AI opportunities; reliable transactional saves and optional stable-character speech.

## Current gameplay acceptance priorities

1. Improve narrative correctness without replacing the Go rules with model judgments. The completed 20m18s campaign demonstrated rise, ownership, retaliation and recovery, but failed story-coherence acceptance. The experimental second-model reviewer also failed on real dialogue and remains offline-only.
2. Continue integrated campaign acceptance for arrival encounters, contact variety, located jobs, business ceasefires, danger pacing, resumable interrupted work and stale-draft protection. Main preview now runs verified release b74c9a2 after isolated save-upgrade, ceasefire and interrupted-job checks; its complete saved state and190 receipts were preserved. These targeted checks do not replace a fresh20–30-minute campaign with coherent stories.
3. Continue public-state headless campaigns and targeted browser playtests for progression, permanent consequences, transactional saves and optional voice behavior. Preserve failure evidence and commit corrections.
4. Maintain the visual handoff and interface contract. Another visual agent has not yet been launched; preparation is complete, visual production itself is not.

## Approved visual expansion — 2026-09-07

The user approved a modern 2D isometric city inspired by the visual approach of Gangsters: Organized Crime. Painted noir realism is the core; warm vintage daylight and amber/crimson nightlife are lighting variations. This replaces the schematic cartoon map as the intended visual destination.

Visual-agent deliverables (project requirements retained; no longer Codex production tasks):
- A coherent small neighborhood assembled from reusable illustrated buildings and street pieces, with consistent scale, perspective, anchors and occlusion.
- Moving cars and pedestrians on authored presentation routes; decorative actors do not require simulation agents.
- Independent marquee lights, window glow and limited atmospheric effects.
- Gameplay-selected destinations remain interactive, and visual travel can be skipped without changing committed outcomes.
- Narrative visual sequences consume backend results. They do not calculate combat, rewards, political outcomes or game time.
- Expandable asset/metadata structure allowing more districts and changed building states.
- Browser visual QA at practical desktop and compact sizes, and motion controls/reduced-motion behavior.

The isolated casino study is an asset test, not completion of the neighborhood requirement. Do not shrink the gameplay goal to an art demo. Keep implementation and validation records in DEVELOPMENT.md and commit coherent changes in this independent repository.
