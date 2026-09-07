# Development record

## 2026-09-07 — initial slice
- Created a separate Black Ledger browser prototype; Afterlight code and saves are untouched.
- Authored a fictional port-city neighborhood, command clock, jobs, contacts, crew, property income, housing, hidden retaliation, injury and death.
- Python reference prototype passed 23 headless tests: command availability, idempotent persistence, stale revisions, hidden information, interruptions, security, tribute, income and a new life in the same city.
- Implemented initial browser presentation with a custom SVG neighborhood and semantic HTML controls. Browser playtest remains pending.
- User selected a strict separation of simulation and presentation. Migrating the small reference core to Go before expanding it. React UI and Pixi city rendering are the recommended presentation layer.
- Keep-awake process launched with `caffeinate -di`; verified macOS idle display/system assertions. This does not override deliberate locking/sleep or lid-close behavior.

## Current acceptance work
- Go command core and SQLite adapter implemented; core/store tests pass with the race detector. HTTP integration coverage remains to be expanded.
- React interface and Pixi map implemented; TypeScript and production build pass. Committed travel playback remains to be implemented.
- Native browser QA, local model encounter generation, optional local voices, failure/reload checks.
- Complete a timed campaign playtest before claiming a 20–30-minute slice.

## Scheduled clock and browser foundation
- Go now jumps to scheduled boundaries and accrues income by interval. Regression checks cover exact income and interruption before later task rewards.
- `go test -race ./...` passes. The command server currently has no automated HTTP tests.
- React/Pixi production build passes. Browser QA exercised travel to Saint Agnes and two courier jobs, reaching the first paused authored conversation. Full campaign, AI encounter and voice QA remain outstanding.
- Renderer retains its active texture until replacement is ready, and disposes removed texture resources; SVG object URLs are revoked even after failed decoding.
- Documented prepared time segments, speculative AI proposals, and the rule that animation speed cannot alter authoritative outcomes.

## Browser campaign and consequence QA
- Played authored opening through first Bluebird Laundry purchase: $55 cash, 10 respect, Day 1 11:55 after acquisition.
- A real local model encounter arrived. Requested speech and observed Speaking status; no browser warnings/errors were reported. Audio quality and late-response cancellation still need a dedicated check.
- Accepted the model encounter then immediately reloaded. Recovery yielded $146 at 13:25: $70 job reward plus $21 passive income, exactly once.
- Traveled back to Saint Agnes and clicked Skip journey. Arrival stayed at 13:50 with $151; skipping only removed the presentation.
- Found model prose claiming a character had been removed during a collection job. Completion records now come from the approved operation, and choices name the actual task. Prompts request concrete noncombat work. This does not constitute semantic validation of all generated dialogue.
- Added HTTP tests for duplicate/stale commands, private plot filtering, closed-event speech, invalid input and foreign origins.
- New-life regression ensures NPC identity persists while personal trust and obsolete director state reset.
- Remaining: full timed campaign, broader dynamic rival incidents, voice prefetch/cancellation QA, responsive visual refinements, and end-to-end death/security/expansion playtests.

## Approved painted-noir direction
- User approved painted noir realism with warm daylight and richer nightlife. Saved both reference studies in docs/art.
- Created a transparent casino asset and standalone /art-study.html with day/night grading and animated marquee lights.
- Inspected night presentation in the browser; transparency and bulb placement are usable for the first study. Production street assembly, window masks and moving actors remain pending.
- Frontend production build passes.

## Animated block and business pressure
- Added docs/GOAL.md with the user's expanded visual acceptance requirements. The app goal was observed paused; its objective/resume fields cannot be changed with the available goal tool. Implementation continues during the current turn.
- Added matching transparent café and sedan assets, plus /street-study.html: moving cars, simple pedestrian silhouettes, independent marquee lights, smoke, motion toggle, reduced-motion support, and a scripted car-arrival preview. Rendering suspends when the document is hidden.
- Browser visual inspection found the initial road overlapped façades; moved traffic and pavements outward. This remains a two-building study, not the integrated gameplay neighborhood. People still use simple procedural silhouettes.
- Added ownership-driven faction pressure. A scheduled demand pauses action time; paying improves relations, refusal or rival backing can schedule hidden sabotage, and damage reduces property income. Available loyal crew mitigate business damage; home guards do not.
- Added regression coverage for interruption timing, private schedules, limited tribute protection, hidden sabotage, crew availability, and no demands without a business.
