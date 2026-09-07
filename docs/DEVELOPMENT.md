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

## Compact viewport guidance
- Verified the prior goal turn made progress through the paid-buyout implementation and regression tests. Continued with a rendered responsive check.
- At 980×800 the existing CSS hid next-opportunity guidance. Kept it beside the city heading at compact desktop widths and placed it in its own row below the view controls below 800px.
- At 390×844 verified the prompt, city illustration, wrapped location controls and latest outcome are visible without horizontal clipping. Property actions stack in one column. At 980×800 verified the restored guidance fits beside the heading.
- These were read-only inspections of the running campaign, not a mobile gameplay completion test. Temporary browser viewport overrides were reset. Production build passed.

## Uninterrupted opening playtest and warning/focus fixes
- Fresh isolated browser campaign measured 16:58:54–17:11:36 UTC on September 7: 12 minutes 42 seconds, without implementation breaks. Earned laundry ownership, recruited Leo, expanded a district, rented an apartment and hired protection. Provoked Bellandi, then survived the first home attack at 60 health. This does not establish the intended 20–30-minute full rise-and-fall arc or browser death/new-life completion.
- Real model offers appeared, but repeatedly used sealed-envelope courier plots and near-identical titles. Recent fiction/operation variety needs a separate director pass. One standard offer was unintentionally accepted by a keyboard event landing on its newly focused first choice; not every decision in this run was intentional.
- Dialogs now focus their non-action container. Browser fixture confirmed Enter leaves a new offer unchanged, Tab reaches Read aloud, and an explicit decline works. This prevents initial focus from selecting an answer. Production build passed.
- Contact warnings previously logged during long rest without yielding control before the attack. They now interrupt at the warning boundary with a paused call. Acknowledging spends no time and leaves the threat active; ordinary travel/security/audience actions remain available. Already-due attacks still resolve immediately.
- Regression verifies unanswered warnings cannot advance time, interrupted rest gives no healing, leaving after the call avoids personal injury while the residence is damaged, and protected defenses still resolve. Updated the HTTP campaign to explicitly acknowledge warnings. Full Go race suite passes.
- Added an exclusively new-file `warning` QA fixture. Browser verified the paused call, acknowledgment at 10:30 with health 60 and $90 unchanged, and the normal travel action afterward. This is deterministic edge-case QA, not earned campaign progress.

## Director memory and operation variety
- Previous goal turn made concrete progress through committed warning interruptions and safer dialog focus. Continued from the current worktree; the Mac's caffeinate process remains live.
- Added bounded persistent memory of the actual last 24 offers, tagged by life, speaker, operation and beneficiary. Rules record offered/in-progress/police-pending/completed/declined/abandoned/interrupted status separately from fictional offer text. Only completion records an authoritative result. Police outcomes update the original job rather than remembering the detective's dialogue as a new offer.
- The model receives this memory and a required least-recent operation, including pending offers. A mismatched operation is rejected before entering the queue. A latest completed current-life job supplies a follow-up contact/allegiance constraint; declines or unresolved work do not become successful callbacks. Original offered dialogue now also enters the player's ledger.
- Regression covers police payment/abandonment, interruption, cloning, queued-operation variety, declined offers and previous-life separation. HTTP model-stub verification checks the actual outgoing memory/brief and rejects an ignored operation. Full Go race suite passed.
- Live qwen3:14b test: first response rejected for an overlong approach label. Added the precise 3–65-character limit to the prompt. Retry generated a laundry shift mediation. Browser completed it at 10:00: cash $135→$190, respect 2→7, with Bellandi credit. Read-only save inspection verified completed memory retained the original offer and canonical result.
- Subsequent live proposals used collection then courier, proving operation variety rather than another consecutive courier. First text was third-person narration; strengthened direct-dialogue instruction. A follow-up acknowledged mediation but copied “loading hours” from a prompt example rather than the actual “shift hours.” Removed that example and explicitly instructed exact subject preservation. This final wording adjustment is not yet verified by a further live generation; semantic continuity remains imperfect. No claim of endless or fully reliable storytelling is made.
- Rendered browser inspection shows the four painted buildings and readable sidebar. Sparse surrounding streets and tiny pedestrians remain visual follow-up work. The overall 20–30-minute completion gate is still open.

## Committed attack scenes
- Added public visual cues to command results for resolved business sabotage, damage to an absent player's residence, and resolved attack choices. Cues contain only location and committed result; they do not expose the hidden faction/plot. Unresolved attacks do not emit outcome scenes.
- The street adapter consumes these cues for a short car arrival, figures approaching, and departure. It never computes attacks or changes game time. Replay and Skip are presentation-only. Automatic playback respects reduced motion and the ambience preference; death/decision dialogs remain the priority. Only currently painted locations can play these sequences.
- First browser screenshot caught the initial car route crossing behind buildings. Constrained event cars to the nearer authored road and mirrored their art on the opposite diagonal. Browser fixture then verified an attack at the laundry: 65 condition, $9 hourly income, 09:00 and $100 after accrued income. Skip returned actions immediately without changing those values. The screenshot verifies the scene caption and initial route; full mid-sequence occlusion/polish remains a visual follow-up.
- Core tests verify cue persistence, absence of private plot information, no premature outcome at an attack decision, and no old cue emitted by the next unrelated command. Full Go race suite and frontend production build passed.

## Browser death, persistence and new-life verification
- Previous turn produced committed attack cues/scenes and browser skip evidence. Continued the earned pacing save on the latest core at port 8799; this resumed test is separate from the earlier uninterrupted 12m42 measurement. The old browser tab had closed, so opened a new tab against the same confirmed live save/server.
- From Day 2 05:30, health 60, cash $168 and respect 28: traveled to The Monarch, repeated its explicitly extreme-risk demand, returned to Ashbury Court and rested. Mara's warning interrupted rest at 08:35 without healing. Acknowledged, deliberately stayed home, then chose Hold the entrance at 10:05. Alex died. These were ordinary UI actions, not fixture state injection.
- Reload preserved the death modal and prevented ordinary actions. Added a compact ending recap from current authoritative values: final respect 29, $969 earned, one property left behind (Bluebird Laundry). Rendered desktop verification shows the recap and new-life button fitting in the modal.
- Used Begin as a new person: Nico Ward appeared at Day 2 18:05 with $90, respect 0, health 100, heat 0. Browser inspection of Bluebird Laundry showed Alex Varga's survivors as owners. Read-only state confirmed no crew or owned property, district 1 retained, and Alex's Day 2 10:05 death record. This now verifies the death/reload/new-life user flow in the earned campaign, alongside existing core/HTTP persistence tests.
- Production build passes. The recap is a presentation-only addition; no game rules changed this pass. A full fresh uninterrupted 20–30-minute run is still unproven, and city/character art remains incomplete.

## Accessible city address book
- Replaced the active schematic Pixi directory view with visible React destination cards. District filters, ownership, current location, locked status and condition are readable without finding hidden map buttons. Painted properties use their actual thumbnails; unpainted destinations use a restrained architectural symbol until their art exists. The painted street remains the main city view.
- Browser mouse selection of Ashbury Court now updates the property panel directly, with no change to money or time. At 980×800 verified filters, two-column cards, guidance and the property panel fit. Separated the scrolling directory from the fixed city heading to avoid cards passing underneath it.
- Traveled from The Mariner to Ashbury Court in the isolated second-life save: 18:05→18:55, cash unchanged at $90. The saved-journey caption appeared, then the normal action buttons became enabled without another command. Directory travel uses a short caption; the street still shows the walking presentation. Restored the temporary viewport override.
- Build passed; active production JS now measures about 227 kB (72 kB gzip), down from about 469 kB for the main bundle plus renderer chunks. The old study modules remain in source but are no longer in the active import graph. Updated architecture/run documentation to reflect React plus the isolated canvas street.
- Full city art and the uninterrupted pacing acceptance test remain outstanding. This pass fixes destination usability, not the remaining art coverage.

## Mercer Exchange painted preview
- Generated a muted 1930s pawn/trading-house storefront matching the existing building camera and material style. The first result included an unwanted explanatory slogan; a second image edit removed it.
- Both outputs painted a checkerboard rather than delivering alpha. Preserved the selected RGB source unchanged and added an explicit SVG silhouette mask for normal UI compositing. This is not a claim that the source PNG is transparent.
- Added a separate preview manifest so illustrated directory/property previews can be introduced without inventing unsafe street positions. Mercer Exchange now appears in the address book and selected-location panel; its gameplay location/actions are unchanged. The street still has four landmarks, and Mercer placement requires the pending neighborhood layout pass.
- Browser screenshot verifies both Mercer previews against the dark UI without an opaque checkerboard rectangle. Production build passed. Source image, silhouette and manifest are stored in the independent repo.

## Prepare queued voices before encounters arrive
- Added a bounded process-local cache (four clips, each at most 8 MB) keyed by encounter, speaker and exact text. Replays reuse synthesis. Voice preparation never advances the clock or mutates the save.
- With voices enabled, the frontend requests preparation when the director has a queued encounter and no current conversation. The preparation endpoint returns readiness only, not the queued text or audio. The normal speech endpoint still accepts only the active event ID and now rechecks it after synthesis, so a decision during generation invalidates the late clip server-side as well as in the browser.
- Tests verify stable speaker/text forwarding, prepared-cache reuse, no early access to queued clips, no save/time change, no access after an encounter ends, and rejection when a conversation ends during synthesis. Go race suite, all five existing browser-audio lifecycle unit tests, and production build pass.
- Added an isolated `voice` fixture. Live local provider preparation took 4.35 seconds. Enabled voices and bought coffee through the browser, causing the queued encounter to arrive; the next inspection showed Speaking. Two cached endpoint reads returned the 346,844-byte clip in 0.50 and 0.40 ms locally. These measure server delivery, not speaker-device latency or perceived voice quality. Declined the encounter and restored the fixture's voice preference to off.
- Authored sudden incidents still synthesize on demand, and cache entries are intentionally lost when the server restarts. This improves queued AI encounters and replay latency; it does not promise zero delay for every event.

## Missing-art fallback and recovery
- Acceptance review identified that failed embedded artwork could leave an unexplained blank canvas. Added a loading/error panel with address-book fallback and artwork retry. The parent accepts status only from its own same-origin frame; manifest/image/script failures are reported or time out visibly.
- Moved the ready message after successful asset loading, so failed artwork cannot accidentally dismiss the error state. Retrying remounts only the street frame, not the campaign. Entering the directory clears presentation-only travel/scene state.
- Found the server silently replaced an explicit BLACK_LEDGER_WEB path with dist whenever dist existed. Corrected configuration precedence and added HTTP verification of an isolated root and explicit missing-asset 404s. Full Go race suite and production build passed.
- Created an isolated frontend copy with every PNG omitted and a separate save on port 8805. Browser immediately showed the unavailable-art explanation. Used Open address book, traveled to Saint Agnes, and completed a courier: 08:00/$90→09:00/$135, respect 2. The game remained readable and actionable with images absent.
- Restored PNGs only in that test copy and used Retry artwork. Browser screenshot verified the illustrated street returned at the same 09:00/$135 state. Existing failed sidebar thumbnails need a normal view remount/reload; the retry specifically reloads the street frame.
- Outstanding acceptance work remains the complete uninterrupted pacing run and visual layout/occlusion polish, including additional district art and clearer actors. Core validation does not substitute for these rendered/playability checks.

## Bounded director correction
- Formatting/validation mistakes previously discarded an entire generation without a correction attempt. The director now permits one additional generation carrying the validation failure and the same operation/contact constraints. It does not relax validation or execute partial output.
- Network/provider errors do not retry through this path. A new life or dead player prevents a correction request from being issued after the first failed attempt. Each request retains the existing timeout; no unbounded retry loop was added.
- HTTP provider-stub tests verify an overlong approach label is corrected into exactly one queued encounter, two invalid responses leave the queue empty, a 503 makes only one provider request, and money/time/revision remain unchanged. The full Go race suite passed.
- This verifies control flow with deterministic provider responses; no live-model quality improvement is claimed from this pass. The correction may still fail, in which case authored play remains available.

## Fresh uninterrupted campaign and resulting fixes — 2026-09-07
- Ran build a72d3e3 on isolated port 8806/save full-playtest-20260907-b.sqlite3 from 18:00:17–18:16:08 UTC (15m51). No implementation edits or state injections during the run. Final state: Day 2 10:40, Alex alive at 100 health, $155, 41 respect, 13 heat, laundry and garage, Leo, Ashbury apartment and one security detail. This is a measured opening/rise/consequence run, not proof of the desired 20–30-minute pacing or of death in this same run.
- Ordinary UI choices earned courier income, recruited Leo, established laundry protection, sought Russo backing against Bellandi, suffered 15-condition sabotage, repaired, expanded districts, rented better housing and hired security. Later paid Bellandi for a temporary understanding and bought garage protection. Read-only ledger inspection confirmed income, wages, repairs and choices persisted.
- Four live AI encounters arrived and completed: market mediation, laundry collection, courier transfer and garage mediation. Voices showed Speaking on later generated offers. Initial generation and its correction both failed with unknown beneficiary; later requests succeeded. Semantic defects remain: a lunch-based option at 02:35, an unsupported promise to discuss lower rent, weak follow-up acknowledgment, and repeated dispute structure. These are not solved by schema validation alone.
- Crew assignment now appears on Leo's People card using the same server-supplied action/availability as location panels. Browser verification on an isolated copy: Day 2 10:40/$155 → 10:55/$165, task due 12:40, button disabled while assigned. Screenshot reviewed at 1280×720.
- Replaced the time-consuming manager inspection with a zero-minute book review including actual/current maximum hourly income. Repairs now report actual restored condition rather than always claiming 40. Alliance feedback identifies the purchased introduction and both standing values without disclosing hidden retaliation.
- Director context now includes exact allowed beneficiary IDs and a readable clock; validation feedback names the rejected value and legal alternatives. Instructions explicitly prohibit time-of-day-dependent approach labels and unsupported rewards. This is a prompt/feedback improvement, not a claim of solved live-model semantic reliability.
- Full Go race suite and production frontend build pass. Added regression coverage for both alliance relationships, inspection without advancing economics/pressure, and correcting an invalid neutral faction value. The pacing target, full art coverage/occlusion and deeper reliable story continuity remain open acceptance work.

## Headless campaign runner — 2026-09-07
- Implemented `go run ./cmd/simulate` with worker, investor and reckless policies consuming a deliberately public-only state projection. Commands use production `core.Execute`; the runner never opens the live save. Optional deterministic validated proposals exercise job/approach paths without model latency. Individual seeds and traces reproduce policy outcomes under the same code version.
- Added tests for reproducibility across all strategy/director combinations, meaningful policy coverage, and independence from hidden retaliation. Full Go race suite passed. No frontend or live campaign changes in this pass.
- QA caught two harness issues before interpreting results: reckless initially kept traveling and missed the attack; fixed it to challenge once then stay home. Consecutive small xorshift seeds also skewed early mortality; use a documented Weyl stride between actual recorded seeds.
- Final baseline: 300 authored campaigns in 8.55 seconds and 300 fixture-provider campaigns in 10.62 seconds, 160-command cap, zero command/invariant errors. Reckless deaths 82/100 in each mode. All investors reached the casino: median 56 commands authored, 45 fixture. All workers/investors survived these specific policies. Investor final median cash $20,609/$23,135 over roughly ten game days. These are policy experiments, not population survival rates or human pacing evidence.
- Baseline summaries stored in `simulation-baseline.json`; usage, seed replay and limitations in `SIMULATION.md`. Primary new balance lead: after completing business expansion, the compliant investor policy becomes reliably wealthy with little danger. Investigate ongoing meaningful pressure and spending decisions, rather than blindly adjusting payouts from one aggregate.
- Remaining goal work includes live director consistency, richer political decisions, post-expansion balance, art coverage and rendered playtests. The runner accelerates rules experiments; it does not replace them or establish that the game is fun.

## Territorial pressure and defiant-policy comparison
- Investigation of the first simulation baseline found a real political selection bug: the highest-income owned business was selected before goodwill was checked. A friendly family in its district therefore suppressed demands by a hostile family in another owned district.
- Pressure now selects the highest-income eligible business among districts whose family still has a claim. Good standing continues to suppress that family's demands; peace across all owned districts remains valid. Existing demand costs, relationship effects, hidden retaliation and early-game pacing are unchanged.
- Regression tests cover a friendly casino alongside hostile laundry territory, the inverse, peace across both districts, and renewed demands after goodwill drops. Added a defiant investment policy, identical economic priorities but refusal of business demands, with explicit policy and reproducibility coverage.
- 100-seed authored comparison at 160 commands: compliant investor median demands increased 5→8, casino still reached at command 56. All runs remained valid and all reached the casino. New defiant comparison: casino at median command 59, 16 demands and 6 repairs versus compliant 56, 8 and 1. Neither policy died; these refusals currently risk businesses, not automatically personal assassination.
- Final cash differs ($19,647 compliant, $14,199 defiant), but so does elapsed game time (13,885 versus 11,245 minutes), so this is not a controlled income-rate comparison. Resistance consumes more decisions and repair effort. The larger post-expansion design still needs meaningful opportunities and political development; this fixes a territorial bug without manufacturing unavoidable punishment for successful diplomacy.
- Full Go race suite passes. No UI layout changes; rendering verification from the previous pass remains separate. Simulation reports are reproducible with the documented seeds and policies.

## Player-initiated diplomacy with both families
- Added an audience at Russo Motor Works and direct destination links from both family cards. The link selects a location; travel and the 45-minute meeting remain explicit game actions. Bellandi audiences retain their existing venue.
- Both meetings offer tribute, a faction courier favor, or leaving. Tribute adds 8 standing (capped), transfers $150 to the actual receiving family and cancels only its current plans. Other families and future demands remain. Old saved audiences without an actor still resolve as Bellandi.
- A favor opens an optional authored proposal through the existing validator, with the meeting's speaker and faction beneficiary. Hearing/declining costs no time or money. Accepted work uses normal reward, approach, heat/police and rival-standing rules. This adds actionable diplomacy; it is not a claim of AI-generated summits or new autonomous faction behavior.
- Tests cover territorial plan cancellation, correct cash recipient, duplicate/insufficient tribute, legacy saves, capped standing, optional favors and normal completion stakes. Full Go race suite and production frontend build pass.
- Browser QA on an isolated copy of the earned campaign: Families → Meet Elena → audience at the current garage. 10:40/$155 → 11:25/$184 from existing business income. Hear favor preserves 11:25/$184. Complete careful route → 12:40/$291, respect 41→44, heat stays 13, Russo +2→+8 and Bellandi 0→−3. Rendered dialogue inspected; all four choices and portrait visible. User campaign was not used for QA.

## Recorded director replay and repeated-story guard
- Added a bounded, strict JSON corpus reader and `-director replay -corpus ...` to the headless runner. Each recorded proposal is validated and queued once; no model calls or save access. Reports include corpus SHA-256 and per-campaign queued count. Seed/version/policy remain necessary for reproduction. Replay is mechanical coverage, not an assertion that the recorded prose fits every alternate campaign.
- A live qwen3:14b sample took 29.69 seconds. It repeated the earned campaign's garage dispute with a different speaker; stored it unedited as a known-bad semantic example. It passed the bounded proposal schema, illustrating why valid JSON is insufficient.
- Added production novelty rejection for repeated current-life memory titles or identical offers (normalized case/whitespace/outer quotes), including previously declined offers. The existing one-attempt correction receives the reason. Prior-life offers do not block new people. Distinct paraphrases and deeper semantic repetition remain unresolved; this is not a general story-quality validator.
- Corpus tests cover size/schema/unknown faction/trailing-data failures, fingerprints, finite consumption, input immutability and reproducibility. Novelty tests cover title/body repeats, distinct tasks and old lives. Full Go race suite passes.
- Replayed the live sample across 100 seeds for each of four policies with a 160-command cap. All command/invariant checks passed. Death before the first eligibility boundary means some reckless runs never queue the sample; queued counts make this visible. No frontend changes this pass.

## Mercer Exchange joins the painted street
- Added the existing Mercer painting as the fifth street landmark, with explicit position, scale, entrance and authored walkway. Runtime canvas compositing uses the existing SVG silhouette to remove its source checkerboard; no source bitmap was altered. Its directory preview and API location remain unchanged.
- Street inspection now uses cached sprite alpha masks and rendering depth, replacing broad rectangular hit boxes. Transparent foreground padding no longer claims another building's click. Extracted the picking function for tests covering overlap, alpha, boundaries, missing masks and input order. The standard frontend test command now includes these tests (7 total, passing).
- Browser QA on the isolated earned campaign shows the masked storefront beside the laundry with its entrance ring. Selected Mercer through the visible street control, traveled from the garage: 12:40/$291 → 13:15/$313, then normal investigation/low-profile actions appeared. Selection itself did not change state. This verifies the UI/API path; silhouette picking has focused unit coverage rather than a recorded physical canvas click in this pass.
- Inspected desktop and 980×800 framing: all five landmarks and selection controls remain visible, with ambience control wrapping at compact width. Restored viewport override. Final module reload produced no browser warnings/errors. Production build and frontend tests pass. Main server serves the updated dist without a core restart.
- Pier 14 still lacks painted street placement; broader district art, actor detail, route/occlusion polish and overall pacing remain open requirements. This pass expands the usable scene without changing simulation rules.

## Travel boundary and arrival feedback
- Fixed completed travel being discarded when a decision arrives at precisely the journey's end. If the full duration elapsed and the player is alive, arrival now persists alongside the decision. Earlier interruptions still keep the previous base location, and now explicitly report elapsed travel and where the player remains. Partial route resumption is not implemented by this change.
- Latest developments no longer filters all travel records and substitutes an old story. A normal arrival shows its destination and duration. A simultaneous danger/political event still wins headline priority, with arrival/interruption available in the additional records.
- Added core tests for a business demand at the exact arrival boundary, location persistence after payment, and an earlier interruption's location/time/explanation. Full Go race suite and frontend production build pass.
- Browser verified current arrival receipts on the isolated campaign: old Mercer arrival correctly displayed after reload; then Saint Agnes travel advanced 13:15/$313 → 13:30/$323 and Latest developments showed Arrived at Saint Agnes / 15 minutes, matching the active location panel. Boundary correctness is covered by core tests rather than claimed from this ordinary-trip browser check.

## Live director history evaluation and remaining semantic gap
- Compared three live generations using isolated copies of the same earned campaign at Day 2 10:40. Original duplicate guard allowed a casino payment story that largely copied the earlier laundry payment offer. A five-word passage overlap check catches that concrete renamed-copy pattern (65% overlap, minimum 20 unique shingles). Exact title/body checks remain. Tests distinguish copied passages from shared genre vocabulary.
- Expanded novelty checks to queued/current offers and recheck current state at commit time, so snapshot age cannot bypass a duplicate that became visible while generation ran. Checks still do not interpret every claim or prove semantic novelty.
- The next live sample paraphrased the same premise and promised an unsupported future benefit. Reduced duplicate prose in model context: older arrangements retain identity/title/operation/status/canonical result, while the full latest connection remains for callbacks. Recent world-change history excludes duplicated story quotes. Saved memory is unchanged; tests verify preservation.
- Third sample was less verbatim but invented weeks of delayed payments and failed to acknowledge the exact latest task. No claim that storytelling consistency is solved. Saved output excerpts in `director-history-evaluation.json` provide a reproducible evaluation record; these are individual samples, not a controlled model-quality benchmark.
- Full Go race suite passes. No user campaign changes during evaluation. Remaining live director work: stronger grounding of proposed circumstances versus past facts, meaningful callbacks, broader varied motives, and measuring a more capable/instruction-following model configuration if available. The goal remains incomplete.

## Installed model comparison
- Verified local Ollama inventory: qwen3:14b and qwen3:4b-instruct. Tested the smaller model through the real app preparation/validation path on another isolated copy of the same earned campaign used in the prior history evaluation.
- First 4B sample became ready in 8.07 seconds including request/poll overhead. Structurally valid, but copied prompt example labels unrelated to its collection task, promised a future market job, and did not acknowledge the exact completed mediation. It does not establish a story-quality improvement over the recorded 14B samples.
- Removed illustrative action labels from the director prompt and explicitly prohibited promises about future job location/subject/timing/reward. Second 4B sample stopped copying those examples and did not promise the next job, but used generic labels and invented last-week repairs without a faithful callback. These remain quality limitations, not solved requirements.
- Saved unedited queued text/choices in `director-model-evaluation.json`. Kept 14B as the default; model selection remains configurable. Focused server tests pass. No gameplay/saves changed in the user campaign during evaluation.

## Crew loyalty and recovery
- Leo now refuses new collection assignments below 30 loyalty. Existing assignments finish normally. Business protection retains its existing 50-loyalty threshold.
- Added a $40 bonus action (15 minutes, up to +25 loyalty, capped at 100), available beside delegation in People and current-location actions. Payment and loyalty commit together before time advances, so interruptions do not consume the payment without the bonus.
- People reports refusal explicitly. Crew do not automatically quit in this implementation; rebuilding loyalty remains possible.
- Added regression coverage for refusal without state mutation, bonus payment/recovery and resumed assignment, unaffordable/max-loyalty rejection, and capping. Full Go race suite, frontend build, and seven frontend tests pass. This change has not yet received a dedicated browser visual playtest.

## Crew UI and simulation follow-through
- Browser playtest on an isolated copy of the earned Day 2 campaign verified People displays both actions. Paying the bonus changed loyalty 65→90, clock 10:40→10:55, and cash $155→$125 ($40 expense offset by $10 owned-business income). Main campaign untouched.
- Investor/defiant public-state policies now prioritize an affordable bonus below 50 loyalty, preserving the ability to protect businesses as well as collect. Added tests for recovery then delegation and unaffordable recovery avoidance. A production midnight-billing test proves missed wages cause assignment refusal. Core/simulation race tests pass.
- Ran 100 authored campaigns (25 each investor, defiant, worker, reckless; up to 160 commands) after the change. Detailed local report: `.runtime/crew-simulation.json`. These ordinary strategies do not deliberately bankrupt their crew; the focused missed-wages test covers that edge rather than claiming the batch demonstrates it.

## Sustained political feuds
- Repeated refusal/backing of rivals at business demands can now create a hidden personal operation when the demanding family's standing falls to -40 or below. First refusal still causes business retaliation only. Each family has at most one active personal operation per life; another provocative decision after resolution can renew the threat.
- Generalized personal retaliation to either existing family. Completed favors that drive either rival to -30 now have the same consequence previously reserved for Bellandi. Contact warnings identify the actual family; bargaining text no longer falsely attributes every attack to Bellandi. Existing family-specific audiences can cancel the corresponding operations.
- Decision details disclose the general risk of a deepening feud, without announcing the existence/timing of a hidden hit. Core remains authoritative; no model can invent an immediate death through these choices.
- Added tests covering first versus sustained refusal for both families, deduplication, Russo warning attribution, actor-specific cancellation, and Russo feuds from favors. Full Go race suite passes.
- Matched 25-seed authored investor/defiant batches (160-command limit) produced zero command errors. Compliant investors: 0/25 deaths, median casino command 56. Defiant investors: 8/25 deaths, all reached casino (median command 62), versus 0/25 deaths before escalation. Summary saved in `feud-simulation-summary.json`. These policies do not proactively change plans after warnings; a human fairness/pacing playtest is still required. Peaceful late-game pressure remains a separate open design issue.
