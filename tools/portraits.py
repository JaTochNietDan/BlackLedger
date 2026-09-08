"""Generate the city's faces once, offline, and bake them into a sprite sheet.

Nothing at runtime depends on this. The game ships PNGs; this is the thing that
made them, kept in the repo so the next person can regenerate or extend the cast
without guessing at the prompt.

Run it through mise, which owns the toolchain:

    mise run art        # once, to build the local venv
    mise run portraits  # to regenerate the sheet

The model is FLUX.1-schnell, quantised to 8-bit and run on the Apple Silicon GPU
through mflux. The first cut of this used SD-turbo and it showed badly beside the
six portraits that were made by hand: flat, plasticky, obviously cheaper. The
brief asks for the interface to be fleshed out, and a face nobody believes is
worse than no face.

The people of Bellwether are not gendered by the rules — the city says "they"
about everybody — so the pool is varied and assigned to a person by a hash of
their id. The same person is always the same face, and any face can belong to
any name.
"""
import os
import subprocess
import sys
import tempfile

from PIL import Image

# The look of the six that were made by hand: oil on board, one hard light from
# the side, a dark olive-brown ground, nobody smiling, nothing modern in frame.
LOOK = (
    "oil painting portrait on board, 1950s, film noir, single hard key light "
    "from one side, deep shadow on the other, dark olive-brown background, "
    "muted desaturated palette, visible brushwork and canvas grain, period "
    "wool suit and soft collar, head and shoulders, looking at the viewer, "
    "unsmiling, sombre, painted by a mid-century portraitist"
)

SUBJECTS = [
    "a gaunt man in his fifties with hollow cheeks",
    "a heavy-set man with a broken nose and thinning hair",
    "a woman in her thirties with dark waved hair",
    "a young man with a thin moustache and slicked hair",
    "an older woman with pinned-back grey hair",
    "a wiry man in his forties with a hard jaw",
    "a broad woman in her fifties with a heavy brow",
    "a boyish man in his twenties with an open face",
    "a hard-faced woman in her forties with tied-back hair",
    "a bald man with heavy brows and a thick neck",
    "a tired man in his thirties with a shaving cut",
    "a well-kept woman in her twenties with pale eyes",
    "a stooped elderly man with white stubble",
    "a freckled woman in her thirties with short hair",
    "a scarred man in his thirties with a crooked mouth",
    "a plump older man in wire spectacles",
    "a severe woman in a high collar with dark eyes",
    "a fair-haired man in his twenties with a long face",
    "a weathered dockworker with a lined face",
    "a neat clerk in his forties with a side parting",
    "a woman with a heavy jaw and a flat stare",
    "a man with sunken eyes and a widow's peak",
    "a matronly woman in her sixties with soft features",
    "a sharp-eyed young woman with a narrow face",
]

COLS, CELL = 6, 256


MODEL = "mflux-community/flux-1-schnell-mflux-q8"


def generate(prompt, seed, path):
    """One face, through the mflux CLI. The CLI is the interface mflux keeps
    stable across versions; its Python module layout is not."""
    subprocess.run([
        os.path.join(os.path.dirname(sys.executable), "mflux-generate"),
        "--model", MODEL, "--base-model", "schnell",
        "--steps", "4", "--height", "512", "--width", "512",
        "--seed", str(seed), "--no-metadata", "--output", path,
        "--prompt", prompt,
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)


def main(out_dir, count):
    work = tempfile.mkdtemp(prefix="faces-")
    faces = []
    for i in range(count):
        subject = SUBJECTS[i % len(SUBJECTS)]
        path = os.path.join(work, f"{i:02d}.png")
        generate(f"{subject}, {LOOK}", 1000 + i * 7, path)
        faces.append(Image.open(path).convert("RGB").resize((CELL, CELL), Image.LANCZOS))
        print(f"  {i + 1}/{count} {subject}", flush=True)

    rows = (len(faces) + COLS - 1) // COLS
    sheet = Image.new("RGB", (COLS * CELL, rows * CELL), (18, 16, 14))
    for i, face in enumerate(faces):
        sheet.paste(face, ((i % COLS) * CELL, (i // COLS) * CELL))
    os.makedirs(out_dir, exist_ok=True)
    path = os.path.join(out_dir, "faces-noir-v1.png")
    sheet.save(path, optimize=True)
    print(f"SHEET {path} {sheet.size} {len(faces)} faces, {COLS} across, {rows} down")


if __name__ == "__main__":
    main(sys.argv[1] if len(sys.argv) > 1 else "public/art",
         int(sys.argv[2]) if len(sys.argv) > 2 else 24)
