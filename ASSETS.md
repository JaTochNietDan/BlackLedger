# Asset provenance and publication status

The project code license does not establish ownership of imported media.
Third-party materials retain their original terms. No AI model weights,
voice models, purchased source packs, or local service installations are bundled.

| Material | Evidence | Publication status |
| --- | --- | --- |
| Locally authored Blender geometry / procedural textures | `tools/`, `public/art/models/README.md` | Original project material; PolyForm Noncommercial 1.0.0 applies to rights held by contributors. |
| Generated portraits, scenes and building images | Local generation tools, `art/textures/README.md`, historical development notes | Inventory/provenance review required before first public push; generated output is not a warranty of exclusive rights. |
| Supplied WAV recordings | `public/audio/README.md`, `public/audio/effects/sources.json` | Public redistribution permission is pending owner confirmation. Permission to use in the game alone does not establish repository redistribution rights. |
| Ground textures and older reference images | `public/art/ground/README.md`, `art/sources/README.md`, `docs/art/` | File-level origin review remains required. |
| npm and Go libraries | Lockfiles and packaged `licenses/` directory | Their own licenses apply; not relicensed under the project license. |

Do not publish the repository history or a game archive until unresolved media
rights are resolved. Removing a file from the current tree does not remove it
from Git history. Record a source URL, author, license, and any required credit
for imported replacements. Do not import assets from the separate Afterlight project.
