# Simulation boundary v1

The simulation is authoritative. UI holds only selection, view, animation and voice preferences. It never computes rules or owns a campaign save.

- GET /api/state: versioned public projection (no RNG, hidden plots or unobserved outcomes).
- POST /api/action: {request_id, revision, kind, target?, event?, choice?}. Exactly-once transactional command. Stale revisions are rejected. Response is the new public projection with committed presentation events.
- POST /api/director: ask the asynchronous director to prepare a bounded proposal. Does not advance game time.
- POST /api/speech: {event}; speak only the active scene's saved text, with a stable voice by character identity. Does not advance game time.

Core is a Go package with no browser, rendering, HTTP or model-runtime dependencies. HTTP/save adapter persists complete state and idempotency receipts together in SQLite. AI and Kokoro are external services, not authorities. React is a replaceable command client. The active isolated Canvas2D street only renders committed events and decorative motion; the earlier PixiJS renderer is not active. Skipping presentation cannot affect outcomes. Headless clients use the same command API.

A migration/reference prototype exists under prototype-python; do not run it as the active game backend. Its original tests establish expected command behavior, to be ported into Go core/HTTP tests.

## Current presentation and director contracts

`last_result` contains `from_location`, `to_location`, elapsed game minutes, and the records produced by the committed command. Record IDs remain stable as the bounded history rolls over. Street travel is a skippable presentation of these endpoints; it never advances the backend clock. The optional `opportunity` is a suggestion derived from public progress, not a disclosure of hidden plots.

Director proposals name an existing speaker, one supported operation (`courier`, `mediation`, or `collection`), and optionally an existing faction ID as `beneficiary`. Go assigns the operation's reward, duration, heat and respect; validates the beneficiary; supplies canonical completion text; and displays the stakes in the choice. Completion gives that faction +6 standing and rivals −3. Rejection or an abandoned police stop grants neither the reward nor faction credit. Previously saved offers retain their promised terms.

At a resulting heat of 15 or more, a completed arrangement pauses at an authored police stop before paying its reward. The saved scene preserves the original arrangement and beneficiary until the player pays to complete it or abandons it. Hidden retaliation scheduling remains private. AI prose is bounded and prompted against inventing outcomes, but is not a complete semantic truth guarantee.

## Save schema v2

The player stores `best_home` for one-time housing-tier progression. Estate ownership uses the same life-specific property owner as businesses, while apartments remain rentals. Loading a v1 save restores a living estate resident's previously omitted deed and initializes housing progression from the current residence. Public command revisions and game time are unchanged by migration; the next save persists the upgraded schema.

## Contextual approaches in AI offers

A proposal may include up to two `approaches`, each with a `method` and short player-facing `label`. Supported methods are `careful` (+30 minutes, −$15 reward, up to 3 less heat) and `press` (15 fewer minutes, +$20, +5 heat). Unknown/duplicate methods and oversized labels are rejected. Standard acceptance and refusal remain available. The saved scene contains authoritative effects for each offered approach; forged choices are rejected. Police interruptions preserve the selected terms, while urgent danger interrupting the work grants no completion reward or faction credit. Business demands pause work for a resume-or-abandon decision. Omitted approaches preserve compatibility with existing offers.

## Director contact progression

New AI requests use established contacts: Mara initially; recruited associates with at least 30 loyalty; family leaders once their family has at least +6 goodwill. Among these, the server prefers the least recently featured speaker, counting current-life arrangement memory and pending offers. Saved continuations preserve their original speaker. The prompt and decoding schema receive this allowed cast, and server validation rejects other speakers with one bounded correction attempt. Other NPC context is not an invitation to impersonate an unavailable contact. This governs new request generation, not audiences or already saved scenes, and does not guarantee semantic correctness of all dialogue.

## Located AI jobs

New model responses require `location`, an existing place ID in an unlocked district. The schema restricts IDs; server validation requires the selected venue's name in the body and rejects explicit locked venue names in title/body/approaches. This is a bounded consistency check, not a semantic proof of the rest of the story. Core validation also rejects supplied inaccessible locations. Authored and legacy proposals may omit location.

Located proposals store the venue in scene `target` and arrangement memory `location`. Action details display its canonical name. Jobs still resolve off-screen within their quoted total duration and return to the player's current base; they do not teleport the player, unlock districts or grant free property access. Reading and rendering remain independent of simulation time. Saved follow-up briefs prefer structured venue memory over place names inferred from old prose.

## Business ceasefires

Audiences now offer a $100 business ceasefire with the represented family. It lasts 1,440 game minutes from purchase (renewals replace the expiry rather than stacking it), removes that family's pending sabotage, and suppresses its business demands and sabotage while active. Personal hits, other families, goodwill and ownership are unaffected. Payment transfers to the family through the ordinary transactional decision path. Agreements expire at the exact saved minute and are cleared on a new life.

`business_truces` is an optional public map of faction ID to active expiry minute. Empty/expired agreements are omitted from that map; private plans are never included. The Families panel shows active terms and the director receives the same public agreement context. This audience option does not authorize arbitrary model-written treaties or territorial divisions.

Routine queued/authored offers wait while discovered current-life threats remain active. They stay saved and become eligible again after the danger clears; hidden plots do not affect this pacing rule. Player-initiated audiences remain available during danger. This changes encounter delivery, not proposal generation or the simulation clock.

## Work paused by business demands

A business demand during an accepted arrangement saves its original scene, selected approach and remaining minutes in `suspended_job`. Resolving the demand opens an authored `resume_job` decision. Resume spends only the remainder, preserving the original reward/heat/standing and applying normal police checks; abandon spends no additional time and grants no reward. Repeated demands can pause the remainder again. A demand exactly at job completion still defers payment until the player chooses resume (zero remaining minutes).

No other activity is committed behind this decision. Routine offers wait while work is suspended. Attacks, urgent warnings and death still fail the operation; this is not a general guarantee of successful work. Death/new life clear suspended work. Older saves without this optional field retain existing behavior. Story memory uses `paused` until resumed or abandoned.

## Generated speaker affiliation and choice script

For new AI jobs, family leaders must name their own family as beneficiary. Independent contacts and crew may bring neutral or either-family work. The generation context supplies allowed speaker/beneficiary pairings; a single eligible leader also narrows the decoding schema. Server validation rejects contradictory pairings and requests a bounded correction. Existing saved offers retain their terms; an inconsistent old completed pairing cannot force a new follow-up. This does not simulate secret betrayals.

The English interface rejects new model approach labels containing non-Latin letters (including the mixed English/Chinese label found in QA), while allowing accents and punctuation. This is a script check, not full language detection or semantic validation.

Before committing a generated offer, the server rechecks property ownership and its speaker's identity, permitted affiliation and eligibility against the current save inside the transaction. A canonical continuation may survive lower goodwill, but its saved completed result must still exist unchanged. Time, money and ordinary property damage may progress while a draft is prepared. A stale draft is discarded without retrying the same obsolete snapshot or undoing player actions; the director becomes available again (or ready if another offer is queued). This is a preparation-time check, not semantic validation or retroactive revalidation of previously queued offers. No public fields change.

New non-neutral generated offers must also mention their beneficiary's full name or short family ID in spoken dialogue. This prevents a reward for one family when only its rival is named. It is a minimum lexical consistency guard, not semantic verification of motives or history. Existing offers are unchanged.

Director input now attributes memory to its historical participant. Current-person arrangements/history are separate from previous-person city history. The callback's original dialogue is labeled as unverified request claims, alongside the saved status/result. Allowed speakers have an explicit relationship to the current addressee. This changes private generation context only: saved memories, public DTOs, hidden plans, and outcome authority are unchanged.

Generated titles, spoken bodies and approach labels reject explicit monetary amounts, durations and recognizable deadline phrases before queuing. A failed draft receives at most one correction; repeated failure leaves the campaign clock/cash unchanged and queues no invalid offer. These lexical checks are not complete semantic verification. Authored choice details still supply actual terms; existing saved offers are not rewritten.

If the model explicitly reports response-token exhaustion (`done_reason: length`), generation fails without a correction request at the same budget. No offer or gameplay action is committed. This is a provider failure, distinct from a complete response rejected for invalid story fields.

`GET /api/health` also returns `build: {revision, modified}` for the running Go binary. Revision is the embedded Git commit when available; otherwise `unknown`. Modified is null when unavailable. This identifies the core binary, not the independently served frontend bundle.

## Relationship, voice, rank, title and attribution checks on generated offers

Five further checks run before a generated offer is queued. Each rejects the draft and requests the same single bounded correction used by the existing lexical guards; repeated failure queues no offer and leaves the clock, cash and campaign unchanged. None of them rewrites prose, and none revalidates offers that are already saved.

A speaker may only claim prior dealings with the player if the save records an arrangement between that speaker and the current life. A declined arrangement counts, because being turned down is still a shared past; a previous protagonist's record never does, so a new person after death is addressed as a stranger. This covers returning greetings, summons back to a place, appeals to a previous occasion, assumed routines and claimed track records, in the title, body and approach labels.

Dialogue must stay inside the city. References to the player as a game role, to the cast, or to the brief itself are rejected, including asides that exclude the listener from the cast. The job brief now carries its restrictions in `constraints_never_spoken` rather than inside the fields the model dramatizes, and the focused prompt states that field is never quoted or referred to.

The player holds no rank in an established family, so dialogue attaching a family rank to the player's name or to "you" is rejected. This targets a speaker transferring their own rank to the listener.

A scene title must exist and must not merely repeat one of the offer's approach labels, so the scene keeps a heading distinct from its actions.

A family may be named freely, but a job's own people and premises may not be attributed to a family that is neither the offer's beneficiary nor the recorded owner of the job's location. Faction words come from the family name and both parts of the leader's name. Ownership is read from saved state, so a property changing hands changes what may be said about it. This prevents neutral work, which moves no goodwill, from implying family standing.

These are lexical and relational checks against saved state. They are not semantic verification of motive, plot or history. Public DTOs, saved memories, hidden plans and outcome authority are unchanged.

## Browser 3D city presentation (2026-09-13)

The active city renderer is now Three.js with locally authored Blender glTF models. City selection, camera motion, pedestrian gait, journey interpolation and event effects are presentation only. The old isometric sprite and address-card renderers remain as source references but are no longer city-view choices.

`street[].vehicle` is an optional public vehicle label for an NPC visibly travelling with an operational car (not dry or damaged). Omitted means on foot. It does not disclose private intentions or change travel timing. Existing journey `progress`, endpoints and remaining minutes retain their meanings. The city interpolates between committed progress observations, then holds; it never advances an NPC to an uncommitted destination. Public property condition produces a persistent darkened building; only committed explosion cues produce transient fire/smoke.

A failed `plant` attempt also emits an `explosion` cue: the charge went off prematurely even when it did not destroy the building. Blast casualties are now chosen from living NPCs physically at the affected premises, excluding travellers. Previously the family-wide casualty selection could kill somebody across town and produce a contradictory death cue. The casualty chance is unchanged; an empty building cannot produce an NPC casualty.

Explicit scene replay can restage all cues in the selected committed result. Ordinary revision refreshes and reloads remain silent. Skip immediately clears the transient effects without posting an action or changing saved time.

Street occupancy is presentation-only: rendered travellers may queue behind their
committed progress, never advance beyond it. Small spatial steps and oriented
vehicle footprints prevent rendered bodies passing through one another. When a
source is physically full, the city reports travellers waiting for departure space.
This does not change Go journey timing, decisions, fuel or saved progress.

Normal travel playback caps pedestrians at 1.8 m/s and cars at 11 m/s in the
authored scene scale. The city offers 1× and 4× travel playback; this changes
only interpolation and gait, never saved game time. Walking cycles follow
distance travelled and stop when occupancy makes a traveller wait.

Police/casualty reenactments reserve presentation space beside the cue's target
building, including the casualty's full fall envelope. Occupied slots delay visual
playback until clear; expiry or Skip releases them. These schematic scene positions
are not authoritative NPC locations and never change event outcomes or time.

3D pedestrian selection follows the existing portrait cast and the player's
one-based saved `face` selection. It is appearance only, with no new gender or
gameplay field. Portrait fallback and city selection share Go's unsigned FNV-1a
ID calculation; named painted identities keep their established appearance.

A `gunfight` cue stages an anonymous schematic shooter; existing cues do not
identify the weapon or shooter, so the renderer does not attribute one to a
named NPC. Co-located `killing` playback waits for the first visual shot. These
poses, timing and muzzle effects never create shots, hits, casualties or time in Go.

Staged city gunfire audio follows rendered muzzle beats. Queueing, skipped/muted
playback and missed frames cannot replay a backlog of gunshots. Skip, expiry and
scene teardown cancel the current shot tail. This audio lifecycle is independent
of saved time and combat outcomes.

Vehicle rendering roots now follow the road/pavement surface beneath them,
including parked and staged police cars. These contact heights, soft shadows and
street fixtures are presentation geometry and do not affect Go travel or collision
rules.

Thorne & Sons (`chapel`) now uses its undertaker-specific building model. The
parked hearse in its yard is decorative premises scenery; it is not a public NPC
journey, player vehicle or gameplay vehicle type. Existing address and command
identifiers remain unchanged.

City keyboard navigation is presentation-only: unmodified arrows pan in camera
screen directions, Q/E rotate, +/− zoom and Home resets. Expanding focuses the
canvas; Escape returns focus to the expansion button. Browser/OS modifier
shortcuts and composing input are left untouched. Camera diagnostics do not
change public simulation state or commands.

Pedestrian presentation roots now ease between road and pavement elevation with
conservative animated-stride clearance. This vertical support and the player's
surface-following marker do not alter authoritative travel or traffic progress.

City scene and journey panels share one responsive overlay in normal and expanded
views. Expanding retains the existing Theatre playback and its headline/cast;
Skip/Go on only dismiss presentation. Public command and save semantics are unchanged.

Theatre scene audio is cancellable: dismissal, cue replacement and navigation
release its scheduled voices. Muting stops active scene sounds without replay on
unmute. This affects presentation only and never advances or reverses a saved event.

Confirmed building detonations start behind authored front-window glazing and burst outward,
with a brief fire/light burst followed by rising smoke and a three-second fade.
Confirmation uses a matching target/minute in public `building_fires`; explicit debug
explosions also use this path. Early accidents and older snapshots without confirmation
retain the exterior burst. Models without window anchors still need authored emitters.
This uses the committed explosion cue; it adds no damage, ignition or physics rules.

Explosion debris is cosmetic: confirmed internal blasts eject instanced masonry
from unobstructed authored windows, clear projecting canopies and settle on the
facade pavement. Saved `building_fires` retains settled rubble until `cleanup_at`,
including the period after extinguishing. Replay temporarily hides the settled
instances while animated fragments play. Other explosions retain exterior debris
that disappears with playback. Animated fragments avoid visible actor footprints and
introduce no collision, inventory, obstruction or damage rule in the simulation.

Rendered explosions now own their onset audio, even when a higher-gravity
casualty supplies the caption. A same-address/minute explosion or gunfight
suppresses that casualty's generic Theatre sound. Late or muted onsets are not
replayed. No new causal or damage field is added to public cues.

Vehicle wheel animation follows actual presentation distance, using the authored
0.37m tyre radius. Placement, waiting and parked states do not advance roll.
Wheel pivots and diagnostics do not change Go travel time or vehicle state.

Front-wheel steering follows presentation route curvature and actual travel
distance. Moving occupancy includes the steered tyre sweep; parked vehicles
and staged police retain straight-wheel occupancy at their existing positions.
These internal presentation distinctions add no public vehicle types or commands.

Front tyres use separate inside/outside steering angles derived from the authored
axle spacing. The inside angle remains within the existing clearance limit;
the optional `frontWheels` canvas diagnostic records lateral pivots and angles.

City renderer teardown releases shared model resources, instance buffers,
shadow targets, decoded image bitmaps and its WebGL context. Canvas metrics
include geometry and texture counts for repeated-mount diagnostics; these are
presentation-only counters, not simulation state or total GPU memory usage.

Unused cells within the browser city grid render as fenced vacant yards. They
are decorative parcels without location IDs, actions, ownership or simulation
state. Their geometry preserves the same building reserve and street clearances;
adding a real location to a cell removes its vacant-yard presentation.

Building condition drives stable localized surface staining. It does not identify
the damage cause or imply an active fire/collapse. Full condition removes the
staining; updates reuse material uniforms and leave shared source textures intact.

The 3D city consumes public `sky.kind` and `sky.wet` for overcast lighting,
fog, falling rain and wet road/pavement materials. Rain uses elapsed presentation
seconds, independently of gameplay time and travel playback speed. Disabling
motion (or OS reduced motion) hides precipitation while retaining wetness.
Drying streets follow `sky.wet` even after rain stops; the view invents no weather.

At night, visible cars on public/presentation journeys illuminate their authored
head/tail lamps and project soft pools ahead of the bumper. Paused journey traffic
keeps its lights on; parked cars and daylight traffic do not. Pools are inexpensive
ground projections, not dynamic shadow-casting lights or gameplay visibility rules.

“Find me” toggles camera following of the rendered player character/car, including
travel playback. It preserves zoom/orbit and changes no selected address or Go
state. Manual pan/reset, touch manipulation, building selection and explicit
address/scene focus release following. “Stop following” leaves the camera in place.

City camera close inspection supports orthographic zoom up to 32. Keyboard pan
uses 8.25 / zoom world metres per press, preserving the initial 5m step while
allowing fine movement at close zoom. These controls change no simulation state.

While following the player, intervening building geometry receives a local screen
cutaway around the character/car. Orthographic sight lines are tested at 10Hz;
the visual opening eases in/out (immediately when motion is disabled). Orbiting
to a clear view or releasing follow restores the facade. Building selection,
collision, condition, static shadows and simulation state are unchanged.

Pedestrian material palettes are stable authored cast presentation, selected from
public portrait choice or the existing identity fallback. They do not encode
attire ownership, wealth, faction or gameplay status. Walking and casualty models
use the same palette selection; anonymous gunfight actors remain anonymous.

When player following is inactive, the same local building cutaway protects the
visible staged cast at the active event address. Its screen opening encloses the
cast together and restores after the staged extras expire. It does not move the
camera, reveal unstaged actors or alter saved events.

Male pedestrian cast silhouettes now select authored full/receding hair, bald
scalp, fedora or cloth-cap groups. The choice follows portrait presentation;
headwear is not a newly owned inventory item. A changed appearance key refreshes
an existing pedestrian even when its base model remains the same.

City selection uses parcel-corner marks. The player's ground marker adapts to
zoom, vehicle size and heading; it represents presentation location only and
changes no selection, collision envelope or simulation state.

Junction crossing markings align with existing pedestrian lanes and connect
pavement islands. They are visual road paint, without new traffic signals,
right-of-way rules or simulation effects.

## Browser event framing

Active public visual cues fit the camera to the presentation envelope at their target address before rendering the staged action. This is a camera-only operation: it changes no command, revision, clock or outcome. Explicit replay switches from the interior to the city. Loading a saved result without replay remains silent.

## Debug scene previews

Append `?city-debug` to the browser URL to expose the scene selector. Gunfight, assassination (shooter plus casualty), explosion, arrest and raid previews run at the selected address. Each creates an immutable presentation-only snapshot with a unique preview world ID; no request is sent to the command API. Stop returns the renderer to the campaign projection without replaying saved results. A new committed revision cancels the preview. Controls are unavailable during travel, active gameplay scenes, loading or disabled motion. Reduced-motion preferences remain respected by the renderer.

## Persistent killing aftermath

The public `aftermath` array describes active observable death scenes independently of `last_result`. Each saved entry has `id`, `target`, a named `victim`, `minute`, `police_at` and `cleanup_at`. New public killing cues for an actually dead NPC create one entry per victim. Initial response tuning is police arrival after five game minutes and cleanup after 180; no wall-clock presentation action advances those deadlines. Entries disappear from the public projection at cleanup, while building condition remains governed by repairs. Repeated reports cannot duplicate bodies or restart an old death's lifetime. Older saves without this field load with no invented historical scenes. The city renders a persistent fallen character and blood pool, then a police car and two uniformed officers from `police_at`, using separate traffic reservations. Active casualty playback suppresses the duplicate persistent body. Response travel, investigation behavior and animated cleanup remain pending.

## First 3D room presentation

Saint Agnes uses a locally authored GLB interior on entry. Public occupants populate up to nine clear standing bays, and clicking an occupant selects their existing room actions. The full roster remains available below the canvas. This changes presentation only; occupant presence, action availability and premises ownership remain public Go projections. Other interiors remain on the existing renderer pending their authored rooms.

## Police presentation cast and detainee identity

Arrest cues may include `detainee: {id, name}`. This is distinct from `actors`, which can identify the arresting detective. Player confinement explicitly identifies `player` as detainee. The renderer must not infer a prisoner from an actor when the field is absent. Raid presentation expands one public cue to three police vehicles and four uniformed officers; arrest presentation uses two vehicles, two officers and the explicit detainee when present. Supporting cast IDs are local presentation identities, not new public events. Raid officers walk within reserved approach corridors and stop short of the measured facade. Vehicle arrival, escort and door-entry choreography remain unfinished.

City and 3D interior canvases accept WASD and arrow-key panning with Q/E orbit; browser-modified shortcuts and IME input retain their normal behavior.

## Scene-first result reveal

During active 3D playback the result strip, unread headline banner and theatre caption/cast/headline remain hidden. The theatre announces only the scene location. The renderer signals completion after its staged effect batch has finished, including any placement waits; results then become visible. Explicit Skip or navigation away ends playback and releases the reveal gate. This does not defer the backend commit or alter its outcome.

## Persistent raid presence

The saved/public `police_presence` array contains observable raid cordons with `id`, `target`, `minute` and `cleanup_at`. A witnessed raid creates one cordon per address, initially lasting 120 game minutes. A later raid at that address extends attendance; duplicate delivery does not restart the clock. Reads return detached active records without changing the save. Older saves have no invented cordons. Playback and debug previews never write these deadlines.

The city renders three parked police vehicles and four officers while attendance is active. Raid playback temporarily owns the cast to avoid duplicate officers; persistent attendance returns after playback. Cleanup releases all scene reservations. This is standing attendance: seamless transitions from the approach, forced entry, search and animated departure remain unfinished.

## Authored raid entry

The tavern GLB exposes `entrance-threshold` and `entrance-door-hinge` nodes with a real vestibule opening. At an available aligned staging bay, the lead raid officer approaches, kicks, waits for the inward-opening leaf to clear, then enters. Raid cast playback lasts seven seconds; other building types retain the approach until their openings are authored. All motion stays in a reserved corridor. This is presentation only and does not modify property damage or command outcomes. Existing active raid presence holds the door open; preview Stop restores campaign presentation. The raid target facade stays visible during entry, while other occluding buildings retain their visibility aids.

## Saved building fire response

Successful building detonations create saved/public `building_fires` records: `id`, `target`, `minute`, `brigade_at`, `extinguished_at`, `cleanup_at`. Initial tuning is response at +10 game minutes, extinguishing at +180 and departure/cleanup at +240. Records remain public after extinguishing until cleanup so the renderer can retain attending crews. A later detonation renews the fire and cleanup deadline while retaining an already dispatched brigade's arrival time. A repeated ignition at the same minute leaves deadlines unchanged. Reads are detached and do not advance time or mutate saves; older saves have no invented fires.

A generic explosion cue, including an early charge accident, does not by itself create this building-fire record. Detonation damage remains governed by existing property condition and repairs; extinguishing cannot restore it. Persistent window fire now consumes this field on ten authored building models; brigade rendering remains pending. The longer response window ensures fire remains after the existing 120-minute planting command.


Window emitters are authored Blender `fire-window-*` nodes. The city draws up to four active windows per burning building, preferring the storey above entrance canopies, with looping flame and smoke particles. Disabled motion freezes particle motion while retaining the observable fire. At `extinguished_at` particles are removed without altering damage. Industrial and specialist models without authored window nodes still need fire emitters; no arbitrary facade positions are invented for them.


The saved fire response now stages a locally authored `fire-engine` from `brigade_at` through `cleanup_at`, including attendance after extinguishing. It reserves a 5.8m by 2.35m parking footprint in side bays offset 10m from the parcel centre, separately from police and public vehicles. This is stationary attendance; crew, driving, hose deployment and extinguishing choreography remain pending.


Fire response attendance now includes two authored firefighter characters facing the affected building. They use separate forecourt reservations and share the brigade's saved cleanup deadline. Their initial standing attendance does not yet represent hose deployment. Scene staging reserves stationary public actors at their known parking/standing destinations before their first visible frame, preventing response vehicles from taking a hidden parked car's bay.


Between brigade arrival and extinguishing, visible firefighters on supported facades hold authored nozzles, with separate hoses connected to the engine's side outlets and animated water arcs ending at authored fire-window positions. Hidden/unavailable crews cannot emit water. Disabled motion freezes stream particles; extinguishing removes hoses/nozzles/streams and resets the working arm pose while brigade attendance continues. This remains initial suppression choreography; deployment, reeling-in, impact spray and adaptive hose routing still need production work.

Condition-driven glazing on the ten standard building exports switches front upper
panes (or front ground panes on single-storey buildings) to authored broken-glass
remnants below 60% condition. Repairs to 60% or above restore intact glazing.
This is visible property wear, not a new ignition or damage rule; extinguishing
and rubble cleanup do not repair windows. Other model families need equivalent variants.

The client retains pre-action property conditions for the matching world/revision
as presentation context. Internal blast playback switches glazing at 90ms; the
context survives an interior-to-city renderer remount. Existing broken panes remain
broken, and snapshots without prior context use their known condition. Debug
explosions temporarily break panes and restore saved condition on Stop/completion.
The synchronized glass sound belongs to the cancellable scene audio lifecycle.

Successful player-directed strikes now attach optional `attacker: {id, name, weapon}`
to the attack cue. `weapon` snapshots the attacker’s tier at the event: 0 unarmed,
1 revolver, 2 pump shotgun, 3 Thompson. Delegated strikes use the selected crew
member’s equipment; later upgrades/seizures must not change this saved cue.
Unarmed successful strikes emit `attack` alongside the victim’s `killing`, while
armed ones emit `gunfight`. Their recorded manner of death matches the weapon.
Other legacy gunfight cues remain anonymous unless they explicitly provide an
attacker; renderers must not infer one from the victim’s `actors` list. Gun model
selection and full coverage of other combat producers remain unfinished.

City gunfight presentation now consumes recorded attacker identity and weapon tier
for revolver, pump-shotgun and Thompson models. A staged attacker temporarily
replaces its public street actor. Anonymous legacy gunfight cues retain a generic
revolver; explicit invalid/unarmed tiers do not invent a gun. Shotgun playback uses
two shots with a pump cycle; Thompson playback uses two three-shot bursts. These
are cosmetic cadence choices, not additional backend damage/ammunition events.
The articulated long-gun rig reaches both grip anchors without changing bone
lengths. Replaying uses a distinct presentation serial while preserving saved cue IDs.

Keyboard panning in the city and Saint Agnes interior now follows held WASD/arrow
keys every animation frame instead of OS key-repeat steps. Speed scales inversely
with zoom, diagonal input is normalized, and key release/focus loss clears movement.
This supersedes the earlier per-press pan distance; command and saved state semantics
are unchanged.


Successful player-directed strikes additionally record optional
`strike: {variant, victim: {id, name}}` on the paired attack and killing cues.
Variants are `back-of-head`, `close-shot`, `burst`, and `close-quarters`.
The back-of-head outcome is eligible for a revolver against a stationary, unarmed
victim with no personal grievance (`sore == 0`), with 65% variation among eligible
successful strikes. Shotguns use close-shot; Thompsons use burst; unarmed strikes
use close-quarters. Selection occurs before death, consumes the existing single
manner-of-death world RNG draw and never rerolls success or changes combat RNG.
Death prose and saved scenario agree. These fields describe a resolved scene,
without introducing ammunition accounting or browser-authoritative damage.
Legacy cues omit them. The browser now uses one shot for back-of-head outcomes;
its approach and head-level aim choreography are still pending integration.

## Shared walk-up assassination playback

A recorded `back-of-head` revolver strike now stages attacker and explicit victim
as one cast, consuming only its matching killing cue. The victim waits while the
attacker approaches from behind, raises the gun, fires once at 3.65 presentation
seconds, and the victim falls after impact. Playback lasts 6.5 seconds. The full
cast movement uses an 8.2m by 1.4m forecourt reservation; an alternate left bay can
clear an existing body. Camera framing uses that actual reserved scene. Small
forward blood droplets are cosmetic. Debug Assassination uses this same timeline
in its private snapshot. All other strike variants retain their existing playback.

The persistent body receives the playback position/yaw as a mounted-renderer hint,
so the handoff keeps its fall position and suppresses duplicate corpses. This hint
is not saved simulation geometry; reloads still select from the existing aftermath
layout, and canonical placement across reloads remains unfinished. Cleanup still
uses the saved deadline. Scene completion/Skip retains the existing news reveal gate.

City gunfire now preloads the three user-provided WAV files under `/audio/guns/`.
The recorded weapon selects the sample, full tails overlap, and Skip/mute/teardown
cancel every active voice. A shared compressor moderates burst peaks. Decode or
asset failures retain immediate synthesized fallback; loading never emits a late
shot. Audio diagnostic counters report loaded models, sampled onsets, active voices
and fallbacks without changing gameplay state. Other effects remain procedural.

During killing replay, the victim's own later police response yields with the body;
other death scenes retain their attendance. Response returns after playback under
its original deadlines. Generic gunfight playback waits for all co-located,
same-minute casualty cues to have visible reservations; explicit strike identity
narrows this match to its recorded victim. This prevents firing at an unstaged
casualty and does not add police, deaths or gameplay time.

### Supplied city effect audio

The browser preloads the ten user-provided effect WAVs once. Explosion/raid
presentation windows are 14/10 seconds to retain recorded tails; no simulation
clock change is involved. Explosion debug snapshots include a private temporary
building fire, removed on Stop. Nearby fire and vehicle audio reads existing
public fires and rendered actor movement. Muting/view disposal cancels owned
voices; no audio completion issues a command. The imported escape sample awaits
moving getaway choreography. Debug `effectAudio` reports loaded/played/active
clips for verification, without changing the public server schema.

### Map-first presentation shell

The city renderer stays mounted while informational/action menus open above it;
menu selection does not issue a command or discard the camera. Entering a building
switches the main scene to its interior; leaving switches back. HUD/menu state is
local presentation state. Existing action requests, revision/idempotency and
server outcomes retain their meaning. Newspaper narration/reveal migration is
not yet implemented by this shell change.

### Published-article narration and post-scene reveal

`POST /api/newspaper/speech` accepts `{story: publishedStoryID}` and returns WAV
narration of that saved article's headline and body. Text is resolved from the
published archive; the client cannot provide speech text. Unknown IDs return404,
changed/removed text during synthesis returns409, and an unavailable voice service
returns503. This endpoint advances no time and does not mutate the save. It uses
a stable Herald narrator profile and the bounded process-local voice cache.

The browser links a scene to an article by exact headline and minute. It prepares
narration during playback when voices are enabled, reveals the paper only after
scene completion (or explicit Skip), and starts the prepared clip on reveal. Close,
voice-off and scene replacement cancel pending work and dispose playback. A missing
article retains the existing scene caption. Preparation failure leaves the paper
readable, with an explicit retry control. Scene replay can reopen its article.

## Committed NPC travel playback

`last_result.street_travel` records public NPC travel observed during a `travel`
command. It is an array (including an explicit empty array) for recorded travel,
empty for other commands in HTTP responses, and may be null or absent in older
saves/servers. Each segment
contains the existing `Journeying` public fields at its start, plus `from_minute`,
`to_minute`, and `end_progress`. `progress` is the start fraction; `minutes` is
remaining scheduled travel at that start. Times never exceed the completed
command's interval. Consecutive portions of the same leg are coalesced. A leg
that finishes during the command ends at progress1; a continuing leg ends at its
observed final fraction. No future departure or private plan is published.

NPC departure and arrival times are now simulation boundaries in `Advance`, so
an NPC can depart, arrive and perform the arrival's existing counter work within
one long action rather than waiting until that action ends. This can change
business income and subsequent city outcomes compared with late settlement.
The renderer samples the recorded legs against the same playback clock used for
the player, including playback-speed changes. Collision waits remain visual;
they never issue commands, change outcomes or advance the backend clock. Skipping
or finishing presentation reconciles to the final public street snapshot.

### Close-quarters strike presentation

An explicit `attack` with `strike.variant == close-quarters` and attacker weapon0
now stages attacker and linked victim together, under the same forecourt reservation
as the walk-up execution. The cosmetic sequence approaches, strikes at3.7/4.15/4.6s,
and falls after the last blow, completing at6.5s. These are presentation beats of
the single committed outcome, not extra combat rolls or damage. Only the exact
same-time, same-place linked killing cue is merged; legacy attacks remain unchanged.
Debug choices cover back-of-head, revolver close-shot, shotgun close-shot, Thompson
burst and unarmed close-quarters; preview receipts carry no real action costs.

### Armed close-range strike choreography

Recorded revolver/shotgun `close-shot` and Thompson `burst` strikes now merge
only their explicitly linked victim into a shared cast. The attacker approaches,
the victim raises their hands, and the weapon aims before firing at3.65s.
Close-shot uses one cosmetic shot; burst uses3.65/3.74/3.83s. The attacker lowers
the weapon, turns and withdraws; the sequence completes at8s. The existing
back-of-head and unarmed sequences remain6.5s. Anonymous legacy gunfights retain
their original timing. These presentation beats do not add damage or ammunition
rules. Audio, flash, recoil, blood and victim fall share the strike timing;
body placement retains the victim's initial facing for the aftermath handoff.

For staged strikes, the mounted aftermath renderer now retains a copied final
joint pose (arms, elbows, legs and knees) alongside position and facing. It samples
the completed visual pose before playback, so Skip also has a settled pose, without
running audio or issuing commands. Model/material objects are never retained by
this snapshot. This is still presentation-only mounted state; a full reload uses
the existing canonical fallback placement and pose.

Saint Agnes interior seating is cosmetic and consumes only its existing public
Presence roster. Explicit bartender roles can stand in the service aisle; seating
never changes NPC activity, location or availability. The HTML roster remains
available beyond the modeled seating capacity. Interior Q/E uses held-key
rotation; neither staging nor camera movement advances simulation time.

Daily accounts now use the same current HourlyIncome calculation as clock payouts,
including actual room footfall and custody collection share. This corrects the
reported daily rate; it does not change the existing operating factors or add
Mariner rental income.
