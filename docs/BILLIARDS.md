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

## Required next work

1. Match rules driven by the actual ordered collision/pocket events: legal break,
   solids/stripes, first contact, rail/pocket requirement, ball in hand, called
   shots and eight-ball victory/loss. Use explicit posted rules, including break
   exceptions. The current official
   [WPA rulebook](https://www.wpapool.com/wp-content/uploads/2026/01/2026.01.02-WPA-Rules.pdf)
   is the reference; distinguish any deliberate house variation.
2. Saved matches and funded individual stakes through exactly-once Go commands;
   no renderer-owned results, hidden odds substitution or minted payouts.
3. A close 3D table with aiming, power, tip position, ball placement, numbered
   rotating balls, cue motion and impact/pocket sound. The hall's initial table
   props have stylized proportions; align the playable model and the six hall
   tables with the solver's 2:1 cloth dimensions and actual ball radius.
4. Opponents must choose and execute physical shots. Tournament scheduling,
   entrants' entry payments, brackets and whole-pool settlement must use real
   participant funds and survive saves/retries/interruption.
5. Airborne/jump/masse and slate impacts are not implemented: this is currently
   a ground-contact solver with 3D angular velocity. It also approximates cushion
   contact at ball-centre height. These are fidelity limitations, not grounds to
   declare the full minigame complete. Calibrate with rendered playtests and
   extend the model as needed.

No live campaign was used for physics QA. Main port 8791 still serves release
16c1c95; tailor/hall source commits await a verified promotion.
