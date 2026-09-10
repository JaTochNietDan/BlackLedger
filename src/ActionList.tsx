import {useMemo, useState} from 'react';
import type {ReactElement} from 'react';
import type {Action, Group, Presence} from './types';
import {placeActions} from './grouping';

// A location was a set of premises with a list of verbs attached. But half of
// what a player does anywhere is done to somebody who happens to be standing
// there, and rendering "Lend Perla Mraz money" as one more card in a column of
// verbs throws away the thing the decision is actually about, which is Perla
// Mraz — who she is, what she thinks of you, and what she already owes.
//
// So the room is the organising idea. The premises first, then the people in
// it with their own work under their own names, then everything else that is
// neither.

// The groups come from the core, which owns them and says so. This file used to
// keep its own copy of the list and the copy had drifted: it was missing
// "people", so any action in that group whose subject was not standing in the
// room vanished from the panel — no button, no reason, nothing. Paying the crew
// a bonus while they are out on collections is the ordinary case of that.
//
// The fallback below is the same one the core applies to an action nobody has
// classified: put it under work rather than lose it.
const lastResort: Group = {id: 'work', title: 'Work', blurb: 'Jobs that pay today'};

function Person({
  who,
  actions,
  render,
}: {
  who: Presence;
  actions: Action[];
  render: (a: Action) => ReactElement;
}) {
  // Open. "We should probably show options like 'buy kerrigan haulage' before
  // you can afford it instead of having it hidden. We probably should just show
  // all hidden options tbh." An action you cannot take yet is a thing to want,
  // and every refusal in this game is a sentence saying what would change it —
  // so a refused card is worth more than an absent one. The toggle stays, as a
  // way to tidy rather than a wall to get past.
  const [open, setOpen] = useState(true);
  const available = actions.filter(a => !a.disabled);
  const blocked = actions.filter(a => a.disabled);
  const notes = [
    who.owes ? `owes $${who.owes.toLocaleString()}${who.overdue ? ' · overdue' : ''}` : '',
    who.known && who.trust !== undefined ? `thinks of you at ${who.trust}` : '',
    who.sore ? `holds ${who.sore} against you` : '',
  ]
    .filter(Boolean)
    .join(' · ');

  return (
    <article
      className={
        'presence' + (who.yours ? ' yours' : '') + (who.overdue || who.sore ? ' sour' : '')
      }
    >
      <header>
        <div>
          <b>{who.name}</b>
          <small>
            {who.standing}
            {who.temperament ? ` · ${who.temperament}` : ''}
          </small>
          {notes && (
            <small className={who.overdue || who.sore ? 'warning' : 'subtle'}>{notes}</small>
          )}
        </div>
      </header>
      {who.says && <p className="said">{who.says}</p>}
      {available.length > 0 && <div className="actions">{available.map(render)}</div>}
      {blocked.length > 0 && (
        <>
          <button className="reveal-blocked" aria-expanded={open} onClick={() => setOpen(o => !o)}>
            {open ? 'Hide' : 'Show'} {blocked.length} you cannot do with {who.name.split(' ')[0]}{' '}
            yet
          </button>
          {open && <div className="actions blocked">{blocked.map(render)}</div>}
        </>
      )}
      {available.length === 0 && blocked.length === 0 && (
        <p className="nothing-here">Nothing to do with them here.</p>
      )}
    </article>
  );
}

export function ActionList({
  actions,
  people,
  render,
  groups,
  here = true,
}: {
  actions: Action[];
  people: Presence[];
  render: (a: Action) => ReactElement;
  groups?: Group[];
  here?: boolean;
}) {
  const [query, setQuery] = useState('');
  // Refusals are shown, and a section can be folded away rather than opened up.
  // Keyed by section so the state is per group, and absent means shown.
  const [hidBlocked, setHidBlocked] = useState<Record<string, boolean>>({});
  const [showRoom, setShowRoom] = useState(false);

  const needle = query.trim().toLowerCase();
  const hay = (a: Action) => (a.label + ' ' + a.detail + ' ' + a.reason).toLowerCase();
  const matches = useMemo(
    () => actions.filter(a => !needle || hay(a).includes(needle)),
    [actions, needle],
  );

  // Travel is what a player reaches for most and it is one button, so it goes
  // to the top rather than into a section of its own at the bottom.
  const lead = matches.filter(a => a.id === 'travel');
  const rest = matches.filter(a => a.id !== 'travel');

  // Anything aimed at somebody standing here belongs under their name.
  const inRoom = new Set(people.map(p => p.id));
  const aimed = rest.filter(a => a.subject && inRoom.has(a.subject));
  const impersonal = rest.filter(a => !a.subject || !inRoom.has(a.subject));

  const present = people
    .map(who => ({who, mine: aimed.filter(a => a.subject === who.id)}))
    .filter(p => p.mine.length > 0 || (!needle && p.who.known));

  const shown = new Set(present.map(p => p.who.id));
  const bystanders = people.filter(p => !shown.has(p.id));

  const sections = placeActions(groups?.length ? groups : [lastResort], impersonal)
    .map(s => ({
      ...s,
      open: s.mine.filter(a => !a.disabled),
      blocked: s.mine.filter(a => a.disabled),
    }))
    .filter(s => s.open.length || s.blocked.length);

  const available = matches.filter(a => !a.disabled).length;

  return (
    <div className="action-list">
      <div className="action-search">
        <input
          type="search"
          value={query}
          placeholder={`Search ${actions.length} actions here…`}
          aria-label="Search the actions available here"
          onChange={e => setQuery(e.target.value)}
        />
        <small>{needle ? `${matches.length} match` : `${available} available`}</small>
      </div>

      {lead.length > 0 && <div className="actions lead">{lead.map(render)}</div>}

      {present.length > 0 && (
        <section className="action-group people-here">
          <h4>
            {here ? 'In the room' : 'Who is there'}
            <span>
              {present.length} of {people.length} you can deal with
            </span>
          </h4>
          {present.map(p => (
            <Person key={p.who.id} who={p.who} actions={p.mine} render={render} />
          ))}
          {bystanders.length > 0 && (
            <>
              <button
                className="reveal-blocked"
                aria-expanded={showRoom}
                onClick={() => setShowRoom(o => !o)}
              >
                {showRoom ? 'Hide' : 'Show'} {bystanders.length}{' '}
                {here ? 'others in the room' : 'others there'}
              </button>
              {showRoom && (
                <ul className="bystanders">
                  {bystanders.map(b => (
                    <li key={b.id}>
                      <b>{b.name}</b>
                      <span>{b.standing}</span>
                    </li>
                  ))}
                </ul>
              )}
            </>
          )}
        </section>
      )}

      {sections.map(s => {
        const showBlocked = !hidBlocked[s.id] || !!needle;
        return (
          <section className="action-group" key={s.id}>
            <h4>
              {s.title}
              <span>{s.blurb}</span>
            </h4>
            {s.open.length > 0 ? (
              <div className="actions">{s.open.map(render)}</div>
            ) : (
              <p className="nothing-here">Nothing available here right now.</p>
            )}
            {s.blocked.length > 0 && (
              <>
                {!needle && (
                  <button
                    className="reveal-blocked"
                    aria-expanded={showBlocked}
                    onClick={() => setHidBlocked(o => ({...o, [s.id]: !o[s.id]}))}
                  >
                    {showBlocked ? 'Hide' : 'Show'} {s.blocked.length} you cannot do yet
                  </button>
                )}
                {showBlocked && <div className="actions blocked">{s.blocked.map(render)}</div>}
              </>
            )}
          </section>
        );
      })}

      {sections.length === 0 && present.length === 0 && (
        <p className="nothing-here">Nothing here matches “{query}”.</p>
      )}
    </div>
  );
}
