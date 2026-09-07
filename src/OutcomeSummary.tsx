import type {Snapshot,Record as CityRecord} from './types';

const importance=(record:CityRecord)=>({death:6,danger:5,intel:4,politics:3,story:2,personal:1}[record.kind]??0);

/** Keep a consequential action visible when routine accounts settle in the same command. */
export function OutcomeSummary({world,onLedger}:{world:Snapshot;onLedger:()=>void}){
 const current=world.last_result?.records||[];
 const records=current.length?current:[...world.history].reverse().slice(0,1);
 if(!records.length)return null;
 const headline=[...records].sort((a,b)=>importance(b)-importance(a))[0];
 const others=records.filter(r=>r.id!==headline.id);
 return <section className={'ticker outcome '+headline.kind} aria-label="Latest developments">
  <div className="outcome-main"><b>{headline.title}</b><p>{headline.text}</p>
   {others.length>0&&<details key={headline.id}><summary>{others.length} more {others.length===1?'development':'developments'} during this action</summary><ul>{others.map(r=><li key={r.id}><b>{r.title}</b><p>{r.text}</p></li>)}</ul></details>}
  </div><button className="plain" onClick={onLedger} aria-label="Open the full ledger">Ledger ↗</button>
 </section>;
}
