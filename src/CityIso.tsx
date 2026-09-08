import type {Snapshot, Place} from './types';
import {blockFor, depth, faces, plan, points, project, TILE} from './iso';
import type {Spotlight} from './CityStreet';

// The city, drawn as a city. Every address the core knows stands where the core
// says it stands; nothing here invents a place or moves one. What this adds is
// the ground, the roads between them, and enough of a building to tell a
// laundry from a police station at a glance.

// The ground each block sits on, so a building never floats.
function Plot({at, w, d, lit}: {at: {x: number; y: number}; w: number; d: number; lit: boolean}) {
  const p = [
    project(at), project({x: at.x + w, y: at.y}),
    project({x: at.x + w, y: at.y + d}), project({x: at.x, y: at.y + d}),
  ];
  return <polygon points={points(p)} className={'iso-plot' + (lit ? ' lit' : '')}/>;
}

export function CityIso({state, selected, onSelect, onEnter, spotlight}: {
  state: Snapshot; selected: string; onSelect: (id: string) => void; onEnter: () => void;
  spotlight?: Spotlight | null;
}) {
  const here = state.player.location;

  // Everything placed first, so the whole city can be sorted back to front
  // before any of it is drawn. Anything that ignores this order draws a
  // building through the one in front of it.
  const placed = state.locations.map(p => {
    const block = blockFor(p.type);
    const at = plan(p.x, p.y);
    return {place: p, block, at, d: depth({x: at.x + block.w / 2, y: at.y + block.d / 2})};
  }).sort((a, b) => a.d - b.d);

  // The extent of the drawing, from the city's own coordinates rather than a
  // guess, with room for the tallest roof.
  const xs = placed.flatMap(b => [project(b.at).x, project({x: b.at.x + b.block.w, y: b.at.y}).x,
                                  project({x: b.at.x, y: b.at.y + b.block.d}).x]);
  const ys = placed.flatMap(b => [project(b.at).y, project({x: b.at.x + b.block.w, y: b.at.y + b.block.d}).y]);
  const pad = 90;
  const minX = Math.min(...xs) - pad, maxX = Math.max(...xs) + pad;
  const minY = Math.min(...ys) - pad * 2.4, maxY = Math.max(...ys) + pad;

  return <div className="city-iso">
    <svg viewBox={`${minX} ${minY} ${maxX - minX} ${maxY - minY}`} role="img"
      aria-label="The city of Bellwether, seen from above">
      <defs>
        <linearGradient id="iso-sky" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stopColor="#16211f"/><stop offset="1" stopColor="#0d1413"/>
        </linearGradient>
      </defs>
      <rect x={minX} y={minY} width={maxX - minX} height={maxY - minY} fill="url(#iso-sky)"/>

      {placed.map(({place, block, at}) => {
        const lit = spotlight?.id === place.id;
        const shut = place.district > state.district;
        return <g key={place.id} data-place={place.id}
          className={'iso-block' + (place.id === selected ? ' picked' : '') + (place.id === here ? ' here' : '')
            + (place.owned ? ' owned' : '') + (lit ? ' lit' : '') + (shut ? ' shut' : '')}
          onClick={() => onSelect(place.id)}
          onDoubleClick={() => { if (place.id === here) onEnter() }}>
          <Plot at={at} w={block.w} d={block.d} lit={lit}/>
          {block.parts.map((part, i) => {
            const f = faces(at, part);
            return <g key={i}>
              <polygon points={points(f.left)} fill={part.left}/>
              <polygon points={points(f.right)} fill={part.right}/>
              <polygon points={points(f.top)} fill={part.top}/>
            </g>;
          })}
          {/* The name, above the roof, so the city can be read as well as looked at. */}
          {(() => {
            const tallest = Math.max(...block.parts.map(p => (p.base || 0) + p.h));
            const centre = project({x: at.x + block.w / 2, y: at.y + block.d / 2});
            return <text className="iso-name" x={centre.x} y={centre.y - tallest * TILE.h - 12}
              textAnchor="middle">{place.name}</text>;
          })()}
        </g>;
      })}
    </svg>
  </div>;
}
