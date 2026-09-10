"""Paint the cars on the forecourt.

"It would also be good if we generated nice car images to show what you're
buying." The lot sold three cars and showed none of them, so a Ford, a Hudson
with a false floor and an armoured Packard were three lines of text that looked
identical on the way past.

Same contract as the other art tools: generated offline, shipped as files,
nothing at runtime depends on a model. Run through mise:

    mise run art     # once
    mise run cars
"""
import hashlib
import os
import subprocess
import sys
import tempfile

from PIL import Image

LOOK = (
    "1930s American automobile, film noir, oil painting on board, side "
    "three-quarter view, parked on a wet forecourt at dusk, overcast light, "
    "dark umber and olive palette, muted and desaturated, visible brushwork, "
    "no people, no text or signage lettering, painted by a mid-century "
    "illustrator"
)

# Keyed by the tier in core/vehicle.go, which is what the interface asks for.
CARS = {
    "1": "a tired black Ford sedan, dull paint, one dented wing, narrow whitewall tyres",
    "2": "a respectable dark green Hudson sedan, deep body, chrome grille, running boards, well kept",
    "3": "a heavy black Packard sedan, thick pillars, small deep-set windows, plated doors, sitting low on its springs",
}

W, H = 768, 512


def seed_for(tier):
    return 9000 + int(hashlib.sha256(("car" + tier).encode()).hexdigest()[:6], 16) % 100000


def generate(prompt, seed, path):
    subprocess.run([
        os.path.join(os.path.dirname(sys.executable), "mflux-generate"),
        "--model", "mflux-community/flux-1-schnell-mflux-q8", "--base-model", "schnell",
        "--steps", "4", "--height", str(H), "--width", str(W),
        "--seed", str(seed), "--no-metadata", "--output", path, "--prompt", prompt,
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def main(out_dir):
    os.makedirs(out_dir, exist_ok=True)
    work = tempfile.mkdtemp(prefix="cars-")
    for i, (tier, what) in enumerate(sorted(CARS.items())):
        out = os.path.join(out_dir, f"car-{tier}-v1.jpg")
        if os.path.exists(out) and "--all" not in sys.argv:
            print(f"  {i + 1}/{len(CARS)} tier {tier} already painted", flush=True)
            continue
        raw = os.path.join(work, f"{tier}.png")
        generate(f"{what}, {LOOK}", seed_for(tier), raw)
        Image.open(raw).convert("RGB").save(out, quality=82, optimize=True)
        print(f"  {i + 1}/{len(CARS)} tier {tier} -> {out}", flush=True)


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art/cars")
