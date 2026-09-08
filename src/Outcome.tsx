import {useEffect,useRef,useState} from 'react';
import type {Snapshot,Record as CityRecord} from './types';

// The result of an action was a strip along the bottom of the screen, below
// the fold, in the place a page puts a status bar. It is not a status: it is
// the answer to the thing the player just decided, and it is the only reason
// they pressed the button. It belongs at the top of whatever pane they pressed
// it in, where their eye already is.

const weight = (r: CityRecord) => ({death: 6, danger: 5, intel: 4, politics: 3, story: 2, personal: 1}[r.kind] ?? 0);
const money = (n: number) => (n < 0 ? '−' : '+') + '$' + Math.abs(n).toLocaleString();
const hours = (m: number) => m >= 120 ? `${Math.round(m / 60)} hours` : m >= 60 ? `${Math.floor(m / 60)}h ${m % 60}m` : `${m} min`;

export function Outcome({world, onLedger}: {world: Snapshot; onLedger: () => void}) {
  const result = world.last_result;
  const records = result?.records || [];
  const [open, setOpen] = useState(false);
  const seen = useRef('');

  // A new result is a new thing to read, so the fold closes and the panel
  // announces itself to anybody listening rather than watching.
  const key = (result?.action || '') + world.revision;
  useEffect(() => { if (key !== seen.current) { seen.current = key; setOpen(false) } }, [key]);

  if (!result || (!records.length && !result.action)) return null;

  const headline = records.length ? [...records].sort((a, b) => weight(b) - weight(a))[0] : null;
  const rest = headline ? records.filter(r => r.id !== headline.id) : [];
  const grave = headline && (headline.kind === 'death' || headline.kind === 'danger');

  // Only the figures that actually moved. A row of zeroes is noise.
  const moved: [string, string, boolean][] = [];
  if (result.cash) moved.push(['', money(result.cash), result.cash < 0]);
  if (result.respect) moved.push(['respect', (result.respect > 0 ? '+' : '−') + Math.abs(result.respect), result.respect < 0]);
  if (result.heat) moved.push(['attention', (result.heat > 0 ? '+' : '−') + Math.abs(result.heat), result.heat > 0]);
  if (result.health) moved.push(['health', (result.health > 0 ? '+' : '−') + Math.abs(result.health), result.health < 0]);
  if (result.elapsed) moved.push(['', hours(result.elapsed), false]);

  return <section className={'outcome-panel' + (grave ? ' grave' : '')} role="status" aria-live="polite">
    <div className="outcome-head">
      <span className="eyebrow">YOU DID THIS</span>
      <b>{result.action || 'Time passed'}</b>
    </div>
    {moved.length > 0 && <ul className="outcome-figures">
      {moved.map(([label, value, bad], i) => <li key={i} className={bad ? 'bad' : 'good'}>
        <b>{value}</b>{label && <small>{label}</small>}
      </li>)}
    </ul>}
    {headline
      ? <div className="outcome-said"><b>{headline.title}</b><p>{headline.text}</p></div>
      : <p className="outcome-quiet">Nothing came of it that anybody wrote down.</p>}
    {rest.length > 0 && <>
      <button className="reveal-blocked" aria-expanded={open} onClick={() => setOpen(o => !o)}>
        {open ? 'Hide' : 'Show'} {rest.length} other {rest.length === 1 ? 'thing that happened' : 'things that happened'} while you did it
      </button>
      {open && <ul className="outcome-rest">{rest.map(r => <li key={r.id} className={r.kind}>
        <b>{r.title}</b><p>{r.text}</p></li>)}</ul>}
    </>}
    <button className="plain outcome-ledger" onClick={onLedger}>The whole ledger ↗</button>
  </section>;
}
