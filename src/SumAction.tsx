import {useState} from 'react';
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
  disabled,
  commit,
}: {
  a: Action;
  money: (n: number) => string;
  disabled: boolean;
  commit: (c: Command) => void;
}) {
  const sum = a.sum!;
  const [typed, setTyped] = useState<number | ''>(sum.preset);
  const amount = typed === '' ? 0 : typed;
  const over = amount > sum.most;
  const under = amount < sum.least;
  const bad = over
    ? `More than the ${money(sum.most)} there is`
    : under
      ? `At least ${money(sum.least)}`
      : '';

  return (
    <div
      className={'action sum-action' + (a.disabled ? ' refused' : '')}
      title={[a.detail, a.disabled ? a.reason : ''].filter(Boolean).join(' — ')}
    >
      <strong>{a.label}</strong>
      <span className="meta">
        {a.minutes ? `${a.minutes} min` : ''}
        <span>
          {money(sum.least)}–{money(sum.most)}
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
            step={5}
            value={typed}
            disabled={a.disabled || disabled}
            aria-label={`${sum.label} — ${money(sum.least)} to ${money(sum.most)}`}
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
        onClick={() => commit({kind: a.id, target: a.target, amount})}
      >
        {bad || `${a.label} — ${money(amount)}`}
      </button>
    </div>
  );
}
