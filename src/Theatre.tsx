import {useEffect,useState} from 'react';
import type {Place,VisualCue} from './types';

// The loudest moments of a campaign arrived as a paragraph in a list, and the
// paper reported them the next morning to a player who never saw them. This is
// where they are played: the camera goes to the building it happened in, the
// thing happens, and the Herald headline comes up afterwards — after the scene
// rather than instead of it.
//
// Nothing here decides anything. The core already committed the event; this
// only shows what it looked like.

// How long to hold. The city says what a moment is worth stopping for, so a
// killing is not played for the same two and a half seconds as a robbery.
const hold = (cue: VisualCue) => 1800 + (cue.gravity || 0) * 320;

export const scenePlate = (kind: string) => `/art/scenes/scene-${kind}-v1.jpg`;

// What still moves when the plate is painted: the flash, the bloom, the lamp.
// The building underneath is only drawn when there is no plate to stand on.
function stage(kind: string, place: Place, t: number, painted: boolean) {
  // t runs 0 → 1 across the scene.
  const lit = kind === 'raid' || kind === 'arrest';
  const flash = kind === 'gunfight' || kind === 'killing';
  const boom = kind === 'explosion';
  const shake = boom ? Math.max(0, 1 - t * 3) * 9 : 0;
  const jitter = shake ? `translate(${(Math.random() - .5) * shake},${(Math.random() - .5) * shake})` : '';

  // A building, flat and black, of roughly the right shape for what it is.
  const tall = place.type === 'casino' ? 150 : place.type === 'home' ? 190 : 165;
  const wide = place.type === 'casino' ? 290 : 230;
  const front = painted ? '' : `<g transform="${jitter}">` +
    `<rect x="${260 - wide / 2}" y="${250 - tall}" width="${wide}" height="${tall}" fill="#0a0d0d"/>` +
    Array.from({length: 3}, (_, r) => Array.from({length: 4}, (_, c) => {
      const on = ((r * 4 + c + place.id.length) % 3 !== 0);
      return `<rect x="${260 - wide / 2 + 22 + c * (wide - 60) / 3.4}" y="${250 - tall + 24 + r * (tall - 70) / 2.6}" width="21" height="26" fill="${on ? '#c9a86a' : '#171c1b'}" opacity="${on ? .55 + (lit ? .2 : 0) : 1}"/>`;
    }).join('')).join('') +
    `<rect x="238" y="196" width="44" height="54" fill="#161b1a"/>` +
    `</g>`;

  const floor = painted ? '' :
    `<rect y="250" width="520" height="110" fill="#101414"/>` +
    `<path d="M0 262h520" stroke="#1b2220" stroke-width="2"/>`;

  let event = '';
  // On a painted plate only light is added. The picture already has the body,
  // the onlookers and the cars in it; drawing another one on top reads as a
  // sticker, where a flash reads as something happening in the room.
  if (flash) {
    const firing = t < .45;
    const n = Math.floor(t * 14);
    if (firing && n % 2 === 0) {
      const x = 190 + (n % 3) * 18, r = 13 - (n % 3) * 2;
      event += `<circle cx="${x}" cy="238" r="${r}" fill="#ffe6a8" opacity=".9"/>`;
      if (painted) event += `<circle cx="${x}" cy="238" r="${r * 6}" fill="#ffe6a8" opacity=".13"/>`;
    }
    if (!painted && t > .4) event += `<ellipse cx="300" cy="290" rx="${Math.min(34, (t - .4) * 90)}" ry="9" fill="#0a0808" opacity=".85"/>`;
    if (!painted && t > .5) event += `<ellipse cx="300" cy="292" rx="${Math.min(46, (t - .5) * 110)}" ry="11" fill="#3a1512" opacity=".55"/>`;
  }
  if (boom) {
    const r = t < .3 ? t * 430 : 130 - (t - .3) * 60;
    event += `<circle cx="260" cy="215" r="${Math.max(0, r)}" fill="#ffcf7a" opacity="${Math.max(0, .85 - t)}"/>`;
    event += `<circle cx="260" cy="215" r="${Math.max(0, r * 1.5)}" fill="#c4531f" opacity="${Math.max(0, .4 - t * .5)}"/>`;
    if (!painted && t > .35) event += `<g fill="#0a0d0d" opacity=".9">${Array.from({length: 9}, (_, i) =>
      `<rect x="${260 + Math.cos(i * 2.2) * (t - .35) * 420}" y="${215 + Math.sin(i * 2.2) * (t - .35) * 240}" width="9" height="7"/>`).join('')}</g>`;
  }
  if (lit) {
    const swing = Math.sin(t * 22) * 40;
    if (!painted) {
      // A car at the kerb with a lamp turning over on the roof.
      event += `<g transform="translate(70,236)"><rect width="118" height="30" rx="6" fill="#0c1010"/>` +
        `<rect x="20" y="-16" width="72" height="20" rx="5" fill="#0c1010"/>` +
        `<circle cx="24" cy="32" r="9" fill="#151b1a"/><circle cx="94" cy="32" r="9" fill="#151b1a"/>` +
        `<circle cx="56" cy="-22" r="6" fill="#e0705c" opacity="${.55 + Math.sin(t * 22) * .4}"/></g>`;
    }
    // The lamp itself sweeps the street either way.
    event += `<path d="M126 214 L${300 + swing} 250 L${230 + swing} 250Z" fill="#e0705c" opacity="${painted ? .2 : .13}"/>` +
      (painted ? `<rect width="520" height="360" fill="#e0705c" opacity="${(.05 + Math.abs(Math.sin(t * 22)) * .07).toFixed(3)}"/>` : '');
  }

  return `<svg class="theatre-stage" viewBox="0 0 520 360" role="img">` +
    (painted ? '' : `<rect width="520" height="360" fill="#070a0a"/>`) + `${floor}${front}${event}` +
    `<rect width="520" height="360" fill="url(#vig)" pointer-events="none" opacity="${painted ? .55 : 1}"/>` +
    `<defs><radialGradient id="vig" cx=".5" cy=".45" r=".75">` +
    `<stop stop-color="#000" stop-opacity="0"/><stop offset="1" stop-color="#000" stop-opacity=".7"/></radialGradient></defs>` +
    `</svg>`;
}

export function Theatre({cue, place, onDone}: {cue: VisualCue; place: Place; onDone: () => void}) {
  const [t, setT] = useState(0);
  const [paper, setPaper] = useState(false);
  // The plate is a backdrop, not a dependency: a kind nobody has painted yet
  // still plays, on the drawn building.
  const [painted, setPainted] = useState(false);
  useEffect(() => {
    setPainted(false);
    const img = new Image();
    img.onload = () => setPainted(true);
    img.src = scenePlate(cue.kind);
  }, [cue.kind]);

  useEffect(() => {
    if (matchMedia('(prefers-reduced-motion: reduce)').matches) { setPaper(true); return }
    const start = performance.now();
    let frame = 0;
    const step = (now: number) => {
      const at = Math.min(1, (now - start) / hold(cue));
      setT(at);
      if (at < 1) frame = requestAnimationFrame(step); else setPaper(true);
    };
    frame = requestAnimationFrame(step);
    return () => cancelAnimationFrame(frame);
  }, [cue.id]);

  return <div className="theatre" role="dialog" aria-label={cue.caption}>
    <div className="theatre-where">
      <span className="eyebrow">{cue.kind === 'arrest' ? 'YOU WERE TAKEN TO' : 'IT HAPPENED AT'}</span>
      <b>{place.name}</b>
    </div>
    <div className={'theatre-frame' + (painted ? ' painted' : '')}
      style={painted ? {backgroundImage: `url(${scenePlate(cue.kind)})`} : undefined}
      dangerouslySetInnerHTML={{__html: stage(cue.kind, place, t, painted)}}/>
    <p className="theatre-caption">{cue.caption}</p>
    {!!cue.actors?.length && <p className="theatre-actors">{cue.actors.join(' · ')}</p>}
    {paper && cue.headline && <div className="theatre-paper">
      <small>THE BELLWETHER HERALD</small>
      <b>{cue.headline}</b>
    </div>}
    <button className="plain theatre-done" onClick={onDone}>{paper ? 'Go on' : 'Skip'} →</button>
  </div>;
}
