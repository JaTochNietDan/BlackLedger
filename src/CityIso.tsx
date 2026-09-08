import {useEffect, useRef} from 'react';
import {Application, Assets, Container, Graphics, Sprite, Text, Texture, TextStyle} from 'pixi.js';
import {Viewport} from 'pixi-viewport';
import type {Snapshot} from './types';
import {along, blockFor, BLOCK, bounds, carriageways, distance, faces, fillers, fillerShape, grid, island, kerbside, lampPosts, middle, mix, nightness, PAVE, plot, project, ROAD, size, TILE, walk} from './iso';
import type {Cell, Vec} from './iso';
import cutouts from '../public/art/iso/isometric.json';
import type {Spotlight} from './CityStreet';

// The city, drawn as a city, on a GPU.
//
// The first version of this was SVG, which was the right way to find out
// whether the projection and the placement were correct — and wrong for what
// comes next: a camera you can drag and zoom, every person in the city moving
// at once, and explosions. Those are what a sprite renderer is for.
//
// PixiJS owns the canvas; React owns everything around it. The bridge is
// deliberately one-way and imperative: the scene is rebuilt from the world when
// the world changes, and it never writes anything back except which address the
// player clicked. Nothing here decides anything — where a building stands and
// who is inside it still come from the core.

const COLOUR = (hex: string) => parseInt(hex.slice(1), 16);

// The painted cut-outs, by address. A building with one is drawn as itself; a
// building without one is drawn as the blocked-out solid it was before, so
// adding an address tomorrow puts a shape on the map rather than a hole.
const painted = new Map((cutouts as {id: string; file: string; w: number; h: number}[]).map(c => [c.id, c]));
// The ordinary buildings between the addresses, painted in the same light as
// the twelve so the city is one place. A block with one of these on it is not
// somewhere the player can go; it is somewhere that exists.
const FILL = [...painted.keys()].filter(id => id.startsWith('fill-')).sort();
const textures = new Map<string, Texture>();

// How far the camera may be pushed. Past these the city either fills the screen
// with one roof or shrinks into the middle of an empty field.
const ZOOM = {min: .45, max: 2.6};

// A car at the kerb, built the same way a building is: small boxes in tile
// space, drawn through the same projection. The first version was drawn by
// hand in screen space and read as a smear at any zoom, because nothing about
// it agreed with the angle everything else is at.
//
// The silhouette is what says 1950 — a long bonnet, an upright cabin set back,
// a short boot — not detail, which disappears at the size a city is drawn at.
function car(colour: number, along: boolean): Graphics {
  const g = new Graphics();
  const L = .62, W = .30;                     // length and width, in tiles
  const dark = shade(colour, .55), light = shade(colour, 1.45);

  // Laid out along x, then mirrored by the caller for the other street.
  const body = [
    {dx: 0, dy: 0, w: L, d: W, h: .11, base: .015, top: colour, left: dark, right: shade(colour, .8)},
    {dx: L * .26, dy: -.005, w: L * .42, d: W + .01, h: .10, base: .125,
     top: light, left: shade(colour, .5), right: shade(colour, .7)},
  ];
  for (const part of body) {
    const f = faces({x: -L / 2, y: -W / 2}, part);
    g.poly(f.left.flatMap(v => [v.x, v.y])).fill(part.left);
    g.poly(f.right.flatMap(v => [v.x, v.y])).fill(part.right);
    g.poly(f.top.flatMap(v => [v.x, v.y])).fill(part.top);
  }
  // Wheels: dark ellipses tucked under each end, which is what stops it
  // floating over the road.
  for (const at of [-L * .3, L * .28]) {
    const w = project({x: at, y: W / 2});
    g.ellipse(w.x, w.y + 1, 3.6, 1.9).fill(0x0d0f10);
  }
  // Headlights, and only a hint of what they throw.
  const nose = project({x: -L / 2, y: 0});
  g.circle(nose.x + 2, nose.y - 5, 1.4).fill({color: 0xf6e6bb, alpha: .9});
  g.poly([nose.x, nose.y - 5, nose.x - 16, nose.y - 9, nose.x - 16, nose.y + 3])
    .fill({color: 0xf0d6a0, alpha: .05});
  void along;
  return g;
}

// A colour lightened or darkened, so one paint job gives a whole car.
function shade(colour: number, by: number): number {
  const r = Math.min(255, Math.round(((colour >> 16) & 255) * by));
  const g = Math.min(255, Math.round(((colour >> 8) & 255) * by));
  const b = Math.min(255, Math.round((colour & 255) * by));
  return (r << 16) | (g << 8) | b;
}

// A person in the street, small enough to belong to a building and clear
// enough to be seen: a coat, a collar and a head. Deliberately not a portrait —
// at the scale a whole city is drawn at, a face is four pixels of mud, and this
// reads as somebody standing there.
function figure(colour: number, yours: boolean, walking: boolean): Graphics {
  const g = new Graphics();
  const h = 17;
  // The coat: narrow at the shoulders, flaring to the pavement.
  g.poly([-3.6, -h + 5, 3.6, -h + 5, 5.2, 0, -5.2, 0]).fill(colour);
  // A collar catching the light, so the figure has a front.
  g.poly([-3.6, -h + 5, 3.6, -h + 5, 2.4, -h + 8.5, -2.4, -h + 8.5]).fill(0xd8cfb4, .5);
  g.circle(0, -h + 2.6, 3.1).fill(0xc9a98a);          // the head
  // A hat, because it is 1950 and because without one a small figure reads as
  // a pin rather than as a man.
  g.ellipse(0, -h + 1.4, 5, 1.5).fill(0x1b1d1f);
  g.poly([-2.9, -h + 1.4, 2.9, -h + 1.4, 2.2, -h - 1.9, -2.2, -h - 1.9]).fill(0x24272a);
  g.ellipse(0, .6, 5.4, 1.7).fill({color: 0x000000, alpha: .32});  // and a shadow to stand in
  if (yours) g.circle(0, -h - 5.2, 1.7).fill(0xd6b77c);            // yours are marked
  if (walking) g.poly([-6.6, -1.2, -4.2, -1.2, -5.4, .8]).fill({color: 0xd6b77c, alpha: .5});
  return g;
}

// What a moment looks like over the building it happened at. Only light and
// shape: the city underneath is a painting, and drawing a body or a car on top
// of it reads as a sticker where a flash reads as something happening.
//
// t runs 0 to 1 across the moment. Nothing here decides anything — the core
// committed the event and named the address; this is what it looked like.
function moment(kind: string, t: number, across: number): Graphics {
  const g = new Graphics();
  const ease = 1 - t;
  switch (kind) {
    case 'explosion': {
      const r = t < .28 ? t * across * 2.6 : across * .73 - (t - .28) * across * .5;
      if (r > 0) {
        g.circle(0, 0, r).fill({color: 0xffcf7a, alpha: Math.max(0, .8 - t)});
        g.circle(0, 0, r * 1.55).fill({color: 0xc4531f, alpha: Math.max(0, .35 - t * .45)});
      }
      if (t > .3) {
        // Debris thrown out and falling: the only moment with anything solid
        // in it, because an explosion without pieces is a lamp.
        for (let i = 0; i < 11; i++) {
          const a = i * 2.1, fly = (t - .3) * across * 1.5;
          g.rect(Math.cos(a) * fly, Math.sin(a) * fly * .55 + (t - .3) * (t - .3) * 260, 5, 4)
            .fill({color: 0x14100d, alpha: ease});
        }
      }
      break;
    }
    case 'killing':
    case 'gunfight': {
      // Muzzle flashes, a few, close together, then nothing.
      if (t < .5 && Math.floor(t * 16) % 2 === 0) {
        const n = Math.floor(t * 16), x = (n % 3 - 1) * across * .12;
        g.circle(x, 0, across * .035).fill({color: 0xffe6a8, alpha: .95});
        g.circle(x, 0, across * .2).fill({color: 0xffe6a8, alpha: .16});
      }
      break;
    }
    case 'raid':
    case 'arrest': {
      // A lamp turning over on a car at the kerb, sweeping the front.
      const swing = Math.sin(t * 26) * across * .3;
      g.poly([0, 0, swing + across * .34, across * .3, swing - across * .1, across * .3])
        .fill({color: 0xe0705c, alpha: .22});
      g.circle(0, 0, across * .05).fill({color: 0xe0705c, alpha: .55 + Math.sin(t * 26) * .35});
      break;
    }
    case 'seizure':
    case 'attack':
    case 'robbery':
    default: {
      // Something happened here: a hard ring that opens once and fades, which
      // is enough for a robbery and not so much that it reads as a fire.
      const r = across * (.18 + t * .5);
      g.circle(0, 0, r).stroke({width: 3, color: 0xd6b77c, alpha: Math.max(0, .7 - t)});
      break;
    }
  }
  return g;
}

export function CityIso({state, selected, onSelect, onEnter, spotlight}: {
  state: Snapshot; selected: string; onSelect: (id: string) => void; onEnter: () => void;
  spotlight?: Spotlight | null;
}) {
  const host = useRef<HTMLDivElement>(null);
  const app = useRef<Application | null>(null);
  const view = useRef<Viewport | null>(null);
  const blocks = useRef<Container | null>(null);
  const effects = useRef<Container | null>(null);
  // Where the camera was before a moment took it somewhere, so it can be given
  // back afterwards rather than leaving the player looking at a rooftop.
  const wasLooking = useRef<{x: number; y: number; scale: number} | null>(null);
  // The handlers change on every render; the scene is built once, so it reads
  // them through a box rather than closing over a stale one.
  const pick = useRef(onSelect); pick.current = onSelect;
  const enter = useRef(onEnter); enter.current = onEnter;

  // Build the renderer once. Rebuilding it on every state change would throw
  // away the camera the player had set, which is the whole point of having one.
  useEffect(() => {
    let dead = false;
    const application = new Application();
    application.init({
      backgroundAlpha: 0, antialias: true, autoDensity: true,
      resolution: Math.min(devicePixelRatio, 2), preference: 'webgl',
      resizeTo: host.current || undefined,
    }).then(() => {
      if (dead) { application.destroy(true); return }
      app.current = application;
      host.current?.append(application.canvas);

      const viewport = new Viewport({
        screenWidth: host.current?.clientWidth || 800,
        screenHeight: host.current?.clientHeight || 560,
        worldWidth: 2600, worldHeight: 1800,
        events: application.renderer.events,
      });
      viewport.drag().pinch().wheel({smooth: 3}).decelerate({friction: .92})
        .clampZoom({minScale: ZOOM.min, maxScale: ZOOM.max});
      application.stage.addChild(viewport);
      view.current = viewport;

      const layer = new Container();
      viewport.addChild(layer);
      blocks.current = layer;
      const above = new Container();
      viewport.addChild(above);
      effects.current = above;
      draw();
      frame();
      // Then the painted city, once the pictures are in. Until they land the
      // blocked-out solids stand in, which is why the first draw happens above
      // rather than waiting on a network.
      Promise.all([...painted.values()].map(async c => {
        try { textures.set(c.id, await Assets.load('/art/' + c.file)) } catch { /* it keeps its block */ }
      })).then(() => { if (!dead) { draw(); frame() } });
    }).catch(() => { /* the card view is still there; see TestTheCardViewIsStillReachable */ });

    const onResize = () => frame();
    window.addEventListener('resize', onResize);

    return () => {
      dead = true;
      window.removeEventListener('resize', onResize);
      app.current?.destroy(true, {children: true});
      app.current = null; view.current = null; blocks.current = null;
    };
  }, []);

  // Redraw the city whenever the world moves. The camera is untouched: a player
  // who has zoomed in on the docks does not want to be thrown back out because
  // an hour passed.
  // Deliberately not spotlight.t: the moment animates in its own layer, and
  // rebuilding twelve buildings sixty times a second to redraw a fireball that
  // is not in them is work for nothing.
  useEffect(draw, [state.revision, state.minute, selected, spotlight?.id]);

  // When a moment starts, take the camera to the building it happened at and
  // remember where the player was looking so it can be given back.
  useEffect(() => {
    const viewport = view.current;
    if (!viewport) return;
    if (!spotlight) {
      const back = wasLooking.current;
      if (back) {
        wasLooking.current = null;
        viewport.animate({time: 520, position: {x: back.x, y: back.y}, scale: back.scale, ease: 'easeInOutSine'});
      }
      return;
    }
    const place = state.locations.find(l => l.id === spotlight.id);
    if (!place) return;
    if (!wasLooking.current) {
      wasLooking.current = {x: viewport.center.x, y: viewport.center.y, scale: viewport.scale.x};
    }
    const cell = grid(state.locations).get(place.id) || {col: 0, row: 0};
    const c = project(middle(cell));
    viewport.animate({
      time: 620, position: {x: c.x, y: c.y - 60},
      scale: Math.min(ZOOM.max, Math.max(1.25, viewport.scale.x)), ease: 'easeInOutSine',
    });
  }, [spotlight?.id]);

  // frame points the camera at the whole city, once, from the size of what was
  // actually drawn rather than from a guessed world rectangle. Called on the
  // first draw and when the window changes shape, never on an ordinary update.
  function frame() {
    const layer = blocks.current, viewport = view.current;
    if (!layer || !viewport) return;
    const b = layer.getLocalBounds();
    if (b.width <= 0 || b.height <= 0) return;
    const w = host.current?.clientWidth || 800, h = host.current?.clientHeight || 560;
    viewport.resize(w, h, b.width, b.height);
    const scale = Math.min(w / (b.width + 120), h / (b.height + 120));
    viewport.setZoom(Math.max(ZOOM.min, Math.min(ZOOM.max, scale)), true);
    viewport.moveCenter(b.x + b.width / 2, b.y + b.height / 2);
  }

  // The moment itself, redrawn every frame it is playing. Kept out of the city
  // layer so the whole city is not rebuilt sixty times a second for it.
  useEffect(() => {
    const above = effects.current;
    if (!above) return;
    above.removeChildren().forEach(c => c.destroy({children: true}));
    if (!spotlight) return;
    const place = state.locations.find(l => l.id === spotlight.id);
    if (!place) return;
    const cell = grid(state.locations).get(place.id) || {col: 0, row: 0};
    const ground = plot(cell);
    const across = (ground.w + ground.d) * (TILE.w / 2);
    const centre = project(middle(cell));
    const art = textures.get(place.id);
    // Over the roof when the building is painted, over the middle of the plot
    // when it is still a blocked-out solid.
    const up = art ? (art.height * (across / art.width)) * .62 : 60;
    const g = moment(spotlight.kind, spotlight.t, across);
    g.position.set(centre.x, centre.y - (up || 60));
    above.addChild(g);
  }, [spotlight?.id, spotlight?.kind, spotlight?.t]);

  function draw() {
    const layer = blocks.current;
    if (!layer) return;
    layer.removeChildren().forEach(c => c.destroy({children: true}));

    const here = state.player.location;
    const cells = grid(state.locations);
    const size = bounds(cells);
    // How dark it is, from the clock the core keeps. Everything below reads
    // this one number rather than deciding for itself what time it is.
    const dark = nightness(state.minute);

    // The ground: one slab under the whole city, so nothing floats and the
    // roads are cut out of something rather than laid on nothing.
    const earth = new Graphics();
    const far = {x: size.cols * BLOCK, y: size.rows * BLOCK};
    const corners = [{x: -.6, y: -.6}, {x: far.x + .6, y: -.6}, {x: far.x + .6, y: far.y + .6}, {x: -.6, y: far.y + .6}]
      .map(project);
    earth.poly(corners.flatMap(c => [c.x, c.y])).fill(mix(0x2a2f2c, 0x0f1416, dark));
    layer.addChild(earth);

    // The carriageways, full width and height, so every junction is square.
    const road = new Graphics();
    for (const way of carriageways(size)) {
      const horizontal = Math.abs(way.b.y - way.a.y) < .001;
      const pad = horizontal ? {x: 0, y: ROAD / 2} : {x: ROAD / 2, y: 0};
      const box = [
        {x: way.a.x - pad.x, y: way.a.y - pad.y}, {x: way.b.x + pad.x, y: way.a.y - pad.y},
        {x: way.b.x + pad.x, y: way.b.y + pad.y}, {x: way.a.x - pad.x, y: way.b.y + pad.y},
      ].map(project);
      road.poly(box.flatMap(c => [c.x, c.y])).fill(mix(0x333a38, 0x161b1c, dark));
    }
    layer.addChild(road);

    // The pavements: the footway inside each block, between kerb and wall,
    // drawn as the island with the building's plot cut out of the middle.
    const pave = new Graphics();
    const kerb = new Graphics();
    for (const cell of cells.values()) {
      const i = island(cell);
      const outer = [{x: i.x, y: i.y}, {x: i.x + i.w, y: i.y}, {x: i.x + i.w, y: i.y + i.d}, {x: i.x, y: i.y + i.d}]
        .map(project);
      pave.poly(outer.flatMap(c => [c.x, c.y])).fill(mix(0x4a514c, 0x252b29, dark));
      kerb.poly(outer.flatMap(c => [c.x, c.y])).stroke({width: 1.6, color: mix(0x5d675f, 0x323b36, dark), alpha: .95});
      // The join between pavement and building, a shade darker so the plot
      // reads as ground the building sits on rather than as more pavement.
      const b = plot(cell);
      const inner = [{x: b.x, y: b.y}, {x: b.x + b.w, y: b.y}, {x: b.x + b.w, y: b.y + b.d}, {x: b.x, y: b.y + b.d}]
        .map(project);
      pave.poly(inner.flatMap(c => [c.x, c.y])).fill(mix(0x3e443f, 0x1e2422, dark));
    }
    layer.addChild(pave, kerb);

    // A broken line down the middle of every carriageway.
    const paint = new Graphics();
    for (const way of carriageways(size)) {
      const length = Math.hypot(way.b.x - way.a.x, way.b.y - way.a.y);
      const dashes = Math.max(2, Math.round(length / .5));
      for (let i = 0; i < dashes; i += 2) {
        const from = project({x: way.a.x + (way.b.x - way.a.x) * (i / dashes), y: way.a.y + (way.b.y - way.a.y) * (i / dashes)});
        const to = project({x: way.a.x + (way.b.x - way.a.x) * ((i + .6) / dashes), y: way.a.y + (way.b.y - way.a.y) * ((i + .6) / dashes)});
        paint.moveTo(from.x, from.y).lineTo(to.x, to.y);
      }
    }
    paint.stroke({width: 1.3, color: 0x6d6a52, alpha: .35});
    layer.addChild(paint);

    // The lamps, at every corner of every block, and the pools they throw.
    const glow = new Graphics();
    const posts = new Graphics();
    for (const foot of lampPosts(size)) {
      const p = project(foot);
      // A lamp burning at noon is the surest sign nothing is looking at the
      // clock, so the pools it throws come up as the light goes down.
      glow.ellipse(p.x, p.y, TILE.w * .40, TILE.h * .40).fill({color: 0xd9b678, alpha: .11 * dark});
      glow.ellipse(p.x, p.y, TILE.w * .21, TILE.h * .21).fill({color: 0xf0d6a0, alpha: .10 * dark});
      const H = 32;
      posts.poly([p.x - 1.4, p.y, p.x + 1.4, p.y, p.x + .8, p.y - H, p.x - .8, p.y - H]).fill(0x1b1f21);
      posts.ellipse(p.x, p.y, 3.2, 1.3).fill(0x14171a);
      posts.rect(p.x - .8, p.y - H - 1, 5, 1.3).fill(0x1b1f21);
      const lx = p.x + 4.4, ly = p.y - H + 1;
      posts.poly([lx - 2.2, ly, lx + 2.2, ly, lx + 1.4, ly + 4.8, lx - 1.4, ly + 4.8])
        .fill({color: mix(0x8b8778, 0xf3dcae, dark), alpha: .3 + .62 * dark});
      posts.poly([lx - 2.6, ly - 1.3, lx + 2.6, ly - 1.3, lx + 2.2, ly, lx - 2.2, ly]).fill(0x22262a);
      glow.poly([lx, ly + 4, lx + 12, p.y + 3, lx - 12, p.y + 3]).fill({color: 0xf0d6a0, alpha: .07 * dark});
    }
    layer.addChild(glow, posts);

    // The blocks nobody lives on. Without these the grid is a scatter of
    // twelve models with holes between them; with them it is a city that
    // happens to have twelve addresses worth knowing. Deliberately plainer
    // than anything the player can walk into: no name, no plot, no click.
    const filler = new Container();
    for (const cell of fillers(cells, size)) {
      const shape = fillerShape(cell);
      const ground = plot(cell);
      const inset = shape.inset;
      const at = {x: ground.x + inset, y: ground.y + inset};
      const w = ground.w - inset * 2, d = ground.d - inset * 2;
      const wall = mix(0x555c53, 0x232a2a, dark);
      const parts = [
        {w, d, h: shape.h, top: mix(0x6a7168, 0x2c3433, dark), left: mix(0x3c423c, 0x171d1e, dark), right: wall},
        // A parapet, a water tank or a stair head, so the roofline varies.
        shape.kind === 0
          ? {dx: w * .18, dy: d * .18, w: w * .3, d: d * .3, h: .28, base: shape.h,
             top: mix(0x5c6359, 0x252c2b, dark), left: mix(0x343a35, 0x141a1b, dark), right: mix(0x474e46, 0x1d2424, dark)}
          : shape.kind === 1
          ? {dx: -.04, dy: -.04, w: w + .08, d: d + .08, h: .1, base: shape.h,
             top: mix(0x4a514a, 0x1f2626, dark), left: mix(0x2f352f, 0x121819, dark), right: mix(0x3b423b, 0x191f20, dark)}
          : {dx: w * .62, dy: d * .2, w: w * .22, d: d * .22, h: .5, base: shape.h,
             top: mix(0x585f55, 0x232a2a, dark), left: mix(0x32382f, 0x131919, dark), right: mix(0x424940, 0x1b2222, dark)},
      ];
      // A painted filler where one has loaded; the drawn solid otherwise, so
      // the block is never empty while the pictures are still arriving.
      const which = FILL.length ? FILL[Math.floor(shape.art * FILL.length) % FILL.length] : '';
      const fillArt = which ? textures.get(which) : undefined;
      const g = new Graphics();
      if (fillArt) {
        const across = (ground.w + ground.d) * (TILE.w / 2) * .82;
        const art = new Sprite(fillArt);
        art.anchor.set(.5, 1);
        art.scale.set(across / fillArt.width);
        const foot = project({x: ground.x + ground.w, y: ground.y + ground.d});
        const mid = project({x: ground.x + ground.w / 2, y: ground.y + ground.d / 2});
        art.position.set(mid.x, foot.y);
        // Held back from the twelve: darker and a little cooler, so an address
        // the player can walk into always reads first.
        art.tint = mix(0xbfc4bd, 0x6f7a80, .35 + dark * .3);
        g.addChild(art);
      } else {
        for (const part of parts) {
          const f = faces(at, part);
          g.poly(f.left.flatMap(v => [v.x, v.y])).fill(part.left);
          g.poly(f.right.flatMap(v => [v.x, v.y])).fill(part.right);
          g.poly(f.top.flatMap(v => [v.x, v.y])).fill(part.top);
        }
      }
      // Windows in rows down both visible faces, some of them lit. Without
      // these a filler is a slab; with them it is a building somebody lives
      // in, which is the whole point of putting it there.
      const rows = fillArt ? 0 : Math.max(2, Math.round(shape.h * 3.4));
      const front = project({x: at.x, y: at.y + d});
      const rightEdge = project({x: at.x + w, y: at.y + d});
      const leftEdge = project({x: at.x, y: at.y});
      for (let row = 0; row < rows; row++) {
        const up = (shape.h * (row + .65) / rows) * TILE.h;
        for (let col = 0; col < 3; col++) {
          const f = (col + .5) / 3;
          const seed = (cell.col * 31 + cell.row * 17 + row * 7 + col * 13) % 9;
          const on = dark > .15 && seed % 3 !== 0;
          const paneR = {x: front.x + (rightEdge.x - front.x) * f, y: front.y + (rightEdge.y - front.y) * f - up};
          const paneL = {x: front.x + (leftEdge.x - front.x) * f, y: front.y + (leftEdge.y - front.y) * f - up};
          const glassOn = mix(0xb9c3b6, 0xe8c184, dark);
          const glassOff = mix(0x39423d, 0x151b1c, dark);
          g.poly([paneR.x - 2, paneR.y - 3.4, paneR.x + 2, paneR.y - 2.2, paneR.x + 2, paneR.y + 1.6, paneR.x - 2, paneR.y + .4])
            .fill({color: on ? glassOn : glassOff, alpha: on ? .5 + .35 * dark : .8});
          g.poly([paneL.x - 2, paneL.y - 2.2, paneL.x + 2, paneL.y - 3.4, paneL.x + 2, paneL.y + .4, paneL.x - 2, paneL.y + 1.6])
            .fill({color: on ? glassOn : glassOff, alpha: on ? .42 + .3 * dark : .75});
        }
      }
      // Haze: the far side of the city is a suggestion, the near side is not.
      const away = distance({x: at.x, y: at.y}, size);
      g.alpha = 1 - away * .35 * (0.35 + dark * .65);
      filler.addChild(g);
    }
    layer.addChild(filler);

    // The traffic, such as it is: cars at the kerb, in the drab colours a
    // 1950s street actually held rather than a paintbox.
    const PAINT = [0x2b3038, 0x3a3129, 0x27302c, 0x40342c, 0x1f2429, 0x4a3b2e];
    const cars = new Container();
    let colour = 0;
    for (const cell of cells.values()) {
      for (const spot of kerbside(cell, 91)) {
        const p = project(spot.at);
        const c = car(PAINT[colour++ % PAINT.length], spot.horizontal);
        c.position.set(p.x, p.y);
        // A car sits along the kerb it is parked at; in this projection that
        // means mirroring the long axis for a street running the other way.
        if (!spot.horizontal) c.scale.x = -1;
        cars.addChild(c);
      }
    }
    layer.addChild(cars);
    const placed = state.locations.map(p => {
      const cell = cells.get(p.id) || {col: 0, row: 0};
      const ground = plot(cell);
      // The block decides the footprint now, not the kind of building: every
      // plot on the grid is the same size, which is what keeps the streets
      // straight. What a place is still decides how it is drawn.
      const block = {...blockFor(p.type), w: ground.w, d: ground.d};
      const at = {x: ground.x, y: ground.y};
      return {p, cell, block, at, d: at.x + block.w / 2 + at.y + block.d / 2};
    }).sort((a, b) => a.d - b.d);

    for (const {p, cell, block, at} of placed) {
      const lit = spotlight?.id === p.id;
      const shut = p.district > state.district;
      const group = new Container();
      group.eventMode = 'static';
      group.cursor = 'pointer';
      group.on('pointertap', () => pick.current(p.id));
      // Hover lifts a building out of the haze; leaving puts it back where it
      // was rather than at full brightness, which would have left every
      // building the mouse crossed permanently nearer than the rest.
      group.on('pointerover', () => { group.alpha = Math.min(1, group.alpha + .25) });
      group.on('pointerout', () => { group.alpha = resting });
      if (p.id === here) group.on('pointertap', () => { if (p.id === here) enter.current() });
      const away = distance(at, size);
      const resting = (shut ? .45 : 1) * (1 - away * .22 * (0.3 + dark * .7));
      group.alpha = resting;

      // The ground it stands on, so nothing floats.
      const ground = new Graphics();
      const corners = [
        project(at), project({x: at.x + block.w, y: at.y}),
        project({x: at.x + block.w, y: at.y + block.d}), project({x: at.x, y: at.y + block.d}),
      ];
      // What the windows throw on the pavement. A painted building is full of
      // lit windows and stood in a pool of nothing; this is the light getting
      // out of it.
      if (dark > .1) {
        const centreNow = project(middle(cell));
        const spill = new Graphics();
        spill.ellipse(centreNow.x, centreNow.y + TILE.h * .35, TILE.w * .78, TILE.h * .78)
          .fill({color: 0xe8bf82, alpha: .07 * dark});
        spill.ellipse(centreNow.x, centreNow.y + TILE.h * .3, TILE.w * .45, TILE.h * .45)
          .fill({color: 0xf2d3a0, alpha: .06 * dark});
        group.addChild(spill);
      }

      const hasArt = textures.has(p.id);
      const marked = lit || p.id === selected || p.id === here || p.owned;
      // A painted building brings its own pavement with it, so the drawn plot
      // is only wanted under a blocked-out one — or under any of them when
      // there is something to say about the ground.
      if (!hasArt || marked) {
        ground.poly(corners.flatMap(c => [c.x, c.y]))
          .fill(hasArt ? {color: lit ? 0xd6b77c : 0x000000, alpha: lit ? .12 : 0} : lit ? 0x2a2f23 : 0x1b2422)
          .stroke({width: p.id === selected ? 2.5 : 1.5, color: lit || p.id === selected ? 0xd6b77c : p.owned ? 0x8f7849 : 0x232e2b});
        group.addChild(ground);
      }

      const centre = project({x: at.x + block.w / 2, y: at.y + block.d / 2});
      // How wide this plot is on screen, which is what a cut-out has to match:
      // the picture is scaled to the ground it stands on, never to itself.
      const across = (block.w + block.d) * (TILE.w / 2);
      let tallest = Math.max(...block.parts.map(q => (q.base || 0) + q.h));

      const texture = textures.get(p.id);
      if (texture) {
        const art = new Sprite(texture);
        art.anchor.set(.5, 1);
        art.scale.set(across / texture.width);
        // Its feet go on the front corner of the plot, where the near edges
        // meet, so the building stands on its own ground rather than floating
        // over the middle of it.
        const foot = project({x: at.x + block.w, y: at.y + block.d});
        art.position.set(centre.x, foot.y);
        group.addChild(art);
        tallest = (art.height / TILE.h) * .6;
      } else {
        // The building, one box at a time, so a piece of it can be replaced.
        for (const part of block.parts) {
          const f = faces(at, part);
          const g = new Graphics();
          g.poly(f.left.flatMap(v => [v.x, v.y])).fill(COLOUR(part.left));
          g.poly(f.right.flatMap(v => [v.x, v.y])).fill(COLOUR(part.right));
          g.poly(f.top.flatMap(v => [v.x, v.y])).fill(COLOUR(part.top));
          group.addChild(g);
        }
      }
      const name = new Text({
        text: p.name,
        style: new TextStyle({
          fontFamily: 'system-ui, sans-serif', fontSize: 13,
          fill: p.id === selected || p.id === here ? 0xd6b77c : 0xcfd4c6,
          stroke: {color: 0x0d1413, width: 4, join: 'round'},
        }),
      });
      name.anchor.set(.5, 1);
      name.position.set(centre.x, centre.y - tallest * TILE.h - 10);
      name.resolution = 2;
      group.addChild(name);

      // Whoever is standing here, on the pavement in front of the building.
      const crowd = p.people || [];
      const island_ = island(cell);
      crowd.slice(0, 12).forEach((who, i) => {
        // Along the pavement in front of the building, and into a second row
        // when the first is full.
        const of = Math.min(crowd.length, 12), per = Math.min(of, 5);
        const across_ = per <= 1 ? .5 : .12 + (i % per) / (per - 1) * .76;
        const rank = Math.floor(i / per);
        const spot = project({
          x: island_.x + island_.w * across_,
          y: island_.y + island_.d - PAVE * (.42 + rank * .55),
        });
        const g = figure(who.yours ? 0x4a4432 : 0x23262a, !!who.yours, false);
        g.position.set(spot.x, spot.y);
        g.eventMode = 'static';
        g.cursor = 'pointer';
        g.on('pointertap', () => pick.current(p.id));
        group.addChild(g);
      });

      layer.addChild(group);
    }

    // And the people crossing the city, drawn last so they pass in front of
    // the buildings they are walking between. Where they are comes from the
    // core: it knows the two addresses and how far along the walk they are.
    for (const j of state.street || []) {
      const from = cells.get(j.from_id), to = cells.get(j.to_id);
      if (!from || !to) continue;
      // Along the streets, not through the buildings.
      const spot = project(along(walk(from, to), j.progress));
      const g = figure(j.yours ? 0x4a4432 : 0x23262a, !!j.yours, true);
      g.position.set(spot.x, spot.y);
      layer.addChild(g);

      const tag = new Text({
        text: j.name.split(' ')[0],
        style: new TextStyle({
          fontFamily: 'system-ui, sans-serif', fontSize: 10,
          fill: j.yours ? 0xd6b77c : 0xb9c0b2,
          stroke: {color: 0x0d1413, width: 3, join: 'round'},
        }),
      });
      tag.anchor.set(.5, 1);
      tag.position.set(spot.x, spot.y - 20);
      tag.resolution = 2;
      layer.addChild(tag);
    }
  }

  return <div className="city-iso" ref={host}>
    {/* The city has to be reachable without a mouse, and without WebGL. */}
    <div className="iso-reader" aria-label="City addresses">
      {state.locations.map(p => <button key={p.id} onClick={() => onSelect(p.id)}>
        {p.name}{p.district > state.district ? ' (not open to you yet)' : ''}
      </button>)}
    </div>
    <p className="iso-hint">Drag to move · scroll to zoom · double-click where you are standing to step inside</p>
  </div>;
}
