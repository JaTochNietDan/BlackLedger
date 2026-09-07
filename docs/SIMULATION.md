# Automated campaign laboratory

Run from this repository with Go installed:

```sh
go run ./cmd/simulate -runs 100 -steps 160 > report.json
go run ./cmd/simulate -runs 100 -steps 160 -director fixture > fixture-report.json
go run ./cmd/simulate -runs 1 -seed 27 -strategies investor -steps 100 -trace > trace.json
```

No server, browser, image, speech service or model is required. The runner creates fresh in-memory worlds and never opens campaign saves. Reports contain per-campaign seeds, action/event counts, final values, command counts at milestones, and aggregate reach/death/error counts. `-trace` includes commands and pre-command minute/cash/health for diagnosis. It can produce large output; use it for a single failing campaign.

## What is exercised

The player policies consume a serialized projection of `World.Public()`, including normal action availability and affordable dialogue choices. They cannot access private plans, RNG, queued offers or director memory. Every player action goes through `core.Execute` with the current revision, exactly as the server's command handler does. Each accepted command checks monotonic time, a single revision increment, nonnegative cash and valid health. This does not test HTTP, database transactions, speech or rendering; those retain separate tests.

- **worker:** ordinary dock shifts, returning home to recover from injuries. Does not pursue ownership.
- **investor:** income, crew, contacts, housing/security and three businesses; delegates when available, pays demands, and repairs severely damaged businesses.
- **defiant:** the same investment priorities, but refuses business demands instead of paying; isolates the consequence of that political stance.
- **reckless:** challenges Bellandi early, then stays home without buying protection. It deliberately exercises the warned-about lethal opening choice.

These are explicit strategies, not claims that a typical player behaves this way. A strategy bug can distort results: inspect traces before changing game balance. Initial development caught a reckless policy that repeatedly traveled instead of remaining at home and therefore avoided the attack.

## Director modes and reproducibility

`authored` runs only the production authored incidents. `fixture` additionally queues simple neutral proposals every 240 game minutes, when eligible, using the production proposal validator and operation rotation. It is a deterministic test provider, not an imitation of live storytelling quality or generation latency. It exercises optional approaches and the same job completion/police paths. Neither mode makes model calls.

The first run uses `-seed` exactly. Later runs add the 32-bit Weyl stride `0x9e3779b9` with wraparound, spreading samples instead of relying on adjacent small xorshift seeds. The actual seed is recorded for each run. Replay one result with `-runs 1 -seed ACTUAL_SEED -strategies POLICY -director MODE -steps LIMIT`. Random record IDs differ, but command outcomes and aggregate reports reproduce under the same game version. Keep the Git commit with a report; future rule changes intentionally change results.

Milestone medians are conditional on reaching that milestone; the accompanying reach count is essential. Final cash is not normalized by elapsed game time, so compare it with game minutes and policy behavior. A command limit is not a wall-clock playtime estimate. Automated runs cannot prove 20–30-minute human pacing, dialogue coherence, tension or UI quality.

Recorded live-proposal replay, more nuanced political strategies and long-term balance experiments are follow-up work. The runner currently stops at first death; browser/core tests separately cover new lives and persistent property history.

## Recorded proposal replay

```sh
go run ./cmd/simulate -runs 100 -steps 160 -director replay -corpus tests/corpora/local-mediation-20260907.json > replay-report.json
```

The corpus is one JSON array of `core.Proposal` objects, bounded to 128 proposals/2 MiB. Unknown fields, invalid operations/speakers/factions/approaches and trailing JSON are rejected before a campaign starts. The report fingerprints the exact input with SHA-256. At eligible 240-minute preparation boundaries, each proposal is queued once in file order, with normal 30-minute arrival delay and production validation. Exhaustion resumes authored-only play; the corpus does not loop. Per-campaign `replay_queued` exposes how much of it was exercised; a character may die before any replay arrives.

This replays **gameplay payloads**, not the original campaign or the model's reasoning. A story generated in one world may be narratively inappropriate in another. Matching original context, operation/contact briefs and semantic novelty remain separate live-director checks. The same corpus/seed/policy/version reproduces mechanical outcomes; it does not establish story quality. No network or model calls occur in replay mode.

The included local-model sample was generated with qwen3:14b from the earned campaign on 2026-09-07. It repeats a prior garage dispute and is intentionally retained as a known-bad story example, not curated writing. Its bounded gameplay proposal is valid. The production director now rejects a matching recent title or identical offer before queueing and uses its normal bounded correction attempt; replay deliberately tests the payload independently of that contextual novelty gate.
