import {useEffect,useState} from 'react';
import type {Snapshot} from './types';
import {pressPlate} from './art';

// The Herald was a rolling list: one flat run of the current life's stories,
// oldest silently evicted, and a previous protagonist's era unreachable. A city
// that remembers ought to let somebody read what it remembered — so the paper
// is an archive of issues now, one a day, and you can walk back through them.

export function Herald({world}: {world: Snapshot}) {
  const issues = world.editions || [];
  const [at, setAt] = useState(0);
  const [showAll, setShowAll] = useState(false);

  // A new day is a new front page: jump to it rather than leaving the reader
  // parked on an old issue they were only browsing.
  useEffect(() => { setAt(0) }, [issues[0]?.day, issues[0]?.count]);

  if (!issues.length) return <div className="paper-sheet"><div className="paper">
    <p className="nothing">No edition has gone to press yet.</p>
  </div></div>;

  const issue = issues[Math.min(at, issues.length - 1)];
  const older = at < issues.length - 1;
  const newer = at > 0;
  const lives = [...new Set(issues.map(i => i.life))];

  return <>
    <div className="issue-bar">
      <button className="plain" disabled={!newer} onClick={() => setAt(a => a - 1)}>← Later</button>
      <div className="issue-which">
        <b>{issue.current ? 'Today’s edition' : `Day ${issue.day}`}</b>
        <small>{issue.dateline} · {issue.count} {issue.count === 1 ? 'story' : 'stories'}{issue.mine ? '' : ` · life ${issue.life}`}</small>
      </div>
      <button className="plain" disabled={!older} onClick={() => setAt(a => a + 1)}>Earlier →</button>
      <button className="plain" aria-expanded={showAll} onClick={() => setShowAll(s => !s)}>
        {showAll ? 'Hide' : 'All'} {issues.length} issues
      </button>
    </div>

    {showAll && <ul className="issue-index">
      {issues.map((i, n) => <li key={i.life + '-' + i.day}>
        <button className={'plain' + (n === at ? ' picked' : '') + (i.mine ? '' : ' past-life')}
          onClick={() => { setAt(n); setShowAll(false) }}>
          <b>{i.current ? 'Today' : 'Day ' + i.day}</b>
          <small>{i.dateline}</small>
          <small>{i.count} {i.count === 1 ? 'story' : 'stories'}{i.mine ? '' : ` · life ${i.life}`}</small>
        </button>
      </li>)}
    </ul>}

    {lives.length > 1 && !issue.mine && <p className="past-life-note">
      This issue was printed before you arrived in Bellwether. The city was already running.
    </p>}

    {/* The sheet and the print are two elements on purpose. The ragged edge is
        a mask on the paper, and a mask is applied after a filter, so a shadow
        on the same element would be cast by the rectangle the mask cut away
        rather than by the torn edge. The shadow belongs to the wrapper. */}
    <div className="paper-sheet">
    <div className="paper">
      <div className="masthead">
        <h1>The Bellwether Herald</h1>
        <div className="rule">
          <span>{issue.dateline}</span>
          <span>{issue.current ? 'Late city edition' : 'Back issue'}</span>
          <span>Day {issue.day} · Five cents</span>
        </div>
      </div>
      <div className="columns">
        {issue.stories.filter(s => s.kind !== 'civic' && s.kind !== 'obituary').map((s, i) => <article className={i === 0 ? 'lead' : ''} key={s.id}>
          {s.subject && <figure className="cut" dangerouslySetInnerHTML={{__html:
            pressPlate(s.kind, s.subject, s.headline) +
            (s.subject.kind !== 'city' ? `<figcaption>${s.subject.name}</figcaption>` : '')}}/>}
          <h2>{s.headline}</h2>
          {s.standfirst && <p className="standfirst">{s.standfirst}</p>}
          <p className="byline">{s.byline} · {s.time}</p>
          <p className="body">{s.body}</p>
        </article>)}
      </div>

      {/* Obituaries. Set apart and set differently, because this is the only
          place in the paper that is about a person rather than an event: a
          rule above it, the name in small capitals, and no halftone cut — a
          picture of the building somebody died at is not an obituary. */}
      {issue.stories.some(s => s.kind === 'obituary') && <div className="obituaries">
        <h3>Obituaries</h3>
        {issue.stories.filter(s => s.kind === 'obituary').map(s => <article key={s.id}>
          <h4>{s.headline}</h4>
          <p>{s.body}</p>
        </article>)}
      </div>}

      {/* The city page: weather, prices, how many people are in the place.
          Set apart from the news because it is not news — it is what a paper
          carries on a day when nothing happened, and a paper that only speaks
          when the city is being violent is not a paper. */}
      {issue.stories.some(s => s.kind === 'civic') && <div className="city-page">
        <h3>The city in brief</h3>
        <div className="briefs">
          {issue.stories.filter(s => s.kind === 'civic').map(s => <article key={s.id}>
            <h4>{s.headline}</h4>
            <p>{s.body}</p>
          </article>)}
        </div>
      </div>}
    </div>
    </div>
  </>;
}
