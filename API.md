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

Explosion presentation now starts at the target building's rendered front bound,
with a brief fire/light burst followed by rising smoke and a three-second fade.
This uses the committed explosion cue; it adds no damage, ignition or physics rules.

Explosion debris is cosmetic: instanced masonry fragments settle on the facade
pavement and disappear with playback. They avoid visible actor footprints and
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
