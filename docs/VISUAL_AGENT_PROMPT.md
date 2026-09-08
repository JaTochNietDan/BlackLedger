You are taking ownership of the visual presentation of Black Ledger, a single-player browser mafia game. Work in the `black-ledger-visuals` Git worktree on `codex/visual-handoff`.

Read `docs/VISUAL_HANDOFF.md`, `docs/VISUAL_ACCEPTANCE.md`, `docs/ART_DIRECTION.md`, `API.md`, and `src/types.ts`. Inspect the actual reference images and running preview before editing. Start it with `./scripts/run-visual-preview.sh damage` on port 8840. The preview uses its own staged save; never use or modify the user's campaign on port 8791.

Improve the painted-noir isometric city and its UI. Prioritize visible compositing artifacts, occlusion, street grounding, pedestrian quality, lighting and consistent presentation. Then add missing district/property/portrait art. Work autonomously within that visual scope and assess your own screenshots and motion. Keep gameplay authoritative in Go; do not invent client-side rules or change the API while reskinning.

The other agent owns gameplay, AI storytelling, economy, saves, functional UX, headless simulations and campaign playtests. Coordinate changes to `src/main.tsx` and `src/types.ts`; prefer extracting presentational components. Record behavioral bugs or API needs for that agent.

Commit coherent progress on your branch. Maintain `docs/VISUAL_DELIVERY.md` with before/after evidence, actual viewport sizes, asset provenance, tests and remaining issues. Do not merge or publish; provide commit hashes for integration. Do not stop at a style mockup: deliver working, inspected presentation through the real game API.
