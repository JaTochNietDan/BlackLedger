// The city's buildings must not stand inside each other.
//
// This exists because they did. A 1.16x overshoot was added to the sprite
// scaling to close the party walls between neighbours, and what it actually
// did was drive every building into the one next to it: roofs cutting through
// roofs, walls hanging over the kerb. It was only ever eyeballed at one zoom,
// and it was wrong. Geometry is checkable, so it is checked here rather than
// looked at.

import test from 'node:test';
import assert from 'node:assert/strict';
import {BLOCK, PAVE, ROAD, addressSlot, island, terrace} from '../.runtime/frontend-test/iso.js';

const SLOTS = 2;                 // must match CityIso.tsx

// Bigger than the city the twelve addresses actually make (six by four), so
// this holds however the addresses move. Whether two buildings collide is a
// property of the blocks themselves, not of how many of them there are.
const cells = [];
for (let col = 0; col < 8; col++)
  for (let row = 0; row < 6; row++) cells.push({col, row});

// A slot's ground, in plan. This is exactly the rectangle the sprite is scaled
// to fill, so two slots overlapping means two buildings overlapping.
const foot = s => ({x0: s.at.x, y0: s.at.y, x1: s.at.x + s.w, y1: s.at.y + s.d});
const over = (a, b) => Math.min(a.x1, b.x1) - Math.max(a.x0, b.x0) > 1e-9
                    && Math.min(a.y1, b.y1) - Math.max(a.y0, b.y0) > 1e-9;

test('no building stands in another building', () => {
  const all = cells.flatMap(c => terrace(c, SLOTS).map(s => ({cell: c, foot: foot(s)})));
  for (let i = 0; i < all.length; i++)
    for (let j = i + 1; j < all.length; j++)
      assert.ok(!over(all[i].foot, all[j].foot),
        `${JSON.stringify(all[i])} overlaps ${JSON.stringify(all[j])}`);
});

test('no building stands in the road or on the pavement', () => {
  for (const cell of cells) {
    const i = island(cell);
    // The buildable ground is the island less its pavement ring; the road is
    // outside that again.
    const inner = {x0: i.x + PAVE, y0: i.y + PAVE, x1: i.x + i.w - PAVE, y1: i.y + i.d - PAVE};
    for (const s of terrace(cell, SLOTS)) {
      const f = foot(s);
      assert.ok(f.x0 >= inner.x0 - 1e-9 && f.x1 <= inner.x1 + 1e-9
             && f.y0 >= inner.y0 - 1e-9 && f.y1 <= inner.y1 + 1e-9,
        `${JSON.stringify(cell)} slot ${JSON.stringify(f)} leaves its pavement ${JSON.stringify(inner)}`);
    }
  }
});

test('the front and back of a terrace are two rows, not one', () => {
  const [front] = terrace({col: 0, row: 0}, SLOTS).filter(s => s.front);
  const [back] = terrace({col: 0, row: 0}, SLOTS).filter(s => !s.front);
  assert.ok(back.at.y + back.d < front.at.y, 'back row runs into the front row');
});

test('the address takes a frontage slot that exists', () => {
  const n = addressSlot(SLOTS);
  const slots = terrace({col: 1, row: 1}, SLOTS);
  assert.ok(n >= 0 && n < slots.length, 'the address has no slot to stand in');
  assert.ok(slots[n].front, 'the address stands at the back of its own block');
});

test('the geometry the city is drawn from is the geometry tested here', () => {
  assert.equal(BLOCK, 3.6);
  assert.equal(ROAD, 1);
  assert.equal(PAVE, .42);
});
