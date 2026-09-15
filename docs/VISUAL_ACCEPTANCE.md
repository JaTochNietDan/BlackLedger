# Visual delivery acceptance

Record evidence and unresolved items; never mark an item passed solely because the build succeeds.

- [ ] Compare assembled block with approved painted-noir references; consistent scale, ground plane, perspective and material treatment.
- [ ] Day/night screenshots at 1440×900 and 980×800, plus 390×844 responsive smoke check. State actual measured viewport sizes.
- [ ] No checkerboards, opaque halos, road seams, floating glows, clipping, illegible signs or obvious character/building depth errors in inspected views.
- [ ] Observe moving people/cars, arrivals and a committed attack sequence; inspect overlap at entrances/corners, not only still frames.
- [ ] Every painted building can be selected through visible art and accessible controls; transparent padding does not steal clicks. Selection does not travel.
- [ ] Unpainted destinations remain clearly accessible through the directory; no false art/identity mappings.
- [ ] Intact/damaged laundry preview and street agree with saved condition. Repair updates both without shifting footprint.
- [ ] Header, landmark controls, known-threat reminder and latest developments remain legible without collisions or horizontal overflow.
- [ ] Test long dialogue/choices, police stop, death recap and new-life portraits. Keyboard focus is visible; modal traps focus without automatically accepting a choice.
- [ ] Skipping/replaying travel or outcomes leaves game state unchanged. Pausing ambience and reduced-motion preferences are respected.
- [ ] Missing art/iframe failure still offers usable directory navigation and retry; no lost campaign state.
- [ ] Browser console inspected; `npm run build` and `npm test` pass. Run `go test -race ./...` if any shared behavioral code changed.
- [ ] Deliver screenshots, commit list, asset provenance, remaining limitations and API proposals. Never claim Steam readiness or full campaign coverage from this visual checklist.


## 2026-09-14 authored interior coverage audit

The 3D revision in GOAL.md supersedes the older painted-only delivery language
above. At f52edbc,22 of33 public addresses map through interiorSettings to an
existing GLB. The24 interior GLBs include two private-room variants; those do
not add public addresses. This is routing/file evidence, not per-room visual
acceptance. Interior.tsx still renders the other11 addresses through its painted
or SVG room branch with HTML occupants. Poker/blackjack/billiards table scenes
are separate and do not establish completion of each venue's general interior.

| Missing authored public room | Address ID |
| --- | --- |
| The Monarch | `club` |
| Kerrigan Haulage | `haulage` |
| The Golden Lily | `goldenlily` |
| Ordway Steam Laundry | `steamworks` |
| The Paper Moon | `burlesque` |
| Ferris Motor Sales | `dealer` |
| Archway Motor Repairs | `archway` |
| Devlin Salvage | `scrapyard` |
| Kessler's Filling Station | `filling` |
| The Viaduct Pumps | `pumps` |

Next room-production groups: casino/club floors (Monarch, Golden Lily,
Paper Moon), industrial premises (haulage and steam laundry), then vehicle
businesses (dealer, repair garage, salvage and filling stations). Reusing props
is appropriate; a generic borrowed room does not prove an address's fixtures,
occupancy, movement clearance, interaction points or incident choreography.

For every added room, verify authored props and materials, camera framing at
normal/compact windows, wall cutaways, actual public occupants, seated/standing
clearance, selectable actions, arrival paths and relevant indoor incidents.
Existing22 rooms still require broader fidelity/incident review. The current
poker desktop review and automated GLB tests are narrower evidence.

### Blue Hour room increment

The casino now adds an authored public room: current coverage is23/33 public
addresses, with10 remaining in the table above. Its original Blender source
exports midnight-blue walls, gaming tables, slots, cashier cage and lounge
seats. Browser review on isolated fixture8998 at1235×1053 shows the full room,
six actual occupants, clear header and hover-menu controls. Automated tests
verify both character rigs, all12 placement slots, seat support, collision
clearance and40 entrance-path samples against exported geometry. Build passes.
Compact windows, room cutaway interaction and casino incident choreography
remain to be reviewed; this is not full visual acceptance.
