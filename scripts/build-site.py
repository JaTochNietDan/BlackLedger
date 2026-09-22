#!/usr/bin/env python3
"""Build the public landing page, excluding game code, saves and local services."""
import argparse
import json
from pathlib import Path
import shutil

ROOT = Path(__file__).resolve().parents[1]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--output', type=Path, default=ROOT / '.runtime/site')
    args = parser.parse_args()
    destination = args.output.resolve()
    if destination == ROOT or ROOT.is_relative_to(destination):
        parser.error('output must not be the repository root or its parent')
    # Never clean arbitrary paths; staging must be new or already a site output.
    if destination.exists():
        if not (destination / '.black-ledger-site').is_file():
            parser.error('refusing to overwrite an unrecognized output directory')
        shutil.rmtree(destination)
    destination.mkdir(parents=True)
    (destination / '.black-ledger-site').touch()
    data = json.loads((ROOT / 'docs/examples/director-encounter.json').read_text())
    assert data['event']['source'] == 'local-ai'
    html = (ROOT / 'website/index.html').read_text()
    assert html.count('__ENCOUNTER_JSON__') == 1
    encoded = json.dumps(data, ensure_ascii=False).replace('<', '\\u003c')
    (destination / 'index.html').write_text(html.replace('__ENCOUNTER_JSON__', encoded))
    for name in ['styles.css', 'site.js']:
        shutil.copyfile(ROOT / 'website' / name, destination / name)
    assets = destination / 'assets'
    assets.mkdir()
    for name in ['bellwether-city.jpg', 'saint-agnes.jpg', 'green-baize.jpg']:
        shutil.copyfile(ROOT / 'docs/screenshots' / name, assets / name)
    shutil.copyfile(ROOT / 'website/favicon.svg', assets / 'favicon.svg')
    (destination / '.nojekyll').touch()
    (destination / 'robots.txt').write_text('User-agent: *\nAllow: /\n')
    print(destination)


if __name__ == '__main__':
    main()
