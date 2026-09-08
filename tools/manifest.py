"""Rebuild the isometric manifest from whatever cut-outs are on disk.

The painter writes the manifest when it finishes, which is no use while it is
still working and no use at all if it stops half way. This reads the directory,
so the interface can be built against a city that is partly painted — the
addresses without a picture keep their blocked-out solid.
"""
import json
import os
import sys

from PIL import Image


def rebuild(out_dir: str) -> int:
    entries = []
    for name in sorted(os.listdir(out_dir)):
        if not name.startswith("iso-") or not name.endswith(".png"):
            continue
        place = name[len("iso-"):].rsplit("-v", 1)[0]
        with Image.open(os.path.join(out_dir, name)) as im:
            entries.append({"id": place, "file": f"iso/{name}", "w": im.width, "h": im.height})
    json.dump(entries, open(os.path.join(out_dir, "isometric.json"), "w"), indent=1)
    return len(entries)


if __name__ == "__main__":
    target = sys.argv[1] if len(sys.argv) > 1 else "public/art/iso"
    print(f"{rebuild(target)} addresses painted")
