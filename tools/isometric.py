"""Paint every address as an isometric cut-out for the city view.

The city view draws the whole of Bellwether from above at a fixed angle. That
needs a different picture from the address book: not a three-quarter view taken
from across the street, which is what tools/exteriors.py paints, but a complete
model of the building seen from one corner, standing on nothing, with the roof
visible — the shape a cut-out has to be to sit on an isometric plot.

Five of these were made by hand for the street study and the other seven never
existed, which is why the city view was blocked out in flat polygons. This
paints all twelve in one pass so they belong to each other: the same light, the
same palette, the same angle. A city assembled from cut-outs drawn at different
angles reads as a collage however good each piece is.

Same contract as the other art tools: generated offline, shipped as files,
nothing at runtime depends on a model.

    mise run art        # once
    mise run isometric
"""
import json
import os
import sys
import tempfile

from PIL import Image

from paint import generate

# The angle is the whole point and it is stated three ways, because a diffusion
# model will drift off a single phrase.
LOOK = (
    "isometric view, 45 degree overhead three-quarter angle, orthographic game asset, "
    "1950s American city building, complete model of the whole building including its roof, "
    "film noir oil painting, dark umber and olive palette, muted and desaturated, "
    "single warm light in the windows, overcast evening, visible brushwork, "
    "isolated object centred on a plain flat black background, "
    "no ground, no street, no scenery, no people, no cars, no text, no lettering"
)

# Each address, described as a building rather than as a place, because the
# model is painting a thing and not a story.
PLACES = {
    "room": "a narrow four-storey lodging house, fire escape zigzagging down the front, shallow pitched roof",
    "bar": "a corner tavern, two storeys, awning over the door, chimney stack, warm lit windows",
    "docks": "a long low waterfront cargo shed, corrugated roof, loading doors, a crane arm at one end",
    "laundry": "a squat single-storey commercial laundry, flat roof, roof vents, wide shopfront window",
    "club": "an ornate two-storey nightclub, stepped parapet, marquee canopy, rooftop sign frame",
    "market": "a grand exchange hall, colonnade along the front, glazed pitched roof, stone steps",
    "apartment": "a respectable four-storey apartment block, bay windows, stone stoop, mansard roof",
    "garage": "a motor repair garage, wide roller door, flat roof, brick front, small office window",
    "casino": "a discreet two-storey gambling club, shuttered upper windows, awning, flat roof",
    "estate": "a large detached house, steep gables, tall chimneys, dormer windows, walled garden edge",
    "precinct": "a grey stone police station, heavy steps, lamps either side of the door, barred windows",
    "herald": "a three-storey newspaper building, tall printing hall windows, loading bay, roof water tank",
}

# The blocks between the addresses that matter. A city with holes in it reads
# as a scatter of models, and a hole filled with a flat grey box reads as a
# placeholder — which is worse, because it looks like a mistake rather than
# like distance. These are ordinary buildings, deliberately plainer than the
# twelve, painted in the same light so they belong to the same city.
FILLERS = {
    "fill-tenement": "a plain five-storey brick tenement, flat roof, fire escape, rows of small windows",
    "fill-warehouse": "a low brick warehouse, shallow pitched roof, loading doors, few windows",
    "fill-terrace": "a row of three joined two-storey shopfronts, awnings, flats above",
    "fill-office": "a narrow four-storey office block, stone facing, regular windows, parapet",
    "fill-works": "a small industrial works, sawtooth roof, brick walls, a short chimney",
    "fill-corner": "a corner building of three storeys, chamfered corner entrance, shop below flats",
}
PLACES.update(FILLERS)

# Terrace rows, drawn to stand shoulder to shoulder. Every building above is a
# detached model carrying its own pavement, which is why the party walls in the
# city meet unevenly however carefully they are placed. These are rows with
# flush ends: the join is painted into the picture rather than faked by pushing
# two models together.
ROWS = {
    "row-brick": "a terrace of four joined 1950s brick buildings in a straight row, flush square ends, shared party walls, shopfronts at street level and flats above",
    "row-stone": "a terrace of four joined stone-faced buildings in a straight row, flush square ends, shared party walls, uniform cornice, regular windows",
    "row-mixed": "a terrace of four joined buildings of different heights in a straight row, flush square ends, shared party walls, brick and render, awnings over the shops",
    "row-low": "a terrace of four joined two-storey shops in a straight row, flush square ends, shared party walls, flat roofs, plain frontages",
}
PLACES.update(ROWS)

# The signage the model invents is nonsense — a marquee reading CIICATCE, a
# newspaper office reading NEWSAPWL — because diffusion cannot spell. Asking
# for "no text" does not stop it, so these are asked for as buildings that
# genuinely have no signs on them.
BLANK = ", plain unlettered facade, blank walls, no signs, no billboards, no marquee, no writing anywhere"

# Square, because a cut-out is placed by its footprint and a square keeps the
# building's own proportions out of the frame's business.
SIZE = 1024

# A fixed order, so every address keeps its own seed for good.
SEEDS = list(PLACES)

# Nudged when a building comes out wrong — the model sometimes paints a floor
# across the frame instead of an isolated object, which tools/inspect_iso.py
# catches. Set ISO_RETRY to move that one building's seed without disturbing
# any of the others.
RETRY = int(os.environ.get("ISO_RETRY", "0"))


def knock_out(path: str) -> None:
    """Take the flat background off, from the corners inwards.

    The model is asked for a plain black field and gives something close to one.
    Flood-filling from the four corners removes it without eating the dark
    parts of the building, which a simple colour threshold would.
    """
    im = Image.open(path).convert("RGBA")
    w, h = im.size
    px = im.load()

    # How near to a corner's colour counts as background. Generous enough for
    # the model's noise, tight enough to stop at a lit window or a wall.
    def near(a, b, tol=42):
        return abs(a[0] - b[0]) <= tol and abs(a[1] - b[1]) <= tol and abs(a[2] - b[2]) <= tol

    seeds = [(0, 0), (w - 1, 0), (0, h - 1), (w - 1, h - 1)]
    field = px[0, 0]
    seen = bytearray(w * h)
    stack = list(seeds)
    while stack:
        x, y = stack.pop()
        if x < 0 or y < 0 or x >= w or y >= h:
            continue
        i = y * w + x
        if seen[i]:
            continue
        if not near(px[x, y], field):
            continue
        seen[i] = 1
        px[x, y] = (0, 0, 0, 0)
        stack.extend(((x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)))

    # Trim to what is left, so the sprite's box is the building's box and the
    # renderer can anchor it to a plot without guessing.
    box = im.getbbox()
    if box:
        im = im.crop(box)
    im.save(path)


def paint(out_dir: str, only: list[str]) -> None:
    os.makedirs(out_dir, exist_ok=True)
    made = []
    for name, described in PLACES.items():
        if only and name not in only:
            continue
        path = os.path.join(out_dir, f"iso-{name}-v1.png")
        prompt = f"{described}, {LOOK}"
        if name in ROWS or os.environ.get("ISO_BLANK"):
            prompt += BLANK
        with tempfile.TemporaryDirectory() as tmp:
            raw = os.path.join(tmp, "out.png")
            print(f"painting {name}…", flush=True)
            # A fixed seed per address, so regenerating one building does not
            # reshuffle the rest of the city.
            generate(prompt, 4100 + SEEDS.index(name) * 13 + RETRY * 977, raw, width=SIZE, height=SIZE)
            Image.open(raw).save(path)
        knock_out(path)
        made.append(name)
        print(f"  wrote {path}", flush=True)

    # A manifest the interface reads, so adding an address tomorrow is a file
    # and a line rather than a code change.
    manifest = os.path.join(out_dir, "isometric.json")
    have = {}
    if os.path.exists(manifest):
        have = {e["id"]: e for e in json.load(open(manifest))}
    for name in made:
        with Image.open(os.path.join(out_dir, f"iso-{name}-v1.png")) as im:
            have[name] = {"id": name, "file": f"iso/iso-{name}-v1.png", "w": im.width, "h": im.height}
    json.dump(sorted(have.values(), key=lambda e: e["id"]), open(manifest, "w"), indent=1)
    print(f"manifest: {len(have)} addresses")


if __name__ == "__main__":
    target = sys.argv[1] if len(sys.argv) > 1 else "public/art/iso"
    paint(target, sys.argv[2:])
