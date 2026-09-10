import {useEffect, useRef, useState} from 'react';
import type {Snapshot, Record as CityRecord} from './types';

// The result of an action was a strip along the bottom of the screen, below
// the fold, in the place a page puts a status bar. It is not a status: it is
// the answer to the thing the player just decided, and it is the only reason
// they pressed the button. It belongs at the top of whatever pane they pressed
// it in, where their eye already is.

const weight = (r: CityRecord) =>
  ({death: 6, danger: 5, intel: 4, politics: 3, story: 2, personal: 1})[r.kind] ?? 0;
const money = (n: number) => (n < 0 ? '−' : '+') + '$' + Math.abs(n).toLocaleString();
// How long something took, in the units a person would use for it. Sitting out
// a sentence the game itself calls "Do the 5 days" was reported as "120 hours",
// which is true and tells the player nothing: past a day, days are the unit.
const hours = (m: number) => {
  if (m >= 1440) {
    const days = Math.floor(m / 1440),
      rest = Math.round((m % 1440) / 60);
    const said = days === 1 ? '1 day' : `${days} days`;
    return rest ? `${said} ${rest}h` : said;
  }
  return m >= 120
    ? `${Math.round(m / 60)} hours`
    : m >= 60
      ? `${Math.floor(m / 60)}h ${m % 60}m`
      : `${m} min`;
};

export function Outcome({world, onLedger}: {world: Snapshot; onLedger: () => void}) {
  const result = world.last_result;
  const records = result?.records || [];
  const [open, setOpen] = useState(false);
  const seen = useRef('');

  // A new result is a new thing to read, so the fold closes and the panel
  // announces itself to anybody listening rather than watching.
  const key = (result?.action || '') + world.revision;
  useEffect(() => {
    if (key !== seen.current) {
      seen.current = key;
      setOpen(false);
    }
  }, [key]);

  // The band keeps its place whether or not anything has happened, so the page
  // never jumps under the player's hand between one action and the next.
  if (!result || (!records.length && !result.action)) {
    return (
      <section className="outcome-band empty" aria-hidden="true">
        <span className="eyebrow">WHAT YOU DID</span>
        <b className="outcome-idle">Nothing yet. The clock is paused.</b>
      </section>
    );
  }

  const headline = records.length ? [...records].sort((a, b) => weight(b) - weight(a))[0] : null;
  const rest = headline ? records.filter(r => r.id !== headline.id) : [];
  const grave = headline && (headline.kind === 'death' || headline.kind === 'danger');

  // Only the figures that actually moved. A row of zeroes is noise.
  const moved: [string, string, boolean][] = [];
  if (result.cash) moved.push(['', money(result.cash), result.cash < 0]);
  if (result.respect)
    moved.push([
      'respect',
      (result.respect > 0 ? '+' : '−') + Math.abs(result.respect),
      result.respect < 0,
    ]);
  if (result.heat)
    moved.push([
      'attention',
      (result.heat > 0 ? '+' : '−') + Math.abs(result.heat),
      result.heat > 0,
    ]);
  if (result.health)
    moved.push([
      'health',
      (result.health > 0 ? '+' : '−') + Math.abs(result.health),
      result.health < 0,
    ]);
  if (result.elapsed) moved.push(['', hours(result.elapsed), false]);

  return (
    <section
      className={'outcome-band' + (grave ? ' grave' : '') + (open ? ' open' : '')}
      role="status"
      aria-live="polite"
    >
      <div className="outcome-line">
        <div className="outcome-did">
          <span className="eyebrow">WHAT YOU DID</span>
          <b>{result.action || 'Time passed'}</b>
        </div>
        {moved.length > 0 && (
          <ul className="outcome-figures">
            {moved.map(([label, value, bad], i) => (
              <li key={i} className={bad ? 'bad' : 'good'}>
                <b>{value}</b>
                {label && <small>{label}</small>}
              </li>
            ))}
          </ul>
        )}
        <div className="outcome-said">
          {headline ? (
            <>
              <b>{headline.title}</b>
              <span>{headline.text}</span>
            </>
          ) : (
            <span className="outcome-quiet">Nothing came of it that anybody wrote down.</span>
          )}
        </div>
        <div className="outcome-tools">
          {rest.length > 0 && (
            <button className="plain" aria-expanded={open} onClick={() => setOpen(o => !o)}>
              {open ? 'Less' : `+${rest.length} more`}
            </button>
          )}
          <button className="plain" onClick={onLedger}>
            Ledger ↗
          </button>
        </div>
      </div>
      {open && rest.length > 0 && (
        <ul className="outcome-rest">
          {rest.map(r => (
            <li key={r.id} className={r.kind}>
              <b>{r.title}</b>
              <p>{r.text}</p>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
