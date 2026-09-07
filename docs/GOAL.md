# Current goal and acceptance requirements

Build and playtest the full Black Ledger single-player mafia vertical slice described in DESIGN.md: a playable 20–30-minute rise from rented housing through contacts, crew and business ownership, into consequential rival incidents; permanent death and new people in a persistent city; validated AI opportunities; reliable transactional saves and optional stable-character speech.

## Approved visual expansion — 2026-09-07

The user approved a modern 2D isometric city inspired by the visual approach of Gangsters: Organized Crime. Painted noir realism is the core; warm vintage daylight and amber/crimson nightlife are lighting variations. This replaces the schematic cartoon map as the intended visual destination.

Required implementation:
- A coherent small neighborhood assembled from reusable illustrated buildings and street pieces, with consistent scale, perspective, anchors and occlusion.
- Moving cars and pedestrians on authored presentation routes; decorative actors do not require simulation agents.
- Independent marquee lights, window glow and limited atmospheric effects.
- Gameplay-selected destinations remain interactive, and visual travel can be skipped without changing committed outcomes.
- Narrative visual sequences consume backend results. They do not calculate combat, rewards, political outcomes or game time.
- Expandable asset/metadata structure allowing more districts and changed building states.
- Browser visual QA at practical desktop and compact sizes, and motion controls/reduced-motion behavior.

The isolated casino study is an asset test, not completion of the neighborhood requirement. Do not shrink the gameplay goal to an art demo. Keep implementation and validation records in DEVELOPMENT.md and commit coherent changes in this independent repository.
