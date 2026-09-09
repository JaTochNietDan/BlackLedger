import {useState} from 'react';
import {Card, pipOf, isRedSuit, knownCard, clothRows, clothColour, outsideBets, wheelAngle, wheelOrder} from './cards';

// The tables, drawn as tables. Blackjack was two numbers in a sentence and
// roulette was a button; both are games somebody sits down to play, and a game
// you play through a paragraph is a game you are being told about.
//
// Nothing here decides anything. The core deals the cards, spins the wheel and
// moves the money; this arranges what it says on a felt. Where the core does
// not know something — a card whose suit it never sent — this draws a card face
// down rather than choosing one.

export interface HandState {playing:boolean; settled?:boolean; won?:boolean; outcome?:string; where?:string; place?:string; stake?:number; player?:number; dealer?:number; cards?:number; mine?:Card[]; theirs?:Card[]}
export interface WheelState {spun:boolean; place?:string; stake?:number; bet?:string; pocket?:number; colour?:string; won?:boolean; pays?:number}

function PlayingCard({card, facedown}:{card?:Card; facedown?:boolean}) {
  if (facedown || !card || !knownCard(card)) {
    return <div className="card facedown" aria-label="face down"><span/></div>;
  }
  const red = isRedSuit(card.suit);
  return (
    <div className={'card' + (red ? ' red' : '')} aria-label={`${card.rank} of ${card.suit}`}>
      <b>{card.rank}</b>
      <i>{pipOf(card.suit)}</i>
      <b className="upside">{card.rank}</b>
    </div>
  );
}

function Row({cards, hidden}:{cards:Card[]; hidden:number}) {
  return (
    <div className="card-row">
      {cards.map((c, i) => <PlayingCard key={i} card={c}/>)}
      {Array.from({length: hidden}, (_, i) => <PlayingCard key={'x' + i} facedown/>)}
    </div>
  );
}

// The blackjack felt. The dealer's second card is not dealt until the player
// stands, so it is drawn face down: that is what is true, not a decoration.
export function CardTable({hand, money, act}:{hand:HandState; money:(n:number)=>string; act:(kind:string)=>void}) {
  if (!hand.playing && !hand.settled) return null;
  const mine = hand.mine ?? [], theirs = hand.theirs ?? [];
  const total = hand.player ?? 0;
  // A settled hand has nothing left to hide: the dealer has turned their card
  // over, and that is the moment the player sat down for.
  const over = !!hand.settled;
  return (
    <div className="felt">
      <div className="felt-head"><span>{hand.place}</span><b>{money(hand.stake ?? 0)} down</b></div>
      <div className="seat">
        <span className="seat-name">Dealer</span>
        <Row cards={theirs} hidden={!over && theirs.length < 2 ? 1 : 0}/>
        <b className="seat-total">{hand.dealer ?? 0}</b>
      </div>
      <div className="seat">
        <span className="seat-name">You</span>
        <Row cards={mine} hidden={0}/>
        <b className={'seat-total' + (total > 21 ? ' warning' : '')}>{total}</b>
      </div>
      {over
        ? <p className={'felt-result' + (hand.won ? ' won' : '')}>{hand.outcome}</p>
        : <div className="felt-actions">
            <button onClick={() => act('hit')} disabled={total > 21}>Another card</button>
            <button onClick={() => act('stand')} disabled={total > 21}>Stand on {total}</button>
          </div>}
      <p className="felt-note">The dealer draws to sixteen and stands on seventeen. A tie gives your money back.</p>
    </div>
  );
}

// The wheel and the cloth. Choosing a bet is a decision, so it is made here and
// sent with the spin; the core is still the only thing that decides where the
// ball lands.
export function Wheel({wheel, stakes, money, spin}:{
  wheel:WheelState; stakes:{id:string; amount:number}[]; money:(n:number)=>string;
  spin:(stakeID:string, bet:string)=>void;
}) {
  const [bet, setBet] = useState('red');
  const landed = wheel.spun ? wheel.pocket ?? 0 : null;
  return (
    <div className="felt wheel-felt">
      <div className="wheel-face" style={{['--drop' as string]: `${landed === null ? 0 : wheelAngle(landed)}deg`}}>
        <div className="wheel-rim">
          {wheelOrder.map(n => (
            <span key={n} className={'pocket ' + clothColour(n) + (landed === n ? ' landed' : '')}
                  style={{['--at' as string]: `${wheelAngle(n)}deg`}}>{n}</span>
          ))}
        </div>
        <div className="wheel-hub">
          {landed === null
            ? <small>Nothing spun yet</small>
            : <><b className={clothColour(landed)}>{landed}</b><small>{wheel.colour}</small></>}
        </div>
      </div>

      {wheel.spun && (
        <p className={'felt-result' + (wheel.won ? ' won' : '')}>
          {wheel.bet} at {money(wheel.stake ?? 0)} — {wheel.won ? `paid ${wheel.pays} to 1` : 'gone'}
        </p>
      )}

      <div className="cloth">
        <button className={'cloth-cell green zero' + (bet === 'number:0' ? ' picked' : '')}
                onClick={() => setBet('number:0')}>0</button>
        <div className="cloth-grid">
          {clothRows().map((row, i) => (
            <div className="cloth-row" key={i}>
              {row.map(n => (
                <button key={n} className={'cloth-cell ' + clothColour(n) + (bet === 'number:' + n ? ' picked' : '')}
                        onClick={() => setBet('number:' + n)}>{n}</button>
              ))}
            </div>
          ))}
        </div>
        <div className="cloth-outside">
          {outsideBets().map(b => (
            <button key={b.id} className={'cloth-cell outside' + (b.wide ? ' wide' : '') + (bet === b.id ? ' picked' : '')}
                    onClick={() => setBet(b.id)}>{b.label}</button>
          ))}
        </div>
      </div>

      <div className="felt-actions">
        {stakes.map(s => (
          <button key={s.id} onClick={() => spin(s.id, bet)}>
            Spin {money(s.amount)}
          </button>
        ))}
      </div>
      <p className="felt-note">
        Thirty-seven pockets. The nought is neither colour and sits in no dozen, so it takes every
        bet on the outside — that is the whole of the house's advantage. Every payout here is the true one.
      </p>
    </div>
  );
}
