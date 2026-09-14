"""One serial, bounded local synthesis process; caller owns render serialization."""
import atexit
import json
import os
import selectors
import subprocess
import time


class VoiceWorker:
    def __init__(self, command, timeout=18):
        self.command = command
        self.timeout = timeout
        self.process = None
        self.sequence = 0
        atexit.register(self.close)

    def close(self):
        process, self.process = self.process, None
        if process is None:
            return
        if process.poll() is None:
            process.terminate()
            try:
                process.wait(timeout=1)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait()
        process.stdin.close()
        process.stdout.close()

    def render(self, text, profile, output):
        if self.process is None or self.process.poll() is not None:
            self.close()
            env = os.environ.copy()
            env.update(HF_HUB_OFFLINE='1', HF_HUB_DISABLE_TELEMETRY='1')
            self.process = subprocess.Popen(self.command, stdin=subprocess.PIPE,
                                            stdout=subprocess.PIPE, stderr=subprocess.DEVNULL,
                                            env=env, bufsize=0)
        process = self.process
        self.sequence += 1
        request = {'id': self.sequence, 'text': text, 'profile': profile, 'output': str(output)}
        try:
            process.stdin.write((json.dumps(request) + '\n').encode())
            process.stdin.flush()
            deadline = time.monotonic() + self.timeout
            line = bytearray()
            with selectors.DefaultSelector() as selector:
                selector.register(process.stdout, selectors.EVENT_READ)
                while b'\n' not in line:
                    remaining = deadline - time.monotonic()
                    if remaining <= 0 or not selector.select(remaining):
                        raise subprocess.TimeoutExpired(self.command, self.timeout)
                    chunk = os.read(process.stdout.fileno(), 4096)
                    if not chunk or len(line) + len(chunk) > 4096:
                        raise OSError('voice worker stopped or invalid response')
                    line.extend(chunk)
            response = json.loads(line)
            if response.get('id') != self.sequence or response.get('ok') is not True:
                raise OSError('voice worker failed')
        except (OSError, ValueError, subprocess.SubprocessError):
            # A timed-out request must never complete into the next request's file.
            self.close()
            raise
