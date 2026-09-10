"""Paint the outside of the buildings that have none.

Five addresses were painted by hand for the street study. The other seven —
the docks, the apartment block, the garage, the Blue Hour, Cypress House, Ward
Street Station and the Herald — showed a wireframe box in the address book,
which made half the city look unfinished.

Same contract as the other art tools: generated offline, shipped as files,
nothing at runtime depends on a model. Run through mise:

    mise run art        # once
    mise run exteriors
"""
import hashlib
import os
import subprocess
import sys
import tempfile

from PIL import Image

LOOK = (
    "1950s American city building exterior, film noir, oil painting on board, "
    "three-quarter view from across the street, overcast evening, single warm "
    "light in the windows, dark umber and olive palette, muted and desaturated, "
    "visible brushwork, no people, no cars, no text or signage lettering, "
    "painted by a mid-century illustrator"
)

FRONTS = {
    "docks": "a waterfront cargo shed and pier, corrugated roof, timber piles, cranes behind",
    "apartment": "a respectable four-storey apartment block, stone stoop, iron railings, bay windows",
    "garage": "a motor repair garage on a corner, wide roller door, flat roof, brick front",
    "casino": "a discreet two-storey gambling club, awning over the door, shuttered upper windows",
    "estate": "a large house set back behind a wall, gables, tall chimneys, gravel drive",
    "precinct": "a grey stone police station, heavy steps, lamp either side of the door, barred windows",
    "herald": "a three-storey newspaper building, tall printing hall windows, loading bay at street level",
    "restaurant": "a small family restaurant on a corner, awning over the window, lace half-curtains, a lit doorway",
    "goldenlily": "a discreet riverside gambling club, stone steps up to a plain door, tall shuttered windows, no signage",
    "steamworks": "a large commercial steam laundry, brick, tall vent stacks, wide loading doors, delivery vans at the kerb",
    "burlesque": "a small revue theatre bar, canopy over the entrance, bulb-lit marquee frame, curtained upper windows",
    "cabstand": "a taxi company yard and office, low brick office, line of parked cabs, fuel pump, wire fence",
    "dealer": "a car dealership forecourt, rows of parked cars, low showroom with wide plate glass, pennant strings overhead",
    "archway": "a motor repair shop under a railway viaduct, brick arch, roller door, cars waiting at the kerb",
    "scrapyard": "a car breaker's yard behind a corrugated fence, stacked wrecks, a crane jib against the sky, weighbridge hut",
    "filling": "a two-pump filling station under a tin canopy, glass-topped pumps, an oil rack and a lit counter window",
    "pumps": "a filling station under a brick viaduct arch, two pumps on a cracked apron, a hand-painted price board",
    "poolhall": "a first-floor billiard hall over a shopfront, long low windows, a stair door at street level",
    "butcher": "a butcher's shop with a tiled front, wide window, delivery van at the kerb, cold store behind",
    "haulage": "a haulage yard behind a wire fence, flatbed trucks, a low office hut, fuel pump",
}

W, H = 768, 512


def seed_for(place):
    """A place's seed comes from its name, not its position in this list.

    It used to be 7000 + index, so adding four addresses to the city shifted
    every index after them and repainted eleven buildings that were already
    finished. A building should look the same tomorrow as it does today.
    """
    return 7000 + int(hashlib.sha256(place.encode()).hexdigest()[:6], 16) % 100000


def generate(prompt, seed, path):
    subprocess.run([
        os.path.join(os.path.dirname(sys.executable), "mflux-generate"),
        "--model", "mflux-community/flux-1-schnell-mflux-q8", "--base-model", "schnell",
        "--steps", "4", "--height", str(H), "--width", str(W),
        "--seed", str(seed), "--no-metadata", "--output", path, "--prompt", prompt,
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def main(out_dir):
    os.makedirs(out_dir, exist_ok=True)
    work = tempfile.mkdtemp(prefix="fronts-")
    for i, (place, what) in enumerate(sorted(FRONTS.items())):
        out = os.path.join(out_dir, f"front-{place}-v1.jpg")
        # Paint what has no picture. A building that is already finished is
        # left alone: running this to add one address must not repaint the
        # rest of the city, and "--all" is there for when that is the point.
        if os.path.exists(out) and "--all" not in sys.argv:
            print(f"  {i + 1}/{len(FRONTS)} {place} already painted", flush=True)
            continue
        raw = os.path.join(work, f"{place}.png")
        generate(f"{what}, {LOOK}", seed_for(place), raw)
        Image.open(raw).convert("RGB").save(out, quality=82, optimize=True)
        print(f"  {i + 1}/{len(FRONTS)} {place} -> {out}", flush=True)


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art/fronts")
