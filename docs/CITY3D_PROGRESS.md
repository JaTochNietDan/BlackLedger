# Browser 3D city — 2026-09-13 checkpoint

The user superseded the older 2D visual handoff. Work is on `codex/browser-city-3d`
in the independent Black Ledger repository. The goal remains active: this is a
working first asset/renderer integration, not a production-quality acceptance claim.

## Implemented

- One active Three.js city, with orthographic orbit/pan/zoom, keyboard controls,
  picking, an HTML address selector, authoritative travel and building entry.
- Expand/return controls and a single scrolling city/sidebar flow at small widths.
- 28 locations on non-overlapping parcels; 18 locally authored Blender GLBs,
  shared physical texture scale, normal-mapped brick and anchored facade signs.
- Public NPC journeys distinguish walking, Ford, Hudson and Packard. Motion
  interpolates committed snapshots then stops with the game clock. The player
  uses the vehicle captured at departure, even if the journey exhausts its fuel.
- Parked player car, articulated gait, smooth heading changes, night windows,
  lamps, ground light pools and public weather fog.
- Transient explosions/fire/smoke, gunfire flashes, casualty falls and police-car
  presentations. Generic attacks do not pretend to be shootings. Property
  condition persists as darkened materials; damage does not invent an ongoing fire.
- Stable cue deduplication, simultaneous events on scene mount, explicit replay,
  skip and reduced-motion handling. GPU resources are released on scene teardown.

## Bugs found through browser play

1. Damage shading erased base colors and whitened the architecture; fixed.
2. Active city mount discarded simultaneous explosion cues; fixed and tested.
3. Bombing a building selected a casualty from the owner's entire family,
   including NPCs across town. Fixed to use living, non-travelling occupants.
   The repeat playtest killed Rosa Erdos at The Monarch, where she actually was,
   instead of Alma Nagy at Russo Motor Works with contradictory explosion prose.
4. Premature charges produced no explosion presentation; fixed in Go and tested.
5. A narrow viewport split city/sidebar into separately clipped rows; changed
   to one scroll flow and added the expanded city view.
6. Replay initially replayed only the caption; now restages the saved cue batch.

## Evidence

Screenshots: `docs/qa/city3d-20260913/` (explosion, compact desktop, phone, night).
These are staged fixture saves, not earned campaign progress.

- `npm test`: 69 passing tests. New checks inspect every exported building bound,
  every pair of current parcels, and 201 samples along every walking/driving route
  against every parcel with actor clearance. They also cover cue lifecycle.
- `npm run build`: passes. Vite still warns about the large Three.js app chunk.
- `go vet ./...`: passes.
- `go test ./... -count=1`: passed during the iteration (sim took 406 seconds).
  After the subsequent blast-locality fix, the complete core and HTTP suites were
  repeated and passed (core 90.5 seconds). The earlier sim pass is not a fresh
  post-locality-fix campaign balance measurement.
- Browser: click-to-select, travel on foot, vehicle destinations, rotation,
  scroll zoom, keyboard pan, expanded city, replay and immediate skip verified.
- Camera, replay and skip retained revision 1 / minute 720 in the blast fixture.
  Both killing and explosion appeared together at `club` in rendered telemetry;
  skip immediately reduced the effect list to zero.
- Observed median ~145 FPS, p95 frame intervals ~7.5–8.3ms with 28 buildings and
  up to 12 NPC journeys plus the player on this Mac's in-app browser. This is a
  local sample, not a hardware-independent 60-FPS guarantee. Telemetry is on the
  canvas `data-metrics`; the debug HUD requires `?city-debug`.
- Responsive tests requested 980×720 and 390×844 through the browser viewport
  capability; the latter reported a 354 CSS-pixel viewport at its current browser
  scale, with document width equal to viewport width. Normal desktop was restored.

## Continue next

### Street furniture follow-up

A nineteenth Blender asset adds slatted timber benches, galvanized litter bins
and cast-iron hydrants. The renderer instances its four material meshes across
the city. A separate pavement band keeps it out of the reserved building bounds
and public journey paths. New tests use the exported asset bounds to check every
placement against all buildings and 201 samples on every walking/driving route.
All 72 frontend tests and the production build pass. Browser orbit inspection
is recorded in `street-furniture.png`; the inspected view measured 145 FPS,
7.1ms p95, 198 draws and 105,790 triangles. Placement is still repetitive and
will benefit from address-specific dressing. Traffic separation beyond opposing
lanes remains an outstanding requirement.

### Architectural follow-up

The Blender source now authors a separate stucco estate with a tiled pitched roof,
closed gables, shutters, side windows, porch and iron boundary railings. Casinos
have an Art Deco crown, vertical fins, red neon blades and a front marquee. Civic
buildings have stepped towers and clock hands driven by the public saved minute.
Fire escapes moved behind the primary frontage. Browser inspection caught and
corrected missing estate side windows and open roof gables during this pass.

All 71 frontend tests and the production build pass after regeneration. Exported
bounds and every current walking/driving route retain their clearance. Daylight
and night browser inspections are recorded in `estate-day.png`, `civic-day.png`
and `casino-night.png` in the evidence directory. The close civic view measured
145 FPS, 7ms p95, 182 draw calls and 55,966 visible triangles on this Mac. This
sample includes frustum culling and is not a worst-case full-city benchmark.

These additions improve archetype recognition; multiple casinos still repeat
the same facade, and city dressing/material richness remains below the requested
final quality. Continue with individual address variation and denser streets.

### Street routing follow-up

Cars now use directed right-hand lanes, a direct route for same-street trips,
and sampled curved junction turns. The full frontend suite passes 71 tests,
including opposite-direction lane separation, steering continuity at every
current route corner, and building clearance at every curve vertex. The browser
night fixture showed the player's Hudson travelling on the new lane at revision
2 / minute 1292, with NPC Ford/Hudson/Packard journeys alongside it. A preceding
travel action also verified the generic property-attack presentation. Expanded
night rendering measured 145 FPS / 7.2ms p95 before those actions.

This separates opposing lanes; same-lane following distances and intersection
right-of-way still need a traffic presentation system. The road map still looks
too sparse and repetitive for final art acceptance.

The follow-up also fixes travel results with event cues: show the completed trip
first, then stage its saved event batch. A fresh `city3d-night` fixture on 8840
verified the Hudson in motion with no effects, followed by the attack at Bluebird
Laundry, both at revision 1 / minute 1276. Expanded view now includes Skip journey.
New commands clear the previous playing cue so an old event cannot follow a new
trip. The preview launcher itself was exercised successfully for this fixture.

The high-quality 1950s noir art requirement remains unfinished. Prioritize
richer architectural silhouettes/materials, streets that read as a lived-in
neighborhood rather than repeated isolated parcels, period-specific pedestrians,
more convincing action choreography and readable event framing. Test every effect
kind, including police and gunfight, in browser fixtures. Footprint/route checks
are not actor-to-actor collision checks: dense traffic separation and animation
pacing still need work. Exercise reduced-motion changes and renderer failure
recovery in the browser. Recheck GPU/resource stability after many transitions.

Keep the simulation authoritative and retain evidence when a playtest exposes a
contradiction. Do not describe this checkpoint as a finished production game.

## Reproduce

`./scripts/run-city3d-preview.sh city3d` creates a fresh isolated test save on
port 8840. `city3d-night` and `city3d-blast` are alternatives. Set
`BLACK_LEDGER_PORT` if that port is occupied. The blast fixture starts inside
The Monarch's address with a charge and deterministic seed: step inside and use
“Put the charge under The Monarch”. It is an intentionally staged action.

No QA command opened `.runtime/campaign.sqlite3`. Pre-existing edits in
`core/mugging.go`, `core/robbery.go`, `core/armed.go`, and `sim/died_test.go` were
preserved and are excluded from this checkpoint commit. Test runs include those
pre-existing working-tree changes.
