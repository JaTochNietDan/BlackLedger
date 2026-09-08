import test from 'node:test';
import assert from 'node:assert/strict';
import {pressPlate} from '../.runtime/frontend-test/art.js';

// The plate is markup the browser has to accept: a malformed one renders as
// nothing at all, and nothing at all is what a missing picture looks like.
const subjects = [
  {kind: 'place', id: 'club', name: 'The Monarch'},
  {kind: 'person', id: 'vittorio', name: 'Vittorio Bellandi'},
  {kind: 'city', name: 'Bellwether'},
];
const kinds = ['killing', 'robbery', 'police', 'business', 'politics', 'war', 'attack'];

test('every subject and kind produces a well-formed plate', () => {
  for (const subject of subjects) {
    for (const kind of kinds) {
      const svg = pressPlate(kind, subject, kind.toUpperCase() + ' AT ' + subject.name);
      assert.ok(svg.startsWith('<svg'), `${kind}/${subject.kind} did not start an svg`);
      assert.ok(svg.endsWith('</svg>'), `${kind}/${subject.kind} did not close`);
      assert.equal((svg.match(/<g/g) || []).length, (svg.match(/<\/g>/g) || []).length,
        `${kind}/${subject.kind} left a group open`);
      assert.ok(svg.includes('viewBox="0 0 112 80"'));
      assert.ok(svg.length > 400, `${kind}/${subject.kind} drew almost nothing`);
    }
  }
});

test('the same story is always the same plate and different stories are not', () => {
  const subject = subjects[0];
  const a = pressPlate('robbery', subject, 'ROBBERY AT THE MONARCH');
  assert.equal(a, pressPlate('robbery', subject, 'ROBBERY AT THE MONARCH'));
  assert.notEqual(a, pressPlate('robbery', subject, 'THEFT AT THE MONARCH'));
});

test('a name with a quote in it cannot break out of the label', () => {
  const svg = pressPlate('killing', {kind: 'person', id: 'x', name: 'O\'Hare "Fingers"'}, 'X KILLED');
  assert.ok(!svg.includes('label="O\'Hare "'), 'the label was not escaped');
  assert.ok(svg.includes('&#39;') || svg.includes('&quot;'));
});
