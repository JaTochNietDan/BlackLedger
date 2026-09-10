import {useEffect, useRef, useState} from 'react';
import {Card, pipOf, isRedSuit, knownCard, clothTable, clothColour, outsideBets, ballAngle, wheelAngle, wheelOrder, wheelPaint} from './cards';

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
// How long the ball is in the air, and how many turns it makes getting there.
// Nothing about the outcome: the core spun the pocket before this component was
// told anything, and the animation only takes its time arriving at it.
// How long the ball is in the air, and how many turns each part makes getting
// there. Nothing about the outcome: the core spun the pocket before this
// component was told anything, and the animation only takes its time arriving
// at it. The head turns a whole number of times so it ends where it started,
// which is what keeps the numbers on it upright at rest.
const FALL = 3200, TURNS = 5, HEAD_TURNS = 3;
// How far out the ball runs, which the stylesheet also has to agree with.
const BALL_TRACK = -80;

// One chip on the cloth. The core takes one bet a spin, so there is one chip:
// putting it somewhere else moves it rather than adding to it.
function Chip({amount, money}:{amount:number; money:(n:number)=>string}) {
  return <span className="chip" aria-hidden="true"><i/>{money(amount).replace('$', '')}</span>;
}

export function Wheel({wheel, stakes, money, spin, turn = 0}:{
  wheel:WheelState; stakes:{id:string; amount:number}[]; money:(n:number)=>string;
  spin:(stakeID:string, bet:string)=>void;
  // The world's revision, so a spin is animated once — and so two spins that
  // land in the same pocket are still two spins.
  turn?:number;
}) {
  const [bet, setBet] = useState('red');
  // Which of the house's stakes the chip is worth. The core offers the amounts;
  // this only says which one is on the cloth.
  const [stake, setStake] = useState(0);
  const [falling, setFalling] = useState(false);
  // The ball and the head are moved with the animation API rather than by a CSS
  // transition on a custom property. The transition was pinned at time zero and
  // never advanced — every render re-wrote the inline style and started it
  // again, so the ball sat at the top of the bowl however long you waited, and
  // an inline transform could not override the stuck transition either. An
  // animation is started once, by the spin, and nothing that re-renders can
  // interrupt it.
  const ball = useRef<HTMLSpanElement>(null);
  const head = useRef<HTMLDivElement>(null);
  const at = useRef({ball: 0, head: 0});
  // Every pocket this sitting has seen, newest first. The view is remembering
  // what the core told it, which is what the board of numbers over a real wheel
  // is: a record, not a prediction.
  const [run, setRun] = useState<number[]>([]);
  const seen = useRef(-1);

  useEffect(() => {
    if (!wheel.spun || turn === seen.current) return;
    seen.current = turn;
    const pocket = wheel.pocket ?? 0;
    const quiet = matchMedia('(prefers-reduced-motion: reduce)').matches;
    const was = at.current;
    const now = {
      ball: ballAngle(was.ball, pocket, quiet ? 0 : TURNS),
      head: quiet ? was.head : was.head - HEAD_TURNS * 360,
    };
    at.current = now;
    setRun(r => [pocket, ...r].slice(0, 14));
    // Where it ends up is written on the element, and the animation only covers
    // the journey there. That order matters: a browser that is not animating —
    // reduced motion, a tab nothing is drawing — shows the ball in the pocket
    // the core spun rather than frozen wherever the flight began.
    // Where it ends up is written on the element, and the animation only covers
    // the journey there. That order matters: a browser that is not animating —
    // reduced motion, or a tab whose clock is frozen because nothing is being
    // drawn — shows the ball in the pocket the core spun instead of holding the
    // first frame of a flight that never finishes.
    const ease = 'cubic-bezier(.12,.58,.16,1)';
    const flights: Animation[] = [];
    const fly = (el: HTMLElement, from: string, to: string) => {
      el.style.transform = to;
      if (!quiet) {
        flights.push(el.animate([{transform: from}, {transform: to}], {duration: FALL, easing: ease}));
      }
    };
    if (ball.current) {
      fly(ball.current, `rotate(${was.ball}deg) translateY(${BALL_TRACK}px)`,
        `rotate(${now.ball}deg) translateY(${BALL_TRACK}px)`);
    }
    if (head.current) {
      fly(head.current, `rotate(${was.head}deg)`, `rotate(${now.head}deg)`);
    }
    if (quiet) return;
    setFalling(true);
    // And the flight is cancelled when its time is up, which uncovers the
    // resting place underneath it. A timer runs even when a timeline does not,
    // so this is what makes a frozen tab still show the right answer.
    const done = setTimeout(() => {
      setFalling(false);
      flights.forEach(f => f.cancel());
    }, FALL);
    return () => clearTimeout(done);
  }, [turn, wheel.spun, wheel.pocket]);

  // While it is in the air the room does not know either. The number is the
  // core's from the moment it was spun; this only holds it back until the ball
  // is in the pocket, the way the table does.
  const landed = wheel.spun && !falling ? wheel.pocket ?? 0 : null;
  const chip = stakes[Math.min(stake, Math.max(0, stakes.length - 1))];
  const on = (id:string) => bet === id;
  const cell = (id:string, label:string|number, cls:string) => (
    <button key={id} className={'cloth-cell ' + cls + (on(id) ? ' picked' : '')}
            aria-pressed={on(id)} onClick={() => setBet(id)}>
      <span>{label}</span>
      {on(id) && chip && <Chip amount={chip.amount} money={money}/>}
    </button>
  );

  return (
    <div className="felt wheel-felt">
      <div className="wheel-table">
        <div className={'wheel-bowl' + (falling ? ' falling' : '')}>
          {/* The wheel head: one slice per pocket, in the pockets' own order,
              turning under the ball the way a real one does. */}
          <div className="wheel-head" ref={head} style={{background: wheelPaint()}}>
            {wheelOrder.map(n => (
              <span key={n} className={'pocket ' + clothColour(n) + (landed === n ? ' landed' : '')}
                    style={{['--at' as string]: `${wheelAngle(n)}deg`}}>{n}</span>
            ))}
            <span className="wheel-cone" aria-hidden="true"/>
          </div>
          {/* The ball rides the track and drops into the pocket the core spun. */}
          <span className="wheel-ball" ref={ball} aria-hidden="true"/>
          <div className="wheel-hub">
            {landed === null
              ? <small>{falling ? 'Round it goes' : 'No more bets'}</small>
              : <><b className={clothColour(landed)}>{landed}</b><small>{wheel.colour}</small></>}
          </div>
        </div>

        <div className="wheel-side">
          {run.length > 0 && <div className="wheel-run" aria-label="What this wheel has done">
            <small>THE LAST OF THEM</small>
            <div>{run.map((n, i) => <span key={i} className={'ran ' + clothColour(n)}>{n}</span>)}</div>
          </div>}
          {wheel.spun && landed !== null && (
            <p className={'felt-result' + (wheel.won ? ' won' : '')}>
              {wheel.bet} at {money(wheel.stake ?? 0)} — {wheel.won ? `paid ${wheel.pays} to 1` : 'gone'}
            </p>
          )}
        </div>
      </div>

      {/* The cloth, laid out the way a table is: the nought down the left, three
          rows of twelve, and the outside along the bottom. Your chip sits on
          whatever you last put it on — the core takes one bet a spin. */}
      <div className="cloth">
        <div className="cloth-numbers">
          {cell('number:0', 0, 'green zero')}
          <div className="cloth-grid">
            {clothTable().map((row, i) => (
              <div className="cloth-row" key={i}>
                {row.map(n => cell('number:' + n, n, clothColour(n)))}
              </div>
            ))}
          </div>
        </div>
        <div className="cloth-outside">
          {outsideBets().map(b => cell(b.id, b.label, 'outside' + (b.wide ? ' wide' : '') +
            (b.id === 'red' || b.id === 'black' ? ' ' + b.id : '')))}
        </div>
      </div>

      <div className="felt-actions chips">
        {stakes.length > 1 && <div className="chip-picker" role="group" aria-label="What the chip is worth">
          {stakes.map((s, i) => (
            <button key={s.id} className={'chip-choice' + (i === stake ? ' chosen' : '')}
                    aria-pressed={i === stake} onClick={() => setStake(i)}>
              <Chip amount={s.amount} money={money}/>
            </button>
          ))}
        </div>}
        {chip && <button className="spin-it" disabled={falling}
                         onClick={() => spin(chip.id, bet)}>
          {falling ? 'The ball is still going' : `Spin — ${money(chip.amount)} on ${bet.replace('number:', 'the ')}`}
        </button>}
      </div>
      <p className="felt-note">
        Thirty-seven pockets. The nought is neither colour and sits in no dozen, so it takes every
        bet on the outside — that is the whole of the house's advantage. Every payout here is the true one.
      </p>
    </div>
  );
}
