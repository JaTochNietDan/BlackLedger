import {useEffect,useState} from 'react';
import type {ReactElement} from 'react';
import type {Action,Coming,Place,Presence} from './types';
import {interiorSVG,paintedRoom,standingSpots,StandingRoom} from './roomart';
import {Portrait} from './Portrait';

// Entering a building should open the building, not fill a column. The room is
// the screen: the inside of the place, the people standing in it as things you
// can click, and the work you can do with whoever you picked. The premises
// themselves keep a fixed strip of their own, in the same order in every
// building, so a player learns once where the roof and the staff and the
// supplies live and never hunts for them again.

// The order premises work is always offered in. A player should find restocking
// in the same place at a laundry as at a casino.
const premisesOrder = ['acquire', 'repair', 'hire', 'layoff', 'restock', 'remedy',
  'inspect', 'operate:standard', 'operate:clean', 'operate:hard', 'post', 'unpost',
  'still', 'dismantle', 'armoury', 'stock_arms', 'bankroll', 'order',
  'fit:door', 'fit:telephone', 'fit:safe', 'fit:cellar'];

function rank(id: string) { const at = premisesOrder.indexOf(id); return at < 0 ? premisesOrder.length : at }

export function Interior({place, people, actions, render, onLeave, comings}: {
  place: Place; people: Presence[]; actions: Action[];
  render: (a: Action) => ReactElement; onLeave: () => void;
  // Who walked in or out while the player was standing here. People move
  // between buildings now, and until this the room's roster simply changed
  // behind their back: somebody they had been talking to was gone, and
  // somebody they had never seen was in the list with no explanation.
  comings?: Coming[];
}) {
  const [picked, setPicked] = useState('');
  // The painted interior is a backdrop, not a dependency: if it is missing the
  // drawn room takes its place rather than leaving a hole.
  const [painted, setPainted] = useState(false);
  useEffect(() => {
    setPainted(false);
    const img = new Image();
    img.onload = () => setPainted(true);
    img.src = paintedRoom(place.id);
  }, [place.id]);

  // Whoever the player was talking to may walk out, be arrested, or die.
  useEffect(() => { if (picked && !people.some(p => p.id === picked)) setPicked('') }, [people, picked]);

  const inRoom = new Set(people.map(p => p.id));
  const personal = actions.filter(a => a.subject && inRoom.has(a.subject));
  const premises = actions.filter(a => !personal.includes(a) && a.group === 'business')
    .sort((a, b) => rank(a.id) - rank(b.id));
  const elsewhere = actions.filter(a => !personal.includes(a) && a.group !== 'business');

  const who = people.find(p => p.id === picked);
  const theirs = personal.filter(a => a.subject === picked);
  const withSomething = new Set(personal.map(a => a.subject!));

  // The floor shows the people the player has something to do with first, then
  // the ones they know. Everybody else is in the roster below and in the count.
  const worth = (p: Presence) => (personal.some(a => a.subject === p.id && !a.disabled) ? 0 : p.yours ? 1 : p.owes || p.sore ? 2 : p.known ? 3 : 4);
  const onFloor = [...people].sort((a, b) => worth(a) - worth(b));

  const traffic = (comings || []).filter(c => c.where === place.id);

  return <div className="interior-stage">
    {!!traffic.length && <div className="room-traffic" role="status">
      {traffic.map(c => <p key={c.id + String(c.leaving)} className={c.leaving ? 'left' : 'came'}>
        <i aria-hidden="true">{c.leaving ? '←' : '→'}</i>{c.note}
      </p>)}
    </div>}
    <div className={'room' + (painted ? ' painted' : '')}
      style={painted ? {backgroundImage: `url(${paintedRoom(place.id)})`} : undefined}>
      <div className="room-plate" dangerouslySetInnerHTML={{__html: interiorSVG(place, painted)}}/>
      {/* The people are drawn over the room in HTML rather than inside the
          picture, so each one can wear their own face. A silhouette with
          nothing on its head could be anybody. */}
      {onFloor.slice(0, standingSpots.length).map((who, i) => {
        const spot = standingSpots[i];
        return <button key={who.id}
          className={'stander' + (who.id === picked ? ' picked' : '') + (who.yours ? ' yours' : '') + (who.overdue || who.sore ? ' sour' : '')}
          style={{left: `${spot.left}%`, bottom: `${spot.bottom}%`, transform: `translateX(-50%) scale(${spot.scale.toFixed(2)})`}}
          aria-pressed={who.id === picked}
          title={`${who.name} — ${who.standing}`}
          onClick={() => setPicked(who.id === picked ? '' : who.id)}>
          <Portrait id={who.id} size="small"/>
          <span className="stander-coat" aria-hidden="true"/>
          <span className="stander-name">{who.name.split(' ')[0]}</span>
        </button>;
      })}
      {onFloor.length > standingSpots.length &&
        <span className="room-rest">and {onFloor.length - standingSpots.length} more in here</span>}
    </div>

    <div className="room-people" role="list">
      {people.map(p => <button key={p.id} role="listitem" className={'room-chip' + (p.id === picked ? ' picked' : '') + (p.yours ? ' yours' : '') + (p.overdue || p.sore ? ' sour' : '')}
        aria-pressed={p.id === picked} onClick={() => setPicked(p.id === picked ? '' : p.id)}>
        <Portrait id={p.id} size="tiny"/>
        <span className="chip-name"><b>{p.name}</b><small>{p.standing}</small></span>
        {withSomething.has(p.id) && <i aria-hidden="true" title="You have business with them">·</i>}
      </button>)}
      {people.length === 0 && <p className="nothing-here">There is nobody here.</p>}
    </div>

    <div className="room-work">
      {who ? <section className="picked-person only">
        <header>
          <Portrait id={who.id} size="small"/>
          <div className="picked-who">
          <b>{who.name}</b>
          <small>{who.standing}{who.temperament ? ` · ${who.temperament}` : ''}</small>
          {(who.owes || who.sore || (who.known && who.trust !== undefined)) && <small className={who.overdue || who.sore ? 'warning' : 'subtle'}>
            {[who.owes ? `owes $${who.owes.toLocaleString()}${who.overdue ? ' · overdue' : ''}` : '',
              who.known && who.trust !== undefined ? `thinks of you at ${who.trust}` : '',
              who.sore ? `holds ${who.sore} against you` : ''].filter(Boolean).join(' · ')}
          </small>}
          </div>
        </header>
        <button className="plain step-away" onClick={() => setPicked('')}>
          ← Step away{premises.length + elsewhere.length > 0 ? ` · ${premises.length + elsewhere.length} other things to do here` : ''}
        </button>
        {theirs.length ? <div className="actions">{theirs.map(render)}</div>
          : <p className="nothing-here">There is nothing to do with them here.</p>}
      </section> : <>
        <p className="room-hint">Pick somebody in the room to deal with them, or use the building itself.</p>

        {premises.length > 0 && <section className="action-group">
          <h4>These premises<span>The same work, in the same order, in every building</span></h4>
          <div className="actions">{premises.map(render)}</div>
        </section>}

        {elsewhere.length > 0 && <section className="action-group">
          <h4>Everything else here<span>Work, standing, money and leaving</span></h4>
          <div className="actions">{elsewhere.map(render)}</div>
        </section>}

        <button className="plain leave-room" onClick={onLeave}>← Back to the street</button>
      </>}
    </div>
  </div>;
}
