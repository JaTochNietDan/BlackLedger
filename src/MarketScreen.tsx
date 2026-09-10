import type {Snapshot} from './types';

// The underground market was a block inside the ledger.
//
// It does not belong there. The ledger answers "what am I worth and what is
// this costing me" — it is accounts. A price is not an account; it is a reason
// to go somewhere. Somebody checking whether moonshine is worth moving today is
// not doing bookkeeping, they are deciding where to spend the afternoon, and
// burying that under the daily upkeep made both jobs harder.
//
// So the market is its own page: what things cost against what they usually
// cost, what is in your hands, what holding it is costing you in attention, and
// where the trade is actually done.

const money = (n: number) => (n < 0 ? '−$' : '$') + Math.abs(Math.floor(n)).toLocaleString();

// Where each good changes hands. The core knows the trade rules; this is the
// one thing the player needs and cannot read off a price.
const counter: Record<string, string> = {
  moonshine: 'The waterfront deals in it. Pier 14 and the docks.',
  cigarettes: 'Mercer Exchange, openly enough.',
  arms: 'Nobody sells these across a counter. Ask at the garage.',
};

export function MarketScreen({world}: {world: Snapshot}) {
  const goods = world.goods || [];
  const held = world.player.stock || {};
  const carrying = goods.filter(g => (held[g.id] || 0) > 0);
  // What the stock in your hands is costing you every day it stays there.
  const attention = carrying.reduce((sum, g) => sum + (held[g.id] || 0) * g.heat, 0);

  return (
    <section className="section-content">
      <div className="eyebrow">PRICES &amp; WHAT YOU ARE HOLDING</div>
      <h1 className="screen-title">The underground market</h1>
      <p className="subtle market-note">
        Prices move whether or not anybody is watching them. Stock is only worth what somebody will
        pay for it today, and every day it sits in your hands it is drawing attention.
      </p>

      <div className="market-board">
        {goods.map(g => {
          const move = g.base ? Math.round(((g.price - g.base) / g.base) * 100) : 0;
          const have = held[g.id] || 0;
          return (
            <article key={g.id} className={'market-good' + (have ? ' holding' : '')}>
              <header>
                <b>{g.name}</b>
                <span className={move > 4 ? 'up' : move < -4 ? 'down' : ''}>
                  {money(g.price)}
                  <small> / {g.unit}</small>
                </span>
              </header>
              <p className="market-move">
                {move === 0
                  ? 'About what it usually goes for'
                  : `${Math.abs(move)}% ${move > 0 ? 'above' : 'below'} the usual ${money(g.base)}`}
              </p>
              <p className="market-holding">
                {have ? (
                  <>
                    <b>{have}</b> in your hands · {money(have * g.price)} at today's price
                  </>
                ) : (
                  <span className="subtle">none held</span>
                )}
              </p>
              <p className="market-where">
                {counter[g.id] || 'Traded quietly, where such things are traded.'}
              </p>
            </article>
          );
        })}
      </div>

      {carrying.length > 0 && (
        <div className="market-risk">
          <h2>What holding it costs</h2>
          <p>
            {carrying.map(g => `${held[g.id]} ${g.name.toLowerCase()}`).join(', ')} —{' '}
            <b>{attention}</b> attention a day, every day, until it moves. Stock found in a raid is
            stock lost, and the raid is likelier the longer it sits.
          </p>
        </div>
      )}

      {goods.length === 0 && (
        <p className="nothing-here">Nothing is being traded that you know of.</p>
      )}
    </section>
  );
}
