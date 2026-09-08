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
  mara: [0, 0], leo: [1, 0], vittorio: [2, 0],
  elena: [0, 1], harlow: [1, 1], 'Alex Varga': [2, 1],
};

// The generated sheet: 6 across, 4 down.
const CAST_COLS = 6, CAST_ROWS = 4, CAST = CAST_COLS * CAST_ROWS;

function faceFor(id: string) {
  let h = 2166136261;
  for (const c of id) { h ^= c.charCodeAt(0); h = Math.imul(h, 16777619) }
  return Math.abs(h) % CAST;
}

export function Portrait({id, size}: {id: string; size?: 'small' | 'tiny'}) {
  const cell = painted[id];
  const cls = 'portrait' + (size ? ' portrait-' + size : '');
  if (cell) {
    return <span aria-hidden="true" className={cls + ' painted-portrait'}
      style={{backgroundPosition: `${cell[0] * 50}% ${cell[1] * 100}%`}}/>;
  }
  const n = faceFor(id);
  return <span aria-hidden="true" className={cls + ' cast-portrait'} style={{
    backgroundPosition: `${(n % CAST_COLS) * 100 / (CAST_COLS - 1)}% ${Math.floor(n / CAST_COLS) * 100 / (CAST_ROWS - 1)}%`,
  }}>
    <span className="drawn-fallback" dangerouslySetInnerHTML={{__html: portrait(id)}}/>
  </span>;
}
