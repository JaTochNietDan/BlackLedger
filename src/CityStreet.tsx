import type {Snapshot} from './types';
import {paintedAsset,paintedFront} from './cityAssets';
import {useEffect,useLayoutEffect,useRef,useState} from 'react';
import {Portrait} from './Portrait';

// The street view showed five buildings. The city has twelve, and two of the
// seven it left out were places this project added itself. Half the city was
// reachable only through the address book, and a view called "the street"
// quietly told the player that the rest of it was not there.
//
// This is every address, painted, in the district it stands in, with the people
// who are actually in it and the hour of the day on it. The animated Old Harbor
// block is still there and still worth looking at; it is one block of a city
// that has three.

const districts = ['Old Harbor', 'Ashbury', 'The Heights'];
const FACES = 3;

// How dark it is. The city runs on a 1440-minute day, so the light is simply
// what the clock says rather than anything the interface decides.
function night(minute: number) {
  const hour = (minute % 1440) / 60;
  if (hour < 5 || hour >= 21) return 1;
  if (hour < 7) return (7 - hour) / 2;
  if (hour > 18) return (hour - 18) / 3;
  return 0;
}

// A moment the city is holding the camera on: which address, what kind of thing
// happened, and how far through it is. The street dims around it and the effect
// plays over the real painted front rather than over a drawing of one.
export type Spotlight = {id: string; kind: string; t: number};

// What light a kind of moment throws. Only light: the painted fronts are the
// city, and drawing a body or a car on top of one reads as a sticker.
function lit(kind: string, t: number) {
  if (kind === 'gunfight' || kind === 'killing') {
    const firing = t < .45 && Math.floor(t * 14) % 2 === 0;
    return {flash: firing ? .85 : 0, glow: '#ffe6a8', shake: 0};
  }
  if (kind === 'explosion') {
    return {flash: Math.max(0, .95 - t * 1.4), glow: '#ffcf7a', shake: Math.max(0, 1 - t * 3) * 7};
  }
  if (kind === 'raid' || kind === 'arrest') {
    return {flash: .18 + Math.abs(Math.sin(t * 22)) * .22, glow: '#e0705c', shake: 0};
  }
  return {flash: .1, glow: '#d6b77c', shake: 0};
}

export function CityStreet({state, selected, onSelect, onEnter, spotlight}: {
  state: Snapshot; selected: string; onSelect: (id: string) => void; onEnter: () => void;
  spotlight?: Spotlight | null;
}) {
  const dark = night(state.minute);
  const here = state.player.location;

  const street = state.street || [];

  // Where each walker stands on the screen: between the front they left and
  // the front they are going to, as far along as they actually are. The fronts
  // are laid out by the browser across three district blocks, so the only
  // honest source for the two ends of a journey is where they were actually
  // drawn. Positions are recomputed when the walkers change and when the
  // window does; nothing animates between actions, because nothing moves
  // between actions — the clock is stopped.
  const stage = useRef<HTMLDivElement | null>(null);
  const [figures, setFigures] = useState<{id: string; name: string; left: number; top: number; yours?: boolean}[]>([]);
  const place = () => {
    const box = stage.current?.getBoundingClientRect();
    if (!box) return;
    setFigures(street.flatMap(j => {
      const from = stage.current!.querySelector(`[data-place="${j.from_id}"]`);
      const to = stage.current!.querySelector(`[data-place="${j.to_id}"]`);
      if (!from || !to) return [];
      const a = from.getBoundingClientRect(), b = to.getBoundingClientRect();
      const at = Math.min(1, Math.max(0, j.progress));
      // Along the line between the two doorways, which is the bottom middle of
      // each front rather than its centre: people walk on the pavement.
      const ax = a.left + a.width / 2, ay = a.bottom;
      const bx = b.left + b.width / 2, by = b.bottom;
      return [{id: j.id, name: j.name, yours: j.yours,
        left: ax + (bx - ax) * at - box.left,
        top: ay + (by - ay) * at - box.top}];
    }));
  };
  useLayoutEffect(place, [state.minute, state.revision, street.length]);
  useEffect(() => {
    window.addEventListener('resize', place);
    return () => window.removeEventListener('resize', place);
  });

  // Bring the address the camera is on into view, once per moment.
  useEffect(() => {
    if (!spotlight) return;
    document.querySelector(`[data-place="${spotlight.id}"]`)
      ?.scrollIntoView({block: 'center', behavior: 'smooth'});
  }, [spotlight?.id]);

  return <div ref={stage} className={'city-street' + (spotlight ? ' watching' : '')}
    style={{'--night': dark.toFixed(2)} as React.CSSProperties}>
    {/* The people actually crossing the city, on the city. */}
    {figures.map(f => <div key={f.id} className={'walker-figure' + (f.yours ? ' yours' : '')}
      style={{left: f.left + 'px', top: f.top + 'px'}} title={f.name}>
      <Portrait id={f.id} size="tiny"/>
      <b className="walker-tag">{f.name.split(' ')[0]}</b>
    </div>)}
    {/* Who is out there. People in this city used to stand at one address for
        life; this is the only place you can watch one of them cross it. */}
    {!!street.length && <section className="street-out" aria-label="People on the street">
      <h3>OUT ON THE STREET</h3>
      <div className="street-out-run">
        {street.map(j => <div className={'walker' + (j.yours ? ' yours' : '')} key={j.id}>
          <Portrait id={j.id} size="tiny"/>
          <div>
            <b>{j.name}</b>
            <small>{j.from} <i>→</i> {j.to}</small>
            <small className="walker-why">{j.because} · {j.minutes} min out</small>
            <span className="walker-track"><i style={{left: (j.progress * 100).toFixed(0) + '%'}}/></span>
          </div>
        </div>)}
      </div>
    </section>}
    {districts.map((name, d) => {
      const places = state.locations.filter(p => p.district === d);
      if (!places.length) return null;
      const shut = d > state.district;
      return <section className={'street-block' + (shut ? ' locked' : '')} key={name}>
        <h3>{name}{shut && <span>Not open to you yet</span>}</h3>
        <div className="street-run">
          {places.map(p => {
            const art = paintedAsset(p.id, p.condition);
            const people = p.people || [];
            const shot = spotlight && spotlight.id === p.id ? lit(spotlight.kind, spotlight.t) : null;
            return <button key={p.id} data-place={p.id} className={'street-front'
                + (shot ? ' lit' : '')
                + (p.id === selected ? ' picked' : '')
                + (p.id === here ? ' here' : '')
                + (p.owned ? ' owned' : '')
                + (paintedFront(p.id) ? ' photo' : ' cutout')}
              aria-pressed={p.id === selected}
              onClick={() => onSelect(p.id)}
              onDoubleClick={() => { if (p.id === here) onEnter() }}>
              <span className="front-art" style={shot && shot.shake
                ? {transform: `translate(${(Math.random() - .5) * shot.shake}px,${(Math.random() - .5) * shot.shake}px)`}
                : undefined}>
                {shot && <span className="front-flash" aria-hidden="true"
                  style={{background: shot.glow, opacity: shot.flash}}/>}
                {art ? <img src={art} alt="" loading="lazy"/> : <span className="front-blank"/>}
                <span className="front-lit" aria-hidden="true"/>
              </span>
              <span className="front-plate">
                <b>{p.name}</b>
                <small>{p.owned ? 'Yours' : p.holder || 'Independent'}</small>
                {p.id !== here && !shut && <small className="front-travel">{p.travel_note}</small>}
                {p.note && <small className={'front-note' + (p.note_warn ? ' warning' : '')}>{p.note}</small>}
                {!!people.length && <span className="front-who">
                  {people.slice(0, FACES).map(w => <Portrait key={w.id} id={w.id} size="tiny"/>)}
                  <i>{people.length}</i>
                </span>}
              </span>
              {p.id === here
                ? <span className="front-here">You are here</span>
                : !shut && <span className="front-away">{p.away} min</span>}
            </button>;
          })}
        </div>
      </section>;
    })}
    <p className="street-note">
      Pick an address to inspect it, and see what the journey costs before you commit to it. You can only step inside the one you are standing in.
      {' '}<a href="/street-study.html" target="_blank" rel="noreferrer">See the Old Harbor block painted →</a>
    </p>
  </div>;
}
