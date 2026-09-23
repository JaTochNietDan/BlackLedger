# Third-party software

Black Ledger uses React, React DOM, Three.js, PixiJS/pixi-viewport and their
transitive dependencies for the browser, and modernc.org/sqlite and its
transitive dependencies for local saves. The Go runtime is included in binaries.

These dependencies retain their upstream licenses. The project’s noncommercial
restriction does not replace or restrict rights granted by those licenses.
Exact versions are recorded in `package-lock.json`, `desktop/package-lock.json`
and `go.sum`. The desktop shell bundles Electron and Chromium; their upstream
license files and Chromium third-party notices are retained in the application
distribution. Desktop packaging tools are build-only dependencies.

The bundled `resources/game` (macOS: `Contents/Resources/game`) contains a
`licenses/` directory with copies of license/notice
files from installed npm packages and downloaded Go modules, a machine-readable
inventory, and the Go license. This intentionally includes build dependencies
as well as runtime dependencies so bundled notices are not silently omitted.
Missing upstream notices and declarations are tracked in docs/LOCAL_AI.md.

Desktop archives also contain Kokoro.js, Transformers.js, ONNX Runtime and their
runtime dependencies. Their notices and inventory are included in `licenses/`.
Ollama and model weights are fetched from their official sources during optional
first-run setup, under their own terms. Pinned versions and checksums are in
`desktop/ai/manifest.json`; their licenses are copied into the installed AI folder.
The Qwen3 and Kokoro weights retain Apache 2.0 terms; Ollama retains MIT terms.
The game's noncommercial restriction does not replace third-party licenses.

Phonemizer.js declares Apache 2.0 but embeds the GPL eSpeak NG engine; sharp also
includes LGPL components. Full corresponding-source/provenance review for these
binary dependencies remains a public binary-release gate in `docs/LOCAL_AI.md`.
Do not treat the npm top-level license strings as complete redistribution approval.
See ASSETS.md for media, which is distinct from software dependencies.

`@pixi/colord` omits the MIT license text from its npm archive. The upstream
notice is retained in `packaging/licenses/colord-LICENSE.md`, obtained from
https://github.com/omgovich/colord/blob/master/LICENSE.md on 2026-09-21.
Esbuild and Rollup platform binary packages use the notices shipped by their
respective parent packages.

`@napi-rs/lzma-linux-x64-gnu` 1.5.1 is an optional Rollup build dependency
present on Linux x64. Its npm archive declares MIT but contains no license text;
neither does its parent `@napi-rs/lzma` 1.5.1 archive or upstream source revision
`f164df92d83e095f195d628b1a68a141ae2eb638` at
https://github.com/Brooooooklyn/lzma/tree/f164df92d83e095f195d628b1a68a141ae2eb638
(checked 2026-09-22). The archive's package.json is retained verbatim as
`LICENSE-declaration.json`, alongside its README and inventory entry. This records
the upstream declaration; it is not a substitute for a missing license notice.
The native compressor is a build dependency, not a shipped game runtime module.
The packaging exception is limited to this version; review notice availability
when updating it.
