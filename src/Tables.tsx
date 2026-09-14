import {planCards} from './blackjackPresentation';
import {BlackjackTable3D} from './BlackjackTable3D';
import {DiceTable3D} from './DiceTable3D';
import {DICE_ROLL_MS} from './dicePresentation';
import {SlotCabinet} from './SlotCabinet';
import {useEffect, useLayoutEffect, useRef, useState} from 'react';
import type {CSSProperties} from 'react';
import {playTable} from './sound';
import {
  Card,
  pipOf,
  isRedSuit,
  knownCard,
  clothTable,
  clothColour,
  outsideBets,
  ballAngle,
  wheelAngle,
  wheelOrder,
  wheelPaint,
  reelWindow,
  drumFaces,
} from './cards';

// The tables, drawn as tables. Blackjack was two numbers in a sentence and
// roulette was a button; both are games somebody sits down to play, and a game
// you play through a paragraph is a game you are being told about.
//
// Nothing here decides anything. The core deals the cards, spins the wheel and
// moves the money; this arranges what it says on a felt. Where the core does
// not know something — a card whose suit it never sent — this draws a card face
// down rather than choosing one.

export interface HandState {
  playing: boolean;
  settled?: boolean;
  won?: boolean;
  outcome?: string;
  where?: string;
  place?: string;
  stake?: number;
  player?: number;
  dealer?: number;
  cards?: number;
  mine?: Card[];
  theirs?: Card[];
}
export interface WheelState {
  spun: boolean;
  place?: string;
  stake?: number;
  bet?: string;
  pocket?: number;
  colour?: string;
  won?: boolean;
  pays?: number;
  back?: number;
  chips?: {bet: string; label: string; amount: number; won: boolean; back: number}[];
}

function PlayingCard({card, facedown}: {card?: Card; facedown?: boolean}) {
  if (facedown || !card || !knownCard(card)) {
    return (
      <div className="card facedown" aria-label="face down">
        <span />
      </div>
    );
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

function Row({cards, hidden}: {cards: Card[]; hidden: number}) {
  return (
    <div className="card-row">
      {cards.map((c, i) => (
        <PlayingCard key={i} card={c} />
      ))}
      {Array.from({length: hidden}, (_, i) => (
        <PlayingCard key={'x' + i} facedown />
      ))}
    </div>
  );
}

// What the player is putting down. One field, in dollars, bounded by what the
// room takes and what they are carrying — the games used to be two buttons at
// two fixed prices, which made the size of a bet a thing the room decided.
export function Money({
  amount,
  limit,
  least,
  cash,
  onChange,
  label,
}: {
  amount: number;
  limit: number;
  least: number;
  cash: number;
  onChange: (n: number) => void;
  label: string;
}) {
  const most = Math.max(least, Math.min(limit, cash));
  return (
    <label className="money-field">
      <span>{label}</span>
      <i>$</i>
      <input
        type="number"
        min={least}
        max={most}
        step={1}
        value={amount}
        aria-label={label}
        onChange={e =>
          onChange(Math.max(least, Math.min(most, Math.floor(Number(e.target.value) || 0))))
        }
      />
      <small>
        {least} to {most}
        {limit > cash ? ' — the house takes ' + limit : ''}
      </small>
    </label>
  );
}

// One chip on the cloth. The core takes one bet a spin, so there is one chip:
// putting it somewhere else moves it rather than adding to it.
function Chip({amount, money}: {amount: number; money: (n: number) => string}) {
  return (
    <span className="chip" aria-hidden="true">
      <i />
      {money(amount).replace('$', '')}
    </span>
  );
}

// The blackjack felt. The dealer's second card is not dealt until the player
// stands, so it is drawn face down: that is what is true, not a decoration.
export function CardTable({
  hand,
  motion = true,
  onPresent,
  money,
  act,
}: {
  hand: HandState;
  motion?: boolean;
  onPresent?: (active:boolean)=>void;
  money: (n: number) => string;
  act: (kind: string) => void;
}) {
  const previous=useRef(hand);
  const [presentation,setPresentation]=useState(()=>({plan:planCards(hand,hand,false),start:0,active:false}));
  const [reduced,setReduced]=useState(()=>matchMedia('(prefers-reduced-motion: reduce)').matches);
  useEffect(()=>{const media=matchMedia('(prefers-reduced-motion: reduce)');const change=()=>setReduced(media.matches);media.addEventListener('change',change);change();return()=>media.removeEventListener('change',change);},[]);
  const animate=motion&&!reduced;
  const handKey=JSON.stringify([hand.playing,hand.settled,hand.mine,hand.theirs]);
  useLayoutEffect(()=>{
    const changed=JSON.stringify(previous.current)!==JSON.stringify(hand);
    const plan=planCards(previous.current,hand,animate&&changed);
    previous.current=hand;
    const active=plan.duration>0;
    setPresentation({plan,start:performance.now(),active});onPresent?.(active);
    if(!active)return;
    const done=setTimeout(()=>{setPresentation(p=>({...p,active:false}));onPresent?.(false);},plan.duration);
    return()=>{clearTimeout(done);onPresent?.(false);};
  },[handKey,animate,onPresent]);
  if (!hand.playing && !hand.settled) return null;
  const mine = hand.mine ?? [],
    theirs = hand.theirs ?? [];
  const total = hand.player ?? 0;
  // A settled hand has nothing left to hide: the dealer has turned their card
  // over, and that is the moment the player sat down for.
  const over = !!hand.settled;
  return (
    <div className="felt card-felt">
      <div className="felt-head">
        <span>{hand.place}</span>
        <b>{money(hand.stake ?? 0)} down</b>
      </div>
      {/* The cloth itself, with the two seats on it and the money in the middle
          of the table where a stake actually sits. */}
      <BlackjackTable3D mine={mine} theirs={theirs} hidden={!over && theirs.length < 2 ? 1 : 0} presentation={presentation}/>
      {presentation.active ? <p role="status">The cards are being dealt…</p> : <div className="baize">
        <div className="seat dealer">
          <span className="seat-name">Dealer</span>
          <span className="hand-description">{theirs.map(c=>knownCard(c)?`${c.rank}${pipOf(c.suit)}`:"Face down").join(" · ")}{!over && theirs.length < 2 ? " · Face down" : ""}</span>
          <b className="seat-total">{hand.dealer ?? 0}</b>
        </div>
        <div className="baize-line" aria-hidden="true">
          {(hand.stake ?? 0) > 0 && <Chip amount={hand.stake ?? 0} money={money} />}
          <span>Blackjack pays 3 to 2</span>
        </div>
        <div className="seat mine">
          <span className="seat-name">You</span>
          <span className="hand-description">{mine.map(c=>knownCard(c)?`${c.rank}${pipOf(c.suit)}`:"Face down").join(" · ")}</span>
          <b className={'seat-total' + (total > 21 ? ' warning' : '')}>{total}</b>
        </div>
      </div>
      }
      {presentation.active ? null : over ? (
        <p className={'felt-result' + (hand.won ? ' won' : '')}>{hand.outcome}</p>
      ) : (
        <div className="felt-actions">
          <button onClick={() => act('hit')} disabled={total > 21}>
            Another card
          </button>
          <button onClick={() => act('stand')} disabled={total > 21}>
            Stand on {total}
          </button>
        </div>
      )}
      <p className="felt-note">
        The dealer draws to sixteen and stands on seventeen. A tie gives your money back.
      </p>
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
const FALL = 3200,
  TURNS = 5,
  HEAD_TURNS = 3;
// How far out the ball runs, which the stylesheet also has to agree with.
const BALL_TRACK = -80;

export function Wheel({
  wheel,
  motion = true,
  onPresent,
  money,
  spin,
  amount,
  least,
  limit,
  cash,
  onAmount,
  refused = '',
  turn = 0,
}: {
  wheel: WheelState;
  motion?: boolean;
  onPresent?: (active: boolean) => void;
  money: (n: number) => string;
  spin: (chips: {bet: string; amount: number}[]) => void;
  // What is going on the cloth, and what the room will take. The player names
  // the figure; the core refuses anything past the house limit whatever this
  // component thinks.
  amount: number;
  least: number;
  limit: number;
  cash: number;
  onAmount: (n: number) => void;
  refused?: string;
  // The world's revision, so a spin is animated once — and so two spins that
  // land in the same pocket are still two spins.
  turn?: number;
}) {
  // What is on the cloth. A real table takes as many chips as you can reach,
  // every one of them settled against the same pocket, and one bet a spin was
  // never how the game works. Clicking a spot adds a chip of whatever is in the
  // field; clicking it again takes one off.
  const [chips, setChips] = useState<Record<string, number>>({});
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
  const seen = useRef(turn);
  const recorded = useRef(-1);
  const [reduced, setReduced] = useState(() => matchMedia('(prefers-reduced-motion: reduce)').matches);
  useEffect(() => {
    const preference = matchMedia('(prefers-reduced-motion: reduce)');
    const changed = () => setReduced(preference.matches);
    preference.addEventListener('change', changed);
    changed();
    return () => preference.removeEventListener('change', changed);
  }, []);
  const animate = motion && !reduced;

  useLayoutEffect(() => {
    if (!wheel.spun) {setFalling(false); onPresent?.(false); return;}
    const fresh = turn !== seen.current;
    seen.current = turn;
    const pocket = wheel.pocket ?? 0;
    const moving = fresh && animate;
    const was = at.current;
    const now = {
      ball: ballAngle(was.ball, pocket, moving ? TURNS : 0),
      head: moving ? was.head - HEAD_TURNS * 360 : was.head,
    };
    at.current = now;
    const record = () => {
      if (recorded.current === turn) return;
      recorded.current = turn;
      setRun(r => [pocket, ...r].slice(0, 14));
    };
    const flights: Animation[] = [];
    const fly = (el: HTMLElement | null, from: string, to: string) => {
      if (!el) return;
      el.style.transform = to;
      if (moving) flights.push(el.animate([{transform: from}, {transform: to}], {
        duration: FALL, easing: 'cubic-bezier(.12,.58,.16,1)',
      }));
    };
    fly(ball.current, `rotate(${was.ball}deg) translateY(${BALL_TRACK}px)`,
      `rotate(${now.ball}deg) translateY(${BALL_TRACK}px)`);
    fly(head.current, `rotate(${was.head}deg)`, `rotate(${now.head}deg)`);
    setFalling(moving);
    onPresent?.(moving);
    if (!moving) {record(); return;}
    const done = setTimeout(() => {
      flights.forEach(f => f.cancel());
      record();
      setFalling(false);
      onPresent?.(false);
    }, FALL);
    return () => {
      clearTimeout(done);
      flights.forEach(f => f.cancel());
      onPresent?.(false);
    };
  }, [turn, wheel.spun, wheel.pocket, animate, onPresent]);

  // While it is in the air the room does not know either. The number is the
  // core's from the moment it was spun; this only holds it back until the ball
  // is in the pocket, the way the table does.
  const landed = wheel.spun && !falling ? (wheel.pocket ?? 0) : null;

  const down = Object.entries(chips).filter(([, n]) => n > 0);
  const total = down.reduce((n, [, v]) => n + v, 0);
  const place = (id: string, by: number) =>
    setChips(c => {
      const next = Math.max(0, (c[id] ?? 0) + by);
      const out = {...c};
      if (next === 0) delete out[id];
      else out[id] = next;
      return out;
    });
  const on = (id: string) => (chips[id] ?? 0) > 0;
  const cell = (id: string, label: string | number, cls: string) => (
    <button
      key={id}
      className={'cloth-cell ' + cls + (on(id) ? ' picked' : '')}
      aria-pressed={on(id)}
      title={
        on(id)
          ? `${money(chips[id])} down — right-click to take one off`
          : `Put ${money(amount)} on ${String(label)}`
      }
      onClick={() => place(id, amount)}
      onContextMenu={e => {
        e.preventDefault();
        place(id, -amount);
      }}
    >
      <span>{label}</span>
      {on(id) && <Chip amount={chips[id]} money={money} />}
    </button>
  );

  return (
    <div className="felt wheel-felt">
      <div className="wheel-layout">
        <div className="wheel-table">
          <div className={'wheel-bowl' + (falling ? ' falling' : '')}>
            {/* The wheel head: one slice per pocket, in the pockets' own order,
              turning under the ball the way a real one does. */}
            <div className="wheel-head" ref={head} style={{background: wheelPaint()}}>
              {wheelOrder.map(n => (
                <span
                  key={n}
                  className={'pocket ' + clothColour(n) + (landed === n ? ' landed' : '')}
                  style={{['--at' as string]: `${wheelAngle(n)}deg`}}
                >
                  {n}
                </span>
              ))}
              <span className="wheel-cone" aria-hidden="true" />
            </div>
            {/* The ball rides the track and drops into the pocket the core spun. */}
            <span className="wheel-ball" ref={ball} aria-hidden="true" />
            <div className="wheel-hub">
              {landed === null ? (
                <small>{falling ? 'Round it goes' : 'No more bets'}</small>
              ) : (
                <>
                  <b className={clothColour(landed)}>{landed}</b>
                  <small>{wheel.colour}</small>
                </>
              )}
            </div>
          </div>

          <div className="wheel-side">
            {run.length > 0 && (
              <div className="wheel-run" aria-label="What this wheel has done">
                <small>THE LAST OF THEM</small>
                <div>
                  {run.map((n, i) => (
                    <span key={i} className={'ran ' + clothColour(n)}>
                      {n}
                    </span>
                  ))}
                </div>
              </div>
            )}
            {wheel.spun && landed !== null && (
              <div className={'felt-result' + (wheel.won ? ' won' : '')}>
                {(wheel.chips ?? []).length > 1 ? (
                  <>
                    <p>
                      {money(wheel.stake ?? 0)} across {(wheel.chips ?? []).length} chips —{' '}
                      {(wheel.back ?? 0) > 0 ? `${money(wheel.back ?? 0)} back` : 'nothing back'}
                    </p>
                    <ul className="chip-run">
                      {(wheel.chips ?? []).map((c, i) => (
                        <li key={i} className={c.won ? 'won' : ''}>
                          {c.label} · {money(c.amount)}
                          {c.won ? ` → ${money(c.back)}` : ''}
                        </li>
                      ))}
                    </ul>
                  </>
                ) : (
                  <p>
                    {wheel.bet} at {money(wheel.stake ?? 0)} —{' '}
                    {wheel.won ? `paid ${wheel.pays} to 1` : 'gone'}
                  </p>
                )}
              </div>
            )}
          </div>
        </div>

        {/* The cloth, laid out the way a table is: the nought down the left,
          three rows of twelve, and the outside along the bottom. It sits beside
          the wheel rather than under it, which is where a real one is, and it
          takes as many chips as you want to put on it. */}
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
            {outsideBets().map(b =>
              cell(
                b.id,
                b.label,
                'outside' +
                  (b.wide ? ' wide' : '') +
                  (b.id === 'red' || b.id === 'black' ? ' ' + b.id : ''),
              ),
            )}
          </div>
        </div>
      </div>

      <div className="felt-actions chips">
        <Money
          label="A chip is worth"
          amount={amount}
          limit={limit}
          least={least}
          cash={cash}
          onChange={onAmount}
        />
        <button className="plain" disabled={down.length === 0} onClick={() => setChips({})}>
          Clear the cloth
        </button>
        <button
          className="spin-it"
          disabled={falling || !!refused || down.length === 0}
          title={refused || (down.length === 0 ? 'Put something on the cloth first' : undefined)}
          onClick={() => spin(down.map(([bet, n]) => ({bet, amount: n})))}
        >
          {falling
            ? 'The ball is still going'
            : down.length === 0
              ? 'No more bets'
              : `Spin — ${money(total)} on ${down.length === 1 ? down[0][0].replace('number:', 'the ') : `${down.length} chips`}`}
        </button>
      </div>
      {refused && <p className="felt-refused">{refused}</p>}
      <p className="felt-note">
        Thirty-seven pockets. The nought is neither colour and sits in no dozen, so it takes every
        bet on the outside — that is the whole of the house's advantage. Every payout here is the
        true one.
      </p>
    </div>
  );
}

export interface MachineState {
  pulled: boolean;
  place?: string;
  stake?: number;
  line?: string[];
  pays?: number;
  won?: boolean;
  stops: number;
  edge: number;
  two_cherries: number;
  one_cherry: number;
  strip: {id: string; face: string; stops: number; pays: number}[];
}

// The bandit. Three drums, a handle, and the paytable painted on the machine —
// which is the core's own strip, so a face the core cannot deal is a face
// nobody can see. The drums roll while the pull is being resolved and stop one
// after another, left to right, the way the real ones do; where they stop is
// the core's answer and nothing else.
export function Machine({
  machine,
  motion = true,
  onPresent,
  money,
  pull,
  amount,
  least,
  limit,
  cash,
  onAmount,
  refused = '',
  turn = 0,
}: {
  machine: MachineState;
  motion?: boolean;
  onPresent?: (active:boolean)=>void;
  money: (n: number) => string;
  pull: (amount: number) => void;
  // What goes in, and the most this machine takes — a tenth of what the tables
  // take, because a bandit is small money by design.
  amount: number;
  least: number;
  limit: number;
  cash: number;
  onAmount: (n: number) => void;
  refused?: string;
  turn?: number;
}) {
  const [reduced,setReduced]=useState(()=>matchMedia('(prefers-reduced-motion: reduce)').matches);
  useEffect(()=>{
    const preference=matchMedia('(prefers-reduced-motion: reduce)');
    const changed=()=>setReduced(preference.matches);
    preference.addEventListener('change',changed);changed();
    return()=>preference.removeEventListener('change',changed);
  },[]);
  const animate=motion&&!reduced;
  const [rolling, setRolling] = useState([false, false, false]);
  // Where each drum's column is sitting. The drums used to shake on the spot
  // with the answer already on them; they travel now, and this is how far each
  // one still has to go: TurnsADrum stops back at the moment the handle drops,
  // nought when it has landed.

  // A saved result is already settled when the machine is first opened.
  const seen = useRef(turn);
  useLayoutEffect(() => {
    if (!machine.pulled || !animate) {
      seen.current=turn;
      setRolling([false,false,false]);onPresent?.(false);
      return;
    }
    if (turn === seen.current) return;
    seen.current = turn;
    setRolling([true, true, true]);onPresent?.(true);
    playTable('handle');
    const stops = [0, 1, 2].map(i =>
      setTimeout(
        () => {
          playTable('reel');
          setRolling(r => r.map((was, at) => (at === i ? false : was)));
          // The tray, once the last drum is down and only if it paid.
          if (i === 2) onPresent?.(false);
          if (i === 2 && machine.won) playTable('coins', (machine.pays ?? 0) * 2);
        },
        700 + i * 450,
      ),
    );
    return () => {stops.forEach(clearTimeout);onPresent?.(false);};
  }, [turn, machine.pulled, animate, onPresent]);

  // Which of the house's machines you are standing at. A nickel machine and a
  // dollar machine are two different machines against the same wall.
  const strip = machine.strip || [];
  const line = machine.line ?? [];

  // What is on the drums. Three faces a drum, the middle one on the payline,
  // taken from the core's own strip so the case cannot show a symbol the odds
  // do not have — and showing where each drum is going to stop from the moment
  // the handle goes down, so nothing changes when it gets there.
  const windows = drumFaces(strip, line, !!machine.pulled);
  // The tray and the line under it wait for the last drum, because what a pull
  // paid is not a thing to announce while the drums are still going.
  const settled = machine.pulled && !rolling.some(Boolean);

  return (
    <div className="felt machine-felt">
      <div className="felt-head">
        <span>{machine.place ?? 'The machine'}</span>
        <b>keeps {machine.edge} in every 100</b>
      </div>

      <SlotCabinet strip={strip} line={line} pulled={!!machine.pulled} rolling={rolling} turn={turn}
        paid={settled&&machine.won?(machine.stake??0)*(machine.pays??0):0} disabled={rolling.some(Boolean)||!!refused} onPull={()=>pull(amount)}/>

      {settled && (
        <p className={'felt-result' + (machine.won ? ' won' : '')}>
          {machine.won
            ? `${windows.map(w => w[1]).join(' · ')} — ${money((machine.stake??0)*(machine.pays??0))} paid (${machine.pays} to 1)`
            : `${windows.map(w => w[1]).join(' · ')} — nothing. The machine keeps it.`}
        </p>
      )}

      <div className="bandit-stakes">
        <Money
          label="Into the slot"
          amount={amount}
          limit={limit}
          least={least}
          cash={cash}
          onChange={onAmount}
        />
        <button
          className="spin-it"
          disabled={rolling.some(Boolean) || !!refused}
          title={refused || undefined}
          onClick={() => pull(amount)}
        >
          {rolling.some(Boolean) ? 'The drums are still going' : `Pull for ${money(amount)}`}
        </button>
      </div>
      {refused && <p className="felt-refused">{refused}</p>}

      <table className="paytable">
        <tbody>
          {strip.map(s => (
            <tr key={s.id}>
              <th>
                {s.face} {s.face} {s.face}
              </th>
              <td>{s.pays} to 1</td>
            </tr>
          ))}
          <tr>
            <th>Two cherries</th>
            <td>{machine.two_cherries} to 1</td>
          </tr>
          <tr>
            <th>One cherry</th>
            <td>{machine.one_cherry} to 1</td>
          </tr>
        </tbody>
      </table>

      <p className="felt-note">
        Three drums of {machine.stops}, the same strip on each. Everything it pays is on the
        machine, and what it keeps is what is left over.
      </p>
    </div>
  );
}

// The dice. The third shape of decision this room offers: the cards ask you
// something on every card, the wheel asks once and then there is nothing to do,
// and this asks once and then makes you sit through a run of throws nobody can
// affect. The point sitting there between throws is the whole game.
export type DiceState = {
  playing: boolean;
  settled: boolean;
  place?: string;
  where?: string;
  bet?: string;
  bet_label?: string;
  stake?: number;
  point?: number;
  dice?: number[];
  total?: number;
  rolls?: number;
  won?: boolean;
  outcome?: string;
  bets?: {id: string; label: string; detail: string}[];
};

export function Craps({
  dice,
  motion = true,
  onPresent,
  money,
  play,
  roll,
  amount,
  least,
  limit,
  cash,
  onAmount,
  refused = '',
  turn = 0,
}: {
  dice: DiceState;
  motion?: boolean;
  onPresent?: (active:boolean)=>void;
  money: (n: number) => string;
  play: (amount: number, bet: string) => void;
  roll: () => void;
  amount: number;
  least: number;
  limit: number;
  cash: number;
  onAmount: (n: number) => void;
  refused?: string;
  turn?: number;
}) {
  const [bet, setBet] = useState('pass');
  const [shaking, setShaking] = useState(false);
  const [reduced,setReduced]=useState(()=>matchMedia('(prefers-reduced-motion: reduce)').matches);
  useEffect(()=>{const preference=matchMedia('(prefers-reduced-motion: reduce)');const changed=()=>setReduced(preference.matches);preference.addEventListener('change',changed);changed();return()=>preference.removeEventListener('change',changed);},[]);
  const animate=motion&&!reduced;
  const seen = useRef(turn);
  useLayoutEffect(() => {
    if (!animate || (!dice.playing && !dice.settled)) {seen.current=turn;setShaking(false);onPresent?.(false);return;}
    if (turn === seen.current) return;
    seen.current = turn;
    setShaking(true);onPresent?.(true);
    const stop = setTimeout(() => {setShaking(false);onPresent?.(false);}, DICE_ROLL_MS);
    return () => {clearTimeout(stop);onPresent?.(false);};
  }, [turn, dice.playing, dice.settled, animate, onPresent]);

  const live = dice.playing;
  const bets = dice.bets ?? [];
  const chosen = bets.find(b => b.id === (live ? dice.bet : bet));
  const faces = dice.dice ?? [];
  const point = dice.point ?? 0;

  return (
    <div className="felt dice-felt">
      <div className="felt-head">
        <span>{dice.place ?? 'The dice'}</span>
        <b>{shaking ? 'Dice in motion' : point > 0 ? `The point is ${point}` : 'Come out'}</b>
      </div>

      <div className="dice-table">
        <DiceTable3D faces={faces} rolling={shaking} turn={turn}/>
        {/* The point, kept where the box is on a real layout: a number that is
          on, and everybody at the table looking at it. */}
        <div className={'point-box' + (point > 0 ? ' on' : '')}>
          <small>POINT</small>
          <b>{!shaking && point > 0 ? point : '—'}</b>
        </div>
      </div>

      {!shaking && (dice.playing || dice.settled) && dice.outcome && (
        <p className={'felt-result' + (dice.settled && dice.won ? ' won' : '')}>{dice.outcome}</p>
      )}

      {live ? (
        <div className="felt-actions dice-actions">
          <p className="felt-note">
            {shaking ? 'The dice are still in motion.' : <>{money(dice.stake ?? 0)} on {(dice.bet_label ?? 'the line').toLowerCase()}. Nothing to
            decide now: the {point} or a seven.</>}
          </p>
          <button className="spin-it" disabled={shaking} onClick={roll}>
            {shaking ? 'The dice are still going' : `Throw again for the ${point}`}
          </button>
        </div>
      ) : (
        <>
          <div className="dice-bets" role="group" aria-label="What to back">
            {bets.map(b => (
              <button
                key={b.id}
                aria-pressed={bet === b.id}
                title={b.detail}
                onClick={() => setBet(b.id)}
              >
                {b.label}
              </button>
            ))}
          </div>
          {chosen && <p className="felt-note">{chosen.detail}</p>}
          <div className="felt-actions">
            <Money
              label="On the line"
              amount={amount}
              limit={limit}
              least={least}
              cash={cash}
              onChange={onAmount}
            />
            <button
              disabled={shaking || !!refused}
              title={refused || undefined}
              onClick={() => play(amount, bet)}
            >
              {shaking ? 'The dice are still going' : `Throw ${money(amount)}`}
            </button>
          </div>
        </>
      )}
      {refused && <p className="felt-refused">{refused}</p>}
    </div>
  );
}

// The back room. Not a felt like the others: there is no house on the other
// side of this table, so nothing is drawn as a dealer and nothing here has an
// edge. What matters on the screen is who is sitting there, what they bought in
// the draw, and what they said when the money went round — because those three
// are the whole of what the player has to read before deciding whether to pay.
// How many stops a drum turns through on a pull, and how tall one stop is on
// the screen. The height has to match the stylesheet: the column is moved by
// whole stops, so a drum that travels 31 pixels a stop comes to rest between
// two symbols.

export interface CardsState {
  place: string;
  ante: number;
  pot: number;
  mine: Card[];
  board: Card[];
  street: string;
  street_name: string;
  hand: string;
  seats: {
    who: string;
    name: string;
    threw: number;
    in: number;
    folded: boolean;
    said: string;
    sore?: number;
    moved?: number;
    stack: number;
    cards?: Card[];
    hand?: string;
  }[];
  bet: number;
  my_bet: number;
  facing: boolean;
  folded: boolean;
  done: boolean;
  outcome: string;
  won: number;
  // The sitting, as opposed to the hand: what is in front of the player, what
  // they brought to the table, how many hands have been dealt since, and
  // whether the night is finished.
  stack: number;
  buy_in: number;
  hands: number;
  over: boolean;
  ended: string;
  up: number;
}

// The back room, at hold'em. Two cards each, five in the middle, and four
// rounds of the only decision at a card table. What the player reads is the
// board, what each seat put in, and what they said when the money went round —
// which is the whole of what there is to go on, because nobody changes a card
// in this game and everybody can see the same five.
export function BackRoom({
  cards,
  money,
  cash,
  act,
}: {
  cards: CardsState;
  money: (n: number) => string;
  cash: number;
  act: (command: {kind: string; choice?: string; amount?: number}) => void;
}) {
  const [bet, setBet] = useState(0);
  const owed = Math.max(0, cards.bet - cards.my_bet);
  return (
    <div className="felt back-room">
      <div className="felt-head">
        <span>{cards.street_name === 'at the showdown' ? 'The showdown' : cards.street_name}</span>
        <b>{money(cards.pot)} in the middle</b>
        {/* What is in front of you, which is now the only money you can bet.
            A table where the pot is the one figure on the screen is a table
            you cannot decide anything at. */}
        <i className="felt-stack">
          {money(cards.stack)} in front of you
          {cards.hands > 1 ? ` · hand ${cards.hands}` : ''}
        </i>
      </div>
      <div className="baize">
        {cards.seats.map(s => (
          <div key={s.who} className={'seat player-seat' + (s.folded ? ' folded' : '')}>
            <span className="seat-name">{s.name}</span>
            <Row cards={s.cards ?? []} hidden={s.cards ? 0 : 2} />
            <small className="seat-said">
              {s.folded ? 'out' : s.said || 'waiting'}
              {s.in > 0 && !s.folded ? ` · ${money(s.in)} in` : ''}
              {s.hand ? ` · ${s.hand}` : ''}
            </small>
            {/* And what they have left. Somebody down to their last two antes
                plays differently, and you can see it coming. */}
            <small className="seat-stack">{money(s.stack)}</small>
            {/* What the night has cost them, and whether they are carrying
                anything about it. The money in this game belongs to somebody,
                so the screen says whose it was and what they think of you. */}
            {cards.done && !!s.moved && (
              <small className={'seat-moved' + (s.moved < 0 ? ' down' : '')}>
                {s.moved > 0 ? `+${money(s.moved)}` : `−${money(-s.moved)}`}
              </small>
            )}
            {!!s.sore && <small className="seat-sore">has something against you</small>}
          </div>
        ))}
        {/* The five in the middle belong to everybody, which is what makes this
            a game of what the other seats do rather than of what they drew. */}
        <div className="board">
          <span className="seat-name">The table</span>
          {/* Defensive on both counts: a board is never sent as null now, and
              a save written before it was published still opens. */}
          <Row cards={cards.board ?? []} hidden={5 - (cards.board ?? []).length} />
        </div>
        <div className="seat mine">
          <span className="seat-name">You</span>
          <Row cards={cards.mine} hidden={0} />
          <small className="seat-said">
            {cards.hand}
            {cards.my_bet > 0 ? ` · ${money(cards.my_bet)} in` : ''}
          </small>
        </div>
      </div>
      {cards.over ? (
        <>
          <p className={'felt-result' + (cards.up > 0 ? ' won' : '')}>{cards.outcome}</p>
          <p className="felt-ended">
            {cards.ended} You brought {money(cards.buy_in)} and played{' '}
            {cards.hands === 1 ? 'one hand' : `${cards.hands} hands`}.
          </p>
        </>
      ) : cards.done ? (
        <>
          <p className={'felt-result' + (cards.won > 0 ? ' won' : '')}>{cards.outcome}</p>
          {/* The night goes on. This is the whole of what a sitting is: the
              next hand is one decision, not getting up and sitting down
              again. */}
          <div className="felt-actions">
            <button className="primary" onClick={() => act({kind: 'deal'})}>
              Deal the next hand
            </button>
            <button onClick={() => act({kind: 'cashout'})}>
              Pick up {money(cards.stack)} and leave
            </button>
          </div>
        </>
      ) : cards.facing ? (
        <div className="felt-actions">
          <button onClick={() => act({kind: 'call'})} disabled={owed > cards.stack}>
            Call the {money(owed)}
          </button>
          <button onClick={() => act({kind: 'fold'})}>Throw the hand in</button>
        </div>
      ) : (
        <div className="felt-actions">
          {/* Bounded by the chips in front of you rather than by your cash:
              nobody at a table reaches into their coat. */}
          <Money
            amount={bet}
            limit={Math.min(500, cards.stack)}
            least={0}
            cash={cards.stack}
            onChange={setBet}
            label="Bet"
          />
          <button onClick={() => act({kind: 'bet', amount: bet})}>
            {bet > 0 ? `Bet ${money(bet)}` : 'Check'}
          </button>
          <button onClick={() => act({kind: 'fold'})}>Throw the hand in</button>
        </div>
      )}
      <p className="felt-note">
        No house in this game. Two cards each and five on the table; the pot is what everybody put
        in, and it goes to the best five cards anybody can make.
      </p>
    </div>
  );
}
