import {useEffect,useRef,useState} from 'react';
import type {Journey} from './TravelPresentation';
import {paintedBuildings,paintedLocations} from './cityAssets';
interface Props {properties:{id:string;condition:number}[];selected:string;minute:number;position:string;journey:Journey|null;onSelect:(id:string)=>void}
/** Isolated presentation adapter. Only whitelisted inspection messages cross this boundary. */
export function StreetScene({properties,selected,minute,position,journey,onSelect}:Props){
 const [motion,setMotion]=useState(localStorage.getItem('black-ledger-motion')!=='off');
 const frame=useRef<HTMLIFrameElement>(null),latest=useRef({properties,selected,minute,position,journey,onSelect,motion});latest.current={properties,selected,minute,position,journey,onSelect,motion};
 function sync(){frame.current?.contentWindow?.postMessage({type:'blackledger:presentation',properties:latest.current.properties.map(p=>({id:p.id,condition:p.condition})),selected:latest.current.selected,minute:latest.current.minute,motion:latest.current.motion,position:latest.current.position,journey:latest.current.journey?{from:latest.current.journey.from.id,to:latest.current.journey.to.id}:null},location.origin)}
 useEffect(()=>{const receive=(e:MessageEvent)=>{if(e.source!==frame.current?.contentWindow||e.origin!==location.origin)return;if(e.data?.type==='blackledger:ready')sync();if(e.data?.type==='blackledger:inspect'&&paintedLocations.includes(e.data.location))latest.current.onSelect(e.data.location)};window.addEventListener('message',receive);return()=>window.removeEventListener('message',receive)},[]);
 useEffect(sync,[properties,selected,minute,motion,position,journey]);
 return <div className="street-scene"><iframe ref={frame} title="Old Harbor illustrated street" src="/street-study.html?embed=1" onLoad={sync}/><div className="street-inspect">{paintedBuildings.map(b=><button key={b.id} aria-pressed={selected===b.id} onClick={()=>onSelect(b.id)}>{b.name}</button>)}<button aria-pressed={motion} onClick={()=>{localStorage.setItem('black-ledger-motion',motion?'off':'on');setMotion(!motion)}}>{motion?'Pause ambience':'Resume ambience'}</button></div></div>
}
