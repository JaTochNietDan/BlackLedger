// What is painted on the drums.
//
// The machine showed the core's own words — "7", "BAR", "CHERRY" — set in a
// row, which is a list of symbols rather than a machine: "when playing the slot
// machine we should show actual images for the stuff on the rollers."
//
// Drawn rather than photographed, the same way the city and the top bar are.
// The core still owns which symbol is on the payline and what it pays; this
// only decides what that symbol looks like, and anything it does not recognise
// falls back to the word, so a symbol added to the strip shows up as itself
// rather than as a blank drum.

const REELS: Record<string, string> = {
  // A red seven with the shoulder a cast seven has.
  seven:
    '<path d="M14 12h22l-11 30h-8l10-23H14z" fill="#c2422f"/>' +
    '<path d="M14 12h22v5H14z" fill="#8e2b1d"/>',
  // A gold bar, seen slightly from above so it reads as a bar and not a box.
  bar:
    '<path d="M8 22l6-5h22l6 5-6 5H14z" fill="#d9a441"/>' +
    '<path d="M8 22v9l6 5h22l6-5v-9l-6 5H14z" fill="#b6832c"/>' +
    '<text x="25" y="30" font-size="8" font-family="serif" text-anchor="middle" fill="#f4e3bd">BAR</text>',
  // A bell with a clapper under it.
  bell:
    '<path d="M25 9c-7 0-11 5-11 12 0 8-2 12-4 15h30c-2-3-4-7-4-15 0-7-4-12-11-12z" fill="#d9a441"/>' +
    '<path d="M25 9c-7 0-11 5-11 12 0 8-2 12-4 15h9c-1-3-2-7-2-15 0-7 2-12 8-12z" fill="#eec773"/>' +
    '<circle cx="25" cy="40" r="3" fill="#8e6a1f"/>',
  // A plum, dark and round, with a leaf.
  plum:
    '<circle cx="25" cy="29" r="13" fill="#6b3d7a"/>' +
    '<path d="M25 16c-3 4-4 9-4 13s1 9 4 13" fill="#57305f"/>' +
    '<path d="M25 16c0-4 3-6 6-7-1 4-3 6-6 7z" fill="#5d8a4a"/>',
  // An orange, with the pith showing where the stalk was.
  orange:
    '<circle cx="25" cy="29" r="13" fill="#d5762a"/>' +
    '<circle cx="25" cy="29" r="13" fill="none" stroke="#a9531a" stroke-width="2"/>' +
    '<circle cx="25" cy="18" r="2" fill="#8d4415"/>' +
    '<path d="M31 20c-2 4-2 14 0 18" fill="none" stroke="#e79a5b" stroke-width="2"/>',
  // A lemon, longer than it is round, with the point at each end.
  lemon:
    '<path d="M25 15c9 0 14 6 14 14s-5 14-14 14-14-6-14-14 5-14 14-14z" fill="#dcc02f"/>' +
    '<path d="M25 15c4 0 7 6 7 14s-3 14-7 14" fill="#e8d566"/>' +
    '<path d="M25 13c1 1 1 2 0 3-1-1-1-2 0-3zm0 30c1 1 1 2 0 3-1-1-1-2 0-3z" fill="#b39a1c"/>',
  // Two cherries on one stem, which is what pays on its own.
  cherry:
    '<path d="M25 12c-4 5-9 9-11 14m11-14c4 5 8 9 10 13" fill="none" stroke="#5d8a4a" stroke-width="2"/>' +
    '<circle cx="14" cy="33" r="8" fill="#b23026"/>' +
    '<circle cx="34" cy="33" r="7" fill="#8f2019"/>' +
    '<circle cx="11" cy="30" r="2" fill="#e07c6f"/>',
};

// reelArt is the drum's painting for a symbol, and nothing for one it has never
// heard of — the caller falls back to the word the core sent.
export function reelArt(id?: string): string {
  return (id && REELS[id]) || '';
}

// reelSymbols is every symbol this file can paint, for a guard to check the
// core's strip against.
export function reelSymbols(): string[] {
  return Object.keys(REELS);
}
