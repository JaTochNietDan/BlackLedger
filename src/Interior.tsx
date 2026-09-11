import {useEffect, useState} from 'react';
import type {ReactElement} from 'react';
import type {Action, Coming, Group, Place, Presence} from './types';
import {interiorSVG, paintedRoom, roomLight, standingSpots, StandingRoom} from './roomart';
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
const premisesOrder = [
  'acquire',
  'repair',
  'hire',
  'layoff',
  'restock',
  'remedy',
  'inspect',
  'operate:standard',
  'operate:clean',
  'operate:hard',
  'post',
  'unpost',
  'still',
  'dismantle',
  'armoury',
  'stock_arms',
  'bankroll',
  'order',
  'fit:door',
  'fit:telephone',
  'fit:safe',
  'fit:cellar',
];

function rank(id: string) {
  const at = premisesOrder.indexOf(id);
  return at < 0 ? premisesOrder.length : at;
}

// One heading and the work under it. What cannot be done yet is folded away:
// a room where eleven of the twenty-six cards are greyed out is a room where
// the player reads eleven refusals to find the two things they can do.
function Work({
  title,
  blurb,
  actions,
  render,
}: {
  title: string;
  blurb: string;
  actions: Action[];
  render: (a: Action) => ReactElement;
}) {
  // What can be done first, what cannot after it, in the one list. The refusals
  // used to sit behind a toggle that started closed, so the room hid half of
  // what it had: "making it so that hidden actions are not hidden anymore, just
  // showed as lower priority in the list … just putting unavailable actions at
  // the end of the list in each subsection."
  const ordered = [...actions.filter(a => !a.disabled), ...actions.filter(a => a.disabled)];
  if (actions.length === 0) return null;
  return (
    <section className="action-group">
      <h4>
        {title}
        <span>{blurb}</span>
      </h4>
      <div className="actions compact">{ordered.map(render)}</div>
    </section>
  );
}

export function Interior({
  place,
  people,
  actions,
  render,
  onLeave,
  onTables,
  felt,
  groups,
  comings,
  minute,
}: {
  place: Place;
  people: Presence[];
  actions: Action[];
  render: (a: Action) => ReactElement;
  onLeave: () => void;
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
  useEffect(() => {
    if (picked && !people.some(p => p.id === picked)) setPicked('');
  }, [people, picked]);

  const inRoom = new Set(people.map(p => p.id));
  const needle = query.trim().toLowerCase();
  const found = actions.filter(
    a => !needle || (a.label + ' ' + a.detail + ' ' + a.reason).toLowerCase().includes(needle),
  );
  const personal = found.filter(a => a.subject && inRoom.has(a.subject));
  const premises = found
    .filter(a => !personal.includes(a) && a.group === 'business')
    .sort((a, b) => rank(a.id) - rank(b.id));
  const elsewhere = found.filter(a => !personal.includes(a) && a.group !== 'business');

  const who = people.find(p => p.id === picked);
  const theirs = personal.filter(a => a.subject === picked);
  const withSomething = new Set(personal.map(a => a.subject!));

  // The floor shows the people the player has something to do with first, then
  // the ones they know. Everybody else is in the roster below and in the count.
  const worth = (p: Presence) =>
    personal.some(a => a.subject === p.id && !a.disabled)
      ? 0
      : p.yours
        ? 1
        : p.owes || p.sore
          ? 2
          : p.known
            ? 3
            : 4;
  const onFloor = [...people].sort((a, b) => worth(a) - worth(b));

  const traffic = (comings || []).filter(c => c.where === place.id);
  // The hour, read once: the room's wash and the people standing in it have to
  // agree, and they only do that if they come off the same number.
  const light = roomLight(minute);

  // Whose room you are standing in. The street panel says this and the room did
  // not, so stepping inside lost the one fact that decides how everything in
  // here should be read: "when inside a building you can't see who the family
  // that owns it (if any) is anymore."
  const held = place.owned ? 'Yours' : place.holder || 'Independent';
  // What is true of this room, as figures rather than as a sentence. The strip
  // used to join them with middots and read as an afterthought; a player
  // standing in a business wants the same three or four facts in the same place
  // in every room, which is what "a more fleshed out display" is asking for.
  const facts: {what: string; is: string; warn?: boolean}[] = [];
  if (place.condition < 100 || place.owned) {
    facts.push({what: 'Condition', is: place.condition + '%', warn: place.condition < 70});
  }
  if (place.owned && typeof place.trading === 'number') {
    facts.push({
      what: 'Working at',
      is: Math.round(place.trading * 100) + '%',
      warn: place.trading < 0.8,
    });
  }
  // The count against what the work takes, because one of them on its own is
  // not a fact anybody can act on: three is right at a laundry and short at a
  // casino, and a counter that has emptied itself reads as nothing at all.
  if (place.owned && place.staff !== undefined) {
    const wants = place.positions ?? place.staff;
    facts.push({
      what: 'On the books',
      is: place.staff < wants ? place.staff + ' of ' + wants : String(place.staff),
      warn: place.staff < wants,
    });
  }
  // What is owed to them. A player who has been short knows the nights are
  // being counted somewhere; this is where.
  if (place.owned && place.unpaid) {
    facts.push({
      what: 'Unpaid',
      is: place.unpaid + (place.unpaid === 1 ? ' night' : ' nights'),
      warn: true,
    });
  }
  if (place.owned && place.income > 0) {
    facts.push({what: 'Earns', is: '$' + place.income + '/hr'});
  }
  // What it pays, now that the wage is a decision rather than a rate. Beside
  // what it earns, because that is the comparison an owner is making.
  if (place.owned && place.wage) {
    facts.push({what: 'Pays', is: '$' + place.wage + '/day'});
  }
  // Who has the keys. The person a rival will come for, and the reason the
  // player does not have to be standing here.
  if (place.owned && place.runs) {
    facts.push({what: 'Run by', is: place.runs});
  }

  return (
    <div className="interior-stage">
      <div className={'room-holder' + (place.owned ? ' yours' : '')}>
        <b>{place.name}</b>
        <span>{held}</span>
        <div className="room-facts">
          {facts.map(f => (
            <i key={f.what}>
              {f.what}
              <b className={f.warn ? 'warning' : ''}>{f.is}</b>
            </i>
          ))}
        </div>
        {place.note && <small className={place.note_warn ? 'warning' : ''}>{place.note}</small>}
      </div>
      {!!traffic.length && (
        <div className="room-traffic" role="status">
          {traffic.map(c => (
            <p key={c.id + String(c.leaving)} className={c.leaving ? 'left' : 'came'}>
              <i aria-hidden="true">{c.leaving ? '←' : '→'}</i>
              {c.note}
            </p>
          ))}
        </div>
      )}
      <div
        className={'room' + (painted ? ' painted' : '')}
        style={painted ? {backgroundImage: `url(${paintedRoom(place.id)})`} : undefined}
      >
        <div
          className="room-plate"
          dangerouslySetInnerHTML={{__html: interiorSVG(place, painted)}}
        />
        {/* The hour, laid over the backdrop rather than baked into it. */}
        <span className="room-light" aria-hidden="true" style={{background: light.wash}} />
        {/* The people are drawn over the room in HTML rather than inside the
          picture, so each one can wear their own face. A silhouette with
          nothing on its head could be anybody. */}
        {onFloor.slice(0, standingSpots.length).map((who, i) => {
          const spot = standingSpots[i];
          return (
            <button
              key={who.id}
              className={
                'stander' +
                (who.id === picked ? ' picked' : '') +
                (who.yours ? ' yours' : '') +
                (who.overdue || who.sore ? ' sour' : '')
              }
              style={{
                left: `${spot.left}%`,
                bottom: `${spot.bottom}%`,
                // The figures are sized against the room they are standing in. When
                // the picture is cropped down to make space for the work below it,
                // everybody in it comes down by the same factor rather than growing
                // into a room half their height.
                transform: `translateX(-50%) scale(calc(var(--fig, 1) * ${spot.scale.toFixed(2)}))`,
                // A person standing in a dark room is dark. The wash over the
                // backdrop used to go under the figures, so at three in the morning
                // the room went dark and everybody in it stayed lit like a shop
                // window. The city dims its people by the same number.
                filter: `drop-shadow(0 6px 10px #000a) brightness(${(1 - light.dark * 0.38).toFixed(2)})`,
              }}
              aria-pressed={who.id === picked}
              title={`${who.name} — ${who.standing}`}
              onClick={() => setPicked(who.id === picked ? '' : who.id)}
            >
              {/* What they are standing on. Nothing in the city is allowed to
              float and neither is anybody in here. */}
              <span className="stander-shadow" aria-hidden="true" />
              <span className="stander-hat" aria-hidden="true" />
              <Portrait id={who.id} face={who.face} size="small" />
              <span className="stander-coat" aria-hidden="true" />
              <span className="stander-name">{who.name.split(' ')[0]}</span>
            </button>
          );
        })}
        {onFloor.length > standingSpots.length && (
          <span className="room-rest">
            and {onFloor.length - standingSpots.length} more in here
          </span>
        )}
      </div>

      <div className="room-people" role="list">
        {people.map(p => (
          <button
            key={p.id}
            role="listitem"
            className={
              'room-chip' +
              (p.id === picked ? ' picked' : '') +
              (p.yours ? ' yours' : '') +
              (p.overdue || p.sore ? ' sour' : '')
            }
            aria-pressed={p.id === picked}
            onClick={() => setPicked(p.id === picked ? '' : p.id)}
          >
            <Portrait id={p.id} face={p.face} size="tiny" />
            <span className="chip-name">
              <b>{p.name}</b>
              <small>{p.standing}</small>
            </span>
            {withSomething.has(p.id) && (
              <i aria-hidden="true" title="You have business with them">
                ·
              </i>
            )}
          </button>
        ))}
        {people.length === 0 && <p className="nothing-here">There is nobody here.</p>}
      </div>

      <div className="room-work">
        {who ? (
          <section className="picked-person only">
            <header>
              <Portrait id={who.id} face={who.face} size="small" />
              <div className="picked-who">
                <b>{who.name}</b>
                <small>
                  {who.standing}
                  {who.temperament ? ` · ${who.temperament}` : ''}
                </small>
                {(who.owes || who.sore || (who.known && who.trust !== undefined)) && (
                  <small className={who.overdue || who.sore ? 'warning' : 'subtle'}>
                    {[
                      who.owes
                        ? `owes $${who.owes.toLocaleString()}${who.overdue ? ' · overdue' : ''}`
                        : '',
                      who.known && who.trust !== undefined ? `thinks of you at ${who.trust}` : '',
                      who.sore ? `holds ${who.sore} against you` : '',
                    ]
                      .filter(Boolean)
                      .join(' · ')}
                  </small>
                )}
              </div>
            </header>
            {/* What they say to your face. The core writes the line out of what
            they are actually carrying; this only prints it. */}
            {who.says && <p className="said">{who.says}</p>}
            <button className="plain step-away" onClick={() => setPicked('')}>
              ← Step away
              {premises.length + elsewhere.length > 0
                ? ` · ${premises.length + elsewhere.length} other things to do here`
                : ''}
            </button>
            {theirs.length ? (
              <div className="actions">{theirs.map(render)}</div>
            ) : (
              <p className="nothing-here">There is nothing to do with them here.</p>
            )}
          </section>
        ) : (
          <>
            <div className="work-head">
              <p className="room-hint">
                Pick somebody in the room to deal with them, or use the building itself.
              </p>
              <div className="work-search">
                <input
                  type="search"
                  value={query}
                  placeholder={`Search ${actions.length} things to do here…`}
                  aria-label="Search what you can do in this room"
                  onChange={e => setQuery(e.target.value)}
                />
                <button className="plain" onClick={onLeave}>
                  ← Back to the street
                </button>
              </div>
            </div>
            {needle && found.length === 0 && (
              <p className="nothing-here">Nothing here matches “{query}”.</p>
            )}

            {onTables && (
              <button className="action primary sit-down-here" onClick={onTables}>
                <span>
                  <strong>{felt ? 'Sit down at the tables' : 'Play the machines'}</strong>
                  <span className="desc">
                    {felt
                      ? 'The cards and the wheel, played out at the table until you get up.'
                      : 'Three drums and a handle, against the wall where they always are.'}
                  </span>
                </span>
              </button>
            )}

            {/* "When inside a building you own the top buttons should probably be
            for owner management and under a separate subtitle for management
            actions." They are the top block already; what was missing was the
            subtitle saying so, and two of the actions that belong in it were
            filed elsewhere by the core until this was written down. */}
            <Work
              title={place.owned ? 'Running ' + place.name : 'These premises'}
              blurb={
                place.owned
                  ? 'Staff, stock, repairs and what the house takes — the work of holding it'
                  : 'The same work, in the same order, in every building'
              }
              actions={premises}
              render={render}
            />

            {/* Everything that is not the premises, in the order the core says the
            work is for, rather than one heading with fifteen cards under it. */}
            {placeActions(
              groups?.length
                ? groups
                : [
                    {
                      id: 'work',
                      title: 'Everything else here',
                      blurb: 'Work, standing, money and leaving',
                    },
                  ],
              elsewhere,
            ).map(g => (
              <Work key={g.id} title={g.title} blurb={g.blurb} actions={g.mine} render={render} />
            ))}

            <button className="plain leave-room" onClick={onLeave}>
              ← Back to the street
            </button>
          </>
        )}
      </div>
    </div>
  );
}
