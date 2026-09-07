# Simulation, director and presentation

The Go core owns truth. A command authorizes a bounded amount of game time. The core advances an event schedule, stops at any player intervention, commits results, and returns public state plus presentation records. React/canvas can take seconds to illustrate those records or skip immediately. Backend time never depends on animation duration, display frame rate, or whether the page is visible.

The core can prepare next options while presentation plays. AI may prepare candidate situations in parallel. Neither may commit future player choices speculatively. A revision and incident ID anchor every response; late proposals for a different life/world are discarded. Private plans are filtered out of the public projection.

Boundary: React components (input/accessibility/layout) + isolated canvas street (art/animation) → versioned command API → Go core → SQLite transactional state and receipts. Model and voice services are replaceable HTTP providers. Headless commands exercise the same rules as the frontend.

Performance policy: prioritize correct event scheduling and bounded work. Do not simulate decorative pedestrians as autonomous NPCs. The Go clock jumps to scheduled boundaries (task completion, attacks, warning times, midnight bills, or the authorized action end), integrating passive income over each interval. It does not loop over every game minute. Do not claim idle animations save CPU if the renderer is needlessly repainting continuously.

Repository policy: this folder is an independent Git repository. Commit coherent changes with validation and maintain a playtest record. Do not place purchased Afterlight assets or private runtime saves in Git.

## Prepared time segments

Five game minutes are not five minutes of backend activity or five minutes of mandatory playback. A command may authorize a five-minute journey. The core resolves only as far as the next intervention, saves the result transactionally, and returns public records. The frontend may illustrate that committed interval briefly or skip it. It cannot award income, resolve combat, or advance the authoritative clock.

If an encounter interrupts that journey after two minutes, the remaining three minutes are not committed. Wait for a fresh player command before resolving them. While a committed segment is displayed, the director may prepare candidate scenes or voice audio using the resulting revision. These are speculative until validated against current facts. Discard or revalidate stale work after a new action; never expose hidden future plans to the browser.

The backend need not tick during playback. It still handles save/API requests and optional generation. Future interactive controls must use committed state or explicit intervention boundaries rather than attempting to undo events already resolved. Playback duration has no effect on game outcomes.
