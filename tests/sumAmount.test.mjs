import {test} from 'node:test';
import assert from 'node:assert/strict';
import {actionUnit, formatActionAmount} from '../.runtime/frontend-test/sumAmount.js';

const money = value => `$${value}`;
const goods = [{id: 'moonshine', unit: 'crate'}, {id: 'cigarettes', unit: 'carton'}];

test('trade quantities use the public goods unit rather than a dollar price', () => {
  assert.equal(formatActionAmount(5, money, actionUnit('buy:moonshine', goods)), '5 crates');
  assert.equal(formatActionAmount(1, money, actionUnit('sell:moonshine', goods)), '1 crate');
  assert.equal(formatActionAmount(3, money, actionUnit('sell:cigarettes', goods)), '3 cartons');
});

test('money decisions keep money labels while unknown trade goods stay quantities', () => {
  for (const id of ['deposit', 'withdraw', 'bankroll', 'wage', 'buy_car:1']) {
    assert.equal(formatActionAmount(5, money, actionUnit(id, goods)), '$5');
  }
  assert.equal(formatActionAmount(2, money, actionUnit('buy:arms')), '2 units');
});
