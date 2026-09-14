import type {Snapshot} from './types';
const money = (n: number) => (n < 0 ? '−$' : '$') + Math.abs(Math.floor(n)).toLocaleString();
export function CityAccounts({world}: {world: Snapshot}) {
 const b=world.books;
 if(!b) return null;
 return <details className="city-accounts"><summary data-shortcut="8" aria-keyshortcuts="8"><span>Daily accounts</span><span>In <b>{money(b.income)}</b></span><span>Out <b>{money(b.costs)}</b></span><span>Net <b className={b.net < 0 ? 'bad' : 'good'}>{money(b.net)}</b></span></summary>
      {b && (
        <div className="books">
          <div className="books-figures">
            <div>
              <span>Coming in</span>
              <b className="good">{money(b.income)}</b>
              <small>
                a day from {b.holdings} {b.holdings === 1 ? 'business' : 'businesses'}
              </small>
            </div>
            <div>
              <span>Going out</span>
              <b className="bad">{money(b.costs)}</b>
              <small>a day, every day</small>
            </div>
            <div>
              <span>Net</span>
              <b className={b.net < 0 ? 'bad' : 'good'}>{money(b.net)}</b>
              <small>{b.net < 0 ? 'you are losing money' : 'a day to the good'}</small>
            </div>
            <div>
              <span>On hand</span>
              <b>{money(b.cash)}</b>
              <small>
                {b.sheltered > 0
                  ? `${money(b.sheltered)} a fine cannot reach`
                  : 'all of it reachable'}
              </small>
            </div>
            {b.lent > 0 && (
              <div>
                <span>Out on the street</span>
                <b>{money(b.lent)}</b>
                <small>{money(b.owed)} due back</small>
              </div>
            )}
            {b.offshore > 0 && (
              <div>
                <span>Outside the city</span>
                <b>{money(b.offshore)}</b>
                <small>survives you</small>
              </div>
            )}
          </div>
          {/* What is already unpaid, above the breakdown, because the day's
              costs are what is owed and this is where it has already gone
              wrong. Wages are the last thing the night gives up and the only
              one with people on the other side of it: a week of this and
              somebody stops coming in, and nobody new takes the job until it
              is paid. */}
          {!!b.behind?.length && (
            <ul className="books-behind" aria-label="Premises behind on wages">
              {b.behind.map(w => (
                <li key={w.id} className={w.shut ? 'warning shut' : 'warning'}>
                  <b>{w.place}</b>
                  <small>
                    {w.nights === 1 ? 'one night unpaid' : w.nights + ' nights unpaid'}
                    {w.shut
                      ? ' · nobody will take the job until it is paid'
                      : ' · ' + w.hands + ' of ' + w.positions + ' still coming in'}
                  </small>
                </li>
              ))}
            </ul>
          )}
          <details className="books-lines">
            <summary>What the {money(b.costs)} a day is</summary>
            <ul>
              {b.lines.map(l => (
                <li key={l.label}>
                  <b>{l.label}</b>
                  {l.detail && <small>{l.detail}</small>}
                  <i>{money(l.amount)}</i>
                </li>
              ))}
              <li className="total">
                <b>Every day</b>
                <i>{money(b.costs)}</i>
              </li>
            </ul>
          </details>
        </div>
      )}

</details>;
}
