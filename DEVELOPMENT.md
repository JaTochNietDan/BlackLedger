
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
