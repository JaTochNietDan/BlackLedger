import {useEffect,useState} from 'react';
import type {ReactElement} from 'react';
import type {Action,Coming,Group,Place,Presence} from './types';
import {interiorSVG,paintedRoom,roomLight,standingSpots,StandingRoom} from './roomart';
import {placeActions} from './grouping';
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


// One heading and the work under it. What cannot be done yet is folded away:
// a room where eleven of the twenty-six cards are greyed out is a room where
// the player reads eleven refusals to find the two things they can do.
function Work({title, blurb, actions, render}: {
  title: string; blurb: string; actions: Action[]; render: (a: Action) => ReactElement;
}) {
  const [open, setOpen] = useState(false);
  const ready = actions.filter(a => !a.disabled);
  const blocked = actions.filter(a => a.disabled);
  if (actions.length === 0) return null;
  return <section className="action-group">
    <h4>{title}<span>{blurb}</span></h4>
    {ready.length > 0
      ? <div className="actions compact">{ready.map(render)}</div>
      : <p className="nothing-here">Nothing here you can do right now.</p>}
    {blocked.length > 0 && <>
      <button className="reveal-blocked" aria-expanded={open} onClick={() => setOpen(o => !o)}>
        {open ? 'Hide' : 'Show'} {blocked.length} you cannot do yet
      </button>
      {open && <div className="actions compact blocked">{blocked.map(render)}</div>}
    </>}
  </section>;
}

export function Interior({place, people, actions, render, onLeave, onTables, felt, groups, comings, minute}: {
  place: Place; people: Presence[]; actions: Action[];
  render: (a: Action) => ReactElement; onLeave: () => void;
  // A room with tables in it offers one way in and the tables take the screen.
  // The cards and the wheel are not premises work to be listed between hiring
  // and restocking: they are a place you sit down.
  onTables?: () => void;
  // Whether that way in leads to a felt or to a wall of machines. The room's own
  // action list has had both taken out of it by the time it reaches here, so it
  // cannot work this out for itself.
  felt?: boolean;
  // The core's own ordering of what work is for. Without it the room had one
  // heading called "everything else here" with fifteen cards under it, which is
  // a list, not an order.
  groups?: Group[];
  // The hour, from the core's own clock — the same number the city outside is
  // lit from, so the inside and the outside are the same place at the same
  // time of day rather than two pictures that happen to share a save.
  minute: number;
  // Who walked in or out while the player was standing here. People move
  // between buildings now, and until this the room's roster simply changed
  // behind their back: somebody they had been talking to was gone, and
  // somebody they had never seen was in the list with no explanation.
  comings?: Coming[];
}) {
  const [picked, setPicked] = useState('');
  // The sidebar's action search came in here with the work. A room with
  // twenty-six things in it needs a way to find one by name.
  const [query, setQuery] = useState('');
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
  const needle = query.trim().toLowerCase();
  const found = actions.filter(a => !needle ||
    (a.label + ' ' + a.detail + ' ' + a.reason).toLowerCase().includes(needle));
  const personal = found.filter(a => a.subject && inRoom.has(a.subject));
  const premises = found.filter(a => !personal.includes(a) && a.group === 'business')
    .sort((a, b) => rank(a.id) - rank(b.id));
  const elsewhere = found.filter(a => !personal.includes(a) && a.group !== 'business');

  const who = people.find(p => p.id === picked);
  const theirs = personal.filter(a => a.subject === picked);
  const withSomething = new Set(personal.map(a => a.subject!));

  // The floor shows the people the player has something to do with first, then
  // the ones they know. Everybody else is in the roster below and in the count.
  const worth = (p: Presence) => (personal.some(a => a.subject === p.id && !a.disabled) ? 0 : p.yours ? 1 : p.owes || p.sore ? 2 : p.known ? 3 : 4);
  const onFloor = [...people].sort((a, b) => worth(a) - worth(b));

  const traffic = (comings || []).filter(c => c.where === place.id);
  // The hour, read once: the room's wash and the people standing in it have to
  // agree, and they only do that if they come off the same number.
  const light = roomLight(minute);

  return <div className="interior-stage">
    {!!traffic.length && <div className="room-traffic" role="status">
      {traffic.map(c => <p key={c.id + String(c.leaving)} className={c.leaving ? 'left' : 'came'}>
        <i aria-hidden="true">{c.leaving ? '←' : '→'}</i>{c.note}
      </p>)}
    </div>}
    <div className={'room' + (painted ? ' painted' : '')}
      style={painted ? {backgroundImage: `url(${paintedRoom(place.id)})`} : undefined}>
      <div className="room-plate" dangerouslySetInnerHTML={{__html: interiorSVG(place, painted)}}/>
      {/* The hour, laid over the backdrop rather than baked into it. */}
      <span className="room-light" aria-hidden="true" style={{background: light.wash}}/>
      {/* The people are drawn over the room in HTML rather than inside the
          picture, so each one can wear their own face. A silhouette with
          nothing on its head could be anybody. */}
      {onFloor.slice(0, standingSpots.length).map((who, i) => {
        const spot = standingSpots[i];
        return <button key={who.id}
          className={'stander' + (who.id === picked ? ' picked' : '') + (who.yours ? ' yours' : '') + (who.overdue || who.sore ? ' sour' : '')}
          style={{
            left: `${spot.left}%`, bottom: `${spot.bottom}%`,
            // The figures are sized against the room they are standing in. When
            // the picture is cropped down to make space for the work below it,
            // everybody in it comes down by the same factor rather than growing
            // into a room half their height.
            transform: `translateX(-50%) scale(calc(var(--fig, 1) * ${spot.scale.toFixed(2)}))`,
            // A person standing in a dark room is dark. The wash over the
            // backdrop used to go under the figures, so at three in the morning
            // the room went dark and everybody in it stayed lit like a shop
            // window. The city dims its people by the same number.
            filter: `drop-shadow(0 6px 10px #000a) brightness(${(1 - light.dark * .38).toFixed(2)})`,
          }}
          aria-pressed={who.id === picked}
          title={`${who.name} — ${who.standing}`}
          onClick={() => setPicked(who.id === picked ? '' : who.id)}>
          {/* What they are standing on. Nothing in the city is allowed to
              float and neither is anybody in here. */}
          <span className="stander-shadow" aria-hidden="true"/>
          <span className="stander-hat" aria-hidden="true"/>
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
        <div className="work-head">
          <p className="room-hint">Pick somebody in the room to deal with them, or use the building itself.</p>
          <div className="work-search">
            <input type="search" value={query} placeholder={`Search ${actions.length} things to do here…`}
              aria-label="Search what you can do in this room" onChange={e => setQuery(e.target.value)}/>
            <button className="plain" onClick={onLeave}>← Back to the street</button>
          </div>
        </div>
        {needle && found.length === 0 && <p className="nothing-here">Nothing here matches “{query}”.</p>}

        {onTables && <button className="action primary sit-down-here" onClick={onTables}>
          <span><strong>{felt ? 'Sit down at the tables' : 'Play the machines'}</strong>
          <span className="desc">{felt
            ? 'The cards and the wheel, played out at the table until you get up.'
            : 'Three drums and a handle, against the wall where they always are.'}</span></span>
        </button>}

        <Work title="These premises" blurb="The same work, in the same order, in every building"
              actions={premises} render={render}/>

        {/* Everything that is not the premises, in the order the core says the
            work is for, rather than one heading with fifteen cards under it. */}
        {placeActions(groups?.length ? groups : [{id: 'work', title: 'Everything else here', blurb: 'Work, standing, money and leaving'}], elsewhere)
          .map(g => <Work key={g.id} title={g.title} blurb={g.blurb} actions={g.mine} render={render}/>)}

        <button className="plain leave-room" onClick={onLeave}>← Back to the street</button>
      </>}
    </div>
  </div>;
}
