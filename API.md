# Simulation boundary v1

The simulation is authoritative. UI holds only selection, view, animation and voice preferences. It never computes rules or owns a campaign save.

- GET /api/state: versioned public projection (no RNG, hidden plots or unobserved outcomes).
- POST /api/action: {request_id, revision, kind, target?, event?, choice?}. Exactly-once transactional command. Stale revisions are rejected. Response is the new public projection with committed presentation events.
- POST /api/director: ask the asynchronous director to prepare a bounded proposal. Does not advance game time.
- POST /api/speech: {event}; speak only the active scene's saved text, with a stable voice by character identity. Does not advance game time.

Core is a Go package with no browser, rendering, HTTP or model-runtime dependencies. HTTP/save adapter persists complete state and idempotency receipts together in SQLite. AI and Kokoro are external services, not authorities. React is a replaceable command client. PixiJS only renders committed events and decorative motion. Skipping presentation cannot affect outcomes. Headless clients use the same command API.

A migration/reference prototype exists under prototype-python; do not run it as the active game backend. Its original tests establish expected command behavior, to be ported into Go core/HTTP tests.
