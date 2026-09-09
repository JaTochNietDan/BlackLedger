// The isometric city.
//
// What stood here before was a grid of painted cards in three rows: twelve
// buildings as pictures in a list, with people drawn as portraits sliding in a
// straight line between two cards. This is the city itself, drawn in plan, from
// the coordinates the core already keeps for every address.
//
// The rule for this module: it decides *shape*, never truth. Where a building
// stands, who is in it and what is happening come from the core. Everything
// here is how to draw that, and it is built out of small parts on purpose —
// a block is a list of boxes, so a better roof, a taller tenement or a piece of
// generated art can replace one part without touching the rest.

export type Vec = {x: number; y: number};

// A 2:1 isometric projection. TILE is the size of one plan cell on screen; the
// city's own coordinates run 130..930 by 95..585, so CELL divides them into a
// grid of about thirteen by nine.
export const TILE = {w: 108, h: 54};
export const CELL = 48;

// plan converts the city's own coordinates into tile space, which is the only
// place the two systems meet.
export const plan = (x: number, y: number): Vec => ({x: x / CELL, y: y / CELL});

// project turns a point in tile space into a point on the screen.
export function project(p: Vec): Vec {
  return {x: (p.x - p.y) * (TILE.w / 2), y: (p.x + p.y) * (TILE.h / 2)};
}

// depth is how near the front of the picture something is. Anything drawn later
// covers what came before, so this is the only thing that decides what a
// building can hide.
export const depth = (p: Vec) => p.x + p.y;

// A part of a building: a box standing on the block's own footprint, in tiles,
// with a height in tiles and its own colour. Everything visible is made of
// these, which is what makes a block swappable a piece at a time.
export type Part = {
  dx?: number; dy?: number;      // offset within the block, in tiles
  w: number; d: number; h: number;
  base?: number;                  // how far off the ground it starts
  top: string; left: string; right: string;
  flat?: boolean;                 // a slab: no walls worth drawing
};

export type Block = {
  w: number; d: number;           // footprint in tiles
  parts: Part[];
  label?: string;
};

// FOOTPRINT is how much ground each kind of place takes, in tiles, kept in one
// table because it is the only thing that decides whether two buildings stand
// on each other. A test reads this literal and the city's own coordinates and
// fails if any two addresses would overlap at the current CELL — which is how
// the first version was caught putting The Mariner through the Herald.
export const FOOTPRINT: Record<string, [number, number]> = {
  casino: [2.2, 1.8],
  market: [2.4, 1.6],
  civic: [1.9, 1.6],
  racket: [1.8, 1.5],
  work: [2.6, 1.4],
  bar: [1.5, 1.4],
  home: [1.5, 1.3],
};

// size is the footprint of a place, falling back to a house for a kind of
// address nobody has drawn yet: the city has to be complete before it is
// pretty.
export const size = (type: string): [number, number] => FOOTPRINT[type] || FOOTPRINT.home;

// The palette. Kept in one place so the whole city can be recoloured for night,
// for a district, or for a moment the camera is holding on.
const stone = {top: '#6a6357', left: '#3f3a33', right: '#2c2823'};
const brick = {top: '#7a5648', left: '#4a332b', right: '#33231e'};
const pale = {top: '#8a8272', left: '#565043', right: '#3a352c'};
const dark = {top: '#4a4a46', left: '#2b2b28', right: '#1d1d1b'};
const roofTile = {top: '#5d4038', left: '#3a2823', right: '#281b18'};
const glass = {top: '#c9a86a', left: '#8a6f3f', right: '#5e4b2a'};

// box is the one primitive: a rectangular solid on the block's footprint.
const box = (w: number, d: number, h: number, c: typeof stone, o: Partial<Part> = {}): Part =>
  ({w, d, h, ...c, ...o});

// Blocks by the kind of place the core says it is. A type it has never heard of
// still gets a building rather than a hole, because the city has to be complete
// before it is pretty.
export function blockFor(type: string): Block {
  switch (type) {
    case 'casino':
      return {w: size('casino')[0], d: size('casino')[1], parts: [
        box(2.2, 1.8, 1.5, stone),
        box(2.2, 1.8, .18, glass, {base: 1.5}),          // a lit band along the top
        box(.9, .5, .5, roofTile, {base: 1.68, dx: .65, dy: .65}),
        box(.24, .24, .7, dark, {base: 2.18, dx: .98, dy: .78}),
      ]};
    case 'market':
      return {w: size('market')[0], d: size('market')[1], parts: [
        box(2.4, 1.6, 1.25, pale),
        box(2.4, 1.6, .12, roofTile, {base: 1.25}),
        box(.3, .3, .55, pale, {base: 1.37, dx: .2, dy: .65}),
        box(.3, .3, .55, pale, {base: 1.37, dx: 1.9, dy: .65}),
      ]};
    case 'civic':
      return {w: size('civic')[0], d: size('civic')[1], parts: [
        box(1.9, 1.6, 1.6, stone),
        box(2.05, 1.75, .1, dark, {base: 1.6, dx: -.07, dy: -.07}),
        box(.7, .7, .45, stone, {base: 1.7, dx: .6, dy: .45}),
      ]};
    case 'racket':
      return {w: size('racket')[0], d: size('racket')[1], parts: [
        box(1.8, 1.5, 1.05, brick),
        box(1.9, 1.6, .1, dark, {base: 1.05, dx: -.05, dy: -.05}),
        box(.5, .34, .3, dark, {base: 1.15, dx: 1.1, dy: .6}),   // a vent on the roof
      ]};
    case 'work':
      return {w: size('work')[0], d: size('work')[1], parts: [
        box(2.6, 1.4, .8, dark),
        box(2.6, 1.4, .22, roofTile, {base: .8}),
        box(.3, .3, 1.1, dark, {base: 1.02, dx: 2.1, dy: .5}),   // a stack
      ]};
    case 'bar':
      return {w: size('bar')[0], d: size('bar')[1], parts: [
        box(1.5, 1.4, 1.15, brick),
        box(1.5, .14, .3, glass, {base: .35, dy: 1.36}),         // a lit window on the street
        box(1.6, 1.5, .09, dark, {base: 1.15, dx: -.05, dy: -.05}),
      ]};
    case 'home':
    default:
      return {w: size('home')[0], d: size('home')[1], parts: [
        box(1.5, 1.3, 1.35, brick),
        box(1.6, 1.4, .1, roofTile, {base: 1.35, dx: -.05, dy: -.05}),
        box(.22, .22, .4, dark, {base: 1.45, dx: 1.05, dy: .3}), // a chimney
      ]};
  }
}

// faces returns the three visible sides of a part as screen polygons, in the
// order they must be drawn. A box is only ever seen from one corner, so three
// faces is the whole of it.
// Solid is the geometry of a box and nothing else. faces only ever needed the
// shape, so a car — whose colours are numbers because Pixi wants numbers —
// goes through the same projection as a building whose colours are strings.
export type Solid = {dx?: number; dy?: number; w: number; d: number; h: number; base?: number};

export function faces(origin: Vec, part: Solid) {
  const bx = origin.x + (part.dx || 0), by = origin.y + (part.dy || 0);
  const base = part.base || 0;
  const lift = (v: Vec, up: number) => ({x: v.x, y: v.y - up * TILE.h});

  const a = project({x: bx, y: by});
  const b = project({x: bx + part.w, y: by});
  const c = project({x: bx + part.w, y: by + part.d});
  const d = project({x: bx, y: by + part.d});

  const top = [a, b, c, d].map(v => lift(v, base + part.h));
  const left = [lift(d, base + part.h), lift(c, base + part.h), lift(c, base), lift(d, base)];
  const right = [lift(c, base + part.h), lift(b, base + part.h), lift(b, base), lift(c, base)];
  return {top, left, right};
}

export const points = (vs: Vec[]) => vs.map(v => `${v.x.toFixed(1)},${v.y.toFixed(1)}`).join(' ');

// Where a person stands on a plot. People gather at the front corner of a
// building rather than in the middle of its roof, and spread along the pavement
// so a crowded address reads as a crowd instead of one figure.
export function standing(at: Vec, block: {w: number; d: number}, i: number, of: number): Vec {
  const spread = Math.min(of, 6);
  const along = spread <= 1 ? .5 : .15 + (i % spread) / (spread - 1) * .7;
  const rank = Math.floor(i / spread);          // a second row behind the first
  return {x: at.x + block.w * along, y: at.y + block.d + .22 + rank * .28};
}

// The door of a building, which is where a journey starts and ends: the near
// corner of its plot, on the pavement.
export const door = (at: Vec, block: {w: number; d: number}): Vec =>
  ({x: at.x + block.w / 2, y: at.y + block.d + .25});

// Somewhere along a walk between two doors.
export const between = (a: Vec, b: Vec, t: number): Vec =>
  ({x: a.x + (b.x - a.x) * t, y: a.y + (b.y - a.y) * t});

// The roof of a block, which is where a moment happens: an explosion goes off
// above a building, not on the pavement in front of it.
export const roof = (at: Vec, block: {w: number; d: number}, height: number): Vec => {
  const c = project({x: at.x + block.w / 2, y: at.y + block.d / 2});
  return {x: c.x, y: c.y - height};
};

// ---------------------------------------------------------------------------
// The streets.
//
// The city was twelve buildings standing in the dark, and the first attempt at
// roads joined every building to its nearest neighbours with L-shaped paths.
// That produced a web: streets at every angle of approach, laid over each
// other, meeting nothing squarely. A city is a grid.
//
// So the ground is a grid, and the buildings are placed on it. The addresses
// keep their own coordinates as the source of where they belong — the core's
// x and y decide which block a building gets — but the blocks themselves are
// regular, the roads run the full width and height of the city, and every
// corner is square. Bellwether's own coordinates fall into a six by four grid
// with no two addresses wanting the same block, which is why this works
// without moving anything by hand.

// BLOCK is the pitch from one street to the next; ROAD is how much of that is
// carriageway; PAVE is the footway inside each block, between kerb and wall.
export const BLOCK = 3.6;
export const ROAD = 1.0;
export const PAVE = .42;

export type Cell = {col: number; row: number};
export type Rect = {x: number; y: number; w: number; d: number};
export type Segment = {a: Vec; b: Vec};

// island is the ground inside one block: everything that is not carriageway.
export const island = (c: Cell): Rect => ({
  x: c.col * BLOCK + ROAD / 2, y: c.row * BLOCK + ROAD / 2,
  w: BLOCK - ROAD, d: BLOCK - ROAD,
});

// plot is where the building itself stands, set back from the kerb by the
// width of the pavement on every side.
export const plot = (c: Cell): Rect => {
  const i = island(c);
  return {x: i.x + PAVE, y: i.y + PAVE, w: i.w - PAVE * 2, d: i.d - PAVE * 2};
};

// The middle of a block, which is where a building is centred and where
// anybody standing at that address is standing.
export const middle = (c: Cell): Vec => {
  const i = island(c);
  return {x: i.x + i.w / 2, y: i.y + i.d / 2};
};

// grid puts every address in a block. Which block comes from the address's own
// coordinates, so the city keeps its shape: the docks stay west, Cypress House
// stays north-east. Two addresses that want the same block push apart rather
// than stacking, though Bellwether as it stands never needs it.
export function grid(places: {id: string; x: number; y: number}[]): Map<string, Cell> {
  const tiles = places.map(p => ({id: p.id, ...plan(p.x, p.y)}));
  const minX = Math.min(...tiles.map(t => t.x)), minY = Math.min(...tiles.map(t => t.y));
  const taken = new Map<string, string>();
  const out = new Map<string, Cell>();
  for (const t of tiles.sort((a, b) => (a.x + a.y) - (b.x + b.y))) {
    let cell = {col: Math.round((t.x - minX) / BLOCK), row: Math.round((t.y - minY) / BLOCK)};
    // If somebody is already there, step outwards until a block is free.
    for (let ring = 1; taken.has(`${cell.col},${cell.row}`) && ring < 6; ring++) {
      const around = [{col: cell.col + ring, row: cell.row}, {col: cell.col, row: cell.row + ring},
                      {col: cell.col - ring, row: cell.row}, {col: cell.col, row: cell.row - ring}];
      const free = around.find(c => c.col >= 0 && c.row >= 0 && !taken.has(`${c.col},${c.row}`));
      if (free) cell = free;
    }
    taken.set(`${cell.col},${cell.row}`, t.id);
    out.set(t.id, cell);
  }
  return out;
}

// The extent of the city in blocks, which is what the roads have to span.
export function bounds(cells: Map<string, Cell>): {cols: number; rows: number} {
  let cols = 0, rows = 0;
  for (const c of cells.values()) { cols = Math.max(cols, c.col); rows = Math.max(rows, c.row) }
  return {cols: cols + 1, rows: rows + 1};
}

// The carriageways, running the full width and height of the city so every
// junction is square and no street stops in the middle of nowhere.
export function carriageways({cols, rows}: {cols: number; rows: number}): Segment[] {
  const out: Segment[] = [];
  const right = cols * BLOCK, bottom = rows * BLOCK;
  for (let col = 0; col <= cols; col++) {
    const x = col * BLOCK;
    out.push({a: {x, y: 0}, b: {x, y: bottom}});
  }
  for (let row = 0; row <= rows; row++) {
    const y = row * BLOCK;
    out.push({a: {x: 0, y}, b: {x: right, y}});
  }
  return out;
}

// A walk from one address to another, along the streets rather than through
// the buildings: out of the block to the nearest corner, along one axis, then
// the other, and in again.
export function walk(from: Cell, to: Cell): Vec[] {
  const start = middle(from), end = middle(to);
  // The junction lines this walk uses: the street on the near side of each.
  const lane = (a: number, b: number) => (a <= b ? Math.max(a, b) * 0 + (a + 1) * BLOCK - BLOCK : a * BLOCK);
  const outX = from.col <= to.col ? (from.col + 1) * BLOCK : from.col * BLOCK;
  const inY = to.row <= from.row ? (to.row + 1) * BLOCK : to.row * BLOCK;
  void lane;
  return [
    start,
    {x: outX, y: start.y},          // out to the street
    {x: outX, y: inY},              // along it
    {x: end.x, y: inY},             // round the corner
    end,                            // and in
  ];
}

// along is a point some fraction of the way through a set of waypoints,
// measured by distance rather than by how many corners there are — so somebody
// crossing a long street does not sprint the last leg.
export function along(path: Vec[], t: number): Vec {
  if (path.length === 0) return {x: 0, y: 0};
  if (path.length === 1) return path[0];
  const legs = [];
  let total = 0;
  for (let i = 0; i + 1 < path.length; i++) {
    const d = Math.hypot(path[i + 1].x - path[i].x, path[i + 1].y - path[i].y);
    legs.push(d); total += d;
  }
  if (total === 0) return path[0];
  let want = Math.min(1, Math.max(0, t)) * total;
  for (let i = 0; i < legs.length; i++) {
    if (want <= legs[i] || i === legs.length - 1) {
      const f = legs[i] === 0 ? 0 : Math.min(1, want / legs[i]);
      return between(path[i], path[i + 1], f);
    }
    want -= legs[i];
  }
  return path[path.length - 1];
}

// The lamps stand at the corners of every block, which is where a city puts
// them and which keeps them square to the grid.
export function lampPosts({cols, rows}: {cols: number; rows: number}): Vec[] {
  const out: Vec[] = [];
  for (let col = 0; col <= cols; col++) {
    for (let row = 0; row <= rows; row++) {
      out.push({x: col * BLOCK + ROAD / 2 + .1, y: row * BLOCK + ROAD / 2 + .1});
    }
  }
  return out;
}

// Where the cars stand: along the kerb of a block, parked rather than driving.
// The clock is stopped between actions, and a car moving while time is not
// would be the view inventing something the simulation has not spent.
export function kerbside(cell: Cell, seed: number): {at: Vec; horizontal: boolean; facing: number}[] {
  const i = island(cell);
  let h = ((cell.col * 73856093) ^ (cell.row * 19349663) ^ seed) >>> 0;
  const next = () => { h = (h * 1664525 + 1013904223) >>> 0; return h / 4294967296 };
  const out: {at: Vec; horizontal: boolean; facing: number}[] = [];
  const slots = 3;
  for (let n = 0; n < slots; n++) {
    if (next() > .5) continue;
    const f = (n + .5) / slots;
    if (next() < .5) {
      out.push({at: {x: i.x + i.w * f, y: i.y - ROAD * .28}, horizontal: true, facing: next() < .5 ? 1 : -1});
    } else {
      out.push({at: {x: i.x - ROAD * .28, y: i.y + i.d * f}, horizontal: false, facing: next() < .5 ? 1 : -1});
    }
  }
  return out;
}

// ---------------------------------------------------------------------------
// The light.
//
// The city was lit the same at noon and at three in the morning, which made it
// a diagram. The clock is the core's, so the light is the core's too: how dark
// it is comes from the hour the world says it is, and everything else — how
// much the windows spill, whether the lamps are burning, how far the haze
// reaches — follows from that one number.

// nightness runs 0 at midday to 1 in the small hours, with the turn happening
// over the hour or so that dusk actually takes.
export function nightness(minute: number): number {
  const hour = (minute % 1440) / 60;
  if (hour < 4.5 || hour >= 20.5) return 1;
  if (hour < 6.5) return (6.5 - hour) / 2;
  if (hour > 18.5) return (hour - 18.5) / 2;
  return 0;
}

// mix blends two packed colours, which is how the ground goes from grey stone
// at noon to blue-black at night without a second palette.
export function mix(a: number, b: number, t: number): number {
  const f = Math.min(1, Math.max(0, t));
  const r = Math.round(((a >> 16) & 255) * (1 - f) + ((b >> 16) & 255) * f);
  const g = Math.round(((a >> 8) & 255) * (1 - f) + ((b >> 8) & 255) * f);
  const c = Math.round((a & 255) * (1 - f) + (b & 255) * f);
  return (r << 16) | (g << 8) | c;
}

// How far back in the picture something is, as 0 (nearest) to 1 (furthest).
// Used to lay haze over the far blocks: in a city at night the next street is
// clear and the one after it is a suggestion.
export function distance(at: Vec, extent: {cols: number; rows: number}): number {
  const deepest = (extent.cols + extent.rows) * BLOCK;
  return deepest <= 0 ? 0 : 1 - Math.min(1, (at.x + at.y) / deepest);
}

// A filler's shape, decided by where it stands so the same block is the same
// building every time the city is drawn.
export function fillerShape(cell: Cell): {h: number; inset: number; kind: number; art: number} {
  let h = ((cell.col * 374761393) ^ (cell.row * 668265263)) >>> 0;
  h = (h ^ (h >> 13)) * 1274126177 >>> 0;
  const r = (n: number) => ((h >> n) & 255) / 255;
  // Shorter and set further back than the addresses that matter. A filler the
  // same size as the Monarch competes with it, and the first version — a plain
  // grey box filling its whole plot — read as a mistake rather than as
  // distance.
  return {h: .55 + r(3) * .85, inset: .34 + r(11) * .2, kind: Math.floor(r(19) * 3), art: r(7)};
}

// ---------------------------------------------------------------------------
// What is on the pavement.
//
// A street with nothing on it but lamps reads as a model of a street. These are
// the things a 1950s pavement actually carried: a hydrant at the kerb, a
// mailbox on the corner, a bin, a bench, a telegraph pole with wires running
// off it. They are placed on the pavement ring only — never in the road, never
// under a building — so nothing can end up somewhere it could not stand.

export type Prop = {kind: 'hydrant' | 'mailbox' | 'bin' | 'bench' | 'pole'; at: Vec; facing: number};

// dressing is what one block carries, decided by where the block is so the same
// corner has the same hydrant every time the city is drawn.
export function dressing(cell: Cell): Prop[] {
  const i = island(cell);
  let h = ((cell.col * 2246822519) ^ (cell.row * 3266489917)) >>> 0;
  const next = () => { h = (h * 1664525 + 1013904223) >>> 0; return h / 4294967296 };
  const out: Prop[] = [];

  // The pavement ring, as four runs a prop can stand on. Kept a little inside
  // the kerb so nothing overhangs the carriageway.
  const inset = PAVE * .42;
  const runs: {from: Vec; to: Vec; facing: number}[] = [
    {from: {x: i.x + inset, y: i.y + inset}, to: {x: i.x + i.w - inset, y: i.y + inset}, facing: -1},
    {from: {x: i.x + inset, y: i.y + i.d - inset}, to: {x: i.x + i.w - inset, y: i.y + i.d - inset}, facing: 1},
    {from: {x: i.x + inset, y: i.y + inset}, to: {x: i.x + inset, y: i.y + i.d - inset}, facing: -1},
    {from: {x: i.x + i.w - inset, y: i.y + inset}, to: {x: i.x + i.w - inset, y: i.y + i.d - inset}, facing: 1},
  ];

  const kinds: Prop['kind'][] = ['hydrant', 'mailbox', 'bin', 'bench'];
  for (const run of runs) {
    // Most stretches of pavement carry nothing. A prop on every one reads as
    // a catalogue rather than as a street.
    if (next() > .62) continue;
    const at = between(run.from, run.to, .2 + next() * .6);
    out.push({kind: kinds[Math.floor(next() * kinds.length)], at, facing: run.facing});
  }

  // A telegraph pole on one corner of some blocks, which is what the wires
  // hang from.
  if (next() < .55) {
    out.push({kind: 'pole', at: {x: i.x + inset * .6, y: i.y + i.d - inset * .6}, facing: 1});
  }
  return out;
}

// The wires: strung between the poles of neighbouring blocks, down the line of
// the street rather than across the city at random.
export function wires(poles: Vec[]): Segment[] {
  const out: Segment[] = [];
  for (const a of poles) {
    for (const b of poles) {
      if (a === b) continue;
      const dx = Math.abs(a.x - b.x), dy = Math.abs(a.y - b.y);
      // One block apart, in line: that is a span. Anything else is not.
      const spanX = dy < .2 && dx > BLOCK - .4 && dx < BLOCK + .4;
      const spanY = dx < .2 && dy > BLOCK - .4 && dy < BLOCK + .4;
      if (!spanX && !spanY) continue;
      if (a.x > b.x || (a.x === b.x && a.y > b.y)) continue;   // once per pair
      out.push({a, b});
    }
  }
  return out;
}

// ---------------------------------------------------------------------------
// Terraces.
//
// One building standing in the middle of its block, with pavement on all four
// sides, is an office park. A city block is a terrace: buildings shoulder to
// shoulder along the frontage, sharing walls, with the yards behind them.
//
// So a block is a row of slots along its street frontage. The address that
// belongs to the block takes a slot and ordinary buildings take the rest, and
// they are drawn back to front so the party walls read as joins rather than as
// gaps.

export type Slot = {at: Vec; w: number; d: number; front: boolean; end: boolean};

// terrace lays a block's frontage out as slots. The front row faces the street
// that runs along the near edge; the back row fills the far edge, so a block
// reads as built-up rather than as one building with a lawn.
export function terrace(cell: Cell, slots = 2): Slot[] {
  const i = island(cell);
  const depth = (i.d - PAVE * 2) * .46;          // how far back a row reaches
  const width = (i.w - PAVE * 2) / slots;
  const out: Slot[] = [];
  for (let n = 0; n < slots; n++) {
    // The near row, along the street the camera looks down.
    out.push({
      at: {x: i.x + PAVE + n * width, y: i.y + i.d - PAVE - depth},
      w: width, d: depth, front: true, end: n === 0 || n === slots - 1,
    });
  }
  for (let n = 0; n < slots; n++) {
    // And the far row, backing onto it.
    out.push({
      at: {x: i.x + PAVE + n * width, y: i.y + PAVE},
      w: width, d: depth, front: false, end: n === 0 || n === slots - 1,
    });
  }
  return out;
}

// Which slot an address takes: the middle of the near row, so the building the
// player came to see faces the street and is never hidden behind another.
export const addressSlot = (_slots = 2) => 0;

// ---------------------------------------------------------------------------
// What is painted on the road.
//
// An empty carriageway with a dashed centre line is a diagram of a road. What
// makes it a street is the paint at the junctions: the bars of a crossing, the
// stop line a car waits behind. All of it is laid on the grid, so it lines up
// with the kerbs by construction rather than by being nudged.

export type Marking = {kind: 'crossing' | 'stop'; at: Vec; along: Vec; width: number};

// Crossings go across the carriageway on the approach to a junction, and stop
// lines sit just behind them. Both are placed from the block grid, so every
// junction in the city is marked the same way.
export function markings({cols, rows}: {cols: number; rows: number}): Marking[] {
  const out: Marking[] = [];
  const back = ROAD * .62;                    // how far from the centre of the junction
  for (let col = 0; col <= cols; col++) {
    for (let row = 0; row <= rows; row++) {
      const x = col * BLOCK, y = row * BLOCK;
      // Only where two carriageways actually meet.
      const hasEast = col < cols, hasSouth = row < rows;
      if (hasSouth) {
        // Crossing the north-south street, on the south side of the junction.
        out.push({kind: 'crossing', at: {x, y: y + back}, along: {x: 1, y: 0}, width: ROAD});
        out.push({kind: 'stop', at: {x, y: y + back + .16}, along: {x: 1, y: 0}, width: ROAD});
      }
      if (hasEast) {
        // And the east-west street, on the east side.
        out.push({kind: 'crossing', at: {x: x + back, y}, along: {x: 0, y: 1}, width: ROAD});
        out.push({kind: 'stop', at: {x: x + back + .16, y}, along: {x: 0, y: 1}, width: ROAD});
      }
    }
  }
  return out;
}

// ---------------------------------------------------------------------------
// Awnings.
//
// A shopfront in 1950 has canvas over it. It is the cheapest thing that turns a
// wall with windows in it into a place of business, and it is the one piece of
// a building that belongs to the street rather than to the block: it hangs out
// over the pavement, which is why it is drawn here from the terrace geometry
// rather than painted into the art.
//
// The rule is the same one the props follow: an awning may reach out over the
// pavement and it may never reach the road. Its projection is a fraction of
// PAVE, so that holds by construction rather than by being checked.

export type Awning = {
  at: Vec;          // the corner of the canopy against the wall
  w: number;        // how much of the frontage it covers
  reach: number;    // how far out over the pavement it hangs
  h: number;        // the height of its back edge, at the top of the shopfront
  drop: number;     // how much lower the front edge is, so rain runs off it
  stripes: number;  // bands of canvas across it
  tone: number;     // which of the awning colours this one is
};

// How far an awning may hang out over the pavement. Well inside the kerb: the
// props already sit at PAVE * .42 and nothing may overhang the carriageway.
const REACH = PAVE * .5;

// awnings gives the canopies on a block's street frontage. Only the front row
// has them — the back row faces the yards — and only some of the shopfronts,
// because a canopy on every one reads as a parade of market stalls.
export function awnings(cell: Cell, slots = 2): Awning[] {
  let h = ((cell.col * 374761393) ^ (cell.row * 668265263) ^ 0x5bf03635) >>> 0;
  const next = () => { h = (h * 1664525 + 1013904223) >>> 0; return h / 4294967296 };
  const out: Awning[] = [];
  for (const slot of terrace(cell, slots)) {
    if (!slot.front) continue;
    if (next() > .58) continue;
    // Not the whole frontage: a shopfront sits between the doorways, so the
    // canvas stops short of the party walls at either end.
    const w = slot.w * (.5 + next() * .22);
    const at = {x: slot.at.x + (slot.w - w) / 2, y: slot.at.y + slot.d};
    out.push({
      at, w, reach: REACH,
      h: .3 + next() * .06,
      drop: .05 + next() * .03,
      stripes: 4 + Math.floor(next() * 4),
      tone: Math.floor(next() * 4),
    });
  }
  return out;
}

// goldenness is the colour of the light rather than the amount of it.
//
// nightness alone made dawn a weaker night: the same blue-black ground, a
// little lighter. What separates six in the morning from three in the morning
// is not brightness, it is temperature — the light comes in low and warm and
// only for an hour or so at each end of the day. This runs 0 through the middle
// of the day and through the small hours, and 1 at the two turns.
export function goldenness(minute: number): number {
  const hour = (minute % 1440) / 60;
  const near = (peak: number) => Math.max(0, 1 - Math.abs(hour - peak) / 1.75);
  return Math.min(1, Math.max(near(6.2), near(19.2)));
}

// ---------------------------------------------------------------------------
// What the city gives off.
//
// A still city is a model of a city. It cannot be fixed with animation — the
// clock is stopped between actions and nothing may move on its own — but a
// chimney with smoke standing over it and a grate with steam coming off it are
// both *still* things in a photograph, and they are what tells you the place is
// occupied. So they are placed here, deterministically, and drawn as shapes
// rather than as particles.

export type Vent = {at: Vec; kind: 'chimney' | 'grate'; height: number; drift: number; size: number};

// vents gives what a block gives off: smoke from a chimney on the roof of one
// of its buildings, steam from a grate in the pavement outside. Never both from
// the same block, and most blocks give off nothing — a city where every roof
// smokes is a foundry.
export function vents(cell: Cell, slots = 2): Vent[] {
  let h = ((cell.col * 2654435761) ^ (cell.row * 40503) ^ 0x9e3779b9) >>> 0;
  const next = () => { h = (h * 1664525 + 1013904223) >>> 0; return h / 4294967296 };
  const out: Vent[] = [];
  const rows = terrace(cell, slots);
  if (next() < .42) {
    // A chimney, standing on the back of one of the roofs. Set back from the
    // street edge so the smoke rises behind the parapet rather than in front
    // of the building's own face.
    const slot = rows[Math.floor(next() * rows.length)];
    out.push({
      at: {x: slot.at.x + slot.w * (.3 + next() * .4), y: slot.at.y + slot.d * .3},
      kind: 'chimney',
      height: 1.5 + next() * .9,        // where the smoke starts, above the ground
      drift: (next() - .5) * .9,        // which way the wind has it
      size: .8 + next() * .5,
    });
  } else if (next() < .3) {
    // Or a grate in the pavement, which is the same idea at ankle height.
    const i = island(cell);
    out.push({
      at: {x: i.x + PAVE * .5 + next() * (i.w - PAVE), y: i.y + i.d - PAVE * .5},
      kind: 'grate', height: 0, drift: (next() - .5) * .5, size: .55 + next() * .4,
    });
  }
  return out;
}

// ---------------------------------------------------------------------------
// The trolley.
//
// One avenue in this city carries a streetcar, which is what a 1950s city has
// instead of the cars it does not have enough of yet. The rails are laid from
// the same grid as the carriageways, so they run down the centre of a street
// the whole length of the city and cross every junction square — there is no
// way for them to end up half on the pavement, because they are not placed,
// they are derived.

// TROLLEY_GAUGE is the distance between the two rails, in tiles. Narrow: a
// streetcar is not a mainline train, and at this scale two lines any further
// apart read as a road marking rather than as track.
export const TROLLEY_GAUGE = .22;

// trolleyAvenue is which north-south street carries it: the middle one, so the
// track runs through the city rather than along its edge.
export const trolleyAvenue = ({cols}: {cols: number}) => Math.max(1, Math.round(cols / 2));

// rails gives the two running rails, as a pair of segments down one avenue.
export function rails(size: {cols: number; rows: number}): Segment[] {
  const x = trolleyAvenue(size) * BLOCK, bottom = size.rows * BLOCK;
  return [
    {a: {x: x - TROLLEY_GAUGE / 2, y: -.5}, b: {x: x - TROLLEY_GAUGE / 2, y: bottom + .5}},
    {a: {x: x + TROLLEY_GAUGE / 2, y: -.5}, b: {x: x + TROLLEY_GAUGE / 2, y: bottom + .5}},
  ];
}

// sleepers gives the cross-ties showing through the setts between the rails.
// Spaced in tiles rather than per block, so they do not line up with the
// junctions and give the track a rhythm of its own.
export function sleepers(size: {cols: number; rows: number}): Segment[] {
  const x = trolleyAvenue(size) * BLOCK, bottom = size.rows * BLOCK;
  const out: Segment[] = [];
  for (let y = 0; y < bottom; y += .34) {
    out.push({a: {x: x - TROLLEY_GAUGE * .8, y}, b: {x: x + TROLLEY_GAUGE * .8, y}});
  }
  return out;
}
