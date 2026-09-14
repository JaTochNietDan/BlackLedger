import {createBilliardsCue} from './billiardsCue';
import {ballTexture} from './billiardsBallTexture';
import {poolRails,poolCushionGeometry} from './billiardsTableGeometry';
import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {TableCamera} from './tableCamera';
import {disposeCityResources} from './city3dResources';
import {decodePoolReplay,poolFramePair,poolPocketCenters,poolAimAngle,poolPlacementHint,PoolTap,poolCueStroke,POOL_CUE_END} from './billiards';
import type {PoolState,PoolReplay} from './billiards';
import {playTable} from './sound';

interface TableProps {
 pool:PoolState;angle:number;top:number;side:number;motion:boolean;locked:boolean;calledBall:number;calledPocket:number;placement:[number,number];
 onPlaying:(playing:boolean)=>void;onAim:(angle:number)=>void;onBall:(ball:number)=>void;onPocket:(pocket:number)=>void;onPlace:(x:number,y:number)=>void;
}
export function BilliardsTable3D(props:TableProps){
 const {pool}=props;
 const host=useRef<HTMLDivElement>(null),live=useRef(props);live.current=props;
 const [error,setError]=useState(''),[playback,setPlayback]=useState(false);
 useEffect(()=>{
  const el=host.current!;let dead=false,renderer:THREE.WebGLRenderer;
  try{renderer=new THREE.WebGLRenderer({antialias:true});}catch{setError('The 3D table could not start.');return;}
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Billiards table. Click cloth to aim or preview cue placement, a ball to call it, or a numbered pocket to call it. Drag to orbit; right drag to pan; wheel to zoom; Home to reset.');el.append(canvas);
  renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  const scene=new THREE.Scene();scene.background=new THREE.Color('#15231f');scene.add(new THREE.HemisphereLight(0xffedcb,0x192c26,2));
  const light=new THREE.DirectionalLight(0xffe0ad,3);light.position.set(-1,5,1);light.castShadow=true;light.shadow.mapSize.set(1024,1024);Object.assign(light.shadow.camera,{left:-2,right:2,top:3,bottom:-3});light.shadow.bias=-.00015;scene.add(light);
  const W=pool.width,L=pool.length,R=pool.radius,H=.78;
  const camera=new THREE.PerspectiveCamera(38,1,.01,30);let dirty=true;
  const view=new TableCamera(camera,canvas,()=>{dirty=true;},new THREE.Vector3(0,H,0),3.6,Math.PI/2);
  const material=(colour:string)=>new THREE.MeshStandardMaterial({color:colour,roughness:.65});
  const wood=material('#39251b'),felt=material('#285b49'),cushion=material('#204a39');
  const grain=document.createElement('canvas');grain.width=512;grain.height=128;const gc=grain.getContext('2d')!;gc.fillStyle='#856346';gc.fillRect(0,0,512,128);
  for(let i=0;i<90;i++){gc.strokeStyle=i%3===0?'#46332055':'#c39b6030';gc.lineWidth=i%4===0?1.5:.5;gc.beginPath();for(let x=0;x<=512;x+=4){const y=(i*13.73)%128+Math.sin(x*.015+i*.63)*2.5+Math.sin(x*.045+i)*.6;x===0?gc.moveTo(x,y):gc.lineTo(x,y);}gc.stroke();}
  const grainMap=new THREE.CanvasTexture(grain);grainMap.colorSpace=THREE.SRGBColorSpace;grainMap.wrapS=grainMap.wrapT=THREE.RepeatWrapping;wood.map=grainMap;wood.color.set('#8c674a');wood.roughness=.32;

  function box(w:number,h:number,d:number,x:number,y:number,z:number,mat:THREE.Material){const m=new THREE.Mesh(new THREE.BoxGeometry(w,h,d),mat);m.position.set(x,y,z);m.castShadow=true;m.receiveShadow=true;scene.add(m);return m;}
  box(12,.1,12,0,-.06,0,material('#252821'));
  for(const x of [-W/2-.10,W/2+.10])for(const z of [-L/4,L/4])box(.18,.28,L/2-.18,x,H-.19,z,wood);
  for(const z of [-L/2-.10,L/2+.10])box(W-.18,.28,.18,0,H-.19,z,wood);
  for(const x of [-W*.36,W*.36])for(const z of [-L*.38,L*.38])box(.16,.54,.16,x,.27,z,wood);
  for(const x of [-W/2-.12,W/2+.12])for(const z of [-L/4,L/4])box(.10,.045,L/2-.19,x,H+.041,z,wood);
  for(const z of [-L/2-.12,L/2+.12])box(W-.17,.045,.10,0,H+.041,z,wood);
  const pockets=poolPocketCenters(W,L);
  // The bed extends beneath the cushions; holes are actual mesh apertures.
  const shape=new THREE.Shape();shape.moveTo(-.13,-.13);shape.lineTo(W+.13,-.13);shape.lineTo(W+.13,L+.13);shape.lineTo(-.13,L+.13);shape.closePath();
  for(const [x,y] of pockets){const hole=new THREE.Path();hole.absarc(x,y,.076,0,Math.PI*2,true);shape.holes.push(hole);
   const cup=new THREE.Mesh(new THREE.CylinderGeometry(.076,.063,.15,32,1,true),material('#171411'));cup.material.side=THREE.DoubleSide;cup.position.set(x-W/2,H-.075,L/2-y);scene.add(cup);
   const bottom=new THREE.Mesh(new THREE.CircleGeometry(.066,32),material('#080b09'));bottom.rotation.x=-Math.PI/2;bottom.position.set(x-W/2,H-.145,L/2-y);scene.add(bottom);
   const lip=new THREE.Mesh(new THREE.TorusGeometry(.077,.006,8,40),material('#352b21'));lip.rotation.x=Math.PI/2;lip.position.set(x-W/2,H+.002,L/2-y);scene.add(lip);

  }
  const bed=new THREE.Mesh(new THREE.ShapeGeometry(shape,48),felt);bed.rotation.x=-Math.PI/2;bed.position.set(-W/2,H,L/2);bed.receiveShadow=true;scene.add(bed);
  for(const segment of poolRails(W,L)){const rail=new THREE.Mesh(poolCushionGeometry(segment,W,L,R,H),cushion);rail.castShadow=true;rail.receiveShadow=true;scene.add(rail);}
  for(const x of [-W/2-.12,W/2+.12])for(let i=1;i<8;i++){if(i===4)continue;const diamond=new THREE.Mesh(new THREE.CircleGeometry(.008,4),material('#e0d0a6'));diamond.rotation.x=-Math.PI/2;diamond.position.set(x,H+.064,L/2-L*i/8);scene.add(diamond);}
  const balls=new Map<number,THREE.Mesh>();const geometry=new THREE.SphereGeometry(R,32,24);
  for(let id=0;id<16;id++){const mesh=new THREE.Mesh(geometry,new THREE.MeshStandardMaterial({map:ballTexture(id),roughness:.22,metalness:0}));mesh.castShadow=true;mesh.receiveShadow=true;mesh.userData.poolBall=id;balls.set(id,mesh);scene.add(mesh);}
  const ghost=new THREE.Mesh(geometry,new THREE.MeshStandardMaterial({color:'#f5edd6',transparent:true,opacity:.65,depthWrite:false}));scene.add(ghost);
  const selection=new THREE.Mesh(new THREE.RingGeometry(R*1.25,R*1.5,48),new THREE.MeshBasicMaterial({color:'#f8d88c',side:THREE.DoubleSide}));selection.rotation.x=-Math.PI/2;scene.add(selection);
  const pocketMarkers=pockets.map(([x,y],i)=>{
   const group=new THREE.Group();group.position.set(x-W/2,H+.006,L/2-y);
   const ring=new THREE.Mesh(new THREE.RingGeometry(.078,.091,48),new THREE.MeshBasicMaterial({color:'#dcc38b',transparent:true,opacity:.5,side:THREE.DoubleSide}));ring.rotation.x=-Math.PI/2;group.add(ring);
   const label=document.createElement('canvas');label.width=128;label.height=128;const c=label.getContext('2d')!;c.fillStyle='#e9d6a8';c.font='bold 76px Georgia';c.textAlign='center';c.textBaseline='middle';c.fillText(String(i+1),64,64);
   const texture=new THREE.CanvasTexture(label);texture.colorSpace=THREE.SRGBColorSpace;const sprite=new THREE.Sprite(new THREE.SpriteMaterial({map:texture,depthWrite:false}));sprite.scale.set(.095,.095,1);sprite.position.y=.04;group.add(sprite);scene.add(group);return {group,ring};
  });
  const headLine=new THREE.Line(new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(-W/2,H+.003,L/4),new THREE.Vector3(W/2,H+.003,L/4)]),new THREE.LineDashedMaterial({color:'#e8d8b6',dashSize:.04,gapSize:.03,transparent:true,opacity:.65}));headLine.computeLineDistances();scene.add(headLine);
  const guide=new THREE.Line(new THREE.BufferGeometry(),new THREE.LineDashedMaterial({color:'#e7d6a5',dashSize:.04,gapSize:.025,transparent:true,opacity:.65}));scene.add(guide);
  const {cue,materials:cueMaterials}=createBilliardsCue();scene.add(cue);
  const cuePose=(origin:THREE.Vector3,angle:number,front:number,top:number,side:number,opacity:number)=>{const direction=new THREE.Vector3(Math.cos(angle),0,-Math.sin(angle)),across=new THREE.Vector3(-Math.sin(angle),0,-Math.cos(angle));cue.position.copy(origin).addScaledVector(direction,front-.714).addScaledVector(across,side);cue.position.y+=top;cue.quaternion.setFromUnitVectors(new THREE.Vector3(0,1,0),direction);for(const m of cueMaterials){m.opacity=opacity;m.depthWrite=opacity>.99;}};
  let tape:PoolReplay|null=null,started=0,loading=false,key='',generation=0,playing=false,eventCursor=0,cueSound=false;
  const preference=matchMedia('(prefers-reduced-motion: reduce)');
  const notify=(v:boolean)=>{if(v!==playing){playing=v;setPlayback(v);live.current.onPlaying(v);}};
  const skip=()=>{generation++;tape=null;loading=false;notify(false);dirty=true;};canvas.addEventListener('pool-skip',skip);
  const ray=new THREE.Raycaster(),pointer=new THREE.Vector2(),hit=new THREE.Vector3(),plane=new THREE.Plane(new THREE.Vector3(0,1,0),-H),tap=new PoolTap();
  const down=(e:PointerEvent)=>tap.begin(e.pointerId,e.clientX,e.clientY,e.isPrimary,e.button);
  const move=(e:PointerEvent)=>tap.move(e.pointerId,e.clientX,e.clientY);
  const cancel=()=>tap.cancel();
  const up=(e:PointerEvent)=>{
   if(!tap.end(e.pointerId,e.clientX,e.clientY))return;
   const v=live.current,p=v.pool;
   if(v.locked||playing||loading||p.unavailable||p.settled||p.turn!==0||p.break_choices.length)return;
   const bounds=canvas.getBoundingClientRect();pointer.set((e.clientX-bounds.left)/bounds.width*2-1,-(e.clientY-bounds.top)/bounds.height*2+1);ray.setFromCamera(pointer,camera);
   if(!ray.ray.intersectPlane(plane,hit))return;
   const point:[number,number]=[hit.x+W/2,L/2-hit.z];
   if(p.ball_in_hand){v.onPlace(...point);return;}
   const pocket=pockets.findIndex(([x,y])=>Math.hypot(point[0]-x,point[1]-y)<.11);
   if(pocket>=0&&!p.breaking){v.onPocket(pocket);return;}
   const picked=ray.intersectObjects([...balls.values()].filter(b=>b.visible),false)[0];
   if(picked){const id=picked.object.userData.poolBall as number;const b=p.balls.find(b=>b.id===id)!;point[0]=b.position[0];point[1]=b.position[1];if(p.legal_balls.includes(id)&&!p.breaking)v.onBall(id);}
   if(point[0]<0||point[0]>W||point[1]<0||point[1]>L)return;
   const cue=p.balls.find(b=>b.id===0)!;const angle=poolAimAngle([cue.position[0],cue.position[1]],point);if(angle!==null)v.onAim(angle);
  };
  canvas.addEventListener('pointerdown',down);canvas.addEventListener('pointermove',move);canvas.addEventListener('pointerup',up);canvas.addEventListener('pointercancel',cancel);canvas.addEventListener('lostpointercapture',cancel);
  const worldRotation=new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(1,0,0),-Math.PI/2),qa=new THREE.Quaternion(),qb=new THREE.Quaternion();
  function pose(id:number,x:number,y:number,z:number,q:number[],pocket:number){const mesh=balls.get(id)!;mesh.visible=pocket<0;mesh.position.set(x-W/2,H+z,L/2-y);qa.fromArray(q);mesh.quaternion.copy(worldRotation).multiply(qa);}
  const resize=()=>{renderer.setSize(el.clientWidth,Math.max(1,el.clientHeight));camera.aspect=el.clientWidth/Math.max(1,el.clientHeight);camera.updateProjectionMatrix();view.resize();dirty=true;};const observer=new ResizeObserver(resize);observer.observe(el);resize();
  let frame=0,lastAngle=NaN,lastPool:PoolState|null=null,lastSelection='';
  const tick=(now:number)=>{
   frame=requestAnimationFrame(tick);const p=live.current.pool;if(lastPool!==p){lastPool=p;dirty=true;}
   const selectionKey=JSON.stringify([live.current.placement,live.current.calledBall,live.current.calledPocket,live.current.locked,live.current.top,live.current.side]);if(lastSelection!==selectionKey){lastSelection=selectionKey;dirty=true;}
   const next=`${p.shots}:${p.replay}`;
   if(next!==key){const first=key==='';key=next;const token=++generation;tape=null;eventCursor=0;cueSound=false;dirty=true;
    if(!first&&p.replay&&live.current.motion&&!preference.matches){loading=true;notify(true);decodePoolReplay(p.replay).then(r=>{if(dead||token!==generation)return;tape=r;started=performance.now();loading=false;dirty=true;}).catch(()=>{if(dead||token!==generation)return;loading=false;notify(false);setError('Replay unavailable. The saved table below is the authoritative result.');});}
    else {loading=false;notify(false);}
   }
   if((!live.current.motion||preference.matches)&&(tape||loading)){tape=null;loading=false;notify(false);generation++;dirty=true;}
   let stroking=false;
   if(tape){const seconds=(now-started)/1000,intent=p.stroke?.intent,animated=!!intent;
    const stroke=poolCueStroke(seconds,intent?.speed??0,R,intent?.top??0,intent?.side??0),elapsed=animated?stroke.ballTime:seconds;
    const initial=tape.frames[0].balls.find(b=>b[0]===0)!;
    if(animated&&stroke.visible){stroking=true;cuePose(new THREE.Vector3(initial[1]-W/2,H+initial[3],L/2-initial[2]),intent.angle??0,stroke.front,intent.top??0,intent.side??0,stroke.opacity);}
    if(animated&&stroke.contact&&!cueSound){cueSound=true;if(seconds<.87)playTable('pool-impact',Math.min(1,(intent.speed??0)/8));}
    const {a,b,mix}=poolFramePair(tape,elapsed);
    const nextBalls=new Map(b.balls.map(ball=>[ball[0],ball]));
    for(const ba of a.balls){const bb=nextBalls.get(ba[0])!;qa.fromArray(ba.slice(4,8));qb.fromArray(bb.slice(4,8));qa.slerp(qb,mix);pose(ba[0],ba[1]+(bb[1]-ba[1])*mix,ba[2]+(bb[2]-ba[2])*mix,ba[3]+(bb[3]-ba[3])*mix,qa.toArray(),ba[8]-1);}
    while((!animated||stroke.contact)&&eventCursor<tape.events.length&&tape.events[eventCursor].Time<=elapsed){const e=tape.events[eventCursor++];if(elapsed-e.Time<.15)playTable(e.Kind==='pocket'?'pool-pocket':'pool-impact',Math.min(1,e.Speed/4));}
    dirty=true;if(elapsed>=tape.duration&&(!animated||seconds>=POOL_CUE_END)){tape=null;notify(false);stroking=false;}
   }
   if(!tape&&!loading)for(const b of p.balls)pose(b.id,...b.position,b.rotation,b.pocket);
   const interactive=!playing&&!live.current.locked&&!p.unavailable&&!p.settled&&p.turn===0&&!p.break_choices.length;
   ghost.visible=interactive&&p.ball_in_hand;headLine.visible=ghost.visible&&p.behind_head_string;
   if(ghost.visible){const [x,y]=live.current.placement;ghost.position.set(x-W/2,H+R,L/2-y);ghost.material.color.set(poolPlacementHint(p,x,y)?'#c3553f':'#f5edd6');balls.get(0)!.visible=false;}
   const called=p.balls.find(b=>b.id===live.current.calledBall);selection.visible=interactive&&!p.breaking&&!p.ball_in_hand&&!!called&&called.pocket<0;if(called)selection.position.set(called.position[0]-W/2,H+.003,L/2-called.position[1]);
   for(let i=0;i<pocketMarkers.length;i++){const marker=pocketMarkers[i];marker.group.visible=interactive&&!p.breaking&&!p.ball_in_hand;marker.ring.material.opacity=i===live.current.calledPocket?1:.35;marker.ring.material.color.set(i===live.current.calledPocket?'#ffd071':'#dcc38b');}
   const aiming=interactive&&!p.settled&&p.turn===0&&!p.ball_in_hand&&!p.break_choices.length;
   guide.visible=aiming;cue.visible=aiming||stroking;
   if(aiming&&(dirty||lastAngle!==live.current.angle)){lastAngle=live.current.angle;const ball=p.balls.find(b=>b.id===0)!;const start=new THREE.Vector3(ball.position[0]-W/2,H+R,L/2-ball.position[1]);const direction=new THREE.Vector3(Math.cos(lastAngle),0,-Math.sin(lastAngle));guide.geometry.dispose();guide.geometry=new THREE.BufferGeometry().setFromPoints([start,start.clone().addScaledVector(direction,.7)]);guide.computeLineDistances();cuePose(start,lastAngle,-R-.06,live.current.top,live.current.side,1);dirty=true;}
   if(dirty){renderer.render(scene,camera);dirty=false;}
  };frame=requestAnimationFrame(tick);
  return ()=>{dead=true;generation++;cancelAnimationFrame(frame);observer.disconnect();canvas.removeEventListener('pool-skip',skip);canvas.removeEventListener('pointerdown',down);canvas.removeEventListener('pointermove',move);canvas.removeEventListener('pointerup',up);canvas.removeEventListener('pointercancel',cancel);canvas.removeEventListener('lostpointercapture',cancel);view.dispose();disposeCityResources([scene]);renderer.dispose();canvas.remove();};
 },[]);
 return <div className="pool-render" ref={host}>{error&&<p role="alert">{error}</p>}{playback&&<button className="pool-skip" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('pool-skip'))}>Skip ball motion</button>}<button className="table-camera-reset" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('table-reset'))}>Reset camera</button></div>;
}
