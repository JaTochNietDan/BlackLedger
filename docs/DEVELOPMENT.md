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

## Live street integration and browser faction test
- The app goal is active again. The expanded visual requirements remain recorded in GOAL.md.
- Integrated the two-building street into the actual React game with same-origin, source-checked, whitelisted inspection messages. Selecting either building changes the real action panel. The wider city remains in the directory pending additional art.
- Sidebar previews now use painted assets for the café and casino. Fixed preview clipping and overlapping header controls found during visual inspection.
- In an independent QA save on port 8793, one-hour wait stopped at 08:30 for the faction demand. Refused it, then advanced two hours: cash $255, respect 12, and laundry damage 35. Browser report did not reveal the hidden attacker's identity.
- The visible QA session on 8792 had additional interactions between observations, so used a separate hidden browser/save for controlled checks.
- Added a visible condition/repair hint for owned businesses and an ambience pause control to the live street.
- Still incomplete: full painted neighborhood, detailed pedestrians, committed travel animation inside street view (currently a saved-arrival banner), voice prefetch/cancellation audit, and full timed campaign playtest.

## Political interactions, progression guidance and speech lifecycle
- Fixed personal-hit deduplication: business sabotage no longer prevents a distinct personal retaliation plot.
- Investigation now reports the actual faction and target, with different protection advice for a business and residence. Added regression checks for Russo attribution and simultaneous plot types.
- Added backend-authored optional next-opportunity guidance from public facts; regression verifies it cannot reveal a hidden personal hit. Damaged businesses take priority over expansion suggestions.
- Extracted voice lifetime/cancellation from the React component into a tested controller. Five Node tests pass for stale responses, character replacement, active cancellation, completion cleanup and unavailable providers. This does not verify subjective voice quality or all browser autoplay behavior.
- Go race suite and frontend build pass.

## Complete API campaign and third painted location
- Full no-cheat HTTP campaign regression passed: 44 committed commands, from initial funds through crew/business ownership, contacts, apartment and security, then feud, death at minute 2400 and a new character. City identity and former-organization property ownership persisted; player money, rank and crew reset. This is automated functional coverage, not a timed human/browser playtest.
- Injuries now reduce intervention survival odds with bounded minimum/maximum probabilities. Added comparative seeded coverage and reran the campaign/race suite successfully.
- Added Bluebird Laundry art and shared public/art/buildings.json used by the illustrated renderer, React location controls and selected-object previews.
- Browser inspection found the laundry initially on the road; adjusted its position to the pavement edge. Increased provisional pedestrian scale relative to vehicles. Surface detail and finished character art are still pending.

## Political job stakes and live campaign checkpoint
- Added validated optional faction beneficiaries to future AI arrangements. Completion adjusts faction standing (+6 beneficiary, −3 rivals), and repeated Bellandi hostility can produce a personal feud. Invalid faction IDs are rejected. Existing saved offers keep their promised terms.
- Rebalanced catalog rewards so special collections are not worse hourly pay than routine courier work. Added a police stop at 15 resulting heat; the reward and faction credit remain pending until completion, and abandonment gives neither. Regression covers interruptions, attribution and exact-once completion.
- A fresh independent browser/save playtest began at 15:53 UTC on September 7. Played earned progression through courier work, real local-AI offers, laundry acquisition, Leo recruitment/delegation, Russo backing, a naturally scheduled sabotage (15 damage with available crew), Ashbury expansion, apartment rental and hired security. The local model produced a Russo mediation with the new beneficiary stakes; completion paid $55 and +5 respect.
- Challenged Bellandi at 23:45 game time, paid $57 midnight expenses, returned home, and attempted to rest. The business demand interrupted at 02:55; after refusing, another rest stopped at 03:45 for a personal attack. Holding the entrance survived at 60 health, losing the hired guard. Reload preserved $86, respect 32, heat 12, and the injury. This run used no state injection, but was interleaved with implementation; wall time is not evidence of a clean 20–30-minute paced playtest. Death/new-life behavior remains covered by the full HTTP regression rather than this surviving browser run.
- Found and fixed the opportunity button's pointer-events blockage. Found routine midnight accounts hiding the consequential confrontation result; the new Latest Developments panel prioritizes consequential records, with other command outcomes in an expandable list and a ledger link.
- Fixed capped-history rollover dropping new presentation records. Added a regression with a full 180-record history.
- Added The Mariner painted sprite as the fourth neighborhood building. Street presentation now includes a player marker, skippable 2.2-second committed travel along authored pavement waypoints, and a reusable asphalt texture. Ground is cached between lighting changes; crossed streets are drawn as a continuous junction.
- Current automated verification: Go race suite, frontend voice tests and production build pass. Remaining art work includes market/docks, finished pedestrian art, richer sidewalk dressing and broader viewport QA. Police decisions have regression coverage but still need an observed browser encounter.

## Portraits and police browser verification
- Added six painted cast portraits in a shared atlas and verified Mara/Alex in the live interface. Other future player identities retain the previous fallback. Compressed the desktop street header to reserve more space for the illustration.
- Added `cmd/qa-fixture`, which exclusively creates a new isolated police-test save and refuses existing output paths. This is a deterministic QA fixture, not natural campaign progression.
- Browser fixture: accepted a Russo courier at heat 14; at 08:45 cash remained $90 and respect 0 while Detective Harlow required a choice. Paid $40; cash became $125, respect 3, heat 7, and Russo standing +6. Expanded the additional-outcome disclosure successfully. Separate core tests cover abandonment and duplicate-decision rejection.
- Day and night street screenshots show the asphalt and continuous junction without browser errors. The small pedestrian art and sparse surroundings still need improvement.

## Housing ownership and progression audit
- Previous goal turn made concrete progress: committed painted character art, a repeatable police QA fixture, and browser-verified police consequences. No blocker is present.
- Found estate purchases only changed the address. Purchasing now records the deed; moving away preserves ownership, returning costs no second purchase, and death transfers the residence to the former organization without inventing rental income. A new stranger cannot buy an already occupied residence through the ordinary purchase action.
- Found repeatable respect rewards from moving between homes. Track the highest housing tier reached in this life; only entering a new tier earns respect. Starting a new life resets that progression.
- Save schema v2 repairs an active v1 estate purchase on load and initializes housing rank from the saved residence without charging money, changing revision or advancing time. The migration does not invent ownership for dead players.
- Regression tests cover purchase/return/death/new-life ownership, repeated rental moves, and legacy save migration. Go race suite and frontend production build passed. Owned residences display condition rather than a misleading zero hourly business income.

## AI-authored approaches with bounded consequences
- Added optional contextual approach labels to model proposals. The director may offer careful preparation or a pressured schedule; Go assigns fixed tradeoffs in time, reward and heat. Unsupported/duplicate approaches and invalid labels are rejected; unoffered choices cannot be forged. Existing offers remain compatible.
- Saved approach effects survive cloning/persistence and police interruptions. Other incidents that interrupt the work still prevent completion credit. Regression covers discreet completion, pressured police interruption/payment, forged choices and business-pressure interruption.
- The authored opening also introduces both approaches, preserving playable choices without the model.
- Verified a real qwen3:14b response through a fresh isolated server: Mara proposed a Bellandi delivery with waiting until the club emptied versus approaching during a busy shift. Chose waiting at 09:00; at 10:15 cash rose from $135 to $195, respect 2 to 5, heat remained zero, and Bellandi received +6 standing.
- Browser inspection found four options overfilled the former narrow dialog. Four-choice scenes now use a wider two-column desktop layout, retaining a single column at compact widths. Visually checked the desktop fixture with all choices and footer visible; no browser errors. Full narrow-viewport and uninterrupted pacing verification remain outstanding.

## Building damage and repair presentation
- Added a damaged Bluebird Laundry storefront and optional manifest damage variants. Public property conditions are passed through the isolated street adapter; the renderer does not calculate damage or repairs.
- Browser QA caught an opaque background in the generated variant. A second generated cutout also lacked alpha. Reused the original verified transparent sprite as a runtime mask, preserving its footprint and removing the rectangle in both street and selected-building preview.
- Extended the isolated fixture tool with `damage`. At 45 condition and $90, the UI showed the damaged storefront and $6/hour. Used the real Repair action: at 09:00 condition was 85, cash $46 (including accrued income), and income $11/hour; both views returned to intact art. Browser error log was empty.
- Production build and script syntax check pass. This implements one damage variant and the reusable state/asset path; variants for other buildings and broader district art remain unfinished.

## Repeat-life business availability
- Found a long-term progression dead end: former organizations permanently occupied the finite businesses after successive deaths. Added an explicit buyout at twice the normal acquisition price for former-organization businesses; ordinary faction property remains unavailable through this action.
- New people still start with $90 and no inherited ownership. Buying out preserves damage and the business's post-death income, and starts the usual political pressure. Opportunity guidance now recognizes buyout candidates rather than silently skipping them.
- Regression verifies insufficient funds are rejected, full payment is charged, changed property state persists, and faction property cannot be seized. The full Go suite passes. This is a rules-level verification; broader repeat-life browser pacing remains to be tested.
