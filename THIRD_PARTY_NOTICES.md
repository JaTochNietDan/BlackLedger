# Third-party software

Black Ledger uses React, React DOM, Three.js, PixiJS/pixi-viewport and their
transitive dependencies for the browser, and modernc.org/sqlite and its
transitive dependencies for local saves. The Go runtime is included in binaries.

These dependencies retain their upstream licenses. The project’s noncommercial
restriction does not replace or restrict rights granted by those licenses.
Exact versions are recorded in `package-lock.json` and `go.sum`.

Release archives contain a `licenses/` directory with copies of license/notice
files from installed npm packages and downloaded Go modules, a machine-readable
inventory, and the Go license. This intentionally includes build dependencies
as well as runtime dependencies so bundled notices are not silently omitted.
The packaging command fails if a dependency's notice cannot be found.

Optional Ollama, model weights, Python voice services and their dependencies are
not included in release archives. Install them separately under their own terms.
See ASSETS.md for media, which is distinct from software dependencies.

`@pixi/colord` omits the MIT license text from its npm archive. The upstream
notice is retained in `packaging/licenses/colord-LICENSE.md`, obtained from
https://github.com/omgovich/colord/blob/master/LICENSE.md on 2026-09-21.
Esbuild and Rollup platform binary packages use the notices shipped by their
respective parent packages.
