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
import {BLOCK, PAVE, ROAD, addressSlot, awnings, goldenness, island, nightness, rails, sleepers, terrace, trolleyAvenue, TROLLEY_GAUGE, vents} from '../.runtime/frontend-test/iso.js';

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

// The awnings hang over the pavement and never over the road, and never over
// each other. Same rule the props follow, checked the same way.
test('an awning hangs over its own pavement and no further', () => {
  let found = 0;
  for (const cell of cells) {
    const i = island(cell);
    for (const a of awnings(cell, SLOTS)) {
      found++;
      // Out from the wall, it stops inside the kerb.
      assert.ok(a.at.y + a.reach <= i.y + i.d + 1e-9,
        `awning on ${JSON.stringify(cell)} reaches the road`);
      // Along the wall, it stays inside its own shopfront.
      const slot = terrace(cell, SLOTS).filter(s => s.front)
        .find(s => a.at.x >= s.at.x - 1e-9 && a.at.x + a.w <= s.at.x + s.w + 1e-9);
      assert.ok(slot, `awning on ${JSON.stringify(cell)} spans more than one building`);
      // And it hangs off the wall it belongs to rather than floating.
      assert.ok(Math.abs(a.at.y - (slot.at.y + slot.d)) < 1e-9, 'awning is not against its wall');
      assert.ok(a.h > a.drop, 'awning front edge is below the pavement');
    }
  }
  assert.ok(found > 8, `only ${found} awnings in the whole city`);
});

test('two awnings never overlap', () => {
  const all = cells.flatMap(c => awnings(c, SLOTS));
  for (let i = 0; i < all.length; i++)
    for (let j = i + 1; j < all.length; j++) {
      const a = all[i], b = all[j];
      const x = Math.min(a.at.x + a.w, b.at.x + b.w) - Math.max(a.at.x, b.at.x);
      const y = Math.min(a.at.y + a.reach, b.at.y + b.reach) - Math.max(a.at.y, b.at.y);
      assert.ok(!(x > 1e-9 && y > 1e-9), 'two awnings occupy the same air');
    }
});

// The light has a temperature as well as a level. Without this, dawn is only a
// weaker night: the same blue-black ground, a little lighter.
test('the light is warm at the turns of the day and nowhere else', () => {
  const at = (h, m = 0) => h * 60 + m;
  assert.equal(goldenness(at(12)), 0, 'midday is not golden hour');
  assert.equal(goldenness(at(3)), 0, 'three in the morning is not golden hour');
  assert.ok(goldenness(at(6, 12)) > .95, 'dawn is not warm');
  assert.ok(goldenness(at(19, 12)) > .95, 'dusk is not warm');
  // And it is the opposite of a brightness curve: dawn and the small hours are
  // both dark, and only one of them is gold.
  assert.ok(nightness(at(3)) === 1 && nightness(at(5, 30)) > 0, 'both are night');
  assert.ok(goldenness(at(5, 30)) > goldenness(at(3)), 'dawn is no warmer than 3am');
  // It wraps with the clock rather than running off the end of a day.
  assert.equal(goldenness(at(6, 12)), goldenness(at(6, 12) + 1440 * 9));
});

// What a block gives off stands on the block, not in the traffic.
test('smoke comes off a roof and steam off a pavement', () => {
  let chimneys = 0, grates = 0;
  for (const cell of cells) {
    const i = island(cell);
    const found = vents(cell, SLOTS);
    assert.ok(found.length <= 1, 'a block gives off one thing at most');
    for (const v of found) {
      assert.ok(v.at.x >= i.x && v.at.x <= i.x + i.w && v.at.y >= i.y && v.at.y <= i.y + i.d,
        'a vent is off its own block');
      if (v.kind === 'chimney') {
        chimneys++;
        // A chimney stands on a building rather than in the yard behind it.
        const on = terrace(cell, SLOTS).some(s =>
          v.at.x >= s.at.x && v.at.x <= s.at.x + s.w && v.at.y >= s.at.y && v.at.y <= s.at.y + s.d);
        assert.ok(on, 'a chimney is not standing on a roof');
        assert.ok(v.height > 1, 'smoke starts below the rooftops');
      } else {
        grates++;
        // A grate is in the pavement: on the block, outside every building.
        const inside = terrace(cell, SLOTS).some(s =>
          v.at.x >= s.at.x && v.at.x <= s.at.x + s.w && v.at.y >= s.at.y && v.at.y <= s.at.y + s.d);
        assert.ok(!inside, 'a grate is under a building');
        assert.equal(v.height, 0, 'steam starts above the pavement');
      }
    }
  }
  assert.ok(chimneys > 5, `only ${chimneys} chimneys in the whole city`);
  assert.ok(grates > 0, 'no steam anywhere');
  // And most of the city gives off nothing: every roof smoking is a foundry.
  assert.ok(chimneys + grates < cells.length * .7, 'the whole city is smoking');
});

// The rails are derived from the same grid as the roads, so they cannot end up
// half on the pavement. This is what "by construction" has to mean to be worth
// anything: the test states the property, the geometry makes it unavoidable.
test('the trolley runs down the middle of a street, not over the kerb', () => {
  const size = {cols: 6, rows: 4};
  const pair = rails(size);
  assert.equal(pair.length, 2, 'a track has two rails');
  const centre = trolleyAvenue(size) * BLOCK;
  for (const r of pair) {
    assert.equal(r.a.x, r.b.x, 'a rail wanders off its street');
    // Inside the carriageway, which is ROAD wide centred on the grid line.
    assert.ok(Math.abs(r.a.x - centre) <= ROAD / 2 - .05, 'a rail is on the pavement');
    // And it runs the whole length of the city rather than stopping in the
    // middle of nowhere, the same rule the carriageways follow.
    assert.ok(r.a.y <= 0 && r.b.y >= size.rows * BLOCK, 'the track stops in mid-air');
  }
  assert.ok(Math.abs((pair[1].a.x - pair[0].a.x) - TROLLEY_GAUGE) < 1e-9, 'the gauge is wrong');
  // The ties stay between the rails.
  for (const s of sleepers(size)) {
    assert.ok(Math.abs(s.a.x - centre) < TROLLEY_GAUGE, 'a sleeper sticks out past the rails');
    assert.equal(s.a.y, s.b.y, 'a sleeper is not square to the track');
  }
  // And no building stands on the track: the avenue it takes is a road.
  for (const cell of cells) {
    const i = island(cell);
    assert.ok(centre <= i.x || centre >= i.x + i.w,
      `the trolley runs through block ${JSON.stringify(cell)}`);
  }
});
