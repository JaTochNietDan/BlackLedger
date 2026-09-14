import {blackjackDealer} from './blackjackDealer';
import './casinoRoom.css';
import {useEffect, useRef, useState} from 'react';
import {playTable, roomTone} from './sound';
import type {Action, Presence, Record as Entry} from './types';
import {CardTable, Craps, Machine, Money, Wheel} from './Tables';
import type {DiceState, HandState, MachineState, WheelState} from './Tables';

// Sitting down at a table is not a thing you do out of the corner of a sidebar.
// The felt used to be drawn in the property panel beside the ownership figures
// and the list of everything else in the building, which reads as a picture of
// a game rather than a game: you dealt a hand in a column, the result went to
// the ledger, and the next screen looked the same as the last one.
//
// A table takes the room. You sit down, the city waits, every hand and every
// spin says what it did on the felt in front of you, and you are at the table
// until you get up and leave. That is what playing one is.
//
// Nothing here decides anything: the core deals, spins and moves the money.

// The actions a table owns. They are lifted out of the room's ordinary list so
// the building does not offer "Take another card" between restocking and hiring.
// The hand's `place` is the room's NAME, for printing on the felt — sending it
// as a command's target asks the core for an address called "The Blue Hour",
// which no city has, and every card played comes back "that action is not
// available here". The target is the id of the room the player is standing in.
// Taking the seat and leaving it belong here too: the room's own list of work
// should not offer "Get up and leave the tables" between restocking and hiring.
export function isTableAction(id: string) {
  return (
    id === 'play' ||
    id === 'wheel' ||
    id === 'pull' ||
    id === 'hit' ||
    id === 'stand' ||
    id === 'sit' ||
    id === 'rise' ||
    id === 'dice' ||
    id === 'roll'
  );
}

type Game = 'cards' | 'wheel' | 'dice' | 'machine';

export function Casino({
  motion,
  place,
  actions,
  people,
  hand,
  wheel,
  dice,
  machine,
  house,
  cash,
  money,
  revision,
  records,
  act,
  onLeave,
}: {
  motion: boolean;
  place: string;
  actions: Action[];
  // Who else is in the room. Not drawn here any more — the room itself shows
  // that, and a row of faces you cannot deal with is furniture on a screen that
  // is meant to be a game. Kept so the takeover can say the floor is busy if it
  // ever has something to say about it.
  people: Presence[];
  hand: HandState;
  // The machines against the wall. A room can have those and no tables at all,
  // which is the whole point of them.
  machine: MachineState;
  // What this room takes on one bet, which is the holder's decision and not
  // this component's.
  house: {
    limit?: number;
    machine?: number;
    least?: number;
    usual?: number;
    pull?: number;
    yours?: boolean;
    high?: number;
  };
  wheel: WheelState;
  // The dice, and the point if one is on.
  dice: DiceState;
  cash: number;
  money: (n: number) => string;
  // The world's revision, so a result is added to the night once and not again
  // on every render.
  revision: number;
  records: Entry[];
  act: (command: {
    kind: string;
    target?: string;
    choice?: string;
    amount?: number;
    chips?: {bet: string; amount: number}[];
  }) => void;
  onLeave: () => void;
}) {
  // A hand that is already dealt is the game you are playing, whatever tab you
  // were last looking at.
  // Which game you are looking at when you sit down. The nav lists Blackjack
  // first, so opening on the wheel showed a screen that did not match the tab
  // that was pressed.
  const [game, setGame] = useState<Game>(
    hand.playing || actions.some(a => a.id === 'play')
      ? 'cards'
      : actions.some(a => a.id === 'wheel')
        ? 'wheel'
        : 'machine',
  );
  // What the player is putting down, in dollars. Theirs to set, up to what the
  // room takes; the core refuses anything past it whatever this says.
  const [bet, setBet] = useState(0);
  const [pull, setPull] = useState(0);
  useEffect(() => {
    if (hand.playing) setGame('cards');
  }, [hand.playing]);
  // A point already on is the game you are playing, whatever tab was last open.
  useEffect(() => {
    if (dice.playing) setGame('dice');
  }, [dice.playing]);
  // A room with nothing but machines opens on them.
  useEffect(() => {
    if (!hand.playing && !actions.some(a => a.id === 'play' || a.id === 'wheel'))
      setGame('machine');
  }, [actions, hand.playing]);

  // The floor, while you are at it. Started when the takeover opens and stopped
  // when it closes: a noise that goes on after you have left the table is a
  // noise nobody asked for.
  useEffect(() => {
    roomTone(true);
    return () => roomTone(false);
  }, []);

  // What the night has done so far. The ledger has all of this and always did;
  // what it did not have is a player watching one hand turn into the next.
  const [presenting,setPresenting]=useState(false);
  const [night, setNight] = useState<Entry[]>([]);
  // Null until the first look, so sitting down is not treated as a thing that
  // just happened: it used to deal the last result into "this sitting" and,
  // once the tables had noises, play a card the moment the room opened.
  const seen = useRef<number | null>(null);
  useEffect(() => {
    if (seen.current === null) {
      seen.current = revision;
      return;
    }
    if (revision === seen.current) return;
    seen.current = revision;
    if (game === 'cards') playTable('card');
    if (game === 'dice') playTable('dice');
    const fresh = records.filter(r => r.kind === 'personal' || r.kind === 'business');
    if (fresh.length) setNight(was => [...fresh, ...was].slice(0, 12));
  }, [revision, records, game]);

  const at = (id: string) => actions.find(a => a.id === id);
  const dealt = hand.playing;
  const bandit = at('pull');
  const shooter = at('dice') || at('roll');
  const cards = at('play');
  const spin = at('wheel');
  // A poolhall with a bandit against the wall is not a casino, and the room
  // should not offer a felt it does not have.
  const tables = !!cards || !!spin;
  const least = house.least ?? 1;
  const limit = house.limit ?? 0;

  return (
    <div className="modal-shade table-shade">
      <section
        className="casino casino-house"
        role="dialog"
        aria-modal="true"
        aria-label={'The tables at ' + place}
      >
        <header className="casino-head">
          <div>
            <div className="eyebrow">{tables ? 'THE GAMING ROOM' : 'THE MACHINES'}</div>
            <h2>{place}</h2>
          </div>
          <div className="casino-purse">
            <small>CASH ON HAND</small>
            <b>{money(cash)}</b>
          </div>
          <button
            className="plain leave-table"
            onClick={onLeave}
            disabled={dealt}
            title={dealt ? 'Finish the hand first' : undefined}
          >
            {dealt ? 'Finish the hand to leave' : 'Get up and leave ↩'}
          </button>
        </header>

        <nav className="casino-games" aria-label="Games">
          {tables && (
            <button aria-pressed={game === 'cards'} onClick={() => setGame('cards')}>
              Blackjack
            </button>
          )}
          {tables && (
            <button
              aria-pressed={game === 'wheel'}
              onClick={() => setGame('wheel')}
              disabled={dealt}
            >
              Roulette
            </button>
          )}
          {!!shooter && (
            <button aria-pressed={game === 'dice'} onClick={() => setGame('dice')} disabled={dealt}>
              Craps
            </button>
          )}
          {!!bandit && (
            <button
              aria-pressed={game === 'machine'}
              onClick={() => setGame('machine')}
              disabled={dealt}
            >
              The machines
            </button>
          )}
        </nav>

        <div className="casino-floor">
          <div className="casino-game">
            {game === 'cards' ? (
              <div className="cards-panel">
                <CardTable hand={hand} place={place} dealer={blackjackDealer(people)} motion={motion} onPresent={setPresenting} money={money} act={k => act({kind: k})} />
                {!dealt && !presenting && (
                  <div className="felt sit-down">
                    <p className="felt-note">
                      {hand.settled
                        ? 'That hand is finished. Put something down and the next one is dealt a card at a time.'
                        : 'Nothing on the table. Put something down and it is dealt a card at a time.'}
                    </p>
                    <div className="felt-actions">
                      <Money
                        label="What you put down"
                        amount={bet || house.usual || least}
                        limit={limit}
                        least={least}
                        cash={cash}
                        onChange={setBet}
                      />
                      <button
                        disabled={cards?.disabled}
                        title={cards?.disabled ? cards.reason : undefined}
                        onClick={() => act({kind: 'play', amount: bet || house.usual || least})}
                      >
                        Deal {money(bet || house.usual || least)}
                      </button>
                    </div>
                    {cards?.disabled && <p className="felt-refused">{cards.reason}</p>}
                    {(bet || house.usual || least) >= (house.high ?? Infinity) && (
                      <p className="felt-note">
                        Money like that and the floor wants to know who you are.
                      </p>
                    )}
                  </div>
                )}
              </div>
            ) : game === 'dice' ? (
              <Craps
                motion={motion}
                onPresent={setPresenting}
                dice={dice}
                money={money}
                turn={revision}
                cash={cash}
                least={least}
                limit={limit}
                amount={bet || house.usual || least}
                onAmount={setBet}
                refused={shooter?.disabled ? shooter.reason : ''}
                play={(amount, choice) => act({kind: 'dice', amount, choice})}
                roll={() => act({kind: 'roll'})}
              />
            ) : game === 'machine' ? (
              <Machine
                motion={motion}
                onPresent={setPresenting}
                machine={machine}
                money={money}
                turn={revision}
                cash={cash}
                least={least}
                limit={house.machine ?? least}
                amount={pull || house.pull || least}
                onAmount={setPull}
                refused={bandit?.disabled ? bandit.reason : ''}
                pull={amount => act({kind: 'pull', amount})}
              />
            ) : (
              <Wheel
                wheel={wheel}
                motion={motion}
                onPresent={setPresenting}
                money={money}
                turn={revision}
                cash={cash}
                least={least}
                limit={limit}
                amount={bet || house.usual || least}
                onAmount={setBet}
                refused={spin?.disabled ? spin.reason : ''}
                spin={chips => act({kind: 'wheel', chips})}
              />
            )}
          </div>

          <aside className="casino-night" aria-label="What the table has done">
            <h4>This sitting</h4>
            {presenting ? <p role="status">The result will be recorded when play settles.</p> : night.length === 0 ? (
              <p className="subtle">Nothing played yet.</p>
            ) : (
              night.map(r => (
                <p key={r.id}>
                  <b>{r.title}</b>
                  {r.text}
                </p>
              ))
            )}
          </aside>
        </div>
      </section>
    </div>
  );
}
