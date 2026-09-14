
## 2026-09-13 — browser 3D city integration

The current user request replaces the previous 2D-only handoff scope. Added the
single Three.js city, a reproducible 18-model Blender pipeline, public journey
vehicle labels, committed event playback and isolated browser fixtures. Browser
QA caught and corrected lost simultaneous explosion cues, remote blast casualties,
material-color loss and compact layout clipping. See `docs/CITY3D_PROGRESS.md`
for validation, screenshots, reproducible previews and the substantial remaining
art, choreography and traffic work. The production-quality goal remains active.

### 2026-09-13 — personal HUD, accounts, and travel presentation

Removed the redundant City navigation item and City briefing overlay. The identity calling card now shows the actual player's portrait, name, life and standing instead of the old title block. Existing threat, commission and opportunity information remains reachable in Guide. Camera controls and latest-result entry use paper/ink styling; the result uses a bounded, wrapping grid instead of a single overflowing flex row. Rail hover no longer transitions between unrelated background images. Daily accounts are a persistent expandable city slip with the complete cost breakdown and arrears; Ledger retains searchable history.

Journey progress now reports actual traffic placement progress to the clock display, throttled to game-minute changes. This interpolates the already committed interval and never advances or rewrites simulation time. Starting travel enables smoothly eased tracking/zoom; returning from an interior mounts the city focused on the player outside. Manual pan cancels tracking. Zoom easing stops after framing, allowing subsequent zoom adjustments.

Validation: frontend suite 211/211, TypeScript and production build passed (existing bundle-size warning remains). CUA browser QA on isolated fixtures at 8878 and 8879, 1280×720: walking clock 10:00→10:01 during motion and exactly 10:15 on arrival; driving 11:28→11:33→11:39→11:44 during motion and exactly 12:44 at arrival, while the saved final minute remained 764. Following camera tracked the moving actor; exit from Ferris refocused the exterior actor at x16/z4.65, zoom8. Ledger contained no financial grid. Expanded result had no horizontal overflow; account breakdown rendered legibly. Local diagnostic FPS 145 during sampled driving, not a cross-device performance guarantee. Evidence: `.runtime/hud-travel-check.json`, `.runtime/hud-revision.png`, `.runtime/hud-revision-tests.log`, `.runtime/hud-revision-build.log`. Compact viewport and physical pointer-hover capture remain unverified; background-transition removal was checked in computed styles. Main campaign untouched. Driver/cabin refinements remain a separate uncommitted workstream pending close-up visual review.

### 2026-09-13 — full-scale seated drivers

Integrated the cabin proportions and seated-character work begun before the HUD revision. Ford, Hudson and Packard journey actors now include their public character's rig and wardrobe behind locally cloned translucent glazing. Drivers are shown during journeys at zoom6+, hidden when distant or parked, and their private materials are released with the actor. No NPC identity or movement is invented. The Blender source now includes steering-wheel grip anchors and revised bench/floor/footwell dimensions; all four articulated car assets regenerated locally. Rigid-limb inverse kinematics places both hands at the actual grips, preserving character scale and limb lengths. Hats are hidden while seated.

Browser review caught the original low hand position and prompted the grip correction. Ford/person and Packard/woman inspected from the driver side with open doors. Evidence `.runtime/packard-driver-seated.png`; the review page now has camera and rig controls and waits for both car and character loading before replacing a model. Tests check both rigs in each civilian car for floor, roof and lateral clearance, pelvis placement and both steering grips within 2.5cm. The former constant front-clearance bound was corrected to each car's authored wheel position. Frontend tests211/211 and production build pass (`.runtime/seated-driver-tests.log`, `.runtime/seated-driver-build.log`); build retains the existing size warning. Previous isolated live city journey samples verified actual driver visibility at close zoom. This is seating and visibility only: boarding, disembarking, animated wheel turning, passenger custody/departure and production-level character detail are not complete. Main campaign untouched.

### 2026-09-13 — keyboard operation and explicit free camera

Manual wheel zoom, keyboard zoom, and pointer camera gestures now turn following off. WASD remains continuous. Q/E now shares the held-key state and advances camera azimuth per frame instead of OS-repeat jumps; opposite rotation keys cancel, key release/focus loss clears input, and long frame gaps are capped. A focused regression checks equivalent one-second rotation at30/60/144FPS, opposite-key cancellation and release.

Added period-styled keyboard reference (`?`) and searchable available actions (`/`, arrows, Enter). Search uses the enabled controls in the active dialog or city/interior, excludes inert and closed disclosure contents, and rechecks connectivity/availability before activation. It uses the existing control handlers and costs, not a second command API. Native form input, IME, modifier shortcuts and repeated action keydowns are protected. Direct keys:1–6 navigation,7 settings,8 accounts,V voices,G travel/enter,B leave,F follow,O overview,Z focus,J directory,T playback speed. Camera keys work without first clicking the canvas. The daily net summary now uses the same green/nonnegative and red/negative colors as the expanded accounts.

Browser QA on isolated8879/8876: wheel zoom changed following=true to false and zoom~8 to5.526; keyboard minus and D released follow; Q changed azimuth through the frame loop; `/`→Ledger→Enter opened Ledger; G entered Mercer; action search left the interior; numeric4 opened Ledger; typing `wasd123?/` stayed in its search box without opening other UI;8 expanded accounts; disabled preview controls were excluded from search. Positive$303 computed green rgb(78,97,57), negative−$15 red rgb(150,62,45). Keyboard guide reviewed at1280×720. No browser error logs in sampled session. No gameplay commands or main-save changes were needed for these checks. Build passes with existing bundle-size warning. Full suite and targeted rotation results recorded in `.runtime/shortcuts-tests.log`; overall production acceptance remains open.

### September 13 scope updates

Read every incoming amended objective through the current source
`/Users/jatochnietdan/.codex/attachments/7c9e06f6-fd97-41d6-a946-eee7fde6c804/goal-objective.md`.
The additions include seated/working interior occupants, cause-specific player
death scenes, people/family backstories and relationship history, and illustrated
contraband/arms. These are outstanding additions to the full goal, alongside the
previous gambling, molotov, drive-by and street-incident requirements.

### 2026-09-13 — shared street playback and account clarity

Travel receipts now retain the public NPC legs observed during the committed interval. NPC departures and arrivals are simulation boundaries; the browser samples recorded legs against the player's shared presentation clock and reconciles to the final public snapshot. Browser QA exposed a deadlock when that clock depended on a blocked player's position: it now continues while individual actors wait for space. Entry remains disabled until the player's visual journey completes. This supersedes the earlier actual-placement clock during recorded journeys; visual collision waits never advance core time.

Imported the supplied newspaper-opening WAV (1.68s, stereo48kHz PCM16), with cancellable playback and procedural fallback. The identity card names the current interior. Daily accounts opens its itemized cost disclosure by default, retaining the option to collapse it. CUA verified these on isolated8881: opening accounts immediately showed Rent/The Mariner/$15 and Staff/$18, with both disclosures open.

Validation: frontend215/215 and production build passed (existing bundle warning); core and store suites passed, and cmd/blackledger passed after updating two stale static UI assertions for the relocated accounts and multiline City3D markup. Targeted tests cover observed departure/arrival, continuous partial legs, empty receipts, and sampling sequential legs. Finer arrival settlement invalidated old net-profit and bounded-history assumptions; their regressions now verify the actual restock debit and observe notifications as they occur. Corrected browser playback completed at minute630, with NPCs continuing while the player yielded; entry remained disabled until arrival. Evidence: .runtime/street-travel-frontend-tests.log, .runtime/street-travel-build.log, .runtime/street-travel-integration-tests.log (initial stale UI assertion failures retained), .runtime/street-travel-http-tests.log, .runtime/street-travel-regression-tests.log. Main campaign untouched. Collision-free choreography across all routes and broader production acceptance remain unfinished.

### 2026-09-13 — assassination preview coverage and close-quarters cast

The debug selector now exposes all current recorded strike/weapon combinations: shot from behind, revolver, shotgun, Thompson and unarmed close quarters. Private preview receipts no longer inherit the preceding action's costs or elapsed time. Existing real outcome cues and API commands are unchanged.

Added an unarmed shared cast for explicit close-quarters strikes, reusing the fully reserved approach space and exact-victim cue merging. The attacker walks into range, delivers three cosmetic blows and the victim falls after the last impact. Synthesized low thuds, restrained camera impulses and the supplied pain recording follow those beats; owned voices cancel on disposal. No weapon or muzzle flash is invented.

Validation:217/217 frontend tests and TypeScript/production build pass; existing bundle-size warning remains. New geometry tests sweep both attacker rigs at60Hz for scene bounds and pavement clearance, verify hand contact at each impact, and prohibit early falls. CUA on isolated8881 verified all menu choices, Thompson weapon3/model, unarmed cast at5.14s with pain recording active, and Stop clearing effects/audio. Browser error log empty; sampled145FPS/p95 7.7ms, not a cross-device guarantee. Save remained revision2/minute630. Logs:.runtime/assassination-preview-tests.log and .runtime/assassination-preview-build.log.

The close-quarters choreography is an initial pass, not final animation quality: victim anticipation/reactions, more varied blows and character detail need improvement. Armed close-shot/burst still use the earlier stationary firing sequences; moving car assassinations, bomb-planter preambles, molotovs and other requested scenes remain outstanding. The full production goal stays active.

### 2026-09-13 — staged armed approaches and withdrawal

Revolver/shotgun close-shot and Thompson burst now stage the recorded attacker and victim together, replacing the earlier independent stationary firing/falling figures. The attacker approaches, raises the actual weapon, fires only after aiming, then turns and walks back out of the reserved space. The victim raises their hands and falls after impact. Long guns carry upright during approach and turns, keeping barrels clear of the adjoining travel lane. Hands retain fixed limb lengths and both authored long-gun grips. Gun samples, muzzle flash, recoil, impact shake, torso-height spatter and pain recording follow the same firing beats. A close-shot fires once; a Thompson strike fires three times. These are cosmetic counts of a resolved strike, not new gameplay damage.

Validation:219/219 frontend tests and production build pass (.runtime/armed-strike-tests.log, .runtime/armed-strike-build.log); bundle warning remains. New tests sweep person/woman attackers with all three guns at60Hz throughout8s, checking reserved bounds, pavement clearance, both grip contacts, muzzle rays into the victim, no pre-shot fall and completed withdrawal. Audio checks prohibit early shots, assert1/1/3 onsets and cancel every voice. Existing execution/unarmed tests remain green. CUA isolated8881: shotgun before-shot screenshot at3.24s, departure/corpse at5.94s with one shot, Thompson at4.14s with three recorded shots, and final revolver build with victim hands raised before gunfire. Stop cleared scene and audio; no browser errors. Save unchanged at revision2/minute630. Sampled145FPS/p95 7.5ms, not broad performance acceptance.

This supersedes the preceding note that close-shot/burst were stationary. Remaining production gaps include richer movement/reaction animation, character fidelity, proper escape routing into the public street and body pose continuity beyond mounted playback. Moving-car attacks, bomb preambles and the rest of the expanded goal remain outstanding. Main campaign was not touched.

### 2026-09-13 — exact strike-to-aftermath pose handoff

Found that the aftermath body preserved its root but reset articulated joints to the source model. It now retains copied final joint quaternions and applies them to the independent aftermath clone. Sampling the finished pose at staging also handles Skip; the actor is reset to its opening pose before playback. Snapshots retain numbers only and are copied again on receipt, preventing later rig updates or caller mutation from changing the body. Cleanup still releases the retained pose. Added a repeatable `strike-unarmed` isolated fixture.

Validation:220/220 frontend tests, TypeScript and production build pass (.runtime/body-handoff-tests.log, .runtime/body-handoff-build.log); existing bundle warning remains. New test compares every mesh world matrix across the handoff for both victim rigs and all five strike variants, including caller-mutation isolation. CUA replay on newly created .runtime/body-handoff-unarmed-20260913.sqlite3 at8882 showed the final victim at x81/z38.35, then the same visible body/arm pose after the Herald was closed; no browser errors. No main campaign mutation. Full-reload pose/placement retention and more natural settled anatomy remain unfinished.

Read the newest objective19b29753-a6c4-4d03-b862-e273c86b1f90. Interior assassinations and plausible wire attacks from behind are added to the outstanding scope; all previous requirements remain active.

### 2026-09-13 — Saint Agnes occupant staging

Replaced the standing grid with deterministic booth/stool placements derived from the existing Blender room's furniture coordinates. Mara receives the first corner booth when present. Only an explicitly public bartender/barman/barmaid role receives the service position behind the bar. Other public guests use separate seated positions and then aisle positions; no decorative people are invented. Articulated legs and resting hands retain their original scale. The allocator supports ten guest positions plus a service worker; excess people remain available in the HTML roster. Q/E in this interior now uses the same held-key frame rotation as the city.

Geometry QA caught opposing booth guests' feet overlapping; seating them diagonally along the benches separated the physical bounds. Tests now cover stable assignment under roster reorder, exclusive positions, role-based service placement, both model rigs' floor clearance, pelvis height and pairwise actor bounds.222/222 frontend tests and production build pass (.runtime/interior-staging-tests.log, .runtime/interior-staging-build.log); bundle warning remains. Browser CUA on isolated8881 showed three seated guests and public bartender Aldo Olsen behind the counter, correct interior identity, Q input rendering, and no browser errors. Sampled232 draws/129232 triangles; the room renders on demand when unchanged. Main campaign untouched.

This is initial static furniture staging for the currently modeled Saint Agnes interior. Animated service, drinking/conversation, arrivals/departures, authored semantic seat anchors, detailed furniture-contact verification and other building interiors remain unfinished. The broader interior/combat/gambling and production-art goals stay active.

### 2026-09-13 — scene control styling and compact interior receipts

Scene preview controls, playback/Skip panels, replay and exterior-return buttons now share the city HUD's paper/ink styling and explicit keyboard focus treatment. Interior latest-entry feedback defaults to a compact disclosure instead of obscuring the premises actions with an expanded fixed receipt; the full result and Ledger control remain inside. No entry is drawn before a result exists.

Production build/TypeScript passed (.runtime/scene-slips-build.log). CUA isolated8881 verified the paper controls, a39px collapsed interior entry, expansion showing the complete result and Ledger link, re-collapse, and no browser errors. Existing gameplay handlers unchanged; main campaign untouched. Other interior/menu styling and compact viewport acceptance remain outstanding.

Read amendment b0f58a61-526b-4ecb-9424-6b007c616bc4: The Mariner must support takeover, actual living NPC tenants, occupancy-based rental income and management obligations. This is the next gameplay priority; all prior visual/simulation requirements remain active.

### 2026-09-13 — income consistency before residential integration

Read the Mariner amendment and its successor903a4109-a113-4a4f-b726-f6ef277f7e41, which additionally requires wealth/standing-appropriate homes for every living NPC. Traced acquisition, workplace/travel state, business obligations, purse settlement and housing charges. Findings and remaining integration are recorded in docs/MARINER_IMPLEMENTATION.md; residential leases and Mariner takeover are not yet implemented.

Fixed a prerequisite accounting discrepancy: Books omitted RoomTrade and CollectionShare while Advance included them. Both now use HourlyIncome, retaining the existing operating factors and fractional payout carry. Targeted exact-hour and account-rate tests pass; full core/store/cmd/blackledger suites pass (109.8s/.128s/.694s), logs .runtime/income-consistency-tests.log and .runtime/income-consistency-integration.log. No save migration or main-campaign mutation. Existing unrelated mugging/robbery edits remain untouched.

### 2026-09-13 — persisted NPC residences

Added NPC Home/Accommodation state, separate from workplace, whereabouts and journeys. New population and v15 migration assign missing homes deterministically using purse/standing, preserve existing residents, reserve the player's place, exclude dead occupants and report housing shortages. Current capacities are24/64/1 at Mariner/Ashbury/Cypress. People cards show the home and accommodation independently of current whereabouts; shortages are visible. These are initial residence assignments, not implemented rent billing or a completed housing economy.

Validation: targeted assignment, wealth, capacity, reorder, death, migration and save round-trip tests pass. Core/store/HTTP integration passed (99.8s/.136s/.677s); existing player-residence ownership and anti-farming regressions passed separately after retaining them in their original test file. Logs .runtime/housing-tests.log, .runtime/housing-integration.log, .runtime/housing-http-tests.log, .runtime/housing-residence-regression.log. TypeScript/build passed (.runtime/housing-build.log, existing bundle warning). CUA isolated8883 (.runtime/housing-20260913.sqlite3) showed the starting86 NPCs with homes and zero shortage, including Mariner rooms, Ashbury apartments and Cypress residence; no command or main-campaign mutation. The prior income correction is also served in that isolated preview.

Remaining: rent payments, homeward travel, changes with long-term wealth, physical capacity/art validation, additional housing, house sales, Mariner management and tenant UI. Latest objective84acfc45 retains these and adds general progression-income expansion and a beginner questline. Full production goal remains active.

### 2026-09-13 — daily cost breakdown without a second disclosure

Daily accounts now contains a permanent cost-breakdown section instead of a nested expandable control. Expanding the accounts always reveals each expense and its explanation, including the Mariner's $15 rent. Preserved the paper heading style. Production build/TypeScript passed (.runtime/accounts-breakdown-build.log); CUA isolated8883 confirmed one expansion shows rent, staff, car and total, with no nested disclosure or clipping at the checked desktop size. No campaign commands issued.
