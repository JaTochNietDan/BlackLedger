import {useMemo, useState} from 'react';
import type {ReactElement} from 'react';
import type {Action, Snapshot, Record as CityRecord} from './types';

// Public history remains searchable here; the city account slip owns finances.
const time = (m: number) =>
  `Day ${Math.floor(m / 1440) + 1} · ${String(Math.floor((m % 1440) / 60)).padStart(2, '0')}:${String(m % 60).padStart(2, '0')}`;

const kinds: [string, string][] = [
  ['danger', 'Trouble'],
  ['politics', 'The city'],
  ['business', 'Business'],
  ['work', 'Work'],
  ['story', 'Arrangements'],
  ['intel', 'Word'],
  ['personal', 'Personal'],
  ['travel', 'Travel'],
  ['death', 'Deaths'],
];

export function LedgerScreen({
  world,
  actions = [],
  render,
}: {
  world: Snapshot;
  // What the police think of you is an account like any other, and the two
  // things that change it are not counter work: a detective who will lose some
  // paperwork does not take the money across a desk at the exchange, and
  // keeping your head down is not somewhere you go. Both used to be printed at
  // the exchange and are now offered wherever the player is standing.
  actions?: Action[];
  render?: (a: Action) => ReactElement;
}) {
  const [query, setQuery] = useState('');
  const [kind, setKind] = useState<string | null>(null);
  const [showAll, setShowAll] = useState(false);

  const needle = query.trim().toLowerCase();
  const history = useMemo(
    () =>
      [...world.history]
        .reverse()
        .filter(
          r =>
            (!kind || r.kind === kind) &&
            (!needle || (r.title + ' ' + r.text).toLowerCase().includes(needle)),
        ),
    [world.history, kind, needle],
  );

  const counts = kinds
    .map(([id, label]) => [id, label, world.history.filter(r => r.kind === id).length] as const)
    .filter(([, , n]) => n > 0);

  const shown = showAll || needle || kind ? history : history.slice(0, 12);

  // Days come out of the record's own minute, so the log reads as a diary.
  let lastDay = -1;

  return (
    <section className="section-content">
      <div className="eyebrow">PRIVATE MEMORANDA</div>
      <h1 className="screen-title">The ledger</h1>


      {!!render && actions.length > 0 && (
        <section className="anywhere-strip" aria-label="What you can do about your attention">
          <h4>
            What the police think
            <span>Neither of these happens at a counter. Do them from wherever you are.</span>
          </h4>
          <div className="actions compact">{actions.map(render)}</div>
        </section>
      )}

      <h2>What happened</h2>
      <div className="people-controls">
        <input
          type="search"
          value={query}
          placeholder={`Search ${world.history.length} entries…`}
          aria-label="Search the ledger"
          onChange={e => setQuery(e.target.value)}
        />
        <div className="people-filters">
          <button aria-pressed={kind === null} onClick={() => setKind(null)}>
            Everything
          </button>
          {counts.map(([id, label, n]) => (
            <button
              key={id}
              aria-pressed={kind === id}
              onClick={() => setKind(kind === id ? null : id)}
            >
              {label} <i>{n}</i>
            </button>
          ))}
        </div>
      </div>

      <div className="log">
        {shown.map(r => {
          const day = Math.floor(r.minute / 1440) + 1;
          const first = day !== lastDay;
          lastDay = day;
          return (
            <div key={r.id}>
              {first && <h3 className="log-day">Day {day}</h3>}
              <article className={'log-row ' + r.kind}>
                <time>
                  {time(r.minute).split(' · ')[1]}
                  <br />
                  Life {r.life}
                </time>
                <div>
                  <h3>
                    {r.title}
                    {(r.count ?? 1) > 1 && <i className="again">{r.count} times</i>}
                  </h3>
                  <p>{r.text}</p>
                </div>
              </article>
            </div>
          );
        })}
        {!shown.length && <p className="nothing-here">Nothing in the ledger matches “{query}”.</p>}
        {!showAll && !needle && !kind && history.length > shown.length && (
          <button className="reveal-blocked" onClick={() => setShowAll(true)}>
            Show {history.length - shown.length} earlier entries
          </button>
        )}
      </div>
    </section>
  );
}
