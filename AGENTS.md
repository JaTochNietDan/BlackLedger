# Black Ledger collaboration

This is the independent browser mafia repository. Do not modify the old Unity/Godot Afterlight project for this task.

Read `docs/GOAL.md` for the user's current amended scope and `API.md` before changing interfaces.

- Codex owns Go gameplay, simulation, AI direction, saves, voice behavior, functional UX and campaign playtests. Report visual defects and integrate reviewed visual commits; visual production has been assigned to a separate agent.
- The visual agent works on `codex/visual-handoff` in `black-ledger-visuals`. Read `docs/VISUAL_HANDOFF.md` and use its preview/acceptance instructions. Rendering, assets, cosmetic animation and UI skin belong there.
- Coordinate `src/main.tsx` and `src/types.ts`; avoid parallel rewrites of shared command handling. Public API meaning stays stable unless an explicit agreed change is documented.
- Never mutate `.runtime/campaign.sqlite3` for QA. Use fresh isolated fixtures or an explicitly identified test save. Main game is port 8791; visual preview defaults to 8840.
- Commit coherent progress and record actual evidence/remaining work. User instructions take precedence over this ownership guidance.
