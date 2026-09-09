import type {Place} from './types';

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
  [128, 352], [260, 356], [392, 350],
  [190, 318], [330, 314],
  [70, 310], [452, 306],
];

// The same positions as fractions of the frame, for the figures the interface
// draws in HTML over the painted room. A silhouette with nothing on its head
// could be anybody; the person you click and the face in the roster have to be
// visibly the same person.
export const standingSpots = marks.map(([x, y]) => ({
  left: (x / 520) * 100,
  bottom: ((360 - y) / 360) * 100,
  // Further back is smaller, on the same scale the drawn room used.
  scale: 0.74 + (y - 306) / 170,
}));

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

// Every address has a painted interior, generated offline (tools/interiors.py)
// and shipped as a JPEG. The drawn room below stays as the fallback: a building
// added tomorrow has somewhere to stand before anybody renders it.
export const paintedRoom = (id: string) => `/art/rooms/room-${id}-v1.jpg`;

export function interiorSVG(place: Place, painted = false) {
  const seed = hash(place.id);
  if (painted) return '';
  return `<svg class="interior" viewBox="0 0 520 360" role="presentation">` +
    `<defs><radialGradient id="int-${esc(place.id)}" cx=".5" cy=".28" r=".8">` +
    `<stop stop-color="#6b5a3f" stop-opacity=".30"/><stop offset="1" stop-color="#000" stop-opacity=".55"/>` +
    `</radialGradient></defs>` +
    room(place.type, seed) +
    `<rect width="520" height="360" fill="url(#int-${esc(place.id)})" pointer-events="none"/>` +
    `</svg>`;
}

// ---------------------------------------------------------------------------
// The light inside.
//
// The city outside knows what hour it is: the ground, the lamps, the window
// spill and the haze all come off one number from the core's own clock. The
// room did not, so stepping inside at three in the morning put the player in
// the same evenly lit room they would have found at noon, and the inside and
// the outside stopped being the same place.
//
// This is the same two functions the city uses, read the same way. It returns
// what to lay over the room rather than drawing anything, so the backdrop —
// painted or drawn — is untouched underneath.

// The extension is explicit because this module is compiled and run directly
// by the frontend tests under plain node, which will not guess it. Vite
// resolves it the same either way.
import {goldenness, nightness} from './iso.js';

export type RoomLight = {
  dark: number;     // how much of the room the night has taken
  gold: number;     // how warm what light there is comes in
  wash: string;     // a CSS gradient to lay over the backdrop
};

export function roomLight(minute: number): RoomLight {
  const dark = nightness(minute), gold = goldenness(minute);
  // At night the room is not uniformly dark: there is a lamp somewhere and the
  // corners go first. At dawn and dusk the light is low and comes in warm from
  // one side, which is a window rather than a lamp.
  const pool = `radial-gradient(120% 90% at 46% 34%,
     rgba(${Math.round(228 - gold * 12)},${Math.round(196 + gold * 8)},${Math.round(150 - gold * 40)},${(0.05 + dark * 0.13).toFixed(3)}) 0%,
     rgba(0,0,0,0) 62%)`;
  const window = gold > 0
    ? `linear-gradient(102deg, rgba(196,124,58,${(gold * 0.17).toFixed(3)}) 0%, rgba(196,124,58,0) 46%)`
    : '';
  const night = `linear-gradient(rgba(8,13,18,${(dark * 0.46).toFixed(3)}), rgba(6,10,14,${(dark * 0.58).toFixed(3)}))`;
  return {dark, gold, wash: [pool, window, night].filter(Boolean).join(', ')};
}
