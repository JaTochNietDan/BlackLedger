// Every address in this city must stand on a block of its own.
//
// "We don't have the map working properly." grid() places an address on the
// block nearest its own coordinates and steps outwards when that block is
// taken — but it only tried the four cardinal neighbours at each ring, gave up
// after five rings, and then wrote itself into the taken map anyway. So two
// addresses could end up in the same block and be drawn standing inside each
// other, which is not something you can see from one screenshot at one zoom.
//
// The city has twenty-five addresses now. It had twelve when this was written.

import test from 'node:test';
import assert from 'node:assert/strict';
import {readFileSync} from 'node:fs';
import {blockFor, grid} from '../.runtime/frontend-test/iso.js';

const places = JSON.parse(readFileSync(new URL('../core/locations.json', import.meta.url), 'utf8'))
  .map(p => ({id: p.id, x: p.x, y: p.y}));

test('the city has addresses to place', () => {
  assert.ok(places.length >= 20, `only ${places.length} addresses`);
});

test('no two addresses share a block', () => {
  const cells = grid(places);
  const seen = new Map();
  const clashes = [];
  for (const [id, cell] of cells) {
    const key = `${cell.col},${cell.row}`;
    if (seen.has(key)) clashes.push(`${seen.get(key)} and ${id} both stand on ${key}`);
    seen.set(key, id);
  }
  assert.deepEqual(clashes, [], clashes.join('; '));
});

test('every address is placed', () => {
  const cells = grid(places);
  assert.equal(cells.size, places.length);
  for (const p of places) assert.ok(cells.has(p.id), `${p.id} is not on the map`);
});

// Seven of the twenty-five addresses are type "work", and every one of them was
// the same low shed with the same stack: the docks, the haulage yard, the cab
// stand, the forecourt, the scrapyard and two filling stations, indistinguishable
// on a map you are supposed to read your own city off.
test('the addresses that do different things do not look identical', () => {
  const shape = p => JSON.stringify(blockFor(p.type, p.id).parts);
  const byShape = new Map();
  for (const p of JSON.parse(readFileSync(new URL('../core/locations.json', import.meta.url), 'utf8'))) {
    const key = shape(p);
    byShape.set(key, [...(byShape.get(key) || []), p.id]);
  }
  // Houses may match houses and rackets may match rackets — a row of shops
  // should look like a row of shops. What must not happen is one silhouette
  // standing for most of a kind of work.
  const worst = [...byShape.values()].sort((a, b) => b.length - a.length)[0];
  assert.ok(worst.length <= 5, `${worst.length} addresses share one silhouette: ${worst.join(', ')}`);
});

test('a burlesque is not drawn as a house', () => {
  const house = JSON.stringify(blockFor('home', 'room').parts);
  assert.notEqual(JSON.stringify(blockFor('burlesque', 'burlesque').parts), house);
});
