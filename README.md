# Black Ledger

Independent single-player mafia game. Does not modify Afterlight or its saves.

## Run

From this directory: `go run ./cmd/blackledger`, then open http://127.0.0.1:8791.

State: `.runtime/campaign.sqlite3`. Override BLACK_LEDGER_DB for isolated tests. The HTTP server binds only localhost. Go owns all game rules and saves; the browser only sends commands and displays public state.

AI: BLACK_LEDGER_OLLAMA (default http://127.0.0.1:11435), BLACK_LEDGER_MODEL (default qwen3:14b). Voice: AFTERLIGHT_DIRECTOR_URL (default http://127.0.0.1:8787). Neither service is required to play authored scenarios.

## Verify

`go test ./...`

`go test ./core -bench=Advance -run='^$' -benchmem`

The archived prototype-python folder is a reference implementation, not the active backend. Its original tests run with `python3 -m unittest discover -s prototype-python/tests -v`.

See DESIGN.md, API.md and docs/ARCHITECTURE.md for scope and system boundaries; docs/DEVELOPMENT.md tracks validation.
