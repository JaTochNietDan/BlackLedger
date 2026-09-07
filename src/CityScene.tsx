import {useEffect,useRef,useState} from 'react';
import {Application,Container,Graphics,Sprite,Texture} from 'pixi.js';
import {citySVG} from './art';
import type {Snapshot} from './types';

/** Replaceable presentation adapter. No simulation rules, rewards, or random outcomes. */
export function CityScene({state,selected,onSelect,hideMarker=false}:{state:Snapshot;selected:string;onSelect:(id:string)=>void;hideMarker?:boolean}){
 const host=useRef<HTMLDivElement>(null),appRef=useRef<Application|null>(null),sequence=useRef(0),callback=useRef(onSelect),[failed,setFailed]=useState(false);callback.current=onSelect;
 useEffect(()=>{let disposed=false;const app=new Application();app.init({backgroundAlpha:0,antialias:true,resolution:Math.min(devicePixelRatio,2),autoDensity:true,autoStart:false,preference:'webgl'}).then(()=>{if(disposed){app.destroy(true);return}appRef.current=app;host.current?.append(app.canvas);draw()}).catch(()=>setFailed(true));
 const observer=new ResizeObserver(()=>draw());if(host.current)observer.observe(host.current);
 async function draw(){if(!appRef.current)return;appRef.current.renderer.resize(host.current?.clientWidth||800,host.current?.clientHeight||560);appRef.current.render()}
 return()=>{disposed=true;sequence.current++;observer.disconnect();if(appRef.current){appRef.current.destroy(true,{children:true,texture:true});appRef.current=null}}},[]);
 useEffect(()=>{let disposed=false;const token=++sequence.current;let timer:number;
 async function update(){const app=appRef.current;if(!app){timer=window.setTimeout(update,30);return}try{let svg=citySVG(state,selected).replace('<svg class="map"','<svg xmlns="http://www.w3.org/2000/svg" width="1620" height="1110" class="map"');if(hideMarker)svg=svg.replace(/<g id="player-marker"[\s\S]*?<\/g>/,'');svg=svg.replace('<defs>','<defs><style>.property-label{font-family:Arial;font-size:13px;fill:#e1ddc9;paint-order:stroke;stroke:#192622;stroke-width:4px}.property-sub{font-family:Arial;font-size:9px;fill:#beaa7c}.district{font-family:Arial;font-size:14px;letter-spacing:4px;fill:#95a18d;opacity:.6}.locked{opacity:.48}.selected .footprint{stroke:#d6b77c;stroke-width:2}.waterline{stroke:#5b7777;stroke-width:1;opacity:.23}</style>');
 const blob=new Blob([svg],{type:'image/svg+xml'}),url=URL.createObjectURL(blob),img=new Image();try{await new Promise<void>((resolve,reject)=>{img.onload=()=>resolve();img.onerror=reject;img.src=url})}finally{URL.revokeObjectURL(url)}if(disposed||sequence.current!==token)return;
 const texture=Texture.from(img);const stage=new Container();const picture=new Sprite(texture);picture.width=1080;picture.height=740;stage.addChild(picture);
 for(const p of state.locations){const hit=new Graphics().rect(p.x-65,p.y-65,130,120).fill({color:0xffffff,alpha:.001});hit.eventMode='static';hit.cursor='pointer';hit.on('pointertap',()=>callback.current(p.id));stage.addChild(hit)}
 app.stage.removeChildren().forEach(c=>c.destroy({children:true,texture:true,textureSource:true}));app.stage.addChild(stage);
 const fit=()=>{if(disposed||!appRef.current)return;const w=host.current?.clientWidth||800,h=host.current?.clientHeight||560;app.renderer.resize(w,h);const scale=Math.min(w/1080,h/740);stage.scale.set(scale);stage.position.set((w-1080*scale)/2,(h-740*scale)/2);app.render()};fit();const observer=new ResizeObserver(fit);if(host.current)observer.observe(host.current);cleanup=()=>observer.disconnect();
 }catch{if(!disposed)setFailed(true)}}
 let cleanup=()=>{};update();return()=>{disposed=true;clearTimeout(timer);cleanup()}},[state,selected,hideMarker]);
 return <div className="city-renderer" ref={host}>{failed&&<div className="map-fallback" dangerouslySetInnerHTML={{__html:citySVG(state,selected)}}/>}<div className="map-accessible" aria-label="City properties">{state.locations.map(p=><button key={p.id} onClick={()=>onSelect(p.id)}>{p.name}{p.locked?' (locked district)':''}</button>)}</div></div>
}
