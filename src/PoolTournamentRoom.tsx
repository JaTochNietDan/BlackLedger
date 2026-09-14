import {useState} from 'react';
import {BilliardsRoom} from './BilliardsRoom';
import type {PoolTournamentState,PoolTournamentNoticeState} from './billiards';
import type {Command} from './types';

function postedTime(minute:number){return `Day ${Math.floor(minute/1440)+1}, ${String(Math.floor(minute%1440/60)).padStart(2,'0')}:${String(minute%60).padStart(2,'0')}`;}
export function PoolTournamentNotice({notice,tournament,busy,act,onOpen}:{notice:PoolTournamentNoticeState|null|undefined;tournament:PoolTournamentState|null|undefined;busy:boolean;act:(c:Command)=>void;onOpen:()=>void}){
 if(!notice)return null;
 return <section className="pool-challenges"><h3>Green Baize tournament notice</h3><p>Every third evening · entry {postedTime(notice.opens)}–{postedTime(notice.closes)} · ${notice.fee} each.</p>
 {tournament&&<button disabled={busy} onClick={onOpen}>View tournament draw</button>}
 {(!tournament||tournament.settled)&&<><p>{notice.entrants.length} eligible opponents here{notice.can_enter?` · $${notice.pot} winner's prize`:''}. {notice.unavailable}</p><button disabled={busy||!notice.can_enter} onClick={()=>{act({kind:'pool_tournament_enter'});onOpen();}}>Pay ${notice.fee} and enter</button></>}
 <small>The champion takes every entry fee. Leaving or withdrawing forfeits your fee. If the hall closes, unforfeited entries are returned; forfeited fees go to the hall if there is no champion.</small></section>;
}
export function PoolTournamentRoom({tournament,busy,motion,act,onClose,initialGame=null}:{tournament:PoolTournamentState;busy:boolean;motion:boolean;act:(c:Command)=>void;onClose:()=>void;initialGame?:number|null}){
 const [selected,setSelected]=useState<number|null>(initialGame);
 const eliminated=tournament.games.some(g=>g.resolved&&g.players.includes(tournament.player_id)&&g.winner!==tournament.player_id);
 const finalRound=Math.max(...tournament.games.map(g=>g.round));
 const paid=tournament.prize_paid??0;
 const status=tournament.voided?'Unforfeited entries returned.':tournament.settled?(paid>0?`$${paid} prize paid${tournament.winner===tournament.player_id?` · $${paid-tournament.fee} profit after your entry`:''}.`:'The tournament is over.'):tournament.withdrawn?'You have withdrawn.':eliminated?'You are out of the draw. You can watch the remaining matches.':'Win each round to reach the final.';
 const game=tournament.games.find(g=>g.index===selected);
 const name=(id:string)=>id===tournament.player_id?'You':tournament.names[id]||'Awaiting qualifier';
 if(game?.table)return <BilliardsRoom key={game.index} pool={game.table} busy={busy} motion={motion} tournament={{names:[name(game.players[0]),name(game.players[1])],observer:game.player_seat<0,onBack:()=>setSelected(null)}} act={c=>act({...c,kind:c.kind.replace('pool_','pool_tournament_'),pool_game:game.index})}/>;
 return <section className="billiards-room" role="dialog" aria-modal="true" aria-label="Green Baize tournament draw"><header><div><small>GREEN BAIZE · TOURNAMENT DRAW</small><h2>{tournament.settled?(tournament.voided?'Tournament cancelled':`${name(tournament.winner)} won the tournament`):`$${tournament.pot} to the champion`}</h2><p>${tournament.fee} per entrant · Eight-ball · {status}</p></div><button disabled={busy} onClick={onClose}>Return to the hall</button></header><div className="pool-controls">
 {tournament.games.map(g=><section className="pool-challenges" key={g.index}><h3>{g.round===finalRound?'Final':`Round ${g.round+1}`} {g.table_number>0?` · Table ${g.table_number}`:''}</h3><p>{name(g.players[0])} &amp; {name(g.players[1])} · {g.resolved?(g.winner?`${name(g.winner)} ${g.round===finalRound?'won the final':g.winner===tournament.player_id?'advance':'advances'}`:'No qualifier'):g.table?'In play':'Awaiting earlier games'}</p>{g.table&&<button disabled={busy} onClick={()=>setSelected(g.index)}>{g.resolved?'View rack':g.player_seat>=0?'Play your match':'Watch this table'}</button>}</section>)}
 {!tournament.settled&&!tournament.withdrawn&&!eliminated&&<button disabled={busy} onClick={()=>act({kind:'pool_tournament_withdraw'})}>Withdraw · forfeit ${tournament.fee} entry</button>}
 <p>{tournament.settled?'The tournament is complete. You can return to the hall or review any played rack.':<>Returning to the hall keeps your entry. Leaving the building forfeits it. Other NPC tables play a shot every two game minutes. Open a table to watch its next stroke while you await your opponent.</>}</p>
 </div></section>;
}
