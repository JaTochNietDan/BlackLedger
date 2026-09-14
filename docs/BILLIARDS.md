# Playable billiards work

The user requires a full-fidelity billiards minigame, individual wagers and
occasional entry-fee tournaments paying the entire pool to the winner. A static
hall or a probabilistic win button does not meet that requirement.

## Shot physics implemented

`billiards` is a deterministic Go package without campaign, HTTP, renderer or RNG
dependencies. It does not currently expose a player action or alter a save.

- Metre-based 1.27 × 2.54 m cloth and 28.575 mm ball radius; six pocket mouths,
  straight cushions, angled pocket facings and rounded endpoint contacts.
- Continuous swept collision times within 1/2400-second friction steps. The
  timestep handles cloth integration, not a discrete overlap-only collision
  approximation. Tied contacts resolve deterministically by sorted ball number
  and authored geometry order.
- Separate sliding, rolling and torsional friction. Slip at the cloth drives both
  velocity and angular velocity, using solid-sphere inertia. Cue contact offset
  produces draw/follow/side spin; impacts use restitution and bounded tangential
  friction, including side-spin transfer at cushions.
- Cue speed 0.05–8 m/s and contact offset at most 0.6 radii. Non-finite inputs,
  overlapping balls, moving tables and pathological configurations are rejected.
  Every shot must settle within a bounded duration; failure returns an error
  rather than silently freezing a still-moving table.
- Replay contains unit quaternion orientation, regular 30 Hz frames and a frame
  at every impact. Renderers must preserve these impact times rather than
  interpolate directly across rebounds. Input arrays are not mutated.
- Eight-ball rack geometry and cue-placement validation are ready for match
  integration. Number order is fixed in physics; any rack randomization belongs
  to the seeded match layer.

The model uses standard rigid-body impulse and sliding/rolling relationships.
Useful primary references consulted: Evan Kiefl's
[physics derivation](https://ekiefl.github.io/2020/04/24/pooltool-theory/) and
[pooltool physics documentation](https://pooltool.readthedocs.io/en/latest/autoapi/pooltool/physics/index.html).
The code is authored here, not copied from pooltool. Friction/restitution values
are initial calibration choices, not measurements of a specific real table.

## Evidence — 2026-09-14

All 15 physics tests pass (`.runtime/billiards-physics-final.log`, 2.531 s).
They check analytic sliding-to-rolling velocity and travel distance, elastic
momentum/energy conservation, dissipative frictional impacts, draw/follow,
side-spin rail rebounds, continuous grazing contacts, six pocket entries, jaw
clipping, spin decay, input immutability, deterministic ordering and bounds.
Thirty-six power/aim/spin combinations exercise complete racks, checking every
replay frame for finite state, table containment and interpenetration. Every
impact must have a replay frame. Halving the production timestep preserves the
chosen cut-shot event sequence and final positions within 0.1 mm.

The first convergence run at 1/600 s moved a final ball by 0.262 mm when halved;
the production interval was reduced to 1/2400 s, retaining the original 0.1 mm
check. This establishes numerical convergence for that fixture, not experimental
validation of all physical coefficients.

The 13-test suite before adding the last jaw/configuration checks also passed
with the race detector (15.370 s). `go vet ./billiards` passed. Full-break benchmark
on Apple M4 Max: 61.90 ms/shot, 517741 bytes and 266 allocations, three iterations
(`.runtime/billiards-physics-bench.log`). This is backend timing only.

## Match adjudication implemented — 2026-09-14

`Match.Play` accepts a seat, physical cue input and a declared ball/pocket (or
safety). It runs the solver itself; callers cannot submit a winning result.
The private adjudicator uses ordered cue contacts, rail impacts and pocket
entries. It supports open-table group assignment only on a legal called pot,
wrong-group/no-contact/no-rail/scratch fouls, continued turns, safeties, eight-ball
wins and premature/wrong-pocket/foul losses. A scratch while shooting at the
eight does not itself lose the rack unless the eight also falls.

Breaks keep groups open, count distinct object balls reaching cushions, and
preserve the appropriate player's choice after an illegal break or a pocketed
eight. Decisions support accepting the position, spotting the eight, taking
cue ball behind the head string and re-racking with the selected breaker.
Spotting finds free space on the long string rather than overlapping another
ball. Head-string placement and actual pre-contact travel are checked. Frozen
cushions count only after the ball has departed; simultaneous legal first
contact is not discarded because another numbered ball is visited first.
Concession records a winner without implementing any money payout itself.

All 33 package tests pass (2.593 s) and `go vet ./billiards` passes. The 18 match
tests cover break decisions, spotting, head-string paths, frozen cushions,
called combinations, early/scratched eights, concession and JSON round trips.
One test restores a serialized match, executes a real physical called-eight
shot, and verifies the win. Another executes a complete physical break through
`Match.Play`. Invalid actions and failed physics leave the match unchanged.
These are package/serialization tests, not HTTP or campaign-save integration.

The [WPA rules](https://wpapool.com/wp-content/uploads/2025/10/2025.09.15-WPA-Rules-NP.pdf)
are the reference for called eight-ball and break options. The implementation
automatically spots the nearest eligible ball when every target lies behind
the head string. The opening breaker is currently supplied by the caller; lag
play is not implemented. Unsupported physical events (jumping off the table,
push shots and equipment/stance fouls) are not claimed to be adjudicated.

## Campaign stake foundation — 2026-09-14

The optional saved `World.Pool` now owns the opponent, life, $10–$500 stake per
player, combined escrow, match, settlement/void flags and latest replay. Starting
checks the actual location, availability and both purses before taking either
stake. A physical win, concession or interruption pays the held amount once;
only profit adds to player earnings. Fire or an unusable hall refunds both
original stakes. A failed shot or placement does not change the ledger.

The active opponent does not leave through an ordinary routine change and
cannot be pruned before settlement. Advance/death/new-life paths reconcile the
stake. Corrupted old-life winning records are not credited to a new protagonist;
a missing winning opponent's escrow stays held rather than being invented or
redirected. Normal lifecycle hooks settle before changing lives.

Replays use base64/zlib JSON with compact pose arrays; they omit velocity and
repeated ball field names but retain all impacts and quaternion orientation.
The tested break encodes to 50388 bytes for 235 frames / 39 events. Round-trip
position and orientation errors stay below 1e-6. Decoder sizes are bounded.

Evidence: all 35 billiards tests pass (2.663 s); all 11 core pool tests pass
(0.244 s); the selected departure/death/new-life/card-seating regression tests
pass (0.262 s). Store tests pass (0.227 s), including a real temporary SQLite
close/reopen between reserving stakes and physically winning, then a second
reopen and repeated reconciliation. `go vet ./billiards ./core ./store` passes.
Only disposable fixtures were used. This proves persistence and settlement
idempotence in the core/store-change path, not request-ID receipt handling:
player-facing command routes and the table UI are still to be connected.

## Physical opponent play — 2026-09-14

`billiards.Opponent` prepares break decisions and ball-in-hand placement, then
plays one actual stroke. Geometry generates direct pots and one-cushion banks;
blocked positions also generate direct safeties and cushion escapes. Candidate
strokes are evaluated through `Match.Play`. Normal-shot previews use 1/600 s
cloth steps; the executed shot always uses the production solver. Skill affects
cue-angle and speed error, not a post-hoc win probability or pocket override.

Breaks evaluate seven physical aims. The initial fixed break was not robust:
perfect rack geometry frequently sent only three distinct object balls to rails,
and one self-play run repeated illegal breaks to its 160-shot limit. Searching
alternative impacts solved this without changing physics, rack layout or the
four-ball rule. Cushion escapes also needed more power and a bounded angle
search to account for dissipative rail rebounds. Safeties are skipped only when
their maximum possible score cannot improve the current evaluated result.

`World.PlayPoolOpponent` derives stable skill and execution seeds from opponent
identity, life and shot number. It cannot act for the player or an absent NPC.
It stores the actual intent, placement/decision and compressed replay, then uses
the existing funded settlement. Planning does not consume either campaign RNG.
The saved optional `last_stroke` also records the player's cue input for future
cue animation. No HTTP or UI path is exposed yet.

Evidence: all 42 billiards tests pass (10.311 s); all 13 core pool tests pass
(0.213 s), store tests pass (0.213 s), and vet passes. The tests re-execute an
advertised bot stroke and require identical physical output, restore and retry
with identical intent/replay, reject wrong-seat/absent/finished play, exercise a
blocked cushion escape, and check eight first-break variations. Full physical
self-play at skill 0.9 finishes seeds 7 and 41 after 18 and 17 strokes, both by a
legally called eight-ball win. The core fixture verifies the NPC's physical win
pays the actual held stake and leaves both world RNGs unchanged.

Measured mid-rack planning/execution on Apple M4 Max, three iterations: 110.66 ms,
2625434 bytes and 1603 allocations per stroke. Before eliminating provably
uncompetitive safety evaluations it was 863.80 ms / 25651618 bytes. These are
backend benchmarks, not browser frame-time measurements. Logs are under
`.runtime/pool-opponent-*`.

## Required next work

1. Connect match rules to campaign commands and post the rules in the playable
   view, including break decisions and any deliberate house variation. Extend
   adjudication as unsupported physical events become available.
2. Keep the implemented command/receipt path as the single authority while
   connecting the browser. Test browser request interruption and restored playback;
   store/HTTP retries and stale revisions now have coverage below.
3. A close 3D table with aiming, power, tip position, ball placement, numbered
   rotating balls, cue motion and impact/pocket sound. The hall's initial table
   props have stylized proportions; align the playable model and the six hall
   tables with the solver's 2:1 cloth dimensions and actual ball radius.
4. Connect the implemented physical opponents to the playable view. Tournament scheduling,
   entrants' entry payments, brackets and whole-pool settlement must use real
   participant funds and survive saves/retries/interruption.
5. Airborne/jump/masse and slate impacts are not implemented: this is currently
   a ground-contact solver with 3D angular velocity. It also approximates cushion
   contact at ball-centre height. These are fidelity limitations, not grounds to
   declare the full minigame complete. Calibrate with rendered playtests and
   extend the model as needed.

No live campaign was used for physics QA. Main port 8791 still serves release
16c1c95; tailor/hall source commits await a verified promotion.

## Command and public-state integration — 2026-09-14

Added the intent-only `pool_*` commands and an explicit local table DTO, documented
in API.md. Commands reserve stakes, place/shoot, play the NPC, resolve breaks,
concede and dismiss finished racks. A live wager requires finish/concession before
switching activities; appearance changes and urgent event decisions remain usable.
The public view copies positions/quaternions, legal targets, decisions, exact cue
intent and the compressed replay without consuming RNG or altering the match.
No UI challenge button is exposed until playable controls are integrated.

Temporary-SQLite tests exercise concurrent same-ID start retries, a physical
winning shot, database close/reopen, byte-identical saved shot receipts, stale
new-ID rejection, exactly one shot/payout and unchanged persisted state. HTTP
checks start/place/break, identical retries, complete public replay, malformed
numeric cue input, missing payload and forged NPC intent. The existing null-shape
guard initially rejected the new nullable `pool`; its contract now explicitly
permits absent rack/previous stroke while the HTTP rack checks verify lists stay
lists. No live campaign was used or promoted.

Final evidence: selected core pool/card/travel/lifecycle tests pass (0.333 s),
full store tests pass (0.271 s), full HTTP/server tests pass (1.286 s), and vet
passes. The expanded HTTP shape check also found null travel/effect lists in
command receipts; command initialization now publishes empty lists. Logs are
`.runtime/pool-api-{core-verified,adapters-verified,vet-final}.log`.

## First playable 3D client — 2026-09-14

`BilliardsRoom` offers real local opponents and held stakes through the existing
command client. A map-sized table has placement, direction, power, top/side spin,
ball/pocket calls, safety, break choices, NPC turns and explicit concession/close.
Controls lock during physical replay; reloading shows the saved final position.
Motion preferences and a skip control apply only to presentation. A late decode
cannot overwrite a newer or skipped stroke. Errors remain above the table overlay.

The Three.js table uses solver cloth dimensions and ball radius, six real mesh
holes, cushions with the solver's segment/jaw coordinates, rotating numbered balls
and a visible cue/direction line. Browser review changed its default camera from
lengthwise to across the table for better use of a wide screen. Orbit, pan, zoom
and reset remain available. This is still a rough table asset: cushion profiles,
wood/pocket detail, surrounding room/players, physical cue stroke and matching hall
props need further work. Numeric placement/aim controls need pointer interaction;
compact layout, motion-off/skip and complete rack UX need further browser review.

All 378 frontend tests pass (14.879s), including zlib decoding, malformed/oversized
replay rejection and preservation of adjacent collision samples. Final TypeScript
and production build pass (existing bundle-size warning). Isolated browser port
8964 uses `.runtime/pool-playable-qa.sqlite3` only. Actual clicks reserved $20 each,
placed the cue, played a missed/illegal player break, then an NPC `hand-behind`
decision and physical safety. The browser displayed moving balls and locked turns;
restoring the page showed the final table. Concession and return left revision 6,
minute 606, cash $5,980 and no pool session. No second stake or payout occurred.

The actual NPC receipt decodes in the browser-compatible module to 324 frames,
156 events, 5.5875 seconds; every final ball position matches the public projection
within 1e-6 metres and all pocket statuses match. Evidence is in
`.runtime/pool-browser-receipt.json` and `.runtime/pool-view-*` logs. Main port 8791
and the live campaign are unchanged. This does not close full-fidelity acceptance.

## Direct table interaction — 2026-09-14

Cloth clicks now set aim, or preview cue placement while in hand. Placement stays
local until the confirm button sends `pool_place`; a translucent cue and dashed
head string show the proposed position and restriction. Invalid previews turn red
and cannot be confirmed. Numeric placement remains under Precise placement.
Clicking a visible ball sets aim toward its centre and calls it if legal; clicking
a numbered pocket selects the matching call. Selected balls/pockets have rings.
These inputs do not predict pot success or replace authoritative Go validation.

A six-pixel movement threshold distinguishes taps from camera drags and retains
movement even if a gesture returns to its start. Secondary buttons, cancellation
and multiple pointers cannot become a cue tap. Orbit/pan/zoom remain available.
Pointer-selected coordinates display three decimals and aim two decimals; the
posted cue inputs use those displayed values. Pocket labels 1–6 map to API 0–5.

All 380 frontend tests pass (14.044s). The final five billiards tests pass after
adding the placement clearance tolerance; final typecheck/build pass (existing
bundle-size warning). Tests cover gesture return-to-origin, multi-pointer/right
button cancellation, aim axes, occupied placement and the strict head boundary.

CUA tab53 on isolated port8964 started a second $20-per-player rack. Two placement
previews left revision7/minute608/ball-in-hand unchanged. Clicking beyond the head
string disabled confirmation. A valid click/confirmation enabled aiming; clicking
the head ball set 89.9839°, and orbit dragging preserved it. The resulting real
break was legal. Leo physically potted a called stripe, retained the turn, then
passed. Clicking the visible blue ball selected 2; clicking pocket marker6 selected
Right middle and retained the ball call. State remained revision11/minute614,
three shots, $5960 player cash and $40 escrow after these local selections.
The replay skip control was observed, but the shot finished before the attempted
click; skip execution and compact/touch browser acceptance remain unverified.
Main campaign untouched. Logs: `.runtime/pool-pointer-{tests,final-tests,
verified-build}.log`. Table/character detail, cue animation, concurrent tables,
tournaments and the retained physics limitations remain open.

## Visible cue strokes — 2026-09-14

Saved player and NPC intent now drives a draw-back, accelerating strike,
follow-through and withdrawal before/alongside the physical replay. The front
of the cue reaches the sphere at the actual top/side contact offset. Cloth
playback begins at contact after a 0.72-second presentation preparation; the
server's frame times and outcomes remain unchanged. The cue fades out by1.12s.
Older replay records without saved intent still play without invented cue input.
A tapered shaft, contrasting butt, ivory-coloured ferrule and chalked tip replace
the original single cylinder. Idle cue placement also reflects selected spin.

Skip, motion-off and stale-decode cancellation apply to the complete stroke and
ball sequence. Contact audio is emitted only when its phase is reached; skipped
or substantially overdue contact is not sounded later. No character pose or
hand/bridge animation is implied: these are still missing.

All382 frontend tests pass (14.120s); the final seven billiards tests pass after
removing a redundant equal-input assertion. They check off-centre tip/sphere
contact, no ball motion before contact, phase continuity and finish visibility.
Final typecheck/build passes with the existing bundle warning. CUA54 at isolated
8964 captured the player cue during committed playback and the NPC cue on its
separate saved direction. Clicking Skip ball motion succeeded while active,
removed the animation and enabled the next-turn control. The subsequent NPC shot
used the ordinary command path. Final fixture revision13/minute618/five shots,
$5960 cash/$40 escrow, exactly13 receipts: skipping added no command or payment.
Evidence logs: `.runtime/pool-cue-{tests,final-tests,final-build}.log`.
Main campaign unchanged. Character choreography, richer table/hall assets,
compact/touch/motion-off acceptance, multiple tables and tournaments remain open.

## Cushion contact geometry and table finish — 2026-09-14

The original centred cushion boxes intruded12.5mm into the playing area. New
closed profiles put their noses exactly on the solver's segment at ball-centre
height, with the remaining material on the outward side. All18 segments include
the short angled pocket facings. Raised wood rails now carry visible sights;
locally generated grain, contrasting cushion slopes, leather-coloured pocket lips
and recessed dark cups replace the thin rail/flat-pocket appearance. Aprons are
split around pocket openings so wood no longer fills their interiors.

All384 frontend tests pass (14.076s). Geometry checks cover every profile vertex
against its contact half-plane, nose height and a raycast against the end cushion.
Final typecheck/build passes (existing bundle warning). CUA55 on isolated8964
visually inspected the revised table without sending a gameplay command; the
fixture remains revision13/minute618. Main campaign unchanged. Evidence logs:
`.runtime/pool-cushion-{tests,final-build}.log`.

This improves the playable table, not the hall's six existing props. Those still
need matching dimensions/finish. Pocket baskets, joined corner detailing, player
bodies/hands, scene surroundings and complete visual acceptance remain unfinished;
the solver's previously documented physical limitations are also unchanged.

## Six hall tables brought to physical scale — 2026-09-14

Regenerated the original Blender hall tables with1.27×2.54m cloth dimensions,
0.78m playing height and0.028575m ball radius. Previously they used2.18×2.70m
cloth,1.136m height and0.053m ball radius. Each table now has the same18 cushion
segments/profiles as the playable view, raised walnut rails/sights, split aprons,
open leather cups and recessed pocket bottoms. Feet and legs fit the lower table.
The six tables retain their room positions and existing spectator/staff staging.
Their five decorative balls are scenery, not concurrent simulated matches.

The targeted hall geometry/clearance test passes (3.684s), including all36 pocket
openings, cloth height and a ball-size probe on each table,24 cushion-dimension
ray checks, both rigs' standing/seated positions and sampled entrance paths.
Final production build/typecheck passes with the existing bundle warning. CUA56
at isolated8965 inspected standard and enlarged/zoomed views with real fixture
occupants. No game command was sent. Fixture `.runtime/poolhall-scale-qa.sqlite3`;
main campaign unchanged. Export/test/build logs use `.runtime/poolhall-dimensions-*`.

The asset grew from3,425,888 to4,961,068 bytes and61,584 to83,312 mesh triangles.
These are GLB geometry totals, not draw-call/FPS measurements; fresh runtime
performance acceptance remains needed. Joined corner/leather detail, character
play/choreography, multiple live tables, tournament play and broader interior
acceptance remain open.

## Saved tournament bracket foundation — 2026-09-14

`billiards.Bracket` stores two-, four- or eight-entrant single-elimination events,
with actual `Match` instances, stable entrant identities, round pairings and table
allocations. Advancement consumes the embedded match's adjudicated winner or its
explicit concession result; it never rolls a separate tournament win chance.
Later rounds wait for both actual winners. Repeated advancement is idempotent.
Initial seed0 stays in seat0 if successful, matching the existing player-first
minigame controls. Later ready rounds choose a genuinely free table, since they
can overlap unfinished opening racks; no two active bracket games share a table.

Tests cover partial-round JSON save/reopen, waiting entrants, all successive
pairings, a completed bracket's repeated reconciliation, invalid/duplicate seeds,
non-aliasing entrant lists, and an early semifinal alongside occupied opening
tables. A physical two-entrant event finishes in13 actual strokes with Leo winning
by a legally called eight. The full billiards suite passes (11.576s) and vet
passes; logs `.runtime/pool-bracket-{suite,tests,final-tests,vet}.log`.

This is internal bracket state only: it is not yet attached to campaign money or
exposed as a playable tournament. Required next integration is scheduled entry,
funded participant deposits, retention/forfeit rules on death/departure/closure,
whole-prize-pool settlement exactly once, command receipts, public bracket and
browser progression. Individual casual matches must also respect tournament table
occupancy when those two systems are joined. No campaign save was used.

## Tournament withdrawal and empty-branch handling — 2026-09-14

Before attaching entry-fee escrow, the bracket now supports explicit withdrawal
and simultaneous withdrawal batches. Active opponents receive a concession;
entrants waiting between rounds cannot return after withdrawing. If both sides
of a branch withdraw, it resolves without a winner and supplies a bye downstream.
An all-withdrawn event has `finished:true` and an empty champion, distinct from an
unfinished bracket. Finished events cannot be rewritten by later departures.

Batches validate every identity before changing eligibility, then advance once.
This is necessary for a shared incident: withdrawing casualties one at a time
could incorrectly crown the last casualty before processing their withdrawal.
Resolved empty branches release their tables even if their historical rack had
not physically finished. Consumers must use the bracket's resolved/finished
status, not infer event progress from rack winner alone.

The full billiards suite passes (11.623s), including saved waiting entrants,
simultaneous empty branches, all-entrant withdrawal, atomic invalid batches,
physical concessions and immutable finished champions. Vet passes. Logs:
`.runtime/pool-withdrawal-{tests,vet}.log`. Campaign entry fees and payouts were
not connected in this increment; this closes a lifecycle gap needed before money
can safely be attached. Core/HTTP/UI integration and scheduling remain open.
No campaign save was used.
