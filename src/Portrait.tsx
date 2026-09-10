import {portrait} from './art';

// Six of this city's people were painted by hand. The rest were a drawing made
// out of their id — better than one repeated face, and still obviously a
// drawing beside the painted six.
//
// There is a generated cast now: two dozen period portraits made offline by a
// local image model (tools/portraits.py) and baked into a sheet, so nothing at
// runtime depends on a model being installed. The people of Bellwether are not
// gendered by the rules — the city says "they" about everybody — so a face is
// assigned by a hash of the person's id. The same person is always the same
// face, and any face can belong to any name. The drawn version remains the
// fallback for anybody the sheet cannot cover.

const painted: {[id: string]: [number, number]} = {
  mara: [0, 0],
  leo: [1, 0],
  vittorio: [2, 0],
  elena: [0, 1],
  harlow: [1, 1],
  'Alex Varga': [2, 1],
};

// The generated sheet: 6 across, 4 down.
const CAST_COLS = 6,
  CAST_ROWS = 4,
  CAST = CAST_COLS * CAST_ROWS;

// The core says which face somebody wears, because the voice they speak in is
// chosen from the same answer and only one side of the wall can be the author
// of it. This is the fallback for anywhere the core has not said — the same
// arithmetic it uses, kept so a portrait never comes out blank.
function faceFor(id: string) {
  let h = 2166136261;
  for (const c of id) {
    h ^= c.charCodeAt(0);
    h = Math.imul(h, 16777619);
  }
  return Math.abs(h) % CAST;
}

// The size of the cast, for anywhere that offers a choice of one. The core
// holds the same number (CastFaces) because the core is what refuses a bad one.
export const CAST_FACES = CAST;

export function Portrait({id, size, face}: {id: string; size?: 'small' | 'tiny'; face?: number}) {
  // A face dealt out by a hash of a name is fair to the cast and can still hand
  // somebody a portrait they do not recognise as themselves. A player who has
  // picked their own overrules that, and overrules a painted one too.
  const cell = face ? undefined : painted[id];
  const cls = 'portrait' + (size ? ' portrait-' + size : '');
  if (cell) {
    return (
      <span
        aria-hidden="true"
        className={cls + ' painted-portrait'}
        style={{backgroundPosition: `${cell[0] * 50}% ${cell[1] * 100}%`}}
      />
    );
  }
  const n = face ? (face - 1) % CAST : faceFor(id);
  // The drawn version is underneath and the generated face is laid over it.
  //
  // It used to be the other way round — the sheet on the element's own
  // background and the drawing as a child behind it — and that does not work:
  // a negative z-index child still paints above its parent's background, so
  // the crude drawing covered the generated face every time and the fallback
  // was what everybody actually saw. Layering it this way means the drawing is
  // seen only when the sheet genuinely fails to load, which is what a fallback
  // is for.
  return (
    <span aria-hidden="true" className={cls + ' cast-portrait'}>
      <span className="drawn-fallback" dangerouslySetInnerHTML={{__html: portrait(id)}} />
      <span
        className="cast-face"
        style={{
          backgroundPosition: `${((n % CAST_COLS) * 100) / (CAST_COLS - 1)}% ${(Math.floor(n / CAST_COLS) * 100) / (CAST_ROWS - 1)}%`,
        }}
      />
    </span>
  );
}

// Where somebody's face sits on the sheet, as a CSS background-position. The
// newspaper needs the same face the rest of the game shows them by, and needs
// it as a value rather than as an element.
export function pressFace(id?: string): string | null {
  if (!id) return null;
  const cell = painted[id];
  if (cell) return `${cell[0] * 50}% ${cell[1] * 100}%`;
  const n = faceFor(id);
  return `${((n % CAST_COLS) * 100) / (CAST_COLS - 1)}% ${(Math.floor(n / CAST_COLS) * 100) / (CAST_ROWS - 1)}%`;
}
