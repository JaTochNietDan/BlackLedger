
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
