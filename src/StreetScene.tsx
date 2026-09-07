import {useEffect,useRef,useState} from 'react';
import {paintedBuildings,paintedLocations} from './cityAssets';
interface Props {selected:string;minute:number;onSelect:(id:string)=>void}
/** Isolated presentation adapter. Only whitelisted inspection messages cross this boundary. */
export function StreetScene({selected,minute,onSelect}:Props){
 const [motion,setMotion]=useState(localStorage.getItem('black-ledger-motion')!=='off');
 const frame=useRef<HTMLIFrameElement>(null),latest=useRef({selected,minute,onSelect,motion});latest.current={selected,minute,onSelect,motion};
 function sync(){frame.current?.contentWindow?.postMessage({type:'blackledger:presentation',selected:latest.current.selected,minute:latest.current.minute,motion:latest.current.motion},location.origin)}
 useEffect(()=>{const receive=(e:MessageEvent)=>{if(e.source!==frame.current?.contentWindow||e.origin!==location.origin)return;if(e.data?.type==='blackledger:ready')sync();if(e.data?.type==='blackledger:inspect'&&paintedLocations.includes(e.data.location))latest.current.onSelect(e.data.location)};window.addEventListener('message',receive);return()=>window.removeEventListener('message',receive)},[]);
 useEffect(sync,[selected,minute,motion]);
 return <div className="street-scene"><iframe ref={frame} title="Old Harbor illustrated street" src="/street-study.html?embed=1" onLoad={sync}/><div className="street-inspect">{paintedBuildings.map(b=><button key={b.id} aria-pressed={selected===b.id} onClick={()=>onSelect(b.id)}>{b.name}</button>)}<button aria-pressed={motion} onClick={()=>{localStorage.setItem('black-ledger-motion',motion?'off':'on');setMotion(!motion)}}>{motion?'Pause ambience':'Resume ambience'}</button></div></div>
}
