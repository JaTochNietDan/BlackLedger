"""Take the wrong words off the buildings.

The model cannot spell. Asked for a nightclub it paints a marquee reading
CIICATCE; asked for a newspaper office it paints NEWSAPWL. Asking for no
lettering does not stop it — that was tried, and the replacement Monarch came
back reading NHTIUUR.

What jars is not that there is a sign. It is that the sign is legibly wrong. A
lit panel with an indistinct glow on it reads as a sign at every zoom this city
is drawn at; a lit panel with misspelt words reads as a mistake.

So this finds the letterforms and softens them: small, high-contrast, saturated
clusters sitting on a bright ground — which is what text on a lit sign is — and
blurs those regions only, leaving windows, brickwork and roofs alone.
"""
import sys

from PIL import Image, ImageFilter


def unletter(path: str, out: str | None = None) -> int:
    im = Image.open(path).convert("RGBA")
    rgb = im.convert("RGB")
    # Edges find letterforms: text is the highest-frequency thing on a facade.
    edges = rgb.convert("L").filter(ImageFilter.FIND_EDGES).filter(ImageFilter.MaxFilter(5))
    px = edges.load()
    src = rgb.load()
    alpha = im.split()[3].load()

    # A pixel is lettering when it sits on a sharp edge, is bright, and is
    # warm — the colour a lit sign is painted in this city.
    mask = Image.new("L", im.size, 0)
    mp = mask.load()
    touched = 0
    for y in range(im.height):
        for x in range(im.width):
            if alpha[x, y] < 20:
                continue
            r, g, b = src[x, y]
            bright = (r + g + b) / 3
            # A sign is far more saturated than stone or brick. The first
            # version asked only for "warm", which is true of every lit window,
            # every cornice catching the light and most of the facade — so it
            # softened the whole building. Painted lettering is the one thing
            # up there that is both very saturated and very bright.
            saturation = (max(r, g, b) - min(r, g, b))
            if px[x, y] > 130 and bright > 120 and saturation > 78 and r > b + 60:
                mp[x, y] = 255
                touched += 1
    if touched == 0:
        return 0
    # Spread the mask just enough to cover a letter rather than outline it.
    mask = mask.filter(ImageFilter.MaxFilter(5)).filter(ImageFilter.GaussianBlur(1.6))
    softened = im.filter(ImageFilter.GaussianBlur(3.0))
    im = Image.composite(softened, im, mask)
    im.save(out or path)
    return touched


# Deliberately not a batch step over the whole art directory. Run over
# everything it softens lit windows and glazing bars — the most attractive
# thing these buildings have — because a bright warm window is not far off a
# bright warm sign. That was tried, and it took the mullions out of the bar's
# windows and dulled the market's colonnade for no gain, since the only two
# buildings with wrong words on them were the Monarch and the Herald. Name the
# files that actually have a garbled sign; leave the rest alone.
if __name__ == "__main__":
    total = 0
    for path in sys.argv[1:]:
        n = unletter(path)
        print(f"{path}: softened {n} pixels of lettering")
        total += n
    sys.exit(0 if total else 1)
