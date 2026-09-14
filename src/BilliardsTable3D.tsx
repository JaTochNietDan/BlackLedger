import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {TableCamera} from './tableCamera';
import {disposeCityResources} from './city3dResources';
import {decodePoolReplay,poolFramePair} from './billiards';
import type {PoolState,PoolReplay} from './billiards';
import {playTable} from './sound';

const colours=['#f2ead8','#e6ad23','#2248a0','#b7312b','#613774','#c66322','#236548','#761e26','#151615'];
function ballTexture(id:number){
 const canvas=document.createElement('canvas');canvas.width=512;canvas.height=256;const c=canvas.getContext('2d')!;
 c.fillStyle=id>8?'#f2ead8':colours[id];c.fillRect(0,0,512,256);
 if(id>8){c.fillStyle=colours[id-8];c.fillRect(0,77,512,102);}
 if(id)for(const x of [128,384]){c.fillStyle='#f7efda';c.beginPath();c.arc(x,128,36,0,Math.PI*2);c.fill();c.fillStyle='#151515';c.font='bold 49px Georgia';c.textAlign='center';c.textBaseline='middle';c.fillText(String(id),x,130);}
 const t=new THREE.CanvasTexture(canvas);t.colorSpace=THREE.SRGBColorSpace;return t;
}
export function BilliardsTable3D({pool,angle,motion,onPlaying}:{pool:PoolState;angle:number;motion:boolean;onPlaying:(playing:boolean)=>void}){
 const host=useRef<HTMLDivElement>(null),live=useRef({pool,angle,motion,onPlaying});live.current={pool,angle,motion,onPlaying};
 const [error,setError]=useState(''),[playback,setPlayback]=useState(false);
 useEffect(()=>{
  const el=host.current!;let dead=false,renderer:THREE.WebGLRenderer;
  try{renderer=new THREE.WebGLRenderer({antialias:true});}catch{setError('The 3D table could not start.');return;}
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Billiards table. Drag to orbit; right drag to pan; wheel to zoom; Home to reset.');el.append(canvas);
  renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  const scene=new THREE.Scene();scene.background=new THREE.Color('#15231f');scene.add(new THREE.HemisphereLight(0xffedcb,0x192c26,2));
  const light=new THREE.DirectionalLight(0xffe0ad,3);light.position.set(-1,5,1);light.castShadow=true;light.shadow.mapSize.set(1024,1024);Object.assign(light.shadow.camera,{left:-2,right:2,top:3,bottom:-3});light.shadow.bias=-.00015;scene.add(light);
  const W=pool.width,L=pool.length,R=pool.radius,H=.78;
  const camera=new THREE.PerspectiveCamera(38,1,.01,30);let dirty=true;
  const view=new TableCamera(camera,canvas,()=>{dirty=true;},new THREE.Vector3(0,H,0),3.6,Math.PI/2);
  const material=(colour:string)=>new THREE.MeshStandardMaterial({color:colour,roughness:.65});
  const wood=material('#39251b'),felt=material('#285b49'),cushion=material('#204a39');
  function box(w:number,h:number,d:number,x:number,y:number,z:number,mat:THREE.Material){const m=new THREE.Mesh(new THREE.BoxGeometry(w,h,d),mat);m.position.set(x,y,z);m.castShadow=true;m.receiveShadow=true;scene.add(m);return m;}
  box(12,.1,12,0,-.06,0,material('#252821'));
  for(const x of [-W/2-.10,W/2+.10])box(.18,.28,L+.28,x,H-.19,0,wood);
  for(const z of [-L/2-.10,L/2+.10])box(W+.38,.28,.18,0,H-.19,z,wood);
  for(const x of [-W*.36,W*.36])for(const z of [-L*.38,L*.38])box(.16,.54,.16,x,.27,z,wood);
  const pockets=[[ -.026,-.026],[ -.026,L+.026],[W+.026,-.026],[W+.026,L+.026],[-.045,L/2],[W+.045,L/2]];
  // The bed extends beneath the cushions; holes are actual mesh apertures.
  const shape=new THREE.Shape();shape.moveTo(-.13,-.13);shape.lineTo(W+.13,-.13);shape.lineTo(W+.13,L+.13);shape.lineTo(-.13,L+.13);shape.closePath();
  for(const [x,y] of pockets){const hole=new THREE.Path();hole.absarc(x,y,.076,0,Math.PI*2,true);shape.holes.push(hole);
   const cup=new THREE.Mesh(new THREE.CylinderGeometry(.076,.063,.15,32,1,true),material('#171411'));cup.position.set(x-W/2,H-.075,L/2-y);scene.add(cup);
  }
  const bed=new THREE.Mesh(new THREE.ShapeGeometry(shape,48),felt);bed.rotation.x=-Math.PI/2;bed.position.set(-W/2,H,L/2);bed.receiveShadow=true;scene.add(bed);
  const rails:number[][]=[[.085,0,W-.085,0],[.085,L,W-.085,L]];
  for(const x of [0,W]){const outside=x===0?-.045:W+.045;rails.push([x,.085,x,L/2-.075],[x,L/2+.075,x,L-.085],[x,L/2-.075,outside,L/2-.065],[x,L/2+.075,outside,L/2+.065]);}
  for(const x of [0,W])for(const y of [0,L]){const sx=x===0?1:-1,sy=y===0?1:-1;rails.push([x+sx*.085,y,x+sx*.049,y-sy*.035],[x,y+sy*.085,x-sx*.035,y+sy*.049]);}
  // Narrow cushion noses follow the solver segments including angled jaws.
  for(const [ax,ay,bx,by] of rails){const rail=box(Math.hypot(bx-ax,by-ay),.045,.025,(ax+bx)/2-W/2,H+.016,L/2-(ay+by)/2,cushion);rail.rotation.y=Math.atan2(by-ay,bx-ax);}
  for(const x of [-W/2-.10,W/2+.10])for(let i=1;i<8;i++){const diamond=new THREE.Mesh(new THREE.CircleGeometry(.008,4),material('#e0d0a6'));diamond.rotation.x=-Math.PI/2;diamond.position.set(x,H-.045,L/2-L*i/8);scene.add(diamond);}
  const balls=new Map<number,THREE.Mesh>();const geometry=new THREE.SphereGeometry(R,32,24);
  for(let id=0;id<16;id++){const mesh=new THREE.Mesh(geometry,new THREE.MeshStandardMaterial({map:ballTexture(id),roughness:.22,metalness:0}));mesh.castShadow=true;mesh.receiveShadow=true;balls.set(id,mesh);scene.add(mesh);}
  const guide=new THREE.Line(new THREE.BufferGeometry(),new THREE.LineDashedMaterial({color:'#e7d6a5',dashSize:.04,gapSize:.025,transparent:true,opacity:.65}));scene.add(guide);
  const cue=new THREE.Mesh(new THREE.CylinderGeometry(.005,.014,1.42,16),material('#b88d54'));scene.add(cue);
  let tape:PoolReplay|null=null,started=0,loading=false,key='',generation=0,playing=false,eventCursor=0;
  const preference=matchMedia('(prefers-reduced-motion: reduce)');
  const notify=(v:boolean)=>{if(v!==playing){playing=v;setPlayback(v);live.current.onPlaying(v);}};
  const skip=()=>{generation++;tape=null;loading=false;notify(false);dirty=true;};canvas.addEventListener('pool-skip',skip);
  const worldRotation=new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(1,0,0),-Math.PI/2),qa=new THREE.Quaternion(),qb=new THREE.Quaternion();
  function pose(id:number,x:number,y:number,z:number,q:number[],pocket:number){const mesh=balls.get(id)!;mesh.visible=pocket<0;mesh.position.set(x-W/2,H+z,L/2-y);qa.fromArray(q);mesh.quaternion.copy(worldRotation).multiply(qa);}
  const resize=()=>{renderer.setSize(el.clientWidth,Math.max(1,el.clientHeight));camera.aspect=el.clientWidth/Math.max(1,el.clientHeight);camera.updateProjectionMatrix();view.resize();dirty=true;};const observer=new ResizeObserver(resize);observer.observe(el);resize();
  let frame=0,lastAngle=NaN,lastPool:PoolState|null=null;
  const tick=(now:number)=>{
   frame=requestAnimationFrame(tick);const p=live.current.pool;if(lastPool!==p){lastPool=p;dirty=true;}
   const next=`${p.shots}:${p.replay}`;
   if(next!==key){const first=key==='';key=next;const token=++generation;tape=null;eventCursor=0;dirty=true;
    if(!first&&p.replay&&live.current.motion&&!preference.matches){loading=true;notify(true);decodePoolReplay(p.replay).then(r=>{if(dead||token!==generation)return;tape=r;started=performance.now();loading=false;dirty=true;}).catch(()=>{if(dead||token!==generation)return;loading=false;notify(false);setError('Replay unavailable. The saved table below is the authoritative result.');});}
    else {loading=false;notify(false);}
   }
   if((!live.current.motion||preference.matches)&&(tape||loading)){tape=null;loading=false;notify(false);generation++;dirty=true;}
   if(tape){const elapsed=(now-started)/1000;const {a,b,mix}=poolFramePair(tape,elapsed);
    const nextBalls=new Map(b.balls.map(ball=>[ball[0],ball]));
    for(const ba of a.balls){const bb=nextBalls.get(ba[0])!;qa.fromArray(ba.slice(4,8));qb.fromArray(bb.slice(4,8));qa.slerp(qb,mix);pose(ba[0],ba[1]+(bb[1]-ba[1])*mix,ba[2]+(bb[2]-ba[2])*mix,ba[3]+(bb[3]-ba[3])*mix,qa.toArray(),ba[8]-1);}
    while(eventCursor<tape.events.length&&tape.events[eventCursor].Time<=elapsed){const e=tape.events[eventCursor++];if(elapsed-e.Time<.15)playTable(e.Kind==='pocket'?'pool-pocket':'pool-impact',Math.min(1,e.Speed/4));}
    dirty=true;if(elapsed>=tape.duration){tape=null;notify(false);}
   }
   if(!tape&&!loading)for(const b of p.balls)pose(b.id,...b.position,b.rotation,b.pocket);
   const aiming=!playing&&!p.settled&&p.turn===0&&!p.ball_in_hand&&!p.break_choices.length;
   guide.visible=cue.visible=aiming;
   if(aiming&&(dirty||lastAngle!==live.current.angle)){lastAngle=live.current.angle;const ball=p.balls.find(b=>b.id===0)!;const start=new THREE.Vector3(ball.position[0]-W/2,H+R,L/2-ball.position[1]);const direction=new THREE.Vector3(Math.cos(lastAngle),0,-Math.sin(lastAngle));guide.geometry.dispose();guide.geometry=new THREE.BufferGeometry().setFromPoints([start,start.clone().addScaledVector(direction,.7)]);guide.computeLineDistances();cue.position.copy(start).addScaledVector(direction,-.79);cue.quaternion.setFromUnitVectors(new THREE.Vector3(0,1,0),direction);dirty=true;}
   if(dirty){renderer.render(scene,camera);dirty=false;}
  };frame=requestAnimationFrame(tick);
  return ()=>{dead=true;generation++;cancelAnimationFrame(frame);observer.disconnect();canvas.removeEventListener('pool-skip',skip);view.dispose();disposeCityResources([scene]);renderer.dispose();canvas.remove();};
 },[]);
 return <div className="pool-render" ref={host}>{error&&<p role="alert">{error}</p>}{playback&&<button className="pool-skip" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('pool-skip'))}>Skip ball motion</button>}<button className="table-camera-reset" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('table-reset'))}>Reset camera</button></div>;
}
