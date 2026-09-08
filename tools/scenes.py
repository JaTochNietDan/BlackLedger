"""Paint the moments the theatre plays.

The theatre takes the player to the building where something happened and shows
it: a killing, an explosion, a raid, ground changing hands. It was drawn — a
black box of a building and a few animated shapes — which was enough to say
where and not enough to make anybody look.

These are backdrops. The animation still plays over the top, out of the cue the
core committed, so what moves is still the truth of the event and only what it
is painted on has changed. Nothing at runtime depends on a model.

    mise run art
    mise run scenes
"""
import os
import subprocess
import sys
import tempfile

from PIL import Image

LOOK = (
    "1950s city at night, film noir, oil painting on board, wet street, hard "
    "shadows, single light source, dark umber and olive palette, muted and "
    "desaturated, visible brushwork, no text or lettering, no faces visible, "
    "painted by a mid-century illustrator"
)

# One plate a kind. They are the aftermath or the middle of the thing, framed
# wide, with nobody identifiable in them — the people who were actually there
# are drawn on top from what the city committed.
SCENES = {
    "killing": "a body under a sheet on a wet pavement outside a bar, police lamp light, onlookers as distant silhouettes",
    "explosion": "the blown-out front of a small city building, smoke, scattered brick and glass, fire in an upper window",
    "gunfight": "an empty street corner after shooting, broken window, shell casings on wet asphalt, one car door open",
    "raid": "police cars at the kerb outside a shop at night, doors open, torch beams into a doorway",
    "seizure": "a shuttered shopfront with a notice pasted on the door, chain and padlock, nobody about",
    "arrest": "the steps of a police station at night under a lamp, a car pulled up, an open rear door",
    "attack": "a smashed shopfront at night, glass across the pavement, one light still burning inside",
    "robbery": "an emptied till and an open back door in a dark shop, papers on the floor",
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
    work = tempfile.mkdtemp(prefix="scenes-")
    for i, (kind, what) in enumerate(sorted(SCENES.items())):
        raw = os.path.join(work, f"{kind}.png")
        generate(f"{what}, {LOOK}", 9000 + i * 23, raw)
        out = os.path.join(out_dir, f"scene-{kind}-v1.jpg")
        Image.open(raw).convert("RGB").save(out, quality=82, optimize=True)
        print(f"  {i + 1}/{len(SCENES)} {kind} -> {out}", flush=True)


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art/scenes")
