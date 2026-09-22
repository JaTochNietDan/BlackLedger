#!/usr/bin/env python3
"""Build an archive with explicit inputs; never bundle saves or local services."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tarfile
import tempfile
import zipfile

ROOT = Path(__file__).resolve().parents[1]
TARGETS = {('windows', 'amd64'), ('darwin', 'amd64'), ('darwin', 'arm64'), ('linux', 'amd64'), ('linux', 'arm64')}


def run(*args, **kwargs):
    return subprocess.check_output(args, cwd=ROOT, text=True, **kwargs).strip()


def notices(destination):
    """Retain license texts for all locked npm packages and Go modules (a superset)."""
    destination.mkdir()
    inventory = []
    lock = json.loads((ROOT / 'package-lock.json').read_text())
    for name, metadata in lock['packages'].items():
        if not name or metadata.get('optional') and not (ROOT / name).exists():
            continue
        directory = ROOT / name
        if not directory.is_dir():
            raise RuntimeError(f'Missing dependency: {name}; run npm ci')
        inventory.append({'package': name, 'version': metadata.get('version'), 'license': metadata.get('license')})
        # Platform binary packages omit a copy; their parent ships the notice.
        if name.startswith('node_modules/@esbuild/'):
            directory = ROOT / 'node_modules/esbuild'
        elif name.startswith('node_modules/@rollup/'):
            directory = ROOT / 'node_modules/rollup'
        elif name == 'node_modules/@pixi/colord':
            directory = ROOT / 'packaging/licenses'
        copy_notices(directory, destination / name.replace('/', '_'))
    subprocess.run(['go', 'mod', 'download'], cwd=ROOT, check=True)
    remaining = run('go', 'list', '-m', '-json', 'all')
    decoder = json.JSONDecoder()
    while remaining.strip():
        module, end = decoder.raw_decode(remaining.lstrip())
        remaining = remaining.lstrip()[end:]
        if module.get('Main'):
            continue
        inventory.append({'module': module['Path'], 'version': module.get('Version')})
        copy_notices(Path(module['Dir']), destination / module['Path'].replace('/', '_'))
    shutil.copyfile(Path(run('go', 'env', 'GOROOT')) / 'LICENSE', destination / 'Go-LICENSE')
    (destination / 'inventory.json').write_text(json.dumps(inventory, indent=2) + '\n')


def copy_notices(source, destination):
    files = [p for p in source.iterdir() if p.is_file() and p.name.lower().startswith(('license', 'copying', 'notice', 'copyright', 'colord-license'))]
    if not files:
        raise RuntimeError(f'No license notice found for {source}')
    destination.mkdir()
    for path in files:
        shutil.copyfile(path, destination / path.name)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--os', choices=['windows', 'darwin', 'linux'], required=True)
    parser.add_argument('--arch', choices=['amd64', 'arm64'], required=True)
    parser.add_argument('--version', required=True)
    parser.add_argument('--output', default='release')
    args = parser.parse_args()
    if (args.os, args.arch) not in TARGETS or not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9._-]{0,79}', args.version):
        parser.error('unsupported target or unsafe version label')
    if not (ROOT / 'dist/index.html').is_file():
        parser.error('run npm ci && npm run build first')
    modified = bool(run('git', 'status', '--porcelain'))
    if args.version.startswith('v') and modified:
        parser.error('versioned releases require a clean committed checkout; use a dev- label for local QA')
    output = Path(args.output).resolve()
    output.mkdir(parents=True, exist_ok=True)
    name = f'black-ledger-{args.version}-{args.os}-{args.arch}'
    suffix = '.zip' if args.os == 'windows' else '.tar.gz'
    archive = output / (name + suffix)
    if archive.exists():
        parser.error(f'refusing to overwrite {archive}')
    with tempfile.TemporaryDirectory(prefix='black-ledger-package-') as temp:
        stage = Path(temp) / 'game'
        stage.mkdir()
        executable = 'blackledger.exe' if args.os == 'windows' else 'blackledger'
        subprocess.run(['go', 'build', '-trimpath', '-ldflags=-s -w', '-o', str(stage / executable), './cmd/blackledger'], cwd=ROOT,
                       env={**os.environ, 'GOOS': args.os, 'GOARCH': args.arch, 'CGO_ENABLED': '0'}, check=True)
        shutil.copytree(ROOT / 'dist', stage / 'dist')
        for file in ['LICENSE', 'NOTICE', 'README.md', 'ASSETS.md', 'THIRD_PARTY_NOTICES.md', 'CONTRIBUTING.md', 'API.md']:
            shutil.copyfile(ROOT / file, stage / file)
        (stage / 'docs').mkdir()
        for file in ['GOAL.md', 'RELEASING.md', 'RELEASE_READINESS.md', 'MEDIA_RIGHTS_REVIEW.md']:
            shutil.copyfile(ROOT / 'docs' / file, stage / 'docs' / file)
        shutil.copyfile(ROOT / 'packaging/PLAY.txt', stage / 'PLAY.txt')
        notices(stage / 'licenses')
        metadata = {'version': args.version, 'os': args.os, 'arch': args.arch,
                    'revision': run('git', 'rev-parse', 'HEAD'),
                    'modified': modified,
                    'go': run('go', 'version'), 'node': run('node', '--version'),
                    'electron': json.loads((ROOT / 'desktop/package.json').read_text())['devDependencies']['electron']}
        (stage / 'build.json').write_text(json.dumps(metadata, indent=2) + '\n')
        if args.os != 'windows':
            (stage / executable).chmod(0o755)
        packaged = Path(temp) / 'packaged'
        subprocess.run(['node', str(ROOT / 'desktop/package.mjs'), str(stage), str(packaged),
                        'win32' if args.os == 'windows' else args.os,
                        'x64' if args.arch == 'amd64' else args.arch], cwd=ROOT, check=True)
        app = Path((packaged / 'package-path.txt').read_text())
        # Packager leaves Electron notices beside the app. Retain them inside
        # the Mac bundle too, so moving only the .app cannot lose attribution.
        resources = app / ('Black Ledger.app/Contents/Resources/game' if args.os == 'darwin' else 'resources/game')
        for source_name, notice_name in [('LICENSE', 'Electron-LICENSE'), ('LICENSES.chromium.html', 'Chromium-LICENSES.html')]:
            source = app / source_name
            if not source.is_file():
                raise RuntimeError(f'Electron notice missing: {source}')
            shutil.copyfile(source, resources / 'licenses' / notice_name)
        (app / 'LICENSE').rename(app / 'LICENSE.electron.txt')
        # Keep instructions and attribution accessible without opening app internals.
        for file in ['PLAY.txt', 'LICENSE', 'NOTICE', 'THIRD_PARTY_NOTICES.md', 'build.json']:
            shutil.copyfile(stage / file, app / file)
        stage = app.rename(Path(temp) / name)
        if args.os == 'windows':
            with zipfile.ZipFile(archive, 'w', zipfile.ZIP_DEFLATED) as bundle:
                for path in sorted(stage.rglob('*')):
                    if path.is_file():
                        bundle.write(path, path.relative_to(stage.parent))
        else:
            with tarfile.open(archive, 'w:gz') as bundle:
                bundle.add(stage, arcname=name)
    with archive.open('rb') as data:
        digest = hashlib.file_digest(data, 'sha256').hexdigest()
    archive.with_name(archive.name + '.sha256').write_text(f'{digest}  {archive.name}\n')
    print(archive)


if __name__ == '__main__':
    main()
