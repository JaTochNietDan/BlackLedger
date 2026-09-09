import {useMemo,useState} from 'react';
import type {Snapshot,Record as CityRecord} from './types';

// The Ledger showed the market, two figures and then sixty rows of
// undifferentiated history: five screens of scrolling for a page whose whole
// job is to answer "what am I worth and what is this costing me". The core now
// adds the books up; this reads them, and turns the history into something a
// person can search rather than a wall they scroll past.
//
// The underground market has since moved out to a page of its own. A price is
// not an account — it is a reason to go somewhere — and somebody checking
// whether moonshine is worth moving today is not doing bookkeeping. What is
// left here is what the player owes, owns and earns, and what they did.
//
// The breakdown of the day's costs stays, because that is genuinely owed money,
// but it folds away: it is the detail behind a figure that is already on the
// page twice, and open by default it pushed the history below the fold.

const money = (n: number) => (n < 0 ? '−$' : '$') + Math.abs(Math.floor(n)).toLocaleString();
const time = (m: number) => `Day ${Math.floor(m / 1440) + 1} · ${String(Math.floor(m % 1440 / 60)).padStart(2, '0')}:${String(m % 60).padStart(2, '0')}`;

const kinds: [string, string][] = [
  ['danger', 'Trouble'], ['politics', 'The city'], ['business', 'Business'],
  ['work', 'Work'], ['story', 'Arrangements'], ['intel', 'Word'],
  ['personal', 'Personal'], ['travel', 'Travel'], ['death', 'Deaths'],
];

export function LedgerScreen({world}: {world: Snapshot}) {
  const b = world.books;
  const [query, setQuery] = useState('');
  const [kind, setKind] = useState<string | null>(null);
  const [showAll, setShowAll] = useState(false);

  const needle = query.trim().toLowerCase();
  const history = useMemo(() => [...world.history].reverse().filter(r =>
    (!kind || r.kind === kind) &&
    (!needle || (r.title + ' ' + r.text).toLowerCase().includes(needle))
  ), [world.history, kind, needle]);

  const counts = kinds.map(([id, label]) =>
    [id, label, world.history.filter(r => r.kind === id).length] as const).filter(([, , n]) => n > 0);

  const shown = showAll || needle || kind ? history : history.slice(0, 12);

  // Days come out of the record's own minute, so the log reads as a diary.
  let lastDay = -1;

  return <section className="section-content">
    <div className="eyebrow">ACCOUNTS &amp; CONSEQUENCES</div>
    <h1 className="screen-title">The ledger</h1>

    {b && <div className="books">
      <div className="books-figures">
        <div><span>Coming in</span><b className="good">{money(b.income)}</b><small>a day from {b.holdings} {b.holdings === 1 ? 'business' : 'businesses'}</small></div>
        <div><span>Going out</span><b className="bad">{money(b.costs)}</b><small>a day, every day</small></div>
        <div><span>Net</span><b className={b.net < 0 ? 'bad' : 'good'}>{money(b.net)}</b><small>{b.net < 0 ? 'you are losing money' : 'a day to the good'}</small></div>
        <div><span>On hand</span><b>{money(b.cash)}</b><small>{b.sheltered > 0 ? `${money(b.sheltered)} a fine cannot reach` : 'all of it reachable'}</small></div>
        {b.lent > 0 && <div><span>Out on the street</span><b>{money(b.lent)}</b><small>{money(b.owed)} due back</small></div>}
        {b.offshore > 0 && <div><span>Outside the city</span><b>{money(b.offshore)}</b><small>survives you</small></div>}
      </div>
      <details className="books-lines">
        <summary>What the {money(b.costs)} a day is</summary>
        <ul>
          {b.lines.map(l => <li key={l.label}>
            <b>{l.label}</b>
            {l.detail && <small>{l.detail}</small>}
            <i>{money(l.amount)}</i>
          </li>)}
          <li className="total"><b>Every day</b><i>{money(b.costs)}</i></li>
        </ul>
      </details>
    </div>}

    <h2>What happened</h2>
    <div className="people-controls">
      <input type="search" value={query} placeholder={`Search ${world.history.length} entries…`}
        aria-label="Search the ledger" onChange={e => setQuery(e.target.value)}/>
      <div className="people-filters">
        <button aria-pressed={kind === null} onClick={() => setKind(null)}>Everything</button>
        {counts.map(([id, label, n]) => <button key={id} aria-pressed={kind === id}
          onClick={() => setKind(kind === id ? null : id)}>{label} <i>{n}</i></button>)}
      </div>
    </div>

    <div className="log">
      {shown.map(r => {
        const day = Math.floor(r.minute / 1440) + 1;
        const first = day !== lastDay;
        lastDay = day;
        return <div key={r.id}>
          {first && <h3 className="log-day">Day {day}</h3>}
          <article className={'log-row ' + r.kind}>
            <time>{time(r.minute).split(' · ')[1]}<br/>Life {r.life}</time>
            <div><h3>{r.title}{(r.count??1)>1&&<i className="again">{r.count} times</i>}</h3><p>{r.text}</p></div>
          </article>
        </div>;
      })}
      {!shown.length && <p className="nothing-here">Nothing in the ledger matches “{query}”.</p>}
      {!showAll && !needle && !kind && history.length > shown.length &&
        <button className="reveal-blocked" onClick={() => setShowAll(true)}>
          Show {history.length - shown.length} earlier entries
        </button>}
    </div>
  </section>;
}
