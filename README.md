# Black Ledger

A single-player 1950s mafia game: start with a rented room, find work, build a
crew, own businesses, and live with the consequences in a persistent city.
The Go simulation owns the rules and transactional SQLite saves; a React/Three.js
browser client presents the city, interiors, encounters and gambling tables.

**Development preview.** Gameplay, visual quality and balance are still evolving.
This is not a finished 1.0 release. The game contains fictional violence, crime,
gambling and permanent character death. No real-money gambling is involved.

## Play a downloaded build

Extract the entire archive, then run `Play.cmd` (Windows), `Play.command`
(macOS), or `./Play.sh` (Linux). The game opens in your default browser. Keep the
launcher terminal running; close it or press Ctrl+C to stop the server.
No Go, Node, Python or AI models are required to play. A modern browser with
WebGL2 and a working graphics driver is required. See the archive's `PLAY.txt`
for saves, startup options and troubleshooting.

Planned distribution targets: Windows x64, macOS Intel and Apple Silicon,
Linux x64 and ARM64. Initial archives are unsigned; macOS notarization and
native installers are not yet implemented. Published downloads will appear in
[GitHub Releases](https://github.com/JaTochNietDan/BlackLedger/releases) once the
publication checklist is cleared. Do not assume releases exist yet.

Saves live in the per-user `BlackLedger` folder under `%AppData%` on Windows,
`~/Library/Application Support` on macOS, or `$XDG_CONFIG_HOME` (normally
`~/.config`) on Linux. Updating the extracted game folder preserves those saves.
Stop the game before backing up its SQLite files. Development-checkout saves
are separate and are never automatically imported.

## Build from source

Install Go 1.23 or newer and Node.js 22 with npm. Python 3.12 is needed only for
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
python3 scripts/package-release.py --os darwin --arch arm64 --version dev-local
python3 scripts/smoke-release.py release/black-ledger-dev-local-darwin-arm64.tar.gz
```

Choose the OS/architecture for your machine to run the archive smoke test.
The packager also supports cross-compilation. It refuses to overwrite an
existing archive. It includes the frontend, assets, executable, launchers,
license notices, build metadata and SHA-256 checksum, never a campaign save.

## Optional AI and voice

Authored gameplay works without either service. AI requests use
`BLACK_LEDGER_OLLAMA` (default `http://127.0.0.1:11435`) and
`BLACK_LEDGER_MODEL` (default `qwen3:14b`). Speech uses the legacy-named
`AFTERLIGHT_DIRECTOR_URL` (default `http://127.0.0.1:8787`). These services receive
game text when enabled. They are separate installations, not bundled downloads.
`BLACK_LEDGER_DIRECTOR_THINK=1` enables experimental reasoning; it is off by default.

`scripts/run-services.sh` is a development helper for a previously provisioned
local `.tools/` installation. It does not install models or make a fresh checkout
self-contained. This repository is independent of the archived Afterlight game.

## License and contributions

Original project code is **source-available under PolyForm Noncommercial 1.0.0**,
with required attribution in [NOTICE](NOTICE). Noncommercial use, modification
and redistribution are permitted subject to [LICENSE](LICENSE). Commercial use
requires a separate license; contact [JaTochNietDan](https://github.com/JaTochNietDan).
This is not an OSI open-source license. Dependencies keep their own licenses;
media scope and unresolved provenance are documented in [ASSETS.md](ASSETS.md).

Read [CONTRIBUTING.md](CONTRIBUTING.md), [API.md](API.md) and [docs/GOAL.md](docs/GOAL.md)
before changing gameplay or interfaces. [docs/RELEASING.md](docs/RELEASING.md)
describes automated builds and the remaining public-release gates.
The archived `prototype-python` is reference material, not the running backend.
