import {CrewOrdersPanel} from './CrewOrdersPanel';
import type {ReactElement} from 'react';
import type {Action, Snapshot, Command} from './types';
import {Portrait} from './Portrait';

// Three thin cards: a strength, a money word and a bare number for standing —
// "+45", "-63" — with nothing about what any of them holds, who they are
// fighting or what the number means. Everything needed was already in the city;
// it was simply never asked for.

export function FamiliesScreen({
  world,
  onMeet,
  act,
  busy=false,
  actions = [],
  render,
}: {
  world: Snapshot;
  act?: (c:Command)=>unknown;
  busy?:boolean;
  onMeet: (id: string) => void;
  // Understandings and what is being said about a family are reached through
  // other people. Neither of them happens at a counter, and both used to be
  // printed at one.
  actions?: Action[];
  render?: (a: Action) => ReactElement;
}) {
  const mine = world.factions.filter(f => f.yours);
  const others = world.factions.filter(f => !f.yours);
  const pacts = new Set((world.pacts || []).map(p => p.id));
  const truce = world.business_truces || {};

  const card = (f: Snapshot['factions'][number]) => {
    const hostile = (f.goodwill ?? 0) <= -35;
    const warm = (f.goodwill ?? 0) >= 35;
    return (
      <article
        className={
          'family-card' + (f.yours ? ' yours' : hostile ? ' hostile' : warm ? ' warm' : '')
        }
        key={f.id}
      >
        <header>
          {/* The player's own organization is led by the player, who is not one
            of the city's people and so has no id to draw from. Their own face
            is the one the rest of the interface already uses. */}
          <Portrait
            id={f.leader_id || (f.yours ? world.player.name : f.id)}
            face={f.yours && !f.leader_id ? world.player.face : undefined}
            size="small"
          />
          <div>
            <b>{f.name}</b>
            <small>{f.leader ? `Led by ${f.leader}` : 'Nobody will say who runs it'}</small>
            <small className={hostile ? 'warning' : f.yours ? 'subtle' : ''}>{f.standing}</small>
          </div>
        </header>

        <div className="family-facts">
          <div>
            <span>Strength</span>
            <b>{f.strength}</b>
          </div>
          <div>
            <span>Money</span>
            <b>{f.money}</b>
          </div>
          <div>
            <span>People</span>
            <b>{f.hands}</b>
          </div>
          <div>
            <span>Ground</span>
            <b>{f.holdings?.length || 0}</b>
          </div>
        </div>

        {f.headquarters && <p className="family-ground"><span>Headquarters</span> {world.locations.find(p=>p.id===f.headquarters)?.name||f.headquarters}</p>}
        {f.yours&&!f.headquarters&&<p className="family-note warning">No usable headquarters deed. Establish your base at an owned business.</p>}

        {!!f.holdings?.length && (
          <p className="family-ground">
            <span>Holds</span> {f.holdings.join(' · ')}
          </p>
        )}

        {!!f.fighting?.length && <p className="family-fighting">{f.fighting.join(' · ')}</p>}

        {pacts.has(f.id) && <p className="family-note">You have an understanding with them.</p>}
        {truce[f.id] !== undefined && (
          <p className="family-note">A ceasefire holds over business, and no further.</p>
        )}

        {f.knowledge < 2 && !f.yours && (
          <p className="family-note subtle">
            {f.knowledge === 0
              ? 'Nobody will talk to you about them. Asking around is done through other people, from wherever you are.'
              : 'What you know of them is second-hand. Somebody inside would tell you more.'}
          </p>
        )}

        {!!render &&
          !f.yours &&
          (() => {
            // Reaching an understanding, ending one, and asking what is being said:
            // all of it is carried by other people, so it is offered here rather
            // than at whatever counter it used to be printed at.
            const theirs = actions.filter(a => a.id.endsWith(':' + f.id));
            return theirs.length > 0 && <div className="actions">{theirs.map(render)}</div>;
          })()}

        {!f.yours && (
          <button className="plain" onClick={() => onMeet(f.id)}>
            Meet {f.leader || 'them'} ↗
          </button>
        )}
      </article>
    );
  };

  return (
    <section className="section-content">
      <div className="eyebrow">POWER WAS HERE BEFORE YOU</div>
      <h1 className="screen-title">The other families</h1>
      <p className="subtle">
        The city does not scale its dangers to your experience. Some doors are better left unopened.
      </p>

      {!world.organization?.named&&<p className="family-note">Form a family at an owned business and choose it as headquarters. You need 25 respect; acquiring businesses no longer forms a family automatically.</p>}

      {!!world.known_threats?.length && (
        <section className="known-threats" style={{margin: '0 0 22px'}}>
          <strong>Reported threats</strong>
          {world.known_threats.map((t, i) => (
            <p key={i}>{t}</p>
          ))}
          <p>
            Business ceasefires do not settle personal threats. An audience may offer separate
            terms.
          </p>
        </section>
      )}

      {mine.length > 0 && (
        <>
          <h2>Yours</h2>
          <div className="family-grid">{mine.map(card)}</div>
          {act&&<CrewOrdersPanel world={world} busy={busy} act={act}/>}
        </>
      )}

      <h2>{mine.length > 0 ? 'Everybody else' : 'The organizations'}</h2>
      <div className="family-grid">{others.map(card)}</div>

      {!!world.conflicts?.length && (
        <>
          <h2>Who is fighting whom</h2>
          <ul className="conflict-list">
            {world.conflicts.map((c, i) => (
              <li key={i} className={c.state}>
                <b>{c.between.join(' and ')}</b>
                <span>
                  {c.state === 'war' ? 'at war' : 'at odds'} since day{' '}
                  {Math.floor(c.since / 1440) + 1}
                </span>
              </li>
            ))}
          </ul>
        </>
      )}
    </section>
  );
}
