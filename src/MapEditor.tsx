import {useEffect, useState} from 'react';
import cutouts from '../public/art/iso/isometric.json';
import {NOTHING, saveLayout} from './layout';
import type {Layout} from './layout';

// The panel beside the map while it is being arranged.
//
// Deliberately small. The map itself is the interface — a slot is chosen by
// clicking its ground — and this is only the list of things that can stand
// there, plus the nudge and the save. Anything more here is a second way to do
// what the map already does better.

const sprites = (cutouts as {id: string; file: string}[]).map(c => c.id);

// How far one press of a nudge moves a building, in screen pixels. Small: this
// is for settling a picture onto its ground, not for moving it somewhere else.
const NUDGE = 2;

export function MapEditor({
  layout,
  slot,
  onChange,
  onClose,
}: {
  layout: Layout;
  slot: string;
  // Takes an update rather than a value. Every edit here is a change to the
  // arrangement as it stands at the moment of the click, and a handler that
  // closes over the arrangement as it stood at the last render will quietly
  // undo whatever happened in between.
  onChange: (update: (prev: Layout) => Layout) => void;
  onClose: () => void;
}) {
  const [saving, setSaving] = useState('');
  const placed = slot ? layout.slots[slot] : undefined;

  const put = (patch: Partial<{sprite: string; dx: number; dy: number}>) => {
    if (!slot) return;
    onChange(prev => ({
      ...prev,
      slots: {...prev.slots, [slot]: {...(prev.slots[slot] || {sprite: NOTHING}), ...patch}},
    }));
  };
  const forget = () => {
    if (!slot) return;
    onChange(prev => {
      const slots = {...prev.slots};
      delete slots[slot]; // back to whatever the city would choose
      return {...prev, slots};
    });
  };

  // The nudge is on the arrow keys, because it is the one thing here that
  // wants to be done a dozen times in a row while looking at the map.
  useEffect(() => {
    if (!slot) return;
    const move = (e: KeyboardEvent) => {
      const by: Record<string, [number, number]> = {
        ArrowLeft: [-NUDGE, 0],
        ArrowRight: [NUDGE, 0],
        ArrowUp: [0, -NUDGE],
        ArrowDown: [0, NUDGE],
      };
      const step = by[e.key];
      if (!step) return;
      e.preventDefault();
      put({dx: (placed?.dx || 0) + step[0], dy: (placed?.dy || 0) + step[1]});
    };
    window.addEventListener('keydown', move);
    return () => window.removeEventListener('keydown', move);
  });

  const save = async () => {
    setSaving('saving…');
    const trouble = await saveLayout(layout);
    setSaving(trouble || 'saved to art/city-layout.json');
    setTimeout(() => setSaving(''), 4000);
  };

  return (
    <aside className="map-editor">
      <header>
        <h4>Arranging the map</h4>
        <button className="plain" onClick={onClose}>
          Done
        </button>
      </header>

      {!slot && (
        <p className="nothing-here">
          Click the ground of any slot on the map. The blue diamonds are the plots buildings stand
          on — the ground rather than the picture, because a building overlaps its neighbours'
          ground and the sprite under the pointer is often not the one you meant.
        </p>
      )}

      {slot && (
        <>
          <p className="editor-slot">
            Block <b>{slot.split(',').slice(0, 2).join(', ')}</b>, slot <b>{slot.split(',')[2]}</b>
          </p>

          <div className="editor-sprites">
            <button
              className={'sprite' + (placed?.sprite === NOTHING ? ' picked' : '')}
              onClick={() => put({sprite: NOTHING})}
            >
              Nothing here
            </button>
            {sprites.map(id => (
              <button
                key={id}
                className={'sprite' + (placed?.sprite === id ? ' picked' : '')}
                onClick={() => put({sprite: id})}
              >
                {id}
              </button>
            ))}
          </div>

          <p className="editor-nudge">
            Nudge with the arrow keys
            {placed && (placed.dx || placed.dy) ? `: ${placed.dx || 0}, ${placed.dy || 0}` : ''}
          </p>
          <button className="plain" onClick={forget}>
            Let the city choose again
          </button>
        </>
      )}

      <footer>
        <button onClick={save}>Save the arrangement</button>
        {saving && <small>{saving}</small>}
        {layout.editable === false && (
          <small>Read only. Start the game with BLACK_LEDGER_EDIT=1 to save.</small>
        )}
      </footer>
    </aside>
  );
}
