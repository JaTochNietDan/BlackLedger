#!/usr/bin/env python3
"""Exercise a native release from an unrelated directory with isolated saves."""
import argparse
from contextlib import contextmanager
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
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
        command = [str(executable), '-desktop', '-browser=false', '-addr', ':0']
        process = subprocess.Popen(command, cwd=root, env=env, stdout=log, stderr=log)
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
    parser.add_argument('--with-ai-voice', action='store_true', help='download the verified speech model and test real synthesis')
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
            parts = file.relative_to(game).parts
            assert not any(part in {'.runtime', '.tools', '.git'} for part in parts)
            if 'node_modules' in parts:
                assert 'app.asar.unpacked' in parts, 'unexpected development dependencies were packaged'
            assert '.sqlite' not in file.name, 'a save was packaged'
        resources = next(game.rglob('game/licenses/inventory.json')).parent.parent
        assert (resources / 'licenses/inventory.json').is_file()
        assert 'Electron' in (resources / 'licenses/Electron-LICENSE').read_text()
        assert (resources / 'licenses/Chromium-LICENSES.html').stat().st_size > 1000
        assert (game / 'LICENSE').is_file() and (game / 'NOTICE').is_file()
        executable = resources / ('blackledger.exe' if os.name == 'nt' else 'blackledger')
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
        # Follow PLAY.txt's stopped-game backup instructions, then restore into
        # a separate save under a replacement extracted game directory. Include
        # WAL sidecars: forced process termination can leave committed data there.
        backup = root / 'restored user config' / 'campaign.sqlite3'
        backup.parent.mkdir()
        original = {}
        for suffix in ('', '-wal', '-shm'):
            source = Path(str(db) + suffix)
            if source.exists():
                original[suffix] = source.read_bytes()
                shutil.copyfile(source, Path(str(backup) + suffix))
        replacement = root / 'replacement game folder'
        relative_executable = executable.relative_to(game)
        game.rename(replacement)
        with server(replacement / relative_executable, root,
                    {**env, 'BLACK_LEDGER_DB': str(backup)}, 'restore.log') as base:
            restored = json.loads(get(base + '/api/state'))
            assert restored['revision'] == committed['revision']
            assert restored['player'] == committed['player']
            assert action(base) == committed, 'restored backup lost the saved receipt'
            command = {**command, 'revision': restored['revision'],
                       'request_id': 'release-smoke-restored-wait'}
            assert action(base)['revision'] == restored['revision'] + 1
        for suffix, data in original.items():
            assert Path(str(db) + suffix).read_bytes() == data, 'restored game touched original save'
        if sys.platform == 'darwin':
            application = replacement / 'Black Ledger.app/Contents/MacOS/BlackLedger'
        else:
            application = replacement / ('BlackLedger.exe' if os.name == 'nt' else 'BlackLedger')
        desktop_save = root / 'desktop smoke'
        desktop_env = {**env, 'BLACK_LEDGER_SMOKE_DIR': str(desktop_save)}
        if args.with_ai_voice:
            desktop_env['BLACK_LEDGER_SMOKE_VOICE'] = '1'
        desktop_env.pop('ELECTRON_RUN_AS_NODE', None)
        for attempt in range(2):
            launched = subprocess.run([str(application), '--smoke-test'], cwd=root, env=desktop_env,
                                      capture_output=True, timeout=360 if args.with_ai_voice else 75)
            assert launched.returncode == 0, launched.stderr.decode(errors='replace')
            assert not (desktop_save / 'error.txt').exists(), (desktop_save / 'error.txt').read_text() if (desktop_save / 'error.txt').exists() else ''
            ready = json.loads((desktop_save / 'ready.json').read_text())
            assert ready['url'].rstrip('/') == ready['origin']
            assert (desktop_save / 'campaign.sqlite3').is_file()
            try:
                get(ready['origin'] + '/api/health')
            except (OSError, urllib.error.URLError):
                pass
            else:
                raise AssertionError('closing the desktop window left the server running')
            if attempt == 0:
                first = ready
            else:
                assert ready['player'] == first['player'] and ready['revision'] == first['revision']
        if args.with_ai_voice:
            assert (desktop_save / 'ai-check/smoke-voice.wav').stat().st_size > 44
            print('PASS: real speech synthesis from packaged runtime and verified downloaded model')
        print('PASS: desktop window, close/relaunch, checksum, archive, unrelated directory, frontend, per-user save, port conflict, command/retry, restart/retry, backup restore, replacement game folder')


if __name__ == '__main__':
    main()
