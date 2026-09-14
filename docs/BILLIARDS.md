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
