"""Isolated offline narrator worker; invoked by speech.py, never plays audio."""
import sys
import os
from pathlib import Path

MODEL = 'mlx-community/Kokoro-82M-bf16'
from voice_cast import resolve
PROFILE_ID = os.environ.get('AFTERLIGHT_KOKORO_PROFILE', '')
PROFILE = resolve(PROFILE_ID) if PROFILE_ID else {'voices': ('bm_george',), 'weight': 1.0, 'speed': 1.0, 'language': 'b'}


def main():
    import numpy as np
    import soundfile as sf
    from mlx_audio.tts.utils import load_model
    text = sys.stdin.read(4001)
    if not text.strip() or len(text) > 1000:
        raise ValueError('invalid narration length')
    from huggingface_hub import snapshot_download
    import mlx.core as mx
    from mlx_audio.tts.models.kokoro.voice import load_voice_tensor
    paths = []
    for voice in PROFILE['voices']:
        voice_root = Path(snapshot_download(MODEL, allow_patterns=[f'voices/{voice}.safetensors'], local_files_only=True))
        path = voice_root / 'voices' / f'{voice}.safetensors'
        if not path.is_file():
            raise FileNotFoundError('local character voice missing')
        paths.append(path)
    voice_path = paths[0]
    if len(paths) == 2:
        weight = PROFILE['weight']
        blend = load_voice_tensor(str(paths[0])) * weight + load_voice_tensor(str(paths[1])) * (1 - weight)
        voice_path = Path(sys.argv[1]).with_suffix('.safetensors')
        mx.save_safetensors(str(voice_path), {'voice': blend})
    model = load_model(MODEL)
    chunks = []
    rate = None
    for result in model.generate(text, voice=str(voice_path), lang_code=PROFILE['language'], speed=PROFILE['speed']):
        if rate is not None and rate != result.sample_rate:
            raise ValueError('inconsistent sample rate')
        rate = result.sample_rate
        chunks.append(np.asarray(result.audio, dtype=np.float32))
    if not chunks:
        raise ValueError('empty narration')
    audio = np.concatenate(chunks)
    if not np.isfinite(audio).all() or len(audio) > rate * 150:
        raise ValueError('invalid narration audio')
    sf.write(Path(sys.argv[1]), audio, rate, subtype='PCM_16')


if __name__ == '__main__':
    main()
