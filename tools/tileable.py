"""Say whether a texture tiles, and crop it so that it does.

A ground texture that does not tile shows its seam on every street in the city,
and the seam is the one fault that is invisible in the source image and
obvious the moment it is laid down. So it is measured rather than trusted.

The test is comparative rather than absolute: a texture tiles when its opposite
edges match each other about as well as any two strips inside it match. A
photograph of asphalt is not uniform, and demanding that its edges match
perfectly would reject every texture ever made.

The fix, where one is possible, is to crop rather than to blend. Regular
material — paving flags, cobbles, brick — repeats at some period, and if the
image contains a whole number of periods it tiles by construction. Blending the
edges together instead smears the courses, which on paving flags is worse than
the seam.

    .venv/bin/python tools/tileable.py art/textures/pavement.jpg --fix
"""
import os
import sys

from PIL import Image


def strip_diff(px, a: int, b: int, run: int, vertical: bool) -> float:
    total = 0
    for i in range(0, run, 3):
        p = px[a, i] if vertical else px[i, a]
        q = px[b, i] if vertical else px[i, b]
        total += abs(p[0] - q[0]) + abs(p[1] - q[1]) + abs(p[2] - q[2])
    return total / max(1, len(range(0, run, 3))) / 3


def report(path: str) -> tuple[float, float, float]:
    im = Image.open(path).convert("RGB")
    w, h = im.size
    px = im.load()
    lr = strip_diff(px, 0, w - 1, h, True)
    tb = strip_diff(px, 0, h - 1, w, False)
    inside = (strip_diff(px, w // 3, 2 * w // 3, h, True)
              + strip_diff(px, h // 3, 2 * h // 3, w, False)) / 2
    return lr, tb, inside


def period(px, size: int, run: int, vertical: bool) -> int:
    """The smallest crop at which the far edge matches the near one again."""
    best, score = size, 1e9
    for cut in range(int(size * 0.55), size):
        d = strip_diff(px, 0, cut - 1, run, vertical)
        if d < score * 0.999:
            best, score = cut, d
    return best


def fix(path: str) -> None:
    im = Image.open(path).convert("RGB")
    w, h = im.size
    px = im.load()
    keep_w = period(px, w, h, True)
    keep_h = period(px, h, w, False)
    im.crop((0, 0, keep_w, keep_h)).save(path, quality=95)
    print(f"cropped {w}x{h} -> {keep_w}x{keep_h}")


if __name__ == "__main__":
    files = [a for a in sys.argv[1:] if not a.startswith("--")]
    for f in files:
        lr, tb, inside = report(f)
        ok = lr <= inside * 1.6 and tb <= inside * 1.6
        print(f"{os.path.basename(f):14s} edges {lr:5.1f} / {tb:5.1f}   inside {inside:5.1f}   "
              f"{'tiles' if ok else 'DOES NOT TILE'}")
        if not ok and "--fix" in sys.argv:
            fix(f)
            lr, tb, inside = report(f)
            print(f"{'':14s} edges {lr:5.1f} / {tb:5.1f}   inside {inside:5.1f}   "
                  f"{'tiles now' if lr <= inside * 1.6 and tb <= inside * 1.6 else 'still does not tile'}")
