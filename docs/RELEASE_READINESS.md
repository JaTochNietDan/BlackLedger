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

## Clean-checkout verification — 2026-09-21

A separate detached worktree at `b10f7b8319b877aeef84efc269102aae4185bcb4`
was created under `.runtime/release-clean-b10f7b8`. It contains committed source
only; none of the original checkout's uncommitted gameplay changes were copied.
`npm ci`, the production frontend build, and all 407 frontend tests passed.
The versioned local Mac ARM64 archive was accepted by the clean-tree guard;
its build.json records the exact commit above and `modified: false`. Its native
archive smoke test passed, including save/retry persistence after restart.
This remains a local QA build, not a tagged or published release.

A focused media review now identifies 23 current files with exact hashes in
`docs/MEDIA_RIGHTS_REVIEW.md`: 14 supplied recordings, four ground textures,
three imported shop originals/derivative and two style references. Commit history
identifies the shop as a Grok image, but describes the ground as photographic
textures without recording a provider or source license. No redistribution
permission has been invented. Older versions and other provenance records still
need review when deciding whether to publish complete Git history.

The GitHub integration was discovered and offered, but is not confirmed installed
or connected. It does not yet supply authenticated repository access.

The clean checkout's `go test -short ./...` completed successfully: core 100.6s,
sim 313.1s, and all other packages passed. This is the short suite; tests that
explicitly skip under `testing.Short()` were not run by that command. Logs are
`.runtime/release-clean-go-tests.log` and `.runtime/release-clean-frontend-tests.log`.
The updated Unix smoke path executes Play.command/Play.sh directly rather than
bypassing the launcher. On the Mac it passed from a path with spaces through
launch, command/retry, shutdown, reload and saved-receipt replay. Windows smoke
still executes the binary directly; validating the .cmd UI remains a Windows
acceptance item. The follow-up `dev-media-review` archive also includes the local
documents directly linked from README/ASSETS; archive membership was verified.

## Campaign amount controls — 2026-09-21

A fresh packaged Mac campaign on isolated port 8879 reached its first owned
business through normal browser controls: courier work, an authored private
job, a laundry order, dock shifts, five-crate trading, a wire, midnight upkeep,
and purchasing Bluebird Laundry. No gameplay fixture or main save was edited.
At revision 23 (day 2, 04:25), the player had $88 cash, $82 offshore, 19 respect,
100 health and the laundry. The first hour after purchase earned $13; the
previous midnight shift paid $75 less $15 upkeep. The books and ownership
controls were usable. The campaign survived a backend process replacement.

This uncovered and fixed two defects:

- Trade amounts were displayed as dollars and incremented by five. They now
  use the existing public goods unit, including accessible labels, range errors
  and submit text, and increment by one. Five moonshine crates cost $180 and
  subsequently sold for $275 in the browser; money controls retained dollars.
- The market wire offer tested the legacy $500 default even when the typed
  minimum was $100. Readiness now tests its affordable preset. A $100 wire from
  $485 left $385 cash and $82 offshore. Backend amount limits, fees and legacy
  zero-amount behavior are unchanged and covered by regression tests.

409 frontend tests, production frontend builds, the server package tests and
focused banking/deposit tests passed; the focused Go race run passed too.
The patched QA build came from the isolated committed-source worktree plus only
these fixes, excluding unrelated local changes. This was a debugging campaign
with development pauses, not the required uninterrupted 20–30-minute acceptance
run. Crew/family progression, major consequences, broader devices and native
Windows/Linux launch acceptance remain open.

Local evidence is in `.runtime/public-campaign-20260921/` (public snapshots at
revisions 6, 14 and 23), `.runtime/release-quantity-*.log` and
`.runtime/release-deposit-tests.log`. The explicitly isolated campaign save is
`/var/folders/2k/kck2sk4n08q6d_w_53nggjj80000gn/T/black-ledger-public-campaign-yzsqt19b/campaign.sqlite3`.
