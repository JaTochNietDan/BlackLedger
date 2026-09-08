import {useEffect, useRef} from 'react';
import {Application, Container, Graphics, Text, TextStyle} from 'pixi.js';
import {Viewport} from 'pixi-viewport';
import type {Snapshot} from './types';
import {blockFor, depth, faces, plan, project, size, TILE} from './iso';
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

// How far the camera may be pushed. Past these the city either fills the screen
// with one roof or shrinks into the middle of an empty field.
const ZOOM = {min: .45, max: 2.6};

export function CityIso({state, selected, onSelect, onEnter, spotlight}: {
  state: Snapshot; selected: string; onSelect: (id: string) => void; onEnter: () => void;
  spotlight?: Spotlight | null;
}) {
  const host = useRef<HTMLDivElement>(null);
  const app = useRef<Application | null>(null);
  const view = useRef<Viewport | null>(null);
  const blocks = useRef<Container | null>(null);
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
      draw();
      frame();
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
  useEffect(draw, [state.revision, state.minute, selected, spotlight?.id, spotlight?.t]);

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

  function draw() {
    const layer = blocks.current;
    if (!layer) return;
    layer.removeChildren().forEach(c => c.destroy({children: true}));

    const here = state.player.location;
    const placed = state.locations.map(p => {
      const block = blockFor(p.type);
      const at = plan(p.x, p.y);
      return {p, block, at, d: depth({x: at.x + block.w / 2, y: at.y + block.d / 2})};
    }).sort((a, b) => a.d - b.d);

    for (const {p, block, at} of placed) {
      const lit = spotlight?.id === p.id;
      const shut = p.district > state.district;
      const group = new Container();
      group.eventMode = 'static';
      group.cursor = 'pointer';
      group.on('pointertap', () => pick.current(p.id));
      group.on('pointerover', () => { group.alpha = 1 });
      group.on('pointerout', () => { group.alpha = shut ? .45 : 1 });
      if (p.id === here) group.on('pointertap', () => { if (p.id === here) enter.current() });
      group.alpha = shut ? .45 : 1;

      // The ground it stands on, so nothing floats.
      const plot = new Graphics();
      const corners = [
        project(at), project({x: at.x + block.w, y: at.y}),
        project({x: at.x + block.w, y: at.y + block.d}), project({x: at.x, y: at.y + block.d}),
      ];
      plot.poly(corners.flatMap(c => [c.x, c.y]))
        .fill(lit ? 0x2a2f23 : 0x1b2422)
        .stroke({width: p.id === selected ? 2.5 : 1.5, color: lit || p.id === selected ? 0xd6b77c : p.owned ? 0x8f7849 : 0x232e2b});
      group.addChild(plot);

      // The building, one box at a time, so a piece of it can be replaced.
      for (const part of block.parts) {
        const f = faces(at, part);
        const g = new Graphics();
        g.poly(f.left.flatMap(v => [v.x, v.y])).fill(COLOUR(part.left));
        g.poly(f.right.flatMap(v => [v.x, v.y])).fill(COLOUR(part.right));
        g.poly(f.top.flatMap(v => [v.x, v.y])).fill(COLOUR(part.top));
        group.addChild(g);
      }

      const tallest = Math.max(...block.parts.map(q => (q.base || 0) + q.h));
      const centre = project({x: at.x + block.w / 2, y: at.y + block.d / 2});
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

      layer.addChild(group);
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
