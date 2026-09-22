#!/usr/bin/env python3
"""Exercise a native release from an unrelated directory with isolated saves."""
import argparse
from contextlib import contextmanager
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tarfile
import tempfile
import time
import urllib.error
import urllib.request
import zipfile


def get(url):
    with urllib.request.urlopen(url, timeout=20) as response:
        return response.read()


@contextmanager
def server(executable, root, env, log_name):
    with (root / log_name).open('w+') as log:
        process = subprocess.Popen([str(executable), '-desktop', '-browser=false', '-addr', ':0'],
                                   cwd=root, env=env, stdout=log, stderr=log)
        try:
            deadline = time.monotonic() + 45
            while time.monotonic() < deadline:
                log.seek(0)
                match = re.search(r'http://127.0.0.1:\d+', log.read())
                if process.poll() is not None:
                    raise RuntimeError('server exited before readiness')
                if match:
                    base = match.group()
                    try:
                        assert json.loads(get(base + '/api/health'))['ok']
                        break
                    except (OSError, urllib.error.URLError):
                        pass
                time.sleep(.1)
            else:
                raise RuntimeError('server readiness timeout')
            yield base
        finally:
            process.terminate()
            try:
                process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait()
            log.seek(0)
            print(log.read())


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('archive', type=Path)
    args = parser.parse_args()
    checksum = args.archive.with_name(args.archive.name + '.sha256').read_text().split()[0]
    with args.archive.open('rb') as data:
        assert hashlib.file_digest(data, 'sha256').hexdigest() == checksum, 'archive checksum mismatch'
    with tempfile.TemporaryDirectory(prefix='black ledger smoke ') as temp:
        root = Path(temp)
        if args.archive.suffix == '.zip':
            with zipfile.ZipFile(args.archive) as archive:
                archive.extractall(root)
        else:
            with tarfile.open(args.archive) as archive:
                archive.extractall(root, filter='data')
        game = next(p for p in root.iterdir() if p.is_dir())
        for file in game.rglob('*'):
            assert not any(part in {'.runtime', '.tools', '.git', 'node_modules'} for part in file.relative_to(game).parts)
            assert '.sqlite' not in file.name, 'a save was packaged'
        assert (game / 'licenses/inventory.json').is_file()
        assert (game / 'LICENSE').is_file() and (game / 'NOTICE').is_file()
        executable = game / ('blackledger.exe' if os.name == 'nt' else 'blackledger')
        config = root / 'test user config'
        home = root / 'test home'
        if sys.platform == 'darwin':
            config = home / 'Library' / 'Application Support'
        db = config / 'BlackLedger' / 'campaign.sqlite3'
        env = {k: v for k, v in os.environ.items() if not k.startswith('BLACK_LEDGER_')}
        env.update({'HOME': str(home), 'APPDATA': str(config), 'XDG_CONFIG_HOME': str(config)})
        with server(executable, root, env, 'server.log') as base:
            html = get(base + '/').decode()
            assets = re.findall(r'(?:src|href)="(/assets/[^"]+)"', html)
            assert assets, 'page contains no frontend bundles'
            for asset in assets:
                assert get(base + asset)
            state = json.loads(get(base + '/api/state'))
            assert db.is_file(), 'per-user save was not created in isolated user config'
            # Occupied ports must fail before opening/mutating a save.
            blocked_db = root / 'must not exist.sqlite3'
            second = subprocess.run([str(executable), '-desktop', '-browser=false', '-addr', ':' + base.rsplit(':', 1)[1], '-db', str(blocked_db)],
                                    cwd=root, env=env, capture_output=True, timeout=10)
            assert second.returncode != 0 and not blocked_db.exists()
            command = {'kind': 'wait', 'target': state['player']['location'],
                       'request_id': 'release-smoke-wait', 'revision': state['revision']}
            def action(url):
                request = urllib.request.Request(url + '/api/action', data=json.dumps(command).encode(), headers={'Content-Type': 'application/json'})
                with urllib.request.urlopen(request, timeout=20) as response:
                    return json.load(response)
            committed = action(base)
            assert committed['revision'] == state['revision'] + 1
            assert action(base) == committed, 'command retry changed a committed outcome'
            assert not (root / '.runtime/campaign.sqlite3').exists()
        # Reopen the same save through the explicit environment override. The
        # receipt and the state must both survive stopping/restarting the process.
        with server(executable, root, {**env, 'BLACK_LEDGER_DB': str(db)}, 'reload.log') as base:
            restored = json.loads(get(base + '/api/state'))
            assert restored['revision'] == committed['revision']
            assert restored['player'] == committed['player']
            assert action(base) == committed, 'retry after restart changed committed state'
        print('PASS: checksum, archive, unrelated directory, frontend, per-user save, port conflict, command/retry, restart/retry')


if __name__ == '__main__':
    main()
