import {portrait} from './art';

// Six of this city's people were painted; everybody else is drawn from their
// own id, deterministically, so a face is always the same face. Lifted out of
// main so the room can use it too — a list of names is a list, a list of faces
// is a room.
const painted: {[id: string]: [number, number]} = {
  mara: [0, 0], leo: [1, 0], vittorio: [2, 0],
  elena: [0, 1], harlow: [1, 1], 'Alex Varga': [2, 1],
};

export function Portrait({id, size}: {id: string; size?: 'small' | 'tiny'}) {
  const cell = painted[id];
  const cls = 'portrait' + (size ? ' portrait-' + size : '');
  return cell
    ? <span aria-hidden="true" className={cls + ' painted-portrait'} style={{backgroundPosition: `${cell[0] * 50}% ${cell[1] * 100}%`}}/>
    : <span aria-hidden="true" className={size ? 'portrait-wrap portrait-' + size : undefined} dangerouslySetInnerHTML={{__html: portrait(id)}}/>;
}
