import test from 'node:test';
import assert from 'node:assert/strict';
import {pipOf, isRedSuit, knownCard, clothRows, clothColour, outsideBets, wheelOrder, wheelAngle} from '../.runtime/frontend-test/cards.js';

test('every suit the core deals has a pip and a colour', () => {
  for (const suit of ['spades', 'hearts', 'diamonds', 'clubs']) {
    assert.ok(pipOf(suit).length > 0, suit + ' has no pip');
  }
  assert.equal(isRedSuit('hearts'), true);
  assert.equal(isRedSuit('spades'), false);
});

test('a card the view does not understand is not guessed at', () => {
  assert.equal(knownCard({rank: 'A', suit: 'spades', value: 11}), true);
  assert.equal(knownCard({rank: 'A', suit: 'swords', value: 11}), false);
  assert.equal(knownCard({rank: '', suit: 'spades', value: 0}), false);
  assert.equal(pipOf('swords'), '');
});

test('the cloth is three columns of twelve, every number once', () => {
  const rows = clothRows();
  assert.equal(rows.length, 12);
  const seen = new Set();
  for (const row of rows) {
    assert.equal(row.length, 3);
    for (const n of row) {
      assert.ok(n >= 1 && n <= 36, n + ' is not on the cloth');
      assert.ok(!seen.has(n), n + ' appears twice');
      seen.add(n);
    }
  }
  assert.equal(seen.size, 36);
});

test('the cloth colours match a real table', () => {
  assert.equal(clothColour(0), 'green');
  let reds = 0, blacks = 0;
  for (let n = 1; n <= 36; n++) {
    if (clothColour(n) === 'red') reds++;
    if (clothColour(n) === 'black') blacks++;
  }
  assert.equal(reds, 18);
  assert.equal(blacks, 18);
  // Not every other number: 10 and 11 are both black, 18 and 19 both red.
  assert.equal(clothColour(10), 'black');
  assert.equal(clothColour(11), 'black');
  assert.equal(clothColour(18), 'red');
  assert.equal(clothColour(19), 'red');
});

test('the outside bets name ids the core takes', () => {
  const known = new Set(['red','black','odd','even','low','high','dozen1','dozen2','dozen3']);
  const bets = outsideBets();
  assert.equal(bets.length, known.size);
  for (const b of bets) assert.ok(known.has(b.id), b.id + ' is not a bet the core accepts');
});

test('the wheel face holds every pocket once and alternates colours', () => {
  assert.equal(wheelOrder.length, 37);
  assert.equal(new Set(wheelOrder).size, 37);
  for (const n of wheelOrder) assert.ok(n >= 0 && n <= 36);
  // Either side of the nought is red then black all the way round.
  for (let i = 1; i < wheelOrder.length; i++) {
    const a = clothColour(wheelOrder[i - 1]), b = clothColour(wheelOrder[i]);
    if (a === 'green' || b === 'green') continue;
    assert.notEqual(a, b, wheelOrder[i - 1] + ' and ' + wheelOrder[i] + ' are both ' + a);
  }
});

test('a pocket sits somewhere on the face and the nought sits at the top', () => {
  assert.equal(wheelAngle(0), 0);
  assert.ok(wheelAngle(26) > 0 && wheelAngle(26) < 360);
  assert.equal(wheelAngle(99), 0);
});
