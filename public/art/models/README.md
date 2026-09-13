# Bellwether realtime models

All 23 GLB files are authored locally in Blender by `tools/export_city3d.py`.
No purchased asset packs or assets from Afterlight are used. The generator is
the editable source; run `.venv-blender/bin/python tools/export_city3d.py` to
rebuild. Blender 5.2.1 LTS (`bpy`) was used for this revision.

13 building archetypes, Ford/Hudson/Packard silhouettes, a police sedan and an
articulated pedestrian, plus a streetside furniture set. Geometry uses metres,
Z-up in Blender and Y-up in glTF.
`manifest.json` records actual exported bounds including cornices, fire escapes,
awnings and cargo. Models must stay within the 17m reserved footprint. The
browser places these on 32m blocks with separate pavements and carriageways.
`sign-anchor` empties place readable address signs onto the authored facades.

Brick base-color and tangent normal maps are deterministic local textures
created by the generator and embedded in the GLBs. Physical UV scale is shared
between buildings. Chrome, glass, canvas, concrete and roof materials use glTF
metallic/roughness materials. Streets reuse the existing repository's ground
textures. Materials are grouped before export to reduce draw calls; pedestrian
limbs and their attached shoes remain separate for browser animation.

The estate uses its own pitched tile roof, gables, porch and shuttered windows.
Casino and civic silhouettes include Art Deco crowns and stepped clock towers.
Clock hands remain separate nodes so the browser can show saved game time.
Rear fire escapes leave the main entrances and signs unobstructed.
The bench, bin and hydrant set is instanced by material along a separate rear
pavement band. Tests check its actual bounds against buildings and all routes.
Four venue variants give The Monarch, Blue Hour, Golden Lily and Paper Moon
separate proportions, brick palettes and neon accents. Their marquee anchors
fit the address signs to the canopy fascia. They share the same physical scale
and reserved parcel envelope as the other buildings.

This is the first realtime asset set, not final visual acceptance. More facade
variety, convincing worn surfaces, better pedestrian anatomy/wardrobe, denser
street dressing and richer district landmarks remain production work.

The pedestrian has a tailored jacket, shirt/cuffs, hands and shaped fedora, with
separate hip, knee and arm joints. The manifest includes bounds sampled across
48 gait phases; traffic occupancy and pavement-clearance tests cover that stride.
