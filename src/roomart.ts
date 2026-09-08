import type {Place,Presence} from './types';

// There is no image model on this machine — Ollama carries only text — so the
// inside of a building is drawn rather than generated, in the same register as
// the newspaper cuts: flat noir shapes, deterministic from the place, no detail
// the eye has to resolve. It is a stage rather than a picture. What matters is
// that the people standing in the room are things on it that can be clicked.

const esc = (s: unknown) => String(s ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] || c));

function hash(seed: string) { let h = 2166136261; for (const c of seed) { h ^= c.charCodeAt(0); h = Math.imul(h, 16777619) } return Math.abs(h) }

// Where people stand in a room. Eighteen figures on one floor is a crowd nobody
// can read: they overlap, the back row is hidden behind the front, and the
// building itself disappears. So the room shows the handful who matter — the
// ones the player has business with, ordered as the core ordered them — and
// says plainly how many more are in there. The roster below is where the whole
// register lives.
const marks: [number, number][] = [
  [148, 338], [260, 344], [372, 336],
  [204, 312], [316, 308],
  [100, 306], [420, 302],
];

// StandingRoom is how many people are drawn before the room says "and N more".
export const StandingRoom = marks.length;

// Each kind of room is a back wall, a floor, and two or three pieces of
// furniture that say what the place is for.
function room(type: string, seed: number) {
  const warm = ['#2c2119', '#2a2420', '#241f1c'][seed % 3];
  const back = `<rect width="520" height="240" fill="${warm}"/>` +
    `<rect y="240" width="520" height="120" fill="#191614"/>` +
    `<path d="M0 240h520" stroke="#0d0b0a" stroke-width="3"/>`;
  const lamps = `<g fill="#e8c98a" opacity=".16">` +
    `<path d="M110 0v46l-38 42h76L110 46z"/><path d="M410 0v46l-38 42h76L410 46z"/></g>` +
    `<g fill="#f0d7a2"><circle cx="110" cy="90" r="7"/><circle cx="410" cy="90" r="7"/></g>`;

  switch (type) {
    case 'bar':
      return back + lamps +
        `<g fill="#3a2c22"><rect x="40" y="118" width="440" height="14"/><rect x="40" y="132" width="440" height="8" fill="#241b15"/></g>` +
        `<g fill="#1f1913"><rect x="60" y="60" width="400" height="58"/></g>` +
        `<g fill="#8a6c44" opacity=".8">${Array.from({length: 18}, (_, i) => `<rect x="${72 + i * 22}" y="${72 + (i % 3) * 5}" width="7" height="${18 - (i % 3) * 4}"/>`).join('')}</g>` +
        `<g fill="#2b211a"><rect x="46" y="246" width="120" height="10"/><rect x="52" y="256" width="8" height="48"/><rect x="152" y="256" width="8" height="48"/></g>`;
    case 'casino':
      return back + lamps +
        `<g fill="#243a2c"><ellipse cx="260" cy="286" rx="150" ry="46"/></g>` +
        `<g fill="#1a2c20"><ellipse cx="260" cy="280" rx="150" ry="46"/></g>` +
        `<g fill="#c9a86a" opacity=".7"><circle cx="210" cy="272" r="6"/><circle cx="238" cy="278" r="6"/><circle cx="292" cy="270" r="6"/></g>` +
        `<g fill="#1f1913"><rect x="70" y="70" width="150" height="48"/><rect x="300" y="70" width="150" height="48"/></g>` +
        `<g fill="#8a6c44" opacity=".55"><rect x="86" y="84" width="118" height="4"/><rect x="316" y="84" width="118" height="4"/></g>`;
    case 'racket':
      return back +
        `<g fill="#242c2b"><rect x="30" y="96" width="130" height="120"/><rect x="180" y="112" width="120" height="104"/><rect x="330" y="88" width="150" height="128"/></g>` +
        `<g fill="#0f1413"><rect x="44" y="110" width="102" height="30"/><rect x="196" y="126" width="88" height="26"/><rect x="346" y="102" width="118" height="32"/></g>` +
        `<g fill="#3b3a30" opacity=".8">${Array.from({length: 7}, (_, i) => `<rect x="${58 + i * 62}" y="248" width="44" height="26"/>`).join('')}</g>`;
    case 'market':
      return back +
        `<g fill="#2b2a22"><rect x="20" y="150" width="480" height="10"/></g>` +
        `<g fill="#332f26">${Array.from({length: 5}, (_, i) => `<path d="M${30 + i * 100} 160h84l-6 40h-72z"/>`).join('')}</g>` +
        `<g fill="#4a4335" opacity=".9">${Array.from({length: 5}, (_, i) => `<path d="M${26 + i * 100} 138h92l-8 14h-76z"/>`).join('')}</g>` +
        `<g fill="#7d6a45" opacity=".5">${Array.from({length: 14}, (_, i) => `<circle cx="${44 + i * 33}" cy="${188 + (i % 3) * 4}" r="5"/>`).join('')}</g>`;
    case 'home':
      return back + lamps +
        `<g fill="#241d18"><rect x="46" y="96" width="150" height="120"/><rect x="330" y="120" width="140" height="96"/></g>` +
        `<g fill="#0e0c0b"><rect x="60" y="110" width="122" height="92"/></g>` +
        `<g fill="#3a2f26"><rect x="200" y="238" width="130" height="18"/><rect x="206" y="256" width="10" height="44"/><rect x="314" y="256" width="10" height="44"/></g>`;
    case 'civic':
      return back +
        `<g fill="#22201d">${Array.from({length: 4}, (_, i) => `<rect x="${52 + i * 116}" y="52" width="34" height="188"/>`).join('')}</g>` +
        `<g fill="#15130f"><rect x="150" y="128" width="220" height="88"/></g>` +
        `<g fill="#2c2823"><rect x="170" y="248" width="180" height="12"/><rect x="176" y="260" width="9" height="42"/><rect x="341" y="260" width="9" height="42"/></g>`;
    default: // work: a dockside shed
      return back +
        `<g fill="#232a2b"><rect x="24" y="80" width="200" height="136"/><rect x="270" y="110" width="210" height="106"/></g>` +
        `<g fill="#101516"><rect x="44" y="100" width="160" height="40"/><rect x="292" y="128" width="166" height="34"/></g>` +
        `<g fill="#37342b" opacity=".85">${Array.from({length: 6}, (_, i) => `<rect x="${40 + i * 74}" y="242" width="52" height="34"/>`).join('')}</g>`;
  }
}

// A figure is a silhouette with a coat and a hat, lit from the front. The only
// thing that varies is build and stance, from the person's own id, so the same
// person is always the same shape in the same room.
function figure(who: Presence, x: number, y: number, selected: boolean, painted = false) {
  const seed = hash(who.id);
  const tall = 60 + (seed % 5) * 4;
  const wide = 17 + (seed % 3) * 2;
  const hat = seed % 4 !== 0;
  const lean = ((seed >> 3) % 3) - 1;
  // In a painted room people read as shapes against the light, not as pale
  // cut-outs laid on top of it. The colour that says who they are moves to the
  // rim, where a single hard light would actually catch them.
  const says = who.yours ? '#d6b77c' : who.sore || who.overdue ? '#c08476' : '#cfc9b6';
  const tone = painted ? '#0b0d0c' : says;
  const rim = painted ? says : 'none';
  return `<g class="figure${selected ? ' picked' : ''}" data-person="${esc(who.id)}" role="button" tabindex="0" aria-label="${esc(who.name)} — ${esc(who.standing)}" transform="translate(${x},${y})">` +
    `<ellipse cx="0" cy="4" rx="${wide + 4}" ry="6" fill="#0b0a09" opacity=".55"/>` +
    `<path d="M${-wide} 2q0-${tall * 0.62} ${wide + lean * 2} -${tall * 0.62}q${wide} 0 ${wide} ${tall * 0.62}z" fill="#12100e"/>` +
    `<path d="M${-wide + 3} 0q0-${tall * 0.58} ${wide + lean * 2 - 3} -${tall * 0.58}q${wide - 3} 0 ${wide - 3} ${tall * 0.58}z" fill="${tone}" opacity="${painted ? '.95' : '.92'}" stroke="${rim}" stroke-width="${painted ? 1.1 : 0}" stroke-opacity=".5"/>` +
    `<circle cx="${lean}" cy="${-tall * 0.62 - 9}" r="9" fill="${tone}" opacity="${painted ? '.95' : '.92'}" stroke="${rim}" stroke-width="${painted ? 1.1 : 0}" stroke-opacity=".5"/>` +
    (hat ? `<path d="M${lean - 15} ${-tall * 0.62 - 12}h30l-4-9h-22z" fill="#12100e"/><path d="M${lean - 17} ${-tall * 0.62 - 11}h34v3h-34z" fill="#12100e"/>` : '') +
    `<circle class="halo" cx="${lean}" cy="${-tall * 0.31}" r="${tall * 0.7}" fill="none"/>` +
    `<title>${esc(who.name)} — ${esc(who.standing)}</title>` +
    `</g>`;
}

// Every address has a painted interior now, generated offline (tools/interiors.py)
// and shipped as a JPEG. The drawn room below stays as the fallback: a building
// added tomorrow has somewhere to stand before anybody renders it.
export const paintedRoom = (id: string) => `/art/rooms/room-${id}-v1.jpg`;

export function interiorSVG(place: Place, people: Presence[], selected: string, painted = false) {
  const seed = hash(place.id);
  const shown = people.slice(0, marks.length);
  // Depth: whoever is further back is smaller and dimmer, so a room of seven
  // reads as a room rather than as a row of stickers.
  return `<svg class="interior" viewBox="0 0 520 360" role="group" aria-label="Inside ${esc(place.name)}">` +
    `<defs><radialGradient id="int-${esc(place.id)}" cx=".5" cy=".28" r=".8">` +
    `<stop stop-color="#6b5a3f" stop-opacity=".30"/><stop offset="1" stop-color="#000" stop-opacity=".55"/>` +
    `</radialGradient></defs>` +
    (painted ? '' : room(place.type, seed)) +
    shown.map((who, i) => {
    const depth = .72 + (marks[i][1] - 300) / 200;
      return `<g transform="translate(${marks[i][0]},${marks[i][1]}) scale(${depth.toFixed(3)}) translate(${-marks[i][0]},${-marks[i][1]})" opacity="${(.62 + depth * .38).toFixed(2)}">` +
        figure(who, marks[i][0], marks[i][1], who.id === selected, painted) + `</g>`;
    }).join('') +
    (painted ? '' : `<rect width="520" height="360" fill="url(#int-${esc(place.id)})" pointer-events="none"/>`) +
    (people.length > shown.length
      ? `<g><rect x="330" y="330" width="180" height="22" fill="#0b0f0ecc"/>` +
        `<text x="500" y="345" text-anchor="end" fill="#c9b98f" font-size="11">` +
        `and ${people.length - shown.length} more in here</text></g>`
      : '') +
    `</svg>`;
}
