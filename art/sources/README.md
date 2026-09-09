# Where the art came from

The originals, as they arrived. `public/art/iso` holds only what the game
serves: a source kept beside it is a megabyte or two shipped to every player
for nothing.

Bring one in with:

    .venv/bin/python tools/importart.py art/sources/<file> <address>

which lifts off a flat background, trims, measures the footprint, and refuses
anything that will not sit on the grid — with a reason that says what to ask
for instead.
