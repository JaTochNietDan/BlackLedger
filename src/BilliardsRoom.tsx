import {useState} from 'react';
import {BilliardsTable3D} from './BilliardsTable3D';
import {poolPockets,poolPlacementHint} from './billiards';
import type {PoolState,PoolOpponent} from './billiards';
import type {Command} from './types';
import './billiards.css';

export function PoolChallenges({opponents,cash,busy,act}:{opponents:PoolOpponent[];cash:number;busy:boolean;act:(c:Command)=>void}){
 const [stake,setStake]=useState(20);
 return <section className="pool-challenges"><h3>A game at Green Baize</h3><label>Your stake $<input aria-label="Billiards stake" type="number" min={10} max={Math.min(500,cash)} step={10} value={stake} onChange={e=>setStake(Number(e.target.value))}/></label>
 {opponents.length===0&&<p>No opponent is on the floor at the moment.</p>}
 {opponents.map(n=><button key={n.id} title={n.unavailable||`Offers up to $${n.max_stake}`} disabled={busy||!!n.unavailable||stake<10||!Number.isInteger(stake)||stake>Math.min(500,cash,n.max_stake)} onClick={()=>act({kind:'pool_start',target:n.id,amount:stake})}>Play {n.name} · ${stake} each</button>)}<small>Eight-ball · both stakes held until the rack ends</small></section>;
}
export function BilliardsRoom({pool,busy,motion,act,tournament}:{pool:PoolState;busy:boolean;motion:boolean;act:(c:Command)=>void;tournament?:{names:[string,string];observer:boolean;onBack:()=>void}}){
 const [angle,setAngle]=useState(90),[speed,setSpeed]=useState(pool.breaking?7.5:4),[top,setTop]=useState(0),[side,setSide]=useState(0),[ball,setBall]=useState(1),[pocket,setPocket]=useState(0),[safety,setSafety]=useState(false),[x,setX]=useState(pool.width/2),[y,setY]=useState(.4),[playing,setPlaying]=useState(false);
 const locked=busy||playing,playerTurn=pool.turn===0&&!tournament?.observer,canPlay=!locked&&!pool.unavailable&&!pool.settled;
 const target=pool.legal_balls.includes(ball)?ball:(pool.legal_balls[0]??8);
 const placementHint=poolPlacementHint(pool,x,y);
 const names=tournament?.names??['You',pool.opponent_name];
 const groups=['Open table','Solids','Stripes'];
 return <section className="billiards-room" role="dialog" aria-modal="true" aria-label="Green Baize billiards">
  <header><div><small>GREEN BAIZE · {tournament?'TOURNAMENT':'EIGHT-BALL'}</small><h2>{names[0]} &amp; {names[1]}</h2><p>{pool.settled?(pool.voided?'Event cancelled':pool.winner>=0?`${names[pool.winner]} won the rack`:'Game ended'):`$${pool.pot} ${tournament?'tournament prize':'in the middle'} · ${names[pool.turn]} to shoot`} · {names[0]}: {groups[pool.groups[0]]}</p></div>
   <button disabled={locked} onClick={()=>tournament?tournament.onBack():act({kind:pool.settled?'pool_close':'pool_concede'})}>{tournament?'Back to the draw':pool.settled?'Return to the hall':`Concede · lose $${pool.stake}`}</button></header>
  <BilliardsTable3D replayLocked={busy} pool={pool} angle={angle*Math.PI/180} top={top/1000} side={side/1000} motion={motion} locked={busy||!!tournament?.observer} calledBall={target} calledPocket={pocket} placement={[x,y]} onPlaying={setPlaying} onAim={r=>setAngle(Number((r*180/Math.PI).toFixed(2)))} onBall={setBall} onPocket={setPocket} onPlace={(x,y)=>{setX(Number(x.toFixed(3)));setY(Number(y.toFixed(3)));}}/>
  <div className="pool-controls">
   <p aria-live="polite">{playing?'Playing the shot…':pool.unavailable||pool.outcome||(tournament?.observer?'Watch the next stroke, or move the camera to inspect the table.':pool.ball_in_hand?'Click the cloth to position your cue ball, then confirm.':'Click to aim; click a ball and a numbered pocket to call your shot.')} </p>
   {canPlay&&playerTurn&&!pool.ball_in_hand&&!pool.break_choices.length&&<p className="pool-pointer-help">Click cloth to aim. Click a ball and a numbered pocket to call them. Drag to move the camera.</p>}
   {!pool.settled&&<fieldset disabled={!canPlay}>
    {!playerTurn?<button onClick={()=>act({kind:'pool_opponent'})}>Watch {names[pool.turn]} play</button>:pool.break_choices.length?<div>{pool.break_choices.map(c=><button key={c.id} onClick={()=>act({kind:'pool_decide',choice:c.id})}>{c.label}</button>)}</div>:pool.ball_in_hand?<div className="pool-inputs">
     <details className="pool-precise-placement"><summary>Precise placement</summary><div className="pool-inputs">
     <label>Across cloth (m)<input type="number" min={pool.radius} max={pool.width-pool.radius} step={.01} value={x} onChange={e=>setX(Number(e.target.value))}/></label>
     <label>From near end (m)<input type="number" min={pool.radius} max={pool.behind_head_string?pool.length/4-.001:pool.length-pool.radius} step={.01} value={y} onChange={e=>setY(Number(e.target.value))}/></label>
     </div></details>
     <span className="pool-placement-note" role="status">{placementHint||'Preview ready — confirm to place the ball.'}</span>
     <button disabled={!!placementHint} onClick={()=>act({kind:'pool_place',pool:{x,y}})}>Place cue ball{pool.behind_head_string?' behind head string':''}</button>
    </div>:<>
     <div className="pool-inputs">
      <label>Aim {angle.toFixed(1)}°<input type="range" min={0} max={359.9} step={.1} value={angle} onChange={e=>setAngle(Number(e.target.value))}/><input aria-label="Precise aim in degrees" type="number" min={0} max={360} step={.1} value={angle} onChange={e=>setAngle(Number(e.target.value))}/></label>
      <label>Cue speed {speed.toFixed(1)} m/s<input type="range" min={.1} max={8} step={.1} value={speed} onChange={e=>setSpeed(Number(e.target.value))}/></label>
      <label>Draw / follow {top} mm<input type="range" min={-12} max={12} step={1} value={top} onChange={e=>setTop(Number(e.target.value))}/></label>
      <label>Left / right spin {side} mm<input type="range" min={-12} max={12} step={1} value={side} onChange={e=>setSide(Number(e.target.value))}/></label>
      {!pool.breaking&&<><label>Called ball<select value={target} onChange={e=>setBall(Number(e.target.value))}>{pool.legal_balls.map(n=><option key={n} value={n}>{n}</option>)}</select></label><label>Called pocket<select value={pocket} onChange={e=>setPocket(Number(e.target.value))}>{poolPockets.map((name,i)=><option key={name} value={i}>{i+1} · {name}</option>)}</select></label><label><input type="checkbox" checked={safety} onChange={e=>setSafety(e.target.checked)}/> Safety</label></>}
      <button onClick={()=>act({kind:'pool_shot',pool:{angle:angle*Math.PI/180,speed,top:top/1000,side:side/1000,ball:target,pocket,safety:!pool.breaking&&safety}})}>{pool.breaking?'Break the rack':'Play shot'}</button>
     </div>
    </>}
   </fieldset>}
   <details><summary>Table controls &amp; house rules</summary><p>Click the cloth to aim or preview cue placement; confirm placement with the button. Click a ball to call it and aim toward its centre. Click a numbered pocket to call it. Drag to orbit, right-drag to pan, wheel to zoom. Focus the table for WASD, Q/E, +/− and Home. Aim 90° points toward the far end. The line ends at the first ball or cushion; the ring marks the cue ball’s centre at contact. It does not predict spin, rebounds or whether a ball will drop.</p><p>Call the ball and pocket after the break. Hit your group first, then pocket a ball or reach a cushion. A safety passes the turn. Fouls give ball in hand; pocketing the eight early, in the wrong pocket or with a scratch loses the rack. Four object balls must reach cushions on the break unless a ball is pocketed. Illegal breaks and break eights offer explicit decisions.</p><p>This table currently models cloth contact, rolling, sliding and spin. Jump shots and massé are not available.</p></details>
  </div>
 </section>;
}
