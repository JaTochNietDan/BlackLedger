import {useEffect, useRef} from 'react';
import {BackRoom} from './Tables';
import type {CardsState} from './Tables';
import {playTable, roomTone} from './sound';
import type {Action} from './types';

// The room behind the poolhall takes the screen.
//
// It did not, and the note from the user was blunt about it: "it just lives in
// a small box above the action bar. That's silly stuff. We need to stop doing
// that in future and always dedicate these games to their own screen." A game
// drawn in a column beside the staffing figures and the supply count is a
// picture of a game. This is the same shape as the casino takeover: you go
// through, the city waits, and you are in there until you come back out.
//
// Nothing here decides anything. The core deals, the core moves the money, and
// the felt only arranges what it was told.

export function BackRoomScene({
  place,
  cards,
  seat,
  cash,
  money,
  act,
  onLeave,
}: {
  place: string;
  // The hand on the table, or nothing when the player has taken a seat and not
  // bought in yet — which is a real state: you can sit down and look at the
  // room before putting money in.
  cards: CardsState | null;
  // The action that starts a hand, so the ante can be named on the felt rather
  // than in the room's list behind it.
  seat: Action | undefined;
  cash: number;
  money: (n: number) => string;
  act: (command: {kind: string; choice?: string; amount?: number}) => void;
  onLeave: () => void;
}) {
  const dealt = !!cards && !cards.done;

  // The room behind the poolhall, out loud. It was the one table in the game
  // with nothing to hear: the casino has had a floor under it and a card on
  // every deal since the tables took the screen, and this had silence.
  //
  // The hum runs while the player is in here and stops when they come back
  // out, the same as the casino floor, because a noise that goes on after you
  // have left the room is a noise nobody asked for.
  useEffect(() => {
    roomTone(true);
    return () => roomTone(false);
  }, []);

  // A card for every card that lands, so the flop sounds like three of them
  // and the turn like one. Off the board the core sent rather than off a timer
  // of the interface's own: what is heard is what happened.
  const board = cards?.board?.length ?? 0;
  const seenBoard = useRef<number | null>(null);
  useEffect(() => {
    if (seenBoard.current === null || board < seenBoard.current) {
      // The first look, and the start of every new hand. Sitting down is not a
      // thing that just happened and a cleared board is not cards being taken
      // off the table.
      seenBoard.current = board;
      return;
    }
    const fresh = board - seenBoard.current;
    seenBoard.current = board;
    for (let i = 0; i < fresh; i++) setTimeout(() => playTable('card'), i * 90);
  }, [board]);

  // And chips when the middle grows, one for each raise's worth rather than
  // one for the lot: a big bet is a longer noise than a call.
  const pot = cards?.pot ?? 0;
  const ante = cards?.ante || 1;
  const seenPot = useRef<number | null>(null);
  useEffect(() => {
    if (seenPot.current === null || pot < seenPot.current) {
      seenPot.current = pot;
      return;
    }
    const added = pot - seenPot.current;
    seenPot.current = pot;
    if (added > 0) playTable('chips', Math.round(added / ante));
  }, [pot, ante]);

  return (
    <div className="modal-shade table-shade">
      <section
        className="casino back-room-scene"
        role="dialog"
        aria-modal="true"
        aria-label={'The back room at ' + place}
      >
        <header className="casino-head">
          <div>
            <div className="eyebrow">YOU ARE IN THE BACK ROOM</div>
            <h2>{place}</h2>
          </div>
          <div className="casino-purse">
            <small>ON YOU</small>
            <b>{money(cash)}</b>
          </div>
          <button
            className="plain leave-table"
            onClick={onLeave}
            disabled={dealt}
            title={dealt ? 'There is money of yours in the middle of the table' : undefined}
          >
            {dealt ? 'Finish the hand to leave' : 'Come back out ↩'}
          </button>
        </header>

        {cards ? (
          <BackRoom cards={cards} money={money} cash={cash} act={act} />
        ) : (
          <div className="back-room-empty">
            <p>
              Five cards each and one draw. There is no house in this game: the pot is what
              everybody put in, and it goes to the best hand at the table.
            </p>
            {seat && !seat.disabled ? (
              <button className="action primary" onClick={() => act({kind: 'cards', amount: 0})}>
                <span>
                  <strong>{seat.label}</strong>
                  <span className="desc">{seat.detail}</span>
                </span>
              </button>
            ) : (
              <p className="warning">{seat?.reason || 'There is no game to sit in on.'}</p>
            )}
          </div>
        )}
      </section>
    </div>
  );
}
