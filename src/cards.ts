// What a playing card looks like, and what a roulette cloth is laid out as.
// Nothing here decides anything: the core deals the cards, spins the wheel and
// settles the money. This is only how those facts are arranged on a table.

export interface Card {
  rank: string;
  suit: string;
  value: number;
}

// The four suits, as the core names them, with the character to print and
// whether it is a red suit. A card the core sends with a suit not in here is
// drawn as a plain rectangle rather than as a guess.
const suits: {[id: string]: {pip: string; red: boolean}} = {
  spades: {pip: '♠', red: false},
  hearts: {pip: '♥', red: true},
  diamonds: {pip: '♦', red: true},
  clubs: {pip: '♣', red: false},
};

export const pipOf = (suit: string) => suits[suit]?.pip ?? '';
export const isRedSuit = (suit: string) => suits[suit]?.red ?? false;

// A card the view knows how to draw. Anything else is shown face down rather
// than invented.
export const knownCard = (c: Card) => !!c && !!suits[c.suit] && !!c.rank;

// The cloth. A roulette table is three columns of twelve with the nought
// across the top, and it has to be laid out in that order or a player who has
// stood at one will not recognise it.
export const clothRows = (): number[][] => {
  const rows: number[][] = [];
  for (let row = 0; row < 12; row++) {
    rows.push([row * 3 + 3, row * 3 + 2, row * 3 + 1]);
  }
  return rows;
};

// The red pockets, which are not every other number. This is the same layout
// the core holds; the view needs it only to colour the cloth, and the core
// remains the authority on what actually won.
const reds = new Set([1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36]);
export const clothColour = (n: number) => (n === 0 ? 'green' : reds.has(n) ? 'red' : 'black');

// The outside bets along the bottom, in the order they sit on a real cloth,
// paired with the bet ids the core accepts.
export const outsideBets = () => [
  {id: 'dozen1', label: '1st 12', wide: true},
  {id: 'dozen2', label: '2nd 12', wide: true},
  {id: 'dozen3', label: '3rd 12', wide: true},
  {id: 'low', label: '1-18', wide: false},
  {id: 'even', label: 'Even', wide: false},
  {id: 'red', label: 'Red', wide: false},
  {id: 'black', label: 'Black', wide: false},
  {id: 'odd', label: 'Odd', wide: false},
  {id: 'high', label: '19-36', wide: false},
];

// Where a pocket sits on the wheel face, in degrees clockwise from the top.
// The order is the wheel's own, which is not the cloth's and not numerical:
// it alternates colours and spreads the low numbers around the rim.
export const wheelOrder = [
  0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14,
  31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26,
];

export const wheelAngle = (pocket: number) => {
  const at = wheelOrder.indexOf(pocket);
  return at < 0 ? 0 : at * (360 / wheelOrder.length);
};

/**
 * Where the ball ends up on the face, in degrees, given where it is now and the
 * pocket the CORE spun. It always travels forward — a ball that jumps backwards
 * to save half a turn is a ball nobody believes — and it always settles exactly
 * on that pocket's own angle, because the animation is not allowed to decide
 * anything. Turns is how many whole revolutions it makes on the way; zero means
 * do not move it at all, which is what a wheel nobody has spun looks like.
 */
export function ballAngle(from: number, pocket: number, turns: number) {
  if (turns <= 0) return from;
  const target = wheelAngle(pocket);
  const now = ((from % 360) + 360) % 360;
  const forward = (((target - now) % 360) + 360) % 360;
  return from + turns * 360 + forward;
}

/**
 * The cloth as a table is actually laid: three rows of twelve, reading left to
 * right and low to high, with the row that pays the third column at the top.
 * clothRows above is the same numbers stood on end, which is what the panel
 * needed when the wheel was a column in a sidebar.
 */
export function clothTable() {
  const rows: number[][] = [[], [], []];
  for (let n = 1; n <= 36; n++) rows[2 - ((n - 1) % 3)].push(n);
  return rows;
}

/**
 * The wheel head, painted in the pockets' own order. One slice per pocket, in
 * the colour that pocket is — the same colours the cloth uses, because they are
 * the same numbers. The core decides what wins; this only says what the thing
 * looks like.
 */
export function wheelPaint() {
  const slice = 360 / wheelOrder.length;
  const ink = {red: '#8e2b2b', black: '#1a1a1a', green: '#1f6b45'};
  const stops = wheelOrder.map((n, i) => {
    const colour = ink[clothColour(n) as keyof typeof ink];
    return `${colour} ${(i * slice).toFixed(3)}deg ${((i + 1) * slice).toFixed(3)}deg`;
  });
  // Half a slice back, so a pocket's own angle points at the middle of it.
  return `conic-gradient(from ${(-slice / 2).toFixed(3)}deg, ${stops.join(', ')})`;
}

// The drums of a bandit.
//
// The machine drew one face per drum: three letters in three boxes, which is a
// picture of a result rather than a machine. A real drum is a strip of faces
// that turns behind a window, and what you see is three of them at a time with
// the payline across the middle one.
//
// The strip is the core's own — the same twenty stops it works the odds out
// over — so a machine cannot show a face the rules do not have.

export type Reel = {id: string; face: string; stops: number; pays: number};

// REEL_WINDOW is how many faces of the strip the case shows at once. Three: the
// one on the line, and the shoulder of the one above and below it, which is
// what makes it read as a drum rather than a card.
export const REEL_WINDOW = 3;

// reelStops expands the core's strip into the actual run of faces on the drum:
// a symbol with four stops appears four times. This is what turns behind the
// window, and its length is the number the core divides by.
export function reelStops(strip: Reel[]): string[] {
  const out: string[] = [];
  for (const s of strip) for (let i = 0; i < s.stops; i++) out.push(s.face);
  return out;
}

// reelWindow is the three faces showing when the drum has stopped with `landed`
// on the payline. The middle one is the result; the other two are its
// neighbours on the strip, which is why a machine feels like it nearly paid.
export function reelWindow(strip: Reel[], landed: string, nudge = 0): string[] {
  const stops = reelStops(strip);
  if (stops.length === 0) return Array(REEL_WINDOW).fill('—');
  // Where on the strip this face sits. A face with several stops has several
  // homes; nudge picks between them so three drums showing the same symbol do
  // not show identical neighbours.
  const homes = stops.map((f, i) => (f === landed ? i : -1)).filter(i => i >= 0);
  const at = homes.length ? homes[Math.abs(nudge) % homes.length] : 0;
  const out: string[] = [];
  for (let i = -1; i <= 1; i++) out.push(stops[(at + i + stops.length * 2) % stops.length]);
  return out;
}

// What the three drums are showing, which is the same question whether they are
// turning or standing still.
//
// The case used to ask "has everything stopped, or is this drum still going?"
// and show the result only then. A drum that had stopped while the others were
// still turning satisfied neither half, so it fell back to the first symbol on
// the strip — and when the last drum came down all three jumped to the real
// result at once. From the outside that is a machine changing its mind: "it
// shows 7-7- as it progresses then at the very end it flips to bell, lemon,
// cherry."
//
// A drum shows where it is going to stop from the moment the handle goes down.
// The blur is the animation's job and the face underneath is already right, so
// there is nothing left to snap to.
export function drumIDs(strip: Reel[], line: string[], pulled: boolean): string[][] {
  const idOf = (face: string) => strip.find(s => s.face === face)?.id ?? '';
  return drumFaces(strip, line, pulled).map(w => w.map(idOf));
}

export function drumFaces(strip: Reel[], line: string[], pulled: boolean): string[][] {
  const faceOf = (id?: string) => strip.find(s => s.id === id)?.face ?? '—';
  return [0, 1, 2].map(i =>
    pulled ? reelWindow(strip, faceOf(line[i]), i) : reelWindow(strip, faceOf(strip[0]?.id), i),
  );
}

// What a drum turns through on the way to where it stops.
//
// The drums did not turn. They showed where they were going to stop from the
// moment the handle went down and shook on the spot for a second, which is a
// machine trembling rather than a machine spinning: "it should scroll through
// the items before displaying the final result. Right now they just shake but
// the result is shown already."
//
// So the column is the three faces it lands on, and behind them a run of the
// strip to travel past. The landing faces are the ones the core sent and they
// are decided before a pixel moves — nothing is picked here, and nothing snaps
// at the end, which is the fault this machine had before the shake.
//
// Walked from a different offset on each drum so three drums do not turn
// through the same symbols in step.
export function drumRun(strip: Reel[], window: string[], i: number, turns: number): string[] {
  const stops = strip.map(s => s.face);
  if (stops.length === 0) return window;
  const out = [...window];
  for (let n = 0; n < turns; n++) {
    out.push(stops[(n * 7 + i * 5 + 1) % stops.length]);
  }
  return out;
}

// And the same run as symbol ids, so the drum is painted rather than spelled.
export function drumRunIDs(strip: Reel[], window: string[], i: number, turns: number): string[] {
  const idOf = (face: string) => strip.find(s => s.face === face)?.id ?? '';
  return drumRun(strip, window, i, turns).map(idOf);
}
