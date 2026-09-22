# Distribution preparation — 2026-09-21

## Implemented

- Attribution/noncommercial code licensing: upstream PolyForm Noncommercial 1.0.0
  plus required NOTICE; public documentation says source-available.
- Browser-game archives for Windows x64, macOS Intel/Apple Silicon and Linux
  x64/ARM64, each containing the Go executable, built frontend/media, launcher,
  dependency notices, build identity and SHA-256 sidecar.
- Desktop startup locates assets beside the executable and saves under the user's
  configuration directory. Source checkout behavior is retained. Explicit `-db`
  / BLACK_LEDGER_DB overrides remain available. The server binds before opening
  a save, and optional browser launch occurs only after successful startup.
- GitHub Actions tests, five native packaging/smoke jobs, retained development
  artifacts, tag-triggered draft prereleases and pinned action revisions.
- Contribution, asset provenance, release and installation documentation.

## Evidence

Local test archives use label `dev-release-check`; they were built from e8f3895
plus the working tree, **not a clean release commit**. Pre-existing uncommitted
work was preserved: core/mugging.go, core/robbery.go, core/armed.go,
src/city3dAftermath.ts, sim/died_test.go and store/rent_recipient_test.go. Archives
record `modified: true`. They are internal QA artifacts, not published releases.

- `npm test`: 407 passed. `npm run build`: passed; existing >500kB chunk warning.
- `npm audit`: no reported dependency vulnerabilities at audit time.
- `go vet ./...`: passed.
- `go test ./cmd/blackledger`: passed after startup changes.
- `go test -race ./cmd/blackledger ./store`: passed (16.8s / 2.7s).
- Initial `go test ./...`: core passed in 391.8s and other listed packages passed;
  sim timed out at 600s with only the pre-existing untracked TestWhatKillsThem
  diagnostic still running (260 campaigns). This is not an all-green full suite.
- `go test -timeout 30m -skip '^TestWhatKillsThem$' ./sim`: passed in 345.1s;
  only that pre-existing untracked diagnostic was excluded.
- Versioned packaging from the dirty working tree was correctly rejected; local
  QA uses an explicit `dev-` label.
- All five target archives compiled locally; each checksum validated, archive
  contents inspected and no campaign/.runtime/.tools content present.
- Native Apple Silicon archive smoke: passed launch from unrelated directory,
  frontend bundles, health/state, isolated per-user default save, port-conflict
  rejection without save creation, exactly-once command replay, and durable state/receipt replay after restart.
- Extracted Apple Silicon game in Codex browser at 1280×720: city loaded,
  Saint Agnes travel/skip advanced exactly 15 minutes, interior loaded; no
  captured warning/error logs. This is not full visual/performance acceptance.
- `cmd/apicheck` against that extracted binary and its fresh isolated save:
  100 commands, 52 distinct command kinds, 3 new lives, no invariant failures;
  request replay matched committed state and stale revision returned HTTP 409.
  71 command kinds were not exercised. Report: `.runtime/release-api-check.json`.
- actionlint 1.7.12: workflow passed. Older 1.7.7 did not recognize the valid
  macos-15-intel runner label, so validation used the current version.
- Gitleaks 8.24.3: scanned 919 commits / 10.40MB text, no leaks reported.
  A separate limited current-tree credential-pattern scan had no findings.
  Scans are evidence, not a guarantee that all history is safe to publish.

## Outstanding

Media redistribution permission/provenance remains unresolved in ASSETS.md.
The available GitHub browser is signed out; no remote exists, no repository has
been created, no source has been pushed, and hosted CI has not run. Windows,
Linux and Intel Mac archives have only cross-compilation evidence locally.
Unsigned/unnotarized packages are not native installers. A fresh sustained
campaign, save-upgrade/backups, wider browser/device testing and the gameplay/
visual gaps in docs/GOAL.md remain part of first-public-preview acceptance.
