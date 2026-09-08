import {useMemo,useState} from 'react';
import type {ReactElement} from 'react';
import type {Action} from './types';

// Ninety actions rendered as one column in whatever order the core built them
// is not a menu. Standing in a bar produced an undifferentiated scroll where
// taking a job, robbing the till, lending a stranger money and paying your own
// man a share all looked identical.
//
// The core now says what each action is for. This renders that: the things you
// can do now, grouped and headed, with everything currently unavailable folded
// away behind a count rather than deleted — a player needs to know that putting
// somebody on a door exists and why they cannot do it yet.

const groupTitles: [string, string, string][] = [
  ['work', 'Work', 'Jobs that pay today'],
  ['business', 'Your premises', 'Keeping what you own earning'],
  ['people', 'People', 'Who works for you, who owes you'],
  ['street', 'The street', 'Work that can go wrong'],
  ['standing', 'Standing', 'Who you are to this city'],
  ['money', 'Money', 'Moving it, hiding it, spending it'],
  ['travel', 'Elsewhere', 'Leaving where you are standing'],
];

export function ActionList({actions, render}: {actions: Action[]; render: (a: Action) => ReactElement}) {
  const [query, setQuery] = useState('');
  const [openBlocked, setOpenBlocked] = useState<Record<string, boolean>>({});

  const needle = query.trim().toLowerCase();
  const matches = useMemo(() => actions.filter(a =>
    !needle || (a.label + ' ' + a.detail + ' ' + a.reason).toLowerCase().includes(needle)
  ), [actions, needle]);

  // Travel is what a player reaches for most and it is one button, so it goes
  // to the top rather than into a section of its own at the bottom.
  const lead = matches.filter(a => a.id === 'travel');
  const rest = matches.filter(a => a.id !== 'travel');

  const sections = groupTitles.map(([id, title, blurb]) => {
    const mine = rest.filter(a => a.group === id);
    return {id, title, blurb, open: mine.filter(a => !a.disabled), blocked: mine.filter(a => a.disabled)};
  }).filter(s => s.open.length || s.blocked.length);

  const available = matches.filter(a => !a.disabled).length;

  return <div className="action-list">
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

    {sections.map(s => {
      const showBlocked = openBlocked[s.id] || !!needle;
      return <section className="action-group" key={s.id}>
        <h4>{s.title}<span>{s.blurb}</span></h4>
        {s.open.length > 0
          ? <div className="actions">{s.open.map(render)}</div>
          : <p className="nothing-here">Nothing available here right now.</p>}
        {s.blocked.length > 0 && <>
          {!needle && <button className="reveal-blocked" aria-expanded={showBlocked} onClick={() => setOpenBlocked(o => ({...o, [s.id]: !o[s.id]}))}>
            {showBlocked ? 'Hide' : 'Show'} {s.blocked.length} you cannot do yet
          </button>}
          {showBlocked && <div className="actions blocked">{s.blocked.map(render)}</div>}
        </>}
      </section>;
    })}

    {sections.length === 0 && <p className="nothing-here">Nothing here matches “{query}”.</p>}
  </div>;
}
