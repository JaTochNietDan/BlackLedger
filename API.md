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
