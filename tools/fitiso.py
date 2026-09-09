"""Measure where each cut-out's building actually stands in its own frame.

The city was drawing every sprite as though three things were true of it: that
the building's base is centred in the image, that the bottom edge of the image
is the near corner of that base, and that the image is exactly as wide as the
base. None of the three is true of a picture somebody generated. The model puts
the building wherever it likes, leaves as much air around it as it likes, and
the widest thing in the frame is usually a cornice or an awning rather than the
footprint. So good art came out misaligned — sunk into the pavement, shoved off
its plot, or scaled to its own overhang instead of its own ground.

Nothing about that is the artwork's fault and no amount of regenerating fixes
it. The building's footprint is measurable, so it is measured here and written
into the manifest, and the view places each sprite by its own feet.

    .venv/bin/python tools/fitiso.py public/art/iso

Writes `anchor` and `base` into isometric.json, and with --proof also writes a
copy of each cut-out with the footprint it found drawn on top, so a fit that is
wrong can be seen rather than argued about.
"""
import json
import os
import sys

from PIL import Image, ImageDraw

# A 2:1 isometric base: for every pixel sideways from the near corner, the base
# edge climbs half a pixel. src/iso.ts draws the city in the same projection —
# TILE is 108 by 54 — so this is not a guess about the art, it is the geometry
# the city is already drawn in.
SLOPE = 0.5

ALPHA = 24          # below this a pixel is background rather than building
DRIFT = 2.4         # how far the profile may stray from the ideal edge, in px


def profile(im: Image.Image) -> list[int | None]:
    """The lowest opaque pixel in each column: the building's underside."""
    alpha = im.getchannel("A").load()
    w, h = im.size
    out: list[int | None] = []
    for x in range(w):
        low = None
        for y in range(h - 1, -1, -1):
            if alpha[x, y] >= ALPHA:
                low = y
                break
        out.append(low)
    return out


def content(im: Image.Image) -> tuple[int, int, int, int]:
    """The box the building actually occupies, ignoring the air around it."""
    box = im.getchannel("A").point(lambda v: 255 if v >= ALPHA else 0).getbbox()
    return box or (0, 0, im.width, im.height)


def fit(im: Image.Image) -> dict:
    """Find the near corner of the building's base, and how wide the base is.

    Two readings, and the strict one only wins when it looks sane.

    The strict reading follows the underside of the silhouette out from its
    lowest point: in this projection a base edge climbs half a pixel for every
    pixel sideways, so the footprint ends where the underside stops following
    that line — which is where a fire escape or a flight of steps begins, and
    those should not count as ground.

    That reading fails on the painted cut-outs, because they were generated with
    a soft cast shadow under them: the lowest pixels are a blurred blob rather
    than two straight edges, and the walk stops after a few pixels. So when it
    comes back with a footprint far too small to be a building, the honest
    fallback is the content box — trimmed of the air the model left around the
    picture, which is most of the misalignment anyway.
    """
    left_px, top_px, right_px, bottom_px = content(im)
    width = right_px - left_px
    low = profile(im)
    solid = [x for x, y in enumerate(low) if y is not None]
    if not solid or width < 2:
        return {}

    deepest = max(low[x] for x in solid)
    flat = [x for x in solid if low[x] == deepest]
    corner_x, corner_y = (flat[0] + flat[-1]) // 2, deepest

    def edge(step: int) -> int:
        x, last = corner_x, corner_x
        while 0 <= x + step < len(low):
            x += step
            y = low[x]
            if y is None:
                break
            want = corner_y - abs(x - corner_x) * SLOPE
            if abs(y - want) > DRIFT:
                break
            last = x
        return last

    strict = edge(1) - edge(-1) + 1
    exact = strict >= width * 0.45      # a real footprint, not a shadow blob
    if not exact:
        # The near corner of the content box, which is what the view anchors to.
        corner_x, corner_y = (left_px + right_px) // 2, bottom_px - 1

    return {
        # Where the building's feet are, as a fraction of the image: what a
        # sprite anchor takes directly.
        "anchor": [round(corner_x / im.width, 5), round((corner_y + 1) / im.height, 5)],
        # How wide its ground is, in pixels of this image. The view scales the
        # sprite so this matches the plot it is standing on, which is why a
        # building with a wide cornice no longer shrinks to fit its overhang.
        "base": strict if exact else width,
        "exact": exact,
        "left": edge(-1) if exact else left_px,
        "right": edge(1) if exact else right_px,
        "corner": [corner_x, corner_y],
    }


def proof(im: Image.Image, found: dict, path: str) -> None:
    """Draw the footprint that was found over the picture it was found in."""
    shot = im.convert("RGBA").copy()
    pen = ImageDraw.Draw(shot)
    cx, cy = found["corner"]
    half = found["base"] / 2
    pen.polygon([(cx, cy), (cx - half, cy - half * SLOPE),
                 (cx, cy - half * SLOPE * 2), (cx + half, cy - half * SLOPE)],
                outline=(255, 0, 255, 255))
    pen.line([(cx - 12, cy), (cx + 12, cy)], fill=(255, 220, 0, 255), width=3)
    shot.save(path)


def run(out_dir: str, proofs: str | None = None) -> int:
    entries = []
    for name in sorted(os.listdir(out_dir)):
        if not name.startswith("iso-") or not name.endswith(".png"):
            continue
        place = name[len("iso-"):].rsplit("-v", 1)[0]
        with Image.open(os.path.join(out_dir, name)) as im:
            im = im.convert("RGBA")
            found = fit(im)
            entry = {"id": place, "file": f"iso/{name}", "w": im.width, "h": im.height}
            if found:
                entry["anchor"] = found["anchor"]
                entry["base"] = found["base"]
                print(f"{place:14s} base {found['base']:4d}px of {im.width:4d}  "
                      f"feet at {found['corner']}  "
                      f"{'measured' if found['exact'] else 'from the content box'}")
                if proofs:
                    os.makedirs(proofs, exist_ok=True)
                    proof(im, found, os.path.join(proofs, name))
            entries.append(entry)
    json.dump(entries, open(os.path.join(out_dir, "isometric.json"), "w"), indent=1)
    return len(entries)


if __name__ == "__main__":
    args = [a for a in sys.argv[1:] if not a.startswith("--")]
    target = args[0] if args else "public/art/iso"
    shots = args[1] if len(args) > 1 else ".runtime/iso-fit"
    run(target, shots if "--proof" in sys.argv else None)
