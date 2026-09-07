# Black Ledger

Independent single-player mafia game. Does not modify Afterlight or its saves.

## Run

From this directory:

```sh
npm ci
npm run build
go run ./cmd/blackledger
```

Then open http://127.0.0.1:8791. Rebuild after frontend edits. React provides the interface and Pixi renders the illustrated city; the Go server serves the compiled frontend from `dist/`.

State: `.runtime/campaign.sqlite3`. Override BLACK_LEDGER_DB for isolated tests. The HTTP server binds only localhost. Go owns all game rules and saves; the browser only sends commands and displays public state.

AI: BLACK_LEDGER_OLLAMA (default http://127.0.0.1:11435), BLACK_LEDGER_MODEL (default qwen3:14b). Voice: AFTERLIGHT_DIRECTOR_URL (default http://127.0.0.1:8787). Neither service is required to play authored scenarios.

## Verify

`go test -race ./...`

`npm test` (speech cancellation and lifetime tests)

`go test ./core -bench=Advance -run='^$' -benchmem`

The archived prototype-python folder is a reference implementation, not the active backend. Its original tests run with `python3 -m unittest discover -s prototype-python/tests -v`.

See DESIGN.md, API.md and docs/ARCHITECTURE.md for scope and system boundaries; docs/DEVELOPMENT.md tracks validation.

### Repeatable browser edge case

Create a fresh, isolated police-stop save (the command refuses an existing path):

```sh
go run ./cmd/qa-fixture .runtime/police-check.sqlite3
BLACK_LEDGER_PORT=8795 BLACK_LEDGER_DB=.runtime/police-check.sqlite3 go run ./cmd/blackledger
```

Accept the courier offer, verify its reward remains unpaid at the police stop, and choose whether to pay or abandon. This fixture is separate from normal progression and must not be used as evidence of an earned campaign run.

For damage/repair presentation QA, append `damage` to the fixture command, using another new output file. The fixture starts at the owned laundry with 45 condition and enough money for one repair. Verify the street and sidebar change together after repairing.
