import {useState} from 'react';
import type {Command,Snapshot} from './types';
import './crewOrders.css';

export function CrewOrdersPanel({world,busy,act}:{world:Snapshot;busy:boolean;act:(c:Command)=>unknown}) {
 const [actor,setActor]=useState(''),[kind,setKind]=useState('restock'),[target,setTarget]=useState('');
 const offers=world.crew_order_offers||[];
 const people=[...new Map(offers.map(o=>[o.actor,o.name])).entries()];
 const selectedActor=people.some(([id])=>id===actor)?actor:people[0]?.[0]||'';
 const choices=offers.filter(o=>o.actor===selectedActor&&o.kind===kind);
 const selected=choices.find(o=>o.target===target)||choices[0];
 const base=world.organization?.headquarters;
 const here=!!base&&world.player.location===base;
 const orders=(world.crew_orders||[]).slice().reverse();
 const active=orders.filter(o=>['outbound','working','returning','guarding'].includes(o.stage));
 const completed=orders.filter(o=>!['outbound','working','returning','guarding'].includes(o.stage));
 const [showCompleted,setShowCompleted]=useState(false);
 const visibleOrders=showCompleted?[...active,...completed]:active;
 return <section className="crew-orders" aria-label="Headquarters operations">
  <div className="eyebrow">THE FAMILY'S WORK</div><h2>Headquarters operations</h2>
  <p>{here?'Orders continue while you attend to other business.':'Issue orders from your headquarters. Existing assignments continue while you are away.'}</p>
  {here&&<>
   {!people.length?<p>Hire an associate or sign people into your family to give orders.</p>:<>
    <div className="crew-order-fields"><label>Operative<select value={selectedActor} onChange={e=>setActor(e.target.value)}>{people.map(([id,name])=><option key={id} value={id}>{name}</option>)}</select></label>
    <label>Operation<select value={kind} onChange={e=>{setKind(e.target.value);setTarget('');}}><option value="restock">Restock a business</option><option value="collections">Make collections</option><option value="guard">Guard a business</option><option value="repair">Repair a property</option><option value="remedy">Resolve business trouble</option><option value="rob">Rob a business</option><option value="bomb">Bomb a business</option><option value="assassinate">Assassinate a known target</option></select></label>
    <label className="crew-order-target">Target<select value={selected?.target||''} onChange={e=>setTarget(e.target.value)}>{choices.map(o=><option key={o.target} value={o.target}>{o.label}</option>)}</select></label></div>
    {selected?<><p className="crew-order-terms">About {selected.minutes} minutes {kind==='guard'?'to take the post · stays until recalled':'including travel'} · ${selected.cost} reserved{selected.charges?` · ${selected.charges} charge committed`:""}{kind==='collections'?' · Payment is carried home; lost if the operative is killed or captured':kind==='bomb'?' · Risk of premature detonation, deaths and retaliation':kind==='rob'?' · Proceeds arrive when the operative returns; injury and retaliation possible':kind==='assassinate'?' · Risk of death, arrest and retaliation':''}</p><button className="crew-order-dispatch" disabled={busy||!!selected.reason} onClick={()=>act({kind:`crew_order:${selected.kind}`,choice:selected.actor,target:selected.target})}>Dispatch {selected.name}</button>{selected.reason&&<p role="status">{selected.reason}</p>}</>:<p>No known targets for this operation.</p>}
   </>}
  </>}
  {orders.length>0&&<div className="crew-order-register" role="region" tabIndex={0} aria-label="Assignment register"><div className="crew-order-register-heading"><b>{active.length} on assignment</b>{completed.length>0&&<button className="plain" aria-expanded={showCompleted} onClick={()=>setShowCompleted(!showCompleted)}>{showCompleted?'Hide':'Show'} completed ({completed.length})</button>}</div>{visibleOrders.map(o=><article key={o.id}><b>{o.name} · {{collections:'Collections',guard:'Guard duty',restock:'Restocking',repair:'Repairs',remedy:'Business trouble',rob:'Robbery',bomb:'Demolition',assassinate:'Assassination'}[o.kind]||o.kind}</b><span>{world.locations.find(p=>p.id===o.place)?.name||o.place} · {{guarding:'On guard',outbound:'On the way',working:'At work',returning:'Returning',done:'Reported back',cancelled:'Cancelled'}[o.stage]||o.stage}{['outbound','working','returning'].includes(o.stage)?` · next stage in ${Math.max(0,o.due-world.minute)} min`:''}</span>{o.result&&<p>{o.result}</p>}{here&&!o.recall&&['outbound','working','guarding'].includes(o.stage)&&<button className="plain" disabled={busy} onClick={()=>act({kind:'crew_recall',choice:o.id})}>Recall</button>}{o.recall&&o.stage!=='done'&&<small>Recall received</small>}</article>)}</div>}
 </section>;
}
