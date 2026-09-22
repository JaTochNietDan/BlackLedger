
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

### 2026-09-13 — ordinary residential journeys

Ordinary residents now schedule overnight journeys to their persisted homes and commute from06:00, preserving workplace Post on home arrival. The clock includes the morning boundary; all journeys use existing room occupancy/public street/command travel recording. Existing noon outings, duties, urgent errands and custody retain priority. Departure staggering remains spread over five hours after a shorter window failed the street-crowding regression. Managers/officials still need a fuller staffed shift system.

Integration uncovered office successors retaining pre-appointment journeys: appointment now cancels those and sets the new Post. Replacement officials also retain duty by role, independent of their original NPC ID. Counter test setup now explicitly clears journeys after its eight-day business warmup; curfew checks compare actual evening crowds rather than assuming staying home must make people less predictable.

Evidence: targeted home/commute/recorded-travel/custody/successor tests passed (.runtime/home-routine-tests-final.log); repaired regression group passed (.runtime/home-routine-regressions.log). Full core/store/HTTP integration passed (.runtime/home-routine-integration-final.log); the final replacement-role guard was checked separately afterward. Earlier failures retained in .runtime/home-routine-integration.log. CUA isolated8884 home-routine-preview-20260913.sqlite3 showed the homebound resident as a person actor and the address book correctly listed a Mariner rented room and Walking to The Mariner; no browser errors. Browser fixture used the earlier short departure offset; final stagger is covered by simulation tests. Main campaign untouched. Rent billing, Mariner management, expanded housing/market and broader goal remain incomplete.

### 2026-09-13 — tenant rent receipts and residential register

Schema16 adds once-per-day NPC budgets and per-property tenant rent accounts. Actual rent replaces bundled lodging; living residents pay while away or travelling, missed payments become arrears, and only cash collected reaches the current property holder. Mariner/Ashbury do not also earn hourly cash. Books/public rate report contracted living-tenant rent. Owner occupants pay no room rent to themselves. An interior paper register lists tenants/rates and, for owners, payments and balances. This does not yet enable Mariner acquisition or operating obligations.

Validation: rent transfer, no double charge/accrual, short paydays/arrears, absent/dead tenants, ownership transfer, save reload, owner occupancy and public-account privacy passed. Full core/store/HTTP suites passed (98.690s/.150s/.689s) in .runtime/rent-integration.log before the register projection and owner waiver; subsequent owner/books/settling tests passed24.116s, projection tests.194s, HTTP.635s. Logs .runtime/rent-owner-tests.log, .runtime/rent-projection-tests.log, .runtime/rent-http-tests-final.log. TypeScript/Vite build passed (.runtime/rent-build.log). Initial fixture assumptions in purse_test repeated weeks at one minute; now advance the day explicitly. HTTP null-shape guard now documents that non-residential places have no register; initial failure retained in .runtime/rent-http-tests.log.

CUA isolated8885 rent-preview-20260913.sqlite3 showed two15-dollar tenants:15 paid,6 paid with9 owing, and precisely21 added to owner cash. Register expanded and all rows remained reachable in the scrolling interior; no browser errors. That preview predates the final owner-rent waiver, verified separately by core tests. Main campaign untouched. Unrelated armed-resistance changes in mugging/robbery remain unstaged; only the PayTheCity rent changes are included.

Remaining: acquisition/pricing, staff/supply/condition obligations, delinquency enforcement, collection from former tenants, richer accommodation/real estate and the broader production scope. Historical receivables stay with premises; no evictions are silently performed.

### 2026-09-13 — buy and operate the Mariner

Schema17 introduces lodging as a trade while keeping the Mariner's home/move-in semantics. The3600-dollar base freehold uses normal acquisition and holding premiums. Buying retains23 starting tenant accounts and brings3 staff positions,18/day wages,40 coal/linen stock,4/day drain,100 restocking and140 boiler remedy. Existing repair, wage, manager and mode controls apply. Service/condition credits reduce room rent; generic standing orders are excluded. Nominal lodging Income is a capacity reference, with actual player/family rent cash received only from tenant payments. Migration does not gift stock to already-established v16 businesses or refill lodging on repetition.

Targeted purchase/service/migration tests passed (.runtime/mariner-tests.log); updated general-trade regressions passed (.runtime/mariner-regressions.log). Legacy nonbusiness-room tests now use Ashbury, which remains unmanaged, and no-price checks inspect the acquisition price rather than a residential move-in price. Added lodging to the trade-benefit guard with its owner-occupancy benefit. Full core/store/HTTP passed99.832s/.187s/.733s (.runtime/mariner-integration-final.log); initial failures retained in mariner-integration.log. TypeScript/Vite build passed (.runtime/mariner-build.log).

CUA isolated8886 mariner-purchase-20260913.sqlite3: normal Buy action charged3600 (12000→8400), retained23 tenants, displayed345/day rent and18/day wages. Midnight settlement left8727, exactly345 receipts minus18 wages; supplies fell40→36. Normal restock charged100 (8627), advanced45minutes, restored stock and paid no duplicate rent. Owner UI showed tenancy register and management controls; no browser errors. The main campaign and unrelated armed-resistance files were untouched.

Remaining: full lease/eviction/debt enforcement, physical accommodation/3D interior fidelity, real-estate sale market, more housing, income progression/questline and all other unfinished production requirements. Current Mariner presentation still uses the old room artwork and general interior action surface.

### 2026-09-13 — property exchange, sale and deed-only purchase

Schema18 adds the property exchange board for Mariner/Cypress, standing broker offers, a normal sell_property action and deed-only Cypress buying. Sales preserve registered occupants and property accounts, permit explicit rent-back for the resident seller, release their posted property assignment and transfer business wage obligations. Payout is65% of base value scaled by condition; it does not increment earned income. A persistent bought_life prevents reacquisition respect farming. The board shows owner/offered status and locates the address while retaining the city as the main view.

Tests cover transfer/payout, retained home/tenants/accounts, duplicate-sale rejection, damage/invalid deeds, purchase without moving occupants and buy-sell-buy cash/respect loss. Full core/store/HTTP passed103.531s/.141s/.683s (.runtime/property-market-integration.log). Final targeted/long-city scan passed1.232s and TypeScript/Vite build passed2.33s (.runtime/property-market-tests-final.log, property-market-build-final.log). Fixed the singular resident label spotted in browser review; the final source/build includes it.

CUA isolated8887 property-market-20260913.sqlite3: board showed2340 broker offer for owned Mariner and3500 Cypress asking price. Inspect address returned to the correct city selection. Normal sale charged30minutes, paid2340 (12000→14340), retained Mariner home and23 tenants, removed rental income/staff bill and showed15/day personal room rent. Market then offered Mariner at3600. No browser errors; main campaign untouched. Unrelated mugging/robbery/armed files remain unstaged.

This is a standing-broker market, not finished NPC-driven real estate: individual buyers/funding, timed listings, additional housing and explicit occupied-house move-in/relocation semantics remain. Cypress deed-only purchase preserves occupants; its legacy combined move-in action still needs that occupancy review. Broad production scope remains active.

### 2026-09-13 — capacity-checked occupied-house move-in

Added a read-only residence plan to move_home offers. Player capacity reservation can rehouse an existing resident into available accommodation, named in the offer; the move is refused if an existing housed person would lose housing. Completion rechecks the quoted allocation, commits residence changes only, and preserves workplaces, physical journeys and old rent debts. Changed arrangements postpone the move and refund its price while preserving the intervening city/time. Market resident counts now say other residents for the player's own home.

Tests: planning does not mutate the world; occupied Cypress move rehouses its resident; no vacancy rejects without mutation; residence changes do not teleport actors or erase arrears. An actual contract death at minute630 during a600→660 move verifies changed occupancy causes postponement/refund while retaining the death/time. Targeted tests passed (.runtime/home-move-final-tests.log); full core/store/HTTP passed103.797s/.154s/.702s (.runtime/home-move-integration.log); build/TypeScript passed (.runtime/home-move-build.log). Full integration preceded the added interruption test and a refund-wording adjustment; targeted interruption test was run against the same behavior.

CUA isolated8888 home-move-20260913.sqlite3: move-in offer explicitly said Commissioner Vance would be rehoused at Ashbury Court. Normal3500 purchase left player at Cypress with4500 cash; public revision1 showed Vance Home apartment while Location remained market, and zero shortage. Main campaign untouched. Automatic insolvency/new-life housing, more residential capacity, tenant notices/evictions and full physical interiors remain outstanding.

### 2026-09-13 — show daily cost breakdown immediately

Moved the itemized daily bill to the top of expanded Daily accounts, ahead of the secondary account figures, with spacing matched to the paper panel. It remains a plain section with no extra disclosure click. TypeScript/Vite build passed (.runtime/accounts-default-build.log). CUA on isolated8888 confirmed opening the panel shows the entire current $90 rent breakdown at scrollTop0, before the figures; no campaign commands were issued. Unrelated armed-resistance edits remain untouched.

### 2026-09-14 — Mercer Court residential expansion

Added district-zero Mercer Court on an existing vacant lot, with48 residential places, shared/private tenancies,12/25 daily NPC rents,120 player lease and25/day private rent. It uses apartment-tier protection, capacity-checked moves, existing residence routines and tenant registers. Existing tenant priority stays intact. Added superintendent/seamstress roles and address-specific event text. No save shape change; current-version campaigns receive the property through ordinary new-address settlement.

Authored its five-storey brick exterior in Blender with lintels, cornice dentils, downpipes, entrance sconces/letterboxes and rooftop laundry, retaining articulated entrance/glazing and rear escapes. Added a Blender lobby with48 numbered letterboxes, stair, bench, dado panels and terrazzo. The lobby is currently displayed as a rendered plate; the exported3D source is available for subsequent interactive interior integration. Reproducible targeted export: `.venv-blender/bin/python tools/export_city3d.py --only=mercer-court` (exports the two models and renders the exterior/interior plates). It leaves other reviewed models untouched.

Browser review exposed ground-level focus clipping tall buildings; Focus address now fits actual model bounds with the existing orthographic framing helper. Full initial integration also exposed missing address art/roles/event text, which were added. Preserved the pre-existing urgent unpaid-wage warning above the daily bill while keeping the bill ahead of secondary account figures and open without extra clicks.


Final evidence: core/store passed120.468s/.127s (.runtime/mercer-integration-final.log), HTTP/frontend-contract tests passed.614s (.runtime/mercer-http-verified.log). Frontend suite passed223/223; final geometry/framing rerun passed21/21 after final exports, covering every building footprint and all pedestrian/vehicle routes plus unchanged pre-existing map positions. TypeScript/Vite passed2.56s (.runtime/mercer-release-verified-build.log). Both rendered plates were visually inspected; the exterior plate was corrected to hide alternate damaged glazing and fit the entire building.

CUA isolated8890 Mercer final fixture (actual file `.runtime/mercer-final-20260914.sqlite3`):116 residents, no housing shortage. New building fully framed; normal70-minute walk reached Mercer Court. Lobby showed its actual occupants and48-tenant register. Normal120 lease named the displaced resident beforehand, completed with7880 cash,47 NPC tenants plus the player, zero shortage and25/day in the expanded cost breakdown. No browser errors. Current lobby is a rendered plate in the legacy interior layout, with simple overlaid occupants; dedicated interactive3D staging and the broader interior redesign remain unfulfilled. Exterior review showed about145FPS locally; this is not a general performance acceptance claim.

Main campaign untouched. Remaining housing work includes emergency/new-life capacity handling, housing assignment on births, individual NPC-funded property markets and lease/debt enforcement. Full high-fidelity visuals, event variants, interiors, progression/questline and campaign acceptance remain active.

### 2026-09-14 — interactive Mercer Court lobby

Mercer Court now loads its authored lobby GLB in the browser rather than the still plate. Interior3D selects the actual room asset/name, camera orientation, lighting and occupant layout for Mercer or Saint Agnes. Its bench has two seated positions; a public superintendent stands by the letterboxes; fourteen general positions plus the service position fit within the clear lobby floor. The renderer reports any additional people as available in the list, rather than silently implying everyone is on the floor. Room walls and mounted fixtures are grouped for camera-dependent cutaways. NPC picking uses the same public people/actions as the ordinary room list; no simulation commands or state semantics changed.

The residents' register now appears directly below the room, before occupants/actions. Moved camera hints out from beneath the persistent accounts HUD, removed an empty-cue length rendering a stray0, and limited service-dependent rent wording to the Mariner. Existing rendered plates were retained; only the lobby GLB/manifest changed for wall grouping.

Validation:224/224 frontend tests passed (.runtime/mercer-interior-tests.log), including actual male/female articulated GLB bounds for every lobby placement, no pairwise overlaps, floor/stair/wall clearance and bench seat support; legacy bar placement checks still pass. HTTP/frontend contract tests passed.618s (.runtime/mercer-interior-http.log). Final TypeScript/Vite build passed (.runtime/mercer-interior-release.log).

CUA isolated8890 at revision2/minute730: actual seamstress seated on the bench and superintendent by the register; selecting the superintendent directly on the3D canvas opened his own actions. Orbit hid both rear/left walls when behind them; zoom increased to1.12 and Home restored1. The47-tenant register is now immediately beneath the canvas. No browser errors. Isolated8891 fresh bar fixture confirmed Saint Agnes still loaded all five public occupants with its original booth/service positions. Reading/entering/selecting/orbiting did not advance game time; main campaign untouched.

Remaining: stronger interior material/detail quality, player staging, arrival/service/idle animations, accurate staging for larger crowds/other rooms, interior combat, and the wider production objectives. These two interactive rooms do not complete the city-wide interior requirement or campaign/performance acceptance.

### 2026-09-14 — player presence in interactive interiors

Passed the existing public player name/face/alive state into both3D rooms. The living player now uses the same portrait-selected model and wardrobe as the street, occupies reserved clear floor space and has a pale identification ring. Clicking the player clears NPC selection and returns to premises actions; it does not fabricate a player-as-NPC command target. NPC roster capacity is unchanged, and overflow counts exclude the player. Appearance/roster rebuilding keys only character identity/appearance/liveness, not cash or other unrelated player stats.

Validation:225/225 frontend tests passed (.runtime/interior-player-final-tests.log), including male/female model bounds against every supported room placement and floor/stair/stool/aisle limits. HTTP/frontend-contract tests passed.649s (.runtime/interior-player-http.log); TypeScript/Vite passed2.16s (.runtime/interior-player-final-build.log). CUA isolated8890 showed player plus both existing lobby occupants, selected the actual superintendent, then clicked the player to restore premises actions. Isolated8891 showed player plus all five existing bar occupants. No gameplay commands were sent and the main campaign was untouched. Browser error check for the lobby was empty.

Remaining: this adds presence, not walking entry/exit or service/idle animation. Other authored rooms, richer visual fidelity, interior combat and the broader city/campaign objectives remain outstanding. Interior draw count grows with the articulated cast (331 draws for the bar/player fixture); continuous-motion performance still needs dedicated profiling and batching work.

### 2026-09-14 — batch the interactive interior cast

Added InteriorCastBatch to instance matching character geometry/material sources across the cast. Original articulated rigs remain intact and provide world transforms; per-instance colors preserve each wardrobe tint. Hidden hair/headwear groups stay hidden. Instance picking resolves to the original NPC/player identity, and unsupported materials remain visible/pickable through the ordinary path. The batch owns only its neutral material copies and instance buffers; prototype geometry/textures stay shared. Disposal restores original mesh visibility and does not dispose shared geometry. Current static interior poses upload once on roster/appearance changes; camera movement does not re-upload unchanged instance buffers. Future pose animation must call batch.update when poses change.

Evidence:228/228 frontend tests passed (.runtime/interior-batch-final-tests.log), including actual authored cast geometry/tint preservation, live joint updates, visibility changes, instance ray picking, resource ownership/disposal and unsupported-material fallback. HTTP/frontend-contract tests passed.704s (.runtime/interior-batch-http.log); final TypeScript/Vite build passed (.runtime/interior-batch-release-build.log).

CUA isolated8891, unchanged six-character bar fixture:331→143 draws (56.8% reduction), unchanged166136 rendered triangles. Clicked the batched barman and received Aldo Olsen's own actions. Leave/re-enter restored the same cast and143 draws; no browser errors. Isolated8890 lobby:177→131 draws with unchanged152112 triangles and the same three actors. Screenshots visually reviewed for appearance and seating. No main campaign writes or simulation commands. Draw-call reduction is measured; general60FPS/crowded-scene GPU acceptance remains unproven.

This optimizes the two current interactive interiors. City traffic/effect batching, full interior animation, other buildings and the wider visual/gameplay production scope remain active.


### 2026-09-14 — bartender counter service animation

Saint Agnes's public bartender now wipes a clear strip of the marble counter with a locally authored Blender cloth. The articulated hand drives the cloth position; iterative contact fitting uses the actual hand geometry without changing arm lengths or moving the feet. The cloth includes folded linen geometry and blue woven borders. Targeted reproducible export: `.venv-blender/bin/python tools/export_city3d.py --only=bar-cloth`. Only an existing public service-role occupant receives this cosmetic activity; no new simulation action or clock advancement is invented.

Passed the existing scenes preference into interactive rooms. Motion pauses when scenes are disabled, the document is hidden, or the system requests reduced motion. Animated poses refresh the character instance batch; static rooms keep render-on-demand behavior.

Validation:230/230 frontend tests passed (.runtime/interior-service-tests.log), including360 sampled poses for each actual male/female GLB, hand/cloth contact, countertop clearance, cloth movement and unchanged feet/root. TypeScript/Vite passed2.53s (.runtime/interior-service-build.log). HTTP command/public-state/campaign tests passed.513s (.runtime/interior-service-http.log). Browser CUA isolated8891 showed all six actors,147 draw calls and167648 triangles. Successive live observations confirmed cloth movement; disabling scenes held both pose time17.5202 and rendered-frame count2519 unchanged across observations, and reenabling resumed motion. No browser errors; the screenshot was reviewed for staging/contact. System reduced-motion behavior is implemented but was not separately toggled in browser QA. Main campaign untouched.

This adds one service animation, not full arrival/exit, crowd idle behavior or interior combat. Other interiors, higher visual fidelity, full event variants, gameplay progression and comprehensive campaign/performance acceptance remain outstanding.


### 2026-09-14 — interior roster and action papers

Replaced the remaining dark-green interior roster and action catalogue surfaces with the existing city paper stock, ink, ruled borders and visible focus treatments. Public character roles are visible again beneath names; selected, player-aligned and hostile states retain distinct borders. Action descriptions are no longer clamped to three lines, keeping costs/consequences readable. Disabled actions remain legible and unavailable. Responsive grids can shrink below their old fixed card minimum; search controls wrap. Styling is scoped to interior controls rather than gambling felt or other scenes.

Corrected roster semantics: its containing group is named People in this building, and each occupant retains native button semantics and aria-pressed rather than overriding the button with listitem. Existing selection/action handlers are unchanged.

Validation: TypeScript/Vite passed2.27s (.runtime/interior-paper-build.log); git diff --check clean. Browser CUA isolated8891 visually reviewed the bar roster and Mara's selected action panel. Her button selected the right character and exposed four actions with pressed=true. Step away restored premises actions. Searching envelope returned the one matching job with its full detail; Tab reached Back to the street with a solid visible focus outline. Zero clipped action descriptions or horizontal document overflow at the checked desktop size; no browser errors. Narrow viewport behavior is implemented but not separately browser-verified in this pass. No gameplay commands or main campaign writes.

Broader room modelling, interior action animation, remaining menu outliers and full campaign/performance acceptance remain active.


### 2026-09-14 — authored interactive Mariner lobby

The Mariner now uses a dedicated Blender-authored boarding-house lobby GLB in the browser. Added oak floorboards with exported deterministic grain texture, painted dado panels/mouldings, reception desk with recessed panels, service bell, ruled rent ledger/pencil, folded linen,24 individually numbered room-key tags/hooks, corridor door/glazing, oak waiting bench, stairs with runner/spindles/rails, radiator and opal sconces. Rear/west walls and their mounted objects use the existing camera cutaway groups. Reproducible targeted export: `.venv-blender/bin/python tools/export_city3d.py --only=interior-mariner`; full export includes it too. No existing models were re-exported.

Added Mariner-specific staging: a public landlady/landlord/receptionist can occupy the actual reception aisle; two bench seats and nine clear floor positions hold other public occupants. A separate player entry position preserves NPC capacity. Overflow remains accessible through the people list. A shared room-placement selector keeps count and renderer assignments aligned. The public rent register and existing purchase/management/person actions remain directly below the room with the revised paper styling. No backend/API changes or invented residents.

Validation:231/231 frontend tests passed (.runtime/mariner-lobby-tests.log), including real male/female GLB geometry for all12 NPC positions plus player: feet above floor, no pairwise collisions, no reception/stair/wall intersections, bench seat support, stable placements and role assignment, and actual room bounds/cutaway groups. TypeScript/Vite passed2.10s (.runtime/mariner-lobby-build.log); HTTP command/public-state/campaign regression checks passed from cache (.runtime/mariner-lobby-http.log); git diff --check clean.

CUA isolated8891 (test save .runtime/interior-regression-20260914.sqlite3) used the normal15-minute walk from Saint Agnes to the Mariner and skipped remaining cosmetic travel playback. At10:15 the actual public boarder Sofia Quintero was seated on the bench; the player was visible near the entrance. Clicking Sofia on the3D canvas selected her own actions. Orbiting behind the room removed its rear wall; Home reset and plus zoom to1.12 worked. Expanded register contained23 tenants; no horizontal overflow at the checked desktop size or browser errors. Initial unselected room rendered131 calls/132592 triangles. There was no public landlady present at this point; her staging is covered by geometry tests, not claimed as browser-observed.

Main campaign untouched. This is the third interactive room, not all interiors; walking entry/exit, richer service/idles, interior combat, higher visual fidelity, gambling, remaining event variants and wider campaign/performance acceptance remain outstanding.


### 2026-09-14 — player walking entrances in interactive rooms

Added InteriorArrival for Saint Agnes, Mercer Court and the Mariner. The player walks along each room's reserved foreground aisle, with distance-based leg/arm swing, eased acceleration/settling and a final turn into the existing facing. Actual shoe geometry fits the gait to floor and entrance mats. Moved the Mariner's final player x from1.3 to1.6 so the resting feet clear the mat edge. Entrance is cosmetic: it neither advances the simulation clock nor blocks actions. Roster rebuilds reuse elapsed entrance time; ordinary public-state updates do not replay it. Hidden documents pause it, while disabled scenes/reduced motion settle immediately. Arrival and service share one cast-batch refresh per changed frame. Completed static rooms return to render-on-demand.

Validation:232/232 frontend tests passed10.962s (.runtime/interior-arrival-final-tests.log). The initial strict zero assertion distinguished negative zero from zero; fixed the assertion to check numerical rest. New geometry checks sample181 frames for each actual male/female model in all three interiors, covering every possible NPC position, furniture, floor/mat clearance, continuous movement, limb activity, final position/facing/rest and repeated completed samples. TypeScript/Vite passed2.16s (.runtime/interior-arrival-build.log); HTTP/public-state/campaign checks passed from cache (.runtime/interior-arrival-http.log); git diff --check clean.

CUA isolated8891 at10:15: Mariner entry observed moving=true at0.208s with player near the threshold, then moving=false at1.8453s at reserved(1.6,3.15). Reentry screenshot captured a stride across the mat; completed pose visually reviewed. Disabling scenes before reentry produced the final pose immediately with only3 renders; preference restored afterward and no browser errors. No gameplay commands or main campaign writes. Reduced-motion system setting is respected in code but was not separately toggled in browser QA.

This is a short player entrance within the current three room models. It does not yet animate leaving, NPC arrivals/departures, traversing upstairs, interacting with furniture or interior combat. The full city-wide interior, animation, visual quality and gameplay/campaign objectives remain outstanding.


### 2026-09-14 — actionable opening income guidance

Reworked NextOpportunity to quote the real destination action's label, duration, fee and details. A read-only value copy supplies the hypothetical destination for action inspection; no real location/state mutation or automatic execution. Early progression now falls back from exhausted/unavailable envelopes to cargo, recommends earning before unaffordable acquisitions/recruitment, looks for another available first business if the laundry is unavailable, and suppresses unavailable driver/repair/expansion/move hints. Existing own people count toward the recruitment step. Custody suppresses street-work suggestions. No reward, economy or save-schema changes. Guide's first milestone now quotes actual45-minute courier duration/readiness; removed its misleading claim that businesses are the only income while elsewhere.

The handbook labels its recommendation Your next move and describes its existing milestones plainly. Browser review caught the old generic card class inheriting playing-card dimensions, making the guidance overflow a tiny rectangle across the whole page. Replaced it with a dedicated paper-note layout and verified its actual bounds against the following content.

Validation: targeted guide/opportunity tests passed.157s (.runtime/guided-income-targeted.log). New regressions cover depleted/replenished envelopes, cargo terms, insufficient/funded property capital, actual recruitment terms, dead driver, custody and read-only state. Full core/store/HTTP suites passed110.666s/.142s/.724s (.runtime/guided-income-regression.log), including hidden-plan guidance checks. Final TypeScript/Vite passed2.83s (.runtime/guided-income-final-build.log); git diff --check clean.

CUA fresh isolated8892 running .runtime/guided-income-server against .runtime/guided-income-20260914.sqlite3: initial90 cash/0 respect at08:00. Handbook showed45 minutes/$45/2 respect and three envelopes. Find the address selected Saint Agnes; normal15-minute walk, skip of remaining travel presentation and normal courier action completed at09:00 with135 cash/2 respect. First milestone checked Done; next move showed2 of6 respect and two remaining envelopes. Corrected note measured760px wide,201px tall, without clipping or overlap; screenshots reviewed and no browser errors. Main campaign untouched.

This improves the existing opening guidance, not a complete authored quest campaign. Richer contacts/jobs, midgame income/progression, full20–30-minute narrative acceptance, visuals/interiors/gambling/events and broad performance acceptance remain active.


### 2026-09-14 — first-business campaign and first-associate milestone

Continued the fresh8892 campaign from the first envelope. The second envelope triggered A favor with a price; the slower route paid its quoted60 and3 respect over75 minutes. The final envelope brought cash285/respect9 at11:45. Guidance switched to cargo when envelopes ran out and correctly named the720 laundry buyout. Normal40-minute travel to Pier14 and six90-minute cargo shifts brought735 cash/respect15 at21:25, with no injury in this run. Normal15-minute travel to Bluebird and its720/60-minute acquisition ended22:40 with15 cash/respect19. Daily accounts showed322/day current income,33/day costs (15 rent+18 staff),289 net; the business milestone completed. All were normal UI commands in .runtime/guided-income-20260914.sqlite3; cosmetic journey playback was skipped where noted, never simulation time.

Found a real progression inconsistency: the next-move hint advised saving90 for the driver, while the crew milestone claimed an organization was required before anybody could join. Guide now checks both initial recruit and subsequent sign-on actions at unlocked destinations without moving the player. Either existing driver crew or organization members complete the milestone. Text distinguishes the two routes. Also renamed the always-available dock job Work a cargo shift, since calling it night cargo at12:25 misdescribed the displayed situation; command/reward/time/risk unchanged.

Targeted guide/promise/opportunity tests passed.224s (.runtime/guide-hiring-targeted.log). New tests verify the initial driver is available before incorporation at90 cash, blocked at89, dead drivers are not offered, guide reads preserve location, and an actual recruit command completes the milestone. Full core/store/HTTP suite result is recorded in .runtime/guide-hiring-regression.log. No frontend implementation changed in this pass.

Copied the isolated campaign using SQLite backup into .runtime/guide-hiring-20260914.sqlite3 and served the new binary at8893. Browser confirmed both hint and milestone now identify insufficient cash, and the cargo label is correct. No browser errors. Main campaign untouched.

Remaining evidence: the opening saving loop required six identical cargo clicks after the authored introduction; this is too repetitive to establish the requested rich progression/quest campaign. More early earning variety and midgame contracts remain needed. The guard-posting milestone can still prefer an irrelevant not-your-business reason across candidate addresses; record for the next guidance correction. This short targeted progression playtest is not the full20–30-minute narrative acceptance. All broader visuals, interiors, event variants, gambling and performance goals remain active.


### 2026-09-14 — limited daily laundry work

Added a hotel linen rush order at Bluebird Laundry as a second limited early-income route before falling back to repeatable cargo. The action pays80 for60 minutes, adds up to one respect below DockName, and adds no job heat. It is available08:00–17:00 to a nonowner when the laundry has staff, two supplies, condition40+ and no operating trouble. It teaches visiting a future business and its operating requirements; an owner instead receives normal business income.

Save schema19 adds a per-property rush_order_day reservation. Acceptance consumes two supplies and reserves the current day before advancing time. Interrupted work pays nothing and retains the reservation/resources already used. It is deliberately not a resumable arrangement. Reservation survives reload, duplicate receipt retries, ownership/life changes; next day reopens subject to operating conditions. Old saves start with no reservation and keep current supplies/condition. No frontend command shape change; Work grouping and normal action UI handle it.

Targeted core/store regressions passed.207s/.175s (.runtime/rush-order-targeted.log); extra version18 save-upgrade/receipt tests passed.131s (.runtime/rush-order-save.log). The initial full suite found its offer-coverage fixture only represented owners, so it could never see this helper job available. Expanded that actual coverage to include a new arrival rather than exempting the action; final targeted checks passed.214s/.151s (.runtime/rush-order-final-targeted.log). Final full core/store/HTTP result is recorded in .runtime/rush-order-final-regression.log.

CUA fresh8894 (.runtime/rush-order-server with .runtime/rush-order-20260914.sqlite3): started90 cash/0 respect at08:00, normal25-minute travel to the laundry, skipped remaining cosmetic travel, entered and reviewed the complete job terms in the paper catalogue. Normal rushorder command ended09:25 with170 cash/1 respect/0 attention, and the action became disabled with Today's hotel order has already been taken. Reload/reentry preserved this refusal. No browser errors; main campaign untouched.

This adds one daily external customer job, not a complete job market or fully modeled hotel/NPC customer budget. Broader earning variety, authored progression/campaign acceptance, 3D laundry/work animation and all other unfinished visual/gameplay goals remain active.

### 2026-09-14 — owned-premises Guide prerequisites

Fixed the campaign defect recorded after the first laundry purchase: the guard and still milestones now consider only owned, unlocked income premises (and eligible still sites). Refusals from unrelated addresses no longer replace the actual next prerequisite. With no eligible property, the Guide explicitly asks for the required acquisition. PostReadiness now distinguishes having no organization members from having members already assigned; the same wording reaches both Guide and action buttons. Posting rules, cost, time and command shape are unchanged.

Targeted Guide tests passed (.runtime/guide-premises-targeted.log, .218s), covering no property, ownership with no guard, unstaffed/funded/underfunded still, actual still construction and actual guard assignment. Full core/store/HTTP regressions passed (111.071s/.148s/.802s), recorded in .runtime/guide-premises-regression.log. Browser QA on8895 used an SQLite backup of the isolated first-business save (.runtime/guide-premises-20260914.sqlite3). At cash15/respect19/22:40 with Bluebird owned, the visible Guide says Sign someone into your organization before assigning a guard and Not enough cash for the still. Main campaign untouched. The previous daily-accounts request was also browser-verified on8894: one expansion immediately shows the15 rent item and total, with zero nested disclosures.

The broader goal remains active: only three interactive 3D interiors exist; laundry interior authoring has not started. Remaining richer interiors, choreography, gambling, campaign narrative and performance acceptance are not established by this guidance correction. Other Guide milestones still use generic location-sensitive refusals and need separate review.

### 2026-09-14 — Bluebird Laundry interactive interior

Authored Bluebird's fourth interactive 3D interior in Blender, replacing its flat room plate: three enamel industrial drum washers with polished rims, rubber seals, perforations, locking handles, numbered controls and supply pipes; glazed wall tiles, quarry floor, steam main, folding/collection counter, cotton-weave linen bundles, shelving, waiting bench and opal work lamps. Model and packed deterministic cotton texture are original local work in tools/export_city3d.py; no third-party assets. Reproducible targeted export: .venv-blender/bin/python tools/export_city3d.py --only=interior-laundry. The GLB is6,325,108 bytes with10x10m floor and4.2m walls. Full exporter also includes it.

Actual public laundry workers/clerks receive the collection-counter spot; other public occupants use the bench and clear floor bays, with overflow retained in the HTML roster. Player entrance uses a separate foreground aisle. Camera control, public person picking, existing actions and reduced-motion arrival behavior share the established room renderer. Equipment is presently at rest: operating-state-linked machine cycles and staff work animation remain outstanding.

Validation: all233 frontend tests passed (.runtime/laundry-tests.log,10.978s); TypeScript/Vite build passed (.runtime/laundry-build.log,2.35s). New geometry tests load the actual GLB and both male/female rigs, test full-occupant spacing, floor support and radial furniture clearance, while entrance tests cover181 frames for each rig. CUA at8894 used the isolated rush-order save, still09:25/cash170/respect1. Actual public Emil Lenz (Laundress) appeared behind the counter; direct canvas click selected street-41 and opened the matching actions. Rotation removed the rear wall correctly, zoom reached1.12, Home reset, browser error log empty. Default two-person room recorded90 draw calls/221,744 triangles; this is a diagnostic count, not a60FPS acceptance claim. Main campaign untouched.

Broader room quality and remaining interiors, service/walking/interior violence, gambling, story/campaign and overall performance goals remain active. This first laundry room does not establish final visual production acceptance.

### 2026-09-14 — public-state laundry machinery motion

Split the three Bluebird washer drums into authored glTF pivot groups, preserving material batching within each. Only perforations and the visible linen load rotate; cabinets, seals, hinges and handles remain fixed. New laundryMotion.ts maps positive public trading to one–three active machines, gated by staff, supplies, condition40+ and no operating trouble. The figure is cosmetic workload, not a new production rule or clock. Missing operation data leaves equipment stopped. Interior receives existing public Place fields; no API or save change.

Each machine rotates at a slightly different fixed speed. Scenes off, reduced-motion preference and hidden pages pause the animation without catch-up; per-frame delta is capped at50ms. Static cast buffers are not rebuilt for machine movement. Individual fabric tumbling/deformation, realistic cycle reversal and staff folding animation remain future work; the present linen rotates with its drum.

Validation:235 frontend tests passed in11.573s (.runtime/laundry-motion-tests.log); TypeScript/Vite build passed in2.41s (.runtime/laundry-motion-build.log). Actual exported geometry was sampled for720 frames: contents stay within the rubber-seal radius, the third machine stays still at a two-machine workload, all fixed furniture matrices remain unchanged, and pause/resume clamps correctly. Existing interior collision, selection and entrance tests remain green. Browser8894 isolated rush-order campaign: three machines visibly moving,101 draw calls/221,664 triangles, no browser errors. Scenes off froze machine seconds19.2557 and rendered frames2771 across checks; restoring Scenes resumed to29.2833/4215. UI-only checks kept cash170, time09:25 and the main campaign untouched. Browser operating-damage scenarios and OS-level reduced-motion switching were not exercised; their gating is unit-tested. Broader performance/visual acceptance remains open.

### 2026-09-14 — laundry worker smooths linen

Added a sewn linen piece on the collection counter and a two-hand smoothing task for the actual public worker occupying that station. The worker stands15cm closer to the rear counter edge, with a separate service clearance allowance; other room spots remain unchanged. LinenPress fits the real hand undersides to the cloth at1.109m, keeps rigid sleeve lengths and planted feet, and returns the arms to rest when public operation becomes unavailable. Scenes/reduced-motion/hidden-page controls pause it; active motion shares the frame's single cast-buffer update with entrance animation. No worker is invented when the public roster lacks one. No API, gameplay time, production or save change.

Initial .runtime/linen-tests.log caught a sleeve penetrating the counter. Changed the elbow bend hint to keep sleeves above the work surface. Final full236 frontend tests passed in10.838s (.runtime/linen-tests-final.log). Broadened the sleeve check to all upper/lower sleeve vertices, reran all three service tests successfully (.runtime/linen-sleeve-check.log,.132s). Final TypeScript/Vite build passed2.13s (.runtime/linen-build-final.log). Both exported character rigs are checked through360 frames for hand/cloth contact, no hand crossing, sleeve clearance, fixed feet, motion amplitude, rest reset and pause/resume delta limits.

Browser8894 isolated rush-order save: public Emil Lenz visibly smooths the cloth, direct canvas selection still opens Emil's actions, zoom and Home work. Scenes off held both linen and machine clocks at19.0546s and rendering at2742 frames across checks; re-enabling resumed both to24.8044s/3570 frames. No browser errors; cash170/time09:25 and main campaign untouched. Static room geometry now6,331,432 bytes. This is a repeated smoothing task, not a full fold/stack or material-deformation simulation. Full working routines, more interiors, character/art quality, interior events, gambling and broader campaign/performance acceptance remain unfinished.

### 2026-09-14 — authored 3D slot cabinet and committed reels

Replaced the flat CSS bandit with an original Blender-authored Lucky Bell cabinet: rounded oxblood/ivory enamel body, metallic face trim, recessed curved reel papers, modeled coin slot, side lever with retained pivot, and a payout tray with coins. Source tools/export_city3d.py; targeted export --only=slot-cabinet. No external asset packs. SlotCabinet prints the existing locally authored reel symbols onto three curved surfaces and scrolls the existing drumRun sequence through them. Outcomes, stakes, probabilities, payout multipliers and commit timing remain Go-owned. The visible payline uses the exact backend symbol IDs, with left-to-right stops on the existing700/1150/1600ms schedule. These are scrolling textures on authored curved papers, not fully rotating twenty-face cylinder meshes.

The modeled lever is clickable, and the existing keyboard-accessible Pull button and stake input remain. Results and payout coins wait for all reels to settle; payout text now states the dollar amount as well as multiplier. Tray coins are illustrative, not a literal count of the money paid. Fixed a pre-existing reload behavior: Machine now treats the saved result as settled on mount rather than replaying its animation/sound. New committed revisions still animate normally. Existing reduced-motion handling remains; a broader Scenes preference integration for gambling and live preference changes is still to review.

Validation: all237 frontend tests passed (.runtime/slot-tests.log,11.197s). Final TypeScript/Vite build after reload fix passed (.runtime/slot-build-reload.log). New test loads the GLB, raycasts every reel at three heights to prove unobscured printing and exact middle UV at the payline, checks upright UV orientation, samples the lever's full swing outside the cabinet, and checks coins stay in the tray. CUA8894 isolated rush-order campaign: normal25-minute travel to Saint Agnes, cosmetic skip, entered and used Play the machines. First direct modeled-lever click wagered5 and returned CHERRY/CHERRY/PLUM with25 paid, cash170→190. Visible line matched the backend; updated tray showed coins. Reload preserved the paid result without spinning. Two further5 pulls returned ORANGE/ORANGE/BELL and PLUM/LEMON/ORANGE, both zero payout. Immediate capture during the final spin showed all three rolling, disabled pull and no result paragraph. Final public state revision7/minute635/cash180 and line plum/lemon/orange matched the settled browser. No browser errors. Main campaign untouched.

Default cabinet uses12 draw calls/7,896 triangles when tray empty; production-wide60FPS remains unproven. Casino outer styling, table occupants, 3D card/dice/roulette tables, full payout motion and remaining game/visual requirements are still outstanding. The goal remains active.

### 2026-09-14 — slot motion preferences and cancellation

Connected the main Scenes preference through Casino to Machine. The slot now listens for live prefers-reduced-motion changes rather than checking only at the start of a pull. Either setting disabling motion clears rolling state, cancels remaining reel/coin timers, and reveals the committed outcome immediately. Re-enabling motion keeps the already-seen revision settled. No gambling arithmetic, command or save changes. Other table games still need a separate preference audit.

Built successfully (.runtime/slot-motion-build.log,2.34s); all237 frontend regressions passed (.runtime/slot-motion-regression.log,11.129s). CUA fixture at8897/.runtime/slot-motion.html exercised the actual Machine component with a visible motion checkbox and a counted fixture-pull callback, without campaign commands. Active pull reported all three rolling/disabled; switching motion off during that spin settled at rendered frame45 with25 paid and one fixture pull. Re-enabling kept frame45 and count1, proving no replay. A second pull with motion disabled settled immediately at frame46/count2. Initial fixture import used a raw Vite dependency path that lacked a named createRoot export; fixed it to normal react-dom/client import before running those checks. OS reduced-motion preference itself was not changed; its shared animate-gate behavior was reviewed but not browser-toggled.

Real isolated game8894: left the slot, disabled Scenes through Settings, entered Saint Agnes and reopened machines. A normal5 pull immediately returned BELL/ORANGE/LEMON, zero payout, all rolling false, four total cabinet renders, and enabled Pull. Cash180→175. Restored Scenes afterward, browser error log empty. Main campaign untouched. The broader game, visual fidelity, gambling table/occupant work and performance requirements remain active.

### 2026-09-14 — game-room stationery and compact scrolling

Restyled Casino's outer room with paper/ink header, cash-on-hand display, integrated leave control, printed game tabs and a session slip. The slot stake field, pull control and paytable use the same paper treatment around the green playing surface. Scoped the sheet to casino-house, leaving the separate back-room screen out of this unreviewed change. Header wording distinguishes a machines-only venue from a gaming room.

The casino floor now owns the main scroll; its game column no longer creates a second nested scrolling area. Desktop session slip remains sticky with its own bounded long-record scroll. At900px and below it follows game content naturally, without reserving a fixed-height row that squeezes the machine. Compact header wraps and gives the leave control a full row at520px.

TypeScript/Vite final build passed2.25s (.runtime/casino-room-build-final.log), diff check clean. Browser8894 at1234x1051 showed the header, cabinet, paper paytable and slip with no document horizontal overflow; computed main floor overflow auto/game visible. A separate iframe fixture (.runtime/casino-responsive.html on8897) exercised the actual isolated game at460x680. Screenshots verified the wrapped leave control, full stake controls, complete paytable and empty/populated session slip by scrolling the floor to the bottom. A normal compact5 pull ended cash175→170; receipt text wrapped fully (ORANGE/LEMON/PLUM, no payout). No browser errors. Main campaign untouched. No new unit tests for this cosmetic layout change; existing game behavior was exercised in the browser.

Blackjack/roulette/craps visual scenes and occupants remain outstanding; their full layouts were not re-accepted by this slot-focused pass. All broader production requirements remain active.

### 2026-09-14 — modeled dice and felt rolling tray

Replaced flat shaking dice icons with two original Blender-authored rounded ivory dice, modeled pips on all six faces, and a walnut/padded felt tray with packed deterministic baize texture. Target export --only=dice-table creates gaming-die.glb and dice-tray.glb; full exporter includes both. New DiceTable3D displays only supplied dice (empty tray before a roll), and PresentedDie animates translation, tumbles and diminishing bounces into the exact committed upper faces. It computes floor clearance from actual exported vertices rather than approximate cube bounds. Gameplay remains Go-authoritative; no reroll, wager, odds, outcome or save change.

Craps uses a shared1100ms presentation duration, settles immediately with Scenes disabled or reduced motion, and listens for live reduced-motion changes. Saved throws remain settled on mount. Its own result paragraph waits for motion to finish; the session-record column can still disclose fresh records during a throw and needs a separate coordinated presentation fix. Full craps layout, patrons/dealer, throwing hands, wall rebounds and realistic rigid-body motion remain outstanding; this is an authored presentation path on a tray.

All238 frontend tests passed (.runtime/dice3d-tests.log,10.892s). The new test loads actual GLB geometry, verifies pip counts1–6, checks all36 face pairs over121 samples each for correct top faces, floor contact, tray bounds and no mutual overlap, and hides invalid/unknown faces. TypeScript/Vite build passed (.runtime/dice3d-build-final.log). Initial browser log warned that PCFSoftShadowMap was removed in the installed Three runtime; switched to its supported PCFShadowMap and verified the final browser log is empty.

CUA8894 isolated campaign: left Saint Agnes, normal25-minute travel to The Monarch with cosmetic skip, entered, sat and selected Craps. Chose5 instead of the default50. First normal throw showed [4,4] while rolling with no result paragraph, then visible upper4/4 and point8. Next roll produced6/2 and made the point, ending cash175/minute725 (5 net profit on this round). Public /api/state matched the dice and outcome exactly. Final reload/reselect opened6/2 settled at three renders with no animation replay. Render diagnostic33 calls/9,396 triangles, not a production60FPS claim. Main campaign untouched; all broader game and visual requirements remain active.

### 2026-09-14 — coordinated gambling result reveal

Added a presentation-state callback from Machine/Craps to Casino. Reel/throw lifecycle changes run in layout effects so the session slip is gated before paint; records collected during a command remain behind a waiting message until the final reel/throw settles. Preference cancellation, game switching and component cleanup release the gate. Dice also suppress the newly supplied point and point-specific continuation note while tumbling. Fixed a second discovered issue: a settled winning/losing dice result previously enabled the next wager during its animation; the new wager button now stays disabled until it ends.

This gates the session slip, local result prose and dice point. The authoritative cash display still updates when Go commits; no gameplay state is deferred or changed. Roulette/card presentation and complete suspense across every possible surface remain separate work.

TypeScript/Vite build passed2.12s (.runtime/casino-reveal-build-final.log); all238 frontend tests passed10.980s (.runtime/casino-reveal-tests.log). CUA8894 isolated Monarch campaign: a5 pass-line throw supplied3/4. Immediate capture showed rolling true, point dash, no outcome paragraph, waiting session slip and disabled next wager; settlement revealed the natural-seven win and both matching records. A5 slot pull supplied7/cherry/lemon: all reels rolling, paid0 presentation and waiting slip initially, then5 paid and its records after settlement. Another slot pull followed immediately by switching to Craps released the gate and showed bell/bell/cherry's committed record, with no stuck waiting state or browser errors. Main campaign untouched. All larger 3D table/occupant, character/art, story and performance goals remain active.

### 2026-09-14 — roulette reveal and animation lifecycle

Roulette now adds the committed pocket to its recent-results strip only when the ball settles, and reports presentation state to Casino so the session slip also waits. The effect runs before paint. It honors Scenes and live reduced-motion preferences, cancels both browser animations and its timer on cleanup, and opens an existing saved spin settled instead of replaying it. Authoritative cash still updates immediately; this does not claim complete suspense across every surface. No backend or API changes.

CUA isolated campaign8894: normal entry into The Monarch and a5 red wager. Immediate observation showed falling=true, no recent-result strip, and the waiting session slip. After settlement, hub/history both showed36 red and the slip reported10 returned, cash185. A fresh browser tab and Roulette selection showed36 settled immediately. Browser error log was empty. Preference branches are implemented but were not independently exercised in this browser pass. All238 frontend tests passed10.751s (.runtime/roulette-reveal-tests.log), TypeScript/Vite build passed (.runtime/roulette-reveal-build.log), and git diff --check passed. Main campaign untouched. Full3D blackjack/roulette tables, gamblers, richer animations and the broader production goal remain unfinished.

### 2026-09-14 — authored 3D blackjack table and public hands

Added reproducible Blender exports blackjack-table.glb and playing-card.glb (targeted exporter --only=blackjack-table, also included in full export). The oval table has walnut apron/legs, rounded oxblood rail, brass reveal/feet, woven baize texture, printed house rules and a four-channel chip rack. Separate rounded card stock has a UV-mapped face. Browser canvas textures print the actual public rank/suit, numeric pip layouts and a patterned back for unknown cards. Dealer concealment follows the existing HandState. These are initial table/card assets: dealer/player bodies, illustrated court cards, dealing motions, player wager chips and richer environment remain unfinished.

BlackjackTable3D replaces the old flat card area, retains visible hand descriptions/totals and keyboard actions, updates from authoritative hands and draws only when assets/state/size change. Cleanup releases resources. No gameplay/API changes. Preserved shared Row/PlayingCard helpers used by the other card game; an initial build caught their accidental removal and they were restored before final verification.

The new geometry test loads the actual GLBs, checks card corner support on the felt and vertical clearance, plus visible UV face orientation. Its initial extreme long-hand sample found cards reaching the leather; narrowed the layout's maximum center span from2.25m to1.7m. Corrected test passes. All239 frontend tests passed12.309s (.runtime/blackjack-tests-final.log), final TypeScript/Vite build passed2.50s (.runtime/blackjack-build.log), diff check clean.

CUA isolated8894 campaign at The Monarch:5 wager dealt5 clubs/5 hearts against8 diamonds and one face-down card. Screenshot verified orientation, pips and readable public summaries after the final asset update. Hit added king of diamonds (20); stand revealed dealer8 diamonds/2 spades/3 spades/5 spades (18), hidden count0, with matching10 return. No browser errors; final table23 draw calls and5 rendered frames. This is static hand presentation, not dealing choreography or a production60FPS benchmark. Main campaign untouched. The full goal remains active.

### 2026-09-14 — blackjack card movement and reveal gating

Added a public-hand presentation plan for opening deals, hits and dealer settlement. Opening cards alternate player/dealer; later draws animate only changed/new cards while existing cards reposition to the updated spread. Cards follow an eased raised path from the dealer side onto the authored felt. Saved hands open at rest. Final poses use their exact target coordinates. CardTable stays mounted while awaiting the initial wager so a new deal is distinguishable from loading a saved hand.

CardTable gates hand descriptions/totals, actions and outcome during movement and reports the same state to Casino, which holds session records and the next-deal controls. Scenes/reduced-motion changes immediately settle, timer/component cleanup releases the gate. Authoritative cash still updates immediately. This remains an initial dealing presentation: no dealer hands yet, and the newly revealed hole card currently arrives from the dealer side rather than physically flipping the existing back. No claims of complete choreography or suspense across every surface.

CUA isolated8894: saved20-v18 hand opened settled. A5 new deal showed J diamonds/2 hearts against7 clubs plus a back, dealing=true, no hand summaries/outcome and waiting session slip. Hit added10 spades; while it moved the bust text remained hidden. After it landed, dealing=false, the matching loss and session record appeared. No browser errors. Final alternating opening order is additionally covered in the pure presentation test. Tests cover exact public cards, unknown back, old-card repositioning/new-card-only dealing, dealer draws, above-felt heights and immediate settled/no-motion poses. Initial strict endpoint comparisons exposed floating-point residue; final poses now return exact destinations. All241 tests passed11.508s (.runtime/blackjack-deal-tests-final.log), TypeScript/Vite build passed (.runtime/blackjack-deal-build-final.log), diff check clean. Main campaign untouched. Full goal remains active.

### 2026-09-14 — physical blackjack hole-card turnover

The presentation plan distinguishes an existing unknown dealer card becoming public from a new draw. It turns that card over at its existing seat (with spread adjustment if the dealer has further draws) instead of dealing it again. Card models now have a separately printed reverse surface so the back remains visible before rotation exposes the public face. Lift follows the scaled card half-width through the rotation, keeping the stock above the cloth; no gameplay randomness or additional card identity is introduced. Dealer hands and the larger gambling scene remain outstanding.

Actual-GLB test samples121 turnover poses, checking stock/felt clearance and face normal at both endpoints; it also checks that no-motion reveal settles immediately. All242 frontend tests passed11.347s (.runtime/blackjack-flip-tests.log), TypeScript/Vite build passed (.runtime/blackjack-flip-build.log). CUA isolated8894: a5 hand showed7 hearts/10 clubs againstJ diamonds plus a back. Standing captured the physical card edge-on with flips=1, dealing=true and no result. Settlement showedJ hearts as the second dealer card,17 against20 and the matching5 loss. Browser errors empty. Main campaign untouched. Restarted caffeinate after the previous process was absent; current command session40950. Full production goal remains active.

### 2026-09-14 — stable blackjack controls and compact hand summary

Kept the hand-summary footprint during dealing, with an overlaid status and visibility/aria gating for the unrevealed totals. Hit/stand controls now remain present but disabled while cards move; the stand label also conceals the new total. Settled outcome text reserves its space while remaining hidden until completion. Compact screens use separate dealer/player rows with wrapping card descriptions and right-aligned totals; removed inherited80px seat minimum heights and extra dealer/player padding that bloated these rows.

Browser fixture on8897 exercises actual CardTable without campaign commands. Turning Scenes-equivalent motion off during an opening deal released the gate immediately. A normal deal retained112px summary height both during and after playback before the compact-padding refinement. A motion-off stand immediately showed17-v20 with dealing=false/flips=0. Final420x680 iframe screenshot verified two compact hand rows and visible Another card/Stand on17 controls, with vertical scrolling available. The fixture needed its own body overflow override because standalone style.css assumes the game's normal scroll container; that is fixture-only. TypeScript/Vite build passed (.runtime/blackjack-stability-build.log), diff check clean. No backend or save changes; full production goal remains active.

### 2026-09-14 — blackjack's actual public croupier

Casino now selects a public room occupant with a dealer/croupier role (deterministic ID order) and passes that person to the card table. Presence alone never makes a patron a gambler. The renderer uses the existing authored person/woman rig and portrait-driven wardrobe, posed behind the dealer rail, and prints the person's name in the table header. Rooms without a known dealer remain unnamed. Appearance reloads on dealer identity/face changes and follows resource cleanup. No simulation/API changes.

First browser render clipped Dante Draga's head; shifted the camera's target upward/back for a populated dealer position. Final CUA8894 screenshot includes the complete head and table, diagnostic dealer street-34 matches the named public croupier,43 draw calls/four settled frames, no initial browser errors. No campaign commands were needed this turn. Actual-rig vertex checks confirm both models avoid the table apron/rail, and selection tests exclude a regular family patron. All243 tests passed11.417s (.runtime/blackjack-dealer-tests.log), TypeScript/Vite build passed (.runtime/blackjack-dealer-build.log), diff check clean. An initial test runner import error was corrected with an explicit .js module specifier. This is a posed dealer; hand-to-card choreography, player seating and higher-fidelity character production remain unfinished. Full goal active; main campaign untouched.

### 2026-09-14 — blackjack table ready before the first bet

CardTable now renders the empty authored table as soon as blackjack is selected, labeled with the current casino and Bets open. It suppresses empty totals/hit/stand controls and does not invent a hidden card until a hand is actually playing. The opening deal can therefore animate on a scene loaded while the player chooses their wager. Found and fixed a resource ownership mistake introduced with the dealer: cloned wardrobe materials shared the card-material cleanup array, so card rebuilds disposed resources still used by the dealer. Wardrobe resources now have their own lifetime and are released only with the scene.

CUA8897 actual-component fixture: before any wager the table reported empty hands, hidden0, dealingfalse, two rendered frames/nine draws and no card actions. Fixture deal then showed dealingtrue on that loaded scene. CUA8894 isolated campaign: a5 hand dealt3 clubs/K clubs against7 clubs plus back, with Dante Draga/street-34 still rendered in the matching wardrobe after card rebuilds. Screenshot verified appearance; finished the hand by standing on13, no browser errors. All243 frontend tests passed11.167s (.runtime/blackjack-ready-tests.log); TypeScript/Vite build passed (.runtime/blackjack-ready-build.log), diff check clean. No main-save mutations. Full game goal remains active, including player seating and coordinated dealer motion.

### 2026-09-14 — player seated at the blackjack table

Added an original Blender gaming-chair export (walnut frame, oxblood cushions, brass foot caps/upholstery tacks; --only=gaming-chair and full-export path). The public player identity now flows from main through Casino/CardTable into the renderer, selecting the same rig/wardrobe as the city/interiors. Only a living supplied player is rendered. The chair/player use a reserved left-front seat, separate from the public dealer and the cards. Measured seated shoe bounds required a .665m cushion top; chair and rigid seated pose share that value. No API/save changes.

Actual-GLB tests verify both person/woman rigs have shoes above the floor within15mm, raycast cushion support at hip position, and no vertices inside the table apron/rail. All244 tests passed10.930s (.runtime/blackjack-player-tests.log). TypeScript/Vite build passed (.runtime/blackjack-player-build.log), diff check clean. First CUA8894 screenshot cropped the player/chair; pulled camera back and retargeted for the two-person scene. Final screenshot shows full chair and player, dealer and unobscured cards, matching public Alex Varga and street-34/Dante Draga. Diagnostic75 draws/three settled frames, not a production FPS claim. No campaign commands this turn; main save untouched. This is static seating, not coordinated hand-to-card choreography. Higher-fidelity character production and the full goal remain unfinished.

### 2026-09-14 — grounded blackjack floor and lighting

Authored gaming-floor.glb:400 alternating parquet strips in five wood tones with physical seams/bevels, plus a woven burgundy/gold rug. Added --only=gaming-floor and full-export inclusion. Rug top is at the existing scene floor height0, parquet below it, so previous furniture/foot clearances remain meaningful. Enabled a warm2048 directional shadow map for table/cards/participants/chair and receiving surfaces, with lower ambient fill. Scene still renders only when dirty or presenting.

First screenshot exposed the common box UV mapping repeating the rug borders through its middle. Changed the rug to one normalized full-face mapping; final screenshot has a continuous border and subtle repeating center motif. Added actual-GLB ray tests for rug support under cast/furniture, height0, exact UV mapping and floor bounds. All245 tests passed11.792s (.runtime/blackjack-floor-tests.log), TypeScript/Vite build passed (.runtime/blackjack-floor-build.log), diff check clean. Final CUA8894 screenshot verified parquet/rug and visible table/chair contact shadows, with no browser errors. Diagnostic156 draws including shadow work/four settled frames; this is not animation FPS acceptance, and shadow cost remains part of broader performance work. No campaign commands or main-save changes. Room walls, richer lighting/decor, coordinated participants and the full production goal remain unfinished.

### 2026-09-14 — blackjack playback resource stability and measured frame pacing

Card mesh/texture rebuilds now depend only on card identities; changing playback timing or reaching settlement no longer disposes/recreates the same cards. A separate presentation-state key still requests the final pose/render, and idle scenes skip pose work as well as drawing. Added read-only diagnostics for build count and visible playback frame intervals (FPS,95th percentile and worst interval). Hidden intervals are excluded by resetting the previous-frame timestamp; metrics are local browser pacing evidence, not GPU timing or a cross-device benchmark.

CUA8894 full loaded scene including both people, floor and shadows: saved hand cardBuilds1, new5 deal cardBuilds2, still2 after settlement.3 spades/A hearts againstQ spades+back played234 frames at143.996FPS, p95 7.1ms/worst8.2ms. Standing on14 revealedQ diamonds: builds3 during and after flip,70 frames at143.960FPS, p95 7.8ms/worst8.2ms. Browser errors empty. All245 tests passed11.783s (.runtime/blackjack-performance-tests.log), TypeScript/Vite build passed (.runtime/blackjack-performance-build.log), diff check clean. Only the isolated8894 campaign was used; main untouched. This verifies one local blackjack sequence above60FPS and does not close whole-game performance or production acceptance. Full goal remains active.

### 2026-09-14 — explicit ordinary gunfire weapon previews

Audited the preview catalogue against core/strike_presentation.go: all currently emitted assassination variants already have entries. Drive-by/wire/interior variants are still missing from the implementation, not merely the menu. Added ordinary shotgun and Thompson gunfight previews; the existing Gunfight now explicitly supplies revolver tier1 instead of relying on an unattributed legacy fallback. Preview cues carry only the synthetic attacker and gun tier, with no invented death/victim or committed records. A proposed ordinary melee preview was removed before finalization after inspecting that path's incomplete choreography.

CUA8894 after normally leaving the tables: shotgun preview diagnostic weapon=shotgun/tier2, Thompson=thompson/tier3, both staged at The Monarch. Shotgun screenshot verified visible held model. Across replays and Stop preview, revision36/minute1160 remained unchanged; stop emptied effects. Recorded audio diagnostic counted four shotgun shots over two runs and six Thompson shots, no synthesized fallbacks, browser errors empty. Main campaign untouched. Added tests for exact tiers, nonfatal cues and unchanged revision/time; all246 tests passed11.567s (.runtime/preview-weapons-tests-final.log), final TypeScript/Vite build passed (.runtime/preview-weapons-build-final.log), diff check clean. Full production goal remains active, especially missing action variants and coordinated choreography.


### 2026-09-14 — authoritative incendiary attack and fire presentation

Added the incendiary action at nonowned businesses:40 dollars/15 minutes,15–30 bounded condition damage, up to10 supplies destroyed, operating trouble,18 capped attention, faction goodwill loss and existing retaliation scheduling. Uses the saved building-fire response without explosive casualties, staff removal or bankroll destruction. API.md documents command/cue semantics and remaining choreography. Added rejection, exact charge/time, target-only damage, attacker identity and JSON persistence checks.

The first isolated8898 run exposed generic attack staging waiting indefinitely behind fire crews. A dedicated incendiary cue removes that unrelated melee actor; 3D fire is owned by the existing persistent window-fire renderer, with no generic doorway particle ring. Added a small fire presentation to the legacy isometric compatibility view after its coverage guard caught the new cue. The broad run also caught two stale source guards for removed CSS slot drums; updated them to inspect the existing3D cabinet's authoritative strip/texture travel path. These source guards supplement, not replace, the existing model and browser checks.

Final CUA8899 copied campaign: cash180→140, attention3→21,19:20→19:35, club condition100→82, revision36→37. Scene ended without Skip and automatically opened ARSON AT THE MONARCH. After folding the paper, screenshot showed window flame/smoke, engine and brigade hose streams. Diagnostic effects empty, saved fire48 particles/four vents. Reopening retained the same fire ID/minute/revision; entering the club showed82% condition and disabled another attack as already burning. Only copied QA saves were changed; main campaign untouched.

Core/store broad run passed115.365s/.146s; cmd initially failed the newly missing compatibility cue. After fixes cmd passed.812s (.runtime/incendiary-interface-tests.log); final targeted core checks passed.154s (.runtime/incendiary-targeted-verified.log). All246 frontend tests passed11.473s (.runtime/incendiary-frontend-tests.log); TypeScript/Vite build passed2.61s (.runtime/incendiary-build.log), retaining the known large-bundle warning. Caffeinate87182 remained live. The preceding daily-accounts verification confirmed the existing immediately visible15-dollar breakdown and required no edit.

Outstanding: authored bottle/throw/flee choreography, debug incendiary staging, risk/balance and interior fire/evacuation behavior, plus the broader full-production objective. Current fire response is not acceptance of a complete molotov scene or whole-game60FPS.


### 2026-09-14 — authored incendiary bottle prop

Built incendiary-bottle.glb in Blender, with a48-sided revolved heel/body/shoulder/neck profile, glossy olive material, curved paper label with packed grain/border texture, physical Bell Reserve1953 lettering and a continuous folded cloth strip with charred tip. The origin is the neck grip; named bottle-grip and bottle-flame empties support the upcoming held/flight presentation. Added targeted export and full-export inclusion.

CUA8897 WebGL close-up initially revealed block-like stacked cloth folds; replaced those with a continuous thickened curved sheet and inspected the corrected front/side view. Bottle is approximately8cm across, base25cm below grip; flame marker above cloth. Actual exported-GLB test passes: grip position, dimensions, outward label raycast in front of glass and under10000 triangles (.runtime viewer is an isolated asset study). No campaign or API changes. Caffeinate87182 confirmed alive. This supplies the required prop but is not yet attached to a live attacker or integrated into throwing/flight/escape choreography; those remain next, along with the full objective.


### 2026-09-14 — held bottle, release and escape motion study

Added CityIncendiary presentation cast: short approach, right-arm windup using the existing two-segment IK, release at3.1s, arced flight to a caller-provided target at3.9s and retreat through6.8s. The held prop follows the actual elbow/hand transform; release caches that position and quaternion so there is no visual snap. The cast is expressed in local metres for future city placement/path reservations. It does not change simulation outcomes.

Browser study on isolated Vite8897 first exposed the body hanging inside the forearm. Orienting the bottle along the forearm exposed the body but crowded its protruding cloth; the final perpendicular grip keeps the bottle body outside the sleeve and makes the label readable in the close-up windup screenshot. Playback reached6.80s normally. This remains a basic rigid-hand pose, not finished finger articulation.

Actual person/woman GLB tests sample held poses under a translated/rotated scene root, check attachment error below1e-6m and bottle-body clearance from the forearm segment, continuous release position/quaternion, flight height, exact target contact, disappearance at impact and escape endpoint. All248 frontend tests pass10.942s (.runtime/incendiary-motion-tests-final.log), TypeScript check passes, diff check clean. No campaign state was changed; existing unrelated core files remain unstaged.

Not yet wired into City3D: traffic reservation for the complete cast path, target selection on the actual facade, flame/impact audio and timed persistent fire/brigade reveal, debug entry and committed-scene browser acceptance. The isolated study does not establish whole-city collision safety or production animation quality. Full goal remains active.


### 2026-09-14 — city incendiary cast integration

Loaded the authored bottle in City3D and wired incendiary cues to CityIncendiary using the actual public attacker rig/wardrobe. Scene slots reserve a6.2m by1.8m forecourt sweep; the player's ordinary street actor is excluded while its staged counterpart owns that reservation. A clear fire-window marker nearest the slot sets the flight endpoint, with entrance fallback. Camera includes the swept path and impact point. Added burning-cloth particles/light, impact glass sound and modest camera impulse. Cast duration6.8s gates normal scene completion/news.

Window fire/rubble are withheld until impact; saved brigade/hoses are deferred until the cast exits to prevent the earlier occupied-frontage stall. Those filters only affect presentation, preserving saved timestamps/condition. Incendiary is now selectable in debug previews with a private synthetic unarmed attacker/fire, no casualties or campaign mutations.

CUA8894 preview at The Monarch: staged true at0.236s, held bottle, window fire empty; after completion effects empty and48 fire particles/four vents. Revision36/minute1160 unchanged. Final build CUA8899 recorded attack replay: actual Alex Varga stages, initial fire/suppression empty; natural completion opens ARSON AT THE MONARCH. After folding newspaper, effects empty and saved fire plus two hose streams restored, revision37/minute1175/cash140 unchanged. A longer read-only browser sample exceeded the CUA evaluation deadline; it was not treated as a failed scene or used as mid-flight evidence. No new campaign commands/main-save writes.

Actual person/woman GLB sweep checks sample181 frames per rig against the reserved footprint. Existing attachment/release/flight tests remain green. All249 tests passed11.393s (.runtime/incendiary-city-tests.log), then the additional private-fire preview test passed in the4-test preview file run. Final TypeScript/Vite build passed2.34s with the known large-bundle warning. Diff check clean.

Remaining acceptance: actual facade/projectile clearance at every address, close-up integrated throwing/flame fidelity, glass/ignition impact detail, smooth escape handoff to the real player and fire evacuation/interior behavior. This is a first integrated scene, not completion of the full production goal or broad60FPS acceptance.


### 2026-09-14 — facade-aware incendiary flight and Mariner fire anchors

Actual exported-building checks exposed the original fixed arc crossing The Paper Moon's canvas awning. Raising the arc instead crossed facade trim. Added once-per-staging flight selection across candidate window targets and several lofts, tracing48 segments with centre and six27cm-offset clearance samples. The endpoint stands30cm outside the window so the bottle body can contact before its neck reaches the pane. Paper Moon selects the flatter0.4m loft in the browser. This is sampled clearance, not continuous full-mesh collision detection. Existing fallback for addresses without a validated window path remains unproven and needs broader work.

The wider check also found Mariner had no authored fire-window markers. Added them in front of its sash panes, above the cross rails, and a targeted --only=mariner export. GLB grows699412→700472 bytes. Its existing geometry is retained.

All three forecourt slots on12 exported models now pass the independent80-segment centre sweep after selecting an envelope-checked path, at the actual18cm building elevation. Models: tavern, monarch, tenement, shop, civic, casino, warehouse, bluehour, goldenlily, papermoon, mariner and mercer-court. CUA8894 Paper Moon diagnostic loft.4/staged true; Mariner loft.8/staged true and after completion48 particles/four vents. Mariner screenshot verified flames at sash windows. Preview revision36/minute1160 remained unchanged and preview was stopped. No campaign/main-save writes; caffeinate87182 alive.

All251 frontend tests pass11.288s (.runtime/incendiary-clearance-tests.log), final TypeScript/Vite build passes2.17s (.runtime/incendiary-clearance-build-final.log) with the known bundle warning; diff check clean. Still outstanding: industrial/remaining special facades, robust no-path handling, continuous bottle/obstacle clearance, detailed impact/escape handoff and scene-start performance cost, plus the unchanged full production objective.


### 2026-09-14 — industrial fire targets and cargo clerestories

Added authored fire-window markers to filling-station office glazing and workshop bay windows, covering filling/garage/dealer assets. Docks and haulage cargo sheds now have three framed high clerestory panes above the loading gate with matching markers. Targeted --only=industrial-fire exports these five assets; full export includes the same changes. Cargo frontage grows2cm for the new panes. Dealer regeneration also incorporates the current shared Ford mesh, extending its existing front bound7.5mm; geometry remains inside its lot.

Extended flight checks to the industrial models at their actual180-degree city rotation: all three candidate slots on17 models (51 sampled paths) have an accepted path and no earlier centre-line masonry contact. Additional marker tests verify three clear street-facing targets per industrial model with actual glazing immediately behind each anchor. These remain sampled flight-envelope/centre tests, not full continuous geometry collision proof.

CUA8894 Kessler's Filling Station staged a target under its canopy (loft.8); afterward effects empty and36 particles/three window vents. Pier14 screenshot shows flame/smoke at the three new high cargo windows after playback, again36 particles and no active effect. Both previews retained revision36/minute1160/cash180 and were stopped normally. No campaign commands or main-save changes. Caffeinate87182 confirmed alive.

All251 frontend tests passed11.362s (.runtime/industrial-fire-tests.log); the subsequently added marker test passed together with the full51-path test (2 tests,3.601s). TypeScript/Vite build passed2.60s (.runtime/industrial-fire-build.log), known bundle warning unchanged, diff check clean. Remaining special residential/undertaker facades, no-path fallback, complete physical collision/impact detail, evacuation and the broader production goal remain outstanding.


### 2026-09-14 — special facade fire coverage

Added authored upper-window fire markers to the undertaker and villa, offset above sash rails, plus --only=special-fire export. Expanded sampled flight checks to all three slots on these models:57 paths across19 models now pass. Additional tests verify four undertaker and three villa markers have actual upper-window glazing behind them. No command/schema changes.

Regeneration corrected villa manifest height10.7513→9.45. Investigated the apparent geometry change by loading the old committed GLB and new export: both actual tops are9.4500006, roof max8.8870216, identical roof vertex count3432. The old manifest was stale; roof/chimney geometry was not shortened.

CUA8894 initially selected Ackerman & Son (pawn/shop), then checked core/locations.json and selected the actual undertaker Thorne & Sons. Its preview stages at chapel and finishes with48 particles/four upper-window vents. Cypress House finishes with36 particles/three vents. Screenshots verify both window fires. Previews stopped normally; revision36/minute1160/cash180 unchanged. No main or QA campaign commands. Caffeinate87182 live.

Full252-test frontend suite passed11.317s (.runtime/special-fire-tests.log); subsequently added special-marker test passed with the complete57-path and industrial-marker tests (3 tests,3.846s). TypeScript/Vite build passed2.66s (.runtime/special-fire-build.log), known bundle warning unchanged; diff check clean. No-path fallback, full continuous collision safety, scene-start cost, detailed impact/escape handoff and broader production acceptance remain open.


### 2026-09-14 — cached incendiary flight selection

Measured identical repeated flight searches before changing them: Node Monarch64–69ms, Paper Moon55–57ms, undertaker15–16ms. Added a per-building WeakMap cache capped at12 start/window combinations, returning cloned result vectors so callers cannot corrupt cached targets. The geometry signature includes mesh world matrices, geometry identity, attribute/index buffer identities and versions (including interleaved backing buffers), and material sides. Building transform/geometry changes clear its paths. Both clear and blocked results are cached.

Tests count actual raycasts to verify identical calls reuse results, returned-vector mutation cannot alter them, moved walls invalidate clear paths, edited geometry invalidates blocked paths, and replacing buffers with the same version invalidates stale geometry. Existing57-path collision checks remain green. All254 frontend tests pass11.576s (.runtime/incendiary-cache-tests-final.log); final TypeScript/Vite build passes2.49s (.runtime/incendiary-cache-build-final.log), known bundle warning unchanged, diff check clean.

CUA8897 actual WebGL model-loader benchmark (private page, no campaign): Monarch first58.7ms, subsequent0–.1ms; Paper Moon first45.2ms, subsequent0–.1ms; undertaker first12.9ms, subsequent0ms at browser timer resolution. These are planning timings, not whole-frame or GPU measurements. Node cached repeats were.03–.10ms. First-use stalls remain unresolved and need incremental/prepared planning; static-building cache is not intended for skinned/deforming geometry. Main and QA campaign state untouched. Caffeinate87182 confirmed alive. Full objective remains active, including no-path fallback and all other unfinished gameplay/visual acceptance.


### 2026-09-14 — finite-segment broad phase for cold flight planning

Flight selection now builds world-space mesh bounds once per cold search and rejects meshes whose first box contact lies beyond the short ray segment. Three's ordinary mesh box check uses the infinite ray, so previously many short segments still reached triangle tests for surfaces farther away. Origin-inside-box segments remain eligible. Exact mesh raycasts still decide collisions for retained candidates; lofts, clearance sample offsets and segment counts are unchanged.

Initial cache test expected a distant mesh to receive raycasts; the new broad phase correctly made that count zero. Updated the cache fixture with a nearby ring whose box overlaps the path but whose opening is clear, retaining meaningful raycast reuse/invalidation checks. Added a transformed-building test that requires zero triangle calls for geometry beyond all segments and a blocked result after moving that same wall into the flight.

CUA8897 first-use planning benchmark: Monarch16.9ms (previous58.7), Paper Moon9.8ms (previous45.2), undertaker5.1ms (previous12.9); cached repeats0–.1ms at timer resolution. Node cold measurements20.8/14.2/2.1ms. These are planning timings only: Monarch still consumes roughly a60FPS frame by itself, so first-use frame-budget work remains open.

All254 existing frontend tests pass11.806s (.runtime/incendiary-bounds-tests.log); subsequently added broad-phase test passes with all facade/cache tests (5 tests,1.225s), including57 sampled authored flight paths. TypeScript/Vite build passes2.37s (.runtime/incendiary-bounds-build.log), known bundle warning unchanged; diff check clean. No campaign writes or API changes. Caffeinate87182 alive. Full goal remains active, including robust no-path handling, continuous physical clearance, choreography fidelity and remaining game systems.

### 2026-09-14 — incendiary actor corridor geometry

Scene selection now checks the complete approach/escape corridor against the target building's authored triangles before accepting an otherwise unoccupied incendiary slot. Mesh bounds prune candidates, but actual triangle intersections decide whether a facade detail obstructs the corridor. All57 existing corridors across19 models pass. A new fence fixture blocks the first slot and verifies selection of another; moving the fence away restores clearance.

All255 existing frontend tests pass11.101s (.runtime/incendiary-corridor-tests.log); subsequently added fence-selection test passes with all facade/cache checks (6 tests,1.218s, .runtime/incendiary-corridor-targeted.log). TypeScript/Vite build passes2.15s (.runtime/incendiary-corridor-build.log), known bundle warning unchanged. CUA8894 Monarch preview staged normally with loft.8; after stopping, effects empty and revision36/minute1160 unchanged. No campaign commands or main-save changes. Caffeinate87182 confirmed alive during verification.

This is a surface-intersection guard for known exterior slots, not generalized solid-volume collision detection. Neighboring props, moving obstacles, fully blocked staging/no-flight behavior, escape handoff and the broader production objective remain unfinished.

### 2026-09-14 — accept staging only with a checked throw

Removed the unchecked first-window/entrance fallback from incendiary staging. Each otherwise available slot now needs both a clear actor corridor and an accepted bottle flight before it is reserved. The selected target and loft are reused for playback. A blocked flight tries the next slot, and missing windows cannot invent a target. The staging helper keeps the cast's local release and supplied window ordering unchanged.

Added a geometric fixture whose wall leaves the walking corridor clear but blocks every loft from the first slot: selection advances to the second. Enlarging the wall blocks all slots and returns no placement; missing window targets also return no flight. Existing57 authored flight/corridor checks remain green. Full257 frontend tests pass11.250s (.runtime/incendiary-staging-flight-tests.log); TypeScript/Vite build passes2.31s (.runtime/incendiary-staging-flight-build.log), known bundle warning unchanged. Diff check clean.

CUA8894 Monarch preview staged with loft.8, completed with effects empty and48 fire particles/four vents, then stopped normally. Revision36/minute1160 stayed unchanged; no campaign commands or main-save changes. Caffeinate87182 alive. Previous goal turn was concrete progress (a39de50); this turn removes an unsafe fallback, not completion of the full goal. Fully blocked scenes still wait and need a recovery strategy that preserves choreography. General obstacle coverage, continuous collision fidelity, escape handoff and the complete production objective remain open.

### 2026-09-14 — smooth incendiary approach and escape

Replaced constant-speed/clamped travel with integrated smooth acceleration and braking, preserving the two-metre approach, four-metre escape, release/impact times and reserved footprint. Limb swing and vertical bounce now fade with speed instead of snapping to neutral at the stopping frame. The throwing arm blends from its lowered IK pose into the relaxed pose during3.75–4s, removing its previous abrupt reset before the escape turn.

Full257 existing frontend tests pass11.143s (.runtime/incendiary-stride-tests.log), including all authored flight/corridor checks. Subsequently added continuity test passes with all three motion tests (124ms): monotonic travel/exact endpoints, speed continuity at acceleration boundaries, and root/joint continuity on person and woman at approach, arm-settle and escape boundaries. Build passes2.16s (.runtime/incendiary-stride-build.log), known bundle warning unchanged; diff check clean.

CUA8894 observed the new approach increasing from x107.00195 at.2423s to107.15719 at.4237s, then107.40673 at.6042s, with the bottle still held. Scene completed with effects empty and48 particles/four window vents; stopped normally. Revision36/minute1160 unchanged; no campaign commands/main-save writes. Caffeinate87182 alive. Previous turn a126863 was concrete progress. This improves staging motion, but does not solve the staged actor's handoff to the ordinary map character, planted-foot animation, fully blocked recovery or the remaining full production objective.

### 2026-09-14 — authored bottle glass at incendiary impact

Authored an irregular curved, thin olive-glass chip in Blender (bottle-shard.glb,1812 bytes), with --only=bottle-shard and full-export support. Twelve instances scatter from the actual checked impact target at3.9s, rotate and fall outward, settle above the pavement and fade by2.4s after contact. They remain hidden before staging/impact. Short per-frame centre ray segments against target-building geometry suppress fragments encountering the facade. Existing effect disposal releases their cloned geometry/material/instance resources.

Added tests for asset dimensions/thickness, impact origin/timing, outward travel, floor bounds and fade over three impact heights. Initial new assertion distinguished negative zero from zero; corrected the coordinate assertion to compare absolute zero. Final259 frontend tests pass11.910s (.runtime/incendiary-glass-tests-final.log). Final TypeScript/Vite build passes2.17s (.runtime/incendiary-glass-build-final.log), existing chunk warning unchanged. Diff check clean.

CUA8894 pre-impact2.486s: zero shown shards; screenshot/diagnostics at4.256s:12 shown, none blocked. Scene completed with effects empty, revision36/minute1160 unchanged, then stopped normally. No campaign commands/main-save writes; caffeinate87182 alive. These small fragments add bottle impact detail, not shattered-window geometry or persistent glass aftermath. Centre-segment suppression is not full swept-fragment collision proof and does not cover neighboring/dynamic obstacles. Escape handoff, blocked-scene recovery and the full production objective remain open. Previous goal turn4897f33 was concrete progress.

### 2026-09-14 — curved glass clearance across slow frames

Replaced a single endpoint chord with short trajectory segments (at most1/120s) over the elapsed interval. Seven centre/axis-offset rays sample7cm chip clearance; finite mesh-box checks prune triangle calls. This closes a demonstrated awning miss when a frame spans a curved portion of the fall. Geometry world bounds update for moved obstacles. Checks remain sampled surface collisions, not a complete swept-volume solver.

Added an awning fixture missed by the old chord but detected by both a single.8s interval and normal60Hz intervals, plus moved-obstacle/pre-impact checks. A separate thin obstacle verifies chip-edge clearance where the centre path remains clear. Full261 frontend tests pass11.332s (.runtime/incendiary-shard-clearance-tests.log); final six motion tests also pass128ms after tightening the edge fixture's depth. Build passes2.15s (.runtime/incendiary-shard-clearance-build.log), known chunk warning unchanged; diff check clean.

CUA8894 Monarch at5.0566s:12 blocked/zero shown fragments after reaching its awning. Local metrics reported145FPS,7.6ms p95,109 draw calls and533626 triangles; this is a single preview observation, not whole-game performance acceptance. Scene completed normally; revision36/minute1160 unchanged, then preview stopped. No campaign/main-save mutations. Caffeinate87182 alive. Previous46f91ac was concrete progress; full goal remains active, including actual glass contact response/persistence, escape handoff, blocked-scene recovery and the many outstanding gameplay/visual requirements.

### 2026-09-14 — authoritative explosion outcome and planter identity

Started the explosion-preamble work by preserving facts the renderer previously lacked. Explosion cues carry optional detonation=planted/premature; player planting captures the existing attacker identity with weapon0. Faction detonations carry planted but no invented individual planter. Legacy cues stay unspecified. Damage, odds, charges and command duration are unchanged. API.md documents the additive field and inference precedence.

Internal blast selection now respects explicit outcome before legacy fire/preview inference, so a premature accident cannot borrow a same-minute building fire and become an internal planted blast. Debug includes Explosion · premature with the appropriate outcome and no invented persistent fire. Normal Explosion records planted. This does not yet implement a character exit or accident-injury animation.

New Go command tests exercise both planting outcomes and verify cue identity/outcome survive full-world JSON save/reload; faction and legacy checks prevent invented identities/preambles. Targeted tests pass.271s. Broader go test ./core ./store ./cmd/... passes: core109.136s, store.164s, blackledger.916s, simulate5.339s and other command packages (.runtime/demolition-cue-go.log). Final263 frontend tests pass14.418s (.runtime/demolition-cue-frontend-final.log); build2.83s (.runtime/demolition-cue-build-final.log), known bundle warning unchanged; diff check clean.

CUA8894 premature preview produced an exterior origin y.25 and zero fire scenes, completed with effects empty/zero persistent fire, and stopped normally. Revision36/minute1160 unchanged; no QA campaign or main-save mutations. Caffeinate87182 alive. Previous fd84d56 was concrete progress. Planter exit choreography, casualty synchronization with the preamble, accident injury/death staging and the entire remaining production objective are still open.

### 2026-09-14 — planter doorway animation review

Added CityPlanter, a doorway-local animation for the recorded planter: waits behind the door's entire swing, walks4.6m through the vestibule, closes the door only after clearing it, turns and walks2m along the pavement. Acceleration/braking and distance-based gait reuse the tested stride profile. The blast-ready beat is6.2s. The caller owns placement and applying the returned door openness. This module is NOT yet connected to live City3D explosion playback.

Committed tools/planter-preview.html makes the animation review reproducible with the existing Vite server. CUA8897 at2.3s showed the character crossing the actual Monarch doorway (local z-.196, door fully open);6.2s showed the character clear at(-2,.02,-2.35), door closed and blast-ready true. Corrected the review floor placement to cover the actual building entrance before final screenshot. No campaign commands/main-save changes; caffeinate87182 alive.

Geometry tests sample187 frames on both person and woman at Monarch and tavern: raised actor bounds do not intersect authored door/masonry triangles; beat checks verify clearance before detonation. This is sampled conservative body bounds, excluding the bottom12cm, not full foot-placement or all-building acceptance. All264 frontend tests pass11.319s (.runtime/planter-motion-tests.log); build2.21s (.runtime/planter-motion-build.log), known chunk warning unchanged; diff check clean.

Next required integration includes traffic reservation, preserving the actual planter's identity, all eligible doorway geometry, blast/fire/glazing/audio/casualty timing and newspaper completion as one scene. Do not claim the preamble is in gameplay until that work is done and tested. The complete production objective remains active; previous534e671 was concrete progress.

### 2026-09-14 — planter doorway coverage and swept traffic reservation

Expanded the actor/triangle path check to all11 current models with animated entrance markers: monarch, tavern, tenement, mercer-court, casino, civic, shop, warehouse, bluehour, goldenlily, papermoon. Both rigs pass187 samples each. The remaining mariner, villa, undertaker and five industrial models lack these animated doorway markers and are not covered.

Added a planter traffic footprint3.4m wide by5.8m deep, centered on the full vestibule/turn/pavement sweep. planterReservation supports translated/rotated placement; availableSceneSlot requires an explicit doorway and rejects an occupied exit. Vertex-level footprint tests cover both rigs across four headings. These helpers are ready for integration but are not yet consumed by live explosion staging.

The committed browser review now has an11-model selector. CUA8897 Paper Moon threshold screenshot shows the exit through its doorway. Found the standalone loader displayed both glazing variants; review now initializes intact glazing through the same buildingGlazing helper used by the game. Selected casino through the UI and verified closed door/clear actor at the6.2s blast-ready beat. No campaign/main-save changes; caffeinate87182 alive.

All265 frontend tests pass11.509s (.runtime/planter-reservation-tests.log); build2.31s (.runtime/planter-reservation-build.log), known chunk warning unchanged; diff check clean. Full scene integration, missing doorway assets, dynamic traffic behavior, aftermath/casualty/audio timing and the complete production objective remain outstanding. Previous969866d was concrete progress; this turn adds verified geometry coverage and the required reservation, not a claim of an in-game preamble.

### 2026-09-14 — planter preamble integrated with committed explosions

City3D now casts the recorded attacker for planted explosions at an authored animated doorway. It reserves the complete exit before casualties, suppresses the ordinary duplicate actor, frames the building/exit and uses a future blast onset after the6.2s preamble. Door movement follows the planter. Blast particles/light/audio/debris and window break/fire wait for onset; same-address/minute casualties remain hidden until the planted blast begins, including when its slot is still unavailable. Fire-brigade staging waits until the scene completes. Existing all-effects completion keeps the newspaper afterward. Normal debug Explosion uses the current player; added Explosion · casualty pairs a synthetic victim with that planter for combined review.

Created a NEW isolated save using qa-fixture city3d-blast: .runtime/planter-live-20260914.sqlite3. Built .runtime/planter-live-server and started it on8900 (session50005). CUA issued the actual Put the charge under The Monarch action once. Revision0→1, minute600→720; cash12000→12027, respect60→68, attention0→22, health100. Result killed Rosa Erdos and committed a planted explosion with Alex Varga as attacker.

At preamble.3194s: planter staged, blastSeconds-5.8806, blast/glass audio counts0, casualty stagedfalse/fall0, no Herald. At blast+2.0704s: blast/glass counts1, casualty fallen and visible,48 fire particles/four vents, planter clear at(110,.2,72.07). Screenshot showed the active scene rather than a newspaper. After the full scene, Herald automatically displayed ROSA ERDOS FOUND DEAD with the explosion cause. Closed it normally: effects empty, saved fire retained, revision1/minute720 unchanged. Separate8894 preview also showed zero fire/audio during its exit and was stopped. No main-campaign writes; caffeinate87182 alive.

All266 frontend tests pass11.953s (.runtime/planter-integration-tests.log), including11-model geometry/footprint coverage and new combined-preview identity/timing linkage checks. Final TypeScript/Vite build passes2.73s (.runtime/planter-integration-build-verified.log), known bundle warning unchanged; Go server build succeeds. No Go implementation changes this turn. Diff check clean.

This now works in gameplay for recorded planted attackers at the11 doorway models. Remaining: eight models without doorway support, unnamed faction planter choreography, premature injury/death scenes, crowded-exit recovery, exact actor handoff, more cinematic close framing and full production acceptance. The full goal stays active; previous c3bd3c5 was concrete progress.

### 2026-09-14 — close planter framing and cancellable blast pullback

Planter staging now fits the doorway/exit first, then smoothly pulls back during4.6–6s for the complete building/blast envelope before6.2s detonation. ScenePullback interpolates target/position and visible world span, completes once and cannot reclaim the view after cancellation. Pointer camera input, wheel, keyboard pan/rotate/zoom, address/whole-city focus and player-follow cancel it. The existing cutaway now follows the actual moving planter/incendiary actor rather than the stationary scene root.

CUA8894 observed exit zoom20.5905 at2.347s and final wide zoom4.96064. A second preview manually switched to Whole city; at7.7219s (blast+1.5219s) zoom remained1 at the city centre. Close-up revealed the old scene-root cutaway lag; after updating its target, screenshot showed the canopy cutaway tracking the character under it. Existing stippled cutaway appearance is still a visual-quality limitation. Final preview completed with effects empty/revision36/minute1160 unchanged; both previews stopped. No campaign/main-save commands; caffeinate87182 alive.

Added projected-envelope tests at narrow/wide aspect ratios, monotonic outward zoom, before-start inactivity and permanent manual cancellation/completion. All267 frontend tests pass11.637s (.runtime/planter-camera-tests.log). Final build after cutaway targeting passes2.38s (.runtime/planter-camera-build-final.log), known chunk warning unchanged; diff check clean. Full production goal remains open, including missing doors, crowded staging, actor handoffs, injury scenarios and remaining gameplay/art requirements. Previous d65b40a was concrete progress.

### 2026-09-14 — pullback framing survives viewport resizing

ScenePullback previously retained zoom fitted to its initial viewport. It now detects frustum-size changes, refits the final blast envelope and adjusts the initial visible span for the resized view. Completion and cancellation are distinct: an automatic completed view can refit on resize, while a manually cancelled one cannot resume. Before-start resizing also preserves the doorway span. Target bounds are copied rather than aliased.

Added tests for narrowing before, during and after the move, checking all projected blast corners and ensuring a resize cannot revive cancellation. All268 frontend tests pass11.764s (.runtime/pullback-resize-tests.log); build2.28s (.runtime/pullback-resize-build.log), known bundle warning unchanged; diff check clean.

Committed tools/pullback-preview.html provides an actual WebGL/model review with controlled canvas resizing. CUA8897 changed to Halfway, narrowed1234→431px, then Blast view: zoom2.0061 and all blast-envelope corners fit. Widening after completion refitted to zoom4.2955 with the envelope still fitting. Screenshot verified the model and envelope inside the narrow canvas. This tests the shared camera helper, not every responsive game HUD layout. No campaign/main-save writes; caffeinate87182 alive. Previous8482d57 was concrete progress. Full production objective remains active, including unfinished doorway assets, crowded exits, actor handoffs and remaining systems.

### 2026-09-14 — Skip actually cancels a scene behind its newspaper

Browser replay exposed a lifecycle bug: Theatre Skip marked the cue finished to reveal its newspaper, but kept passing that cue as active to City3D. The explosion/actors/audio could therefore continue behind the article. City3D now receives an active cue only while scenePending is true. The article retains its separate playing reference, so cancelling presentation does not dismiss the newspaper or alter the recorded event.

CUA8900 replayed the existing isolated planted explosion. Before the fix, effects were still running after Skip (planter16.096s/blast+9.896s). After the fix, a controlled Skip at preamble1.3472s changed the open door(-pi/2) to0 and effects to empty, while the article appeared; blast audio had still been0 immediately before cancellation. Revision1/minute720 unchanged. Closed Herald and replayed again: preamble restarted at.2374s/audio0, later finished naturally with effects empty, door0 and Herald visible. No new gameplay command or main-save write.

All268 frontend tests pass12.762s (.runtime/scene-skip-tests.log); TypeScript/Vite build2.36s (.runtime/scene-skip-build.log), known chunk warning unchanged; diff check clean. Caffeinate87182 alive. Previous9a4a5e0 was concrete progress. Full production goal remains active, including missing doorway models, crowded exits, handoffs and the other unfinished gameplay/visual requirements.

### 2026-09-14 — Give the Mariner a working planter exit

Replaced the Mariner's solid entrance masonry/foundation with two wings, a 3m recessed vestibule, lintel and paving. Its existing painted door, panels, glazing and brass pull now share an authored hinge; the threshold marker enables the existing city planter sequence without a renderer special case. Retained the porch, windows, roof and overall footprint. The local Blender export updates mariner.glb and its measured manifest (720624 bytes, paving minimum height -.08). No external assets or gameplay/interface changes.

Extended the exact-triangle building-clearance test to the Mariner: both person rigs, 187 sampled frames each, with the door animated. All 12 doorway models pass, as do the translated/rotated reservation tests. This checks actor bounds above the bottom12cm against building triangles; it is not a full swept-volume/foot-contact proof. CUA8897 tab227 visually checked threshold2.3s and blast-ready6.2s: the actor clears the porch and the door closes. CUA8900 tab228 tested Explosion at The Mariner using the private debug snapshot: preamble1.4441s showed an open door and zero blast audio; blast+2.5146s showed the actor outside, door closed, one audio blast and window fire. Natural completion and Stop both left effects empty/door0. Revision1/minute720 stayed unchanged; no gameplay command or main-save write.

All268 frontend tests pass in11.459s (.runtime/mariner-door-tests.log); TypeScript/Vite build2.30s (.runtime/mariner-door-build.log), existing chunk warning unchanged. Caffeinate87182 confirmed alive. The preceding Daily Accounts response only revalidated existing code (no new progress); this turn advances the actual building-exit requirement. Remaining: seven missing doorway models, crowded staging recovery, actor handoff, higher-quality assets/animation, Mariner breakable glazing, and the broader production goal. Full goal remains active.

### 2026-09-14 — Author the funeral parlor's offset exit

The undertaker model now has a real3m vestibule at its existing offset front entrance, independent masonry/foundation wings, an animated paneled oak door with attached brass pull, and a threshold marker. The raised sill intersected the walking cast in the first geometry check at2.1s; made the stone threshold flush and extended the leaf to its floor. Preserved the long display window, porch-free frontage, roof, hearse and yard. Targeted local Blender --only=undertaker export updates its GLB/manifest (970704 bytes); no gameplay change. Added the model to the shared exit review and collision suite. Lowered only the review's artificial floor5mm to prevent coplanar flicker with authored paving.

Both rigs now pass187 sampled frames against all13 doorway models. An experimental full-foot bounds check also found the existing shared gait dipping a sole to y=-.0048 at.633s; the established test still excludes the bottom12cm. This is explicit remaining gait work, not proof of complete foot contact. Frontend268 tests pass12.283s (.runtime/undertaker-door-tests.log); build2.54s (.runtime/undertaker-door-build.log), known chunk warning unchanged; diff check clean.

CUA8897 tab229 inspected the actual offset doorway at2.3s. CUA8900 tab230 first selected Ackerman by mistake (existing pawn/shop model); stopped that private preview and selected the authoritative chapel address Thorne & Sons. Its preamble1.4444s showed the chapel door open, actor x180.5 and audio0; blast+5.9039s showed the actor outside at(178.5,.2,102.1), door0, audio1 and window fire. Revision1/minute720 stayed unchanged. No gameplay command or main-save write. Caffeinate87182 remains alive. Previous05b631d was concrete progress. Six doorway models, accurate planted-foot gait, crowded exits/handoffs, and the full production objective remain outstanding.
Natural completion was also observed in tab230 with effects empty, chapel door0 and revision1/minute720 unchanged; then Stop restored the private preview.

### 2026-09-14 — Ground the planter's animated soles

CityPlanter now samples its two authored shoe meshes after joint posing and adjusts the actor vertically so the lowest sole remains5mm above the local pavement. Cached mesh references/matrices/vector avoid per-frame allocations and whole-character bounds. This removes the observed sole dip and airborne settling caused by the former fixed bob. Exit distances, headings, doorway timing and detonation timing remain unchanged. This is vertical support on a flat local surface, not horizontal planted-foot IK or stair traversal; those remain required animation work.

New geometry regression samples both actual rigs at120Hz through6.8s in translated and rotated roots, checks both shoes, sole contact and bounded vertical changes. The undertaker's full character bounds now also pass the building-triangle collision test without its old bottom12cm exclusion; other models retain their prior conservative clearance test. All269 frontend tests pass11.562s (.runtime/planter-foot-tests.log); build2.32s (.runtime/planter-foot-build.log), known chunk warning unchanged. After strengthening the undertaker test, its full three-test suite passes2.275s; diff check clean. CUA8897 tab231 reviewed threshold2.3s, running2.6509s and natural6.8s completion/closed door, then inspected6.2s settled pose. No campaign operation. Caffeinate87182 confirmed alive. Prior0fd858d was concrete progress; the full objective, six missing doorway models, horizontal foot sliding, other actor gait paths, staging/handoffs and all remaining production requirements stay active.

### 2026-09-14 — Working filling-station office entrance

Added a glazed cream-enamel office door beside the station's existing window band, with an inward hinge, independent rails/panel/glazing/pull, jambs, recessed passage, rear/side closure and paving. The right hinge permits the established opening rotation with this model's opposite authored frontage. Test and browser review apply the actual city rotation(pi); city scene handling needs no special case. Re-exported local Blender industrial assets; only filling.glb/its manifest changed (315084 bytes). Existing pumps, canopy and office windows remain.

The sampled full exit clears office masonry, door, pumps and columns for both rigs. All14 doorway models remain in the test; the established bottom12cm exclusion still applies to filling. Frontend269 tests pass11.468s (.runtime/filling-door-tests.log); build2.32s (.runtime/filling-door-build.log), existing chunk warning unchanged. CUA8897 tab232 inspected threshold2.3s; the review camera can obscure the actor behind a pump, which is occlusion rather than tested geometry overlap. CUA8900 tab233 at Kessler's Filling Station: preamble1.4444s had the filling door open, actor(139,.2053,18.4534), audio0; blast+8.2229s had door0, actor outside(137,.18,15.05), audio1. Revision1/minute720 unchanged, no gameplay commands/main-save writes. The canopy hides much of the window fire at the default wide camera angle; presentation visibility remains unfinished. Caffeinate87182 confirmed alive. Prior a53a6f6 was concrete progress. Five models still lack working exits, and all broader production requirements remain active.
Tab233 natural completion also showed effects empty/door0 with revision1/minute720 unchanged; then stopped the private preview.

### 2026-09-14 — Cargo-shed pedestrian exits

Dock and haulage cargo sheds now have a hinged steel pedestrian wicket within the existing loading gate. Independent masonry wings, passage lintel/rear/paving and split gate framing leave a real exit; braces and pull move with the leaf. Retained the clerestory, zinc roof, yard crates and dock derrick. Local Blender industrial export changed only docks.glb(457596 bytes), haulage.glb(449192 bytes) and their measured manifest. The shared city animation uses their normal pi building rotation, matching the expanded geometry tests/review fixture.

Both rigs pass sampled door/building/yard geometry checks across16 doorway models. Established bottom12cm exclusion remains for these two models; no claim of swept-volume collision or complete planted-foot gait. All269 frontend tests pass11.450s (.runtime/cargo-door-tests.log); build2.29s (.runtime/cargo-door-build.log), known chunk warning unchanged. CUA8897 tab234 inspected dock threshold2.3s. CUA8900 tab235 saw both Pier14 and Kerrigan Haulage preambles at1.4444s with their correct doors open/audio0, followed by natural completion with doors0/effects empty. A controlled haulage replay captured7.4709s/blast+1.2709: actor outside(81,.18,106.07), closed door, audio1. Stop cleared all effects/open doors. Revision1/minute720 unchanged throughout; no gameplay or main-save writes. Caffeinate87182 confirmed alive; diff check clean. Previous593323d was concrete progress. Three models still lack exits (garage, dealer, villa); asset fidelity, gait, staging/handoffs and the broader full objective remain unfinished and active.

### 2026-09-14 — Garage and dealership pedestrian exits

Authored an inward-opening steel wicket in the central service bay for garage/dealer models. Split the bay panels/seams and masonry around a real passage, preserving upper glazing, other bays, zinc roof/vents and the dealership's Ford display car. Added the entrance marker, jambs, reinforcing ribs and attached handle. Local Blender export changed only garage.glb(572180 bytes), dealer.glb(707356 bytes) and their manifests. The shared animation uses the actual pi city rotation; the exit moves away from the display car.

Both rigs pass sampled clearance against all18 doorway models, including the full dealership/car geometry (existing bottom12cm exclusion remains). All269 tests pass11.866s (.runtime/workshop-door-tests.log); build2.28s (.runtime/workshop-door-build.log), known chunk warning unchanged; diff check clean. CUA8897 tab236 inspected the dealership threshold2.3s. CUA8900 tab237 checked Russo Motor Works and Ferris Motor Sales in one timed UI sequence each:1.4443s showed the correct open door/audio0;7.7774/7.7769s showed actor outside, all doors closed, audio1. Stop cleared each private preview. Revision1/minute720 stayed unchanged; no gameplay/main-save writes. Caffeinate87182 confirmed alive. Prior4323f2c was concrete progress. The villa remains the final missing doorway; stair-aware gait, scene handoffs, asset fidelity and the broader full production objective remain outstanding and active.

### 2026-09-14 — Prepare the villa's raised entrance

Corrected the villa's overlapping ground-floor central window/door, replaced solid entrance masonry with a recessed vestibule, and authored a paneled hinged door/brass pull/jambs. Its existing foundation now meets a60cm porch landing and three consistent20cm risers. Preserved the raised entrance, columns, upper windows and roof rather than flattening away the stair requirement. Local Blender villa export is446648 bytes; measured manifest updated. The special-fire export left undertaker unchanged.

The entrance uses an entrance-landing marker, deliberately not entrance-threshold: the city planter remains disabled until stair-aware movement is implemented. This is incomplete entrance work, not completion of villa choreography. The shared review includes villa and caps playback at2.5s/on the landing, with Landing only labeling. Both actual rigs pass sampled door/masonry clearance through2.5s; exact downward geometry rays verify the interior floor, landing and60/40/20cm treads. Existing18 full exits remain covered. All270 tests pass11.603s (.runtime/villa-entry-tests.log); build passes (.runtime/villa-entry-build.log), known chunk warning unchanged. CUA8897 tab238 visually inspected the doorway/landing/steps at2.3s. No campaign operations. Caffeinate87182 confirmed alive and diff check clean. Prior70982f4 was concrete progress. The next required piece is faithful stair descent with grounded feet; full villa playback and the broader production objective remain active and unfinished.

### 2026-09-14 — Independent stair foot placement

Added CityLegs, a two-segment solver for the actual person/woman hip/knee rig, with separate ankle groups preserving authored shoe transforms. Reachable actor-local ankle targets bend knees forward while compensating ankle rotation to keep soles level. Unreachable/degenerate targets return false without changing the pose; segment lengths remain fixed. This is a prerequisite for stair descent, not yet integrated into city playback. The villa remains excluded from active planter playback pending continuous descent/clearance verification.

Geometry tests cover both rigs, both legs, nine ankle targets each under translated/rotated roots, actual ankle positions, level soles, actual shoe vertex heights and unchanged pose on unreachable input. All271 frontend tests pass11.451s (.runtime/stair-legs-tests.log); build2.19s (.runtime/stair-legs-build.log), existing chunk warning unchanged. After adding explicit sole-vertex assertions, the five planter/leg tests pass2.698s. Browser review tools/stair-leg-preview.html uses the real villa and three static stair-contact poses. CUA8897 tab239 inspected upper, lower and pavement transition poses; all reported both ankle targets reached. No gameplay/main-save writes. Caffeinate87182 alive, diff check clean. Previous d10cc07 was concrete progress. Continuous foot lifts, planted support, body transfer, stair clearance and animation integration remain next; broader production objective stays active.

### 2026-09-14 — Continuous stair descent prototype

CityStairDescent now drives six alternating foot transfers over the villa's three20cm risers/40cm treads. Each foot lifts before advancing past the nosing and lowers afterward; its partner remains fixed. Body motion excludes the swing-foot lift, correcting an initial reach failure at.708s. The3.9s step-together descent uses CityLegs for level soles, fixed segment lengths and explicit reach feedback. It remains a review-only motion component; villa city playback is still disabled pending approach/exit transitions and complete cast clearance. The stance is too crouched at rest and still needs a natural standing transition and polish before integration.

Both actual rigs pass120Hz samples through3.9s: all ankle targets reachable, support target stationary, all shoe vertices above the stair surface, final feet on pavement. Surface levels use the villa's previously geometry-verified tread profile; this test is not full-leg/building collision or a continuous swept-volume solver. All272 frontend tests pass11.817s (.runtime/stair-descent-final-tests.log); final build passes (.runtime/stair-descent-final-build.log), known chunk warning unchanged. CUA8897 tab240 played tools/stair-leg-preview.html: at1.7475s both targets reached mid-step; natural3.9s completed with both targets reached and feet on pavement. Browser screenshot confirms the remaining crouched resting pose. No gameplay/save writes. Caffeinate87182 alive; diff check clean. Prior b7233ff was concrete progress; full production objective remains active, including natural transitions, villa integration and all other unfinished work.

### 2026-09-14 — Upright stair preparation and recovery

The stair motion now begins upright, lowers into its working stance with both feet planted over.35s, and recovers standing height over.4s after both feet reach the pavement. Total review sequence is4.65s. This fixes the crouched resting pose observed in the previous browser review without lifting the planted feet or changing stair positions. Seek buttons use the adjusted timeline.

All272 frontend tests pass11.592s (.runtime/stair-settle-tests.log); build2.30s (.runtime/stair-settle-build.log), existing chunk warning unchanged. Extended120Hz coverage over the full new duration, body-position continuity and final upright height. Subsequently strengthened tread-clearance checks from shoes to all actual actor mesh vertices for both rigs; six focused planter/stair tests pass4.256s. This is vertex clearance against the previously verified tread profile, not all triangle intersections or whole-building clearance. CUA8897 tab241 played the full4.65s sequence and confirmed both targets reached, feet planted and the final upright stance. Initial screenshot arrived before assets drew and is not visual evidence of the starting pose. No gameplay or save operations; Caffeinate87182 alive; diff check clean. Prior538551e was concrete progress. Doorway approach/turn/departure transitions and city integration remain necessary before villa playback can be enabled. Full production goal remains active.

### 2026-09-14 — Join villa approach to stair descent

CityVillaExit now composes a3.9s doorway approach ending at the first tread with the4.65s stair descent, stopping upright on the pavement at8.55s. The same actor changes presentation parents; replay seeks reset the walking/ankle pose correctly. CityPlanter accepts optional exit/leave distances while retaining existing defaults. The combined test exposed a zero-distance incendiaryStride division by zero that produced NaN joints when sideways departure was disabled; zero distance now returns stationary distance/weight. The broken browser pose was reproduced in tab242, then corrected and rechecked.

Both rigs pass120Hz samples across the combined sequence: finite/reachable poses, bounded body-position changes, every actor vertex above the landing/tread profile, final completion/closed door and seek back to the start. All273 tests pass11.420s (.runtime/villa-sequence-final-tests.log); final build passes (.runtime/villa-sequence-final-build.log), known chunk warning unchanged; diff check clean. CUA8897 tab243 played the corrected sequence, inspected the4.48s stair handoff and natural8.55s completion upright on the pavement with door0. No gameplay/save operations; Caffeinate87182 alive. Prior4807b65 was concrete progress. The component remains review-only: final turn/departure, full body/door/column triangle clearance, city reservations/camera timing and detonation integration remain required. Full production goal remains active.

### 2026-09-14 — Integrate the complete villa planter sequence

Enabled the villa's entrance-landing cast in City3D. Its12.9s preamble combines doorway approach, stair descent, one planted forward clearance step, then a turn and2m departure. A first direct walking transition clipped a heel into the last tread at8.692s; the flat clearance step fixes that before rotation. CityStairDescent accepts a flat single-step profile and CityLegs reuses existing ankle groups across these phases. The common planter reservation is6.4m long, centered25cm toward the pavement, to include the extended departure. Existing flat preambles remain6.2s. Renderer blast/casualty/fire/audio gating and diagnostics now use each cast duration; camera pullback runs1.6s–.2s before the respective detonation. No backend outcome/time changes.

All273 frontend tests pass11.898s (.runtime/villa-city-final-tests.log); build2.32s (.runtime/villa-city-final-build.log), known chunk warning unchanged. The final combined test extends to12.9s for both rigs, checks all vertex tread clearance/body continuity/replay, and now verifies all vertices remain within the actual reservation. Seven focused tests pass4.738s (.runtime/villa-reservation-tests.log). This is sampled vertex/tread/footprint validation; whole-building triangle/swept-volume and natural animation acceptance remain broader work.

CUA8900 tab244 at Cypress House:6.8683s showed the actor descending, door closed and audio0 (blast still6.0317s away);13.8686s showed the departed actor(174,.18,8.17), closed door and audio1. Natural completion left no transient effects/open doors. Stop restored the private preview. Revision1/minute720 unchanged, no gameplay/main-save writes. Caffeinate87182 alive, diff check clean. Prior6553d66 was concrete progress. All19 current building models now have a planter preamble path; that coverage does not establish full production quality. Crowded staging recovery, actor handoffs, unsupported unnamed faction planters, premature injury/death choreography, gait/asset quality and all remaining full-objective requirements stay active.

### 2026-09-14 — Let pending planter exits clear of moving traffic

StreetTraffic accepts optional pending scene footprints. New moving arrivals cannot enter them, while existing visible occupants inside may leave on their existing committed route. City3D supplies the footprints of unstaged planter casts; once staged, the ordinary scene reservation takes over. This prevents a stream of arrivals repeatedly occupying a requested exit. Cancelling/completing removes the hold. Stationary occupants and aftermath are unchanged and can still block indefinitely; this does not solve that distinct recovery requirement or other action-cast stalls.

New deterministic regression starts one pedestrian inside and another approaching, checks no overlaps and no new entrant through900 updates, confirms the original occupant exits, then verifies travel resumes after release. All274 frontend tests pass11.416s (.runtime/pending-scene-tests.log); build2.40s (.runtime/pending-scene-build.log), known chunk warning unchanged. CUA8897 tab245 used tools/pending-scene-preview.html with the actual traffic scheduler and3D actor models: held state had original occupant atx70/outside and follower atx47.3933/outside; after Release scene, follower advanced tox56.5179. This isolated deterministic review verifies scheduler behavior, not a complete crowded campaign. No gameplay/save writes. Previous c11208e was concrete progress. Full production goal and stationary/crowded recovery, actor handoffs, animation/asset fidelity and other unfinished requirements remain active.

### 2026-09-14 — Stationary pavement clearance for pending exits

StreetTraffic can now move a stationary pedestrian along the same frontage to clear a pending planter footprint. It searches up to6m in either direction, checks the swept path in12.5cm samples against traffic/other held exits, stays away from intersections, and advances at the ordinary1.8m/s presentation speed. The saved request/location/progress remain unchanged. The resulting stance persists after release rather than snapping back. Parked vehicles remain fixed. City3D uses the actual stationary placement for scene occupancy and enables the walking gait during this move. This only applies on the known frontage pavement lane; fully occupied pavement, stationary vehicles, bodies and interrupted journey progress can still prevent staging. Current debug planted explosions replace the player with the planter, so this particular blocker requires the isolated two-person review rather than claiming coverage from an ordinary player-plant replay.

All276 frontend tests pass11.678s (.runtime/stationary-yield-tests.log); build2.62s (.runtime/stationary-yield-build.log), existing chunk warning unchanged. Regression covers30/60/144Hz bounded displacement, unchanged saved request/progress, retained stance on release, selecting the other side of an occupied frontage and keeping a parked car fixed. CUA8897 tab247 reviewed tools/pending-scene-preview.html?stationary with actual3D people: the blocker cleared the footprint to x45.33246, nearby person stayed x44, both visible/outside the held area; Release retained the cleared stance. This deterministic review does not prove full crowded campaign or final gait fidelity. No campaign writes. Caffeinate87182 confirmed live; diff check clean. The previous response supplied fresh browser evidence that the requested accounts behavior was already complete; no implementation was needed there. Full production objective remains active, including staging recovery, actor handoffs, all unfinished action variants, asset/animation quality and campaign acceptance.

### 2026-09-14 — Rejoin journeys from the visible pavement stance

The previous stationary-clearance work exposed a handoff defect: replacing its stationary request with a real journey reset the actor to the route start in one frame. StreetTraffic now remembers the original stance only when a cosmetic sidestep actually moves, and connects a subsequent same-model journey from that same origin back to its physical start. The connector advances at pedestrian speed with12.5cm collision substeps. Saved journey progress remains zero until rejoining; blockers stop the connector without hiding the person. A changed model or different origin retains existing handling. This does not solve general scene-cast handoffs, vehicle boarding or unrelated route transitions.

All278 frontend tests pass11.416s (.runtime/yield-handoff-tests.log); build2.32s (.runtime/yield-handoff-build.log), existing chunk warning unchanged. New regressions check30/60/144Hz bounded positions, zero journey progress before rejoining, actual passage through the start, onward travel, and waiting/resuming at an occupied start with no overlap. CUA8897 tab248 used the actual traffic scheduler/3D person review: cleared stance x45.32904; Start journey then showed x47.04138/progress0 during rejoining, followed by x59.5893/progress.965775 on the committed route. The review rig's basic gait is not final animation acceptance. No gameplay/save writes. Caffeinate87182 live; diff check clean. Prior3fd2ab9 was concrete progress. Full production objective remains active, including remaining crowded staging, named NPC preambles, general handoffs, richer assets/actions and campaign acceptance.

### 2026-09-14 — Assign present faction planters and save their escape

DemolitionDay now delegates its resolved blast to factionDetonation. It selects a living, free faction member actually inside the target, with no existing departure schedule and a valid different home, preferring skill then stable ID. That person begins a normal saved homeward journey before damage/casualty resolution. The explosion cue captures their ID/name with weapon0, enabling the existing planter preamble and excluding the departing actor from indoor blast casualties. Damage, daily bombing chance and charge price are unchanged. When no such member is present, the existing anonymous blast remains: dispatching somebody from elsewhere is still outstanding. No new save fields/endpoints.

The faction-planter QA scenario creates a fresh isolated save, places one existing Bellandi member in player-owned Bluebird Laundry, and calls the actual daily bombing rule until it resolves (bounded2000 attempts without clock steps). This is a forced test setup, not proof of ordinary campaign frequency. .runtime/faction-planter-20260914.sqlite3 is served by .runtime/faction-planter-server on8901/session54473. First launch used an incompatible host-prefixed address and exited; corrected launch with :8901 is live. Main campaign untouched.

CUA8901 tab249 replayed the recorded event: Vittorio Bellandi appeared as the captured attacker,1.8403s showed the open laundry door/audio0;20.0341s showed the exited actor/closed doors/audio1. Natural completion removed transient effects; revision0/minute600 stayed unchanged. This exposed remaining handoff/UI defects: the normal escaping actor waits at the player's occupied departure afterward, and this replay exposes a Herald link instead of automatically opening the newspaper. Those are not claimed fixed. Targeted Go tests pass.301s (.runtime/faction-planter-tests.log), including unavailable-person rejection, saved escape/arrival and64 fatality seeds proving the escaped planter survives while actual occupants can die. QA command compiles; backend binary builds. Full Go suite evidence recorded below on completion. Caffeinate87182 confirmed live. Prior7f1949e was concrete progress; full production objective remains active.

The same browser check located and fixed the newspaper defect: detonate filed EXPLOSION DESTROYS / ONE DEAD IN EXPLOSION but gave its cue EXPLOSION AT, breaking the frontend's deliberate exact headline/minute match. The cue now uses the already-filed headline. The64-seed regression checks matching news for every fatal/nonfatal case. Expanded focused core tests pass1.037s, including the existing faction-demolition campaign test. Fresh .runtime/faction-paper-20260914.sqlite3 on8902/session90173 (new rebuilt binary) verified the fix through CUA tab250:1.6455s preamble/audio0 and no paper; natural completion removed effects and automatically displayed the correct styled article and Fold away control. Screenshot confirms the reveal; the initial case-sensitive mixed-case text probe returned false because the rendered edition text is uppercase, corrected probe returned true. Revision0/minute600 unchanged, no gameplay commands. Occupied-departure handoff remains outstanding.

Full go test ./... remains live at this checkpoint (session6405, go PID48617, sim.test PID48694 confirmed active at6m26s). Core passed126.785s; command packages passed. The simulation package includes the pre-existing untracked TestWhatKillsThem, which runs260 campaigns, and is still consuming CPU. Do not restart or claim a full-suite pass until this same handle completes. The final headline-only change is covered by the later focused tests and rebuilt/reviewed browser fixture. Other dirty gameplay files remain untouched and unstaged.

### 2026-09-14 — Clear occupied pedestrian departures

Extended stationary pavement clearance to waiting pedestrian departure footprints as well as pending planter reservations. A standing person can now walk out of a waiting traveller's start space using the existing bounded, collision-checked frontage move. This addresses the concrete blocked departure observed in the named faction bombing review. It does not move parked cars, advance an uncommitted journey, or solve the remaining jump from the scene cast's final position to its ordinary street actor.

All279 frontend tests pass11.947s (.runtime/departure-clearance-tests.log); build2.50s (.runtime/departure-clearance-build.log), existing chunk warning unchanged. New30/60/144Hz regression keeps a real waiting departure at committed progress0 until room is made, checks no overlap/hidden stationary person, then advances the departure after its requested progress changes. CUA8902 tab251 reviewed the existing isolated faction fixture: Vittorio visible at(48,68.65), player stepped to(46.6734,68.65), waiting[]. Replaying the full recorded explosion and allowing natural completion retained both visible actors, waiting[], and automatically opened the newspaper. Revision0/minute600 unchanged. No campaign commands or main-save writes. Caffeinate87182 live; diff check clean. Prior8939cf1 was concrete progress. Full production objective remains active, especially continuous scene-to-street handoffs, dispatching faction planters from other premises, remaining action variants, asset fidelity and campaign acceptance.

The prior full Go run is now terminal: session6405 exited1 after sim553.517s. Core/command/store packages passed, but TestNothingAPlayedCitySaysIsMalformed reported a real magpie-seed2 sentence: “They went out to  and came back hurt.” Trace points to itWentWrong in core/strike.go passing an empty task description into HandHurt. The pre-existing long death test was not the reported failure. This finding is being addressed separately; no full-suite pass is claimed.

### 2026-09-14 — Restore the task in failed delegated strike reports

The full-suite prose failure was caused by two itWentWrong calls supplying an empty description to HandHurt. Both now pass the actual target and address (go after NAME at PLACE). This preserves injury, loyalty, death odds and command outcomes while making the injury/death account meaningful. A64-seed direct failed-delegation regression requires returning-injured reports to include both names. Initial test compilation referred to Records rather than the world's History field; corrected before validation. Focused delegation/strike tests pass.622s (.runtime/delegated-strike-prose-tests.log). The formerly failing39-campaign TestNothingAPlayedCitySaysIsMalformed is rerunning independently, with no need to repeat unrelated long tests; result to be recorded on completion. Existing dirty armed/mugging/robbery/death-test files remain untouched.

The exact previously failing39-campaign prose gate now passes145.865s (.runtime/delegated-strike-campaign-prose.log; session96994 completed exit0). This closes that observed malformed-sentence failure. Other packages and simulation tests passed in the preceding full run; this was a targeted rerun of its sole reported failure after the correction, not a newly repeated full-suite invocation. Full production-game acceptance remains incomplete.

### 2026-09-14 — Do not report a dead delegated attacker escaping

A follow-up action-coherence check reproduced another failed-strike contradiction: itWentWrong's escape branch calls HandHurt with a severe injury, which can kill the crew member, then unconditionally claimed they got out. The new128-seed regression failed at seed23 before the fix. After HandHurt, the branch now returns when the crew has been removed; the existing death account remains and no false escape account follows. Surviving crew retain the existing escape account. No outcome probabilities, injury, loyalty or death rules changed. Focused failed-strike/delegation/personal-strike tests pass.720s (.runtime/crew-survival-report-tests.log), including assertions for both fatal and surviving outcomes. Go tests compile the changed package; no frontend changes required.

Also checked the premature-explosion path rather than assuming last turn's newspaper fix applied to it: its existing Report/Witness headlines already match. CUA's old tab251 handle was unavailable, so opened fresh isolated8902 tab252. At Saint Agnes, Explosion · premature showed one blast sound, glassBreaks0, no planter and no window-burst list. Stop removed all transient effects; revision0/minute600 stayed unchanged. No gameplay/save writes. This confirms separation from the successful preamble, not completed injury/death choreography, which remains outstanding. Caffeinate87182 confirmed alive; diff check clean. Prior2426586/7f480bc were concrete progress, and the previous failing39-campaign prose rerun completed successfully. The full production objective remains active with all unfinished animation, asset, dispatch, handoff and campaign requirements intact.

### 2026-09-14 — Record premature accident injury for faithful replay

Added optional CueAccident/VisualCue.accident with actual health_lost and fatal. Plant captures the health decrease after armour/clamping and the zero-health result before calling DieOf or advancing command time. The event's existing attacker identifies who suffered it. Successful detonations/legacy records omit the field; missing data does not mean survival. Updated the browser contract and API documentation. No gameplay damage, odds, duration or save migration changes. This is the causal input for injury/death choreography, not the completed visual sequence; the renderer does not yet consume the new field.

128-seed regression exercises both fatal and surviving accidents, checks actual health loss, changes the current player's health/life status, and verifies serialized cue replay retains the original outcome. Existing command-result/save test now checks accident presence only on premature blasts and preserves its values through the full World JSON round trip. Focused demolition tests pass.441s (.runtime/accident-outcome-tests.log); command/API/store boundary tests pass (.runtime/accident-outcome-boundary-tests.log: server.728s/store.169s). TypeScript/production build2.42s (.runtime/accident-outcome-build.log), existing chunk warning unchanged. No browser visual claim or campaign operations for this contract-only change. Caffeinate87182 live; diff check clean. Priorc98658e was concrete progress. Next required work includes staging the injured/fatal actor from this recorded outcome, then integrating its animation without clipping, duplicate actors or premature newspaper/death overlays. Full production objective remains active and unchanged.

### 2026-09-14 — Grounded accident fall and survivor recovery study

Added review-only CityAccident using the actual person/woman joint rigs. A short fall transitions to a prone resting pose; the supplied fatal outcome prevents recovery, while a survivor begins sitting up after2s. Total study4s. Posed mesh vertices determine ground contact, including coat/hands, under translated/rotated roots. Timeline seeking resets the pose. This is a first motion study, not city-integrated accident choreography or final animation quality. Supporting hands, believable body weight, impact effects, casualty aftermath, collision-reserved staging and player-death overlay timing remain outstanding.

The first browser pass (CUA8897 tab253) revealed an implausible suspended survivor pose caused by the lower-leg bend. Corrected the recovery to extend the legs and reviewed it in tab254, including playback at2.4496s. Added actual cast/receive shadows to the review for clearer contact assessment. Fatal variants remain prone. All280 frontend tests pass11.408s (.runtime/accident-motion-tests.log); build2.48s (.runtime/accident-motion-build.log), existing chunk warning unchanged. New120Hz four-second samples cover both actual rigs and both outcomes, lowest mesh contact5mm above ground, finite bounds, continuous root positions and replay reset. These do not prove self-collision, whole-body support mechanics or scene/building clearance. tools/accident-preview.html remains an isolated browser review; no gameplay/save operations. Prior492c108 was concrete progress. Full production objective remains active.

### 2026-09-14 — Support the accident recovery with fixed hand targets

CityAccident now introduces wrist pivots preserving the authored hand transforms and uses the existing two-segment arm solver to bring survivor hands beside the body before the torso rises. From1.95s onward their root-space targets remain fixed at x±.43/y.084/z-.6; wrists compensate the changing actor/arm orientation so hands stay flat. Fatal poses do not brace/recover. Seeking resets arm/elbow/wrist rotations before reconstructing the pose.

Clearance tests first caught a brief hand/cuff intersection during approach, then a cuff dipping during recovery. Added a small lifted approach and changed the elbow bend hint into the root-space upward direction, keeping the forearm above the floor; increasing hand height alone was insufficient. Both actual rigs now pass the existing120Hz whole-mesh ground/continuity check and a new fixed-wrist/no-palm-penetration/replay-reset check through recovery. The support target has a small mesh clearance rather than exact skin contact; full anatomical support/self-collision and sleeve deformation remain quality work. All281 frontend tests pass11.572s (.runtime/accident-support-tests.log); production build passes (.runtime/accident-support-build.log), existing chunk warning unchanged.

CUA8897 tab255 inspected initial supported motion at2.4491s; tab256 inspected corrected final poses after the upward elbow adjustment. The review remains intentionally separate from city playback; no campaign operations. Caffeinate87182 live; diff check clean. Priorc284218 was concrete progress. City integration, staged impact/aftermath, death-overlay timing, scene-to-street handoffs, asset refinement and all other full production requirements remain active and unfinished.

### 2026-09-14 — Reserve the complete accident fall footprint

Added an accident traffic footprint (3m length/2.2m width), offset45cm toward the fall, and three candidate frontage staging slots. Slots orient the fall along the pavement, outside the road and conservative17m building envelope. Existing occupied-slot checks reject collisions and return no slot when all candidates are occupied. This prepares CityAccident for city integration; it does not yet instantiate the cast from a cue or handle its impact/death overlay/aftermath.

Measured both actual rigs through480 frames: maximum lateral reach about.6493m and longitudinal bounds -1.4482..+.8324m. New regression checks every posed mesh vertex against the actual rotated/transformed reservation for both outcomes, both rigs and four headings; verifies alternate-slot selection, fully occupied rejection, and rectangle clearance from road/building envelope. All282 frontend tests pass11.445s (.runtime/accident-reservation-tests.log); build passes (.runtime/accident-reservation-build.log), known chunk warning unchanged. These are sampled geometry and conservative-envelope checks, not arbitrary prop/building triangle or continuous swept-collision proof.

Added reservation outlines to tools/accident-preview.html. CUA8897 tab257 played the motion and inspected containment at1.8451s. No gameplay/save operations. Caffeinate87182 live; diff check clean. Prior51fe01a was concrete progress. Full production goal remains active, with integration of accident animation/impact/outcome, scene handoffs, asset quality and all other unfinished scope unchanged.

### 2026-09-14 — Stage recorded accident outcomes in the city

City3D now creates CityAccident only for a premature explosion with both a captured attacker and accident outcome. It uses the actor's model/wardrobe, reserves an accident frontage slot, orients the fall along the pavement, frames the full cast, and suppresses the duplicate normal actor. Blast/light/debris originate at the staged accident rather than the unrelated building window. Unstaged casts and their debris wait for traffic clearance. Fatal actors remain prone; survivors begin the supported recovery. Existing explosion lifetime/cleanup remains; legacy unknown-outcome cues keep their prior blast presentation. The renderer still needs persistent accident aftermath, seamless handoff and actual player-death overlay timing.

Debug now supplies an explicit nonfatal injury for Explosion · premature and adds Explosion · fatal accident. Neither private preview alters player health or invents building fire. Preview regression verifies both outcome fields and no fire; all282 frontend tests pass12.286s (.runtime/accident-city-tests.log); final build2.62s (.runtime/accident-city-build.log), known chunk warning unchanged. Final diagnostic blast origin reports the actual staged origin.

CUA8902 tab258 at Saint Agnes: surviving cast2.3471s was staged, recovering, audio1; fatal cast6.5143s remained prone at rotation-π/2/audio1. Replayed survivor9.7152s reached the supported final rotation-.7708, visually inspected in the city. An earlier survivor screenshot was taken after natural cleanup and is not evidence of its pose. Stop cleared the private previews. Revision0/minute600 unchanged; no gameplay/main-save writes. These checks exercise the city debug path, not an actual fatal player command behind the death overlay. Caffeinate87182 alive; diff check clean. Priorda7988f was concrete progress. Full production objective remains active, including the accident follow-through, higher asset/animation quality and all other outstanding requirements.

### 2026-09-14 — Reveal the memorial after the fatal scene and paper

Deferred the player memorial until no recorded scene or scene newspaper remains. Focus now follows the memorial's actual appearance, including when the newspaper closes; a loaded dead save still opens and focuses the memorial immediately. Gameplay stays inert while the player is dead. No simulation or public-interface change.

Added fatal-charge to qa-fixture: probes disposable seeded worlds, then saves the untouched pre-command setup, refusing existing output paths as before. Built a fresh backend and served only .runtime/fatal-charge-20260914.sqlite3 on8903 (session11004). CUA tab260 entered The Monarch and committed the actual charge command. At5.8875s the fatal player cast was staged/prone with one blast, and neither memorial nor newspaper covered it. After natural completion the paper was the sole dialog; folding it away mounted/focused the memorial. Fresh tab261 opened the saved death directly with memorial focus. No new life was started; main campaign untouched.

Build passes2.42s (.runtime/death-handoff-build.log), existing chunk warning unchanged. Focused authoritative charge outcome tests pass.376s (.runtime/death-handoff-core.log); diff check clean. Browser validation supplies the UI sequencing evidence. The previous accounts turn verified the requested one-click breakdown, requiring no further edit. Caffeinate87182 remains live.

Remaining observed defects: fatal premature-blast prose still claims somebody escaped, and the in-scene Skip button inherits the dead shell's inert state. Natural completion and newspaper controls work, but skip accessibility needs a separate presentation-control boundary. Persistent accident aftermath and all wider production requirements remain unfinished; this is not whole-goal acceptance.

### 2026-09-14 — Accessible fatal-scene controls and truthful blast reporting

Fatal scene Theatre controls now render through a portal outside the inert gameplay shell, in the existing newspaper styling. Skip receives initial keyboard focus; the controls disappear when the scene newspaper or memorial takes over. Gameplay remains inert. Fatal premature charges file a one-death headline and report instead of claiming the dead player escaped; the cue uses the same headline for article linkage. Survivors retain their existing escape report. Removed the generic attack standfirst's unsupported claims of extensive building damage and arrest prospects.

Extended the 128-seed accident regression to require a matching article, correct fatal/surviving report, and no invented extensive damage. Charge/newsroom tests pass37.017s (.runtime/fatal-report-tests.log); frontend build passes2.33s (.runtime/fatal-controls-build.log), known chunk warning unchanged. No public schema or gameplay odds changed. Diff check clean.

Fresh isolated fatal-controls-20260914.sqlite3 on8904 (server session52633), CUA tab262: actual charge command produced fatal scene with focus on Skip, no inert ancestor on that control, gameplay shell still inert, memorial absent. Pressing Enter opened ONE DEAD IN EXPLOSION AT THE MONARCH with the corrected fatal body and removed the portal. Effects cleared; revision1/minute600 stayed unchanged by Skip. Folding the paper mounted/focused the memorial. The browser check preceded the standfirst correction; that final copy change is covered by the Go projection regression, not a rebuilt browser backend. Main campaign untouched; no new life started. Caffeinate87182 live. Priorf2e1c69 was concrete progress. Persistent accident aftermath, visual fidelity, and the remaining full production scope are still open.

### 2026-09-14 — Persist fatal charge accident bodies

Fatal player charge accidents now leave an authoritative aftermath entry with life-specific victim identity, captured name/face, cause charge-accident, and existing +5-minute police/+180-minute cleanup deadlines. Survivors leave no body. The renderer selects the captured rig/wardrobe, reconstructs CityAccident's final fatal pose, and retains the staged accident reservation through natural completion. Active real accident playback suppresses its duplicate body; private previews do not suppress a real saved death. Reloads choose available frontage slots, as exact rendered coordinates are not saved. API.md documents the optional fields and limits.

Regression traverses every named mesh transform for both actual character rigs to compare final animation and persistent body, checks reservation retention, police entry and cleanup. All283 frontend tests pass11.615s (.runtime/accident-aftermath-tests.log). Initial TypeScript failure on a union's unknown face field was corrected with an explicit number check. Final build passes2.57s (.runtime/accident-aftermath-build.log), known chunk warning unchanged. Focused Go accident/aftermath tests pass.312s (.runtime/accident-aftermath-core.log), including both outcomes over128 seeds, saved appearance after reload/player replacement, exact cleanup, and cue/article linkage.

Fresh isolated accident-aftermath-20260914.sqlite3 on8905 (server session27318), CUA263 actual fatal charge: at.1471s one staged fatal effect, no duplicate aftermath. After natural completion effects were empty and the retained body was visible at108,70.35 behind the newspaper, revision1/minute600. Fresh tab264 loaded the same persistent body behind the memorial at the same available slot. No main campaign writes or new-life command. Browser diagnostics establish persistence/placement; rig tests establish matching transforms. Caffeinate87182 live, diff check clean. Prior97b2832 was concrete progress.

Police arrival/cleanup are still stationary time-based presentation, not animated response. A dead campaign is paused, and starting another life advances eight hours, so this does not yet provide a watched post-death police investigation. Survivor handoff, fuller gore/aftermath choreography, asset fidelity, and the full remaining production scope stay open.

### 2026-09-14 — Ground posed strike victims and persistent bodies

Appearance tracing found no current NPC palette mismatch: the live and retained rigs use the same identity-derived wardrobe. Geometry inspection instead found the fixed casualty height left posed strike victims visibly floating: sampled final rig minima .195m/person and .2374m/woman relative to the scene root, which itself sits at.2m. Added groundCharacter to seat visible posed mesh vertices at the existing character contact plane. Assassination variants, generic casualty falls, and non-accident persistent bodies use it; accident grounding remains its existing implementation. The helper handles transformed parents and ignores hidden wardrobe geometry.

Added sampled contact/vertical-continuity checks for both actual rigs and four strike variants at30Hz, with translated/rotated roots. The existing handoff test initially failed because its live rig exposed every authored hair/hat while the aftermath applied wardrobe visibility; corrected that fixture to dress the live actor as the real renderer does. Final all284 frontend tests pass11.414s (.runtime/body-ground-tests.log); build2.23s (.runtime/body-ground-build.log), known chunk warning unchanged. No Go or public-interface change. Measurements target the existing .205m contact plane and do not establish arbitrary terrain or full anatomical collision fidelity.

CUA265 on isolated body-ground-20260914.sqlite3 port8906 (session75463), replayed the recorded shotgun strike. At.2431s the staged scene had no shots yet; after natural completion the paper appeared, then closing it exposed the retained Mara body at81,38.35. Screenshot inspected: body and blood remain on pavement; revision0 unchanged. This checks a recorded fixture replay, not a new gameplay command. Main save untouched, caffeinate87182 live, diff check clean. Priorf1d8323 was concrete progress. Survivor recovery/handoff, response choreography, higher-quality assets, and all remaining production scope stay open.

### 2026-09-14 — Recover from a surviving charge accident

Extended the survivor motion from supported sitting through gathering the legs, releasing the palms, crouching, and standing by7.5s. Fatal motion still ends prone at4s. Added ankle pivots via the existing leg rig: soles level during the gathering phase and remain stationary throughout the final5.3–7.5s rise. The original14-second explosion scene lifetime is unchanged. The standalone motion review now has a time input and the full recovery timeline.

Initial crouch review at5.3s exposed tilted feet, corrected before accepting the rise. Reservation regression then caught forward silhouette overflow at4.9667s and5.0667s. Measured maximum forward reach1.1613m/person and1.0815m/woman for the intermediate foot anchor; moved the gathered sole center to.32m so the full recovery fits the existing reservation. Expanded contact/continuity and four-heading reservation checks through each outcome's full duration. New test checks every shoe vertex remains fixed during the final rise for both actual rigs and that all principal joints finish upright and reset for replay. All285 frontend tests pass11.717s (.runtime/accident-rise-tests.log); build2.39s (.runtime/accident-rise-build.log), known chunk warning unchanged. This verifies sampled geometry/support, not anatomical self-collision or film-quality motion.

CUA8897 tabs266–268 inspected the initial crouch, leveled-feet revision, and7.5s upright survivors beside prone fatalities. CUA8902 tab269 private premature-blast preview at Saint Agnes: first screenshot came after cleanup and is not pose evidence; replay sampled2.2424s recovering and10.3398s upright (rotation0, one blast), with the latter screenshot inspected in the city. Stop cleared effects; revision0/minute600 unchanged. No gameplay/main-save writes. Caffeinate87182 live; diff check clean. Prior189ca2d was concrete progress.

Scene-to-normal-actor positioning still needs a seamless handoff; the recovery does not solve that relocation. Injury-specific gait, richer body/face animation, asset fidelity, animated emergency response and all wider production scope remain open.

### 2026-09-14 — Transfer recovered pedestrians back into street traffic

Added a validated same-frontage adoption operation to StreetTraffic. A completed or skipped real surviving player accident hands its final upright world position and heading to the stationary pedestrian entry, preserving its canonical route origin for later travel. The next route connects first to the walking lane and then along the same frontage to the origin, using existing collision/substep checks; saved progress stays zero through the connector. Natural cleanup immediately exposes the normal actor at the transferred position. Fatal, preview, moving, vehicle and different-location cases do not use this path. API.md documents presentation-only behavior and scope.

Added survivor-charge isolated fixture alongside fatal-charge, selecting the untouched pre-command setup by disposable seeded probes. Tests cover adopted stance retention, continuous departure at30/60/144Hz, occupied-lane waiting/resumption, and rejecting cars/unrelated frontages. All287 frontend tests pass11.392s (.runtime/survivor-handoff-tests.log), build2.23s (.runtime/survivor-handoff-build.log), known chunk warning unchanged. Fixture compiled/ran successfully. No public schema or simulation rule change.

CUA270, fresh .runtime/survivor-handoff-20260914.sqlite3 on8907 (session72051): actual charge resolved with4 health, rev1/minute720. Natural completion produced the newspaper and ordinary player at108.2499999997,70.35, matching the recovered frontage position rather than canonical112,68.65. After closing the paper, committed travel to Saint Agnes: first sample retained x108.25 while z moved continuously to69.9126; later sample103.0633,59.35 on the route, waiting[]. Revision2/minute745 is the committed travel outcome; no playback action changed it. Main save untouched. Caffeinate87182 live; diff check clean. Priorcbcee35 was concrete progress.

This establishes pedestrian survivor handoff, not all actor/vehicle handoffs. Skip during a fall intentionally resolves to the final upright stance; dedicated browser Skip handoff coverage remains pending. Fully blocked lanes, turn smoothing, injury-specific gait, police/cleanup choreography, high-quality assets and the full remaining production objective stay open.

### 2026-09-14 — Verify early Skip and replay after survivor handoff

Closed the browser Skip coverage gap from0fb44aa using a fresh isolated survivor-skip-20260914.sqlite3 on8908 (server session26881). CUA271 committed the actual surviving charge. At.1458s the staged player accident was still starting; immediate Skip cleared all effects and opened the scene newspaper, preserving the recovered pedestrian at108.2499999997,70.35. Revision1/minute720 remained the committed outcome. Closing the paper, replaying the same recorded event, skipping again and closing the second paper left exactly one visible player at that position, no effects, no waiting entries and no open dialog. No gameplay/main-save changes beyond the isolated charge; no new life or extra travel.

No implementation adjustment was needed for this path, so no redundant build or unit-suite run. Existing287-test/build evidence remains from0fb44aa; this turn adds actual early-Skip and repeated-replay evidence that was previously missing. Caffeinate87182 verified live. Previous0fb44aa was concrete progress; this verification completes its outstanding browser Skip check. Vehicle representation during pedestrian action scenes, broader handoffs, animation/asset fidelity, emergency response and all remaining production requirements are still unfinished.

### 2026-09-14 — Keep the owned car after a fatal scene

Tracing confirmed ordinary on-foot player scenes already retain a separate player-car actor. The actual defect was its alive guard: a committed player death removed a still-existing saved vehicle immediately. Removed that guard; the empty parked car now follows authoritative vehicle/location state through death, and disappears when that state changes. No simulation, ownership, damage or public-schema changes. Added fatal-charge-car to the isolated fixture scenarios with a Ford and known vehicle condition/fuel.

Build passes2.20s (.runtime/fatal-car-build.log), known chunk warning unchanged; fixture compiles and creates a fresh save. CUA272 on .runtime/fatal-car-20260914.sqlite3, port8909/session81000: pre-command Ford at121.6,80; actual fatal charge at.1319s retained that same car, no driver, zero wheel phase/steering, waiting[]. After natural completion the persistent body at108,70.35 and the unchanged car were both present with no effects. Closed the newspaper and started a new life only in this explicitly isolated save: Life2 Nico Ward, rev2/minute1080, car absent and expired aftermath empty. This validates both persistence and cleanup against saved state; it does not animate abandoned-car collection or vehicle damage.

No redundant broad unit run for this one-guard presentation change; browser coverage and build supply the relevant evidence. Main campaign untouched; caffeinate87182 live and diff check clean. Priorb7d980f completed outstanding Skip QA. Full asset/animation fidelity, richer vehicle interactions, emergency response, campaign acceptance and the wider production objective remain unfinished.

### 2026-09-14 — Align emergency response with the visible journey clock

Added an actual-mesh reservation regression for both victim rigs, four strike variants, police car and officers. All visible body/response vertices fit their claimed traffic rectangles, which remain pairwise clear. This verified the footprint concern; timing inspection found a separate real defect: CityAftermath used the committed endpoint minute while travel and the HUD showed an earlier interpolated minute.

Journey presentation now retains pre-command aftermath, police-presence and fire records, merges newer endpoint records by ID, and evaluates attendance/expiry at the displayed travel time. Fire visuals, suppression, rubble and raid-door presence use that same minute. This preserves records that expire within a journey without changing simulation state. Records missing from both endpoints cannot be reconstructed; this does not establish a complete event history or historical building/lighting playback. API.md records the boundary.

CUA273 on the existing isolated body-ground-20260914.sqlite3 port8906: committed travel from Saint Agnes to The Monarch, endpoint505. At displayed480.0852 and484.6229 only the body was present. At486.9835, after the scheduled485 response, the body, one police car and two officers were visible. No early police response; only the isolated travel command advanced the save. Main campaign untouched.

Final all290 frontend tests pass13.178s (.runtime/response-clock-tests.log), including exact police/cleanup thresholds across merged endpoints and a fire absent from the final snapshot through brigade arrival, extinction and cleanup. Build passes2.44s (.runtime/response-clock-build.log), known chunk warning unchanged. Fire boundary coverage is component-level; the actual browser timing check covered police attendance. Caffeinate87182 live; diff check clean. Priorb699357 was concrete progress. Animated arrival/investigation/cleanup, full event-history playback, higher-quality assets and all remaining production scope stay open.

### 2026-09-14 — Keep city clocks and lighting on travel time

Moved day/night lighting, window emission, streetlight pools and vehicle lamps onto presentationMinute. Clock hands now update from that same interpolated minute through cached object references. Lighting/material traversal is cached by day/night/weather and invalidated for a new snapshot, avoiding a whole-city traversal every frame. This retains the existing6a.m./8p.m. lighting thresholds; gradual twilight and historical weather reconstruction remain open.

Added city3d-dusk isolated fixture at19:55 with walking travel and no fixture sabotage. Removed unintended redundant dusk blocks from unrelated fixture branches during final diff review; final fixture compile succeeds (.runtime/city-clock-fixture.log). Build passes2.32s (.runtime/city-clock-build.log), known chunk warning unchanged. No public-interface or gameplay change; no redundant broad test run for this presentation scheduling change.

CUA274 on fresh .runtime/city-clock-20260914.sqlite3, port8910/session93824: committed Saint Agnes→The Monarch travel ends at1220. At displayed1195.0857 and1199.2622, lighting remained day; both hour/minute hand pairs advanced with those minutes. At1201.8247 night lighting was active and the minute hands wrapped to.19108rad, consistent with20:01.8247. Screenshot inspected at the clock tower after nightfall. Only this isolated travel mutated state; main campaign untouched. Caffeinate87182 live, diff check clean. Prior55730b6 was concrete progress. Lighting aesthetics, complete historical event/weather playback, animated response, asset quality and the full remaining production scope stay open.

### 2026-09-14 — Gradual twilight during city travel

Revalidated the previous Daily accounts acknowledgment: the requested visible breakdown was already implemented, so that turn made no new goal progress. Read the full amended objective and continued with the observed abrupt day/night transition.

City sunlight, ambient light, sky/fog colors, window emission, streetlight pools and moving vehicle lamps now blend over 05:30–06:30 and 19:30–20:30 using the interpolated travel minute. Twilight caches luminous material references on snapshot/weather changes rather than traversing every building each frame; background/fog objects are reused. No simulation, public API or save changes. Committed endpoint weather is still used; historical weather reconstruction remains unfinished.

All 292 frontend tests pass in 11.494s (.runtime/twilight-tests.log), including dawn/dusk endpoints, minute continuity, midnight/negative-day normalization, weather independence and stable paused inputs. Build passes in 2.24s (.runtime/twilight-build.log), existing chunk warning unchanged. Caffeinate PID87182 confirmed live.

CUA276, isolated .runtime/city-twilight-20260914.sqlite3 on port8911/session76140: actual Saint Agnes→The Monarch walking command committed minute1220. At displayed1195 night amount .3762/sun2.128; sampled1197.72–1198.16 smoothly advanced .4431→.4541, sun1.9370→1.9059. At1204.49 amount .6114/sun1.4574, confirming the former20:00 switch no longer jumps to full night. Screenshot inspected at19:58. Skip moved presentation to1220 with amount.9259/sun.5611, stable across1.2s with revision1 unchanged. Main campaign untouched. Dawn continuity is unit-tested; browser travel covered dusk. The screenshot's145FPS counter is not a full performance acceptance result. High-quality assets, animated response, historical event playback and the broader production scope remain open.

### 2026-09-14 — Keep watched characters readable through foreground scenery

Prior231fe61 was concrete verified progress. Re-read the full amended goal. The sunset screenshot showed a partly hidden walker; inspection found torso-only building rays, permanent2% cutaway grain, and street trees outside the cutaway system.

Building visibility now samples torso plus head/foot corners across the pedestrian outline, ignoring hidden objects, hidden ancestor variants and hidden materials. The opening's center clears completely once settled while its edge retains the existing dither. Street furniture uses cloned cutaway-capable materials on the existing instanced meshes: fragments inside the opening dissolve only up to the watched cast's farthest projected outline depth. Scenery behind that depth remains rendered. No public API, gameplay or save changes. Large vehicle outlines and a full range of scene/camera acceptance remain to be verified.

First browser pass CUA277 port8912 at1200.63 reproduced the remaining tree obstruction despite the station cutaway; extended the implementation based on that evidence. CUA278 port8913 then showed the walker unobstructed through both station and tree. Final build CUA279, fresh .runtime/cutaway-final-20260914.sqlite3 on8914/session38915: actual Saint Agnes→Monarch travel, at1199.48 player72.00092/.2/27.35 visible with precinct cutaway; screenshot inspected. At1212.26 Stop following cleared the cutaway set and restored foreground scenery, revision1 unchanged by the camera control. The character paused near the crossing before moving onward; waiting[] at the sampled pause, so route pacing warrants a separate investigation rather than assuming a collision defect. Only these identified QA saves were changed; main campaign untouched.

All294 frontend tests pass11.837s (.runtime/character-cutaway-tests.log): actual thin roof geometries hide head/feet despite clear torso, invisible variant/material exclusion, cutaway depth uniforms and restoration, plus prior coverage. Final build2.52s (.runtime/character-cutaway-build.log), known chunk warning unchanged; final follow-up expanded screen/depth bounds to the already-tested outline, followed by browser verification. Diff check clean. Caffeinate87182 live. Screenshot145FPS is only a momentary counter, not broad performance acceptance. Asset quality, full historical travel/event playback, response choreography and the remaining production goal remain open.

### 2026-09-14 — Explain visible travel stops from actual occupancy

Prior8c8a5d5 was concrete progress. Re-read the full amended objective and traced the apparent pause from the cutaway QA using the public port8914 snapshot, saved locally as .runtime/cutaway-state.json. A deterministic replay (.runtime/trace-walk.mjs, .runtime/walk-trace.log) reproduces the player stopping at72.02/27.35 near junction64/32. The blockers are Gio/person-12's Packard, followed by Elena's Packard; these are actual junction reservations, not a lost walking animation. Existing waiting[] only meant an actor not admitted to the street, so it could be empty while a visible actor yielded.

StreetTraffic now returns the actual blockedBy ID for a blocked admission, connector or moving route step, without changing occupancy rules, routes, saved time or movement speed. A natural pause at the requested timeline fraction is not labeled blocked. City3D exposes trafficBlocks for inspection and reports state changes to the HUD. During a travel obstruction the clock reads 'Waiting for a clear path' and the journey panel explains the wait; both revert when it clears. This wording also covers parked vehicles or scene reservations without falsely calling them moving traffic. Skip remains presentation-only.

All294 existing frontend tests passed11.452s (.runtime/traffic-wait-tests.log). Added regression coverage for a visible actor stopping behind a parked Ford without overlap, reporting the correct obstacle, resuming after its removal and not misreporting a timeline hold. All25 traffic tests including the addition pass.180s (.runtime/traffic-wait-focused.log). Final build2.27s (.runtime/traffic-wait-build.log), existing chunk warning unchanged. Caffeinate87182 confirmed live; diff check clean.

Remaining pacing defect confirmed by the deterministic replay: at the nominal68.833s/1220 endpoint the player has only reached about68% after yielding. The current shared travel clock reaches the committed endpoint before collision-delayed walking completes. The feedback change does not solve that; collision-aware travel scheduling and full historical event timing remain required, along with the broader production goal.

CUA280 on fresh .runtime/traffic-wait-20260914.sqlite3 port8915/session93924: committed Saint Agnes→Monarch. At displayed1205.153 trafficBlocks identified player blockedBy elena; HUD showed 'Waiting for a clear path' and the journey explanation, screenshot inspected. At1218.810 blocks[] and the moving player66.644/68.65, both wait messages cleared; revision remained1. Main campaign untouched. Browser verifies stop/resume feedback, not resolution of the documented endpoint timing defect.

### 2026-09-14 — Reserve displayed travel time for delayed walking

Prior2856913 was concrete progress. Re-read the full amended objective and addressed its confirmed endpoint pacing defect. advanceJourneyClock limits the shared clock's rate using both nominal duration and the player's remaining route distance/physical speed. During a yield, the clock and NPCs continue advancing, avoiding a stopped-clock traffic deadlock; delays leave displayed time for walking afterward. Arrival completes the clock; skipped/reduced-motion journeys retain their existing immediate reconciliation. The player keeps the same physical speed/collision limits. No saved duration or Go outcome changes. API.md documents the presentation boundary.

Regression simulates a real cornered walking route against crossing Packard traffic at30/60/144FPS and1×/4×. It verifies positive clock movement during actual obstruction, no overlap, no premature clock completion, successful physical arrival, and completion-time agreement within1s across frame/speed settings. Separate checks cover unobstructed pacing and paused/negative elapsed input. All297 frontend tests pass11.362s (.runtime/journey-clock-tests.log). Build passes2.19s (.runtime/journey-clock-build.log), known chunk warning unchanged.

The original public QA snapshot replay (.runtime/trace-walk-adaptive.mjs, .runtime/walk-adaptive-trace.log) now arrives after96.217s with player progress1 and clock.9999999 immediately before arrival, versus the old clock finishing at68.833s with only68% of the route rendered. Fixed the diagnostic harness to clamp the player's requested progress at1, matching the actual renderer, before recording that result.

CUA281, fresh .runtime/journey-clock-20260914.sqlite3 on8916/session25013: actual Saint Agnes→Monarch25-minute command. Departure1195.088 had clock.003536/player.003430 and matching lighting minute. At1203.625, player.139621 was blocked by person-12 while clock.345001 still advanced. Changed the real travel speed control1×→4×; shared clock continued monotonically. Browser then recorded physical arrival progress1, endpoint1220, elapsed43394ms; HUD20:20 paused, lighting1220, revision1. Main campaign untouched. Caffeinate87182 live; diff check clean. Individual NPC collision-delay reconciliation, fully blocked route recovery, full historical events/weather and the broad production scope remain open.

### 2026-09-14 — Complete delayed NPC legs in order during travel

Prioreda7cff was concrete progress. Re-read the amended goal. A replay of the ordinary dusk fixture showed its four completed NPC legs disappearing only millimeters from their endpoints; it did not prove a large pop on that route. Code inspection nevertheless confirmed that streetAt removed actors solely by the recorded deadline and could replace a delayed leg with an already-progressed successor.

Added a per-presentation StreetPlayback queue. Completed observed legs stay requested at progress1 until their physical actor arrives. Successors wait for that handoff and their first traffic request uses the observed starting fraction, avoiding a jump down the next road when their scheduled departure has passed. Partial observed legs keep their recorded endpoint fraction. This changes no Go time or outcomes. API.md explicitly records that player-arrival/Skip/snapshot reconciliation still clears the presentation queue; retaining unfinished NPCs across that boundary remains open.

All299 then-existing frontend tests pass11.410s (.runtime/npc-leg-tests.log). Added a real StreetTraffic collision/queued-successor test with a parked obstacle, verifying continuous movement through both legs and no teleport at the delayed handoff; its first fixture placed the obstacle too far away to be reached within180 frames, corrected to107/16 before final all8 playback tests passed.217s (.runtime/npc-leg-focused.log). Final build2.18s (.runtime/npc-leg-build.log), known chunk warning unchanged. No new gameplay fixture code was needed; used the existing city3d-traffic scenario.

CUA282 on .runtime/npc-leg-20260914.sqlite3 port8917/session86898: actual tailor→Mariner72-minute journey; inspected traffic queues and final arrival, but missed the NPC deadline on the first observation. Repeated with fresh .runtime/npc-leg-watch-20260914.sqlite3, port8918/session68253, CUA283 and continuous bounded read-only samples. All12 NPC records finish at minute665. At665.047, eleven delayed cars remained visible, including editor at130.078/-1.6. At666.360 Elena had physically arrived and disappeared; Leo disappeared by667.711, Mara by669.028, mayor by670.342. Other cars continued smoothly throughout, including editor at105.169/-1.6 at670.999. This proves retention beyond the deadline and individual completion during the active journey. It also exposes the remaining final reconciliation gap: several cars are still far from arrival when the player's committed endpoint672 is reached. Main campaign untouched. Diff check clean. Full persistent NPC handoff, historical events/weather, high-quality assets and the broad production goal remain open.

### 2026-09-14 — Preserve NPC catch-up after normal player arrival

Prior81b5d24 was concrete progress. Re-read the amended goal and fixed its browser-confirmed final reconciliation gap. A naturally completed, animated journey now retains its StreetPlayback queue, actor instances and traffic occupancy, owned by world/life/revision. At the fixed committed minute, delayed NPCs finish observed legs and disappear individually on arrival; partial records stop at their recorded fraction. No new simulation advances while the clock is paused. Skip, new authoritative state, new scene or disabled motion discards the remainder and reconciles. Actor assignment now clears obsolete timelineProgress/legKey so subsequent ordinary snapshots cannot inherit stale per-leg requests. API.md updated coherently.

All300 frontend tests pass11.583s (.runtime/npc-remainder-tests.log). Build2.40s (.runtime/npc-remainder-build.log), known chunk warning unchanged. Existing queued-leg tests cover physical progression at a fixed minute; browser checks cover the React/renderer lifecycle transitions added here.

CUA284, fresh .runtime/npc-remainder-20260914.sqlite3, port8919/session6094: actual tailor→Mariner72-minute journey. Before arrival at670.891 there were7 NPC cars, editor105.598/-1.6. At player arrival672 (progress1, elapsed23562ms), settlingStreet true retained4 then3 cars while editor moved84.606→63.584 at fixed672/revision1. Later all NPCs reached their doors, count0 and settlingStreet false; HUD remained11:12 paused, revision1. This verifies continuation without an additional command or clock advance.

CUA285, fresh .runtime/npc-remainder-interrupt-20260914.sqlite3, port8920/session60543: caught5 pending cars at672/revision1, then issued a real new travel to Saint Agnes. Revision2 at displayed672.754 had settlingStreet false and no stale queue actors. Explicit Skip finished at681/revision2, no journey and no settling queue. Main campaign untouched. Disabled-motion and new-scene invalidation are implemented but not separately browser-exercised in this pass. Caffeinate87182 live; diff check clean. Fully blocked route recovery, rich historical event/condition playback, assets and the remaining production scope remain open.

### 2026-09-14 — Verify motion-off handoff and clarify settings controls

Priorc84a83f was concrete progress. Re-read the amended goal. Completed the pending disabled-motion browser gate using fresh .runtime/npc-motion-off-20260914.sqlite3 on8921/session68145. CUA286: actual tailor→Mariner journey, caught5 remaining NPC cars at672/revision1 with settlingStreet true. Opened Settings and disabled City animation: count0, settling false, minute672/revision1 unchanged. Re-enabled it: count0 and no revived queue, same saved minute/revision. Restored the original enabled animation preference. Main campaign untouched.

The existing setting called 'Scenes' also controls traffic and travel, so renamed it City animation and accurately described instant destinations/results when disabled. Added meaningful accessible switch names for City animation, Sound effects and Voices; previously the latter controls only announced On/Off. CUA287 verified all three final names and preserved states, and the revised settings layout screenshot was inspected. Updated debug preview guidance to name City animation only when it is disabled, instead of always saying scenes must be enabled.

Final build2.50s (.runtime/npc-motion-off-build.log), known chunk warning unchanged. No redundant test suite for copy/accessibility-name changes; the underlying motion handoff was exercised through actual browser controls. Final debug hint copy was source-checked and built after the settings screenshot. Caffeinate87182 live; diff check clean. New-scene invalidation remains an unexercised browser gate; broad assets, event/response fidelity, historical playback and the full production objective remain unfinished.

### 2026-09-14 — Complete the scene takeover browser gate

Priorb5c30f3 was concrete progress. Re-read the amended goal and completed the remaining new-scene invalidation check without changing implementation. CUA288, fresh .runtime/npc-scene-handoff-20260914.sqlite3 on8922/session66471: actual tailor→Mariner journey, caught5 pending NPC cars at672/revision1 with settlingStreet true. Started a private gunfight preview at the selected address. After.9s, settlingStreet false, only player/player-car remained in the ordinary actor set, and exactly the preview gunfight was active. Stop preview cleared effects and did not revive the prior NPC queue. Minute672/revision1 stayed unchanged throughout the preview. Main campaign untouched.

No build/test rerun for a verification-only gate; the tested distribution is the already-built b5c30f3 frontend. Caffeinate87182 live. Existing unrelated core/mugging.go, core/robbery.go, core/armed.go and sim/died_test.go remain untouched. The pending NPC handoff gates are now exercised, so stop repeating that QA absent a new concern. Next implementation priority: the requested building drive-by action, which is still absent from core command dispatch (incendiary exists). Broader art quality, response choreography, historical event playback and campaign acceptance remain unfinished; this is not whole-goal completion.

### 2026-09-14 — Building drive-by simulation foundation

The previous accounts acknowledgment was no progress: source and earlier browser evidence already showed the requested default breakdown. Re-read the complete amended objective and began the missing building drive-by feature. Added a core resolution with distinct player shooter/first-crew driver, existing delegation eligibility plus physical co-location, operational player car/fuel, equipped firearm, property/funds/life/custody checks, bounded weapon-dependent damage and supply loss, attention and actual-owner retaliation. Optional drive_by cue data captures the driver's authoritative identity, car tier/name and condition before/after; attacker captures the player's event-time weapon. Documented the staged contract in API.md and matching optional TypeScript field.

Focused core regression passes (.runtime/building-driveby-regression.log, .305s):96 seeded gun/car combinations, invalid attempts leave serialized world unchanged, dead/held/absent/travelling/busy/disloyal drivers, invalid targets, fuel/condition/heat limits, unchanged NPCs/staff/bankroll/other property, no fire, owner-specific retaliation and saved LastResult replay facts surviving later equipment/name/condition changes. Related arson/fuel/delegation/witness tests also pass. Test development caught an incorrect seed type and the transient VisualCues save assumption; replay test now uses the actual LastResult persistence path. Inspection also established NPC.Hurt means their damaged car, not personal injury; no erroneous injury gate remains, and a test ensures a crew member's damaged/empty car cannot block driving the player's operational car.

TypeScript/Vite build passes (.runtime/building-driveby-build.log,2.40s; existing chunk warning). No browser behavior changed or browser acceptance claimed in this core step. No action listing, dispatch or workshop kind is exposed yet: the moving car, separately seated driver/shooter, weapon-specific fire and facade hits, escape, footprint/traffic protection, damage timing and delayed newspaper still need implementation and real browser verification before the feature is playable. Initial balance is provisional. Main campaign and unrelated dirty files untouched; caffeinate87182 confirmed live. Broad production goal remains active.

### 2026-09-14 — Moving drive-by cast and corrected car doors

Prior91e94ed was concrete core progress. Re-read the full amended objective and implemented CityBuildingDriveBy as a reusable browser scene cast: actual two-seat car anchors, full-size driver at wheel grips, passenger pelvis planted in its cushion, forward-facing legs independent of torso turn, open passenger window, gun-specific firing beats/recoil/pump movement, facade-directed aim, rolling wheels and a continuously paced approach/slow pass/accelerating exit. It covers6.2 seconds and19.2 metres, from localx10 to−9.2. The cast takes supplied models and a target; it does not calculate property damage or advance the game clock.

Browser close-up revealed existing car door panels stopped above the footwell, exposing shoes through a large gap. Extended full-height door skins and upholstery in the Blender source and re-exported Ford, Hudson, Packard and police via .venv-blender/bin/python tools/export_city3d.py --only=cars (.runtime/driveby-cars-export.log). Geometry bounds/manifest remain coherent; other models were not regenerated. These repairs affect ordinary city cars too.

Added tools/building-driveby-review.html, an isolated motion study served by the existing8897 Vite process. CUA289 inspected Ford/revolver and shotgun poses, then Packard/woman/Thompson. Corrected low gun placement, unreachable long-gun support grips, revolver rear-window-pillar interference, and rotated hip origins that could cross the legs. The study now sizes its canvas on render, fixing the in-app browser's initially zero viewport. Browser Play samples:1.9972s/x4.0056,3.7611s/x.4778,5.5179s/x−5.4951,6.2s/x−9.2; passenger returns forward during escape, driver stays at wheel. No campaign API or save involved.

Final frontend suite319/319 passes11.654s (.runtime/driveby-motion-final-tests.log); TypeScript/Vite build2.27s passes (.runtime/driveby-motion-build.log, existing chunk warning). Added19 checks cover continuous bounded velocity at30/60/144FPS and18 car/rig/gun combinations using actual GLBs: unit-scale occupants, cushion/hip anchors, legs within cabin, hands reaching weapon grips and barrels crossing the open window clear of sill/roof/pillars at every shot. Existing car seating and scene tests remain passing. This does not prove all body/stock self-collision or final art fidelity.

Still not a playable action: integrate the cast into City3D with a reserved swept road footprint, ordinary actor/car suppression and restoration, facade impact visuals, weapon/car audio, captured condition timing, framing and delayed newspaper; then expose command/action/debug preview and browser-test actual committed results. No whole-city overlap, performance or production-quality claim is made by this isolated study. Unrelated armed-resistance files untouched; caffeinate87182 live. Broad goal remains active.

### 2026-09-14 — Drive-by road reservation and traffic clearance

Prior2c692f0 was progress. Read the full amended goal again and added a driveby-building scene slot on the actual near road lane, with a26m by3.6m swept footprint centered to cover the complete car/passenger/gun movement. There is one road candidate and no pavement fallback. Existing scene admission can wait for it to clear; existing pending-reservation traffic rules let current occupants leave while holding new arrivals.

Extended the actual-GLB checks to encompass the entire visible cast and car in that footprint at entry, firing beats and exit across all18 car/rig/weapon combinations. Added a traffic integration check proving a car already inside leaves, an incoming car stays outside, scene admission becomes possible, and traffic resumes after release. A separate placement check confirms an occupied near lane refuses admission while opposite-lane cars and sidewalk pedestrians do not overlap the reservation. Full321-test frontend suite passes11.798s (.runtime/driveby-clearance-tests.log); TypeScript/Vite build passes (.runtime/driveby-clearance-build.log), existing chunk warning only.

CUA290 on isolated8897 study reviewed the longest car with an outline drawn from the same reservation function. The initial study framing clipped the outline; widened its playback camera and visually confirmed the complete entry/exit footprint fits the viewport. This is a study camera, not final city framing. No campaign reads/writes. Caffeinate87182 live; unrelated files untouched.

Next: connect this reserved cast to City3D and its pending traffic mechanism, including normal cast replacement/restoration, facade hits and sound, condition reveal and news sequencing, then expose/test the real command. The new slot is not yet requested by live city effects, and the full production goal remains unfinished.

### 2026-09-14 — Preserve scene orientation during traffic admission

Priore6b4e0c was progress. Re-read the full amended objective and inspected City3D integration. Found a shared clearance defect before connecting the drive-by: scene reservations enter StreetTraffic as a single-point route, but onRoute discarded that point's authored heading and returned zero. The drive-by's east-west26m sweep would therefore become north-south after admission; rotated accident/planter reservations were also affected. Point now carries an optional heading, preserved for a single stationary point; moving routes still derive heading from their path.

Added an admission-level regression using the exact scene request shape, proving the drive-by retains its orientation, blocks a car near the actual end of its sweep, does not claim a point perpendicular to the road, and holds incoming traffic outside the admitted footprint. This closes a gap in the earlier pending-only checks. Full322-test frontend suite passes11.353s (.runtime/scene-reservation-orientation-tests.log); final TypeScript/Vite build passes (.runtime/scene-reservation-orientation-build.log), existing chunk warning only.

Added authored/admitted reservation poses to existing city-debug effect diagnostics. CUA291 on existing isolated8922 preview: actual Explosion · premature scene at Saint Agnes showed authored and traffic-admitted pose bothx75.55,z38.35,heading1.5707963267948966, stagedtrue at1.6388s. Stop preview removed all effects; campaign remained minute672/revision1. Main save untouched. Caffeinate87182 confirmed live. This turn fixes the shared admission defect; drive-by City3D cast/effects/command integration remains the next implementation step, not completed by this regression check. Broad production goal remains active.

### 2026-09-14 — Drive-by integrated into city preview playback

Prior7111b3e was progress. Read the complete amended objective and connected CityBuildingDriveBy to City3D effect creation, road admission/pending reservations, vehicle surface height, full sweep framing, gun-specific audio/flash/smoke and recoil, facade dust, vehicle approach/departure audio, watched-character cutaways, duration and cleanup. The cast consumes captured driver/car/weapon data. Ordinary player/car and pedestrian driver are excluded while their scene cast owns the space; a driver's separate parked car is retained. Added Building drive-by to the private debug scene menu with distinct preview driver/shooter and Packard/Thompson; it does not change campaign property condition.

CUA292 on existing isolated8922: admitted east-west reservationx80.4,z34,headingπ/2, car on actual road atz33.6/y−.115. Playback samples1.743s/0shots,3.2498s/5shots,4.7568s/8shots, then no effects after6.2s. CUA293 reviewed final dust/vehicle-audio build at1.9443s/0shots and3.0553s/4shots. Stop preview cleared all effects and restored ordinary player/player-car; minute672/revision1 unchanged. Screenshot inspected city framing and moving car. Audio scheduling counters are evidence of triggered sounds, not a listening-quality acceptance. The final parked-driver-car retention refinement was built after this private-cast check; real crew/car suppression still requires a committed-action fixture.

Full frontend suite322/322 passes11.641s (.runtime/driveby-city-tests.log). Additional focused preview test passes with the existing seven (.runtime/driveby-city-preview-tests.log), verifying distinct cast, captured equipment and unchanged campaign condition. Final TypeScript/Vite build passes (.runtime/driveby-city-final-build.log), existing chunk warning only. API.md now accurately distinguishes playable city preview from still-unexposed command. No main-save changes; unrelated armed-resistance files untouched; caffeinate87182 live.

Remaining before gameplay exposure: choose verified real facade hit points for all target building models, reveal captured property condition at gunfire, wire the transactional10-minute action, ensure recorded-result newspaper waits for completion, and test actual player/crew/car replacement/restoration including Skip/motion-off. The current dust is a simple effect and does not meet the full requested impact fidelity. Preview currently uses one gun/car combination; other variants have model/geometry tests but need city coverage. Broader assets, interior violence, police response, historical replay and campaign acceptance remain unfinished. Goal stays active.

### 2026-09-14 — Drive-by damage revealed with firing

Prior51c0e4b was progress. Read the full amended goal and connected captured drive-by condition to material wear and glazing. A pure presentation helper distributes the recorded loss across the weapon's firing beats, never rolling more damage. Pending/pre-fire scenes hold their captured starting condition. Once effects finish or are cancelled, the current snapshot's condition is restored; real commands will retain their committed damage while private previews restore the campaign. Material uniforms update only when the displayed condition changes, not through every mesh on every frame. Debug diagnostics expose committed/rendered condition separately.

Private drive-by previews now derive their starting condition from the actual target rather than resetting damaged buildings to100. Tests cover all three weapon schedules, low remaining condition, zero damage, exact endpoints and preservation of an already-damaged preview target. Full326-test suite passes11.466s (.runtime/driveby-damage-tests.log); additional final preview assertions pass (.runtime/driveby-damage-preview-tests.log); TypeScript/Vite build passes (.runtime/driveby-damage-build.log), existing chunk warning only.

CUA294 on isolated8922: Saint Agnes rendered100 at1.6944s/0shots,86 at3.1527s/4shots,72 at4.611s/8shots and through escape, then100 after natural completion. Separate cancellation after rendered72 immediately restored100. Committed condition100, minute672/revision1 unchanged. Main campaign untouched; caffeinate87182 live. This verifies private preview sequencing, not yet a committed attack/newspaper. Next required work remains verified facade hits across buildings and real action dispatch/acceptance with player/crew/car handoff and news after the scene. Broader production goal remains active.

### 2026-09-14 — Drive-by shots hit authored surfaces

Priora9929e3 was progress. Re-read the complete amended objective and replaced bounding-box aim points with first-visible-surface raycasts against actual building meshes. Hidden geometry variants and invisible materials are ignored. Scene admission requires a reachable authored target; impact dust uses the first surface from the actual muzzle sampled at the firing beat, including after a delayed render frame. A miss shows no invented dust impact. There are at most five target-search rays at admission and one impact ray per newly observed shot, rather than per-particle/frame mesh searches.

Actual GLB coverage across all19 active building models and all three weapons verifies every firing beat hits geometry and clears the passenger window. This exposed a small revolver rear-pillar collision for deeper industrial targets; moved its grip3cm forward relative to the opening and rechecked hand reach/default cast clearances. Added explicit hidden-variant/invisible-material/miss tests. Full346-test suite passes11.525s (.runtime/driveby-facade-final-tests.log); TypeScript/Vite build passes (.runtime/driveby-facade-final-build.log), existing chunk warning only.

CUA295 on isolated8922: Kessler's drive-by at3.2568s/5shots hit actualx144,y1.685,z17.3900; Saint Agnes at3.054s/4shots hitx80,y1.685,z41.8975. Authored/admitted road reservations matched, car stayed atroadheight−.115. Stop restored zero effects; campaign stayed672/revision1. Industrial screenshot inspected. Its canopy can still obscure small dust puffs, and the current simplified effects/art are not final impact fidelity. No main-save changes; unrelated armed-resistance files untouched; caffeinate87182 live.

Next implementation is the actual10-minute command/action and recorded newspaper sequencing, with a fresh isolated player/crew/car fixture verifying normal completion, replay, Skip/motion-off, expense/fuel/damage persistence and actor restoration. The drive-by is still debug-only; do not repeat completed geometry gates absent a new concern. Broad production goal remains active.

### 2026-09-14 — Playable building drive-by and recorded newspaper

Prioreca3384 was progress. Read the amended objective and wired driveby-building into normal action projection, The Street grouping, command dispatch and witness gravity. The action discloses35 dollars/petrol, equipped gun, available co-located first crew member, operational own car, damage and24 attention. Resolution pays/burns fuel once; ordinary command timing advances10 minutes and persists the scene in LastResult. Updated API.md to describe the actual playable contract rather than the earlier staged foundation.

Added a meaningful command-path regression proving one payment/fuel burn,10-minute advancement, captured event minute/condition and save retention, plus stale-revision rejection. Focused drive-by/arson/witness/grouping regressions pass1.032s (.runtime/driveby-command-tests.log). No frontend code changed or redundant frontend suite/build: the previously verified346-test distribution is used by the new Go server.

Added reproducible building-driveby QA scenario: go run ./cmd/qa-fixture .runtime/building-driveby-live-20260914.sqlite3 building-driveby. Fresh player/Leo at The Monarch with Packard/Thompson,30 fuel and20000 cash. Built .runtime/building-driveby-server and served only that isolated save on8923/session46012. Initial launch used a host:port value and failed terminally because this CLI requires :port; corrected to -addr :8923. Main campaign untouched.

CUA296 entered The Monarch and clicked the real Drive past and shoot up The Monarch action. Scene automatically returned to city. At1.6186s/0shots committed condition74/rendered100, at3.0769s/4shots rendered87, at4.5351s/8shots rendered74. Ordinary player/player-car were absent during cast playback and restored afterward. After escape the newspaper appeared with GUNFIRE AT THE MONARCH and the matching passing-car attack text. End state:610/revision1,19965 cash,24 attention,29 fuel, property condition74, car condition100. Public saved cue identifies Alex Varga/weapon3, Leo Carver and Packard, before100/after74, event minute600. No extra revision/time from presentation.

Normal committed flow is now exercised. Next acceptance: reload/replay the saved result, Skip and City-animation-off behavior on explicitly isolated saves, ensuring no extra cost/time and consistent damaged building/actor restoration. Broader weapon/car combinations have geometry/core tests but need campaign coverage; NPC chronology during non-travel actions, richer effects, visual quality and whole-campaign production acceptance remain unfinished. Caffeinate87182 live; unrelated armed-resistance files untouched. Goal remains active.

### 2026-09-14 — Preserve event newspapers with animation disabled

Read the complete amended objective. The preceding daily-accounts turn verified existing behavior but made no implementation progress; resumed the outstanding drive-by presentation gate. Inspection of the live isolated8924 fixture confirmed the attack committed at610/revision1 but remained inside The Monarch without automatically showing its newspaper. The command handler incorrectly gated the entire result flow on motion, and turning motion off discarded the active cue.

Changed main.tsx to retain the recorded cue and mark it finished immediately when City animation is off, both for new actions and manual replay. Turning animation off mid-scene likewise completes presentation instead of discarding the newspaper. Normal animated scenes keep their existing completion timing. No simulation/API payload changes. TypeScript/Vite build passes (.runtime/motion-off-news-build.log).

Fresh isolated fixture .runtime/motion-off-news-fix-20260914.sqlite3 on8925/session44311, CUA299: disabled City animation through Settings, entered The Monarch and executed the real drive-by. Newspaper immediately showed GUNFIRE AT THE MONARCH; zero effects and committed/rendered condition74. Motion-off replay opened the same article immediately. Enabled animation, replayed, and disabled it mid-scene: newspaper retained, zero effects,610/revision1. Re-enabled animation afterward: no resurrected scene, ordinary player/player-car restored. Read-only public snapshots on8923 (prior replay/Skip),8924 (original motion-off defect),8925 (fix) all confirm cash19965/fuel29/heat24/condition74/minute610/revision1. Main save untouched, caffeinate87182 live, unrelated armed-resistance changes left alone.

These gates verify presentation does not spend again or lose the article under disabled motion. They do not establish broad visual fidelity, all weapon/car campaign combinations, non-travel NPC chronology, police cleanup animation or whole-campaign production acceptance. The full goal remains active.

### 2026-09-14 — Behind-the-head shooter visibly withdraws

Prior f303441 was implementation progress. Re-read the complete amended objective and inspected assassination staging. The existing close-quarters variant is a frontal beating, not a wire attack; wire/interior variants remain unfinished. The behind-the-head shooter previously stopped beside the victim until the cast disappeared. Added a lowered-weapon turn followed by a1.35m/s return along the already reserved4.07m approach. Turning finishes before walking starts. Its scene now lasts8.5 seconds so withdrawal completes; close-quarters retains6.5 seconds and other armed variants retain8. No event outcome, shot timing, victim pose or traffic footprint changes.

Added regression coverage for facing/lowered gun before departure, continuous bounded return speed and completion at the path origin. Existing actual person/woman/revolver GLB tests now sweep the longer execution through every frame, verifying ground, lamp/reservation clearance, hand contact and exact persistent-body handoff. All347 frontend tests pass12.194s (.runtime/execution-withdrawal-tests.log); TypeScript/Vite build passes (.runtime/execution-withdrawal-build.log).

CUA300 on isolated8925 private Assassination · shot from behind preview: attacker reachedx112.07 by4.479s, withdrew tox111.275 at5.889s and109.363 at7.305s; effects cleared after8.5 seconds. Minute610/revision1 unchanged. Inspected a second playback screenshot during retreat: walking shooter and prone victim remain framed on the forecourt, gun lowered. Characters still have the acknowledged simplified/blocky look; this is a choreography correction, not art-quality acceptance. Ordinary nearby bystanders, eventual cast-to-world actor handoff, wire/interior attacks, richer police response and broader production acceptance remain unfinished. Caffeinate87182 remains live; main save and unrelated armed-resistance edits untouched.

### 2026-09-14 — Remove duplicate assassination actors and retain player departure stance

Prior f006290 was progress. Read the full amended objective and inspected the ordinary-city/scene handoff. Staged assassinations did not exclude their attacker from ordinary traffic/rendering, unlike arson and drive-by casts. They now suppress matching pedestrian attacker/victim IDs during playback; parked cars retain their independent presence. A real, stationary pedestrian player at the recorded address adopts the assassination's final position and facing through StreetTraffic.adoptFrontage. Private previews, vehicles, moving players, different addresses and dead players do not receive that stance. The existing collision-checked frontage connector handles the next journey. Scene positions remain unsaved presentation state.

Build passes (.runtime/assassin-handoff-build.log).32 focused assassination and traffic tests pass1.419s (.runtime/assassin-handoff-tests.log), including actual GLB cast geometry/body handoff, continuous frontage departure at30/60/144Hz, blocker yielding and rejection of cars/unrelated frontages. No new simulation rules or tests that merely mirror the UI condition.

Fresh isolated strike-revolver fixture .runtime/assassin-handoff-20260914.sqlite3, server8926/session78290, CUA301. Replayed the actual recorded player/Mara assassination: ordinary player absent throughout approach, shot and withdrawal; at completion player restored exactlyx76,y.2,z38.35, newspaper present and minute480/revision0 unchanged. Closed paper and travelled to The Mariner. Subsequent samples connected continuously fromx76,z38.35 to the pavement z36.65, thenx77.101,78.289,79.476 before joining the route atx80, with no traffic blocks. This command advanced the committed state to495/revision1; the debug presentation minute also showed495 during connector movement, which needs a separate follow-up against the requested paced walking clock (do not call all travel timing complete). Skip/revision-interruption coverage for this new handoff and NPC/vehicle actor handoffs remain to exercise. Caffeinate87182 live; main save and unrelated armed-resistance edits untouched. Broad production goal remains active.

### 2026-09-14 — Correct clock diagnosis; verify skipped assassination departure

Prior1a212f2 was implementation progress. Read the full amended objective and revalidated its reported clock concern against current code. The prior browser probe read data-presentation.minute, which is the committed simulation endpoint, rather than streetMinute or the visible HUD. That observation did not demonstrate an actual clock-pacing bug. No timing-code change was warranted.

Fresh isolated strike-revolver fixture .runtime/assassin-clock-20260914.sqlite3 on8927/session1979, CUA302: replayed and skipped before the shot. Player restoredx76,y.2,z38.35, zero effects, newspaper present, committed/displayed minute480. Closed newspaper and travelled to The Mariner. While walking the frontage connector at zero route progress, displayed streetMinute advanced480.348,480.688,481.023,481.347; HUD08:00 then08:01. Route progress then began at.00496 with streetMinute481.668. Later observed08:08/street488.208 at route progress.493. Arrival took31283ms, endedx48,y.2,z36.65, HUD08:15/street495/committed495/revision1, with journey controls cleared. Thus the connector and remaining route share the paced display clock; only Go's committed result holds the endpoint from the outset.

This completes the new Skip handoff gate and resolves the suspected clock defect as an incorrect diagnostic-field interpretation. Added the assassination handoff behavior to API.md. No runtime edits or repeated tests/build required: unchanged1a212f2 distribution used throughout. Next useful implementation/acceptance work concerns NPC/vehicle cast handoffs, interruption by a new revision, wire/interior assassination staging and the remaining larger production scope, not reworking already verified clock pacing. Caffeinate87182 live; main save and unrelated armed-resistance work untouched. Goal remains active.

### 2026-09-14 — First 3D poker scene and close card-table cameras

The new explicit visual-production request supersedes the historical handoff restriction. Added a procedural Three.js poker table with wood rail, felt, physical card stock, shared readable face textures, pot/contribution chips and public seat labels. Only API-exposed cards are drawn face up; folded seats remove their cards from the 3D table. The full public showdown remains in the text summary. Poker and blackjack take the full screen width, with table renders around60% viewport height and compact authoritative summaries. Removed the remaining760px blackjack panel cap. Enlarged card corner ranks. This is code-authored geometry, not a new Blender/AI asset.

Shared TableCamera provides pointer orbit, right-drag pan, wheel zoom, focused-canvas WASD/arrows, Q/E, +/- and Home/reset. Elevation/distance/pan are bounded; resize adjusts the framing without resetting a deliberately moved camera. Listeners and GPU resources are disposed on exit. Blackjack retains existing committed dealing/flip presentation. Poker currently updates cards immediately and still needs dealing choreography, seated cast and richer authored room detail. Procedural poker materials are a first scene and do not establish final noir fidelity.

Added reproducible `poker` qa-fixture scenario. Fresh .runtime/poker-3d-20260914.sqlite3 served8930/session38132; main save untouched. Browser played600 buy-in from preflop through flop/turn/river/showdown:10♥5♥, board6♥J♠J♦4♥2♣, pot357, final stack477, on-person cash4400. Camera Q/W changed coordinates while preserving cards; reload retained showdown. Inspected default1235x1051 and compact600x850 viewport; all five board cards and player hand framed, controls remain available by scrolling on compact screens. Temporary viewport override reset. An existing projection discrepancy was observed: NPC seat hand summaries say high card at showdown while the result correctly reports Magda's full house. This needs a separate core projection fix, not client-side hand evaluation.

Fresh .runtime/blackjack-camera-20260914.sqlite3 served8931/session94191. Browser entered The Blue Hour, sat, dealt50: playerK♣7♦, dealer9♥ plus hidden card. Inspected final full-width close camera, Q/W movement, reset and stand. Cash3950 before stand/4050 after committed win; saved-hand reload did not redeal. Camera diagnostics144 draw calls in the inspected blackjack frame, poker68–70 during play; these are not measured FPS acceptance. No production campaign was mutated. Caffeinate session24710 running during work.

Validation: TypeScript/Vite build passes (.runtime/card-room-build.log, existing chunk warning). All347 frontend tests pass12.415s (.runtime/card-room-tests.log). Focused core Poker/BackRoom/Sitting/Card tests pass26.557s (.runtime/card-room-core-tests.log). Pre-existing core/mugging.go, core/robbery.go, core/armed.go, sim/died_test.go and src/city3dAftermath.ts changes left untouched and excluded from this commit.

Next work: fix NPC showdown labels; deepen poker cast/materials/dealing; neighborhood-driven estate pricing and NPC deed exchange; burglary and present-at-home targeting; remaining modeled interiors, income progression and map housing expansion. Existing housing assignment/Mariner income/basic deed sales are foundations, not completion of the amended goal. Goal remains active.

### 2026-09-14 — Residential prices respond to neighborhood violence

Previous5fe46ab was implementation progress. Revalidated the amended goal and existing deed paths; Mariner/Cypress purchases and broker offers had static bases. Added saved district pressure from committed located violence, a60% value floor, and one percentage point per quiet game day recovery with fractional-day retention. Market projection/notice, Mariner acquisition premiums, Cypress standalone deed and buy-and-move action costs, readiness and broker offers share the index. UI reads consume no randomness or time. Legacy saves start neutral; pressure is world state across lives. Neighborhoods currently follow simulation districts, and the existing two traded properties remain the only player saleable deeds. NPC deed trading/new individual housing remain unfinished.

Focused regression suite passes (.runtime/neighborhood-prices-tests.log): local shock versus unaffected district, quiet-time recovery, JSON reload, partial-day accumulation, neutral legacy state, read purity, price floor/spread, action/market/actual-payment agreement on both Cypress purchase paths, and Mariner acquisition/sale accounting. TypeScript/Vite build passes (.runtime/neighborhood-prices-build.log).

Fresh .runtime/property-prices-20260914.sqlite3 building-driveby fixture served8932/session30117 using .runtime/property-prices-server. Actual POST driveby-building at club advanced600/revision0 to610/revision1. Mariner asking3600→3384, offer2340→2199, index100→94 while its condition stayed100; Cypress asking3500/offer2275/index100 unchanged. Browser inspected the Market and its explicit6% local-violence discount text. Only isolated save modified. Caffeinate processes16576/87182 confirmed running. Full Go suite running when this note was written; terminal result appended below before commit.

Additional gates: idempotent HTTP retry of the same attack retained revision1/minute610/index94, with no second market shock. Added and passed actual new_life command coverage proving neighborhood pressure survives the protagonist. Final focused price/housing suite0.230s. The broad run found one outdated action-availability fixture: no firearm/co-located driver in any of its three states made driveby-building permanently appear refused. Equipped its established state and placed its actual crew driver alongside the player; both broad offered-action checks then passed0.907s (.runtime/offered-equipped-tests.log), including actual execution. This is fixture coverage for the already implemented command; no readiness rule was relaxed.

### 2026-09-14 — Correct opponent poker showdown labels

Fixed the concrete prior browser defect: CardsDescription ranked only an NPC's two hole cards with the old five-card evaluator, so the winner's public label said high card despite a full-house result. It now uses BestOfSeven with the shared board, as settlement and the player's summary already do. The projection still exposes neither cards nor rank before Done. Regression reproduces2♥J♣ on6♥J♠J♦4♥2♣ and requires a full-house label; passes0.157s (.runtime/showdown-description-tests.log). No payout, deck, decision or visibility timing rule changed.

Full-suite handoff: `go test ./...` is still live as exec session58737, sim.test PID18967 verified at6m elapsed (CPU active). Do not restart it: poll that session and read .runtime/neighborhood-prices-go-all.log next turn. Completed packages cmd/apicheck, cmd/blackledger, cmd/playtest, cmd/simulate passed. Core completed154.390s with only the offered-action fixture failure described above, then its corrected tests passed. Store/sim terminal results remain unproven. The complete command predates the small showdown-label fix, which has its own passing focused regression. Broad goal remains active; next housing work is individual deed stock/NPC transactions, with burglary/home targeting and remaining interiors still pending. Unrelated armed-resistance/aftermath changes remain unstaged.

### 2026-09-14 — Apartment ownership, resale and funded NPC deed trading

Previous7c11122/f857aa0 were progress. Continued the amended goal by adding stable individual deeds within the existing64 Ashbury/48 Mercer residential places. This is subdivision of existing housing capacity, not new map buildings. Initial assignment/migration preserves homes, money and actor locations; residence move previews stay pure. Existing NPC owner-occupants cannot be evicted to accommodate a player move. Purchases end rent without granting the whole building; moving retains the deed; sale restores tenancy without displacement. Base1200/1800 prices use neighborhood pressure and65% broker spread. Asset sales do not earn progress.

NPCs can buy their rented flat with savings and buy deeds from cash-poor NPC owners. One daily funded transfer conserves NPC-to-NPC money and leaves current travel untouched. Rent goes to the deed holder from available tenant cash, with owner-occupants paying no rent. Books includes contracted apartment rent and owned-flat holdings. Numbered home addresses and owner-occupant labels appear in public people information. Market shows the player's rented/owned flats with prices, rent and an address link. Private NPC sellers are not offered to the player automatically. Added buy/sell apartment business grouping after initial browser revealed fallback Work placement.

Fresh apartments qa-fixture .runtime/apartments-20260914.sqlite3 on8933/session11505. Browser entered Mercer Court and bought current flat47:6000→4800 cash,25→0 daily housing,08:00→09:00. Market showed own deed/no rent/780 broker offer. Reload retained deed and4800 cash. Sold through real UI:5580 cash,25 rent,09:30, same home/address and the buy option returned. Main campaign untouched. The final grouping/reference-price refinements were compiled and covered by focused tests after that server build; its observed earlier action placement was Work and is not claimed as final grouping browser evidence.

Validation: final core/store apartment/housing/rent/books/property/migration/save/offered-action/group checks pass (.runtime/apartment-final-tests.log); TypeScript/Vite build passes (.runtime/apartment-build.log, existing chunk warning). Tests cover stable112-unit registry, read-only move planning, exact purchase/sale payment, retained deed after moving, duplicate-sale rejection, NPC self-purchase then money-conserving private transfer, rent paid to buyer from real cash, unpaid-rent cap/exactly-once collection, JSON persistence, no new-life inheritance, and blocking owner eviction. Broader campaign pacing, inherited deed distribution, personal apartment interiors, detached homes and NPC move-up preferences remain unfinished. The core state is more expressive but does not complete the whole housing/visual goal.

Resolved prior full-suite handle58737: terminal failure, no longer running. Core's only reported failure was the previously corrected offered-action fixture; store passed0.156s. Sim timed out at600.087s solely with pre-existing untracked sim/died_test.go TestWhatKillsThem still running (260 long campaigns). Its timeout is not a passing campaign result. That user-owned diagnostic was inspected but left unmodified; do not silently restart it. Targeted current checks replace no claim of full-suite success. Other pre-existing armed-resistance/aftermath files remain untouched.

Next substantial work: player/NPC homes rendered as actual private interiors, burglary and home-presence targeting, housing expansion beyond existing capacity and more NPC relocation choices, alongside poker cast/dealing/material quality. Keep the full goal active.

### 2026-09-14 — Finite household burglary and present-at-home attacks

Previousdaa3335 was progress. Added private household cash funded from actual NPC purses at daily settlement, capped60/day and1000 total, with daily idempotence and withdrawals for living costs. Apartment purchases and poor-seller checks include household cash so savings are not stranded. No new migration loot or public stash amounts. Added named45-minute burglary actions at known residents' home addresses: finite stash transfer, empty-home zero payout, lower odds in occupied homes,7/14 attention, injury/clothing damage/fatal risk on failure, identification-based grudges/family retaliation and committed news/robbery pressure. This does not steal wallets, other tenants' cash or deed-holder assets.

Added home_strike action for a physically present resident, using existing Strike resolution. Departed/travelling/held/absent targets fail its residence check. Optional strike.setting=home captures context before death for future interior choreography. Existing generic street attacks remain distinct. These commands are playable backend/UX progress; private-room entry, burglary animation and home-specific assassination rendering are not complete. Current city/newspaper playback remains in use, and its newspaper proprietor fallback portrait is not a household illustration.

Fresh burglary fixture .runtime/burglary-20260914.sqlite3 on8934/session24793. Player at Mercer Court, known Mara away at Saint Agnes, private stash180 and pocket90. Browser entered Mercer and found the named action under The Street despite Mara not being in the room. Actual action:1000→1180 cash,0→7 attention,08:00→08:45. Existing city sequence completed, then BREAK-IN AT MERCER COURT appeared with the committed break-in story. Main campaign untouched. The later integration allowing savings to fund apartments has focused coverage and was compiled after this burglary server; the observed theft mechanics were unchanged.

Validation: TypeScript/Vite build passes (.runtime/burglary-build.log); focused core/store/sim checks cover burglary, household funding/cap/daily repeat, finite payout, wallet/other-stash preservation, failed occupied entry/injury/grudge, private snapshot boundaries, JSON retention, household moving, home-target absence/travel/death rejection, committed home setting, apartment funding from savings, property pricing, migration/saves, offered/grouped actions, campaign reproducibility and policy hidden-threat isolation (.runtime/home-crime-final-tests.log; final results recorded below). Unknown/dead/no-home, low-health and crew exclusions retain ordinary action guards; broad art quality is not established by these core tests.

Next work remains private residential rooms and additional 3D interiors, playable home scene staging, expanded housing capacity/detached homes and NPC relocation preferences, plus richer poker/table cast/dealing/materials. Whole goal remains active. Pre-existing armed-resistance/aftermath changes and untracked long death diagnostic remain untouched. Caffeinate16576/87182 verified active during this turn.

Final focused results: core3.226s, store0.184s, sim8.723s, all pass. The previous full-suite timeout diagnostic was not rerun.

## September 14 — Cypress House drawing room and study

Added the estate's own Blender-authored interior, rather than a renamed lobby.
The model has staggered parquet, raised wall panels/cornices, sash windows and
curtains, fireplace/mantel/clock, bookcase with individual volumes, leather sofa,
velvet armchairs, rug, reading lamps, newspaper table, study desk/radio and drinks
sideboard. Only the new GLB/manifest entry was exported. Both walls cut away when
the camera passes outside them. Existing orbit, keyboard pan/zoom/reset and
reduced-motion entrance behavior remain available.

Seven public occupants fit five measured cushions and two clear standing bays;
the player has a separately tested entrance route. Overflow remains visible in
the complete person list. No new public-state contract or invented NPC occupancy.
This is a drawing room, not completion of every room in the estate or every city
interior. Fitted house upgrades and burglary/home-assassination choreography are
still outstanding, as are broader texture/character fidelity improvements.

Evidence:
- `npm run build`: pass (`.runtime/cypress-build.log`).
- `npm test`: all 348 pass (`.runtime/cypress-frontend-tests.log`), including a new
  actual-GLB test of both character rigs, cushion support, cast separation, clear
  standing bays and a 41-sample entrance sweep. Floor-board seams exposed a gap
  in the first geometry check; a continuous substrate now supports those seams.
- Isolated `cypress` QA fixture added. The first fixture omitted district access;
  corrected fixture `.runtime/cypress-interior-v2-20260914.sqlite3` serves 8936.
  Main `.runtime/campaign.sqlite3` was never changed for QA.
- Browser at 1280×720 and 820×740: room fits, eight rendered actors (seven public
  plus player), two overflow occupants explicitly counted and listed. Completed
  entrance, zoom 1→1.12 and Home reset, pointer orbit with wall cutaways, and
  direct click on seated Mara selected her actual action panel. Clock remained
  17:00 and cash $12,000 throughout presentation checks. Property upgrade/sale,
  rest and security commands remain visible in the corrected fixture.
- Observed default scene: 169 draw calls / 191,440 triangles with eight actors;
  this is a scene measurement, not a hardware performance acceptance claim.

Other requested interiors, richer poker cast/dealing, expanded residential map,
additional progression income and integrated campaign acceptance remain open.

## September 14 — poker cast and public card movement

Added seated poker participants using the existing locally authored rigs and
chairs, with public presence costumes, plus the player's foreground seat.
Adjusted the close camera after browser inspection found clipped opponent heads.
Added a pure public-state animation plan: round-robin initial deal, new board
cards, and physical-backed showdown flips. Unchanged cards remain settled;
folded hands leave the felt; new hands discard old exposed cards. Restored games
never replay a deal. Starting from the empty sitting passes an explicit initial
deal flag through the newly mounted hand component. Motion-off/reduced-motion
settles immediately; idle presentation does not force continuous rendering.

Evidence: `npm test` passed all 351 tests (`.runtime/poker-cast-tests.log`),
including new privacy, round-robin ordering, flop-only movement, reveal, folded
hand removal, new-hand and motion-off plan checks. `npm run build` passes
(`.runtime/poker-cast-build.log`). Fresh isolated poker fixture
`.runtime/poker-cast-20260914.sqlite3`, port 8937: browser played through flop,
turn, river and showdown, then dealt hand 2. Four actors followed the actual
roster change. During turn and showdown, `data-poker.dealing` was true; exposed
card counts stayed zero until showdown and returned to zero next hand. Reload
of hand 2 showed `dealing:false` with the saved cards. Desktop 1280×720 and compact
820×740 screenshots reviewed; player cards remain unobstructed, public hand
summary available, opponent heads included. No main-save QA mutation.

Still incomplete: a physical dealer-hand/chip-pushing performance, richer poker
room dressing, more expressive/anatomically detailed character models, all other
requested interiors, housing expansion and remaining progression/campaign work.

## September 14 — buy-to-let apartment progression

Players can now buy broker-owned rental investments without moving into each
flat first. At each apartment address, a bounded board offers three tenanted
flats and one vacancy in addition to the player's current home/holdings. The
board replenishes after a sale, while privately owned NPC homes remain protected.
The market screen exposes tenant, contracted income, price and district access;
it distinguishes buying one's home from a vacant investment. Purchases preserve
resident/address/journey and only transfer the deed. Existing rent collection
pays actual available tenant cash, not an invented guaranteed income. Existing
crime-sensitive valuations and broker resale spread continue to apply.

Capacity evidence: focused housing/Mercer tests pass, including the established
30-new-arrival scenario with no shortage or displacement. This pass adds a
usable property-income path, not additional building geometry or capacity.

Validation:
- Focused apartment/broker/housing/Mercer/property/save tests: core 2.067s, store
  .123s pass (`.runtime/apartment-investment-final-tests.log`). New command test
  checks purchase principal, preserved home and tenant location, collection of
  only $7 from a $7 purse, and tenant retention on sale. Listing test checks
  bounded/read-only ordering, replenishment and private-owner refusal.
- Frontend production build passes (`.runtime/apartment-investment-build.log`).
- Fresh isolated `.runtime/rental-investment-20260914.sqlite3` on port 8938:
  browser bought Mercer unit 1 for $1200; cash 6000→4800, time 08:00→09:00,
  player still rents unit 47, Zoltan Toth remains tenant of unit 1. Accounts
  show $12 contracted income; market shows the owned deed/$780 sale offer,
  and broker unit 4 replenishes the board. Vacant flats show no income.
- Main campaign save untouched. Full simulation campaign/balance acceptance,
  additional earning activities, additional residential districts as population
  requires, and remaining 3D interiors/action choreography are still open.

## September 14 — Ashbury Court entrance hall

Added a distinct Blender-authored residential hall: stone floor and pilasters,
64 numbered mailboxes, walnut reception counter/register/bell, brass lift gate,
upholstered bench, notices, plants and stone staircase. Browser review caught a
stair terminating against the rear wall; the final model has a supported landing
and open exit beyond the cutaway. Walls and attached fixtures cut away with orbit.
Public occupants use the authored bench and clear floor; a concierge role, when
present, gets the reception aisle. No invented decorative NPC or new capacity.

Validation: production build passes (`.runtime/ashbury-build.log`), complete
frontend suite passes 352 tests (`.runtime/ashbury-tests.log`). An additional
focused final check passes both Ashbury tests (`.runtime/ashbury-geometry-tests.log`):
both rigs remain separated, seats have actual cushion support, standing bays have
floor/clearance; all 12 stair treads, landing and upper opening checked by rays
against the exported GLB. The added stair test ran after the 352-test suite.

Fresh isolated `ashbury` fixture at `.runtime/ashbury-interior-20260914.sqlite3`,
port 8939. Browser reviewed 1280×720 and 820×740, completed entry with six public
occupants plus player, no omissions, and selected Mara directly in the canvas.
The registry still shows 16 tenants / 64 places; property/burglary actions remain
available. Cash6000/time08:00 unchanged through presentation checks. Observed
143 draw calls / 253,928 triangles for the seven-character hall. This is a scene
measurement, not a performance acceptance claim. Main campaign save untouched.

Remaining: private flat rooms, other commercial/civic/industrial interiors,
interior residential crime choreography, higher character/material fidelity,
additional income activities and integrated campaign/balance acceptance.

## September 14 — long-run housing progression defects

An unattended economy audit found zero NPC deeds across three 30-day runs.
Ordinary households hit the fixed $1000 savings ceiling below the normal $1800
Ashbury asking price. Replaced that ceiling for numbered-flat residents with
max($1000, local flat price + $500). Deposits remain actual purse transfers, capped
at $60/day; savings above a subsequently falling limit are not destroyed. This
allows earned savings to fund deeds and also makes older affluent households
worth more to a burglar, through actual accumulated cash rather than a loot roll.

Extending the simulation to 90 days exposed another real defect: later NPC
replacements had no homes despite vacancies (19/20/20 unhoused in seeds 7/27/61).
New-game population settlement worked, but later daily creation paths skipped it.
Added housing/flat settlement after daily officials/roles are filled and before
committed commands increment revision. Existing residences and current journeys
remain stable; no new stock was needed for these populations.

Final three 90-day runs: seed7 108 living, zero shortage, 8 living NPC-owned deeds;
seed27 103 living, zero shortage, 4 deeds; seed61 99 living, zero shortage, 4 deeds.
Recorded `.runtime/housing-audit-final.log`. These are unattended economy runs:
player observer starts with $100000 to avoid rent bankruptcy, initial tasks/plots/
contracts are cleared, and player event prompts are dismissed between real-clock
240-minute advances. NPC purses, wages, conflicts and purchases are not seeded
with extra money. This does not replace played campaign acceptance.

Added permanent three-seed 90-day regression checking actual NPC ownership,
zero shortage and building capacity, plus committed-command newcomer housing.
Focused housing/household/apartment/broker/property/save checks pass (core2.805s,
store.170s, `.runtime/housing-progression-tests.log`). No main-save QA mutation.
Broader visual fidelity, remaining interiors/indoor scenes, extra income activity
and full campaign/balance acceptance remain outstanding.

## September 14 — private apartment view

Added a locally authored furnished bedsit, reached from the player's current
Ashbury/Mercer home. The private scene contains the player only; the building
list remains available, and selecting a person returns to the public hall.
Changing home/place resets private-view state. The room has bed/linen/radio,
living furniture, dining chairs/table and a detailed enamel kitchenette. It is
currently a shared floor plan; differentiated unit layouts, fitted upgrades,
bathroom interior and private-home action choreography remain outstanding.

Browser inspection found the first return control under the absolute Back to
city button; final spacing clears it at desktop and compact widths. Private
caption removes person-picking advice and the action prompt explains returning
to the hall. No command is sent to enter or leave this local view.

Validation: build pass `.runtime/private-flat-build.log`; complete frontend suite
354/354 pass `.runtime/private-flat-tests.log`. After dining chairs were added,
final targeted GLB test passes `.runtime/private-flat-geometry-final.log`: both
player rigs have a supported, unobstructed 41-sample entrance path; private
staging rejects public occupants. Browser used existing isolated Ashbury fixture
8939 (no main save), reviewed 1280×720 and820×740. Observed private cast exactly
['player']; clicking Mara in building list returned to the Ashbury hall with
Mara selected. Cash6000/time08:00 unchanged. Remaining broader scope is active.

## September 14 — clean release integration check

Prepared a clean archive of f1cfbc0, excluding the pre-existing uncommitted
armed/mugging/robbery/aftermath work. Frontend build and all354 tests pass there.
A read-only SQLite backup of the main campaign was upgraded on isolated port8940:
Jamie Moretti, life11/revision2115/minute160095, complete player object preserved.
All2117 receipt contents match by SHA256, not only count. Migration14→19 adds
housing/presentation fields, Mercer Court and one required population role.
Main stored-state digest still matches its pre-QA digest. No listener was on8791.

Full clean Go package checks caught Apartments' `omitempty` tag violating the
project list contract; removed it. The apartment market API itself was already
an explicit list. Release integration continues after remaining checks finish.

The full clean core run additionally caught the fee guard's static call-site
count (43) lagging the new rental-investment offer (44). Updated that count;
the command-level investment test verifies the actual single principal debit.
The original clean core package completed in170.226s with this sole failure;
its corrected recheck and the remaining full simulation run are tracked below.

## September 14 — verified release running on8791

Release b4b7aee is now serving the main campaign. All clean packages passed:
cmd/blackledger corrected recheck1.103s, core corrected full recheck195.778s,
sim394.306s, store0.128s, cmd/simulate7.115s and cmd/playtest0.076s.
Frontend clean build and all354 tests passed. These checks excluded the existing
uncommitted armed/mugging/robbery/aftermath files and sim/died_test.go.

Go1.23's VCS detection ignored the nested worktree .git file and reported the
outer checkout's dirty state. The final binary instead comes from a clean local
clone at `.runtime/release-live-b4b7aee`; /api/health verifies revisionb4b7aee
and modified:false. Its frontend is the tested release build. Main port8791 had
no listener before startup; no running campaign process was replaced.

Before startup, made `.runtime/campaign-pre-b4b7aee.sqlite3` with SQLite backup
and verified integrity. Main startup performed the reviewed v14→19 migration.
Read-only verification confirms its complete resulting state equals the isolated
upgrade copy, including revision2115, player object and all2117 receipt contents.
Live state SHA256: db1471392f15b1bc3381f5456e169e6a1ba67eb4d000626f379190b352d1ec55.
Evidence: `.runtime/release-live-verification.json`, clean test/build logs and
`.runtime/release-live-identity.log`. No gameplay QA action used the main save.

This release integrates poker/table cameras, apartment ownership and trading,
crime-sensitive prices, actual household savings/burglary/home-strike rules,
rental investments, Cypress/Ashbury/private-flat interiors and population housing
settlement. It does not close the goal: indoor assault/burglary choreography,
remaining commercial/civic/industrial interiors, broader income progression,
visual fidelity acceptance and a fresh coherent played campaign remain.

## September 14 — Fassano Meats interactive shop interior

Replaced the butcher's flat plate with original Blender-authored 3D shop dressing:
glazed wall tiles/borders, quarry floor, refrigerated display with glass/rails and
meat trays, end-grain block, knife rack, wrapping-paper stand, mechanical scale
with dial/ticks/needle, cold-room door and shop lettering. Source is
`tools/butcher_interior.py`; reproducible targeted export:
`.venv-blender/bin/python tools/export_city3d.py --only=interior-butcher`.
Full export includes it; other reviewed model files were not regenerated.
The model is1,224,432 bytes. No third-party assets or invented sale prices.

Actual public occupants fill nine customer positions and a reserved service
position for a present butcher/shopkeeper/clerk. The complete roster remains
available. The separate player entrance clears the fullest roster. Tests use
the actual GLB and both character rigs for floor support, furniture clearance,
non-overlap and sampled arrival. All355 frontend tests passed before adding the
final entrance check; both final butcher geometry tests pass, build passes.
Evidence: `.runtime/butcher-tests.log`, `butcher-geometry-final.log`,
`butcher-final-build.log` and `butcher-export.log`.

Isolated fixture `qa-fixture ... butcher`, port8941, browser verified1280x720
and820x740: room loads, camera zoom/orbit/reset work, clicking the rendered
butcher selects the correct public person/actions. The first preview opened
before the export reached dist and correctly showed the fallback; a fresh
completed build loads successfully. No main-save QA commands. Main8791 still
serves releasedb4b7aee; this shop is currently in the development preview.

Remaining: cold-room interior/work animations, richer surface wear and prop
variation, and the other commercial/civic/industrial rooms. This is one further
interactive destination, not completion of the full interiors/visual goal.

## September 14 — enlarged interiors and explicit room lighting

Added an enlarged room toggle to every authored 3D interior, retaining the same
renderer, camera and cast while increasing available screen height. Standard
view and canvas Escape restore the prior layout; actions remain below the scene.
Browser checked1280x720 and820x740 on isolated8941: expansion, Escape staying
inside, complete compact framing and readable toggle. Corrected the initial
pressed-state contrast after visual inspection. No main-save actions.

Replaced increasingly nested room fallbacks with typed interiorSettings, covering
all eight current room types, model/name/framing, fixture positions and wall
thresholds. This fixes Fassano Meats inheriting Mercer Court lights outside its
walls and incorrect cutaway distances. Its lamps now align with the two authored
pendants, and cutaways use its actual9m room boundaries. Existing rooms retain
their reviewed settings. Registration checks reject unknown/prototype keys and
verify exported models; frontend tests and build evidence are in
`.runtime/interior-settings-tests.log` and `interior-settings-final-build.log`.

This improves available interior viewing space; it does not create the remaining
rooms or finish their planned activities/choreography. Main8791 remains on the
previous verified release; these changes are on the development preview.

## September 14 — Russo Motor Works 3D workshop

Replaced Russo's flat interior with an original Blender workshop: masonry,
concrete slabs, inspection lift, tool drawers/bench/vice and hanging spanners,
engine/stand/manifold/pulley, two-tier tyre rack, oil drums/pumps, compressor,
service desk/ledger and suspended fluorescent fixtures. Source is
`tools/garage_interior.py`; targeted export:
`.venv-blender/bin/python tools/export_city3d.py --only=interior-garage`.
Full exporter includes it. Scene dressing does not impersonate a player's car
or stock; actual vehicle positioning and service animations remain unfinished.

Public staging reserves a mechanic position beside the engine, nine customer
positions and a separate player entrance. Actual GLB/both-rig checks prove floor
support, furniture clearance, non-overlap and full-roster entrance clearance.
All359 frontend tests pass (`.runtime/garage-final-tests.log`); final geometry
and registration checks pass after suspension-cable export, and build passes
(`garage-final-geometry.log`, `garage-reviewed-build.log`). The initial test
attempt preceded GLB export completion; its missing-asset failures were rerun
after completion. The registration rejection example was changed from garage
(now supported) to missing-room-id.

Isolated garage-interior QA fixture on8942: browser1280x720 and820x740, standard
and enlarged views, orbit and rendered-mechanic selection checked. Observed140
draw calls/260840 triangles with eight public occupants plus player, prior to
four suspension cables. This is rendering-load evidence, not a60FPS guarantee.
Added cables after visual inspection. Moved the room-size button above the
camera caption to prevent overlap with fixed Back to City when scrolled.
Final browser inspection confirms separate controls and suspended fixtures.

Main8791 remains the verifiedb4b7aee release; these changes are in development.
Remaining: player vehicle/service staging, more surface wear, other business and
civic interiors, indoor crimes and the broader campaign/visual acceptance.

## September 14 — first private-flat home attack playback

Connected explicit home strikes at apartment/mercercourt to HomeStrikeScene,
reusing the full-size CityAssassination articulated approach/weapon/fall routines
inside the furnished flat. The clear front aisle fits all four variants without
rescaling rigs. Recorded identity/face/wardrobe, actual weapon model, muzzle
pulse, timed shot/impact and pain audio, spatter and fall are presentation only.
Scene clock starts after room/cast assets load, completion drives the existing
newspaper timing, and motion/reduced-motion settles immediately. Cleanup stops
audio and disposes resources; final frames redraw on resize. No core outcomes
or main campaign state are changed by playback.

Added explicit home-cue routing with same-target/minute/victim linkage when the
selected headline is a death cue. It refuses street settings, unsupported homes
and unrelated deaths. Replay remounts with the existing replay serial. Fresh
isolated home-strike fixture uses the real home_strike command, not a fabricated
outcome; served at8943. Browser observed apartment action, completed fall then
newspaper, and compact820x740 framing. Raised the camera's focus toward the
front aisle after the first desktop frame placed the body too close to controls.
Final compact frame shows body and attacker clear of the playback band.

All361 frontend tests pass (`.runtime/home-strike-tests.log`), including sampled
full-size room bounds/furniture clearance for melee/revolver/shotgun/Thompson.
Build passes (`home-strike-reviewed-build.log`). Browser replay preserved the
fixture's displayed revision-independent clock/cash/health; no gameplay action
was issued from the browser. Main8791 remains the prior verified release.

Remaining: Cypress/boarding-house indoor crime choreography, burglary scenes,
persistent indoor bodies/blood/police instead of the existing city aftermath,
and broader visual/campaign acceptance. Audio events are wired to shared tested
samples; subjective listening acceptance is not established by this check.

## September 14 — Cypress home attack interior

Extended recorded home-strike routing to estate and generalized the scene's
per-home model/aisle/focus metadata. Cypress uses its actual drawing room, with
the full-sized attacker/victim path across the foreground at glTFz3.4, clear of
the sofa, chairs and sideboard. Flat staging remains unchanged. Estate camera
focus moves forward to2.7 so the action has clearance above the playback band.

All361 frontend tests pass (`.runtime/estate-strike-tests.log`); the geometry
coverage now samples approach/impact/fall for melee, revolver, shotgun and
Thompson in BOTH furnished homes. Final routing/geometry checks pass in
`estate-strike-geometry-final.log`; build passes in `estate-strike-final-build.log`.
Added estate-strike fixture using the actual home_strike command, isolated8944.
Browser1280x720: initial cast, approach/fall and correct Cypress room observed;
820x740: both participants fully visible on approach and newspaper follows the
scene. No main campaign QA actions. Main8791 remains the prior release.

Still unfinished: private Mariner room/home playback, persistent indoor crime
aftermath, burglary choreography, remaining interiors and broader acceptance.

## September 14 — private Mariner lodging and indoor home attacks

Authored a distinct6x6m lodging room in `tools/lodging_room.py`: iron bed/frame,
linen/blanket, bedside drawers/lamp, sash window/curtains, writing desk/paper/chair,
washstand/basin/jug, trunk and coat hooks. Targeted export:
`.venv-blender/bin/python tools/export_city3d.py --only=interior-lodging-room`.
Full exporter includes it. Original local assets only. The public lobby's roster
is never placed in this private room; the player can open it at their Mariner
home and return to the hall. Camera/enlargement and existing commands remain.

Mariner recorded home attacks use the new room and a diagonal aisle (root yawπ/4)
without shrinking people or furniture. Spatter now transforms through the cast
root so particles follow the diagonal too. Sampled geometry tests cover all four
weapon variants in lodging, flat and Cypress, plus both-rig grounded private-room
arrival and exclusion of public lobby occupants. All361 frontend tests passed
before the final entrance test; the final3 home geometry/routing tests pass.
Build passes. Evidence: `.runtime/lodging-room-tests.log`,
`lodging-final-geometry.log`, `lodging-final-build.log`, `lodging-room-export.log`.

Isolated lodging-strike fixture uses the real home_strike command, port8945.
Browser checked recorded initial cast/room, completion→newspaper, and then
Mariner lobby→Go to your room at820x740. The latter shows only the player, with
return-to-hall control and actual public roster below. Reduced the small bedside
lamp from generic12 to1.4 after observed wall blowout; final build includes it.
The lower intensity has not yet received a fresh browser screenshot check.
No main-save QA actions. Main8791 still serves the prior verified release.

All current residential addresses now select their appropriate indoor attack
room, but this does not finish crimes/interiors: burglary animation, persistent
indoor aftermath/police, other destinations and full visual/campaign acceptance
remain. The new private lodging also does not depict player-fitted upgrades.

## September 14 — authoritative burglary presentation outcomes

Added optional CueBurglary to the existing completed robbery cue, capturing the
intruder/resident, witnessed historical occupancy, success, actual cash taken,
actual injury, identification and fatality. No private remaining savings are
exposed, and no extra random draw or gameplay outcome changed. This prevents a
future indoor scene from treating an empty successful search as a failed attack
or animating an escape after a fatal confrontation. Legacy robbery cues omit it.

Tests execute real burglary outcomes:180cash once, empty successful search,
occupied failed attempt with untouched stash/actual injury/identification,
fatal attempt with no escape, and persistence of the committed result after
serialization and subsequent savings/movement changes. Initial save test used
the intentionally unsaved transient VisualCues buffer; corrected it to Execute
and persisted LastResult, matching actual command flow.

Focused burglary/home core tests pass0.962s, store0.154s, server contract tests
0.900s; frontend build passes. Evidence: `.runtime/burglary-presentation-tests.log`,
`burglary-contract-tests.log`, `burglary-contract-build.log`. No main-save QA.
This is required presentation data; the dedicated indoor search/confrontation/
escape animation remains unfinished, alongside persistent indoor aftermath and
remaining interiors. Main8791 remains the prior verified release.

## September 14 — working residential drawers for search staging

Replaced solid bedside props in private-flat and Mariner lodging models with
hollow cabinet cases and physically sliding upper trays. Reusable original
Blender source `tools/residential_storage.py`; grouped as burglary-drawer,
travel+.34m glTFZ, preserving case/lamp/radio positions. Re-exported only the
two interiors through their existing --only targets. Geometry tests verify a
real accessible tray floor, whole-drawer travel and closed-position restoration.

Private player rooms now offer local opening/closing to review and interact
with these props. This is presentation only, not a household inventory or loot
command. Smooth motion follows the current preference; reduced/off settles.
Browser isolated8945 verified visible open tray and fully closed case. Also
verified the previously reduced Mariner lamp intensity after loading fresh
assets. The closing screenshot exposed dark hover contrast; added explicit
light hover colours to drawer/enlargement controls in the final build.

All364 frontend tests pass (`.runtime/storage-tests.log`), build passes
(`storage-final-build.log`). Exports: storage-flat-export.log and
storage-lodging-export.log. Main campaign untouched for QA and main8791 remains
on the prior release. Recorded burglary approach/search/outcome animation is
still unfinished; these assets remove the solid-furniture obstacle to staging
that interaction. Remaining interiors/aftermath/campaign work remain active.

## September 14 — recorded unattended burglary search

Connected successful absent-resident burglary cues to the existing home renderer
for private flats and Mariner lodging. BurglarySearch walks an authored furniture-
clear route, opens the actual bedside drawer, leans/reaches over the tray, shows
a small cash bundle only for actual taken>0, closes the drawer and exits. The
bundle follows the actual hand during the short pocketing gesture. Empty searches
never show it. No resident is invented. Occupied/failed/fatal/estate and generic
robbery cues are deliberately still pending their distinct scenes.

Added burglary-search isolated fixture using the real committed burgle action;
new QA server8946 includes CueBurglary (older release binaries would discard the
new optional saved cue). Browser observed entry, approach, exit and report flow.
Geometry tests found the initial hand stopped short of the tray; moved the stance
closer and corrected Euler rotation order to YXZ so the lean is toward the
cabinet. Actual exported hand centers now project inside the open tray for both
rigs and both rooms. Sampled full routes verify floor support and furniture
clearance; checks also prove empty searches have no money and unsupported
outcomes do not select this scene. Final corrected hand contact has geometry
evidence; the timed browser capture caught retreat rather than the brief search.

All366 frontend tests pass (`.runtime/burglary-search-final-tests.log`), final
build passes (`burglary-search-final-build.log`), and dedicated contact tests pass
(`burglary-search-contact.log`). No gameplay commands from browser replay and no
main-save QA. Main8791 remains the prior verified release. Search/confrontation
sound design, occupied/failure/fatal and Cypress burglary choreography, indoor
persistent aftermath and remaining interiors/campaign acceptance remain open.

## September 14 — clean interior/crime release now on8791

Released e4270de from clean local clone `.runtime/release-e4270de`, excluding
pre-existing uncommitted armed/mugging/robbery/aftermath work and sim/died_test.go.
Clean frontend366/366 tests pass12.946s; build passes. Targeted burglary/home/fee
core checks0.942s, store0.158s, server contracts0.869s pass; qa-fixture builds.
Binary health identity reports e4270de74b610a28c76cd641015b08024e2539b5,
modified:false. The old server PID43340 was verified before replacement.

Created SQLite backup `.runtime/campaign-pre-e4270de.sqlite3` and tested a copy
on8947 before main startup. The first byte comparison found24 added fields:
car/drove initialization for12 later-created NPCs. An independent copy opened
by the OLD b4b7aee server on8948 produced the exact same state digest as the new
release, establishing unchanged existing SettleCars loader behavior. No new
migration/regression was inferred from that expected startup normalization.

Both copied releases and the live release match SHA256
63a214cca7228db606d159e88e4dd2df2eae332f37ec4083ffc9235eb1b5c25d.
Complete player state, revision2115, minute160095 and all2117 receipt contents
are preserved. Receipt identity used sorted id/result JSON contents, not count
alone. Main-save access for QA was read-only; actual main startup followed copied
verification and a backup. Evidence under `.runtime/release-e4270de`: check-before,
check-diff, compatibility and live-verification JSON; sibling release-e4270de
build/test/identity logs. Main session is71997, serving its own clean dist/binary.

This release brings the butcher/workshop/private-lodging interiors, enlarged
room controls, working residential drawers, indoor home strikes and successful
unattended burglary searches into the main game. Occupied/failed/fatal/Cypress
burglary scenes, persistent indoor aftermath, remaining interiors, richer income
progression and broader visual/played-campaign acceptance are still unfinished.

### Resident-funded household repair jobs — September 14

Added householdwork at The Mariner, Ashbury Court and Mercer Court. A named living
resident commissions 90 minutes of small domestic repairs for $45, funded by their
purse/savings, retaining $100. One booking per building/day; 08:00–16:30 acceptance,
sound building required. Start reserves the building's saved day; completion pays
only if the original customer remains alive, resident and funded. Interruptions
pay nothing and keep the booking taken. No structural repairs or deed/tenant
changes are implied. Mariner makes this accessible from the starting district.

Isolated household-work fixture on8949: browser visited Mariner and completed the
actual work command; cash90→135, respect0→1, clock08:00→09:30; repeat action became
disabled with the used-booking explanation. Main8791 remains on releasee4270de;
no main-save QA commands were made. Work currently resolves through the ordinary
local result, with no bespoke handyman animation. The offer remains below the
premises section in the current interior layout; improving paid-work prominence
for newcomers remains useful UX work.

Verification repeated in `.runtime/household-work-clean`, an archive of committed
HEAD with only this change copied in, excluding pre-existing armed/robbery/mugging/
aftermath/simulation edits. Targeted household/rush/group/price core tests pass
1.134s; full store0.160s and HTTP server0.900s suites pass; qa-fixture builds.
Tests cover read purity, starting-district access, exact conserved pocket/savings
payments, customer death/move/poverty, remote/closed/damaged rejection, completion
and interruption, repeat prevention, reload/new-life persistence and next-day use.
The broader goal remains active; this is an income addition, not final visual or
campaign acceptance.

### Keep paid work discoverable inside buildings — September 14

Visitor interiors now show Work before acquisition/upgrades. Owned buildings keep
the requested management-first ordering. A paid action with a present NPC subject
stays on the shared work list as well as that person's separate interaction view;
previously the customer's arrival removed the job from the shared list. Search,
disabled reasons, other groups and authoritative command handling remain intact.

Extended isolated household fixtures for a customer physically present at Mariner
and a player-owned Mariner. Browser8950 shows Mara in the roster and her $45 job
under Work before These premises; selecting her retains the same job. Browser8951
shows Running The Mariner before Work. Visitor screenshot confirmed readable
placement. No gameplay actions were needed in either fixture; main8791 unchanged.
Frontend366/366 checks pass12.750s; production build passes2.79s. Logs:
.runtime/interior-work-tests.log and .runtime/interior-work-build.log. This run used
the current checkout (including the pre-existing aftermath edits); those unrelated
files are excluded from this commit. These changes still await clean main release
integration alongside9198e41. Remaining 3D interiors and scene fidelity stay open.

### Vittoria's 3D dining room — September 14

Replaced the restaurant's flat fallback with an original Blender dining-room GLB:
10×11m terrazzo floor, walnut wainscoting and mouldings, four set tables with facing
oxblood banquettes, plate settings/cutlery/napkins/goblets and carnations, pendant
lamps, wine sideboard/racks, service hatch with stacked plates and coffee urn,
private panelled door and original geometric still-life wall decorations. Source
`tools/restaurant_interior.py`, targeted/full export integrated in export_city3d.py;
2.1MB interior-restaurant.glb plus manifest. No external asset pack.

Registered independent room lamps, cutaways, framing and entry. Public diners use
eight authored seats and two waiting positions. Cook/cellarman/two waiters have
separate service stations; additional staff remain in the roster instead of being
posed as dining customers. The first browser pass exposed that fallback seating
issue and the second pass verified its correction. No decorative fake occupants.

Geometry checks use both character rigs: full14-person roster plus reserved player
position, seat cushion support, disjoint occupant bounds, standing furniture
clearance and clear supported entrance. All367 frontend checks pass12.997s; final
build passes2.81s. Browser8952 isolated restaurant-interior fixture shows nine
actual NPCs plus player, omitted0, 157draw calls and246296triangles (not an FPS or
cross-device performance guarantee). Inspected standard/enlarged view and zoom.

At820×740 the enlarged dining room initially clipped its left wall. Interior3D
resize now preserves a minimum horizontal field as aspect narrows while keeping
authored vertical framing on wide screens; it preserves user zoom. Fresh browser
review at820×740 shows the complete room with visible side margins. Desktop
viewport restored. Main8791 and its campaign remain unchanged; this and the prior
household-work/UX commits await clean release integration.

Remaining: detailed kitchen/private back room behind the doors, dining/serving
choreography, richer glass/material/character fidelity and the other unbuilt
interiors. This is a furnished dining-room integration, not final whole-game art
or campaign acceptance. Export/build/test logs are .runtime/restaurant-*.log.

### Clean release8ba57c1 integrated on main8791 — September 14

Built a clean local clone at `.runtime/release-8ba57c1`, excluding unrelated dirty
armed/robbery/mugging/aftermath/simulation files. Frontend367/367 pass13.090s;
production build passes2.64s. Targeted household/rush/group/price core checks
pass1.266s; store0.164s and full HTTP server0.878s pass; qa-fixture builds. Final
Go binary stamps revision8ba57c136d28dc7f313bbf58fbdd5477ba42a5c4, modified:false.
An initially untracked verification JSON was moved under ignored .runtime before
rebuilding the binary; the final source tree is clean.

Created SQLite backup `.runtime/campaign-pre-8ba57c1.sqlite3` and verified integrity.
Loaded its separate campaign-check.sqlite3 copy on8953 using the release binary
and dist. Complete state and all receipts were byte-identical to the backup.
Rechecked live state against that backup, verified old mainPID60866's executable,
then replaced it with the release on8791 (session83919). Live startup remains
byte-identical: revision2115, minute160095, full player preserved, 2117receipts.
StateSHA63a214cca7228db606d159e88e4dd2df2eae332f37ec4083ffc9235eb1b5c25d;
receiptSHA5f9dc235c96c51baa248fb740553ffe5030ce9a683d3fadd2c94eede2979d30b.
Health confirms clean new revision; index and restaurantGLB both serve200.

Evidence: release-8ba57c1/.runtime/{before,compatibility,live-verification}.json;
sibling release-8ba57c1-{frontend,build,core,server,identity}.log. QA used isolated
copies/fixtures and read-only main inspection; no gameplay command touched main.
The live release includes9198e41 household jobs,133a6a4 work-list UX and8ba57c1
Vittoria's dining room/narrow camera correction. More interiors, scene animation,
material/character fidelity and complete played-campaign acceptance remain open.

### Cypress study burglary search — September 14

Successful unattended estate burglaries now use the existing committed-result
indoor search timeline. Replaced the study desk's solid left pedestal with a
hollow case and physically sliding drawer, exported interior-cypress GLB from its
original Blender source. Added a route along the clear east aisle to that drawer,
then back out; no change to cash, odds, occupancy or Go command resolution. The
estate camera frames the study route instead of the home-strike foreground, and
was moved to the opposite side after the first review obscured the working hand.

All368 frontend checks pass13.504s; build passes2.62s. Search geometry now covers
both rigs in Cypress as well as flat/lodging: floor/furniture clearance over the
full approach/exit, real sliding tray, hand centre inside tray horizontal bounds
and at its opening height, actual-positive-only cash, closed drawer after exit.
Existing Cypress seating and all residential assassination checks still pass.

Isolated estate-search fixture on8954 executes actual burgle:mara at estate and
records success, absent resident, taken180, health_lost0, identifiedfalse. Browser
replay stages the search and reveals the newspaper after completion (~13.5s).
First camera captured open-drawer search at6.987s; final angle was inspected during
approach and exit. Exact final-angle hand contact is geometry-verified, not claimed
as a paused visual inspection. Replay used no additional gameplay command; main
8791 remains the clean8ba57c1 release. Export/build/test logs cypress-search-*.log.

Occupied, failed and fatal burglary choreography, search audio, persistent indoor
aftermath and more lifelike hand/body motion remain unfinished. No final visual
acceptance is claimed; this completes basic unattended search coverage across the
four residential addresses while the larger goal remains active.

### Continuous burglary direction changes — September 14

Replaced instantaneous waypoint yaw changes with deterministic look-ahead heading
samples along the same furniture-tested routes. Ease into facing the drawer during
the final approach; reduce limb swing at arrival/exit; after closing the drawer,
take an0.8s turning step before departing. Preserve the exit heading instead of
snapping back toward the furniture when movement finishes. This affects only the
successful unattended search presentation, not the committed outcome or game time.

New frame-by-frame checks at60Hz bound heading changes below0.12radians through all
corners, arrival, turn and exit for each residential route, and verify backwards
seeking reconstructs the same orientation. Existing both-rig route/tray/contact/
empty-cash checks pass. All369 frontend checks pass13.031s; build passes2.63s.
Browser8954 fresh replay confirms the final Cypress camera and open drawer at8.483s;
no new gameplay command. Logs .runtime/search-turns-{tests,build}.log. Main8791
remains8ba57c1 pending this and126357b's clean integration. These are procedural
movement improvements; lifelike locomotion and occupied/failure scenes are still
unfinished, and no full motion-capture-quality acceptance is implied.

### Remove home-scene shadow banding — September 14

Investigated the conspicuous concentric pattern previously described as repetitive
wood grain. Cypress GLB contains zero images and23 plain materials; the pattern
also crossed its rug. Root cause was HomeStrikeScene's unbounded default shadow
camera depth (far500) and zero depth/normal bias, producing self-shadow artifacts.
Set near0.1/far25 to cover the authored rooms, bias−0.0003, normalBias0.015, and
explicit PCF filtering. Retained1024shadow resolution and existing asset geometry.

Browser8954 fresh Cypress burglary replay directly confirms removal of the rings
from parquet, rug, upholstered chairs and cabinet while preserving cast shadows
under furniture and the actor. Production build passes2.59s; no new tests for this
reversible lighting-only adjustment. Existing369 frontend checks passed in the
preceding movement commit. Main8791 remains8ba57c1 and unchanged; this renderer
fix joins126357b/0678fd8 for later clean release integration. Log:
.runtime/home-shadow-build.log. Material/character fidelity still needs further
work; correcting this rendering defect is not final visual acceptance.

### Mercer Exchange 3D hall — September 14

Added an original10×10m Blender interior for market: stone floor/cornices/pilasters,
three brass message windows with writing pads/envelopes, pigeonhole archives,
a four-seat reading table and benches with newspapers, pinned noticeboards, a
forms desk, opal pendants and a readable Mercer Exchange plaque. Source
`tools/exchange_interior.py`; targeted/full export wired in export_city3d.py;
interior-exchange.glb (~1.4MB) and manifest registered. Decorative notices carry
no invented live quotes, offers or financial amounts. Existing game actions remain
authoritative; the location is an information/trading venue, not a new stock game.

Independent placements offer four seated readers, six standing visitors and three
role-based clerk/broker/teller stations; player has a reserved foreground entry.
Overflow remains in the actual roster and the caption reports it. Isolated8955
exchange-interior fixture shows11NPCs plus player, with7more in the list;149draw
calls/257348triangles (not a measured FPS or cross-device guarantee). Reviewed
standard and enlarged room. First pass found lettering crossing a pilaster;
replaced it with a foreground plaque and verified the full name is readable.

Geometry checks cover both rigs, full13NPC staging plus player, disjoint occupant
bounds, actual cushion support, standing clearance and the supported entrance.
All370 frontend tests pass13.189s; final plaque export's targeted room test passes
0.845s; final build passes2.62s. Logs .runtime/exchange-{export,tests,build,
final-geometry}.log. QA used isolated save and UI-only entry/enlargement; main8791
remains8ba57c1. Public cast and command meaning unchanged. More interiors, richer
materials/characters, physical document interactions and whole-campaign acceptance
remain open. This is a furnished room integration, not final visual acceptance.

### Clean release8ada776 integrated on main8791 — September 14

Released126357b Cypress search,0678fd8 movement smoothing,6fcda54 shadow correction
and8ada776 Mercer Exchange from clean clone `.runtime/release-8ada776`. Pre-existing
uncommitted armed/robbery/mugging/aftermath/simulation edits remain excluded. All370
frontend checks pass13.219s; build2.70s, store0.220s, server1.186s pass; qa-fixture
builds. Core/store/API sources are unchanged relative to8ba57c1. Final Go stamp:
8ada77631378fd11dbc2c5a07a0cbcc03baa0eb0, modified:false; clone source is clean.

Backed up main through SQLite's backup API to campaign-pre-8ada776.sqlite3 and
verified integrity. Loaded a separate copy on8956; full saved state and receipts
remain byte-identical. Rechecked main against that snapshot, verified oldPID69534's
executable, then replaced it with the release on8791/session19129. Live verification
again preserves complete state/player, revision2115, minute160095, all2117receipts.
StateSHA63a214cca7228db606d159e88e4dd2df2eae332f37ec4083ffc9235eb1b5c25d;
receiptSHA5f9dc235c96c51baa248fb740553ffe5030ce9a683d3fadd2c94eede2979d30b.
Health confirms the clean revision; index, exchangeGLB and CypressGLB return200.

Evidence under release-8ada776/.runtime: before,compatibility,live-verification JSON;
sibling release-8ada776-{frontend,build,server,identity}.log. Main inspection for QA
was read-only; actual release startup followed copied verification. No gameplay
command was sent to main. Remaining civic/commercial interiors, occupied/failure
burglary scenes, sound/aftermath, richer character/material fidelity and complete
played-campaign acceptance remain open; the full goal is still active.

### Ward Street Station 3D booking room — September 14

Added original Blender public booking-room interior:10×11m terrazzo floor, green
painted dado, timber sergeant's counter with open ruled register/nameplate/bell,
rotary telephone with receiver/dial/cord, report desk with mechanical typewriter,
steel lockers, reports board, waiting bench, pendant lamps and secured holding
area door. All locally authored in tools/precinct_interior.py; targeted/full
export integrated, interior-precinct.glb (~1.3MB) and manifest registered. Cells
beyond the holding door remain unmodeled, with no invented prisoner cast.

Placements use actual public people: three bench seats, eight standing visitor
positions, two police staff stations and separate player entry. Desk-sergeant
roles receive priority for the booking counter. Other officers beyond the two
stations remain in the list. Both-rig tests cover full13NPC staging plus player,
seat support, disjoint occupant bounds, standing furniture clearance and the
supported entrance route. All371 frontend tests pass13.302s; build passes2.73s.

Isolated precinct-interior fixture8957: browser reviewed standard/enlarged room;
actual sergeant and commissioner at the two desks, three seated visitors, one
standing visitor, player at entry;153draw calls/164968triangles (not an FPS or
cross-device guarantee). Signage readable, no observed cast/furniture overlap.
Main8791 remains8ada776; no QA gameplay command touched main. Logs
.runtime/precinct-{export,tests,build}.log. Remaining holding cells/custody staging,
physical desk interactions, other missing interiors and full visual/campaign
acceptance remain open. Unrelated dirty files excluded from this commit.

### Riverside housing capacity and residential map expansion — September 14

A one-year isolated economy probe exposed the finite stock failure: seeds 7, 41
and 97 ended with 98, 164 and 149 unhoused living NPCs respectively. The previous
137-place city also left 264 people unhoused at the 400-person roster ceiling.
Added Riverside Courts on the eastern map: four 68-unit, nine-storey wings with full-height entrances, a communal garden,
a furnished 3D entrance hall, and 272 stable apartment deeds. Total accommodation
is now 409. Existing homes, journeys and numbered deeds survive settlement.

Riverside supports $100 leases, $20 private/$10 shared daily rent, $1,000 base
unit prices with existing local crime pressure and resale spreads, actual funded
NPC trades, investment rent, private flat views, repair bookings, and existing
home-crime scenes. Current-version save loading now settles added housing and
deeds immediately; a crowded-save regression checks no cash, clock, revision or
player changes, preservation of established homes, and repeated-load idempotence.
The resident register gained a name/accommodation search for large buildings.

Evidence: three 365-day simulations now have zero shortage at every residential
command boundary and 25, 10 and 9 living NPC-owned deeds respectively after adding Riverside caretaker/porter jobs. The standalone Advance-only diagnostic can briefly show
unassigned arrivals between residential settlements; committed actions perform
settlement before publishing their result. A separate full-400-person fixture
has no shortage or capacity violations and preserves old residents and deeds.
Targeted transaction tests exercise leasing, ownership ending rent, selling
without displacement, cash-funded rent and NPC resale, and neighborhood-local
price changes. These are economy simulations, not a newly played campaign.

Browser QA on isolated port 8958: all four wings render inside their street lot;
Focus address frames the complex. The lobby shows six real NPCs and the player,
with 139 draw calls / 159,440 triangles before the final noticeboard leg detail.
The expanded lobby fits at 820×740. Both exported actor rigs pass seat support,
furniture/actor separation and entrance-path geometry checks. Full frontend suite
passed 372 tests; production TypeScript/Vite build passed. The first full rules
run identified two obsolete capacity fixtures, updated to deliberately fill the
expanded stock; their focused rerun passed. Additional address checks caught missing caretaker/porter trades, local death descriptions and fallback-front registration; these were added and their focused tests passed. A repair-booking subject is validated by residence, since bookings remain on the building board during the resident’s absence. Server and store suites passed. A final broad simulation rerun is still in progress; live promotion is recorded separately. The main campaign was not used for these tests.

Outstanding: other venue interiors, richer occupied/failed burglary choreography,
further visual detail and progression/campaign acceptance remain in the active
goal. This housing expansion does not close the overall request.

Release checkpoint for 3105991: a clean shared clone at
`.runtime/release-3105991` built successfully with `modified:false`. Its complete
residential/home/property/rent/household test selection passed in 21.711s; store
and server suites passed in 0.445s and 3.420s. Final authored assets passed all 20
selected geometry/route/settings tests, and the final build passed. Browser
inspection confirmed the corrected nine-storey wings and full-height doors;
resident search reduced the 264-person register to Wanda Costa alone.

The exact live-save backup was opened only by candidate port 8959. Comparison
proved that only `properties.riverside` and 272 appended deeds changed; all old
112 deeds, all NPCs, player state, minute 160095, revision 2115, and all 2,117
receipts were preserved. Main state remains unchanged. Candidate startup's
expected state SHA-256 is
`2c407ead0ac62316d20a3c104857cf9741b8979fc0f1db7a7d505550b103b0ae`.
Evidence is in the candidate's ignored `.runtime/compatibility.json` and
`before.json`. An initial path mistake started an empty isolated compatibility
save; it was stopped and retained as `empty-start.sqlite3`, then the actual
backup was created and verified byte-for-byte before candidate startup.

The broad working-tree simulation suite timed out at 600s, with an unrelated
untracked 260-campaign sweep among its running tests. That is not a pass. Clean
release simulation validation has been restarted with a 20-minute limit; the
working-tree full core rerun is also still running. Main remains on release
8ada776 pending this verification. Do not describe 3105991 as live yet.

### Ackerman & Son pawnshop interior — September 14

Added an original 8×9-metre 3D shop: oak valuation counter, watch trays, balance
scales, register and pledge book, cameras, shelf radios, wall clocks, framed
pictures and pledged trunks. These are scenery; no invented stock or prices
replace actual pawn actions. Three distinct staff stations place the pawnbroker
at the watch tray, valuer at the scales and counter clerk at the register. Six
customer positions and a separate entrance preserve the public aisle; excess
staff or visitors remain accessible in the existing people list.

`interior-pawn.glb` is 1.1 MB with separately cutaway side/rear walls. The regular
room camera, map-sized expansion, picking, motion controls and pawn/business
commands remain in use. Isolated browser QA at port 8960 showed eight actual NPCs
and the player (one extra pawnbroker in the list), 141 calls and 189,088 triangles.
The expanded room fits 820×740. Both rigs pass furniture clearance, actor spacing
and all sampled entrance positions; role checks bind clerk and valuer to their
correct workstations. All 373 frontend tests passed in 24.631s; the final staff
placement's geometry test and production build passed after that refinement.

The broad core run for Riverside finished in 549.708s with one false-positive
prose test: “Greta Berger of Alex Varga’s people is…” was parsed as though “people”
were the subject. Commit 6dd6e29 preserves the full clause for the already-existing
agreement checker; the exact failing campaign now passes (0.711s). All other
core tests in that run passed. This changes only the test, not game prose or rules.
The clean release simulation suite remains in progress; no live promotion yet.

### Verified live promotion — 16c1c95

The clean 3105991 simulation suite passed in 397.132s. Release 16c1c95 adds only
the pawnshop presentation/QA fixture and corrected prose test beyond that tested
ruleset. Its clean build reports `modified:false`; the corrected naming test,
store suite and server suite passed in 0.376s, 0.187s and 1.138s respectively.

A fresh SQLite backup `.runtime/campaign-pre-16c1c95.sqlite3` preceded candidate
startup on isolated port 8961. Candidate comparison again proved that only the
Riverside property and its 272 new deeds are added, preserving every existing
field and all receipts. Port 8791 now serves the clean release from
`.runtime/release-16c1c95`, including Ward Street Station, Riverside Courts and
Ackerman & Son. No QA command was executed against the live campaign.

Read-only live verification matches the expected upgraded state hash exactly:
`2c407ead0ac62316d20a3c104857cf9741b8979fc0f1db7a7d505550b103b0ae`.
Revision 2115, minute 160095, all 2,117 receipts, and Jamie Moretti (cash $2,976,
health 100, home Ashbury Court, location Fassano Meats) are preserved. There are
384 apartment deeds. The release's `.runtime/compatibility.json` and
`.runtime/live-verification.json` hold the assertions and clean health response.
The broader interior/scenes/progression goal remains active.

### 2026-09-14 — billiards command and receipt boundary

Connected funded physical billiards to `pool_start/place/shot/opponent/decide/
concede/close` commands and a read-only local table/opponent projection. The
client submits cue intent; Go owns physics, rules, opponent inputs and settlement.
Active wagers require explicit concession before another activity. Event decisions
and appearance changes remain available; interrupted racks can be conceded.
API.md documents coordinates, pocket indexing, replay encoding and all command
inputs. The playable view and tournaments remain unfinished.

Core tests cover real player/NPC strokes, wrong-turn and missing/forged inputs,
interruption, explicit close/concession and immutable/local public views. Temporary
SQLite checks concurrent duplicate stake requests, a real winning shot, restart,
byte-identical winning receipt retry, stale new-ID rejection and one payout. HTTP
checks start/place/break, duplicates, replay, malformed input and response shapes.
The extended shape check exposed null travel/effect lists in action receipts;
command initialization now makes those empty lists, consistent with the existing
contract. Nullable rack/previous stroke are explicitly documented in that guard.

Final core pool/card/travel/lifecycle selection passes (0.333s), full store tests
pass (0.271s), full server tests pass (1.286s), and vet passes for billiards/core/
store/server. Evidence: `.runtime/pool-api-core-verified.log`,
`.runtime/pool-api-adapters-verified.log`, `.runtime/pool-api-vet-final.log`.
No frontend changes or live campaign QA/promotion in this increment. Next is the
close playable 3D table and controls, with physical replay and corrected table
proportions; full fidelity, multiple active tables and tournaments remain open.

### 2026-09-14 — first playable 3D billiards client

Added funded local challenges and a full-screen 3D billiards table using the
existing exactly-once command client. Aim/power/tip controls, placement, calls,
safety, break decisions, physical NPC turns, concession and close now connect to
Go. Numbered spheres rotate through the compressed physical replay; collision
samples are preserved. The playable table uses the solver's cloth/radius/jaw
coordinates. Camera review changed the default view across the table. The room's
main activity area now offers opponents, after browser QA found an initial offer
placement reachable only in property details. Error notices display over the game.

378 frontend tests pass (14.879s); final typecheck/build pass with the existing
bundle warning. Isolated 8964 browser play covers funded start, placement, player
stroke/replay locking, NPC reply, restore, camera reset, concession and close.
The actual NPC replay has 324 frames and 156 events over 5.5875s; all 16 final ball
positions/pocket states match the public result (positions within 1e-6m). The
fixture ends revision6/minute606/$5980/pool=null. Main save was not used.

This is a playable first pass, not full visual/physics acceptance. Next: pointer
placement/aiming, pocket picking, cue and player animation, richer table/pocket
materials, corrected six-table hall props, compact/motion/skip browser checks and
full-rack playtests. Concurrent tables and whole-pool tournaments remain open.
See docs/BILLIARDS.md for evidence and retained physical fidelity limitations.

### 2026-09-14 — direct billiards aiming, placement and calls

Added cloth clicking for aim and cue-placement preview, direct legal-ball calls,
numbered pocket picking, selected-call rings, ghost cue and visible head string.
Placement requires confirmation and gives local bounds/overlap feedback before
Go validates it. Numeric placement remains available in a disclosure. Camera
movement retains its gestures; tap detection rejects drags, return-to-origin
movement, secondary buttons and multiple pointers.

380 frontend tests pass (14.044s); final five billiards tests and typecheck/build
pass after small input-formatting/clearance changes. Isolated browser8964 proves
preview clicks do not advance the save, invalid head placement disables confirm,
valid placement commits, a head-ball click aims the cue, and a drag leaves aim
unchanged. The actual pointer-aimed break was legal. Following two physical NPC
shots, clicks selected ball2 and pocket6 with matching controls. The fixture now
holds revision11/minute614/$5960/$40 escrow/three shots. Main save unchanged.
Skip was visible but finished before the test click; compact, touch and successful
skip checks remain. See docs/BILLIARDS.md for details and broader unfinished work.

### 2026-09-14 — committed billiards cue animation

Added a detailed tapered cue and saved-intent draw-back/strike/follow-through/
withdrawal for both player and NPC strokes. Top/side offsets place the chalked tip
on the correct sphere surface. Physical replay starts at visible cue contact;
all original impact samples remain intact. Skip and motion preference cancellation
cover both cue and balls. Idle cue orientation reflects the spin controls too.

382 frontend tests pass (14.120s), final billiards tests pass, and final production
build/typecheck pass with the existing bundle warning. Isolated CUA54 captured
player and NPC cue playback; successful active replay skip unlocked the next turn.
The fixture ends revision13/minute618/five shots/$5960 cash/$40 escrow, with exactly
13 receipts and no extra skip command. Main campaign unchanged. Human bridge/hand
and body choreography, table/hall detail, tournament/concurrent-table work and
remaining physics fidelity acceptance are still open. See docs/BILLIARDS.md.

### 2026-09-14 — billiards cushions matched to collision geometry

Replaced centred boxes (12.5mm visual intrusion) with closed cushion profiles whose
nose lies on the authoritative segment at ball-centre height. Added raised wood
rails with grain and sights, contrasting cushion slopes, pocket lips and recessed
dark cups; split the aprons so they do not block the pocket openings.

384 frontend tests pass (14.076s), including all18 profiles' contact half-planes,
nose heights and a rendered end-cushion raycast. Final typecheck/build passes with
the existing bundle warning. CUA55 inspected the table at isolated8964 with no
command or save change (revision13/minute618). Main campaign untouched. The hall
props still require matching dimensions, and corner/pocket detail, character
choreography, tournaments and full-fidelity acceptance remain open.

### 2026-09-14 — hall tables matched to the playable minigame

Re-authored all six Blender table props to use the solver's1.27×2.54m cloth,
0.78m playing height and28.575mm ball radius. Added matching cushion/jaw profiles,
raised wood rails/sights, split aprons, open cups and pocket lips; resized legs.
Existing room positions, lights and NPC staging remain usable.

Targeted geometry/clearance test passes (3.684s):36 open pockets, six cloth/ball
probes,24 cushion-dimension rays, both rigs and sampled entrance paths. Final
build/typecheck passes with existing bundle warning. CUA56 on fresh isolated8965
inspected standard and enlarged close views without a gameplay command. Main
campaign untouched. Model4.96MB/83,312 mesh triangles versus3.43MB/61,584 previously;
these are asset totals, not verified runtime performance. Multiple active matches,
tournaments, character play animation and broader interior fidelity remain open.

### 2026-09-14 — saved physical billiards tournament bracket

Added single-elimination bracket state for2/4/8 distinct entrants. Every game is
an actual eight-ball match; later pairings derive only from adjudicated match or
concession winners. Save/reopen and repeat advancement preserve progress. Table
allocation accounts for overlapping rounds instead of assuming earlier rounds
finished together; the regression verifies an early semifinal does not displace
an unfinished opening game.

The full billiards suite passes (11.576s), vet passes, and a physical bracket
fixture produces a legal called-eight champion after13 strokes. No campaign or
HTTP interface changed. Entry fees, occasional scheduling, participant lifecycle,
full-pool payout and playable bracket integration remain next; this foundation
alone is not a completed tournament feature. Main campaign untouched.

### 2026-09-14 — tournament withdrawals before fee integration

Added withdrawal-aware bracket progression, including saved waiting entrants,
empty branches/byes and an explicit finished-without-champion state. Simultaneous
withdrawals validate and apply as a batch to avoid awarding a tournament to a
second casualty while processing the first. Resolved branches release tables;
finished champions remain immutable.

Full billiards tests pass (11.623s), vet passes. Tests exercise partial-round
save/reopen, active concessions, all-entrant withdrawal, malformed batch atomicity
and finished-event protection. No campaign money or live save changed. Entry-fee,
payout, scheduling and UI integration remain outstanding; withdrawal lifecycle
was completed first so those integrations can distinguish winners from void events.

### 2026-09-14 — funded tournament entries and whole-pool settlement

Added internal tournament campaign state and atomic fee reservation for4/8 actual
entrants. Reconciliation pays the champion the entire escrow, records only player
profit as earnings, handles withdrawal/closure, protects prior-life money, pins
eligible NPCs and retains unpaid identities. Cancelled events return unforfeited
fees; forfeited fees go to the hall till. One-stroke NPC execution uses real
physics and saves replay/intent. Scheduling and command/UI integration remain off.

Pool/tournament tests pass (0.373s), store suite (0.287s), server suite (1.190s),
and vet pass. Tests verify physical player/NPC final payouts, money conservation,
invalid-entry atomicity, repeated settlement, refunds/forfeitures, simultaneous
withdrawals and actual death/new-life isolation. Temporary SQLite reopen checks
preserve deposits, rounds and exactly one prize across restart. Main save unused.
Next: scheduled entry, public bracket/player commands, receipt retries and browser
progression; tournament funds are not exposed to normal player actions yet.

### Verified live promotion — 89cc61b

Port 8791 now serves clean release 89cc61b8a8dfad7a2aafc5428d071b88136b022d
from `.runtime/release-89cc61b` (`modified:false`). This includes playable
physical billiards, public multi-table tournaments, owner-set entry fees and
house cuts, optional owner participation, funded field previews, unattended
NPC match progression, hall table selection and the lossless GLB transfer fix.

Clean checkout validation: 388 frontend tests passed, production build passed,
Go vet passed, and full billiards/core/store/server suites passed in
12.346s/320.209s/0.439s/2.216s. Logs are `.runtime/release-89cc61b-{frontend,build,go}.log`.
The unrelated dirty gameplay/aftermath files were excluded from the release.

Before promotion, SQLite backup `.runtime/pre-89cc61b-campaign.sqlite3` was
loaded through the new server on isolated port 8971. Recursive comparison
allowed only `last_pool_tournament_slot:0` and the existing last-result
comings/cues null-to-empty-list normalization. Every other saved value and
all 2,134 receipts were preserved. Browser review showed the city and Fassano
Meats interior loading correctly on that copy, with actual occupants.
No gameplay QA commands were submitted to the live campaign.

The unchanged live save was rechecked immediately before stopping the old
server. Post-startup state matched the verified candidate byte for byte:
SHA256 `4cca8f4fe0be17920d8e010f7a1fb2023a0ea7843eeb445e49ffb17759916af5`.
Revision 2132, minute 160825, life 11 Jamie Moretti, cash $2,625 and all receipts
were retained. Release-local `.runtime/compatibility.json` and
`.runtime/live-verification.json` record the assertions and clean health data.

This is release progress, not completion of the broad goal. Full human
shooting/bridge animation, airborne billiards physics, remaining planned
interiors and broader visual/campaign acceptance remain open.

### 2026-09-14 — inspectable billiards replays

Previous e89b18b was verified release progress. Added explicit last-shot replay,
pause/resume, quarter/half/normal speed and return-to-saved-table controls to
the close billiards view. These consume the existing encoded physics tape and
cue intent without submitting commands. Spectators and finished tournament
racks can replay too. Automatic playback still respects motion preferences;
explicit replay is user initiated, and changing motion preference interrupts it.
The presentation clock preserves elapsed position across speed/pause changes.
Paused replay stops repeated GPU rendering while camera movement still redraws.

389 frontend tests passed, including continuity across pause, speed changes,
restart and cue contact. TypeScript/production build passed; logs:
`.runtime/pool-replay-controls-{tests,build}.log`. Browser tab 66 on isolated
8968 reviewed the saved championship shot at quarter speed, paused it, verified
unchanged ball positions across screenshots, orbited while paused and returned
to the authoritative resting layout. The fixture retained revision 1/minute
1082/one receipt and state SHA256
`ec1c28526c72844fee90d05ba316ecd194ba796ee4a2bb084ad68372cd21ba48`.
No gameplay commands were issued. Main 8791 remains clean release 89cc61b;
these new controls are currently in the working preview only.

Close-table character shooting/bridge animation remains absent: the cue still
moves without a person. A complete pose must keep feet outside the cabinet,
solve both hands to the shaft and handle shots needing a rest. This turn adds
replay inspection, not that character animation or airborne/slate physics.
The broad interiors, art and campaign goal remains active.

### 2026-09-14 — Herald newsroom interior

Previous acc8d2d was committed replay UX progress. Replaced the Herald's flat
interior fallback with an original Blender newsroom: four copy desks with
individual typewriter keys, platens and paper, green desk lamps, wooden chairs,
filing cabinets, pigeonholes, clipping board, wire-service machine, pendant
lights and a readable masthead. Source `tools/herald_interior.py` supports
selective/full export via export_city3d.py. The GLB is 2,004,732 bytes.

Registered Herald settings, authored two seated editorial positions and five
visitor positions, and retained the normal people/action selection. Only
publicly present occupants are rendered. Added herald-interior isolated fixture.
390 frontend tests and production build pass, including real-geometry tests
for both rigs, chair support, actor separation, standing furniture clearance
and the sampled entrance path. Logs: `.runtime/herald-{export,tests,build}.log`.

Browser tabs 67/68 on isolated port 8972 reviewed standard and enlarged room
views, the final seated staff and direct character selection (Mara's current
fixture identity/actions selected from her rendered model). Fixture
`.runtime/herald-interior-qa.sqlite3` remains revision 0/minute 600 with no
receipts; no gameplay commands issued and main save untouched. Main 8791
continues to serve release 89cc61b; this room is preview-only until promotion.

This is the city desk, not all three storeys of the Herald. Typing animation,
upper offices, other unfinished interiors and full billiards character/physics
fidelity remain open. The broad goal remains active.

### 2026-09-14 — occupied newsroom typing

Previous 5cecc1d was room-production progress. Moved desk chairs 30cm toward
the keyboards so the seated rigs can reach without stretched arms. Added
NewsroomTyping to actual occupied editorial spots: alternating hand presses,
reading pauses and staggered timing. Both rigid sleeves solve toward the
physical keyboard; hand undersides remain above its keys. The cast batch
updates with the articulated pose. No fabricated news or gameplay time.
Motion-disabled and hidden-tab states freeze the cosmetic clock.

391 frontend tests pass, including 150 sampled poses for each actual rig,
keyboard bounds/height, unchanged segment scales and disabled/hidden motion.
Existing chair-support, actor separation and entrance-clearance checks still
pass. Production build passes. Logs `.runtime/herald-typing-tests-final.log`,
`.runtime/herald-typing-build.log`, `.runtime/herald-typing-export.log`.
CUA tabs 69/70 reviewed close keyboard reach on isolated 8972; the first raised
elbow hint was corrected to bend back toward the body and the final build was
visually inspected. This is hand/arm motion, not individually articulated
fingers, mechanical key depression or moving carriage animation.
No main campaign commands or release promotion; 8791 stays at 89cc61b.
Broader interiors and billiards fidelity remain unfinished.

### 2026-09-14 — Thorne & Sons public chapel

Previous d693d60 was committed newsroom animation progress. Added an original
Blender chapel interior for Thorne & Sons: walnut panelling, cornices, eight
upholstered visitor chairs, reception desk/ledger, brass fittings, flowers,
sconces and a closed display casket on a bier. This is scenery and does not
assert a dead NPC or create funeral simulation. Added selective/full model
export, chapel registration, occupied-seat staging and isolated chapel fixture.

392 frontend tests pass; production build passes. Tests load real exported
geometry and both character rigs, checking seat support, actor separation,
reception clearance and the sampled entrance path. The first check caught
a carpet-top mismatch with the shared walking floor; corrected the authored
carpet height and reran successfully. Logs `.runtime/chapel-tests.log`,
`.runtime/chapel-tests-final.log`, `.runtime/chapel-build.log`,
`.runtime/chapel-export.log`.
Browser tab 71 at isolated8973 reviewed the room, seated occupants and direct
receptionist selection into the correct existing action list. Fixture remains
revision0/minute600/no receipts; no gameplay commands or main-save writes.
Main8791 is still89cc61b; this room awaits release promotion.

The public chapel does not include the rear preparation room/hearse yard.
Other interiors, broader visual acceptance and full billiards physics/character
animation remain unfinished; the goal remains active.

### 2026-09-14 — correct cue spin units

Previous3c6cf3f was interior-production progress. Physics review found a real
public/solver mismatch: API and UI supplied metres, but Shoot multiplied as
though Top/Side were normalized radii and allowed offsets up to0.6m. Corrected
the angular impulse to2.5*v*offset/R² and the limit to0.6R (17.145mm). Updated
older solver fixtures to express their intended radius fractions in metres,
preserving the physical scenarios instead of weakening assertions.

Full billiards suite passed11.875s, core Pool command/tournament selection
passed8.662s, store/server suites passed0.283s/1.890s and vet passed. Logs
`.runtime/pool-spin-units-{tests,core,adapters}.log`. New checks verify initial
angular momentum along three cue directions, combined UI maximum offsets,
out-of-range rejection, and a command-level12mm draw that reverses cue-ball
travel after object contact. Rejection preserves state; saved intent retains
metres and the resulting replay decodes. No frontend changes or browser claims
for this solver correction. No live save mutations or promotion.

Existing replay tapes remain historical outcomes, with no retroactive
resimulation/payout changes. Centre hits are unchanged. This fixes spin
strength, not missing airborne/jump/massé/slate physics or shooting characters.
The broader goal remains active.

### 2026-09-14 — physical opponent position play

Previous2bc008a corrected cue units. Opponents now preview draw/follow on their
four strongest pot lines, then evaluate the actual settled cue ball for another
clear pot. The extra candidates use ±8/12mm top offsets, the same physical
preview/rules as existing shots, and full-resolution final execution. A small
spin penalty preserves centre hits when outcomes and position are equivalent.

The focused fixture chooses12mm draw, legally pockets its called ball and
leaves a better next-pot line than an otherwise identical centre strike.
Re-executing the advertised intent reproduces the complete match/physics result.
Full billiards passed17.930s, core Pool/tournament checks25.899s, store/server
0.302s/1.936s; vet passed. Three-iteration mid-rack benchmark measured248ms/op
and7.13MB/op on this machine under concurrent tests. Logs are
`.runtime/pool-bot-position-{tests,core,adapters,bench}.log`. This is added
planning cost, not a whole-game performance certification.
No browser/live promotion or main campaign commands. Side-spin planning,
airborne physics, shooting-character animation and the remaining broad
interior/gameplay goal are still unfinished.

### Clean release candidate — 01fb5cf (validation still running)

Prepared isolated clean clone `.runtime/release-01fb5cf` at
01fb5cf613d01a45b48e441ffb012742f1a595bc, excluding the unrelated dirty
gameplay/aftermath work. Clean frontend392 tests pass15.222s, production build
passes2.60s and Go vet passes. Binary health reports modified:false.
Fresh SQLite backup `.runtime/pre-01fb5cf-campaign.sqlite3` and candidate copy
`.runtime/candidate-01fb5cf.sqlite3` preserve all state bytes and2,134 receipts.
Candidate on8974 (session49933) matches SHA256
4cca8f4fe0be17920d8e010f7a1fb2023a0ea7843eeb445e49ffb17759916af5;
revision2132/minute160825/Jamie Moretti life11/$2,625 retained.
Release-local `.runtime/compatibility.json` records exact assertions.
Browser72 reviewed the clean city view on the copy, without commands.

Full billiards/core/store/server suite is STILL RUNNING in session25227.
Billiards passed18.540s; core PID62583 was confirmed active with CPU use.
Poll this same handle; do not restart or claim success until its terminal
result is known. Logs `.runtime/release-01fb5cf-{go,frontend,build,vet}.log`.
Build session49560 completed successfully. Main8791 remains release89cc61b
(PID54030); no promotion or main-save mutation yet. Recheck live save against
the backup immediately before any eventual restart. Broad goal remains active.

### Verified live promotion — 01fb5cf

Previous71efe7f prepared the clean candidate and verified a live test wait.
Resolved the same full-suite session25227 successfully: billiards18.540s,
core281.208s, store0.385s and server2.402s. Together with392 frontend tests,
production build and vet, this completes the release checks.

Verified compressed candidate transfers of Herald, chapel and pool-hall GLBs
against the clean source files byte for byte; release-local
`.runtime/model-verification.json` records lengths and SHA256. Live save and
all receipts were rechecked against the fresh backup immediately before
stopping only the old main server. Port8791 now serves clean
01fb5cf613d01a45b48e441ffb012742f1a595bc from `.runtime/release-01fb5cf`.
Health reports modified:false. New main session57534.

Post-startup state remains byte-identical to the backup (SHA256
4cca8f4fe0be17920d8e010f7a1fb2023a0ea7843eeb445e49ffb17759916af5),
with all2,134 receipts unchanged. Revision2132/minute160825/life11 Jamie
Moretti/$2,625 preserved. Release-local `.runtime/live-verification.json`
records the assertions. No QA gameplay commands were executed against main.

Live additions: inspectable slow-motion billiards replays, corrected metre-based
cue spin, physical opponent draw/follow position planning, Herald newsroom
with occupied typing desks, and Thorne & Sons public chapel. Remaining rooms,
full shooting characters, airborne billiards physics and broader art/campaign
acceptance still prevent completion of the overall goal.

### 2026-09-14 — apartment exchange usability

Previous d71f2a3 completed verified live release progress. Added all/owned/for-sale
listing filters, address/number/resident/owner search, address/value/rent sorting
and an apartment holdings summary. Broker offers remain separate from scheduled
rent and its tenant-cash limitation. Individual apartments now explain the
public neighborhood discount already used to calculate their prices.

The new neighborhood_index projection is documented in API.md; a core test
verifies a recorded explosion's10% discount, per-building consistency and read
purity. Apartment/property selection passed0.187s, store/server0.297s/1.883s,
392 frontend tests and production build passed. Logs
`.runtime/apartment-exchange-{core,adapters,frontend,build}.log`.
Browser73 on isolated8975 reviewed the period market layout, empty owned-deed
filter (0 of10) and case-insensitive MERCER search (5 of10). No gameplay commands
were issued; no main save writes. This is existing-market UX, not new housing
stock or a claim that broader progression/visual requirements are complete.
Main8791 remains clean01fb5cf; these changes are preview-only.

### 2026-09-14 — city-wide NPC apartment investors

Previous b9d9903 improved housing-market inspection. Removed the private deed
buyer's arbitrary restriction to residents of the three apartment buildings.
Living housed NPCs, including Mariner tenants and house residents, can now buy
a cash-poor NPC's apartment investment. Stable identity ordering, real household
funds, one daily transaction, existing tenancy and original housing/journeys
remain in force. This does not force a move or create new wealth.

A Mariner-investor regression verifies the exact buyer-to-seller payment,
conserved household funds, retained seller tenancy, unchanged buyer home and
active journey, and subsequent tenant-funded rent paid to the new landlord.
Apartment/housing/private-trade checks passed7.190s, store/server0.325s/1.897s;
vet passed. Logs `.runtime/citywide-landlords-{core,adapters}.log`.
No frontend/interface changes, browser claims, main-save writes or promotion.
Main8791 remains01fb5cf. NPC move-up housing behavior, other interiors and
full billiards fidelity remain open within the broad goal.

### 2026-09-14 — Pier 14 cargo shed

Previous34651aa expanded funded NPC investment. Added an original Blender
receiving/dispatch shed for Pier14: corrugated walls and columns, timber crates
on pallets, freight tags, mechanical scale/dial, telephone and shipping ledger,
hand truck, rope coils and hanging industrial lamps. Crates are neutral scenery,
not authoritative stock counts. Added selective/full export, room settings,
dispatch/aisle cast positions and a clear foreground entrance route.

393 frontend tests passed14.632s, including real-geometry checks for both rigs,
occupant separation, equipment clearance and sampled entrance motion.
Production build passed. Logs `.runtime/docks-{export,tests,build}.log`.
Browser74 on isolated8976 reviewed the room, actual fixture occupants and
clerk selection into the matching existing action list. Fixture
`.runtime/docks-interior-qa.sqlite3` remains revision0/minute600/no receipts;
no gameplay commands/main-save writes.

This is the transit shed, not the full waterfront/loading yard or animated
cargo work. Those, other missing interiors and full billiards fidelity remain
open. Main8791 stays01fb5cf; this room awaits release promotion.

### 2026-09-14 — smooth table camera keys

Poker, blackjack and billiards now share continuous held-key pan/orbit using
city camera input helpers. WASD/arrows move between keyboard repeat events;
Q/E rotate continuously. Key release, focus loss, visibility changes and Home
clear held input. Disposal cancels the added animation loop. Updated poker
help and billiards accessibility instructions; gameplay commands are unchanged.

394 frontend tests passed (14.693s), including the real TableCamera and
OrbitControls with controlled animation frames for movement, release, blur,
reset and disposal. Test compilation now uses bundler module resolution for
Three's addon export. Production build passed (2.69s). Logs:
`.runtime/table-camera-held-{tests,build}.log`. Browser75 on isolated8968
opened the completed tournament final and confirmed close table rendering;
held-key behavior is automated-test evidence, not a browser hold-duration test.
No gameplay commands or main-save QA writes. Main8791 remains01fb5cf;
this change and preceding interior/property work await promotion. Broad
interior and billiards fidelity work remains open.

### 2026-09-14 — Vance Cab Company dispatch office

Previous e045e79 was progress on shared table camera movement. Added an original
10m dispatch office for the cab company, replacing its non-3D fallback: cream
plaster/green dado, tiled linoleum, telephone counter with rotary dials, fare
ledger and pigeonholes, schematic route board, lockers, coffee urn, driver
bench, clock and pendants. The route board is decorative and does not assert
live cab positions/jobs. Added selective/full Blender export and manifest,
room camera/lighting settings, dispatcher position and eight visitor positions
(two seated), and an isolated QA fixture. Existing business/person actions
remain authoritative and unchanged.

395 frontend tests passed (14.415s); real GLB/character tests cover stable
placement, both rigs, bench support, separation, standing furniture clearance
and sampled entrance path/floor support. Production build passed (3.02s).
Logs `.runtime/cabstand-{export,tests,build}.log`. Browser77 on isolated8977
reviewed standard and enlarged views at1280x720 and selected Mara at the
counter into her matching Dispatcher action list. Fixture
`.runtime/cabstand-interior-qa.sqlite3` remains revision0/minute600/receipts0.
No main-save QA writes. This is the dispatch office; fleet yard, active
telephone/dispatch animation, remaining interiors and billiards fidelity
remain open. Main8791 stays01fb5cf; this and preceding changes await promotion.

### 2026-09-14 — NPC purchases that improve their home

Previous7c00441 completed the Vance dispatch interior. ApartmentDay now has a
funded voluntary purchase fallback after existing resident/private sales: a
renter can buy a vacant broker apartment in a higher residential tier, preferring
Ashbury, then Mercer, then Riverside as actual funds permit. The purchase leaves
$500 in purse/savings, skips ruined premises, transfers that exact deed and
releases only the buyer's old unit. Owner-occupants remain settled, other tenants
are not displaced, and workplace/current journey/old rent debt remain intact.
Uses existing neighborhood-adjusted asking prices and daily simulation hook;
no new save fields or public command shapes. API contract updated.

Targeted core apartment/housing/home tests pass7.814s, including new funded
purchase/savings, rent-free ownership, reserve boundary, deed protection,
old-tenancy release, debt and journey preservation cases. Store tests pass.303s,
HTTP server tests1.916s, Go vet passes. Logs
`.runtime/npc-home-purchase-{tests,adapters}.log`. No browser claims or main-save
writes. Initial new test fixture needed its savings map initialized; corrected
before passing verification. Main8791 remains01fb5cf; this awaits promotion.
Downsizing, richer moving preferences and the remaining interior/billiards
fidelity work remain open within the full goal.

### 2026-09-14 — finish old-home walks without losing the workplace

Previous0c4c07c added voluntary funded home purchases. Follow-up inspection
found midnight schedules journeys before apartment trading: a buyer's old-home
arrival could overwrite their job and strand them at the former residence.
Arrival now recognizes the existing saved home-journey reason, preserves Post,
finishes the original walk, and reuses ordinary per-NPC departure planning from
the actual arrival address when Home changed. Night arrivals walk to the new
home; after dawn the normal work/urgent-duty priorities apply. No teleport,
new save fields or frontend-calculated route. SetOut delegates unchanged
scheduling rules to the same helper.

New regressions exercise purchase during an active walk, onward departure,
new-home arrival and an old-home arrival after dawn. Broader home/housing,
errand/routine/journey/arrival, custody and pool tournament core tests pass
28.181s; store.303s, HTTP server1.863s and Go vet pass. Logs
`.runtime/home-move-journey-{tests,adapters}.log`. No main-save QA writes or
browser claims. Main8791 remains01fb5cf; release promotion and broad interior,
housing preference and billiards fidelity work remain open.

### 2026-09-14 — clean 3e7ee26 release staged

Prepared clean clone `.runtime/release-3e7ee26` at
3e7ee26a10338b0a306b84989f472edd8274fd21, excluding unrelated dirty work.
395 frontend tests pass16.255s; production build3.59s and Go vet pass.
Binary health reports modified:false. Candidate8978 session97631 uses only
`.runtime/candidate-3e7ee26.sqlite3`, copied from fresh SQLite backup
`.runtime/pre-3e7ee26-campaign.sqlite3`. Complete state bytes and all2134 receipts
match; revision2132/minute160825/life11 Jamie Moretti/$2625 are unchanged.
SHA2564cca8f4fe0be17920d8e010f7a1fb2023a0ea7843eeb445e49ffb17759916af5.
Release-local compatibility.json records verification. Cabstand/docks/poolhall
compressed GLB transfers match source byte-for-byte (model-verification.json).
Browser78 reviewed copied campaign and apartment exchange without commands.

Full Go suite remains RUNNING in session45900, core PID74340 verified active
with CPU use after five minutes; billiards passed18.829s. Keep this same test
handle; do not restart it or claim full-suite success without its terminal
result. Logs `.runtime/release-3e7ee26-{go,frontend,build,vet}.log`.
Main8791 is still01fb5cf (PID64096), not promoted. Recheck live save/receipts
against the backup immediately before eventual restart. This turn staged and
verified the candidate and verified an active test wait; the broad goal remains
open.

### 2026-09-14 — verified live promotion of3e7ee26

Previous0308210 staged the candidate and verified the ongoing full-suite wait.
The same session45900 completed successfully: billiards18.829s,
core422.587s, store.388s, HTTP server2.570s. Together with395 frontend
tests, production build, vet, browser and model-transfer checks, release
validation passed.

Rechecked live state and every receipt against the fresh backup immediately
before stopping the verified01fb5cf process(PID64096). Main8791 now serves
clean3e7ee26a10338b0a306b84989f472edd8274fd21 from
`.runtime/release-3e7ee26`, session5478, health modified:false.
After startup, complete saved state bytes and2134 receipts still match exactly:
SHA2564cca8f4fe0be17920d8e010f7a1fb2023a0ea7843eeb445e49ffb17759916af5;
revision2132/minute160825/life11 Jamie Moretti/$2625 preserved.
Release-local `.runtime/live-verification.json` records the assertions.
No gameplay QA commands were run against main.

Now live: apartment exchange search/filters/holdings and neighborhood context,
citywide NPC property investors, funded vacant-home purchases with onward
journeys, smooth held table-camera controls, Pier14 cargo shed and Vance cab
dispatch office. Other interiors, full billiards shooting/airborne physics,
richer housing preferences and overall art/campaign acceptance remain open.

### 2026-09-14 — billiards first-contact aiming guide

Previous495f4e1 promoted the verified housing/interior release. Replaced the
fixed.7m aiming line with a geometric first-contact guide: circle sweeps against
object balls and the shared physical cushion segments/endpoints, including
pocket jaws. A cloth ring marks the cue-ball centre at contact. Pocket openings
end beyond the bed. The helper reads current balls without mutating state;
it does not predict spin, rebounds, pockets or game outcomes. Added that
explanation to table help. Actual shots remain Go-authoritative.

399 frontend tests pass14.081s, including straight/cut hits, misses, pocketed
balls, rail-before-ball, jaw/opening geometry, touching-ball direction and
read-only behavior. Build passes2.71s. Logs
`.runtime/billiards-aim-{tests,build}.log`. Browser79 on isolated8979 started a
$20 rack against Mara and confirmed initial cue placement, then checked the
line/contact ring against the front rack ball and a cushion with a different
local aim. Fixture `.runtime/billiards-aim-qa.sqlite3` revision2/minute602/two
receipts; no stroke committed and no main-save QA writes. Main8791 remains
3e7ee26; this awaits promotion. Full shooting character/bridge pose, cue-rail
clearance, airborne physics and remaining interiors still need work.

### 2026-09-14 — poker cards-first camera preset

Previouse5be150 improved billiards aim feedback. Poker now opens at a closer
3.35m camera frame centred between the hand and board, with a Whole table /
Read cards toggle for the former4.6m view. Pan/orbit/zoom remain available;
Reset returns to the selected frame. Shared TableCamera.frameView clears held
keys when changing frame and keeps responsive aspect compensation. Corrected
selected/hover button contrast after browser review.

399 frontend tests pass14.638s, including frame selection clearing held input
and Home restoring the selected centre/distance. Final build passes3.00s.
Logs `.runtime/poker-framing-{tests,build}.log`. Browsers80–82 on isolated8980
compared the former view with close/whole-table views, including the player hand
and flop together at1280x720. One Check command dealt the fixture flop; camera
toggles remained local. No main-save QA writes. Main8791 remains3e7ee26;
this and the aiming guide await promotion. Further interior production and
full billiards fidelity remain open.

### 2026-09-14 — matching blackjack camera views

Previous1541f0e added poker cards-first framing. Blackjack now opens with a
3m close frame and the same Whole table / Read cards toggle (wide4.1m), free
camera controls and selected-frame reset. Camera buttons share styling across
both games. Added configurable TableCamera aspect fit, default unchanged for
other tables; blackjack uses1.15 to contain long hands on narrow screens.
Initial projection check exposed clipping at360x440; corrected before final
verification. Camera changes still clear held keys through frameView.

400 frontend tests pass14.813s, including both rows with2/7/12 cards at1280x440,
800x600 and360x440 using the actual blackjack camera settings. Build passes2.71s.
Logs `.runtime/blackjack-framing-{tests,build}.log`. Browser83 on isolated8981
entered the tables and dealt a$50 fixture hand, reviewing the exposed dealer9,
hidden card and playerK/7 in both camera views; button contrast is readable.
No main-save QA writes. Main8791 remains3e7ee26. Card-camera and billiards-guide
changes await promotion; broader interior and full billiards work remain open.

### 2026-09-14 — camera shortcuts follow the active game view

Previous9567a5f aligned blackjack camera framing. Shortcut fallback now searches
visible poker, blackjack, pool, interior and city canvases within the current
modal scope, instead of only the city canvas. Betting/text inputs retain their
native keys. Global data-shortcut navigation stays restricted to the body;
a modal without a camera cannot move an underlying scene. Updated keyboard
reference for pan/reset. First browser attempt exposed the old early modal
return; corrected routing to search within that modal rather than bypass it.

Production build passes2.64s (`.runtime/active-camera-build.log`). Browser85 on
isolated8980 verified zoom and Home from the focused poker framing button,
with focus transferred to the poker canvas inside the game modal. Browser84
also confirmed the existing action-search shortcut. No gameplay commands or
main-save QA writes in this turn. Main8791 remains3e7ee26; camera/aim changes
await promotion. Remaining interiors and full billiards fidelity remain open.

### September 14 — rental mobility, apartment accounts and contextual G

Added a daily, deterministic review of vacant rentals. Households with funds can
upgrade; strained tenants can seek cheaper accommodation. At most one moves per
day, no resident is displaced, owner-occupants remain, and player-held vacancies
are eligible without a deed transfer. Existing travel and workplace fields stay
unchanged. Rent collection uses the existing real-cash payment path. The market
ownership overview adds owned/vacant shortcuts and current-day collections and
current-tenant arrears alongside rent and broker values. This is an initial
management overview; adjustable leases and a dedicated property paper remain.

G now enters, leaves or skips the journey using the visible contextual action;
the existing travel G remains. Leaving is unavailable during journey playback so
it cannot compete with skipping. Input and modal shortcut guards remain.

Scene playback now explicitly requeues a selected cue at playback start regardless
of result revision. Its clock stops in hidden tabs and caps stalled-frame advance
at 100ms, preserving opening beats on return. Full explosion/fire-response browser
acceptance remains; these timing fixes alone do not establish full scene fidelity.

Evidence: targeted rental/home-purchase Go tests passed (0.146s), including a
player-owned vacancy taking a tenant, $35 rent reaching the landlord, public
collection data, preserved journey/deed, and refusal of unfunded/owner moves.
401 frontend tests passed (14.262s); production build passed (2.64s). Isolated
apartments fixture on 8982, browser tab86, verified overview controls and vacancy
labels. No gameplay commands were sent to main8791 and no main save was changed.
These changes are not yet promoted. Existing unrelated armed/robbery/aftermath
work was left unstaged. Death-service businesses and individual proprietor
purchase/succession remain outstanding after the ownership audit.

### September 14 — scene-first interior shell

Replaced the legacy stacked interior layout with a viewport-filling room and
one-at-a-time floating Building, People, Business and Residents panels. Reused
all authoritative actions, searches, tenancy accounts and private-room controls.
Selecting an actor opens their business panel; Escape closes the panel and G
retains contextual exit. Added paper/brass styling, readable ownership header,
scrollable panel bounds, and compact-width rules. Daily accounts moves to the
lower corner while inside, with panel space reserved above it. Removed the
redundant enlarge control because the standard interior now uses the full stage.

Browser evidence: isolated apartments fixture8982, tabs87/88. Initial screenshot
exposed accounts covering the new header; moved accounts and rechecked at1280x720:
Mercer Court title/ownership and private-room control visible, room centered,
business paper readable with its own scrolling. Escape and G returned to city.
Final footer-space correction passed production build; small-width CSS is present
but phone-size visual acceptance and other venue panels remain to be checked.
No main save commands or deployment; this is preview-only. No new tests for this
reversible layout change; TypeScript and production build passed.

### September 14 — recorded assassinations inside their venues

Successful named player/crew strikes now capture interior versus home/travelling
context before Kill. React resolves the selected attack or explicitly linked
victim headline to the interior scene. Actual venue models come from existing
interiorSettings, private homes retain their authored rooms, and missing rooms
show an explicit result fallback. City3D no longer stages indoor cues outside.
No outcome, time or combat RNG changes. Generic NPC deaths without captured
attacker/scenario metadata remain outside this first named-strike implementation.

Venue staging tests the recorded actor path for nearby room geometry and floor
support, including body bounds through the fall. The browser revealed a body too
near the floor edge on the first version; added body-corner floor checks and
replayed successfully. These sampled checks do not replace all venue/choreography
visual acceptance. Interior playback now caps frame advances and pauses while
hidden, like the corrected city scene clock.

Evidence:403 frontend tests passed14.461s, including explicit venue routing,
linked victim identity, legacy omissions, missing models, and clear placements
inside Saint Agnes and the cab office. Targeted strike/home tests passed0.168s;
production build2.71s. Isolated strike-revolver fixture8983, tabs89/90, replayed
inside Saint Agnes; final screenshot showed the casualty entirely on the interior
floor. Main8791/save untouched. Preview-only; other weapons/venues and generic
contract/NPC incident coverage still require broader acceptance.

### September 14 — funded civilian proprietors and ownership on death

Previous goal turn classified as progress: committed interior assassination
routing and verified an isolated replay. Continued the business-ownership audit.
Added deterministic daily NPC purchases for ordinary civilian trades, with actual
household spending and a reserve. Added personal business settlements using
condition/capacity/custom, wages, stock and repairs. Saved settlement stamp
prevents repeat daily revenue and a buyer gets no retroactive purchase-day income.

NPC-owned properties are now protected from Unheld family expansion and holder
labels resolve their owner's name. Kill releases the individual's deeds directly
to purchase while preserving staff/condition; organizational deeds stay with the
family. A stale dead personal owner is also handled by daily reconciliation.
Gambling and other special financial trades still need proprietor integration;
new morgue/cemetery/crematorium venues and richer insolvency remain outstanding.

Evidence: targeted proprietor, family succession and expansion tests pass0.154s.
They cover funded buying, working reserve, dead buyer rejection, acquisition and
expansion protection, actual daily income/wage arithmetic, repeat settlement,
personal-versus-family deed release, duplicate death idempotence and leadership
succession retaining a family business. An initial test caught a misplaced name
formatter edit; fixed it and reran. No main save commands/deployment. Unrelated
armed/robbery/aftermath edits remain unstaged. Broad campaign validation is next.

### September 14 — special-trade proprietor accounting

Previous goal turn was progress: committed funded personal business ownership and
started clean full tests for cc4a6ea. Revalidated that exact test process/session;
it is still running, with billiards and command packages passing so far. No
replacement run was started for that checkout.

Expanded proprietor acquisition to all priced trading premises plus the Mariner.
Casino purchases separately fund an initial float from the household; personal
casino nights use it, daily accounts fund a low float only above the household
reserve and draw excess above the full-float target. A failed personal casino
cannot penalize the player's respect. Existing family casino semantics remain.

Added non-player deed-based settlement for card/dice/roulette/slot outcomes,
NPC floor bets, back-room fees, tournament cuts, fuel, dealer and tailor margins,
vehicle disposal and garage repair receipts. Player-owned branches keep existing
semantics. The Mariner receives real tenant rent and operating expenses, without
a second abstract rental-income payment. NPC ownership no longer leaves those
receipts in an ownerless sink. New death-service venues remain outstanding.

Evidence: new funds/purchase/casino/reputation/hall-cut/fuel tests pass0.256s;
proprietor/rental/Mariner selection passes0.526s; vet passes. The cut test resolves
an escrow-funded bracket and checks the NPC owner receives $40 exactly once;
fuel checks both payer and proprietor; lodging checks actual transfer and no
phantom income. Prior casino/tournament/floor selection passed58.597s. Expanded
integration selection is still running in session27100 (log
.runtime/business-owner-integration-tests.log); clean cc4a6ea full suite remains
session71514 (.runtime/check-cc4a6ea-tests.log). Neither is treated as complete.
These changes are not deployed and no main-save commands were issued.

### September 14 — release preparation and purchase eligibility

Previous turn was progress: expanded gaming/service owner accounting and added
fund-transfer tests. Both earlier long-running test handles71514/27100 were
revalidated live this turn; no test restart was inferred from their quiet logs.
Added an early household-funds eligibility check before scanning holdings and
listings: every purchase already requires more than the $500 minimum reserve.
Targeted proprietor/Mariner/tournament-owner checks pass0.257s. Preparing a clean
release and copied-save compatibility evidence while broader tests finish.

### September 14 — clean release gate found failures; no promotion

The exact cc4a6ea full run (session71514) terminated with failure. It exposed
plural-family SendWord prose ("people has") and ownership-reference scan failures
for living NPC proprietors. Both core and sim also hit the default10-minute
package timeout; this was confirmed terminal, not inferred from quiet output.
Expanded owner integration session27100 passed313.238s.

Corrected SendWord agreement for has/have and is/are. The reference audit now
recognizes living individual owners while retaining errors for missing and dead
owners; a regression test checks all three. The apartment listing path appeared
in the timeout stacks. Cached the life-specific owner ID per listing/market/rent
read and avoided NPC lookups for vacant units or already-filled broker groups.
No listing eligibility or ordering policy was intentionally changed.

Clean07e41c4 candidate8984 uses only a copied campaign.403 frontend tests passed
56.886s under test contention; build8.78s. Backup revision2182/minute163810/life11
and all2184 receipts are preserved byte-for-byte after loading and public reads:
SHA256a12b3718954540bb3cec26d7ad952c9c1ebacf876194939c46350b184e5ed0d4.
Release-local compatibility.json records assertions. The first binary reported
modified:true because node_modules was a symlink rather than an ignored directory;
added that symlink to local git exclude, rebuilt and restarted only candidate8984.
Current candidate session37894 reports modified:false. It predates these gate
fixes and is NOT approved for promotion. Main8791 remains untouched.
Targeted apartment/reference/one-family prose checks pass2.707s after the fixes.
A new clean full suite with a longer package timeout is required before release.

### September 14 — copied campaign reveals unreachable rental tiers

Previous goal turn was progress: corrected release-gate failures and launched
clean aa9b617 full tests. Same handle10830 revalidated RUNNING; do not restart it.
Prepared clean aa9b617 binary/UI on8985, session58569, copied save only. Production
build4.06s; frontend sources are identical to07e41c4, whose403 tests passed.
Compatibility.json verifies clean build metadata, all saved bytes and2184 receipts
preserved. Browser91 entered Riverside via G and inspected the market without
sending gameplay commands. Six owned Riverside apartments were correctly listed
as vacant. Interior header and controls were clear; collapsed Latest entry still
covers part of the bottom camera-help caption, a remaining cosmetic defect.

The tenant eligibility inspection caught a real logic error: wealth bands selected
Mercer/Ashbury before a household could ever satisfy Riverside/Mercer upgrade
reserves. Replaced the bands with direct affordable-tier selection. Added shared
leases, already supported by NPCRent, after private options, using a smaller
one-week reserve so ordinary room tenants can access affordable housing. Private
leases retain two weeks; household strain permits cheaper downgrades. No deposits
or artificial cash are introduced.

Rental/home-purchase tests pass0.255s, covering boundaries69/70/84/105 for shared
leases and280/349/350/489/490 for private tiers, plus middle-tier downsizing. A
rental-only in-memory projection of the copied campaign, seven daily reviews and
no budget adjustments, changes owned Riverside occupancy0→1. This proves current
eligible demand, not a seven-day whole-campaign forecast or guaranteed occupancy.
Source/output in.runtime/rental-projection and rental-projection-result.log.
The running aa9b617 suite predates this rental correction. Candidate8985 also
predates it. Main8791 was untouched; nothing promoted this turn.

### September 14 — clearer room footer and cinematic safe area

Previous goal turn was progress: fixed unreachable rental tiers and verified
eligible demand from a copied campaign. Revalidated the ongoing aa9b617 clean
suite via session10830; still active, no replacement started.

Replaced the permanent room camera caption with an accessible Room controls
summary/panel. It sits above footer accounts/latest-entry panels, with private
room drawer controls offset above it. Recorded interior scenes now respect the
same HUD/rail margins as rooms; their playback card is compact at lower right,
with accounts in the lower left rather than covering the top of the scene.

Evidence: production build6.09s. Isolated browser93 opened Room controls in
Mercer Court at1280x720; screenshot verified readable expanded help clear of
accounts. Browser94 replayed the Saint Agnes fixture and verified the scene,
actors and compact playback card clear of the main HUD. Browser92 had loaded
the prior bundle during build and was not used as acceptance evidence. No new
tests for this reversible layout change; compact-width visual acceptance remains.
No gameplay commands were sent to main8791, and no release was promoted.

### September 14 — release validation and even casino nights

The clean aa9b617 full suite finished: sim passed498.253s; core ran573.270s
and failed only TestTheCityReadsAfterThePlayerDies (seed219), which found a
casino report saying "$0 of it kept". Other packages passed. This is a terminal
failure, not a timed-out or still-running suite. Corrected an even night's
report to say the house broke even; a high roller who breaks even no longer
gets a false "good night" report. No economic calculation changed.
Targeted death-reading and seasonal house-edge tests passed0.797s in the
working tree. Clean bb58284 frontend build passed4.00s and targeted rental,
apartment, NPC home purchase and ownership-reference tests passed0.261s.
Nothing has yet been promoted to8791; latest fixes require clean validation.

### September 14 — promote4daa6c8 with preserved campaign

Clean release-even-night validation: all403 frontend tests passed15.102s;
production build2.62s; targeted death-reading/seasonal-house-edge, rental,
apartment, NPC home purchase and ownership-reference tests passed1.187s.
These checks cover the changes after the aa9b617 full run; the earlier full
run's only failure was the now-passing death-reading report. A second complete
suite was not claimed or run for this prose correction.

Candidate8986 loaded a fresh SQLite backup of the actual campaign. Read-only
health/state calls preserved all raw campaign bytes and all2184 ordered receipt
pairs. Clean build4daa6c822320b80504ff272317924ac2aa1195eb, modifiedfalse.
Browser95 entered Riverside with G and visually verified the large room,
floating navigation, readable header and footer accounts at1235x1053.

Rechecked live3e7ee26 and the complete save against the tested backup immediately
before stopping its verifiedPID76368. Started4daa6c8 on8791 with the original
campaign path (session57392). Post-start health/state and database comparison
verified clean revision and exact preservation: revision2182, minute163810,
life11 Jamie Moretti,2184 receipts, state SHA256
 a12b3718954540bb3cec26d7ad952c9c1ebacf876194939c46350b184e5ed0d4.
Backup:.runtime/pre-4daa6c8-campaign.sqlite3. No QA actions touched the live save.

Now live: closer card cameras, contextual G and fresh incident clock, dynamic
rental moves and ownership overview, independent NPC proprietors and funded
business accounts, interior shell/footer improvements and recorded named
indoor assassination routing. Broader goal remains open: generic NPC/contract
interior incidents, missing venue interiors, death-service businesses, fuller
billiards animation/physics, compact visual QA and fresh campaign acceptance.

### September 14 — funded funeral accounts

Previous goal turn was progress: promoted4daa6c8 with a byte-preserved campaign.
Rechecked the worktree; unrelated mugging/robbery/armed/aftermath/sim work remains
untouched. The next death-business step connects actual payment to ownership.

Ordinary deaths settle an affordable funeral from the dead person's purse and
savings, then available family cash. Full service260, external costs90 with
proportional smaller services; under25 means an unpaid burial. Only the margin
reaches the player/family/individual proprietor. The existing Buried marker
prevents duplicates. Player crew remain unburied pending the player's existing
choice; arranging their funeral now pays another parlour's proprietor without
charging twice. Existing saves receive no retrospective bills.

Focused funeral tests passed0.210s. Wider death/killing/succession/proprietor
integration results are recorded in.runtime/funeral-death-integration.log.
Added explicit estate/family split, all three owner types, poor means, repeat
settlement, crew choice/payment, savings use and no family overdraft coverage.
No live release or live-save QA commands;4daa6c8 remains on8791.
New morgue/cemetery/crematorium venues and their authored interiors remain open.

Final funeral validation: the added savings test initially failed because its
fixture wrote to a nil savings map. Initialized that fixture explicitly and
reran the full focused death/funeral/ownership selection successfully (see
.runtime/funeral-death-integration-final.log). The earlier integration selection
passed9.085s before that additional test was added. This is targeted simulation
coverage, not a new complete-suite or campaign acceptance claim.

### September 14 — death-service premises and funded supply chain

Previous turn was progress: funded funeral settlement committed and targeted
integration passed. Rechecked the worktree; preserved unrelated work.
Added Bellwether Mortuary, Oak Ridge Cemetery and Stillwater Crematorium along
the eastern map column, each with price, trade, staffing/supplies, trouble/repair
and income for retainer/grounds work. Standard deed and save-repair mechanisms
apply. Their existing funeral allowance now funds care plus one disposition;
only provider margins enter owner accounts. This does not levy an extra bill.
Closed, damaged, unstaffed or unsupplied premises cannot receive service trade.
Disposition is stable by NPC ID pending actual personal wishes/choice.

Targeted funeral, service, new-place/save-repair and ownership-eligibility tests
passed0.207s (.runtime/death-services-final-tests.log). They verify215 total
owner profit from a260 bill when all involved premises are player-owned,
separate burial/cremation custom, family and individual owner margins, closed
providers, duplicate prevention and old-save initialization. Distinct 3D assets
and interior staging for these addresses are still pending; currently they use
the generic work-building map representation. Nothing promoted to8791.

### September 14 — authored mortuary receiving room

Previous turn was progress: added funded death-service premises. Revalidated
full clean86dd0d0 suite session55129; still running. cmd/blackledger has already
reported missing painted fronts and fallback interiors for the three new
addresses. These are actual remaining release defects, not a passed full run.
Do not restart the still-live suite or promote this incomplete asset set.

Authored original Blender mortuary receiving room: six cold cabinet doors with
gaskets/hinges/latches, jade wall tiles, terrazzo floor, stainless empty trolley
and basin, oak registry desk/telephone/register, filing drawers, visitor chairs,
lettering and milk-glass task lights. Exported interior-mortuary.glb1241656bytes;
reproducible via .venv-blender/bin/python tools/export_city3d.py
--only=interior-mortuary. Registered room camera/cutaway settings and reserved
occupant positions. Only an actual public attendant/clerk gets the registry.
Recorded indoor strikes can now use this room through existing room routing.

Isolated fixture .runtime/mortuary-interior-qa.sqlite3 on8987 (session60850),
new mortuary-interior scenario. Browser96 G-entered and visually inspected at
1235x1053: six public occupants plus player clear of furnishings, registry
attendant behind desk, receiving cabinets/trolley and header readable. Current
room geometry and controls are in the main stage; no live8791 commands.
Production build passed4.59s before the attendant-only correction; final build
and frontend tests are recorded in.runtime/mortuary-final-build.log and
.runtime/mortuary-ui-tests.log. Expanded the existing indoor-strike placement
check to the actual mortuary GLB and added public registry assignment coverage.
Distinct exterior, fallback art, cemetery and crematorium interiors remain.

Final build passed5.64s. Full frontend run finished402passed/2failed: the
new registry test lacked its import (corrected), and city-grid reports shared
generic silhouettes for new service addresses. The silhouette failure is a
real pending exterior task and was not suppressed. Rechecked targeted staging
and indoor-strike tests after the import correction; results in
.runtime/mortuary-targeted-ui.log. No claim that all frontend tests passed.

### September 14 — distinct death-service exteriors and fallback plates

Previous turn made progress with the mortuary interior. Authored original
Blender mortuary receiving-house facade, cemetery grounds/lodge/gate piers/iron
rails/headstones, and crematorium with furnace chimney and memorial planters.
Registered all three models in city planning AND the explicit city load list;
Browser97 caught the omitted load entries (undefined clone) before acceptance.
The initial cemetery crossing had coplanar faces; separated the cross-path
surface and re-exported/rerendered after visually observing the dark rectangle.

All models and fallback JPEG plates come from the repository's authored geometry
and generated brick/roof materials. tools/death_service_exteriors.py and
--only=death-service-exteriors reproduce the GLBs;
tools/render_death_service_plates.py reproduces three fronts and mortuary room
plate at960x720. Visually inspected all three final fronts. Updated legacy
isometric silhouettes to represent the receiving house, headstones/lodge and
chimney. The initial direct node grid check had used stale compiled files;
the final npm test recompiled sources and passed all404 tests24.811s.
Final production build2.77s. Front-art coverage Go test passed0.193s.
Browser98 at1280x720 successfully loaded and Z-framed the mortuary, showing its
facade, lot/footpaths, entrance and city label; adjacent cemetery visible.

Full clean86dd0d0 Go run session55129 is now TERMINAL FAILED. sim passed415.830s;
core416.695s failed on missing new-address worker roles, manner descriptions,
counter material and missing three reach-test cases. cmd/blackledger failed on
new front/interior plates. Fronts and mortuary interior plate now supplied;
cemetery/crematorium interior plates and their 3D rooms remain, along with those
simulation/content integration tasks. Do not restart this terminal handle or
call this a passing full suite. No release promoted; live8791 remains4daa6c8.

### September 14 — integrate service venues into the living city

Previous turn was progress: distinct exteriors, rendered fallback plates and
all404 frontend tests passed. Rechecked the worktree and preserved unrelated
changes. Addressed the four core failures from the terminal86dd0d0 full run:
added mortuary/grounds/furnace worker roles, location-specific manner text,
trade-specific counter information, and three actual ordinary-death benefit
cases in the every-trade reach sweep. Tests now exercise owner profit from
Kill rather than directly invoking the payment helper.

Replaced the undertaker counter's inaccurate all-city-deaths claim with actual
saved receipts, also used by all three new providers. Bounded last256 service
records retain paid/margin/place/person/minute across reload and ownership.
No historical backfill. Counter summarizes its recent actual book; own-parlour
crew arrangements record costs without a margin. Each finalized death remains
exactly once through existing Buried handling.

Targeted prior-failure plus funeral, death-service and JSON-reload checks passed
0.947s (.runtime/death-integration-final-tests.log). Earlier initial integration
selection passed0.945s. These are focused checks, not a new full-suite claim.
Cemetery and crematorium 3D interiors/fallback room art remain release blockers;
no main8791 release or QA mutation occurred. Broad goal remains incomplete.

### September 14 — cemetery and crematorium interiors

Previous turn was progress: integrated workers and actual saved service books.
Rechecked the worktree; preserved unrelated work. Authored original Blender
Oak Ridge plot office (plot map, archive drawers, grounds tools, registry and
visitor benches) and Stillwater remembrance/furnace room (closed furnace doors,
flues/gauges, urn shelves, registry and benches). Registered room models,
camera/cutaway settings, public occupant positions and player entry points.
Both now participate in existing named indoor-strike room selection.
Export path: --only=memorial-interiors; fallback plates rendered from those GLBs.
Restored the unchanged older front/mortuary plates after the batch renderer
regenerated them, avoiding unrelated render-noise changes in this commit.

Initial browser99 cemetery/100 crematorium fixtures exposed a title obscured
by the plot map and more public people than reserved positions. Lowered the
map, expanded the clear aisle positions, re-exported/rebuilt. Browser101 final
cemetery check at1280x720 shows readable title and spaced occupants. Reviewed
crematorium browser scene and rendered fallback plate; expanded actual-geometry
checks cover both rooms with11 public occupants plus player, pairwise body
bounds and furnishing clearance at three heights. Existing indoor-strike
placement test now also loads both room GLBs.

All405 frontend tests passed17.879s; final build3.85s. Go front/interior plate
coverage passed0.199s. Fixtures:.runtime/cemetery-room-qa.sqlite3 on8988
(session44565) and crematorium-room-qa.sqlite3 on8989(session93862). These are
isolated staged fixtures, not campaign acceptance. No live8791 commands or
release promotion. Broader goal remains open, including remaining city venues,
full generic indoor incident choreography and continued billiards fidelity.

### September 14 — explicit family formation and headquarters foundation

Latest user amendment now recorded in docs/HEADQUARTERS.md and docs/GOAL.md.
Formation requires one owned usable business, 25 respect, independence and an
explicit action at the chosen address. No daily auto-formation. Headquarters
survives saves and family succession; NPC families replace lost bases while the
player chooses a replacement. Known family cards display the base. This is the
foundation, not the requested comprehensive crew dispatch system: named actors,
travel, reserved resources and operation-specific resolution remain to implement.

Focused formation/migration/lost-deed checks passed 0.442s; succession checks
0.163s; inherited-organization integration checks 3.909s. Frontend build 5.36s.
Isolated headquarters fixture on 8990 (.runtime/headquarters-qa.sqlite3), browser
102: entered laundry, formed through Business action, API revision 1 recorded
headquarters laundry, Families displayed Bluebird Laundry. No live-save action.

Earlier clean c123527 full suite terminated: sim passed 484.227s, core failed
routine predictability (87% vs 90%), cmd failed missing three front manifest
entries and an omitted receipt-list tag. These remain tracked failures, not a
passing release. Live 8791 remains the earlier 4daa6c8 build.

### September 14 — named combat actors before queued operations

Added stable Hand.ID resolution for original associates and signed family
members. Strike readiness, weapons, scene attacker, trust penalties and death/
capture now follow that actor. Fixed the existing fatal/captured strike branches
that cleared the entire associate roster. Legacy no-ID synchronous callers keep
their default; other crime adapters and the operations queue remain pending.
No new dispatch buttons have been exposed prematurely. Existing unrelated
robbery/mugging edits remain untouched.

Named regression tests reorder two associates, charge the selected member's
trust, reject held/departed members and exercise death/capture across 100 seeds
without losing the other associate or hurting the player. Broader selected
combat/delegation/scene checks passed 2.870s. First regression used adjacent raw
LCG states and exercised only death; corrected its seed distribution rather
than weakening the death-and-capture assertion. No live release.

### September 14 — death-service integration coverage corrections

Registered mortuary, cemetery and crematorium fronts in the existing fallback
manifest so the already-authored plates are actually reachable. Removed the
optional tag from the saved service receipt list to comply with the core list
contract; no public receipt disclosure was added. The three formerly failing
front/list/wire-shape checks pass 0.186s. Routine predictability remains under
investigation; no full-suite pass or release is claimed.

### September 14 — authoritative individual takeovers and fixture migration

Routine investigation found multiple living NPCs with “Runs The Mariner”: the
legacy ambitious takeover set a role/location but never the deed. It now writes
the individual owner, refuses an existing proprietor, preserves competing deeds,
and clears the actor's obsolete journey. Only trading businesses qualify. This
is the existing criminal takeover path, distinct from funded proprietor buys.
Regression verifies exclusive ownership and death returning the deed to market.

Additional older family test setup helpers now explicitly Incorporate before
daily updates; negative tests and daily simulation loops retain OrganizationDay
without formation. Selected takeover/proprietor/headquarters/crew checks pass
0.435s. Routine predictability still fails (86% after the ownership correction);
logs show housing moves, changes in work and actual takeovers disrupting the old
whole-fortnight metric. No threshold was lowered and no passing result claimed.
Temporary diagnostic tests removed; evidence remains in .runtime/routine-context.log.
Comprehensive queued dispatch and operations UI are still outstanding.

### September 14 — saved crew orders, first adapters and interior controls

Previous goal turn was progress: headquarters, named combat and ownership fixes
committed. Added saved named assignments for restock/assassinate with outbound,
work, return, result, recall and cash escrow. Actors use normal NPC journeys;
routines, ambition, guard selection and other person actions respect their
assignment. Player identity/location is never swapped. Last-seen addresses
prevent target tracking through hidden movement. Restock shares its existing
effect; assassination shares odds, consequences and actual attacker cues.

Interior Business menu and Families page expose authoritative named offers and
assignment register. Browser103, isolated .runtime/crew-orders-qa.sqlite3 on8991
(session62165): enter laundry, choose Russo Motor Works, dispatch Leo; cash
reserved130 (one dollar ordinary income during five-minute issuance), actor left
the room, outbound35min/Recall visible, player remained inside. Screenshot
review at1235x1053 exposed low disabled-button contrast, corrected in CSS.
Build2.88s after correction. Compact-width visual acceptance still outstanding.

Core/order/cue checks passed0.401s and HTTP retry/list projection checks0.198s
before the additional hidden-address regression. Final privacy and frontend
suite results recorded separately below when terminal. No main8791 actions.
Remaining scope: robbery/bombing/all practical adapters, equipment/cargo budgets,
full family-succession continuation, broader unavailable-actor integration and
complete UI/scene acceptance. Do not treat two adapters as full dispatch.

Prior clean25b8530 full run session9082 is terminal FAILED: core501.113s, sim
528.113s. Failures: routine86%; room staff count3 vs named0 onday10; taking a
family no longer implicitly incorporates; HTTP progression's fixed laundry
acquisition encounters an existing owner; long publicans buy no stock. These
predate this order implementation and remain release blockers. No restart or
promotion based on that run.

Final order/hidden-address/core checks passed0.414s; HTTP idempotency and wire
list checks0.217s (.runtime/crew-orders-privacy-tests.log). All405 frontend tests
passed19.982s (.runtime/crew-orders-frontend-tests.log). Cancellation leaves a
living operative's current journey intact, preventing a dismissed/new-life
actor from teleporting back to their street-leg origin. This final small
adjustment requires the focused cancellation checks again before release.

Final cancellation regression explicitly retains the living actor's destination
after new-life cancellation; focused core0.192s and HTTP0.152s pass. Diff
whitespace check passes. Goal remains active; no release promotion.

### September 14 — named bombing orders and charge custody

Previous turn made progress on saved orders. Added bombing to the same catalogue
and lifecycle. One charge is reserved once, persisted through reload, consumed
on an actual attempt and returned once after recall/invalidated target. Captured
or dead operatives lose the charge; cancellation across issuing lives never
credits the new player. Arrival and work commitment recheck the deed/destruction.
Shared resolvePlant preserves ordinary player RNG/consequences; named attempts
use operative odds and injury/death, with the correct attacker on blast cues.
Successful planters are excluded from collateral casualties rather than remaining
in the victim pool while the scene shows them getting clear.

Regression covers both outcomes across50 seeds, player position/health unchanged,
charge accounting, recall through clone/reload, changed deed and custody. Selected
order/demolition/cue/HTTP checks pass (.runtime/crew-bomb-final-tests.log); initial
expanded core selection1.093s. Frontend build2.90s. No full-suite pass claimed.
Browser104 isolated .runtime/crew-bomb-qa.sqlite3 on8992(session17306): Families
order selector showed bombing,170min andonecharge terms; dispatch recorded Leo
outbound20min after issuance and exposed Recall. Corrected dispatch log to show
the place's actual name rather than internal ID observed in that check. Full
end-to-end bombing scene and compact-layout acceptance remain outstanding.

Unrelated dirty robbery/mugging/aftermath files preserved. Main8791 untouched.
Robbery, other practical adapters and known broad campaign failures remain open;
this is the third adapter, not completion of the requested operations system.

### September 14 — campaign integration after ownership and headquarters

Previous turn made progress on named bombing. Rechecked worktree and retained
unrelated changes. Repaired three failures from the earlier full run:
NPC criminal takeovers now immediately assign named workers, like player deed
acquisition, because they occur after BusinessDay. EmptyChairs no longer invents
a full staff count before checking a hiring cooldown/unpaid wages; reopening
hires real people. Taking over an existing family explicitly incorporates the
transferred deeds and retains its headquarters, rather than relying on the
removed daily auto-formation. The existing takeover test now checks that base.

Core staffing/takeover/counter selection passes2.959s. Added targeted closed-hiring
regression. HTTP full rise/fall/new-life route now shops for an available deed
and rechecks after travel rather than assuming laundry stays unowned; its final
inheritance assertion follows the actual acquired deed. Passed0.878s. Initial
rerun caught that remaining hardcoded laundry assertion and was corrected.
Publican first-business search now reads the already-public owner and excludes
held deeds. TestThePublicanActuallyRunsWhatItBuys passes1.712s; an earlier
TestPublican regex selected no tests and is not counted as validation.

Routine predictability and long-publican manual-restocking coverage still have
no passing rerun. Other jobs, full scene/layout checks and broad original goal
remain incomplete. Live8791 was not changed or used for QA. A fresh clean full
run will check the accumulated headquarters/job changes and these fixes.

### September 14 — visible, consistently priced manager restocking

Previous turn was progress on ownership/staffing and explicit leadership.
The clean c06c36c full run remains live (session47957); current output has a
Families prop-order source-check failure, repaired by keeping world/actions
adjacent without changing behavior. Focused source-check passes0.284s.

Managers previously debited the undiscounted trade price and refilled supplies
without any receipt, despite the comment promising shared purchasing rules.
They now use Restock readiness/payment/effect: the player's haulage discount
applies and the ledger records actual spending. Managers in custody or on a
headquarters assignment cannot simultaneously purchase supplies. Regression
verifies exact discounted debit, one receipt, no duplicate full-stock charge
and no purchase from custody. Focused manager checks pass0.389s.

The campaign report now counts stock-purchase records in public last-result
receipts (observed_restocks). Long-campaign coverage uses these purchases,
including managers and operatives, instead of requiring the protagonist to
press Restock personally. It still requires actual purchases and wage-setting.
The short publican check now asserts observed receipts and passes5.024s. This
metric counts observed receipt entries, not a complete lifetime accounting book.

Initial broad manager/charge/running selection remains live session7693, log
.runtime/manager-restock-tests.log; no completion claimed. A separate clean long
publican run will validate the revised purchase measurement. No live8791 action
or promotion; remaining campaign, job and visual scope stays open.

### September 14 — temporary errands preserve workplaces

Previous turn made progress on manager purchasing and observed campaign
receipts. Investigation confirmed Arrivals treated nearly every daytime
arrival as a new job: petrol, vehicle repair/purchase, grudges and collections
could overwrite Post. These temporary errands now retain the existing post.
Service visits finish their existing purchase/repair effects and reconsider an
ordinary return journey from the actual counter, without teleporting. Named
headquarters journeys and home journeys retain their prior exclusions.

Tests cover all six temporary errand labels and a funded fuel purchase followed
by an actual delayed return to the original job. Selected fuel/garage/forecourt,
crew-order and morning/evening mobility checks pass0.387s. The broader same-hour
routine metric remains86%; its failure was reproduced and not weakened. This
change fixes a concrete workplace bug but does not claim that entire metric.

Full c06c36c test session47957 and clean9aec451 long publican session57475 both
re-polled as live; no restart. Prior broader manager/charge/running selection
completed successfully139.025s (.runtime/manager-restock-tests.log). No main8791
actions or promotion. Broader goal and comprehensive dispatch remain open.

### September 14 — named robbery orders and loot custody

Previous turn was progress on routine detours. Extended the shared robbery
handler to a named saved order, including arrival/commit checks, 60min work and
loot held until return. Player location/health are not substituted for the actor.
Death/capture loses carried loot; surviving operatives retain it if authority
ends. Faction/personal-owner payouts now debit actual bounded funds instead of
paying more than an owner has. Unowned daily trade retains its existing modeled
takings. Personal proprietors remember the robbery. Shared robbery attempts
now emit a scene cue with their actual actor and neighborhood crime pressure.

The robbery file had pre-existing armed-resistance additions. Kept an exact
patch in .runtime/pre-crew-robbery.patch; staged the HEAD base plus only this
task's edits via index blob. Verified those unrelated added lines remain
byte-for-byte in the unstaged diff. Other unrelated files untouched.

Named loot tests cover delayed payment, clone/reload, no second payout, custody,
new ownership and personal proprietor funds/grievance. Core selected robbery/
actor checks1.336s and HTTP0.211s before final proprietor branch; frontend
build2.92s. Final proprietor and clean-check evidence follows when complete.
Actual complete robbery scene choreography and compact UI checks remain open.

Prior full c06c36c session47957 terminated FAILED: core542.498s, sim574.874s;
only routine86%, the now-corrected Families prop-order check and old manual
restocking metric failed. Clean9aec451 long purchase-observation run session57475
terminated PASS232.821s. Do not restart these terminal handles. No main8791
actions/promotion; comprehensive jobs and broad visual/game scope incomplete.

Robbery follow-up: personal-proprietor regression passed0.195s. Clean07ce79b
checkout excluding unrelated local edits passed selected core1.370s and
HTTP0.213s (.runtime/crew-robbery-clean.log), terminal session54541. Remaining
unstaged robbery diff is exclusively the pre-existing armed-resistance work.


### Headquarters property maintenance adapters

Added named repair and remedy assignments to the saved dispatch catalogue and
operation picker. Both reserve their actual budget, travel, work for 60 minutes,
and return. Deed changes, already-completed work and recall refund the unused
reservation. The work shares personal repair/remedy effects without charging a
second time. Owned residences can receive repair orders as well as businesses.

Evidence: `.runtime/crew-property-tests.log` passes targeted crew, repair,
remedy, HQ and workplace UX checks (core 0.413s, HTTP package 0.207s).
New tests cover exact reservation with an empty remaining wallet, JSON reload
mid-work, exactly-once settlement, lost deeds before arrival/during work,
already-resolved problems and working recall. `.runtime/crew-property-build.log`
passes TypeScript/Vite (3.00s; existing bundle-size advisory).

Not promoted to port8791. Remaining: broader practical job coverage, custody
and succession continuation, compact-layout/remote-scene review, the previously
recorded routine-predictability full-suite failure, and fresh campaign acceptance.
The broad visual/game goal remains active. Unrelated armed-resistance and
scene-aftermath work is preserved unstaged.


### Assignment funds follow family succession

Previous goal turn made authoritative progress (f3d8ded property orders). This
turn fixes custody of assignment money across player-family succession: sworn
members recall unfinished work, return naturally, and settle reserved funds or
already-committed robbery proceeds into the successor faction. The saved estate
recipient prevents a new protagonist receiving the funds. Captured/dead carriers
lose their funds; a surviving carrier retains them if the family dissolves.
Personal associates remain outside sworn-family inheritance. Unused demolition
charges are removed from circulation, pending broader equipment custody work.

`.runtime/crew-succession-final.log` passes the focused crew/HQ/succession tests
(core 0.440s, command package 0.225s). New cases cover outbound, working and
returning stages, reload plus new life, completed robbery proceeds, loss by death
or capture, dissolved family and recall of an unfinished assassination. No live
campaign was mutated or release promoted. Full offensive-order continuation
under successor authority remains unfinished, along with broader job coverage,
visual acceptance and the previously recorded full-suite routine failure.


### Compact headquarters overlay review

Previous turn progressed with successor-family assignment custody (1e24bf2).
Current isolated browser review at1280×720 found the three stacked selectors
pushed the dispatch button below the interior overlay. The form now places
operative and operation side by side, retains a full-width target, and tightens
spacing on short windows. Screenshot review shows the quote, dispatch control
and disabled reason together while the laundry scene remains visible. Assignment
history now keeps active orders first, folds completed reports behind a counted
toggle, and uses readable operation/stage labels. The register is keyboard
focusable with a visible focus indicator.

QA: fresh `.runtime/hq-layout-qa.sqlite3`, server8993 (session62494), browser107
`?layout=compact`. Browser interaction dispatched Leo to restock Russo Motor
Works: cash6000→5871 (130 reserved plus1 ordinary income), time08:00→08:05,
busy reason and1 active assignment displayed. No actions taken on8791.
`.runtime/hq-layout-final-build.log` passed TypeScript/Vite in2.58s before the
final focus-only accessibility change; final build recorded separately below.
Desktop widths below1280, complete report-toggle interaction and other interior
panels still need acceptance; this targeted check is not full UI acceptance.

Final `.runtime/hq-layout-reviewed-build.log` passed in2.70s, including keyboard
focus changes. No live release promoted.


### Routine acceptance follows changes of home and work

Previous turn progressed with verified compact HQ layout (50e1966). Traced the
remaining routine-test failure: the old measure mixes fourteen days of changing
homes, workplaces, roles and family positions into a single expected address.
The settled-routine test now groups consecutive unchanged contexts and requires
four distinct days at each sampled hour. It retains the90% reliable person-hour
threshold and75% modal-address criterion, excludes residents who never move,
requires at least12 repeated person-hours per moving resident in aggregate, and
fails if the world stops advancing during sampling. Housing/job changes remain
fully simulated rather than frozen to satisfy the check. Other city movement,
evening, stable-haunt and wartime tests remain separate.

`.runtime/routine-settled-tests.log` passes (core0.793s):2307 repeated person-hours,
81 moving residents,211 contextual episodes. Diagnostic raw evidence remains in
`.runtime/routine-episodes.log`; temporary diagnostic source removed. This is a
measurement correction, not a claim that all NPC routing is perfect. Fresh clean
full-suite validation is next; no live campaign mutation or promotion.

Clean checkout `.runtime/check-settled-routines` verified empty git status at
bc169160cdf07255e15f3916152f3f226f1f2d61. Full `go test ./... -count=1 -timeout 20m`
is running as session22511, output `.runtime/settled-routines-clean-full.log`.
Poll that handle to terminal before assessing or restarting it.


### Poker close-camera browser acceptance

Previous turn progressed with settled-routine validation and clean suite launch.
Fresh `.runtime/poker-readable-qa.sqlite3` on8994 (server session3443) loaded a
funded hand at Saint Agnes. Browser108,1235×1053 screenshot review verified the
main-scene poker table, readable10♥/5♥ player cards, and6♥/J♠/J♦ flop after a
real Check command. Pot48→168, player stack588, and call48 correspond to the
public hand state. Whole-table framing showed the actual seated roster while
keeping the cards present. A folded player had no3D cards but was incorrectly
labelled Two hidden cards below; changed that to Cards folded. Browser109
`?review=folded` confirmed the rebuilt summary, while active opponents retained
hidden-card labels. `.runtime/poker-readable-build.log` passed in5.53s.

This proves this desktop hand/view, not full poker visual acceptance. Character
fidelity, compact-window review, other card games and full scene choreography
remain open. No8791 actions or release. Clean Go full-suite session22511 remains
running (polled this turn); do not restart. Frontend test session37161 also
started in the clean checkout; initial optional dependency-link setup used an
incorrect relative path, so inspect the test's terminal result before retrying.


### Guard custody, location and family continuity

Previous turn progressed with poker browser verification and folded-card fix.
Found and fixed two guard invariants before broader HQ guard dispatch: PostedAt
validated every guard against the current protagonist rather than the deed's
family, and defence/raid-victim selection could count a living guard standing
at another address. Inherited guards now remain attached to their family deed;
only guards physically present defend or get caught there. Unpost checks owner
authority, and guide completion counts only the player's own posts. Public
posting and premise text distinguish arrival toward the post from other absence.

`.runtime/guard-custody-final.log` passes focused posting, succession, guide and
raid-defence checks (core4.536s). Added cases for successor/new-life retention,
no new-player dismissal or guide credit, loss of deed, absent guard defence/raid
exclusion, public absence and resuming protection on return. No live mutation or
promotion. Named HQ guard dispatch and fuller guard scheduling remain open.

Clean frontend session37161 is terminal PASS:405 tests,142.494s, in
`.runtime/settled-routines-clean-frontend.log`. Dependencies resolved via the
parent checkout despite the earlier failed optional symlink setup. Go full-suite
session22511 remains active on clean bc16916; its result will not cover these
later guard changes. Continue polling the same handle.


### Named headquarters guard orders

Previous turn progressed with physical guard/succession fixes (9ca2ad7). Added
crew_order:guard through the standard saved assignment catalogue: named signed
member, owned business, journey,45-minute setup, persistent guarding stage and
recall home. Guarding is active for busy-state exclusion but has no due clock,
avoiding one-minute boundary polling. Recheck deed/post availability at arrival
and commitment. Recall, relief, destroyed/lost premises and member loss clear
only this operative's posting. Established posts persist into successor-family
authority. UI lists On guard among active orders and quotes time to take the post.

`.runtime/crew-guard-final.log` passes crew/HQ/posting checks (core0.464s,
command package0.210s). Added travel/no-remote-defence, setup, routine exclusion,
reload/recall, succession, lost deed before and after setup, and departed-member
cases. `.runtime/crew-guard-build.log` passes TypeScript/Vite4.29s. Browser
acceptance of this new operation remains pending, as do temporary hired-associate
guarding and broader practical order coverage. No live release or save changes.
Clean full Go run22511 remains active on bc16916; do not restart it. That run
predates later guard changes. Clean frontend37161 is already terminal PASS405.


### Hired associates can guard, with browser acceptance

Previous turn progressed with saved guard orders (4c395cd). Removed the signed-
member-only restriction for HQ guarding; existing hired associates can now take
posts without being forced into family membership. Their defence/public guard
reliability uses paid loyalty, while sworn successors use trust. An employee's
personal contract ends on employer death; signed-family posts retain succession.
The prior departure test now removes both employment and membership explicitly.

Fresh `.runtime/guard-dispatch-qa.sqlite3`, backend8995/session18136, browser110:
Leo received Guard Bluebird Laundry, entered setup with41 minutes remaining,
and reached On guard after an hour's action. The building showed38 defence;
Recall immediately removed protection and entered Returning. This exposed a
same-address phantom street journey; fixed crewOrderJourney to keep the operative
present during its one-minute dispatch/report processing. This final fix has
regression coverage but was not rebuilt into that browser server yet.
`.runtime/associate-guard-final.log` passes core crew/HQ/guard tests0.428s.

Clean full Go session22511 is terminal PASS, explicitly polled: core602.792s,
sim630.138s, all other packages passed in`.runtime/settled-routines-clean-full.log`.
It tested clean bc16916, before subsequent guard changes; it must not be reported
as a full pass of the current HEAD. Clean frontend37161 previously passed405.
No live promotion; full latest integration/campaign/visual acceptance remain.


### Shared availability for legacy guard and delegation paths

Previous turn progressed with hired guarding and confirmed clean full-suite pass.
Legacy direct crime delegation now uses named-hand readiness, including existing
guard posts and named legacy errands. Legacy guard nomination excludes actual
travellers and named errands, and posting replaces an unstarted routine journey
so a newly posted guard does not immediately leave on an old schedule. Existing
HQ queue checks remain shared rather than bypassed. No additional loyalty floor
was imposed on legacy guard nominations.

`.runtime/guard-availability-tests.log` passes0.409s. Added legacy-post→delegation,
relief, traveller/task→post nomination checks; expanded final run additionally
covers pending routine cancellation. Final test handle63917 targets Door/Post/
Delegation/Crew cases; inspect terminal result before starting full integration.
No live release or campaign changes.

Final availability test63917 is terminal PASS, core10.843s, recorded in
`.runtime/guard-availability-final.log`.

Clean full integration started atbf4dd93 in`.runtime/check-guard-integration`,
`GOMAXPROCS=4 go test -p 2 ./... -count=1 -timeout 25m`, session16200,
log`.runtime/guard-integration-full.log`. Poll this handle before restarting.
Prior clean full session22511 and frontend37161 are terminal PASS; do not poll them.


### Campaign-copy release preparation found same-version HQ migration gap

Previous turn progressed with shared guard availability. Main8791 health still
reports clean4daa6c8. Created isolated SQLite backup
`.runtime/guard-release-campaign-copy.sqlite3`; baseline manifest records raw
state SHAa12b3718954540bb3cec26d7ad952c9c1ebacf876194939c46350b184e5ed0d4,
2184 receipts, revision2182, minute163810, life11, version19. Clean bf4dd93
candidate8996/session54324 built in`.runtime/check-guard-integration` (frontend
`.runtime/guard-release-build.log`3.04s). Startup added mortuary/cemetery/crematorium
properties and empty crew_orders/death_service_receipts; all other top-level
state, including player/time/revision, and ordered receipt bytes were unchanged.
The raw-state checksum therefore changed as expected; no live mutation occurred.

Candidate public state exposed an existing family with blank headquarters:
SettleHeadquarters was gated behind a save-version migration even though current
version19 saves predate the feature. Store decode now settles missing bases
regardless of version. Full store tests pass0.401s in
`.runtime/headquarters-upgrade-store-tests.log`; new cases cover missing base,
read-only DB behavior, explicit choice/lost deed preservation and no automatic
formation of independent property owners. Candidate8996 predates this fix.
Clean Go integration16200 remains running on bf4dd93; do not restart it, and do
not treat its result as covering this later store fix. Live promotion remains
pending; the broad goal and missing gameplay/visual acceptance remain active.

Clean f52edbc store tests also pass0.470s; build session73988 is terminal PASS.
Isolated upgraded-copy server8997/session78482 uses that clean binary. Public
Jamie Moretti family headquarters now resolves to butcher (Fassano Meats).
Against the pre-fix candidate, only factions changed; player, clock, revision
and every ordered receipt are equal. Evidence:
`.runtime/headquarters-upgrade-candidate-evidence.json`. No main save written.


### Existing campaign management and request retry acceptance

Previous turn progressed with same-version HQ migration and isolated upgrade
comparison. Exercised the upgraded campaign copy on8997 through normal API
commands: travel to Fassano headquarters, observe authoritative no-dispatch
reasons (low loyalty or existing guard duty), travel to chapel, relieve Anton
Iordan's existing post, return to headquarters, dispatch Anton to restock Fassano
for200 reserved, then let two ordinary hours pass. Order reached done with
Supplies delivered, zero reserved funds, and45 supplies at the premises. The
identical dispatch request was retried and returned byte-identical output with
one active order. Evidence `.runtime/campaign-copy-dispatch-evidence.json`.
This is a bounded existing-save management sequence, not a full20–30-minute
campaign acceptance. Copied campaign advanced2182→2189; live raw state and all
2184 ordered receipts remain equal to the original copy baseline.

Clean f52edbc frontend session32557 is terminal PASS405 tests in38.449s,
`.runtime/headquarters-upgrade-frontend.log`. Full guard integration16200 remains
active on bf4dd93; last polled this turn. This run predates the separately tested
store migration fix. Do not restart it. No release promoted, and broader job,
interior, billiards physics and visual acceptance requirements remain open.


### Release validation and remaining interior coverage inventory

Previous turn progressed through an existing-save dispatch/retry/settlement
sequence. Clean f52edbc API and store run72778 is terminal PASS: command package
2.749s, store0.443s, frontend build3.99s. Logs:
`.runtime/headquarters-upgrade-api-tests.log` and
`.runtime/headquarters-upgrade-release-build.log`. Full integration16200 remains
active on bf4dd93 (processes confirmed live this turn); no restart or promotion.

Reviewed all33 public addresses against interiorSettings and on-disk GLBs:
22 mapped rooms,11 older painted/SVG rooms.24 interior GLBs include2 private
variants. Exact missing address inventory and per-room completion requirements
are now in docs/VISUAL_ACCEPTANCE.md. This is a scope audit, not visual approval
of the22 existing rooms; it prevents3D gambling-table views being mistaken for
complete venue interiors. Broader goal remains active.


### Blue Hour authored interior and completed guard integration

Added original Blender casino room source and reproducible single/full exports,
registered the GLB, staged cashier/slot/table/lounge occupants deterministically,
and routed entrance movement through a tested clear lane. Added isolated
casino-interior QA fixture. Browser8998 at1235×1053 shows the full room with six
actual occupants and unobscured controls. Geometry tests cover both rigs, all12
placements, seat support and40 entrance samples; PASS1 test in2.697s,
`.runtime/blue-hour-placement-tests.log`. Frontend build PASS2.59s in
`.runtime/blue-hour-final-build.log`. Original asset provenance is the checked-in
Blender code. Compact and incident review remain outstanding.

Full clean bf4dd93 guard integration session16200 is now terminal PASS: all Go
packages, core530.209s and sim499.772s. Final output is in
`.runtime/guard-integration-full.log`; do not restart or poll this completed run.
It predates the separately passing same-version HQ store fix and this room.
No main campaign QA writes or live release in this increment. Broad goal active.


### Headquarters budget custody correction

Found a mismatch between forfeited loot/charges and remotely refunded unused
job cash when an operative died or was captured. Reserved funds now follow the
carrier consistently; missing, dead or imprisoned operatives clear all carried
order balances without crediting player, estate or personal purse. Existing
living-return/recall and successor settlement behavior is preserved. New tests
exercise outbound, working and recalled return legs for death, capture and
missing actors, with reloads before and after cancellation to check exact-once
settlement. Targeted crew/succession tests PASS0.441s, session69063 terminal,
`.runtime/crew-custody-budget-tests.log`. Updated the current headquarters
implementation summary; remaining command coverage is still explicit.


### Named headquarters collection rounds

Previous goal turn made authoritative progress with the Blue Hour room and
budget-custody fix. Added collections to saved named headquarters orders and
its UI selector. Each owned usable business supports one outstanding round;
legacy task rounds cannot overlap headquarters collections. Work retains the
existing120-minute/$65 supplemental activity, adds real outbound/return travel,
and keeps payment in order loot until the operative returns. This is modeled
work revenue, not another withdrawal of already-paid business income. Rechecks
at arrival/completion cover ownership, destruction and trouble. Recall before
completion earns nothing; carrier loss forfeits proceeds through shared custody.

Tests cover actual Execute with a signed member other than the first associate,
public quotes, legacy/duplicate exclusion, reload before work/return, exact-once
payment, recall, changed deed, ruin, trouble and death/capture after earning.
Targeted core crew/collection/succession tests PASS0.428s (25118 terminal),
`.runtime/crew-collections-tests.log`. Final frontend build PASS2.57s (98875
terminal), `.runtime/crew-collections-build.log`; git diff check passes.
Browser acceptance of this selector addition and full clean integration remain
pending. No main save writes or local release. Broader goal remains active;
legacy standalone collections still exist for pre-family progression and need
further lifecycle consolidation, while additional practical adapters remain.


### Collection dispatch browser acceptance and button contrast

Previous turn implemented named collections and passing targeted tests. Created
clean checkout `.runtime/check-collections-c36e51a` at c36e51a, built binary and
frontend, and created fresh crew-orders QA save. Server8999/session82471 serves
that isolated build. Browser112 at1235×1053: entered laundry with G, selected
Make collections/Russo Motor Works, saw200-minute quote and65-on-return terms,
dispatched Leo, observed08:05 and outbound order with35-minute remaining leg.
The register and actual departure reflected the command; player stayed at HQ.

The enabled dispatch button looked like a grey disabled control. Added explicit
dark-green enabled styling, cream type, inset border, hover and distinct muted
disabled state. Root frontend build20923 terminal PASS3.27s; log
`.runtime/crew-dispatch-contrast-build.log`. Separate clean-backend/root-frontend
preview9000/session20532 uses another fresh fixture. Browser113 at1280×720 shows
all collection selectors, payment terms and the enabled button above the panel
fold, with unobscured header and visible interior. This verifies those desktop
sizes, not phone or every operation/long register layout.

Clean full Go integration session17026 is still RUNNING, freshly polled after
browser review. Command GOMAXPROCS=4 go test -p2 ./... -count=1 -timeout25m;
log `.runtime/check-collections-c36e51a/.runtime/full-go.log`. It covers c36e51a
before this cosmetic change. Do not restart on observation timeout. Main8791
and its campaign have not been modified. Broader goal remains active.


### Casino attack staging and collection return acceptance

Previous turn verified collections dispatch and improved button contrast. The
same isolated8999 campaign now completed its collection order through four
ordinary wait API commands: working at545/605, returning at665 carrying65,
done at725 with zero carried loot. Evidence
`.runtime/collections-browser-settlement.json`; main campaign untouched.

Added casino-strike fixture through the actual Strike handler, and actual-GLB
coverage for all four assassination variants in Blue Hour. Initial browser replay
showed actors alive at the start inside the casino, but staging chose a rear
corner beside slot machines. Placement now ranks valid candidates by the cast's
midpoint distance from room center before testing the full approach/fall. It
retains the floor/furniture safety checks and avoids preferring corners solely
because of iteration order. Home/venue tests PASS6 in7.892s, build PASS4.26s;
logs `.runtime/casino-strike-tests.log` and `.runtime/casino-strike-build.log`.
Session70630 terminal PASS. Fixture creation81626 terminal PASS.

Isolated9001/session6787 uses clean c36e51a backend and current frontend with
`.runtime/casino-strike-qa.sqlite3`. Browser115 at1280×720 visibly shows central
approach with both actors standing, then victim fallen on clear floor and the
attacker withdrawing. Browser114 prior1235×1053 screenshot documents the old
rear-corner framing. No outside substitution or load error. This is bounded
casino revolver playback evidence, not all-venue or all-variant visual approval.
Clean Go integration17026 remains RUNNING, freshly polled; do not restart.
Broad goal, remaining room production and other job adapters remain active.


### Owned apartment accounts at the front of Market

Previous turn improved casino attack staging and verified collected proceeds.
Reviewed actual RentalMoveDay: daily rotation of households checks affordability,
upgrades/downsizes, vacancies, condition>=60 and no trouble without transferring
deeds or displacing tenants. Occupancy is conditional, not guaranteed. Added a
prominent ownership overview before broker listings: grouped by building with
owned/tenanted/vacant counts, scheduled rent, today's actual collections, current
tenant arrears and broker offers. Links filter the existing deed register and
inspect the address. Existing public projection supplies all values; no API or
simulation changes. Added apartment-portfolio fixture with six owned units and
explicit synthetic rent-account values for visual QA.

Build7362 terminal PASS5.64s, `.runtime/apartment-portfolio-build.log`. Fresh
fixture30542 terminal PASS. Server9002/session22719 serves clean c36e51a backend
and root frontend against `.runtime/apartment-portfolio-qa.sqlite3`. Browser116
at1235×1053 shows two clear account cards,6 deeds/2 vacancies/$12 collected/$36
arrears, matching fixture accounts. This is layout evidence, not a live-campaign
rent report. Main save untouched. Full clean Go test17026 last poll still live;
do not restart. More housing-management controls and broad goal remain open.


### Rent recipient attribution after purchasing an occupied flat

Previous turn added and visually verified the grouped ownership overview.
Found that rent_paid_today used the tenant account without checking who received
it, so a newly purchased flat could claim the former landlord's payment. Rent
accounts now retain paid_to; public totals count only actual player receipts.
Current-day legacy positive receipts without recipient identity return null and
individual rows explain the missing attribution. Unknown values are excluded
from totals. No retroactive money transfer or inferred previous ownership.
Updated synthetic portfolio fixture to name its receipt recipient.

Tests exercise old-owner collection, actual BuyApartment, reload, next-day
payment and missing legacy recipient. Targeted rent/move/investment tests
PASS0.270s; build PASS7.67s, session43564 terminal. Logs
`.runtime/rent-recipient-tests.log` and `.runtime/rent-recipient-build.log`.
Clean checkout `.runtime/check-portfolio-1cba625` build54155 terminal PASS;
full frontend86445 remains RUNNING (freshly polled), log in that checkout's
`.runtime/frontend-tests.log`. It predates this accounting correction. Clean
c36e51a full Go17026 remains RUNNING at last poll. Do not restart these handles.
Prior portfolio browser116 was no longer available when revisited; no browser
navigation acceptance claimed this turn. No live release or main save writes.

## 2026-09-21 — source publication and distribution foundation

Prepared PolyForm Noncommercial licensing with required attribution (the user
requires no free commercial use), cross-platform browser-game archives, isolated
per-user desktop saves, launchers, dependency notices/checksums, and GitHub CI
with five native package jobs and tag-triggered draft prereleases. The license
means source-available, not OSI open source. Existing uncommitted gameplay and
visual work was preserved and is not included in the release-setup commit.

All five targets cross-compiled; native Apple Silicon archive launch, assets,
command/retry, restart/receipt persistence, per-user save and occupied-port
checks passed. City and Saint Agnes interior rendered in the browser. A 100-command
isolated API run had no invariant failures. Frontend 407/407, production build,
Go vet, server/save race tests and workflow validation passed. Initial full Go
run timed out in the pre-existing untracked 260-campaign diagnostic; core passed.
Gitleaks reported no leaks across 919 commits. Detailed evidence and remaining
publication/acceptance gates are in docs/RELEASE_READINESS.md and docs/RELEASING.md.
Nothing was published; media permission and GitHub authentication are unresolved.
The main campaign and port 8791 were not used or modified by these checks.

### 2026-09-21 — verify a clean release checkout and narrow media review

Created a detached QA worktree at b10f7b8, preserving the original dirty worktree.
Fresh npm install/build, 407 frontend tests, Go short suite across all packages
(core100.6s, sim313.1s), clean-tree versioned Mac packaging and native archive
restart/retry checks passed. build.json confirms b10f7b8 with modified:false.
This removes dependence on uncommitted local work from the packaging evidence;
it does not replace full/native hosted CI or campaign acceptance.

Added docs/MEDIA_RIGHTS_REVIEW.md with exact hashes and history evidence for23
files needing source/permission records. The four ground photographs have no
recorded provider/license; the imported shop's commit identifies Grok provenance.
No asset was removed or rights declared without evidence. Packaged linked docs
and exercised the actual Unix launcher from a path with spaces. Follow-up archive
smoke and document membership checks passed. GitHub integration was offered but
is not confirmed installed/connected; no publication or main-save access occurred.

### 2026-09-21 — fix amount controls found during release campaign

An isolated browser campaign reached first ownership through jobs and trade.
Fixed trade quantities incorrectly rendered as dollars/stepped by five, using
public goods units with one-unit steps. Fixed wire offer eligibility below the
legacy $500 lot while retaining the $100 minimum and command validation. Browser
verification covered five-crate buy/sell, a $100 wire, midnight upkeep, business
purchase/income and persistence across restart. 409 frontend tests, production
builds, server tests and focused banking/deposit tests (including race) passed.
Detailed evidence and limitations are in docs/RELEASE_READINESS.md. This session
included debugging pauses and does not sign off the full campaign acceptance.
The main save and unrelated working changes were preserved.

### 2026-09-21 — verify archive backup restoration

Extended native smoke checks to restore stopped-game saves and receipts from a
replacement game folder, advance the restored copy, and ensure the original
save is unchanged. Mac ARM64 passed on an isolated build of 11c5eb9. Added pinned
actionlint 1.7.12 workflow validation; it passes. No hosted/native Windows/Linux
run is claimed. See release readiness for exact evidence and remaining gates.

### 2026-09-21 — integrate remaining work on main

Reviewed the pending armed-resistance rules, death-cause diagnostic, police
arrival presentation and legacy rent receipt regression. Added focused coverage
for fatal personal robberies, absent-player safety and armour. Updated action
copy to disclose lethal resistance. The death diagnostic now skips with -short;
its complete 260-campaign report was not rerun in this integration pass.

The arrival test exposed an immediate jump to the end of the traffic route.
The first reservation now starts at zero; subsequent frames permit movement,
and officers wait for the car to arrive. The traffic integration test confirms
completion and officer visibility. This is automated scene coverage, not a new
visual/browser approval. All 410 frontend tests and the production build pass;
focused armed-resistance race tests and Go vet pass. The simulation diagnostic
compiles and its explicit short-mode skip passes. The full core and store suites also passed. Logs: .runtime/remaining-*.

### 2026-09-21 — connect published repository

The owner pushed main to the public JaTochNietDan/BlackLedger repository.
Verified remote main at 37dfc3d and corrected the previously proposed repository
name in attribution, download links and release docs. GitHub initially reported
zero registered workflows/runs despite build.yml being present. A follow-up
push will exercise the branch build trigger. Public availability does not by
itself resolve the media provenance questions or prove native build acceptance.
