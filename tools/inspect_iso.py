"""Find cut-outs that are not cut out.

The model is asked for a building standing on nothing. Sometimes it paints a
floor across the whole frame instead, and the flood fill from the corners
stops at it because a pale floor is not the black field it was told to remove.
The result is a building with a slab stuck to it, which on the map reads as a
rectangle of nothing beside the building.

A cut-out is suspect when its bottom edge is opaque from side to side: a
building has a footprint, a floor has the frame.
"""
import json
import os
import sys

from PIL import Image


def suspect(path: str) -> tuple[bool, str]:
    im = Image.open(path).convert("RGBA")
    w, h = im.size
    px = im.load()
    bottom = sum(1 for x in range(w) if px[x, h - 1][3] > 128)
    edges = sum(1 for y in range(h) if px[0, y][3] > 128) + sum(1 for y in range(h) if px[w - 1, y][3] > 128)
    share = bottom / w
    if share > .92:
        return True, f"opaque across {share:.0%} of its bottom edge: a floor, not a footprint"
    if edges > h * .5:
        return True, f"opaque down {edges / (2 * h):.0%} of both sides: the frame was painted in"
    return False, f"bottom {share:.0%}"


def main(out_dir: str) -> int:
    bad = 0
    for entry in json.load(open(os.path.join(out_dir, "isometric.json"))):
        path = os.path.join(os.path.dirname(out_dir), entry["file"])
        wrong, why = suspect(path)
        print(f"{'SUSPECT' if wrong else 'ok     '} {entry['id']:<10} {why}")
        bad += wrong
    return bad


if __name__ == "__main__":
    sys.exit(1 if main(sys.argv[1] if len(sys.argv) > 1 else "public/art/iso") else 0)
