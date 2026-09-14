# Current goal and workstream ownership

## User amendment — playable billiards (2026-09-14)

The pool hall must support playing billiards across its tables, individual games
with wagers, and occasional entry-fee tournaments whose winner receives the
whole prize pool. The user explicitly requires a full-fidelity billiards
minigame with proper physics simulation. An authored room with static balls does
not satisfy this requirement. Eight-ball is the initial implementation choice;
physics, match rules, funded stakes and physical opponent play now exist in
the backend. Exactly-once commands, public table/replay projection and HTTP
retry coverage are implemented. A first playable 3D view has isolated browser
coverage for funded start, placement, physical player/NPC strokes and concession.
Visual refinement, fuller controls/animation, multiple tables and tournaments
remain to be implemented and verified. This extends all earlier requirements.

## Highest priority amendment — map-first UI

The persistent HUD (left navigation and top bar) must now receive the same
period styling as the improved menus (latest September 13 amendment).

Menus must also match the 1950s mafia setting: replace generic menu boxes with
styled in-game surfaces and integrated controls (September 13 amendment).

The city is the full-screen persistent game view. Menus open over the map; they
do not close it. Entering a building replaces the city with the interior as the
main scene. After event animation, show a styled newspaper article and its
appearance sound. If voices are enabled, prepare narrator audio ahead of reveal
so it starts promptly. Preserve all earlier simulation/visual requirements.


## User amendment — detail, interiors and event direction (2026-09-13)

Current objective source: `/Users/jatochnietdan/.codex/attachments/84acfc45-069d-4a1a-ace4-a4a5eee74946/goal-objective.md` (read September 13; supersedes the earlier 35be00d7 objective).

The expanded scope also requires much richer gambling tables and slot-machine presentation, occupants and animations; consistent styling across every menu; debug access to every assassination variant; explosion preambles showing the planter leaving; new authoritative molotov and building drive-by actions with visible damage and escape; and attacks against the player during street travel. Also stage interior occupants at chairs and counters, render cause-specific player deaths, track and display personal/family backstories and interactions, and illustrate moonshine, untaxed cigarettes and crated arms. The latest amendment additionally requires assassinations inside buildings and plausible wire attacks from behind, with animation fidelity sufficient to support them. These additions remain outstanding and do not replace earlier requirements.

The user requires substantially more detailed, high-quality 1950s noir gangster architecture and characters; current blocky assets do not meet visual acceptance. Model detailed building interiors for display on entry, preserving the room's actual public occupants and available gameplay actions. For committed simulation events, automatically open the city and frame the complete action before presenting its animation. Provide a debug mode with selectable action scenes for visual review without mutating the campaign. The user also requests gore and persistent aftermath: bodies, blood and police response should remain until cleanup, while damaged buildings retain their authoritative condition until repaired. Explosions must originate inside the target building with window fire/smoke persisting until a fire-brigade response extinguishes them. Raids need multiple police units and officers; arrests must show officers taking a person into custody. Gun models must represent the actual weapon used by each participant, with NPC ownership of varied guns contributing to their power as player equipment does. Support WASD panning alongside the arrow keys. The user further requires high-fidelity, impactful action choreography with crisp sound, appropriate camera shake/impact, and variety. Raids must show police forcing entry and lingering afterward. Assassinations should support simulation-grounded variants such as a moving car drawing alongside a walking victim for a drive-by; all scenes should visibly play out rather than remain static cast arrangements. Reveal news/results only after the action scene finishes, preserving suspense during playback. Layer contextual audio into action timelines: cries for help, shouts/swearing, vehicle approach, gunfire, screams and tyre-screech escapes; feud-specific utterances should fit the actual public context. Keep the Mac awake during active work. Validate visual quality and performance in the browser as production work proceeds; local high FPS alone does not establish production readiness.

The latest amendment prioritizes detailed assassination choreography: an unsuspecting standing target, an assassin walking up from behind, and one shot to the back of the head. Support multiple scenarios selected from simulation variables and variation; do not substitute static firing poses for the complete approach, attack and aftermath. All earlier quality and simulation requirements remain in scope. The next amendment also explicitly requests assassination blood spatter and realistic, loud, intense gun sounds. The user separately requested smooth continuous WASD/arrow panning.

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

## September 13 follow-up — personal HUD and travel continuity

The user requests removal of redundant City navigation/briefing, a period-styled player identity, matching camera controls and latest-result treatment, no hover flash or clipped outcome, and separation of income/expenses from Ledger into persistent city accounts. Travel should visibly advance the clock at its playback pace, automatically zoom/follow the player, and return the camera to the player outside after leaving an interior. Implemented desktop checks are recorded in DEVELOPMENT.md; this does not close the broader production game goal.

The latest amendment requires takeover and management of The Mariner as a lodging business. Living NPC tenants must actually rent rooms, rental income must derive from those occupants, and ownership must bring business obligations comparable to other premises. This gameplay expansion remains outstanding.

The subsequent amendment requires all living NPCs to have homes appropriate to their wealth and standing, with additional housing where capacity requires it. Mariner leases must therefore be part of city-wide residential assignment, not an isolated occupancy counter.

Latest amendments additionally require a real-estate market with houses changing hands, player purchases and sales, more residential housing, broader income opportunities across progression, and a guiding beginner questline. These remain outstanding.

## September 14 amendment — card scenes and residential city

The current user explicitly authorizes continued visual production: poker must
join the other 3D games, with close overlooking cameras and readable cards,
large scenes taking the map's screen area, and camera movement. Continue overall
stylization, feeling and UX. Expand earning opportunities through progression,
housing stock for every living NPC, individual ownership/purchases/sales and NPC
house trading, with prices reacting to neighborhood crime. Add burglary and
assassination at a target's home when they are present. Expand the map as needed
and complete the planned 3D interiors. This supersedes the historical visual
production restriction; no separate visual agent was launched in this session.
First 3D poker/table-camera progress is recorded in DEVELOPMENT.md. These broader
simulation and interior requirements remain outstanding.
