# Contributing

Black Ledger is an unfinished, source-available game under PolyForm Noncommercial
1.0.0. Contributions must be your own work or compatible material with documented
redistribution rights. Include provenance and required notices for imported assets
or code. Do not submit purchased packs, private saves, credentials or model weights.

Read docs/GOAL.md for scope and API.md before changing interfaces. Go owns gameplay,
randomness, saves and outcomes; the browser sends commands and presents public
state. Keep hidden simulation state out of API responses. Explain player-visible
behavior, tests and remaining limitations in pull requests. Coordinate shared
src/main.tsx and src/types.ts changes with other active work.

Run `go vet ./...`, `go test ./...`, `npm test` and `npm run build`. CI also runs
Go's race detector and native package smoke tests. Test in a fresh temporary
save with an explicit `-db` argument. Never use `.runtime/campaign.sqlite3` for QA.

Report bugs with build.json (or /api/health), OS, browser, steps, expected versus
actual behavior and relevant logs. Redact private paths and credentials. Do not
upload a campaign save without reviewing its contents.

Submitting a contribution licenses that contribution under the repository's
applicable terms; it does not transfer copyright or grant extra commercial rights.
A separate contributor agreement would be needed if the maintainer wants to
commercially relicense someone else's contribution. No such agreement exists yet.
