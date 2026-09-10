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

## Warning-to-negotiation browser QA
- Added a repeatable `russo-warning` isolated fixture (apartment, discovered Russo operation, $300). Browser verified acknowledgment kept 10:30 paused; Families → Meet Elena → travel took 20 minutes; audience took 45 minutes; tribute left $150 and cancelled the threat. An hour wait reached 12:35 with health 100, beyond the former noon attack. This is a focused fixture, not evidence of an uninterrupted full campaign.
- Travel receipts buried the warning during that sequence. Added a separate known-threat reminder in the city with a Families shortcut, suppressing ordinary opportunity advertising while intelligence is active. Only discovered current-life plans are exposed; no due times, private IDs or undiscovered operations. Cancelled/resolved plans disappear. Old API snapshots remain compatible through an optional field.
- Core tests verify undiscovered/previous-life exclusion and cancellation cleanup. Full backend race suite and production frontend build pass. Browser visual QA caught overlap with the absolute city heading; corrected the reminder's spacing and rechecked the rendered result. The reminder remains visible after reload. The street framing still has substantial space above the buildings and could be tightened in a later art/layout pass.

## Recorded contact context and bounded follow-up scheduling
- Prepared follow-up scenes now store an explicit connection to a completed current-life arrangement: record ID, title, and canonical result. The model cannot supply this object. The dialogue displays it separately as completed work with this contact; voice narration still reads only the actual dialogue.
- Arrangement memory retains the parent link through save serialization. An already queued/current follow-up reserves its parent, and completed follow-ups release the next brief from mandatory same-contact continuation. This creates space for fresh threads rather than endlessly forcing the latest speaker. It does not guarantee that a fresh generation chooses another contact.
- Added tests for parent reservation, chain termination and serialization. Full Go race suite and frontend build pass.
- Live default-model generation on an isolated earned campaign saved the exact garage mediation result. The generated body again proposed a generic laundry collection without a coherent garage callback. This remains an explicit director-quality failure; the context panel does not solve semantic grounding. Browser verified the completed-record card renders separately and all four choices remain visible. Label avoids claiming the new prose causally follows from the earlier job.

## Focused director brief experiment
- Checked the actual local model runtime: qwen3:14b loaded with a 32768-token context window; its template supports the requested no-thinking mode. No small-context configuration error identified.
- Ran two focused direct requests against the same completed garage-tool mediation. First (5.77s, 464 prompt tokens) acknowledged the garage but invented $2,500 left behind and overlong labels. After explicitly tightening constraints, second (4.47s, 561 prompt tokens) correctly acknowledged the tool dispute and provided shorter labels, with a new garage payment request. These were iterative exploratory samples, not controlled comparisons. Raw parsed outputs/metrics saved in `director-focused-evaluation.json`.
- Added optional `BLACK_LEDGER_DIRECTOR_BRIEF=focused`, an embedded shorter general prompt and context restricted to current player identity, actual contact/places/property/faction data, up to eight arrangement summaries, and the required completed connection. All existing operation, identity, novelty and consequence validation remains in place; bounded correction retains its feedback.
- Default remains the existing full brief. The optional production mode has context tests and server tests passing but has not yet undergone a live multi-scenario campaign evaluation. The two direct samples were narrower than that mode's final generalized context, so their quality/timing must not be attributed to it.

## Focused mode through the production server
- Tested focused mode through the real server's preparation, validation and queue path on another isolated copy of the earned campaign. Full-city focused context still produced a generic laundry payment with no callback.
- Narrowed follow-up context to the actual speaker and places explicitly named in the previous offer/title, carrying current ownership/condition. Fresh requests still see the full city. Removed unrelated voice/color/render coordinates and historical offer summaries from this experimental brief; up to eight prior titles remain for novelty. Tests verify locality, current player ownership, fresh-city access and preservation of mandatory operation/correction data.
- Second live production sample remained at the garage, but did not acknowledge the completed tool mediation. It also contradicted creditor/debtor direction: mechanic expects payment, yet a choice demands money from the mechanic. Outputs retained in `director-focused-production-evaluation.json`.
- Focused mode remains opt-in and the live user campaign stays on the unchanged full default. Server tests pass. This experiment motivates an explicit narrative-role brief (payer/payee, task subject, factual callback) rather than continuing to add prose-only constraints. That is not implemented yet.

## Explicit experimental job roles and factual acknowledgement
- Added a structured narrative brief for focused mode: location, proposed premise, player task, and source/recipient roles for deliveries/collections. Collection explicitly moves payment from customer to establishment manager; mediation remains a peaceful shared-access negotiation. These are premises for the proposed operation, not historical events or new simulated cash transfers.
- Follow-up briefs supply a short operation-specific factual opening based on completed work. Focused mode requires that opening before accepting a proposal, using the existing single-correction limit. Remaining prose still needs evaluation; this check does not establish semantic consistency of the rest of the story.
- Added tests for collection direction/location, no invented fresh-job callback, and a server-level missing-opening correction that queues exactly one linked offer. Server race tests pass.
- Live role-brief run passed the opening check but failed contact identity. Replaced vague connection-validation feedback with exact expected/received speaker and beneficiary IDs. A second live run corrected a supplied display name (Mara Bell) to its required ID (mara), but then failed approach-label validation; no offer was queued and authored play remained available. This is correct failure handling, not a successful live narrative test.
- The repeated structural mistakes suggest using a constrained JSON response schema to keep identity/enum/length errors from consuming the correction budget. That has not yet been implemented. Focused mode remains experimental; no user save was used for these tests.

## Constrained director response schema
- Focused mode now passes an actual JSON schema as Ollama's `format`, following https://docs.ollama.com/capabilities/structured-outputs. Limits operation and contact/faction IDs to the current brief, required fields, allowed approach methods, item counts and text lengths. Core validation still rejects duplicate methods, invalid UTF-8 byte limits, stale state, novelty failures and unsupported consequences; schema constraints are not a semantic guarantee.
- Tests verify fixed follow-up identity versus fresh NPC choices and the actual focused HTTP request carrying a schema. Server race tests pass. Full default remains unchanged.
- Live schema-backed generation succeeded: prior mediation acknowledged, correct garage location, payment moves customer → manager, valid IDs and labels. It copied a negative instruction from the task brief, so that sentence was removed from the brief afterward. No claim of a retest of that small wording removal yet.
- Browser playtest: travel to Saint Agnes, coffee, accepted careful collection; reached a police interruption at 14:00 before reward. Paid $40 and received the correct $115 reward afterward. Saved authoritative completion and final state in `director-schema-evaluation.json`. This is one functional AI sequence, not proof of broad story variety or the full campaign pacing target.

## Response schema rollout
- Continued the schema-backed focused campaign after the completed garage collection. Next courier offer had no false prior-work connection. It remained structurally valid, though its recipient drifted from the supplied brief; focused narrative mode is not being promoted.
- Tested the response schema with the existing full prompt on another isolated earned-campaign copy. It queued a valid collection with short labels and correct IDs. The body was still a generic laundry request without a faithful callback. Evidence in `director-schema-rollout-evaluation.json` distinguishes mechanical reliability from prose quality.
- Enabled the same response schema for both modes; full remains the default narrative prompt. All authoritative validation, one correction limit and authored fallback remain. Full backend race suite passes; default-mode HTTP tests now assert that requests actually carry the schema. No changes to player economics or saves.

## Visual artifact and framing pass
- Replaced absolute city-header placement and accumulated artwork padding with a grid header in normal flow. Street/directory controls, opportunity prompts and known-threat reminders no longer reserve overlapping vertical offsets. The illustrated iframe follows the actual canvas aspect ratio; landmark controls remain in normal flow below it.
- Addressed crossing-road outlines being drawn into the junction by painting all outlines before asphalt. The painted asphalt was not seamless, so the renderer now maps it once over the ground instead of repeating a sheared pattern; reduced texture opacity to keep the enlarged grain restrained. Source images were not changed.
- Pedestrians now receive the nighttime brightness treatment; reverse-facing cars place their headlight glow at the correct end. Property panels without painted art now use the same restrained architectural-symbol treatment as the directory instead of the old green cartoon building.
- Browser inspected main-game daytime layout (no game commands issued), standalone night lighting, and final isolated street view. Road texture seams disappeared in the isolated screenshot; no browser errors/warnings. Production build and seven frontend tests pass. A viewport override did not produce the intended compact screenshot, so compact-width verification is not claimed for this pass. Pedestrian artwork/occlusion and missing district building art remain open visual work.

## Visual-agent handoff prepared
- User reassigned visual production to Claude/another agent; Codex retains gameplay, functional UX, AI stories, simulation and playtesting. Recorded amended ownership in `docs/GOAL.md` and root `AGENTS.md`; retained visual project requirements under the separate owner.
- Added `VISUAL_HANDOFF.md` with approved reference images, current defects, file map, presentation/API boundaries, asset/mask rules and integration instructions; `VISUAL_ACCEPTANCE.md` with visual/motion/accessibility checks; and ready-to-use `VISUAL_AGENT_PROMPT.md`.
- Created sibling Git worktree `mafia-game-visuals`, branch `codex/visual-handoff`, at c5700b7. Claude CLI is installed; no agent was launched. Main and visual worktrees are separate and clean.
- Verified `scripts/run-visual-preview.sh damage` from the new worktree: lockfile dependencies installed, frontend built, fresh isolated damaged-laundry save created, own server running on port 8840. HTTP checks confirmed owned laundry at 45 condition and successful frontend/art responses. Frontend tests pass. User campaign/port 8791 were not used for this fixture.
- Visual agent delivers commits/evidence for review, not direct merges. Parallel functional changes to shared React shell/types need coordination. The original goal controller text cannot be rewritten using its available status-only tool; project goal documents carry the user's amendment without falsely marking the gameplay objective complete.

## Active gameplay goal and arrival encounters — 2026-09-07

- The goal controller returned no existing goal, so created an active goal with the user's amended gameplay ownership. Visual production remains assigned to a separate agent; the handoff is prepared, not dispatched. Confirmed the existing `caffeinate -di` process is running.
- Prepared encounters now surface at completed travel boundaries. Previously an offer due during travel remained queued until an unrelated local action. Urgent incidents still take precedence, and delivery consumes neither extra time nor a second command.
- Added regression tests for already-due and exactly-on-arrival offers, immutable original state, no repeat after declining, and preserving business-pressure incidents ahead of queued offers.
- Validation: full `go test -race ./...` passed. Fifty authored campaigns (25 investor, 25 defiant; up to 160 commands) completed in 2.03 seconds with zero errors; investor deaths 0/25, defiant 8/25, matching the previous baseline. Authored runs do not validate live-model narrative quality. The arrival change has not yet been deployed to the user's running server; the user's save was not modified.
- Next: exercise prepared arrival encounters in an isolated browser campaign, then improve narrative continuity and run the complete rise-and-consequence playtest. The overall goal remains active.

## Arrival browser QA and director contact progression — 2026-09-07

- Previous goal turn classified as progress: committed arrival timing change and activated amended goal. This turn verified it in an isolated browser save on 8841: traveled Saint Agnes→Mariner at 08:15 with no early encounter, returned at 08:30 and saw the prepared scene without another action. Reload preserved the scene. Declining left cash $90, minute 510 and location bar unchanged. Inspected the rendered dialog; all choices visible, no browser errors/warnings.
- New generated offers now use established contacts and prefer those least recently featured, counting pending scenes. Mara is available initially, recruits at 30+ loyalty, family leaders at +6 goodwill. Mandatory saved follow-ups retain their original speaker. Added schema/prompt constraints and server-side rejection so a model ignoring the schema cannot introduce an unavailable speaker.
- Added tests for recruitment, low loyalty, relationship changes, pending-scene recency, old-life isolation, follow-up preservation and bounded correction of unavailable speaker output. Full Go race suite passed; reran server race tests after adding the correction test and they passed.
- Added safe `contact` QA fixture. Live qwen3:14b full-prompt generation on isolated port 8842 chose Leo and produced a mediation lead rather than another Mara offer. Browser accepted it: after travel and coffee, completed at 09:45 with cash $135, respect 11, health 100, Russo +6 and Bellandi −3. Evidence: `director-contact-evaluation.json`. One sample does not establish broad narrative quality.
- Open design gap: generated jobs mention locations outside the unlocked district; abstract job resolution currently does not require access to the prose location. Continued coherence, outcome variety and full uninterrupted campaign remain open. User server/save unchanged; new binaries are isolated QA only.

## Job location contract and live mismatch evidence — 2026-09-07

- Prior goal turn was progress: committed contact progression and browser-tested a generated Leo job. This turn addresses its inaccessible-location gap with typed model `location`, unlocked-place schema restrictions, core validation, canonical venue names in choice terms, and saved arrangement location. Old unlocated offers remain compatible.
- Focused briefs prefer structured location memory and cannot use old prose to unlock a district. The server checks that dialogue names the selected venue and rejects explicit locked venue names. Tests cover district expansion, unknown/missing/locked IDs, venue mismatch correction, saved completion, unchanged base/quoted duration, and old-offer compatibility. Full Go race suite passes.
- First live qwen3:14b sample selected laundry but described a garage. Added mandatory selected-name validation after observing this; preserved evidence in `director-location-evaluation.json`. Second sample correctly selected and named Mercer Exchange. Browser verified readable venue labels, accepted careful mediation, and reached 10:15 with $120 at Saint Agnes; saved completed memory retains market location. No browser errors/warnings.
- Story semantics remain incomplete: second sample proposes a territorial split between families that current mediation mechanics do not implement. Location correctness is not evidence of outcome coherence. Next work should address story promises versus supported consequences, then perform the complete fresh campaign playtest. No claim of overall completion.
- Changes are in isolated QA binaries (8843 initial / 8844 checked); user campaign/server untouched. Keep-awake process remains active. Visual production remains assigned to the other workstream.

## Temporary business ceasefires — 2026-09-07

- Previous goal turn was progress: job-location validation and live semantic mismatch evidence. Added an actual bounded political agreement through audiences, rather than treating model-written territorial divisions as accomplished facts.
- $100 buys 24 game hours without the negotiating family's business demands or sabotage. Removes only that family's pending sabotage, preserves personal hits/other-family plans/standing/ownership, expires exactly at the saved minute, renews from now without stacking, and clears on a new life. The Families panel displays its expiry and the director receives active agreement context.
- Regression tests cover payment atomicity, save cloning, scope separation, suppression and exact expiry, unaffordable/repeated purchases, and new-life isolation. Full Go race suite, frontend production build and seven frontend tests pass.
- Isolated browser QA on 8845: acknowledged Russo's personal warning, traveled to garage, held audience, bought ceasefire at Day 1 11:35. Cash $300→$200; panel visibly shows expiry Day 2 11:35 and personal-threat exclusion. Read-only save/API verification confirms Russo's personal hit and known warning remain. This fixture does not prove the full campaign arc.
- Minimal functional additions in shared `src/main.tsx`/`src/types.ts` recorded in visual handoff notes; no art or renderer changes. User save/server not changed; building dist makes the optional new status rendering available on the shared frontend, without enabling backend agreements on the old server.
- Open: model-written jobs can still promise unsupported political consequences. The new audience agreement is real, but generated jobs do not automatically create agreements. Need story/mechanical scope alignment and a full uninterrupted campaign, not merely more local tests. Goal remains active.

## Diplomatic campaign policy and fresh browser run — 2026-09-07

- Prior turn was progress: implemented/tested business ceasefires. Added a public-state diplomat simulation policy to exercise their integration with progression, including avoiding repurchase while active. Regression tests verify negotiations coexist with casino progression. Full Go race suite passes.
- 100 authored campaigns (50 investor + 50 diplomat, up to 220 commands) completed in 6.40 seconds, zero errors/deaths, all reached casino. Diplomat median casino command 61 vs investor 56. Policy results are not proof of optimal balance or live AI coherence.
- Started fresh browser campaign on isolated 8846 at 21:50:40 UTC, save `.runtime/campaign-review-20260907-c.sqlite3`, binary `.runtime/blackledger-business-truce`, server session 55135. Read Guide; completed two courier shifts and careful authored delivery; traveled and acquired laundry. Current checkpoint is in `campaign-review-20260907-c.json`. Continue this campaign rather than resetting; no full-duration or full-arc claim yet. User save untouched.

## Fresh 20-minute rise-and-consequence playtest — 2026-09-07

- Continued the current save on 8846; original tab 44 had closed, so reopened as tab45 without restarting the live server. Prior goal turn was progress (simulation implementation and opening campaign actions). This turn exercised the campaign rather than changing runtime mechanics.
- Whole run: 21:50:40–22:10:58 UTC, 20m18s elapsed including reading/reporting and task continuation. Fresh save, browser actions only, unchanged executable, no state injection/rerolls. This is not 20 minutes of continuous clicking and it includes a tab reopen. Evidence/history in `campaign-review-20260907-c.json`.
- Earned: laundry, Leo, two contacts, Ashbury access, apartment, paid overnight upkeep, two security details. Took Russo backing and suffered business sabotage; refused a later demand, provoked Bellandi, received a warning during rest, defended at home and survived at60health with both guards lost. Recovered to85health, funded repairs through work/collections, restored laundry, paid renewed demand and returned home. Final Day2 16:25: $159, respect36, health85, heat3, one repaired business, apartment, Leo65loyalty, no hired guards. No death occurred in this run; separate death/new-life tests remain separate evidence.
- No browser errors/warnings. The core economic/consequence loop works through this arc. Narrative acceptance fails: unsupported new ownership becomes follow-up premise, Elena offers a Bellandi-benefiting job without a grounded reason, labels sometimes mix languages, dialogue promises unimplemented deadlines/favors. Prior canonical results are visible but do not prevent prose corruption.
- Priorities for the next gameplay pass: prevent casual offers during known danger; preserve progress for routine conversation interruptions or clearly separate irreversible operational failure; constrain faction-speaker affiliation and unsupported historical assertions; expose guard contribution in sabotage outcomes. These findings are evidence for changes, not claims that the goal is complete. Keep this save for recovery/regression playtests.

## Known-danger pacing and defensive feedback — 2026-09-07

- Prior turn was progress: completed/documented a fresh 20m18s campaign, including repeatable interruption and defense-feedback findings. Addressed two findings without changing the preserved campaign's running executable or save.
- Routine offers now wait while current-life discovered threats exist. Queue entries remain saved and unconsumed; hidden plots and old-life threats do not suppress offers. Explicit audience actions still work. Regression tests cover postponement, save/release, timing/privacy and chosen negotiations.
- Sabotage results now name the available defending colonist and report actual condition loss prevented, clamped to remaining property condition. Found and fixed the minimum-damage edge case: a weak attack below 10 strength no longer becomes stronger because a defender is present. Tests cover ordinary defense, nearly destroyed property, weak attacks and zero damage.
- Full `go test -race ./...` passes. No frontend changes or new visual production. New behavior is source/test verified and not yet deployed to the user's server. Preserved full campaign remains on 8846. Remaining priorities: interrupted operations and unsupported AI assertions/affiliation, followed by live verification of the combined fixes. Goal remains active.

## Resume work after routine demands — 2026-09-07

- Prior goal turn made progress on known-danger pacing/defense feedback. Addressed the fresh campaign's lost-job-progress finding: business demands now suspend the original arrangement and its remaining minutes, then offer resume or abandon after the demand decision. Urgent warnings/attacks remain operational failures.
- Preserves selected approach, quoted reward, heat, location memory and police checks. Pending work is saved, blocks casual offers, clears on death/new life, and supports repeated interruptions and zero-minute remainder. No extra economic reward is granted before completion.
- Full Go race suite passed. Added further regression for repeated demands and reran core tests successfully. Coverage includes all approaches, exact-boundary completion, save cloning, replay rejection, abandonment, police payment and warning failures.
- Isolated browser fixture 8847, `.runtime/paused-job-qa-20260907.sqlite3`, session41413: paid the demand, inspected the remaining15-minute/$75 decision, reloaded, resumed. Finished at08:45 with$225 ($200 start−$60 demand+$75 job+$10 business income), no active event and no browser errors/warnings. UI uses existing scene components; no art changes.
- Preserved the 20-minute campaign executable/save unchanged. New code remains isolated QA/source, pending integration with the remaining AI consistency fixes. Goal remains active.

## Director affiliation and mixed-script choices — 2026-09-07

- Prior goal turn was progress: resumable arrangements verified in tests/browser. Added structural prevention for the campaign's Elena→Bellandi contradiction. Family leaders' ordinary offers must benefit their own family; fixers/crew retain both-family/neutral options. Context provides per-speaker IDs, single-speaker schema restricts beneficiary, and server validation catches providers ignoring schema.
- Old contradictory saves remain intact but are not selected as mandatory follow-up connections. Tests cover both leaders, neutral rejection for leaders, independent contacts, valid/invalid legacy connections, decoding constraints and correction-before-queue.
- Initial full test run caught changed correction wording lacking the actionable phrase expected by existing validation tests. Preserved actionable wording; full race suite then passed. Added mixed-script label rejection for the observed `Use the back通道` failure, keeping Latin accents; server race tests pass after that addition.
- Live qwen3:14b check on isolated 8848 (`leader-affiliation-qa-20260907.sqlite3`, session84030) produced Elena/Russo shift-hours mediation at Bluebird Laundry. Recorded `director-affiliation-evaluation.json`. This sample predates the later script guard and does not prove arbitrary narrative facts or promises are correct.
- Main user save/server unchanged; overall goal remains active. Remaining semantic issues include invented ownership and callbacks treating prior claims as facts. Prepared visual handoff ownership unchanged.

## Semantic-review feasibility — 2026-09-07

- Prior goal turn was progress: structural affiliation/script validation and local generation evidence. Tested a separate local-model reviewer before imposing another generation call on the game. Added reproducible `scripts/evaluate-story-review.py`; it talks only to the local model and never reads/writes campaign saves.
- First six curated cases: 5/6 correct. Missed a completed-job claim after an explicitly declined courier. Renamed the ambiguous future `completion` fact to `allowed_result_if_current_job_completes` and supplied an explicit prior `completed` boolean. Re-evaluation: 6/6; separate six-case holdout: 6/6. Preserved all three results (`story-review-feasibility.json`, `story-review-explicit-facts.json`, `story-review-holdout.json`).
- Holdout covered paused work falsely claimed complete, truthful acknowledgement of a decline, a new supervisor without changed ownership, property promised as reward, an invented death attributed to an old delivery, and a legitimate future payment. Accepted cases approximately0.4–0.6s; rejections1.5–2.7s in these local samples. These are curated samples under short context, not statistical production benchmarks.
- Reviewer explanations are not wholly reliable: one correctly rejected paused-job case also falsely criticized negotiating loading hours, which is allowed. Initial ownership rejection similarly invented a detail about mediation participants. This matters because unfiltered critique could mislead a correction prompt.
- No production gate enabled and no claim of solved coherence. Next: integrate an opt-in reviewer with authoritative state/conditional effects clearly separated, measure its false rejections against real proposals, and bound correction/failure behavior. The remaining goal requires coherent live stories, not simply a passing synthetic suite. User campaign and preserved playtest unchanged; keep-awake verified.

## Actual campaign story-review evaluation
- Evaluated seven untouched local-AI offer bodies from the preserved 20m18s campaign, plus two marked authored controls. Per-case facts describe the offer time; current-job final statuses and expected labels are withheld. Two ambiguous cases remain unscored. The save was opened read-only.
- The simple reviewer accepted everything (3/7 scored matches). An evidence-first audit found the labeled bad offers but rejected both valid controls and produced one invalid excerpt (4/7 matches, one invalid response). Some matching rejections used incorrect reasons. Neither variant is suitable as a production gate.
- Added reproducible corpus/results and `--cases` / `--audit` modes; Python compilation and diff checks pass. Raw parsed output is now retained for future invalid responses; the recorded audit predates that last diagnostic improvement. See STORY_REVIEW.md for limitations and next experiments.
- Confirmed the active task goal matches the visual-handoff scope and clarified current gameplay acceptance priorities in GOAL.md. Keep-awake process remains active.

## Discard drafts prepared against obsolete world facts
- Revalidate property ownership, speaker identity/affiliation and fresh-contact eligibility inside the save transaction before adding a generated offer. Ordinary clock/cash/damage changes are allowed; completed continuations can survive lower goodwill if their canonical saved result remains unchanged.
- Stale drafts return a distinct error and do not retry the obsolete snapshot. The background director returns to available, or ready if an earlier offer remains queued. Player actions are preserved. Existing queued offers are not retroactively checked by this change; semantic claims in current dialogue remain a separate unresolved issue.
- Full `go test -race ./...` passes. Regression tests simulate a property acquisition during the model HTTP request and verify no offer is queued, no purchase is overwritten and no obsolete retry is made; also cover departed/disloyal contacts, changed identity/leadership, legitimate progress and canonical follow-ups.
- No user server restart or user-save mutation. This source change requires an integrated isolated playtest before claiming it in the running main preview.

## Integrated fresh campaign: rise, interruption, ceasefire and death
- Prior goal turn was progress: stale-draft protection and real-campaign semantic-review failures committed. Built isolated `.runtime/blackledger-integrated-40ca80b` on8849/session2647, fresh `.runtime/integrated-campaign-20260907.sqlite3`. Build and frontend passed; all seven frontend tests passed. Keep-awake process30528 confirmed active.
- Played normal browser commands from Alex's$90 opening through first AI job, Bluebird ownership, Leo recruitment, warning contacts, sabotage defense, Ashbury expansion/apartment, wages, home guards and further AI work. One demand naturally paused a collection with5minutes remaining; reload/resume led to police payment and exactly one original$130 reward.
- Bought business ceasefire and verified public expiry; provocation still brought a personal threat. Explicitly prepared another AI encounter during known danger, confirmed public ready state, then rest led to the attack rather than routine dialogue. Holding the entrance with two guards was fatal. Reload preserved death; new Nico starts$90 with no old crew/security/assets/truce/queued offers while Alex's former business and history persist.
- First-life wall duration667.68seconds (11m08s), including reading/reporting. This is an integrated regression campaign, not another20-minute run. See integrated-campaign-20260907.json for initial/death/new-life snapshots, invariants and findings. No hidden schedules were used to choose play actions; private read-only verification was done only after the first life ended and new life began.
- Narrative acceptance still fails: weak explicit callbacks, generic supplier backstories and favors, crew addressing player as employer, and an invented agreement with a rival outfit. Mechanical integration passing does not prove coherent evolving stories.

## Fixes from integrated playtest
- Police-delayed completion now uses the original saved job title in the ledger, retaining a neutral fallback when legacy saves lack the linked memory. Regression checks completion after cloning the saved police decision; full Go race suite passes. No prior history is rewritten.
- Families now shows public known threats above the agreement cards and avoids generic relationship copy that implied no active situation. People labels faction leaders by their organization and the detective as city authority. No new hidden information or API field. Existing threat style needed a local margin override outside the City layout; verified the corrected placement in the browser.
- Browser verification used the existing8845 Russo-warning/ceasefire fixture without advancing it. Threat and ceasefire are visible together, faction affiliations correct, no warnings/errors. Frontend build and seven tests pass. Shared `src/main.tsx` functional changes are noted for the visual owner.
- User save and main Go process remain untouched. The shared frontend bundle was rebuilt after the integrated campaign finished; it can appear on reload. The main Go process still predates recent gameplay changes.

## Reasoning-mode narrative experiment
- Prior goal turn was progress: completed integrated campaign through death/new life and fixed known-threat/history feedback. Rechecked source361bc92, clean worktree and keep-awake.
- Added offline reviewer `--think` with4,096-token limit and provider counters. Same real-campaign corpus/simple prompt:5/7 scored matches, no malformed responses, versus3/7 fast mode. Missed wrong-family leader and unsupported deadline; two ambiguous cases unscored. Median21.00seconds, range13.55–45.72seconds. Not suitable as a production gate; evidence saved in story-review-campaign-thinking.json.
- Added opt-in `BLACK_LEDGER_DIRECTOR_THINK=1`; default remains fast. Reasoning gets4,096 tokens,180-second asynchronous request limit and128KiB response cap; fast remains700/100seconds/20,000bytes. Structural validation, bounded correction and stale-save checks unchanged. Test covers both options,30KB non-dialogue draft, final-content-only storage, and unchanged cash/clock. Full Go race suite passed before final response-limit adjustment; command-server race tests passed after it.
- Live initial attempt8850 on an isolated copy of earned Day2 10:40 state timed out at100seconds with no queued job. Kept that evidence. Second8851 copy/binary with extended bound prepared one job in40.08seconds observed polling time. No claim that timeout change caused faster generation: output is stochastic.
- The generated collection preserves neutral Mara/garage context and recalls the tool mediation, unlike prior samples that lost that connection. Still muddled deposit/payment premise and weak causal language; not proof of general coherence. Saved proposal, previous memory, failed attempt and successful result in director-thinking-evaluation.json.
- Browser8851: wait60, cautious collection120, reload police decision, pay40. Cash344, respect44, heat5 atDay2 13:40; original job title retained in ledger. Source baseline and user campaign untouched. No second-model review enabled.

## Restore useful location context on reload
- Browser playtest showed reloading an ongoing job reset inspection to Saint Agnes even while the player was at the garage. Boot now selects the current public player location and uses the directory for locations without street art. A new person with no completed courier work in the starting room still gets Saint Agnes as the first lead. This is frontend selection only, not travel or a save mutation.
- Production frontend build passes. Browser8851 opened directly at owned Russo Motor Works; reopened preserved newcomer8849 and verified Nico's Saint Agnes lead remains. Original tab47 had closed; reopened its live server rather than restarting it. Shared functional main.tsx change noted for visual owner.

## Live director scenario suite and faction consistency
- Added an explicitly opted-in local-provider suite with four isolated saves: crew after a declined job, family leader work, a completed tool-dispute callback, and a new person after death. Reports retain the supplied facts and actual prose; ordinary tests skip provider calls. Report paths must be new to preserve previous evidence.
- Qwen3:14b reasoning-mode run: four mechanically valid offers in 138.83 seconds overall (24.39–52.21 seconds each). This was a narrative failure: wrong-family dialogue, missing callback, unsupported deadlines and an invented prior agreement. See director-thinking-scenarios.json. Reasoning remains opt-in.
- Added a narrow guard requiring non-neutral spoken offers to name their structured beneficiary, with word boundaries and punctuation normalization. A Russo reward with Bellandi-only dialogue now triggers the bounded correction path. Both prompts explain the requirement. This does not prove semantic correctness if both families are mentioned. Tests cover the observed failure, punctuation, substring rejection, neutral work, and successful correction.
- Full Go race suite passed after the guard/prompt changes. Installed official Ollama qwen3.5:9b locally for an isolated comparison; production defaults unchanged. Four non-thinking cases completed in54.39seconds (9.01–17.89seconds each), all mechanically ready, with two venue correction attempts. See director-qwen35-fast-scenarios.json.
- The alternative also fails narrative acceptance: leader addresses the player as Elena and invents a previous room-dispute meeting; tool callback becomes an unsupported second payment with a noon deadline; new-life offer recognizes Nico but invents a sundown deadline. Faster here is not evidence of better storytelling. Modes and prompt/guard revisions differ between runs, so this is an exploratory comparison, not a controlled ranking. No user save accessed or default model changed.

## Attribute director history and dialogue roles
- Previous turn was progress: committed bb5eb3b with live narrative failure evidence and the faction mention guard. Verified clean starting tree and active caffeinate.
- Separated current-person history/arrangements from previous-person city history and labeled each with its participant. Unknown historical lives remain explicitly unknown rather than being assigned to the current player. Retained the city history and canonical outcomes.
- Callback prose is now explicitly named original_request_claims_not_verified; it identifies the prior subject without becoming evidence that every allegation occurred. Saved offer/result records remain unchanged. Explicit addressee and speaker relationship metadata distinguishes the player, their employee and a family leader. This private context decoration also applies to focused mode, which now receives the attributed history in addition to its brief. No public API or save migration.
- Tests exercise death/new life attribution, preserved history/results, unchanged original memory, and crew-to-boss dialogue roles. Full Go race suite passed; command-server race suite repeated after a helper reuse and stronger history assertion. Live Qwen3:14b reasoning scenarios are evaluated separately below.
- Live callbacks still failed despite clearer provenance: the tool-dispute continuation invented a $25 payment, a particular mechanic owing it, and a noon deadline. The generated body is preserved verbatim in director-attributed-scenarios.json; this is not accepted storytelling. The next narrow guard should reject explicit unsupplied monetary/time terms while broader narrative evaluation remains open.
- Four attributed-context scenarios finished in254.39seconds (23.20–113.48seconds each), mechanically ready with two venue corrections. New-life offer addressed Nico correctly and did not claim Alex's old work; its aggressive choice nevertheless changed the objective to exclusive access. This suite fails coherence overall. Manual assessment and observed retry reasons accompany the final offers; rejected first-attempt prose was not captured.
- Main server, default model/mode and user campaign remain untouched. The clearer context is committed source, not evidence that all current AI offers are coherent or that the main preview runs this source.

## Reject unsupported spoken payment and time terms
- Previous turn was progress:2e2d4ff attributed current/former protagonist context and recorded continued narrative failures. Rechecked clean source and active caffeinate.
- Added pre-queue lexical checks for explicit currency amounts, numeric/spelled durations, named times, clock times and common closing/end-of-period deadlines in title/body/approach labels. Preserves ordinary place numbers, shift-hour disputes and action ordering. Uses the existing single correction, with no prose rewriting or gameplay mutation. This is not a semantic guarantee.
- Regression cases include the observed $25/noon offer, spelled amounts, currency codes, choice deadlines, acceptable Pier14/shift wording, successful correction and repeated failure with unchanged clock/cash. Initial tests exposed missing spelled clock times and a test fixture using the wrong requested operation; fixed both. Full Go race suite passed.
- Default fast14b live suite completed in39.82seconds:3ready,1rejected after both attempts omitted the beneficiary name. The callback still invented Tom, lost its prior subject and used 'before the end of the day'; this deadline phrasing was added to the check and command-server race suite passed afterward. The raw report predates that final deadline fix. See director-terms-fast-scenarios.json; this is a failed narrative suite, not a successful acceptance run. Main server/save untouched.
- Preparing a larger local candidate after repeated identity/callback failures across14b and9b. Official Ollama qwen3.5:35b-a3b tag lists24GB Q4_K_M (https://ollama.com/library/qwen3.5:35b-a3b); local disk had160.7GiB free. Download is an experiment, not a default-model switch or quality claim.

## Preserve rejected live-provider evidence
- Added a test-only local proxy to the opted-in scenario harness, recording every completed provider attempt's final content, HTTP status, token/timing counters and request fingerprint. Leaves provider status/body and cancellation behavior intact; omits reasoning text. This fixes the earlier diagnostic gap where only final queued prose survived.
- Ordinary tests still skip live models; a recorder regression verifies a503 response passes through unchanged while its final content is retained and reasoning excluded. Command-server race suite passed. Brief mode is now recorded alongside model/thinking mode.

## Larger local model and response-budget handling
- Installed qwen3.5:35b-a3b (22.23GiB local artifact, digest3460ffeede5453ead027dbd2f821b12ad0aa3de54630971993babdb2165221f7). Runtime reported21.8GiB loaded with32,768 context; no default model change. Source for the official tag: https://ollama.com/library/qwen3.5:35b-a3b.
- Fast trial atf29be37: four mechanically ready cases, six provider attempts,37.79seconds overall. Better explicit boss/Nico addressing in these samples, but leader invented a prior toll agreement and callback both invented previously secured parts and reversed debtor/creditor. This fails narrative acceptance. Complete attempt evidence in director-qwen35-moe-fast-scenarios.json.
- Reasoning trial: first response exhausted4,096 tokens with empty content in63.04seconds; the correction then hit180-second timeout. Second case logged another incomplete JSON response, but its exact done_reason was not captured. Deliberately stopped only the identified benchmark child71575 after436.45seconds and verified terminal exit; remaining cases are untested. Partial evidence is explicitly marked interrupted in director-qwen35-moe-thinking-scenarios.json. No claim that the second error was definitely token exhaustion.
- Model done_reason=length now returns an ordinary provider-budget error rather than triggering a futile same-budget story correction. Both fast and reasoning regression cases verify one call, no offer, unchanged clock/cash. Ordinary malformed complete responses still use bounded correction.
- Fixed the test recorder's cancelled-attempt race by awaiting active handlers before snapshotting attempts. Cancellation regression initially exposed a fake provider that did not consume its POST body; draining that body allows cancellation to reach it correctly. Targeted race tests now verify the cancelled attempt is retained. Older partial report retains its missing timeout-attempt row honestly.
- Narrative quality remains incomplete; the next gameplay integration pass should also verify the current accumulated core/UI changes in an isolated latest-source build, rather than allowing model comparisons to indefinitely postpone playable integration.

## Latest-source integration and upgrade validation
- Previous turn was progress:41a554f completed price/time guards, failure-budget handling and larger-model evidence. Verified clean source and caffeinate, built the accumulated core and production frontend. Full Go race suite, seven frontend tests and TypeScript/Vite build pass.
- Created .runtime/integrated-upgrade-20260907.sqlite3 via SQLite backup from the live save, without changing it. Latest core41a554f on8852 preserved id/revision190/life2/minute11365/Nico stats/factions/ownership/history/event; only optional business_truces public data was added.
- Browser8852: Russo audience45minutes,100-dollar business ceasefire, reload. Cash665 atDay8 22:10, standing+9 unchanged, ceasefire visible throughDay9 22:10.
- Fresh paused-job fixture8853: paid60 to settle demand, reloaded15-minute remainder, resumed to08:45/cash225/respect3/heat3. Ledger retains original A sealed message and single75-dollar completion. No live-save QA commands. Evidence in integrated-upgrade-20260907.json.
- Families had a confusing inherited-city title and unmarked second-person death cause from Alex's life while playing Nico. Changed to The city before you and labeled quoted final record; browser confirmed. Shared functional copy change noted for visual owner.
- Added read-only Go build identity to health response, so preview/core staleness can be checked directly. Unknown metadata is explicit; it does not identify frontend assets. Main preview update is prepared only after upgrade testing; deployment verification follows separately.

## Main preview updated to verified release b74c9a2
- Built from clean b74c9a2; health reports exact commit b74c9a2b2695c10da262b47cbf079f8915fd3d1d and modified=false. Launched8854 on another fresh copy of the live save and verified all prior public values match; only optional business_truces data is added.
- Created SQLite backup .runtime/backups/campaign-before-b74c9a2-20260907.sqlite3. Identified and stopped only previous main process57905 (.runtime/blackledger-schema-default), then started .runtime/blackledger-release-b74c9a2 on8791 against the original save. No QA actions were sent to that save.
- Verified the entire campaign state string and all190 receipt rows by SHA-256 before/after startup: byte-for-byte preserved. Nico remains revision190/life2/Day8 21:25/cash765/health100/respect76/heat9 at the garage with Cypress House as residence. Evidence in release-b74c9a2-20260907.json. Backup is retained; old binary is retained.
- Opened main preview in browser, confirmed current location selection, existing voice preference and no warnings/errors. Screenshot shows a readable directory at the current in-app size; it does not prove compact-device acceptance. Missing garage artwork and multiple vertical scrolling regions are visual-owner follow-ups, not completed visual production.
- Main now includes accumulated gameplay/narrative-validation fixes. AI default remains qwen3:14b fast; existing saved offers remain as stored. Narrative-coherence and full campaign acceptance remain open.

## Relocated to ~/Workarea/black-ledger and archived Afterlight
- User direction: move the mafia project into `~/Workarea` (their projects directory) and archive Afterlight so it is no longer around. Visual production was not resumed; this turn is repository organization only, with no gameplay, rules, API or save-format change.
- Stopped all 44 leftover processes first: 43 fixture/preview servers spanning ports 8792-8842 plus the 8791 main preview, then the private Ollama on 11435 and the voice service on 8787. SQLite integrity check passed after shutdown and a `.backup` snapshot was taken beforehand at `.runtime/backups/campaign-before-move-20260907-204752.sqlite3`.
- `mafia-game` moved to `~/Workarea/black-ledger`; `mafia-game-visuals` moved to `~/Workarea/black-ledger-visuals` and its worktree link was fixed with `git worktree repair`. Both worktrees resolve; `codex/visual-handoff` still has zero commits ahead of main. Requested worktree removal was declined by the sandbox, so the checkout was preserved rather than deleted.
- The project is now self-contained. `.tools/{ollama,models,voice-venv,audio-source}`, `director/` and the `.runtime/{speech,voice-cast}` caches moved in from the Afterlight folder, which had been hosting them. `.tools/` is git-ignored; `director/` is now tracked here. Added `scripts/run-services.sh` to start Ollama on 11435 and the voice service on 8787 from the project's own copies.
- Verified after the move: `go test -race ./...` passes for all four packages, seven frontend tests pass, the production build succeeds, `ollama list` shows all four models, qwen3:14b inference returns in 4.5s, and the voice service synthesizes a real 24kHz PCM WAV.
- Main preview restarted on 8791 from the identical `.runtime/blackledger-release-b74c9a2` binary rather than a rebuild, so the move stayed a pure relocation. Health still reports b74c9a2 with modified=false, and the campaign is unchanged: Nico Ward, revision 190, life 2, minute 11365, cash 765, health 100, respect 76, heat 9, at the garage, 190 receipts.
- Afterlight archived to `~/Workarea/_archive/afterlight-20260907.tar.gz` (153MB) with `afterlight-README.md`. Full Git history and all 933 tracked files; `gzip -t` and `git fsck` pass and a test extraction restored a tree with no differences against HEAD. Excluded only regenerable bulk: 4.4GB of purchased Unity packs re-importable from `~/Downloads`, Unity/Godot build output and caches, the Unity installer payload and the old `.tools` installers. The 8.9GB source folder was then removed; `~/Documents/ChatGPT/` is now empty.
- Left untouched: Afterlight saves in `~/Library/Application Support/Afterlight/` and the playable bundles in `~/Applications/`. Those bundles still launch but can no longer reach the AI director or voice, which now belong to this project.
- Codex's in-flight `cmd/blackledger/director_brief.go` mediation-premise edit and two uncommitted evaluation JSON reports were preserved as-is, not committed with this reorganization.

## Narrative guards for invented relationships, out-of-world voice and misassigned rank
- Assessed the four live proposals in director-qwen35-moe-focused-scenarios.json, which the previous session generated on qwen3.5:35b-a3b but never reviewed. Two are acceptable; two are not, in ways the existing guards did not detect.
- `leader_new_request`: Elena Russo opened "Alex, as leader of the Russo Outfit, I need you to step in", attaching her own rank to the listener, and then recited the brief's cast restriction aloud as "Neither is you, but the tension threatens our operations". The mediation brief's own sentences were being dramatized as dialogue.
- `new_person_after_death`: Mara greeted a brand-new life-two protagonist with "I need you back at The Mariner", inheriting the dead Alex's standing, and titled the scene with its own approach label, "Quietly listen to both sides".
- Root cause of the recitation is structural: constraints lived inside `proposed_premise` and `player_task`, the two fields the model is told to turn into prose. Added `constraints_never_spoken` to the brief and moved the cast and no-payment restrictions there, so the constraints survive instead of being deleted. The focused prompt now states that field is never quoted or referred to. The courier brief also stopped calling a character "the NPC", which was teaching the model out-of-world vocabulary.
- Added four guards in the existing one-file-per-check style, all wired into generateAttempt before the save boundary. `validateEarnedFamiliarity` rejects claimed prior dealings unless the speaker has an arrangement with the player in the current life, counting declined offers but never a previous protagonist's record. `validateInWorldVoice` rejects references to the player, the cast or the brief. `validatePlayerRole` rejects dialogue that gives the player a family rank they cannot hold. `validateSceneTitle` rejects an empty title or one that repeats an approach label.
- Evidence: `TestRecordedProposalsAreJudgedByCurrentGuards` replays the recorded run with no provider. Both defective proposals are now rejected, by four distinct guards with specific corrections; both acceptable proposals still pass every guard, so the checks do not fire on ordinary model prose. Unit tests additionally cover the exact observed strings and the near-miss phrasings that must stay legal, including "bring it back to the manager", "report back", "neither of them will budge" and a speaker naming their own rank.
- Full `go test -race ./...`, `go vet`, seven frontend tests and the production build pass. These are lexical and relational checks, not a semantic guarantee; a live suite and a fresh campaign are still required before claiming narrative acceptance.
- This resolves the previous session's uncommitted mediation-brief edit, which deleted both constraint sentences to stop the leak. Deleting them removes the constraint as well as the symptom, so the constraints were relocated instead.
- Chrome browser control became available mid-session and was used to check the running preview. The street canvas appearing blank and stuck behind "Opening Old Harbor…" was an artifact of screenshotting a backgrounded tab: `tick()` returns early while `document.hidden`, and an existing `visibilitychange` handler repaints on return. With the tab visible the loader clears and the canvas paints fully. No street defect was found and no visual-owned file was changed.

## Neutral work must not hand its people to a family
- Ran a fresh isolated campaign on 8860 against clean 20d1c76 with focused briefs, playing normally: travel to Saint Agnes, the authored courier job, then a requested encounter. This exercised the new guards in the real generation path rather than in unit tests only.
- The first generated offer was much improved: "Storeroom Dispute" titles the situation instead of repeating an approach label, opens with an earned callback to the courier job Mara actually gave, invents no amount or deadline, and breaks no frame. The previous session's failure modes did not recur.
- It still contained a defect no guard caught. The offer is mechanically neutral, with an empty beneficiary and no faction goodwill, yet the body describes "two Bellandi staff" at Saint Agnes, which the saved properties show is owned by independent. Dialogue implied family standing that the reward never moves.
- Added `validateFactionAttribution`: a family may be named freely, but the job's own people may not be given to a family that is neither the beneficiary nor the recorded owner of the job's location. Faction words come from the family name and both parts of the leader's name, so "Bellandi staff", "the Bellandi's men" and "Vittorio's people" are all covered. Ownership is read from saved state, so a property that changes hands changes what may be said about it.
- Tests cover the observed live body verbatim, attribution allowed for the beneficiary and for the owning family, ownership transfer, and ordinary mentions that must stay legal such as "Keep the Bellandi Family out of this" and "Vittorio Bellandi drinks here". The same live offer is also asserted to still pass the voice, title, terms and role guards, so the new check is not masking a broader rejection.
- Regenerated on the same save with the guard active: the model produced a clean neutral collection offer on the first attempt in about ten seconds, with no rejection logged. The guards constrain the director without starving it. Existing queued offers are not rechecked retroactively, so the earlier defective offer remains in that playtest save by design.
- Full `go test -race ./...` and `go vet` pass. This remains a lexical and relational check against saved ownership, not semantic proof.

## A city that fights its own wars
- Added `docs/LIVING_WORLD.md` as the working plan for the dynamic world, with an inbox section the user owns for unsorted ideas. Entries stay there in their original wording until genuinely implemented.
- Organizations now hold two properties each and have relations with each other, not only with the player. `Conflict` records hostility and a state of cold, feud or war between a pair; `FactionTurn` runs twice a day from `Advance`, and `FamilyDay` settles their books once a day.
- Wars are attritional. At war each side raids the other's weakest holding; a raid damages condition, cuts the defender's power and cash, and a holding below 15 condition held by the weaker side changes hands permanently. Every transfer is logged as public city news.
- Tuning took several passes, each measured over 200 campaigns of 180 days rather than guessed. A one-way hostility ratchet made war certain in 200 of 200 campaigns; zero-mean drift floored at zero made it impossible at 0 of 200; without hysteresis wars flickered off the same day they began and cost nobody ground. The committed model uses mean reversion toward an uneasy normal, hysteresis between entering and leaving war, roughly one raid a day, and exhaustion while fighting. Current behaviour: war in 123 of 200 campaigns, 10 by day 15 and 34 by day 30, 17 holdings changing hands, 3 organizations reduced to nothing. `TestCityConflictStaysVaried` asserts wide bounds around this so a later tuning change cannot quietly make the city degenerate in either direction.
- Ambient simulation draws from a reserved `WorldRandom` stream. Off-screen family politics previously shifted the player's own odds, which broke a seeded security test; player actions keep the original stream and that test passes unchanged.
- Removed the assumption that the city holds exactly two families. `Rival` returns the strongest other organization, provocation answers whoever actually holds the premises, and acquiring property angers its previous owner rather than a hardcoded index. Acquiring the garage no longer angers Russo specifically, because district 1 is unclaimed; that is now a consequence of holdings rather than a special case.
- Two player levers. `sabotage` was already committed; `incite` spends $25 and a reputation for discretion to make one family believe another moved against them, hardening that quarrel by 18 and sometimes starting a war the player is not part of. A story that fails to hold costs standing with the family it was told to. Both actions share one readiness check between the offered action and the committed command.
- The public API now reports `conflicts` (who is feuding or at war, and since when) and a readable `holder` per location. The interface showed ownership through a hardcoded Bellandi check, so Russo-held property would have displayed as independent; it now uses the reported holder. The Families screen lists the city's quarrels.
- Full `go test -race ./...`, `go vet`, seven frontend tests and the production build pass. A 200-command `cmd/apicheck` run on a fresh save reports no invariant failures, with idempotency and stale-revision rejection verified.
- Not yet built: factions created or destroyed during play, war affecting the player's own holdings directly, and director context carrying the conflict state. These are the next layers in LIVING_WORLD.md.

## Organizations the city creates for itself
- Weakened or embattled families now lose people. `Splinter` breaks a new organization out of a parent that holds at least two properties and is either beaten below three quarters of its peak strength or already at war. The breakaway takes the parent's weakest holding, a third of its strength and a sixth of its money, gets a generated name and a named leader who exists as an NPC the director can speak through, and starts at hostility 62 with its parent: a feud, not yet a war, so the parent has to decide whether it can afford to take the ground back.
- A family never loses its last holding this way, and healthy families at peace do not split. Names are checked against every living organization and person before use, so two organizations cannot share one.
- `dissolve` removes organizations holding nothing with no strength left, taking their quarrels and plots with them, while never reducing the city below two organizations. The people who led them remain as ordinary names.
- Measured over 200 campaigns of 180 days at the committed rate: war in 129, 18 by day 15 and 49 by day 30, 78 holdings changing hands, 72 new organizations formed across 51 campaigns, and 5 cities left with fewer than two organizations holding anything. The balance test now asserts churn as well as variety: the city must never routinely empty out, and it must be able to grow rather than only shrink.
- Fixed a real interaction between the two systems before it could reach a player. `factionWords` in the attribution guard took the first word of a family name to identify it. Organizations formed during play are named like "the Falcone Crew", so that word was "the", and the guard would have rejected ordinary dialogue containing "the men" or "the staff". Articles and the words every organization shares are now ignored, and identity comes from the distinctive parts of the name and the leader's. Regression covers both the false positives and that a new organization is still recognised by its own name.
- Full `go test -race ./...`, `go vet` and `gofmt` clean. A 60-campaign simulation on investor and defiant strategies reports no errors.

## The director is told what the city is doing
- Generation context now carries `city_conflicts` (who is feuding or at war, and since when) and `organization_holdings` (what each organization currently holds, by the names a character would say aloud). Both prompts state that these may be referred to because the simulation committed them, and that a war, truce, betrayal, breakaway or change of ownership may not be invented or resolved by the model. Organizations formed during play appear in these lists like any other, and speaker eligibility and beneficiary resolution already iterated factions generically, so a new leader becomes speakable without further change.
- The focused brief gained a `current_situation` field describing the standing quarrels the request happens inside.
- Honest result: the supplied context is correct and the guards accept the output, but the model does not yet build on it. Two live generations at the situated build produced ordinary, clean mediation and collection offers that ignored the feud entirely. Supplying world state is necessary but not sufficient; the operation set itself is still courier, mediation and collection, none of which is about a conflict. Scenarios that arise from the war need operations that only exist because of it. That is the next piece of work, not a claim about this one.

## Everyone with a name is a person who can die
- NPCs gained an organization, a place they are usually found, a standing inside that organization, ambition, skill and mortality. The field is `Dead` rather than `Alive` so saves written before people could die read back as living. Each established family is now seeded with a lieutenant and a soldier, so there is always someone who could replace the person above them.
- `Kill` ends a life for a stated cause and is the same function whoever did it and whoever it was. Losing ordinary members costs an organization strength; losing its leader triggers `Succeed`, which promotes the strongest surviving member, gives them the leader's standing and weakens the organization while the change settles. An organization whose people are all dead becomes a shell rather than silently continuing.
- Wars now reach people. A successful raid has an 18% chance of killing someone on the defending side, weighted toward the ranks that do the work but never excluding the top; `TestViolenceReachesEveryRank` asserts that leaders are sometimes caught and that soldiers are hit more often than leaders.
- `ConsiderInternalMove` gives an ambitious lieutenant in a failing organization a reason to take it by force. Nobody moves on a leader who is still winning. Over 400 forced attempts, 172 succeeded and 228 cost the challenger their life, so it is a real gamble rather than a scripted promotion. A breakaway now also takes people with it and its leader is a full person with rank, ambition and skill.
- Every person holds a distinct voice drawn from the local synthesis service's roster, chosen so no two living characters share one, and it stays with them as the city changes. This is the first part of the user's request that living NPCs each have their own persistent voice; portraits, style and individual behaviour remain open.
- Fixed a consequence before it could reach a player: the director selected speakers from the full NPC list without checking whether they were alive, so a murdered fixer or a dead family head could still have been the voice of a new offer. Eligibility now requires a living person, and a successor can speak for the family in their place. The public roster reports the living; the dead persist in the city's history.
- Balance after adding deaths, measured over 200 campaigns of 180 days: war in 121, 12 by day 15 and 32 by day 30, 288 holdings changing hands, 59 new organizations across 44 campaigns, and 7 cities left with fewer than two organizations. More violent than before, which is intended, and still not degenerate.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass. A 150-command `cmd/apicheck` run reports no invariant failures and the API lists eight named people with roles, organizations and locations.

## Upgrading an existing campaign to the living world
- The user's campaign was begun before the city had holdings, people or quarrels. Loading it on the current build would have left the new systems inert: Russo held nothing, so no war could start; no NPC had an organization or rank, so nobody could be promoted or killed for position; and no conflict existed to escalate.
- Added save version 3 and `MigrateLivingWorld`, run from `decode` for any save below the current version. It claims only premises whose owner is still `independent`, so a property the player owns or that a dead protagonist's estate still holds is never taken. It backfills holding income, faction peak strength, leaders' organization and rank, a lieutenant and a soldier per organization, a place and a voice for everyone, and the standing rivalry between the two established families. It is idempotent and changes no player value.
- `Antagonize` now sets the conflict's state from its hostility using the same hysteresis as the daily turn, so state and hostility can no longer disagree between turns. Previously a quarrel seeded at hostility 50 still read as cold until a turn ran, in a new campaign as well as a migrated one.
- Verified on a copy of the live save before touching it. Every identity and player field was byte-identical across the upgrade: id, life, minute 21425, revision 302, district, name, cash 5830, health, respect 183, heat, location, home, alive, earned, job count, contacts and security. Property changed only where it was unowned: bar and market to Russo, docks to Bellandi, while the player's casino and estate and Alex Varga's laundry and garage were untouched. History was preserved.
- Backup taken first at `.runtime/backups/campaign-before-living-20260907-223707.sqlite3` with 302 receipts. The live preview on 8791 was then stopped and restarted from the new build against the real save; the campaign reads back identically and the city now shows the two families feuding with nine named people in it. A generation on the migrated copy produced a clean, guard-passing offer addressed to Nico.
- Voices are now cast to match how a character is named. Generated people were previously given a voice by name hash alone, which put a woman's voice on Aldo Novak in the migrated save. Names and voices are kept in matched pools, checked across 480 generated people, and no two living characters share a voice.

## How a business is run became a decision
- Owning premises was passive income. Each trading business the player owns now has a standing operating mode, offered at the property alongside inspect and repair, costing no game time because it is a decision about the future rather than an errand.
- Run it clean takes 0.7 of what the premises could earn, wears nothing, draws no police attention and makes a family demand arrive later. Run it as usual is 1.0 and unremarkable. Skim what it will bear takes 1.5, adds a point of heat and a point of wear every day, and brings family demands sooner because the premises that are visibly earning are the ones worth demanding a share of. `BusinessPressure` now weighs targets by what they actually take rather than their nominal income, and its interval shortens with `SkimNotice`.
- Saves written before this keep earning exactly what they earned: an empty mode is the ordinary way, asserted by `TestOlderSavesKeepEarningWhatTheyEarned`.
- The first tuning was a trap rather than a bargain, and measuring caught it. At 3 condition of wear a day, a skimmed business destroyed itself inside 40 days and earned 22,431 against 24,936 for running clean, because income scales with condition. Wear is now 1 a day and ordinary trade wears nothing. Over 120 campaigns of 40 days: clean earns 24,936 at no heat and full condition, usual earns 35,880, skimming earns 43,266 at 80 heat and 40 condition lost. Skimming pays about 1.7 times running clean and requires repairs, lying low and an answer for the families, which is the intended shape. `TestSkimmingIsABargainNotAFreeLunch` holds those relationships without pinning the exact numbers.
- Full `go test -race ./...`, `go vet` and `gofmt` clean. A 150-command `cmd/apicheck` run on a fresh save reports no invariant failures with idempotency and stale-revision rejection verified.

## The underground trade
- Added a contraband market as the high-variance income path from the user's inbox. Two goods trade: moonshine at a base of $40 a crate, which draws police attention while held, and untaxed cigarettes at $22 a case, which are quiet and thinner. Prices move twice a day toward their base with noise, bounded to a third and triple of base, so a run of good prices is temporary and a crash recovers.
- Mercer Exchange deals in everything; Pier 14 deals only in moonshine. Trade happens in fixed lots of five because the interface offers plain actions rather than a quantity field. Buying costs cash and adds a point of attention; selling moves the entire holding at the current price. Positions scale with capital, since a lot can be bought repeatedly.
- Holding stock is what costs. Attention accrues every day it is carried, capped at six a day, and stops the moment it is sold. A police stop now searches the player and seizes everything they are carrying whichever way the stop is settled, which is what makes moving goods quickly the point rather than hoarding them.
- War makes goods dearer: while any two organizations are fighting, scarcity pushes prices up, so the city's own conflicts reach the player's ledger even when the quarrel is not theirs.
- Measured over 150 campaigns: buying and selling on a fixed rhythm while ignoring the price nets minus 32, while buying under base and selling over it nets plus 222, with stock genuinely held on 5,205 ticks. Profit comes from judgement rather than from the market existing, which is what `TestTheTradeRewardsJudgementRatherThanExistence` holds.
- Save version 4 seeds the market for existing campaigns at reference prices, so nobody inherits a windfall. A save with no recorded stock is carrying nothing.
- Verified through the API on a fresh save: a day-one player with $90 is correctly refused a $200 lot, and after six courier jobs a full cycle bought five crates at $33 for $165, raised attention to 1, offered the sell action with live pricing, and committed. Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass, and a 150-command `cmd/apicheck` run reports no invariant failures.
- The Ledger screen lists current prices against their usual level and what is being carried.

## A name and a price
- Added commissioned killings, the inbox's hitman system, built on the living-persons layer so anyone alive can be a target. Offered at Mercer Exchange, where information is already bought, and gated on having a contact who will carry it.
- Two steps, so the price is seen before anything is committed: choose a name from everyone alive except the player's own crew, then choose who does the work. Someone who needs the money costs the base rate, is caught half the time it goes wrong and does not hold their tongue; a professional costs 2.2 times and is caught a quarter of the time; a specialist costs 4.5 times and is rarely caught at all. What a name costs rises with the target's standing, competence and the strength of whoever protects them, so a family head is dear and a soldier is not.
- Nothing is decided when the money changes hands. A contract is committed with a due time between three and thirteen hours out and resolves on the clock, which is what makes a hit something that lands while the player is doing something else. The public projection never lists contracts, their timing or their targets; `TestAContractIsNeverVisibleToThePlayer` asserts the projection carries no contracts key, id or due time.
- How it was done is described by where the target was found: shot at the counter of Saint Agnes, held under a press at the laundry, crushed under a car at the garage, gone into the water off Pier 14. Two variants per location, asserted to differ and to name the place.
- A failed attempt is the real risk. The hitman may be taken alive and questioned, which puts 25 heat on the player and costs 45 standing with the target's organization along with a retaliation. Measured over 400 failed cheap attempts, 182 were traced back and 183 got away clean, matching the tier's stated capture rate. Over 300 attempts on a family head, the cheapest option killed 27 and a specialist killed 82, so paying more buys something without buying certainty.
- It runs both ways. `ConsiderFactionContracts` lets an organization at war buy the same service against the person at the top of the other side, priced and resolved by the same rules, and it will not spend money it does not have. Of 300 cities at war, 237 saw an organization commission a killing and 54 of those landed. Killing a leader through a contract promotes their deputy exactly as any other death does.
- `TestAMigratedCampaignCanActuallyFightAWar` was rewritten to measure across 40 seeds rather than one. Adding faction contracts drew from the world random stream and shifted that single seed, which is test fragility rather than a behaviour change; the balance test showed 359 seizures across 200 campaigns throughout. It now reports 199 holdings changing hands across 40 migrated cities and still asserts the player's residence is never a spoil.
- City balance after contracts: war in 130 of 200 campaigns, 9 by day 15 and 35 by day 30, 359 holdings changing hands, 65 new organizations across 49 campaigns, 4 cities left with fewer than two organizations.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass, and a 150-command `cmd/apicheck` run reports no invariant failures.

## Robbery, in both directions
- Added robbery from the user's inbox, the crudest money in the city and the least deniable. Offered at any premises with takings the player does not own, gated on being fit enough to try it.
- A successful robbery pays the day's cash scaled by the premises' condition, adds 10 heat and 2 respect, and wears the place slightly. If a family owns it, the money comes out of their cash, their standing with the player falls 25 and they answer it. A failure costs 10 to 30 health, 15 heat and standing, and can kill. Odds start near even, improve with reputation and an available loyal crew member, and fall against premises belonging to a strong organization.
- The city robs back. Twice a day, exposure is computed from what the player is visibly worth taking: carrying stock is the clearest signal, a great deal of cash adds to it, and an available loyal crew member or security at home halves it. Carrying goods is seized outright; otherwise cash is taken. Measured over 400 exposed players, 53 were robbed on a given tick, so carrying twenty crates and four thousand dollars around is dangerous without being certain doom.
- Nobody is robbed by a nobody. The thief is always a living, named person with a place in the city, weighted toward low standing and high ambition, and never the head of an organization in person. `TestNobodyIsRobbedByANobody` asserts that over 200 draws.
- Whether the player learns who did it depends on whether they built anyone who would tell them: with two or more contacts the name comes back within a day, and without them there is no face worth describing. That is the information network paying for itself, and it matches the user's note that not every situation has an answer the player can see.
- City balance is unchanged by this: war in 130 of 200 campaigns, 359 holdings changing hands, 65 new organizations across 49, 4 cities left with fewer than two organizations.
- Full `go test -race ./...`, `go vet`, `gofmt` clean, and a 150-command `cmd/apicheck` run reports no invariant failures. Driven through the API on a fresh save, robbing The Monarch took cash from $550 to $1,101 and raised heat from 0 to 10.

## The Bellwether Herald
- Added the newspaper from the user's inbox. Stories are filed at the moment events are committed rather than summarised afterwards, so a headline always refers to something that actually happened. War breaking out, a holding changing hands, an organization collapsing or splitting, a killing, an attempt that failed, an arrest, a robbery and damage to premises all run.
- The paper prints only what the city can see. It never names who arranged anything, and `TestThePaperNeverPrintsWhatIsHidden` asserts that a commissioned killing and a pending plot leave no trace in the edition, by id, by target or by the word contract. The player's own crimes are reported as unsolved: `TestThePlayersOwnCrimesAreReportedWithoutTheirName` checks the report never contains the player's name and reads as an open case, while their private history still records what they did.
- A new person picks up a paper that has been running the whole time, but the archive they can read starts when they arrive. The archive is bounded at 60 stories.
- Fixed a real leak found while writing the tests: `Contracts` was not cleared on death, unlike `Plots`, so a dead protagonist's paid arrangements accumulated in the save forever, filtered out of resolution but never removed.

## Making a commissioned killing legible
- The user reported taking out a hit and seeing nothing come of it. Investigating their save showed the contract had in fact resolved: day 21 at 17:29, "An attempt that failed". The system worked and the feedback did not. Two things were wrong.
- The player could not see what they had paid for. Contracts were hidden from the public projection entirely, which was the wrong line: the secrecy that matters is other people's arrangements, not your own. `PendingArrangements` now reports the player's own outstanding contracts with the target, who was hired and what it cost, and never the due time or anyone else's arrangement.
- The resolution did not say it was theirs. A contract the player paid for now reports as "Your arrangement is settled" or "Your arrangement failed", naming the fee, who was hired and, on success, how it was done.
- On the user's further note that headlines deserve more than a quiet log line: a failed attempt is now itself a headline, since a shooting that misses is still a shooting, and an arrest after one is bigger news again. The interface carries an unread count on the Herald tab and puts the newest headline on the city screen as a notice with a link to the paper. Read state lives in the browser, so reading the paper costs no game time and mutates no save.
- Verified through the API end to end on a fresh save: sixteen courier jobs and two contacts, then the contract flow offered three tiers with the two expensive ones correctly disabled at $790, the arrangement appeared under the player's own arrangements at $360 against Mara Bell, and waiting resolved it into the headline "ATTEMPT ON THE LIFE OF MARA BELL" with the arrangement cleared.
- On the user's note that the game should contain swearing as characteristic of the genre: both director prompts now state that characters may swear when it fits the person and the moment, kept in character rather than decorative, and never in place of the specifics of a request.

## Businesses with a character of their own
- The operating mode gave every business the same three decisions. This gives two of them something only they can do, which is the inbox's request that owning premises feel like more than a rate of return.
- A casino is a room you can walk into with money. Two stakes are offered at any casino the player does not own, since you cannot win money from your own house. Winning heavily comes out of the owner's cash and, if a family owns the room, costs standing with them and draws attention. A wrecked room runs no games.
- The first payout table was an exploit rather than a casino. Measured over 6,000 sessions it returned 8.9% more than it took, which is infinite money for anyone willing to sit there. Rebalanced and re-measured: the house now keeps 6.5% of everything staked, and of 2,000 sessions 1,269 lost against 715 won, so a night at the tables is usually a night down. `TestTheHouseKeepsItsEdge` holds the edge between 1% and 25%, so it can never invert again or become robbery.
- A laundry takes cash and gives back paperwork. Running takings through an owned cash business clears police attention for a cut, wearing the premises slightly, capped by their condition, with a day needed between rounds. A business below 40 condition cannot explain anything, and a casino is not a laundry. This is the first thing in the game that trades money for heat directly, which gives the crime systems somewhere to land.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass, and a 150-command `cmd/apicheck` run reports no invariant failures. Verified through the API: both stakes offered at The Monarch with correct costs, and a session committed.

## The police become a third force
- Heat was a number that only went up and did exactly one thing: trigger a police stop during a job above 15. It never faded, and nothing else in the game read it. Laundering had just made it a currency, so it needed to be worth spending money on.
- Attention now fades by 1 a day when nothing new is added, deliberately slower than a hard-run business generates, so choosing to skim still costs something rather than being quietly absorbed. `TestAttentionFadesWhenNothingIsAdded` asserts that relationship directly so the two systems cannot drift apart.
- Past 45 attention the police stop watching and start arriving, with the chance rising as it climbs. A raid seizes everything being carried, fines the player, wrecks the condition of their best-earning business and then eases off, because they got what they came for. Past 80 they stop taking money and take the premises: the business reverts to independent and is simply gone. Of 400 days, zero raids happened below the threshold and 35 happened at high attention, so it is a risk rather than a schedule.
- Detective Harlow can be paid to lose paperwork, at a price that rises with what there is to lose and clearing more when the player has contacts to route it through. Above 70 attention nobody will be seen taking it: money stops working before the danger does.
- The same pressure applies to organizations. A family at war is fined by the police every day it keeps fighting, and the raids are reported in the Herald. An organization at peace is left alone. This is the user's principle that everything runs both ways, applied to the law.
- Rebalanced the operating modes against this, because adding decay silently made skimming free. A hard-run business now generates 3 attention a day against 1 that fades. The old balance test asserted that skimming simply pays more, which is no longer the design: it pays more only if the attention is handled. Rewritten to measure three players over the same 120 campaigns of 40 days. Running clean earns 21,366 at no attention. Skimming and ignoring it earns 18,830, which is worse than clean, sits at 53 attention, wears the premises down to 29 and loses the business outright in 14 of 120 campaigns. Skimming and laundering it away earns 32,090 at no attention with the business intact in all 120. Money is available to someone who handles the consequences and not to someone who does not.
- Full `go test -race ./...`, `go vet`, `gofmt` clean; city balance unchanged at war in 130 of 200 campaigns; a 200-command `cmd/apicheck` run reports no invariant failures.

## Arms, and a correction to a number already reported
- Added the underground arms trade from the inbox, tied to the same list's request that weapons and armour affect survivability. Arms come off a boat at Pier 14 and nowhere else. Three weapons and two grades of armour, bought one step at a time, each costing more than the last.
- Equipment is not a stat line. A weapon shifts the odds of violence the player starts, by seven points a tier, in both sabotage and robbery. Armour reduces what violence costs when it goes wrong, never below a third of the injury, and improves the odds of surviving an unwarned attack at home. Of 400 unwarned attacks, an unarmoured player died 310 times and an armoured one 237: meaningful help, not safety.
- Carrying hardware is itself a reason to be looked at, and a police raid now seizes it along with the stock. Arms are handled separately from goods because they are evidence rather than merchandise.
- A measurement flaw found and fixed, affecting a number already recorded above. This generator produces nearly identical first draws for consecutive seeds: across seeds 1 to 400 the first value only moves from 0.236 to 0.39. Tests that measured an absolute rate off a first draw were therefore reporting the seed, not the behaviour. It surfaced as armour appearing to do nothing, with 400 of 400 dying either way, because every case fell in the same branch.
- Rate-measuring tests now stride their seeds the way `cmd/simulate` does. Most numbers barely moved, because they were paired comparisons under identical seeds: the hitman tiers went from 27 and 82 to 25 and 85 of 300, tracing a failed attempt from 182 and 183 to 179 and 188, faction contracts from 237 and 54 to 239 and 53, robbery from 53 to 49 of 400. One number was materially wrong and is corrected here: police raids at high attention are 93 of 400 days, roughly 23% and consistent with the intended chance, not the 35 recorded in the previous entry.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass, and a 150-command `cmd/apicheck` run reports no invariant failures.

## Work that only a conflict could create
- The recorded gap was that supplying the conflict state was necessary but not sufficient, because every operation was conflict-agnostic: courier, mediation and collection exist in any city on any day, so the director could be handed a war and still write an errand. Four operations now exist only because of something the simulation committed.
- `escort` and `warning` exist while two organizations are at war. `recovery` exists for three days after ground changed hands, drawn from the city's own news rather than a hidden record. `settlement` exists for three days after a leader is replaced and their organization still stands. Each carries the committed fact that justifies it, so the speaker can refer to a war or a seizure that actually happened.
- A quiet city offers none of them, asserted directly. In a city at war the rotation stays mixed: over twelve offers, six conflict jobs and six ordinary, with all five kinds appearing.
- Three defects found by running it live rather than reasoning about it.
- First, the tie-break. On a fresh campaign every operation has been seen zero times and the ordinary rotation held the tie, so a city already at war could never lead with a war story. The first live generation produced an ordinary mediation. Situational work now wins ties.
- Second, the model accepted the operation and wrote an errand anyway: a war-derived warning that never mentioned the war. `validateSituationalGrounding` now requires the body to name at least one subject of the fact the work exists because of, in the same spirit as the beneficiary mention check. The next generation produced "A Message Across the War... Bellandi to Russo", after one correction round.
- Third, and not the model's fault: the accept choice rendered blank. Its label came from a map covering only the three original operations, so every new operation produced an unpressable button and an empty ledger outcome. Labels and outcomes are now generated by functions with a fallback, and a test asserts every operation including unknown ones has words for both. `validateApproachLabels` separately rejects a generated approach with no usable label, which the model had also produced.
- Final live check on a city at war: operation `warning`, titled "Message to the other side", body grounded in finding somebody on the Russo side, four labelled choices and a proper outcome recorded.
- Full `go test -race ./...`, `go vet` and `gofmt` clean.

## Testing what had not been tested, and two charges taken twice
- This iteration added nothing to the game. Eight systems had landed in quick succession and the point was to find what they break together.
- The first finding was about the evidence rather than the game. `cmd/apicheck` climbs a progression ladder, and every system added since it was written is optional: contraband, operating modes, sabotage, incite, robbery, contracts, laundering, tables, bribery and arms are never required to make progress, so the harness had never once tried them. Six iterations of "no invariant failures" therefore covered fifteen kinds of command and said nothing about the rest. The claim was weaker than it sounded and is corrected here.
- The harness now mixes in ventures, prefers a system it has not yet exercised in the current run, follows multi-step events through rather than always taking the exit, and reports coverage alongside failures. A clean run that never tried a system has not tested it, and it now says so.
- That immediately found a real defect, in two places. The command layer charges what an action declares as its cost, and the choice layer charges what a choice declares. Handlers for arms, tables, laundering, bribery, contraband and contracts also charged, so the price was taken twice. Because commands are atomic the money was never actually lost, but the visible effect was that any of these was refused whenever the player's cash sat between the price and twice the price: a player with exactly $240 could not buy a $220 revolver, and one who could afford it was overcharged.
- Fixed by making the payment happen exactly once for each. The dynamically priced handlers keep charging and their actions declare no cost, with the price stated in the detail text and affordability still gated by their own readiness checks. `Commission` stops charging, because the choice that reaches it already carries the fee. `TestNothingIsChargedTwice` executes each through the real command path with cash set to just over the price, which is the case that failed.
- After the fix, a 700-command campaign over four lives exercised 32 kinds of command with no invariant failures, leaving only the bribe untried because its cost climbs steeply with attention and the campaign never held that much cash at once. Full `go test -race ./...`, `go vet` and `gofmt` clean.

## A city that could stop being a city
- A 1500-command campaign over six lives found a state the balance tests were not looking for. The Bellandi Family held every trading property in the city at power 100 with 59,671 in the bank, and the Russo Outfit held nothing, earned nothing, had no cash, and sat at power 58 which was climbing back toward its peak.
- The cause was in `FamilyDay`: an organization with no holdings had its strength drift toward its peak, because the target defaulted to peak when there was nothing to compute a target from. Strength comes from holdings, so an organization that holds nothing now fades by 4 a day instead of recovering to the strength it had when it owned half the waterfront.
- That produced a second, worse state. The faded organization dropped below the strength `considerReestablish` required, `dissolve` refuses to reduce the city below two organizations, and the result was one family owning everything beside a shell that could never act again. A city with one real organization cannot have a war, so nothing further could happen in it.
- Two changes. Recovery now needs people rather than strength: an organization reduced to nothing still has members who want somewhere to work, and it moves onto premises nobody holds, never the player's and never a rival's. And a family with nobody left to fight fractures from the inside, so total domination becomes the start of the next cycle rather than the end of the city.
- Measured over 200 campaigns of 180 days: organizations left holding nothing fell from 11 to 0, and cities left with fewer than two organizations from 4 to 0, with the war rate unchanged at 130. A repeat 1500-command campaign showed the cycle running: the Russo Outfit was destroyed and replaced by a breakaway, and the Herald recorded "END OF THE MARCHETTI COMBINE", "ALDO HALE KILLED" and "SAINT AGNES CHANGES HANDS" in one campaign.
- Full `go test -race ./...`, `go vet` and `gofmt` clean; 32 kinds of command exercised with no invariant failures.

## Something to leave behind
- Added the offshore account from the inbox, including the user's later note that reaching it as a new start should be difficult and need capital. It is the first mechanic that touches permanent death, so it was built to sharpen that rather than soften it.
- Everything a person owns still dies with them, and a new arrival still starts with ninety dollars, a rented room and no protection. What survives is money deliberately sent out of the city beforehand: wired in lots of $500 with the arrangement taking 18%, so banking is a loss taken on purpose. It earns nothing sitting there.
- The balance lives on the world rather than the player, so `new_life` resets the person and leaves it standing. Whether a given person can reach it is theirs alone and resets with every life. Establishing that the account is yours costs $400 and four hours, deliberately more than the $90 a new arrival is given, so an inheritance has to be earned before it can be collected and a small balance is not worth reaching at all.
- Tests hold the properties that matter: banking costs a cut and never accrues, the balance survives a death while access does not, a new person still starts with 90 and cannot draw on it, a penniless arrival is refused, an empty account offers neither access nor withdrawal, none of it is arranged anywhere but where money and paper already change hands, and a round trip through the account can never end with more money than it started with.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass. Verified through the API on a fresh save: twenty-two courier jobs to $1,065, then $500 wired out with $410 arriving, with access and withdrawal correctly refused until there was something to reach and somebody to prove it belonged to. `cmd/apicheck` now exercises all three actions, and the Ledger shows the balance and whether it answers to the current person.

## People with a day of their own
- Until now everyone except the player and the organizations was inert. People had ambition and a standing and did nothing with either: everything that happened to them was done to them. `PeopleDay` gives each living person one chance a day to pursue something, with the chance driven by their ambition so the people who want more are the ones the city hears about.
- They act through the same rules. Somebody helping themselves to a business takes the same money out of the same books, and the Herald reports it exactly as it reports the player doing it. Whoever is at the top of an organization does not do this personally, and nobody robs their own people.
- An unaffiliated person with the makings of an operator takes over premises nobody holds and starts being somebody: their standing rises and their role becomes running the place. Never the player's premises and never an organization's.
- The player's business is a target like any other, which is the half of this the player feels. Targets are weighted by what a place actually takes rather than always the richest, which was the first version: it meant the player was never robbed unless they owned the best business in the city, and the measurement showed exactly zero player robberies across 200 campaigns. Weighted, the player is robbed in 116 of 300 campaigns over 40 days, falling to 99 with security and a loyal crew.
- City balance is unchanged: war in 130 of 200 campaigns, no organization left holding nothing, no city reduced below two organizations.
- What a 1200-command campaign produced with nobody writing it: Mara Bell, the fixer the player starts with, ended up running Russo Motor Works. The Russo Outfit ceased to exist. The Bellandi Family split, and two organizations that did not exist at the start, the Delano Syndicate and the Amato Syndicate, each ended with a leader and a lieutenant of their own. The Herald carried the split, the end of Russo, open war on the waterfront and police pressure on three organizations.
- Full `go test -race ./...`, `go vet` and `gofmt` clean, 31 kinds of command exercised with no invariant failures. Five actions went untried in that run because the campaign never held the cash for them at the right moment.

## The inside of a business
- A business was a number that produced money while the player was elsewhere. Each trading business now has an inside: people who work it and have to be paid, stock or a float that runs down as it trades, and its own kind of trouble that stops it earning until somebody deals with it.
- What goes wrong is specific to the trade. A press breaks at the laundry, parts walk out of the garage store, and a dealer at the casino is working with somebody on the floor. Each has its own remedy, its own supplies and its own wage bill, asserted to be distinct.
- Capacity multiplies income alongside condition and operating mode. Short-handed, out of supplies or in trouble, a business earns less but never stops: a place limping is more interesting than one shut. Trouble finds the businesses nobody is watching, at 282 of 300 neglected against 139 of 300 well run.
- Wages are now part of what a day costs, so staffing is a real decision rather than a free upgrade, and letting somebody go is a genuine option when money is short.
- Two tests failed on the first build and both were right to. Businesses defaulted to no staff and no supplies, so acquiring one handed over an empty shell earning a sixth of its potential, and campaign income collapsed. A trading business is now a going concern from the moment the city exists: it has people and stock before anybody buys it, and ownership changes who answers for that rather than whether it exists. Save version 5 does the same for existing campaigns.
- The balance test then failed for a better reason: its managed player only laundered and repaired, which is no longer what managing means. With restocking, staffing and dealing with trouble included, passive clean ownership earns 8,123 over 120 campaigns of 40 days, careless skimming earns 6,560, and actively running the businesses earns 27,132. Attention is worth more than three times passive ownership, which is the shape the inbox asked for.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass, and a 900-command campaign reports no invariant failures. The property panel shows what a business is working at, who is on the books, what it has left and whether something is wrong.

## A source of your own
- The moonshine on the market came from nowhere: it could be bought and sold but never made. A still turns a business the player already owns into a source, which is the difference between trading on somebody else's supply and having one of your own.
- It can be hidden in a laundry or a garage, never on a casino floor and never in premises the player does not own, and somebody has to be working there to run it. It consumes the same supplies the business runs on, produces crates scaled by how well the business is working, and stops when there is nowhere left to put the output.
- The stock is the risk, not the still. Output arrives as goods the player is holding, and holding draws attention every day through the rule that already governs anything carried, on top of the attention the still draws by existing. A raid that finds one costs double the fine plus four hundred, takes the still, wrecks the premises and is reported in the Herald. Taking it out voluntarily is a way to stop being worth watching.
- People who rob premises would rather have crates than a till, so a still is a reason for the city's own thieves to come for the player: crates were taken in 92 of 300 campaigns over 30 days.
- Measured over 120 campaigns of 40 days: no still earns 7,329. A still whose owner never moves the stock earns 241 and is seized in 120 of 120 campaigns, because attention climbs until the police arrive. A still whose owner sells earns 13,594 at 41 attention, just under the threshold where raids begin. It nearly doubles income for somebody who moves what it makes and costs everything to somebody who does not.
- Full `go test -race ./...`, `go vet`, `gofmt`, seven frontend tests and the production build pass.

## The director learns the economy
- The economy the last few iterations built was invisible to the director: it knew about wars and holdings but nothing about the player's businesses, their stock or their attention, so it could not write from the situation the player was actually living in.
- `player_situation` now carries what a contact would plausibly know and the player can already see: what they run, what each place is working at, how it is staffed, what it is short of, what has gone wrong in it, how much stock they are holding, their police attention and whether they are armed. Both prompts state that a contact may refer to any of it and may not invent a business, a shortage or a quantity it does not show.
- Two operations exist because of that economy. `supply` exists when one of the player's businesses is out of what it runs on, short-handed or in trouble, and names the place and the problem. `distribution` exists when the player is holding ten units or more, and says how much. Neither can be selected by a player who owns nothing and carries nothing.
- Verified live against a player owning a laundry that was out of supplies, running a still and holding 22 crates. The director selected `distribution` and wrote from the actual position: stock sitting where it should not be, a buyer at the other end.
- That first generation exposed a defect the schema could not prevent. Approach labels are capped at 45 characters and the model wrote up to the cap, producing "Drive directly to the exchange and hand it to", cut mid-phrase and rendering as a broken button. A length limit cannot stop a dangling phrase, so `validateApproachLabels` now rejects a label ending on a preposition, article or conjunction, and both prompts state that labels are a few words comfortably inside the limit. The next generation produced four complete labels.
- Full `go test -race ./...`, `go vet` and `gofmt` clean.

## What you are wearing
- Respect is earned and cannot be bought. How the player looks can be, and in this city it is read first, so standing now has a second, purchasable half. Attire runs from working clothes through a pressed suit and a tailored one to bespoke, at $190, $680 and $1,900, worth 5, 11 and 19 presence.
- `Presence()` is respect plus what the clothes are currently worth, and it is what rooms judge. The high tables at a casino now want to see 14 before they will seat somebody, which a nobody in bespoke clears and a nobody in an off-the-rack suit does not. A tribute at an audience buys a quarter of the player's attire in extra standing on top of the flat eight. Intimidating a business off its takings reads presence rather than respect alone.
- It is deliberately a consumable rather than a purchase. Clothes lose two condition a day simply from being worn in this city, thirty to a beating, twenty to a search and fifty-five to an attack at the residence. Below forty condition a good suit is worth nothing at all, because a ruined one reads worse than honest working clothes, and the player is told the moment it crosses that line.
- The cost runs the other way too: a tailored suit draws one police attention a day and bespoke draws two, because dressing above your visible means is a question waiting to be asked. A suit already in rags draws none, which is one small mercy for somebody hiding.
- Pressing restores 45 condition for $25 anywhere the player lives, and for nothing at a laundry of their own — the only return a laundry ever gives that is not money. A suit taken below shabby never comes back past 80, so replacing it eventually becomes the cheaper answer.
- Measured over 200 campaigns: a neglected tailored suit stops being worth anything after 31 days, so it is a season's expense rather than a day-one purchase. Half a year of keeping one presentable at home costs 8 pressings and $200 and draws 180 police attention.
- Verified over the HTTP stack against fresh isolated saves. A 220-command campaign exercised 24 command kinds including `dress` with no invariant failures; a scripted campaign bought a pressed suit at the exchange, wore it to 96 over forty hours at home, and pressed it back to 100 for $25, with the snapshot reporting attire, condition, standing and presence at each step.
- Save version 6 puts an existing campaign's suit in good order rather than in rags. Clothes die with their owner: a new arrival inherits nothing.
- The director is told what the player is wearing and whether it is kept, so a contact can read the room the way anybody else would.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc` and the production build pass.

## A casino is a float, not a number
- A casino was premises with an hourly number on it, which is the one business in this city that number describes least well. A room is a float: the money behind the tables is what the house is willing to lose in a night, and it is the only reason anybody with real money walks in.
- An owned casino now runs a night every day. What crosses the tables is set by the room's condition, its staffing and stock, how hard it is being run, and above all by what is behind the tables — a thin float is not a secret, and the people who bet seriously know exactly which rooms can cover them. The house keeps six percent of that turnover in expectation, which is the same edge a player faces at somebody else's tables, and loses on 148 of 400 individual nights.
- A high roller arrives on nights the room can attract one. What they turn over is set by the room rather than by the float, because there are only so many people in this city who bet like this and they bet the same way wherever they sit down. So a thin float sees one rarely and cannot survive it; a deep one sees them often and barely notices. Over 30 days a $400 float cannot pay a winner in 41 of 400 campaigns; a $4,000 float in 4.
- A house that cannot pay loses the whole float, twenty-five condition and five respect, and the Herald runs it, because word of that travels faster than anything else in this city.
- Money won stays behind the tables. Putting $250 in and drawing $250 out are separate decisions, and every lot drawn is action the room can no longer attract, which is the tension the whole system exists for.
- **A deliberate rebalance, recorded because it changes an existing campaign.** The casino's hourly income used to include the tables; it is now the floor take alone and drops from $48 to $18 an hour, with save version 7 correcting an existing campaign to match. Over 40 days the floor take is $17,280. Dark tables add nothing. A $500 float earns $2,514 across 120 campaigns and cannot pay in 16 of them. A $3,000 float earns $10,044 and cannot pay in 1. So a passively owned casino is worth substantially less than it was, an actively funded one earns most of the difference back, and an under-funded one is the worst of both.
- Verified over the HTTP stack against a fresh isolated save: a campaign reopened the casino, found it earning $18 an hour with dark tables, funded the float to $1,509 to cover $2,400 of action a night, played a night out and drew $250 back into hand. A 200-command campaign exercised 32 command kinds including `bankroll` and `draw` with no invariant failures.
- The director is told what is behind the tables and treats a thin float as trouble at that business, so a contact can mention it the way anybody in the city would.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` (no errors across four policies) pass.

## Something to drive
- The clock is the one thing nobody in this city can buy more of: every journey is time in which businesses earn, rivals move and arrangements come due. A car buys that time back, and it is the most visible thing a person can own.
- Three of them, sold at Russo Motor Works and nowhere else: a used Ford at $620, a Hudson with a false floor at $1,750, an armoured Packard at $4,400. A journey takes 74%, 60% or 50% of what it takes on foot, scaled by how well the car is running rather than switched off at a threshold. Measured over 24 campaigns of 16 days working a fixed round of jobs: 5,064 jobs and $11,003 on foot against 5,736 jobs and $12,444 driving a Hudson — thirteen percent more work done, and the upkeep already paid out of it.
- The false floor is the second half of the trade. Stock under it draws no attention and a search does not find it, so the contraband system now distinguishes what is held from what is exposed. Over the limit is over the limit: with a 20-unit floor and 32 units held, a search takes exactly the 12 in the open.
- It costs on three sides. Upkeep of $4 to $12 a day joins the housing, security, crew and wages on the daily bill — halved if the player owns the motor works, where servicing is also free, which is the second business in the game to pay its owner in something other than money. Condition falls two a day, fifteen when a robbery goes wrong, forty in an attack at the residence, and below thirty a car is worth nothing to anybody. And a car outside is a thing witnesses describe: attention across 300 robberies runs 3,905 on foot, 4,381 in a Ford, 4,857 in a Hudson and 5,333 in a Packard, with the owner's goodwill falling further by the same measure. The best car in the city is not always the right car to commit a crime in.
- A warrant that turns up a false floor takes the car and everything in it, in 140 of 400 raids on somebody carrying with one.
- Verified over the HTTP stack against a fresh isolated save: a campaign expanded into Ashbury, bought up to the Hudson, watched its daily bill go from $15 to $22 and the walk to Saint Agnes fall from 50 minutes to 30, wore the car to 96 over a few days and had it serviced back to 100. A 200-command campaign exercised 34 command kinds including `service` with no invariant failures.
- Save version 8 puts an existing campaign's car in working order. Nobody inherits one.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` (no errors across four policies) pass.

## A place with something in it
- A residence was a rent line and a number of guards; everything else a person might do with the place they live in did not exist. Four comforts can now be fitted, each costing money once, costing money every day after, and doing exactly one thing the rest of the game already cares about.
- **A door that holds** ($420, $2/day) counts as another guard. **A telephone in the hall** ($300, $5/day) counts as another contact, for warnings and for hearing whose name is on something. **A safe behind the panelling** ($560, $1/day) puts up to $900 beyond a fine and beyond a thief. **A dry cellar** ($480, $2/day) hides 25 units of stock from the attention holding it draws — and a warrant served on the residence finds all of it, so it is a hiding place with a specific failure rather than a free one.
- They belong to the building, not to the person. Moving out leaves them, and the move action now names what would be left behind, which turns a housing tier from a strict upgrade into a decision about what you are giving up. They also survive their owner: the next arrival starts in the same rented room and the safe is still in the wall, with none of the money that was in it.
- Measured across 500 campaigns attacked at home: 112 survive with nothing on the door and 202 with steel behind it. Across 400 raids at 95 attention: $110 left in hand without a safe, $900 with one. Holding 20 units for 20 days: 20,000 police attention in the open across 200 campaigns, none of it in a cellar. With one contact and nobody hired, a telephone turns 0 warnings into 500.
- **A defect the balance test caught and the design deserved.** The first build had the door counted in `Guard()`, and `Guard() > 0` is the check that decides whether an attack at home arrives as a scene or as a killing — so a $420 door made the player survive 500 of 500 attacks. `Watchers()` now separates the people who could raise an alarm from the building that slows somebody down: a door cannot shout. It instead cuts the odds of being killed without warning by eighteen points, which takes survival from 22% to 40% and leaves it lethal.
- A cellar and a false floor are separate hiding places that stack, and a search of one does not find the other: with 45 units held behind a 25-unit cellar and a 20-unit floor, a warrant on the residence takes exactly 25.
- Verified over the HTTP stack against a fresh isolated save: a campaign fitted a door, a telephone and a safe to its rented room, watched the daily bill go from $15 to $23, its reach from 0 to 1 and $900 pass out of reach, and was told the doorway already had a door in it. A 200-command campaign exercised 35 command kinds including all four fittings with no invariant failures.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Somebody asks you for something
- Everything the player could do until now either resolved the moment they committed to it or paid out by the hour. Nothing in the game ever asked them for something and waited. A commission is standing work an organization wants done, with a three-day deadline, a reward and a cost for failing.
- It is deliberately not a script. The objective is a condition the world can check for itself — be carrying twenty-five crates at a named place, put a named property at forty-five percent or worse, take nine hundred dollars out of a named organization from wherever their money stood when the work was taken, see a named person dead, be worth a stated amount to the city. How the player gets there is their business, so every system built so far is a way of answering.
- What is asked for is derived entirely from the asking organization's situation. A family at war asks for a person, and only at war, because it is the heaviest thing anybody asks and it pays like it. A family whose rival has more money than they do asks for the money; one whose rival's strength is the ground they hold asks for the ground. A family poorer than its neighbours asks for stock it can sell. A quiet, solvent one asks the player to become somebody worth being seen with. Across 300 cities run for sixty days on their own, organizations asked for damage 118 times, a drain 121, standing 32, a delivery 16 and a killing 13.
- The drain records what the target's money stood at when the work was accepted, so what somebody else takes out of a rival's pocket is never credited to the player. A delivery hands the crates over, and at $1,150 for twenty-five it is worth taking in 314 of 400 campaigns and a loss in 75 — the market is part of the judgement.
- Three at a time, one per organization, and people whose goodwill with the player is below -20 do not hand them work. Nobody inherits an obligation.
- **A defect `go vet` caught before it could reach a save.** The struct grouped `Due, Life int` and `Pay, Respect, Goodwill, Penalty int` under single json tags, so five fields would have serialised under the wrong names and collided. Each field now carries its own tag, and a round-trip test writes a live commission out and reads it back field by field.
- Verified over the HTTP stack against a fresh isolated save: a campaign walked into the Mercer Exchange, was asked by Elena Russo to take nine hundred dollars out of the Bellandi Family, accepted it, watched the snapshot report "$0 of $900 taken out of Bellandi Family" with 4,290 minutes left, found Russo would not ask for a second thing while the first was outstanding, and three days later read "Nothing came of it" with the standing gone. A 220-command campaign exercised 41 command kinds with no invariant failures.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Somewhere else to be
- Bellwether is not the only place on the map, and the things it cannot give you are worth days of your life to go and get. Three destinations are bookable at the exchange, each with one reason to make the trip. Rockridge is two days and buys moonshine at fifty-five percent of the city price. Kingsport is three days and arranges the account abroad in person for $150 rather than the $400 a wire costs. Halloway is four days somewhere nobody has heard of you: attention falls twelve a day and you come back knowing somebody new.
- **This is not a second map.** A journey is one committed decision with a long price. The clock runs a day at a time while the player is out of the city, so businesses go unwatched, trouble finds the ones nobody is minding, commissions run down toward their deadlines, families split and rivals help themselves. The player is in transit for the duration, so nothing that asks where they are standing can be satisfied while they are two hundred miles from it.
- What you can bring back is what you have somewhere to put: twenty crates in your pockets, plus a car's false floor and a cellar at home. Measured across 300 trips each: carrying it in your pockets turns $150,000 of trips into $143,200 of stock, which is a loss, because you are jumped on the road 121 times in 300. With a false floor the same 300 trips turn $282,000 into $383,200, because the exposure that decides whether somebody picks you out of the street is now what is visible on you rather than everything you hold. A smuggling run does not pay until you have somewhere to hide the load, which is what a car was for all along.
- The one thing a journey buys that nothing else does is being elsewhere. With a hit arranged and nobody watching the door, 0 of 300 campaigns survive staying home and 300 of 300 survive leaving town — and the house is damaged in all 300, which is the price. Leaving the city is a real answer to a week you were not going to survive, and it costs exactly what four days away costs.
- Verified over the HTTP stack against a fresh isolated save: a campaign rested to full health, was quoted $680 all in for Rockridge, committed, and came back two days later having been jumped on the road, to a city where a family had split, Saint Agnes had run into trouble and somebody had helped themselves to the Russo garage.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## What people have not forgiven
- People in this city rob each other, get passed over for promotion and belong to organizations at war, and until now none of it left a mark: nobody remembered who had done what to them, and nobody ever settled it.
- A grudge is one person's memory of one specific thing another person did. They come from events the simulation actually committed — somebody helped themselves to a place you answer for, somebody took the job you thought was yours when the person above you died — they accumulate, and they fade by one a day, because people let things go slowly. Being in an organization at war is worth 35 on top, computed rather than stored, because it stops being true the day the war ends.
- At weight 55 the holder stops brooding and moves. Ambition decides who is willing to, skill against skill decides whether it works, and a failed attempt gets the wrong person killed three times in ten. Measured across 300 cities with the heaviest possible grudge: settled in 299 inside ten days, with somebody dead in 185 of them.
- **This is where private history becomes public war.** Every settled grudge between two organizations raises the hostility between them, in all 299 of 299 measured. One killing is not enough to start a war on its own — it takes a few — which is exactly right: the player reads about a killing in the Herald, then another, and then two families they had no quarrel with are shooting at each other for a reason nobody told them.
- A city left entirely alone for 120 days across 200 runs: 106 people dead, 158 cities where nobody died at all, and 16 live wars at the end. It is a risk rather than a certainty, which is the difference between a living city and a churning one. The most grievances any of 100 cities carried was 3 against a cap of 60, so the save stays bounded.
- The player hears about any of it only with two contacts or a telephone in the hall, and hears about it below the weight at which anybody acts — the useful thing is knowing two people have a problem before one of them settles it.
- **A cosmetic defect this surfaced.** A successor's role already carries their organization's name, so the Herald was reporting "They were Head of the Russo Outfit of Russo Outfit". Fixed in `describeStanding`.
- Verified over the HTTP stack against a fresh isolated save: thirty-nine days of a campaign doing nothing but building contacts and letting the clock run produced "Elena Russo has not forgiven Sofia Hale for what happened at Mercer Exchange" — a grievance from a robbery the player had no part in — alongside a splinter organization called the Brenner Company that had gone to war with Russo, and a Herald carrying OPEN WAR ON THE WATERFRONT and SPLIT IN RUSSO OUTFIT. A 200-command campaign exercised 33 command kinds with no invariant failures.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## People you come to know
- Everybody in this city had a name, a job, an ambition and a number for how good they were at things, and none of it made any two of them behave differently: a person with 80 ambition and a person with 80 skill did the same thing in the same situation, only more or less often.
- A temperament is the part of somebody that decides *how* they do a thing rather than whether. Hot-headed people move on a grievance long before it is wise to, and move badly. Careful people wait, and are rarely the ones who end up in the river. Grasping people take more than they came for. Loyal people will not move against their own, whatever they are owed. Vain people want to be seen doing it, which is how anybody ever finds out.
- It is derived from the person's name rather than stored, so it costs nothing in a save, it is the same person every time the campaign is loaded, and it cannot drift. The first hash — a weighted sum of the letters — put two of the five temperaments on more than a third of the city each; a proper mix replaced it, and a test now fails if any one kind of person is more than half the city.
- It changes outcomes rather than labels. With the same grievance, the same ambition and the same skill, across 400 cities and twenty days: hot-headed people acted in all 400 and killed their mark 157 times; careful people acted in 387 and killed theirs 196, so waiting is worth thirteen points of success rate. Grasping people took $35,157 across 300 robberies against $24,335 for careful ones. A vain thief is named by the player's contacts a full contact earlier than anybody else. A loyal person carrying the heaviest possible grudge against somebody in their own organization never acts on it at all, and acts on the same grudge against an outsider immediately.
- **A stranger is a stranger.** The player knows whoever runs an organization, anybody they have dealt with, and — with three contacts — most of the city. Everybody else is a face. What they know about somebody they do know is a dossier assembled from records the city actually committed: not a hidden stat block, just what was reported.
- Verified over the HTTP stack against a fresh isolated save: on day one the player knew four people; eighteen days later they knew eight, including "Ivo Duarte · Careful · Russo Outfit · known for: Taken from Pier 14, Somebody tried it at Russo Motor Works", and could be told that Vittorio Bellandi had not forgiven Ivo Duarte for what happened at The Monarch. A 200-command campaign exercised 30 command kinds with no invariant failures.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A charge under the floor
- Sabotage is a crew breaking things in the dark and everybody pretending afterwards that they do not know who. A charge is not that. It is unmistakable, it is in the paper, it kills people who happened to be standing there, and it is the difference between a message and a declaration.
- One is bought off a boat at Pier 14 for $850 — dearer than the best gun in the city — and holding it costs four police attention a day, which is more than anything else in the game costs to simply hold. Two is the most anybody carries. A search that finds one is a prosecution rather than a fine.
- Planting it needs 20 presence, because nobody who is nobody gets under a rival's floor, and the charge is spent whether it works or not. When it works the building loses 45 to 75 condition, its stock and a member of its staff, its still, and two thirds of anything behind its tables; the owning organization loses power and thirty dollars a point of damage; and somebody who worked there dies in 139 of 400 measured blasts. When it fails it goes off with the player under it — verified live, where it killed Alex Varga outright and put EXPLOSION AT THE MONARCH in the Herald with a description of a man leaving on foot.
- **It goes both ways, and neither direction needs the player.** An organization at war with money to spare will occasionally stop raiding a rival's premises and simply destroy them, and one whose standing with the player has fallen to -60 will do it to the player's best business. `detonate` is the same function whoever set the charge, so a bombing is a bombing whichever end of it you are on. Measured over 60 days: 214 of 300 cities saw an organization destroy premises outright, and with one organization at -80 standing and no war to fight, the player's laundry was wrecked in 157 of 300.
- The first rate was 6% a day per organization, which bombed 297 of 300 cities and made wreckage into weather. It is 1.2% now, and gated on the organization holding $1,200 spare.
- Verified over the HTTP stack against a fresh isolated save, both outcomes: one attempt killed the player, and the next campaign put The Monarch from 100 condition to 43 under the headline EXPLOSION DESTROYS THE MONARCH, with the Bellandi Family's standing at -50 and their strength down from 90 to 76.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A room nobody wanted to be in
- Two organizations that will not speak to each other will sometimes speak in front of somebody they both owe. That somebody can be the player, and it is the only thing in this game that can end a war without either side losing one.
- A sit-down is called on neutral ground — Saint Agnes, which nobody holds — for $220 that buys the room and the guarantees whether or not anybody agrees to anything. It needs 25 presence, because neither of them would cross the street for a nobody, and neither side will sit in a room arranged by somebody whose standing with them is below -25.
- Four things can be said in it, and every branch is decided from state the world already holds. Pressing them to settle takes 45 hostility off the quarrel, gives both sides 22 standing with the player and puts a truce in the Herald. Letting them talk cools it by six and costs nothing. Coming down on one side buys 18 standing there, costs 30 on the other and brings an operation against the player. Leaving costs eight with both.
- **Unless somebody in that room had already decided.** Whether a meeting is an ambush is worked out before anybody sits down, from what the two organizations have been doing to each other and from the temperament of the people leading them: hostility at 62 or above *and* a hot-headed or vain leader on one side. Both halves have to be true, so a hot-headed man in a cooling quarrel is still just a man in a bad mood. Pressing an ambush kills somebody on each side, injures the player by 35 to 75 before armour, ruins their suit, and takes their crewman away one time in three. Listening to one is safe.
- The player only knows in advance if somebody warned them. With two contacts or a telephone in the hall, every ambush is flagged, and the scene itself reads differently: *one of them brought more men than the room needs, and Mara caught your eye on the way in*. Without that, the first they know is when somebody stands up.
- Measured over 300 cities left to run 90 days on their own: 244 had a quarrel worth calling a room over, and 29 of those rooms had somebody in them who had already decided — about one in eight. A stress fixture forcing hostility between 55 and 95 pressed 200 rooms: the quarrel changed state in 75 and the player was hurt in 167, which is what the top of the range costs.
- Verified over the HTTP stack against a fresh isolated save: a campaign called the Bellandi Family and the Russo Outfit to a room, pressed them, and settled it — TRUCE BETWEEN TWO FAMILIES in the Herald, the conflict off the public list entirely, and the player out of the evening at 30 health.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## The week after
- A coup already existed: somebody below the leader of a failing organization took their place by force, the leader died, and the organization lost twelve strength. That was the whole of it. Nobody chose a side, nobody was counted afterwards, and the family that had just killed its own head went back to business the following morning.
- Who moves is now decided the way everything else in this city is. A loyal person will not move on their own leader whatever they are owed. A hot-headed one moves early. And a grievance somebody has actually been carrying is worth more than any amount of ambition — a challenger with 40 ambition and a real reason is more willing than one with 100 and none.
- **Ambition is now a source of grievance in its own right**, which is what made the whole system reachable. An ambitious lieutenant resents whoever is above them by two a day against one a day of forgetting, so it takes about a month and a half to become a reason. This is the only grievance in the city that does not need somebody to have died first, and without it the measured coup rate over 200 cities and 120 days was zero.
- The aftermath is the part everybody remembers. Everybody left has a view on the man who did it — in 400 of 400 measured aftermaths. Somebody who backed the wrong side is dealt with rather than forgiven in 175 of 400. And in 139 of 400, somebody walks out rather than answer to him and takes premises with them, which is the splinter machinery doing what it was built for. The Herald runs UPHEAVAL IN whichever family it was.
- Measured over 200 cities left alone for 120 days: 46 saw an organization turn on its own leadership, 5 ended with more organizations than they started with, and none ended with fewer than two holding anything. Upheaval is an event rather than the weather.
- Verified over the HTTP stack against a fresh isolated save, sixty days of a campaign doing nothing but building contacts and letting the clock run: the Russo Outfit gone entirely, the Doyle Company and the Serra Company in the city that had never had them, Luca Lenz dead, and a Herald carrying ONE DEAD IN EXPLOSION AT PIER 14, ATTEMPT ON THE LIFE OF VITTORIO BELLANDI, SPLIT IN BELLANDI FAMILY and OPEN WAR ON THE WATERFRONT.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## People with salaries
- A detective who would lose a file for money already existed. What did not is the rest of the building. Commissioner Vance and Mayor Ellis Crane are now people in the city like anybody else — they have a voice, a temperament and a location, they can be resented and they can be killed — and each can be put on a standing arrangement rather than paid a one-off favour.
- The commissioner costs $900 to open and $45 a day, and raids begin twenty attention later than they otherwise would while nothing is ever forfeited. The mayor costs $1,400 and $60 a day, and every business earns a fifth more. Both cut the player loose the moment their attention passes a stated ceiling — 78 for the commissioner, 65 for the mayor — and keep the opening payment. Nobody in that building goes down with anybody.
- **They are corrupt to whoever is paying most.** An organization richer than the player by $4,000 *and* whose standing with them is below -20 quietly takes the arrangement's benefit away while the player keeps paying for it. Money alone is not enough, because every family in this city has more of it than a man starting out — it takes money and a reason. Keeping the families sweet is now a way of protecting something the player bought from somebody else entirely.
- Measured over 60 days: one laundry nets $-9,338 without the mayor and $-10,842 with him, so he does not pay for himself on the strength of one business. Three businesses net $14,998 and $17,637, so he does. Held deliberately at 60 attention for 30 days across 300 campaigns, a player was raided 289 times without a commissioner and none at all with one.
- Killing a man with a title is the loudest thing that can happen in this city, and the same rule covers the player doing it and anybody else: 35 attention on the player, eight strength and $1,200 off every organization in the city whether they had anything to do with it or not, the arrangement gone, and CITY REELS AS COMMISSIONER VANCE IS KILLED in the Herald.
- **A defect this surfaced immediately.** Adding two officials as ordinary people put them into `PeopleDay`, and a mayor started going out at night taking tills — a balance test caught it as "something happens in every city every time, which is not a city, it is a treadmill". Officials are people who can be resented and killed, and are not people who rob businesses.
- **A second defect, in a test rather than the code, worth recording because the finding is real.** The licence test first read as the mayor *reducing* income. He does not: the retainer had made that run insolvent first, and an unpaid daily bill costs a business more than a licence is worth. With both runs solvent and income measured per elapsed minute rather than per campaign — an encounter pauses the clock and the two runs do not stop in the same place — a laundry earns exactly the advertised twenty percent more.
- Verified over the HTTP stack against a fresh isolated save: a campaign opened both arrangements at the exchange, watched its daily bill go from $15 to $120, found both listed as not outbid, and ended one to take the bill back to $60.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.
- Save version 9 puts them into a campaign that began before there was anybody in that building. The first restart of the live game reported an empty cast, because migration only runs when the stored version is behind: adding people to the world is not enough on its own.

## Whose hands were on it
- Every violent thing the player could do, they had to do standing there. A man with a crew does not do that, and the reason he has a crew is so he does not have to. Robbery and sabotage now each come in two forms, and the same rules resolve both.
- **Your own hands** bring your standing, your gun and your armour to it. `HandEdge` is presence over 400 plus the weapon edge, which is what the odds were built on all along. All of the risk is yours: the beating, the ruined suit, the damaged car, and the death.
- **Somebody sent** brings their loyalty to it and nothing of your name — loyalty over 500, less eight points, which is meaningfully worse. In exchange the player takes no physical risk at all, draws six less police attention because somebody else was the man described, and earns a fifth of the standing, because a man who sends people is respected less than a man who goes. That last part is most of why anybody goes.
- The risk does not disappear, it moves. A job that goes wrong costs the man who was sent 25 loyalty, and a job that goes badly wrong kills him: 81 times in 400 measured, and never for a setback below the threshold. He then goes through `Kill` like anybody else, so his death reaches the Herald and the city's records rather than quietly emptying a slice.
- Measured across 500 robberies at the same casino: going yourself takes $86,522, earns 446 standing, draws 6,385 attention and hurts you 277 times. Sending Leo takes $61,272, earns 157 standing, draws 3,715 attention, hurts you never, and costs you Leo 24 times.
- Nobody goes out on this kind of errand below 40 loyalty or while already on an assignment, and both forms refuse for the reasons that apply to both — a player who can barely stand is offered neither.
- Verified over the HTTP stack against a fresh isolated save: a campaign recruited Leo, was offered both forms at The Monarch with the trade stated in each, sent him, and came away $271 up, one standing, four attention and entirely unhurt.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## The room under the floor
- Arms could be bought and never sold. The player could carry a gun and nobody else in the city could buy one, so a war between two organizations was fought with whatever they already had and cost them nothing but people.
- Crated arms are now a good like any other, at $165 a crate, three attention a unit — the heaviest thing in the game to hold — and they come off a boat at Pier 14 and are not sold across the exchange's public floor.
- An armoury is a room built under a laundry or a garage the player owns for $1,200. It holds sixty crates, and the customer is a war: any organization currently fighting somebody buys one to three crates a day at 235% of the waterfront price, gets two strength a crate for it, and thinks better of the player for selling. Only one room, because everything in one place is bad enough.
- **The room is as loud as what is in it.** Three attention a day plus one for every twelve crates, so a full room draws eight and the only way to make it quieter is to sell it down. The first build used a flat three, which meant a player could earn $19,000 in a fortnight and never pass thirty attention — the balance test caught it as money with no danger attached.
- Arming one side is something the other side finds out about, 210 times in 300 measured days of selling into a war, costing eight standing with them. And a warrant that finds the room takes the crates, the room, thirty-five condition off the premises and twenty more attention, with ARMS CACHE SEIZED in the Herald.
- Measured over 200 campaigns of 45 days, starting with a full room and never once bribing anybody or lying low: selling into a war earns $21,360 against $6,987 for running the same laundry quietly — three times the money — while all 200 passed the attention the police act on and 177 lost the room to a warrant. It is the most profitable thing in this game and the most difficult to explain to anybody.
- Verified over the HTTP stack against a fresh isolated save: a campaign acquired the Bluebird, built the room, bought twenty crates at $143 on the waterfront, stowed them, and over the following days sold six of them to the two organizations quarrelling in that city at $296 and rising, for $2,237.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.
- Save version 10, because a good added after a campaign began was missing from it entirely and a market the player cannot see is a market that does not exist. Migration now adds any good the save is missing at its reference price, so this cannot recur for the next one.

## Taking it off somebody
- The city could rob premises and the player could rob premises, and nobody in it could rob a person. That is the crudest thing anybody does here and it was the one thing missing.
- The target has to be somebody the player knows and who is standing where they are, and never their own crewman. What a person is carrying comes from what they are: thirty dollars plus twice their standing, plus a share of their organization's money scaled by how far up it they are. A man with a title carries three times that, and a vain man half again, because it is on him and it shows.
- **Standing makes you better at this and worse at getting away with it.** Presence improves the odds like it improves everything else, and above 35 of it the man you robbed can name you: his organization's standing with you falls by 25 instead of 12 and they send somebody. This is the only place in the game where being known costs rather than pays, and sending your crewman instead puts his face there rather than yours — the grievance is his to carry.
- Measured across 400 attempts each, at the same casino with the same run of luck: a soldier gives way 160 times for $130 apiece, a boss 113 times for $430. Going after the soldier hurts you 240 times, the boss 287. Who you pick is the whole decision.
- Robbing a man with a title is not robbing a man: twenty more attention, and whatever arrangement you had with that building is over.
- **A defect worth recording.** A grudge needs a person to hold it against, and the player is not a person in the city's records — so the first build blamed the crewman for something the player had done themselves, and the live check reported *Vittorio Bellandi has not forgiven Leo Carver* for a robbery Leo was not at. When the player does it themselves the answer is their organization's standing, and nothing else.
- Verified over the HTTP stack against a fresh isolated save: a campaign walked up to Vittorio Bellandi at The Monarch, was told he had about $461 on him, took it, and read *He knows your face* with the Bellandi Family's standing at -25.
- Money behind the panelling at home is also no longer in a pocket in the street: the safe now protects against being robbed as well as against a fine.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## How it was done
- A killing was a sentence the code wrote: killed at a place, by an organization, over something. It said what happened and never how, so twelve different deaths in a campaign read as the same death twelve times. There was a table of methods, and only commissioned hits ever used it.
- Every death in the city now has a manner, composed from state the simulation committed before anybody died: where the person was standing, what hour it was, and the kind of person who came for them. Ten locations with four ways of dying each, seven times of day, and a signature for each temperament — a hot-headed killer leaves *it was not quiet and it was not quick*, a careful one *nobody heard it and nobody has said anything since*, a grasping one *their pockets were empty when they were found*.
- Every call site that knew who did it now passes them: the grudge that was settled, the coup and its purge, the man sent to do something for somebody else, the sit-down that was never a meeting, the war that reached people. A commissioned hit adds what the fee bought — a cheap one leaves a great deal behind, a specialist leaves nothing and was never going to.
- Measured over 300 cities left alone for 120 days: 332 deaths in 280 distinct descriptions, with the commonest used six times.
- **Two defects the live check caught, both in the writing rather than the rules.** The place was appended after the method, which turned a clause into nonsense — *walked twenty feet before anybody caught him at The Monarch*. The place and hour are stated first now, so an entry is free to end on a clause. And that same entry called Sofia Hale *him*: nobody in these descriptions is a he or a she any more, because the city does not know and does not need to.
- Verified over the HTTP stack against a fresh isolated save, sixty-four days of a campaign doing nothing but letting the clock run: *At Mercer Exchange, in the small hours: Ivo Duarte was beaten in a storeroom with the door shut on it*, and *At The Monarch, in the small hours: Vittorio Bellandi was knifed in the crowd, and got twenty feet before anybody noticed. Whoever did it had done it before.*
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A city with people in it
- A city with eight people in it is a cast, not a city. Every system built so far — grudges, coups, muggings, wars that kill people — was operating on a population small enough that one bad month emptied it.
- **The enabling change was names.** Fourteen first names and twelve surnames gave 168 possible people, which is a hard ceiling, and `newPersonName` gives up after sixty attempts so it started failing long before that. Sixty first names and forty surnames give 2,400.
- A new city now holds 48 people: nine to each organization with two lieutenants apiece, twenty-six who answer to nobody, the officials, and the fixed cast. The street has trades — barmen, dockers, croupiers, a photographer, a newspaperman — because somewhere to be and something to lose is all it takes to be worth robbing, killing or recruiting.
- Organizations that lose people take more off the street rather than inventing them, so the two populations move between each other: *Franca Sabbatini signs on with the Bellandi Family. They were a stallholder a week ago.* They stop at nine and never grow past it on their own.
- **Filling the streets had to make the city deeper rather than four times busier.** Each person's daily chance now scales by how many people there are, so a city of eight and a city of forty-eight produce about the same number of incidents a week — 1.8 across 300 measured cities — and what changes is who they happen to. The old test asked whether any city was ever silent for twenty days; in a city of eight that was a fair question, and in a city of forty-eight a three-week silence would mean a dead one, so it measures the rate now.
- **Two more things this broke, both of which were the old scale showing.** Voices are shared, because there are more people than voices — the guarantee that matters is that a voice belongs to a person for life and that no one voice covers a third of the room. And a network of five contacts made the player know all sixty people in the city, which is not knowing anybody: contacts are worth the people your contacts would talk about, so a laundress stays a stranger until you deal with her.
- The save stays bounded: 200 people is the ceiling, the dead are forgotten once no contract, grudge or commission refers to them, and the worst any of 60 cities reached after 200 days was 111.
- Verified over the HTTP stack against a fresh isolated save: 48 people on the first morning with four of them known, and by day 39 sixty people, thirty of them known, twenty-six in organizations and thirty on the street, with the Herald carrying two killings and a war.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass. Save version 11 fills out a city that was eight people.

## On the way there
- Travel was minutes passing. The city was settling grudges, raiding premises and robbing people during those minutes, and the player never saw any of it because none of it happened where they were.
- A completed journey now sometimes puts them at the scene. Nothing here invents an incident: it shows one the city already had a reason for, in the order of how hard it would be to miss. A war in progress is the loudest thing on any street. A grievance heavy enough to be about to be settled is two men in a doorway, one of them doing all the talking. Police attention past half the raid threshold is a car that was not there when you set out and is still there when you look back. A quiet city has quiet streets, verified across 300 journeys through one.
- **Whether the player understands what they are looking at depends on whether anybody would tell them.** With two contacts a shooting is *the Bellandi Family and the Russo Outfit, and no doubt about which was which*; without them it is *two cars and a lot of shouting, and somebody face down in the road when the street filled in again. Nobody you know can tell you what it was about.* A man shot in the street for a reason nobody explains is the point rather than a failure.
- A war can reach somebody who was only passing: three times in ten of the shootings crossed, something comes off a wall and finds them for 18 before armour. Two men arguing in a doorway never shoot a passer-by.
- Measured across 7,200 journeys through 120 cities at open war — the worst case rather than the usual one: 16% of journeys were an incident, 338 of them injured the player, and none killed one outright.
- Verified over the HTTP stack against a fresh isolated save: a campaign with five contacts crossed three shootings between Saint Agnes and Pier 14 while the Bellandi Family and the Russo Outfit were at war, was told exactly who was shooting at whom, and lost 18 health to one of them.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## What the director was never told
- Twenty systems had been built since the director last got new vocabulary. It knew about wars, holdings and the player's businesses, and nothing about arms, commissions, arrangements with city hall, grudges, journeys or the fact that the city now holds fifty people.
- Four situations now exist that could not have before, each with the committed fact that justifies it. **Consignment**: an organization that is fighting wants what is under the player's floor and does not want to be seen at a laundry in daylight. **Grievance**: two people the player knows have a quarrel far enough along that somebody is going to get hurt. **Obligation**: work the player promised somebody, with the deadline running out and the progress stated. **Warning off**: an organization that has stopped complaining about the player, which is worse than complaining. None of them can be offered by a player with nothing going on, verified.
- The director is also handed what the player owes and to whom, what is under their floor, and the shape of the city's population.
- **The defect this iteration existed to find.** The narrative brief — the premise, the roles, the player's task and the constraints for every operation — was only ever attached to the *focused* prompt, and the default is the full one. In the live game nothing had ever told the model what an operation is. Asked for a consignment, it wrote a story about collecting a debt for a delivery of fabric and labelled it *Hand the crates over*. The brief now goes on both prompts, and a coverage test fails if any offerable operation lacks a premise, a task or constraints — which immediately found that `courier` had no constraints at all, so nothing had stopped a courier scene from deciding what was in the envelope.
- Verified live against the local model, on a fresh isolated save with a laundry, an armoury holding twenty crates and two organizations fighting: *"Take those crates of arms out of Bluebird Laundry and get them to the docks. We don't want to be seen there during the day, and the Bellandi Family doesn't like to make a scene. I'll have someone waiting at Pier 14 to take them from you. Just make sure the crates don't get opened and no one sees you carrying them."* Before the fix, the same request produced a story about a merchant's debt.
- One rejection along the way was the attribution guard doing its job: the model gave the work to "Bellandi men" while declaring it neutral. The consignment brief now states that the organization named in the fact is the beneficiary, and the warning-off brief that nobody speaks for the other side.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A hand at a time
- Sitting down at the tables was one roll and a paragraph. The house kept its edge, the money moved, and the player made no decision after the one to sit down — which is the only decision that does not belong to a gambler.
- A hand is now several decisions with the arithmetic anybody at a table faces: what you are showing, what the dealer is showing, and whether one more card is worth it. Draw or stand, an ace comes down from eleven to one rather than killing you, the dealer draws to sixteen and stands on seventeen, a tie gives your money back, and a two-card twenty-one pays half again.
- **The house's edge comes from the rules of the room rather than from a weighted table**, and it lands where it has to: 6.3% of everything staked across 20,000 hands played to seventeen, against the 6% the casino's own books have always modelled. The game the player sits down to and the game an owner's float simulates are now the same game. Most of that edge is one rule — a hand that goes over is dead before the dealer plays at all, so the player busts first and loses even when the dealer would have gone over too.
- The stake is taken when the cards are dealt, because that is when it is on the table. A hand that has not been settled lives on the world, so a save written mid-hand comes back to the same hand. Nothing settles twice, no second hand can be dealt over the first, and the money moves through the same books a session always moved it through — including the part where winning heavily in somebody else's room is noticed by the somebody.
- Verified over the HTTP stack against a fresh isolated save: three hands at The Monarch — drew from 14 and went over, drew from 15 to 20 and was paid, and stood on 18 against the dealer's 18 for a stand-off that gave the money back.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A picture under every headline
- A headline with nothing beside it is a line of text. The paper the inbox asked for has a coarse black and white block under every story.
- **Nothing here generates an image, and nothing here needs to.** No model weights are installed on this machine, and downloading some would be a decision the player has not asked for. What the core owes the presentation is the truth of what the picture should be *of*, and that can be read out of the story the city already filed rather than kept as a second record that could disagree with the first: headlines name premises and people in capitals. Premises win when a headline names both, because ROBBERY AT THE MONARCH is a picture of The Monarch, and the dead still get theirs.
- Across 60 cities left to run 120 days, 1,464 stories: 1,238 pictures of premises, 103 of people, 123 of the city when it was neither. Every single story had something to draw.
- The plate itself is drawn as SVG in the presentation layer: a blocky silhouette screened with halftone dots the way a 1930s press would have printed it, deterministic from the headline so the same story always carries the same picture, with the story's kind adding to it — a police car under a raid, a struck-through figure under a killing, a broken line over damaged premises. Three frontend tests cover it: every subject and kind produces well-formed markup with balanced groups, the same story is always the same plate and a different one is not, and a name with a quote in it cannot break out of the label.
- Verified over the HTTP stack against a fresh isolated save, twenty-six days of a city running on its own: ROBBERY AT THE BLUE HOUR pictured as The Blue Hour, FAUSTO TOTH TAKES OVER BLUEBIRD LAUNDRY as the Bluebird, and OPEN WAR ON THE WATERFRONT as the city.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, `npm test` and the production build pass.

## Nobody holds a job for ever
- Mara is the fixer at Saint Agnes, Leo drives, Harlow is the detective who keeps turning up. The game referred to all three by name — `Speaker: "mara"` — which made all three immortal by accident: killing Mara would have left every scene that speaks through her without a speaker, and Harlow only existed at all once a police stop had happened.
- A role is a job in this city rather than a person. Everything now reaches these three through the job: the fixer's warnings, the contract conversations, the sit-down, the attack at the residence, the detective's stop and his bribe. When whoever holds one is gone, somebody else is doing it by the end of the day — a different name, a different voice, and no memory of anything the player did for the last one. The first holder is still the one the campaign was written with: an earlier build promoted a stranger on the first morning and lost Detective Harlow entirely.
- **The city can now rob its own people, which is what makes these three mortal rather than merely killable.** Nothing in the simulation had ever reached somebody who is in no organization and runs no premises. An ambitious person may now take what somebody standing next to them is carrying, by the same arithmetic the player faces, leaving the victim with a grievance — and a grievance is how anybody in this city ends up dead.
- Measured across 120 cities: a charge under Saint Agnes passed 240 jobs to somebody new and left none unfilled, while 120 stayed with the person who started with them. A successor arrives at zero trust and is a stranger until the player deals with them.
- **Two idempotence defects this surfaced.** Filling a job takes somebody off the street, and the street was being topped up first — so the city grew by one every time it was counted, and running the migration twice added a person. Jobs are filled before the street is counted now. And somebody doing one of these jobs is no longer on the street at all, so no family recruits the detective and nothing replaces the fixer with the fixer.
- Verified over the HTTP stack: a 200-command campaign ended with Mara and Harlow still in place and *Ugo Lenz* driving, because Leo had not survived it.
- Save version 12. Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, `npm test`, the production build and `cmd/simulate` pass.

## The outside of a business
- A casino had a float, a nightly handle and somebody with real money at the table. A laundry had a number that went up. Both were called businesses.
- Custom is the legitimate trade a place has built up: the people who bring their washing there rather than somewhere else, and the cars that come back. It starts at fifty, rises two a day while the place works at eighty percent of its potential with nothing wrong, falls three a day while something is, falls two a day while it is being skimmed hard, and multiplies everything the place earns — half at nothing, half again at the ceiling.
- **The books and the shop pull against each other, which is the decision the whole thing exists for.** A round through the books is worth twelve trade, because the machines are always busy and the regulars go elsewhere, and a police search in daylight is worth twenty. Measured over 120 campaigns of 60 days: a laundry whose books are never used holds 30% trade and earns $9,341; one used every four days falls to 1% and earns $4,959. The cheapest way to clear attention is no longer free.
- The reward for the other way of playing it is a standing order: somebody respectable wants the same work every week, for $26 a day, offered only to a place holding 65% trade and kept only while it works at eighty percent. Losing one costs fifteen trade on top of the money. Over 90 days a laundry left alone earns $11,787; one kept stocked, staffed and out of trouble with an order on the books earns $47,755.
- **A zero-value defect worth recording.** A stored zero meant both "a save written before any of this" and "no trade at all", so a shop that decayed to nothing read as a fresh fifty and started decaying again. The floor is one rather than nothing: a shop nobody goes into any more still has a door on it.
- Verified over the HTTP stack: trade reaches the API on every trading business, started at 50%, rose to 52% over three clean days and then fell to 37% across a fortnight in which a family sabotaged the premises to 5% condition and left trouble standing in it — which is the mechanic working rather than failing.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## When the city starts calling you something

*The first system in this game the player did not ask for. The inbox they wrote is finished; this is the answer to a question it left open — the city could do everything to an organization and nothing to the player, because the player was not one.*

- Declaring war on it, raiding its holdings, weighing its strength against a neighbour's, deciding it is weak enough to move on: none of it could be done to the player, who stood outside the machinery that runs everybody else. Two premises and 25 respect is the point at which that stops.
- **Nothing here invents an identity.** The player's ownership string has always been `player:1`, so an organization filed under that id already holds their premises: `FamilyHoldings` needed no change at all. Their strength is what they actually have — eight a holding plus its condition, up to thirty for their name, three a guard, an eighth of each crewman's loyalty — recomputed daily and capped at a hundred.
- Everybody already in the city gets a relationship with the new organization, starting where they already stood: forty, less half of whatever they thought of the player. Somebody at -60 goodwill starts a feud on the first morning.
- **A raid on the player is a raid.** `contest` takes their premises' condition, six of the shop's trade, and the money out of their own pocket rather than out of a number that is rewritten every morning. It reaches their crew, because they have no soldiers of their own to lose, and a crewman it reaches dies through `Kill` like anybody else. A holding it cannot defend changes hands to whoever took it.
- Measured over 150 campaigns of 90 days, with the worst case forced — every family at war with them from the first day: a proprietor loses no money and no premises; an organization loses money in all 150, premises in all 150, and 293 of 300 holdings in total. Being somebody is worse than being nobody, and ending a war is now something worth paying for.
- **Two things this could have quietly broken.** `FamilyDay` repairs every organization's holdings eight condition a morning, which would have repaired the player's businesses for free and undercut the repair action; the player's entry is skipped. And `dissolve` ends an organization that holds nothing and has no strength left, which would have deleted the player's the first bad month; theirs ends when they do, and takes its quarrels with it.
- Verified over the HTTP stack: a campaign took the Bluebird and the Blue Hour, and by the next morning was *Alex Varga's people* at 72 strength, in a feud with both established families and listed alongside a splinter that had formed the same week.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## People who answer to you
- Being filed with the others gave the player holdings, a name, a strength and quarrels, and one man. Every other organization in this city has nine people with standing, ambition and an opinion of whoever is above them.
- Somebody standing in front of you can be asked: $140 up front and $14 a day. They become a person in the city who answers to you, which is a real thing rather than a slot — `Members` finds them, so a raid on your premises reaches them and kills them like anybody else, and what they think of you adds to what your organization is worth in a fight. Six is as many as answer to one name at this level.
- **What they think of you is the whole system.** It drifts up a point a day while the bill is met and falls six a day when it is not. Below twenty they start looking for somewhere else to be, and what they do about it depends on who they are: most simply are not there one morning, but somebody ambitious and not the loyal kind walks out with your best business and every arrangement you have. A loyal man never does, however far down he goes — verified across 400 days of neglect.
- Measured over 200 campaigns of 60 days at war with everybody: alone you keep 152 of 400 holdings; three people kept and paid keep 326; three people you stop paying keep 157 and 599 of 600 of them walk out. People are worth more than double what you hold, and worth nothing at all if the bill is not met.
- Nobody signs on with a man — only with something that has a name — and nobody signs on who answers to somebody else, holds one of the city's jobs, or has a title.
- **A defect the live check caught rather than any test.** Signing somebody on required knowing them, and nothing in the game makes a civilian known — so across four locations, forty jobs and two acquisitions, the action never once appeared. Knowing somebody is not required to offer them work: walking up to a man in a bar and asking is the whole of it, and signing him on is what makes him somebody you know.
- Verified over the HTTP stack: a campaign became *Alex Varga's people* at 72 strength, put Ugo Lenz on at Saint Agnes, and watched the daily bill go from $88 to $102.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Taking ground
- The city could take the player's ground and the player could take nobody's. Sabotage damages a holding and a charge wrecks one, and neither has ever moved a name on a deed — so a war with the player was something they could survive or end, never win.
- A move is the same raid the city runs, with the player's organization as the attacker. `contest` was split so the target can be named: the city always goes for the weakest thing somebody holds, and the player picks. That is the only difference between what they can do and what is done to them. Everything else — strength against strength, the damage, the casualties, a holding that cannot be defended changing hands — is one function serving both sides.
- It needs an organization rather than a man, somebody you are past being civil with, thirty strength, being there, and being able to stand. It draws fourteen attention, costs the holder twenty-five standing and hardens the quarrel by twelve. Taking ground is worth ten respect.
- **Driven off, it costs the people who went.** Somebody who answers to you dies first; then your crewman; and only if there is nobody else does it come out of the player's own health. Across 400 moves on a weakened holding: 204 took it, 196 were driven off, and the player was hurt in 196.
- **The honest number.** Over 200 campaigns of 60 days at war with the strongest family in the city from the first morning, holding two shops: sitting still ends with more ground in none of them and less in 117; fighting ends with more in 6, less in 161, and hurts you in all 200. A war can now be won. This is not the war to pick.
- Verified over the HTTP stack: a campaign became *Alex Varga's people*, moved on Saint Agnes against the Russo Outfit, was driven off — *a raid repelled at Saint Agnes*, *you went in with them* — and came out at 73 health with fourteen more attention, having crossed a shooting on the way there.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## What you can find out
- The player could read a rival's exact strength, their exact standing and the exact money in their books, off a screen, having never spoken to anybody inside them. Nothing in the fiction accounted for it: a man who has never met anybody in an organization does not know what is in its accounts. The snapshot was handing over `w.Factions` whole.
- What you can find out now depends on who you know, in four steps. A stranger gets a name and nothing else — not even who runs it. A network of two contacts gets the leader and a description: *as much as anybody*, *a good deal*, *something*, *not much*. Somebody inside it you have actually dealt with gets their strength as a number. And asking around — $90 in the right pockets, an afternoon, current for about a week — gets the books.
- What they think of the player is always known, because the player can always tell how they are being treated, and their own organization needs no enquiry.
- **This is the first thing in the game that rewards dealing with people rather than spending money**, and the first that can be out of date: what you were told a week ago is what an organization looked like a week ago, and the enquiry lapses rather than persisting.
- Measured across 200 played campaigns of 60 days: of 433 organizations, 190 were known not at all, 223 by reputation, 20 through somebody inside, and none completely without somebody having asked. Most of the city is something you have not asked about.
- Verified over the HTTP stack: a campaign on its first morning was told two names and nothing else; with five contacts it learned who ran them and that the Bellandi Family had *as much as anybody* and the Russo Outfit *a good deal*; and $90 later the Bellandi Family were *90 of a hundred*.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Standing with somebody
- The player could change what an organization was worth in exactly one way: by hurting it. There was no version of this city's politics in which two people who both have a problem with a third stand together, which is most of what politics is.
- An understanding costs $700 to open and $30 a day to hold, and needs four things: being an organization rather than a man, +15 standing with them, a common enemy — somebody you are both past being civil with — and knowing enough about them to be sure what you are agreeing to, which means the intelligence built last iteration is now a prerequisite for the diplomacy built this one.
- While it holds, neither of you moves on the other in either direction, and somebody who stands with you may already be there when a raid comes: 175 of 400 measured, scaled by their strength, which drives the attacker off and costs them three strength.
- **It is protection bought with money and paid for in enemies.** Signing puts fifteen hostility on you with whoever they are fighting, and every day standing with them puts two more on everybody they are at odds with. Miss the tribute and it lapses — nobody stands with a man who cannot pay for it. End it yourself and they take twenty standing off you, because they remember which side did that.
- Measured over 200 campaigns of 60 days at war with a ninety-strength family: alone you keep 266 holdings and spend $274,170; with an understanding you keep 399 and spend $386,540. Half again the ground for two fifths more money.
- **What the measurement did not support.** The enemy's hostility ends *lower* with an understanding than without — 4,414 against 6,823 — because a war burns itself out faster than two a day accumulates, and an ally at the door takes strength off the attacker. What making an enemy costs is immediate rather than cumulative, and the test says so rather than claiming otherwise.
- Verified over the HTTP stack: a campaign became an organization and was told exactly why it could not yet reach an understanding with either family — *they think of you at +12. It takes +15 before anybody would consider it* — which is a standing it can go and earn.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## What you leave behind
- A protagonist died and their premises passed to their estate, which was a string in an ownership field and nothing else. **Their people kept a faction id that no longer resolved to anything** — orphaned, unfindable, still counted as somebody's men by a lookup returning nil. That was a bug, and the fix is the feature.
- What actually happens when a man who ran something is killed is that somebody who worked for him is running it by the end of the week. The strongest of the player's people now takes it over: a new organization under a name the player never chose, holding their premises, keeping their quarrels — the conflicts are re-pointed rather than dropped, so whoever was at war with the man is at war with what he left. It is worth thirty strength less than it was, because most of what it was worth was him.
- Somebody who had nobody leaves nothing, and their people are released rather than left answering to a dead man. Either way nothing in the save points at an id that resolves to nothing.
- The next protagonist arrives into a city containing their predecessor's organization and can do anything to it they could do to any other: ask around about it, deal with it, stand with it, or take it back. Measured over 150 inherited organizations left to run 90 days: 134 were still standing at the end and 16 had been destroyed entirely, which is the same attrition any organization faces.
- Verified over the HTTP stack, and the run tells the story better than the summary does. Alex Varga became an organization, put Ugo Lenz on, went after the Bellandi Family at their own casino until the police took his businesses at eighty attention, and was killed on the street by a war he was not part of. Nico Ward arrived with $90 into a city that now contains the Bellandi Family, the Russo Outfit, the Vance Crew — and *Ugo Lenz's people*.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.
- A 200-command campaign against a fresh save exercised 22 command kinds including three deaths and three new lives, with no invariant failures.

## The other career
- Every path this game offered started with a man buying premises. There was no version of it in which somebody arrives with ninety dollars, goes to work for one of the families, and comes up through it — which is the oldest story this city has.
- Answering to somebody takes +20 standing with them, 8 respect, and the conversation happening where they actually are. It pays $22 a day rather than costing one, out of their money rather than out of nowhere: a family that cannot meet it does not pay, and nobody stays where the money stops. Their quarrels become yours the same afternoon — everybody they are fighting takes fifteen standing off you.
- **Coming up is the point.** Work done for whoever you answer to — an arrangement completed for them, a commission settled for them — is standing inside them. Three pieces makes you a soldier, six a lieutenant, and a lieutenant takes four percent of what their organization's holdings bring in each day: $22 at the bottom against $71 for the same family, and less when they lose ground, because it is a share of what they actually hold.
- Nobody comes up through somebody else's organization while they have one of their own, and nobody in service starts one — the two careers exclude each other, which is what makes it a choice rather than an accumulation. Walking out costs forty-five standing and they answer it.
- Measured over 150 campaigns of 45 days from the same start: a wage earns $1,980 and a laundry earns $8,455. It is a slower living, and it is the one that needs no capital, draws no attention, forfeits nothing to the police and comes with somebody standing behind you.
- Verified over the HTTP stack: a campaign paid tribute at The Monarch until the Bellandi Family were at +48, went to work for them, and was reported back as an *Associate of Bellandi Family, $22 a day, 6 more jobs to come up*.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## The chair
- Coming up stopped at lieutenant, which made the second career a wage with a ceiling rather than a story. This city has exactly one thing above lieutenant.
- Everything needed already existed: the city runs coups on itself every season, weighing a challenger's poise against a leader's. This is the same move with the player as the challenger, and what it wins is not a promotion — it is the organization. Its premises become theirs, its quarrels are re-pointed at them, and its people choose: six in ten stay for the man who did it and the rest are out of work by the morning.
- It needs to be a lieutenant, in the room with him, worth 40 presence and able to stand, and it needs the organization to have fallen to three fifths of what it was — nobody moves on the head of something that is winning. **Which means the arc is: go to work for them, come up by hurting their enemies, and then quietly undermine them until the moment arrives.** The refusal says so in as many words: *Bellandi Family is doing too well for anybody to move on the man running it.*
- The odds are what the player brings against what he has: their presence, every piece of work they have done, whatever they are carrying, and a loyal crewman, against his own competence and half of what the organization can put behind him. A hundred respect and a Thompson takes it 139 times in 400 against 78 at the bare minimum.
- Measured over 500 attempts made at the bottom of what the rules allow — 50 to 89 health, which is where somebody who has just come up usually is: 105 took it, 255 were thrown out and no longer one of theirs, and 140 did not leave the room.
- **A gap the live check found.** Work only counted from an arrangement with a beneficiary or a settled commission, and most players are never offered one. Harm done to somebody's enemies is work done for them, so sabotage and a move against an organization they are at odds with now count. Verified: a campaign went to work for the Bellandi Family, sabotaged the Russo Outfit six times, and went from Associate on $22 a day to Lieutenant on $71.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## The city's own temperature
- Attention is a number attached to one man. Nothing in this city got harder because of what everybody in it had been doing: a campaign of bombings, killings and daylight robberies left the police exactly as interested in week twenty as in week one.
- Scrutiny is the city's temperature, and it is derived from the one thing the whole game already agrees on — what the Herald carried. A killing is worth fourteen, an explosion twelve, a war eight, a police story six, down to nothing for a shop changing hands. A day the paper carried nothing worth noticing takes a point off.
- Past sixty the police stop waiting, and it is announced once and lifts once: *It is not about you. Between the killings, the explosions and whatever was in the paper this morning, the police have stopped waiting for a reason.* Raids then begin up to twenty attention sooner than they otherwise would — held ten below the ordinary threshold, a quiet city comes to the door 0 times in 400 and a city at its worst 29. An arrangement with somebody in that building costs a quarter more and their ceiling falls, because a man with a career to protect wants more for the risk and will be seen with fewer people.
- **The first weights were wrong and the measurement said so.** They produced no crackdown at all in 150 cities at open war over four months, because the city files a story every few days and a point a day of cooling ate all of it. At the corrected weights: a city at peace peaks at 20, ends at 4, and sees a crackdown in 1 campaign of 150; a city at war peaks at 48, ends at 32, and sees one in 51.
- Verified over the HTTP stack: a quiet campaign in a city with nothing worse than a feud sat at zero for twenty-six days, with the Herald carrying nothing but robberies and shops changing hands. That is the mechanic working — the temperature is what everybody has been doing, and in that city nobody had been doing much.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Somebody on the door
- The player's people were an abstraction: signing somebody on added to what the organization was worth in a fight and nothing else. There was no way to put one of them anywhere, so a business being taken apart was something the player could only watch and pay for.
- Putting somebody on a door takes a morning and one of your people, and it is exclusive — a man can only be in one place, so a second business means a second man, and the readiness says which: *Everybody who answers to you is standing somewhere already.* Somebody who dies or stops answering to you comes off the door by themselves rather than leaving a save pointing at a name that resolves to nothing.
- What they are worth is what they are worth: twelve, plus a quarter of their poise and a tenth of how much they mean it. It counts twice against the street and once against an organization coming for the ground, because turning away a thief is a different job from holding a door against thirty men.
- Measured on the same laundry four hundred times: an empty door turns away 213 raids, a manned one 269; an empty one loses its takings to the street 328 times, a manned one 306.
- **They are also the one standing in it.** A raid that gets through now reaches whoever is on the door before it reaches anybody else. Over 300 raids on a manned laundry that cost the man on the door his life 20 times.
- **What the measurement did not support.** Over 200 campaigns of 60 days at war, posting people costs *fewer* men, not more — 57 against 88 — because being harder to raid means fewer raids get through to reach anybody. The premises end at 28,411 condition against 22,278 and 363 holdings against 318. What posting actually costs is the man himself: somebody standing on a door is not available for anything else. The balance test asserts the direction the arithmetic produces rather than the one the flavour suggested.
- Verified over the HTTP stack against a fresh isolated `doorman` fixture (new in `cmd/qa-fixture`): the laundry offered *Put somebody on the door — worth about 27*, the command committed, and the state came back with `posted: {name: Hedda Gruber, trust: 45, worth: 38}` and the action replaced by *Take Hedda Gruber off the door*. The field survived commit, reload and six more days. `cmd/apicheck` ran 200 commands against that save exercising 38 command kinds with no invariant failures.
- Also added `core/roundtrip_test.go`: a reflection-based guard that fills every field of `World`, `Player` and every stored record type with a distinctive value, writes it, reads it back and compares field by field. It exists because `go vet` had just caught five fields of `Commission` sharing two JSON tags — data silently lost on save. Proved it catches that class by temporarily giving two `Pact` fields one tag: *Pact.Since did not survive: wrote 11, read 0*.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Being taken in
- The police in this city took money, stock and premises, and never once took a person. Their own newspaper printed ARRESTS EXPECTED after every raid and no arrest ever followed. Attention could end a career and could not end a week.
- A search that turns up something a prosecutor could use — a still, a room of crates, explosives, stock nobody has papers for — is no longer a fine. Confiscation and a charge are not alternatives: they take the still *and then* somebody answers for it, so the weight of the evidence is read before the search removes it and the arrest is the last thing that happens. A retained commissioner buys the whole outcome, which is what that arrangement was always for.
- The arrest is a decision with three answers. **Go with them:** two to twelve days, longer for somebody they have been watching. **Give them one of your people:** he is gone for the same term and every other man who works for you takes 25 trust off you for it. **Find out what the officer wants:** four figures, and he now knows exactly what you are worth.
- **What a cell actually is.** Not death — you keep everything you had. The map offers nothing anywhere except the station, no work is brought to you, the police do not come back for a man they already have, and nobody robs you in the street. Meanwhile wages come out, supplies run to nothing, and **the takings do not reach you at all** unless one of your own people is standing on the door, which is the second reason that job exists.
- Measured over 200 campaigns of seven days at war: a week you can use ends with $1,833,331 and 12,127 supplies across the run; the same week in a cell ends with $972,316 and nothing. A sentence costs about half of what a week is worth and leaves every business empty.
- Three ways out of the same five days, over the same 200 campaigns: **serving** ends at 28,000 respect and $2,371,884 having spent 1,440,000 minutes; **a lawyer** at 24,000 respect, $2,140,000 and no time at all; **talking** at 15,000 respect, $2,400,000, and −18,000 standing across every organization in the city, with everybody who works for you at zero trust. Time is the cheapest thing a man in a cell has and the only one the street pays for.
- **What the measurement did not support.** Neither ground nor condition moves — 369 holdings against 376, and condition within a percent either way. Seven days is not long enough for a war to take premises off anybody, and the two arms run the clock in different numbers of calls, which is enough to move a figure that small. Those numbers are recorded and not asserted.
- The city added a location for it, Ward Street Station, with its own ways of dying in it and one thing to do when you are not the one inside: bail out somebody of yours, which he notices.
- Verified over the HTTP stack against a fresh isolated `arrest` fixture, and the run tells it better than the summary does. Alex Varga was charged over the still at the Bluebird, went with them, and did five days. In that week his wages kept coming out, **Sofia Doyle stopped being paid properly, signed on with somebody else and took Russo Motor Works with her**, and he came out to $875 less, 20 respect more, a business gone and both stockrooms empty. `cmd/apicheck` ran 200 commands over 34 kinds against that save with no invariant failures.
- Save version 13, with a migration that gives any existing campaign a record for a place the city did not have when it began.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Money you are owed
- Money in this city only ever went one way. The families lean on the player for a share, the player pays tribute, the police fine, the street robs — and in all of it there was never a single person who owed the player anything. Cash sat in a pocket doing nothing between purchases, and the oldest business this trade has did not exist.
- A loan is 35% over seven days, made in person, to somebody standing in front of you. **The size is the skill**: a loan is never more than three quarters of what somebody could service, and what somebody could service is what their position in this city is worth — nine days of what they carry. A lieutenant of a rich family can take $1,700; a docker can take the $120 floor. Never more than 40% of the player's capital is out at once.
- Nearly everybody pays, which is why the business works. **A default is a man whose circumstances changed**, not a man who was always poor: the floor loan is serviceable by anybody alive, so the debts that go bad are the ones where a family stopped paying its people, or where somebody who could find it decided not to — nerve is willingness, and a lender with nothing to lean on is somebody a certain kind of man simply stops answering.
- Then there are three answers to the same moment, and all three are in the room. Over 300 real defaults: **leaning** recovers $421,330 and 25,165 respect for 1,500 attention and 13,500 held against you; **another week** recovers nothing now, leaves 1,800 trust and turns $1,686 into $2,107, which is how this trade actually makes its money; **writing it off** recovers nothing ever, costs standing, and buys 12,300 trust — the one thing money cannot.
- Measured over 200 campaigns of 45 days, walking to every conversation and paying for the journey out of the same day: money in a pocket ends at $423,536; money on the street ends at $1,166,743, with 2,356 attention and 103 debts still standing.
- **Two silent bugs the work found.** `Resent` refuses to record anything against the player by design — the player's answer to an organization is goodwill and to a crewman is loyalty — so `Resent(id, "", …)` was a no-op, in the new lending code *and* in last iteration's fall-guy branch, which had been quietly recording nothing. The fix is the missing feature: `Aggrieve` is the one number somebody can hold against the protagonist personally. It does not announce itself. It arrives as an ordinary robbery in the street with a name attached — a man who was put against a wall is picked as the thief 18 times in 400 against 9 for a stranger — and it fades a point a day, because nothing here is held forever.
- **A functional-UI call, recorded.** Fifteen strangers in one room rendered as fifteen identical buttons is not a choice, it is a wall. New lending is capped at four offers a room; anybody who already owes you is always listed, however many that is.
- Verified over the HTTP stack against a fresh isolated `debt` fixture: $1,249 went to Hedda Gruber, her position collapsed while she was carrying it, the collection recovered $1,416 of $1,686 for 5 respect and 5 attention, the Herald carried ASSAULT REPORTED AT MERCER EXCHANGE with nobody named, and she now holds 45 against the player with the reason attached. `cmd/apicheck` ran 250 commands over 28 kinds with no invariant failures, and now matches named work by prefix so lending is exercised at all.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## A man on the top floor
- The Bellwether Herald was the one institution in this city nobody could touch. It reported everything and answered to no one — the player read their own work in it with no name attached and could do nothing about the fact that it was there at all. Meanwhile the city's entire temperature is derived from what the paper carried, so **the one thing nobody could reach was the thing that decided how hard everything else was.**
- The Herald is now a place on the map with an editor in it, retained by the same machinery as the commissioner and the mayor: an opening payment, $38 a day, a ceiling above which a man with a career will not be seen with you, and a person who can be outbid, killed, or replaced. Officials are now found where they actually work rather than all of them in the market.
- **Pulling a story** takes it out of the paper, and the city's interest in it goes with it, because the city's interest was never anything but what it read. Only today's edition, and only stories about premises of yours or about the police. Every one is a permanent fact about you held by somebody on a weekly retainer: of 400 stories pulled, the arrangement came apart 53 times.
- **Running something about somebody else** costs every one of their premises 14 trade and their organization 4 strength — and files a story, which heats the city like any other. A paper full of crime is a paper full of crime whoever it is about. Of 400 stories arranged about a rival, they worked out who paid for it 81 times.
- **A paragraph about a local businessman** is 6 respect and 5 off what the police think: the cheapest standing in this city and the only kind nobody was hurt for. Everything the player puts in the paper shares a two-day cooldown, because a newspaper that carried the same man's arrangements twice in a week would not be a newspaper.
- Measured over 200 campaigns of 40 days in a loud city: reading the paper like everybody else ends at 17,187 scrutiny with 181 crackdowns; a man at the Herald ends at 15,190 with 158, and costs $451,852 across the run.
- Measured over 200 campaigns of 30 days: leaving a rival alone ends with their trade at 19,400, the city at 3,998 scrutiny and them at 0 standing. A month of stories ends with their trade at 8,532, **the city at 8,940 scrutiny — more than double — and them at −14,825**. The mechanic works on a city the player also lives in, which is the whole point of it.
- **A ceiling the first measurement demanded.** Without one, a month of stories put every rival premises on the floor: 411 trade against 19,400. A newspaper can talk a business down and cannot talk it out of existence, so press damage now stops at 25 trade and the strength loss stops with it. That turned a delete button into a lever.
- Verified over the HTTP stack against a fresh isolated `herald` fixture: a raid was in the morning edition at 45 city scrutiny and 35 attention; the arrangement was opened, the story pulled — the paper empty, scrutiny 39, attention 31 — and then a story run about the Bellandi Family, which put POLICE EXAMINE ACCOUNTS AT PIER 14 on the front and emptied both of their places. **`cmd/apicheck` run against that fixture exercised 71 command kinds in 260 commands with no invariant failures**, the widest coverage the harness has reached, including `retain:editor`, `spike`, `smear`, and the posting and lending work from the last two iterations.
- Save version 14. The Herald has its own ways of dying in it, like everywhere else.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## The page nobody opened
- **The interface had been blank since `45d21d6`.** A `useEffect` was written below the early return that renders while the world is still loading, so the hook count changed between the first render and the second: React error #310, a white page, every boot. Four systems shipped on top of it — posting, custody, lending, the newsroom — each verified over HTTP and through the simulator, each recorded as working, and none of them reachable by the person the work is for.
- Every check in this repository passed the whole time, because none of them ever opened the page. `go test`, `go vet`, `tsc` and the production build are all blind to a hook-order violation; it is a runtime error and only a browser sees it. **The user found it, which is the wrong way round.**
- Fixed by moving the hook above the early return, where every other hook already was.
- **Two guards so it cannot happen quietly again.** `TestNoHookIsWrittenBelowAnEarlyReturn` reads the interface sources and fails on a hook below a top-level early return in a component — narrow on purpose, aimed at the exact mistake. Checked against the shipped file: *main.tsx:52 calls a hook below the early return at line 47*.
- And the page is now served `Cache-Control: no-store`, with the hashed assets beside it marked immutable. Without it the browser keeps yesterday's HTML, which names a bundle that no longer exists, so a fixed build stays broken for a player who did nothing wrong — which is exactly what happened while this was being diagnosed.
- The real lesson is in the process rather than the code: from here, the page gets opened and its console read as part of verifying an iteration, not just the API.

## A menu instead of a wall
- The interface direction changed: the inbox now asks for the interface to be built out and hydrated, for actions to be findable, for the paper to look like a paper, and eventually for the city's events to play out on screen. This is the first iteration under that direction.
- **Ninety actions in one column is not a menu.** Standing in a bar produced an undifferentiated scroll where taking a job, robbing the till, lending a stranger money and paying your own man a share all looked identical and sat in whatever order the switch statement happened to be written in.
- Grouping belongs in the core, not the interface: the core knows what an action is *for*, the interface only knows its id. `GroupOf` answers in one place — work, your premises, people, the street, standing, money, elsewhere — so a new action is grouped the day it is written. **205 actions are offered across the city and every one of them lands in a group the interface renders**, which is what `TestEveryActionBelongsSomewhere` asserts; an action nobody classified still reaches the player rather than vanishing.
- The sidebar renders that: a search box over everything available here, travel lifted to the top because it is what a player reaches for most, headed sections with a line saying what each is for, and everything currently unavailable folded behind *Show 4 you cannot do yet* rather than deleted — a player needs to know that putting somebody on a door exists and why they cannot do it yet.
- **The Herald is now a newspaper.** Masthead, double rule, dateline, late city edition, five cents; two columns with a rule between them; a lead story with a drop cap; halftone cuts with captions. None of it is new truth — it is the truth already committed, set the way the city would have set it.
- The text is hydrated from the same records. Every story kind now carries a **standfirst** in the register of a paper that has to sell itself on a newsstand and cannot say what everybody knows — *"The accused was not represented. The department would not say how the name came to it."* — and a **byline** naming the desk that filed it. `TestEveryKindOfStoryIsSetLikeANewspaper` fails if any of the thirteen kinds falls through to the generic line or is filed by nobody in particular. The paper also carries a date that advances and rolls over month ends.
- **What the browser check found.** The street illustration was reported as unavailable and the whole stage was blank — but only in automation: a backgrounded tab stops `requestAnimationFrame`, so the study never paints and never completes its readiness handshake. Forcing visibility showed it painting correctly. The art was never broken. Two things came out of it: verifying visually now means forcing visibility first, and one stray error no longer condemns the stage — the iframe is remounted once before the player is told anything is wrong.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass, and both screens were opened in a browser and read.

## Who is in the room
- Grouping the actions by kind was only half the problem. Half of what a player does anywhere is done to **somebody who happens to be standing there** — lending them money, collecting it, putting them on, putting them out, taking what they are carrying — and rendering *Lend Perla Mraz money* as one more card in a column of verbs throws away the thing the decision is actually about, which is Perla Mraz: who she is, what she thinks of you, and what she already owes.
- Every action now says **who it is about**. Most of it comes free: work aimed at a person already carries their id after the colon, so the core knows without anybody restating it. The handful that put the name in the label instead — taking what somebody is carrying, recruiting, buying Mara a coffee, sending Leo out — say so explicitly. Measured across the city: **76 actions aimed at a person, every one of them naming who**, and no action claims to be about somebody who is not in the room.
- `PeopleHere` is the other half: everybody standing where the player is, ordered the way a person would notice them — your own first, then anybody who owes you, then people you know, then the rest. Each carries the one line that matters (*Yours · on the door at Bluebird Laundry*, *Head of Russo Outfit*, *In a cell at Ward Street Station*) and, **only for people the player has earned the right to know**, what they think of you, what they hold against you, and what they owe. A stranger's character is not on display, and a test enforces that.
- The sidebar is built around the room now: what you can do here, then the people in it with their work under their own names, then everything that is neither. A debtor's card reads *Stallholder · Careful · owes $1,686 · overdue · thinks of you at 1* with collect, extend and write-off beneath it, and the others in the room are listed by name behind a fold so the city is visibly populated rather than reduced to the people you have business with.
- **This is a waypoint, not the destination.** The next step, from the same conversation: entering a building should open the building, not fill a sidebar — an interior with the people in it as things you can click, and a standardized set of premises actions in the same place in every building so a player learns where to look once. The work here is what that needs: the core already answers *who is here* and *who is this about*.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass, and the screens were opened in a browser and read.

## Stepping inside
- The sidebar was the wrong container and would not have held. The actions grow with every system, and a column that grows with them is a column nobody can read. The direction from the same conversation: **entering a building should open the building** — an interior, the people in it as things you can click, and a standardized set of premises actions in the same place every time.
- *Step inside* is now a view of its own. The room is the whole screen: the interior on the left, everybody standing in it drawn as a figure that can be clicked, the roster of who is present beneath it, and the work in a pane on the right. Picking somebody — from the room or from the roster — puts their name at the top of the pane with what they are to you underneath, and their work under that. The sidebar is not rendered at all while you are inside, because it would only repeat the room.
- **Premises actions are standardized.** Repair, staff, supplies, trouble, operating mode, the door, the still, the safe — always offered in the same order under *These premises*, so a player learns once where to look and never hunts again. Everything that is neither a person nor the building falls under *Everything else here*.
- **There is no image model on this machine.** Ollama carries `qwen3.5:35b-a3b`, `qwen3.5:9b`, `qwen3:14b` and `qwen3:4b-instruct` — text only — and there is no ComfyUI or Automatic1111 installed on any of the usual ports. So the inside of a building is drawn rather than generated, in the same register as the newspaper cuts: flat noir shapes, deterministic from the place, no detail the eye has to resolve. It is a stage rather than a picture, and the point of it is that the people on it can be clicked. Real generated interiors and portraits need something installed; the drawn version is what can be honestly shipped today.
- Each figure is deterministic from the person's own id — build, stance, whether they wear a hat — so the same person is the same shape in the same room every time, and the colour says what they are: gold for your own people, red for anybody who is owed money or holds something against you.
- Verified in a browser against a fresh isolated fixture: fifteen people in the Mercer Exchange, all fifteen drawn and clickable, and picking the woman who owes $1,686 puts *Stallholder · Careful · owes $1,686 · overdue · thinks of you at 1* at the top of the pane with collect, extend and write off beneath it.
- **A note on verifying visually.** A backgrounded tab stops `requestAnimationFrame` and Chrome reports `document.hidden`, so anything canvas-driven looks blank in automation. Visibility has to be forced before a screenshot means anything. That is what made the street illustration look broken when it never was.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass.

## Faces, and one thing at a time
- **Portraits in the room.** A list of names is a list; a list of faces is a room. The portrait was lifted out of `main.tsx` into its own component so the room could use it: six of this city's people are painted, and everybody else is drawn deterministically from their own id, so a face is always the same face.
- **A bug the move exposed.** The painted sprite sheet was reached by `.person .painted-portrait` — scoped to a container that exists nowhere in the room — so Elena Russo and Detective Harlow rendered as empty grey squares in the one place a face is most worth having. The sprite belongs to the portrait, not to whatever it happens to sit inside.
- **A person's pane is theirs alone.** Picking somebody now shows everything you can do with *them* and nothing else. A list that mixes *collect what she owes* with *restock the laundry* is asking the player to do the sorting. The building's own work — the standardized premises strip and everything else here — is what the pane holds when nobody is picked, and the way back is a line that says what is waiting: *← Step away · 22 other things to do here*.
- **Every chip the same size.** The roster is a grid rather than a ragged row, so a card's width says nothing about the length of somebody's name. Measured in the browser: fifteen people, one width.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass, verified in a browser against an isolated fixture.

## What you just did
- The result of an action was a strip along the bottom of the screen, below the fold, in the place a page puts a status bar. It is not a status. It is the answer to the thing the player just decided and the only reason they pressed the button, and it was the least prominent element on the page.
- **The core was not reporting enough to make a better panel possible.** `Result` carried the records the command wrote and nothing about the command itself — not what was chosen, not what it cost — so the interface could show the city's minutes and never the player's decision. It now carries the label the player actually read on the button, the command kind, and what moved: cash, respect, attention and health, **measured across the whole command** rather than announced by whichever rule happened to move a number. Whatever the day charged while the clock ran is part of what the decision cost, and the player is told once instead of diffing two screens.
- The panel sits at the top of the pane the action was taken in, where the eye already is: *YOU DID THIS · Collect from Hedda Gruber · +$1,416 · +5 respect · +5 attention · 1h · Collected at Mercer Exchange*. Gains read green, attention and injury read red, and only figures that actually moved are shown, because a row of zeroes is noise. A death or a danger turns the whole panel red. Everything else that happened while the clock ran folds away behind a count rather than competing with the thing the player asked for.
- The toast that used to repeat the same headline in the corner is gone. Toasts are for what the panel cannot say: an error, or a recovered action after a dropped connection.
- Three tests hold it: the result names the action and the deltas agree with the books; the cost is measured across the whole command; and every ordinary action reports itself, since a blank there is a screen that says something happened and will not say what.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass, verified in a browser against an isolated fixture.

## The camera goes there
- The city told the player what happened in prose and never once showed them. Six places emitted a visual cue and all six called it an *attack*; **a killing, a police raid, a seizure of ground and an arrest — the four loudest things in this game — produced nothing to look at at all.** The most dramatic moments of a campaign arrived as a paragraph in a list, and the Herald reported them the next morning to somebody who never saw them.
- A cue is now the core saying what a moment looked like and where: the kind, the building, **who was in it**, and **the headline the paper will carry**, so the scene and the story cannot disagree — a test fails if a cue promises a headline the Herald never ran. It also carries a *gravity*, the same ordering the city uses for its own temperature, so the interface never has to guess which of five things in one command is the one to show.
- Killings, raids, seizures, arrests, bombings, gunfights and repelled raids all emit one now, and nothing can be sent to a building that does not exist — a moment with nowhere to happen is dropped rather than pointed at an address the map has never heard of.
- **The theatre** takes the player there. The camera goes to the building, the thing plays — muzzle flashes and a shape on the pavement, an explosion that throws the windows out and shakes the frame, a police lamp turning over on a car at the kerb — and then the Herald headline comes up *after* the scene rather than instead of it. It is drawn, because there is still no image model on this machine, and it is skippable and silent under reduced motion.
- **The result moved to a band of its own.** Putting it at the top of whichever pane was under the player's hand was better than the bottom of the page and still wrong: it moved when the layout moved. It now sits in a reserved strip directly under the top bar that keeps its place whether or not anything has happened — *WHAT YOU DID · Let an hour pass · +$10 · 1h · Broken glass at Bluebird Laundry* — so the page never jumps between one action and the next, and it is visible from every screen rather than only the one the action was taken from. It turns red for a death or a danger.
- Verified in a browser against isolated fixtures: an attack on the Bluebird played in the theatre, the band went red, and the same event read correctly in both. `cmd/apicheck` exercised **75 kinds of command in 240 with no invariant failures**, the widest yet.
- Full `go test -race ./...`, `go vet`, `gofmt`, `tsc`, the production build and `cmd/simulate` pass.

## Back issues, and a paper worth keeping
- Asked whether the Herald could show past editions, the answer was no: it was one flat run of the current life's stories, oldest silently evicted at sixty, with a previous protagonist's era unreachable.
- **The archive.** The paper is issues now — one a day, newest first, every life the save still remembers — and the reader can walk back through them or pick one out of an index. A predecessor's paper is readable and marked as theirs, which is the premise of a city that remembers: the player who comes next inherits the news as well as the streets. Within an issue **the biggest story of the day leads rather than the latest**, by the same weighting the city uses to decide how hard it is looking, so a shop changing hands does not lead over a killing because it happened at four in the afternoon. The cap went from 60 to 240 — measured at a little over half a story a day, sixty was about three months, which is one life and nothing across several.
- **What the question actually uncovered.** The live campaign had **one story in twenty-one days**. Backups confirmed it: news 0, then 1, while the history filled to its cap. That looked like data loss, and it was not — `Clone`, the JSON tag and the eviction are all correct, and a 62-day probe filed 23 stories with none lost. Running the *real save* forward thirty days into a war filed fourteen. **The paper was empty because the city was quiet, and every route into it was violence** — a killing, a raid, a robbery, a war. Three weeks of ordinary play produced nothing to read, and nothing to archive.
- **So the paper comes out every day now.** The ordinary edition is trade, civic business and whatever the town is talking about — the price of a good against what it usually fetches, who still holds which post, the district's population off the rating rolls, how hard the police say they are looking, a business visibly in poor repair. All of it read out of state the city already holds rather than invented, so the paper can be trusted the way the rest of the game is. It fills only the space a thin day left: **a morning with a killing in it is never padded with the price of coal**, and a month of ordinary editions moves the city's temperature by exactly zero, because a column about coal is not a reason for anybody to look harder at anybody.
- Measured: forty quiet days now produce forty issues, none of them empty, and never the same lead two days running.
- Verified in a browser against a copy of the live campaign advanced a day: today's edition, then *Day 22 · Back issue* leading on MARA BELL FOUND DEAD with the drop cap and the halftone cut, reached with one click.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass.

## A city you can see is inhabited
- The brief's ultimate goal is a city where every living person is on screen doing what they are actually doing. Two things stood between the game and that, and both were user-friendliness problems of exactly the kind the brief asks to be cleaned up.
- **The address book listed buildings and never once said who was in them.** A city of fifty people read as ten empty addresses. Every card now carries the faces of who is standing there, how many, and what the first of them is doing — *Mercer Exchange · 18 here · 14 not shown · running Mercer Exchange*. Four faces and a count, never all eighteen: the point is to see the place is inhabited, not to read a register on the way past. The register lives inside.
- **The core now says what each person is doing**, read out of the same state as everything else and never invented: standing on the door at the Bluebird, held at Ward Street Station, running The Monarch, taking your money and answering your calls, carrying money that is yours, avoiding you and not doing it well, armed and expecting trouble, looking for a way up. Measured across a full campaign: **fifty people placed, twenty distinguishable things being done**, none of them blank — and what somebody is doing follows what is true of them, so putting a man on a door changes his line and so does the police taking him.
- **Eighteen figures on one floor was a crowd nobody could read.** They overlapped, the back row hid behind the front and the building disappeared. The room draws seven now — the people the player has business with first, then their own, then anybody who owes them or holds something against them, then the ones they know — scaled and dimmed by depth so it reads as a room rather than a row of stickers, with *and 11 more in here* said plainly on the floor. The full register is the roster underneath, where it belongs.
- Verified in a browser against a copy of the live campaign: eleven addresses showing their population, and the Mercer Exchange floor with seven readable figures instead of eighteen stacked ones.
- Full `go test ./...`, `go vet`, `gofmt`, `tsc` and the production build pass.

## Faces worth looking at
- Six of this city's people were painted by hand and the other forty-four were a drawing made out of their id. The drawing was improved first — head shape, hairline, brow, hat, collar, six skin tones and six hair colours instead of one repeated man in four tints — and it was still obviously a drawing beside the painted six.
- **A local image model now makes the cast.** `tools/portraits.py` generates two dozen period portraits offline and bakes them into a sheet; nothing at runtime depends on a model being installed, because the game ships the PNGs. The people of Bellwether are not gendered by the rules — the city says *they* about everybody — so a face is assigned by a hash of the person's id: the same person is always the same face, and any face can belong to any name.
- **The first attempt was not good enough and was thrown away.** SD-turbo generated a full cast in a minute and it read as flat and plasticky next to the hand-made six. A face nobody believes is worse than no face. The second attempt is FLUX.1-schnell at 8-bit on the Apple Silicon GPU through mflux — about fifteen seconds a face, six minutes for the cast, and it sits alongside the originals.
- **Two things had to be worked around.** FLUX.1-schnell upstream is gated behind an account acceptance nobody here can grant, so the weights come from the ungated `mflux-community` 8-bit quantisation. And mflux's Python module layout moves between versions while its CLI does not, so the generator drives the CLI.
- **Tooling moved to mise, on request.** `mise.toml` pins Go, Node and Python and carries the tasks that matter — `build`, `verify`, `race`, `simulate`, `art`, `portraits` — so the whole build is declared in the repository rather than assumed from whatever is on the machine. The image toolchain installs into a project-local `.venv` rather than the user's global site-packages, which is where the first attempt wrongly put it and where it collided with numpy and huggingface-hub versions it did not own. Homebrew is not used.
- Verified in a browser: eighteen people in the Mercer Exchange, seventeen of them wearing a generated face and one the hand-painted original, all distinct.
- `mise run verify` — gofmt, vet, the full Go suite, tsc and the production build — passes.

## Rooms worth standing in
- The inside of a building was flat procedural shapes: a back wall, a floor and two or three blocks that said what the place was for. It was a stage rather than a room, and the room is the screen the player spends most of their time on.
- **Every address has a painted interior now**, generated offline by the same local model as the faces (`tools/interiors.py`, `mise run interiors`) — a bar with booths and bottles, a market hall under a glass roof, a laundry of presses and steam, a station front office, the Herald's composing room. Twelve rooms, 760KB all in, shipped as JPEGs because they are backdrops behind figures rather than things anybody reads. The rooms are painted **empty**: who is standing in one is drawn over the top out of what the core says, so a room can never disagree with the city about who is in it.
- **The drawn room stays as the fallback.** A building added tomorrow has somewhere to stand before anybody paints it, and a missing file degrades to the procedural version rather than to a hole.
- **The figures had to change to survive the backdrop.** Pale shapes that read fine on a flat stage floated on top of a painting like stickers. Against a painted room a person is a silhouette with the light catching one edge, so the colour that says who they are — gold for your own, red for anybody owed money or holding something — moved to the rim where a hard light would actually catch it, and they stand lower, on the floor of the picture rather than in the middle of the air.
- Verified in a browser: the Mercer Exchange as a glass-roofed hall with seven people standing in it, and the roster of eighteen underneath.
- `mise run verify` passes.

## The seven addresses nobody had painted
- Five buildings were painted by hand for the street study. The other seven — the docks, the apartment block, the garage, the Blue Hour, Cypress House, Ward Street Station and the Herald — showed a wireframe box in the address book, which made half the city look unfinished. Two of those seven were places this project added itself and never went back to.
- All seven have a painted front now (`tools/exteriors.py`, `mise run exteriors`), 396KB for the set. Every card in the address book carries a picture; **no wireframes and no broken images**, checked in the browser rather than assumed.
- **The two kinds of picture are now drawn differently on purpose.** A hand-painted cut-out is a model of a building and sits inside its frame; a generated street view is a picture taken from across the road and fills it. Squeezing the second into a frame built for the first made Pier 14 look like a mistake rather than a photograph, so street views fill their frame and take a gradient at the foot of the large version. The wireframe stays as the last resort, so a building added tomorrow still has a card rather than a hole.
- `mise run verify` passes.

## The moments, painted
- The theatre took the player to the building where something happened and showed it as a black box and a few animated shapes. Enough to say *where*, not enough to make anybody look.
- **Eight painted plates now**, one a kind — a body under a sheet on a wet pavement, the blown-out front of a building, police cars at a kerb with torch beams into a doorway, a shuttered shopfront with a notice on the door. 512KB for the set. They are the aftermath or the middle of the thing, framed wide, with **nobody identifiable in them**, because the people who were actually there are the city's business and not the picture's.
- **What moves is still drawn, and had to change.** The plate already contains the body, the onlookers and the cars; drawing another one on top of it reads as a sticker. So on a painted plate the overlay adds only *light* — the muzzle flash with a bloom around it, the explosion's heat, a police lamp sweeping the street and washing the whole frame. The drawn building and its debris remain for any kind nobody has painted, so an event added tomorrow still plays.
- Verified in a browser against a fresh isolated fixture: an attack on the Bluebird played on `scene-attack-v1.jpg` with the caption and the result band beneath it.
- The city's art is 24MB all in: two dozen faces, twelve interiors, seven fronts, eight event plates and the hand-painted originals.
- `mise run verify` passes.

## Fifty names is not a decision
- Opening the People screen on the live campaign turned up something worse in front of it: **the scene that asks for a name offered every living person in the city except the player's own crew.** Fifty rows, identical, in one modal — a wall rather than a decision, and a wall that included a laundress the protagonist had never heard of. The numbering rendered them "09", "010", "011", because it was a literal zero and an index rather than a padded number.
- **A name is something you have a reason to say.** The list is now the people this protagonist actually knows, plus anybody who has given them a reason whether they know them or not: somebody carrying a grudge against them, somebody who owes them and has stopped paying. Their own people are never on it — that is what dismissing somebody is for. It is ordered by how much reason there is, so the name the player is most likely to be thinking of is the first one they read.
- Measured on the live campaign: **fifty names became twenty-four**, led by the two family heads, then the fixer, the driver, the commissioner, the mayor, the detective and the editor. A test holds the shape — never more than half the city, never the player's own, never a complete stranger — and proves the ordering by giving a stranger a grudge and watching them become the first name offered.
- The scene itself now caps its list so the question and the speaker stay on screen while the answers scroll, and sets them in two columns on a wide screen. The numbering is padded.
- `mise run verify` passes, and `cmd/apicheck` ran 220 commands over 69 kinds against a fresh fixture with no invariant failures.

## The people in your life, and then everybody else
- The People screen rendered every living soul as an identical card in whatever order the save happened to hold them: **53 cards, 4.4 screens of scrolling**, no search, no grouping, and the man who works for you indistinguishable from a docker he has never met.
- **The core now answers the question the screen was actually asking.** `Everyone` is the whole city as the player sees it, and each person carries *why they matter* — yours, your crew, they owe you, bad blood, a name everybody knows, an organization, the street — which is also the order they should be read in. `PeopleHere` was rewritten to use the same projection, so a person cannot describe themselves one way in a room and another way on a list.
- Measured on the live campaign: **4.4 screens became 1.5**, 53 cards became 27 with the street's 23 folded behind a line that says how many there are. Your three people first, with what each is doing and what they think of you; then Leo; then the seven names everybody knows; then the organizations; then everybody else. A search runs across all fifty, and each group can be isolated with one click.
- A stranger still gives nothing away — no temperament, no trust, no grudge — but always has a standing and something they are doing, because a person with neither is furniture. Tests hold both: everybody appears exactly once, in group order, with the man who works for you first; and a stranger's character is never on display.
- **A stale scene the fix cannot reach.** `w.Event` is persisted with its choices baked in, so the live campaign is still holding a fifty-name contract scene generated before the previous iteration's fix. It renders once more as it was written; every scene opened after it is built by the new rule.
- `mise run verify` passes.

## What am I worth and what is this costing me
- The Ledger answered its own question with two figures and sixty rows of undifferentiated history — **4.9 screens of scrolling**. `DailyCost` added nine things together and returned one number, so a player losing money had to guess which of the nine it was.
- **The books are the core's arithmetic, not the interface's.** `Books` reports what comes in a day and what goes out, broken into lines that say what each is made of — *Staff, 11 hands across your premises, $97* — plus what is on hand, what a fine cannot reach, what is out on the street against what is due back, and what is outside the city. A test holds the two things that make it trustworthy: **the lines add to the total, and the total is the number the clock actually charges**. Nothing that costs nothing is listed, because a page of zeroes is noise.
- The history is searchable and filterable by kind, grouped under the day it happened, showing the last twelve with the rest behind a line that says how many. **4.9 screens became 1.7.**
- On the live campaign it reads: coming in $1,140 a day from 4 businesses, going out $206, net $934, $3,778 on hand and all of it reachable — then rent, security, crew, staff and your own people, adding to $206.
- `mise run verify` passes.

## The organizations, said plainly
- Families was the opposite problem to the other screens: not a wall but a thin one. Three cards, each carrying a strength, a money word and **a bare number for standing — "+45", "−63" — with nothing saying what either meant**, nothing about what any organization holds, and nothing about who is fighting whom. All of it was already in the city and simply never asked for.
- An organization now reports what it holds by name — premises are the most public thing it has, so they need no informant — who it is at war or at odds with, what its number means in words (*Hostile*, *They think well of you*, *You are as good as one of theirs*), and its leader's face. **How many people answer to it still needs somebody inside**, which a test enforces: what you cannot count, you are not told.
- The player's own organization is listed first and says *Yours* rather than reporting a number about how much it likes itself. Its card carries the player's own face, because its leader is the player, who is not one of the city's people and has no id to draw from.
- A section beneath says who is fighting whom and since which day, out of the conflicts the city already tracks.
- On the live campaign: *Nico Ward's people · 100 of a hundred · $4,202 · 3 people · 4 holdings · Bluebird Laundry, Russo Motor Works, The Blue Hour, Cypress House · at war with Russo Outfit*, then the Bellandi Family who think well of you and the Russo Outfit who are hostile, and the two wars beneath.
- `mise run verify` passes.

## A guide that cannot rot
- The Guide was prose written when this game had eight actions, and it had rotted into something actively misleading. It still told the player that *"broader autonomous family politics"* was future work — in a build where two families fight their own wars, hold ground, run coups and collapse entirely. It described death wrongly, saying what you built "becomes independent", when a protagonist's estate now passes to the strongest of their own people as an organization the next one can deal with or fight. It mentioned none of the twenty-odd systems added since. **A guide that misdescribes the game is worse than no guide.**
- So it is not prose about the game any more. It is **the game reporting on itself**: twelve things a player might be doing, each answered by the same readiness function that answers the button. Somewhere to start, somebody who knows people, premises, a name of your own, people who answer to you, somebody on the door, money on the street, a still, somebody in the building, an understanding, somebody else's ladder, the chair. Each is open, done, or refused **in the game's own words** — *Nobody is taking anybody on*, *You answer to nobody*.
- It cannot go stale, and a test proves it: change the world and the guide changes with it. A new arrival is told to buy premises; a man with two businesses is not; put somebody on a door and the guide has noticed by the next read. Another test holds that no step is ever blank — closed without a reason is a screen that says nothing.
- Beneath it, the five rules that genuinely do not change, which are the only part of a guide safe to write down once.
- On the live campaign: five done, five open, two shut with their reasons.
- `mise run verify` passes.

## Eight figures that never said what they meant
- The top bar grew a stat at a time to eight of them and never once explained any. A new player reads *PRESENCE 253* and has no way to learn what it does short of dying of it. Worse, the city's temperature was drawn as a number above the words *"The city is ordinary"*, which parses as two ordinary cities.
- **The explanations belong in the core, because they are about rules and they contain the thresholds.** Written as copy in the front end they would rot exactly the way the Guide rotted. `Dashboard` reports each figure with its label, its value, and a sentence carrying the real numbers — *"Past 45 they come to the door; past 80 they take the premises"*, *"Past 60 the police stop waiting and raids begin up to 20 attention sooner, for everybody"*. A test fails if any explanation stops containing its own constant.
- Every figure now says whether it is currently worth worrying about, so low health and high attention colour themselves rather than waiting to be noticed, and a healthy solvent player is warned about nothing.
- The city reads as a state with its number beside it, in **the same word the city description already uses** — a test holds the two together, because the bar saying *ordinary* while the city said *watchful* is exactly the kind of drift nobody notices.
- **A regression the browser caught.** Moving money formatting into the core dropped the thousands separator the interface used to add: `$3778`. The core writes sums the way people read them now, with a test over seven cases.
- **The theatre holds for as long as the moment is worth.** It played a robbery and a killing for exactly the same 2.6 seconds. The city already knew one was worth more than the other — `Gravity` has ordered these since the theatre was built — so a killing now holds 4.7 seconds, an explosion 4.4, a robbery 2.4, and anything nobody has classified 1.8. Bounded so nothing is too short to read or outstays its welcome.
- `mise run verify` passes.

## The whole city, not one block of it
- The street view showed **five buildings**. The city has twelve, and two of the seven it left out — Ward Street Station and the Bellwether Herald — were places this project added itself and never went back to. Half the city was reachable only through the address book, and a view called *the street* quietly told the player the rest of it was not there.
- **The city view is every address now**, painted, grouped into the district it stands in, each carrying who is actually in it as faces and a count, and lit by the hour the clock says — the same 1440-minute day everything else runs on, so the fronts darken through the evening without the interface deciding anything. *You are here* is marked, districts you have not opened are dimmed and say so, and double-clicking the one you are standing in steps inside.
- The animated Old Harbor block is not thrown away. It was always one block of a city with three, and it is linked from the foot of the view for anybody who wants to look at it; the component that embedded it as though it were the whole city is gone.
- **Three guards so the art cannot fall behind the city again.** Every address must have a picture, every address must have an interior to stand in, and every moment the theatre can play must have a plate — each failing with the `mise` task that fixes it. These are what stop the next location being added with a wireframe box, which is exactly how five of twelve happened.
- Verified in a browser: three districts, twelve fronts, twelve pictures, no blanks and no broken images.
- `mise run verify` passes.

## What each place is doing, and one door instead of two
- **Two views onto the same twelve addresses.** Once the city view showed every address with its people, the address book was a second door to the same room: the same places, the same faces, one of them painted and one of them a list. The address book is gone and its district headings live on in the city view, which is the only city screen now. One toggle instead of two, and one fewer thing to learn.
- **A place says what it is doing.** The city knew that one address was out of soap, one had a press broken, one had a still running in the back and one had a man on the door, and showed none of it. Six badges on a card would be a wall again, so a place gets the same discipline the people got: **one line, the most important true thing about it right now**, in the order a proprietor would worry — the police can take it, trouble, out of supplies, short-handed, wants repair, a still running, somebody on the door, the regulars have gone, and otherwise how well it is trading.
- **A passer-by only sees what is visible from the street.** Somebody else's casino can be out of stock, short-handed, in trouble and running a still and it reads as nothing at all — but a boarded window is not a secret, so a wrecked building says so whoever owns it. A test holds both halves.
- Verified in a browser against a copy of the live campaign: twelve fronts, five of them with something to say, two of those flagged — *Visibly in poor repair* at The Mariner, *Wants repair at 55%* at Cypress House.
- `mise run verify` passes.

## The person you click is the person in the roster
- The people standing in a room were anonymous silhouettes drawn inside the picture. Every one of them had a painted face two inches below in the roster, and no way to tell which shape was which — so the room was a diagram of how many people were present rather than a picture of who.
- **The figures moved out of the SVG and into the interface**, which is the only way each can wear its own face. A person on the floor is now their portrait over a coat, sized and dimmed by how far back they stand, named on hover, and clicking them selects the same person the roster does. What was drawn stays drawn: the procedural room is still the fallback under a painted interior that has not been generated.
- The colour that says what somebody is to you moved to the coat — gold for your own people, red for anybody owed money or holding something — so it still reads at a glance without a legend.
- Verified in a browser: seven people on the Mercer Exchange floor, seven painted faces, every one inside the frame, *and 11 more in here* said plainly, and clicking Elena Russo on the floor opens Elena Russo in the pane.
- `mise run verify` passes.

## What the journey costs, before you commit to it
- Twelve addresses on one screen and no sense of how far any of them was. A place across town and one on the next corner looked identical, and the journey time only appeared **after** the place had been selected and the travel action read — which is to say, after the player had already decided where they were going. The clock is the scarcest thing in this game, and a cost belongs on the thing it is a cost of.
- Every address now carries how long it takes to get there from where the player is standing, **by whatever they actually travel by**, in minutes and in the words somebody who lives here would use: *Round the corner*, *A short walk*, *The other side of the district*, *Across town*, *The far side of the city*. A test holds that twelve addresses read as at least three different distances, or the words say nothing.
- Somebody with a car sees what the car is buying them: *35 min driving, 52 on foot*. A test proves a car never makes a journey longer and makes at least one shorter — and that distance follows where the player is standing rather than being a fixed table, by walking them across the city and measuring the same trip in reverse.
- **A misleading field name found on the way.** `Player.CarWear` is the car's *condition*, not its wear: zero is a wreck and a hundred is a car that runs. The first version of the test set it to zero to mean "no wear" and got a car that would not start. Noted where it will be read next.
- On the live campaign, standing at Mercer Exchange: Saint Agnes round the corner at 15 minutes, The Monarch a short walk at 20, Pier 14 the other side of the district at 60.
- `mise run verify` passes.

## The people the scene is about
- The theatre named who was in a moment and then drew nobody: the cue carried *Detective Harlow* as a string and the scene put the words under the picture. **A scene about somebody that cannot show them is a scene about nobody.**
- A cue's actors carry their **ids** alongside their names now, because a face is drawn from an id. The theatre shows them under the plate: the arrest is Ward Street Station, the caption, Detective Harlow's face, and then the Herald headline. A test holds that every name a cue carries resolves to a person the city can find and draw — a cue naming somebody who does not exist would fail rather than rendering a blank square.
- **What the scene deliberately does not show.** A killing names its victim and nobody else. The city knows who arranged it and the paper is careful never to print that; showing the killer's face in the theatre would leak exactly what the newspaper is written to withhold. The scene stops where the paper stops.
- An anonymous event still names nobody, which is correct: a window put in at the Bluebird by people who left no proof shows the plate, the caption, and no faces at all. Verified both ways in a browser.
- A `killing` fixture was added for QA, and finding that a fixture-time killing produces no scene taught something worth writing down: **`VisualCues` is `json:"-"`**, so a cue only exists inside the command that made it. A fixture can set up the state that will cause a moment, never the moment itself.
- `mise run verify` passes.

## What the dead leave behind
- The death screen was the last screen with no interface pass, and it was **stating something that had stopped being true**: *"Your properties pass to your former organization. Another life begins without your money, rank, or authority."* Succession has not worked that way since `succession_estate.go` — the estate goes to the strongest of the player's *own* people as a new organization. The most dramatic moment in the game was reading a rule that had been replaced.
- Worse, `Die()` was computing `inherited := w.Inherit()` and then throwing it away with `_ = inherited`. The city decided who took over and then never told anybody, including the dead.
- `Death` carries an `Estate` now, and `core/epitaph.go` answers the screen from committed facts: what they died of, how long they lived, what they earned, **what became of what they built**, what still stands in their name, what the paper actually printed about them, and what the next one starts with. Three tests: an estate that exists is named, somebody who had nobody leaves nothing, and everything the screen says is traceable to a record rather than to a sentence written once.
- **A silent bug found by the test, not the screen.** `Inherit()` returned the *successor's name* where `EstateName` expects the estate **id**, so the lookup returned `""` and `Death.Estate` was always empty. It now returns `"estate:"+successor.ID`.
- Verified against an isolated `dead` fixture on port 8892: *Hedda Gruber's people took it over. It holds their premises, keeps their quarrels, and is worth rather less than it was, because most of what it was worth was them.* — under four headlines the city had genuinely printed, including the raid on the Bluebird.
- **The new-life flow was reporting a bereavement as a transaction.** Clicking through to the next life showed the result band saying *−$8,770, −65 RESPECT*. Those deltas are measured across a change of protagonist: a man arriving in the city with ninety dollars was being told he had lost nine thousand. He had not; that happened to somebody else. `new_life` now reports no deltas, held by a test.
- `mise run verify` passes. 100-run simulation unchanged from baseline: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors. `cmd/apicheck` clean.

## A setting the game obeyed and the player could not reach
- The Settings screen had never had a pass. It was two cards of developer prose — *"Model output proposes bounded opportunities"* — and one button.
- Underneath it was something worse than untidy. `black-ledger-motion` gated **the theatre**: the modal that takes the screen after a killing and holds it for several seconds. The commit path read it on every action. **Nothing in the interface had ever written it.** The most intrusive thing in the game could be switched off only by opening developer tools.
- A guard test now holds the rule: every stored preference the interface *reads* must also be *written* somewhere a player can click. Scratch storage is exempt by name — a command being replayed after a refresh, how far through the paper the reader has got — and the test fails if it finds no preferences at all, so it cannot pass by reading nothing. It caught `black-ledger-motion` before the fix, as intended.
- Settings is now four rows in one frame, each one decision: **Scenes** (with the switch that was missing), **Voices**, **The storyteller**, **This life**. Same shape, control always in the same place. The model paragraph survives, demoted and honest: it writes encounters and cannot touch money, time, injuries, property or death.
- Turning scenes off also cancels anything already playing and stops the travel animation, and it takes effect immediately rather than on the next reload — the gate reads a ref, not storage.
- **Verified both ways in a browser** against an isolated fixture on 8892, by wrapping `fetch` so a committed action came back carrying a high-gravity cue: with scenes off the cue arrived and no `.theatre` appeared; with scenes on the same path produced *IT HAPPENED AT — Saint Agnes*. No console errors on load or during either run.
- A stray `src/styles.css` was created and folded back into `src/style.css` before commit; the project has one stylesheet and now still does.
- `mise run verify` passes.

## The terms of a deal, as figures rather than as a sentence
- The scene modal is the surface a player sees most and had the oldest styling in the game. The real fault was not the styling: **the one decision the scene exists to pose was the only thing on it that could not be read.** Each way of doing a job carried its terms as prose — `$140 · 75 minutes · +6 respect · +4 heat · At 15 heat, police may stop completion. · Russo Outfit standing +6; rival standing −3.` — and a scene offering three approaches printed that same run-on string three times, with the words that actually differed buried in the middle of it.
- `Choice` now carries `Pay`, `Minutes`, `Respect` and `Heat` as **numbers**, taken from the same `Effect` the command will apply. A test holds that every choice's advertised figures equal the effect the job pays — a scene must not quote a price the city will not honour — and that the trade is real in both directions: pressing pays more, costs more attention and takes less time than being careful, or it is not a choice.
- Conditions that hold **however** the job is done — the venue, the police at 15 heat, who gains standing — moved to `Scene.Conditions` and are stated once above the choices. Two existing tests asserted the player was told these things; they still assert it, at the place it is now said.
- `Choice.Reason` replaces the interface's hardcoded `' · Not enough cash'`. That guess happened to be correct, which is the dangerous kind of wrong — the core is the only thing that knows why something is refused, and it now says so with the actual figures: *Not enough cash: this takes $150 and you are holding $40*.
- The modal renders the terms as the same scannable chips the rest of the game uses, attention in red against gold. Verified in a browser on two isolated fixtures: a new `offer` fixture (a courier job with two approaches — **$75/45min/+3 attention** against **$60/75min/no attention** against **$95/30min/+8 attention**, a comparison that can be made at a glance) and a new `audience` fixture (Bellandi asking for $150 from a player holding $40, both priced choices refused in the city's own words). No console errors on either.
- `mise run verify` passes. `cmd/apicheck` exercised 19 kinds of command with no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## Nobody in this city walked anywhere
- The largest thing still missing from a city meant to be lived in: **its people stood still.** A person had an address and kept it for life. The woman who ran the laundry was inside the laundry at every hour of every day; a man who had spent six weeks deciding he hated somebody across town never once went to look at him. The only thing that ever moved anybody was a takeover, and it moved them instantly — in one building, and in the same minute, in another.
- `core/errands.go` gives people errands. `NPC` gained `Heading`, `Arrives` and `Errand`; `Location` deliberately stays where they set off from, so everything that reasons about where a person belongs keeps working while they are out. `SetOut` runs twice a day from the clock, `Arrivals` on **every** clock step — a journey is shorter than a day and the player should be able to watch one finish. Journeys are on foot: nobody but the player has a car.
- Reasons to be elsewhere, in the order a person would weigh them: the work that feeds you (whoever runs a place should be in it), the office you hold, **the man you have not forgiven** (the first time in this game a grudge makes anybody do anything of their own accord), and minding your family's ground.
- **Two false starts, both caught by measuring rather than by looking.** First the city never moved at all: everybody starts exactly where they belong, and a rule with no reason to fire is a system that ships as decoration. Adding "mind your family's holdings" produced the opposite — *street occupied 100% of samples*, because a man who walked across the city to mind a holding became its only minder and was then summoned straight back to his leader. **The summons rule was removed**, and its epitaph is in the source: nobody moves without a reason that ends.
- Measured at scale over twenty cities and three weeks each: **40 journeys completed, somebody on the street in 1% of 3,360 samples, never more than 2 at once.** The city settles — which is correct. Movement follows change: a test holds that when a holding changes hands, somebody walks to it, because that is the war becoming visible on the street rather than only in the ledger.
- The city view gained an **OUT ON THE STREET** band — portrait, where from, where to, why, minutes still to walk, and a marker showing how far along — shown only when somebody is actually out.
- Verified in a browser on an isolated `street` fixture: Otto Weiss walking from The Monarch to Bluebird Laundry and Hedda Palma from Mercer Exchange to The Monarch, both because ground had changed hands. Buying Mara a coffee spent 30 minutes; **both walks finished, both arrived, and Otto appears among the people at Bluebird Laundry.** No console errors.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged from baseline: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors — people walking changes who is in a room, not what the city is worth.

## Taking the camera there, instead of away
- The stated goal for a loud moment is that *the camera is taken there* — to the building it happened in — and the headline comes up afterwards. What the theatre actually did was wash the entire city view to near-black (`.theatre{position:absolute;inset:0;background:#050808f2}`) and draw **its own building**: a flat black rectangle with lit windows, a squad car and a body, in an SVG, on top of a city that already has twelve painted addresses in it. That is the camera being taken *away* from the city and pointed at a drawing of it.
- The city is the stage now. `CityStreet` takes a `spotlight` — which address, what kind of moment, how far through it is — dims every other front, lights that one, scrolls it into view, and plays the effect **over the real painted front**. Only light is added: a flash for a killing, a shockwave and a shake for an explosion, a police lamp sweeping for a raid or an arrest. Drawing a body on top of a painted photograph reads as a sticker; light reads as something happening.
- `Theatre` no longer draws a building at all. It is a band along the bottom — the place, the painted plate for that kind of moment, the caption, whoever the cue names, and the Herald headline after the hold — with `pointer-events:none` on the band and `auto` on its controls, so the rest of the city stays visible and usable while a moment plays.
- A text guard in `cmd/blackledger` fails if the theatre goes back to covering the city with an opaque wash that swallows its clicks, or if the street loses the ability to single out one address. It failed on the old CSS before the change, as intended.
- **A dead mechanism found while verifying.** "Replay recorded scene" set a `sequence` state that rendered a one-line strip reading *Recorded event · <caption>* and nothing else — it never drove the camera, so replaying a moment showed no moment. The control now replays the highest-gravity cue through the camera properly, and the `sequence` state is gone rather than left as a path that looks like it does something.
- **A false alarm worth recording.** `.theatre-frame{display:none}` looked at first like a rule that had been silently hiding the whole stage. It is inside `@media(prefers-reduced-motion:reduce)` and is correct; the grep that found it had stripped the media query. Checked before claiming it.
- Verified in a browser on an isolated `arrest` fixture: after going with them, the city view scrolls to **Ward Street Station**, that front is lit and gold-bordered while the other eleven dim to 28%, Detective Harlow's face is under the caption, and *MAN CHARGED AFTER DISTRICT SEARCHES* arrives after the hold. Ten consecutive samples of the lamp's opacity gave **ten distinct values** — the light really is sweeping the painted front rather than sitting at a fixed tint. Dismissing puts the city back to full brightness. No console errors.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## A room that changes behind your back
- People move between buildings now, and a player standing in a room never noticed. Somebody they had been talking to was simply not in the list any more; somebody they had never seen was in it, with no explanation. That is a list changing behind your back rather than a place you are standing in.
- `Coming` records somebody arriving in or leaving **the room the player is in**, and only that room: everything else in this city happens at a distance and is read about, and a room the player is not in changes without remark, which is what a room you are not in does. A test stands the player at the Herald — an address with nobody in it who has anywhere to be — and fails if they are told about anything at all.
- It rides on the command like `VisualCues` do: `World.Comings` is `json:"-"`, cleared at the top of every command, and copied onto the `Result`, so the interface reads it beside what the player did. The existing round-trip guard caught it immediately and it is now listed as deliberately unstored, with the reason.
- The line is rendered over the room in the Interior view: *← Leo Carver leaves Bluebird Laundry, due at Saint Agnes.*
- **A bug that only opening the page could have caught.** The traffic line rendered perfectly — 1,284 pixels down the page, below the fold, invisible. `.interior-stage` is a grid whose children are placed explicitly by `grid-column`/`grid-row`, so a new child was auto-placed after everything else. A DOM test asserting the element exists would have passed while no player ever saw it. It now claims the room's own cell and sits over the top-left of the picture.
- **A test that was wrong before the code was.** The first version stood the player at `club` and asserted nothing would be reported — then failed because somebody left The Monarch, which is where the player was standing. The code was right and the test had picked the wrong room.
- Verified in a browser on a new `room` fixture (the player inside their own laundry at 11:55, five minutes short of the half-day when the city's people set off): one ordinary action crossed noon, Leo Carver left for Saint Agnes, the line appeared over the room, and he was gone from the floor. No console errors.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## The interface was confidently wrong about where people were
- Giving the city's people journeys left a lie behind it. A person keeps the address they set off from until they arrive — deliberately, so everything reasoning about where somebody belongs keeps working — but the People screen reads that address straight off the record. So it said **Mara Bell is at Saint Agnes** while she was somewhere between there and the garage, and a player could spend forty minutes crossing the city to a room she was not in. An interface that is confidently wrong is worse than one that says nothing.
- `Presence` gained `Walking` and `Minutes`. Somebody on the street is reported as *On the way to Russo Motor Works*, and **`WhereID` now points at where they are going** rather than the door they walked out of — that is the only place they could be met, because reaching a man in the street is not something this game models. `doingNow` says it before anything else: a man on the street is not on a door, not on duty and not behind a desk, whatever his job is.
- **A worse half of the same bug, found by writing the test.** The player could still *act* on somebody who had walked out — buy Mara a coffee while she was halfway across the city. Every action about a person is written at the place it belongs to, and not one of the forty call sites checked whether the person was still standing in it. Rather than adding a check to each, one rule now sweeps the finished list: any action whose subject is out walking is disabled, with the reason in the city's words — *Mara Bell is out on the street, walking to Russo Motor Works — 10 minutes out*. The next action about a person cannot forget.
- Verified in a browser on two isolated fixtures. On `street`, the People screen shows two walkers in gold with *↗ Walking to Bluebird Laundry, 15 minutes out*, each card on screen. On a new `gone` fixture — the player standing in the bar with Mara, five minutes short of the half-day she is due elsewhere — one action sends her off, and afterwards all three actions aimed at her carry the same honest refusal, while the room itself says *← Mara Bell leaves Saint Agnes, opening up at Russo Motor Works* and *→ Ennio Zanetti comes in*. No console errors.
- **What the player is told, and where.** The room says somebody left, the street band says who is out and where they are going, and the People screen says what they are doing and how long until they get there. The disabled actions are the backstop rather than the announcement: the interface hides subject actions for people who are not in the room, so their reason is what the core answers with rather than what the player reads first.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## The people crossing the city, drawn on it
- The last piece of the living-city picture. People have crossed the city over real time for several iterations, and the street only *listed* them: names and minutes, in a box, above a picture of the city they were supposedly walking through.
- Each walker is now drawn on the street between the front they left and the front they are going to, at the fraction of the way they have actually got. The fronts are laid out by the browser across three district blocks, so the only honest source for the two ends of a journey is **where they were actually drawn** — the positions come from measuring the two fronts and interpolating between their doorways (bottom centre, because people walk on the pavement), recomputed when the walkers change and when the window does.
- **Nothing animates between actions, because nothing moves between actions.** The clock in this game is stopped until the player commits to something; a figure gliding along the pavement while the clock is paused would be the interface inventing time the simulation has not spent.
- **A bug found by looking rather than by asserting.** Both figures were placed perfectly and both faces were blank squares. `Portrait` renders a `<span>`, so the label rule `.walker-figure span{…background:#0d1413d9…}` matched the portrait itself and the `background` **shorthand reset its face to none**. The label has its own class now, and the reason is in the stylesheet where the next person will read it.
- **The hook guard cried wolf, and was made precise rather than dodged.** `TestNoHookIsWrittenBelowAnEarlyReturn` failed on `if (!from || !to) return [];` — six spaces deep inside a callback, where a return cannot change how many hooks run. Reformatting the code to slip past a guard is the wrong instinct, so the guard was pinned to the one or two spaces of indentation a component body actually uses, and `TestTheHookGuardKnowsWhichReturnsMatter` now tests the regular expression directly against both shapes. Checked against `git show 45d21d6:src/main.tsx`: the tightened pattern still matches the exact line that caused the blank page.
- **Verified geometrically, not by presence.** For each walker the browser was asked where the figure actually is relative to the two fronts: Otto Weiss **0.40** of the way from The Monarch to Bluebird Laundry against a state progress of **0.40**, Hedda Palma **0.50** from Mercer Exchange to The Monarch against **0.50**. Faces confirmed painted (`faces-noir-v1.png`, `600% 400%`) and both figures on screen. No console errors.
- The OUT ON THE STREET band stays. It carries what a figure on a pavement has no room for — why they are going and how long they have left — and is the readable version of the same fact for anybody not reading the picture.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## Your own people were the only ones who could be in two places at once
- **A guess that turned out to be wrong, checked before acting on it.** The assumption was that the city's people mind their family's ground and the player's people do not, because the player's premises are owned differently. They are not: `PlayerOrganizationID()` returns the same `player:N` string a property carries as its owner, so the minding rule has always reached the player's people. That is now held by a test instead of by luck.
- What was actually broken was **posting**. `Post` set the man's location to the door in the same minute the order was given, however far away he was standing — the player's own people were the only ones in this city who could be in two places at once, and the door was defended from the moment of the *decision* rather than from the moment somebody was standing in it.
- Sending somebody to a door is a journey now. He walks, he is on the street with everybody else and the city can say why (*sent to stand on the door at Bluebird Laundry*), and `PostingDefenceAt` returns **0** until he arrives. Somebody already standing in the place is on the door at once, because there is nowhere for him to walk.
- The interface stops claiming a door it does not have. The address reads *Hedda Gruber is on the way, 25 minutes out* rather than *is on the door*, and the button says what it will cost before it is pressed: *Hedda Gruber is at The Monarch, 25 minutes away. The door is worth nothing until they get there.*
- **Four existing tests were asserting the teleport.** They posted a man and immediately measured the door. They now post and let him walk, through one `postAndArrive` helper — the guarantee they were written for (a man on a door is worth something) is unchanged; only the timing is. With him actually in the doorway the measured effect is the same as before: **400 raids on the same laundry — an empty door turns away 213, a manned one 269**; and of 400 attempts by the street, an empty laundry loses its takings 328 times against 306 with somebody on it.
- Verified end to end in a browser on a new `post` fixture (the player in their laundry, their one man across the district at The Monarch): the order sends her walking, she is drawn on the street between the two fronts, the laundry says she is on the way and is worth **0**, and after the clock runs she is *on the door* and worth **38**. No console errors.
- A `placeName` helper already existed in `press.go`; the duplicate this iteration briefly added was removed rather than left as a second answer to the same question.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged: defiant 52 deaths, investor 0, reckless 82, worker 0, 0 errors.

## "Leo heads out", and Leo did not head anywhere
- The log line said he headed out. He stood exactly where he had been standing, the money appeared two hours later, and the city reported that he had come back from somewhere he never went. Everybody else in this city walks; the man the player pays was the last one who did not have to, and he is the one they have most reason to watch.
- Collections have an address now — the player's best-earning premises, or the corner the work comes from when they own nothing, because **work with no address cannot be watched**. He walks there, does two hours of doors, and walks back to wherever he set off from. `TaskHome` remembers where that was, so the city puts him back rather than leaving him wherever the work happened to be.
- The button says what it will cost before it is pressed: *Leo Carver walks to Bluebird Laundry — 25 minutes — and is not here while he is doing it.*
- **A difficulty change, measured and explained rather than waved through.** The 100-run simulation moved on one profile: defiant deaths **52 → 58**. That is a seeded simulation, so it is not noise in the sampling sense, but it could still have been a diverged random stream — so it was checked at 600 runs, where it holds: **317 → 348 deaths (52.8% → 58.0%, about 2.5σ)**. The other three profiles are unchanged.
- The **mechanism was found rather than assumed**, by probing the action list at each stage: while he is walking, `crew_bonus` is refused by last iteration's rule — *Leo Carver is out on the street, walking to Bluebird Laundry — 10 minutes out* — and becomes available again the moment he arrives. The simulated defiant player tops up loyalty whenever it falls below 50 and now sometimes cannot, so loyalty slips, collections get refused, and the policy falls through to riskier work. **This is the rule working**: you cannot hand a man a bonus while he is halfway across the city, and sending him out is now a decision with a cost rather than a button that prints money.
- While he is doing the round the city says *Collecting at Bluebird Laundry* rather than *Driver, on duty at Bluebird Laundry* — what he is actually doing outranks what his job is called, which is the same rule already applied to people on the street.
- Verified in a browser on a new `round` fixture, through the whole arc: **Walking to Bluebird Laundry, 10 minutes out** (drawn on the street, face painted) → **Collecting at Bluebird Laundry** (nobody on the street) → **Walking to Saint Agnes, 10 minutes out — coming back from the round** (drawn again), with the task cleared and the round paid. No console errors.
- One looseness noticed and left alone rather than quietly changed: `crew_bonus` is offered wherever the player is standing, not only where the man is. That predates this work.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures.

## Rechecking the guide against a game that had changed under it
- The Guide was rebuilt once already, as the game reporting on itself rather than prose about it, and that structure held: every milestone is still answered by the same readiness function that answers the button, so none of them had gone stale. **What had not been checked is the part of it that is still written down once** — the five rules, and the sentences describing each step.
- **The city walks now and the guide never said so.** A new player could read every rule, cross town to the address the screen named, and find the room empty. There is a sixth rule: *People are not furniture. They walk between buildings on their own errands, and anybody out on the street is not at the address they left and cannot be dealt with until they arrive — including your own, when you send them somewhere.* A test fails if the rules stop mentioning walking or the street.
- **The rules quoted numbers as literals.** *"Past 45 the police come to the door; past 80 they take the premises"* was typed, not computed, next to `RaidThreshold` and `ForfeitThreshold` — the only thresholds this game promises are public. They are computed now, and a test holds that the rules mention both constants and that **no number appears in them that is not one the game actually uses**. Confirmed the guard bites by reverting the sentence to a drifted literal: it failed four ways, naming both missing thresholds and both invented ones.
- **The first thing a new player is told is now the first thing a new player gets.** The courier's $45 lived as a literal in four places — the payment, the log line, the button and the guide sentence — beside an unrelated `45` for how long the job takes, which is the shape of a number that eventually stops agreeing with itself. `CourierPay`, `CourierRespect` and `CourierMinutes` separate them, and a test reads the guide, reads the button, then does the job and checks all three agree.
- The door step now states the cost the last iteration gave it: *They have to walk there first, and the door is worth nothing until they arrive.*
- Verified in a browser on an isolated fixture: the Guide renders on screen with the new rule at the bottom, the recomputed thresholds reading 45 and 80, and the opening step quoting $45 and 2 respect. No console errors.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures. 100-run simulation unchanged from the previous iteration's baseline: defiant 58 deaths, investor 0, reckless 82, worker 0, 0 errors.

## One question, asked once: is this person here to be dealt with?
- The sweep added when people started walking asked exactly one thing — is this person out on the street? — and every other way of being unreachable went on being ignored. Probing the action list at the edges turned up three:
  - **A dead man stayed on the books.** `Kill` never touched `Player.Crew`, so after the city shot Leo Carver the crew list still held him, the buttons still offered him work, and the order was still accepted. Under the previous iteration's code that paid $65 from a corpse; under this one it silently did nothing, which is not better.
  - **A man in a police cell could be sent out on collections** and paid a bonus through the bars.
  - **A man across the city could be dealt with**, while the same man one step onto the pavement could not.
- `OutOfReach(id)` is now the single question, asked of every action with a subject after the list is built: dead, held at Ward Street Station, or out on the street, each in the city's own words. `Kill` takes the dead off the player's books, because somebody who is dead is not on anybody's.
- **A rule considered and rejected, recorded so it stays rejected.** The obvious symmetry — you must be in the same room — was implemented, broke four existing tests and the opening (Leo starts at the bar, the player at the Mariner), and was taken out again: requiring presence is *a change to how the game plays*, not a correction of something it was claiming falsely. The line that does hold is physical, and is now a test of its own: **a man at an address can be reached; a man between two addresses is nowhere.**
- **Fixing the crew list uncovered an older bug it had been hiding.** *Recruit Leo Carver* was refused only because he was already in the crew — so the moment death removed him, the game offered to hire him again. It also named him whatever had happened to him, in a city where nobody holds a job for ever. The button now names **whoever actually drives** and is an action about that person, so the same question is asked of it: with Leo dead and the seat empty it reads *Leo Carver is dead*; by the next morning `FillRoles` has given the job to somebody else and it reads *Recruit Ugo Lenz*.
- Verified in a browser on a new `bereaved` fixture, both halves: the crew list empty, the recruit action refused with *Leo Carver is dead* and no clickable button anywhere on the page; then a day later, **Recruit Ugo Lenz**, available. No console errors.
- **Balance unchanged, and checked rather than assumed**: 100 runs identical to the previous baseline, and the defiant profile held at 600 runs — **348 deaths, exactly as before**. Crew deaths are rare enough in simulated play that removing a bonus the player should never have had does not move the numbers.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures.

## The jobs where nobody is named
- The reach question added last iteration is asked of every action that names a person in its id. **The delegated jobs name nobody.** `rob:crew`, `mug:crew` and `sabotage:crew` are aimed at a *place*; the man who actually does them appears nowhere in the id, so the subject parser read "crew", found no such person, and the sweep skipped all three.
- The result, found by probing rather than by reading: **a player could send somebody out of a police cell to rob a business while the police were still holding him**, and go on paying his wages for the privilege. The same gap let a man already crossing the city on an errand of his own be sent somewhere else.
- All three pass through one readiness function, `DelegateReadiness`, and it now asks `OutOfReach` about the man it is going to send. One place, because that is where all of them already meet.
- Verified in a browser on a new `inside` fixture: with Leo Carver in a cell, all five crew actions are refused with *Leo Carver is being held at Ward Street Station*, and no clickable button for any of them appears on the page. **The refusal is not a matter of the interface hiding a button** — posting `mug:crew` straight to `/api/action`, past the interface entirely, is refused by the core with the same sentence. No console errors.
- **Measured, and nothing moved**: 100 runs identical across all four profiles, and the reckless profile — the one that actually uses delegated crime — checked at 600 runs with the change and without it: **469 deaths either way**. A man of the player's is rarely in a cell in simulated play, so closing a hole they should never have had costs them nothing measurable.
- **Noticed and deliberately left alone.** Every command failure answers HTTP **409**, the same status as a genuine stale-revision conflict, so a client cannot tell "refresh, the city moved" from "you cannot do that". Nothing the player sees is wrong — the interface drops the pending command and shows the sentence either way — and changing the status is an API change with no player-visible gain, so it is recorded here rather than done quietly.
- **A second candidate examined and not claimed.** A crew member who joined a family would stay on the player's books, and the two lists would disagree. Nothing in the city currently puts a player's man into a faction, so this is a shape the code allows rather than a bug the game has; it is written down here instead of being fixed on a guess.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures.

## The books counted money that was in the ground
- Continuing the probe: put a person in an unusual state and read what every screen then claims. This tick took the two screens that had never been audited.
- **Two candidates examined and cleared, which is worth as much as a fix.** The Families screen reports no leader and no headcount for either family in a fresh game — that looked wrong until the source said otherwise: both are gated on `Intelligence`, and *"how many people answer to them is a thing you need somebody inside to count"*. It is a designed secret, not a bug. And killing a family's head does not leave the screen naming a corpse: succession renames the organization in the same instant, verified with the player's network raised so the name is actually visible.
- **The Ledger was wrong.** Lend somebody money, have the city shoot them, and the books went on counting the debt as an asset: *OUT ON THE STREET — $374, $504 due back*, owed by a man who was dead. `LoanDay` has always known what to do about it — *"$249 went out and whoever was carrying it is not carrying anything now"* — but it runs at **midnight**, so between the killing and the next morning the ledger claimed money that would never arrive. That window is exactly when a player looks at their books.
- `WriteOff` closes a dead debtor's loan from `Kill`, at the moment it becomes true, beside the crew list which is already emptied there. `LoanDay`'s branch stays as a safety net for saves written before this, rather than being deleted on the assumption that nothing can reach it.
- Verified in a browser on two isolated fixtures, one lent and one lent-then-killed: alive, the Ledger reads **OUT ON THE STREET $374 · $504 due back**; dead, the line is gone, `owed` and `lent` are **0**, the book is empty, and the history carries *Nothing to collect*. The screen was on-screen and measured in both. No console errors.
- **The timing claim is the unit test's**, not the browser's: a fixture cannot show a transition it performs before the page loads. The test holds that the books are correct in the same call that kills the man, and the browser shows that the corrected core renders correctly. Said plainly here rather than implied.
- Measured: 100 runs unchanged across all four profiles — the money was never collectable, so only the moment of the write-off moved.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures.

## A guide that told a man in a cell to go and buy premises
- The Guide was rebuilt once as *the game reporting on itself* — every milestone answered by the same readiness function that answers the button — with the promise, written into its own source, that *"it cannot tell you something the rules do not"*. Probing it with the player in an unusual state showed the promise was not kept.
- **Locked in a cell at Ward Street Station**, the Guide said *You can do this now* against carrying envelopes at Saint Agnes, buying Mara a coffee, and buying premises — while the action list two clicks away correctly offered exactly three things: sit it out, pay a lawyer, or name somebody. **Dead, it went on offering the first two.**
- The cause is worth stating precisely, because the structure was right and the question was wrong: a readiness function answers about a **rule** — whether a still can be built, whether anybody would take a loan — and knows nothing about whether the player is alive or standing in a cell. The action list asks that first and returns early; the Guide never asked it at all.
- It asks now, before any rule: dead closes everything with *This life is over*; held closes everything with *You are being held at Ward Street Station, 5 days to go*, counted from the same `DaysLeft` the cell's own buttons use.
- **Verified in a browser, both directions.** In the cell: every step closed, nothing on the page reading *You can do this now*, and the reason rendered under the first milestone. Then *Do the 5 days* was actually pressed — and afterwards four milestones are open again, no step mentions Ward Street, and the page says *You can do this now* once more. The screen was measured on-screen at both ends. No console errors.
- **A difference that was the rules working, not a bug.** The first version of the test released the player and expected the same four things open as before, and got three: released *at the station*, a man cannot lend money to somebody standing at the bar. The test now puts him back where he was and says why in a comment, rather than the rule being loosened to make a test pass.
- Measured: 100 runs unchanged across all four profiles — the Guide is a report, and reports do not play the game.
- `mise run verify` passes. `cmd/apicheck`: no invariant failures.

## The top bar had nothing to say about being in a cell
- The top bar is the one thing always on screen and the place a player looks to know their situation. Probed with the player locked up, it showed **six figures identical to a free man's**: cash, respect, presence, health, attention, the city. The most important fact about the next five days appeared nowhere on it.
- Worse, one of those figures was actively wrong. *Cash on hand — what you can spend today* is not what a man in a cell can spend, and nothing on the bar suggested otherwise.
- The bar now leads with **Held · 5 days · Ward Street Station**, marked as a warning, with its own icon rather than the city skyline every unrecognised id falls back to.
- **The meaning was measured before it was written.** Rather than assert that the city keeps going, five days were served in a test with two premises running: cash **5,860 → 5,740**, a net cost of $120, and respect **25 → 45**, the 20 the cell's own button promises. So the tooltip says what was observed: *your businesses go on earning, your people go on being paid, and the day costs what it costs whether you are there or not.*
- **A guard for a class of silent failure.** Every top-bar figure is drawn by looking its id up in the icon table, and an id with no entry falls back to the city skyline — so a new figure silently wears the wrong picture and nothing fails. A test now walks the dashboard, with the player confined so the new stat is present, and fails if any id has no path of its own.
- **Verified in a browser, and the first measurement was wrong.** Every stat reported off-screen, including cash — because the page was still scrolled down from the arrest scene, not because anything was hidden. Scrolled to the top, the figure measures 98×57 at y=14, renders in the warning colour, and carries the meaning as its tooltip; a zoomed capture shows the bars icon and *5 days / HELD / Ward Street Station* first in the row. Recorded because a check that fails for the wrong reason is worth as much as one that passes for the right one.
- Measured: 100 runs unchanged across all four profiles, and `cmd/apicheck` clean — the dashboard is a report.
- `mise run verify` passes.

## A wrecked business that said it was trading at everything it could
- Probing premises through extreme states — no staff, no supply, no income, lost to a rival, the player at 100 attention — turned up one that reads wrong in the ordinary course of play, not just at the edges.
- **The figure was measured against the wrong things.** *Trading at 100% of what it could* counted staffing, supply and custom, and left out the **condition of the building**, which the clock uses to scale every dollar the place earns. Measured over a day, one laundry, same staff and supply: **condition 100 earns 305, condition 80 earns 237, condition 60 earns 169, condition 30 earns 67** — and at 80 and 60 the game said it was trading at everything it could.
- `Trading(id)` is the share of what a place could earn, condition included, and the note and the sidebar's *Working at* both read it. A test measures gross earnings at four conditions and fails if the claim and the takings differ by more than six points.
- **A measurement that was wrong before the code was.** The first version compared *net* cash over a day, so the fixed rent and wages — the same whatever state the building is in — made a wrecked laundry look worse than it traded, and the 30% case failed by eight points. Comparing what the place actually earned rather than what was left over fixed the test, not the formula.
- **The panel was arguing with itself.** *Working at 80%* sat beside *Trade — earns 100% of what it could*: the same sentence about a different number, the first about everything and the second about how many regulars a place keeps. The Trade row now says *regulars are worth 100% of ordinary takings*, and a guard fails if the phrase "of what it could" reappears in the property panel — that claim belongs to the one figure the core computes.
- **Two things probed and cleared.** The Herald keeps its archive correctly across a death: in a new life the current edition is empty and the previous life's issue is still there, marked not-current and not-yours. And premises at zero condition, zero staff, zero supply and 100 attention each say the right thing in the right order of urgency.
- Verified in a browser on a new `worn` fixture (a laundry at 80%): *Trading at 80% of what it could* on the street, **Working at 80%** in the panel measured on screen at y=575, and zero occurrences of the contradicting phrase. No console errors.
- Measured: 100 runs unchanged across all four profiles, `cmd/apicheck` clean — the arithmetic was already right; only the number the game showed was wrong.
- **And the fix was wrong in the live game before it was right.** Deployed to the campaign on 8791, three of the player's businesses read *Trading at 102% of what it could* — a percentage of a ceiling, above the ceiling. Custom above the middle lifts a place past an ordinary day, so "what it could" was never a maximum; the sentence now reads **of an ordinary day**, which is what the number has always measured. Caught by restarting the live game and reading it, not by any test.

## Five days reported as a hundred and twenty hours
- Probes that **cleared**, recorded because clearing is worth as much as fixing: the People screen's *$14 a day* per man matches the books exactly (3 men, `Your own people = 42`); the Ledger's income claim is accurate to within two dollars a day at both 100% and 60% condition, which is what showed that the earlier fault was in `Capacity`, not in the books; and *looking for somewhere else to be* is honoured — in 40 campaigns where wages stopped, somebody went restless in 40 and somebody actually left in 40.
- What was wrong was small and constant. The result band reported elapsed time in hours past two hours, so **sitting out a sentence the game itself calls "Do the 5 days" was reported as "120 hours"**. True, and it tells the player nothing.
- Past a day, days are the unit. Verified in a browser by actually serving the sentence: the button said *Do the 5 days* and the band read **5 days**, measured on screen at y=86; an ordinary thirty-minute errand taken straight afterwards still read **30 min**.
- The Go test beside it is a text guard and is labelled as one in its own comment: it can only say the formatter knows what 1440 minutes is. That the band reads "5 days" is the browser's claim, not the test's.

## The city, drawn as a city (slice one: ground, placement, blocks)
- The user asked directly whether the isometric city view had been started, and the honest answer was **no**. What stood in for it was `CityStreet.tsx`: twelve painted cards in three district rows, with walkers as portrait chips sliding in a straight line between two cards. `CityScene.tsx` was 23 lines of orphaned dead code imported by nothing.
- This is the first slice of the real thing. `src/iso.ts` holds the projection and nothing else that matters: a 2:1 isometric transform, a plan-to-tile mapping from **the coordinates the core already keeps for every address**, a depth function, and a block made of small boxes. The rule written into its head: *it decides shape, never truth.* Where a building stands, who is in it and what is happening all still come from the core.
- `src/CityIso.tsx` draws all twelve addresses, depth-sorted back to front, each on its own plot, named above the roof, selectable, with the player's address and the spotlit one marked. The card view stays reachable and the choice is remembered.
- **Two overlaps found by measuring, not by looking.** The first spacing put The Mariner and Saint Agnes through the Bellwether Herald. Buildings standing on each other is the one fault no amount of art fixes, so the footprints moved into a single table and `TestNoTwoBuildingsStandOnTheSameGround` now reads that table and `CELL` straight out of the renderer, checks them against the city's real coordinates, and **fails loudly if it cannot parse them** — a guard that quietly finds nothing to check is worse than none. Confirmed it bites by restoring the old spacing: it named both overlaps exactly.
- **A verification that was wrong before the code was.** The browser check for back-to-front drawing reported the order broken; it was comparing plan `x+y` while the renderer sorts on block centres including footprint. Re-measured with the formula the renderer actually uses, the order is correct.
- Verified in a browser: twelve blocks drawn for twelve addresses, none missing, no overlapping footprints, drawn back to front, clicking Pier 14 moved the panel from Saint Agnes to Pier 14, and the player's address is marked. No console errors.
- **What this slice deliberately is not.** The blocks are flat-shaded polygons, not the high-fidelity isometric art the goal calls for. That is the next slice, and the count is worth writing down: **five addresses have true isometric cut-outs** (`buildings.json`: Saint Agnes, The Monarch, Bluebird Laundry, The Mariner, Mercer Exchange) and **seven have street-elevation photographs**, which are the wrong projection for a city seen from above and cannot be dropped onto this grid. Those seven need generating before the city can be painted rather than blocked out.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged — the view reads the world and never writes it.

## Choosing a renderer, and finding one already in the box
- The user asked whether an existing JS framework should handle the isometric work. It should, and **the right one was already a dependency nobody was using**: `pixi.js` 8.20.1 sat in `package.json` while `grep` found zero occurrences of it in the built bundle, because the only file importing it — `src/CityScene.tsx`, 23 lines — was imported by nothing.
- Worth recording what that dead file did, because it is the trap to avoid: it rasterised the whole SVG map to **one texture** and laid invisible hit rectangles over it. That is PixiJS used as an image blitter — no scene graph, no per-building depth, no camera. Deleted.
- **What was surveyed.** Phaser 4 (April 2026) has isometric tilemaps built in (`Tilemap.Orientation.ISOMETRIC`) and renders on PixiJS underneath; Excalibur and melonJS cover similar ground; the isometric-specific libraries — Isomer, obelisk.js, the Phaser 3 isometric plugin — are small and long unmaintained.
- **The decision, and why.** PixiJS v8 with `pixi-viewport` 6.0.3, and no game framework above it.
  - The isometric part is **ten lines of arithmetic** that already exist in `src/iso.ts` and are covered by the overlap test. Importing an engine to get a projection would be paying a great deal for something already written.
  - Phaser's isometric support is **tilemap-shaped**. This city is twelve buildings placed at the coordinates the core keeps for them, not a grid of tiles, so most of what Phaser offers would go unused while it took ownership of the loop, the canvas and the input.
  - What was actually missing from the SVG version is exactly Pixi's job: a real camera, hundreds of sprites at once, and filters for explosions and light.
  - `@pixi/react` was considered and **not** used. The scene is state-driven rather than component-shaped, and the documented friction between that renderer and `pixi-viewport` buys nothing here. React keeps every piece of UI chrome; Pixi owns one canvas.
- The camera is real: drag, wheel, pinch, deceleration, and a zoom clamped between .45 and 2.6. It frames the whole city on first draw **from the bounds of what was actually drawn**, not from a guessed rectangle, and it re-frames only when the window changes shape — a player who has zoomed in on the docks stays there when an hour passes, and a test fails if that is ever wired to an ordinary update.
- The city stays reachable without a mouse or WebGL: the same twelve addresses as focusable buttons, tested for.
- Verified in a browser by driving it: a wheel event zoomed in, a pointer drag moved the city, and the screenshot afterwards shows The Monarch and Mercer Exchange filling the pane where a moment earlier the whole city fitted. Twelve reader buttons, 1952×1200 backing store, no console errors.
- **Still placeholders.** The buildings remain flat-shaded boxes. The renderer swap changes what is possible, not what it looks like; the art is the next slice.

## The city, painted
- Twelve isometric cut-outs, generated in one pass so they belong to each other: the same angle, the same palette, the same overcast evening light, each a complete building including its roof, standing on its own pavement with the background knocked out to transparency.
- **Why one pass rather than seven.** Five addresses already had hand-made cut-outs and seven had street elevations, so the obvious move was to paint the missing seven. That would have produced a collage — a city assembled from pictures drawn at different angles reads as a collage however good each piece is. All twelve were repainted together instead, and the reasoning is in the tool's own head.
- `tools/isometric.py` paints them, `tools/paint.py` holds the one way to call the model, and `tools/manifest.py` rebuilds the manifest from whatever is on disk. The shared caller exists because **the first version of the isometric tool invented a `--path` flag mflux does not have and died on the first building** — the working invocation was already in `tools/exteriors.py`, so now there is one of it rather than two that can drift.
- **A quality gate, because one of the twelve came out wrong.** The Bluebird Laundry arrived with a pale floor painted across the whole frame rather than an isolated building, so the flood fill from the corners stopped at it and the map showed a black rectangle stuck to the laundry. `tools/inspect_iso.py` calls a cut-out suspect when its bottom edge is opaque from side to side or its sides are opaque down half their length — *a building has a footprint, a floor has the frame*. It named the laundry and cleared the other eleven, which is exactly what was seen on screen. `ISO_RETRY` nudges one building's seed without disturbing any other; one retry fixed it, and the gate then passed twelve out of twelve.
- The renderer draws the picture where there is one and the blocked-out solid where there is not, so an address added tomorrow appears on the map rather than leaving a hole. Cut-outs are scaled to **the ground they stand on** — the projected width of their own plot — never to their own pixel size, and their feet are placed on the near corner of the plot so a building stands on its ground rather than floating over the middle of it. A painted building brings its own pavement, so the drawn plot is only kept under a blocked-out one or when there is something to say about the ground.
- A test holds that every address in `core.Locations` has a cut-out in the manifest **and that the file is actually there**, because a missing one degrades to a flat solid — correct behaviour that looks like a bug beside eleven painted buildings.
- Verified in a browser: all twelve painted and placed, the laundry artifact gone after the retry, no console errors. 7.8MB of art for the whole city.
- **Still to come, and not pretended otherwise:** nobody is in the streets yet. The people, their journeys and the events all exist in the core and none of them are drawn in this view.

## The people, in the city
- Forty-eight people drawn where the core says they are: on the pavement in front of the address they are standing in, spread along it and into a second row so a crowded building reads as a crowd, and the two who are out walking placed **between the two doors at the fraction of the way they have actually got**. Where they are is entirely the core's answer — `PeopleHere` for the standing, `OnTheStreet` for the walking.
- **A deliberate choice about characters, recorded because it constrains everything after it.** These are not portraits. At the scale a whole city is drawn at a face is four pixels of mud, so a person is a coat, a collar, a hat and a shadow, in two colours, with their own people marked. Fully animated per-character sprites are where AI generation stops being reliable; this stays inside what the style can carry, and it reads as a man standing on a pavement — which is the whole requirement. The first version had no hat and read as a pin.
- Walkers are drawn **after** every building, so they pass in front of the city rather than through the middle of it, and they carry their first name so the street can be read at a glance.
- Verified in a browser at two zooms: from above, figures on pavements across the city with Otto and Hedda named between buildings; zoomed in on The Monarch, men in hats and coats standing under the marquee. 48 people, 2 walkers at 40% and 50% of their journeys, matching the API exactly. No console errors.
- **One flaw worth naming rather than hiding.** The generated buildings carry invented signage — The Monarch's marquee reads *CIICATCE*. Diffusion models cannot spell, and the prompt asking for no lettering did not stop it. It is legible as a sign and nonsense as a word. Fixable by regenerating with the sign masked out or by painting signage separately; not fixed today.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## The moment happens in the city, not over it
- The theatre used to take the camera to a building by dimming a row of cards. Now the camera actually goes there: when a cue arrives, the viewport **animates to the address the core named**, zooms in if it was further out than 1.25, and plays the moment over that building. When it ends the camera is **given back** to exactly where the player had it, position and zoom — a test fails if that is ever dropped, because being left staring at a rooftop across town is worse than not moving at all.
- Each kind of moment is light and shape over the real painted building, never a drawing of one: a fireball and a shockwave with debris thrown out and falling for an explosion, muzzle flashes for a killing or a gunfight, a police lamp sweeping the front for a raid or an arrest, and a hard ring that opens once for anything else. A test reads the kinds straight out of `core/witness.go` and fails if the city draws nothing for one — **it caught `robbery` reaching only the generic fallback**, which is now named deliberately rather than by accident.
- The theatre's painted plate is dropped when the city view is the stage. A stock picture of a police station in front of the actual police station is one picture too many.
- **Two bugs found while verifying, both real.** *Replaying a moment forced the card view* — a leftover from when that was the only view — so replay silently threw the player out of the city; it now keeps whichever city they were looking at. And the city was being **rebuilt on every frame of a moment**: twelve buildings, sixty times a second, to animate a fireball drawn in a different layer. The redraw no longer depends on the moment's clock.
- **How this was actually verified, including what did not work.** Four screenshots in a row showed no fireball, and the temptation was to conclude the effect was broken. It was not: instrumenting the renderer showed the graphic present, visible, alpha 1, 493×493 at canvas (488, 194), with `t` animating 0.21 → 0.71. Reading the rendered frame back through Pixi's extractor gave **RGBA(49, 23, 8, 67)** at the shockwave — the exact alpha the arithmetic predicts for that radius at that instant. The screenshots were simply landing after a 4.4-second moment had finished; a round trip to the browser is slower than the thing being photographed. Holding a moment permanently in flight produced the picture: a fireball engulfing The Monarch with its shockwave ring around it, over the painted building, in the city view.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged — the view reads the world and never writes it.

## The city makes a noise
- The last thing on the stated goal that had not been started. An explosion, a shot, a police lamp turning over at the kerb — played once when a moment starts, over the building it happened at.
- **Synthesised, not sampled, and the reasoning matters.** A gunshot is a noise burst with a four-millisecond attack and a short tail; an explosion is low-passed noise sweeping 900Hz to 90Hz over a second with a sine thump under it; a police lamp is two square tones alternating. Nothing is downloaded, nothing is licensed, nothing waits on a model that cannot be relied on for audio, and the whole file is a few hundred bytes rather than a few megabytes. It is also the honest match for the art: a 1950s city heard through a newsreel rather than a field recording.
- Sound has its own control in Settings, beside Scenes and Voices — required by this repository's own guard that **every stored preference the interface reads must be writable from somewhere a player can click**. A browser that refuses local storage defaults to on rather than to silence, and a test holds that too.
- **Verified by counting, because I cannot hear it.** The page's `AudioContext` was wrapped and every source it started was counted while moments were driven through the interface. Each kind produced exactly the graph its code specifies: an **arrest 4 oscillators** (the four siren tones), an **explosion 1 noise buffer and 1 oscillator** (the blast and the thump under it), a **killing 2 buffers** (two shots), a **gunfight 5 buffers** (five shots) — and with the setting turned off, **0 and 0**. One `AudioContext` for the session, built on the first sound because a browser will not let a page make a noise before somebody has touched it.
- The five Settings rows were measured on screen: Scenes, Sound, Voices, The storyteller, This life. No console errors.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## Streets, pavements, lamps and cars — and a grid, after getting it wrong once
- **The first attempt at roads was wrong and is worth recording.** It joined every building to its three nearest neighbours with L-shaped routed paths, checking each elbow against the footprints. The geometry was sound and the result was a web: streets at every angle of approach, laid over one another, meeting nothing squarely. A city is a grid, and routed paths between scattered points will never look like one.
- So the ground is a grid and the buildings are placed on it. The addresses still decide **where they belong** — an address's own coordinates choose its block, so the docks stay west and Cypress House stays north-east — but the blocks are regular, the carriageways run the full width and height of the city, and every junction is square. Bellwether's twelve addresses fall into a **six by four grid with no two wanting the same block**, which is why this needed nothing moved by hand.
- What is on the ground now: one slab under the whole city so nothing floats; carriageways cut across it with broken centre lines; a pavement inside every block with a kerb line around it and a darker plot where the building stands; a cast-iron lamp standard at every block corner with an arm, a lantern and the pool it throws; and cars parked at the kerb.
- **The cars were drawn twice.** The first were hand-drawn in screen space and read as smears at any zoom, because nothing about them agreed with the angle everything else is at. They are built from small boxes in tile space now and go through the same `faces` projection as the buildings — which meant widening that function to take geometry without colours, since a car's colours are numbers and a building's are strings.
- **The clock is stopped between actions, so the cars are parked and nobody idles.** A car moving while time is not would be the view inventing something the simulation has not spent, which is the same rule the walkers already follow.
- People walk the streets now rather than crossing the map in a straight line: out of their block to the nearest carriageway, along it, round the corner and in. **Measured in the browser at twenty points across each journey**: on the street for everything between 15% and 85%, and off it only for the first and last tenth — which is a person crossing their own forecourt to reach the road, and is correct.
- **A test replaced rather than kept.** The old guard checked the addresses' raw coordinates against a footprint table for overlaps. That was right when buildings stood wherever their coordinates put them; it stopped meaning anything the moment the grid decided placement, and a test that cannot fail is worse than none. It now checks the thing that actually matters — that every address gets its own block — by doing the same binning the renderer does.
- Verified in a browser at two zooms: from above, a clean grid of blocks with aligned streets and lit corners; zoomed in, a parked car reading as a car, the lamp standards, the kerbs, and people on the pavement outside The Monarch. No console errors.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## Night, and the blocks between the addresses
- **The light comes from the clock.** The city was lit identically at noon and at three in the morning, which made it a diagram. `nightness(minute)` runs 0 at midday to 1 in the small hours with the turn happening over the hour dusk actually takes, and everything reads that one number: the ground goes from grey stone to blue-black, the pavements and kerbs with it, the lamps come up as the light goes down, and the pools they throw appear with them. A lamp burning at noon is the surest sign nothing is looking at the clock.
- Painted buildings spill warm light onto their own pavement after dusk, and haze lies over the far side of the grid — the next street is clear and the one past it is a suggestion. Both scale with the hour.
- **The empty blocks are built on.** Twelve addresses in a six by four grid left twelve holes, and a city with holes reads as a scatter of models.
- **The first fillers were wrong and are worth recording.** They were flat grey boxes drawn from the same solids as the placeholder blocks. Zoomed in, one filled the frame as a featureless slab with three orange squares floating on it — which is *worse* than a hole, because a placeholder reads as a mistake rather than as distance. Six ordinary buildings were generated instead — a tenement, a warehouse, a terrace of shopfronts, an office block, a small works and a corner block — in the same pass and the same light as the twelve, and all six passed the cut-out check first time.
- They are held back deliberately: tinted cooler and darker, scaled smaller, never named and never clickable, so an address the player can walk into always reads first. The drawn solid survives as the fallback for a block whose picture has not loaded, and it grew rows of lit windows down both faces so that even the fallback is a building rather than a slab.
- A test holds that the filler set exists, that every file the manifest names is on disk, and that there are **at least four** — a row of blocks showing the same building repeated is its own kind of placeholder.
- Verified in a browser at two zooms: from above, every block built on with the twelve still reading first; zoomed in, painted warehouses and terraces with lit windows, a car at the kerb, lamps lighting the street, and everything inside its own block with the carriageways clear. No console errors.
- **A bug caught while wiring the haze.** Hovering a building set its opacity to 1 and leaving it set it back to *full* brightness rather than to its hazed value, which would have left every building the mouse crossed permanently nearer than the rest. Hover now lifts from the resting value and returns to it.
- 11MB of art for the whole city. `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## What a pavement carries, and people who look like people
- Hydrants, mailboxes, bins, benches and telegraph poles with wires strung block to block. Each is a few boxes in tile space through the same projection as the buildings and the cars — anything drawn in screen space stops agreeing with the angle the moment the camera moves, which is the mistake the first cars made.
- **Placement is provable rather than eyeballed.** Props are laid on the pavement ring only, inset from the kerb, so nothing can stand in a carriageway or under a building. Checked in the browser by recomputing all **63 props** against their own blocks: **none in a road, none under a building**. A text guard holds that the placement still comes from the block's pavement ring and is still inset by the pavement's width — if that stops being true, props can wander.
- **Two things that looked like bugs and were.** The first hydrant read as a red cube hanging in the street: it had no contact shadow, so it floated. Every prop has one now. And it was pillar-box red in a city of muted olive and umber — the only saturated thing on the street, which made it read as a mistake rather than as a hydrant. It is oxide red now. Props were also a size too large and competing with the people; all of them came down.
- **People stand like people.** They were spaced evenly along the pavement, which read as a fence. They gather in twos and threes now with gaps between the knots, at slightly different depths, and every second one is turned the other way so a knot looks like a conversation rather than a queue. Walkers face their own direction of travel, taken from where they were a moment earlier on their own path.
- **The player is on the map.** They are not in any room's list — the core keeps the protagonist apart from the city's own people — so the one figure that mattered most was the only one not drawn. They stand at whatever address they are at, inside a gold ring on the pavement, findable in a crowd without reading a name.
- Verified in a browser zoomed in: a telegraph pole with its crossarms, a hydrant at the kerb with a shadow under it, people in knots outside The Monarch, and the gold ring under the player at Saint Agnes. No console errors.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## Terraces, because a city is not an office park
- **The user was right and the model was wrong.** One building centred in its block with pavement on all four sides is a suburban office park. A city block is a **terrace**: buildings shoulder to shoulder along the frontage, sharing party walls, with the yards behind.
- A block is now a row of slots — three across the near frontage and three across the far one. The address that belongs to the block takes the middle front slot, so the building the player came to see always faces the street and is never hidden; ordinary buildings take every other slot. Sixty-odd buildings now stand where twelve stood in fields.
- Sprites are scaled a little wider than their slot, which is what closes the party walls: without the overshoot each building sat in its own stripe of pavement and the terrace read as a shelf of models.
- **Six buildings across sixty slots reads as wallpaper**, so the set is stepped through by position rather than picked at random — the same picture cannot land next door to itself — and neighbours vary slightly in tone and height, so a terrace looks like buildings put up at different times rather than one building stamped along the street.
- **A test replaced, and a function deleted.** The guard asserted that fillers went on the blocks the addresses left empty, which was true of the old model and meaningless under the new one; it now checks that blocks are laid out as terraces and that nothing builds on the slot an address takes. `fillers()` — which found whole empty blocks — was superseded and removed rather than left imported and unused.
- Verified in a browser at two zooms: from above, a dense city of terraced blocks with varied rooflines; zoomed in, The Monarch with neighbours butted against it on both sides. No console errors.
- **Honest about what is still wrong.** The party walls meet unevenly, because every sprite is a detached building carrying its own pavement — the art was generated as twelve models and six fillers, not as terrace rows meant to tile. The buildings still carry garbled signage. Both want better art rather than better placement: FLUX.1-schnell at four steps was the right choice for getting a coherent city quickly and is the wrong one for finishing it.
- `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.

## Four attempts at better art, two that worked
- The plan for this iteration was to fix the two things placement cannot fix: party walls that meet unevenly, and buildings carrying garbled signage. Most of it failed, and the failures are the useful part of the record.
- **Terrace rows: failed, and the output was kept for something else.** Four sprites were asked for as *"a terrace of four joined buildings in a straight row, flush square ends, shared party walls"*. All four came back as buildings bent around a **corner** — schnell at four steps does not follow that kind of structural instruction. They are good buildings and they are not rows, so they are used as **corner blocks at the ends of terraces**, which is where a corner building belongs, and named `row-*` with a comment saying exactly what happened rather than being passed off as what was asked for.
- **Signage suppression by prompt: failed.** The Monarch and the Herald were regenerated with *"plain unlettered facade, blank walls, no signs, no billboards, no marquee, no writing anywhere"* appended. The replacement Monarch came back reading **NHTIUUR**. Diffusion cannot spell and cannot be told not to write.
- **Softening the letterforms: works, but only when aimed.** `tools/unletter.py` finds small, very saturated, very bright, high-contrast clusters — what painted lettering is — and blurs those regions only. On the Monarch it took the marquee down to an indistinct warm glow and left the cornice, the balcony, the windows and the brickwork untouched. **Run across the whole directory it does real damage**: it took the mullions out of the bar's windows and dulled the market's colonnade, because a bright warm window is not far from a bright warm sign. That pass was reverted and the tool now says in its own head that it is not a batch step.
- **A regeneration that had to be thrown away.** The replacement Herald came back with a bite taken out of its facade: a dark bay matched the background closely enough that the knock-out flood-fill walked into it. Reverted to the committed version.
- **A guard attempted, measured, and not shipped.** A hole detector was written for that failure — transparent pixels the background cannot reach — and it passes the holed Herald, because the bite is connected to the outside edge rather than enclosed. A raggedness metric was tried instead and **measured across all 22 buildings**: the bad Herald scores 18.9 while the docks, which is fine, scores 33.7, and `fill-corner` 35.5. Fire escapes and chimneys raise perimeter legitimately. The metric does not separate good from bad, so it was not shipped — a check that cannot fail on the case it was written for is worse than none.
- Where that leaves the signage: the one sign the user actually named is gone. Smaller lettering survives on both buildings, less saturated than the filter's threshold.
- Verified in a browser: a dense city of terraced blocks with corner buildings at the ends of rows, the Monarch's marquee now an indistinct lit panel. No console errors. `mise run verify` passes, `cmd/apicheck` clean, 100 runs unchanged.
- **What this iteration says about the tools.** Every one of these failures is the same failure: FLUX.1-schnell at four steps is a fast, low-fidelity tier that cannot follow structural instructions or refrain from writing. The city was built with it because it produced twelve coherent buildings in eight minutes. Finishing the city wants FLUX-dev at twenty or thirty steps — mflux supports it as a built-in model — which is a **9 to 17GB download onto a disk that is 89% full**, and that is the user's call rather than a decision to take quietly in a loop.

## Road markings, and a claim that had to be withdrawn
- Crossings and stop lines at every junction, laid from the block grid so they line up with the kerbs by construction. Worn paint rather than fresh — it has been on the road a while.
- **A correction.** The previous iteration reported that the garbled sign the user actually named was gone. That was true of the Monarch's *marquee* and false of the city: the six filler buildings and four corner blocks carry their own invented signage — *AIUT HEFTER*, *TUALE*, *Tut Calle* — and those repeat across sixty-odd slots, so the city had **more** wrong text on it than before the fillers existed, not less. The aimed filter has now been run over all ten of those, which takes the worst of it out and leaves the brickwork and lit windows intact.
- **And a second correction, because the first fix was still not the whole truth.** The Monarch's *rooftop billboard* still reads CIICATC. The filter finds saturated warm lettering; that sign is painted pale on dark green and falls under every threshold that keeps it from eating windows. Two claims about this being fixed have now been wrong, so: **the signage problem is not solved.** What is true is that the loudest examples are softened.
- **A false alarm chased properly.** The Monarch appeared to still show CIICATCE in the browser after the art was fixed, which looked like a caching bug — regenerated art never reaching a returning player would be a real fault. Fetched the file from the running server and hashed it: **790,812 bytes, sha256 b0343ff4d62fc942**, byte-identical to disk and to `dist`. No caching bug. What was on screen was a different building's sign.
- The honest read on all of this: filtering is the wrong tool for a model that cannot spell. It removes legible wrongness where the paint happens to be saturated, and cannot touch anything else without softening the windows that make these buildings worth looking at. This wants regeneration at higher fidelity, which is the FLUX-dev download already recorded as the user's call.
- Verified in a browser: crossings and stop lines visible at the junctions, terraced blocks, no console errors. `mise run verify` passes, `cmd/apicheck` clean.

## The buildings were standing inside each other
- The user said the blocks looked terrible and that buildings were overlapping, and asked whether that had been worked on. **It had not.** Props were verified against the roads and the buildings — 63 of them, none in a road, none under a building — and every address was verified to hold its own block, but **building-to-building overlap was never checked once**. Worse, it was caused deliberately: a `* 1.16` overshoot was added to the sprite scaling to close the party walls, and an entry in this file two sections above describes it approvingly. Multiplying a sprite past its own plot *is* overlapping; it was eyeballed at one zoom and called good.
- **The overshoot is gone**, in both places it appeared — the filler terraces and the twelve painted addresses. A building is now scaled to exactly the ground it stands on. A hairline gap between two buildings reads as two buildings; an overlap reads as broken.
- **The height distortion is gone too.** `art.scale.y *= .93 + warmth * .16` varied each neighbour's height for "variety" and did it by breaking the aspect ratio of art that was painted correctly — a second artefact stacked on the first. Tone still varies between neighbours; proportions no longer do.
- **Three buildings a frontage became two.** Six buildings crammed into a 3.6-tile block gave each one a slot narrower than the art wanted, which is why they were reaching into each other in the first place. `terrace()` now cuts two slots a side, four to a block, and the address takes the near-left frontage.
- **Measured, not eyeballed.** `tests/city-blocks.test.mjs` runs the real `iso.ts` over a grid larger than the city and asserts that no two building footprints intersect, that none leaves its pavement ring for the road, that the front and back rows are two rows, and that the address's slot exists and faces the street. `iso.ts` is pure geometry with no Pixi imports, so this is the module the city is actually drawn from rather than a copy of it.
- **A guard on the cause, not just the symptom.** The footprint test passes on the *old* code too — the slots never overlapped, the pictures drawn on them did. So `TestNoBuildingIsDrawnWiderThanItsGround` reads `CityIso.tsx` and fails if a sprite is scaled by anything but its own plot, or if `scale.y` is touched independently again.
- Verified in a browser at three zooms, including hard in on the block behind The Monarch: four buildings to a block meeting at their walls, roofs no longer cutting through neighbouring roofs, nothing hanging over a kerb. No console errors.
- `mise run verify` passes, `npm test` 15/15, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.
- Still true and still unsolved: the garbled signage (CIICATCE is plainly readable at this zoom), which wants regeneration at higher fidelity rather than any amount of placement work.

## Awnings over the shopfronts
- A canvas awning is the cheapest thing that turns a wall with windows in it into a place of business, and it is the one part of a building that belongs to the street rather than the block: it hangs out over the pavement. So it is drawn from the terrace geometry in `awnings()` rather than painted into the art — the generated buildings cannot be asked for one, and if they could it would be stuck to the wall rather than reaching past it.
- Sloped down toward the street so rain runs off, in bands of canvas across the frontage, with a valance hanging off the front lip so it reads as cloth rather than as a shelf, and a shadow on the pavement so it does not float. Four era colourways over cream.
- **It cannot overhang the road by construction**, not by being checked afterwards: its projection is `PAVE * .5`, so the kerb is out of reach. `tests/city-blocks.test.mjs` holds it anyway — an awning stays inside its own shopfront, hangs off the wall it belongs to rather than floating in the middle of a plot, never reaches past the kerb, and no two occupy the same air. 25 awnings across the six-by-four city.
- **Verified they were actually being drawn, rather than assumed.** At night the first version was too dark to find in a screenshot, which is indistinguishable from not rendering. They were filled magenta and photographed: **20 canopies, one per non-address shopfront, each against its building's street wall on the pavement** — matching the 20 the geometry predicts (of 25, five sit on slots the painted addresses take and are not drawn). Then the fill was restored and lightened, because an awning sits under a lamp standard and over a lit window and has no business being as dark as a roof.
- Verified in a browser at two zooms after the fix: striped canopies over the shopfronts either side of Bluebird Laundry, on the pavement, not in the road. No console errors.
- `mise run verify` passes, `npm test` 17/17, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## Dawn, and a light that has a colour as well as a level
- The palettes had only ever been seen at 20:10 and 12:12. **Looked at the two that had not been**, on isolated copies of the QA save with the clock moved: 05:40, 03:10 and 12:12. Dawn was wrong — `nightness` alone makes six in the morning *a weaker three in the morning*, the same blue-black ground a little lighter, when what actually separates them is that the light comes in low and warm for an hour at each end of the day.
- `goldenness(minute)` is that: 0 through the middle of the day and through the small hours, 1 at the two turns. It rides alongside `nightness` rather than replacing it — level and temperature are different things, and the test holds exactly that: 05:30 and 03:10 are both night, and only one of them is gold.
- **Stone takes the low sun; asphalt barely does.** The first attempt warmed ground, pavement and carriageway equally and the whole city went one flat brown with the road grid dissolved into it. The roads now take a ninth of the warmth the stone does, so the grid still reads at dawn. Painted buildings take none — warming the art as a whole washes it out.
- Verified in a browser at three times of day on three isolated saves: **06:15 warm stone with the road grid still dark and legible; 03:10 cold blue-black with the lamp pools; 12:12 unchanged**, which is what `goldenness` returning 0 there guarantees arithmetically and what the screenshot confirms. No console errors at any of them.
- **An observation, not a claim of progress:** midday still reads as overcast rather than as daylight. That is the day palette itself and it is older than this change — `goldenness` is 0 at noon so nothing here touched it. Recorded rather than quietly adjusted.
- QA note: the clock was moved by patching the `minute` field of the state JSON in a **copy** of the fixture save, never the live campaign. Nothing in the core was changed to make the view testable.
- `mise run verify` passes, `npm test` 18/18, `cmd/apicheck` reports no invariant failures.

## Midday that reads as midday
- Recorded last iteration and fixed here: noon looked like an overcast dusk. The day ends of the palette were nearly as dark as the night ends — ground `0x2a2f2c`, pavement `0x4a514c` — the painted buildings were greyed by 30% even at noon, and the depth haze was a third as thick at midday as it is at midnight.
- Only the **day** ends moved: ground, carriageway, pavement, kerb and the plot join all lightened, the buildings' grey wash at noon cut from .30 to .10, and the haze at noon from .35 to .16 of its night depth.
- **Night is arithmetically untouched**, not just eyeballed as unchanged: every one of these is `mix(day, night, dark)` with the night argument the same as before, and both rebalanced pairs sum to the same value at `dark = 1` — `.3 + .28` and `.1 + .48` are both .58, `.35 + .65` and `.16 + .84` are both 1. There is nothing to re-verify at 03:10 because the code takes the identical branch.
- Verified in a browser at 12:12 on an isolated save, at two zooms: pale stone pavements, crossings and stop lines legible against dark asphalt, awnings visible over the shopfronts, blocks reading as separate. No console errors.
- Still true and still art rather than lighting: several cut-outs have their windows painted lit, so a few buildings glow at noon. That wants regeneration, not a tint.

## Smoke off the chimneys, steam off the grates
- A still city is a model of a city, and this one cannot be fixed with animation: the clock is stopped between actions and nothing may move on its own. But a chimney with smoke standing over it is a *still* thing in a photograph, and it is what says the place is occupied. `vents()` places them from the block geometry — smoke from a chimney on one of a block's roofs, or steam from a grate in the pavement outside, never both, and most blocks give off nothing, because a city where every roof smokes is a foundry.
- The test holds where they can be: a chimney stands on a building rather than in the yard behind it, a grate is on the block but under no building, a block gives off one thing at most, and fewer than 70% of blocks give off anything. **11 chimneys and 4 grates** across the six-by-four city.
- **Proved they render rather than assuming it**, the same way the awnings were: filled magenta and counted in a screenshot — 11 tall plumes off roofs and 4 low ones at pavement level, matching the geometry exactly.
- **The first version was wrong and the screenshot is why.** Six evenly spaced ellipses stack into a column of visible grey rings — a drill bit, not smoke. It now takes twenty small overlapping puffs per chimney, each nudged off the centre line, widening faster near the top where a plume is losing its shape. That reads as smoke.
- Verified in a browser at 12:12 and at 03:10: ragged plumes drifting off the rooftops, visible against both the pale daylight ground and the dark. No console errors.
- `mise run verify` passes, `npm test` 19/19, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## A streetcar down the middle avenue
- `rails()` and `sleepers()` derive the track from the same grid the carriageways come from, so it runs down the centre of one avenue for the whole length of the city and crosses every junction square. It cannot end up half on the pavement because it is not placed anywhere — it is computed from `trolleyAvenue(size) * BLOCK`.
- The test states that as a property rather than trusting it: both rails are parallel to the avenue, within half a carriageway of its centre line, spanning the full height of the city, exactly `TROLLEY_GAUGE` apart, with the ties between them and no block anywhere on the track.
- The dashed centre line is suppressed on that one street. A painted centre line and a pair of rails down the same tarmac is two things claiming the middle of the road.
- **Proved it renders and lands where the geometry says**: rails filled magenta and photographed — one straight track down one avenue, corner to corner, between the blocks and through the junctions, partly occluded by the buildings in front of it, which is the depth sort behaving. Then restored to polished-steel grey, a little brighter at night because a rail head that is used is the one thing in a dark street that catches light.
- Verified in a browser at 12:12: the double line with its ties reads as track in the avenue beside Saint Agnes. No console errors.
- `mise run verify` passes, `npm test` 20/20, `cmd/apicheck` reports no invariant failures.

## The room now knows what time it is
- The Interior had not been touched since the city was rebuilt, and the clearest way it had come loose was this: the city outside reads one number off the core's clock for its ground, lamps, window spill and haze, and the room read nothing. Stepping inside at three in the morning put the player in the room they would have found at noon.
- `roomLight(minute)` in `roomart.ts` imports **the same `nightness` and `goldenness` the city uses** — not a second copy of the rule — and returns a wash to lay over the backdrop: a lamp pool that deepens as the night does, a low warm light through a window at the two turns of the day, and a flat darkening that is zero at noon. The backdrop itself, painted or drawn, is untouched underneath.
- The test asserts the agreement rather than the appearance: the room's darkness *equals* `nightness` at the same minute and its warmth equals `goldenness`, noon adds no darkening at all, and dawn is dark and warm at once while noon has no sunset in it.
- Verified in a browser on two isolated saves: Saint Agnes at **03:10** is a dark room with a lamp pool over the table; at **12:12** it is the painted room with barely a wash on it. No console errors.
- **What this does not do**, since the loop note called the Interior "disconnected": the room is still a painted 4:3 backdrop in a different register from the isometric city, and the figures on it are still portraits on cones. Sharing the clock is one real connection, not the whole of one. Saying otherwise would be dressing it up.
- `mise run verify` passes, `npm test` 21/21, `cmd/apicheck` reports no invariant failures.

## Weather, and where a fact is allowed to live
- The city view could have decided for itself that it was raining. It is not allowed to: what the sky is doing is a fact about the world, and a fact invented by a view exists only where somebody is looking — a different sky in the street than in the room, and none at all in the newspaper. So `core/sky.go` owns it and `Public()` publishes it, exactly like the clock.
- **Nothing in the rules turns on it, and that is deliberate rather than unfinished.** The fact is established first and cheaply so that when working a door in the rain is worth a modifier, the modifier has somewhere to live and every screen already agrees about which day it got wet. 100 runs confirm no balance movement: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors — the same numbers as before the core changed.
- `SkyOn(day, seed)` is derived rather than stored, so there is no save migration and no way for it to drift: the same day of the same campaign is the same day however many times it is loaded, and two campaigns get different fortnights. The Go tests hold all of that, plus that every kind happens, that most days are dry, and that **the streets are still wet the day after rain** — a street does not dry the moment the rain stops, which is the only reason `Sky` carries a number rather than a flag.
- What the view does with it: wet asphalt goes darker and cooler and takes less of the low sun; every lamp gets a second copy of itself smeared down the road, because a reflection stretches toward whoever is looking at it; and fog — the one weather that changes how far you can see — deepens the depth haze rather than touching the ground.
- Verified in a browser on isolated saves at three skies the core itself chose: **day 7 rain at 21:30** (dark roads, lamp reflections lying in them), **day 40 fog at 20:00** (the far blocks fading out well before the near ones), and the dry nights already photographed for comparison. The API was read directly to confirm the core is producing varied weather across days 4-44, including `day 16: clear, wet 0.45` — the drying rule visible in the data rather than only in the test.
- `mise run verify` passes, `npm test` 21/21, `cmd/apicheck` reports no invariant failures.

## The people in the room, and a fallback that was covering the faces
- **A bug found by looking closely.** `Portrait` put the generated face sheet on the element's own background and the crude drawn fallback in a child at `z-index:-1`. A negative z-index child still paints **above its parent's background**, so the drawing covered the generated face every time: the fallback was what everybody actually saw, everywhere in the game. The two are now layered the other way round — drawing underneath, generated face over it — so the drawing appears only when the sheet genuinely fails to load, which is what a fallback is for. Confirmed in the browser: five figures in Saint Agnes, five generated faces where two were blank drawn ovals before.
- The figures were portraits in rectangles standing on cones. They now have a contact shadow on the floor — the same thing that stopped the street props floating — an oval crop so the painted face reads as a head rather than a picture, shoulders instead of a point, and a coat lit from the side the room is lit from.
- **They also dim with the room.** The wash went under the figures, so at three in the morning the room went dark and everybody in it stayed lit like a shop window. They now take the same `roomLight().dark` the backdrop does.
- **Two things were tried and taken out**, which is most of what this iteration was: a hat, which over a rectangular portrait reads as a picture in a box rather than a head; and a bright collar, which read as a bow tie on every single person. Both are recorded in the CSS where they were, so they are not tried again.
- Verified in a browser at 12:12 and 03:10. No console errors. `mise run verify` passes, `npm test` 21/21.

## Good art that came out misaligned, and whose fault that was
- Generated buildings that look right on their own were landing wrong on the map: sunk into the pavement, shoved off their plot, or scaled down to nothing. That is not the artwork's fault and no amount of regenerating fixes it. **The view was assuming three things about every cut-out that are never true of a generated picture**: that the building's base is centred in the image, that the bottom edge of the image is the near corner of that base, and that the image is exactly as wide as the base. A model leaves whatever air it likes around the building, and the widest thing in the frame is usually a cornice or an awning rather than the footprint.
- `tools/fitiso.py` measures the footprint instead. In this projection a base edge climbs half a pixel for every pixel sideways, so it follows the underside of the silhouette out from its lowest point and stops where the underside stops following that line — which is exactly where a fire escape or a flight of steps begins, and those are not ground. The near corner and the base width go into `isometric.json`, and `CityIso.tsx` anchors each sprite at its own feet and scales it so its own ground matches the ground it is given.
- **It falls back rather than failing.** On the twelve painted addresses the strict measurement collapses, because those were generated with a soft cast shadow: the lowest pixels are a blurred blob, not two straight edges, and the walk stops after nine pixels. When the measured footprint comes back too small to be a building, it uses the content box instead — the picture trimmed of its air, which is most of the misalignment anyway. Every cut-out reports which reading it got.
- **`--proof` draws the footprint it found onto a copy of the picture**, so a bad fit can be looked at rather than argued about. On a Blender render, which has no cast shadow, the measured diamond traces the building's base exactly: `bar base 868px of 1024`.
- Existing art is unchanged by this, and that is the expected result rather than a disappointment: `tools/isometric.py` already trims its output, so the content box is the whole image and the new anchor is the old one. Verified in a browser at 12:12 — the city renders exactly as before, no console errors. What changes is that art from anywhere else now lands correctly without being hand-cropped first.
- `mise run verify` passes, `npm test` 21/21.

## A map editor, and what it is deliberately not allowed to do
- The city needs arranging by hand — which picture stands where — and two people need to be able to do it. So there is an editor, and the thing that makes it shared is not the interface but the file: `art/city-layout.json`, plain and in the repository. One person drags in the browser, the other edits the JSON, and it arrives as an ordinary diff either way.
- **It cannot move a building off its plot, and that is the point.** Where each address stands comes from the core's own coordinates; the grid, the blocks and the terrace slots are computed. The arrangement only says which sprite fills which slot and how many pixels it is nudged. A layout file able to disagree with the core about where The Monarch is would be a second city, and the first thing to go wrong would be a building standing in a road that the core says is a road.
- **Tiled was considered and not used.** It is a desktop application rather than a JS library, and it does support isometric maps. What it cannot do is know any of this city's constraints: that a building belongs inside a block's pavement ring, that a slot holds one building, that an address's position is derived, or where a sprite's measured footprint is. It would mean hand-aligning against a grid it does not understand, a converter, and a second source of truth.
- The slot is chosen by clicking **the ground rather than the picture**. A building overlaps its neighbours' ground by design, so the sprite under the pointer is frequently not the one whose slot was meant.
- **Saving is off unless asked for.** `BLACK_LEDGER_EDIT=1` enables the write endpoint; without it the same build reads the arrangement and refuses to write, with a message saying how to turn it on. Writing to the repository from a web request has no business being reachable from a running game.
- **A hazard found and fixed before it shipped.** The panel's edits took the arrangement as it stood at the last render, so two edits in quick succession would have had the second quietly undo the first. They take an update now instead of a value.
- Verified in a browser end to end on a clean load: click a slot, assign the sprite, save. `Block 4, 0, slot 0` produced exactly one entry, `"4,0,0": {"sprite": "bar"}`, and nothing else. A build without the flag returns *this build is not editable* and still serves the arrangement for reading.
- **Honest about the two stray entries seen while testing this**: they were mine, not the editor's — the file was reset *after* the page had already loaded it, so the browser was saving what it had read a minute earlier. Reproduced clean and the entries did not recur.
- `mise run verify` passes, `npm test` 21/21.

## Real ground, and the geometry left alone
- The roads and pavements are now photographs of a surface rather than flat colour — but the shape of every carriageway, pavement, kerb and plot is still computed from the grid. That is the whole point: junctions stay square, crossings still line up with the kerbs, the trolley rails still run down the exact centre of an avenue. Only the material is art. Placing road *sprites* would have meant straights, junctions, tees and corners all matching pixel-perfectly at their edges, and any that did not tile would show a seam on every street.
- **Laid in the ground plane, not pasted over the screen.** One tile east projects to `(TILE.w/2, TILE.h/2)` and one tile south to `(-TILE.w/2, TILE.h/2)`, so the texture's own axes are mapped onto those two vectors. The give-away for getting this wrong is a surface that does not turn the corner at a junction.
- **One texture did not tile, and it was measured rather than eyeballed.** `tools/tileable.py` compares opposite edges against ordinary interior variation — a photograph of asphalt is not uniform, so demanding a perfect edge match would reject every texture ever made. Asphalt, kerb and cobbles passed. **Pavement failed badly: its top and bottom edges differed by 66 against an interior variation of 11**, so the flag courses would not have lined up. Cropped to its own period, 1408 to 1057x1051, and it tiles. The fix is a crop rather than a blend, because blending edges together smears the courses, which on paving flags is worse than the seam.
- **The first scale was wrong and the browser said so.** A repeat of three flags across 1.5 tiles made a single flag 0.5 of a tile — wider than the 0.42-tile footway it was paving. The spans are now set so several flags fit across the pavement.
- The trolley avenue is cobbled rather than asphalted, since it already had rails down it.
- Verified in a browser at 12:12, zoomed in and out: flag courses running with the grid, asphalt grain on the carriageways, setts on the trolley avenue, no seams at the junctions. No console errors.
- `mise run verify` passes, `npm test` 21/21.

## The page is the viewport
- From the inbox, and stated twice in the Interface principles: the page must not scroll, only the areas inside it that are meant to. A screen that grows past the bottom of the window takes the header with it, and the header carries the clock, the money and the rail — the three things that have to be visible at all times.
- The shell is now the viewport height and `body` cannot scroll at all. Everything under the header takes the space that is left and scrolls inside itself: the sidebar, the room's work panel, the destination directory, the guide, the log. The city canvas stops guessing at `calc(100vh - 162px)` and simply fills what remains, which also means it no longer disagrees with the header when the header wraps.
- **Measured rather than eyeballed, on all six screens.** Document and body scroll height minus client height is **0 everywhere** — City, People, Families, Ledger, Herald, Guide. And nothing is lost to the change: a sweep for elements with `overflow:hidden` whose content exceeds their box found **none** on any screen, so no content became unreachable in exchange for the fixed height.
- One thing the browser showed that the numbers did not: on the guide, the scrolling element was the 750px reading column itself, so the scrollbar sat in the middle of the page beside a wide empty margin. The column is now inside the scroller rather than being it.
- `mise run verify` passes, `npm test` 21/21.

## The Herald as a piece of paper
- From the inbox. The typography was already right — masthead, double rule, columns with a rule between them, drop cap, halftone cut. What was wrong is that it was a rectangle of flat colour, so it read as a panel with newspaper styling rather than as newsprint.
- Three things make a sheet look like a sheet, and all three are drawn rather than photographed, so they cost nothing to load and scale to any size of issue: **a surface** (fractal-noise fibre at low contrast, enough to break up a flat fill and not enough to read as texture under 13px type), **an edge that was cut rather than ruled** (turbulence displacing the sheet's own outline), and **the memory of having been folded** (a long fold down the middle and a half fold across it, each a bright crease between two soft shadows). Plus handled edges browner than the middle, and a fifth of a degree of rotation, because paper does not lie flat.
- **The sheet and the print are two elements, and that is not tidiness.** CSS applies a filter *before* a mask, so a drop-shadow on the masked element would be cast by the rectangle the mask cut away — a hard rectangular shadow around a ragged sheet. The shadow belongs to the wrapper.
- **One correction from looking at it.** The first displacement was 9 and read as *torn*. Newsprint is cut, badly, on a machine that has been running all night: the edge wanders by a millimetre, it is not ripped. Down to 5, and the noise frequency evened out so the top and bottom edges are cut like the sides rather than left straight.
- Verified in a browser on an isolated copy of a day-22 save: the ragged trim, the fold running down through the masthead, the fibre in the stock, the shadow following the real edge. No console errors. `mise run verify` passes, `npm test` 21/21.
- QA note: that save had a pending event blocking the screen, so the copy's `event` was cleared in its own JSON. The live campaign was not touched.

## The city page
- From the inbox: filler about the city in the newspaper, to make it feel lived in. The Herald only ever printed things that had just happened to somebody — a killing, a robbery, a business changing hands. A real paper is mostly not that. Without the ordinary copy the city only ever speaks when it is being violent.
- Six sources, all of them facts the world already holds and all of them things the city could see for itself: **what the sky is doing** (and whether the streets are still wet from yesterday), **what a good is fetching against what it usually fetches**, **how hard the police are looking**, **how many people live here and how many of them answer to a family**, **premises standing with nobody in them**, and **Sunday**. Nothing on the page reveals who arranged anything, which is the paper's standing rule.
- **A decision worth recording: this is composed from facts rather than written by the director.** The inbox asked for the director to generate it. Filler has to appear every single day, including days when nothing happened, and a page that depends on a language model being installed and answering is a page that is sometimes blank. It is also archived: yesterday's paper has to carry yesterday's weather, and a model asked today would write today's. So the facts compose it, deterministically per day and per campaign, and the director layer is the next slice rather than the foundation.
- **Filed on the day turn, not composed on demand.** A back number carries the weather it was printed with. Generating on read would have last week's paper reporting this morning's rain.
- **A side effect caught before it shipped.** Two briefs a day would have quietly evicted the news: the archive is bounded at 240 stories, so filler alone would fill it in four months and a player reading back would find nothing but prices. `trimNews` now drops the oldest brief before it drops any real story, and the test files 480 days of filler and checks a killing from the first day is still in the paper.
- Verified in a browser on an isolated copy of a day-24 save driven forward by `cmd/apicheck`: a killing leading, a robbery under it, and **THE CITY IN BRIEF** below a double rule carrying *More officers on the streets* — which is true, that city's attention was 43 and watchful — and *A clear day*, which was the weather the core chose. No console errors.
- The tests hold that the page prints two different briefs on a quiet day, that it carries no scrutiny (`civic` is weight 0, so the ordinary edition costs the city nothing), that it never leads over a killing, that the same day of the same campaign always prints the same page, and that when it rained the paper said so.
- `mise run verify` passes, `npm test` 21/21, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## The director sets the city page, and is not allowed to lie in it
- The second half of the newspaper inbox entry: the filler should be written by the director. The fact-composed page shipped first and deliberately, because a page that depends on a model answering is sometimes a blank page. This adds the model on top of it: one brief a day is offered to the local model to be set in the register of a 1953 city paper, and the facts are the fallback whenever it is absent, slow or wrong.
- **The check is about addition, not quality.** Nothing in `core.AcceptPolish` can tell whether prose is any good; it can tell whether the prose is about the same facts. A rewrite is refused if it contains **a number that was not in the original** — the dangerous invention, because a reader takes a number as fact and acts on it — or **a capitalised name that was not there**, or if it addresses the reader, or runs past two or three sentences. A model asked to make a paragraph about the weather more colourful will put a named commissioner and a count of closed streets into it, and in a game whose premise is that the paper prints only what the city could see, an invented fact is not a blemish but a lie the player has no way to detect.
- A refusal is final and is recorded on the story. Asking the same model the same question again gets the same answer, and the paper is not going to sit there re-asking about Tuesday's weather forever.
- **Verified against the real model**, not only against the guards. Ran an isolated day-24 save with `qwen3:4b-instruct` and drove it forward with `cmd/apicheck`. Four briefs were filed and one came back accepted:
  - as filed: *"Cloud over the city and no sign of it lifting. The forecast, for what it has been worth lately, says the same again tomorrow."*
  - as printed: *"Cloud covers the city. The forecast says the same for tomorrow. It has not changed."*
  - Clipped, dry, adds nothing — and in the browser it sits beside a fact-composed brief in the same issue with no seam between them.
- The guards are tested without the model at all: a rewrite that invents hours and streets is refused, one that introduces Commissioner Vance into a weather report is refused, one that says "you" is refused, a model explaining itself before the copy has the preamble stripped, and a good rewrite is taken.
- Nothing waits on it. The request goes in the background after an action, shares the director's lock so two model calls cannot queue behind each other, and the interface never mentions it: the reader has no idea a rewrite was attempted.
- `mise run verify` passes, `npm test` 21/21, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## The workshop
- From the inbox: a debug mode on a separate port for playing moments from a menu. Every visual moment in this game is the end of a chain of decisions — to see what a killing looks like you have to get somebody killed, which means playing until somebody is worth killing. That is right for a player and useless for whoever is building it, and a moment that is hard to reach is a moment that gets looked at once.
- **Separate is the strongest form of separate**: its own binary (`cmd/playtest`), its own port (8790), its own page (`playtest.html`, a second Vite entry). There is no flag on the game that turns it on, no route to find, and nothing to leave switched on by accident in something a player runs — the game's binary contains none of it. `mise run workshop`.
- **It cannot open a save.** It builds a scratch world from a seed and rebuilds it per request, and it refuses to start at all if `BLACK_LEDGER_DB` is set, saying so. The value of a workshop is that you can do anything in it, which is only true if nothing in it matters.
- **The menu is the city's own list.** `core.MomentKinds()` is exported and the workshop reads gravity and hold from the core, so a kind added to the game turns up in the workshop without anybody remembering to add it. Three tests hold that: every kind the city can show is offered, they are offered heaviest first, and each holds for at least the floor.
- **A scrub control, added because the browser could not catch a moment.** A moment is over in three to four seconds, so a screenshot lands after it. Holding it still runs no clock and plays no sound: it puts the drawing at whatever instant you ask for and leaves it there. That is what finally proved the graphic renders.
- **A bug the looking found, and it was mine.** The workshop passed the *cue's* id as the spotlight id; `CityIso` looks that up in `state.locations` as an **address**. The lookup failed and nothing was drawn at all — which looks exactly like a moment that has no graphic. The game passes `playing.target`. Fixed, and the note is in the source so it is not repeated.
- **And something the workshop immediately showed that nobody had seen**: the explosion at 25% is a fireball covering four city blocks. It is drawn at `across * 2.6` where `across` is one plot. That is almost certainly too big now the city is a real grid rather than three rows of cards — recorded here rather than changed in the same slice.
- Verified in a browser: 12 addresses, 8 kinds with their gravity and hold, the city rendering behind, the theatre band with caption, actors and Herald headline, and the explosion held at 25%. No console errors. `mise run verify` passes, `npm test` 21/21.

## The ledger is accounts again, and the market has its own page
- The last of the inbox. The Ledger carried the underground market, the day's costs itemised, and the whole history, and the user's note was that the market "probably doesn't belong in there" and the daily cost maybe not either.
- **The market moved out to a page of its own.** A price is not an account. The ledger answers *what am I worth and what is this costing me*; somebody checking whether moonshine is worth moving today is not doing bookkeeping, they are deciding where to spend the afternoon. Burying one under the other made both harder. The new page is a board rather than a table — each good a card, because what is being compared is not columns of figures but "is this worth moving".
- It says more than the old block did, all of it from the world: the price against what the good **usually** costs as a percentage either way, what is in your hands and **what that is worth at today's price**, where the trade is actually done, and a band naming what holding it costs — *14 moonshine, 3 crated arms — 23 attention a day, every day, until it moves*.
- **The daily cost stayed, because it is genuinely owed money**, but it folds away. It is the detail behind a figure already on the page twice, and open by default it pushed the history below the fold. The summary carries the number, so nothing is hidden: *What the $206 a day is*.
- **What the ledger looks like now**: four figures, one disclosure, and then *What happened* with its search and its filters, all above the fold. Before, the search was three screens down.
- Verified in a browser on an isolated day-22 save with stock patched into it so the holding path was actually exercised: three goods, the risk band reading 23 attention a day, the ledger showing one heading where it had three, and the cost breakdown opening to six lines beginning *Rent $35*. No console errors.
- `mise run verify` passes, `npm test` 21/21. No core change, so balance is untouched by construction — this slice is entirely presentation.
- **The inbox is now empty.** Everything the user left in it has been built.

## Every moment was the wrong size, and one of them was invisible
- The workshop was built to look at moments, and the first thing it showed was that none of them had been looked at since the city stopped being three rows of cards. Everything is scaled from `across`, the width of the plot a moment happens on, and that number shrank by about two thirds when the city became a real grid. Nobody noticed, because a moment is over in four seconds.
- **The explosion covered four blocks.** It was drawn at `across * 2.6`, so a building going up put a fireball over the neighbours' roofs. `reach()` in `iso.ts` now owns how far every kind extends, as a multiple of its own plot, and it is tested: **no moment may reach past 0.6 of a block**, because a moment wider than the block it happens on stops saying *here* and starts saying *everywhere*, and the whole job of a moment is to tell the player where to look.
- **The test caught a second fault the ceiling would not have.** My first retune left a police lamp covering more ground than a building exploding. An explosion is the loudest thing in this game and a shot is the quietest, so the ordering is now asserted too: explosion > raid > killing, and the explosion has to come back down rather than staying open.
- **A moment was being drawn half a block behind the building it happened at.** The position came from `middle(cell)`, the centre of the *block*, but an address stands on the front-left slot of a terrace. It now uses the same slot the building was drawn on, and `SLOTS` moved to module scope so the moment layer and the city layer cannot disagree about the geometry.
- **Everything happened above the roof, including the things that happen on the pavement.** A man shot at the door, a police lamp on a car at the kerb and a robbery at a till were all played in the air over the chimney. `liftOf()` puts each kind where it belongs as a fraction of the building's own height: an explosion inside it at .55, a shot at street level at .1.
- **And the muzzle flash was invisible — drawing perfectly, and too small to see.** Held still, it was a three-pixel core at sixteen per cent alpha. That is the same as not existing. It was only provable by filling it magenta: the magenta landed exactly where it should, at the shopfront door, which told me the position was right and the *size* was the fault. Flash, police lamp and the plain ring were all rescaled and brightened.
- Verified by holding each kind still in the workshop and looking: the explosion engulfing Saint Agnes rather than the district, a white muzzle flash at the shopfront door, a red lamp sweeping the facade. No console errors.
- `mise run verify` passes, `npm test` 23/23. No core change: this is entirely the view, so balance is untouched by construction.
- The user added a line to `docs/LIVING_WORLD.md` overnight — that they expect the agent to think of these things, build them, and playtest them without being asked. Its formatting was repaired; the wording is theirs and untouched.

## Obituaries
- A killing is reported the same day, as it should be: a body was found, police say enquiries are continuing. That is the crime desk, and it is about what happened rather than about who it happened to. A paper carries the other thing the next morning, and this game did not.
- The Herald now runs an obituary the day after somebody the city knew dies. **It is the only place in the game where a person is described as a life rather than as a threat, an asset or an obstacle** — what they were, where those who had business with them had it, what carries on without them, and how many people stood below them and will be told by somebody.
- **It knows nothing it should not.** An obituary never says who arranged anything, and a test kills somebody with a named killer and asserts the killer's name never appears, along with "arranged", "ordered by" and "on behalf of".
- **Not everybody gets one.** A soldier nobody had heard of does not get a column; a man with a title always does, as does anybody who stood over other people, ran premises, or whom the player had actually dealt with. The dead are marked seen whether or not a column is written, so somebody the city would not have noticed is passed over once rather than reconsidered every morning forever.
- **It costs the city nothing.** The killing was counted yesterday; counting the death twice would have the city look hardest at the people who are mourned most.
- Two facts had to be added to the core to make this possible at all: the city recorded *that* somebody was dead and never *when*, so nothing could ask "who died yesterday" — which is the question a paper asks every morning. `w.Dead` turned out to hold only the player's deaths.
- **Three copy faults, all caught by reading the output rather than by a test passing**: "1 people in the same organization"; "They were Russo boss of Russo Outfit", because the shared helper appends the family to a role that already carries it; and an estate line naming a holding across the city while the man was found at the door of one he ran, which reads as the paper picking a building at random — which is what it was doing. All three now have tests.
- **A QA fault of mine, worth writing down.** A fresh copy of a save kept showing the *previous* run's obituaries: copying `foo.sqlite3` while leaving a stale `foo.sqlite3-wal` beside it replays the old write-ahead log over the new copy. Delete the `-wal` and `-shm` with it. I nearly recorded the old wording as a live bug in the new code.
- Verified in a browser on an isolated day-23 save: a front page carrying *SPLIT IN RUSSO OUTFIT*, then OBITUARIES centred under a rule for **VITTORIO BELLANDI** — *"They were Bellandi boss. Those who had business with them had it at The Monarch."* — then the city page below it. The weather brief in that same issue had been rewritten by the director, so all three layers of the paper were visible at once. No console errors.
- `mise run verify` passes, `npm test` 23/23, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## Measuring whether the living world is actually alive
- `docs/LIVING_WORLD.md` asks for these by name and they did not exist: *wars started, holdings changed hands, factions created and destroyed, and how often a run is affected by a conflict it had no part in*. A simulation reporting only deaths and cash cannot say whether layers 2, 3 and 4 do anything — and the same document sets the bar: a system is finished when a long simulation shows it producing varied, non-degenerate outcomes.
- `cmd/simulate` now reports a `city` block per strategy: each measure as a total **and** as how many runs saw any of it at all, because a world producing one war in a hundred campaigns is not a living world and the total alone would hide that.
- **The first version of the measure lied, and the shape of the lie is the useful part.** It counted about one new organization per campaign for the strategies that form one and almost none for the others — a number about the strategies, not the city. It was counting the player naming their own outfit as the city making a new family. Excluded now, and the test holds it: the player's organization forming or dissolving is the player playing, and the player taking premises is not the city moving them.
- **What 400 campaigns actually show**, and it is not what the document intends:

| | total | runs with any, of 400 |
|---|---|---|
| wars started | 52 | **43** |
| wars the player was no party to | 33 | 25 |
| holdings changed hands between families | 34 | 24 |
| organizations created | 40 | 36 |
| organizations destroyed | 5 | **4** |
| player hurt during somebody else's war | 58 | 12 |

- Read plainly: **a war starts in about one campaign in nine, and an organization is destroyed in one campaign in a hundred.** The layers are not broken — they fire, and they vary by strategy in ways that make sense, with `defiant` seeing eleven wars of which *none* were between other families, while `investor` and `worker` mostly watch other people's. But the city is sparse: nine campaigns in ten never see it do anything to itself.
- That is a finding, not a fix. Making the city more active is a balance change and deserves its own slice with before-and-after numbers, which is what this measure now makes possible. **The baseline is recorded above.**
- `mise run verify` passes, `npm test` 23/23. Balance unchanged: this only observes — defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors, exactly as before.

## A correction: the city is alive, and my measurement was too short
- The previous entry concluded that the living world barely fires — a war in 43 campaigns of 400, an organization destroyed in 4. **That conclusion was wrong, and the commit message carrying it (`f02725d`) is wrong too.** It is left in the history rather than rewritten, and this is the correction.
- The fault was the horizon. `cmd/simulate` defaults to 200 commands, which is a **median of 12.6 game days**, and for `reckless` **0.2 days**. Families escalate, split and fall over weeks. Asking whether the city moved in a fortnight and concluding that it does not move is a fact about the run, not about the city.
- **Re-measured at 1000 commands**, where the surviving strategies reach 60 to 70 game days:

| | at 200 commands (12.6-day median) | at 1000 commands (`investor`, 71.6-day median) |
|---|---|---|
| runs where a war started | 15 of 100 | **38 of 40** |
| runs where an organization was destroyed | 1 of 100 | **25 of 40** |
| runs where holdings changed hands | 7 of 100 | **30 of 40** |
| wars the player was no party to | 8 | 90 |

- `worker` tells the same story: a war in 28 of 40 runs, holdings moving in 26. Across all 160 long runs a war starts in 54% and an organization dies in 28% — and that average is dragged down by the two strategies that get themselves killed in the first week, not by a quiet city. **The living world does what the document asks. Nobody had looked at it over a long enough campaign.**
- So the fix was to the measurement, not the world. The summary now reports **`median_game_days`** rather than only minutes, and **`city_measures_meaningful`** per strategy, and the run prints a warning naming any strategy that ended before the twenty-day horizon: *"a campaign has to run past about 20 days before families have time to escalate, split or fall, and these ended sooner."* At the default every strategy is flagged, which is exactly the guard that would have stopped the wrong conclusion.
- **A separate finding worth its own look later**: `reckless` has a median campaign of **0.2 game days** — under five hours — and dies in 31 of 40 long runs. Dying is that strategy's job, but dying before lunch on the first day suggests something degenerate in its opening rather than a hard game.
- No world change: the city was never touched. `mise run verify` passes, `npm test` 23/23, balance identical at defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## What the street says about the families
- The previous entry proved the city is busy: over a long campaign a war starts in most runs, organizations are destroyed in a quarter of them, holdings move between families constantly. **None of which the player could feel.** To know the Russo Outfit had been losing for a fortnight you had to open a screen and compare a number to a number you had not written down. A city that is busy behind glass is not a living city.
- The paper now notices. When a family's strength moves more than `FortuneShift` since the paper last mentioned them, the city page carries it — *BELLANDI FAMILY SAID TO BE STRUGGLING*, "a poor few weeks by the reckoning of people who watch such things" — and the wording changes with what they are actually still holding: nothing at all, one premises left, or nobody at the family saying what went wrong.
- **It never prints the number**, and a test enforces that: no digit may appear in a fortunes brief. "Power 41" is a statistic; a paper writes about what people have noticed. The other three tests hold that a fall is noticed, that **the same slide is not reported every morning for a week** — the family remembers what was last said about it — that a small drift is not news, and that the player's own organization is never reported to the player, because it is not news to them.
- One thing this needed in the core: `Reported` on a family. Saves written before it carry none, and the first day simply records where the family stands rather than announcing a change from nothing.
- Verified in a browser on an isolated save with a family knocked down thirty points, as a lost war would: the brief sits in **THE CITY IN BRIEF** beside the population and the weather, in a back issue from day 23. No console errors.
- Also resolved, a false alarm of mine from the previous entry: `reckless` has a median campaign of 0.2 days because the strategy provokes a family on its first move and then rests at home without security. Dying quickly is what it is for. Nothing degenerate.
- `mise run verify` passes, `npm test` 23/23, `cmd/apicheck` reports no invariant failures, 100 runs unchanged: defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## The whole arc of a quarrel, not just the middle of it
- The paper announced every war and never once said one was over. As far as a reader could tell, every war ever declared was still running. And every war was announced under the same headline — *OPEN WAR ON THE WATERFRONT* — so somebody who read the paper twice could not tell that the second one was a different war.
- A quarrel now reads as an arc. **It hardens** (*BAD BLOOD BETWEEN BELLANDI FAMILY AND RUSSO OUTFIT*, on the city page, costing the city nothing), **it becomes a war** (named for both families rather than for the waterfront), and **it stops** — reported quietly, as politics rather than as war, because if the ending carried the same weight as the beginning the city would look hardest at the moment the shooting stopped. A test holds that ordering.
- **How it stopped is not one story.** A war that burned out and a war that finished somebody are the same transition in the model and completely different in the city, so `howItEnded` reads the ground: if a side is holding nothing, it names them and says whether that is the end of them is a question nobody is asking out loud; otherwise both are smaller than they were and both are still here.
- **A bug the API run caught that four passing tests did not.** A war does not only end by going cold — it usually decays into a *feud* first, and every test I had written was looking at war-to-cold. The first version therefore printed *"bad blood between them"* on the day the shooting stopped: the hardening story, told at the moment of the ending. Arriving at a feud means two opposite things depending on where it came from, and the code now asks. There is a test for it, and it is the case the tests had missed.
- Verified against a running save with a war set one point above the threshold it ends at: the paper carried *THE FIGHTING STOPS BETWEEN BELLANDI FAMILY AND RUSSO OUTFIT — "Both are smaller than they were, and both are still here."* `cmd/apicheck` reports no invariant failures.
- `mise run verify` passes, `npm test` 23/23, balance unchanged at defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors. Long-horizon run (20 × 4 at 1000 commands) consistent with the recorded baseline: wars in 40 of 80 runs, organizations destroyed in 18.

## What leads the paper is not what the police notice
- **A regression I introduced tonight, found by looking at a real save.** The city screen's news banner showed `newspaper[0]` — the *newest* story. Once the city page began filing the weather every morning, the newest story was usually the weather. Confirmed on a driven save: the newest story was **A CLEAR DAY**, so the banner would have announced cloud cover over a man shot that afternoon.
- Underneath it was an older fault. The banner and the edition lead were both ranked by `scrutinyWeight`, which is not a measure of how big a story is — it is **what a story costs the player in police attention**, and it does not cover half the kinds the paper files. An arrest, an attempt on somebody's life, a family splitting and a family moving back onto premises all score zero there, correctly, because none of them is a reason for the city to look harder at the player. That zero was also deciding what leads an edition, so **an arrest could never lead the paper**.
- `Newsworthiness` is now a separate ordering covering every kind the paper files, from a killing at 100 down to the city page at 5. Kept deliberately apart: changing what a story costs in attention changes the balance of the game, and changing what leads the paper changes only what the player is told first. A test asserts the separation by pinning the scrutiny values a killing and an arrest carry, so a future edit to the front page cannot quietly move the balance.
- The banner now shows the biggest unread story rather than the latest.
- Verified in the browser on that same save: the banner reads **IVO COSTA KILLED** where it would have read *A CLEAR DAY*. No console errors.
- `mise run verify` passes, `npm test` 23/23, `cmd/apicheck` reports no invariant failures, balance unchanged at defiant 58 / investor 0 / reckless 82 / worker 0, 0 errors.

## The paper stops printing the same story seven times

I drove a QA save (a copy of a backup, with the stale `-wal` and `-shm` deleted
beside it) forward to day 42 with eight `cmd/apicheck` runs, dumped all 28
editions in order, and read them end to end as a story. LIVING_WORLD.md says a
system is finished when the results still make sense read back that way, and
nobody had done it.

Most of it read well. A war started, a killing avenged a killing and said so
("It was over what happened to Luca Toth"), a family split, the Russo outfit
ended, the dead got obituaries, and the fighting stopped. That arc is the thing
the simulation is for and it survived being read as prose.

One issue was ruined. The same robbery at Saint Agnes was printed seven times in
a single edition, each time in full, and each copy carried a doubled police
line: "Police have asked anybody who saw it to come forward. Police have made no
arrest and are appealing for anyone who saw the incident at Saint Agnes." No
paper does that. It writes one piece saying it happened seven times, and that
piece is a better story than any of the seven.

`Report` now collapses a repeat into the story already filed for that kind and
headline today, in this life, bumping a `Count` and rewriting the body with a
run sentence appropriate to the kind — a robbery "happened three times in the
same day, which residents say is not the usual run of things"; police "were back
three times before the day was out". A different headline is still a different
story, yesterday's story is never today's, and `unattributed()` no longer
prepends "Police" to a sentence that already names them.

Evidence: four new properties in `core/repeats_test.go` (a run collapses, two
distinct robberies stay two stories, a day boundary breaks a run, the police
line is not printed twice). Two older archive tests filed one identical headline
hundreds of times to test boundedness; that is now a run of length N, so they
file distinct headlines instead, which is what a bounded-archive test actually
means. `mise run verify` and `npm test` pass. `mise run simulate` is unchanged
at defiant 58 / investor 0 / reckless 82 / worker 0 deaths, 0 errors.

## Five things the paper got wrong about English

Having read one campaign's paper end to end, I read another, and this time
checked every sentence rather than the story. The world was right in all of it.
The English was not.

- `the Rizzo Crew has ceased to operate` — the splinter naming forms bake a
  lower-case article into the name, so any sentence starting with one started in
  lower case.
- `Cesare Ferro's people has had a poor few weeks` — the player's own
  organization is named with a plural noun and took a singular verb.
- `Crated arms is fetching more than it did` — two of the three goods have
  plural names.
- `They were Lieutenant.` — the role went in with no article, in the obituary
  and, from a second function, in the killing story.
- `SUNDAY IN BELLWETHER` under a dateline reading Monday, April 13 — the city
  page did its own weekday arithmetic and was one day out.

The first four are one problem: prose is assembled from names and roles that
nothing checks the grammar of. `core/names.go` now holds `Leads` (capitalise an
article the name carries), `Agree` (verb form for an organization name) and
`article` (a/an), and the newspaper sites use them. Whether a good's name is
plural is deliberately *not* stored in the save — it is a fact about the word,
not the market, and a campaign begun before tonight would otherwise print the
old sentence forever. The Sunday fault is now impossible: `Weekday` and
`IsSunday` are the only place that knows what day it is, and `Dateline` uses
them too.

One more thing came out of the same reading: a protagonist dying at midday split
the archive, so DAY 41 appeared twice with the same dateline. One city, one
paper, one issue a day; `Editions` groups by day alone and the issue takes the
latest life it carries.

Evidence: six properties in `core/copy_test.go` and one in
`core/repeats_test.go`, each stating the fault it came from. Then the real
check — a QA save driven eight days forward and every sentence of all eighteen
issues scanned for the five patterns, plus any body starting in lower case:
zero. `mise run verify` and `npm test` pass, `mise run simulate` unchanged at
58 / 0 / 82 / 0, 0 errors.

## The prose is now checked by the thing that drives the game

Fixing the paper's grammar by hand found the same faults again in the ledger,
which the player reads far more often: "the Lindqvist Combine no longer holds
anything worth defending", "Cesare Ferro's people moved against Saint Agnes and
was driven off". Fixing those found a third site — the warning that a family
"has people asking where you sleep". At that point the pattern was the point:
the fault is not in any one sentence, it is that nothing checks the sentences.

Two things came out of it.

`cmd/apicheck` now reads everything the city has written down. After a run it
scans every ledger record and every story in every issue for a body beginning in
lower case, a plural organization name with a singular verb, a plural good with
a singular verb, a role used without an article, a doubled police line, a day
printed twice, and a Sunday page under a weekday dateline. Copy faults are
invariant failures now, reported with the sentence, and every future run checks
them for free. This is the right home for it: the faults only exist in assembled
prose, and apicheck is the only thing that assembles prose from real play.

And the root cause turned out to be one line. Splinter names were minted as
"the %s Crew" while every seeded family is "Bellandi Family" or "Russo Outfit",
so a splinter carried an article into the start of any sentence it began.
Removing the article from `splinterForms` fixes every site at once and always
will. `Leads` stays, because saves written before tonight still hold the old
names and the player still reads those records.

Evidence: fourteen `apicheck` runs against a fresh campaign, zero invariant
failures. Run against a save written by the older build, the same check reports
two — correctly, because those sentences are still in that player's ledger.
`mise run verify`, `npm test` and `mise run simulate` unchanged at 58 / 0 / 82 /
0, 0 errors.

## The city keeps hours

I measured before building anything. Over a simulated week, sampled every hour,
**forty-seven of fifty people never moved at all** and three moved once between
them. The three busiest places at three in the morning were the three busiest
places at three in the afternoon, by the same margin. The city's population was
furniture.

This was not for want of a movement system. `core/errands.go` is good: a man
crosses town to settle a grudge, a family sends somebody to mind a holding it
has just lost. But its own comment says it — "nobody in this city moves without
a reason that ends" — and every reason it knows is a reason that ends. Once
everybody had a reason to be where they already stood, nobody moved again for
the rest of the campaign.

What was missing was the ordinary reason. The day ends and people go out.

`core/routine.go` gives the city a shift. In the first half of the day people
are at their posts; in the second they are where they drink. Where somebody
drinks is derived from their id rather than stored, so it never drifts and
survives every save ever written — the whole point is that Ivo Costa is at The
Blue Hour of an evening and is there every evening for the rest of his life.
Anybody with work to do keeps doing it: a role holder is on duty, a family
member with ground to mind is minding it, and a leader is not found propping up
a bar. The resolution is a half-day because that is the resolution the clock
has, and pretending to finer grain would be a lie told by this file rather than
a fact about the city.

**A stampede is not a shift change.** The first version emptied every building
at the same minute and put thirty-nine of fifty people on the street at once —
caught immediately by an existing test that had been written for a world where
walking was rare. People now leave over about five hours, at a time fixed per
person, so the one who always leaves early always leaves early. Deciding to go
and going are no longer the same minute, which is also more truthful: the room
notices somebody leave at the moment they walk out of it.

That change quietly broke the test that caught it. Staggered departures meant
the test's one sample per half-day always landed before anybody had set off, so
it saw an empty street and passed on nothing. It now steps through the half-day
an hour at a time, and asserts a floor as well as a ceiling so it cannot go
blind again.

Measured after, same probe, two weeks:

| | before | after |
|---|---|---|
| people who moved in a week | 3 of 50 | 36 of 50 |
| busiest room at 09:00 | Mercer Exchange | Mercer Exchange |
| busiest room at 21:00 | Mercer Exchange | Saint Agnes |
| most people on the street at once | 0 | 10 |

And the promise the whole thing rests on: of people who move, **860 of 864
person-hours find them where they usually are at that hour**. A rhythm nobody
can learn is just noise.

Balance moved slightly and in the right direction: defiant deaths **58 → 53**,
reckless unchanged at 82, investor and worker still 0, no errors. A city where
people are not standing still all day is a marginally less reliable shooting
gallery. New baseline is 53 / 0 / 82 / 0.

Also fixed, seen in the live game's own result banner: "The a professional you
paid $2501 to reach Elena Russo did not finish it." The contract tiers are
labelled for buttons — "A professional", "Someone who needs the money" — and a
label carrying its own article cannot be given another one. Tiers now have a
separate noun for prose, and `cmd/apicheck` checks for the whole class.

Evidence: six properties in `core/routine_test.go`, one in `core/copy_test.go`,
ten clean `apicheck` runs on a fresh campaign, `mise run verify` and `npm test`
green.

Not yet looked at: a bar holds thirteen people at nine at night, and I have not
seen what the room panel does with that. The list is already filtered to who
you can deal with and it lives in the scrolling sidebar rather than the main
viewport, so it is probably fine, but probably is not looked at.

## The room says whether it is busy

Two things from standing in the game and looking.

**The open question from last time is answered: no.** A bar holds fourteen
people at ten at night now, and I had not seen what the room panel does with
that. It does the right thing already. The header reads "IN THE ROOM · 3 of 14
you can deal with", three cards render, and the rest sit behind "Show 11 others
in the room". No overcrowding, no main-viewport scroll, no change made. Worth
recording as a negative result: I went looking for a problem I had created and
there was not one.

**What was actually wrong is that none of the new rhythm reached the player.**
Standing in Saint Agnes at ten at night with fourteen people in it read exactly
like standing in it at nine in the morning with two. The only trace of the
difference was a number in the corner of a filter, and a number with nothing to
compare it against says nothing. Having built a city that keeps hours, I had
left it invisible from inside.

`core/crowd.go` gives the room a sentence about itself, and only when there is
something to say. Bands are read as a share of everybody still living, so the
description keeps meaning something as the population rises and falls, and the
evening wording differs from the daytime wording because a full bar at ten at
night is a different fact from a crowded market at nine. A room with an
unremarkable number of people in it says nothing at all — a note on every room
is a note the player stops reading.

Read out of a real save at ten at night:

```
Saint Agnes    14  Saint Agnes is full tonight. 14 people, and the ones by the
                   door are watching who comes in.
Pier 14         1  Nearly empty at this hour: one other person, and nobody else
                   worth counting.
The Mariner     0  There is nobody else in here.
```

The note is a description and nothing else. One of the five properties in
`core/crowd_test.go` exists only to hold that line: filling a room with the
entire city must not change what can be done in it. Another checks that
somebody out on the street is not counted into the room they have left.

`cmd/apicheck` now scans the room notes along with the ledger and the paper,
since they are assembled prose like everything else.

Evidence: five properties in `core/crowd_test.go`, ten clean apicheck runs on a
fresh campaign, `mise run verify` and `npm test` green, `mise run simulate`
unchanged at 53 / 0 / 82 / 0.

## Coming home

A trip out of Bellwether takes two to four days and the clock runs the whole
time. I measured whether that means anything: forty campaigns, driven to day
twelve, then four days away.

| in four days away | in 40 campaigns |
|---|---|
| the paper printed something | 40 |
| a family's power moved by more than 8 | 3 |
| somebody died | 2 |
| a holding changed hands, a family formed or fell, a quarrel started or ended | 1 each |

So the city does run, and the structural changes are properly rare over four
days, which is right. But the record filed on the player's return read "4 days
gone. The city did not wait" — a claim with nothing in it. Everything that had
happened was sitting in the Herald, filed under days the player had no reason to
go back and read. Coming home is the moment that information is worth most, and
the game was silent at exactly that moment.

The return record now carries what the paper carried. The headlines are printed
as the paper printed them, in the paper's own order of importance, because
sentence-casing them would lower-case the names and paraphrasing them would be
this file inventing news. Civic filler is dropped — somebody back from four days
away does not need to be told it rained. Read out of a real save:

```
3 days gone. While you were gone the paper carried: POLICE PRESSURE ON RUSSO
OUTFIT · MERCER EXCHANGE CHANGES HANDS · RUSSO OUTFIT MOVES INTO BLUEBIRD
LAUNDRY. And 1 other story.
```

Half the trips come back to "Nothing in the paper you needed to be here for",
which is honest rather than a failure: over four days at day twelve, half the
time nothing above civic filler happened.

**Two faults found by reading around the thing I was building.** The first was
in the line next to it: "Whatever was arranged for you happened 1 times to a
locked door." The second only appeared once the summary existed: "POLICE
PRESSURE ON RUSSO OUTFIT · POLICE PRESSURE ON RUSSO OUTFIT" — the paper
collapses repeats within a day, a trip spans several, and two days of the same
headline is one thing to be told.

**And a fault in the checker itself.** I added a count-agreement pattern to
`cmd/apicheck` as a plain substring, and it reported "11 people in here, which
for this hour is a crowd" as a fault four runs out of ten. A check that cries
wolf is worse than no check. It is a bounded regexp now, and fourteen
consecutive runs on a fresh campaign are clean.

Two other sites where a count of one was reachable are fixed: a one-day sentence
served, and one day of somebody else's bail. The rest of the `%d days` sites in
the codebase are not patched blind — the check will surface them if real play
ever produces a one.

Not covered, stated plainly: `cmd/apicheck` has never once taken a trip. All
three appear in its "never tried" list every run, because the fare is out of
reach at the cash its play reaches. This slice was verified by driving a save to
day fourteen, setting the cash directly in the save file, and taking all three
trips through the HTTP API by hand.

Evidence: seven properties in `core/away_test.go`, fourteen clean apicheck runs,
`mise run verify` and `npm test` green, `mise run simulate` unchanged at
53 / 0 / 82 / 0.

## What the harness never sees, and a correction

Last night I wrote that `cmd/apicheck` never takes a trip "because the fare is
out of reach at the cash its play reaches." **That was wrong.** I had not
checked; I inferred it from having seen a disabled trip button next to $80. The
game's own recorded reason, across fifteen runs, is `You are in no condition to
travel` — the harness plays recklessly, ends up hurt, and health is what stops
it. In several runs the trip was offered *and enabled* and the harness simply
did not take it, for a reason that turned out to be its own.

Chasing that down turned up three things.

**A bug that made every JSON report useless.** The coverage list was computed
after the report was written to disk, so `never_tried` was empty in every report
file ever produced while the terminal printed the real one. Anyone reading the
reports rather than watching the run would have concluded coverage was perfect.

**A coverage list that could not tell you anything.** "56 never tried" does not
distinguish a system the game never offered from one it offered and greyed out
from one it offered, enabled, and the harness walked past — and those want
opposite fixes. The report now separates them and carries the game's own reason
for the greying:

| why a venture went untried in 15 runs | before | after |
|---|---|---|
| offered, enabled, and never taken | 10 | 8 |
| offered but always out of reach | 31 | 27 |
| never appeared in any action list | 15 | 13 |
| kinds exercised at least once | 39 | 47 |

**And a real defect in the harness, which that split found.** The driver already
preferred an untried venture over a repeat — but it scanned the venture list
from the top every single turn, so an untried venture near the front won every
time and anything late in the list was starved no matter how often the game
offered it. All three trips are near the end. The scan now starts at a different
point each turn.

That took kinds exercised from 39 to 47 across fifteen runs, and moved eight
systems out of the untried list, `lie_low`, `hire`, `mug:crew`, `arms:armour`,
`buy:cigarettes` and both operating modes among them. Fifteen runs, zero
invariant failures, so nothing newly exercised was broken.

What is left is honest and not a harness problem. Twenty-seven systems are gated
behind wealth or standing the harness's play never reaches — `car`, `still`,
`armoury`, the three retainers, both service contracts, the safe and the cellar.
Thirteen more never appear at all because they need a state it never gets into:
`bankroll` and `draw` want a casino, `sell:arms` wants arms in the ground,
`lawyer` wants an arrest, `unpost` wants somebody posted. Those are the expensive
half of the game and they are still unverified over the HTTP stack. Saying so
precisely is worth more than the vague version I published yesterday.

`mise run verify`, `npm test` green, `mise run simulate` unchanged at
53 / 0 / 82 / 0.

## A harness that stays on its feet

The measurement first, because last time I guessed and was wrong. Fifteen runs
of eighty commands each, and I printed what state the player was actually left
in: the harness had reached **life fourteen**. It was dying every run or two and
starting again from nothing. Cash never passed $400, respect never passed 40,
heat pegged at 100 repeatedly, and it owned a laundry twice and nothing else.

Then I asked the game what killed it. Thirteen deaths, and **twelve carried the
same line: "Caught on the street by somebody else's war."** It travelled
relentlessly looking for systems it had not tried, at any health, through any
war. `travel` was its most-used command by a factor of two.

So the reason half the game was unverified was not that those systems are
expensive. It is that the harness never lived long enough to afford anything.

Two lines fix it: rest when below 55 health, lie low when above 85 attention.
That is not timidity, and it costs nothing in coverage — resting is itself a
command, and a player who is hurt goes home rather than walking across town.

| across 15 runs | before | after |
|---|---|---|
| lives used (deaths) | 14 | 5 |
| kinds exercised | 47 | 66 |
| best cash reached | ~$400 | $1091 |
| best respect reached | 38 | 102 |
| premises owned | a laundry, twice | laundry, garage, casino |

**Then it found four real bugs, which is the entire point.** Living long enough
to own a casino and name an organization put the harness into sentences nothing
had ever produced before, and the copy scan caught all of them: "Franca
Sabbatini's people has people asking where you sleep", "Franca Sabbatini's
people has taken Saint Agnes from Falcone Crew", "Franca Sabbatini's people is
short of people", and — in the same line — "They were soldier a week ago", the
missing-article fault in a fourth site. All four are the same two helpers
applied: `Leads` and `Agree`, and `article`.

**And a fault in the checker, for the second time tonight.** Four of the twelve
reports were "Violence between Brenner Company and Franca Sabbatini's people has
escalated beyond the usual", which is correct English — the subject is the
violence, not the people. The rule was a plain substring. It is now anchored at
a sentence start and refuses to cross a comma or a second party, so it catches
the fault and leaves the correct sentence alone.

Both copy rules now have their own test in `cmd/apicheck`, with the sentences
they wrongly reported written down as cases. A check that cries wolf gets
ignored, and an ignored check is worse than none.

Fifteen runs after the fixes: 66 kinds exercised, zero invariant failures.
`mise run verify`, `npm test` green, `mise run simulate` unchanged at
53 / 0 / 82 / 0 — the harness changes touch nothing the game does.

Still unverified over HTTP, and honestly so: 19 systems behind wealth or
standing even a surviving harness does not reach, and 12 that never appear
because it never enters the state, `lawyer` wanting an arrest and `sell:arms`
wanting arms in the ground among them.

## Coming out

A sentence runs from two days to twelve, which is three times the longest trip
out of the city, and the record filed on release said only that whatever it cost
"happened while you were not there to watch it." The same empty claim the trips
made, behind a different door. `WhatYouMissed` was written to be reused, so this
should have been a two-line change. It was, and then reading the output found
three faults, one of them in my own measurement.

**The first measurement was wrong and I nearly published it.** I measured that
forty campaigns out of forty came out of a two-day sentence to real news, which
was too good to be true and was. The paper files the story of the player's own
arrest at the exact minute the door shuts, so every single release was reporting
`MAN CHARGED AFTER DISTRICT SEARCHES` back to the man it had charged. The minute
you leave is a minute you were there for. With that excluded, the honest curve:

| days inside | of 60 sentences, came out to real news |
|---|---|
| 2 | 21 |
| 4 | 33 |
| 6 | 40 |
| 8 | 51 |
| 12 | 57 |

That shape is right: two days and it is a coin flip, twelve days and the city
has almost always moved.

**Then the same fault at the other end.** Talking your way out files `CHARGES
DROPPED AFTER COOPERATION` at the moment you walk through the door, and the
release record read that back too. The window is open at both ends now — the
minute it shut and the minute it opened are both minutes the player was present
for.

**And a line I had been reading past.** "Whatever what was found at the laundry
cost you happened while you were not there to watch it." The reason phrase was
spliced straight after "Whatever". It reads "You went in for what was found at
the laundry, and whatever that cost you happened while you were not there to
watch it" now.

Read out of a real save, eight days served:

```
Nobody meets you. You went in for what was found at the laundry, and whatever
that cost you happened while you were not there to watch it. While you were
gone the paper carried: SOFIA DOYLE FOUND DEAD · END OF BRENNER COMPANY ·
POLICE PRESSURE ON BELLANDI FAMILY. And 3 other stories.
```

**Second bird: three systems verified over HTTP for the first time.** `sit_out`,
`lawyer` and `talk` were all in the "never appeared in any action list" group,
and this confirms exactly why — they exist only inside a cell, and the harness
has never been arrested. All three were driven through the API against a save
with the sentence forced into it, and all three work. That does not close the
coverage gap, but it does mean those three are no longer unverified.

Evidence: three more properties in `core/away_test.go`, twelve clean apicheck
runs, `mise run verify` and `npm test` green, `mise run simulate` unchanged at
53 / 0 / 82 / 0. Six existing tests in that file were filing their stories at
the same minute as the return and had to advance the clock the way real play
does — the exclusion is at both ends, and they were leaning on neither.

## The arms loop, driven end to end

The armoury is the most profitable system in the game and the most dangerous,
and no part of it had ever been exercised over HTTP. `armoury`, `buy:arms`,
`stock_arms` and `sell:arms` were all in the harness's untried list. I forced a
save into a state where the player owned the laundry with money in hand, and
drove the whole loop through the API by hand.

It works, and the shape of it is right:

```
A room under Bluebird Laundry: $1200 of brick, board and a door that locks
from the outside.
A quiet purchase: 5 crates of arms for $1100, at $220 each.
Crates under Bluebird Laundry: 20 crates off your back and into the room.
Somebody came to Bluebird Laundry: Bellandi Family took 2 crates and left
$1230. They are 39 strong now, and there are 18 crates left under the floor.
Falcone Crew took 3 crates and left $1845.
Rizzo Crew took 2 crates and left $1230.
```

Bought at $220, sold at $615, seven crates gone in one day, and attention up
from four to seven. That is the system doing exactly what its own comment says
it should. Selling is passive by design — the wars come to the door — which is
why `sell:arms` never appears at the armoury; it is the ordinary waterfront sale
of what the player is carrying, and it works too.

**Two copy faults found by reading the screen while doing it.**

The sidebar read `5 families buying` as `5 at war`, when two were at war and
three were merely feuding. The number was honest — a feuding family buys — and
the word was not.

And the ledger read `5 crates of Crated arms for $1100`. The market lists a good
by a name fit for a price board, and that name cannot follow a count of units.

**I fixed the second one wrong, in exactly the way I had already written down
tonight not to.** My first version put the sentence form in a field on `Good`,
which is stored in the save — so every campaign begun before tonight kept
reading "5 crates of crated arms", because its saved goods had no such field. I
had hit this on the plural names hours earlier, documented it as "a fact about
the word, not the market", and then did it again. It is keyed by id now, and the
test loads a good with the old name to prove an old save reads correctly.

`cmd/apicheck` scans for the price-board names after a unit count.

Evidence: the loop driven end to end through the API against a forced save, one
new property in `core/copy_test.go`, twelve clean apicheck runs, `mise run
verify` and `npm test` green, `mise run simulate` unchanged at 53 / 0 / 82 / 0.

## Reading the paper a second time

Eight copy fixes had landed since anyone last read a campaign's paper end to
end, and that read is what found most of them. So I drove a fresh save to day
forty-seven, forty-six issues and two hundred and forty stories, and read it in
order.

It reads. Russo Outfit dies over ten days, and the sequence is legible as a
story: an attempt on Vittorio Bellandi, an explosion at Mercer Exchange, Ennio
Zanetti knifed in the crowd *over what happened to Perla Moreau*, Elena Russo
shot at the loading doors having *gone after Vittorio Bellandi over what
happened to Ennio Zanetti*, END OF RUSSO OUTFIT, and the obituaries following a
day behind each death. That is the thing the simulation is for and it survives
being read as prose.

**Four faults, and two of them were mine from earlier tonight.**

`They were a head of the Russo Outfit.` The fix for "They were Lieutenant" gave
every role an article, and a role that already carries its own complement takes
none. A titled office is now recognised as one.

`There were twice such incidents before the day was out.` The run-collapsing
used a frequency where a count belongs. It says "two such incidents" now, while
the robbery sentence keeps the frequency, because that sentence wants one.

`2 people in the same organization stood below them.` A 1953 newspaper spells
small numbers, and prose with a numeral in it reads like a report from a
machine, which it was.

`Franca Sabbatini's people now controls Bluebird Laundry.` A plural name with a
singular verb in a site the plural-subject scan did not cover, because its verb
list stopped at "holds".

**And one that is not a copy fault at all.** `DETECTIVE HARLOW TAKES OVER THE
BLUE HOUR` — the city detective had walked off his beat and seized a casino.
The guard on that excluded the officials and the heads of organizations but not
the people holding the city's standing jobs. `keepsPost`, written for the NPC
routines, is exactly that predicate, and it is the guard now. Premises still
turn over: five in the next fifty-three-day campaign, all ordinary people.
Balance moved defiant 53 → 52, which is a real behavioural change and a small
one. New baseline 52 / 0 / 82 / 0.

**The scan cried wolf for the third time, and this one is worth writing down.**
Widening the plural-subject verb list, I added "moved" and "came" — and it
immediately reported "Nico Ward's people moved against Saint Agnes and were
driven off", which is correct English, because a past tense is the same for
singular and plural. Every verb on that list must be one that exists *only* in
the singular. Both sentences are regression cases now.

Evidence: four new properties in `core/copy_test.go`, twenty clean apicheck runs
against a fresh campaign, a fifty-three-day paper scanned for all nine known
copy patterns with zero hits, `mise run verify` and `npm test` green.

## The influence system, driven end to end

The three officials are what the endgame is made of and none of it had been
exercised over HTTP. `retain:commissioner`, `retain:mayor`, `retain:editor`,
`puff`, `spike` and `bribe` were all in the harness's untried list. I forced a
save into a state with money, respect and two premises, and drove the lot
through the API.

It works, and the writing is the best in the game:

```
An arrangement with Commissioner Vance: $900 to open it and $45 a day to keep
it. Files go to the bottom of piles.
Page five, with a photograph: A paragraph about a local businessman. Worth 6
respect and 5 off what the police think, and it cost nobody anything.
It does not run: "AUTHORITIES SEIZE THE BLUE HOUR" was set and is not in
tomorrow's paper. As far as this city is concerned it did not happen, and one
more person knows it did.
The arrangement with Editor Sam Rourke ends: You stop paying. He does not
argue, which tells you what it was worth to him.
Somebody at the paper talked: The arrangement at the Herald is over. Nobody
says why, and nobody at that desk will take a call from you again.
```

`puff` really does put a story in the next issue — `LOCAL BUSINESSMAN BACKS
DISTRICT TRADE` appeared in the paper the following day — and `spike` really
does take one out. The whole loop is verified.

**One fault, and it hit all three officials.** `You know Mayor Ellis Crane now:
They answered to nobody.` An official belongs to no organization, and the
function that describes where somebody stands checked the organization first
and returned before it ever looked at the office. A man who is the mayor is
described by the office. There is exactly one of each, so it takes "the": *They
were the police commissioner*, *They were editor of the Bellwether Herald*.

**A correction to my own note.** I had listed `press` with the newspaper actions
in the running plan for the influence system. It is not one — `press` is having
your clothes cleaned and pressed. It belongs with `dress`, and grouping it by
the name alone was a guess I did not check.

Also confirmed in passing, from the same run's paper: `Crated arms are fetching
more than they did` — the plural fix from earlier tonight is holding in a live
campaign.

Evidence: the loop driven end to end through the API against a forced save, one
new property in `core/copy_test.go`, three more patterns on the `cmd/apicheck`
copy scan, fifteen clean runs, `mise run verify` and `npm test` green, `mise run
simulate` unchanged at 52 / 0 / 82 / 0.

## Reading the ledger

The paper has been read end to end twice. The ledger's "What happened" log never
had been, and the player reads it far more often. So I drove a save to day ten
and read all sixty records in order.

It reads well as a story — a crew member sent out and coming back hurt, a
business skimmed until the family notices, a war on the street that eventually
kills you — and it turned up one fault, in forty-nine places at once.

```
Stella Iordan was not as easy as he looked.
$298 off Stella Iordan, and a man with your standing is not somebody anybody
has to describe twice.
```

The city hands out names of every kind. Stella, Franca, Perla, Elena, Mara and
Ida sit beside Nico, Leo and Otto, and **nothing anywhere in the save records a
gender**, because the game has never had a reason to. The prose assumed one
anyway. Two paragraphs further down the same screen, the obituaries were saying
"They were a soldier of Russo Outfit."

That is the argument on its own: the game contradicting itself about the same
person in the same issue. The convention already existed and simply had not been
applied. Forty-nine lines across sixteen files now use it — the crew, the marks,
the officials, the people in a cell, and the player, whose own name is generated
the same way. Where the line was a period idiom about nobody in particular, it
stayed idiomatic without picking a gender: "the street respects somebody who
does their time quietly".

**The check reads the source, not the output.** Every other copy rule tonight
lives in `cmd/apicheck`, but these faults are spread across sixteen files and
most need a state no single campaign reaches — the ledger read found them in
one campaign only because that campaign happened to mug somebody. The test walks
every non-test file in `core` and fails on a gendered pronoun inside a string
literal, naming the file, the line and the text.

And I checked the check, which is the lesson of the last three ticks: I put
"he looked" back into `mugging.go`, confirmed the test failed and printed the
offending line, then restored it. A test that cannot fail proves nothing.

Evidence: `core/pronoun_test.go`, twelve clean apicheck runs, a fresh campaign's
sixty ledger records scanned for gendered pronouns with zero hits, `mise run
verify` and `npm test` green, `mise run simulate` unchanged at 52 / 0 / 82 / 0.
One string in `src/playtest.tsx` fixed alongside.

## Taking a family, and the rest of the gendered prose

The takeover path — going to work for a family, coming up to lieutenant, and
then moving on whoever is in the chair — had never been exercised over HTTP.
Neither had `serve`, `sitdown` or `pact`. I forced a save and drove it.

It works, and the payoff line is one of the best in the game:

```
Vittorio Bellandi is dead: At The Monarch, late, with the city quiet:
Vittorio Bellandi was shot in the doorway while the band kept playing. They
had run Bellandi Family, and somebody who worked for them had come up far
enough to want it.

It is yours: Vittorio Bellandi is dead and Bellandi Family is a name nobody
uses now. You hold four of its premises, eight of its people stayed and two
would not, and every quarrel it was in is yours.
```

The gates are real and behaved correctly: the family has to have fallen to a
third of its peak first, and being refused said so plainly.

**A negative result worth recording.** `serve` looked as though it was missing
from the game — it appeared at no location in the whole payload. It is not
missing. Only the room the player is standing in carries its full action list,
and `serve` is offered at the family's *home*, which is its highest-earning
holding. I nearly filed this as a bug. The rule from earlier tonight held: do
not conclude a cause you have not checked.

**Two copy faults, the second one mine, from the fix for the first.** The
takeover record read `8 of its people stayed and 0 would not` — a zero written
as a figure. I fixed it with the number-speller, and the next run printed `You
hold no of its premises`, because "no" reads in "no people" and never after a
preposition. Both were found by reading the same line twice.

**And the gendered prose, finished.** The pronoun check from the last slice
caught pronouns; it did not catch nouns. `Nobody makes an arrangement like this
with a man`, `3 men` for a crew of three, `MAN CHARGED AFTER DISTRICT SEARCHES`,
`A man was robbed at Saint Agnes`, `You are your own man again`. Twenty-six more
lines across fourteen files. Where a specific person was meant, the city now
names them without guessing; where it was period idiom about an anonymous crowd,
it kept the voice without picking a gender — `Bellandi Family had people
waiting`, `Two of them in a doorway`, `They came armed, damaged your residence
and left`.

**The check cried wolf a fourth time,** on `pluralNames = {"people",
"Brothers", "Boys"}` — a list of name endings the game matches against, not
words it shows anybody. Data that is not prose now carries a visible marker, and
the marker has to be justified where it is written.

**And I nearly trusted a check that could not fail.** I broke `countOf` to prove
the test caught it, saw a pass, and only on looking properly found my edit had
never applied. The second attempt asserted the marker was present before
replacing it, and the test failed exactly as it should. "Prove the check can
fail" means proving the break happened too.

Evidence: the path driven end to end through the API, three new properties, a
fresh forty-day campaign with zero gendered stories and zero gendered ledger
records, fifteen clean apicheck runs, `mise run verify` and `npm test` green,
`mise run simulate` unchanged at 52 / 0 / 82 / 0.

## Reading what the model writes

The director's scenes are the one part of the prose a model writes rather than
the code, and nobody had read a batch of them. I drove a save with Ollama
running and read what came back.

They are good. This is the whole of one, as the player gets it:

```
A supply mission for Bluebird Laundry
"The Bellandi Family needs soap and coal for Bluebird Laundry, and I need you
to fetch it. This isn't just about keeping the business running — it's about
keeping our people employed and our reputation intact."
  Bluebird Laundry · At 15 heat, police may stop completion.
  Bellandi Family standing +6; rival standing −3.
    Fetch what it needs          $95 · 60 min · +4 respect · +3 heat
    Negotiate with the supplier  $80 · 90 min · +4 respect
    Demand the delivery         $115 · 45 min · +4 respect · +8 heat
    Decline the arrangement      No cost or time
```

The model wrote the situation and the three labels. Every figure beside them
came from the core, and the three approaches are genuinely different bargains.

**A negative result, and the fourth time this week I nearly filed a bug from a
partial dump.** My first read printed only each choice's `detail`, saw it empty
on the AI scenes, and I was ready to call it a hole in the game's own principle
that no cost should be hidden. It is not. The terms are carried as figures
rather than prose, deliberately, with the reason written down beside the struct:
"a scene exists to make the player compare two or three approaches, and a run-on
sentence is the one form those numbers cannot be compared in." The view renders
them as a row. Checking the payload and the renderer before concluding took two
minutes and saved a wrong fix.

**One real bug, and it is a bad one.** The validator turned a draft down for
repeating an earlier arrangement, the retry repeated it again, and the game set
the director to `offline` with "Local AI unavailable or proposal rejected."

Two different situations behind one word, and the wrong one of the two: the
model was running and answering. Worse, `offline` is terminal — nothing asks
again. A single repeated title ended the AI for the rest of the campaign. My own
reading script then ran to **day 147 with no scenes at all** and I assumed the
model was broken.

A refused draft now leaves the director available and says what was refused, in
one sentence rather than the paragraph the validator writes for the model.
`offline` is reserved for a model that could not be reached.

Verified both ways against real behaviour rather than only in a test: a live
request still comes back ready, and pointing `BLACK_LEDGER_OLLAMA` at a dead
port gives "The local model could not be reached. Authored play is unaffected."

Evidence: two properties in `cmd/blackledger/director_status_test.go`, both
branches driven live, twelve clean apicheck runs, `mise run verify` and `npm
test` green, `mise run simulate` unchanged at 52 / 0 / 82 / 0.

## What the model actually gets past the guard

The city page is the only prose in the game a model is allowed to touch, and the
log said `City page: rewrite refused` over and over without ever saying which
rule turned it down. That is the same fault as the director's status line from
the last slice: a refusal nobody can read is a refusal nobody can act on. It now
names the rule, and `PolishRefusal` exposes the reason without duplicating the
logic.

Then I measured, which found a real hole.

**The guard checked digits and not words.** A model asked to write like a 1953
city paper writes like one, and city papers spell their numbers. All three of
these went straight into the paper as fact:

```
Seventeen arrests were made in the district overnight, police said.
The forecast says the same again tomorrow, with rain expected for three days.
A dozen shops closed early because of it.
```

The file's own comment says an invented fact "is not a blemish, it is a lie the
player has no way to detect — they would read that eleven people were arrested
and believe it", and that is exactly what was getting through. Spelled
quantities are checked now. "One" is deliberately excluded: in this register it
is almost always a pronoun, and refusing every rewrite containing it would
refuse nearly all of them for nothing.

I found this by writing a test case for the digit rule and picking a bad
example, which passed when it should not have.

**A measurement I nearly published was measuring my own bug.** My first probe
ran the real model over the real briefs and reported three of nine accepted,
five of them refused for coming back empty. That is not the game: the server
sets `think: false` and my probe did not, so qwen3:14b spent its whole budget
reasoning and returned nothing. Matching the server's request:

| of nine real briefs | first probe (wrong) | matching the server |
|---|---|---|
| accepted | 3 | 7 |
| model returned nothing | 5 | 0 |
| came back unchanged | 0 | 2 |

Seven in nine, and the accepted rewrites are better than the originals:
"Additional officers have been assigned to the district. The department will not
say why or for how long."

**And one thing I was wrong about while looking.** Instrumenting the polish path
during an automated run showed 80 of 81 attempts returning immediately because
the model lock was busy, and I was ready to write that the feature never runs.
It is not a fault. The attempt fires on every committed action, one rewrite is
in flight at a time by design, and `cmd/apicheck` commits a whole campaign
faster than one rewrite completes. A long-running server driven at a human pace
logged twenty-seven completed rewrites. The queue is doing what its comment says
it does.

Evidence: two properties in `core/polish_test.go`, the empty-response case now
named separately from a short one, nine real briefs through the live model,
twelve clean apicheck runs, `mise run verify` and `npm test` green, `mise run
simulate` unchanged at 52 / 0 / 82 / 0.

## What survives you

The account abroad is the game's answer to its own hardest rule — everything a
person owns dies with them — and none of it had been driven over HTTP. Nor had
laundering, the still, restocking, inspecting or dismantling. All seven were in
the harness's unreached list.

I drove the whole arc against a forced save, including the part that only
happens once a protagonist is dead.

```
Money leaves the city: $500 sent out, $410 of it arrives. The account holds
$1640 and answers to nobody here, including you if anything happens.
```

Then the protagonist died. The next arrival, Frankie Vale, started with the
usual ninety dollars, and the market offered:

| | |
|---|---|
| deposit | You need $500 to send out at once |
| offshore_access | You need $400 to establish that it is yours |
| withdraw | The account does not answer to you yet |

$1640 sitting there and out of reach, exactly as the file's comment says it
should be: "reaching it as a stranger costs money a stranger does not have."
Given $900, the new life paid the $400, took the money home, and the ledger said
`It comes home: $1640 back in the city and in your hands, where anybody can take
it from you.` The central promise of the system works.

Worth noting because it surprised me: the person who sent the money out has to
pay for access too. That is not an oversight — the deposit record says it at the
time, "answers to nobody here, **including you** if anything happens" — but it
is a sharper rule than I expected and it is stated plainly at the moment it
matters.

Laundering and the rest of the premises work read correctly as well:

```
The books absorb it: $348 through Bluebird Laundry. Police attention falls by
14, to 56. The premises take a little more wear each time, and trade at
Bluebird Laundry is down to 1%.
A still at Bluebird Laundry: $450 of copper and pipe in the back. It makes its
own stock now, and stock has to be moved.
The still comes out of Bluebird Laundry: Copper and pipe out through the back
door. There is less to find here now.
```

**No bugs, which is the result.** Seven systems driven end to end and every line
read: no invented figure, no assumed gender, no count of one taking a plural, no
cost left unstated. After a night of finding a fault in almost everything I
looked at, a system that comes through clean is worth writing down as clean.

Evidence: `deposit`, `offshore_access`, `withdraw`, `launder`, `still`,
`restock`, `inspect` and `dismantle` all driven through the API, the inheritance
verified across a real death, twelve clean apicheck runs, `mise run verify` and
`npm test` green.

## The room nobody wanted to be in

The sitdown is the last big system nothing had driven over HTTP, and the only
thing in the game that can end a war without either side losing it. I forced a
war, gave the player standing with both families, and drove all four choices and
both kinds of evening.

It works, and it is the best-designed thing in here. The scene tells you which
evening you are in before you commit, if somebody warned you:

```
hostility 20  "Bellandi Family and Russo Outfit, in the same room, because you
               asked. Nobody has said anything yet."
hostility 90  "...One of them brought more people than the room needs, and Mara
               caught your eye on the way in."
```

All four choices resolve differently and all four read well. Pressing a room
that came to talk ends the war outright — `It is settled: they shook on it in
front of you, which means it holds as long as you do` — and the conflict list
came back empty. Pressing a room that did not:

```
It was never a meeting: Somebody stood up before anybody had finished a
sentence. You went out through the kitchen with 45 less health than you came
in with.
```

**One fault, and it took reading the line next to the one I was checking.** The
paper carried this:

```
At The Monarch, late, with the city quiet: Ennio Zanetti was knifed in the
crowd. A meeting between Bellandi Family and Russo Outfit went the way
somebody had already decided.
```

The meeting was at Saint Agnes. Casualties are drawn from anywhere in the
organization, and the city reports a death where the person was standing, so
somebody killed at a sitdown was reported dying three streets away — in a room
the player could see they were not in.

Whoever dies at a sitdown came to it, so they are put in the room before they
are killed in it. The killing description is chosen from the location too, so
the fix corrected the prose as well:

```
At Saint Agnes, late, with the city quiet: Ennio Zanetti was held down over a
table until the room emptied.
```

**The test refused to be vacuous twice.** My first version skipped when the room
did not turn, which proves nothing, so it now asserts the fixture is a trap
before it starts and fails if nobody dies. Then I broke the fix to check the
test caught it, and it named the person and the wrong room.

Evidence: all four choices and both evenings driven through the API, one
property in `core/sitdown_test.go` proven to fail when the fix is removed,
twelve clean apicheck runs, `mise run verify` and `npm test` green, `mise run
simulate` unchanged at 52 / 0 / 82 / 0.

## Reading the paper a third time

The last two end-to-end reads found four faults each, so I drove a fresh
campaign to day fifty-two and read it again.

**Every fix from tonight is holding.** No gendered assumption, no plural name
with a singular verb, no numeral in prose, no role without an article, no
price-board name after a unit count. `Somebody was robbed at Saint Agnes`,
`Crated arms are fetching more than they did`, `Two people in the same
organization stood below them`, `They were a landlady`. The English is clean
across two hundred and forty stories.

**One fault, and it is the day-to-day version of one I fixed this morning.**
Within a day the paper collapses repeats. Across days it repeated itself
verbatim:

| ran in seven days | |
|---|---|
| THE CITY COUNTED | 4 times, always "Some 48 people… of whom 18…" |
| PREMISES STANDING EMPTY | 5 times, always the Mariner |
| MORE OFFICERS ON THE STREETS | 5 times |

Day 30's entire issue was two stories, both of them yesterday's. The pick was
randomised per day and the comment above it claimed consecutive days would not
repeat, but nothing ever looked at what had already been printed. Filler that
has nothing new to say has nothing to say: a civic brief now waits three days
before it may run again.

**I broke the fix while making it and caught it by re-reading my own code.**
Rewriting the loop to skip recent briefs, I dropped the guard that stopped a
page running the same item twice in one issue. Walking the candidate list once
from a per-day starting point does both jobs.

**And one repeat survived, which found a second site.** After the fix, fifty-
eight days produced exactly one repeat: `BELLANDI FAMILY SAID TO BE STRUGGLING`
on consecutive days. That story is filed by the family-fortunes code, not the
city page, so it never saw the new rule. A family's fortunes are news; the same
sentence about them two days running is not.

Measured on a fresh fifty-nine-day campaign after both fixes: **zero civic
repeats inside the rest period, zero empty issues, median three stories an
issue.**

Evidence: one property in `core/citypage_test.go`, proven to fail when the rest
rule is removed (it names the headline and the two days), twenty-two clean
apicheck runs, `mise run verify` and `npm test` green, `mise run simulate`
unchanged at 52 / 0 / 82 / 0.

## The guide cannot tell you something the rules do not, except once

The Guide's own header promises that it "asks the game the same question the
buttons ask, so it cannot tell you something the rules do not". Reading it in
the browser for the first time, it lives up to that almost everywhere: every
locked step gives the game's own reason for being locked, and the thresholds in
`Rules that do not change` are interpolated from the constants the police use,
so they cannot drift.

Almost. The standing it takes to buy premises was written out three times as a
bare `6` — in the rule, in the opportunity that suggests it, and in the guide
line itself. Any one of the three could have changed and the page would have
gone on saying six. It is one constant now, and the test breaks if the guide
states a figure of its own: hardcoding `Earn 6 respect first` back while moving
the rule to nine makes it fail with *"the guide says "Earn 6 respect first"
while the rule uses 9"*.

That was the only figure in the file not derived from a constant.

## What the director actually gets past the validator

Nobody had taken this number, and it was not possible until a refused draft
stopped killing the director for the rest of the campaign. Eighteen requests
against the live model on a mid-campaign save:

| of 18 requests | |
|---|---|
| produced a usable scene | 0 |
| refused by the validator | 11 |
| lost to the world changing while the model wrote | 7 |

**The seven are largely my harness and I will not claim otherwise.** My loop
advances the clock to make a prepared scene arrive, so the world moves while the
model is writing. In real play the clock is paused while the player decides,
which is exactly when the director is asked, so context changes are much rarer
than this makes them look.

**The eleven are real, and they are not random.** The validator's reasons, from
the server's own log:

```
6x  body must explicitly name the selected venue "Bluebird Laundry" and
    describe the task there
2x  body includes an unsupported deadline "before the end of the day"
2x  body assigns this job's people to a family through "Bellandi men", but the
    work is credited to Nico Ward's people
1x  director context changed during preparation
```

The dominant failure is the model not naming the venue in the dialogue. The
prompt does ask for it — `Name that venue accurately in the dialogue` — on line
22 of a 32-line brief. The validator is right to refuse; the instruction is
buried.

**And the correction pass is weaker than it looks.** A rejected first draft is
retried once with the reason fed back. Measured over the same run: nine first
drafts rejected, eight still rejected after the retry. **The retry rescued one
in nine.**

I am not changing the prompt on this evidence. Moving an instruction and
declaring victory without a controlled before-and-after is the mistake I made
twice tonight measuring my own harness. What is worth recording is the shape:
the model fails the same way repeatedly, the failure is specific and stated, and
the retry barely helps. A campaign earlier tonight on a simpler save did get
scenes through, so this is a rate on one state, not a universal one.

## The ledger was eating its own history

I read the ledger end to end a second time. Every fix from the first read is
holding: `They are already at Saint Agnes`, `They are available again`, `Stella
Iordan helped themselves to $74` — no assumed gender anywhere in sixty records.

Then I counted, and found something worse than a copy fault.

| campaign length | records held | distinct | exact repeats |
|---|---|---|---|
| day 12 | 60 | 37 | 38% |
| day 28 | 60 | 44 | 26% |
| **day 58** | **60** | **10** | **83%** |

At day fifty-eight, **forty-nine of the sixty records the ledger holds were the
same sentence**: "Envelope delivered — Mara pays $45. A small favor, completed
without questions." Repetition had not merely made the log unreadable. The
archive is bounded, so the repeats had pushed every notable thing that ever
happened out of it. A campaign with a war, three deaths and a family collapse in
it had a permanent record consisting of one errand, forty-nine times.

The paper learned this within a day, hours earlier this morning. The ledger had
it across a whole campaign and nobody had counted.

The same thing happening again today now collapses onto the record already
there, with a count. Two details mattered:

- The record is **removed and re-appended** rather than updated in place. The
  result panel after every action is built by diffing record ids, so a
  collapsed repeat has to take a fresh id or the player commits an action and
  is told nothing happened.
- A different day, a different life, or a different amount is a different
  event, and stays one.

| after, same seeds | records | distinct | exact repeats |
|---|---|---|---|
| day 12 | 60 | 53 | 11% |
| day 58 | 60 | 39 | 35% |

Thirty-nine distinct records over fifty-eight days instead of ten. The
remaining repeats are across days, which is correct: the count is per day, and
`x32` on one day beside `x31` on another is two true facts.

The Ledger screen shows the count beside the title.

**One test had to change and it is the same shape as this morning's.** A test
filled the log with two hundred copies of one line to force a rollover. Two
hundred copies is now one record, so it fills nothing — it uses distinct lines,
which is what filling a log means.

Evidence: three properties in `core/ledger_test.go`, one of them holding the
line that a collapsed repeat must still read as a new outcome, fourteen clean
apicheck runs, `mise run verify` and `npm test` green, `mise run simulate`
unchanged at 52 / 0 / 82 / 0.

## Counting the paper, and a correction I made to myself before publishing it

The ledger was fixed by counting it rather than reading it, so I counted the
Herald the same way. Fifty-eight days, 198 stories:

| | |
|---|---|
| civic filler | 117 of 198 (59%) |
| issues carrying no news at all | 26 of 58 |
| real stories per issue, median | 1 |

Days 29 to 43 were a near-unbroken run of issues containing nothing but the
weather and a standing notice about an empty building. My first conclusion was
that the living world goes quiet once one organization is on top — in that
campaign the player's people were at power 100 holding three of the best
premises, every conflict was theirs, and there had been one death in fifty-eight
days.

**That conclusion was wrong, and I caught it by running the right instrument
before writing it down.** `mise run simulate -runs 20 -steps 1000` puts twenty
campaigns past the horizon where the living-world measures mean anything:

| across 20 campaigns, ~85 game days | |
|---|---|
| wars started | 70, in 17 of 20 runs |
| families created | 19, in 14 runs |
| families destroyed | 17, in 12 runs |
| holdings changed hands | 29, in 16 runs |

The city is not degenerate at length. The quiet campaign was one dominated city,
not a dead simulation, and I would have published the opposite from a single
run. That is the fourth time tonight a measurement of mine was really a
measurement of one harness's play.

**What is true is that the paper has too few things to say.** The city has kept
hours since this morning — people are at their posts through the morning and in
the bars and clubs after midday — and the paper had never once mentioned it,
which for a paper printed in this city is a strange thing to miss. It now
reports the evening when the evening is worth reporting, from a real count taken
at midnight before anybody sets off for the day shift.

**And an honest negative result: it did not fix the filler share.** After
adding it, the same fifty-eight day campaign is 60% civic — the new brief simply
takes a slot another one would have had. The ratio is set by how much news the
city makes, not by how many things the page can say, and adding briefs will
never change it. Worth knowing before anybody adds more.

My own test caught a copy fault in it before it shipped: the first version
printed "fifty of the district were in one of them", a figure in prose, which
this paper does not do. It describes the share instead.

## Unresolved: a sitdown offered and then refused

`cmd/apicheck` now reaches the sitdown, and on two runs out of twenty-two it
reported:

```
HTTP 409 on sitdown: {"error":"One of them would not sit in a room you arranged"}
```

The game listed the action as available and then refused the command. That
breaks the rule the interface runs on — the reason a thing cannot be done is
supposed to be knowable before committing to it.

`CallSitdown` runs **after** the clock advances, deliberately, "so the evening
actually costs the evening", so three hours pass between the readiness check
that enabled the button and the check that refused it.

I could not reproduce it. Four hundred sitdowns in a single-quarrel fixture
never refused; three hundred with a second quarrel whose family would not sit
with the player never once saw the worst quarrel change under the player during
those three hours. Both hypotheses are ruled out, and I am not shipping a fix
for a cause I have not found. Recorded here with the reproduction that does
work: drive a fresh save with twenty-two `apicheck` runs and watch for a 409 on
`sitdown`.

## The meeting that was called off after the player had paid for it

Yesterday I recorded a bug I could not find. `cmd/apicheck` reported, twice in
twenty-two runs against a fresh save:

```
HTTP 409 on sitdown: "One of them would not sit in a room you arranged"
```

The game listed the action as available and then refused the command, which
breaks the rule the whole interface runs on: the reason a thing cannot be done
is knowable before you commit to it.

**Three guesses failed, across seven hundred attempts.** Goodwill drifting
during the three hours the room takes: four hundred sitdowns in a single-quarrel
fixture, never refused. The worst quarrel changing under the player: three
hundred with a second quarrel whose family would not sit, never once. And a
third I had not noticed was impossible — my fixture starts at minute 480, so a
180-minute advance never crosses the twelve-hour boundary where organizations
reconsider each other, and could not have reproduced anything.

**Then I stopped guessing and made the refusal say which door was shut.** One
run:

```
Doyle Crew would not sit in a room you arranged.
They think of you at -26 and it takes -25.
```

One point. The families in the failing save were sitting exactly on the
boundary, which no fixture had put them on, and something in those three hours
moved one of them a single point. The first hypothesis was right; my test could
never have produced the value it needed.

This is the third time tonight that a message which could not tell two cases
apart was itself the bug. The director said "offline" for a model that was
answering; the city page said "refused" without saying by which rule; and this
said "one of them" when it knew the name and the number.

The room is arranged when the player commits and the guarantees are paid. Three
hours later it opens, and the meeting that opens is the one that was agreed:
`CallSitdownAs` takes the quarrel captured before the clock moved. The evening
still costs the evening, which was the point of advancing first.

**Verified with the instrument that found it.** Sixty-six `apicheck` runs across
three fresh saves, where twenty-two used to produce two failures: clean. Two
properties in `core/sitdown_test.go`, one holding that an agreed meeting is not
called off by a point, one holding that a real refusal names the family and the
figure. `mise run verify`, `npm test` green, `mise run simulate` unchanged at
52 / 0 / 82 / 0.

## An action the game offers is an action the game accepts

The sitdown fault was not going to be unique. Every action whose work is done
after the clock advances can refuse for a reason that became true during the
hours it took, and there are ten of them in the command path: `sit_out`,
`lawyer`, `talk`, `post`, `unpost`, `operate:*`, `incite`, `move`, `sabotage`
and the sitdown.

So I wrote the rule down as a property instead of chasing them one at a time.
`TestEveryOfferedActionIsAccepted` walks forty campaigns through varied money,
standing, health and clock positions, and for every action the game lists as
*available* it executes that action and requires it to succeed. **2,264 offers
across 36 distinct actions, all accepted.**

**And then it failed to catch the thing it was written for, which is worth more
than the test.** I reverted the sitdown fix and ran it: green. Random states
almost never sit a value exactly on a threshold, and that is precisely where
this class lives — a family at −26 against a limit of −25.

So I added a second pass that puts the world *on* the boundaries: every faction
at exactly the standing a sitdown needs, cash at exactly the fee, respect at
exactly the requirement, and the clock started so the hours an action takes will
cross the point where organizations reconsider each other. Reverted the fix
again: **still green.**

I could not make a fixture reproduce it. The one-point drift needs a live
campaign's accumulated history — business demands, retaliation, standing moved
by things the player did days ago — and a world built at rest does not have any.

**So the honest account is this.** The new test asserts a real invariant over
2,264 offers and would catch a new action added with a mismatched gate, which is
the common form of this fault. It does **not** cover the timing-drift form. The
only instrument that has ever caught that is `cmd/apicheck` driving a real save,
which found it originally and now runs clean over eighty-eight runs across four
fresh campaigns.

A test that passes when you break the code is not evidence, and saying which
part of the class it covers is the difference between a guard and a comfort.

Evidence: `core/offered_test.go`, twenty-two more clean apicheck runs, `mise run
verify` and `npm test` green, `mise run simulate` unchanged at 52 / 0 / 82 / 0.

## Making the only instrument that catches it keep the evidence

`cmd/apicheck` is the only thing that has ever caught the timing-drift class —
an action the game lists as available and then refuses. It found the sitdown by
accident, and its report carried the message and nothing else. By the time
anyone read that report the state which produced it was gone, which is why
finding the cause took three failed hypotheses and seven hundred attempts.

A refusal now brings its own evidence: the action as the game offered it
(including whether it was listed as available and any reason shown), the minute,
the location, cash, health, attention, respect, and **the player's standing with
every organization** — which is the thing that moved under the sitdown and the
one thing the report would never have shown.

**I could not verify it against the real fault, and say so.** I reverted the
sitdown fix and ran sixty-six `apicheck` runs across three fresh saves to make
the 409 happen again with the new capture in place. It did not recur — the
original was two occurrences in twenty-two runs on one save's particular
history, and fresh campaigns did not reproduce it. So the end-to-end capture of
a live 409 is unproven.

What is proven is the capture itself, directly: `asOffered` is given a snapshot
where `sitdown` is enabled and `rob` is disabled with a reason, and one family
sits at −26. It records the enabled action as offered, the disabled one as
disabled with its reason, the world around both, and the standing. Blanking the
standing line makes the test fail with *"standing was not recorded"*.

That is a smaller claim than "the report now explains the bug", and it is the
one I can stand behind. The next time this class fires, the evidence will be in
the report instead of in a state nobody can get back.

Evidence: one property in `cmd/apicheck/copy_test.go`, proven to fail when the
capture is removed, fourteen clean apicheck runs, `mise run verify` and `npm
test` green, `mise run simulate` unchanged at 52 / 0 / 82 / 0.

## Reading two screens nobody had opened

The Market and Settings screens had never been looked at. The Market is clean:
three goods, each with the price, how far it has moved from usual, what is held,
and where it trades. Nothing to fix.

The Settings screen verified a fix and then showed me a fault I had put there
myself this morning.

**The verification.** The live save was still displaying `Local AI unavailable
or proposal rejected` — the message replaced hours ago, kept because the
director's detail is stored state and that campaign had not asked since. Pressing
*Prepare an encounter* replaced it, and the new behaviour worked end to end in
the browser for the first time: the status came back **available** rather than
the terminal `offline`, with a reason and an invitation to try again. That fix
had only ever been checked over HTTP.

**The fault.** What it actually said was:

```
The last draft was turned down: This work exists only because of a committed
event, and the dialogue never refers to it. You can ask for another.
```

That is the validator's own prose, written to instruct a model, shown verbatim
to a person who cannot act on it. I piped it there this morning while fixing the
message above it.

The distinction that matters to a player is the one the status already carries:
the model could not be reached, or it answered and what it wrote could not be
let in. Which rule caught it is a developer's question, and it is already in the
server log — where the last five faults were found. The screen says what
happened in the game's voice now, and the log keeps the detail.

**And the test I wrote for it was case-blind.** It scanned the player-facing
line for words like "proposal" and "dialogue", but `firstSentence` capitalises,
so restoring the old message left the test green. Lower-casing before the scan
makes it fail properly: *"the player is shown the validator's own words"*. Two
ticks ago the break did not compile; this time it compiled and the check was
looking for the wrong string. Confirming the break applies is not the same as
confirming the check can see it.

Evidence: `cmd/blackledger/director_status_test.go`, proven to fail when the
validator's words are restored, twelve clean apicheck runs, `mise run verify`
and `npm test` green, `mise run simulate` unchanged at 52 / 0 / 82 / 0.

## The scene, and a button for stopping something that was not happening

The scene modal had never been looked at. Getting one on screen took forcing the
state: no authored encounter fired in ninety actions of ordinary play, so I gave
the player a business and set the next pressure to the current minute.

It reads very well — portrait, name, office, the demand in the speaker's own
voice, three numbered choices each stating its cost and its consequence, TIME
PAUSED in the corner and "Decisions are saved immediately" at the foot. Nothing
wrong with any of it.

Underneath the dialogue were two buttons: **Read aloud** and **Stop voice**,
always, both. A control for stopping something that was not happening, sitting
next to the control for starting it.

**I nearly filed a second fault that was not one.** The label beside the speaker
looked like dead state — `setSpeech` appeared nowhere in the file. It is passed
by reference as the voice player's `status` callback, so the label is live and
says "Preparing voice…", "Speaking…", "Replay voice" as it goes. Reading the
wiring before writing it up saved a wrong fix, which is the fifth time that has
happened tonight.

That live status is also the fix. The player already reports what it is doing,
so `speaking()` moved into `voice.ts` beside the states it reads, and the stop
appears only while a reading is being prepared or is playing. Putting it there
rather than in the view means it cannot drift from those strings without the
existing voice tests failing too.

The test walks a real reading: idle, preparing, speaking, ended, and failed.
Making `speaking` return true always fails it with *"stop offered before
anything was asked for"*.

Also confirmed while there: the voice toggle in Settings governs whether a scene
reads itself aloud on opening, not whether the button exists. Offering a manual
reading with auto-play off is deliberate, and correct.

Evidence: one test in `tests/voice.test.mjs` proven to fail when the predicate
is broken, the modal re-read in the browser with the stop control gone, twelve
clean apicheck runs, `mise run verify` and `npm test` green, `mise run simulate`
unchanged at 52 / 0 / 82 / 0.

## The death screen, and the promise the whole game rests on

The death screen is the one thing everything else builds toward and nobody had
looked at it. I drove a campaign to a protagonist worth killing — day 30, three
premises, six people, 481 respect — and then killed him properly, on the street,
in somebody else's war.

It reads as well as anything in the game. `THE CITY CONTINUES` above the fold,
the name and the cause in two lines, then what he lived to, what he was worth
and what he earned; what became of the thing he built; the four headlines the
paper carried; what the next one gets; and a single way forward.

**And it verified the Guide's central promise for the first time.** The Guide
says *"Death is permanent. What you built passes to the strongest of your own
people and becomes an organization you can deal with, or fight."* That is
exactly what happened: a new faction appeared, `Franca Sabbatini's people`,
power 70, holding two of the three premises, led by one of the dead man's own —
and the ledger said so in its own words:

```
Franca Sabbatini is running it now: Everything that answered to Nico Ward
answers to Franca Sabbatini by the end of the week, under a name Nico Ward
never chose.
```

**Three things looked like faults and none of them were.**

The epitaph showed the *previous* protagonist. That was my forcing method: I had
set `alive=false` in the save, which never runs the death path that builds it.
Killing him through the game produced the right one.

The `standing` list — what is left standing in the dead man's name — was empty
for a man with three premises. It is correct: somebody inherited them, so
nothing stands in his name. The view renders that list when it is not empty, and
I checked before writing it up.

And the dead protagonist's own organization is still in the faction list at full
strength holding nothing, dissolved only when the next life begins. It is never
seen, because the death screen blocks every other screen until the player begins
again.

That is three near-misses in one slice, and seven tonight. Reading the wiring
before writing up a fault is now the most valuable habit of the night.

**What was actually missing is a test for the other branch.** Dying with nobody
to inherit says *"They had nobody. Whatever they held stands with no one to
answer for it"* — and if the list behind that sentence were ever empty, the
sentence would be a claim with nothing behind it. There is a property for it
now; emptying the list makes it fail and prints the whole epitaph.

Evidence: `core/epitaph_test.go`, proven to fail when the list is removed, the
screen read in a browser against a real death, twelve clean apicheck runs, `mise
run verify` and `npm test` green, `mise run simulate` unchanged at 52/0/82/0.

## Do not call a room the player is not standing in "the room"

The first thing a new player ever sees had never been read. I started a fresh
life on 8894 and looked: day 1, 08:00, Alex Varga, $90, standing in The Mariner,
with the opening opportunity pointing at Saint Agnes. The sidebar showed Saint
Agnes, which is right — the panel follows the selected place, not the player.
Its people section was headed "IN THE ROOM · 2 of 4 you can deal with", which is
not. The player is not in that room. The same heading also promised "2 others in
the room" behind the expander.

The wiring says why. `ActionList` printed "In the room" as fixed text and had no
way to know whose room it was drawing, while its parent in `src/main.tsx`
already computes `l.id === p.location` for the eyebrow above it — that is how the
panel says "YOU ARE HERE" or "NEIGHBORHOOD DIRECTORY". The knowledge existed one
level up and was not passed down.

`ActionList` now takes an optional `here`, defaulting to true, and `main.tsx`
passes the comparison it already makes. A room you stand in still reads "In the
room" and "others in the room"; a room you are only looking at reads "Who is
there" and "others there".

Evidence: the opening read in a browser on a fresh save, before and after — the
Saint Agnes panel now reads "WHO IS THERE · 2 of 4 you can deal with" and "Show
2 others there" while the player stands in The Mariner. Twelve apicheck runs, no
invariant failures. `mise run verify` and `npm test` green, `mise run simulate`
unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

Not changed, and recorded as a judgement rather than a fault: the opening ledger
record carries the game's premise ("A room. A name. No protection.") and the
City screen does not show it, because the clock has not moved and the banner
correctly reads "Nothing yet. The clock is paused." Where the premise line
belongs is a design decision, not a bug.

## Hire the person the button named

Twelve apicheck runs against one save reached day 58 and stopped teaching me
anything: the last five runs spent 53, 73, 68, 60 and 76 of their eighty
commands on `courier`, standing in the same bar. Reading that room's action list
is what found this. Two refusals sat next to each other and contradicted each
other. `delegate` was disabled because "Leo Carver is dead". `recruit` was
disabled because "Leo is already in your crew". Both cannot be true, and killing
a man takes him off the crew list, so he had been hired after he died.

The recruit button had already been taught that nobody drives forever. Its label
names whoever holds the job, and the reachability sweep is asked about that same
person, which is what stops it offering to hire a corpse. The effect was never
taught: it appended a hardcoded `Crew{"leo", "Leo Carver", 65}` whatever the
button said. So once the city gave the wheel to somebody else — which it does on
the daily tick, by design — the label offered the new driver, the gate cleared
the new driver, and the game put dead Leo Carver on the books. The refusal beside
it was hardcoded the same way and printed Leo's name over whoever was actually
in the crew, which is how a crew of one dead man went unread for fifty-eight
days.

Both now name the person. `recruit` hires `w.Holder("driver")` and refuses if
nobody drives; the crew refusal prints the name of whoever is on the books.

Evidence: `core/recruit_dead_test.go` states the property that the person hired
is the person the button named, and failed before the fix with the exact shape
of the bug — "the button offered Zoltan Esposito and the game hired Leo Carver
(leo)". Verified over HTTP on the day-58 save that produced the fault: with Leo
dead and the crew emptied, the bar offered "Recruit Ugo Lenz", the command hired
Ugo Lenz, the refusal then read "Ugo Lenz is already in your crew", and
`delegate` came back enabled. `mise run verify` and `npm test` green, `mise run
simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors.

Recorded and not fixed, because it is the harness and not the game: apicheck
checks `earning` before travel, so a late campaign standing in a room that
offers `courier` forever never leaves it. That is why coverage collapsed. It is
also why the fault was found, so the ordering stays for now.

## A seat across from whoever runs the family now

The hardcoded hire suggested a class rather than a bug: a label taught something
its effect was not. Auditing for the same split found the mirror image of it —
two scenes whose *speaker* was hardcoded while everything around them asked the
city who held the job.

Every other authored scene names its speaker by role. The fixer's offers use
`w.HolderID("fixer")`, the police stop uses `w.HolderID("detective")`, and both
follow the job when the person holding it dies. The audience and the business
demand named `"vittorio"` and `"elena"` — the two people who happened to lead
the two families on the first morning. Families change hands constantly; that is
most of what the living world does. Across twenty long runs it destroys
seventeen families and creates nineteen. So the moment a leader was killed and
the strongest survivor took over, the successor's demand arrived in a dead
predecessor's voice, and the player was seated across a table from a corpse.

The demand was worse than the audience. It picked its speaker by *position* in
the faction list, `actor == 1`, so the identity was not even wrong in a stable
way.

Both now ask `w.Leader(faction)`, which already existed and already skips the
dead. That opened a second question the hardcoding had hidden: a family can lose
everybody. The view falls back to the first person in the city when it cannot
find a scene's speaker, so an audience opening with no leader would have seated a
stranger and let them set a family's terms. Neither scene opens now without
somebody to speak, and `AudienceReadiness` disables the action with "Nobody is
left to speak for the Bellandi Family" rather than spending the player's evening
on nothing.

Evidence: `core/speaker_test.go` states that the person in the chair is the
person who runs the family. Both cases failed before the fix with the shape of
the bug — "Rosa Marchetti leads bellandi now, but the chair holds Vittorio
Bellandi (dead: true)" — and the audience case was re-broken afterwards to
confirm the check can still see it. A third test covers the family with nobody
left. Verified over HTTP on a driven save: with Vittorio dead and Rosa Marchetti
leading the Bellandi Family, requesting an audience at the club returned a scene
spoken by `heir`, Rosa Marchetti. `mise run verify` and `npm test` green, `mise
run simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors.

Ninth near-miss, recorded so I do not chase it again: in that HTTP check Rosa's
role still read "Soldier" rather than "Head of the Bellandi Family". That is my
forcing method, not the game. I edited the save directly instead of going
through `Succeed`, which is what sets the title.

Still unexamined in this class, and worth the next pass: the view's silent
fallback from a missing speaker to `world.npcs[0]`. Core no longer emits one, so
nothing reaches it today, but a fallback that turns "nobody said this" into "a
specific real person said this" is the wrong shape for a game whose core is the
only source of truth.

## A scene nobody said is attributed to nobody

The previous slice guarded core against emitting a scene with no speaker, and
noted the reason: the view resolved its speaker with
`world.npcs.find(...)||world.npcs[0]`. An id the city did not know became
whoever happened to be first in the list — a real person, with their portrait,
their name and their job, delivering a family's terms they never set. Core is
the only source of truth for who spoke. A fallback that turns "nobody said this"
into "this specific person said this" is the view inventing a fact.

`speakerOf` in `src/voice.ts` returns the person or null, and the scene panel
renders the attribution block only when there is somebody to attribute it to.
The body, the choices and the voice control are unchanged: an unattributed scene
still plays, it just does not borrow a face.

Evidence: a test in `tests/voice.test.mjs` asserts null for an unknown id, an
empty id, a missing id and an empty city. Re-adding the `??people[0]` fallback
makes it fail — confirmed, not assumed. Read in a browser against the forced
succession save from the previous slice: the audience at the club renders "A
seat across from Bellandi" with Rosa Marchetti's portrait, name and job above
the body, and the two choices the player cannot afford carry their own refusals
("Not enough cash: this takes $150 and you are holding $75"). `mise run verify`
green, `npm test` 25 pass 0 fail.

## Whose street this is, decided by ground rather than by slice position

Continuing the audit into positional indexing found a crash, not a wording
fault. The business demand chose which family was collecting with two literal
indices: `Factions[0]` for the player's first district, `Factions[1]` for
anything beyond it. `Dissolve` removes a family from that slice, and destroying
families is most of what the living world does — twenty-one of them across
eighty long campaigns in the run below. So those indices stopped meaning the
families they were written for as soon as the city did its job, and once one of
the two originals fell, `Factions[1]` was not an index into anything. A player
who owned an earning business outside their first district and had seen a family
fall crashed the collection. Nothing in the server recovers, so the request died
mid-commit.

`Claimants` replaces both indices, and has to hold two things at once. Ground is
the honest answer to whose street this is, so a family holding premises in a
district collects there. But separate districts need separate claimants, or
buying one family off would silence every demand in the city. That second
property is not mine — it is already stated by three tests written long before
this change, and the first version of the fix broke all three by handing every
district to whoever was strongest. Ground alone is not enough because the
opening city is a two-two tie in the first district and nobody at all holds the
second. So a district nobody holds goes to the strongest family not already
collecting from the player, and only falls back to one that is when there is
nobody else.

Evidence: `core/pressure_family_test.go` reproduces the crash — "the city
panicked collecting from apartment: runtime error: index out of range [1] with
length 1" — and still does when the positional version is put back, confirmed
by re-breaking it. Two further tests state that a family taking a district over
inherits the claim, that no families means no demands, and that districts get
distinct claimants. The three pre-existing truce and territorial-claim tests
pass unmodified; that they pass is the evidence the design property survived,
and I did not touch them.

Verified over HTTP on a driven save with Bellandi dissolved, one family left,
and the player owning the garage beyond their first district: the collection
that used to crash produced "A claim on your earnings" from the Russo Outfit,
spoken by Elena Russo.

Measured before and after over eighty long campaigns on the same seeds
(`-runs 20 -steps 1000`), so the political layer can be seen not to have moved:

| measure | before | after |
| --- | --- | --- |
| wars started | 106 | 112 |
| holdings changed hands | 50 | 52 |
| families created | 38 | 39 |
| families destroyed | 18 | 21 |
| deaths | 20 / 0 / 20 / 0 | 20 / 0 / 20 / 0 |
| errors | 0 | 0 |

`mise run verify` and `npm test` green, `mise run simulate` unchanged at defiant
52 / investor 0 / reckless 82 / worker 0, 0 errors.

Correcting an earlier record: the living-world baseline written up as "70 wars,
19 families created, 17 destroyed, 29 holdings changed hands across 20 runs" is
not comparable to the table above. That figure counted a subset; the run here
aggregates all eighty campaigns, twenty per strategy. The numbers are the same
world measured differently.

## The city must not forget who the player is talking to

Sweeping the same class further — a lookup whose result is used without asking
whether it found anything — found one live crash and two near-misses. The
near-misses first, so they are not chased again. `migrate.go` guards its
`Factions[0]`/`Factions[1]` pair with `len(w.Factions) >= 2`, which is correct.
`yourpeople.go` reads `w.NPC(id).Trust` and `n.Faction` straight off the lookup,
but every caller passes through `LetGoReadiness`, which returns a refusal when
the person is nil. Both fine.

The real one is `PrunePeople`, which keeps a long save bounded by forgetting the
dead once nothing refers to them. What counts as a reference was a list, and the
list did not include the person standing in front of the player. A scene names
its speaker by id and by nothing else. So: a proposal is open, the person who
made it is killed, the prune forgets them, and declining the offer prints their
name straight off the lookup — `w.NPC(e.Speaker).Name` on nil. Nothing in the
server recovers. The same gap covered a suspended arrangement, which remembers
who it is with the same way.

Both are references now. The refusal is also written so it does not need the
lookup at all: a save from before this change can already have lost a speaker,
and no migration puts somebody back, so it reads "You decline the proposal" when
there is nobody left to name.

Evidence: `core/prune_test.go` reproduces the crash — "declining crashed the
city: runtime error: invalid memory address or nil pointer dereference" — and
three tests state the properties. Each was re-broken afterwards: removing the
scene-speaker reference fails the prune test, and restoring the unguarded
lookup fails the refusal test. `mise run verify` and `npm test` green, `mise run
simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors, and sixteen apicheck runs against one save with no invariant failures.

Two negative results, recorded as results.

The loan book does not need protecting here, and the first version of that test
was measuring its own setup. `Kill` writes off a current-life loan at the moment
of death, and a loan carried over from an earlier life is already read through a
nil check in both `LoanDay` and `LoanDescription`. I removed the clause I had
added for it.

And this crash is not one a player hits today. `PrunePeople` only runs once the
city holds more than 150 people. Driven to day 47 over sixteen apicheck runs,
the city held 52 living people and 3 dead, and that number is stable rather than
growing — the population settles around fifty and stays there. So the prune has
almost certainly never run in any campaign anybody has played, and the crash is
reachable only in a campaign far longer than any I have driven. The reference
list was incomplete either way, which is the reason to fix it; the severity is
not what it first looked like.

## A refusal that knows where the hand is should say so

Read the six rooms I had never read side by side, driven to day 18 on a fresh
save: garage, docks, market, club, precinct, casino. The pair that did not agree
was the tables. The club and the casino both refused to deal, both with the same
sentence — "There is a hand on the table already" — and only one of them offered
anything to do about it. The hand was at the club, so the club showed "Take
another card" and "Stand on 12" beside its refusal. The casino showed the
refusal alone, and did not say where the hand was, while the world had been
storing the room on the hand the whole time.

The refusal now carries it: "There is a hand of yours still on the table at The
Monarch" when the player is somewhere else, and "You are in the middle of a
hand" where the controls to finish it are already on screen.

Evidence: `core/table_elsewhere_test.go` states both halves and failed before the
fix with "the refusal does not say where the hand is". Re-broken afterwards to
confirm the check can see it. Verified over HTTP on the day-18 save that showed
the problem: standing in The Blue Hour, the small tables refuse with "There is a
hand of yours still on the table at The Monarch". `mise run verify` and `npm
test` green.

Qualifying this honestly, because the browser said something the API did not.
The player was not without the information. Reading the same save on screen, the
player card in the sidebar already read "At the tables in The Monarch: showing
12, dealer shows 10, $50 down". So this is two parts of one screen disagreeing
about how much they will tell you, not a fact the game was withholding. The
refusal is where a player looks to find out why a button will not work, and it
now agrees with the card above it.

Also noted while reading, and deliberately not changed, because they are
questions about design rather than faults. The crew is named inconsistently
across one room: `rob:crew` reads "Send Leo Carver for the till" while
`delegate` reads "Send Leo on collections" and `crew_bonus` reads "Pay Leo a
bonus". And `sabotage` names the family it hurts while `sabotage:crew` names the
premises it hits — "Move against Russo Outfit yourself" beside "Send Leo Carver
against Russo Motor Works" — so a player cannot tell from the labels whether the
two do the same thing. Both are worth settling, and settling them is a choice
about what these actions are, not a correction.

## One crew, one name, and a button that said it sent somebody it did not

Reading the last six rooms side by side — market, laundry, herald, apartment,
room, estate — turned up two things in a single room, and one of them corrected
a fix I had just made.

**The crew was two people.** On a driven save the market showed "Send Bela Havel
for the till", "Send Bela Havel after Elena Russo" and "Send Bela Havel against
Mercer Exchange", and directly beside them "Send Leo on collections", "Pay Leo a
bonus", and a refusal reading "Leo refuses assignments below 30 loyalty". Bela
Havel was the crew. Leo Carver was a name written into the strings. Three
buttons in that block ask the city through `CrewHands`, and one line of the
collections description already reads `n.Name`; four others did not. Fixing the
hire to take whoever actually drives is what made this visible — before that,
the crew was always Leo, so the hardcoding never showed. Three more sites
carried it: the task list read "Leo · collections", and two ledger entries read
"Leo returns". Those two only fire when the person is gone, so they now say "The
round is finished" rather than naming somebody who is not there. The opening
opportunity named Leo too, and now names whoever drives.

**A button said it sent somebody it did not, and my first correction was also
wrong.** `sabotage` was labelled "Move against Bellandi Family yourself" and
described as "Send your crew against The Monarch" — the other button. I read the
label, the "a failed attempt injures you" clause, and `Sabotage` calling
`SabotageBy` with `OwnHands`, concluded the player goes alone, and rewrote the
description to say so. That was wrong. `SabotageReadiness` requires a crew on
both halves and refuses on their loyalty on both, `sabotageChance` reads that
loyalty on both, and the failure text settles it: this half says "You and Leo
left without reaching anything" while the other says the crew "went in without
you". The player never goes alone. The difference is whether they go at all. The
description now reads "Go in with your crew", and the test I wrote had to be
narrowed with it — forbidding "your crew" in the self description would have
been forbidding the truth.

The pair also now names the same target on both halves. Robbery and mugging each
name one subject on both sides — the till, the person — so the player reads them
as one choice about who carries it out. Sabotage named the family on one and the
premises on the other, which reads as two different acts. Both name the premises
now; the family is still in the description, where the other two pairs keep the
same kind of detail.

Evidence: `core/own_hands_test.go` holds three properties. Fifteen self-and-crew
pairs are checked across the whole city and only sabotage failed; both halves
must name the same target; and one crew must be called by one name. All three
failed before the fixes with the exact strings above, and re-breaking the code
fails them again. Verified over HTTP on the save that showed the fault: with
Bela Havel in the crew, the room reads "Send Bela Havel on collections", "Pay
Bela Havel a bonus", and "Bela Havel refuses assignments below 30 loyalty".
`mise run verify` and `npm test` green, `mise run simulate` unchanged at defiant
52 / investor 0 / reckless 82 / worker 0, 0 errors.

Thirteenth near-miss, and the first where I nearly published a wrong correction
rather than a wrong fault. Reading a label and an effect signature was not
enough; the failure text was where the truth was.

## The city counts out loud, so make the noun follow the number

Two things this pass, and the first one is a negative result that closes a
question I carried forward.

**The empty laundry was an artefact, not a hole.** Reading rooms on a driven
save, the laundry once returned no actions at all, not even the option to wait.
Hunting it directly — travelling at random and flagging any state where the
player's own room came back empty — reproduced it in twenty-seven moves and
explained it: a scene had opened during the journey, and while a scene is open
every room in the city correctly returns nothing, exactly as the death modal
does. My reading script checked for an open scene at the start of each room and
the scene opened in the middle of one.

The interrupted journey itself is handled honestly, which I checked because a
travel that returns success without moving the player is the shape of a lie. It
is not one. The result carries `from_location` and `to_location` both set to
where the player still stands, and the ledger entry reads "The journey to Mercer
Exchange was interrupted after 25 minutes. You remain based at Cypress House."
That also closes the travel banner, one of the two screens I had never read.
Fourteenth near-miss.

**Then the detail text of every action, side by side — eighty-four of them,
across all twelve rooms.** The labels had been read this way; the descriptions
had not, and the last real bug lived in one. This pass found the city counting
out loud and getting it wrong. Three buttons at the casino read "$164 due in 1
days", and the Herald read "1 stories in today's paper are about you or about
the police" — the noun wrong and the verb with it.

`counted` joins the grammar helpers in `core/names.go`, which are all keyed to
code rather than to save state for the reason established earlier tonight: a
fact about a word stored in a save means old campaigns keep old wording. Four
sites use it now, and the Herald agrees its verb as well.

Evidence: `core/counted_test.go` states three properties. The helper agrees at
zero, one, two and twenty-one; nothing the player can read in any room puts a
plural noun after a one, scanned across 168 actions in a world holding a loan
due tomorrow; and the Herald says one story is about you and two stories are.
All failed before the fix, and re-breaking both sites fails them again. Verified
over HTTP on the save that showed the fault: the three casino buttons now read
"$164 due in 1 day". `mise run verify` and `npm test` green, `mise run simulate`
unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

Found in the same reading and not fixed, recorded so it is not lost. Three
descriptions quote figures of zero: a loan the player cannot afford to make
reads "$0 out, $0 back inside 7 days", bringing money home with no account reads
"Brings $0 back into the city", and servicing a car the player does not own
reads "currently 0 of 100". All three are disabled with refusals that already
explain the situation, so nobody is misled about what they can do — but a
description of a transaction that cannot happen should not quote its terms as
zero. That is the next slice in this pass.

## Terms of zero are not terms

The three descriptions the detail pass left behind. A loan the player has no
room on their book to make read "$0 out, $0 back inside 7 days". Bringing money
home with no account read "Brings $0 back into the city". Servicing a car the
player does not own read "Restores up to 55 condition, currently 0 of 100".

Stating the severity plainly, because it is not what a nil-dereference is:
nobody is misled about what they can do. All three actions are disabled, and
their refusals already say why — "You have $262 out already, which is as much as
you can afford to be owed", "The account does not answer to you yet", "There is
nothing of yours to work on". The fault is narrower. A description exists to
tell the player what an action would do if they took it, and a figure of zero
tells them nothing while looking like a figure. The car one is the worst of the
three: "currently 0 of 100" reads as a wreck sitting in the yard rather than as
no car at all.

Each now says what is true when there is nothing to quote — the loan describes
its rate and term without a sum, the account says it brings back whatever is out
there, the garage says what a garage does — and quotes the figures unchanged the
moment there are any.

Evidence: `core/zero_terms_test.go` scans every action in every room for a
player who owns nothing and finds no dollar figure of zero, across 104 actions;
a second test covers the car, which has no dollar sign; and a third states the
other half, that a player with $1,200 offshore, a car in the yard and room on
their book still sees every figure. All three failed before the fix and fail
again when the zero-quoting forms are put back. Verified over HTTP on a driven
save with an empty account and no car: the market reads "Money out at 35% back
inside 7 days, when there is room on your book for it" and "Brings whatever is
out there back into the city", and the garage reads "What a garage does, once
there is something of yours in it". `mise run verify` and `npm test` green,
`mise run simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker
0, 0 errors.

## A possessive stacked on a family that names itself after a person

Read the Ledger screen in a browser — the last screen I had never read. The
running totals and the filters are sound: "$0 a day from 0 businesses" agrees,
and the eight filter counts sum exactly to the sixty entries the search box
offers. One entry did not read.

> Sofia Doyle helped themselves to $133 of Franca Sabbatini's people money at
> Bluebird Laundry.

Splinter families are named after whoever broke away, so their names already end
in a possessive plural. Putting "money" after one stacks a second possessive on
the first. The seeded families are why this survived: "$133 of Bellandi Family
money" is ordinary attributive English, and every campaign starts with two
families that read fine. It only breaks on the names the living world makes for
itself, and making families is most of what the living world does — thirty-nine
created across eighty long campaigns.

The line now says what was taken without a possessive at all, and agrees its
verb with the name through `Agree`, the helper already used for exactly this
class of organization name:

| name | reads |
| --- | --- |
| Bellandi Family | of what Bellandi Family keeps at Bluebird Laundry |
| Franca Sabbatini's people | of what Franca Sabbatini's people keep at Bluebird Laundry |
| the Duarte Brothers | of what the Duarte Brothers keep at Bluebird Laundry |

Evidence: `core/possessive_test.go` runs the line through four shapes of family
name and checks both the possessive and the verb. It fails when the old form is
put back, on the verb for the singular names and on the possessive for the
plural ones. `mise run verify` and `npm test` green, `mise run simulate`
unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

Also checked on that screen and correct, so it is not chased again: a commission
in progress reads "$0 of $900 taken out of Franca Sabbatini's people", which
looks like the zero-figure fault fixed above and is not — it is progress against
a target, and nothing has been taken yet.

## Prose written for two families, in a city that makes its own

The possessive found on the Ledger was one instance of a shape, so this pass
swept it. An organization inherited by somebody is named "<Person>'s people" and
is grammatically plural. The two families every campaign starts with —
"Bellandi Family", "Russo Outfit" — are singular. Any sentence written against
the openers alone is wrong for every family the world creates by that path, and
the world creates far more families than it starts with.

I did not trust a grep for this. `core/agreement_test.go` builds a city with a
plural-named family holding two premises, sets it at war with both openers, runs
four thousand half-hours, and then reads everything the city wrote down: 180
ledger entries, 240 newspaper stories, every action label, description and
refusal in all twelve rooms, and every commission brief. Eleven sites were
wrong, in the ledger, the paper, the commission board and a button description:
the family knew, wanted, believed, kept, had, was and did, where it should have
known, wanted, believed, kept, had, were and done. All eleven now go through
`Agree`, the helper already written for this.

The scan caught something my grep had not, and also reported something that was
not a fault. "Violence between Bellandi Family and Otto Reiss's people has
escalated" is correct English — the violence has escalated, not the people — and
the first version of the check flagged it. A plural name inside a prepositional
phrase does not govern the verb after it, so the check now reads clause by
clause and ignores a name that follows "between". Had I fixed what it reported,
I would have broken a correct sentence.

Evidence: the scan passes on a driven world and fails when any of the sites is
put back. A second test states the other half, which matters just as much —
"Bellandi Family does not keep people who leave in their good books" must keep
its singular verb, and leaving a seeded family end to end confirms it does.
Fixing agreement by making everything plural would read exactly as wrong for the
families most players actually meet. `mise run verify` and `npm test` green,
`mise run simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker
0, 0 errors.

One site that looks like the others and is correct, recorded so it is not
"fixed" later: the Herald's inquiry story reads "premises associated with X are
the subject of an inquiry into their books. Y has not responded to the
newspaper's questions", where Y is the leader rather than the family. A person
takes a singular verb whatever their organization is called.

## The same harness, pointed at three more kinds of name

The agreement sweep proved a method rather than fixing one bug: build a city
with the awkward case in it, run it forward, then read everything it wrote.
`core/naming_scan_test.go` generalises that to the other name shapes the world
can produce, and to a second hazard.

**A name that carries its own article.** Splinter families come out as "the
Duarte Brothers", and a name beginning with a lower-case "the" starts a sentence
in lower case unless something capitalises it. `Leads` exists for exactly this
and is used in eight places. Two more had missed it, and the second one is the
reason this test runs five cities rather than one. The first city wrote "the
Duarte Brothers struck Saint Agnes, a holding of Russo Outfit" and nothing else.
Only a different seed ever wrote "the Duarte Brothers and Falcone Brothers have
stopped short of destroying each other", because a sentence is only written when
the world happens to do the thing that writes it. One run's luck is not
evidence; the same line appears at three call sites and all three are fixed.

**Two clean results, recorded as results.** A family whose leader's name ends in
s — "Otto Reiss" — is written the same way everywhere: 21 possessives across 22
passages, all "Reiss's", none in the other convention. And the player's own
organization, which is named "<Player>'s people" and appears in prose a rival's
never does, reads correctly in all 31 passages that name it: no lower-case
sentence start, no verb disagreement. Neither needed a change.

Evidence: the article scan fails on both sites when `Leads` is removed, and
reports the two separately because they come from different seeds. `mise run
verify` and `npm test` green, `mise run simulate` unchanged at defiant 52 /
investor 0 / reckless 82 / worker 0, 0 errors.

What this pass is really worth is the harness. `writings(w)` collects everything
the city has written — ledger titles and bodies, newspaper headlines and
bodies, every commission brief, and every label, description and refusal in all
twelve rooms — and `live(t, w, n)` runs the world forward to produce it. Any
future question of the form "does the game read correctly when X" is now a
fixture plus a scan, run over several seeds.

## The same harness, pointed at states instead of names

`core/state_scan_test.go` builds cities in awkward conditions — a family holding
nothing, a family with nobody left, a city down to one organization, a player
with nothing, and a world run five times longer than anything measured before —
and reads everything each one writes against a set of checks rather than one.
Nine thousand passages. Three faults, and two false alarms I did not act on.

**A ledger entry that reported losses already taken.** Every midnight the player
could not cover the bills, it wrote "Security leaves; your residence is now a
rented room. Unpaid crew lose loyalty." After the first such night there is no
security to leave, the residence is already a rented room, and the crew's
loyalty is already nothing. The entry now names only what it actually took, and
when there is nothing left says so: "You could not cover the bills. There is
nothing left to take, which is its own kind of trouble."

**A splinter announcement that agreed one family's verb and not the other's.**
"Bruno Duarte's people wants it back" — the line already ran the new family's
name through `Agree` and left the parent's alone.

**A refusal that could not tell two cases apart**, the sixth of that class.
Lending is capped at a share of everything the player has, out and in hand, so a
full book and an empty pocket look identical to the arithmetic. Both were
refused with "You have $0 out already, which is as much as you can afford to be
owed", which for a player holding nothing is true and says nothing. A player
with no money is now told "Nobody lends money they do not have. It takes $300 in
hand to put $120 on the street."

**The two I did not act on.** The scan reported "Janos Kovac of Vera Kohl's
people is not to see the end of the week" as a disagreement. The subject is
Janos Kovac; the plural name sits in a prepositional phrase and governs nothing.
That is the second construction of this kind, after "Violence between A and B
has escalated", and the check now skips a clause where the name follows "of" as
well as "between". Skipping a whole clause containing "of" can mask a real
fault, and that is the price of not manufacturing false ones.

The other was mine. Over four hundred game days, 154 of 180 ledger entries were
the same headline — but this harness keeps a broke player alive far past where a
campaign ends, so the bills fail every midnight by construction. That is what I
built, not what the game does. I removed the assertion rather than "fixing" the
game to satisfy it, and pointed the test at the newspaper instead, which is the
living world's own voice and owes nothing to a player standing still.

That gave a new measurement worth keeping. Over 20,000 half-hours, about four
hundred and seventeen game days, the paper printed **240 stories with 117
distinct headlines**, the most repeated appearing **17 times**. The city does not
run out of things to say at length.

It also corrects something I published two slices ago. I wrote that the people
prune "has almost certainly never run in any real campaign" because a city holds
about 52 people at day 47. That holds for the horizons a campaign actually
reaches, but the population does grow: at day 417 this world held 153 people,
which is past the threshold of 150. The prune is reachable, just not soon.

Evidence: each of the three fixes fails its test when reverted, confirmed one at
a time. `mise run verify` and `npm test` green, `mise run simulate` unchanged at
defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

## Agreement and capitalisation are two jobs at the same site

Four more states through the reading harness: the player in a cell, the player
dead, a city holding as many organizations as it allows, and two families whose
names collide in the same sentence. Two read clean, two did not, and what they
found is a lesson about the previous slice rather than a new class of fault.

**Custody and death read clean.** In a cell, everywhere but the precinct
correctly offers nothing, and the three ways out read properly. After the player
dies, the epitaph and everything around it read properly. 1,294 and 962 passages
respectively, across three cities each. Neither needed a change.

**The other two found eight sentences beginning in lower case**, all of the same
shape: a family named "the Duarte Brothers" starting a sentence. But five of
those sites are ones I had edited *in the previous slice*, when I gave them verb
agreement and did not notice they also began sentences. `Agree` and `Leads` are
two independent requirements that land on the same interpolation, and I had
applied one and moved on.

So rather than trust the scan alone, I listed every call to `Agree` in the core
and asked of each whether its name begins a sentence. Twenty-six sites: sixteen
already capitalised, seven correctly mid-sentence, and five that needed it and
had not been reached by any scan yet. The scan and the list found different
things — the scan found what the world happened to write, the list found what it
could write — and neither would have been enough alone.

Ten sites in total this slice: three war reports, three commission briefs, and
the four sentences about leaving service, a rival learning where weapons came
from, a family under new leadership, and a family that has stopped complaining.

Evidence: nine scans now pass across the states and names tried so far,
`core/state_scan2_test.go` adds the four new ones, and reverting any fix fails
its scan. The other half is stated too: `Leads` must leave an ordinary name
untouched, checked against six name shapes and end to end by leaving a seeded
family, or every campaign's two openers would be mangled to fix the ones the
world invents. `mise run verify` and `npm test` green, `mise run simulate`
unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

## Bailing somebody out is not a conversation with them

The prose sweeps had gone quiet, so I went back to the actions whose label, gate
and effect had never been read together. Bail was the first, and it could never
be used.

There is one rule in the action list, applied once over everything aimed at a
person: somebody who cannot be dealt with is not offered. Its own comment
explains why it lives in one place — "doing it here rather than at forty call
sites means the next action about a person cannot forget" — and it is right
about that. It stops the player sending a man on collections from a police cell,
which is a fault this project has fixed before.

Bail is the one action that exists *because* the person is in a cell. The sweep
saw an action aimed at somebody in custody and disabled it with "Otto Reiss is
being held at Ward Street Station", which is the reason the button is offered.
So the button was dead in every campaign, and the only way to reach the mechanic
at all was to post the command directly to the API, which is how I had
"verified" custody earlier tonight without noticing the button never worked.

The rule now has exactly one exception, named and commented, and a test holds it
to one action wide: everything else aimed at a man in a cell stays refused, four
of four in the check.

The same audit found the smaller thing beside it. The ledger entry bail writes
has always counted correctly — "for the day still on them" — while the button's
description said "for the 1 days still on them". The effect knew better than the
label, which is the reverse of the split I had been looking for all night.

Evidence: `core/bail_test.go` states that bail is offered when it can be paid,
refused on the button when it cannot, counts one day as one day, and reaches
into a cell where nothing else does. Each fails when its fix is reverted, checked
one at a time. Verified over HTTP on a save driven to day 28 with one of the
player's people held for two days: the precinct offered "Bail out Otto Reiss",
"$640 for the 2 days still on them", the command took $640, and he walked out.
`mise run verify` and `npm test` green, `mise run simulate` unchanged at defiant
52 / investor 0 / reckless 82 / worker 0, 0 errors.

Correcting my own record: bail appears in this document's list of things
verified by hand over HTTP. That was true of the command and not of the button,
and I did not check the difference at the time. Posting a command directly is
not evidence that a player can reach it.

## Which buttons can never be pressed

Bail was dead in every campaign and nothing failed, because no test ever asked
the obvious question of it: is there any state at all in which this button is
live? `core/reachable_test.go` asks that of every button the game can produce.

It builds cities in six conditions — a rich established player with a crew, an
organization and somebody of theirs in a cell; the same city aged two thousand
half-hours; a beginner with just enough standing to open the next district;
somebody serving a family; nobody's man who is well thought of everywhere; and a
proprietor with the press and the police on a retainer, short-handed premises,
the police interested in him, and money abroad he has not yet reached. Across
five seeds it records every action id offered and every id offered *enabled*,
and reports anything that appears only in the first list.

The first run named twenty-seven kinds that were never live. Sixteen of those
were my fixtures rather than the game, and building the right state cleared
them one at a time. Two of those corrections are worth recording because I had
the game's own rules wrong: a retainer is held by the **role**, not by the
person in it, so storing the editor's id retains nobody and the whole newspaper
went on refusing; and a business the player takes over arrives fully staffed, so
hiring is correctly refused until somebody leaves.

Seventy-eight kinds are now confirmed pressable, including every action this
document had listed as never reached in play — bail, bribe, charge, move,
remedy, leaving a family's service, and the trip upriver among them. Four
remain, and each was read and explains itself with a real condition rather than
a fault: a takeover needs the family weak and you in the room with its leader,
a pact needs standing the sweep's player has not earned with that particular
family, and two are per-family variants whose base form is live.

The regression is the point. Sixteen of the hardest-to-reach buttons are named
in the test and must stay live; reverting last slice's bail fix makes it fail
with "bail:* is offered 5 times and never live: no state in this sweep can press
it", which is what should have happened months ago. `mise run verify` and `npm
test` green, `mise run simulate` unchanged at defiant 52 / investor 0 / reckless
82 / worker 0, 0 errors.

Recording the shape plainly, because it is the seventh this audit has produced
and the only one that hides in what a test does *not* ask: a button offered and
never live looks exactly like a button correctly refused, and only counting
across many states tells them apart.

## Twenty-nine and seventeen out of fifty-two

Back to the browser after sixteen fixes, and a second reading of the People
screen. Its opening line:

> 52 people live in this city and you know 33 of them. 29 answer to an
> organization and 17 to nobody.

Twenty-nine and seventeen make forty-six. Six people were in neither figure: the
four officials, the fixer, and the player's own crew before the player is an
organization. Somebody doing one of the city's jobs answers to the city rather
than to a family or to nobody, and the sentence had no room for them.

It also disagreed with the screen underneath it. The filters file each person in
exactly one place and their counts do add to the whole city — crew 1, they owe
you 1, names everybody knows 10, organizations 23, the street 17. So the header
said twenty-nine answer to an organization directly above a chip reading
"Organizations 23", because that chip does not hold the family heads it files
under names everybody knows. Both numbers are correct under their own
definitions and a reader cannot reconcile them.

The summary now carries a third figure and the line reads "29 answer to an
organization, 6 hold one of the city's jobs, and 17 answer to nobody", which
comes to fifty-two.

Evidence: `core/population_test.go` states that the three figures add to the
population, in a new city, in one run four thousand half-hours, and in one where
the player is an organization so their crew answers to somebody. It fails when
the new figure is removed, in all three. Read on screen against a save driven to
day 32.

One thing worth recording about the reading itself. The first look at the fixed
screen showed "undefined hold one of the city's jobs", because the page was the
new build and the server behind it was not. That is not a fault in the game and
it is exactly why the screen gets read rather than the code: a field added in
one place and not served from the other looks fine in every test and wrong to
the only person who matters.

`mise run verify` and `npm test` green, `mise run simulate` unchanged at defiant
52 / investor 0 / reckless 82 / worker 0, 0 errors.

## A card that drops a column rather than admitting it does not know

Second reading of the Families screen. Each card lays out four figures —
strength, money, people, ground — and two of the five had three. The People cell
was simply not there.

Strength and money already degrade properly. With nobody inside a family they
read "not much"; with nobody to ask at all they read "nobody will say". People
did not degrade, it vanished, so the same ignorance was expressed two different
ways on one card and the reader saw a table whose columns change from row to
row. It also could not tell two things apart: a family the player knows well
that has nobody left looked exactly like a family they know nothing about.

The people figure now behaves like the other two. Nobody to ask reads "nobody
will say", an impression reads "a handful", "a fair few" or "a great many", and
somebody inside gives the count — including "nobody left" when that is the count.

A crowd gets its own words rather than borrowing the money scale. My first
version reused `roughly`, and eight people came out as "as much as anybody",
which is a strange thing to say about eight people.

Fifteenth near-miss, and this one nearly reached the write-up. I read the
screen, saw two cards missing the People column, checked the API and found
`people` absent for exactly those two, and concluded those families had nobody
left. They had eight and two. The field is omitted below the knowledge level
that reveals it, and the count was never zero. The fix is the same either way —
a column that disappears is the fault — but the reason I would have published
was wrong.

Also checked on that screen and correct, so it is not chased again: every
declared quarrel is reciprocated. Bellandi lists four, and each of the four
lists Bellandi back, with "at war" and "at odds" matching on both sides.

Evidence: `core/families_card_test.go` states that every card carries a people
figure at every knowledge level, that a family with nobody left does not read
like a family nobody will discuss, and that the crowd words are not the money
words. All fail when the change is reverted. Read on screen against a save
driven to day 32. `mise run verify` and `npm test` green, `mise run simulate`
unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

## The panel kept its own copy of a list the core owns

Second reading of the room panel, with the arithmetic check that has been
paying. The search box offered "Search 10 actions here" and the sections held
nine. One action was not on the screen at all.

The core owns the list of groups an action can belong to, complete with titles
and blurbs, and says so in a comment: "Groups is the ordered list, for anything
that renders them." Nothing rendered them. The panel kept its own copy of the
list, and the copy was missing **people**. So every action in that group whose
subject was not standing in the room fell out of the panel entirely — not
greyed out with a reason, gone. Paying the crew a bonus while they are away on
collections is the ordinary case, and the player could see neither the button
nor why it was unavailable. Lending, collecting, extending and writing off a
debt all live in that group too, and all vanish the moment the other person
walks out of the room.

The panel now renders from what the core sends, so the two cannot drift again,
and anything carrying a group this build has not heard of joins the first
section rather than falling out — the same fallback the core applies to an
action nobody has classified.

Sixteenth near-miss, and it mattered. Bail is in that same group, and I fixed
bail's reachability the night before. My first thought on seeing this was that
the fix had never actually reached the screen. It had: somebody held at the
station is in that room, so bail attaches to their card and renders. Checking
before writing saved a retraction.

Evidence: `core/grouping_view_test.go` states from the core's side that every
group an action can carry has a title, checked against every action the
reachability sweep can produce across six kinds of city. `tests/grouping.test.mjs`
states the panel's half: every action lands in a section, a group this build has
not heard of is not dropped, nothing is placed twice, and no groups at all
places nothing rather than guessing. Both halves fail when broken — removing the
fallback loses the unknown group, and deleting "people" from the list reproduces
the original bug exactly. Read on screen against a save driven to day 32, where
the panel now shows a People section with the crew bonus behind "Show 1 you
cannot do yet", and the sections add to ten.

Also checked on that screen and correct: the books add up, income nought against
costs of thirty-two, itemised as rent fifteen, crew twelve and the house five.
`mise run verify` and `npm test` green (29 tests), `mise run simulate` unchanged
at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

## Seventeen promised, six delivered, eleven struck off unread

The arithmetic check, fourth time, now on the Herald. The badge on the rail read
17. The button read "Read the paper (17)". Opening it showed "Today's edition ·
Friday, April 3, 1953 · 6 stories", and the page held exactly six. The other
eleven were the issues from the two days before.

Worse than the mismatch: opening the paper marked the newest story as seen, and
the count was computed as that story's position in the list — so all seventeen
became read while six were on screen. Eleven stories were struck off without
ever being shown.

It also grows without bound. The list is every story of the current life, so a
campaign at day four hundred would open with a badge in the hundreds, which is
not a number anybody can act on.

The count is now what opening the paper will actually show: the stories in the
latest issue the player has not seen. The archive is untouched and still one
click away behind "All 32 issues" — it is the prompt that had to be honest. The
button says "Read today's paper (6)".

Evidence: `src/paper.ts` holds the count and `tests/paper.test.mjs` states five
properties, including the one that was broken — a story from an older issue must
not mark today read. Restoring the whole-life count fails two of them. Read on
screen against a save driven to day 32, with the seen-marker cleared first so
the paper was genuinely unread: the button offered six and the issue held six.
`mise run verify` and `npm test` green (34 tests), `mise run simulate` unchanged
at defiant 52 / investor 0 / reckless 82 / worker 0, 0 errors.

Checked on the way past and correct, so it is not chased again: the Ledger's
seven filter chips still sum exactly to its sixty entries, and three successive
repair entries reconcile — restored by 40 to 50%, by 40 to 90%, by 10 to 100%.

One thing read and left alone, because it is a judgement rather than a fault.
Two robbery stories ran side by side with the same standfirst, "The proprietor
of X declined to be photographed. Officers have asked witnesses to come
forward", and the second body then repeated the same appeal. It reads like a
template because it is one. Giving the paper more ways to say this is a writing
task, not a correction, and it is worth doing deliberately rather than in
passing.

## The guide page makes a promise, so test the promise

The guide opens by describing itself: "This page is not written down anywhere.
It asks the game the same question the buttons ask, so it cannot tell you
something the rules do not." That is a claim about the code, and it is testable.
It was false in two ways.

**It offered what every button refused.** With $51 and 33 respect, the guide
said "Premises of your own — You can do this now", and every door in the city
said "Not enough cash". The step re-derived its own gate from respect alone,
because acquiring premises was the one action with no readiness function: the
button computed the standing requirement inline and let the generic cost check
add the refusal. So the guide could not ask the same question, and invented a
worse one. `AcquireReadiness` now exists and both use it.

**It preferred its own words to the game's.** Each of the five searches the
guide runs — over people, places, families, officials — reports the briefest
refusal any candidate gave, and each was seeded with a sentence for the case
where there are no candidates at all. Those sentences are short, so they won
against every real reason. A player with no organization was told "There is
nobody in this city for that yet" while the buttons in front of them said
"Nobody signs on with one person. They sign on with something that has a name."
One helper now holds all five, and the invented sentence appears only when there
was genuinely nothing to ask.

Seventeenth near-miss, and it was my test rather than the game. My first version
of the second check required every guide refusal to be a sentence some button
was showing. That is stricter than the promise: the guide quotes readiness
functions, and a lending refusal like "That conversation happens in person" is a
real rule that no button ever displays, because lending is only offered for
somebody already in the room. I rewrote the test to state what was actually
wrong — the invented sentence must not beat a real one — rather than bending the
game to a standard I had made up.

A second thing I introduced and then corrected by reading: the new guard reads
"There is nothing here to take over", which is a sentence about a place. At page
level the guide is not standing anywhere, and on a save where every earning
premises was owned it surfaced as the shortest reason. The premises step now
asks only of places that are businesses, and reads "This property belongs to
another organization".

Evidence: `core/guide_promise_test.go` states both halves. Reverting either fix
fails its own test — putting the re-derived premises gate back reproduces the
"$51 and it says you can do this now" fault exactly, and letting the fallback
compete again reproduces the invented sentence. Read over HTTP against a save
driven to day 31. `mise run verify` and `npm test` green (34 tests), `mise run
simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors.

Checked and clean this pass, recorded so they are not redone: the guide states
no total of its own, so there is no arithmetic to reconcile there; the market
lists three goods with a price and a reference price each and nothing that
sums; and the district names in the view are not a copy of a core table — core
has no such table, it names them in prose.

## The rules the game prints about itself

The guide prints six rules as prose. Each is a claim about the code, and the
page beside them had already turned out to be wrong about itself twice, so I
asked three of them rather than assuming.

**"The clock only moves when you commit to something."** Read as a property: an
action declaring no time must not move the clock. Twenty-four free actions were
taken across two cities and every one of them was honest, except the three trips
out of the city. Those declare zero minutes and cost two, three and four days.
The reason is sound — a trip runs the clock a day at a time so that what happens
while the player is away happens to a city they are not standing in — but zero
is also what the panel prints the cost from, so the longest actions in the game
showed no time at all beside "Rest for four hours · 240 min". The description
said "2 days away" and the button said nothing.

An action can now declare time its effect will spend itself, and the trips read
"2 days away" where everything else reads its minutes. The command layer still
spends nothing, so nobody is sent away for a fortnight.

**"Attention is public and so are the thresholds."** The sentence quotes two
constants; it quotes the right ones, and they are in the right order. Clean.

**"Anybody out on the street cannot be dealt with until they arrive."** True, and
kept in two places rather than one: somebody walking is not in the room, so most
work about them is never offered, and anything that still is gets refused. The
test states both halves, including that the same work becomes possible again
when they stop walking — otherwise it would pass for somebody nobody can ever
deal with.

**One more hardcoded name, in a class I had declared clean.** Paying the crew a
bonus wrote "A share for Leo" whoever was actually on the books. My earlier
sweep missed it because the pattern looked for a quote before the name or a
space after it, and this one sits at the end of a string. Correcting the record:
that sweep was not clean, and the lesson is that a grep anchored on delimiters
misses the cases at a string's edge.

Eighteenth near-miss, and the second time the same trap has caught me: my first
fixture for the street rule picked the fixer, who is refused earlier and
correctly for a reason that has nothing to do with the street. The second picked
somebody I had appended to the city, who never receives a lending offer because
only the first few people in a room do. The test now takes whoever the game is
actually offering.

Evidence: `core/stated_rules_test.go` holds the three properties. Reverting the
trips to a bare zero fails the first; the others fail if the constants or the
street rule move. `mise run verify` and `npm test` green (34 tests), `mise run
simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors. Verified over HTTP on a driven save: the three trips report 2,880, 4,320
and 5,760 minutes away where they used to report nothing.

## The player was taking a third more than anybody else

Fourth of the six rules the game prints about itself: "Everything anybody in
this city does is done by the same rules you use. There is no separate
arithmetic for them."

The randomness is deliberately two streams, one reserved for events the player
is not party to so that families quarrelling off-screen cannot shift the odds of
a decision being made. That is separate randomness for a stated reason, and it
is right. The arithmetic is a different question, and it was two functions: a
family moving on a rival runs `contestAt`, the player moving on the same
premises runs `SabotageBy`.

Measured over four hundred attacks each, on identical premises:

| what an attack costs | by the player | by a family |
| --- | --- | --- |
| condition off the building | 34.6 | 32.6 |
| money off the family | $689 | $483 |

Six percent apart on the damage, which is two similar formulas drawing similar
numbers. Thirty percent apart on the money, because they used different
constants — twenty a point when the player did it, fifteen when anybody else
did. The player was taking a third more out of a family than the city could take
out of anybody, which is exactly the privilege the rule says does not exist.

One constant now, named, at the city's figure rather than the player's. The
money figures are 6% apart, which is the damage difference and nothing else.

The balance did not move. `mise run simulate` is unchanged at defiant 52 /
investor 0 / reckless 82 / worker 0, 0 errors, and the long horizon over eighty
campaigns is identical on every measure: 112 wars started, 74 elsewhere, 52
holdings changed hands, 39 families created, 21 destroyed, 268 hurt in somebody
else's war. That is what should happen — the constant governs what the player's
own sabotage takes, and the simulated strategies rarely reach for it.

Evidence: `core/same_arithmetic_test.go` measures both and fails when either
constant is put back — restoring the twenty reproduces "30% apart" exactly.

A flaw in my own check, recorded because it is the kind that hides: the first
version of the comparison found the larger of the two figures inside an `else`,
so it was skipped whenever the family's figure was the larger. It could only
fail in one direction, and it happened to be the direction that was green.
Breaking it again does not fail now, because with the fix in place both figures
are close either way — so the flaw is latent rather than visible, and the
arithmetic is tested directly instead of pretending otherwise.

## The last two rules the game states, and no fault in either

The remaining two of the six. Both hold, and the work was almost entirely in
getting my own harness out of the way.

**"What you built passes to the strongest of your own people and becomes an
organization you can deal with, or fight."** The passing had been verified once.
The second half never had. A player with three premises and three people of
their own is killed, and the next life arrives to find all three premises and
all three people under one name with a leader at its head — and the city offers
the new arrival every ordinary way of dealing with it: asking around about it,
reaching an understanding, going to work for it, running something about it in
the paper. Its premises can be moved on like anybody else's. Nothing here needed
changing.

**"The city does not scale to you. An organization at ninety strength will kill
you on your first day if you give it a reason."** A hundred and twenty new
arrivals provoked the opening family and went home. Seventy-seven died. Forty-
three were still standing two days later. Both halves are asserted: a majority
must die, or the sentence is a warning the city does not mean, and some must
live, or the provocation is not a risk but a way of ending the game.

**Two near-misses in one test, both mine.** The first version cleared any scene
that opened rather than answering it, and an attack arrives as a scene — so
discarding it discarded the attack, and nobody died. The second version answered
properly and still nobody died, because the player was left standing in the
nightclub they had provoked for two days running, and `Attack` checks whether
the player is at home before anything else. The car pulls up outside an empty
room. A person is at home at night; the test now puts them there.

That is the nineteenth and twentieth near-miss, and both had the same shape: a
harness that avoided the thing it was measuring. It is worth saying plainly that
"nobody died in a hundred and twenty cities" was twice a fact about my code.

All six stated rules are now tested. One was false and is fixed — the player was
taking a third more money out of a family than the city could take out of
anybody. Five hold. Evidence: `core/inheritance_test.go` and
`core/stated_rules_test.go`; both of these fail when broken, one when the
holdings stop passing and one when the city stops killing. `mise run verify` and
`npm test` green, `mise run simulate` unchanged at defiant 52 / investor 0 /
reckless 82 / worker 0, 0 errors.

## Work that pays its own fee showed no price

Auditing the actions that had never had their label, gate and effect read
together. The sitdown reads "Call Vittorio Bellandi and Elena Russo to a room",
its description says "$220 for the room and the guarantees", and its button
showed no money at all.

The reason is the same one the trips had with time. The panel prints a price
from the declared cost, and the command layer *charges* the declared cost — so
an effect that pays its own fee has to declare nothing, or the money goes out
twice. Calling a sitdown pays `SitdownFee` itself, so it declared zero, so the
price vanished from the only place a player looks for one.

An action can now name a fee the command layer will not take, the same way it
can name time the effect will spend. Six actions that pay their own way say what
they cost: the sitdown, pulling a story, planting a paragraph, reaching an
account abroad, a still, and a room under the floor.

A false start worth recording. My first attempt measured this by pressing every
action a rich player could press and comparing cash before and after. Twenty-four
of twenty-seven came back "charges silently", which is nonsense: those actions
advance the clock, so the difference includes a day's income, wages and rent. It
also flagged lending, where the money is not a fee at all but a sum going out on
the street. The measurement was of my own harness, which is the twenty-first time
that has happened tonight. The real audit was reading every `w.Pay(` in the core
against the declared cost at the button that leads to it.

Evidence: `core/asks_test.go` states both halves — the fee is on the button, and
declaring it does not charge it twice. Removing it from any of the six fails the
first; the second watches a still cost $450 and take $408 once the hours' own
income and rent had landed, which is one fee and not two. `mise run verify` and
`npm test` green, `mise run simulate` unchanged at defiant 52 / investor 0 /
reckless 82 / worker 0, 0 errors.

Scope stated honestly: this covers the six whose fee is a plain constant at the
point the button is built. Others charge a figure that is computed further in —
a retainer's opening payment, laundering, a car, a suit, arms, restocking — and
each would need its figure lifted to the button. They are the same fault and are
not fixed here.

## The rest of the silent prices

Finishing the previous slice, which fixed six of them and named the rest. Twenty
actions in the game pay their own fee, and fourteen were still showing no price:
an arrangement with an official, laundering, a car, a suit, a weapon, armour, a
charge off a boat, an understanding with the detective, a pact, a bankroll,
hiring, restocking, a remedy, asking around about a family, and running
something about one in the paper.

Read over HTTP on a driven save, every one of them now names what it takes and
declares no cost, so the engine takes it once:

| where | what it asks |
| --- | --- |
| an arrangement with the commissioner | $1,476 |
| an understanding with the detective | $2,920 |
| a pact with a family | $700 |
| a charge off a boat | $850 |
| a used Ford | $620 |
| asking around about a family | $90 |

The regression guard is worth more than the fix. My first version of it asserted
that "at least twelve" actions named a fee — and silencing one still left
twenty-four, so it passed when I broke the code. It names the twenty now, and
both failure modes fail it: making one go quiet, and making one declare a cost
the engine would charge on top of its own fee. That is the twenty-second time a
check of mine has needed the same lesson, and the first time I have caught it by
running the break before writing the slice up rather than after.

Evidence: `core/asks_test.go`. `mise run verify` and `npm test` green, `mise run
simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0, 0
errors, with identical median cash — as it should be, since nothing about what
anything costs has changed, only what the button says.

The shape is closed. Twenty actions pay their own way, all twenty say so, and a
new one that does not will not be caught by any test from outside — that is
stated plainly rather than implied, because the guard protects what exists and
cannot know what has not been written.

## Bail was charged twice

I had recorded, and written into the previous slice, that the engine charges a
declared cost only in the final switch, so the actions carrying a person's name
after a colon were exempt. That was wrong. The line that takes the money sits
above the whole prefix chain and runs for every action in the game.

Bailing somebody out declared the fee on the button and paid it again inside the
effect. A three-day bail costs $960 and took $1,920. The test measures the
pocket, not the ledger, and subtracts what the same hour costs in a second world
that does not bail anyone, so the fee stands alone:

```
bail of $960 took $1920 out of the pocket (spent 1864, the hour itself costs -56)
```

Fixing it broke an older test, which is the best thing that happened here. A
declared cost also disables the button when the player cannot afford it, and a
display-only price does not. The affordability gate had been riding on the same
field as the charge, and moving one silently dropped the other. The readiness
function is now passed explicitly.

The narrow fix is worth less than the rule, so the rule is a test: nothing the
prefix chain handles may declare a cost. The thirteen prefixes are written out
by hand rather than derived from the file the rule is about, because a list
built from the code would agree with the code no matter what the code said.

Two silent prices turned up in the same reading. Signing somebody on and paying
somebody a share both pay their own way and both showed no money on the button.
They were missed the first time because the earlier guard matched exact action
names and these carry a person's id. It now matches on the prefix, and the test
world has somebody of the player's in a bar and somebody of theirs in a cell so
that all three are actually exercised. Both failure modes fail:

```
share: pays its own fee and names no price
sign:  ...also declares a cost of 140, so the money would go out twice
```

A share cost sixty dollars in three places — the gate, the payment and the
button — as a bare number in each. It is `ShareCost` now.

Scope: this covers actions reachable in a well-supplied city and a cell. Twenty-
five of the twenty-six named are exercised; a remedy is not offered anywhere in
that city, so it is checked by reading and not by the test.

Evidence: `core/double_charge_test.go`, `core/asks_test.go`. `mise run verify`
and `npm test` green at 34, `mise run simulate` unchanged at defiant 52 /
investor 0 / reckless 82 / worker 0, 0 errors, cash 7122 / 14019 / 90 / 12360.

## Families had money and no bills

From the inbox: track what people and families are worth, and let it drive what
they do. The first question was whether family money was a number worth
deciding from. It was not.

Families collected from their holdings every morning and paid for nothing. Over
a simulated year not one organization in the city was ever short of anything:

| day | poorest family | median | richest |
| --- | --- | --- | --- |
| 30 | $3,869 | $26,100 | $45,440 |
| 120 | $14,745 | $86,400 | $153,938 |
| 400 | $62,223 | $232,094 | $547,077 |

So "a family with little money" was a state the world could not reach, and a
director keyed to wealth would have been reading a number that only ever went
up.

A family now pays for its day. Wages go by strength rather than by the handful
of people who happen to have names, because a family of ninety runs more than
nine people, and paying only the named ones would have made a family of ninety
and a family of ten cost the same. Repairs are bought rather than granted: the
morning used to put eight points of condition back on every damaged property in
the city and take nothing for it. Poverty is now reachable, and between one and
four families are short at any time.

The starvation attack, tested end to end, did not work. A family stripped of
every business faded on a timer at four points a day, and that ran faster than
its money did, so it shrank to nothing without ever failing to pay anybody. The
losing of the businesses was never the cause of anything. An organization now
holds together for exactly as long as its money lasts and comes apart at the
first payday it cannot meet:

| what was done to them | after | strength | worst run short |
| --- | --- | --- | --- |
| left alone | $8,912 | 70 | none |
| every business taken | $0 | 0 | 11 days |
| wrecked once | $10,678 | 70 | none |
| wrecked every week | $40 | 12 | 5 days |

A single wrecking is absorbed: they buy the damage back over about a fortnight.
That is measured against a control family run through the same weeks untouched,
because comparing them to their own past only says the year went on.

I changed an existing test rather than the code it guarded, which needs saying.
It asserted that a landless family loses strength on the first morning. It does
not any more, deliberately: with money in the safe it holds together, and the
money is now the reason it comes apart rather than the calendar. The rule the
test was written for still holds, and the test still checks it, one month later
instead of one morning.

Two prose bugs surfaced that had nothing to do with money. Different families
now survive, so names the seeded cities never produced reached sentences that
had never seen them: "the Duarte Brothers had come for The Monarch" opened three
death notices and a war declaration in lower case, and a man a minute away was
"1 minutes out" in three places. Both helpers already existed and neither was
being called. That is the fifth time a name the world invented has broken prose
written for the names it ships with.

Scope: the balance simulation is unchanged at defiant 52 / investor 0 / reckless
82 / worker 0 with identical median cash, and that is not evidence the change is
harmless. Those campaigns run seven to thirteen game days. Family finances take
about a month to say anything, so the simulation is too short to see this at
all. The evidence for the change is the probe and the starvation tests, not the
baseline holding still.

Individual people still have no money. Mugging invents a purse from the family's
cash at the moment somebody is robbed rather than taking what a person actually
has. That is the next slice.

Evidence: `core/family_money_test.go`, `core/starve_test.go`,
`core/money_probe_test.go`. `mise run verify` and `npm test` green at 34.

## People have money of their own

The second half of the inbox item: everybody in the city is a person like the
player, so they should have money like one. They had none. What somebody was
carrying was worked out at the moment they were robbed, from their rank and
their family's cash, and nothing was ever taken off the person. The same man
could be robbed every day for a year and be carrying the same amount every time.

A person now owns a purse. The old derivation survives as what a person of that
standing would carry on an ordinary day: it is what the city settles a newcomer
onto, and what their week pays them toward. Robbing somebody empties it.

The reason this matters is that it joins the two halves of the inbox together. A
family that misses payday does not pay its people, so starving a family reaches
the men in it rather than stopping at a number on a screen:

| the family | between five people | broke |
| --- | --- | --- |
| paying its way | $1,359 | none |
| everything taken, six weeks | $640 | four of five |

Officials are on the city's books rather than a family's and are paid whoever is
struggling.

Three things went wrong worth recording. The first arithmetic I wrote paid
people a seventh of their standing and charged them twelve dollars a day to
live. An ordinary man carries about thirty dollars, a seventh of which is four,
so three quarters of the city was broke inside four months and the median person
had nothing. Earnings now cover the day and make up a seventh of whatever the
week has taken, so somebody robbed on Monday is himself by the weekend. Thirty
to forty people in eight hundred are broke, and they are the ones in starved
families, which is the point.

The second: a family that had collapsed entirely owed nothing, so it stopped
being short and started paying its five remaining men out of an empty safe. The
bill was charged on strength alone, and the named people were not in it. They
are now. That is the third time the same reset has bitten in one night.

The third: lending asks what somebody could find over a week, and its own
comment says explicitly that this is "not what is in their pocket tonight". It
was calling the function that now returns exactly that. Its comment was right
and the code had quietly become wrong about it.

Four tests were changed rather than the code they guarded, which needs saying.
Two asserted that rank and a rich family decide what is in a pocket. That rule
moved rather than went away: it is now about what somebody is paid toward, and
both tests check it there. A balance test promotes a man after the fact and had
to settle his pocket onto the standing he now has, or both ranks were measured
carrying whatever the same man happened to have. One asserted that a landless
family loses strength on the first morning, covered in the previous slice.

## A sweep instead of a sentence a week

Four family names beginning a sentence in lower case were fixed today, one at a
time, each surfacing only because the world happened to write that sentence. The
existing scan gives one invented family an article-carrying name and waits.

The new sweep renames every organization in five cities to carry a lower case
article, ages them, renames whatever the world invented in the meantime, and
reads every passage. It found one more site nobody had reached, in the brief for
having somebody killed. Then it read 5,724 passages clean. Reverting one fix
makes it fail, so it can see what it claims to check.

A grep for format strings that start a sentence with a name suggested seventy
candidate sites. The sweep found one. The grep was over-reporting, and the sweep
is the honest measure of the two.

Evidence: `core/purse_test.go`, `core/lowercase_sweep_test.go`,
`core/money_probe_test.go`. `mise run verify` and `npm test` green at 34, `mise
run simulate` unchanged at defiant 52 / investor 0 / reckless 82 / worker 0 with
identical median cash — and as before, those campaigns are too short to see any
of this.

## The director can see who is struggling

The last part of the inbox item: family money should drive what the director
writes. It could not, because the director had never been told about money. The
full brief ships the whole faction record, so a family's cash was in the JSON as
a bare integer, and neither prompt mentioned it. Searching the director for the
word "money" or "cash" returned nothing at all.

A bare integer would not have helped anyway. Ten thousand dollars is a fortune
or a fortnight depending entirely on what a day costs them, and nothing knew
what a day cost them: the bill was a line inside the morning that charged it.

So the bill is a function now, and three things ask it. How an organization is
placed is one of four states — comfortable, getting by, struggling, or cannot
pay its people — decided by how long it could go on paying everybody if the
money stopped tomorrow. The words belong to the core. What it means to be
struggling is a fact about the world, not a phrase a prompt or a screen invents.

Both briefs now explain what the four states mean for behaviour, close to what
the inbox asked for: a comfortable organization is patient, protects what it has
and is worth taking from; a struggling one presses harder and moves on somebody
else's ground; one that cannot pay its people is desperate and unpredictable.

One judgement worth recording. A family living within its income is not given a
countdown. Ninety days of cover is not a thought anybody in the world is having,
and handing the model that number invites a speaker to talk about three months
of runway nobody is counting. The field is simply absent for them.

Three tests, and all of them fail when the code is broken: renaming a state in
the prompt is caught, and so is the context sending a word the prompt never
explains. There is also a test that the bill described is the bill the morning
actually charges, so the two cannot drift apart.

Not done, and next: the controlled comparison. Whether the model actually writes
differently for a struggling family is a claim I have not tested, and I am not
making it. All that is established is that the information reaches it and the
prompt explains it.

Also still open: nothing in the core reads a person's poverty yet. `Broke` is
consulted by one button. A broke man and a comfortable one behave identically.

Evidence: `core/placed_test.go`, `cmd/blackledger/director_money_test.go`. All
gates green, `npm test` 34.

## The director reads the money picture and writes the same thing anyway

I shipped the money picture last slice and wrote explicitly that whether the
model uses it was untested and unclaimed. This is that test, and the answer is
no.

One fixed city on one seed, one family driven into each of two states, eight
requests each way against qwen3:14b, everything identical except whether
`organization_money` was in the context:

| | scenes reaching for a word about money |
| --- | --- |
| told, family cannot pay | 2 of 8 |
| not told, family cannot pay | 3 of 8 |
| told, family comfortable | 5 of 8 |
| not told, family comfortable | 5 of 8 |

Being told changes nothing. It is very slightly the wrong way round for the
broke family, which at eight runs is noise rather than a finding.

The measure is weak and that has to be said before the result is trusted. It
counts words like wages, owe and short, which appear in perfectly ordinary job
offers that have nothing to do with a family's finances — which is almost
certainly why the comfortable family scores highest of the four. A better
measure would ask what the speaker wants rather than which words they used.

A hypothesis I have not tested: the full brief ships every faction record, every
person, every property and the whole recent history, so one small key is easy to
lose in it. The focused brief exists and sends far less. Comparing the two on
this same question is the obvious next experiment, and it is not evidence yet.

What this does not overturn: the four states are still the right thing for the
core to know, and everything in the game that reasons about a family can now ask
how it is placed. The claim that fails is only the one about the model.

The context assembly was extracted into `directorContext` so this comparison
sends the real payload rather than a copy that could quietly fall out of step
with it.

Evidence: `cmd/blackledger/director_compare_test.go`, which is skipped unless
`BLACK_LEDGER_COMPARE` is set because it needs a model and takes seven minutes.

## Four more businesses, and the guards that stopped me shipping them badly

From the inbox: expand the businesses. The city had three you could own — a
laundry, a garage and a casino. It now has seven. A restaurant with a back room
nobody books, a billiard hall where somebody is already running a book, a
butcher with a cold room and a licence to move refrigerated goods at any hour,
and a haulage yard with four trucks.

Each has its own trouble, because a business is only interesting while there is
a kind of trouble the player has not met. The restaurant fails an inspection.
The pool hall loses its own money to somebody's book in the back. The butcher's
cold room dies overnight and takes a week of stock with it. The haulage yard
has a driver talking to Ward Street.

I was wrong about the view and should say so. I read the asset code, found a
wireframe fallback with a comment saying a building added tomorrow would still
have a card rather than a hole, and concluded adding places was safe. It is not.
Three tests require every address to have a painted front, an interior, and a
block of its own on the map, and they refuse the fallback. Reading the comment
was not reading the wiring.

Five guards caught real faults, which is the whole argument for having them:

- Two addresses had no picture and no interior.
- Two of my coordinates landed on blocks already taken, so Vittoria's would have
  been pushed off Saint Agnes and the butcher off Russo Motor Works.
- Every place needs its own ways of dying, and four had none.
- "One man off the books" tripped the scan for prose that assumes a gender
  nothing in the game records.

Painting them turned up a fault in the art tools worth more than the paintings.
A place's seed was its position in a sorted list, so adding four addresses
shifted every index after them and repainted eleven finished buildings. The
seed comes from the place's name now, and both tools skip anything already
painted unless asked for everything. A building should look the same tomorrow as
it does today.

The balance simulation barely moved: one death fewer for the defiant strategy
and $147 more for the investor. That is not evidence the new businesses work.
Those campaigns run seven to thirteen days and never get near the $980 a haulage
yard costs, so the simulation hardly touches them. They are tested directly
instead, over the whole trade table rather than over the four, so a fifth cannot
be added without the test noticing: every trade must be buyable, earn something,
open already staffed and supplied, suffer when run down, and have trouble and a
remedy of its own. Both failure modes fail.

| | before | after |
| --- | --- | --- |
| ownable businesses | 3 | 7 |
| defiant deaths | 52 | 51 |
| investor median cash | $14,019 | $14,166 |

Evidence: `core/businesses_test.go`. All gates green, `npm test` 34.

## The businesses commit crashed the live game

Committed, restarted the live save, and it died on the first request. Not a
regression I found by reading: the server returned nothing and the log had a nil
dereference in the call that reads the world.

A campaign that predates an address has no record of it, and everything that
walks the location list finds nothing where a property should be. There is a
repair for exactly this, and it has been there a long time. It never ran,
because it sits behind `Version < SaveVersion` — and adding a business does not
change the shape of a save, only its contents, so nothing had bumped the number.
The live campaign was already at the current version and was therefore skipped.

The repair is not a version migration and is no longer gated like one. It runs
on every load. Bumping the version would have fixed today's crash and left the
next person to add a place to discover the same thing the same way.

Two things came out of it that are worth more than the fix. What a place earns
was a switch inside world creation, so an address added later was worth nothing
forever in a campaign already running — it is one table now, read by both the
new city and the repair. And a business arriving in an old save now comes
staffed and stocked, because a business is a going concern before anybody buys
it, which is already the rule everywhere else.

Verified against the save that actually crashed, not a fixture: it loads, and
the four businesses arrive earning 20, 16, 28 and 40 with their people and
stock in place.

Three tests, and I had to write a fourth. The first three sit in the core and
call the repair directly, so commenting out the call in the load path left them
all passing — the wiring was untested, which is exactly the seam that broke.
The store test loads a campaign at the current version with an address struck
out and asks whether it comes back. Both breaks now fail.

Evidence: `core/new_places_test.go`, `store/store_test.go`.

## A wheel, and an edge that can be stated rather than tuned

From the inbox: a proper game of roulette or blackjack rather than text buttons.
Blackjack mechanics already existed and are good. Roulette did not exist
anywhere — no wheel, no bets, no payouts, nothing in the core and nothing in the
view. Core is the only source of truth, so the wheel is built and tested before
anything is drawn.

The two games are worth having side by side because they are opposite. A hand of
cards is a decision the player keeps making. A spin is one decision made before
anything happens and then nothing to do about it.

The important part is the edge. Every payout is the true one: a number pays 35
to 1 against 36 other pockets, red pays even against eighteen others, a dozen
pays 2 to 1. Nothing is shaded to make the house win. The whole advantage is the
nought, which is neither colour, neither odd nor even, and in no dozen, so it
takes every bet on the outside of the cloth.

That makes the fairness measurable rather than asserted. Walk every bet over
every pocket exactly once: each returns 36 for 37 staked. The same for all of
them, which is what tells you nothing has been tuned. Three deliberate breaks
fail it — letting red take the nought, paying a dozen even money, and paying a
winner without returning the stake.

The red pockets are written out rather than computed. They are not every other
number and there is no arithmetic that produces them.

Two seams needed testing beyond the function. The bet rides on the command's
choice field, not its target: target names the room and is what the action
lookup searches, so a bet put there would send the game looking for a wheel in a
place called red. And the world sent the interface a hand and no wheel, so a
spin happened and nothing outside the ledger could say what the ball did.

Verified over HTTP on a driven save, not only in tests. Eleven spins across
every kind of bet, money moving the right way each time, and the colours
matching the real layout: four black, one red, eight black.

One thing I cannot explain and am not claiming. In the first run, the response
to the very first spin carried no wheel while the eight after it did. Three
deliberate repeats did not reproduce it. I have not established a cause.

Next: the table itself, for both games — cards, felt and a wheel.

Evidence: `core/roulette.go`, `core/roulette_test.go`. All gates green,
`npm test` 34.

## The tables are tables now

From the inbox, the user's own words: a proper game of roulette or blackjack
rather than text buttons, "think actual playing cards, table, etc". Both games
are now played on a felt.

The core had to change first, and the reason is the whole architecture. A hand
of blackjack was two totals and a count — the player showing nineteen, the
dealer showing ten. There was no such thing in this game as the nine of hearts.
Drawing playing cards from that would have been the view inventing facts, which
is the one thing it may never do.

So the deck has faces. A hand keeps the cards it was dealt, and what the hand is
worth is added up from them and from nowhere else, so the number on the screen
is always the sum of the cards beside it. That is a test rather than an
intention: it walks a hand card by card and compares the two every step.

Giving the cards faces uncovered a rules bug that had been there all along.
Softening looked only at the card just drawn, so an ace already in the hand
could never come down later. An ace, a five and a ten is sixteen at any table in
the world; this game called it twenty-six and took the money. It also meant
four aces could not all come down.

That fix moved the balance and it should be said plainly. The house edge over
twenty thousand hands went from 6.3% to 4.6%, because hands that used to be
taken as busts are now played out. It remains inside the tolerance the books
assume, and the cause is a rule being right rather than a number being tuned.

What is drawn: a felt with the dealer's row and the player's row, real cards
with rank and suit, and the dealer's hole card face down — because it genuinely
has not been dealt yet, not as decoration. For roulette, a wheel with the
pockets in the wheel's own order, which alternates colours and is not the
cloth's order and not numerical, and a cloth in three columns of twelve with the
dozens and outside bets along the bottom.

The view's own layout facts are tested: the cloth holds every number once,
eighteen red and eighteen black with ten and eleven both black and eighteen and
nineteen both red, the wheel face holds all thirty-seven pockets and alternates
colours all the way round, and every outside bet names an id the core accepts. A
card whose suit the core did not send is drawn face down rather than guessed.

Two things caught in the browser rather than in a test. The same game was
offered twice in one room, once on the felt and again as a text button below it;
the buttons are hidden only where the felt is actually offering that game, so a
refusal the player needs to read is never swallowed with them. And an existing
guard refused the file: it resets its hook count on `function Name`, and I had
written arrow constants, so it read a hook in one component as sitting below
another component's early return. Matching the surrounding convention was the
right fix rather than weakening the rule.

Played end to end in the browser. A straight-up bet on 14 for $50, the ball in
9 red, the wheel showing it and the cloth saying "gone". A hand dealt at the
same table: ten and nine of hearts against a jack of clubs, nineteen against
ten, with the hole card face down.

Evidence: `src/Tables.tsx`, `src/cards.ts`, `tests/cards.test.mjs`,
`core/cards.go`, `core/cards_test.go`. All gates green, `npm test` 41.

## Money presses on a quarrel

The last open piece of the inbox item, in the user's own words: a family with a
lot of money is stable but more of a target, and one with little money is
struggling, aggressive and unpredictable. The four states existed and nothing in
the world read them. Only the director's context did, and the director
demonstrably ignores it.

An organization that cannot meet its wages now presses a quarrel harder, because
being careful has stopped paying. One with a great deal of money and premises to
hold is worth moving on, which is the other half of the same fact. Money nobody
can reach is not a temptation: a rich family holding no ground adds nothing.

The measurement is the part worth reading. My first attempt ran two cities for a
month, identical except for money, and reported a clean result in both
directions before a line of the mechanism existed. Money decides which branches
a family takes, and every branch changes what is drawn from the world's shared
random stream, so the two cities were not running the same run of luck at all.
A month of simulation cannot see this. The pressure is a function of its own and
is measured directly, with no randomness in it.

Then the wiring, which is the fault this project keeps being caught by. My first
wiring test called the function directly and passed with the line in the city's
turn commented out. The real one runs two cities with the *same* cash on both
sides, so nothing that reads cash takes a different branch, and pins the world's
stream to the same value every turn. The only difference is whether payday was
missed. It fails when the connection is cut.

Applied whole, the pressure was a shove rather than a thumb: up to six against a
drift running from minus three to plus three, twice a day. Every city in two
hundred went to war, which an existing guard caught immediately. It is divided
now, so two desperate organizations move faster than two comfortable ones and
that is all it does.

Two other things the change surfaced. A family can now be destroyed inside a
month, and a balance test reached straight into a lookup for one that was no
longer there; that is the second-oldest fault shape in this list. And the
lowercase sweep found one more sentence, in the line the paper prints when a war
ends with somebody holding nothing.

I raised a balance test's sample from 120 campaigns to 600 while the pressure
was still too strong, and then reverted it: with the pressure divided, that test
returns exactly the numbers it returned before, 114 of 120 and the same median
attention. The sample was never the problem. Recording it because the reasoning
was published in a comment for a while and was wrong.

`mise run simulate` is unchanged at defiant 51 / investor 0 / reckless 82 /
worker 0 with $10 of movement in one median, which is expected: those campaigns
run seven to thirteen days and a family's money takes a month to say anything.

Evidence: `core/desperation_test.go`. All gates green, `npm test` 41.


## What the inbox asked for, and what came of it

Both entries the user wrote are now built, so they come out of the inbox and are
recorded here in their original wording. The rule is that an entry stays in the
inbox until it genuinely exists in the game, which is what keeps that list
honest.

> Something to add to the inbox, we should be tracking money that other living
> NPCs have, their families have etc. remember it's a living world and all the
> other living NPCs are just like us so we have to ensure that's all fleshed out.
> Then their money can also play into their decision making from the AI director.
> I'd imagine that a family with a lot of money would be more stable but also more
> of a target while a family with little money might be struggling and also more
> aggressive and unpredictable as they're trying to survive. This type of stuff
> >
> Following on from fleshing that out I'd imagine that there's a way to wage war
> against a family by starving it financially by decimating businesses or taking
> them over that they get money from and then they won't be able to pay wages for
> their people and it'll weaken them. This is a whole section of the game I want
> you to flesh out in a continuous loop. I also want you to start expanding the
> businesses available in the game. Also did we manage to build out the specific
> unique UI for gambling? So that looks like a proper game of roulette or
> blackjack you're playing and not just text buttons on a screen, think actual
> playing cards, table, etc

What was built, in the order it happened: families got a daily bill so their
money meant something; the starvation war was made to work, because a stripped
family used to fade on a timer that outran its money and so losing the
businesses caused nothing; everybody in the city got a purse of their own; the
director was told how each family is placed, and then measured and found to
ignore it; four businesses were added, each with its own kind of trouble; the
gambling tables were built with real cards and a real wheel; and money was made
to press on a quarrel.

Two of those are worth remembering for what they cost. Adding the businesses
crashed the live save, because the repair that gives an old campaign a new
address sat behind a version check that adding content does not trip. And the
money-pressure measurement reported a clean result in both directions before the
mechanism existed, because the two cities being compared were never running the
same run of luck.

## A family ends and its people are still on its books

The inbox is empty, so this is my own choosing. When an organization is
destroyed the city writes its obituary, drops it from the world and forgets it.
Nine people go on answering to it. Their record names a family that no longer
exists, everything that tries to look it up finds nothing, and they keep
standing in rooms working for nobody.

The fix was already written. `orphan` exists, its comment says it is for exactly
this — "so nobody is left answering to an id that resolves to nothing" — and it
was only ever called when the *player's* organization ended. Neither path that
destroys a family called it. That is a rule applied to one case and not to the
identical one beside it, which is now the twenty-fourth shape on the list.

Both paths call it now, and the ledger says what became of them, because an
organization ending is something that happens to people rather than to a row in
a table.

Two prose faults came out of it, both mine and both caught by guards that were
already there. My first wording put men on the street, and nobody in this city
has a gender the game ever recorded. And the line the paper prints when a war
ends with somebody holding nothing had no verb agreement, so a family called
somebody's people "is not holding anything any more" — I fixed the
capitalisation of that same sentence one slice ago and did not read the verb.
More families reach that state now, which is why it surfaced today.

Evidence: `core/orphans_test.go`, and both calls verified by removing them. All
gates green, `npm test` 41.

## A business belongs to a kind, so the city can hold two of them

From the inbox: multiple casinos, multiples of businesses, strip clubs. The
first of those was not a content change. A business's trade — how many hands it
needs, what it runs on, what goes wrong in it — was keyed by street address, so
the city could hold exactly one laundry and exactly one casino. A second would
have needed its own copy of the same rules under a different key, and the two
would have drifted apart the first time either was touched.

An address now says what kind of business it is, and the trade belongs to the
kind. What an address *earns* stays with the address, because that is genuinely
a fact about the premises: the Golden Lily is quieter than the Blue Hour and the
Ordway is bigger than the Bluebird.

Four new addresses, and for the first time two of them are second helpings
rather than new kinds: a second gambling house, a second laundry, a revue bar
where somebody outside is leaning on the dancers for a cut, and a cab company
with two cars off the road. The city has twenty addresses now and eleven ownable
businesses across nine kinds.

The sweep for rules written about an address where they meant a kind was the
real work. Two were live: a crate room could be built under one laundry and not
under the other, and the player's own laundry would press their clothes for
nothing while the second one would not. Both are about kinds now, and both fail
when reverted to an address.

Three tests walked the trade table as though its keys were places, which stopped
being true. They walk the addresses instead. That is the same property asked of
the thing that now owns it — there are more businesses in the city than there
are kinds of business — and it is a test change rather than a code change, so it
is recorded as one.

Guards caught four things without help: four addresses with no painted front,
four with no interior, four with no ways of dying, and one line of mine that
walked somebody out "between two men" in a game that records nobody's gender.

| | before | after |
| --- | --- | --- |
| addresses | 16 | 20 |
| ownable businesses | 7 | 11 |
| kinds of business | 7 | 9 |
| kinds with more than one address | 0 | 2 |

`mise run simulate` is unchanged at defiant 51 / investor 0 / reckless 82 /
worker 0, with two medians moving by tens of dollars. Still to come from this
inbox entry: many more people living in the city, and more ways a business
touches the simulation.

Evidence: `core/kinds_test.go`. All gates green, `npm test` 41.

## Half the city had nobody in it

From the inbox: a lot more characters living here. The number was the wrong
thing to look at first. The street trades — the people who are not in this
business and live here anyway — name the address each of them works at, and that
list was written when the city had ten addresses. It has twenty. Eight of them
had nobody in them at all: a butcher with no butcher, a cab company with no
drivers, a revue bar with nobody on the stage.

So the city is inhabited before it is enlarged. Twenty-three new trades across
those eight addresses, and the street is sixty-two people rather than
twenty-six, which is roughly one for each way there is to make a living here.
That was always what the number meant.

| | before | after |
| --- | --- | --- |
| people per city, day 30 | 53 | 87 |
| people per city, day 400 | 119 | 170 |
| addresses with nobody working there | 8 of 20 | 0 |
| ways to make a living | 26 | 49 |

Nobody shares a name with anybody else in a city of eighty-eight, and every city
reaches the street it means to have — both tested, because naming somebody is
sixty attempts at a random pairing and running out is silent: the city would
simply be smaller than it intended and nothing would say so.

Two guards fired. One of my new trades was a delivery boy and another a cold
store man, in a game that records nobody's gender. And the check on how many
people are on the street at once was written as "more than twelve, which is most
of the city" when the city held about fifty. Twelve stopped meaning a quarter of
anybody. It is a proportion now, which is what its own message always claimed it
was — a test change, and recorded as one.

The costs, measured rather than assumed. The core suite went from 45 seconds to
94: the prose scans read every passage the city writes, and a city with seventy
per cent more people in it writes more. The state the interface reads is 80KB.
The busiest room holds twenty people at once, and the room list already hides
bystanders behind a count, but I have not looked at a crowded room in the
browser yet and am not claiming it reads well.

Evidence: `core/population_places_test.go`. All gates green, `npm test` 41.

## Three newspapermen in one room

I published last slice that the busiest room now holds twenty people, that I had
not looked at one in the browser, and that I was not claiming it read well. This
is that look, and it turned up something a test would not have.

The room itself reads fine. The core already says "20 people in here, which for
this hour is a crowd", offers the seven you can actually deal with, and folds
the other thirteen behind a count. Opened, they are a two-column list of names
and trades. Nothing needed fixing there.

What the list showed was three newspapermen standing in the market together, and
two stallholders. A civilian's trade was drawn at random from the table each
time, with replacement, so sixty-two people across forty-nine ways of earning a
living gave four or five of one job while a third of the city's jobs had nobody
doing them at all. A room where everybody is the same thing is a room of one
person repeated, which is a shorter city than it looks — and it is exactly the
failure mode of raising a population without watching what fills it.

Trades are dealt now rather than drawn: the least-taken jobs are found and one
of those is chosen, so every way of making a living here is somebody's before
any of them is a second person's. Which of the open ones is still chance, so two
seeds are not the same city in the same order. Nobody's job is doubled more than
twice, and no trade goes unfilled.

I also owe a correction. I wrote that the suite went from 45 seconds to 94
"because the prose scans read every passage a bigger city writes". That is not
where the time is. Timing every test individually, the ten slowest are all
many-campaign balance tests — the slowest is eleven seconds and no scan is near
the top. The cost is that every simulated day now moves eighty-seven people
instead of fifty-three, spread across everything that runs campaigns. No scan
needs weakening, which is just as well, because they have caught real faults all
night.

Evidence: `core/population_places_test.go`, verified by reverting to the random
draw. All gates green, `npm test` 41.

## What a business is worth beyond its takings

The last part of the inbox entry: how businesses interact with the city
simulation. A trade decided what a business cost and what went wrong in it, and
nothing else, so owning a cab company was owning a butcher with different words.

A trade now says what it is worth as a front. Laundering absorbed exactly
fourteen points of police attention whatever the player owned, scaled by
condition and nothing else. It is the trade's number now: a laundry is where the
word comes from and absorbs eighteen, a casino handles more loose cash in a
night than a laundry sees in a week and absorbs twenty, a yard full of trucks
explains six.

The better find was underneath it. The gate asked what KIND OF ROOM a place was
rather than what trade was run in it, so a casino could not launder a dollar —
and a casino is most of the reason anybody owns one. That was two gates, not
one: the readiness function and the button that leads to it, written separately
and both asking the room. Fixing the readiness alone would have left the player
refused a button the rules said they could press, and never told why. There is a
test over every address that asks whether the button is offered exactly where
the trade allows it.

| the front | absorbs, in good condition |
| --- | --- |
| casino | 20 |
| laundry | 18 |
| restaurant | 16 |
| burlesque | 14 |
| poolhall | 11 |
| garage | 9 |
| butcher | 8 |
| cabs | 7 |
| haulage | 6 |

Two test changes, both recorded. An existing test asserted "a casino is not a
laundry" and refused one; that was the room type talking, and it now asserts
what survives — a rival's premises are not your books, a rented room with no
trade in it is not a front, and a casino absorbs more than a laundry. And I
wrote that last clause backwards first, asserting the laundry was better, then
corrected it against the design rather than the other way round.

The balance moved and it should: `mise run simulate` gives defiant $7,223
against $7,122 and worker $12,510 against $12,384, with the investor down $209.
Deaths are unchanged at 51 / 0 / 82 / 0 and there are no errors. Attention is
easier to clear for a player who owns the right business and harder for one who
does not, which is the point.

Evidence: `core/cover_test.go`, four breaks verified. All gates green.

## What you own decides how fast the city forgets you

The companion to what a business hides: what it costs you to be seen owning it.
A business was pure upside once bought — takings, and now cover — and what a
place *is* never cost the owner anything.

A revue bar with a late licence and a gambling house are watched in a way a
laundry is not. Attention now fades every night for somebody the city has no
standing reason to watch, and every second or fourth night for somebody who owns
the rooms people are seen going into. Over twenty nights, a laundry lets twenty
points fade and a revue bar five.

The shape of it was forced by arithmetic twice over, and both are worth writing
down because the obvious designs are both wrong.

It cannot be attention ADDED each day. Attention fades one point a day, so
anything adding more than that climbs without limit: past forty-five the police
arrive, past eighty they take the premises, and somebody who bought a burlesque
and did nothing else would lose it inside a month with no way to stop it. That
is the unbounded-figure fault, and I would have shipped it if I had not checked
the decay rate before picking numbers.

Nor can the nightly fade simply be made BIGGER for clean owners. I tried that
first — a base of three, reduced by what you own — and an existing test caught
it immediately: the fade is deliberately slower than a hard-run business
generates, so a base of three would absorb skimming and make the operating
decision free. Fading less OFTEN leaves the rate exactly right for a clean owner
and slower for a watched one.

The measurement was worthless the first time. I called the day ten times without
advancing the clock, and since cooling happens on the Nth night, that was the
same night ten times: it reported no difference at all between a laundry and a
revue bar. That is the fourth harness this month that measured itself.

One naming note. An operating mode already has a `Notice`, meaning what a family
notices about your takings, so the trade's field is `Watched`. Two fields named
Notice meaning different things in the same package is a fault waiting to be
written.

The attention guard still holds: a careless skimmer keeps the business in 116 of
120 campaigns rather than all 120, so forfeiture is still a real risk. `mise run
simulate` is unchanged in every figure, which is expected — those campaigns run
seven to thirteen days and rarely own a watched business at all.

Evidence: `core/notice_test.go`, two breaks verified. One existing test changed
with the reason: it compared against a flat constant and now asks the world what
tonight's fade is, because that is no longer the same for everybody.

## A place for things to sit without being looked at

The third shape of what a business is. Cover and Watched are numbers a trade
contributes; this is something the player can only do because of what they own.

The first candidate did not survive reading the wiring, and that is worth
recording. A cab company whose drivers see where people go sounds right, and it
would have been worth nothing: the city already tells the player where all
eighty-six people are, whether they are walking, where to, and how many minutes
out. There is no finding anybody to be done. Building it would have been an
action that revealed what the screen already said.

What the player genuinely cannot do is hold contraband without it being seen.
Attention accrues every day for every unit not out of sight, and out of sight
meant a false floor in a car or a cellar under the house. A yard full of trucks
and a cold room are places things sit without being looked at, and that is what
those trades are for. Five crates of moonshine draw five points a day with
nowhere to put them, one point behind a cab yard, and nothing at all behind a
haulage yard.

The measurement was wrong first time, in a way worth naming. I used eight crates
of arms, which is twenty-four points of attention before a cap of six — so every
case measured the cap, and a cab yard and no yard at all both came back as six.
Five crates of moonshine sits under the cap and the difference is plain.

Underneath it, another address named where a kind was meant: a car costs half to
keep when the player owns "the garage", written when the city could hold exactly
one. It asks the kind now, through `OwnsKind`.

That fix is currently latent and I am not claiming otherwise. There is still one
garage, so swapping the fix back out leaves the garage test passing — it cannot
be caught failing. What can be proved is the thing the fix rests on, and it is
proved over the kinds the city genuinely has two of: owning the Golden Lily
counts as owning a casino, and breaking `OwnsKind` fails that immediately.

| the trade | keeps out of sight |
| --- | --- |
| haulage | 7 |
| butcher | 4 |
| cabs | 4 |
| garage | 3 |
| laundry | 2 |
| restaurant | 1 |
| casino, poolhall, burlesque | nothing |

Evidence: `core/hides_test.go`, two breaks verified and one honestly reported as
unprovable until the city has a second garage.

## Cars come from somewhere now

From the inbox: a dealership the player and the city both buy from, so the city
links to itself. A car came from nowhere. The player pressed a button at the
motor works, the money left their pocket, and a car existed. Nobody sold it and
nobody was paid for it.

There is a forecourt now, and a second motor works under the viaduct. Buying a
car happens at the forecourt, and a quarter of the price stays with whoever
holds it: a family that owns the lot takes $155 on a sale, and a player who owns
it buys at $465 rather than $620, because the margin never leaves their pocket.

Twenty-two addresses, thirteen ownable businesses, ten kinds.

The change found a fault of the kind that only shows when a world grows. One
function answered two questions — where a car is SOLD and where a car is WORKED
ON — and nobody noticed while the answer to both was the motor works. Moving
sales to a forecourt moved servicing with them, and a garage could no longer
touch a car. Two questions, two functions, and its own comment had been telling
me for months: "a motor works, which is the one place in this city that has
any", which was true of a city with nowhere to buy a car.

Last slice I recorded that the `OwnsKind("garage")` fix was latent and could not
be caught failing, because the city had one garage and "the garage" and "any
garage" were the same thing. It has two now, so both garage rules are checked on
each of them alone, and both fail when reverted. That debt is paid.

Three existing tests were repointed with the reason: they bought cars at the
motor works, which was the rule, and the rule moved. One was renamed to say what
it now guards.

Balance unchanged in every figure, which is expected — those campaigns rarely
buy a car at all, let alone own the lot.

Not built, and it is the half the inbox actually asked for: people in this city
still do not own cars. Nobody buys one, nobody loses one, and the forecourt has
no customers but the player. `LoseCar` exists for the player alone. That is the
next slice and it is where the link becomes real.

Evidence: `core/dealer_test.go`, four breaks verified.

## The city drives

Both car entries in the inbox stand on one thing that did not exist: nobody in
this city owned a car but the player. The forecourt had exactly one customer,
there was no car to steal parts from, and a garage had nothing to repair. That
is why this slice is the foundation and not the feature.

People own cars now, and who drives is who could afford to. Fifteen of
eighty-six on day one: every official, and everybody from soldier upward.
Nine of nine people of standing drive and none of the sixty-five on the street
do, which is what a 1950s city looks like.

And the forecourt has customers. Somebody who would drive, has none and can
afford one buys, and whoever holds the lot takes the same quarter the player
pays — the same transaction seen from the other side. One sale a day at most,
because a city where everybody replaces a car on the same morning is a city
where nothing was ever taken from anybody.

Two things I built deliberately that are worth naming.

A car carries the minute it was got, not just a tier. Without it the settling
pass cannot tell somebody who never had a car from somebody whose car was taken
last night, and would quietly hand the second one a replacement. That is tested:
take a driver's car away, settle the city, and they are still walking.

And the repair runs on every load, like the one for new addresses. A campaign
that predates cars would otherwise have no driver in it for the rest of its
life, and the forecourt would sell nothing forever.

Three breaks verified, and all three are seams that have caught me before: the
day not calling the trade, the margin never reaching the lot's owner, and
everybody in the city driving.

Balance unchanged in every figure. Those campaigns are seven to thirteen days
and this is a slow trade.

Next, and now possible for the first time: cars being taken and wrecked. A raid
already reaches the people at a place, so the demand the inbox asks for has
somewhere to come from — and once cars are taken, stolen parts have a seller and
a garage has a reason to want more of it happening.

Evidence: `core/citycars_test.go`.

## A car in pieces is four links in one act

From the inbox: people steal car parts and sell them to garages, and a garage
does better when there is more of it about. Both halves needed a car that
belonged to somebody, which the city got last slice.

Taking a car apart now does four things at once, which is the point. Whoever
drove it is walking. The parts have a buyer and pay $95 a tier, half again with
a garage of your own to take them to rather than selling them on at whatever is
offered. Every garage in the city picks up trade, because more cars going to
pieces is more work about. And the person you did it to holds it against you
personally, so a man who lost his car is a man who comes looking later.

The forecourt closes the circle without anything new being written: he is
somebody who would drive and has none, so in a week or two he buys another and
the lot takes its quarter.

Two things I got right by reading rather than guessing. A garage doing better
needed no new counter — `Custom` already multiplies a place's earnings and is
capped, so this cannot run away. And my first version made the owner resent
whoever was standing there rather than the player: `Resent` records a grudge
between two people in the city, and `Aggrieve` is the one that reaches the
protagonist. The test caught it as a grudge of zero.

Four breaks verified, including the command not being wired, which is the seam
that has caught me more often than any other.

Balance unchanged in every figure. The simulated strategies never strip a car,
so this is invisible to them, and that is expected rather than reassuring.

Not built: a scrapyard, and the city taking cars off each other in raids. A raid
already reaches the people standing at a place, so the second is a short step.

Evidence: `core/parts_test.go`, `core/parts.go`.

## Cars are destroyed during operations, and the yards take them

The inbox asked for cars destroyed during operations, and for a scrapyard that
links in the way garages do. Both are in. A raid already reached the people
standing at a place; it reaches what is parked outside them now. Whoever loses
one is somebody who would drive and has none, so within a fortnight they buy
another and the forecourt takes its quarter. A war is good business for
anybody holding a lot.

Devlin Salvage is the eleventh kind: six acres of what the city used to drive,
stacked four high. A car going to pieces lifts a scrapyard's trade the way it
lifts a garage's, through the same capped mechanism, so neither can run away.

The measurement was wrong first, and this one is worth recording because it
passed convincingly. My first version counted cars among a family's surviving
members before and after a month of war, and reported twenty-two lost before a
line of the mechanism existed. Somebody killed in a raid leaves the member list,
so it was counting deaths. Following named people who are alive at both ends
gives the honest figure: three of fifty-one survivors lost what they drove.

Twenty-three addresses, fourteen ownable businesses, eleven kinds.

Balance unchanged in every figure, for the same reason as the last three slices:
the simulated strategies never buy or lose a car. That is expected and it is not
reassuring, and I keep saying so because it would be easy to read as evidence.

Still open from the inbox, and it is the good half of a line I have not used:
"more car repairs to be made from broken windows from theft". Stripping takes a
car outright, so a garage gets parts but never repair work. A car left damaged
rather than gone would give the owner a reason to visit a garage, which is a
better link than the one I built.

Evidence: `core/wrecks_test.go`, two breaks verified.

## The inbox, and how it is read

`docs/LIVING_WORLD.md` now separates three things that used to sit in one list:
standing instructions about how the work is done, open entries the game does not
answer yet, and answered ones. A standing instruction can never be finished and
must never be ticked off; an open entry stays in the user's exact words until
the game actually does it; an answered entry keeps its original wording and a
line saying where the work went.

Every line I write in those sections begins with an em dash. Anything without
one is the user's own words, and they are never edited, summarised or split
across sections — a message that carries both an idea and a standing instruction
stays whole, where its main purpose puts it, and my line says so.

## Money moves in the figure you type

Four actions moved money in lots somebody else chose: $250 behind the tables,
$250 off them, $500 out of the city, and the whole offshore account home in one
go. The inbox asked for the opposite — "as dynamic and user settable as
possible" — and the gambling side already had the shape, because a stake has
been typed into `Command.Amount` since the tables were built.

So `Bankroll`, `Draw`, `Deposit` and `Withdraw` take an amount. Zero still means
what the fixed lot meant, which is what keeps every existing caller, save and
simulation honest: nothing that never names a figure behaves differently.

The bounds belong to the core, not the panel. `Action.Sum` carries `Least`,
`Most`, `Preset` and a label, and the panel renders a number field from it and
sends what is in it. So the withdraw field caps at what is actually out there,
the draw field caps at the float, and the funding field caps at cash on hand —
none of it computed twice in two languages.

The card had to stop being a `<button>`. A number input inside a button is not
clickable and is not valid HTML, so an action carrying a Sum renders as
`SumAction`: the same card, a field, an "All" shortcut, and its own commit
button that prints the figure it is about to move.

| Action | Least | Most | Starts on |
|---|---|---|---|
| bankroll | 25 | cash in hand | 250 |
| draw | 25 | what is behind the tables | 250 or the float |
| deposit | 100 | cash in hand | 500 |
| withdraw | 25 | what is offshore | all of it |

One test change, and the reason. `TestEveryPricedActionKeepsItsPrice` required
every action that pays its own fee to print a price on the card. An action whose
figure the player types has no one number to print before they have said how
much, so the guard now exempts a card carrying a Sum from that and checks its
bounds instead. The half of that guard that matters — `Cost` must be zero or the
money leaves twice — is unchanged and still covers all of them.

Verified by breaking it: `BankrollSum` made to ignore the typed figure fails
both the amount test and the refusal test, in the right direction.

Evidence: `core/typed_amounts_test.go`, one break verified. The browser
playtest did not happen this slice: the Chrome extension is not connected, so
the card was exercised over the API instead — $1,337 behind the tables moved
cash 3150 → 1826 and the float 0 → 1337, and the draw field's ceiling followed
it to 1337.

## A sitting is something the world knows about

From the inbox: "merely sitting down to play the games already started playing a
game, which is wrong."

It was, and the cause is worth writing down because it will happen again in
another shape. The felt is saved state. The last hand, the last spin and where
the three drums stopped are all kept on the world, deliberately, so a game can
be drawn rather than described. Whether the player was *at* the tables was not
kept anywhere: it was a boolean in one React file. So opening the takeover drew
whatever the last sitting had left behind — cards face up, a ball in a pocket,
three drums showing a line — and that reads exactly like a game that started
without you.

The view was not wrong about anything. It had nothing to be right with.

So sitting down is a decision the core owns. `World.Seated` names the room,
`Sit` takes the seat and clears the table, `Rise` gets up and clears it again,
and `RiseReadiness` refuses while a hand is live because the money is already
down. `Advance` calls `LeaveTable`, so walking out of the room ends the sitting
whatever moved you — a journey, a scene, an errand somebody else ran.

The room offers a seat rather than a game. "Sit down at the tables" is an action
the core writes, with "Play the machines" in a room that has a bandit and no
felt, and the panel's button commits it instead of flipping its own flag. The
takeover opens on the game its own first tab names, which is why sitting down
used to show the wheel under a nav that said Blackjack.

Verified over the API on an isolated fixture, four behaviours in one run:

| Step | seated | wheel | drums |
|---|---|---|---|
| sit | casino | clean | clean |
| spin, then pull | casino | shown | shown |
| rise | none | clean | clean |
| sit again | casino | clean | clean |

Getting up mid-hand is refused with "Finish the hand first", and travelling to
the bar left seated empty.

On the other half of the same message, "you changed it back to play the nickle
machine": that label went in 5f5e30d, which is before the message was written,
and `slotStakes` is now referenced by nothing the game can offer. A stale
bundle in the browser is the only way to still see it. Saying so is not a
dismissal — it is the thing to check first next time a fixed label reappears.

Evidence: `core/sitting_test.go`, one break verified — disabling the clear in
`clearTable` fails the sitting test in the right direction.

## A hit comes to where you are

From the inbox: "Can people only make attempts on your life while you're at
home? They always seem to hit my home when I'm not there and they are coming
after me."

They could not. `Attack` returned early unless the player was standing in their
own residence, took forty-five condition off the house and left. So anybody who
spent the day working was unkillable, and the punishment for having a life was
a broken door.

The rule now asks where you are and what that place is worth.

| Where | What it is worth |
|---|---|
| Home | the door, as before |
| A room, per stranger in it | 5%, up to 25% |
| A room, per one of yours in it | 10%, up to 20% |
| Any room, in total | never more than 35% |
| The street between two addresses | nothing, and 10% against you |

A room also warns you when it is busy enough — three people, enough that a
stranger walking in with his hand in his coat is something somebody says out
loud — and a warned player gets the escape, defend or bargain scene wherever
they are standing rather than only at home.

Two places they cannot walk into, and both leave them the house: a police cell,
and anywhere outside the city. The second needed the world to say so. A walk
between two addresses and a week in Halloway both parked the player in
"transit", which is how the first version of this quietly made leaving town
worthless — the test that measures what a journey buys went to nought of three
hundred. `World.Abroad` names the destination while the player is out of the
city, so the two cases stopped being the same string.

A warning still buys what it always bought. Somebody told you they were coming
and told you to avoid home, so that is where they go: being elsewhere costs you
the door instead of your life, which is a decision rather than a free pass.

**Two test changes, both deliberate.** `TestAbsentPlayer` asserted the exact
behaviour the inbox complained about — a player in a bar comes home to a house
at 55 — so the test changed and the code did not. It now asserts the opposite,
and `TestWarningAllowsLeavingBeforeHit` still guards the house-instead case that
a warning creates. `TestNobodyInThisCityHasAGenderTheGameNeverGaveThem` caught
"Two men come through the door" in the new scene, which is the guard doing its
job; it reads "Two of them" now.

**Measured, three arms rather than two.** Counting survivors of a warned attack
measures the warning, not the room: a scene raised is a player still standing,
and that test passes whatever cover is worth. So the measurement uses a room
below the number of people it takes to warn you.

| Where, 400 unwarned attacks | Survived |
|---|---|
| A bar with two people in it | 131 |
| An empty bar | 90 |
| The street | 48 |

Deleting cover collapses the first two to 90 and 90; deleting the street penalty
collapses the last two to 90 and 90. Each half of the rule fails the test on its
own, in the right direction.

**Balance is unchanged**, and here that is a real measurement rather than an
absence of one. Deaths 50/0/82/0 and median cash 7223/13903/90/12360, identical
to the recorded baseline. Making an unwarned attack certainly fatal moves
reckless from 82 deaths to 100 and nothing else, so the simulation does reach
the new path — it is just that the strategies which get shot at are shot at
home.

Evidence: `core/ambush_test.go`, `core/ambush.go`, two breaks verified.

## A place remembers being robbed

From the inbox: "It seems like you can rob places or take from people's cars
multiple times in a row, that should probably be tracked and time limited etc."

You could. Nothing about a robbery was remembered by the place it happened to,
so the same till could be emptied every forty-five minutes for as long as
somebody felt like standing there.

What stops it is not a cooldown bolted onto a button. It is that the place
changes. A shop robbed on Tuesday does not leave Wednesday's takings in the same
drawer, and a row of cars that lost one to a man with a bag of tools gets looked
at from the window afterwards. Two days for a till, one night for a street, and
both wear off — a place that can never be robbed again has been deleted rather
than defended. A failed attempt counts the same as a successful one, and if
anything counts for more: somebody shouted, and everybody in there remembers.

The refusal says when it will be worth trying again, in the words somebody would
use rather than in minutes, and the panel says what a passer-by would see:
"Wary since the last time, and keeping its money elsewhere."

**The car half was already true and I had not noticed.** Taking a car apart sets
the owner's car to none, so they cannot be stripped again; and `sellCarTo` puts
them back on the road when they next walk onto the forecourt with the price in
their pocket. Mugging is the same shape already — the purse is emptied and
refills on its own. So the only two gaps were the till and the street.

**The simulation could not see any of this, so I gave it eyes.** No strategy in
`sim` ever robbed anything, which means the rule could have been added or
deleted and every number in the report would have stayed where it was. There is
a `thief` policy now. Its first version died in six commands — it robbed from
the first minute, took two beatings and never built anything — which measures
nothing, so it builds the ordinary way first and starts taking tills once it has
a name and somebody on the door.

It also needed the city to say what it knows. `Actions` for a room the player is
not standing in returns travel and nothing else, so a policy could not tell a
cold till from a warm one and just walked back and forth between two addresses
until it died. The snapshot carries `shy` and `curtains` per location now, which
is the same fact the panel prints.

| thief, 100 runs | deaths | median cash | median days |
|---|---|---|---|
| With the rule | 82 | 586 | 2.8 |
| With `Shy` returning zero | 86 | 443 | 2.5 |

A thief who can stand in one room and take the same till all night takes more
beatings, dies more often and ends with less. The rule is worth having rather
than merely restrictive.

**New balance baseline, five strategies**: deaths 0/0/50/82/82 and median cash
12360/13903/7223/90/586 for worker, investor, defiant, reckless and thief.
The first four are unchanged.

Evidence: `core/cooling_test.go`, two breaks verified, plus the simulation
figures above.

## Work that says what it is, and sits where it belongs

Four complaints from the inbox, none of them about what an action does and all
of them about how it is labelled or where it is filed.

**"Is buying a business called 'establish protection'?"** It was, and it should
not have been: this is buying the premises outright, and the game has a
protection racket elsewhere that the label described instead. It reads "Buy The
Green Baize" now, and the description says what you are taking on — two
positions to keep filled at five dollars a day each, cloth and chalk and drink
to buy in, repairs when it needs them — rather than only what it earns.

That change failed six existing tests, correctly. I had put the price into the
description, and every address that is not for sale carries a cost of zero, so
the city filled up with "$0 for the freehold". The panel prints the price from
`Cost` already. This is the fault shape the docs list under "a figure that says
nothing while looking like one", and the guard written for it caught mine.

**"Why does it seem like you can send Leo Carver on collections in practically
every single building's action menu?"** Because it was offered in every one. An
order to your own man is work that belongs to you rather than to whatever
counter you happen to be at, and `Action.Anywhere` has existed for exactly this
since the placement audit. Sending him on collections and paying him a bonus
both follow the player now, and the interface files them under the person.

**"I don't think moving against a business yourself should require respect."**
It does not any more. Nothing about a reputation stops a man walking into a
warehouse with a crowbar, and what decides whether he gets out again is his
crew, their loyalty and the family's strength — all asked already. A name is for
other people doing things on your account, so the requirement moved to
`SendAgainstReadiness`, which is the version where you are asking one of your
own to go in without you.

**"When inside a building you can't see who the family that owns it is."** The
street panel says it and the room did not, so stepping inside lost the fact that
decides how everything else in there should be read. The room names itself and
its holder at the top, with its condition and how it is trading.

Balance unchanged across all five strategies. Three breaks verified, one per
core rule; the wording is checked by asserting the label and what the
description has to mention.

Evidence: `core/placement_test.go`, and the six guards that caught the zero.

## The gate stopped costing three minutes

From the inbox, filed as a standing instruction: "Ensure efficiency of
development loops by increasing efficiency of your workflow in any way that you
can accomplish."

Two things were paid for on every single tick.

**The core suite ran everything in sequence.** 945 tests, none of them parallel,
and the fifty-four balance tests are most of the time: the slowest is nearly
sixteen seconds on its own and the top ten come to about a minute. They are also
the tests with the least reason to be sequential — each one builds its own
worlds from its own seeds and none of them writes package state. `t.Parallel()`
on all fifty-four takes the suite from 140s to 90s.

That is worth checking properly rather than trusting: three consecutive clean
runs, and a full `-race` run over the parallel suite, also clean. Cutting
iterations would have been the other way to make balance tests faster, and it
would have bought speed with evidence.

**`verify` ran its steps one at a time**, so the machine sat on one core for
three minutes while the core suite went and then started the type check. There
is a `gate` task now: the core suite and the API suites start in the background,
and the format check, vet, `tsc`, the build and the node tests run while they
go. The whole gate finishes in the time the core suite takes on its own.

| | Before | After |
|---|---|---|
| `go test ./core` | 140s | 90s |
| A full gate | ~200s serial | 90s |

`verify` is unchanged and still there for anything that wants the steps in
order. Nothing about what is checked has changed — this is the same work, done
at the same time as itself.

One more thing the gate was not doing: `gofmt -l` prints names and exits zero,
so a check that has quietly listed eleven unformatted files for a while passed
every time it ran. It fails now, and the eleven are formatted.

## A lead has a table, and the chair holds whoever is sitting in it

Two complaints in one message, and they turned out to be the same fault.

"Vitor Bellendi is always at the kessler filling station for some reason." He
was, and so was every other family lead, wherever the city happened to leave
them. `keepsPost` exempts a lead from the evening — "the leader of a family is
not found propping up a bar" — and nothing else ever gave one a reason to walk.
So a lead who went out once for petrol stood at that pump for the rest of the
game. A lead holds court at the family's own seat now, its best-earning ground,
and goes back to it.

| Over 40 cities, a week each | Hours at their own seat | Elsewhere | Leads who drifted |
|---|---|---|---|
| With the rule | 11,808 | 1,632 | 0 of 80 |
| Without it | 5,722 | 7,718 | 71 of 80 |

"I request a sit down with the controlling family and it sits me down with the
family lead but the lead is not present at this location." The audience was
addressed rather than peopled: the club meant Bellandi and the garage meant
Russo, whatever either family was actually doing, and the chair held the lead
wherever they were. It happens where somebody who can settle terms is standing
now — the lead when the lead is there, otherwise the most senior of theirs in
the room, and nobody below a lieutenant is authorised to agree to anything. A
lieutenant sitting in for an absent lead says so in the scene. Asking in a room
with none of them in it is refused, and the refusal says where their lead is.

That also unhooked the sit-down from two addresses. It is offered anywhere the
family's speaker is, which is what makes the first fix visible: you go to where
the lead is, and the lead is somewhere that means something.

**Four test changes, all for the same reason and all stated.** Four tests opened
an audience at an address and assumed a family, which is exactly the assumption
being removed; they seat a speaker in the room first. A fifth, about a holding
changing hands, looked only at who was still walking — and the person with the
best reason to go is now the lead, whose seat is that ground, so she sets off
first and has arrived by the time the last of them leaves. It accepts somebody
standing in it who was not there before.

Balance: deaths 0/0/50/82/83 against 0/0/50/82/82, median cash
12378/14050/7230/90/586 against 12360/13903/7223/90/586. One extra death in a
hundred thief runs and cash within a percent, which is the city's people
standing in slightly different rooms rather than a change in what anything pays.

Evidence: `core/seat_test.go`, three breaks verified.

## The dice, and a lesson about sample size

From the inbox: "It may be worth adding a couple more casino games."

Craps, which is the game the period actually had, and which is worth adding for
a reason beyond being a third thing to press. The three games in that room are
three different shapes of decision. A hand of cards asks you something on every
card. The wheel asks once and then there is nothing to do. The dice ask once and
then make you sit through a run of throws nobody at the table can affect — the
point is on, and it is the point or the seven, and everybody watches.

Three bets, all real: the pass line, the don't with the twelve barred, and the
field. As with the wheel, the edge is arithmetic rather than a fudge factor. The
pass line wins 244 times in 495, which is 1.414%, and the don't is the same game
the other way at 1.364%.

**Measured, and the first measurement was wrong.** At 200,000 decisions the pass
line came out at 1.75%. That looked like the world's linear congruential
generator failing on consecutive draws — the lattice problem, and craps is the
only game in the house that resolves over a sequence of rolls rather than one.
So I mixed the draw with an avalanche step, and it came out at 1.12%. Then the
same measurement at five million decisions:

| | Edge over 5,000,000 decisions |
|---|---|
| Two plain draws | 1.402% |
| One draw, mixed | 1.431% |
| The truth | 1.414% |

Both inside one standard error. The standard error at 200,000 decisions is
0.22%, so the number that sent me looking for a fault was 1.5 of those from the
true value, which is nothing at all. The generator was never the problem. The
mixing is gone and the test takes two million decisions, where the error is
0.07% and the claim it makes is one the sample can support. The single-roll
distribution was right the whole time — sevens at 16.64% against a true 16.67%
over 600,000 throws — which should have been the clue.

**A break that broke nothing.** Flipping `Won` on the come-out craps branch left
the pass-line edge identical to the last dollar, because `Won` is a label and
the money comes from the argument to `settleDice`. The test that noticed was the
one checking what the table says, not the one checking the edge. Breaking the
payout instead moved the edge to -20.8%, which is the break in the right
direction.

The table draws real dice with real pips and a point box that lights when a
number is on. The stickman's call is period and conditional: a three is "craps"
on the come-out and just a three once a point is on, and the totals are spoken
rather than printed — "8, 8" is a table read out by a machine.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/dice_test.go`, one break verified after a first one that did not
bite, and the whole flow driven over the API on an isolated fixture.

## The cloth takes as many chips as you put on it

From the inbox, three things in one message: the betting layout belongs beside
the wheel rather than under it, the chip should be worth whatever you type up to
the house limit, and "you can also place multiple bets in roulette, on different
numbers, combinations etc, like the real game by putting down chips on each one
you want to bet on."

The middle one was already built. The other two were not, and the third is the
one that matters: one bet a spin was never how the game works. What a table
takes is a cloth covered in chips, every one of them settled against the same
pocket, and that is most of what makes roulette a game rather than a coin toss
with extra numbers.

`Chip` is a bet and what is on it. `SpinChips` takes a whole cloth, charges the
total before the ball drops, settles each chip against the one pocket and fills
in what each was worth. `PlayWheel` is now one chip in a slice of one, so every
existing caller, save and test behaves exactly as before — and a spin still
reads by `Bet` and `Down` for anything written before the table took a cloth.

The house limit is per bet, the way a real table's is: two chips at the limit
are two bets and both stand. What is refused is more on the cloth than the
player has, however it is spread about.

On the table itself, clicking a spot lays a chip down, clicking again stacks
another, and right-clicking takes one off. After the ball drops the panel lists
every chip and which of them came in, rather than describing one of them. The
cloth moved to the right of the wheel, which is where a real one is and which
also stopped the spin button sitting below the fold on a short window.

Verified over the API on an isolated fixture: $20 on red, $5 straight up on 17
and $10 on even, all in one spin. Pocket 9, red — red paid $40 back, the other
two went, and the player finished $5 up on $35 down.

Two breaks verified: settling every chip against the first chip's bet, and
letting the total past the player's cash.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/chips_test.go`.

## Plate, and the one place it is worth anything

From the inbox: "You should probably also be able to outfit your car with
protection like armor etc at a garage which will help you survive attacks when
traversing out in the streets."

The top of the range has been called an armoured Packard since the day the list
was written, and no rule in the game had ever read that word. The plate was a
sentence in a description.

A garage fits it now, in two stages — the doors, then glass that has stopped
things — at $900 and four hours each. It is fitted to the car rather than to the
person: a new one off the lot is bare, whatever the last one was carrying, and
the Packard is sold with both stages on because that is what it is.

What makes this worth building rather than a number going up is where it
applies. The street between two addresses is the one stretch of the city with no
walls, no door and nobody who knows you — which is exactly what the ambush work
established a few slices ago, and exactly where a car is. Plate is the only
cover out there, and it is worth nothing at all to somebody who walks in a door
after you. A car parked outside a bar protects nobody in the bar.

| 600 unwarned hits on the street | Survived |
|---|---|
| No plate | 71 |
| Doors plated | 149 |
| Doors and glass | 225 |

It costs speed. Plate is weight, and each stage takes 9% of what the car was
worth as a car — so the fastest thing in the city is a bare Ford and the safest
is a Packard you cannot hurry in.

Two breaks verified, one per half: cover on the street, and drag on the pace.
The gender guard caught "a man who walks in a door after you" in the new
description, which is the second time it has earned its place this session.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/plate_test.go`.

## A car for one of your own

The other half of the same inbox line: "We could probably also extend to be able
to provide cars, and armor the cars for people in our family to keep them more
protected from attacks."

What makes this worth building rather than a number on somebody else's sheet is
where it lands. Sending one of your own after somebody is the one decision in
the game whose entire point is that it is not you taking the risk, and when it
goes wrong there are three endings: they are killed, they are taken alive and
your name comes out of it, or they get out with nothing. The difference between
the second and the third is whether there was something running at the kerb.

| 500 jobs that went wrong | Got away |
|---|---|
| On foot | 120 |
| A car at the kerb | 175 |
| A plated car | 231 |

Buying is on the forecourt, for anybody of yours standing on it with you, at the
lot price — less the margin if the lot is yours, which is the same rule the
player's own car already followed. Plating is at the garage, the same two stages
and the same price as your own. Both nudge their trust: somebody paying for that
is not a thing people forget.

Two breaks verified, one per half: the car, and the plate on top of it.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/theirs_test.go`.

## People say what they are carrying

Two inbox lines, and they turned out to be one system.

"I don't see a lot of cursing from characters in this game, we should increase
that since it's with the mafia style. Characters should be able to make threats
to you too, I have not seen that yet."

"I see 'bad blood' red box at the top of the page but it never goes away and
it's annoying."

The city already knew who had a reason to dislike the player, how badly, and
what it was about — `n.Sore` and `n.SoreAt` have been committed state for a
while — and did nothing with it but sort a list and draw a banner.

Somebody carrying something against you, standing in the same room as you, says
so. Three levels, and the closer they are to the weight at which people actually
move the less polite it gets. Temperament shapes it: a careful man does not make
speeches about it, a vain one wants the room to hear, a loyal one mentions his
people. Of 498 lines from people past that weight, 289 carried an oath.

Nothing here invents a grievance. Every line is built from the weight the
simulation committed and the reason it recorded, and only from people the player
actually knows — a stranger with a grievance is a stranger, and this game does
not put words in the mouths of people it has not introduced.

**And the banner was a standing state drawn as news.** Bad blood fades a point a
day, so a grudge heavy enough to be worth mentioning sat at the top of the page
for six weeks. What is worth hearing is that two people have just fallen out, so
it is said for three days and then stops. The grudge goes on being true
underneath — the test checks that the headline going away does not throw the
grudge away with it — and the person holding it now tells the player themselves.

**A break that did not bite, and the fix.** Deleting the worst of the three
levels left the test passing, because it only asked whether the top of the scale
differed from the bottom and the top fell through to the middle. It asks for
three distinct voices now, and the same break fails it.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/threats_test.go`, two breaks verified.

## The street is not scenery any more, so the bar says what it is

From the inbox: "the little bar that explains that you're traveling between
buildings is at the bottom of the page often below the fold."

It was, and it said the name of the place and a number of minutes. That was fine
while the street was scenery. It is not scenery now: the ambush work made the
street between two addresses the one stretch of the city with no walls, no door
and nobody who knows you, and the plating work exists entirely because of that.

`Crossing` is the core stating what a journey is before the player sets off:
where to, how long, on foot or driving, how much plate is on the car, and
whether anybody with a price on the player is known to be looking. The bar sits
at the top of the city pane where the eye already is, and goes red when somebody
is out looking for you.

The map half of that message is still open.

**And a piece of bookkeeping worth writing down.** The inbox entry about typing
your own stake read "not started" for far longer than it was true — the work
landed in 5f5e30d and I never moved the entry. Checking it properly turned up
that half of it genuinely was not finished: the holder's button to set the house
limit was offered with no field on it, so it sent an amount of nothing and
`SetLimit` refused it every single time. A button that cannot be pressed
successfully is worse than no button, because the player spends the trip finding
out. It has a field now, bounded by the same floor and ceiling the rule uses.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/crossing_test.go`, two breaks verified.

## What a car is worth, said in minutes, with a picture of it

The last two things on the cars message: "nice car images to show what you're
buying and stats information about the speed of the car relevant to what it
gives to you."

The forecourt quoted a percentage — "journeys take 74% of the time they take on
foot" — which is true and tells nobody anything, because nobody walks a
percentage. `CarWorth` times the longest walk from where the player is standing,
both ways, and the card names the road: "Pier 14 from here is 40 minutes on foot
and 30 in this."

The longest walk, because that is where a car earns its money. And the figure
includes the weight of whatever plate is on it, using the same arithmetic
`World.Pace` runs on the car the player owns — a quote that does not agree with
what happens is a lie told slowly.

Three cars were also three lines of text that looked identical on the way past.
`tools/cars.py` paints them on the same contract as every other picture in this
game: generated offline, shipped as files, nothing at runtime depending on a
model. The action carries the tier it is selling so the panel draws that car
rather than parsing the label.

Two breaks verified: the plate weight dropping out of the quote, and the road
never being named.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/worth_test.go`, `public/art/cars/`.

## The paper prints a picture of who it is about

From the inbox: "We should try to improve the images being displayed on the
newspaper. Having a portrait of an affected person or building would be great."

The cut was a drawn silhouette — a generic man, or a generic pair of houses —
picked by the story's kind. Meanwhile the city has a painted face for everybody
in it and a painted front or cut-out for its addresses, and the story already
carries what it is about: `Subject` names the person or the place by id.

So the cut is that picture, screened and inked: greyscale, hard contrast, and a
dot grid over it, which is what a halftone is. It reads as something a press
ran rather than a photograph pasted into a 1930s page.

**One thing to get right.** The first version laid the picture over the drawn
plate with multiply blending, keeping the plate as a frame. That does not work:
the plate's silhouette is solid black, so the generic man would have shown
through the actual man's face whatever the blending mode. The two are never
drawn together now — a picture when there is one, the plate when there is not,
and the frame belongs to whichever is showing.

No core change and no new art: this is the paper using pictures the game already
ships.

Evidence: the live edition carries person and place subjects with real ids, so
both paths are exercised by the current save.

## A guard that could not fail, and the hundred and sixty actions behind it

From the inbox: "When inside a building you own the top buttons should probably
be for owner management and under a separate subtitle for management actions."

The block was already at the top. What it had no subtitle saying was what it
was, and inside a place of the player's it reads "Running The Blue Bird" over
"staff, stock, repairs and what the house takes" now.

Going to look at what belonged in that block found something bigger.
`TestEveryActionBelongsSomewhere` asked whether every offered action's group is
one the interface renders. `GroupOf` falls back to "work" for anything nobody
classified, and "work" is a group the interface renders — so the test could
never fail, whatever was forgotten. It is the fault shape already written down
in the loop's own notes as "a guard that cannot fail", and it had been sitting
in the file that exists to prevent exactly this.

What was actually filed under "Jobs that pay today":

| Action | Belongs |
|---|---|
| `strike:*`, `send:*` — all the violence in the game | the street |
| `limit`, `draw` — the house limit and the float | your premises |
| `plate`, `fill`, `service` — the car you own | standing |
| `strip` — taking a car apart | the street |
| `sit`, `rise`, `play`, `pull`, `wheel`, `dice`, `roll`, `hit`, `stand` | the tables |
| `car:*`, `plate:*` — a car for one of your own | people |

A hundred and sixty-odd action ids in a full campaign, most of them the two
violence prefixes. There is a "tables" group now, `draw` stopped being filed as
work that can go wrong, and everything about the player's own car sits with
buying one.

`Classified` reports whether a group was chosen rather than defaulted, and
`TestNoOfferedActionIsThereByDefault` fails on the fallback. The old guard stays
— it checks a different thing, that nothing lands in a group the interface would
not draw — and the fallback stays too, because an action nobody classified must
still reach the player.

Balance unchanged: deaths 0/0/50/82/83, median cash 12378/14050/7230/90/586.

Evidence: `core/grouping_test.go`, one break verified.

## A map you can read your own city off

"We don't have the map working properly" is the vaguest line in the inbox, so
the first job was finding something checkable in it.

Two things, and both are checked now.

**Fourteen of twenty-five addresses were drawn as two silhouettes.** Seven are
type "work" — the docks, the haulage yard, the cab stand, the forecourt, the
scrapyard and two filling stations — and every one was the same low shed with
the same stack. Seven more are "racket" and were the same brick shopfront. A
city that is a grid of copies is a city you cannot find anything on, which is
what a map is for.

`blockFor` takes the address as well as its kind now, and the ones that do
different things look different: pumps under a canopy, a crane jib over the
scrapyard, plate glass on the forecourt, a line of cabs at the office, a flatbed
in the haulage yard, roller doors on both garages, an awning over the butcher,
long upstairs windows on the poolhall, vent stacks on the steam laundry. The
burlesque had no shape at all — its kind was not in the table, so a revue
theatre was drawn as a terraced house — and it has a lit marquee now.

The test does not demand that every building be unique. A row of shops should
look like a row of shops. It demands that no one silhouette stand for most of a
kind of work, which is what had happened.

**And a fault that has not bitten yet.** `grid` puts an address on the block
nearest its own coordinates and steps outwards when that block is taken — but
it tried only the four cardinal neighbours at each ring, gave up after five
rings, and then wrote itself into the taken map regardless. Two addresses could
end up in the same block and be drawn standing inside each other. The current
twenty-five do not collide, which is why nobody has seen it; the guard is there
for the twenty-sixth.

Evidence: `tests/city-grid.test.mjs`, three new tests, 50 node tests passing.

## The machine, made properly, and the noise of the room

"I also want you to flesh out the slot machine a lot more, make it much nicer
like you did for blackjack and roulette. Right now it looks scraggy." And, from
the same message: "We should also add ambient sounds and sounds to the slot
machines and whatnot."

The machine was three letters in three boxes, which is a picture of a result
rather than a machine. A real drum is a strip of faces turning behind a window,
and what you see is three of them at a time with the payline across the middle.

`reelStops` expands the core's own strip — a symbol with four stops appears four
times, and the length of the run is the number the odds divide by — and
`reelWindow` returns the three faces showing when a drum has stopped on a given
symbol. Its neighbours are its actual neighbours on the strip, which is why a
machine feels like it nearly paid. A symbol with several stops has several
homes, so three drums showing the same face need not show the same shoulders.
Five tests, including one that a drum can only ever show a face the core has.

Around it: a crown, the window with the line across it, the handle down the
side, and a tray at the bottom that says what fell into it.

**The noise.** Same contract as the rest of `src/sound.ts`: synthesised, nothing
downloaded, nothing licensed. A handle is a spring and a clunk, a drum stopping
is a wooden knock, a payout is a run of coins whose length is what it paid.
There is a floor tone under the room while the player is at the tables, started
when the takeover opens and stopped when it closes — a noise that goes on after
you have left the table is a noise nobody asked for. A card for a hand, a rattle
for the dice.

**And a thing the noise exposed.** The takeover's "this sitting" panel keyed off
the world's revision with a sentinel of -1, so the first look counted as
something that had just happened: it dealt the previous result into the list,
and once there were noises it played a card the moment the room opened. The
sentinel is null until the first look now, and the first look only records where
it is starting from.

Evidence: `tests/cards.test.mjs`, five new tests, 55 node tests passing.

## Ackerman & Son, pawnbrokers

The inbox is clear and the queue is done, so this is the standing instruction:
"there is always another business, and always more people to put in the city",
and "find ways to link them together."

A pawnbroker is the trade that sits at the join between things that already
exist and had nowhere to happen. What gets taken off the street has to turn into
money somewhere. Somebody who is short has to turn what they own into money and
hope to get it back. Neither had an address.

**Both directions are real.** Robbing a till, taking somebody's pockets and
stripping a car all lift the shop's trade — the same shape the garage already
has with car theft, which is the link the inbox liked. And the counter lends
against the player's own car or the suit off their back: 35% of what it cost
new, scaled by its condition, and 130% of that to get it back inside a week.
After that it is in the window and it is gone. Pawning the suit costs the
standing it was buying, which is the decision.

A ticket that runs out is kept rather than thrown away, marked sold, so the
counter can tell the player what became of the thing instead of shrugging at
them.

**The checklist caught everything I forgot**, which is what it is for. Three
guards failed on the first run: no painted front, no interior to stand in, and
"The Golden Lily and Ackerman & Son both want block 3,1 — one of them will be
pushed off its own coordinates". Moved to a free block, and painted. The front
came back with PAWNBROKE'S STOP across the fascia in letters the model invented,
so the prompt asks for a plain unlettered board and it was repainted.

Twenty-six addresses, fifteen ownable businesses, twelve kinds. Three more
people in the city: a pawnbroker, a counter clerk and a valuer.

Balance: deaths 0/0/50/82/79 against 0/0/50/82/83, median cash
12378/14050/7230/90/567 against .../586. The thief moved by four deaths in a
hundred and by three percent of a small number, which is a twenty-sixth address
changing which room a policy walks into rather than anything paying differently.

Evidence: `core/pawn_test.go`, and the three address guards that failed first.

## Watching the city instead of the player

`docs/LIVING_WORLD.md` sets the standard for these systems: "a system is not
finished because it runs. It is finished when a long simulation shows it
producing varied, non-degenerate outcomes." Every report this harness has ever
printed says `factions_destroyed: 0`.

That is not because nothing falls. Left alone for sixty days, ten organizations
fell across twenty-five cities. The measurement was watching the wrong thing: a
campaign follows one protagonist, and the median campaign lasts between three
and fourteen days — the reckless policy lasts five hours. Families escalate,
split and collapse over weeks, and no run was alive long enough to see any of
it. The number was reporting the lifespan of the player, not the health of the
city.

So `sim.City` runs a city with nobody in it. No policy and no commands: the
clock turns and the families do whatever the rules make them do. Any scene is
cleared rather than answered, because this is a measurement of the city and not
of what the director would have said about it.

| 12 cities, 60 days each | |
|---|---|
| Organizations formed | 11 |
| Organizations fell | 7 |
| Wars started | 28 |
| Wars settled | 25 |
| Holdings changed hands | 22 |
| People killed | 32 |

Two tests hold it to the standard. One says a city must keep moving: something
must fall, something must be formed, and property must change hands. The other
says it must not degenerate: no city may end with nobody alive, with no
organizations left, or with ninety percent of the held property in one pair of
hands.

**That second number was wrong, and the correction is the more interesting
result.** I reported the largest holding in any of the twelve as eighty percent
and said it was heading towards collapse. It was not a holding. The measure
counted every distinct value of `Owner`, and "independent" is a value of
`Owner` — the placeholder for a shop that answers to nobody. Seventeen of
twenty-one earning addresses carry it, so the instrument reported four fifths of
the city in one pair of hands and the hands were nobody's. A figure that says
nothing while looking like one, in the very slice that exists to stop that.

Counting only owners that resolve to a live organization:

| Days | Largest organization | All organizations together |
|---|---|---|
| 60 | 19% | 19% |
| 120 | 23% | 20% |
| 240 | 28% | 23% |
| 480 | 28% | 25% |

Nothing is concentrating. The opposite is true and it is the real finding: after
a year and a third of simulated time, organizations hold a quarter of the city
and three quarters of it answers to nobody at all, permanently. Layer 1 of
`docs/LIVING_WORLD.md` says families own income-earning property; layer 4 says
war creates the openings a player exploits. A map three quarters of which nobody
is fighting over is as dead as a map somebody has won, and that is the next
thing to fix rather than to measure.

`mise run simulate` reports all of it under `city_alone`, so a living-world
claim made in this file can be checked against a run rather than asserted.

Evidence: `sim/city_test.go`, one break verified.

## A city somebody else is competing for

The corrected measure said the families held nineteen percent of the city after
two months and twenty-five after sixteen. Three quarters of it answered to
nobody, permanently.

Nothing in the rules ever took an unheld shop. Property moved between families
in a war, a splinter walked off with one holding, and a family that had lost
everything could start again on unheld ground. A family in good order never
grew. Layer 1 of `docs/LIVING_WORLD.md` says organizations own income-earning
property and layer 4 says war creates the openings a player exploits; neither is
true of a map nobody wants.

`ConsiderExpansion` is a family in good order leaning on somewhere that answers
to nobody: forty-five power, three hundred and twenty of their own money, the
nearest unheld address to their own seat so a family grows outward from where it
already is. One a turn, off the world's own stream, because the player is not
party to it. And a family already holding a quarter of the city stops — past
that they are not expanding, they are winning.

**That cap was not enough on its own**, and the measurement said so: once there
was more to take, war seizures carried one family to sixty-one percent by eight
months and it stayed there. So the other half of this slice is that a family
which has won too much fractures from the inside. Splintering was only ever
weakness, fighting, or having nobody left to fight — which is half of why
anybody breaks away. The other half is a lieutenant looking at how much there is
and how little of it is his.

| 12 cities | Organizations hold | Largest one |
|---|---|---|
| 60 days | 52% | 28% |
| 240 days | 65% | 57% |

Holdings changing hands over a season went from 22 to 104. Organizations formed
went from 11 to 15.

Two breaks verified, one per half: with expansion disabled the city sits at
nineteen percent again, and with the overgrown rule removed one family reaches
sixty-one percent and the guard fails.

**And the prose guard earned its keep again.** The headline read VERA KOHL'S
PEOPLE TAKES OVER THE PAPER MOON. A family name can be plural, `Agree` exists
for exactly that, and the scan that reads every passage in the game found it in
four different city states.

Balance: deaths 0/0/51/82/77, median cash 12405/14156/7230/90/547.

Evidence: `core/expansion.go`, `sim/city_test.go`.

## The last thing moving in somebody else's units

`Lot = 5`, and the comment above it said why: "the interface offers plain
actions rather than a quantity field, so trade happens in fixed lots." That was
true when it was written. It stopped being true the night the typed-amount work
landed, and nobody went back for it — so the one high-variance income path in
the game, the one whose whole decision is *how much do you dare carry*, had that
decision made by a constant.

Buying and selling take a number now. Nothing named still means the lot for a
purchase and all of it for a sale, which is what selling has always meant here:
the risk ends when the last of it is gone, so every existing caller, save and
simulation behaves as it did.

The field knows both real limits — what a person can carry, which is their
pockets plus whatever the car hides, and what they can pay for. A card should
not offer a number the rule will refuse, and the refusals say which limit was
hit rather than "not enough cash": "You can carry 3 more crates", "That is
$420", "You are carrying 2 crates".

Two breaks verified, one per side: a purchase that ignores the typed number
buys four with the price of three, and a sale that always empties the stock
sells twelve when five were asked for.

Balance unchanged: deaths 0/0/51/82/77, median cash 12405/14156/7230/90/547.
The simulated policies do not trade contraband, so that is an absence of
evidence rather than evidence — the same thing the car work had to say.

**And a check for harm I might have done last tick.** Families now take unheld
premises, which could quietly narrow the player's path to owning anything. It
did not: the investor and defiant policies still reach the laundry, the garage
and the casino in 100 runs of 100. The thief moved from 28 to 25 and 25 to 23,
which is noise for a policy that dies at fifty commands.

Evidence: `core/trade_amount_test.go`.

## Somebody who actually trades, and what they found

Last slice let contraband move in whatever figure the player types, and had to
admit the balance numbers said nothing about it: no policy in the harness had
ever bought a crate of anything. The main high-variance income path in the game,
by the design document's own account, was entirely unmeasured.

Two things had to exist first. The view carried premises and people and not the
one number the trade is decided on, so `Goods` is on it now. And `View.most`
lets a policy ask a card how much it will take — the core already works out what
a person can carry and what they can pay for, and a policy that recalculates
that is a second opinion waiting to disagree. Mine did: the first version asked
for four crates with room for one and the campaign ended on the refusal.

**What the smuggler found is not what the design claims.**

| 60 campaigns | |
|---|---|
| Purchases | 152 |
| Sales | 26 |
| Median days | 3.8 |
| Median final cash | 3,570 |

Six purchases for every sale, and a third of the investor's money. Contraband is
not the main high-variance income path; it is a slow way to end up holding
stock. The reason is structural rather than a matter of tuning. There is one
price series for the whole city, so buying at the docks and selling at the
market is the same price — there is no route, only a wait. Prices do swing
widely, 55% to 177% of base over four months, but a campaign is a few days long,
so a cycle rarely completes inside one. Meanwhile every day the stock is held
draws attention.

Layer 6 of `docs/LIVING_WORLD.md` asks for "sources, routes, and buyers" and
"profit from risk taken knowingly". What exists is a price that moves and a
pocket that holds. Routes are the missing half, and that is a design slice of
its own rather than something to tune.

Balance, six strategies: deaths 0/0/51/82/77/0, median cash
12405/14156/7230/90/547/3595.

Evidence: `sim/campaign.go`, and the 152-against-26 count above.

## A route, which is two prices and a walk between them

Last slice measured the underground trade and found it was not one: 152
purchases against 26 sales over sixty campaigns, because there was one price for
the whole city and buying at the docks and selling at the market was the same
transaction done twice. There was no route, only a wait.

A floor is dearer or cheaper than the city's price by a fact about the floor.
The waterfront is where it comes ashore — moonshine at 78% and arms at 82% —
and the exchange is where the buyers are, at 118% for moonshine. The city's
price still moves under all of it, so a good week and a bad week are still real;
what is new is that the two ends of the city disagree about what a crate is
worth, permanently.

And the difference is paid for by the walk. That is the point of doing it now
rather than earlier: the street between two addresses is the one stretch of this
city where nothing covers anybody, which is why a car can be plated and why the
crossing bar says what you are carrying it in. A route is a reason to be out
there with something worth taking.

The card quotes both ends: "$31 each here. Mercer Exchange pays $47."

| 60 campaigns | Before | After |
|---|---|---|
| Purchases | 152 | 3,470 |
| Sales | 26 | 3,463 |
| Median final cash | 3,570 | 5,577 |

It is a trade now rather than a way to end up holding stock. It is still worth
less than owning premises — 5,500 against the investor's 14,156 — which seems
right for money you have to carry through the street to collect, and no
smuggler died in sixty campaigns, which says the risk is still mostly attention
rather than violence. Whether that is the right shape for "the main
high-variance income path" is the next question about it.

**Two test changes, with the reason.** Two tests asserted that a trade moves
`Good.Price`. The price is a fact about the floor now, so they ask the floor
they are standing on. One break verified: with the spread deleted, carrying four
crates across the city turns $160 into $160.

Balance, six strategies: deaths 0/0/51/82/77/0, median cash
12405/14156/7230/90/547/5500.

Evidence: `core/routes_test.go`.

## What the city takes, as against what the player keeps

The route made the underground trade a trade. The open question it left was
whether the risk is the right shape: no smuggler died in sixty campaigns, and a
report of deaths and cash cannot tell a safe policy from one whose money is
taken rather than whose life is.

So the report carries three more numbers: the attention a run ends on, whether
it ended hurt, and how many times the goods were taken. A seizure is written in
the record rather than in a counter — it is a thing that happened — so the
harness reads it back out of the history.

| 100 campaigns | Deaths | Median cash | Mean heat | Seizures | Hurt |
|---|---|---|---|---|---|
| worker | 0 | 12,405 | 0 | 0 | 100 |
| investor | 0 | 14,156 | 0 | 0 | 0 |
| defiant | 51 | 7,230 | 0 | 0 | 63 |
| reckless | 82 | 90 | 0 | 0 | 82 |
| thief | 77 | 547 | 57 | 0 | 95 |
| smuggler | 0 | 5,500 | 37 | 14 | 0 |

The three paths are three different dangers, which is what they should be. The
investor is untouched. The thief is the violent one — fifty-seven attention,
ninety-five runs in a hundred ending hurt, seventy-seven of them dead. The
smuggler carries real attention and loses the goods fourteen times in a hundred,
and is never once hurt.

**One of three risks exists.** Layer 6 of `docs/LIVING_WORLD.md` asks for
"seizure, informants, a rival who wants the route". Seizure is there and it
bites. Informants and a rival who wants the route are not built, and now that
the route exists a rival wanting it is a thing the faction system could act on —
which is also what the layer means by contraband being "a common cause of war".
That is the next slice on this system and it is written down rather than
guessed at.

The worker being hurt in every campaign and holding no attention is worth a
second look on its own: that is a policy that spends its life on the docks.

Balance, six strategies: deaths 0/0/51/82/77/0, median cash
12405/14156/7230/90/547/5500.

Evidence: `sim/campaign.go`, `cmd/simulate/main.go`.

## A rival who wants the route

Of the three risks layer 6 asks for — "seizure, informants, a rival who wants
the route" — one existed. A policy running the trade over a hundred campaigns
carried thirty-seven attention, lost the goods fourteen times and was never once
hurt, because nobody in the city had an opinion about a man carrying crates
through it.

Selling is what gets somebody noticed. Buying is a man with money; selling is a
man with a trade, so `Player.Runs` builds on loads that cross the city and fades
two a day, because a trade you have stopped running stops being your trade. Past
thirty units the word is out, and the strongest organization that is not the
player's decides a trade worth running is a trade worth taking.

They come for the load rather than for the man. If he is carrying, they take it
in the street and leave him standing, and they are richer and stronger for it.
If he is not, they say what they came to say and the word cools by half. Either
way it costs the player standing and it costs that family goodwill, so answering
it is a quarrel that already exists — which is what the layer means by
contraband being "a common cause of war".

| smuggler, 100 campaigns | Before | After |
|---|---|---|
| Seizures | 14 | 65 |
| Median final cash | 5,500 | 5,123 |
| Deaths | 0 | 0 |

Still worth running and no longer unopposed. Nobody else's numbers moved.

**A break that landed in the wrong place, and the fix.** I first broke the
seizure by zeroing the count it returns rather than the call itself, and the
test passed: the goods were still gone, and the test was asking about the goods.
Breaking the call instead fails it with "they came for the load and left 6 of
it". The gender guard also caught "taken from a man on foot", which is the third
time tonight it has earned its place.

Informants are the risk still missing.

Balance: deaths 0/0/51/82/77/0, median cash 12405/14156/7230/90/547/5123.

Evidence: `core/route_rival_test.go`, two breaks verified.

## Somebody talks

The last of the three risks layer 6 asks the underground trade to carry, and the
only one that is not about wanting what you have. An informant is somebody who
wants you finished and has found a cheaper way to do it than a gun.

It had to be a person. The rule this game already applies to the people who rob
the player — "named people with a place in the city, not anonymous thieves, and
they can be answered" — holds here, or the whole thing is a dice roll wearing a
hat. So an informant is somebody already carrying a grudge the simulation
committed, over a reason it recorded, and saying it costs them most of the
grudge: telling the police is what they had to get off their chest.

Whether the name comes back is the question the city already asks about a
robbery. With contacts you hear who and why. Without them you hear that Ward
Street knows more than anybody there worked out for themselves, and that is all
you ever hear.

**The first version was too narrow, and the measurement said so.** I tied it to
the trade alone, and in a hundred campaigns nobody ever picked up a telephone —
because the policy that runs the route makes no enemies, and the policy that
makes enemies never trades. What an informant actually needs is somebody who
hates them and something to point at: a trade being run, crates in their hands,
or a file the police already have open.

**And a limit worth stating.** The mechanism is unit-tested with two breaks
verified — no heat added when somebody talks, and the name never coming back —
but the campaign figures barely move, because neither policy in the harness
satisfies both halves at once. A policy that both trades and makes enemies is
the probe this needs, and it does not exist yet. What is written above is what
the tests show, not what the campaign numbers show.

Balance: deaths 0/0/51/82/76/0, median cash 12405/14156/7230/90/567/5160.

Evidence: `core/informant_test.go`.
