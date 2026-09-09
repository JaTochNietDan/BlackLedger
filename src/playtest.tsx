import {StrictMode, useEffect, useMemo, useState} from 'react';
import {createRoot} from 'react-dom/client';
import {CityIso} from './CityIso';
import {Theatre} from './Theatre';
import {playMoment, setSound, soundOn} from './sound';
import type {Snapshot, VisualCue} from './types';
import './style.css';

// The workshop.
//
// Every visual moment in this game is the end of a chain of decisions: to see
// what a killing looks like you have to get somebody killed, which means
// playing until somebody is worth killing. That is the right way round for a
// player and a useless way round for whoever is building the thing. A moment
// that is hard to reach is a moment that gets looked at once.
//
// So this is the same components, driven by hand. It runs on its own port,
// from its own binary, against a scratch world it makes for itself — it cannot
// reach a campaign, and there is nothing here to find in a shipped game because
// the game does not serve this page.

type Kind = {kind: string; gravity: number; hold: number};

function Workshop() {
  const [world, setWorld] = useState<Snapshot | null>(null);
  const [kinds, setKinds] = useState<Kind[]>([]);
  const [where, setWhere] = useState('');
  const [cue, setCue] = useState<VisualCue | null>(null);
  const [beat, setBeat] = useState(0);
  const [loud, setLoud] = useState(soundOn());
  // Holding a moment still.
  //
  // A moment is over in three or four seconds, which is right in the game and
  // useless for looking at one: by the time you have seen the explosion you
  // cannot say what the second half of it did. Scrubbing plays no sound and
  // runs no clock — it puts the drawing at whatever instant you ask for and
  // leaves it there.
  const [held, setHeld] = useState<number | null>(null);
  const [note, setNote] = useState('');

  useEffect(() => {
    fetch('/api/state').then(r => r.json()).then((w: Snapshot) => {
      setWorld(w);
      setWhere(w.locations[0]?.id || '');
    }).catch(() => setNote('The workshop server is not answering.'));
    fetch('/api/kinds').then(r => r.json()).then(setKinds).catch(() => {});
  }, []);

  const place = useMemo(
    () => world?.locations.find(l => l.id === where) || world?.locations[0],
    [world, where]);

  // Playing a moment is the only thing this page does: build the cue the core
  // would have built, and hand it to the same theatre the game uses.
  const play = (kind: string) => {
    if (!place) return;
    setBeat(0);
    setHeld(null);
    setCue({
      id: 'workshop-' + kind + '-' + Date.now(),
      kind, target: place.id,
      caption: captionFor(kind, place.name),
      headline: headlineFor(kind, place.name),
      actors: (world?.cast || []).slice(0, 2).map(c => ({id: c.id, name: c.name})),
      gravity: kinds.find(k => k.kind === kind)?.gravity || 0,
    });
  };

  if (!world || !place) {
    return <div className="workshop-empty">{note || 'Reading the scratch world…'}</div>;
  }

  return <div className="workshop">
    <aside className="workshop-panel">
      <header>
        <div className="eyebrow">BLACK LEDGER</div>
        <h1>Workshop</h1>
        <p>Not the game. A scratch world, on its own port, so a moment can be
          looked at without playing far enough to cause one.</p>
      </header>

      <section>
        <h2>Where</h2>
        <div className="workshop-places">
          {world.locations.map(l => <button key={l.id}
            aria-pressed={l.id === where}
            onClick={() => setWhere(l.id)}>{l.name}</button>)}
        </div>
      </section>

      <section>
        <h2>What happens</h2>
        <p className="workshop-hint">Ordered the way the city orders them: what it
          is worth stopping for.</p>
        <div className="workshop-kinds">
          {kinds.map(k => <button key={k.kind} className="workshop-kind"
            onClick={() => play(k.kind)}>
            <b>{k.kind}</b>
            <small>gravity {k.gravity} · {(k.hold / 1000).toFixed(1)}s</small>
          </button>)}
        </div>
      </section>

      <section>
        <h2>Hold it still</h2>
        <p className="workshop-hint">A moment is over in four seconds. This puts
          the drawing at one instant and leaves it there.</p>
        <div className="workshop-scrub">
          <input type="range" min={0} max={100} value={Math.round((held ?? beat) * 100)}
            aria-label="How far through the moment"
            onChange={e => setHeld(Number(e.target.value) / 100)}/>
          <b>{Math.round((held ?? beat) * 100)}%</b>
        </div>
        <button className="plain" disabled={held === null} onClick={() => setHeld(null)}>
          Let it run again
        </button>
      </section>

      <section>
        <h2>Sound</h2>
        <button className="plain" aria-pressed={loud}
          onClick={() => { setSound(!loud); setLoud(!loud) }}>
          {loud ? 'Sound on' : 'Sound off'}
        </button>
        <button className="plain" disabled={!cue} onClick={() => cue && playMoment(cue.kind)}>
          Play the sound again
        </button>
      </section>

      {cue && <section>
        <h2>Playing</h2>
        <p className="workshop-hint">{cue.kind} at {place.name} — {(beat * 100).toFixed(0)}% through</p>
      </section>}
    </aside>

    <main className="workshop-stage">
      {/* The spotlight's id is the ADDRESS the moment happens at, not the cue's
          own id — CityIso looks it up in state.locations. Passing the cue id
          drew nothing at all, silently, which looks exactly like a moment that
          has no graphic. */}
      <CityIso state={world} selected={where} onSelect={setWhere} onEnter={() => {}}
        spotlight={cue ? {id: cue.target, kind: cue.kind, t: held ?? beat} : null}/>
      {cue && place && held === null && <Theatre cue={cue} place={place} plate={false}
        onProgress={setBeat} onDone={() => { setCue(null); setBeat(0) }}/>}
    </main>
  </div>;
}

// The words a moment carries. The core writes these from what actually
// happened; here they are stand-ins that name the place, so the layout is
// exercised with real-length copy rather than with "test".
function captionFor(kind: string, place: string): string {
  switch (kind) {
    case 'killing': return `A man was put down against the wall at ${place}.`;
    case 'explosion': return `Something went up at ${place} and took the windows with it.`;
    case 'gunfight': return `Shots traded across the front of ${place}.`;
    case 'raid': return `Police came through the doors at ${place}.`;
    case 'arrest': return `They took them out of ${place} in front of everybody.`;
    case 'seizure': return `${place} changed hands, and not quietly.`;
    case 'attack': return `Somebody was worked over outside ${place}.`;
    default: return `The day's takings went out of the back of ${place}.`;
  }
}

function headlineFor(kind: string, place: string): string {
  return `${kind.toUpperCase()} AT ${place.toUpperCase()}`;
}

createRoot(document.getElementById('root')!).render(<StrictMode><Workshop/></StrictMode>);
