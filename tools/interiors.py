"""Generate the inside of every building once, offline.

Same contract as tools/portraits.py: nothing at runtime depends on a model
being installed, because the game ships the PNGs. Run through mise:

    mise run art        # once, to build the local venv
    mise run interiors  # to regenerate the rooms

The rooms are painted empty. The people who are standing in one are drawn over
the top by the interface, out of who the core says is actually there, so a room
never disagrees with the city about who is in it.
"""
import os
import subprocess
import sys
import tempfile

from PIL import Image

LOOK = (
    "1950s interior, film noir, oil painting on board, single hard light source, "
    "deep shadow, dark olive-brown and umber palette, muted and desaturated, "
    "visible brushwork, empty room with no people, wide establishing view from "
    "eye level, period fittings, painted by a mid-century illustrator"
)

# One room per address. The words are what the place is, not what it looks like
# in some other game — a laundry is presses and steam, a station is tile and a
# counter, and the Herald is three floors of typewriters.
ROOMS = {
    "room": "a cheap rented room above a bar, iron bedstead, washstand, one window with a drawn blind",
    "bar": "a small city bar, long wooden counter, bottles on a back shelf, booths along one wall",
    "docks": "a dockside cargo shed, crates and pallets, chain hoist, wide loading doors open to grey water",
    "laundry": "a commercial laundry back room, steam presses, hanging sheets, tiled floor, brass pipework",
    "club": "a plush casino floor, roulette table, brass rail, heavy curtains, low chandeliers",
    "market": "a covered market exchange hall, wooden stalls, hanging scales, high dusty windows",
    "apartment": "a respectable apartment sitting room, upholstered chairs, a mantelpiece, patterned wallpaper",
    "garage": "a motor repair garage, a car on a lift, tool boards, oil drums, concrete floor",
    "casino": "an upstairs gambling room, card tables under green shaded lamps, panelled walls",
    "estate": "the drawing room of a large house, tall windows, heavy drapes, a grand fireplace",
    "precinct": "a police station front office, tiled floor, high wooden counter, notice boards, barred inner door",
    "herald": "a newspaper composing room, rows of typewriters, paper spikes, hanging lamps, printing press beyond",
}

W, H = 832, 512


def generate(prompt, seed, path):
    subprocess.run([
        os.path.join(os.path.dirname(sys.executable), "mflux-generate"),
        "--model", "mflux-community/flux-1-schnell-mflux-q8", "--base-model", "schnell",
        "--steps", "4", "--height", str(H), "--width", str(W),
        "--seed", str(seed), "--no-metadata", "--output", path, "--prompt", prompt,
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def main(out_dir):
    os.makedirs(out_dir, exist_ok=True)
    work = tempfile.mkdtemp(prefix="rooms-")
    for i, (place, what) in enumerate(sorted(ROOMS.items())):
        raw = os.path.join(work, f"{place}.png")
        generate(f"{what}, {LOOK}", 4000 + i * 13, raw)
        out = os.path.join(out_dir, f"room-{place}-v1.jpg")
        # Saved as JPEG: these are backdrops behind figures, and a megabyte a
        # room would be four times the weight of the whole rest of the game.
        Image.open(raw).convert("RGB").resize((W, H), Image.LANCZOS).save(
            out, quality=82, optimize=True)
        print(f"  {i + 1}/{len(ROOMS)} {place} -> {out}", flush=True)


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art/rooms")
