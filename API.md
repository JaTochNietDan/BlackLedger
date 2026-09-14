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

Explosion cues now optionally carry `detonation: "planted" | "premature"`.
Successful `detonate` calls record `planted`; an early player charge accident
records `premature`. Player planting outcomes also capture the existing `attacker`
identity with weapon0. Faction detonations select the most skilled available
faction member already inside the target (ID breaks ties), provided that person
has a valid different home and no custody, death or existing departure schedule.
They immediately begin a saved journey home before the blast, excluding them
from indoor casualties; the cue captures their identity with weapon0. No eligible
member leaves the individual identity absent. This does not yet dispatch a
planter from another address. Legacy cues omit this field.
This records the resolved cause for presentation without changing damage, charge
consumption, command duration, odds or fire rules. It survives result/save replay.
The renderer prioritizes explicit detonation outcome over matching fire records;
legacy cues retain their existing inference. Debug includes a separate premature
explosion with no invented building fire. Planted explosions with a recorded attacker and an authored animated doorway now
play a planter exit before detonation (6.2s on flat entrances,12.9s at the villa). The city reserves the full exit,
uses the recorded person, and delays blast audio/light/debris, glazing damage,
new fire and co-located casualty playback until the exit completes. Fire-brigade
staging follows the scene and the newspaper waits for completion. This currently
covers all19 current building models, including the villa landing/stairs; unnamed faction planters
and premature injury/escape choreography remain unsupported. Legacy explosions
retain their existing playback. Debug also includes Explosion · casualty for
reviewing the combined sequence without a campaign command.

The planter preamble starts with a doorway view, then pulls back from1.6s to.2s
before its detonation to fit the building and blast envelope. The villa adds a
planted clearance step beyond the last tread before its turn/departure; the
reserved path includes this extra pavement space. Camera pan,
rotation, zoom, address/whole-city focus or enabling player-follow cancels that
automatic pullback. This is presentation only and does not alter scene timing.

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
Confirmation prefers explicit `detonation: "planted"`; legacy cues use a matching
target/minute in public `building_fires` or a legacy debug explosion. Explicit
`premature` accidents retain the exterior burst even if a fire already exists.
Older snapshots without confirmation also retain the exterior burst. Models without window anchors still need authored emitters.
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
replayed. Audio ownership does not change damage rules.

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

## Save schema v15 — NPC residences

NPCs persist optional `home` (location ID) and `accommodation` (display category),
independent of Location, Post and Heading. New-city population and migration assign
missing homes in deterministic wealth/standing order, preserving existing valid
residents. `Presence` adds `home_id`, `home_name`, and `accommodation`; these identify
a residence, not the NPC's current whereabouts. Dead NPCs never count as residents.

Initial capacity is24 residents at The Mariner,64 at Ashbury Court and one private
resident at Cypress House, with a place reserved for the player at their home.
Ashbury supports shared flats for lower standing and private apartments for higher
standing; the private estate requires high wealth and cannot admit a stranger into
a player/family-owned house. Capacity shortages leave homes unassigned and appear
in the public `housing_shortage` count and People screen. Additional housing remains
planned. Existing income, rent payments, ownership and movement rules are unchanged:
this residence assignment is not yet a paid lease or NPC return-home schedule.

## Residential routines

Ordinary NPCs with a home and workplace now schedule homeward journeys at midnight, with deterministic departure staggering, and reconsider work at 06:00. Existing noon evening routines remain. Duties, custody and urgent errands retain priority; this is not a complete shift system for officials or business managers. Home arrivals preserve the separate workplace. These trips use the existing public street and recorded travel segments; no frontend clock or endpoint changes. Taking office clears a successor's previous journey and sets their new workplace.

## Save schema v16 — residential rent accounts

NPC `budget_day` guards salary and living-expense settlement once per calendar day. Property `rents` maps tenant IDs to accounts (`day`, current `due`, most recent `paid`, outstanding `arrears`, lifetime `collected`). Day keys are minute/1440+1. Missing old-save accounts start on the next settlement without historical bills. Rent replaces the former bundled lodging expense: Mariner/shared flats cost15 daily, private Ashbury apartments35, plus4 other necessities. Existing salary recovery covers the selected accommodation; missed family payroll can leave arrears.

Only living residents pay; work/travel does not remove tenancy. Cash transfers once to the current owner, and property receivables remain with the premises on transfer. Independent-owner payments leave the modeled cash pool. Old debts remain on the previous premises if somebody moves; debt enforcement/eviction is not implemented. Mariner/Ashbury earn no duplicate hourly cash. Books and public income show the contracted living-tenant rate, which is a forecast rather than guaranteed collection. A player owning their Mariner/Ashbury residence pays no room rent to themselves; Cypress household upkeep remains unchanged.

Locations have `rent_register` (null outside Mariner/Ashbury): capacity, occupied NPC count, daily contracted rate and a non-null tenants list with ID/name/accommodation/daily rate. The player's reserved place is not an NPC tenant. Owner views additionally include each existing account by value; nonowners do not receive account balances. The interior exposes the register. Mariner purchase and operating controls are still outstanding.

## Save schema v17 — Mariner freehold and lodging operations

The Mariner retains type `home` and its zero move-in price; kind is now `lodging`. Acquisition uses a separate3600-dollar base freehold, with existing additional-holding and former-organization premiums. Normal `acquire` readiness, command validation, staffing and property ownership apply. Holding the Mariner counts toward the price of later holdings. Existing leases/accounts survive purchase.

The lodging trade has3 staff at6 dollars per day,40 supplies (coal and linen),4 daily drain,100-dollar restocking and140-dollar boiler remedy. Normal repair, wage, manager, operating-mode and protection systems apply. No generic standing-order contract is offered or accepted. Mariner rent is the integer room base15 multiplied by condition, service capacity capped at full service, and operating mode. Service failure/damage therefore lowers actual tenant charges; zero condition produces no current rent. Accounts/public forecast/register/inspection use this current rate. Staffing/stock obligations remain even with empty rooms.

Property Income15 is a capacity reference for business eligibility, never an additional Mariner cash stream. Player hourly accrual and family operating-cash accrual exclude residential rent, whose cash enters through tenant payment only. Family income forecasts include current living-tenant rent. The schema upgrade establishes the lodging trade without resetting ownership/accounts or restocking existing v16 businesses. Later repeated migration cannot refill an already introduced lodging trade.

## Save schema v18 — property exchange and deed transfers

`property_market` is a non-null array of current Mariner/Cypress entries: id/name, owned/available, holder, asking price, broker offer, condition, registered-resident count, player-home flag and district lock. Prices are authoritative; the UI locates the address for normal action handling. Asking/reference prices are distinct from broker bids. This first market uses standing broker offers, not individually funded NPC bids or timed auctions.

`sell_property` is offered for owned Mariner/Cypress deeds with positive value. It immediately transfers ownership to the independent market and pays65% of the base deed price scaled by condition, then spends30 minutes on closing. Acquisition premiums are never refunded. Fixtures, staff and tenant accounts remain; the player's posted guard is released from the property assignment. A resident seller stays at the same home/location as a renter under the displayed daily charge. Proceeds increase cash, not earned-income/job progress. Normal revision/idempotency handling applies.

`buy_residence` purchases the Cypress deed for3500 without moving the player or existing residents, followed by60 minutes of paperwork. Moving home remains a separate action. The existing combined Cypress purchase/move action remains compatible. Mariner purchases continue through `acquire`. Property `bought_life` prevents repeated Mariner reacquisition from awarding respect again in the same life; old owned properties receive the current life during migration. None of this implements tenant eviction, an individual NPC property budget or additional residential geography.

## Capacity-checked move-in and disclosed rehousing

`move_home` now derives a residence plan on copied NPC records before offering the action. It reserves the player's new place, preserves existing residents where capacity allows, and names any required rehousing addresses in the action detail. A move that would leave a previously housed living NPC without accommodation is disabled. The plan neither moves people physically nor changes cash, posts, active journeys or rent accounts during reads.

At completion the command checks the plan again. If occupancy/allocation changed during paperwork, it preserves elapsed time and intervening simulation events, refunds any move payment, and reports postponement without changing the player's home/title or granting housing progress. Otherwise residence/accommodation changes are logged and committed with the move. Previous property rent debts remain on those premises; NPC physical movement continues through normal journeys. Deed-only purchases still preserve all residences.

This is a rehousing rule for voluntary move-in, not a complete tenancy-law, eviction or emergency-housing system. Automatic insolvency/new-life housing and insufficient city-wide capacity still require further work alongside additional housing.

## Additional residential address — Mercer Court

`mercercourt` is a district-zero rental residence with 48 resident places. Existing homes retain priority; missing homes may use its shared flats ($12/day) or private apartments ($25/day), based on the same standing threshold as Ashbury. A player lease costs $120 and $25/day, grants apartment-tier protection/progression and uses the normal capacity-checked `move_home` command. It does not transfer the building deed. The address has ordinary public occupants, rent-register projection and street journeys. Old saves receive its independent property record through `SettleNewPlaces`; no save field or command shape changes. Its map lot fills an existing vacancy without relocating any previous address.


## Action-grounded next-step guidance

The existing `opportunity` shape is unchanged. Its suggested work now comes from the current action projection at the destination, including actual timing, price, availability and detail. A local value copy changes only the inspected player location; reading guidance does not move the real player, alter the save or advance time. Arrival can change the action's availability, so the frontend still selects an address and requires normal travel/action handling.

Early guidance routes to available envelope work, then cargo work when envelopes are unavailable; it suggests earning capital before an unaffordable purchase or hire, and offers only currently available repair, recruitment, district expansion and move-in actions. It does not reveal private plans. Custody, death and pending events suppress these street-work hints. These are public next-step suggestions and existing milestone progress, not a new persisted quest/reward system.


The existing guide crew milestone now includes the initial driver recruitment as well as later organization sign-ons. It inspects actual available hiring actions at unlocked addresses, and either kind of hired crew completes it. The cargo action is labeled Work a cargo shift because it is available by day as well as at night; command ID, pay, duration and risk are unchanged.


## Save schema v19 — daily hotel linen order

Property `rush_order_day` is an optional one-based calendar-day reservation; zero/missing means no order has been taken. Existing saves upgrade without historical reservations or resetting supplies/condition. The reservation belongs to the property and survives a new life or ownership transfer.

`rushorder` at Bluebird Laundry is a60-minute helper job paying80 dollars on uninterrupted completion, with up to one respect below the existing DockName cap. It requires being at the address, not owning it, starting between08:00 and17:00 inclusive, condition at least40, no operating trouble, at least one staff member and two supply units. The normal command boundary checks life, custody, pending events, revision and idempotency.

Acceptance reserves today's order and consumes two of the premises' supplies before time advances. Interruption pays nothing and leaves the order/supplies used; it is not a resumable arrangement. A successful job does not add attention or count toward the fixer's envelope-job milestones. Another order is available the next day subject to the same operating requirements. Owners receive their business's ordinary income rather than this helper payment.

The action appears under Work with its actual terms/refusal. Opening earning guidance tries it after unavailable envelopes and before repeatable cargo work. No new frontend command shape is required. This is an external customer order, not a new fully modeled hotel business or NPC-paid contract market.


## Incendiary property attack

`incendiary` targets a business at the player's current address. It requires a nonowned property with positive income and condition, no active unextinguished fire, and40 dollars. Normal action command life/custody/pending-event/revision/idempotency checks apply. The method charges40 once; the projected generic action cost remains zero, with the actual supplies cost disclosed in its detail. The command advances15 minutes after resolution.

The attack removes15–30 condition points, bounded by remaining condition, removes up to10 supply units and marks operating trouble. It adds18 player attention (capped100); a faction owner loses25 goodwill and receives the existing retaliation scheduling call. It does not directly kill occupants, remove staff, destroy bankroll, consume an explosive charge or grant respect. This is the initial balance, not a completed broader arson economy.

The existing saved building-fire lifecycle supplies brigade arrival, extinguishing and cleanup; extinguishing does not repair condition. A new `incendiary` visual cue has gravity8, the target address and an attacker identifying the player with weapon0. Its news is an attack-category arson report. No save schema or endpoint shape changes. The browser stages an authored held-bottle approach, throw, flight and escape before revealing the newspaper. A reserved forecourt footprint protects the actor path. Presentation hides window fire until impact and defers the saved brigade response until the cast exits; authoritative fire timestamps and damage remain unchanged. An Incendiary debug scene uses a private snapshot. Facade/path clearance across all addresses, impact detail and animation quality remain under development.

## Building drive-by property attack

The `driveby-building` action calls `World.BuildingDriveBy` against a nonowned operating business
at the player's current address. The player must be alive, free, carrying an
equipped firearm (tier1–3), and have an operational personal car with petrol.
The first crew member must satisfy existing delegation availability and loyalty
rules and already be at that address to drive. The driver's own car is irrelevant.
Wrecked premises and buildings with an active unextinguished fire are refused.

Initial balance:35 dollars,24 attention capped at100, and the fuel consumed by
10 minutes of driving. Revolver damage is14–19 condition points, shotgun20–25,
Thompson26–31, capped by remaining condition. Supply loss is floor(actual damage/3),
bounded at zero; operating trouble is set. A faction owner loses25 goodwill and
receives the existing retaliation scheduling call. The attack itself does not
start a fire, create casualties, remove staff/bankroll, consume charges or grant
respect. The method does not advance the clock; command dispatch advances the
quoted10 minutes once through the normal transactional action path. The action
appears in The Street with its35-dollar/fuel cost, crew/car requirements and
attention disclosed; its generic action cost is zero because the method pays.

The resulting `driveby-building` cue has gravity7 and captures the player in
`attacker`, the driver in `actors`, and optional `drive_by` data:
`{driver: {id, name}, vehicle, vehicle_tier, condition_before, condition_after}`.
Names, equipment and actual damage refer to the event time, including after
serialization in `last_result.cues`. The attack-category newspaper headline is
shared with the cue. Legacy and unrelated cues omit `drive_by`.

City3D plays the captured car and distinct driver/shooter through a
reserved road sweep, with timed gunfire, facade dust and vehicle audio; the
Building drive-by debug preview uses a private Packard/Thompson cast without
changing campaign condition. Material wear and glazing reveal the captured
condition loss across the firing beats; before firing they show the recorded
starting condition. Scene completion or cancellation restores the current
snapshot's condition (including undoing private preview damage). Aim points and
impacts use the first visible authored building surface along a ray, rather
than its bounding box; shot impacts sample the muzzle at the firing beat.
The recorded newspaper is revealed after scene completion, through the existing
result/newspaper flow. Normal completion, replay, Skip and disabled-animation
presentation have been exercised on isolated campaigns. With City animation off,
the recorded newspaper opens immediately; disabling animation during a scene
finishes its presentation and retains its newspaper. Re-enabling animation does
not restart that scene. Broader visual acceptance remains unfinished.


Pending planter exits now hold their future footprint against newly arriving moving
traffic. Existing occupants can continue along their committed routes out of that
space; the scene still waits for actual clearance before starting. Releasing or
cancelling the scene releases the hold. This advances no simulation time and allows a stationary pedestrian on the frontage pavement to walk aside when a
collision-free space is available. This cosmetic movement preserves the saved
location, activity and clock; the resulting stance persists until a new route is
assigned. It stays on the same frontage, respects other traffic and held exits,
and does not move parked vehicles or aftermath. Fully occupied frontages and
other stationary blockers still require separate staging recovery.

After a cosmetic pavement sidestep, a new pedestrian journey from the same
original start first walks back along that frontage under traffic collision
checks. Presentation journey progress remains zero until the start is reached.
This preserves continuity without changing the saved route, elapsed minutes or
outcome. Different models/starts retain their existing transition behavior.

Planted explosion cues now use the exact headline filed by their damage/casualty
report, retaining the same minute. This lets automatic newspaper presentation
match both fatal and nonfatal blasts without guessing among unrelated articles.
Previously those cues used a shorter headline which prevented automatic reveal.

The same cosmetic pavement clearance now also responds to a pedestrian waiting
at an occupied departure. It uses the waiting actor's traffic footprint, keeps
both saved journeys/locations unchanged, and does not advance a zero-progress
departure merely to free space. Vehicle departures are not included.

Premature explosion cues now optionally include `accident: {health_lost, fatal}`.
This records the player's actual health decrease after armour and the zero-health
fatal result at detonation, before later command time, healing or a new life.
It applies to the recorded attacker. Successful planted blasts and legacy cues
omit it; absence must not be interpreted as survival. Existing result/save JSON
retains the value without a migration or altered damage/odds/time rules. The
city renderer uses this event outcome for a reserved frontage fall: fatal actors
stay prone, while survivors begin a supported sit-up. Blast effects originate
at the staged actor, with no invented building fire or broken windows. The
normal attacker is suppressed during the cast. Debug offers separate surviving
and fatal accident previews. Persistent aftermath, seamless handoffs and actual
player-death overlay timing still require further integration.

## Fatal charge accident aftermath

A fatal player charge accident now creates a saved aftermath entry alongside its
explosion cue. Optional `cause: "charge-accident"` identifies the fall pose and
`face` preserves the victim's portrait selection; `victim.name` is captured at
death and `victim.id` is life-specific (`player:<life>`). Existing NPC records and
legacy saves are unchanged. Police and cleanup retain the existing five- and
180-game-minute deadlines. Surviving accidents create no body.

The city suppresses the persistent body during fatal accident playback and
retains its staged reservation for the handoff. The persistent cast reconstructs
the same final fall pose from the saved appearance. Reloads choose an available
accident frontage slot; exact presentation coordinates are not saved. Blood and
response reservations remain presentation-only. Animated police arrival and
cleanup remain unfinished; a paused dead campaign does not advance their clock.

## Survivor presentation handoff

A completed or skipped recorded assassination can also transfer a surviving,
stationary pedestrian player at the recorded address to the attacker's final
position and facing. Matching ordinary pedestrian attackers and victims are
suppressed while the shared cast is active. Parked vehicles retain their own
presence. Private previews, moving actors and different addresses do not adopt
the stance. A following journey walks from this stance to the frontage and saved
route under normal collision checks. Its displayed clock continues through that
connection; the committed snapshot already contains the arrival minute.

A completed or skipped real charge-accident scene can transfer a surviving,
stationary pedestrian player to its final rendered frontage position. This
changes no location, health, time, route progress or save. The next committed
journey first connects that stance to the same frontage's walking lane and then
to the saved route start, using ordinary traffic collision checks. Moving actors,
vehicles, different locations, fatalities and private previews do not use this
handoff. Browser reload restores the canonical position; stance coordinates are
not persisted. Fully blocked connectors wait rather than teleporting.

## Emergency-response timing during travel playback

The browser retains the pre-travel public aftermath, police-presence and
building-fire records for the active journey, merging them by ID with the
committed result (newer records win). Body visibility, police/brigade attendance,
fire extinction, rubble cleanup and held-open raid doors use the interpolated
travel minute. A record absent at the journey's end can therefore remain visible
until its known deadline. Skip uses the committed snapshot immediately. This is
presentation memory only; no command, deadline or save changes. It cannot recreate
an event absent from both endpoint snapshots, and does not replay every historical
building condition or lighting change.

### Collision-delayed journey clock

During recorded travel playback, the browser budgets the remaining displayed
minutes against the player's remaining route distance at its physical speed.
The clock can continue while the player yields, so other recorded travellers
and emergency-response deadlines keep advancing on that same minute. Traffic
delays stretch the presentation rather than consuming all displayed time before
the player finishes walking. Reaching the destination completes the clock;
Skip and reduced motion still reconcile immediately to the committed endpoint.
No travel duration, collision outcome or save is changed in Go. Individual NPC
collision delays and events absent from the endpoint records still require
broader historical-playback handling.

Recorded NPC legs are now played in order. If an actor is delayed by cosmetic
occupancy after its recorded arrival minute, its route remains visible until
physical arrival. A later recorded leg waits for that handoff and is admitted
at its observed starting fraction, then catches up within normal movement and
collision limits. Partial legs stop at their last observed fraction. After a
normal player arrival, the remaining NPC queue keeps its actors and occupancy
until physical completion, using the fixed committed minute. This finishes
already observed travel; it does not simulate new actions while the clock is
paused. The queue belongs to one world, life and revision. Skip, disabled motion,
a new snapshot or a new scene reconciles to the authoritative endpoint. A new
actor assignment clears obsolete per-leg progress before using its new route.

## Card-table camera presentation

Poker now draws a 3D table directly from the public cards projection. Missing
opponent cards remain face down; the renderer never inspects the private deck or
calculates a hand. Folded seats clear their rendered cards. Poker currently
reconciles immediately to each committed street; blackjack retains its existing
dealing sequence. Both table cameras allow bounded local orbit/pan/zoom/reset,
including focused-canvas keyboard input, without commands or saved camera state.
Text summaries and transactional betting actions remain available below the scene.

## Neighborhood residential prices

The optional public `property_market[].neighborhood_index` is a percentage of
normal local residential value (60–100, default100 for older servers). Mariner
freehold acquisition and Cypress deed/buy-and-move actions use this index. Broker
offers apply the same index before their existing65% spread and condition factor;
Mariner acquisition premiums still apply. Rental bills are unchanged.

Core records pressure when a committed, located Witness incident occurs: robbery2,
attack3, gunfight/vehicle gunfire6, incendiary7, killing8, explosion10 percentage
points, capped at40 total. District boundaries define neighborhoods for this first
market model. Quiet game time recovers one point per1,440 minutes; fractional-day
recovery is retained when another incident occurs. A public integer discount
rounds remaining pressure up, so an incident does not lose a whole point after
one minute. Multiple kinds of incident at one address can contribute separately.
Unknown locations, nonviolent political cues, private plans, UI previews and
presentation replay do not contribute. No new random draws or elapsed time occur
when calculating/reading prices.

Saved `property_pressure` stores pressure units and last incident minute by
district. Missing data starts at normal prices without reconstructing old crime
from newspaper prose. This additive field survives JSON/SQLite saves and lives;
market reads apply decay without mutating it. Existing command revision and
exactly-once checks continue to protect transactional prices. NPC deed transfers,
new individual homes and broader market demand are still pending.

## Separately owned apartments

World saves now include optional `apartments`: stable numbered deeds for the64
existing Ashbury Court places and48 Mercer Court places. These subdivide existing
residential capacity; they are not112 new buildings or additional beds. Each deed
has `id`, `building`, `number`, `owner` (independent, NPC ID, or player life ID),
and optional `resident`. Committed housing reconciliation retains existing unit
assignments, clears departed/dead residents and prefers a returning resident's
vacant owned flat. It does not move their current location or journey. Housing
move previews never mutate this registry. A player move cannot displace an NPC
who owns the apartment they occupy.

`buy_apartment:<id>` buys the current rented flat only when the independent broker
owns it, for60 minutes and the quoted internal payment. `sell_apartment:<id>` sells
one of this life's flats at its address for30 minutes and65% of current asking.
Actions remain on the normal revision/request-ID path and are grouped as business.
Base prices are1200 Mercer/1800 Ashbury, multiplied by neighborhood index. Ownership
ends that resident's rent but grants no building freehold. Moving retains the deed;
selling retains the resident and restores their rent. Sales are asset proceeds,
not earned-income progress.

At midnight after PeopleDay, at most one funded NPC purchase occurs. A resident
with price+200 cash can buy their broker-owned flat or buy from a cash-poor NPC
owner (under100 cash); otherwise an NPC with price+500 can buy from such an owner
as an investment. NPC-to-NPC sale conserves their combined purses. Existing
residents retain their tenancy and journeys; no actor is teleported. Deceased
NPC/previous-player deeds become broker stock through this daily process. New
protagonists never inherit prior player deeds. This is a first property market,
not yet mortgage, bidding, probate-beneficiary or voluntary moving-house AI.

NPC rents go to the unit owner, capped by tenant cash through the existing daily
rent account. Player-owned occupied flats contribute contracted rent to Books;
actual receipts use Earn. An NPC owner renting to the player receives the housing
part of a fully paid daily bill. Owner-occupants pay no rent (NPCs still pay other
living costs). Repairs/common-part economics remain on the existing building.

Optional public `apartment_market` lists the player's rented flat and owned deeds,
with ID/building/number/address, owned/home/available flags, owner/resident names,
asking/offer and daily_rent. Market shows reference prices when a private owner is
not selling. Existing Presence.home_name now includes the flat number, and
Accommodation says Owned apartment for NPC owner-occupants. HomeID remains the
building destination. No private travel or plans are added to this projection.

## Household cash, burglary and home attacks

Private saved `household_savings` maps NPC IDs to `{cash, day}`. It is not included
in public snapshots. At midnight, living residents can withdraw enough for living
costs before PeopleDay, then save20% of pocket money above the larger of50 or three
days' living costs, at most60 per day and1000 total. Actual purses fund every
transfer. Repeated savings settlement within a day cannot deposit twice. Cash
moves with its resident rather than with an apartment deed. Legacy saves begin
with no household cash; no historical loot is invented. Apartment affordability
and cash-poor seller checks now include this saved cash, spending purse first and
then savings. Invalid/unfunded household payments mutate neither balance.

`burgle:<npc-id>` is a45-minute personal action offered at a known living resident's
home. It excludes the player's crew/organization and requires at least25 health
and freedom from custody. It does not reveal savings or occupancy through its
readiness message. Success takes only that household's existing cash, credits
Earn, and adds7 attention; an empty account pays nothing. Failure retains the
stash, adds14 attention, damages health/clothing and can be fatal. Occupancy lowers
the success chance; it means physically at Home, not travelling or held. A present
resident identifies the intruder; witnesses may identify a failed intruder when
the resident is absent. Identification creates a personal grievance and the
existing family retaliation consequences. The recorded robbery/news cue feeds
neighborhood prices; replay never performs the crime again.

`home_strike:<npc-id>` replaces the generic personal strike label for a resident
physically at home. It revalidates home address and presence at execution, then
uses existing Strike outcomes/weapon/health/time rules. Optional
`VisualCue.strike.setting: "home"` captures the location context before the victim
is killed or moves; generic street scenarios omit it. Already-departed residents
cannot be attacked via this home command. This is authoritative action support;
private-room 3D assassination/burglary choreography remains unfinished, and the
current browser still uses the existing city result/news presentation.

Cypress House now uses its own locally authored 3D drawing room and study on
entry (`estate`). The room seats five existing public occupants and provides two
standing positions, plus the player's reserved entry aisle. Additional public
occupants remain in the complete person list. Seating, entrance motion, camera
orbit/zoom/pan, and cutaway walls are presentation only; they do not assign homes,
advance time or change action availability. The room does not yet illustrate
purchased house fixtures or stage residential assaults inside it.

Poker's 3D presentation now seats only the public `cards.seats` participants,
using matching public presence faces when available, plus the current player.
Cards deal round-robin from a common deck position; new board cards slide in and
newly public showdown cards turn over against a physical back face. No private
card is inferred. Folded hands leave the 3D felt; the public text summary remains
available. A restored sitting opens settled, while a new hand animates. The
motion setting and reduced-motion preference settle the presentation immediately.
These effects do not delay or compute any gameplay command or payout.

The apartment broker now supports rental investment purchases as well as buying
one's current rented home. `buy_apartment:<id>` remains a 60-minute, internally
paid command, offered at the unit's address; it transfers only an independent
broker's deed. Occupied units keep their resident and route future actual rent
payments to the new owner; vacant units yield no rent. Existing pricing,
neighborhood pressure, sale spread and NPC private-owner protections apply.
`apartment_market` now includes up to three occupied and one vacant broker flat
per building, plus all player holdings and their current unit. Listings replenish
as units leave the broker board. The projection includes optional `locked` for
district access. Listing reads are deterministic and never mutate residences.

Ashbury Court (`apartment`) now displays its own 3D entrance hall, with three
bench seats, nine clear standing positions and a reserved concierge position
when a public occupant has that role. The player uses a separate entrance aisle.
Mailboxes represent the existing 64 numbered units; they do not create units or
assign residents. The lift and upper stair opening are scenic access cues, not
new commands. The complete public roster and residents' register remain below
the render; no private-flat interior or indoor assault playback is implied.

Household savings now reserve enough for local home ownership: residents with a
numbered apartment may save up to `max(1000, current apartment price + 500)`;
other households retain the $1000 ceiling. Deposits still come from actual purse
cash, at most $60 per day, and a valuation decline never confiscates existing
savings. All saved household cash remains private and available to the existing
purchase, bills and burglary rules.

The daily population update and committed-command boundary now settle missing
housing assignments before refreshing numbered flats. Later replacements and
recruits therefore receive homes when capacity exists. Reads remain mutation-free;
established residents keep their homes and assignment does not teleport them.

At the player's current Ashbury or Mercer home, the interior offers a local
private-apartment view and a return to the entrance hall. This presentation
switch does not change authoritative location, time, tenancy or command state.
Only the player appears in the private flat; public building occupants are not
projected into it. Selecting someone in the building list returns to the hall.
The control disappears when this is no longer the player's home. The shared
furnished bedsit does not yet depict fitted upgrades or a bathroom interior.

Fassano Meats (`butcher`) now displays an authored 3D shop interior. Public
occupants use nine customer positions and one reserved butcher/shopkeeper/clerk
position behind the counter; excess occupants remain in the complete roster.
The player has a separate clear entrance aisle. Shop dressing, display trays,
scale and cold-room door are cosmetic, not public inventory or new actions.
Existing selection, camera controls and authoritative building commands apply.

All registered 3D interiors offer an enlarged room view. This resizes the render
in the existing page; people and commands remain below it. The control or Escape
while the canvas is focused restores standard height. Camera and selection are
retained; no simulation command or time advancement occurs. Interior model,
lighting and cutaway metadata are centralized in `src/interiorSettings.ts`.

Russo Motor Works (`garage`) now uses an authored 3D workshop with nine customer
positions, one reserved public mechanic position and a separate player entrance.
The lift, engine stand, tools, tyres and service desk are cosmetic room dressing;
they do not assign or service any vehicle or alter inventory. Existing garage
commands remain authoritative, and the full public roster remains available.

Explicit home strikes at Ashbury or Mercer Court now play inside the shared
private-flat 3D scene. The renderer uses the recorded attacker, victim, weapon
and variant; a selected death cue may use only its same-place/same-minute attack
with an explicitly matching victim. Asset loading precedes animation timing;
completion feeds the existing newspaper reveal. Replay sends no command. Motion
off/reduced motion settles the recorded result immediately once assets load.
Other residential addresses retain their prior presentation for now. Persistent
indoor aftermath and burglary playback are not implemented by this scene.
