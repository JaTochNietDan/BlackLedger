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

### Staged gunshot audio follow-up

City gunfire now drives sound from the rendered shot beat. Four muzzle pulses
share one timing definition; each may start one cancellable Web Audio shot.
Queued scenes make no requests, stalled/background frames consume missed beats
without replaying a backlog, and muted frames stop any current tail and consume
the beat silently. Skip, expiry, unmount and a world/life change dispose the
handler. Graphics interruption stops current audio and disables further animated
scene playback until reload. These cleanup paths do not modify Go state.

The Theatre suppresses its older sound burst for a city gunfight and for a
co-located killing whose gunfight is staged by the city. Other kinds still use
their existing audio. Shot source/filter/gain nodes disconnect on natural end
or cancellation. Audio allocation failure returns silently so it cannot kill the
renderer's animation loop. Other scene sounds still need cancellation/timing
review; this pass specifically covers staged gunfire.

Evidence: 100 frontend tests, build and the HTTP package pass. Tests cover
30/60/144 FPS beat counts, queue/skip/mute/unmute/frame stalls, idempotent disposal,
and a mocked Web Audio graph's scheduling/disconnection and allocation failure.
Logs: `.runtime/city3d-audio-tests.log`, `.runtime/city3d-build.log`, and
`.runtime/city3d-audio-go.log`.

Isolated browser replay on 8855 requested four shot sounds, then Skip cleared the
scene. An early Skip after one shot left no effects after another 1.1 seconds.
Sound Off replay retained the animated scene with zero shot-audio requests; the
original Sound On preference was restored. All samples stayed revision 0 /
minute 480. The diagnostics count sound handles started; audible output was not
independently recorded, and broad browser/OS audio latency acceptance is pending.

Next priorities remain the requested visual quality: richer street/block art,
character surfaces and varied action/vehicle choreography, along with mobile,
reduced-motion/context-recovery and prolonged resource acceptance. The production
quality goal remains active.


### Gunfire choreography follow-up

Gunfight cues now stage an anonymous articulated shooter with a locally authored
Blender revolver (blued frame, cylinder/barrel, walnut grip, open trigger guard,
sights and a muzzle anchor). The three-second sequence raises the arm, fires four
brief muzzle flashes with recoil and smoke, then lowers the weapon. Light and
particles follow the exported muzzle through the animation. The shooter uses the
same forecourt reservations as casualties and prefers a slot behind the first
casualty along the firing direction. This is a schematic reenactment: existing
Go gunfight cues do not identify a shooter or weapon, so the renderer does not
assign the figure to a named NPC or invent hit outcomes.

Co-located killing cues hold their fall until the first gunfire moment. Shooter
slots are allocated first so a full casualty batch cannot consume every slot
while waiting for a shooter to stage. Other killing cues retain their existing
fall timing. The detailed position/shot count is presentation, not saved combat.

Evidence: 96 frontend tests pass, including four firing pulses at 30/60/144 FPS,
casualty delay, and actual parsed GLB geometry through every aim/recoil frame
checked against the reserved footprint and pavement. `npm run build` and
`go test ./cmd/blackledger ./cmd/qa-fixture` pass. Logs:
`.runtime/city3d-gun-tests.log`, `.runtime/city3d-gun-go.log`,
`.runtime/city3d-gun-export.log` and `.runtime/city3d-build.log`.
There are 25 exported models. Existing person/woman binaries were preserved when
the full exporter rewrote those unchanged models.

Browser fixtures on 8854 (`gunfight`) and 8855 (`gunfight-killing`) use core Witness
and, for the paired case, core Kill to produce replayable records. These are
explicit isolated presentation fixtures, not a campaign combat acceptance test.
The paired replay showed the shooter at (77,38.35), Mara at (80,38.35), and the
player clear at (80,36.65). Before the shot, the arm was raised and fall was zero;
a later sample showed the fall at -1.327 radians. Skip cleared both effects at
revision 0 / minute 480. Screenshots: `qa/city3d-20260913/gunfight.png` and
`gunfight-casualty.png`. The post-scene local view measured 145 FPS / 7.1ms p95,
65 draws and 78,460 triangles; this is not a busy combat performance guarantee.

Still pending: sound currently follows the older Theatre audio timing; align it
with staged muzzle pulses and cancellation. Add varied, convincing multi-actor
combat and vehicle choreography, richer character materials, and broader browser
acceptance for accessibility, resource stability and dense scenes. Final art and
production acceptance remain incomplete.


### Cast appearance follow-up

Added a second locally authored Blender pedestrian: waved/pinned chestnut hair,
a fitted jacket, blouse collar, lapel pin and slacks, with the same independent
arm/hip/knee joints as the fedora model. The renderer selects these two base
silhouettes from the existing portrait cast mapping in `core/voices.go`, including
Mara/Elena and the player's saved face choice. Walking, observed arrivals and
casualty reenactments use the same selection. Both models retain the pedestrian
speed and occupancy envelope. There are now 24 exported GLBs; the new figure is
266,848 bytes. These are two shared base appearances, not individual likenesses;
clothing/skin/age variation and higher-quality character surfaces remain pending.

The old Portrait fallback used signed JavaScript hash arithmetic and absolute
value; Go's FaceOf uses unsigned FNV-1a bytes. Portrait and 3D selection now share
the unsigned calculation. Explicit saved face selections are unchanged. Tests
compare all 24 cast entries and painted identity mappings with the actual Go
source, include FNV vectors, inspect the real GLB articulated nodes and check
both models' exported walking/falling bounds. All 93 frontend tests, the build
and the HTTP package pass. Evidence logs: `.runtime/city3d-cast-tests.log`,
`.runtime/city3d-cast-export.log`, `.runtime/city3d-build.log`, and
`.runtime/city3d-cast-go.log`.

Browser evidence: the isolated killing replay on 8852 now renders Mara with the
waved-hair figure (`qa/city3d-20260913/cast-casualty.png`). In the existing isolated
walking fixture on 8850, selecting Face 3 saved revision 4 without changing minute
645, and the city showed the new player model plus both NPC model types. A normal
walk to Saint Agnes completed at revision 5 / minute 660 in 17,778ms, with visible
articulation and clear pavement (`qa/city3d-20260913/cast-walking.png`). Its zoomed
post-arrival view measured 145 FPS / 8ms p95 / 59 draws / 77,996 triangles. This is
a local sparse-view sample, not broad performance acceptance.


### Event occupancy and framing follow-up

Police and casualty extras now share presentation occupancy with ordinary actors.
Each address has six side bays for police and five forecourt slots for casualties;
the casualty reservation encloses the complete standing-to-fallen motion. A full
scene waits for a free slot instead of stacking bodies. Playback starts when the
slot is available and releases its reservation on expiry or Skip. Falling roots
lift enough to keep the exported person above the pavement. Police flash a red
roof beacon instead of spraying particles around their car. Front-side bays are
preferred because the original left-side default was hidden behind the building.

The browser exposed the old full-width theatre gradient covering the city and
camera controls. City captions now occupy a compact corner panel; expanding the
city during playback works and shows the existing expanded caption/skip control.
Portraits remain in the normal caption, but 3D actors are still generic. An indoor
killing is represented schematically at the building forecourt; this is not yet
character-specific or indoor action choreography.

Evidence: 89 frontend tests pass, including every current route sampled against
all event slots, building/lamp/furniture clearance, slot overflow with a parked
car, the exported character bounds throughout its fall, and traffic waiting for
an event reservation then proceeding after removal. `npm run build` and
`go test ./cmd/blackledger ./cmd/qa-fixture` pass. Logs are in
`.runtime/city3d-scenes-tests.log`, `.runtime/city3d-build.log` and
`.runtime/city3d-scenes-go.log`.

Isolated browser fixtures: killing on 8852 replayed the actual core-produced cue
at Saint Agnes, root (80,38.35), with the player clear at (80,36.65). It remained
revision 0 / minute 480 through replay. Arrest on 8853 used the actual “Go with
them” choice, then replayed the police scene at Ward Street Station (89.6,10).
Skip immediately removed its effect and retained revision 1 / minute 480.
Screenshots: `qa/city3d-20260913/casualty-scene.png` and `police-scene.png`.
The zoomed post-scene view measured 145 FPS / 7.4ms p95 / 63 calls / 72,476
triangles; this is an idle local sample, not a busy-effect performance guarantee.
The preview launcher accepts `killing` and `arrest`; the killing fixture now
preserves its real core cue in LastResult for explicit replay after loading.

Still pending: distinct character appearances, convincing shooters/assassins and
vehicle arrival/departure choreography, richer city art, mobile caption recheck,
reduced-motion/renderer-failure browser checks and prolonged resource stability.


### Movement pacing follow-up

Normal presentation now caps walking at 1.8m/s and cars at 11m/s, replacing the
shared 80m/s ceiling. Walking phase follows actual distance (1.15m per cycle)
instead of wall-clock time, so queued people stop stepping. A 1×/4× travel
control advances the movement clock and occupancy together; saved time and
event results remain unchanged. Arrival completion still follows the endpoint.

Browser evidence: a 32m walking leg completed at 17784ms. Toggling 4× retained
revision 2 / minute 630; the return route, including road crossings, completed
at 7034ms at revision 3 / minute 645. Tests verify physical distance at both
rates and 30/144 FPS. Existing completion tests now allow time appropriate to
their path lengths at the slower speeds. All 85 tests, build and HTTP checks
pass. The browser sample remained 145 FPS / 7.1ms p95. Narrow-screen scene
captions have extra clearance beneath the expanded control row; that CSS change
has not yet had a fresh mobile viewport test.

Continue event choreography and integration of event actors with traffic, plus
same-direction pedestrian passing, varied character appearance and richer city
surfaces. The more readable walking pace is still a basic articulated cycle,
not final foot planting or animation acceptance.

### Pedestrian model and walking-route follow-up

The Blender pedestrian now has a fitted jacket, lapels, pockets, shirt/tie,
cuffs, hands and a shaped fedora. Hip, knee and arm joints produce opposing
left/right steps; the previous mesh-order phase assignment could swing both
legs together. Animated bounds sampled at 48 phases are recorded in the model
manifest and checked against traffic occupancy and pavement height. Close zoom
now reaches 12x, with labels scaling down to retain readable screen size.

Browser play exposed an indefinite walking queue: opposing pedestrians shared
one path. Walkers now use directional pavements with explicit source/destination
crossings. Parcel pitch increased from 28m to 32m to leave room for the stride,
street furniture and parking; building scale and Go travel timing are unchanged.
Lamps moved clear of the walking envelope. A new regression then exposed overly
broad junction reservation: a turn outside the junction blocked parallel traffic.
The check now inspects only the path ahead inside that crossing.

The original blocked trip was repeated in a fresh `city3d-walk` fixture and
completed at progress 1 / 2403ms, revision 1 / minute 615. Controls became
available without Skip. `pedestrian-walk.png` records the playtest. All 84
frontend tests, the production build and HTTP tests pass. The inspected close
view measured 145 FPS / 7ms p95, 155 draws / 98,514 triangles.

The figures still share one outfit and appearance. Gait articulation is better,
but snapshot playback compresses walking distances too aggressively; improve
movement pacing next. Also check same-direction pedestrians behind snapshot-held
travellers, plus event actors against ordinary occupancy. These are still open
acceptance items, alongside the broader art and scene-choreography work.

### Individual venue follow-up

Four new Blender exports replace the shared casino model at The Monarch, Blue
Hour, Golden Lily and Paper Moon. They vary height, width, depth, textured brick
palette and neon colour while retaining the district's Art Deco vocabulary.
Marquee anchors place the names on the front fascia; bulbs hang beneath it.
Rooftop utilities that conflicted with the stepped crowns were removed.

The asset set now contains 23 GLBs. All 81 frontend tests and the production
build pass; footprint and route tests include the new exports. Browser day/night
comparisons are `venues-day.png` and `venues-night.png`. The inspected night
view measured 145 FPS / 7.6ms p95, 223 draws and 118,170 triangles. Repeated
shops, tenements and industrial sheds remain visually repetitive; street layout,
surface wear, pedestrian quality and action choreography still need improvement.

### Multi-way junction follow-up

A four-way test exposed a real deadlock: all four cars stopped at progress
0.4733 and remained there. Junction occupancy now reserves space before bodies
enter the crossing. Straight parallel/opposing lanes can share the reservation;
turning and perpendicular traffic waits outside it. The regression confirms
collision-free completion at 30, 60 and 144 updates per second. All 81 frontend
tests, production build and HTTP tests pass.

`city3d-junction` stages twelve cars across four approaches. Browser actions
advanced the fixture from minute 600 through 609 to 618, with opposing traffic
crossing and the other approach yielding. `junction.png` records the view.
Observed telemetry was 145 FPS / 7.6ms p95, 328 draws / 155,392 triangles. This
checks one busy junction arrangement, not universal gridlock freedom. Resume
address-specific architecture/materials and action choreography next; large
sections of the city still repeat the same archetypes.

### Queued player arrival follow-up

Player travel now ends when rendered occupancy reaches the endpoint, replacing
the unconditional 2.4-second timer. Skip remains available, leaving the city
cancels presentation, and a failed renderer completes the pending presentation.
Travel initiated from another view opens the city so it cannot block controls
behind an unmounted renderer. Journey completion keys include the world ID.

The traffic fixture now starts the player at its shared route's departure. In
browser, a long trip completed with all twelve NPC cars arriving. A return trip
recorded `data-arrival` progress 1 at **3555ms**, revision 2 / minute 764; after
completion the player changed to the pedestrian and parked Hudson, retaining
that revision/time. The new long-route test confirms occupancy remains in transit
at 2.4 seconds and later reaches its endpoint. All 80 frontend tests, the build,
and HTTP tests pass. Renderer failure recovery and OS reduced-motion changes
still need direct browser fault/setting tests; their handling is implemented but
not claimed verified by this trip.

### Traffic occupancy follow-up

The renderer now maintains oriented actor footprints and advances in small spatial
steps. Same-lane cars queue, crossing bodies yield, and a full departure reports
its waiting population. Positions remain at or behind public committed progress.
New journeys can reset progress without inheriting an old trip. The player's
parked car moved out of the travel lane to a separate side position; tests cover
its clearance from every building, furniture band and route.

`city3d-traffic` stages twelve operational cars on one route. Browser evidence
shows all twelve in a separated queue at revision 0 / minute 600, then advancing
through a turn after a real travel command at revision 1 / minute 609. The player
Hudson travelled concurrently. `traffic-queue.png` records the resulting view.
Observed telemetry was 145 FPS / 7.1ms p95, 317 draws / 151,686 triangles before
the close view settled. All 79 frontend tests and the production build pass;
HTTP tests pass. Tests include twelve-car monotonic following, capacity overflow
and recovery, opposing lanes, crossing safety and completion, journey reset,
exported vehicle dimensions and parking clearance.

This is not full traffic acceptance: test more crowded multi-way intersections,
extend player playback to wait for an actual queued arrival, and give waiting
departures a stronger visual indication than the current status count. The
parking position also needs an authored recessed kerb/bay. Event-specific actors
still need integration with ordinary street occupancy.

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
