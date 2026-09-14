import {KeyboardShortcuts} from './KeyboardShortcuts';
import {CityAccounts} from './CityAccounts';
import {SceneNewspaper} from './SceneNewspaper';
import {MapMenu} from './MapMenu';
import './mapFirst.css';
import './encounterStyle.css';
import './cityHud.css';
import {cityOwnsAudio} from './city3dEvents';
import type {VisualCue} from './types';
import {paintedCar} from './cityAssets';
import {SumAction} from './SumAction';
import {useCallback, useEffect, useRef, useState} from 'react';
import {createRoot} from 'react-dom/client';
import {createPortal} from 'react-dom';
import {City3D} from './City3D';
import {VoicePlayer, speaking, speakerOf} from './voice';
import {unreadInLatest,articleForScene} from './paper';
import {paintedAsset, paintedFront, paintedMask} from './cityAssets';
import {type Journey} from './TravelPresentation';
import {icon, pressPlate} from './art';
import {ActionList} from './ActionList';
import {Interior} from './Interior';
import {MarketScreen} from './MarketScreen';
import {Portrait, CAST_FACES} from './Portrait';
import {Casino, isTableAction} from './Casino';
import {BackRoomScene} from './BackRoomScene';

// The verbs of a hand in the back room. Drawn on the table itself, so the room's
// ordinary list must not offer "Throw the hand in" between hiring and restocking.
function isBackRoomAction(id: string) {
  return id === 'change' || id === 'bet' || id === 'call' || id === 'fold';
}
import {Outcome} from './Outcome';
import {PeopleScreen} from './PeopleScreen';
import {LedgerScreen} from './LedgerScreen';
import {FamiliesScreen} from './FamiliesScreen';
import {Theatre} from './Theatre';
import {playMoment, setSound, soundOn} from './sound';
import {Herald} from './Herald';
import type {Snapshot, Command, Action, Place} from './types';
import './style.css';
const money = (n: number) => '$' + Math.floor(n).toLocaleString();
const time = (m: number) =>
  `Day ${Math.floor(m / 1440) + 1} · ${String(Math.floor((m % 1440) / 60)).padStart(2, '0')}:${String(m % 60).padStart(2, '0')}`;
function Icon({id}: {id: string}) {
  return <span aria-hidden="true" dangerouslySetInnerHTML={{__html: icon(id)}} />;
}
class RequestError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
  }
}
async function api<T>(path: string, payload?: unknown): Promise<T> {
  const r = await fetch(
    '/api/' + path,
    payload !== undefined
      ? {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload),
        }
      : {},
  );
  const data = await r.json();
  if (!r.ok) throw new RequestError(data.error || 'Request failed', r.status);
  return data;
}
// Which city the player last chose to look at. Boot used to force the card
// view on every load, so the isometric city could be picked and then quietly
// taken away again by the next refresh.
function remembered(): 'iso' {
  return 'iso';
}
function App() {
  const [world, setWorld] = useState<Snapshot | null>(null),
    [tab, setTab] = useState('city'),
    [selected, setSelected] = useState('bar'),
    [busy, setBusy] = useState(false),
    [notice, setNotice] = useState(''),
    [voice, setVoice] = useState(localStorage.getItem('black-ledger-voice') === 'yes'),
    [speech, setSpeech] = useState('Read aloud'),
    [error, setError] = useState('');
  useEffect(()=>{if(tab==='news')return playMoment('newspaper');},[tab]);
  const [cityView, setCityView] = useState<'interior' | 'iso'>(remembered);
  // Whether the player is sitting at a table. A game takes the whole screen and
  // holds it until they get up: playing one out of the corner of a sidebar, with
  // the building's staff and supplies beside it, is being shown a game rather
  // than playing one.
  // Whether the player is at the tables is the world's fact, not this file's.
  // It used to be local state, which meant the takeover could open on a table
  // still showing the last hand, the last spin and where the drums stopped —
  // all of it saved state — and that reads as a game that started without you.
  // Sitting down and getting up are commands now, and the core clears the felt.
  const atTable = !!world && world.seated === world.player.location;
  // Which kind of seat it is. The back room is a game with no house in it and
  // no other table beside it, so it gets its own screen rather than a tab in
  // the casino's.
  // Which of the two things in the room they sat down to. It has to be what
  // the player chose rather than what the room holds: a bar and a poolhall each
  // have a wall of machines and a room behind the room, so a screen picked from
  // the address put somebody who asked for the machines into a hand of cards.
  const inTheBackRoom = atTable && world!.seated_to === 'back';
  const [sound, setSoundOn] = useState(soundOn);
  const [motion, setMotion] = useState(() => {
    try {
      return localStorage.getItem('black-ledger-motion') !== 'off';
    } catch {
      return true;
    }
  });
  const [newsSeen, setNewsSeen] = useState<string>(() => {
    try {
      return localStorage.getItem('black-ledger-news-seen') || '';
    } catch {
      return '';
    }
  });
  const [journey, setJourney] = useState<Journey | null>(null);
  const [journeyProgress, setJourneyProgress] = useState(0);
  // The moment the city thought was worth taking the player to.
  const [playing, setPlaying] = useState<VisualCue | null>(null);
  // How far through the moment the camera is holding on. Driven by the theatre,
  // read by the street, which lights the building while it happens.
  const [beat, setBeat] = useState(0);
  const [finishedCue, setFinishedCue] = useState('');
  const scenePending = !!playing && finishedCue !== playing.id;
  const sceneArticle=playing?articleForScene(playing,world?.newspaper||[]):undefined;
  const newspaperVisible=!!sceneArticle&&!!playing&&finishedCue===playing.id;
  // Let the recorded cause of death finish, then its newspaper, before the
  // memorial takes focus. A loaded death without playback opens immediately.
  const deathVisible=!!world&&!world.player.alive&&!scenePending&&!newspaperVisible;
  const scene = useRef<HTMLElement | null>(null),
    latest = useRef(world),
    busyRef = useRef(false);
  const sceneReplay=useRef(0);
  const sceneConditions=useRef<{world:string;revision:number;conditions:Record<string,number>}|null>(null);
  latest.current = world;
  const voicePlayer = useRef<VoicePlayer | null>(null);
  if (!voicePlayer.current)
    voicePlayer.current = new VoicePlayer({
      load: async (event, signal) => {
        const r = await fetch('/api/speech', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({event}),
          signal,
        });
        if (!r.ok) throw Error('Voice unavailable');
        return r.blob();
      },
      audio: blob => {
        const url = URL.createObjectURL(blob),
          player = new Audio(url);
        let disposed = false;
        return {
          play: () => player.play(),
          pause: () => player.pause(),
          dispose: () => {
            if (!disposed) {
              disposed = true;
              URL.revokeObjectURL(url);
              player.onended = null;
            }
          },
          onEnded: fn => {
            player.onended = fn;
          },
        };
      },
      status: setSpeech,
      unavailable: () => setNotice('Voice is unavailable. You can continue reading.'),
    });
  const stopVoice = useCallback(() => voicePlayer.current!.stop(), []);
  // Every hook must run on every render, including the one where the world is
  // still null. Sitting below the early return below made the hook count change
  // between the first render and the second, which is React error #310 and a
  // blank page on every boot.
  useEffect(() => {
    if (tab !== 'news') return;
    const latest = world?.newspaper?.[0]?.id;
    if (latest && latest !== newsSeen) {
      try {
        localStorage.setItem('black-ledger-news-seen', latest);
      } catch {}
      setNewsSeen(latest);
    }
  }, [tab, world?.newspaper?.[0]?.id, newsSeen]);
  const speak = useCallback(async () => {
    const id = latest.current?.event?.id;
    if (id) await voicePlayer.current!.speak(id);
  }, []);
  useEffect(() => {
    let canceled = false;
    async function boot() {
      try {
        const pending = localStorage.getItem('black-ledger-pending');
        if (pending) {
          try {
            await api('action', JSON.parse(pending));
            localStorage.removeItem('black-ledger-pending');
          } catch (err) {
            if (err instanceof RequestError && [400, 409].includes(err.status))
              localStorage.removeItem('black-ledger-pending');
            else throw err;
          }
        }
        const w = await api<Snapshot>('state');
        if (!canceled) {
          setWorld(w);
          const destination =
            w.player.location === 'room' && w.player.job_count === 0 ? 'bar' : w.player.location;
          setSelected(destination);
          setCityView(remembered());
        }
      } catch (err) {
        setError((err as Error).message);
      }
    }
    boot();
    return () => {
      canceled = true;
      stopVoice();
    };
  }, [stopVoice]);
  useEffect(() => {
    if (tab !== 'city' || cityView !== 'iso') setJourney(null);
  }, [tab, cityView]);
  useEffect(() => {
    if (!notice) return;
    const id = setTimeout(() => setNotice(''), 6000);
    return () => clearTimeout(id);
  }, [notice]);
  useEffect(() => {
    if (deathVisible) scene.current?.focus();
  }, [deathVisible]);
  useEffect(() => {
    stopVoice();
    if (world?.event) {
      scene.current?.focus();
      if (voice) speak();
    }
  }, [world?.event?.id, stopVoice, speak]);
  useEffect(() => {
    if (!voice || world?.event || world?.director.status !== 'ready') return;
    const abort = new AbortController();
    fetch('/api/speech/prepare', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: '{}',
      signal: abort.signal,
    }).catch(() => {});
    return () => abort.abort();
  }, [voice, world?.event?.id, world?.director.status]);
  useEffect(() => {
    const id = setInterval(async () => {
      if (busyRef.current || latest.current?.director.status !== 'writing') return;
      try {
        const next = await api<Snapshot>('state');
        setWorld(w => (w && w.revision === next.revision ? {...w, director: next.director} : w));
      } catch {}
    }, 3000);
    return () => clearInterval(id);
  }, []);
  async function prepare() {
    try {
      await api('director', {});
      setWorld(await api<Snapshot>('state'));
    } catch (err) {
      setNotice((err as Error).message);
    }
  }
  async function commit(command: Command) {
    if (!world || busyRef.current || journey) return;
    stopVoice();
    busyRef.current = true;
    setBusy(true);
    try {
      const pending = localStorage.getItem('black-ledger-pending');
      if (pending) {
        await api('action', JSON.parse(pending));
        localStorage.removeItem('black-ledger-pending');
        setWorld(await api<Snapshot>('state'));
        setNotice('Recovered your previous action. Please choose your next action again.');
        return;
      }
      const payload = {...command, request_id: crypto.randomUUID(), revision: world.revision};
      localStorage.setItem('black-ledger-pending', JSON.stringify(payload));
      const next = await api<Snapshot>('action', payload);
      sceneConditions.current={world:next.id,revision:next.revision,conditions:Object.fromEntries(world.locations.map(p=>[p.id,p.condition]))};
      localStorage.removeItem('black-ledger-pending');
      setWorld(next);
      setPlaying(null);
      const cues = next.last_result?.cues || [];
      const worst = cues.length
        ? [...cues].sort((a, b) => (b.gravity || 0) - (a.gravity || 0))[0]
        : undefined;
      if (worst && !next.event && motionRef.current) {
        setBeat(0);
        setFinishedCue('');
        setPlaying(worst);
        setCityView('iso');
        setTab('city');
        setSelected(worst.target);
      }
      if (
        command.kind === 'travel' &&
        !next.event &&
        next.player.alive &&
        motionRef.current &&
        !matchMedia('(prefers-reduced-motion: reduce)').matches
      ) {
        const from = world.locations.find(l => l.id === world.player.location),
          to = next.locations.find(l => l.id === next.player.location);
        if (from && to && from.id !== to.id) {
          setTab('city');
          setCityView('iso');
          setJourneyProgress(0);
          setJourney({
            street: next.last_result?.street_travel ?? undefined,
            fromMinute: world.minute,
            from,
            to,
            minutes: next.last_result?.elapsed || 0,
            driving: !!world.vehicle?.running,
            vehicle: world.vehicle?.car,
          });
        }
      }
      if (command.kind === 'new_life')
        setSelected(
          'bar',
        ); /* The result has a panel of its own now; a toast repeating its headline is
   the same news twice. Toasts are for what the panel cannot say: errors, and
   a recovered action. */
      if (
        next.player.job_count >= 2 &&
        !['writing', 'ready'].includes(next.director.status) &&
        next.minute - next.director.last_request > 180
      ) {
        await api('director', {});
        setWorld(await api<Snapshot>('state'));
      }
    } catch (err) {
      if (err instanceof RequestError && [400, 409].includes(err.status))
        localStorage.removeItem('black-ledger-pending');
      setNotice((err as Error).message);
      try {
        setWorld(await api<Snapshot>('state'));
      } catch {}
    } finally {
      busyRef.current = false;
      setBusy(false);
    }
  }
  // A hand already on the table is a game in progress, so the table opens
  // itself. It is also what makes leaving the only way out: the core keeps the
  // hand, so walking away from the screen would only hide it.
  const motionRef = useRef(motion);
  motionRef.current = motion;
  function toggleMotion() {
    const on = !motion;
    setMotion(on);
    try {
      localStorage.setItem('black-ledger-motion', on ? 'on' : 'off');
    } catch {}
    if (!on) {
      setPlaying(null);
      setJourney(null);
    }
  }
  function toggleSound() {
    const on = !sound;
    setSoundOn(on);
    setSound(on);
  }
  function toggleVoice() {
    const enabled = !voice;
    setVoice(enabled);
    localStorage.setItem('black-ledger-voice', enabled ? 'yes' : 'no');
    if (enabled) speak();
    else stopVoice();
  }
  if (!world)
    return (
      <div className="loading">
        BLACK LEDGER<span>{error || 'Opening the books…'}</span>
        {error && <button onClick={() => location.reload()}>Reconnect</button>}
      </div>
    );
  const p = world.player,
    locationInfo = world.locations.find(l => l.id === selected) || world.locations[0],
    event = world.event,
    npc = speakerOf(world.npcs, event?.speaker);
  function actionButton(a: Action) {
    if (a.sum)
      return (
        <SumAction key={a.id} a={a} money={money} disabled={busy || !!journey} commit={commit} />
      );
    // What you are buying, drawn. Three cars were three lines of text that
    // looked identical on the way past.
    const car = a.tier ? paintedCar(a.tier) : null;
    if (car)
      return (
        <button
          key={a.id}
          className="action car-card"
          title={[a.detail, a.disabled ? a.reason : ''].filter(Boolean).join(' — ')}
          disabled={a.disabled || busy || !!journey}
          onClick={() => commit({kind: a.id, target: a.target, choice: a.choice})}
        >
          <img src={car} alt="" loading="lazy" />
          <strong>{a.label}</strong>
          <span className="meta">
            {a.minutes ? `${a.minutes} min` : ''}
            {a.cost > 0 ? (
              <span>{money(a.cost)}</span>
            ) : a.asks ? (
              <span>{money(a.asks)}</span>
            ) : null}
          </span>
          <span className="desc">{a.reason || a.detail}</span>
        </button>
      );
    return (
      <button
        key={a.id}
        className={`action ${a.id === 'provoke' ? 'danger' : a.id === 'travel' ? 'primary' : ''}`}
        title={[a.detail, a.disabled ? a.reason : ''].filter(Boolean).join(' — ')}
        disabled={a.disabled || busy || !!journey}
        onClick={() => commit({kind: a.id, target: a.target, choice: a.choice})}
      >
        <strong>
          {a.label}
          {a.id === 'travel' ? ' ↗' : ''}
        </strong>
        <span className="meta">
          {a.minutes ? `${a.minutes} min` : a.away ? `${Math.round(a.away / 1440)} days away` : ''}
          {a.cost > 0 ? <span>{money(a.cost)}</span> : a.asks ? <span>{money(a.asks)}</span> : null}
        </span>
        <span className="desc">{a.reason || a.detail}</span>
      </button>
    );
  }
  function propertyPanel(l: Place) {
    return (
      <aside className="sidebar">
        <div className="person">
          <Portrait id={p.name} face={p.face} />
          <div>
            <b>{p.name}</b>
            <small>
              {p.respect < 6
                ? 'An unknown face'
                : p.crew.length
                  ? 'Crew leader'
                  : 'Neighborhood operator'}{' '}
              · Life {world!.life}
            </small>
            {world!.hand?.playing && (
              <small className="warning">
                At the tables in {world!.hand.place}: showing {world!.hand.player}, dealer shows{' '}
                {world!.hand.dealer}, {money(world!.hand.stake ?? 0)} down
              </small>
            )}
            {world!.armoury?.held && (
              <small className="warning">
                {world!.armoury.crates} of {world!.armoury.capacity} crates under{' '}
                {world!.armoury.place} · {world!.armoury.attention} attention a day ·{' '}
                {world!.armoury.buyers} families buying
              </small>
            )}
            {world!.service?.serving && (
              <small className="subtle">
                {world!.service.title} of {world!.service.name} · ${world!.service.pay}/day ·{' '}
                {world!.service.next
                  ? `${world!.service.next} more jobs to come up`
                  : 'as high as they go'}
              </small>
            )}
            {!!world!.pacts?.length && (
              <small className="subtle">
                Standing with{' '}
                {world!.pacts.map(p => p.name + (p.strength ? ` (${p.power})` : '')).join(' · ')} ·
                ${world!.pacts.reduce((n, p) => n + p.tribute, 0)}/day
              </small>
            )}
            {!!world!.retainers?.length && (
              <small className="subtle">
                Paying{' '}
                {world!.retainers.map(r => r.name + (r.outbid ? ' (outbid)' : '')).join(' · ')} · $
                {world!.retainers.reduce((n, r) => n + r.retainer, 0)}/day
              </small>
            )}
            {!!world!.arms?.charges && (
              <small className="warning">
                Carrying {world!.arms.charges} charge{world!.arms.charges > 1 ? 's' : ''} ·{' '}
                {world!.arms.charges * 4} attention a day
              </small>
            )}
            {!!world!.residence?.comforts?.length && (
              <small className="subtle">
                Home: {world!.residence!.comforts.map(c => c.label).join(' · ')} · $
                {world!.residence!.upkeep}/day
                {world!.residence!.sheltered > 0
                  ? ` · ${money(world!.residence!.sheltered)} out of reach`
                  : ''}
              </small>
            )}
            {world!.vehicle && world!.vehicle.car !== 'On foot and by streetcar' && (
              <small className="subtle">
                {world!.vehicle.car} · {world!.vehicle.condition}%
                {world!.vehicle.tank
                  ? ` · petrol ${world!.vehicle.fuel} of ${world!.vehicle.tank}`
                  : ''}
                {world!.vehicle.running
                  ? ` · $${world!.vehicle.upkeep}/day${world!.vehicle.concealed ? ` · hides ${world!.vehicle.concealed} units` : ''}${world!.vehicle.plate ? ` · plated ${world!.vehicle.plate} of ${world!.vehicle.plate_max}` : ''}`
                  : ' · will not start'}
              </small>
            )}
            {world!.appearance && (
              <small className="subtle">
                {world!.appearance.attire}
                {world!.appearance.standing > 0
                  ? ` · ${world!.appearance.condition}% kept · +${world!.appearance.standing} presence`
                  : world!.appearance.condition < 100
                    ? ` · ${world!.appearance.condition}% kept · worth nothing until it is put right`
                    : ''}
              </small>
            )}
            <div className="bar">
              <i style={{width: `${p.health}%`}} />
            </div>
          </div>
        </div>
        <div className="eyebrow">
          {l.locked
            ? 'BEYOND YOUR REACH'
            : l.id === p.location
              ? 'YOU ARE HERE'
              : 'NEIGHBORHOOD DIRECTORY'}{' '}
          / {['Old Harbor', 'Ashbury', 'The Heights'][l.district]}
        </div>
        <h2>{l.name}</h2>
        <p className="subtle">{l.blurb}</p>
        {l.id === p.location && l.room && (
          <p className="room-note">
            <Icon id="crew" />
            {l.room}
          </p>
        )}
        <div className={'building-art' + (paintedFront(l.id) ? ' street-front' : '')}>
          {paintedAsset(l.id) ? (
            <img
              src={paintedAsset(l.id, l.condition)!}
              alt=""
              style={
                paintedMask(l.id)
                  ? {
                      maskImage: `url(${paintedMask(l.id)})`,
                      WebkitMaskImage: `url(${paintedMask(l.id)})`,
                      maskSize: 'contain',
                      WebkitMaskSize: 'contain',
                      maskPosition: 'center',
                      WebkitMaskPosition: 'center',
                      maskRepeat: 'no-repeat',
                      WebkitMaskRepeat: 'no-repeat',
                    }
                  : undefined
              }
            />
          ) : (
            <div className="property-art-pending">
              <svg viewBox="0 0 80 64" aria-hidden="true">
                <path
                  d="M12 56V22L40 8l28 14v34H12Zm18 0V38h20v18M22 27h7m22 0h7M22 34h7m22 0h7M8 57h64"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                />
              </svg>
              <span>{l.name}</span>
            </div>
          )}
        </div>
        <div className="fact-grid">
          <div>
            <span>Ownership</span>
            <b>{l.holder || (l.owned ? 'Your organization' : 'Independent')}</b>
          </div>
          {l.owned && l.income > 0 && typeof l.trading === 'number' && (
            <div>
              <span>Working at</span>
              <b className={l.trading < 0.8 ? 'warning' : ''}>{Math.round(l.trading * 100)}%</b>
              <small>
                {l.hands?.length
                  ? l.hands.map(h => h.name + (h.here ? '' : ' (out)')).join(' · ')
                  : l.staff !== undefined
                    ? `${l.staff} on the books`
                    : ''}
                {l.supply !== undefined ? ` · ${l.supply} supplies` : ''}
                {l.trouble ? ' · trouble' : ''}
                {l.still ? ' · still running' : ''}
              </small>
            </div>
          )}
          <div>
            <span>{l.owned && l.income > 0 ? 'Hourly income' : 'Condition'}</span>
            <b>
              {l.owned && l.income > 0 ? money((l.income * l.condition) / 100) : l.condition + '%'}
            </b>
          </div>
          {l.owned && l.trade && (
            <div>
              <span>Trade</span>
              <b className={l.trade.custom < 40 ? 'warning' : ''}>{l.trade.custom}%</b>
              <small>
                {l.trade.order
                  ? `standing order · ${money(l.trade.order_pays)}/day`
                  : `regulars are worth ${Math.round(l.trade.multiplier * 100)}% of ordinary takings`}
              </small>
            </div>
          )}
          {l.owned && l.type === 'casino' && (
            <div>
              <span>Behind the tables</span>
              <b className={(l.bankroll ?? 0) < 500 ? 'warning' : ''}>{money(l.bankroll ?? 0)}</b>
              <small>
                {(l.bankroll ?? 0) === 0
                  ? 'The tables are dark'
                  : `covers about ${money(l.handle ?? 0)} of action a night`}
              </small>
            </div>
          )}
          {/* The poolhall is a racket rather than a casino, so none of the
              float rows reached it — and it is the one room that runs a card
              game and charges for the seat. Its money is a till: nothing is
              covered out of it, and what the table takes goes into it. */}
          {l.owned && l.id === 'poolhall' && (
            <div>
              <span>In the till</span>
              <b>{money(l.bankroll ?? 0)}</b>
              <small>
                {(l.bankroll ?? 0) === 0
                  ? 'Nothing in it yet'
                  : "the room's own money, and what the table has taken for the seat"}
              </small>
            </div>
          )}
          {l.owned && l.income > 0 && (
            <div>
              <span>Condition</span>
              <b className={l.condition < 70 ? 'warning' : ''}>
                {l.condition}%{l.condition < 100 ? ' · Repairs available' : ''}
              </b>
            </div>
          )}
        </div>
        {(() => {
          const seat = l.actions.find(a => a.id === 'sit');
          return (
            l.id === p.location &&
            seat && (
              <button
                className="action primary sit-down-here"
                disabled={seat.disabled || busy || !!journey}
                title={seat.disabled ? seat.reason : undefined}
                onClick={() => commit({kind: 'sit', target: l.id})}
              >
                <span>
                  <strong>{seat.label}</strong>
                  <span className="desc">{seat.disabled ? seat.reason : seat.detail}</span>
                </span>
              </button>
            )
          );
        })()}
        {l.id === p.location ? (
          <div className="here-instead">
            <button
              className="action primary"
              onClick={() => {
                setSelected(p.location);
                setCityView('interior');
              }}
            >
              <strong>Step inside {l.name} ↗</strong>
              <span className="desc">
                {l.actions.filter(a => !isTableAction(a.id) && !a.anywhere && !a.disabled).length}{' '}
                things you can do in here, and {(l.people || []).length}{' '}
                {(l.people || []).length === 1 ? 'person' : 'people'} standing in it.
              </span>
            </button>
            <p className="subtle">
              The room is where the work is. This column is for reading the city from where you are.
            </p>
          </div>
        ) : (
          <ActionList
            actions={l.actions}
            people={l.people || []}
            render={actionButton}
            groups={world!.groups}
            here={false}
          />
        )}
        <div className="bottom-note">
          <Icon id="clock" /> Decisions pause the clock. Commitments advance it.
        </div>
      </aside>
    );
  }
  // The work that belongs to the player rather than to the room they are in. The
  // core marks it; this is where it is filed, which is beside the people and the
  // families it is actually about.
  const anywhere = (
    world?.locations.find(l => l.id === world.player.location)?.actions || []
  ).filter(a => a.anywhere);
  const unreadNews = unreadInLatest(world?.newspaper || [], newsSeen);
  // The banner announces the biggest unread story, not the newest one. Once the
  // city page started filing the weather every morning, "newest" meant the
  // banner would announce cloud cover over a man being shot the same afternoon.
  const headline = (() => {
    const paper = world?.newspaper || [];
    const fresh = paper.slice(0, Math.max(unreadNews, 1));
    return fresh.reduce((best, s) => ((s.weight || 0) > (best.weight || 0) ? s : best), fresh[0]);
  })();
  function content(view=tab) {
    const w = world!;
    if (view === 'city') {
      const inside = cityView === 'interior' && locationInfo.id === p.location;
      const theatre = playing && !journey && (
        <Theatre
          focusOnStart={!p.alive}
          cue={playing}
          finished={cityView==='iso' ? finishedCue===playing.id : undefined}
          place={w.locations.find(l => l.id === playing.target) || w.locations[0]}
          onProgress={setBeat}
          plate={cityView !== 'iso'}
          stagedAudio={cityView === 'iso' && cityOwnsAudio(playing, w.last_result?.cues || [])}
          onDone={() => {if(sceneArticle&&playing)setFinishedCue(playing.id);else setPlaying(null);}}
        />
      );
      // Presentation controls remain usable after death; the gameplay shell
      // stays inert. Hide the controls while the newspaper owns focus.
      const sceneOverlay = !p.alive
        ? theatre && !newspaperVisible && !deathVisible && createPortal(
          <div className="map-first fatal-scene-playback"><div className="city3d-story">{theatre}</div></div>, document.body)
        : theatre;
      const journeyOverlay = journey &&
        (() => {
          const cross = w.locations.find(l => l.id === journey.to.id)?.crossing;
          return (
            <div
              className={'street-journey' + (cross?.warned ? ' warned' : '')}
              role="status"
            >
              <div>
                <strong>Crossing to {journey.to.name}</strong>
                <span>
                  {journey.minutes} minutes{' '}
                  {cross?.driving
                    ? `driving${cross.plate ? ` · ${cross.plate} of ${cross.plate_max} plated` : ' · no plate'}`
                    : 'on foot'}
                </span>
                {cross?.note && <small>{cross.note}</small>}
              </div>
              <button onClick={() => setJourney(null)}>Skip journey →</button>
            </div>
          );
        })();
      return (
        <div className={'workspace city-workspace' + (inside ? ' inside' : '')}>
          <section className="city-pane">
            <div className="city-stage">
              {inside&&<button data-shortcut="b" aria-keyshortcuts="B" className="map-leave-building" onClick={()=>setCityView('iso')}>Back to city ↗</button>}
              {inside && sceneOverlay}
              {cityView === 'interior' && locationInfo.id === p.location ? (
                <Interior
                  motion={motion}
                  player={p}
                  place={locationInfo}
                  people={locationInfo.people || []}
                  actions={locationInfo.actions.filter(
                    a => !isTableAction(a.id) && !a.anywhere && !isBackRoomAction(a.id),
                  )}
                  onTables={
                    locationInfo.actions.some(a => a.id === 'sit' && !a.disabled)
                      ? () => commit({kind: 'sit', target: locationInfo.id})
                      : undefined
                  }
                  felt={locationInfo.actions.some(a => a.id === 'play' || a.id === 'wheel')}
                  groups={w.groups}
                  comings={w.last_result?.comings}
                  minute={w.minute}
                  render={actionButton}
                  onLeave={() => setCityView('iso')}
                />
              ) : (
                <City3D
                  immersive
                  state={w}
                  replaySerial={sceneReplay.current}
                  beforeConditions={sceneConditions.current?.world===w.id&&sceneConditions.current.revision===w.revision?sceneConditions.current.conditions:undefined}
                  overlay={sceneOverlay || journeyOverlay}
                  activeCue={journey || !scenePending ? null : playing}
                  onSceneDone={setFinishedCue}
                  onJourneyDone={() => setJourney(null)}
                  onJourneyProgress={setJourneyProgress}
                  selected={selected}
                  onSelect={setSelected}
                  onTravel={id => commit({kind: 'travel', target: id})}
                  onEnter={() => {
                    setSelected(p.location);
                    setCityView('interior');
                  }}
                  motion={motion}
                  journey={journey}
                  busy={busy}
                />
              )}
              {inside && journeyOverlay}
            </div>
            {!playing && !w.event && !!w.last_result?.cues?.length && (
              <button
                className="replay-scene"
                onClick={() => {
                  const cues = w.last_result?.cues || [];
                  const cue = [...cues].sort((a, b) => (b.gravity || 0) - (a.gravity || 0))[0];
                  if (cue) {
                    sceneReplay.current++;
                    setCityView('iso');
                    setSelected(cue.target);
                    setBeat(0);
                    setFinishedCue('');
                    setPlaying(cue);
                  }
                }}
              >
                Replay recorded scene ↻
              </button>
            )}
          </section>
          {cityView === 'interior' && locationInfo.id === p.location
            ? null
            : propertyPanel(locationInfo)}
        </div>
      );
    }
    if (view === 'market') return <MarketScreen world={w} onFind={id=>{setSelected(id);setTab('city');setCityView('iso');}} />;
    if (view === 'ledger')
      return (
        <LedgerScreen
          world={w}
          render={actionButton}
          actions={anywhere.filter(a => a.id === 'bribe' || a.id === 'lie_low')}
        />
      );
    if (view === 'crew')
      return (
        <PeopleScreen
          world={w}
          actions={anywhere}
          render={actionButton}
          at={p.location}
          here={(w.locations.find(l => l.id === p.location)?.actions || []).filter(
            a => !!a.subject,
          )}
          onFind={id => {
            setSelected(id);
            setTab('city');
            setCityView('iso');
          }}
        />
      );
    if (view === 'families')
      return (
        <FamiliesScreen
          world={w}
          actions={anywhere}
          render={actionButton}
          onMeet={id => {
            const seat = id === 'bellandi' ? 'club' : 'garage';
            setSelected(seat);
            setTab('city');
            setCityView('iso');
          }}
        />
      );
    if (view === 'news')
      return (
        <section className="section-content">
          {!!w.arrangements?.length && (
            <>
              <div className="eyebrow">PAID FOR, NOT YET DONE</div>
              <h1 className="screen-title">Your arrangements</h1>
              {w.arrangements.map((a, i) => (
                <article className="card" key={i} style={{marginBottom: 14}}>
                  <h2 style={{margin: '0 0 6px'}}>{a.target}</h2>
                  <p>
                    {a.hired} · {money(a.paid)} paid. {a.status}
                  </p>
                </article>
              ))}
            </>
          )}
          <div className="eyebrow">THE BELLWETHER HERALD</div>
          <h1 className="screen-title">What the city is reading</h1>
          <p className="subtle">
            The paper prints what can be seen. It does not know who arranged anything, and reading
            it costs no time.
          </p>
          <Herald world={w} />
        </section>
      );
    if (view === 'settings')
      return (
        <section className="section-content">
          <div className="eyebrow">HOW THIS PLAYS</div>
          <h1 className="screen-title">Settings</h1>
          <div className="settings">
            <div className="setting">
              <div>
                <h3>Scenes</h3>
                <p>
                  When something happens that the city would remember — a killing, an arrest, a fire
                  — the game takes you there and holds for a moment before the headline. Turn this
                  off and the result is reported in words only.
                </p>
              </div>
              <button
                className={'toggle' + (motion ? ' on' : '')}
                role="switch"
                aria-checked={motion}
                onClick={toggleMotion}
              >
                <i />
                {motion ? 'Shown' : 'Off'}
              </button>
            </div>
            <div className="setting">
              <div>
                <h3>Sound</h3>
                <p>
                  An explosion, a shot, a police lamp turning over at the kerb. Recorded and synthesized effects, played when their matching actions occur.
                </p>
              </div>
              <button
                className={'toggle' + (sound ? ' on' : '')}
                role="switch"
                aria-checked={sound}
                onClick={toggleSound}
              >
                <i />
                {sound ? 'On' : 'Off'}
              </button>
            </div>
            <div className="setting">
              <div>
                <h3>Voices</h3>
                <p>
                  Named characters keep the same voice for as long as they live. Deciding anything
                  cuts them off, including a line still being spoken.
                </p>
              </div>
              <button
                className={'toggle' + (voice ? ' on' : '')}
                role="switch"
                aria-checked={voice}
                onClick={toggleVoice}
              >
                <i />
                {voice ? 'On' : 'Off'}
              </button>
            </div>
            <div className="setting">
              <div>
                <h3>The storyteller</h3>
                <p>{w.director.detail}</p>
                <p className="subtle">
                  A model writes the encounters. It cannot touch money, time, injuries, property or
                  death — the simulation owns all of those, and the game is playable with the model
                  switched off.
                </p>
              </div>
              <button
                className="plain"
                onClick={prepare}
                disabled={w.director.status === 'writing'}
              >
                {w.director.status === 'writing' ? 'Writing…' : 'Prepare an encounter'}
              </button>
            </div>
            <div className="setting face-setting">
              <div>
                <h3>Your face</h3>
                <p>
                  Nobody in Bellwether is gendered by the rules — the city says "they" about
                  everybody — so a portrait is dealt out by a hash of your name, and it can hand you
                  somebody you do not recognise as yourself. Pick your own. It is saved with the
                  life.
                </p>
                <div className="face-choices">
                  {Array.from({length: CAST_FACES}, (_, i) => i + 1).map(n => (
                    <button
                      key={n}
                      className={'face-choice' + (p.face === n ? ' chosen' : '')}
                      aria-pressed={p.face === n}
                      aria-label={'Face ' + n}
                      onClick={() => commit({kind: 'face', choice: String(n)})}
                    >
                      <Portrait id={p.name} face={n} size="small" />
                    </button>
                  ))}
                </div>
              </div>
              <button
                className="plain"
                disabled={!p.face}
                onClick={() => commit({kind: 'face', choice: '0'})}
              >
                Let the city decide
              </button>
            </div>
            <div className="setting">
              <div>
                <h3>This life</h3>
                <p>
                  Life {w.life} · {time(w.minute)} · saved after every action. There is one save and
                  no way back: whatever happens to {p.name} has happened.
                </p>
              </div>
              <span className="setting-note">Ironman</span>
            </div>
          </div>
        </section>
      );
    return (
      <section className="section-content help">
        {w.opportunity && <section className="guide-next" aria-label="Your next move"><div className="eyebrow">Your next move</div><h2>{w.opportunity.title}</h2><p>{w.opportunity.detail}</p><button onClick={()=>{setSelected(w.opportunity!.target);setTab('city');setCityView('iso');}}>Find the address ↗</button></section>}
            {!!w.grudges?.length && (
              <section className="known-threats" aria-label="What people are saying">
                <strong>Bad blood</strong>
                {w.grudges.map((g, i) => (
                  <p key={i}>
                    {g.holder} has not forgiven {g.against} for {g.because}.
                  </p>
                ))}
              </section>
            )}
            {!!w.commissions?.length && (
              <section className="known-threats" aria-label="Work you have taken on">
                <strong>What you owe people</strong>
                {w.commissions.map(c => (
                  <p key={c.id}>
                    <b>
                      {c.giver} · {c.patron}
                    </b>{' '}
                    — {c.brief} <i>{c.met ? 'Ready to settle.' : c.progress}</i> {money(c.pay)} ·{' '}
                    {Math.round(c.minutes_left / 60)}h left
                  </p>
                ))}
              </section>
            )}
            {!!w.known_threats?.length && (
              <section className="known-threats" aria-label="Known threats">
                <strong>Word on the street</strong>
                {w.known_threats.map((threat, i) => (
                  <p key={i}>{threat}</p>
                ))}
                <button className="plain" onClick={() => setTab('families')}>
                  Consider negotiations ↗
                </button>
              </section>
            )}

        <div className="eyebrow">WHERE YOU STAND</div>
        <h1 className="screen-title">What you can do, and what you cannot yet</h1>
        <p className="subtle">
          Follow the next move above, or choose your own path. The milestones below track what
          you have built in this life. Travel takes time; check the terms again when you arrive.
        </p>
        <ol className="guide-steps">
          {(w.guide || []).map(s => (
            <li key={s.title} className={s.done ? 'done' : s.open ? 'open' : 'shut'}>
              <div className="guide-mark" aria-hidden="true">
                {s.done ? '✓' : s.open ? '›' : '·'}
              </div>
              <div>
                <b>{s.title}</b>
                <p>{s.what}</p>
                {s.done ? (
                  <small className="subtle">Done.</small>
                ) : s.open ? (
                  <small className="ready">You can do this now.</small>
                ) : (
                  <small className="warning">{s.reason}</small>
                )}
              </div>
            </li>
          ))}
        </ol>
        <h2>Rules that do not change</h2>
        <ul className="guide-rules">
          {(w.rules || []).map((r, i) => (
            <li key={i}>{r}</li>
          ))}
        </ul>
      </section>
    );
  }
  return (
    <>
      <div className={'shell map-first ' + (busy ? 'busy' : '')} inert={!!event || !p.alive || atTable || newspaperVisible}>
        <nav className="rail" aria-label="Main navigation" inert={tab!=='city'}>
          <div className="monogram">
            <span>B</span>
          </div>
          {[
            ['crew', 'People'],
            ['families', 'Families'],
            ['market', 'Market'],
            ['ledger', 'Ledger'],
            ['news', 'Herald'],
            ['help', 'Guide'],
          ].map(([id, label], shortcutIndex) => (
            <button
              key={id}
              data-shortcut={String(shortcutIndex+1)} aria-keyshortcuts={String(shortcutIndex+1)}
              className={
                (tab === id ? 'active' : '') + (id === 'news' && unreadNews > 0 ? ' has-news' : '')
              }
              aria-current={tab === id ? 'page' : undefined}
              aria-label={
                id === 'news' && unreadNews > 0 ? `${label}, ${unreadNews} unread` : label
              }
              onClick={() => {setTab(id);if(id==='city')setCityView('iso');}}
            >
              <Icon id={id} />
              {label}<kbd aria-hidden="true">{shortcutIndex+1}</kbd>
              {id === 'news' && unreadNews > 0 ? <i className="news-count">{unreadNews}</i> : null}
            </button>
          ))}
          <button
            className="bottom" data-shortcut="v" aria-keyshortcuts="V"
            onClick={toggleVoice}
            aria-label={(voice ? 'Disable' : 'Enable') + ' voice acting'}
          >
            <Icon id={voice ? 'voice' : 'mute'} />
            Voice {voice ? 'on' : 'off'}
          </button>
          <button data-shortcut="7" aria-keyshortcuts="7" onClick={() => setTab('settings')} aria-label="Settings">
            <Icon id="settings" />
            Settings
          </button>
        </nav>
        <main className="page">
          <header className="topbar" inert={tab!=='city'}>
            <div className="hud-identity" aria-label={`Playing as ${p.name}`}><Portrait id={p.name} face={p.face} size="small"/><div><small>Bellwether · Life {world.life}</small><strong>{p.name}</strong><span className={cityView==='interior'?'hud-location':undefined}>{cityView==='interior' ? `Inside ${world.locations.find(place=>place.id===p.location)?.name||p.location}` : p.crew.length ? "Crew leader" : p.respect < 6 ? "An unknown face" : "Neighborhood operator"}</span></div></div>
            <div className="stats">
              {(world.dashboard || []).map(s => (
                <div className={'stat' + (s.warn ? ' warning' : '')} key={s.id} title={s.meaning}>
                  <Icon id={s.id} />
                  <b>{s.value}</b>
                  <small>
                    {s.label}
                    {s.note ? <i>{s.note}</i> : null}
                  </small>
                </div>
              ))}
              <div className="stat clock">
                <b>{time(journey ? world.minute - journey.minutes + Math.round(journey.minutes * journeyProgress) : world.minute)}</b>
                <small>{journey ? (journey.driving ? 'Driving through Bellwether' : 'Walking through Bellwether') : busy ? 'Resolving…' : 'Clock paused · awaiting your action'}</small>
              </div>
            </div>
          </header>
          <CityAccounts world={world}/>
          <div className="map-main-scene" inert={tab!=='city'}>{content('city')}</div>
          {!scenePending && !journey && world.last_result && (cityView==='interior' ? <details className="map-outcome map-outcome-folded"><summary>Latest entry <span>＋</span></summary><Outcome world={world} onLedger={()=>setTab('ledger')}/></details> : <div className="map-outcome"><Outcome world={world} onLedger={() => setTab('ledger')} /></div>)}
          {tab!=='city'&&<MapMenu edition={tab} title={({crew:'People',families:'Families',market:'Market',ledger:'Ledger',news:'The Bellwether Herald',settings:'Settings',help:'Guide'} as Record<string,string>)[tab]||tab} onClose={()=>setTab('city')}>{content()}</MapMenu>}
        </main>
      </div>
      {playing&&sceneArticle&&!event&&<SceneNewspaper key={`${playing.id}:${sceneReplay.current}`} article={sceneArticle} visible={newspaperVisible} voice={voice} onClose={()=>{setPlaying(null);setTab('city');}}/>}
      {atTable && inTheBackRoom && !event && p.alive && (
        <BackRoomScene
          place={world.locations.find(l => l.id === p.location)?.name || 'the back room'}
          cards={world.cards ?? null}
          seat={(world.locations.find(l => l.id === p.location)?.actions || []).find(
            a => a.id === 'cards',
          )}
          cash={p.cash}
          money={money}
          act={c => commit({target: p.location, ...c})}
          onLeave={() => commit({kind: 'rise', target: p.location})}
        />
      )}
      {atTable && !inTheBackRoom && !event && p.alive && (
        <Casino
          player={p}
          motion={motion}
          place={world.locations.find(l => l.id === p.location)?.name || 'the tables'}
          actions={(world.locations.find(l => l.id === p.location)?.actions || []).filter(a =>
            isTableAction(a.id),
          )}
          people={world.locations.find(l => l.id === p.location)?.people || []}
          hand={world.hand ?? {playing: false}}
          wheel={world.wheel ?? {spun: false}}
          dice={world.dice ?? {playing: false, settled: false}}
          machine={
            world.machine ?? {
              pulled: false,
              stops: 20,
              edge: 0,
              two_cherries: 0,
              one_cherry: 0,
              strip: [],
            }
          }
          house={world.house ?? {games: false}}
          cash={p.cash}
          money={money}
          revision={world.revision}
          records={world.last_result?.records || []}
          act={c => commit({target: p.location, ...c})}
          onLeave={() => commit({kind: 'rise', target: p.location})}
        />
      )}
      {(event || deathVisible) && (
        <div className="modal-shade encounter-shade">
          <section
            ref={scene}
            tabIndex={-1}
            className={'scene encounter-surface' + (!p.alive?' memorial-sheet':' private-meeting') + (event && event.choices.length > 3 ? ' extended' : '')}
            role="dialog"
            aria-modal="true"
            aria-labelledby="scene-title"
            onKeyDown={e => {
              if (e.key === 'Escape') stopVoice();
              if (e.key === 'Tab') {
                const buttons = [
                  ...scene.current!.querySelectorAll<HTMLButtonElement>('button:not(:disabled)'),
                ];
                if (!buttons.length) return;
                const first = buttons[0],
                  last = buttons.at(-1)!;
                if (
                  e.shiftKey &&
                  (document.activeElement === first || document.activeElement === scene.current)
                ) {
                  last.focus();
                  e.preventDefault();
                } else if (!e.shiftKey && document.activeElement === last) {
                  first.focus();
                  e.preventDefault();
                }
              }
            }}
          >
            {!p.alive ? (
              <>
                {(() => {
                  const e = world.epitaph;
                  return (
                    <>
                      <div className="eyebrow">THE CITY CONTINUES</div>
                      <div className="death-seal">✦</div>
                      <h2 id="scene-title">{p.name} is dead.</h2>
                      <p className="spoken">{e?.cause || world.dead.at(-1)?.cause}</p>
                      <div className="life-recap" aria-label="This life in numbers">
                        <div>
                          <small>Lived to</small>
                          <b>Day {e?.day ?? Math.floor(world.minute / 1440) + 1}</b>
                        </div>
                        <div>
                          <small>Final respect</small>
                          <b>{e?.respect ?? p.respect}</b>
                        </div>
                        <div>
                          <small>Money earned</small>
                          <b>{money(e?.earned ?? p.earned ?? 0)}</b>
                        </div>
                      </div>
                      {e && (
                        <p className="legacy-line">
                          <small>WHAT BECAME OF IT</small>
                          {e.became}
                        </p>
                      )}
                      {!!e?.standing?.length && (
                        <p className="legacy-properties">
                          Still standing in their name: {e.standing.join(' · ')}
                        </p>
                      )}
                      {!!e?.headlines?.length && (
                        <div className="legacy-press">
                          <small>WHAT THE PAPER CARRIED</small>
                          {e.headlines.map((h: string, n: number) => (
                            <b key={n}>{h}</b>
                          ))}
                        </div>
                      )}
                      {e && (
                        <p className="legacy-line">
                          <small>WHAT THE NEXT ONE GETS</small>
                          {e.inherits}
                        </p>
                      )}
                      <button
                        className="action primary"
                        disabled={busy}
                        onClick={() => commit({kind: 'new_life'})}
                      >
                        Begin as a new person <Icon id="arrow" />
                      </button>
                      <p className="source">
                        Life {world.life} · {time(world.minute)} · Outcome committed
                      </p>
                    </>
                  );
                })()}
              </>
            ) : (
              event && (
                <>
                 <div className="meeting-page">
                  <div className="eyebrow">
                    <span>
                      {event.kind === 'attack' ? 'A MOMENT TO ACT' : 'A PRIVATE CONVERSATION'}
                    </span>
                    <span>TIME PAUSED</span>
                  </div>
                  <h2 id="scene-title">{event.title}</h2>
                  {npc && (
                    <div className="person">
                      <Portrait id={npc.id} />
                      <div>
                        <b>{npc.name}</b>
                        <small>{npc.role}</small>
                      </div>
                    </div>
                  )}
                  {event.connection && (
                    <aside className="story-connection" aria-label="Previous arrangement">
                      <small>COMPLETED WORK WITH THIS CONTACT</small>
                      <b>{event.connection.title}</b>
                      <p>{event.connection.result}</p>
                    </aside>
                  )}
                  <p className="spoken">{event.body}</p>
                  <div className="speech-control">
                    <button className="plain" onClick={speak}>
                      <Icon id="voice" /> {speech}
                    </button>
                    {speaking(speech) && (
                      <button className="plain" onClick={stopVoice}>
                        Stop voice
                      </button>
                    )}
                  </div>
                  {event.conditions ? (
                    <p className="scene-conditions">
                      <small>WHATEVER YOU CHOOSE</small>
                      {event.conditions}
                    </p>
                  ) : null}
                 </div>
                 <div className="meeting-decisions">
                  <div className="meeting-reply-label">Your reply</div>
                  <div className="choices">
                    {event.choices.map((c, i) => (
                      <button
                        key={c.id}
                        data-choice
                        className="action"
                        disabled={c.disabled || busy}
                        onClick={() => commit({kind: 'choice', choice: c.id, event: event.id})}
                      >
                        <span className="number">{String(i + 1).padStart(2, '0')}</span>
                        <span>
                          <strong>{c.label}</strong>
                          {c.pay || c.minutes ? (
                            <span className="terms">
                              {c.pay ? <b>{money(c.pay)}</b> : null}
                              {c.minutes ? <b>{c.minutes} min</b> : null}
                              {c.respect ? <b>+{c.respect} respect</b> : null}
                              {c.heat ? <b className="cost">+{c.heat} attention</b> : null}
                            </span>
                          ) : null}
                          {c.detail ? <span className="desc">{c.detail}</span> : null}
                          {c.disabled && c.reason ? (
                            <span className="desc refused">{c.reason}</span>
                          ) : null}
                        </span>
                      </button>
                    ))}
                  </div>
                  <p className="source">
                    {event.source === 'local-ai' ? 'Local AI encounter' : 'City encounter'} ·
                    Decisions are saved immediately.
                  </p>
                 </div>
                </>
              )
            )}
          </section>
        </div>
      )}
      <KeyboardShortcuts/>
      {notice && (
        <div id="toast" role="status" aria-live="polite" style={{display: 'block'}}>
          {notice}
        </div>
      )}
    </>
  );
}
createRoot(document.getElementById('root')!).render(<App />);
