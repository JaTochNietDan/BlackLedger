
### September 13 — supplied effect recordings

Copied ten user-provided WAVs unchanged, with provenance and checksums. Connected
explosion, raid, arrest siren, authored door contact, shooting casualty pain and
explosion panic. Added nearby persistent fire/vehicle idle loops and one approach
sample per rendered vehicle journey; loop seams crossfade decoded copies. Imported
escape clip awaits moving getaway choreography. Explosion/raid windows now 14/10s.
Browser 8873 private previews verified ten decoded assets, simultaneous explosion /
panic / fire, cancellation to zero active voices, and raid recording playback;
revision 0/minute480 unchanged. Door-kick did not trigger in this Saint Agnes
fixture; the cause and contact path still need an
unobstructed browser test. Vehicle/pain playback also needs targeted browser QA.
Audio unit checks pass, including tail duration, stereo seam and mute/restart.
Frontend build passes. No subjective listening or production-ready claim.

### September 13 — occupied raid doorway

Resolved the missing door contact above: the first officer used a lateral bay
when a pre-existing casualty occupied the central forecourt. Added an optional
shorter doorway-aligned corridor derived from the authored threshold, retaining
conservative traffic reservations. Browser8873 now shows three police cars/four
officers, one recorded door kick, an open door before entry and the original body
in place. Screenshot reports145FPS/321draws; this is local evidence only. Stop
restores the private preview and clears audio. Evidence: raid-door-corridor.json
and PNG. All190 frontend tests pass; build passes. Other building models still
need authored working entrances; hero character detail remains insufficient.

### September 13 — rounded tailored character geometry

Replaced rectangular jacket and sleeve/trouser meshes with shaped 24-sided
cross sections, shoulder caps, fitted waists and tapered cloth profiles. Regenerated
person/woman GLBs using local Blender and retained existing packed wool textures
and joint coordinates. Asset growth is approximately20KB each. All190 frontend
tests pass, including actual-model weapon-hand placement, walking surfaces and
full assassination cast/fall envelopes; frontend build passes. Browser8873 private
assassination review shows the updated silhouettes at145FPS/143draws. Screenshot
and diagnostic evidence: tailored-characters.png/json. This is an incremental
silhouette pass; toy-like anatomy, hands, faces and clothing variety still fall
short of final visual acceptance. No campaign command was issued.

### September 13 — working entrances across standard venues

Extended the authored hinged doorway/vestibule to ten standard building exports,
covering sixteen mapped addresses. Cleared front windows and piers away from the
opening. Added doorway-distance-aware staging for shallower buildings so the
breaching officer reaches the actual leaf. All200 frontend tests pass, including
closed/open passage rays, floor support and parcel bounds for every changed GLB;
build passes. Browser8873 Blue Hour and Paper Moon private raids each opened their
door and triggered one recorded kick. Blue Hour screenshot145FPS/277draws. Evidence:
bluehour-entry.png and venue-entries.json. No campaign commands issued. Specialist
buildings remain without working entrances; full interior and entry art variety
are still unfinished.

### September 13 — custody pose foundation

Replaced the detainee's identical rigid arm tilt with two-segment arm placement
that brings both wrists behind the waist over1.15s. Added a locally authored
steel handcuff GLB, revealed after the hands reach its position. Actual person and
woman rig tests verify wrist alignment and the full animation's reservation
bounds. All201 tests and build pass. Browser8873 arrest preview rendered the new
pose with two police cars and two officers; wide framing prevents detailed cuff
inspection, so no close-up visual acceptance is claimed. Evidence custody-pose
PNG/JSON. The arrest remains incomplete: an officer must approach/apply the
restraint and escort the detainee along a collision-checked shared path into a
vehicle with an opening door. Current work only supplies the restrained pose
and accessory for that sequence. No campaign command issued.

### September 13 — paired custody contact

Added shared officer/detainee choreography and a4.4m swept reservation. Officer
approach precedes wrist contact and visible handcuffs. Regenerated the police
asset with current elbow joints; actual-model tests verify both prisoner rigs,
contact accuracy, ground clearance and complete swept bounds. Arrest presentation
now lasts6s. Browser8873 private close-up verified the paired pose beside existing
aftermath; screenshot145FPS/167draws. All202 tests and build pass. Evidence:
paired-custody PNG/JSON. Escort/car-door/departure remain unfinished. No campaign
command issued. User then reprioritized map-first UI; continuing that migration.

### September 13 — map-first shell migration

Responded to the amended highest-priority objective by keeping one city component
mounted under menu overlays. Navigation and dashboard are now in-game HUD layers;
People/Families/Market/Ledger/Herald/Guide/Settings use a contained dialog with
focus handling, Close and Escape. City briefing preserves opportunities, public
threats, commissions and headline notices in an expandable HUD panel. Entering
Saint Agnes replaces the city canvas with its interior; Back to city restores
the city. Menus no longer cancel active scene playback merely by navigation.
Browser1280x720 verified People and Families overlays, unchanged camera coordinates
across open/close, a single background canvas, and interior replacement with the
status bar remaining onscreen. Fixed menu focus scrolling the shell and navigation
overflow found during QA. All202 existing frontend tests pass; layout correctness
is supported by browser evidence, not those simulation/geometry tests. Narrow
viewport and full action flows still need acceptance. Automatic post-scene
newspaper reveal, entrance sound and pre-generated narrator audio remain the next
priority; no completion claim for that requirement.

### September 13 — post-scene newspaper and narrator preparation

Added an automatic single-article Herald overlay after the city reports scene
completion; explicit Skip also reveals it. Exact headline/minute matching avoids
substituting unrelated weather or another event. A short cancellable paper-rustle
sound accompanies reveal. Enabled voices start a cancellable article-ID request
during animation and play the returned blob on reveal. Added server endpoint
restricted to saved published text, stable narrator profile, bounded cache and
post-synthesis stale-text check. Tests verify read-only state, unknown-ID rejection
and cache reuse. Browser8874 isolated strike fixture showed no article at1.4s and
an automatic paper after6.5s while retaining one city canvas. The speech service
on8787 is offline: browser fallback displayed Narration unavailable and Close
worked. Actual synthesized narrator quality/latency remains unverified; cache
behavior was tested with an isolated HTTP provider fixture, not real TTS.

Browser review exposed and fixed the newspaper using the generated face sheet for
hand-painted cast members, and a standfirst treating a person as a location.
Screenshot now shows Mara's correct portrait and corrected location phrasing.
All203 frontend tests pass, targeted newspaper core and speech HTTP tests pass,
and build passes. Evidence newspaper-reveal.png/json. No campaign actions issued;
QA save is .runtime/newspaper-reveal-20260913.sqlite3 on8874. Remaining work includes
real narrator playback acceptance, mobile UI, full event coverage and the larger
visual/choreography/interior requirements.

### Standalone newspaper and real narration — September 13

Removed the generic green menu frame and duplicate title from the scene edition.
The newsprint itself is now the modal, with a stamped Fold away control above the
masthead and matching narration controls inside its ruled footer. Keyboard focus
is contained, Escape closes it, and closing restores prior focus. Paper scrolls
on smaller viewports; reduced-motion preference is retained. Fixed paused audio
resuming without updating its status; loading/playing disables redundant Read
aloud requests.

Production build passes. Browser QA on isolated8874 at1280×720 confirms zero
map-menu frames, all buttons inside paper, entire sheet within viewport, and one
city canvas retained through reveal and close. Pause reports Narration paused;
Resume reading reports Narrating. Screenshot: standalone-newspaper.png. No save
commands were issued. The temporary QA tab was closed.

Real local narration on8787 now plays: cold synthesis prepared in8418ms, beginning
1915ms after reveal; cached preparation19ms and reveal-to-play7ms in the timing
fixture (narrator-timing.json). Standalone UI replay measured106ms preparation and
22ms reveal-to-play. These are playback API timings, not an auditory quality
review. Cold-start latency still needs improvement; production-wide acceptance
and the broader city/interior requirements remain open.

### Reuse the local narrator model — September 13

Previous goal turn classified as progress: standalone newspaper committed and
browser-verified. Continued the full goal by reducing preparation latency for new
articles. Kokoro now runs in one persistent subprocess under the existing serial
render lock, retaining its model across requests while resolving each character's
own voice recipe. Requests have matching IDs, bounded responses and an18-second
deadline. Timeout, process exit or invalid reply terminates/reaps that worker;
next request starts a fresh one. Cache bounds, offline assets, text validation
and WAV validation remain intact. No public interface change.

Four process-level tests pass: reuse across distinct articles, timeout recovery,
crash recovery and mismatched-response rejection. Python compilation passes.
Actual uncached synthesis of two different articles took3978ms then346ms in the
same PID, both valid RIFF outputs (narrator-worker-timing.json). An isolated local
voice service on8788 with an empty cache and game preview8875 prepared the saved
article in3202ms; browser playback began26ms after reveal, while preserving the
city canvas (narrator-worker-browser.json). QA only replayed a committed fixture;
no campaign commands/save mutations. Temporary browser closed and voice preference
restored. Main local voice service8787 restarted with this implementation.

This improves new-article reuse; it does not prove latency on every machine or
long article, nor auditory quality. Broader cinematic, interior and visual-quality
acceptance remains unfinished.

### In-game document menus — September 13 amendment

Read the amended objective409771b4-1304-4bc4-ae42-631508727ff1 and updated
GOAL.md. Previous turn was progress (persistent voice worker and real playback
verification). Reprioritized menu presentation before further custody choreography.

The six navigation menus now use period document surfaces: a contact book, family
dossiers, a trade circular, private ledger, handbook and control desk. Added cloth
spines, index tabs, textured stock, ink rules and an integrated Put away control.
Each edition keeps its original content/actions and focus/close handling. Local
palette overrides cover cards, search/filter controls, settings switches, account
figures, risk text and unavailable actions. Browser review caught residual dark
panels and low-contrast disabled labels; these were corrected before acceptance.
The archive also receives a paper reading-room treatment; its full content was not
included in this six-menu QA pass.

Production build passes. CUA reviewed all six menus at1280×720 on isolated8875:
one city canvas retained, menu within viewport, visible close button and no
horizontal body overflow for each. Recorded document-menus.json, menu-contacts.png
and menu-ledger.png. Closed temporary QA tab. No gameplay/save commands issued.
Responsive rules were added but compact viewport acceptance remains unverified.
Dialogue and other remaining game surfaces, fuller menu identity/illustration,
interiors and cinematic requirements remain unfinished under the full goal.

### Conversation folios and memorial sheets — September 13

Read the current409771b4 goal; previous turn was progress (document-menu commit
and browser evidence). Continued the menu amendment into mandatory encounters.
Conversations now present an open paper folio with speaker portrait, message,
voice control and shared conditions on one leaf, with separate numbered reply
slips beside it. Paper choices retain all costs, durations, attention and refusal
reasons. Compact layout stacks the leaves with one outer scroll owner. Removed
the inherited nested choice scroller after browser inspection clipped reply four.
Death records use a black-edged memorial sheet with integrated new-life action.
No command semantics, voice handling, modal focus behavior or API were changed.

Build passes. At1280×720 CUA verified all four offer choices fully visible,
no horizontal overflow and one city canvas behind the folio. Declined the offer
on new isolated .runtime/encounter-style-20260913.sqlite3 (8876): revision1,
minute480, cash900, respect25, health100 and no pending event; UI returned to the
same city canvas. This was the only QA command, and never touched the campaign.
New isolated memorial fixture on8877 shows continuation button fully within the
viewport and one city canvas; no new-life command issued. Screenshots saved as
conversation-folio.png and memorial-sheet.png. Both temporary tabs closed.
Compact viewport and varied long encounters still need acceptance; broader city
assets, interiors, choreography and performance requirements remain unfinished.
