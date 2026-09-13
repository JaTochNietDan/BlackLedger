# Bellwether realtime models

All 18 GLB files are authored locally in Blender by `tools/export_city3d.py`.
No purchased asset packs or assets from Afterlight are used. The generator is
the editable source; run `.venv-blender/bin/python tools/export_city3d.py` to
rebuild. Blender 5.2.1 LTS (`bpy`) was used for this revision.

13 building archetypes, Ford/Hudson/Packard silhouettes, a police sedan and an
articulated pedestrian. Geometry uses metres, Z-up in Blender and Y-up in glTF.
`manifest.json` records actual exported bounds including cornices, fire escapes,
awnings and cargo. Models must stay within the 17m reserved footprint. The
browser places these on 28m blocks with separate pavements and carriageways.
`sign-anchor` empties place readable address signs onto the authored facades.

Brick base-color and tangent normal maps are deterministic local textures
created by the generator and embedded in the GLBs. Physical UV scale is shared
between buildings. Chrome, glass, canvas, concrete and roof materials use glTF
metallic/roughness materials. Streets reuse the existing repository's ground
textures. Materials are grouped before export to reduce draw calls; pedestrian
limbs and their attached shoes remain separate for browser animation.

This is the first realtime asset set, not final visual acceptance. More facade
variety, convincing worn surfaces, better pedestrian anatomy/wardrobe, denser
street dressing and richer district landmarks remain production work.
