import {poolHostPreview} from './billiards';
import {useState} from 'react';
import {BilliardsRoom} from './BilliardsRoom';
import type {PoolTournamentState,PoolTournamentNoticeState} from './billiards';
import type {Command} from './types';

function postedTime(minute:number){return `Day ${Math.floor(minute/1440)+1}, ${String(Math.floor(minute%1440/60)).padStart(2,'0')}:${String(minute%60).padStart(2,'0')}`;}
export function PoolTournamentNotice({notice,tournament,busy,act,onOpen,cash}:{notice:PoolTournamentNoticeState|null|undefined;tournament:PoolTournamentState|null|undefined;busy:boolean;act:(c:Command)=>void;onOpen:()=>void;cash:number}){
 const [fee,setFee]=useState(25),[cut,setCut]=useState(10),[enter,setEnter]=useState(false);
 if(!notice)return null;
 const preview=poolHostPreview(notice.host_offers??[],fee,cut,enter,cash);
 return <section className="pool-challenges"><h3>Green Baize tournament notice</h3><p>Every third evening · entry {postedTime(notice.opens)}–{postedTime(notice.closes)} · ${notice.fee} each.</p>
 {tournament&&<button disabled={busy} onClick={onOpen}>View tournament draw</button>}
 {(!tournament||tournament.settled)&&<><p>{notice.entrants.length} selected opponents for the public draw{notice.can_enter?` · $${notice.pot} winner's prize`:''}. {notice.unavailable}</p><button disabled={busy||!notice.can_enter} onClick={()=>{act({kind:'pool_tournament_enter'});onOpen();}}>Pay ${notice.fee} and enter</button></>}
 {notice.can_host&&(!tournament||tournament.settled)&&<fieldset disabled={busy}><legend>Arrange a tournament as proprietor</legend><div className="pool-inputs"><label>Entry fee $<input type="number" min={10} max={500} value={fee} onChange={e=>setFee(Number(e.target.value))}/></label><label>House cut %<input type="number" min={0} max={50} value={cut} onChange={e=>setCut(Number(e.target.value))}/></label><label><input type="checkbox" checked={enter} onChange={e=>setEnter(e.target.checked)}/>Enter and pay the same fee</label><button disabled={!!preview.reason} onClick={()=>{act({kind:'pool_tournament_host',pool_host:{fee,cut_percent:cut,enter}});onOpen();}}>Start owner tournament</button></div><p role="status">{preview.reason||`${preview.count} entrants · $${preview.gross} total entries · $${preview.prize} prize · $${preview.cut} house cut.`}</p>{!preview.reason&&<p>{enter?'You, ':''}{preview.entrants.map(n=>n.name).join(', ')}</p>}<small>Uses funded players already in the hall. Higher fees may leave too few entrants. The cut is paid only when a champion wins; as owner you receive it as income. Settings are fixed once play begins.</small></fieldset>}
 <small>Scheduled public tournaments pay every entry fee to the champion. Owner tournaments deduct the posted cut. Leaving or withdrawing forfeits your fee. If the hall closes, unforfeited entries are returned; forfeited fees go to the hall if there is no champion.</small></section>;
}
export function PoolTournamentRoom({tournament,busy,motion,act,onClose,initialGame=null}:{tournament:PoolTournamentState;busy:boolean;motion:boolean;act:(c:Command)=>void;onClose:()=>void;initialGame?:number|null}){
 const [selected,setSelected]=useState<number|null>(initialGame);
 const eliminated=tournament.games.some(g=>g.resolved&&g.players.includes(tournament.player_id)&&g.winner!==tournament.player_id);
 const finalRound=Math.max(...tournament.games.map(g=>g.round));
 const paid=tournament.prize_paid??0;
 const status=tournament.voided?'Unforfeited entries returned.':tournament.settled?(paid>0?`$${paid} prize paid${tournament.winner===tournament.player_id?` · $${paid-tournament.fee} profit after your entry`:''}.`:'The tournament is over.'):tournament.entered===false?'You are hosting. The field plays without you.':tournament.withdrawn?'You have withdrawn.':eliminated?'You are out of the draw. You can watch the remaining matches.':'Win each round to reach the final.';
 const game=tournament.games.find(g=>g.index===selected);
 const name=(id:string)=>id===tournament.player_id?'You':tournament.names[id]||'Awaiting qualifier';
 if(game?.table)return <BilliardsRoom key={game.index} pool={game.table} busy={busy} motion={motion} tournament={{names:[name(game.players[0]),name(game.players[1])],observer:game.player_seat<0,onBack:()=>setSelected(null)}} act={c=>act({...c,kind:c.kind.replace('pool_','pool_tournament_'),pool_game:game.index})}/>;
 return <section className="billiards-room" role="dialog" aria-modal="true" aria-label="Green Baize tournament draw"><header><div><small>GREEN BAIZE · TOURNAMENT DRAW</small><h2>{tournament.settled?(tournament.voided?'Tournament cancelled':`${name(tournament.winner)} won the tournament`):`$${tournament.pot} to the champion`}</h2><p>${tournament.fee} per entrant · {tournament.house_cut_percent??0}% house cut · Eight-ball · {status}{(tournament.house_cut_paid??0)>0&&` House cut paid: $${tournament.house_cut_paid}.`}</p></div><button disabled={busy} onClick={onClose}>Return to the hall</button></header><div className="pool-controls">
 {tournament.games.map(g=><section className="pool-challenges" key={g.index}><h3>{g.round===finalRound?'Final':`Round ${g.round+1}`} {g.table_number>0?` · Table ${g.table_number}`:''}</h3><p>{name(g.players[0])} &amp; {name(g.players[1])} · {g.resolved?(g.winner?`${name(g.winner)} ${g.round===finalRound?'won the final':g.winner===tournament.player_id?'advance':'advances'}`:'No qualifier'):g.table?'In play':'Awaiting earlier games'}</p>{g.table&&<button disabled={busy} onClick={()=>setSelected(g.index)}>{g.resolved?'View rack':g.player_seat>=0?'Play your match':'Watch this table'}</button>}</section>)}
 {!tournament.settled&&tournament.games.some(g=>!g.resolved&&g.table&&g.player_seat<0)&&<button disabled={busy} onClick={()=>act({kind:'pool_tournament_wait'})}>Let tables play · 10 min</button>}
 {tournament.entered!==false&&!tournament.settled&&!tournament.withdrawn&&!eliminated&&<button disabled={busy} onClick={()=>act({kind:'pool_tournament_withdraw'})}>Withdraw · forfeit ${tournament.fee} entry</button>}
 <p>{tournament.settled?'The tournament is complete. You can return to the hall or review any played rack.':<>{tournament.entered===false?'You may leave the building; the tournament continues as game time passes.':'Returning to the hall keeps your entry. Leaving the building forfeits it.'} Other NPC tables play a shot every two game minutes. Open a table to watch its next stroke and follow the draw.</>}</p>
 </div></section>;
}
