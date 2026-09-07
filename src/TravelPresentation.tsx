import {useEffect,useRef,useState} from 'react';
import type {Place} from './types';
export interface Journey {from:Place;to:Place;minutes:number}
/** Plays a committed journey. Skipping changes no game state and sends no command. */
export function TravelPresentation({journey,onDone}:{journey:Journey;onDone:()=>void}) {
 const [progress,setProgress]=useState(0),done=useRef(onDone);done.current=onDone;
 useEffect(()=>{let frame=0;const start=performance.now();const tick=(now:number)=>{const p=Math.min(1,(now-start)/2200);setProgress(p);if(p<1)frame=requestAnimationFrame(tick);else done.current()};frame=requestAnimationFrame(tick);return()=>cancelAnimationFrame(frame)},[journey]);
 const {from,to,minutes}=journey;
 const start={x:from.x,y:from.y+24},end={x:to.x,y:to.y+24};
 // A schematic route, not autonomous pathfinding or simulation.
 const corner={x:end.x,y:start.y};const a=Math.abs(corner.x-start.x),b=Math.abs(end.y-corner.y),distance=(a+b)*progress;
 const x=distance<=a?start.x+Math.sign(corner.x-start.x)*distance:end.x;
 const y=distance<=a?start.y:corner.y+Math.sign(end.y-corner.y)*(distance-a);
 return <div className="travel-presentation"><svg viewBox="0 0 1080 740" aria-hidden="true"><path d={`M${start.x} ${start.y} L${corner.x} ${corner.y} L${end.x} ${end.y}`} fill="none" stroke="#e1c083" strokeWidth="2" strokeDasharray="5 7"/><g transform={`translate(${x},${y})`}><circle r="15" fill="#142522" stroke="#f2d398" strokeWidth="2"/><circle cy="-4" r="4" fill="#f2d398"/><path d="M-6 7Q0-4 6 7" stroke="#f2d398" strokeWidth="3" fill="none"/></g></svg><div className="travel-caption" role="status"><div><b>{to.name}</b><small>{minutes} game minutes · journey saved</small></div><button className="plain" onClick={()=>done.current()}>Skip journey →</button></div></div>
}
