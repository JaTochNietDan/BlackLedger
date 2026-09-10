import {useEffect, useRef, useState} from 'react';
import type {Action, Presence, Record as Entry} from './types';
import {CardTable, Machine, Wheel} from './Tables';
import type {HandState, MachineState, WheelState} from './Tables';

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
export function isTableAction(id: string) {
  return id.startsWith('play:') || id.startsWith('wheel:') || id.startsWith('pull:') ||
    id === 'hit' || id === 'stand';
}

export function hasTables(actions: Action[]) {
  return actions.some(a => isTableAction(a.id));
}

type Game = 'cards' | 'wheel' | 'machine';

export function Casino({place, actions, people, hand, wheel, machine, cash, money, revision, records, act, onLeave}: {
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
  wheel: WheelState;
  cash: number;
  money: (n: number) => string;
  // The world's revision, so a result is added to the night once and not again
  // on every render.
  revision: number;
  records: Entry[];
  act: (command: {kind: string; target?: string; choice?: string}) => void;
  onLeave: () => void;
}) {
  // A hand that is already dealt is the game you are playing, whatever tab you
  // were last looking at.
  const [game, setGame] = useState<Game>(hand.playing ? 'cards' : 'wheel');
  useEffect(() => { if (hand.playing) setGame('cards') }, [hand.playing]);
  // A room with nothing but machines opens on them.
  useEffect(() => {
    if (!hand.playing && !actions.some(a => a.id.startsWith('play:') || a.id.startsWith('wheel:'))) setGame('machine');
  }, [actions, hand.playing]);

  // What the night has done so far. The ledger has all of this and always did;
  // what it did not have is a player watching one hand turn into the next.
  const [night, setNight] = useState<Entry[]>([]);
  const seen = useRef(-1);
  useEffect(() => {
    if (revision === seen.current) return;
    seen.current = revision;
    const fresh = records.filter(r => r.kind === 'personal' || r.kind === 'business');
    if (fresh.length) setNight(was => [...fresh, ...was].slice(0, 12));
  }, [revision, records]);

  const stakes = (prefix: string) => actions.filter(a => a.id.startsWith(prefix));
  const dealt = hand.playing;
  const machines = stakes('pull:');
  // A poolhall with a bandit against the wall is not a casino, and the room
  // should not offer a felt it does not have.
  const tables = actions.some(a => a.id.startsWith('play:') || a.id.startsWith('wheel:'));

  return <div className="modal-shade table-shade">
    <section className="casino" role="dialog" aria-modal="true" aria-label={'The tables at ' + place}>
      <header className="casino-head">
        <div>
          <div className="eyebrow">YOU ARE AT THE TABLES</div>
          <h2>{place}</h2>
        </div>
        <div className="casino-purse"><small>ON YOU</small><b>{money(cash)}</b></div>
        <button className="plain leave-table" onClick={onLeave} disabled={dealt}
                title={dealt ? 'Finish the hand first' : undefined}>
          {dealt ? 'Finish the hand to leave' : 'Get up and leave ↩'}
        </button>
      </header>

      <nav className="casino-games" aria-label="Games">
        {tables && <button aria-pressed={game === 'cards'} onClick={() => setGame('cards')}>Blackjack</button>}
        {tables && <button aria-pressed={game === 'wheel'} onClick={() => setGame('wheel')} disabled={dealt}>Roulette</button>}
        {machines.length > 0 && <button aria-pressed={game === 'machine'} onClick={() => setGame('machine')}
          disabled={dealt}>The machines</button>}
      </nav>

      <div className="casino-floor">
        <div className="casino-game">
          {game === 'cards'
            ? <div className="cards-panel">
                {(dealt || hand.settled) && <CardTable hand={hand} money={money} act={k => act({kind: k})}/>}
                {!dealt && <div className="felt sit-down">
                  <p className="felt-note">{hand.settled
                    ? 'That hand is finished. Put something down and the next one is dealt a card at a time.'
                    : 'Nothing on the table. Put something down and it is dealt a card at a time.'}</p>
                  <div className="felt-actions">
                    {stakes('play:').map(a => <button key={a.id} disabled={a.disabled} title={a.disabled ? a.reason : undefined}
                      onClick={() => act({kind: a.id})}>{a.label}</button>)}
                  </div>
                  {stakes('play:').every(a => a.disabled) && <p className="felt-refused">{stakes('play:')[0]?.reason}</p>}
                </div>}
              </div>
            : game === 'machine'
            ? <Machine machine={machine} money={money} turn={revision}
                       stakes={machines.filter(a => !a.disabled).map(a => ({id: a.id.slice(5), amount: a.asks ?? 0}))}
                       pull={id => act({kind: 'pull:' + id})}/>
            : <Wheel wheel={wheel} money={money} turn={revision}
                     stakes={stakes('wheel:').filter(a => !a.disabled).map(a => ({id: a.id.slice(6), amount: a.asks ?? 0}))}
                     spin={(stake, bet) => act({kind: 'wheel:' + stake, choice: bet})}/>}
        </div>

        <aside className="casino-night" aria-label="What the table has done">
          <h4>This sitting</h4>
          {night.length === 0
            ? <p className="subtle">Nothing played yet.</p>
            : night.map(r => <p key={r.id}><b>{r.title}</b>{r.text}</p>)}
        </aside>
      </div>
    </section>
  </div>;
}
