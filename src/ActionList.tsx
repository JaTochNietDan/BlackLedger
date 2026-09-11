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
  // "We should probably show options like 'buy kerrigan haulage' before you can
  // afford it instead of having it hidden. We probably should just show all
  // hidden options tbh." An action you cannot take yet is a thing to want, and
  // every refusal in this game is a sentence saying what would change it — so a
  // refused card is worth more than an absent one.
  //
  // The toggle that stood in front of them is gone: "making it so that hidden
  // actions are not hidden anymore, just showed as lower priority in the list."
  // What you can do comes first and what you cannot follows it, in the same
  // list, in the same place it would have been.
  const ordered = [...actions.filter(a => !a.disabled), ...actions.filter(a => a.disabled)];
  const notes = [
    who.owes ? `owes $${who.owes.toLocaleString()}${who.overdue ? ' · overdue' : ''}` : '',
    who.known && who.trust !== undefined ? `thinks of you at ${who.trust}` : '',
    who.sore ? `holds ${who.sore} against you` : '',
    // What the player put in their hand. It changes what sending them does and
    // it was bought and paid for, and the only way to know they had it was to
    // remember buying it.
    who.carrying ? `carrying ${who.carrying.toLowerCase()}` : '',
    // And what the player put them in. It shifts a job that went wrong away
    // from the two endings nobody wants.
    who.driving ? `driving ${who.driving.toLowerCase()}` : '',
    // Somebody of yours who has got far enough down that they are thinking
    // about where else they could be. A man can walk out of here with one of
    // your businesses and the only warning was a number.
    who.restless ? 'thinking about leaving' : '',
  ]
    .filter(Boolean)
    .join(' · ');

  return (
    <article
      className={
        'presence' +
        (who.yours ? ' yours' : '') +
        (who.overdue || who.sore || who.restless ? ' sour' : '')
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
            <small className={who.overdue || who.sore || who.restless ? 'warning' : 'subtle'}>
              {notes}
            </small>
          )}
        </div>
      </header>
      {who.says && <p className="said">{who.says}</p>}
      {ordered.length > 0 && <div className="actions">{ordered.map(render)}</div>}
      {ordered.length === 0 && <p className="nothing-here">Nothing to do with them here.</p>}
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

  // Each section in the order it was already in, and inside it what can be done
  // before what cannot. Nothing moves section: "not changing location of sub
  // sections, just putting unavailable actions at the end of the list in each
  // subsection."
  const sections = placeActions(groups?.length ? groups : [lastResort], impersonal)
    .map(s => ({
      ...s,
      ordered: [...s.mine.filter(a => !a.disabled), ...s.mine.filter(a => a.disabled)],
    }))
    .filter(s => s.ordered.length > 0);

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

      {sections.map(s => (
        <section className="action-group" key={s.id}>
          <h4>
            {s.title}
            <span>{s.blurb}</span>
          </h4>
          <div className="actions">{s.ordered.map(render)}</div>
        </section>
      ))}

      {sections.length === 0 && present.length === 0 && (
        <p className="nothing-here">Nothing here matches “{query}”.</p>
      )}
    </div>
  );
}
