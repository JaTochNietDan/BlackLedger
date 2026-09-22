# Media rights review before public distribution

This is a focused list of material needing owner/source evidence. Hashes identify
the exact files reviewed; they do not establish a license. No files were removed
or replaced, and no permission was inferred from their presence in the repository.

The 14 WAV files are documented as user-supplied in `public/audio/README.md`.
The record authorizes use in Black Ledger but does not state repository/archive
redistribution terms. Confirm original ownership or provide the source license.

Commit `079f0a7` introduced the four ground images as photographic textures.
Its development record describes cropping/tiling, without naming a provider,
photographer, source URL or license. Their filenames alone cannot establish rights.

Commit `11b4190` calls the imported shop image a Grok building. Confirm the
original generation/source and applicable output rights; the cutout and live
isometric sprite derive from that original. The early style images are recorded
as project references in `docs/ART_DIRECTION.md`, without file-level provider
or source-license records. Keep them within the review even though they are not
needed by the current runtime.

For each group, record source/creator, source URL or receipt if available,
license/version, required credit, and whether modified copies may be distributed
in both GitHub source and binary game archives. Private receipts need not be
published; record the permission evidence and public attribution instead.

## Supplied sound recordings

| File | SHA-256 |
| --- | --- |
| `public/audio/effects/door-kick.wav` | `cc09152afb223736fdc0c1f718eea7c772b47f7b240acb13b217120f1abcf98b` |
| `public/audio/effects/drive-away.wav` | `e3705691ecb395cf1b9713b499b0da47caf02a16e4d98b9fadd7b174db75d602` |
| `public/audio/effects/engine-idle.wav` | `b273c9a853ad2691d3707bd8230d26bd872485d2cd10328725bcb93d654db3f1` |
| `public/audio/effects/explosion.wav` | `b511683672956ec58d5fbef3c4e3e9100b7de8e1102bea39c5030e85d8fb2f43` |
| `public/audio/effects/fire.wav` | `6d0a5117f849eb24c5f3bc4a0f4e573c96f7ab30c479037c8cfe340d9f9420b8` |
| `public/audio/effects/newspaper.wav` | `901a30471fcad6a4d337e1f80912a1a0143fb80a8764ce93a4a07513b43100a9` |
| `public/audio/effects/pain.wav` | `ed26ef1e799ecc33a0afde7d138328593efaca30c3fbf941a985ddd893de3eae` |
| `public/audio/effects/panic.wav` | `084c56bfeb807e39c9d91472dcdd350f39debb14f61bcb0b7daeb7eeef0023f8` |
| `public/audio/effects/raid.wav` | `e8dac86267d853268f06f19e97564b01862bd402c0748815c7311b2d6ae9d7a5` |
| `public/audio/effects/siren.wav` | `182f3f91aeba1486b54f89a77e0778f230b5071fffa9a71316fa75efca69e675` |
| `public/audio/effects/vehicle-approach.wav` | `8f7716e50708f85868b0eecec55d46bd6ba9a08ce35cc47b5865b41b05ef8681` |
| `public/audio/guns/revolver.wav` | `3f9beed919a410a25a63638a6b208427696dfb70df9f2fe8407a6a34be2954e0` |
| `public/audio/guns/shotgun.wav` | `718503005981fb15ee1daa5829f91935cbd054b45011b255bc2eb2a9d99c4319` |
| `public/audio/guns/thompson.wav` | `77f06115d264cccab707bc9496e4c19c5f9267b7c87a74d1b9a1fd5951764a2f` |

## Ground textures with no recorded source license

| File | SHA-256 |
| --- | --- |
| `public/art/ground/asphalt.jpg` | `ec5b3e4b640eef6ae33ec5a66a343bc78bbd29ef711b5b31c6133b7b50034ca2` |
| `public/art/ground/cobbles.jpg` | `8cab1d013fb38ad7ea203b5220e0e9dd1db5322526c13bf481b13819333dc67a` |
| `public/art/ground/kerb.jpg` | `b9c1d04cae2e4e447de44430be259c2ffb2812f45001a543fa609627e18000bd` |
| `public/art/ground/pavement.jpg` | `54a4ec08acb9c0de7f83310e863fd55480902ec4d1c06b41b554d250d7c519c7` |

## Imported originals and derivative

| File | SHA-256 |
| --- | --- |
| `art/sources/shop-corner-on-magenta.jpg` | `761831e87058dfe43ca7b27c4a297449ed977d84948db12720780e1fb0adb74a` |
| `art/sources/shop-corner-cutout.png` | `c46e8f06df367dadfc0e96cadf15a08196dfe060a2aba4fbdf48c8d36bf68b69` |
| `public/art/iso/iso-bar-v2.png` | `da48bc7f60813db2ad821a2cfeba5f8a49cf63df4269824cdd420886929c44df` |

## Early style references

| File | SHA-256 |
| --- | --- |
| `docs/art/city-direction-v1.png` | `9e7830c836d6b8b69d404cc55f14785fc27cd1bba56e1cf33d8def9744613e26` |
| `docs/art/style-comparison-v1.png` | `793664aaf4f1a43166392054ce713af4a06e97688e82b53d9880e8319f28297c` |

## Other media

Locally authored models have editable Blender source under `tools/` and provenance
notes in `public/art/models/README.md`. Local FLUX generation tools document
portraits, cars, scenes, frontages and interiors; the six death-service frontage/
interior plates are rendered from local GLBs by `tools/render_death_service_plates.py`.
These records provide a stronger origin trail, but should not be mistaken for a
guarantee of exclusive rights in generated imagery. Third-party software licenses
are handled separately in `THIRD_PARTY_NOTICES.md`.

This list covers selected current-tree files, not every historical image version.
Before publishing full Git history, apply any unresolved source restrictions to
earlier versions too. See `ASSETS.md` for the overall publication status.
