# Experimental story review

The production director still uses structural validation. A second model pass is being evaluated to catch unsupported ownership changes, false memories and promised effects outside the operation contract.

Run from the gameplay repository:

```sh
python3 scripts/evaluate-story-review.py --output /tmp/story-review.json
python3 scripts/evaluate-story-review.py --holdout --output /tmp/story-review-holdout.json
```

The script requires only Python's standard library and the local Ollama endpoint (default11435, qwen3:14b). `--url` and `--model` override those settings. It does not interact with any game save. The output is written incrementally so a timeout does not lose earlier cases. Expected answers are withheld from the model. Current output includes the prompt and facts for reproduction.

Initial evaluation missed a false completed-work claim. Explicitly separating conditional future results from past completed work fixed that case; six original and six holdout verdicts then matched expectations. The holdout was not used to rewrite the prompt. Correct verdicts sometimes had an additional incorrect reason, so reasons cannot be treated as authoritative repair instructions.

Before enabling in gameplay:

- Derive review facts from saved ownership, identities, actual job status/results, known threats and the validated operation's future effects. Keep offer claims separate from facts.
- Test actual generated proposals from a full campaign, including acceptable fictional disputes, to measure false rejections.
- Bound review latency and retries; a reviewer outage must not freeze the clock or mutate the save.
- Revalidate life/revision-relevant facts at commit, and retain the existing structural checks.
- Keep evidence of both rejection and acceptance. These twelve curated cases are insufficient to declare arbitrary stories coherent.
