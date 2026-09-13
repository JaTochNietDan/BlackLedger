import {useEffect, useRef, useState} from 'react';
import type {Place, VisualCue} from './types';
import {Portrait} from './Portrait';
import {playMoment} from './sound';

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

// The stage used to be a drawing: a flat black building with lit windows, a
// squad car, a body on the pavement, all rendered into an SVG that covered the
// city view. The city view already has twelve painted addresses in it. Taking
// the camera "there" and then hiding there behind a curtain was the wrong idea
// twice over, and this component no longer draws a building at all — CityStreet
// lights the real one, and what is left here is what a camera cannot say:
// the caption, who was in it, and the headline afterwards.

export function Theatre({
  cue,
  place,
  onDone,
  onProgress,
  plate = true,
  stagedGunfire = false,
}: {
  cue: VisualCue;
  place: Place;
  onDone: () => void;
  // Whether to show the painted plate for this kind of moment. In the city
  // view the building it happened at is on screen behind this band, so a stock
  // picture of a police station in front of the actual police station is one
  // picture too many.
  plate?: boolean;
  stagedGunfire?: boolean;
  // How far through the moment is, reported outward every frame so the street
  // can light the building while it happens.
  onProgress?: (t: number) => void;
}) {
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

  // The band reports where something happened, and the column it sits in can be
  // taller than the window — so on a short screen it opened below the fold and
  // the player was told nothing at all. It brings itself into view.
  const band = useRef<HTMLDivElement>(null);
  useEffect(() => {
    band.current?.scrollIntoView({block: 'nearest', behavior: 'smooth'});
  }, [cue.id]);

  // The noise the city makes, once, at the top of the moment — not on every
  // frame, and not again when the same moment is replayed mid-flight.
  useEffect(() => {
    if (!stagedGunfire) return playMoment(cue.kind);
  }, [cue.id, stagedGunfire]);

  useEffect(() => {
    if (matchMedia('(prefers-reduced-motion: reduce)').matches) {
      setPaper(true);
      onProgress?.(1);
      return;
    }
    const start = performance.now();
    let frame = 0;
    const step = (now: number) => {
      const at = Math.min(1, (now - start) / hold(cue));
      setT(at);
      onProgress?.(at);
      if (at < 1) frame = requestAnimationFrame(step);
      else setPaper(true);
    };
    frame = requestAnimationFrame(step);
    return () => cancelAnimationFrame(frame);
  }, [cue.id]);

  return (
    <div
      className={plate ? 'theatre' : 'theatre theatre-city'}
      ref={band}
      role="status"
      aria-label={cue.caption}
    >
      <div className="theatre-where">
        <span className="eyebrow">
          {cue.kind === 'arrest' ? 'YOU WERE TAKEN TO' : 'IT HAPPENED AT'}
        </span>
        <b>{place.name}</b>
      </div>
      {/* The painted plate for this kind of moment, when there is one, small and
        beside the caption rather than instead of the city. */}
      {plate && painted && (
        <div className="theatre-plate" style={{backgroundImage: `url(${scenePlate(cue.kind)})`}} />
      )}
      <p className="theatre-caption">{cue.caption}</p>
      {!!cue.actors?.length && (
        <div className="theatre-cast">
          {/* The core names who was in it. A scene about somebody that cannot show
          them is a scene about nobody. */}
          {cue.actors.map(a => (
            <span key={a.id} className="theatre-face">
              <Portrait id={a.id} size="small" />
              <small>{a.name}</small>
            </span>
          ))}
        </div>
      )}
      {paper && cue.headline && (
        <div className="theatre-paper">
          <small>THE BELLWETHER HERALD</small>
          <b>{cue.headline}</b>
        </div>
      )}
      <button className="plain theatre-done" onClick={onDone}>
        {paper ? 'Go on' : 'Skip'} →
      </button>
    </div>
  );
}
