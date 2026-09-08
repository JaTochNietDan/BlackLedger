import type {Snapshot} from './types';
import {paintedAsset,paintedFront} from './cityAssets';
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

export function CityStreet({state, selected, onSelect, onEnter}: {
  state: Snapshot; selected: string; onSelect: (id: string) => void; onEnter: () => void;
}) {
  const dark = night(state.minute);
  const here = state.player.location;

  const street = state.street || [];

  return <div className="city-street" style={{'--night': dark.toFixed(2)} as React.CSSProperties}>
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
            return <button key={p.id} className={'street-front'
                + (p.id === selected ? ' picked' : '')
                + (p.id === here ? ' here' : '')
                + (p.owned ? ' owned' : '')
                + (paintedFront(p.id) ? ' photo' : ' cutout')}
              aria-pressed={p.id === selected}
              onClick={() => onSelect(p.id)}
              onDoubleClick={() => { if (p.id === here) onEnter() }}>
              <span className="front-art">
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
