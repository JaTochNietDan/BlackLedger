"""Take a picture somebody generated and make it a cut-out this city can use.

Art from outside arrives framed however the model felt like framing it: on a
flat background rather than on nothing, with the building anywhere in the frame,
sometimes running off the edge of it. This does the mechanical part — lifts the
background off, trims the air, measures the footprint — and refuses the ones
that cannot be fixed mechanically, with a reason, rather than importing
something that will look wrong in the city and leave nobody any the wiser.

    .venv/bin/python tools/importart.py speakeasy.png bar

The second argument is the address the picture is for.
"""
import os
import shutil
import sys

from PIL import Image

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from fitiso import ALPHA, SLOPE, fit                     # noqa: E402
from isometric import knock_out                          # noqa: E402

# How much clear air a picture needs around its building. A building that runs
# off the edge of its own frame has been cropped, and the part that is missing
# is usually the corner of its base — which is the one part the city needs, so
# it can stand the building on its own feet.
MARGIN = 6

# How far the two base edges may differ before the picture is not in this
# city's projection at all. The view draws in 2:1: every base edge falls half a
# pixel for every pixel sideways. A picture drawn in perspective, or square on,
# or from a different height, will not sit on the grid however it is scaled, and
# no amount of fitting rescues it.
SKEW = 0.22


def edges(im: Image.Image) -> tuple[float, float]:
    """The slope of each base edge, out from the lowest point of the building."""
    al = im.getchannel("A").load()
    w, h = im.size
    low = []
    for x in range(w):
        y = None
        for yy in range(h - 1, -1, -1):
            if al[x, yy] >= ALPHA:
                y = yy
                break
        low.append(y)
    solid = [x for x, y in enumerate(low) if y is not None]
    deepest = max(low[x] for x in solid)
    flat = [x for x in solid if low[x] == deepest]
    cx = (flat[0] + flat[-1]) // 2

    def slope(step: int) -> float:
        xs, ys = [], []
        for k in range(15, 200):
            x = cx + step * k
            if not (0 <= x < w) or low[x] is None:
                break
            xs.append(k)
            ys.append(deepest - low[x])
        if len(xs) < 30:
            return float("nan")
        n = len(xs)
        sx, sy = sum(xs), sum(ys)
        sxx = sum(v * v for v in xs)
        sxy = sum(a * b for a, b in zip(xs, ys))
        return (n * sxy - sx * sy) / (n * sxx - sx * sx)

    return slope(-1), slope(1)


def check(path: str) -> list[str]:
    """Everything wrong with this picture, in words that say what to do."""
    im = Image.open(path).convert("RGBA")
    w, h = im.size
    box = im.getchannel("A").point(lambda v: 255 if v >= ALPHA else 0).getbbox()
    if not box:
        return ["the picture is empty"]
    faults = []

    touching = []
    if box[0] < MARGIN: touching.append("left")
    if box[1] < MARGIN: touching.append("top")
    if box[2] > w - MARGIN: touching.append("right")
    if box[3] > h - MARGIN: touching.append("bottom")
    if touching:
        faults.append(
            f"the building runs off the {', '.join(touching)} of its own frame. "
            "It has been cropped, and what is missing at the bottom is the corner "
            "of its base — the one part the city needs in order to stand it on its "
            "own feet. Ask for the whole building inside the frame with clear air "
            "around it.")

    left, right = edges(im)
    if left == left and right == right:                  # not NaN
        if abs(abs(left) - SLOPE) > SKEW or abs(abs(right) - SLOPE) > SKEW \
                or abs(abs(left) - abs(right)) > SKEW:
            faults.append(
                f"the base edges fall at {abs(left):.2f} and {abs(right):.2f} where "
                f"this city's projection wants {SLOPE:.2f} and {SLOPE:.2f}. The "
                "picture is not in 2:1 isometric — it is in perspective, or turned, "
                "or seen from a different height — and it will not sit on the grid "
                "however it is scaled. Ask for orthographic 2:1 isometric with no "
                "perspective and both faces equal.")
    return faults


def bring_in(source: str, place: str, out_dir: str = "public/art/iso") -> bool:
    target = os.path.join(out_dir, f"iso-{place}-v2.png")
    shutil.copy(source, target)
    knock_out(target)                                    # the flat field comes off
    faults = check(target)
    if faults:
        os.remove(target)
        print(f"not imported: {os.path.basename(source)}")
        for f in faults:
            print(f"  - {f}")
        return False
    with Image.open(target) as im:
        found = fit(im.convert("RGBA"))
    print(f"imported as {target}")
    print(f"  base {found['base']}px, feet at {found['corner']}, "
          f"{'measured' if found['exact'] else 'from the content box'}")
    print("  now run: mise run fit-art")
    return True


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print(__doc__)
        raise SystemExit(2)
    raise SystemExit(0 if bring_in(sys.argv[1], sys.argv[2]) else 1)
