import {useEffect, useRef} from 'react';
import {Application, Assets, Container, Graphics, Sprite, Text, Texture, TextStyle} from 'pixi.js';
import {Viewport} from 'pixi-viewport';
import type {Snapshot} from './types';
import {addressSlot, along, awnings, blockFor, BLOCK, bounds, carriageways, distance, dressing, faces, fillerShape, goldenness, grid, island, kerbside, lampPosts, markings, middle, mix, nightness, PAVE, plot, project, ROAD, size, rails, sleepers, terrace, TILE, trolleyAvenue, TROLLEY_GAUGE, vents, walk, wires} from './iso';
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
type Cutout = {id: string; file: string; w: number; h: number; anchor?: number[]; base?: number};
const painted = new Map((cutouts as Cutout[]).map(c => [c.id, c]));

// Where a cut-out's building actually stands inside its own picture.
//
// The city used to assume all three of these: that the base is centred in the
// image, that the bottom edge of the image is the near corner of that base, and
// that the image is as wide as the base. None of them is true of a generated
// picture — the model leaves whatever air it likes around the building and the
// widest thing in the frame is usually a cornice — so good art came out
// misaligned: sunk into the pavement, shoved off its plot, or scaled to its own
// overhang. tools/fitiso.py measures the footprint and writes it into the
// manifest; this reads it, and falls back to the old assumption for any
// cut-out that has not been measured yet.
const feet = (id: string, texture: Texture) => {
  const cut = painted.get(id);
  return {
    anchor: cut?.anchor ?? [.5, 1],
    base: cut?.base ?? texture.width,
  };
};
// The ordinary buildings between the addresses, painted in the same light as
// the twelve so the city is one place. A block with one of these on it is not
// somewhere the player can go; it is somewhere that exists.
const FILL = [...painted.keys()].filter(id => id.startsWith('fill-')).sort();
// The corner blocks. These were asked for as terrace rows with flush ends and
// came back as buildings bent around a corner — which is not what was wanted
// and is exactly what a grid city needs at the end of a terrace, so they are
// used there rather than thrown away or pretended to be rows.
const CORNERS = [...painted.keys()].filter(id => id.startsWith('row-')).sort();
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

// The things a pavement carries. Each is a few boxes in tile space, through
// the same projection as everything else, because anything drawn in screen
// space stops agreeing with the angle the moment the camera moves.
function prop(kind: string, dark: number, facing: number): Graphics {
  const g = new Graphics();
  // A contact shadow first, under everything. Without one a prop reads as a
  // shape floating over the pavement rather than as an object standing on it —
  // the first hydrant looked like a red cube hanging in the street.
  if (kind !== 'pole') g.ellipse(0, 1, 5.5, 2).fill({color: 0x000000, alpha: .3 + .12 * dark});
  else g.ellipse(0, 1, 4, 1.6).fill({color: 0x000000, alpha: .35 + .12 * dark});
  const solid = (w: number, d: number, h: number, top: number, left: number, right: number, base = 0, dx = 0, dy = 0) => {
    const f = faces({x: -w / 2 + dx, y: -d / 2 + dy}, {w, d, h, base});
    g.poly(f.left.flatMap(v => [v.x, v.y])).fill(left);
    g.poly(f.right.flatMap(v => [v.x, v.y])).fill(right);
    g.poly(f.top.flatMap(v => [v.x, v.y])).fill(top);
  };
  switch (kind) {
    case 'hydrant': {
      // Oxide red, not pillar-box: everything else in this city is a muted
      // olive or umber, and a bright hydrant was the only saturated thing on
      // the street, which made it read as a bug rather than as a hydrant.
      const paint = mix(0x6e3a2c, 0x412219, dark);
      solid(.065, .065, .10, mix(0x84493a, 0x4e2a20, dark), mix(0x5e2a20, 0x3c1a14, dark), paint);
      solid(.10, .036, .022, paint, mix(0x5e2a20, 0x3c1a14, dark), paint, .07);   // the arms
      break;
    }
    case 'mailbox': {
      const paint = mix(0x2f4a3c, 0x1b2b23, dark);
      solid(.095, .08, .16, mix(0x3d5c4b, 0x22362c, dark), mix(0x1f3227, 0x121d17, dark), paint);
      break;
    }
    case 'bin': {
      solid(.08, .08, .11, mix(0x4a4b43, 0x24261f, dark), mix(0x2a2b25, 0x14150f, dark), mix(0x3a3b33, 0x1c1e18, dark));
      break;
    }
    case 'bench': {
      const wood = mix(0x5a4632, 0x2e2419, dark);
      solid(.25, .065, .03, wood, mix(0x33281c, 0x1a140e, dark), mix(0x453626, 0x231b13, dark), .045);
      solid(.25, .022, .065, wood, mix(0x33281c, 0x1a140e, dark), mix(0x453626, 0x231b13, dark), .075, 0, -.022);
      break;
    }
    case 'pole': {
      const timber = mix(0x4a3f31, 0x241f18, dark);
      solid(.07, .07, 1.35, mix(0x5c5040, 0x2c261e, dark), mix(0x2e271f, 0x171310, dark), timber);
      // The crossarm the wires run off.
      const top = project({x: 0, y: 0});
      g.rect(top.x - 11, top.y - 1.35 * TILE.h - 4, 22, 2).fill(timber);
      g.rect(top.x - 11, top.y - 1.35 * TILE.h - 12, 22, 2).fill(timber);
      break;
    }
  }
  void facing;
  return g;
}

// A person in the street, small enough to belong to a building and clear
// enough to be seen: a coat, a collar and a head. Deliberately not a portrait —
// at the scale a whole city is drawn at, a face is four pixels of mud, and this
// reads as somebody standing there.
function figure(colour: number, yours: boolean, walking: boolean, you = false): Graphics {
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
  // And the player is not one of yours: a ring on the ground under them, so
  // they can be found in a crowd without reading a name.
  if (you) {
    g.ellipse(0, 1, 8.5, 3).stroke({width: 1.6, color: 0xd6b77c, alpha: .85});
    g.circle(0, -h - 5.2, 2.2).fill(0xf0d6a0);
  }
  void walking;
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
    // And what colour that light is. Low sun at either end of the day, nothing
    // in the middle of it and nothing in the small hours, so dawn stops being
    // a weaker night and becomes its own hour.
    const gold = goldenness(state.minute);
    // Anything the light falls on flatly — ground, road, pavement — takes the
    // temperature of it. Buildings do not: they are painted, and warming them
    // as a whole washes the art out.
    // Stone takes the low sun; asphalt barely does, which is what keeps the
    // roads reading as roads at dawn instead of the whole city going one
    // flat brown.
    // And what the sky is doing, which comes from the core: the view is not
    // allowed to decide it is raining. Wet streets outlast the rain itself, so
    // this is a number rather than a flag.
    const wet = state.sky?.wet ?? 0;
    // Fog is the one weather that changes how far the player can see, so it is
    // the one that touches the depth haze rather than the ground.
    const murk = state.sky?.kind === 'fog' ? .42 : 0;
    const sunlit = (c: number) => mix(c, 0xa9713f, gold * .2);
    const tarmac = (c: number) => mix(mix(c, 0x141a20, wet * .38), 0x6b4f3c, gold * .09 * (1 - wet * .6));

    // The ground: one slab under the whole city, so nothing floats and the
    // roads are cut out of something rather than laid on nothing.
    const earth = new Graphics();
    const far = {x: size.cols * BLOCK, y: size.rows * BLOCK};
    const corners = [{x: -.6, y: -.6}, {x: far.x + .6, y: -.6}, {x: far.x + .6, y: far.y + .6}, {x: -.6, y: far.y + .6}]
      .map(project);
    earth.poly(corners.flatMap(c => [c.x, c.y])).fill(sunlit(mix(0x565c50, 0x0f1416, dark)));
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
      road.poly(box.flatMap(c => [c.x, c.y])).fill(tarmac(mix(0x4b514e, 0x161b1c, dark)));
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
      pave.poly(outer.flatMap(c => [c.x, c.y])).fill(sunlit(mix(0x7c8175, 0x252b29, dark)));
      kerb.poly(outer.flatMap(c => [c.x, c.y])).stroke({width: 1.6, color: mix(0x939a8b, 0x323b36, dark), alpha: .95});
      // The join between pavement and building, a shade darker so the plot
      // reads as ground the building sits on rather than as more pavement.
      const b = plot(cell);
      const inner = [{x: b.x, y: b.y}, {x: b.x + b.w, y: b.y}, {x: b.x + b.w, y: b.y + b.d}, {x: b.x, y: b.y + b.d}]
        .map(project);
      pave.poly(inner.flatMap(c => [c.x, c.y])).fill(sunlit(mix(0x6c7266, 0x1e2422, dark)));
    }
    layer.addChild(pave, kerb);

    // A broken line down the middle of every carriageway — except the one the
    // trolley runs down, where the track is what is down the middle.
    const paint = new Graphics();
    const avenue = trolleyAvenue(size) * BLOCK;
    for (const way of carriageways(size)) {
      if (Math.abs(way.a.x - avenue) < .001 && Math.abs(way.b.x - avenue) < .001) continue;
      const length = Math.hypot(way.b.x - way.a.x, way.b.y - way.a.y);
      const dashes = Math.max(2, Math.round(length / .5));
      for (let i = 0; i < dashes; i += 2) {
        const from = project({x: way.a.x + (way.b.x - way.a.x) * (i / dashes), y: way.a.y + (way.b.y - way.a.y) * (i / dashes)});
        const to = project({x: way.a.x + (way.b.x - way.a.x) * ((i + .6) / dashes), y: way.a.y + (way.b.y - way.a.y) * ((i + .6) / dashes)});
        paint.moveTo(from.x, from.y).lineTo(to.x, to.y);
      }
    }
    paint.stroke({width: 1.3, color: 0x6d6a52, alpha: .35});

    // The trolley track: two running rails down one avenue, with the ties
    // showing through the setts between them. Laid from the same grid as the
    // carriageways, which is why it crosses every junction square and cannot
    // end up half on the pavement.
    const track = new Graphics();
    for (const tie of sleepers(size)) {
      const a = project(tie.a), b = project(tie.b);
      track.moveTo(a.x, a.y).lineTo(b.x, b.y);
    }
    track.stroke({width: 1.4, color: mix(0x4d4a3f, 0x22231f, dark), alpha: .55});
    for (const rail of rails(size)) {
      const a = project(rail.a), b = project(rail.b);
      // The rail head is polished by use, so it catches whatever light there
      // is — the one thing in this street that is brighter at night.
      track.moveTo(a.x, a.y).lineTo(b.x, b.y);
    }
    track.stroke({width: 1.6, color: mix(0x9a9c93, 0xb4af9f, dark), alpha: .62 + dark * .26});
    layer.addChild(track);

    // The paint at the junctions: the bars of a crossing and the line a car
    // waits behind. Laid from the same grid as the kerbs, so it lines up by
    // construction rather than by being nudged into place.
    const road_paint = new Graphics();
    for (const mark of markings(size)) {
      const across = {x: -mark.along.y, y: mark.along.x};   // square to the street
      if (mark.kind === 'crossing') {
        const bars = 5;
        for (let i = 0; i < bars; i++) {
          const t = (i + .5) / bars - .5;
          const centre = {
            x: mark.at.x + across.x * t * mark.width * .82,
            y: mark.at.y + across.y * t * mark.width * .82,
          };
          const half = .085, long = .30;
          const corners = [
            {x: centre.x - across.x * half - mark.along.x * long, y: centre.y - across.y * half - mark.along.y * long},
            {x: centre.x + across.x * half - mark.along.x * long, y: centre.y + across.y * half - mark.along.y * long},
            {x: centre.x + across.x * half + mark.along.x * long, y: centre.y + across.y * half + mark.along.y * long},
            {x: centre.x - across.x * half + mark.along.x * long, y: centre.y - across.y * half + mark.along.y * long},
          ].map(project);
          road_paint.poly(corners.flatMap(c => [c.x, c.y]));
        }
      } else {
        const half = mark.width * .44, thick = .05;
        const corners = [
          {x: mark.at.x - across.x * half - mark.along.x * thick, y: mark.at.y - across.y * half - mark.along.y * thick},
          {x: mark.at.x + across.x * half - mark.along.x * thick, y: mark.at.y + across.y * half - mark.along.y * thick},
          {x: mark.at.x + across.x * half + mark.along.x * thick, y: mark.at.y + across.y * half + mark.along.y * thick},
          {x: mark.at.x - across.x * half + mark.along.x * thick, y: mark.at.y - across.y * half + mark.along.y * thick},
        ].map(project);
        road_paint.poly(corners.flatMap(c => [c.x, c.y]));
      }
    }
    // Worn paint, not fresh: it has been on the road a while.
    road_paint.fill({color: mix(0xd8d2b8, 0x9a957f, .35 + dark * .3), alpha: .34});
    layer.addChild(paint, road_paint);

    // The lamps, at every corner of every block, and the pools they throw.
    const glow = new Graphics();
    const posts = new Graphics();
    for (const foot of lampPosts(size)) {
      const p = project(foot);
      // A lamp burning at noon is the surest sign nothing is looking at the
      // clock, so the pools it throws come up as the light goes down.
      glow.ellipse(p.x, p.y, TILE.w * .40, TILE.h * .40).fill({color: 0xd9b678, alpha: .11 * dark});
      glow.ellipse(p.x, p.y, TILE.w * .21, TILE.h * .21).fill({color: 0xf0d6a0, alpha: .10 * dark});
      // On a wet road the lamp is twice: the pool it throws, and the smear of
      // itself lying in the water. A reflection stretches toward whoever is
      // looking at it, which in this projection is straight down the screen.
      if (wet > 0) {
        glow.ellipse(p.x, p.y + TILE.h * .55, TILE.w * .035, TILE.h * 1.05)
          .fill({color: 0xf0d6a0, alpha: .1 * wet * (.25 + dark * .75)});
        glow.ellipse(p.x, p.y + TILE.h * .3, TILE.w * .09, TILE.h * .5)
          .fill({color: 0xd9b678, alpha: .07 * wet * (.25 + dark * .75)});
      }
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

    // What the pavements carry, and the wires overhead.
    const clutter = new Container();
    const poles: {x: number; y: number}[] = [];
    for (let col = 0; col < size.cols; col++) {
      for (let row = 0; row < size.rows; row++) {
        for (const item of dressing({col, row})) {
          const p = project(item.at);
          const g = prop(item.kind, dark, item.facing);
          g.position.set(p.x, p.y);
          clutter.addChild(g);
          if (item.kind === 'pole') poles.push(item.at);
        }
      }
    }
    // The wires: strung block to block, sagging the way a wire does.
    const strung = new Graphics();
    for (const span of wires(poles)) {
      const a = project(span.a), b = project(span.b);
      const lift = 1.35 * TILE.h + 8;
      strung.moveTo(a.x, a.y - lift)
        .quadraticCurveTo((a.x + b.x) / 2, (a.y + b.y) / 2 - lift + 9, b.x, b.y - lift);
      strung.moveTo(a.x, a.y - lift + 8)
        .quadraticCurveTo((a.x + b.x) / 2, (a.y + b.y) / 2 - lift + 17, b.x, b.y - lift + 8);
    }
    strung.stroke({width: 1, color: mix(0x4a4f48, 0x15191a, dark), alpha: .75});
    layer.addChild(clutter, strung);

    // A canvas awning over a shopfront, hanging out over the pavement. Drawn
    // from the terrace geometry rather than painted into the art, because it
    // belongs to the street: it is the one part of a building that reaches
    // past its own wall.
    const CANVAS = [[0x2f4436, 0xd8cdb4], [0x6b2f2c, 0xd8cdb4],
                    [0x8a6a2c, 0xe0d5bb], [0x2c3d55, 0xd2c9b2]];
    const canopy = (a: ReturnType<typeof awnings>[number]) => {
      const g = new Graphics();
      const lift = (v: Vec, up: number) => ({x: v.x, y: v.y - up * TILE.h});
      const [band, pale] = CANVAS[a.tone % CANVAS.length];
      // The shadow it throws on the pavement, which is what stops it floating.
      const shade = [{x: a.at.x, y: a.at.y}, {x: a.at.x + a.w, y: a.at.y},
                     {x: a.at.x + a.w, y: a.at.y + a.reach}, {x: a.at.x, y: a.at.y + a.reach}]
        .map(v => project({x: v.x + .04, y: v.y + .04}));
      g.poly(shade.flatMap(v => [v.x, v.y])).fill({color: 0x0b0f10, alpha: .3 - dark * .16});
      // The canvas itself, in bands across the frontage, sloping down to the
      // street so the rain runs off it.
      const step = a.w / a.stripes;
      for (let k = 0; k < a.stripes; k++) {
        const x0 = a.at.x + k * step, x1 = x0 + step;
        const quad = [
          lift(project({x: x0, y: a.at.y}), a.h),
          lift(project({x: x1, y: a.at.y}), a.h),
          lift(project({x: x1, y: a.at.y + a.reach}), a.h - a.drop),
          lift(project({x: x0, y: a.at.y + a.reach}), a.h - a.drop),
        ];
        // Not darkened as far as a roof is: an awning sits under a lamp and
        // over a lit window, which is the whole reason a shopfront has one.
        g.poly(quad.flatMap(v => [v.x, v.y]))
          .fill(mix(k % 2 ? pale : band, 0x1d2124, .1 + dark * .34));
      }
      // The valance hanging off the front lip, which is what makes it read as
      // cloth rather than as a shelf.
      const fl = lift(project({x: a.at.x, y: a.at.y + a.reach}), a.h - a.drop);
      const fr = lift(project({x: a.at.x + a.w, y: a.at.y + a.reach}), a.h - a.drop);
      const hang = .055 * TILE.h;
      g.poly([fl.x, fl.y, fr.x, fr.y, fr.x, fr.y + hang, fl.x, fl.y + hang])
        .fill(mix(band, 0x0f1416, .22 + dark * .4));
      return g;
    };

    // Smoke off a chimney, steam off a grate. Still, not animated: the clock
    // is stopped between actions and nothing in this city moves on its own. A
    // plume in a photograph is a still shape, and it is the thing that says
    // somebody is in there.
    const plume = (v: ReturnType<typeof vents>[number]) => {
      const g = new Graphics();
      const foot = project(v.at);
      // Many small overlapping puffs, not a few big ones. Six evenly spaced
      // ellipses stack into a column of visible grey rings — a drill bit, not
      // smoke. What reads as smoke is enough of them that no single edge shows,
      // each one nudged off the centre line so the column is ragged.
      const puffs = v.kind === 'chimney' ? 20 : 12;
      const rise = v.kind === 'chimney' ? 1.9 : .7;
      let h = Math.abs(Math.round(v.at.x * 733 + v.at.y * 971)) >>> 0;
      const next = () => { h = (h * 1664525 + 1013904223) >>> 0; return h / 4294967296 };
      for (let k = 0; k < puffs; k++) {
        const t = (k + 1) / puffs;
        const up = (v.height + t * rise) * TILE.h;
        // Widening as it goes, and faster near the top where it is losing its
        // shape rather than holding a column.
        const wide = v.size * TILE.w * .05 * (1 + t * t * 3.4 + t);
        const wobble = (next() - .5) * wide * .55;
        g.ellipse(foot.x + v.drift * t * t * TILE.w * .5 + wobble,
                  foot.y - up + (next() - .5) * wide * .3, wide, wide * .66)
          .fill({color: mix(0xb9bdb8, 0x8b9296, dark), alpha: (1 - t) * (v.kind === 'chimney' ? .085 : .07)});
      }
      return g;
    };

    // Every slot in every block that the addresses do not take. A block is a
    // terrace: the address takes one frontage slot and ordinary buildings take
    // the rest, shoulder to shoulder, so the city is built up rather than
    // twelve models in twelve fields.
    const SLOTS = 2;
    const filler = new Container();
    const rows: {cell: Cell; slot: ReturnType<typeof terrace>[number]; depth: number}[] = [];
    const addressAt = new Map<string, string>();      // "col,row,index" -> id
    const hung = new Map<string, ReturnType<typeof awnings>>();
    const smoking = new Map<string, ReturnType<typeof vents>>();
    for (const [id, cell] of cells) addressAt.set(`${cell.col},${cell.row},${addressSlot(SLOTS)}`, id);
    for (let col = 0; col < size.cols; col++) {
      for (let row = 0; row < size.rows; row++) {
        hung.set(`${col},${row}`, awnings({col, row}, SLOTS));
        smoking.set(`${col},${row}`, vents({col, row}, SLOTS));
        terrace({col, row}, SLOTS).forEach((slot, index) => {
          const key = `${col},${row},${index}`;
          if (slot.front && addressAt.has(key)) return;   // the address builds here
          rows.push({cell: {col, row}, slot, depth: slot.at.x + slot.w / 2 + slot.at.y + slot.d / 2});
        });
      }
    }
    rows.sort((a, b) => a.depth - b.depth);
    // The frontmost slot of each block, which is where its plumes are hung.
    const nearest = new Map<string, number>();
    for (const r of rows) {
      const key = `${r.cell.col},${r.cell.row}`;
      nearest.set(key, Math.max(nearest.get(key) ?? -Infinity, r.depth));
    }
    for (const {cell, slot} of rows) {
      const shape = fillerShape({col: cell.col * 7 + Math.round(slot.at.x * 3), row: cell.row * 5 + Math.round(slot.at.y * 3)});
      // Which building this is. Stepping through the set by position rather
      // than picking at random stops the same picture landing next door to
      // itself, which with six buildings and sixty-odd slots it otherwise does
      // constantly and reads as wallpaper.
      const order = Math.round(slot.at.x * 4 + slot.at.y * 7 + cell.col + cell.row * 2);
      // A corner block at the ends of a row, an ordinary building in between —
      // which is how a street is actually built.
      const onEnd = slot.end;
      const set = onEnd && CORNERS.length ? CORNERS : FILL;
      const which = set.length ? set[((order % set.length) + set.length) % set.length] : '';
      const fillArt = which ? textures.get(which) : undefined;
      const g = new Graphics();
      if (fillArt) {
        // Scaled to the slot it fills, so neighbours meet at their walls.
        // No overshoot. It was added to close the party walls and its actual
        // effect was to drive every sprite into its neighbour: roofs cutting
        // through roofs and buildings hanging over the kerb. A hairline gap
        // between two buildings reads as two buildings; an overlap reads as
        // broken.
        const across = (slot.w + slot.d) * (TILE.w / 2);
        const art = new Sprite(fillArt);
        const fit = feet(which, fillArt);
        art.anchor.set(fit.anchor[0], fit.anchor[1]);
        art.scale.set(across / fit.base);
        const foot = project({x: slot.at.x + slot.w, y: slot.at.y + slot.d});
        art.position.set(foot.x, foot.y);
        // Tone varies between neighbours; proportions do not. Scaling the
        // height alone squashed and stretched buildings that were drawn
        // correctly, which is its own artefact on top of the overlapping.
        const warmth = ((order * 37) % 7) / 7;
        art.tint = mix(mix(0xd8d4c6, 0xc4ced2, warmth), 0x6f7a80, .1 + dark * .48);
        g.addChild(art);
      } else {
        const inset = .04;
        const at = {x: slot.at.x + inset, y: slot.at.y + inset};
        const w = slot.w - inset * 2, d = slot.d - inset * 2;
        const f = faces(at, {w, d, h: shape.h});
        g.poly(f.left.flatMap(v => [v.x, v.y])).fill(mix(0x3c423c, 0x171d1e, dark));
        g.poly(f.right.flatMap(v => [v.x, v.y])).fill(mix(0x555c53, 0x232a2a, dark));
        g.poly(f.top.flatMap(v => [v.x, v.y])).fill(mix(0x6a7168, 0x2c3433, dark));
      }
      const away = distance(slot.at, size);
      g.alpha = 1 - away * .35 * (0.16 + dark * .84) - away * murk;
      filler.addChild(g);
      // And the canvas over its shopfront, drawn straight after the building it
      // hangs on so it can never end up behind it.
      for (const a of hung.get(`${cell.col},${cell.row}`) || []) {
        if (a.at.x < slot.at.x - 1e-9 || a.at.x + a.w > slot.at.x + slot.w + 1e-9) continue;
        const shop = canopy(a);
        shop.alpha = g.alpha;
        filler.addChild(shop);
      }
      // The block's plumes go down after its nearest building, so smoke stands
      // over its own roofs and still passes behind anything in front of it.
      if (nearest.get(`${cell.col},${cell.row}`) === slot.at.x + slot.w / 2 + slot.at.y + slot.d / 2) {
        for (const v of smoking.get(`${cell.col},${cell.row}`) || []) filler.addChild(plume(v));
      }
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
      // The address stands in the middle of its block's frontage, in a slot the
      // width of its neighbours, so it belongs to the terrace rather than
      // sitting in a field of its own.
      const slot = terrace(cell, SLOTS)[addressSlot(SLOTS)];
      const block = {...blockFor(p.type), w: slot.w, d: slot.d};
      const at = {x: slot.at.x, y: slot.at.y};
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
      const resting = (shut ? .45 : 1) * Math.max(.15, 1 - away * .22 * (0.14 + dark * .86) - away * murk);
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
      // How wide this slot is on screen, which is what a cut-out has to match:
      // the picture is scaled to the ground it stands on, never to itself, and
      // to exactly that ground, so it cannot lean into its neighbour's.
      const across = (block.w + block.d) * (TILE.w / 2);
      let tallest = Math.max(...block.parts.map(q => (q.base || 0) + q.h));

      const texture = textures.get(p.id);
      if (texture) {
        const art = new Sprite(texture);
        const fit = feet(p.id, texture);
        art.anchor.set(fit.anchor[0], fit.anchor[1]);
        // Scaled so the building's own ground matches the ground it is given,
        // rather than so its picture matches the plot's width.
        art.scale.set(across / fit.base);
        // And stood on the near corner of the plot, which is the point the
        // anchor above names — the corner of the building's own base.
        const foot = project({x: at.x + block.w, y: at.y + block.d});
        art.position.set(foot.x, foot.y);
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
        // People stand in twos and threes, not in a line. The gap between one
        // knot and the next is what makes a pavement read as people rather
        // than as a row of pegs — an even spacing looked like a fence.
        const of = Math.min(crowd.length, 12), per = Math.min(of, 5);
        const knot = Math.floor(i / 2), inKnot = i % 2;
        const wobble = ((i * 2654435761) % 1000) / 1000;
        const across_ = Math.min(.9, Math.max(.08,
          .10 + (knot % 3) * .31 + inKnot * .055 + wobble * .05));
        const rank = Math.floor(i / per) + (wobble > .7 ? 1 : 0);
        const spot = project({
          x: island_.x + island_.w * across_,
          y: island_.y + island_.d - PAVE * (.36 + rank * .5 + wobble * .18),
        });
        const g = figure(who.yours ? 0x4a4432 : 0x23262a, !!who.yours, false);
        // Half of them turned the other way, so a knot of people looks like a
        // conversation rather than a queue.
        if (inKnot === 1) g.scale.x = -1;
        g.position.set(spot.x, spot.y);
        g.eventMode = 'static';
        g.cursor = 'pointer';
        g.on('pointertap', () => pick.current(p.id));
        group.addChild(g);
      });

      // And the player, standing at whatever address they are at. They are not
      // in the room's list — the core keeps them apart from the city's own
      // people — so without this the one figure that matters most is the only
      // one not on the map.
      if (p.id === here && state.player.alive) {
        const you = figure(0x2f3a2c, false, false, true);
        const at_ = project({
          x: island_.x + island_.w * .5,
          y: island_.y + island_.d - PAVE * .18,
        });
        you.position.set(at_.x, at_.y);
        group.addChild(you);
      }

      layer.addChild(group);
    }

    // And the people crossing the city, drawn last so they pass in front of
    // the buildings they are walking between. Where they are comes from the
    // core: it knows the two addresses and how far along the walk they are.
    for (const j of state.street || []) {
      const from = cells.get(j.from_id), to = cells.get(j.to_id);
      if (!from || !to) continue;
      // Along the streets, not through the buildings.
      const path = walk(from, to);
      const spot = project(along(path, j.progress));
      // Which way they are pointing: from where they were a moment ago to
      // where they are now, so a figure faces its own direction of travel.
      const behind = project(along(path, Math.max(0, j.progress - .04)));
      const g = figure(j.yours ? 0x4a4432 : 0x23262a, !!j.yours, true);
      if (spot.x < behind.x) g.scale.x = -1;
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
