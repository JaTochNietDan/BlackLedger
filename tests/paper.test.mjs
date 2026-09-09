import {test} from 'node:test';import assert from 'node:assert/strict';
import {unreadInLatest} from '../.runtime/frontend-test/paper.js';

// Newest first, which is how the core sends it.
const paper = [
  {id:'f',day:32},{id:'e',day:32},{id:'d',day:32},
  {id:'c',day:31},{id:'b',day:31},
  {id:'a',day:30},
];

test('an unread paper counts only the latest issue', () => {
  // Six stories on file, three of them in today's issue. The badge used to say
  // six and the screen showed three.
  assert.equal(unreadInLatest(paper, null), 3);
});

test('having read today, nothing is outstanding', () => {
  assert.equal(unreadInLatest(paper, 'f'), 0);
});

test('part of today read leaves the rest', () => {
  assert.equal(unreadInLatest(paper, 'e'), 1);
  assert.equal(unreadInLatest(paper, 'd'), 2);
});

test('a story from an older issue does not mark today read', () => {
  // This is the fault: opening the paper on day 30 used to strike off days 31
  // and 32 as well.
  assert.equal(unreadInLatest(paper, 'c'), 3);
  assert.equal(unreadInLatest(paper, 'a'), 3);
});

test('an empty paper is nothing to read', () => {
  assert.equal(unreadInLatest([], null), 0);
  assert.equal(unreadInLatest([], 'f'), 0);
});
