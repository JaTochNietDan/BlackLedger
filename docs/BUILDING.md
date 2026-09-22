# Development

Install Go 1.23 or newer and Node.js 22.12 or newer with npm. Python 3.12 is needed only for
packaging. CI uses the current stable Go release. From the repository root:

```sh
npm ci
npm run build
go run ./cmd/blackledger
```

Open `http://127.0.0.1:8791`. Rebuild after frontend edits. The server binds only
to localhost; this is a local single-player application, not an internet-facing
multiplayer server. The source workflow saves to `.runtime/campaign.sqlite3`.
Override it with `-db /path/to/save.sqlite3` or `BLACK_LEDGER_DB`.
Use a different port with `-addr :8795` or `BLACK_LEDGER_PORT`.

```sh
go vet ./...
go test ./...
npm test
npm ci --prefix desktop
npm test --prefix desktop
python3 scripts/package-release.py --os darwin --arch arm64 --version dev-local
python3 scripts/smoke-release.py release/black-ledger-dev-local-darwin-arm64.tar.gz
```

Choose the OS/architecture for your machine to run the archive smoke test.
Desktop packaging requires the target OS and architecture because AI speech
includes native dependencies. It refuses to overwrite an
existing archive. It includes the frontend, assets, desktop app and game server,
license notices, build metadata and SHA-256 checksum, never a campaign save.

## Saves

Saves live in the per-user `BlackLedger` folder under `%AppData%` on Windows,
`~/Library/Application Support` on macOS, or `$XDG_CONFIG_HOME` (normally
`~/.config`) on Linux. Updating the extracted game folder preserves those saves.
Stop the game before backing up its SQLite files. Development-checkout saves
are separate and are never automatically imported.

## Local AI services for development

For source-checkout development, AI requests use
`BLACK_LEDGER_OLLAMA` (default `http://127.0.0.1:11435`) and
`BLACK_LEDGER_MODEL` (default `qwen3:14b`). Speech uses the legacy-named
`AFTERLIGHT_DIRECTOR_URL` (default `http://127.0.0.1:8787`). These services receive
game text when enabled. These environment settings apply to developer servers;
the desktop app selects its own private local services.
`BLACK_LEDGER_DIRECTOR_THINK=1` enables experimental reasoning; it is off by default.

`scripts/run-services.sh` is a development helper for a previously provisioned
local `.tools/` installation. It does not install models or make a fresh checkout
self-contained. This repository is independent of the archived Afterlight game.

See [LOCAL_AI.md](LOCAL_AI.md) for desktop first-run setup and [RELEASING.md](RELEASING.md) for distribution checks.
