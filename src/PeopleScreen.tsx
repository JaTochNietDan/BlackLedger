import type {ReactElement} from 'react';
import {useMemo,useState} from 'react';
import type {Action,Snapshot,Presence} from './types';
import {Portrait} from './Portrait';

// Fifty-three cards and four screens of scrolling, in whatever order the save
// happened to hold them, with the man who works for you indistinguishable from
// a docker he has never met. The core now says why each person matters and
// what they are doing; this reads that back as the people in your life first
// and the rest of the city behind them.

const groups: [string, string, string][] = [
  ['yours', 'Your people', 'They answer to you and they cost you every day'],
  ['crew', 'Your crew', 'The people who came up with you'],
  ['owes', 'They owe you', 'Money out with a name on it'],
  ['sore', 'Bad blood', 'They are carrying something against you'],
  ['job', 'Names everybody knows', 'The people who hold something in this city'],
  ['organization', 'Organizations', 'Who answers to whom'],
  ['street', 'The street', 'Everybody else in Bellwether'],
];

const money = (n: number) => '$' + Math.floor(n).toLocaleString();

function Card({who}: {who: Presence}) {
  return <article className={'person-card' + (who.yours ? ' yours' : '') + (who.overdue || who.sore ? ' sour' : '') + (who.walking ? ' walking' : '')}>
    <Portrait id={who.id} size="small"/>
    <div className="person-of">
      <b>{who.name}</b>
      <small>{who.standing}{who.temperament ? ` · ${who.temperament}` : ''}</small>
      <small className="doing">{who.walking && <i className="on-street" aria-hidden="true">↗</i>}{who.doing}</small>
      {(who.owes || who.sore || (who.known && who.trust !== undefined)) &&
        <small className={who.overdue || who.sore ? 'warning' : 'subtle'}>
          {[who.owes ? `owes ${money(who.owes)}${who.overdue ? ' · overdue' : ''}` : '',
            who.known && who.trust !== undefined ? `thinks of you at ${who.trust}` : '',
            who.sore ? `holds ${who.sore} against you` : ''].filter(Boolean).join(' · ')}
        </small>}
    </div>
  </article>;
}

export function PeopleScreen({world, actions = [], render}: {
  world: Snapshot;
  // Work about people rather than about a building: putting a price on a name,
  // asking what is being said. It used to be printed at the exchange, which
  // meant crossing the city to reach a decision about somebody who was never
  // there.
  actions?: Action[];
  render?: (a: Action) => ReactElement;
}) {
  const everyone = world.everyone || [];
  const [query, setQuery] = useState('');
  const [only, setOnly] = useState<string | null>(null);
  const [openStreet, setOpenStreet] = useState(false);

  const needle = query.trim().toLowerCase();
  const matches = useMemo(() => everyone.filter(p => !needle ||
    (p.name + ' ' + p.standing + ' ' + (p.doing || '') + ' ' + (p.where || '')).toLowerCase().includes(needle)
  ), [everyone, needle]);

  const sections = groups
    .map(([id, title, blurb]) => ({id, title, blurb, people: matches.filter(p => p.because === id)}))
    .filter(s => s.people.length && (!only || only === s.id));

  const counts = groups.map(([id]) => [id, everyone.filter(p => p.because === id).length] as const)
    .filter(([, n]) => n > 0);

  return <section className="section-content">
    <div className="eyebrow">PEOPLE ARE YOUR EMPIRE</div>
    <h1 className="screen-title">Names worth knowing</h1>
    <p className="subtle">
      {world.population
        ? `${world.population.living} people live in this city and you know ${world.population.known} of them. ${world.population.organized} answer to an organization, ${world.population.jobs} hold one of the city's jobs, and ${world.population.street} answer to nobody.`
        : 'Everybody in Bellwether, and what they are doing about it.'}
    </p>

    {!!render && actions.length > 0 && <section className="anywhere-strip" aria-label="What you can do about people">
      <h4>Whatever room you are in<span>These are about people, not about premises</span></h4>
      <div className="actions">{actions.filter(a => a.id === 'contract' || a.id === 'investigate').map(render)}</div>
    </section>}

    <div className="people-controls">
      <input type="search" value={query} placeholder={`Search ${everyone.length} people…`}
        aria-label="Search everybody in the city" onChange={e => setQuery(e.target.value)}/>
      <div className="people-filters">
        <button aria-pressed={only === null} onClick={() => setOnly(null)}>Everybody</button>
        {counts.map(([id, n]) => {
          const title = groups.find(g => g[0] === id)![1];
          return <button key={id} aria-pressed={only === id} onClick={() => setOnly(only === id ? null : id)}>
            {title} <i>{n}</i>
          </button>;
        })}
      </div>
    </div>

    {sections.map(s => {
      // The street is most of the city and none of it is urgent. It stays
      // folded unless it is asked for, so the screen opens on the people the
      // player actually has something to do with.
      const folded = s.id === 'street' && !needle && only !== 'street' && !openStreet;
      return <section className="people-group" key={s.id}>
        <h2>{s.title}<span>{s.blurb}</span><i>{s.people.length}</i></h2>
        {folded
          ? <button className="reveal-blocked" onClick={() => setOpenStreet(true)}>
              Show {s.people.length} more people in Bellwether
            </button>
          : <div className="person-grid">{s.people.map(p => <Card key={p.id} who={p}/>)}</div>}
      </section>;
    })}

    {sections.length === 0 && <p className="nothing-here">Nobody in this city matches “{query}”.</p>}
  </section>;
}
