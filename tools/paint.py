"""One way to call the painter.

Every art tool in here shells out to the same local model with the same
arguments, and the first version of tools/isometric.py got them wrong — it
invented a --path flag mflux does not have and failed on the first building.
The invocation lives here now so there is one of it.

FLUX.1-schnell upstream is gated behind an account acceptance, so the weights
come from the ungated mflux-community 8-bit quantisation.
"""
import os
import subprocess
import sys

MODEL = "mflux-community/flux-1-schnell-mflux-q8"


def generate(prompt: str, seed: int, path: str, width: int = 768, height: int = 512, steps: int = 4) -> None:
    subprocess.run([
        os.path.join(os.path.dirname(sys.executable), "mflux-generate"),
        "--model", MODEL, "--base-model", "schnell",
        "--steps", str(steps), "--height", str(height), "--width", str(width),
        "--seed", str(seed), "--no-metadata", "--output", path, "--prompt", prompt,
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
