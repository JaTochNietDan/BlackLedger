# Experimental story review

The production director still uses structural validation. A second model pass is being evaluated to catch unsupported ownership changes, false memories and promised effects outside the operation contract.

Run from the gameplay repository:

```sh
python3 scripts/evaluate-story-review.py --output /tmp/story-review.json
python3 scripts/evaluate-story-review.py --holdout --output /tmp/story-review-holdout.json
python3 scripts/evaluate-story-review.py --cases docs/story-review-campaign-cases.json --output /tmp/campaign-review.json
python3 scripts/evaluate-story-review.py --audit --cases docs/story-review-campaign-cases.json --output /tmp/campaign-audit.json
```

The script requires only Python's standard library and the local Ollama endpoint (default11435, qwen3:14b). `--url` and `--model` override those settings. It does not interact with any game save. The output is written incrementally so a timeout does not lose earlier cases. Expected answers are withheld from the model. Current output includes the prompt and facts for reproduction.

Initial evaluation missed a false completed-work claim. Explicitly separating conditional future results from past completed work fixed that case; six original and six holdout verdicts then matched expectations. The holdout was not used to rewrite the prompt. Correct verdicts sometimes had an additional incorrect reason, so reasons cannot be treated as authoritative repair instructions.

## Real campaign evaluation: do not enable this gate

`story-review-campaign-cases.json` contains seven unmodified local-AI offer bodies from the completed campaign and two explicitly marked authored controls. Per-case facts reflect ownership and canonical parent completion at the time of the offer, not the final player's balances or relationships. Final arrangement status for the current offer is withheld. Expected labels and human assessments are also withheld. Two ambiguous cases (factions disputing access, and a colloquial promise of a favor) remain unscored rather than forcing a debatable binary answer.

The original reviewer accepted every case: **3/7 scored matches**, missing both false-owner stories, the contradictory leader allegiance and a nonexistent delivery deadline. See `story-review-campaign-results.json`. The earlier twelve short cases did not generalize to these actual dialogues.

The `--audit` experiment asks for verbatim evidence before the verdict. It detected the four labeled bad encounters but rejected both valid controls and returned a nonverbatim excerpt for the valid courier: **4/7 matches, one invalid response**. Several explanations were wrong even where the verdict matched. For example, it objected to a proposal serving Bellandi because the job had not yet happened, instead of identifying Elena's Russo affiliation. It called explicit confirmation of unchanged ownership a violation, while admitting in the explanation that it was not a contradiction. See `story-review-campaign-audit.json`.

Neither variant is enabled in the game. A binary match alone is insufficient: reasons, false rejections, malformed responses and ambiguity must all be evaluated. The audit was designed after seeing the first campaign results, so this corpus is now development data, not a fresh holdout. The script now preserves raw parsed responses before validating them; the saved audit run predates that addition and retains only its error for the nonverbatim case.

Next experiments should separate mechanical claims from harmless dialogue and test on fresh generated encounters before runtime integration. Existing deterministic affiliation, location, operation and save checks remain in force. No user campaign was changed by either evaluation.

Before enabling in gameplay:

- Derive review facts from saved ownership, identities, actual job status/results, known threats and the validated operation's future effects. Keep offer claims separate from facts.
- Test actual generated proposals from a full campaign, including acceptable fictional disputes, to measure false rejections.
- Bound review latency and retries; a reviewer outage must not freeze the clock or mutate the save.
- Revalidate life/revision-relevant facts at commit, and retain the existing structural checks.
- Keep evidence of both rejection and acceptance. These twelve curated cases are insufficient to declare arbitrary stories coherent.
