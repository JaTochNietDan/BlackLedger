import {useState} from 'react';
import {formatActionAmount} from './sumAmount';
import type {Action, Command} from './types';

// Money into and out of a business used to move in lots somebody else chose:
// $250 behind the tables, $500 out of the city, and the whole account home in
// one go. A card with a number on it is a decision the player makes rather
// than one they accept, so any action the core marks with a sum gets a field.
//
// It cannot be a button with an input inside it — a form control inside a
// button is not clickable and not valid — so the card is a div and the commit
// is its own button at the bottom.
export function SumAction({
  a,
  money,
  unit,
  disabled,
  commit,
}: {
  a: Action;
  money: (n: number) => string;
  unit?: string;
  disabled: boolean;
  commit: (c: Command) => void;
}) {
  const sum = a.sum!;
  const format = (amount: number) => formatActionAmount(amount, money, unit);
  const [typed, setTyped] = useState<number | ''>(sum.preset);
  const amount = typed === '' ? 0 : typed;
  const over = amount > sum.most;
  const under = amount < sum.least;
  // What is wrong with the figure, said the same way whatever the field is for.
  // "More than the $18 there is" reads as a statement about the player's
  // pocket, which is true of a stake and false of a wage: the ceiling on what a
  // laundry pays a hand is what the trade will carry, not what is in the till.
  const bad = over ? `${format(sum.most)} is the most` : under ? `At least ${format(sum.least)}` : '';

  return (
    <div
      className={'action sum-action' + (a.disabled ? ' refused' : '')}
      title={[a.detail, a.disabled ? a.reason : ''].filter(Boolean).join(' — ')}
    >
      <strong>{a.label}</strong>
      <span className="meta">
        {a.minutes ? `${a.minutes} min` : ''}
        <span>
          {format(sum.least)}–{format(sum.most)}
        </span>
      </span>
      <span className="desc">{a.reason || a.detail}</span>
      <div className="sum-field">
        <label>
          <span>{sum.label}</span>
          <input
            type="number"
            inputMode="numeric"
            min={sum.least}
            max={sum.most}
            step={unit ? 1 : 5}
            value={typed}
            disabled={a.disabled || disabled}
            aria-label={`${sum.label} — ${format(sum.least)} to ${format(sum.most)}`}
            onChange={e =>
              setTyped(e.target.value === '' ? '' : Math.floor(Number(e.target.value)))
            }
          />
        </label>
        <button
          type="button"
          className="plain"
          disabled={a.disabled || disabled}
          onClick={() => setTyped(sum.most)}
        >
          All
        </button>
      </div>
      <button
        type="button"
        className="action-commit"
        disabled={a.disabled || disabled || !!bad}
        onClick={() => commit({kind: a.id, target: a.target, choice: a.choice, amount})}
      >
        {bad || `${a.label} — ${format(amount)}`}
      </button>
    </div>
  );
}
