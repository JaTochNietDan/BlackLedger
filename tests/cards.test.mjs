import test from 'node:test';
import assert from 'node:assert/strict';
import {drumFaces, pipOf, isRedSuit, knownCard, clothRows, clothColour, outsideBets, wheelOrder, wheelAngle, ballAngle, clothTable, wheelPaint, reelStops, reelWindow, REEL_WINDOW, drumRun, drumRunIDs} from '../.runtime/frontend-test/cards.js';

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

test('the ball always travels forward and stops on the pocket the core chose', () => {
  // The wheel used to snap: the pocket the core picked was simply outlined. The
  // ball goes round to it now, and the one thing the animation must never do is
  // land anywhere but on the number the core actually spun.
  const turns = 4;
  for (const [from, pocket] of [[0, 26], [37.5, 0], [1440, 32], [-15, 3], [123.4, 15]]) {
    const to = ballAngle(from, pocket, turns);
    assert.ok(to > from, `the ball went backwards from ${from} to ${to}`);
    assert.ok(to - from >= turns * 360, `the ball barely moved: ${to - from} degrees`);
    assert.ok(to - from < (turns + 1) * 360, `the ball spun for ever: ${to - from} degrees`);
    const settled = ((to % 360) + 360) % 360;
    const want = wheelAngle(pocket);
    assert.ok(Math.abs(settled - want) < 0.001, `settled at ${settled}, the pocket is at ${want}`);
  }
});

test('a wheel nobody has spun yet does not move the ball', () => {
  assert.equal(ballAngle(500, 0, 0), 500);
});

test('the cloth is laid out the way a table is: three rows of twelve', () => {
  const rows = clothTable();
  assert.equal(rows.length, 3);
  for (const row of rows) assert.equal(row.length, 12);
  const flat = rows.flat();
  assert.equal(new Set(flat).size, 36);
  for (let n = 1; n <= 36; n++) assert.ok(flat.includes(n), n + ' is not on the cloth');
  // The top row is the one that pays the third column: 3, 6, 9 and so on.
  for (const n of rows[0]) assert.equal(n % 3, 0);
  for (const n of rows[2]) assert.equal(n % 3, 1);
  // And it reads left to right, low to high, like every table ever built.
  assert.deepEqual(rows[2].slice(0, 3), [1, 4, 7]);
});

test('the wheel is painted in the pockets own order and colours', () => {
  const paint = wheelPaint();
  // One slice per pocket — two angles each — plus the one the gradient starts
  // from, which is what turns a pocket's own angle into the middle of its slice.
  assert.equal((paint.match(/deg/g) || []).length, wheelOrder.length * 2 + 1);
  assert.ok(paint.startsWith('conic-gradient('));
  // One green slice, named once: the nought, and nothing else on the wheel.
  const greens = (paint.match(/#1f6b45/g) || []).length;
  assert.equal(greens, 1, 'the nought is the only green pocket');
});

// The drums.
test('a drum shows the face it landed on, on the line', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 2, pays: 50},
    {id: 'cherry', face: 'CHERRY', stops: 2, pays: 25},
  ];
  for (const s of strip) {
    const win = reelWindow(strip, s.face);
    assert.equal(win.length, REEL_WINDOW);
    assert.equal(win[1], s.face, `${s.face} was not on the line: ${win.join('/')}`);
  }
});

test('a drum only ever shows faces the core has', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bell', face: 'BELL', stops: 3, pays: 25},
  ];
  const faces = new Set(strip.map(s => s.face));
  for (let nudge = 0; nudge < 6; nudge++)
    for (const f of reelWindow(strip, 'BELL', nudge))
      assert.ok(faces.has(f), `the drum showed ${f}, which is not on the strip`);
});

test('the strip is as long as the odds say it is', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 2, pays: 50},
    {id: 'bell', face: 'BELL', stops: 17, pays: 25},
  ];
  assert.equal(reelStops(strip).length, 20);
});

test('three drums on the same face need not show the same shoulders', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 4, pays: 50},
    {id: 'bell', face: 'BELL', stops: 15, pays: 25},
  ];
  const runs = new Set([0, 1, 2].map(n => reelWindow(strip, 'BAR', n).join('/')));
  assert.ok(runs.size > 1, 'every drum showing BAR looked identical');
});

test('an empty strip still fills the window', () => {
  assert.equal(reelWindow([], 'anything').length, REEL_WINDOW);
});

test('a drum shows where it is going to stop from the moment the handle goes down', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 60},
    {id: 'bell', face: '🔔', stops: 4, pays: 18},
    {id: 'lemon', face: '🍋', stops: 7, pays: 0},
    {id: 'cherry', face: '🍒', stops: 8, pays: 4},
  ];
  // The three the core landed on, which is deliberately not the strip's first
  // face: that is what the case used to fall back to.
  const line = ['bell', 'lemon', 'cherry'];
  const rolled = drumFaces(strip, line, true);
  assert.deepEqual(
    rolled.map(w => w[1]),
    ['🔔', '🍋', '🍒'],
    'the payline does not show what the core landed on',
  );
  // Every drum reads the same whether its neighbours have stopped or not: the
  // case used to fall back to the first symbol on the strip for a drum that had
  // stopped while the others turned, and then all three jumped to the result
  // when the last one came down.
  for (let i = 0; i < 3; i++) {
    assert.deepEqual(rolled[i], drumFaces(strip, line, true)[i], 'a drum changed its mind');
  }
  // And before the handle goes down it is the strip's own first face, not a
  // result nobody has pulled for.
  const idle = drumFaces(strip, line, false);
  assert.deepEqual(
    idle.map(w => w[1]),
    ['7', '7', '7'],
    'a machine nobody has pulled is not sitting on its own first face',
  );
});

// The drum travels to where it stops. It used to show its answer from the
// moment the handle went down and shake on the spot, and before that it showed
// one thing and flipped to another at the end. The run has to begin with the
// faces it lands on, so neither can happen: the answer is already there and
// the travelling is behind it.
test('a drum runs behind the faces it lands on', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 2, pays: 50},
    {id: 'cherry', face: 'CHERRY', stops: 2, pays: 25},
    {id: 'bell', face: 'BELL', stops: 3, pays: 20},
  ];
  const window = ['BELL', '7', 'BAR'];
  const run = drumRun(strip, window, 0, 14);
  assert.equal(run.length, window.length + 14, 'the run is the window plus what it travels');
  assert.deepEqual(run.slice(0, 3), window, 'the drum lands on something other than its window');
  for (const face of run) assert.ok(strip.some(s => s.face === face), 'a face not on the strip: ' + face);
});

test('three drums do not turn through the same symbols in step', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 2, pays: 50},
    {id: 'cherry', face: 'CHERRY', stops: 2, pays: 25},
    {id: 'bell', face: 'BELL', stops: 3, pays: 20},
  ];
  const window = ['BELL', '7', 'BAR'];
  const runs = [0, 1, 2].map(i => drumRun(strip, window, i, 14).slice(3).join(','));
  assert.notEqual(runs[0], runs[1], 'the first two drums turn through the same symbols');
  assert.notEqual(runs[1], runs[2], 'the last two drums turn through the same symbols');
});

test('the run is painted by the same ids as the faces', () => {
  const strip = [
    {id: 'seven', face: '7', stops: 1, pays: 100},
    {id: 'bar', face: 'BAR', stops: 2, pays: 50},
    {id: 'cherry', face: 'CHERRY', stops: 2, pays: 25},
    {id: 'bell', face: 'BELL', stops: 3, pays: 20},
  ];
  const window = ['BELL', '7', 'BAR'];
  const faces = drumRun(strip, window, 1, 6);
  const ids = drumRunIDs(strip, window, 1, 6);
  assert.equal(ids.length, faces.length);
  faces.forEach((face, i) => {
    const want = strip.find(s => s.face === face);
    assert.equal(ids[i], want ? want.id : '', 'the drum is painted as something it is not');
  });
});
