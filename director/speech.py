"""Bounded local Kokoro narration and system NPC-voice rendering. Text is data, never a shell command."""
import hashlib
try:
    from .voice_cast import assignment, resolve
    from .voice_worker import VoiceWorker
except ImportError:
    from voice_cast import assignment, resolve
    from voice_worker import VoiceWorker
import re
import subprocess
import tempfile
import threading
from pathlib import Path

VOICES = {'narrator': 'Daniel', 'warm': 'Samantha (English (US))'}
CACHE = Path(__file__).resolve().parents[1] / '.runtime' / 'speech'
_render_lock = threading.Lock()
MAX_FILES = 128
ROOT = Path(__file__).resolve().parents[1]
VOICE_PYTHON = ROOT / ".tools" / "voice-venv" / "bin" / "python"
_worker = VoiceWorker([str(VOICE_PYTHON), str(ROOT / 'director/kokoro_voice.py'), '--worker'])
NARRATOR_VERSION = "kokoro-82m-bf16-bm_george-speed1-pcm16-v1"


def normalize(payload):
    text = payload.get('text')
    if not isinstance(text, str) or not text.strip() or len(text) > 1000:
        raise ValueError('speech text must be 1–1000 characters')
    # Apple's embedded speech commands use [[...]]. Strip them, including
    # unmatched brackets, so generated prose cannot select voices or rates.
    text = re.sub(r'\[\[.*?\]\]', '', text, flags=re.S)
    text = ' '.join(''.join(' ' if c.isspace() else c for c in text if (c.isprintable() or c.isspace()) and c not in '[]').split())
    if not text:
        raise ValueError('empty speech')
    voice = payload.get('voice', 'narrator')
    if not isinstance(voice, str) or voice not in VOICES:
        raise ValueError('unknown voice')
    return text, voice


def render(payload):
    text, voice = normalize(payload)
    profile_id = (payload.get('profile') or assignment(payload['speaker'])) if voice == 'warm' and payload.get('speaker') else None
    profile = resolve(profile_id) if profile_id else None
    character = profile['voices'][0] if profile else None
    backend = NARRATOR_VERSION if voice == 'narrator' else 'kokoro-cast-v2-' + profile_id if character else 'system-v1'
    key = hashlib.sha256((backend + '|' + voice + '|' + text).encode()).hexdigest()
    # Serial rendering bounds CPU/subprocess use; extra callers retry later.
    if not _render_lock.acquire(timeout=0.1):
        raise BlockingIOError('speech busy')
    try:
        CACHE.mkdir(parents=True, exist_ok=True)
        target = CACHE / (key + '.wav')
        if target.exists():
            target.touch()
            return target.read_bytes()
        with tempfile.TemporaryDirectory(prefix='render-', dir=CACHE) as folder:
            aiff = Path(folder) / 'voice.aiff'
            wav = Path(folder) / 'voice.wav'
            if voice == 'narrator' or character:
                _worker.render(text, profile_id, wav)
            else:
                subprocess.run(['/usr/bin/say', '-v', VOICES[voice], '-r', '165',
                                '-o', str(aiff)], input=text.encode(),
                               stdout=subprocess.DEVNULL, stderr=subprocess.PIPE,
                               timeout=15, check=True)
                subprocess.run(['/usr/bin/afconvert', '-f', 'WAVE', '-d', 'LEI16',
                                str(aiff), str(wav)], stdout=subprocess.DEVNULL,
                               stderr=subprocess.PIPE, timeout=5, check=True)
            data = wav.read_bytes()
            if not data.startswith(b'RIFF') or len(data) > 8_000_000:
                raise OSError('invalid rendered audio')
            wav.replace(target)
        for old in sorted(CACHE.glob('*.wav'), key=lambda p: p.stat().st_mtime)[:-MAX_FILES]:
            old.unlink(missing_ok=True)
        return data
    finally:
        _render_lock.release()
