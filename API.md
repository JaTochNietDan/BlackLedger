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

A proposal may include up to two `approaches`, each with a `method` and short player-facing `label`. Supported methods are `careful` (+30 minutes, −$15 reward, up to 3 less heat) and `press` (15 fewer minutes, +$20, +5 heat). Unknown/duplicate methods and oversized labels are rejected. Standard acceptance and refusal remain available. The saved scene contains authoritative effects for each offered approach; forged choices are rejected. Police interruptions preserve the selected terms, while another event interrupting the work grants no completion reward or faction credit. Omitted approaches preserve compatibility with existing offers.

## Director contact progression

New AI requests use established contacts: Mara initially; recruited associates with at least 30 loyalty; family leaders once their family has at least +6 goodwill. Among these, the server prefers the least recently featured speaker, counting current-life arrangement memory and pending offers. Saved continuations preserve their original speaker. The prompt and decoding schema receive this allowed cast, and server validation rejects other speakers with one bounded correction attempt. Other NPC context is not an invitation to impersonate an unavailable contact. This governs new request generation, not audiences or already saved scenes, and does not guarantee semantic correctness of all dialogue.

## Located AI jobs

New model responses require `location`, an existing place ID in an unlocked district. The schema restricts IDs; server validation requires the selected venue's name in the body and rejects explicit locked venue names in title/body/approaches. This is a bounded consistency check, not a semantic proof of the rest of the story. Core validation also rejects supplied inaccessible locations. Authored and legacy proposals may omit location.

Located proposals store the venue in scene `target` and arrangement memory `location`. Action details display its canonical name. Jobs still resolve off-screen within their quoted total duration and return to the player's current base; they do not teleport the player, unlock districts or grant free property access. Reading and rendering remain independent of simulation time. Saved follow-up briefs prefer structured venue memory over place names inferred from old prose.
