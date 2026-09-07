# Simulation, director and presentation

The Go core owns truth. A command authorizes a bounded amount of game time. The core advances an event schedule, stops at any player intervention, commits results, and returns public state plus presentation records. React/Pixi can take seconds to illustrate those records or skip immediately. Backend time never depends on animation duration, display frame rate, or whether the page is visible.

The core can prepare next options while presentation plays. AI may prepare candidate situations in parallel. Neither may commit future player choices speculatively. A revision and incident ID anchor every response; late proposals for a different life/world are discarded. Private plans are filtered out of the public projection.

Boundary: React components (input/accessibility/layout) + Pixi scene (art/animation) → versioned command API → Go core → SQLite transactional state and receipts. Model and voice services are replaceable HTTP providers. Headless commands exercise the same rules as the frontend.

Performance policy: prioritize correct event scheduling and bounded work. Do not simulate decorative pedestrians as autonomous NPCs. The reference minute-step implementation is bounded by current short actions; replace with a scheduled-event clock before large-city scaling and benchmark the result. Do not claim idle animations save CPU if the renderer is needlessly repainting continuously.

Repository policy: this folder is an independent Git repository. Commit coherent changes with validation and maintain a playtest record. Do not place purchased Afterlight assets or private runtime saves in Git.
