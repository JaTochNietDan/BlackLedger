
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
