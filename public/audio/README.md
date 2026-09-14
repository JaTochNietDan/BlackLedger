# Game sound sources

The user supplied the three gunshot WAVs on September 13, 2026, explicitly for
use in Black Ledger. Copies are byte-for-byte originals; Downloads files are
untouched. No third-party sound library or license is claimed here.

| Game asset | Provided filename | Format | Duration | SHA-256 |
|---|---|---|---|---|
| guns/revolver.wav | revolver.wav | 48kHz stereo PCM16 | 1.243750s | 3f9beed919a410a25a63638a6b208427696dfb70df9f2fe8407a6a34be2954e0 |
| guns/shotgun.wav | shotgun.wav | 48kHz stereo PCM16 | 1.000000s | 718503005981fb15ee1daa5829f91935cbd054b45011b255bc2eb2a9d99c4319 |
| guns/thompson.wav | gunshot_thompson.wav | 48kHz stereo PCM16 | 1.047729s | 77f06115d264cccab707bc9496e4c19c5f9267b7c87a74d1b9a1fd5951764a2f |

All have immediate signal onset; no silence trimming or resampling was required.
Runtime gains are .75/.65/.8 respectively. A shared compressor moderates overlapping
burst peaks. Full sample tails overlap; Skip/mute/navigation cancel active voices.
The revolver and shotgun originals contain a few full-scale samples; headroom in
the mix prevents simply adding them at unity, but does not restore clipped source
data. Final subjective sound/mix review remains part of production acceptance.

## Additional effects — September 13

Ten additional WAVs supplied by the user were copied unchanged from Downloads.
`effects/sources.json` records original filenames, formats, durations and SHA-256.
Explosion and raid use their full recordings in 14s/10s scene windows. Arrest uses
siren; authored door contact uses door-kick; matched shooting casualties use pain;
explosions layer panic after one second. Nearby active building fires use crackle.
Visible nearby vehicles use an approach sample once per observed journey and idle
when stationary. Fire/idle have a 50ms decoded-buffer seam crossfade; source files
remain unchanged. Loops stop on mute, hidden tab, context loss or view disposal.
The drive-away recording is imported for the still-unfinished moving getaway
choreography; it is not played over stationary cars. Subjective mix review remains.

The user's `newspaper opening.wav` is imported unchanged as `effects/newspaper.wav`
(1.68s, stereo48kHz PCM16). It plays at gain0.5 on newspaper reveal and when opening
Herald, respecting the sound preference and cancellation on closing. The previous
procedural rustle remains a fallback if the recording has not loaded.
