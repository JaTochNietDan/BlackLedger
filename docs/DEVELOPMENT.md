# Development record

## 2026-09-07 — initial slice
- Created a separate Black Ledger browser prototype; Afterlight code and saves are untouched.
- Authored a fictional port-city neighborhood, command clock, jobs, contacts, crew, property income, housing, hidden retaliation, injury and death.
- Python reference prototype passed 23 headless tests: command availability, idempotent persistence, stale revisions, hidden information, interruptions, security, tribute, income and a new life in the same city.
- Implemented initial browser presentation with a custom SVG neighborhood and semantic HTML controls. Browser playtest remains pending.
- User selected a strict separation of simulation and presentation. Migrating the small reference core to Go before expanding it. React UI and Pixi city rendering are the recommended presentation layer.
- Keep-awake process launched with `caffeinate -di`; verified macOS idle display/system assertions. This does not override deliberate locking/sleep or lid-close behavior.

## Current acceptance work
- Go command core and SQLite adapter implemented; core/store tests pass with the race detector. HTTP integration coverage remains to be expanded.
- React interface and Pixi map implemented; TypeScript and production build pass. Committed travel playback remains to be implemented.
- Native browser QA, local model encounter generation, optional local voices, failure/reload checks.
- Complete a timed campaign playtest before claiming a 20–30-minute slice.

## Scheduled clock and browser foundation
- Go now jumps to scheduled boundaries and accrues income by interval. Regression checks cover exact income and interruption before later task rewards.
- `go test -race ./...` passes. The command server currently has no automated HTTP tests.
- React/Pixi production build passes. Browser QA exercised travel to Saint Agnes and two courier jobs, reaching the first paused authored conversation. Full campaign, AI encounter and voice QA remain outstanding.
- Renderer retains its active texture until replacement is ready, and disposes removed texture resources; SVG object URLs are revoked even after failed decoding.
- Documented prepared time segments, speculative AI proposals, and the rule that animation speed cannot alter authoritative outcomes.
