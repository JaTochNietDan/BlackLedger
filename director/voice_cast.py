"""Versioned, offline voice recipes. IDs must never be reassigned after shipping."""
import hashlib

GROUPS = (
    ('af_alloy','af_aoede','af_bella','af_heart','af_jessica','af_kore','af_nicole','af_nova','af_river','af_sarah'),
    ('am_echo','am_eric','am_fenrir','am_liam','am_michael','am_onyx','am_puck'),
    ('bf_alice','bf_emma','bf_isabella','bf_lily'),
    ('bm_daniel','bm_fable','bm_lewis'),
)
PROFILES = {}
for group in GROUPS:
    for i, voice in enumerate(group):
        for blend in (False, True):
            profile_id = f'cast2-{len(PROFILES):02d}'
            PROFILES[profile_id] = {
                'voices': (voice, group[(i+1) % len(group)]) if blend else (voice,),
                'weight': 0.7 if blend else 1.0,
                'speed': (0.96, 1.0, 1.04)[i % 3] if blend else 1.0,
                'language': voice[0],
            }


def assignment(speaker):
    if not isinstance(speaker, str) or not speaker.strip() or len(speaker) > 120:
        raise ValueError('invalid speaker')
    identity = ' '.join(speaker.lower().split())
    index = int.from_bytes(hashlib.sha256(identity.encode()).digest()[:4], 'big') % len(PROFILES)
    return f'cast2-{index:02d}'


def resolve(profile_id):
    if not isinstance(profile_id, str) or profile_id not in PROFILES:
        raise ValueError('unknown voice profile')
    return PROFILES[profile_id]
