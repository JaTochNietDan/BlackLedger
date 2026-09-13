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

### Expanded police cast and WASD — September 13

- Raid cues now stage three cars and four uniformed officers; arrest cues stage two cars, two officers and an explicitly identified detainee. The detainee's arms move behind their back. Supporting actors use the existing spacing/reservation, cutaway and camera-envelope paths. This is still a short staged tableau; arrivals, escort, handcuff contact, search and door-entry choreography remain unfinished.
- Found that existing player-arrest cue actors identify the detective, not the prisoner. Added an optional explicit `detainee` field to player confinement; absent identity never substitutes the detective. Debug arrest uses a clearly synthetic detainee. Targeted Go custody/arrest tests pass. Renderer tests cover cast counts, pairwise spacing and detainee identity.
- Added requested WASD panning alongside arrows to the city and first 3D interior. Q/E retains orbit; shortcut modifiers/IME remain excluded. Browser 8860 W–D–S–A returned the camera to x96/z80 within floating-point tolerance. No gameplay commands in this preview run; save stayed revision 3/minute 660.
- Evidence: `raid-group.png`, `arrest-group.png`, `police-cast.json`, `wasd.json`. Both full casts staged without waiting; stopping removed all effects; captured logs were clear. Production build passes. All 153 frontend tests pass (`.runtime/police-cast-tests.log`). The latest user amendment explicitly calls for high-fidelity action timelines, forced-entry raids with lingering presence, crisp impact sound/camera feedback and simulation-grounded drive-by variants. Current tableaux do not satisfy that requirement. Broader art, interiors, fire brigade, weapon fidelity and production acceptance remain active.

### Uniformed crime-scene officers — September 13

- Added a Blender-authored `police-officer.glb` with woven navy uniform, peaked cap, brass cap/shield badges, breast pockets, epaulettes, duty belt, buckle, closed holster and utility pouch. The existing articulated cast rig is retained; uniform pigment is baked into its packed texture for reliable glTF export.
- Persistent killing aftermath now adds two officers alongside its police car at the backend response deadline. Officers face the victim and reserve separate bays; placement waits if available space is occupied. Tests check all pairwise response reservations, replay suppression/cleanup, and actual exported officer bounds through 32 headings.
- The 149-test frontend suite and build passed; the subsequently added actual-officer geometry test also passed. Fresh isolated 8861/save `.runtime/officers-20260913-171058-6418.sqlite3` was advanced by one hour-pass action to revision 1/minute 540. Browser inspection shows body at x80/z38.35, officers at x77 and x83/z38.35, car at x89.6/z42, all visible with no waiting. Sample: 145 FPS/8.5ms p95, 153 draws/530,981 triangles, no captured warnings/errors. Main save untouched.
- Evidence: `uniformed-officers.png` and `.json`. This is stationary crime-scene attendance, not completion of multi-car raids, visible arrests, patrol travel, investigation/cleanup animation or final character art. Those requirements and the broader city/interior/fire/weapon work remain active.

### Interior camera cutaways and idle rendering — September 13

- Split Saint Agnes's wall assemblies into authored Blender groups, retaining material batching within each wall. Camera-facing obstruction now hides the left/back wall assembly independently, including its framed decoration, while furniture remains in the room. Returning to the original angle restores the walls.
- Added Q/E or arrow-key orbit, +/− zoom (0.7–3), and Home reset. The room draws only when controls, size, occupants, selection, game minute or loaded models change. Diagnostics report actual rendered-frame count instead of misleading idle FPS.
- Production build and the exported interior geometry test pass, including both named wall groups, floor support and all standing-bay clearance checks. Browser verification on isolated 8860: frame count stayed at 3 across an 800ms idle observation; reverse orbit hid both walls; zoom capped at 3; Home restored zoom 1 and both walls. No captured warnings/errors. No gameplay commands or main-save changes.
- Evidence: `interior-cutaway.png`, `interior-controls.json`. Initial view is 264 draws/128,272 triangles; reversed view 248/115,576. This improves room inspection and idle GPU use, but does not complete compact-layout, context-loss, broader hardware or overall art-quality acceptance. Other interiors and the outstanding event/response/weapon requirements remain active.

### First furnished 3D interior — Saint Agnes — September 13

- Authored `interior-saint-agnes.glb` in Blender (1,923,660 bytes): walnut wall panelling/cornices, mosaic floor, mirrored bottle display, cupboard doors, marble bar, brass foot rail and stools, oxblood upholstered booths, café tables/cups, copper espresso boiler, mechanical register and pendant lamps. The deterministic source is part of the canonical model exporter.
- Entering Saint Agnes now opens an orbitable/zoomable Three.js room. Up to nine public occupants stand in clear aisle bays with the existing cast models/palettes; selecting a model selects the same existing room action panel. The full accessible roster and premises actions remain below. Other locations still use their existing backdrops.
- All 149 frontend tests and production build pass. Actual exported geometry tests check room bounds, floor support under shoe footprints and radial furniture clearance at five heights around every bay. Browser testing on isolated 8860 selected Leo by native canvas click, confirmed his room panel, and measured steady 144 FPS, 255 draws/128,352 triangles for five occupants. Fixed a deprecated shadow-map warning; final browser logs are clear. Initial asset/shader warmup sampled 22 FPS, so loading performance remains to improve. No gameplay command in this interior test; save stayed revision 3/minute 660.
- Evidence: `saint-agnes-interior.png` and `.json`. This is the first room, not final interior acceptance. Other building interiors, player presence, seating/idle behavior, richer characters, more convincing surface detail, robust graphics failure handling, compact view and repeated-entry resource QA remain unfinished. Event/aftermath/weapon production requirements remain active as well.

### Persistent aftermath rendered — September 13

- The city now consumes saved `aftermath`: a fallen cast character and irregular dark blood pool remain after the short animation, with a police car appearing from the backend response deadline. Bodies and vehicles reserve separate bays in the existing traffic system. Active casualty playback suppresses duplicate persistent bodies; cleanup removes the scene and releases private costume materials.
- All 148 frontend tests and production build pass. New lifecycle tests cover persistent object reuse, separate body/car reservations, playback suppression and exact deadline cleanup. Browser QA uses fresh isolated port 8860/save `.runtime/aftermath-20260913-165531-863.sqlite3`. At minute 480 the body is visible without police; after a real hour-pass command and city re-entry, revision 1/minute 540 shows body x80/z38.35 and police x89.6/z42. Two further hour-pass commands reach revision 3/minute 660, where aftermath is empty. No captured browser warnings/errors. Main campaign untouched.
- Evidence: `aftermath-body.png`, `aftermath-police.png`, `aftermath.json`. This is an initial visual aftermath pass: the blood pool is not complete gore art; police officers, multi-unit response, arrival/cleanup animation and stable staging through replay still need work. Detailed interiors, richer assets, building fire response and weapon fidelity also remain unfinished.

### Saved aftermath lifecycle — September 13

- Added saved/public `aftermath` records for known victims of public killing cues. Each has a stable cue/victim identity, address, occurrence minute, police-arrival deadline (+5 game minutes) and cleanup deadline (+180). Reads return detached active entries, without mutating the save or advancing time. Repeated reports cannot duplicate bodies or restart old deaths. Old saves produce an empty projection rather than fabricated historical scenes.
- Found official deaths returned before creating a killing cue. They now emit that public cue after the existing official consequences, allowing the same aftermath path. Tests cover ordinary/official deaths, immutable projection, result replacement, JSON save round-trip, exact cleanup boundary and stale reports. Full core/HTTP suites passed (118.497s/0.951s) on the initial lifecycle change; after the official/stale-report corrections, focused aftermath/official/killing checks and TypeScript check passed.
- This is backend groundwork, not completed visual aftermath. Bodies, blood, officers, multi-car response and cleanup animation remain to integrate and browser-test. Further user amendments require interior-origin building explosions with window fire/smoke and brigade extinguishing, larger raids, visible custody actions and weapon models matching actual NPC equipment. Existing uncommitted armed-resistance work was inspected but left untouched.

### Debug action previews — September 13

- `?city-debug` now exposes a scene selector, Play and Stop controls for gunfight, assassination, explosion, arrest and raid at the selected address. All use the existing renderer, staging reservations, audio, particles and event camera fitting. The preview has its own immutable world identity and never sends a campaign command. New revisions cancel previews; stop restores the original projection silently.
- Browser testing caught inherited `pointer-events: none` swallowing clicks on the new toolbar. Fixed its pointer handling, then ran all five scenes successfully on isolated 8859. Assassination staged shooter and victim together; explosion spawned 12 fragments; arrest/raid staged police cars. Stop left no effects and the save remained revision 0/minute 480. Captured 145 FPS/8.3ms p95 with no warnings/errors; this is a single local sample.
- Evidence: `debug-scenes.json`, `debug-explosion.png`. Production build and 147 frontend tests pass, including immutable campaign input and silent restoration after preview. Preview art is still the existing limited scene art: police currently use the car presentation, with officers/cleanup/gore still to implement. Richer models, detailed interiors and persistent aftermath remain unfinished.

### Event camera fitting — September 13

- Active committed cues now fit a world-space action envelope into the orthographic camera, independently of prior zoom. The envelope includes alternative staging bays for same-address actors and, for explosions, the building plus a particle margin. Camera orientation is preserved; replay explicitly opens the city from an interior.
- All 146 frontend tests and production build pass. New projection tests check every envelope corner across narrow/square/wide viewports, four rotations and both previous zoom extremes. Browser replay on isolated 8859 verified the gunfight and casualty together, then entered Saint Agnes and replayed again: city reopened at zoom 6.9055, target x80/z38.35, both effects staged, revision 0/minute 480 unchanged. Captured sample 145 FPS/7.7ms p95, no browser warnings/errors.
- Evidence: `event-framing.png` and `event-framing.json`. Explosion/police framing still needs browser inspection; simultaneous events at separate addresses, caption occlusion at compact sizes and scene timing during slow asset loading remain to address. The new user-requested debug scene selector and detailed 3D interiors remain pending. A further amendment adds gore and persistent bodies/police until cleanup. `core/world.go` stores NPC `DiedAt`; public aftermath projection and cleanup lifecycle still need design and implementation, rather than retaining transient cue objects indefinitely. Existing architecture/characters still do not meet the requested quality standard.

### Harbour frontage and revised visual acceptance — September 13

- Added locally authored Blender timber landing and masonry quay assets at Pier 14. Individual deck boards, grain/normal textures, submerged piles, iron collars, bollards and rubber fenders give the harbour a distinct construction. Instanced retaining-wall sections leave the landing opening clear. Animated normal-mapped water extends beyond the western city edge; motion settings stop the water clock.
- Added a paved waterfront promenade. Browser inspection caught stretched paving UVs; corrected them to physical scale before the production build. The landing is selectable as the existing Pier 14 destination. Travel remains tied to the public entrance and existing routes.
- All 145 frontend tests and production build pass. New tests inspect actual GLB bounds, deck elevation, quay depth and route clearance. Existing bundle-size warning remains. Browser evidence in `harbour-landing.png`, `harbour-quay.png` and `harbour.json` uses isolated port 8858. Final sampled view: 145 FPS, 7.4ms p95, 167 draws, 574,736 triangles, one actor/28 buildings, 234 geometries/130 textures; no captured browser warnings/errors. This is a local sample, not broad hardware acceptance.
- The amended objective explicitly rejects current blockiness and adds detailed 3D interiors plus automatic event framing. Existing interiors are painted backdrops with functional public occupants/actions. Current event focus selects a parcel but preserves zoom; complete scene fitting remains to implement. Interior production, richer architecture/characters, harbour activity and final quality acceptance remain unfinished.

### Painted pedestrian crossings — September 13

- Added paired crossing lines at junctions, centred on existing walking lanes and connecting raised pavement islands. The 6×5 grid gets 98 crossings/196 line segments; exterior sides without destination pavement are omitted. Crossing paint shares the existing instanced lane-marking draw and adds no texture or draw call.
- All 143 frontend tests and production build pass. New geometry checks prove every painted corner is on the road surface, both crossing ends meet existing pavement within the city, there are no duplicate segments, and every interior junction has paired marks around all four walking lanes. Existing route/asset clearance tests also pass; bundle-size warning remains.
- Isolated 8858 browser playtest walked Saint Agnes → The Mariner. Captured the player at x65.43938/z27.35 between the junction crossing lines, then arrival at revision 8/minute 645, progress 1, 28,116ms presentation time. Final three-actor view sampled 145 FPS / 7.4ms p95, 78 draws and 508,582 triangles, with no captured warnings/errors. Main save untouched.
- Evidence: `docs/qa/city3d-20260913/pedestrian-crossing.png`, `crossing-markings.json`. Markings change no routing or right-of-way rules. Existing routes can also cross mid-block at their endpoints; broader pedestrian choreography, dense traffic acceptance and final city art quality remain unfinished.

### Quieter parcel and player markers — September 13

- Replaced the large building-selection circle with muted brass parcel-corner marks, rendered in one eight-instance draw. The marks stay within the parcel pavement and no longer cut across the architectural silhouette as a broad loop.
- The pedestrian marker shrinks at close zoom and retains a bounded minimum screen presence at wider views. While driving it scales to the vehicle envelope and rotates with its heading; on arrival it returns to the pedestrian size. Both marker materials avoid tone-mapping washout and depth writes.
- All 141 frontend tests and production build pass. Browser QA on isolated 8858 verified whole-city/close/building-scale readability and a native canvas click selecting The Mariner (`room`). On isolated 8857, drove Thorne & Sons → Ruttledge & Vance → Thorne & Sons and captured the Hudson marker through a junction turn. Final arrival was revision 4/minute 1350, 5,034ms presentation time, progress 1. Main save untouched.
- Evidence: `docs/qa/city3d-20260913/player-marker-close.png`, `parcel-selection.png`, `vehicle-marker.png`, `city-markers.json`. Final night view sampled 145 FPS / 7.1ms p95, 83 draws and 509,024 triangles, two actors/28 buildings, with no captured warnings/errors. This is local browser evidence, not broad hardware or accessibility acceptance. Existing bundle-size warning remains.
- The city still needs denser, less repetitive architecture and higher-quality character/event art. These marker changes improve readability; they do not establish final visual acceptance.

### Authored hair and headwear silhouettes — September 13

- The Blender male model now contains separately selectable side-parted hair, receding hair, fedora and cloth-cap groups. Bald portraits hide both hair groups. Crown surfaces have subtle authored combing relief; headwear selection is mutually exclusive. Fixed named cast retain their previous fedora presentation; generated portrait choices use their corresponding hair/headwear direction.
- Runtime appearance-key tracking refreshes a pedestrian when its portrait palette/silhouette changes even if the base model remains the same. Geometry and texture resources remain shared; private tinted materials are released with replaced actors.
- All 141 frontend tests and production build pass. Actual GLB tests cover each selected silhouette, mutually exclusive visibility and existing walking clearance. Rest and 48-pose motion bounds are exactly unchanged; the male GLB adds 39,916 bytes. Existing bundle-size warning remains.
- Browser inspection on isolated 8858 used the actual Settings portrait controls for faces 4, 1, 10 and 19, showing full hair, receding hair, bald scalp and cloth cap at zoom 32. Evidence: `docs/qa/city3d-20260913/hair-full.png`, `hair-face-1.png`, `hair-face-10.png`, `hair-face-19.png`, `cap-walking.png`. Main save untouched.
- The cap-wearing player completed The Mariner → Saint Agnes at revision 6/minute 630 in 17,784ms presentation time. The final eight-actor view sampled 145 FPS / 7.8ms p95, 55 draws and 496,676 triangles, with no captured warnings/errors. Restored the fixture’s original portrait using Settings afterward; `hair-variants.json` records the checks. This is local browser evidence, not broad hardware acceptance.
- These silhouettes improve cast variety, but faces, body proportions and hair/clothing detail remain simplified. Female hairstyle geometry, additional period clothing and full production character art acceptance remain unfinished.

### Staged-event visibility and cast acceptance — September 13

- Fresh combined shooting/killing browser fixture confirmed the victim's cast palette and fall timing, then exposed a visibility defect: orbiting behind Saint Agnes hid both event actors completely. Player-only cutaways did not protect the staged scene.
- Extended local building cutaways to the visible staged cast at the active event address when player following is inactive. Parallel sight lines test each staged actor; the screen opening encloses the group. Camera orientation and saved event state remain unchanged, and unstaged actors are not revealed. Cutaways restore after extras expire.
- All 140 frontend tests and production build pass. Browser evidence includes the actual before/fixed reverse angle and restored facade: `docs/qa/city3d-20260913/event-occluded-before.png`, `event-cutaway-after.png`, `event-cutaway-restored.png`. The fixed scene shows Mara falling after gunfire begins while the anonymous shooter remains visible.
- Three additional replays on isolated port 8859 returned to zero effects/cutaways and identical 227 geometry / 126 texture counts; final views measured 145 FPS, 7.9–8.3ms p95, 126 draws and 531,296 triangles. No captured warnings/errors. Evidence: `event-cast-cutaway.json`, `cast-event-start.png`, `cast-event-fall.png` in the same directory. Fixture: `.runtime/gunfight-killing-20260913-160027-83302.sqlite3`; revision 0/minute 480 stayed unchanged throughout presentation replay. Main save untouched.
- This closes the prior cast palette browser-check gap for the combined casualty scene. Broader event/address combinations, police staging visibility, explosions without cast extras and street-furniture occlusion still need acceptance. Choreography, character geometry and overall art quality remain below final production acceptance.

### Stable pedestrian cast palettes — September 13

- Added authored muted suit, skin, hair, hat and shirt palettes using the shipped noir portrait sheet as colour direction. Public face choice selects the generated palette; fixed painted identities retain stable choices and generated IDs use the existing cast hash. This is appearance presentation, not an assertion about equipment, wealth or faction. Geometry, animation joints and clearance remain unchanged.
- Each actor clones only its tinted materials, sharing the existing weave/normal maps and geometry. Suit pieces reuse one private material per source. Casualty extras use the same identity palette; generic shooters remain anonymous. Removal and event expiry/cancellation dispose private materials without disposing shared textures.
- All 140 frontend tests and production build pass. New checks cover portrait selection, deterministic fallback, invalid face values, material sharing within one actor, isolation between actors, and shared texture survival on material disposal. Existing bundle-size warning remains.
- Browser QA used fresh isolated walking save `.runtime/city3d-walk-20260913-155609-81907.sqlite3` on 8858. Inspected player and Mara at close zoom, then walked Saint Agnes → The Mariner: revision 1/minute 615, 28,117ms presentation time, completion progress 1. Started with 13 visible pedestrians; moving evidence includes both models. Final 11-actor view sampled 145 FPS / 7.1ms p95, 72 draws, 511,180 triangles, with no captured warnings/errors. Main save untouched.
- Evidence: `docs/qa/city3d-20260913/cast-mara.png`, `cast-walking.png`, `cast-wardrobe.json`. This does not establish broader hardware or full event acceptance. Body proportions and clothing/hairstyle geometry still repeat; palette variation alone does not meet final character art quality. Casualty palette integration was not separately browser-playtested in this pass.

### Slate and mineral-felt roof materials — September 13

- Twelve Blender building exports now carry purpose-authored roof colour, normal and roughness maps. Thorne & Sons and The Mariner use staggered slate courses at a consistent 2.5m UV tile scale, mirrored on opposing slopes. Removed the funeral roof's oversized course bars. Ten flat-roof models use mineral-surfaced felt beneath their existing seams and rooftop equipment.
- Browser inspection caught and corrected sRGB encoding that initially made the slate too dark. Slate maps retain 256px detail; subtle felt maps use 128px to limit repeated embedded-image download cost. The final exports add 877,212 bytes total, with every exported building bound unchanged.
- Evidence: `docs/qa/city3d-20260913/slate-roof.png`, `mariner-slate.png`, `mineral-roof.png`, `roof-materials.json`. Close slate inspection on isolated port 8847 sampled 145 FPS / 7.7ms p95, 93 draws and 527,946 triangles, three actors/28 buildings. This is local idle-view evidence, not broad device acceptance. Main save untouched; no gameplay command issued.
- All 138 frontend tests and the production build pass after final regeneration. The final flat-roof browser sample measured 145 FPS / 7.1ms p95, 41 draws and 514,398 triangles, with no captured warnings/errors. The existing build bundle-size warning remains.
- Materials improve surface readability, but architecture still repeats, roofs need address-specific weathering/detail, and the overall requested production art quality remains unfinished.

### Follow-camera building cutaways — September 13

- Followed players/cars now receive a local, softly dithered opening through intervening buildings. Orthographic sight lines use actual mesh intersections at 10Hz, with cached building bounds as an initial filter. Private condition materials carry cutaway uniforms; normal opacity/depth behavior and shared source materials remain intact. Fades restore on clear orbit or follow release; motion-off uses immediate transitions.
- All 137 frontend tests and the production build pass. New checks exercise off-centre orthographic sight lines, intervening versus behind/beside geometry, private material isolation and full restoration without shader recompilation. The existing bundle-size warning remains.
- Browser verification reproduced the obscured Thorne & Sons view on isolated port 8847 at zoom 32. Its `chapel` location ID was the sole blocker. The screenshot shows the player through a local roof/facade opening; stopping follow restored the solid facade, and orbiting to a clear side removed the blocker while retaining follow. Evidence: `docs/qa/city3d-20260913/player-cutaway.json`, `player-cutaway.png`, `player-cutaway-restored.png`.
- This three-actor close view sampled 145 FPS / 7.7ms p95, 54 draws, 498,680 triangles, with no captured warnings/errors. Main save untouched; no gameplay command issued. This does not establish dense moving-scene or broader hardware acceptance. The cutaway currently handles buildings while following the player; street furniture, event targets and NPC focus need broader visibility design. Character anatomy/variety and production art quality remain unfinished.


### Pedestrian fabric and close inspection — September 13

- Both Blender pedestrian models now embed woven wool colour, normal and roughness maps, with a consistent 20cm UV tile across the articulated clothing. Restrained eyes, brows and a mouth seam add facial definition. Exported rest and walking bounds remain exactly unchanged; the pair adds about 112KB.
- Maximum orthographic zoom increases from 12 to 32. Keyboard pan scales inversely with zoom, retaining the original 5m step at initial zoom and allowing a measured 0.2578125m step at maximum zoom. Manual pan releases player following as before.
- All 135 frontend tests and production build pass. The shooter geometry test now uses the same texture-free material loading approach as other Node geometry tests; actual embedded textures were inspected in the browser. The existing bundle-size warning remains.
- Browser evidence on isolated port 8847: `docs/qa/city3d-20260913/pedestrian-wool.png` and `pedestrian-wool.json`. Close view sampled 145 FPS, 8.2ms p95, 53 draws and 509,074 triangles with three actors/28 buildings; no captured warnings/errors. No gameplay command or main save access. This is a local idle-view sample, not a dense animation or hardware acceptance benchmark.
- Character silhouettes remain visibly simplified, with only two clothing/appearance archetypes. Facial anatomy, identity variety and broader final art quality remain unfinished. Buildings can occlude a followed actor at some camera angles; automatic occlusion handling remains open.

### Follow the rendered player (2026-09-13)

- “Find me” now toggles following of the actual rendered player/car instead of centering only the player's address. The camera keeps its offset, zoom and orbit while translating with the actor. Keyboard zoom/rotation and mouse orbit preserve following; manual pan/reset, touch interaction, building selection and explicit address/scene focus release it. The button exposes its pressed state and becomes “Stop following.”
- All 134 frontend tests and production build pass. Browser verification on isolated 8857 drove The Mariner → Thorne & Sons, reaching revision 2/minute 1332 in 20,742ms presentation time. All 109 recorded samples retained following; camera target versus visible player position had zero measured x/z error. After reloading the final build, keyboard zoom/rotation and mouse orbit retained tracking, while keyboard pan and address focus released it.
- Evidence: `docs/qa/city3d-20260913/player-camera-follow.json` and `player-camera-follow.png`. No captured browser warnings/errors; main save untouched. Touch release is implemented but not exercised on a physical touch device. NPC following, broader device acceptance and final city art/interaction quality remain unfinished.

### Night vehicle lamps and road pools (2026-09-13)

- Night journey vehicles now illuminate their authored head/tail lamps and project two soft forward pools through one shared instanced draw. Lamp materials are private to each actor and released when it leaves the scene, so parked cars can switch off without changing other vehicles or source models. Paused public journeys retain lamps; parked cars and daylight traffic do not. Pool capacity grows with the actor count, releasing replaced instance buffers.
- All 134 frontend tests and production build pass. Tests cover bumper-relative pool positions through all headings and current vehicle lengths, and transparent tapered mask boundaries. Browser testing caught an initial mistake that treated a paused journey as parked; the final rule uses the journey path rather than frame-to-frame movement.
- New isolated night server 8857 uses `.runtime/city3d-night-20260913-152215-71362.sqlite3`. Initial 14-actor/28-building view showed 18 pools at 145 FPS, 7.9ms p95, 391 draws, 604,466 triangles. Saint Agnes → The Mariner drive reached revision 1/minute 1269 in 2,917ms presentation time; 15 sampled travel frames had player lamps on, then the parked Hudson had lamps off. Daylight rainy fixture 8856 showed zero pools/illuminated actors. No captured console warnings/errors; main save untouched.
- Evidence: `docs/qa/city3d-20260913/night-headlights.png` and `night-headlights-metrics.json`. These are ground light approximations; they do not cast vehicle headlight shadows onto buildings. Staged-police projections, richer wet reflections, broad-device performance and final city art acceptance remain unfinished.

### Backend weather in the 3D city (2026-09-13)

- Public `sky.kind` now drives overcast/rain lighting and fog range; `sky.wet` darkens and reduces roughness on road/pavement materials, including the backend's drying-day wetness. Rain uses one LineSegments draw with a reusable 1,800-streak vertex buffer. It follows elapsed presentation seconds independently of Go time and travel speed. Motion-off and OS reduced-motion paths hide precipitation while retaining wet surfaces. Teardown now collects line/point geometry and materials as well as meshes.
- Added `city3d-rain` to the isolated QA launcher. It chooses a campaign ID whose real `World.Sky()` returns rain, preserving the public weather contract rather than overriding the response. Port 8856 uses `.runtime/city3d-rain-20260913-151635-69573.sqlite3`; main save untouched.
- All 132 frontend tests and production build pass; QA fixture/server builds and startup pass. Tests cover rain versus residual wetness, bounded/finite buffers, elapsed-time equivalence at 30/60/144 FPS and precipitation resource disposal. Browser motion-off/on returned `rainVisible=false/true` with `wet=1` unchanged; original motion preference restored. No gameplay commands issued.
- Evidence: `docs/qa/city3d-20260913/rain-city.png` and `rain-city-metrics.json`. Expanded 14-actor/28-building local view sampled 145 FPS, 7.2ms p95, 422 draws, 625,730 triangles with no captured warnings/errors. Rain visibility and muted daylight are verified; night rain, splashes, reflections, broad-device performance and final city art acceptance remain unfinished.

### Localized condition staining (2026-09-13)

- Replaced uniform whole-building darkening with stable surface-coordinate staining driven by the public condition value. Spatial variation preserves readable facade detail; it implies neither ongoing fire nor structural collapse. Full repairs remove staining without altering source textures. Healthy buildings bypass noise calculations, and condition changes update existing uniforms instead of recompiling materials.
- All 129 frontend tests and production build pass. New tests cover clamping, repair restoration, uniform reuse and isolation between damaged/healthy clones sharing source textures. Browser shader compilation and close inspection passed on isolated port 8843 (Monarch condition 38) and port 8847 (condition 100), with no captured warnings/errors. Evidence: `docs/qa/city3d-20260913/condition-stains-38.png`, `condition-stains-100.png`, `condition-stains-metrics.json`.
- Damaged-view local sample: 145 FPS, 7.4ms p95, 186 draws, 582,816 triangles, two actors/28 buildings. Repair restoration is unit-tested; the browser comparison uses separate healthy/damaged fixtures, not a repair command. Main save untouched. Broken-window/rubble variants, richer structural damage and broader device/art acceptance remain unfinished.

### Distinct Mariner lodging house (2026-09-13)

- The player's starting home now uses a dedicated Blender model instead of the generic tenement: three storeys of weathered brick, limestone courses, sash windows, a single sheltered front entrance and lodging sign, closed brick gables, pitched slate roof/ridge, chimney pots and rear iron escape. Its existing `room` location, travel endpoints and gameplay remain unchanged.
- All 127 frontend tests and production build pass. Actual GLB raycasts verify the exposed front door and slate roof coverage including the ridge; existing tests cover every building footprint and route, parking/event bays and blast-fragment pavement clearance with the new model.
- Isolated port 8847 browser inspection covered front/rear orbit and close zoom, then native coordinate clicks selected Saint Agnes and The Mariner (`bar` → `room`) without a gameplay command. Evidence: `docs/qa/city3d-20260913/mariner-lodging.png`, `mariner-rear.png`, `mariner-metrics.json`. Local idle view: 145 FPS, 7.1ms p95, 141 draws, 545,216 triangles, three actors/28 buildings, no captured warnings/errors. Main save untouched.
- This removes one repeated silhouette. Other repeated commercial buildings, oversize pavement, stronger district character, richer event choreography and broader performance/art acceptance remain unfinished.

### Unused parcels as fenced yards (2026-09-13)

- The two unused cells in the current 6×5 grid now render locally authored Blender vacant yards: textured earth/gravel, silvered timber boards with irregular tops and gaps, posts/rails/fixings, and sparse folded weeds. The material groups are instanced across unused cells. Kerb stones, drains and manhole treatment now cover these parcels as well as occupied blocks.
- Vacancies are derived from unoccupied grid coordinates, with no location IDs or gameplay actions. Real addresses retain their positions and commands; a newly occupied cell no longer receives yard scenery. Yard geometry remains inside the existing 17m building reserve.
- All 126 frontend tests and production build pass. Tests cover exact vacant/occupied partitioning and every current pedestrian/vehicle route against yard bounds, with tyre/drain checks extended to unused cells. Browser orbit/zoom inspection on isolated port 8847 is recorded in `docs/qa/city3d-20260913/vacant-yards.png` and `vacant-yards-metrics.json`: 145 FPS, 8.2ms p95, 169 draws, 573,968 triangles, three actors/28 buildings, no captured console warnings/errors. Main save untouched and no gameplay command issued.
- The empty-cell placeholder appearance is resolved, but repeated architecture and excessive paved frontage around occupied buildings remain. This does not establish final art quality or broad hardware performance.

### Guarded pavement trees (2026-09-13)

- The instanced street-furniture set now includes one compact, Blender-authored tree beside the rear bench band: tapered trunk, branching, iron guard, radial grate, soil opening and 540 individually folded leaf blades. Three locally authored 128px leaf-vein textures provide restrained foliage variation. This adds decorative planting without simulation entities or commands.
- The complete exported canopy remains outside the building envelope, walking/vehicle paths, parked cars and event bays. All 125 frontend tests and production build pass; existing all-route/all-lot tests use the enlarged exported bounds, and a new GLB test verifies leaf geometry and embedded UV textures.
- Browser orbit/zoom inspection on isolated port 8847 is recorded in `docs/qa/city3d-20260913/street-trees.png` and `street-trees-metrics.json`. Three actors/28 buildings sampled 145 FPS, 8.2ms p95, 159 draws and 531,760 triangles, with no captured warnings/errors. This is a local idle-view measurement; broader hardware and dense moving-scene acceptance remain open. Main save untouched; no gameplay command issued.
- Trees soften a small part of the empty pavement. The repetitive parcel layout, oversized paved frontage, architectural variety and final art quality remain unresolved; this is not acceptance of the finished city.

### Corrugated industrial roofing (2026-09-13)

- Garage/dealer workshops and dock/haulage sheds now have closed corrugated roof meshes with real 7cm ridges, embedded 256px weathered zinc colour/normal textures, and physical-scale UVs. Workshop vents have flashing and rain caps. Roof parts join by material during Blender export, avoiding a draw per ridge.
- Regeneration exposed a latent dealer recipe bug from articulated car pivots: translating both parents and children put display wheels outside the lot. The static display Ford now explicitly uses the non-articulated recipe. The failed footprint and blast-pavement tests prompted this correction; all 124 frontend tests and production build now pass. New checks inspect exported texture slots, ridge heights and local footprint; existing checks cover all building extents and route/event clearances.
- Browser close inspection on isolated port 8847 shows the roof profiles and vents (`docs/qa/city3d-20260913/corrugated-workshop.png`). Local close-view sample: 145 FPS, 8ms p95, 167 draws, 426,896 triangles, three actors and 28 buildings. Four rebuilt assets add approximately 966KB combined and eight texture allocations; broader loading/device acceptance remains open. Metrics/log evidence is `corrugated-workshop-metrics.json`. No gameplay command was issued; main save untouched.
- The city still needs stronger architectural variety, less empty paved frontage, more convincing streets and richer event choreography. These roof improvements do not establish the requested final art quality.

### Individual front-wheel steering (2026-09-13)

- Inside and outside front wheels now use distinct Ackermann angles from the authored axle spacing and lateral pivots. The central steering limit is derived from the inside tyre's existing 0.5-radian limit, preserving the 2.35m moving clearance envelope. Parked wheels remain straight. This is presentation geometry; Go travel outcomes and timing are unchanged.
- All 123 frontend tests and production build pass. Tests check a shared turn centre for both directions and every vehicle, mirror symmetry, and actual GLB track/wheelbase, combined roll/steering bounds and tyre contact.
- Isolated port 8848 browser journey Saint Agnes → Thorne & Sons reached revision 3/minute 668 in 17,831ms presentation time. Among 196 moving-wheel samples, peak wheel angle was 0.499989 radians and maximum discrepancy between the two calculated turn centres was 7.11e-15m. Parked angles returned to zero; no browser warnings/errors were captured. Evidence: `docs/qa/city3d-20260913/ackermann-steering.json`. Final idle sample was 145 FPS/8.1ms p95 with two actors; this is not a moving-load or broad-device performance claim. Suspension, differential wheel roll, richer vehicle bodies and broader art acceptance remain unfinished. Main save untouched.

### Renderer resource teardown (2026-09-13)

- City unmount now disposes instanced buffers, shadow targets and decoded image bitmaps alongside shared geometry, materials and textures. One collection covers the scene and loaded prototypes to avoid repeated disposal of shared assets; late model loads use the same cleanup. The old renderer explicitly releases its WebGL context after cancelling animation and listeners.
- All 122 frontend tests and production build pass. Disposal-event tests cover shared resources, shadow maps, instances, bitmap deduplication and late-loading models. The existing build chunk-size warning remains.
- Eighteen street/interior cycles on isolated port 8847 retained 200 geometries and 70 textures per mounted city, with zero city canvases inside and no captured browser errors or warnings. Evidence: `docs/qa/city3d-20260913/renderer-resource-cycles.json`. These are renderer counters, not a measurement of OS GPU memory. This three-actor fixture does not establish long-session or broad-device performance; larger scenes, richer art and wider production acceptance remain open. No gameplay commands were issued and the main save was untouched.

### Front-wheel steering and tyre sweep clearance (2026-09-13)

- Front wheel pivots now steer through bends using curvature averaged across the wheelbase and distance-based easing. Roll and steering use YXZ order, rear wheels remain aligned, and parked assignments reset steering to straight. The rendered steering limit is ±0.5 radians.
- Moving traffic reserves 2.35m width for the tyre sweep; straight-wheel parked cars and staged police use their 2.15m envelope. An initial attempt to move parking outward failed pedestrian-lane clearance; the final change preserves existing parking/bay positions and explicitly distinguishes stationary occupancy. All 120 tests and production build pass, including every-route parking/event clearance and actual GLB bounds through combined roll/steering angles.
- Isolated port 8847 drive Ackerman & Son → Thorne & Sons reached revision 5/minute 791 in 6,194ms presentation time. Browser samples peaked at 0.499897 radians through bends and eased to approximately zero on the final straight; parked steering was zero after arrival. Evidence: `docs/qa/city3d-20260913/steering-travel-samples.json`.
- Port 8848 expanded fixture with 14 actors and 28 buildings sampled 145 FPS, 7.2ms p95, 506 draws and 537,080 triangles with no captured console errors (`steering-traffic-metrics.json`). Main save untouched. Suspension, individual inside/outside wheel angles, richer vehicle bodies and wider hardware/art-quality acceptance remain unfinished.

### Distance-driven vehicle wheels (2026-09-13)

- Ford, Hudson, Packard and police models now have four separate wheel pivots, 32-sided tyres/whitewalls, valves and hub bolts. Wheel submeshes are joined by parent/material to retain animation without a draw per bolt. Static undertaker hearse generation explicitly retains its non-articulated form.
- Wheel roll follows actual rendered distance at a 0.37m tyre radius. Queued/parked cars do not accumulate roll, initial placement does not imply travelled distance, and reduced-motion presentation does not animate wheels. Steering angle and suspension remain unfinished.
- All 117 frontend tests and production build pass. Geometry tests rotate every exported wheel through a full turn, checking four pivots, material grouping, footprint and tyre contact; phase tests cover 30/60/144 FPS and waits. Browser drive on isolated port 8847 reached Ackerman & Son at revision 4/minute 782 in 5,034ms presentation time. Twenty-one straight sampled segments matched distance/radius phase to 1.8e-15 radians; parked phase remained zero on repeated reads. Evidence: `docs/qa/city3d-20260913/wheel-travel-samples.json` and `wheel-model-contact.png`.
- Expanded port 8848 traffic fixture (14 actors, 28 buildings) sampled 145 FPS, 7.7ms p95, 506 draws, 537,080 triangles locally. Evidence: `wheel-traffic-metrics.json`. This does not establish broad hardware performance. The main save was untouched. Vehicle body polish, steering, suspension and wider city-quality acceptance remain open.

### Explosion audio under casualty captions (2026-09-13)

- Rendered explosion effects now own one cancellable blast sound at their visible onset, independently of which cue wins the caption priority. A casualty caption suppresses its generic killing sound when the same address and minute already contain a rendered explosion or gunfight. This fixes the charge casualty replay previously producing generic shots while the city rendered an explosion.
- Late/muted onset is consumed without a backlog or replay on unmute; Skip, graphics interruption and teardown use the existing audio disposal path. The Theatre ownership flag is now `stagedAudio`, covering both blast and gunfire.
- All 115 frontend tests and production build pass. New tests check location/minute matching and single playback at 30/60/144 FPS, late onset, mute and idempotent disposal. Isolated port 8843 browser replay showed the Rosa Erdos casualty caption alongside one explosion audio start, then clean Skip and no console errors. Evidence: `docs/qa/city3d-20260913/blast-audio-pair.json`. No gameplay command or clock advancement occurred; this is an audio start diagnostic, not an acoustic recording.
- The cue schema does not express a general causal graph. Same-address/minute association matches these committed batches but does not prove every possible compound incident. Broader event choreography, sound design and visual-quality acceptance remain unfinished.

### Instanced masonry debris (2026-09-13)

- Added a Blender-authored chipped clay fragment with embedded masonry colour/normal textures. Each explosion uses twelve instances in short tumbling arcs; they settle with rotation-aware surface support, scatter across the facade pavement and fade with the smoke. No debris persists as an invented gameplay obstacle or additional damage.
- Fragment envelopes remain separated and off roads across every current building. Fragments overlapping visible traffic or staged character footprints are suppressed, preserving the committed outcome. Per-effect geometry/material/instance allocations are released on timeout, Skip, world reset and scene teardown.
- All 113 frontend tests and production build pass. Added tests cover deterministic paths, separated envelopes, settled transforms, actor exclusion at all headings, exported fragment bounds and every facade's pavement clearance. Browser replay on isolated port 8843 inspected flight/settling, verified natural cleanup and immediate Skip cleanup, and reported no console errors. Evidence: `docs/qa/city3d-20260913/blast-debris-airborne.png` and `blast-debris-settled.png`. The saved noon clock was unchanged.
- A local close-view sample during the earlier narrower scatter pass measured 145 FPS/8.3ms p95; this does not establish broad hardware or multi-blast performance. Richer structural damage, varied fragments and final explosion/art-quality acceptance remain unfinished.

### Facade-anchored blast and rising smoke (2026-09-13)

- Explosions now originate at the target model's actual exposed front bound instead of its pedestrian arrival point. The Monarch browser replay reports origin (112,0.25,72.24), 3.59m nearer the building than the former curb arrival point. A short light pulse/fire burst gives way to slower rising smoke; billows use a locally generated irregular alpha texture and fade while expanding. No new gameplay consequences or persistent fire are inferred.
- All 111 frontend tests and production build passed. New tests cover bounded visible particles, finite trajectories, fire ending before smoke, three-second completion, soft texture edges and late opacity fade. Existing event/traffic/asset tests remain passing.
- Isolated port 8843 browser replay inspected close-range fire, rising smoke and final fade, then verified an empty effect list and no captured console errors. Evidence: `docs/qa/city3d-20260913/blast-fire-origin.png`, `blast-rising-smoke.png`, `blast-smoke-fade.png`. Replay left the saved clock at noon. Smoke uses the existing 32-instance effect budget and one additional shared texture.
- Debris, richer facade damage, volumetric smoke and broader multi-event performance acceptance remain unfinished. This is an improvement to origin and timing, not final explosion-quality acceptance or a claim that every visual effect is complete.

### Cancellable explosion and siren audio (2026-09-13)

- Theatre now releases its sound when the cue changes or the component unmounts. Scene sounds own their active and future scheduled voices; explosion noise/thump, siren pulses, legacy gunfire and knocks stop and disconnect on dismissal. Natural endings release their graphs, and muting cancels every active scene without replaying them on unmute. Partial construction failures cancel voices that already started.
- All 109 frontend tests and production build passed. New tests exercise the actual audio functions with mocked Web Audio nodes: all supported scene kinds, future siren scheduling, repeated cancellation, independent overlapping scenes, natural endings, global mute and an oscillator failure after explosion noise starts. Existing visible city gunshot cancellation tests still pass. Logs: `.runtime/city3d-scene-audio-tests.log` and `city3d-scene-audio-build.log`.
- Browser integration checks on isolated fixtures: port 8843 recorded casualty scene removed on immediate Skip; port 8853 police scene removed on navigation to People. Both had no captured console errors and neither check issued a gameplay command. Audio graph cancellation is established by the node tests and effect cleanup; no claim of an acoustic recording or listening test is made. Room ambience/table effects have separate lifecycles and were not changed. Broader visual polish and event choreography remain unfinished.

### Coopered rooftop water tanks (2026-09-13)

- Replaced the simple iron cylinders on tenement and shop roofs with Blender-authored timber cisterns: 32 separate staves, steel hoops, cross-braced stands and support girders, conical caps/vents, and rung ladders. Deterministic embedded 256px cedar colour/normal textures add wood grain; a browser review prompted a lighter weathered wood tone. Casino and civic crowns remain tank-free.
- Regenerated only `tenement.glb` and `shop.glb` and their manifest entries. Existing facade footprint and public building/command identifiers remain unchanged. Heights increased by 0.88m; labels use actual exported geometry bounds. These are decorative roof structures, not water-management gameplay.
- All 108 frontend tests and production build passed. New actual-GLB checks verify timber stays inside its clear roof bay, sits above its support, faces outward and carries embedded textures; exact geometry matches updated manifest bounds. Existing all-building footprint/route checks passed. Front/reverse browser views are recorded in `docs/qa/city3d-20260913/rooftop-tanks-front.png` and `rooftop-tanks-rear.png`.
- Whole-city local sample on isolated port 8850: 145 FPS, 7.7ms p95, 272 draws and 490,666 triangles for 28 buildings/one actor. This is one local Mac observation, not broad hardware acceptance. Street density, more differentiated architecture, characters, event choreography and final visual-quality acceptance remain unfinished.

### Shared city scene and travel panels (2026-09-13)

- The normal city no longer positions Theatre or the travel banner outside its responsive controls. Both now occupy the same flowing story slot as expanded playback. Expanding preserves the mounted Theatre, its audio/progress, cast and headline instead of covering it with a duplicate caption. Interior Theatre remains in its existing stage.
- Story content has a bounded scroll region, wrapping cast rows and a sticky Skip/Go on control. Travel retains its saved duration, walking/driving details and crossing warning. The redundant City3D dismissal props and old fixed event offsets were removed.
- Browser replay on isolated port 8855: normal city at 1024×768 left clear space between toolbar, scene and address; expanded scene panel had no measured intersections at 320×568, 390×844, 800×600, 844×390 and 1280×720. Keyboard navigation scrolled Go on into view; dismissal removed the scene panel. Evidence: `docs/qa/city3d-20260913/scene-panel-phone.png` and `scene-panel-layout.json`. Replay did not advance the fixture's 08:00 clock.
- Isolated port 8850 travel test showed exactly one banner before and after expansion, retained inside the city story slot, then removed it on Skip. Only that test save advanced. All 107 frontend tests and production build passed. Default viewport restored; main save untouched. Extremely short portrait windows, unusual long captions and full touch accessibility remain unverified; this is not final visual-quality acceptance.

### Responsive city controls (2026-09-13)

- Heading, camera controls and expanded scene/journey actions now share a flowing grid. Container queries use the city panel width, including when the gameplay sidebar is present. Wrapped instructions and buttons push subsequent actions down instead of depending on fixed pixel offsets. Short expanded landscape views omit the visible shortcut hints and address description; canvas accessibility instructions and the actions remain available.
- Final production build passed. Actual browser panel rectangles at 320×568, 390×844, 800×600, 844×390 and 1280×720 show no intersections or out-of-viewport panels in idle expanded state. Evidence: `docs/qa/city3d-20260913/responsive-layout.json` and `responsive-phone.png`. Travel scene actions were additionally inspected at 320×568 and 844×390 with clear separation from the destination panel. Isolated port 8850 walk advanced only its test save to minute 690; main save untouched. Default browser viewport restored.
- Embedded desktop layout was visually inspected and keyboard Escape/Enter still closes/reopens the city. This check does not cover every possible long event caption, browser text scale, touch gesture or extremely short window. Collapsed-city external Theatre placement still needs its own narrow-panel event audit. Broader visual quality acceptance remains open.

### Pedestrian road elevation (2026-09-13)

- Moving pedestrians now follow asphalt and pavement elevations instead of staying at pavement height across every road. A 40cm smoothstep on the road side of the curb provides a continuous rise/descent; a conservative 68cm stride radius encloses both articulated models and completes the rise before their shoes reach the raised surface. Existing distance-driven gait remains intact.
- The player location ring follows the actual surface independently of the character's stride support. Go journey duration, routing and traffic progress are unchanged.
- Isolated port 8850 browser walk Saint Agnes → The Mariner completed at revision 6/minute 675, progress 1 after 28,116ms of presentation. Sampled walking root ranged from −0.01976m on asphalt to 0.265m on pavement, and settled at 0.2m at (48,36.65). Evidence: `docs/qa/city3d-20260913/walking-ground-samples.json` and `walking-road-contact.png`. The screenshot is a distant scene view; the numeric samples provide the stronger height evidence.
- All 107 frontend tests and production build pass. New geometry checks sweep both exported animated bounds across positive/negative X/Z road edges and camera-independent character headings, asserting surface clearance and continuous bounded height. This is conservative body support, not individual foot planting or authored curb-step animation; those and broader architectural/character polish remain unfinished.

### Keyboard camera and expanded-city focus (2026-09-13)

- Arrow-key pan now follows the camera's horizontal screen axes after rotation. Modified browser/OS shortcuts and composing input are ignored by camera handling. Visible keyboard instructions cover rotate, pan, zoom, reset and return.
- Expanding focuses the canvas; Shift+Tab wraps to the address directory, Tab returns to the canvas, and Escape closes the expanded city and restores focus to its button. Camera diagnostics expose position, target and zoom for read-only browser verification.
- Isolated port 8847 browser check: Enter opened/focused the canvas; Q then Up moved target from (96,80) to (93.5004,84.3304); + changed zoom 1.65 to 1.815; Home restored zoom 1 and target (96,80). Both focus-wrap directions and Escape/Enter reopening passed. Screenshot: `docs/qa/city3d-20260913/keyboard-navigation.png`. No gameplay command or save mutation during this check.
- All 106 frontend tests, production build and HTTP package tests passed. Unit coverage checks modifier/IME handling and pan direction/distance through a complete camera orbit. Desktop layout inspected at 1280×720; new mobile hint spacing still needs a browser viewport check. This is a navigation pass, not production-quality visual acceptance; architectural density, character polish, curb stepping and wider device/performance checks remain.


### Thorne & Sons architectural follow-up

The current location data describes an undertaker with a brass plate, long empty
window and rear hearse yard; its older church silhouette did not fit. Thorne &
Sons now uses a dedicated two-storey brick funeral premises with sandstone trim,
a long framed display window, oak entry, brass plate/pulls, upstairs sashes, slate
gables and chimneys. A gated rear coach yard includes a static period hearse,
authored by extending the existing Packard model in Blender. This is decorative
business scenery, not an authoritative travelling NPC or a new purchasable car.
The old chapel model remains on disk but is no longer loaded for this address.
There are 27 exported model files, with 26 in the active renderer catalog.

The yard has locally generated periodic gravel colour and normal maps embedded
in the GLB, with physical-scale UVs. The new asset is 798,024 bytes and its exact
Blender bounds are (-7.7,-7.715,-0.03) to (7.7,8.01,8.87), within the reserved
parcel. Export bounds now use actual transformed vertices rather than rotated
local bounding-box corners, which had overstated this roof's height. Runtime
building labels likewise use precise geometry bounds. Other existing model files
were preserved during the focused export.

Evidence: 104 frontend tests, build and HTTP package pass. The active model is
covered by the existing parcel/route checks, actual GLB raycasts at its front door
and across the 2.2m rear gate corridor, exact manifest/geometry agreement and
embedded colour/normal texture checks. The final denser gateway scan also passed
as a targeted geometry rerun. Logs: `.runtime/city3d-undertaker-tests.log`,
`.runtime/city3d-undertaker-geometry-tests.log`, `.runtime/city3d-build.log`,
`.runtime/city3d-undertaker-go.log`, `.runtime/city3d-undertaker-export.log`.

Browser front and rear inspection used the existing isolated fixture on 8847 and
left revision 3 / minute 773 unchanged. Screenshots:
`qa/city3d-20260913/undertaker-front.png` and `undertaker-yard.png`. A whole-city
sample measured 145 FPS / 7.1ms p95 / 271 draws / 440,000 triangles with 28 buildings
and two actor objects. That remains a local sparse-traffic sample, not broad
performance acceptance. More location-specific architecture, richer character
surfaces, pedestrian kerb transitions and full accessibility/resource/device
acceptance remain unfinished; the goal remains active.


### Street surfaces and vehicle grounding follow-up

Added a locally authored Blender street-bed set: individually jointed kerbstones
with occasional replacement stones, cast-iron manhole covers and recessed drain
grilles. Four shared material groups are instanced across the 28 occupied lots;
the new GLB is 743,712 bytes. Kerb corners are trimmed to avoid overlapping stone
runs. Surface detail stays below pedestrian sole clearance, with manholes/drains
outside the sampled tyre tracks. There are now 26 authored GLBs.

Moving cars previously inherited the pedestrian root height, leaving tyres about
32cm above the asphalt. Vehicle roots now follow asphalt (-0.1m) or raised pavement
(0.17m), accounting for the exported 2cm tyre offset and a 5mm surface gap. Parked
cars and staged police use the same helper. A shared soft contact-shadow texture
now follows each vehicle; the selection ring and police beacon follow the adjusted
height. Pedestrian kerb-crossing height transitions still need authored step motion.

A close-up exposed a separate visible defect at Thorne & Sons: the chapel door
was buried behind the tower wall, despite the renderer's correct façade rotation.
The Blender doorway now sits in front of that wall, with oak panels, stone jambs,
threshold and brass pulls. A raycast against the actual exported GLB verifies the
visible wooden entry at three positions. The broader building still needs richer,
location-specific architecture; fixing the missing door is not full art acceptance.

Evidence: 103 frontend tests pass, including sampled current driving routes,
parked-vehicle height, kerb/sole clearance, tyre-track/manhole/drain separation and
the GLB entry raycasts. Build and HTTP package pass. Logs:
`.runtime/city3d-surfaces-tests.log`, `.runtime/city3d-build.log`,
`.runtime/city3d-surfaces-go.log` and `.runtime/city3d-street-bed-export.log`.
Existing unchanged model binaries were preserved after the full export.

Browser on the existing isolated traffic fixture 8847 showed the parked Hudson
at y=0.155, and a drive to Thorne & Sons at y=-0.115, revision 3 / minute 773. The
corrected doorway, contact shadow and final street surfaces are captured in
`qa/city3d-20260913/street-surfaces.png`. Whole-city view measured 145 FPS / 8ms p95,
263 draws and 436,064 triangles with 28 buildings and two actor objects. This is a
local sparse-traffic sample; heavy traffic and broader device acceptance remain
pending. The main campaign save was not opened or modified.


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
