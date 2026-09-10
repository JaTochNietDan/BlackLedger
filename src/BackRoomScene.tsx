import {BackRoom} from './Tables';
import type {CardsState} from './Tables';
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
