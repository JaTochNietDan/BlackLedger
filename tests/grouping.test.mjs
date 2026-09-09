import {test} from 'node:test';import assert from 'node:assert/strict';
import {placeActions} from '../.runtime/frontend-test/grouping.js';

// The core's own list, as it sends it. The panel used to keep a copy of this
// and the copy was missing "people".
const core = [
  {id:'work',title:'Work',blurb:'Jobs that pay today.'},
  {id:'business',title:'Your premises',blurb:'Keeping what you own earning.'},
  {id:'people',title:'People',blurb:'Who works for you, who owes you, who you know.'},
  {id:'street',title:'The street',blurb:'Work that can go wrong, and hurt somebody.'},
  {id:'standing',title:'Standing',blurb:'Who you are to this city.'},
  {id:'money',title:'Money',blurb:'Moving it, hiding it, or putting it somewhere else.'},
  {id:'travel',title:'Elsewhere',blurb:'Leaving where you are standing.'},
];
const count = placed => placed.reduce((n,s) => n + s.mine.length, 0);

test('every action lands in a section', () => {
  const actions = [
    {id:'delegate',group:'work'}, {id:'fit:door',group:'business'},
    {id:'crew_bonus',group:'people'}, {id:'rob',group:'street'},
    {id:'press',group:'standing'}, {id:'deposit',group:'money'},
    {id:'wait',group:'travel'},
  ];
  const placed = placeActions(core, actions);
  assert.equal(count(placed), actions.length);
  assert.deepEqual(placed.find(s => s.id === 'people').mine.map(a => a.id), ['crew_bonus']);
});

test('a group this build has not heard of is not dropped', () => {
  // A new group added to the core and not yet known here: the panel must show
  // the action somewhere rather than losing it, which is what happened to
  // "people" for as long as the list was copied by hand.
  const actions = [{id:'delegate',group:'work'}, {id:'something',group:'unheard-of'}];
  const placed = placeActions(core, actions);
  assert.equal(count(placed), 2);
  assert.deepEqual(placed[0].mine.map(a => a.id), ['delegate','something']);
});

test('no action is placed twice', () => {
  const actions = [{id:'a',group:'work'},{id:'b',group:'people'},{id:'c',group:'nowhere'}];
  const placed = placeActions(core, actions);
  const seen = placed.flatMap(s => s.mine.map(a => a.id));
  assert.equal(seen.length, new Set(seen).size);
  assert.equal(seen.length, 3);
});

test('with no groups at all it places nothing rather than guessing', () => {
  assert.deepEqual(placeActions([], [{id:'a',group:'work'}]), []);
});
