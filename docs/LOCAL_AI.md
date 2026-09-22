# Local AI in desktop distributions

First launch offers Download and enable AI or Play without AI. Game → AI setup
reopens setup. No separate commands, Node/Python install, API keys or accounts
are required. Skipping is remembered. Download failures show Retry; cancellation
keeps partial files. Existing campaigns remain separate from setup.

The app installs Ollama 0.34.2 and the current Qwen3 14B model, plus Kokoro 82M
ONNX q8 speech. The story weight blob is 9,276,184,896 bytes; speech weights are
92,361,116 bytes. Runtime downloads vary by platform (approximately 159 MB on
Mac and 1.5 GB on Windows/Linux). Allow 10–12 GB of network transfer plus extra
installation space (about 20 GB free is a conservative starting point). Memory needs are substantial; older/low-memory systems may
not generate within the game's deadlines. Authored encounters remain playable.
No minimum hardware performance is certified by the current Mac test.

AI files live at `<default save directory>/ai`, independent of the app archive
and any explicit campaign override. `BLACK_LEDGER_AI_HOME` selects an isolated
setup directory for development/QA. Never point QA at a player's installation.
Model/runtime URLs, sizes, SHA-256 checksums and immutable voice model revision
are recorded in `desktop/ai/manifest.json`. Downloads are streamed, resume with
HTTP Range, and must verify before installation. Setup reserves extraction space.
A matching installed manifest and file sizes skip downloading on later launches;
reopening setup and installing again also rechecks full hashes.

Services bind only to loopback on free ports. Go receives the managed AI proxy
address for director and speech. The proxy permits only the installed chat model
and bounded speech requests, and rejects browser Origin headers. External URLs are not opened by the game window.
The AI host owns Ollama and shuts its process tree down when its parent pipe
closes; quitting stops Go first, then AI. Existing developer services are untouched.
On Windows, built-in tar.exe extracts the verified Ollama ZIP before any native
speech addon loads. The AI host adds the archive's MSVC DLL folder to its private
process search path, so it does not install a system-wide redistributable. The
required DLLs were confirmed in the pinned official Windows archive; clean-machine
Windows acceptance is still required (CI runners already have development tools).
No cloud API is used for inference. Network access is required for initial
GitHub/Ollama/Hugging Face downloads. Installed voice inference disables remote
model lookup. Speech keeps the 48 cast2 profiles, blend weights and speaking
rates, with bm_george for narration. Different inference implementations and
quantization mean waveforms will not exactly match the old Python service.

## Validation and limits

Recorded September 21, 2026 on Apple Silicon:

- Actual verified Ollama archive and Kokoro weights downloaded and installed into
  `.runtime/managed-ai-qa`. Existing matching Qwen blobs were reused via isolated
  hardlinks and fully SHA-256 checked, avoiding another 9.3 GB transfer. This
  verifies existing-blob reuse, not a complete fresh large-blob network download.
- The managed host produced narrator and cast2-01 WAVs: 6.175 and 5.775 seconds
  of 24 kHz mono PCM, about two seconds generation each for the sample sentence.
- An isolated Go game requested a real Qwen encounter. The first proposal failed
  venue validation, the correction passed, and public director state became
  `ready`. `.runtime/managed-ai-game.log` records this; gameplay authority stayed
  in Go. This is a targeted test, not full narrative-quality acceptance.
- Unit tests cover download resume, ignored Range, corrupted payload rejection,
  cancellation, voice mapping/PCM, readiness and shutdown. Native archive smoke
  loads phonemization and ONNX dependencies, opens the desktop window twice,
  verifies saves/retries/backup restoration and checks server shutdown.
- Packaged Mac first-run UI completed installation and opened the isolated
  campaign. Skip/recovery opened the city without AI. The packaged director
  produced another validated encounter, and Go `/api/speech/prepare` returned
  `ready:true` for its character voice in 9.92 seconds. Normal app quit stopped
  the Go server, AI host, Ollama and loaded model runner; all three listening
  ports closed. `.runtime/packaged-ai-integration.log` records integration.
- Each of the five CI targets runs the dependency and desktop smoke checks,
  downloads the verified 92 MB speech model and synthesizes a real WAV.
  CI does not download the 9.3 GB story model. Full first-run AI generation on Windows,
  Intel macOS and Linux still needs native acceptance; Mac success does not prove it.

## Dependency publication review

Upstream model/runtime license texts and downloaded notice sources are retained
in `desktop/ai/licenses`. Desktop dependency notices accompany the distribution;
`licenses/inventory.json` records exact versions. Original game licensing does
not restrict rights separately granted by third-party components.

Before publishing binaries, finish corresponding-source/provenance review for
Phonemizer.js's embedded eSpeak NG (GPL) and sharp's libvips dependencies (LGPL),
including the appropriate source distribution and any necessary separation of
components. Do not infer all embedded binary terms from npm's Apache/MIT labels.
`guid-typescript` 1.0.9 declares ISC in its published package but omits a license
text, including in its upstream repository; its declaration is retained rather
than inventing a copyright holder. Resolve this missing notice or remove that
transitive dependency before public binary release. These gates are additional
to existing asset rights, signing and gameplay acceptance in RELEASING.md.
