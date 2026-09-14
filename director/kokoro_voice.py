"""Isolated offline narrator worker; invoked by speech.py, never plays audio."""
import sys
import os
import json
from contextlib import redirect_stdout
from pathlib import Path

MODEL = 'mlx-community/Kokoro-82M-bf16'
from voice_cast import resolve
DEFAULT_PROFILE = {'voices': ('bm_george',), 'weight': 1.0, 'speed': 1.0, 'language': 'b'}
_model = None


def render(text, profile_id, output):
    global _model
    profile = resolve(profile_id) if profile_id else DEFAULT_PROFILE
    import numpy as np
    import soundfile as sf
    from mlx_audio.tts.utils import load_model
    if not text.strip() or len(text) > 1000:
        raise ValueError('invalid narration length')
    from huggingface_hub import snapshot_download
    import mlx.core as mx
    from mlx_audio.tts.models.kokoro.voice import load_voice_tensor
    paths = []
    for voice in profile['voices']:
        voice_root = Path(snapshot_download(MODEL, allow_patterns=[f'voices/{voice}.safetensors'], local_files_only=True))
        path = voice_root / 'voices' / f'{voice}.safetensors'
        if not path.is_file():
            raise FileNotFoundError('local character voice missing')
        paths.append(path)
    voice_path = paths[0]
    if len(paths) == 2:
        weight = profile['weight']
        blend = load_voice_tensor(str(paths[0])) * weight + load_voice_tensor(str(paths[1])) * (1 - weight)
        voice_path = Path(output).with_suffix('.safetensors')
        mx.save_safetensors(str(voice_path), {'voice': blend})
    if _model is None:
        _model = load_model(MODEL)
    chunks = []
    rate = None
    for result in _model.generate(text, voice=str(voice_path), lang_code=profile['language'], speed=profile['speed']):
        if rate is not None and rate != result.sample_rate:
            raise ValueError('inconsistent sample rate')
        rate = result.sample_rate
        chunks.append(np.asarray(result.audio, dtype=np.float32))
    if not chunks:
        raise ValueError('empty narration')
    audio = np.concatenate(chunks)
    if not np.isfinite(audio).all() or len(audio) > rate * 150:
        raise ValueError('invalid narration audio')
    sf.write(Path(output), audio, rate, subtype='PCM_16')


def main():
    if '--worker' not in sys.argv:
        render(sys.stdin.read(4001), os.environ.get('AFTERLIGHT_KOKORO_PROFILE', ''), sys.argv[1])
        return
    # Only protocol replies go to stdout; third-party model logs go to stderr.
    for line in sys.stdin:
        request = json.loads(line)
        try:
            with redirect_stdout(sys.stderr):
                render(request['text'], request.get('profile'), request['output'])
            response = {'id': request['id'], 'ok': True}
        except Exception:
            response = {'id': request.get('id'), 'ok': False}
        print(json.dumps(response), flush=True)


if __name__ == '__main__':
    main()
