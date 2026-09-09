# The ground

Served to the browser, unlike `art/textures`, which holds surfaces for the
offline Blender renderer. These four are laid down live by `src/CityIso.tsx`
under the roads and pavements the grid computes — the geometry is unchanged,
the surface is real material.

| file | goes on |
|------|---------|
| `asphalt.jpg`  | the carriageways |
| `pavement.jpg` | the footway inside each block |
| `kerb.jpg`     | the edge between them |
| `cobbles.jpg`  | the avenue the streetcar runs down |

Check one tiles before trusting it:

    .venv/bin/python tools/tileable.py public/art/ground/pavement.jpg --fix

A texture that does not tile shows its seam on every street in the city, and
that seam is invisible in the source image and obvious the moment it is laid
down. The rules a texture has to follow are in `art/textures/README.md`; the
one that matters most is flat lighting, because the city lights itself and the
light changes with the hour.
