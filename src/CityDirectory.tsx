import {useState} from 'react';
import type {Snapshot} from './types';
import {paintedAsset,paintedBuildings} from './cityAssets';

const districts=['Old Harbor','Ashbury','The Heights'];
/** An ordinary accessible destination browser. Selecting never commits travel. */
export function CityDirectory({state,selected,onSelect}:{state:Snapshot;selected:string;onSelect:(id:string)=>void}){
 const [district,setDistrict]=useState<number|null>(null);
 const locations=state.locations.filter(p=>district===null||p.district===district);
 return <section className="destination-directory" aria-label="City properties">
  <div className="directory-intro"><span>THE ADDRESS BOOK</span><p>Select a destination to inspect its opportunities. Travel is a separate decision.</p></div>
  <div className="district-filters" aria-label="Filter by district"><button aria-pressed={district===null} onClick={()=>setDistrict(null)}>All districts</button>{districts.map((name,i)=><button key={name} aria-pressed={district===i} onClick={()=>setDistrict(i)}>{name}{i>state.district?' · Locked':''}</button>)}</div>
  <div className="destination-grid">{locations.map(p=>{
   const asset=paintedAsset(p.id,p.condition),original=paintedBuildings.find(b=>b.id===p.id);
   const owner=p.owned?'Your organization':p.owner.startsWith('former:')?p.owner.slice(7)+'’s survivors':state.factions.find(f=>f.id===p.owner)?.name||'Independent';
   return <button key={p.id} className={'destination-card'+(p.locked?' locked-destination':'')} aria-pressed={selected===p.id} onClick={()=>onSelect(p.id)}>
    <span className="destination-image" aria-hidden="true">{asset?<img src={asset} alt="" style={original?{maskImage:`url(/art/${original.file})`,maskSize:'contain',maskPosition:'center',maskRepeat:'no-repeat'}:undefined}/>:<svg viewBox="0 0 80 64"><path d="M12 56V22L40 8l28 14v34H12Zm18 0V38h20v18M22 27h7m22 0h7M22 34h7m22 0h7M8 57h64" fill="none" stroke="currentColor" strokeWidth="2"/></svg>}</span>
    <span className="destination-copy"><small>{districts[p.district]||'New district'} · {p.type}</small><strong>{p.name}</strong><span className={p.owned?'owned-label':''}>{owner}</span><span className="destination-status">{p.locked?'District not yet open':p.id===state.player.location?'You are here':p.id===state.player.home?'Your residence':'Inspect destination →'}{p.condition<100&&` · ${p.condition}% condition`}</span></span>
   </button>
  })}</div>
 </section>
}
