# Surfaces

Four images clothe the whole city. They are multiplied by each building's own
palette colour, so one brick serves a street of differently coloured buildings —
you need one brick, not a red one and a brown one and an olive one.

| file | goes on |
|------|---------|
| `brick.png`  | the storeys above the shopfront, and chimney stacks |
| `stone.png`  | cornices, sills, lintels, parapets, plinths |
| `roof.png`   | the flat roof |
| `render.png` | the shopfront wall below the fascia |

PNG, JPG or WebP. **A missing file is not an error** — that surface falls back
to its flat palette colour, so a half-finished set still renders a whole city.

## The rules a texture has to follow

1. **Flat lighting. No shadows, no highlights, no lighting of any kind.**
   Blender does the light, and the light has to be the same on every building
   in the city. A texture with its own shadows baked into it fights the sun and
   the building goes muddy. This is the one that gets got wrong.
2. **Seamless and tileable.** Every surface is box-projected, so a texture that
   does not tile shows its join on every wall.
3. **Shot square on**, no perspective. It is a material, not a photograph of a
   wall taken from an angle.
4. **Fairly desaturated.** The colour comes from the building's palette; the
   texture supplies the grain, the courses and the dirt. A vivid texture
   fights the palette instead of taking it.
5. **512 or 1024 square.** Bigger buys nothing at the size a building is drawn.

A prompt that satisfies all five:

> old weathered red brick wall, worn mortar courses, soot stained, seamless
> tileable texture, flat even lighting, no shadows, no highlights, no
> perspective, photographed straight on, uniform, fills the whole frame

## Testing one

Drop it in as `brick.png` and run:

    mise run city-art

Ten seconds for the whole set. `tools/blendercity.py` holds `TILE_SIZE`, which
is how much wall one tile covers — the reason brick is the same size on the
tenement as on the tavern.

## What is here now

Placeholders, generated locally with FLUX-schnell at four steps to prove the
pipeline end to end. They are not good: `render.png` came back as stone blocks
rather than stucco, and `roof.png` reads as slate rather than tar and felt.
They are meant to be replaced.
