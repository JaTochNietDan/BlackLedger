# Development record

## 2026-09-07 — initial slice
- Created a separate Black Ledger browser prototype; Afterlight code and saves are untouched.
- Authored a fictional port-city neighborhood, command clock, jobs, contacts, crew, property income, housing, hidden retaliation, injury and death.
- Python reference prototype passed 23 headless tests: command availability, idempotent persistence, stale revisions, hidden information, interruptions, security, tribute, income and a new life in the same city.
- Implemented initial browser presentation with a custom SVG neighborhood and semantic HTML controls. Browser playtest remains pending.
- User selected a strict separation of simulation and presentation. Migrating the small reference core to Go before expanding it. React UI and Pixi city rendering are the recommended presentation layer.
- Keep-awake process launched with `caffeinate -di`; verified macOS idle display/system assertions. This does not override deliberate locking/sleep or lid-close behavior.

## Current acceptance work
- Finish authoritative Go command core and SQLite HTTP adapter; port tests.
- Replace imperative frontend shell with React components, then integrate Pixi for map presentation.
- Native browser QA, local model encounter generation, optional local voices, failure/reload checks.
- Complete a timed campaign playtest before claiming a 20–30-minute slice.
