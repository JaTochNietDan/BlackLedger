"""Generate the inside of every building once, offline.

Same contract as tools/portraits.py: nothing at runtime depends on a model
being installed, because the game ships the PNGs. Run through mise:

    mise run art        # once, to build the local venv
    mise run interiors  # to regenerate the rooms

The rooms are painted empty. The people who are standing in one are drawn over
the top by the interface, out of who the core says is actually there, so a room
never disagrees with the city about who is in it.
"""
import hashlib
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
    "restaurant": "a small restaurant dining room, tables with white cloths, a curtained door to a back room, wall lamps",
    "goldenlily": "a small private gambling room, one green baize table, low shaded lamp, panelled walls, heavy curtains",
    "steamworks": "a commercial laundry floor, rows of washing machines and mangles, steam, drying racks on rails",
    "burlesque": "a revue bar interior, a small stage with a curtain, round tables, bar along one wall, footlights",
    "cabstand": "a taxi dispatch office, a wall of hooks and route cards, a radio set on the counter, cabs through the window",
    "poolhall": "a billiard hall, three tables under low hanging lamps, cue racks on the wall, a payphone in the corner",
    "butcher": "a butcher shop interior, marble counter, hooks and rails, a heavy cold room door at the back",
    "haulage": "a haulage yard office, a wall of route boards and keys, a counter, trucks visible through the window",
}

W, H = 832, 512


def seed_for(place):
    """A place's seed comes from its name, not its position in this list.

    It used to be 4000 + index, so adding four addresses to the city shifted
    every index after them and would have repainted sixteen rooms that were
    already finished. A room should look the same tomorrow as it does today.
    """
    return 4000 + int(hashlib.sha256(place.encode()).hexdigest()[:6], 16) % 100000


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
        out = os.path.join(out_dir, f"room-{place}-v1.jpg")
        # Paint what has no picture; leave a finished room alone. See the note
        # in tools/exteriors.py.
        if os.path.exists(out) and "--all" not in sys.argv:
            print(f"  {i + 1}/{len(ROOMS)} {place} already painted", flush=True)
            continue
        raw = os.path.join(work, f"{place}.png")
        generate(f"{what}, {LOOK}", seed_for(place), raw)
        # Saved as JPEG: these are backdrops behind figures, and a megabyte a
        # room would be four times the weight of the whole rest of the game.
        Image.open(raw).convert("RGB").resize((W, H), Image.LANCZOS).save(
            out, quality=82, optimize=True)
        print(f"  {i + 1}/{len(ROOMS)} {place} -> {out}", flush=True)


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art/rooms")
