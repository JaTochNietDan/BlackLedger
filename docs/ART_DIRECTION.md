# Approved visual direction

User selected painted noir realism as the core, warm vintage daylight, and richer amber/crimson nightlife. Reference: `art/style-comparison-v1.png`, panel 1 for materials and overall realism, panel 4 for daylight, panel 3 for night atmosphere. Panel 2's graphic-novel ink treatment is not the selected style. The first city painting is `art/city-direction-v1.png`.

## Asset approach

Use modular 2D architecture with a consistent orthographic isometric view, realistic proportions and restrained painted surface detail. Author streets, buildings, vehicles and people separately. Layer signage, window light, smoke, weather and damage over static assets. The backend owns gameplay outcomes; animation only depicts committed results or cosmetic activity.

Start with a small coherent block before scaling production. Every building needs a footprint, ground anchor, entrance point, hit polygon, depth-sort anchor and light anchors. Establish these through an assembled block test; generated images alone do not guarantee matching perspective, dimensions or seamless street joins.

Avoid generating a full animation as independent images. Use moving sprites and procedural effects for cars, pedestrians, bulbs and smoke. Important story actors receive presentation tracks from committed events. Decorative traffic has no authority over gameplay.

## First building study

`public/art/casino-noir-v1.png`: generated transparent casino, 1536 × 1024 RGBA. Corner samples verified transparent. `public/art-study.html` overlays independently animated marquee bulbs and demonstrates a simple day/night grade. Preview at `/art-study.html` after `npm run build`.

This is an isolated asset/presentation study. Global image grading is not a complete lighting system; proper window masks, foreground occlusion, street tiles and standardized production dimensions remain to be developed. Do not describe this as a complete modular city.

Keep source art unchanged. Store footprint and lighting coordinates as metadata when integrating production assets. Generated art is never stored only in a user-specific tool cache once referenced by the project.

## Street assembly study

`/street-study.html` uses café, casino and sedan PNGs with scripted paths and depth sorting. It exercises separate building, ground, actor and light layers. Road art and pedestrian figures are provisional. Building positions, entrance anchors and presentation walkways now live in `public/art/buildings.json`, shared by the live street and React controls. Marquee bulb anchors remain in the renderer. The Mariner and laundry are integrated; Mercer Exchange is also integrated with an explicit silhouette mask; docks and additional district art remain pending.

## Painted character atlas

`public/art/cast-noir-v1.png` is a 3 × 2 atlas of square portraits: Mara, Leo, Vittorio / Elena, Harlow, Alex. React selects stable atlas cells by identity; it does not generate faces during play. Future player identities still use the procedural fallback until additional cast art is authored. The raw atlas is preserved without destructive crops.

## Saved building condition

The building manifest supports `damage: {file, below}`. Bluebird Laundry currently switches to its damaged storefront below 70 condition. Both the street and selected-property preview consume the saved public condition. The variant preserves the original footprint by using the intact sprite's alpha as a runtime rendering mask (Canvas destination-in for the street, CSS mask for the preview). Source images remain unmodified. The first generated damage asset had an opaque halo; a second attempted cutout rendered a checkerboard, so neither generated background is trusted as transparency. The original verified alpha defines the displayed silhouette.
