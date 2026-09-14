
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
