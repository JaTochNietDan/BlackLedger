# Simulation boundary v1

The simulation is authoritative. UI holds only selection, view, animation and voice preferences. It never computes rules or owns a campaign save.

- GET /api/state: versioned public projection (no RNG, hidden plots or unobserved outcomes).
- POST /api/action: {request_id, revision, kind, target?, event?, choice?}. Exactly-once transactional command. Stale revisions are rejected. Response is the new public projection with committed presentation events.
- POST /api/director: ask the asynchronous director to prepare a bounded proposal. Does not advance game time.
- POST /api/speech: {event}; speak only the active scene's saved text, with a stable voice by character identity. Does not advance game time.

Core is a Go package with no browser, rendering, HTTP or model-runtime dependencies. HTTP/save adapter persists complete state and idempotency receipts together in SQLite. AI and Kokoro are external services, not authorities. React is a replaceable command client. PixiJS only renders committed events and decorative motion. Skipping presentation cannot affect outcomes. Headless clients use the same command API.

A migration/reference prototype exists under prototype-python; do not run it as the active game backend. Its original tests establish expected command behavior, to be ported into Go core/HTTP tests.

## Current presentation and director contracts

`last_result` contains `from_location`, `to_location`, elapsed game minutes, and the records produced by the committed command. Record IDs remain stable as the bounded history rolls over. Street travel is a skippable presentation of these endpoints; it never advances the backend clock. The optional `opportunity` is a suggestion derived from public progress, not a disclosure of hidden plots.

Director proposals name an existing speaker, one supported operation (`courier`, `mediation`, or `collection`), and optionally an existing faction ID as `beneficiary`. Go assigns the operation's reward, duration, heat and respect; validates the beneficiary; supplies canonical completion text; and displays the stakes in the choice. Completion gives that faction +6 standing and rivals −3. Rejection or an abandoned police stop grants neither the reward nor faction credit. Previously saved offers retain their promised terms.

At a resulting heat of 15 or more, a completed arrangement pauses at an authored police stop before paying its reward. The saved scene preserves the original arrangement and beneficiary until the player pays to complete it or abandons it. Hidden retaliation scheduling remains private. AI prose is bounded and prompted against inventing outcomes, but is not a complete semantic truth guarantee.
