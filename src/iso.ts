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
export function faces(origin: Vec, part: Part) {
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
