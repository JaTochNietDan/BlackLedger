import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {poseInteriorOccupant} from './interiorStaging';
import {pedestrianModel} from './city3dCast';
import {dressPedestrian,wardrobe} from './city3dWardrobe';
import {planPoker,pokerCardPose,pokerSeatX} from './pokerPresentation';
import type {Person,Presence} from './types';
import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {cardTexture} from './cardTexture';
import {TableCamera} from './tableCamera';
import {disposeCityResources} from './city3dResources';
import type {CardsState} from './Tables';
import type {Card} from './cards';
import './cardRoom.css';

// The projection is the only source of exposed cards, names and chip amounts.
export function PokerTable3D({cards,player,people=[],motion=true,animateOnMount=false}:{cards:CardsState|null;player?:Pick<Person,"name"|"face"|"alive">;people?:Presence[];motion?:boolean;animateOnMount?:boolean}){
 const host=useRef<HTMLDivElement>(null),latest=useRef(cards);latest.current=cards;
 const castProps=useRef({player,people,motion});castProps.current={player,people,motion};
 const [failed,setFailed]=useState(false),[wide,setWide]=useState(false);
 useEffect(()=>{
  const el=host.current!;let dead=false;let renderer:THREE.WebGLRenderer;
  try{renderer=new THREE.WebGLRenderer({antialias:true});}catch{setFailed(true);return;}
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Poker table. Your cards are nearest you; shared cards are in the center.');el.append(canvas);
  const scene=new THREE.Scene();scene.background=new THREE.Color('#10251e');
  scene.add(new THREE.HemisphereLight(0xffeed0,0x182821,2.5));
  const light=new THREE.DirectionalLight(0xffddaa,3);light.position.set(-2,6,3);light.castShadow=true;light.shadow.mapSize.set(1024,1024);Object.assign(light.shadow.camera,{left:-3,right:3,top:3,bottom:-3});light.shadow.bias=-.0002;scene.add(light);
  const camera=new THREE.PerspectiveCamera(38,1,.01,30);let dirty=true;
  const view=new TableCamera(camera,canvas,()=>{dirty=true;},new THREE.Vector3(0,.924,.20),3.35);
  const framing=(event:Event)=>{const whole=(event as CustomEvent<boolean>).detail;view.frameView(whole?new THREE.Vector3(0,1.08,-.12):new THREE.Vector3(0,.924,.20),whole?4.6:3.35);};
  canvas.addEventListener('poker-frame',framing);
  function oval(radius:number,height:number,y:number,color:string,scaleX:number){
   const mesh=new THREE.Mesh(new THREE.CylinderGeometry(radius,radius,height,96),new THREE.MeshStandardMaterial({color,roughness:.75}));mesh.scale.x=scaleX;mesh.position.y=y;mesh.castShadow=true;mesh.receiveShadow=true;scene.add(mesh);return mesh;
  }
  oval(1.42,.17,.7,'#342217',1.48);oval(1.39,.07,.81,'#815b32',1.48);oval(1.32,.085,.85,'#291e18',1.48);oval(1.22,.025,.889,'#285744',1.52);
  const floor=new THREE.Mesh(new THREE.PlaneGeometry(25,25),new THREE.MeshStandardMaterial({color:'#232921',roughness:1}));floor.rotation.x=-Math.PI/2;scene.add(floor);
  for(const x of [-1.1,1.1]){const leg=new THREE.Mesh(new THREE.CylinderGeometry(.14,.24,.7,16),new THREE.MeshStandardMaterial({color:'#392719'}));leg.position.set(x,.35,0);scene.add(leg);}
  let contents=new THREE.Group();scene.add(contents);let key='';
  let plan=planPoker(null,null,false),started=0,previous:CardsState|null=null,initialized=animateOnMount;
  let cardObjects:THREE.Group[]=[];
  const preference=matchMedia('(prefers-reduced-motion: reduce)');
  const addCard=(card:Card|undefined)=>{
   const object=new THREE.Group();
   const stock=new THREE.Mesh(new THREE.BoxGeometry(.265,.009,.38),new THREE.MeshStandardMaterial({color:'#d5c9ae'}));stock.castShadow=true;object.add(stock);
   for(const back of [false,true]){
    const texture=cardTexture(back?undefined:card);texture.flipY=true;
    const face=new THREE.Mesh(new THREE.PlaneGeometry(.257,.372),new THREE.MeshStandardMaterial({map:texture,roughness:.86}));
    face.rotation.x=back?Math.PI/2:-Math.PI/2;face.position.y=back?-.006:.006;object.add(face);
   }
   contents.add(object);cardObjects.push(object);
  };
  const models=new Map<string,THREE.Group>(),costumes:THREE.Material[]=[];
  let cast=new THREE.Group(),castKey='';scene.add(cast);
  const loader=new GLTFLoader();
  Promise.all(['person','woman','gaming-chair'].map(async name=>{
   const g=await loader.loadAsync(`/art/models/${name}.glb`);
   if(dead){disposeCityResources([g.scene]);return;}models.set(name,g.scene);
  })).then(()=>{dirty=true;}).catch(()=>{});
  const label=(text:string,x:number,z:number,width=1)=>{
   const bitmap=document.createElement('canvas');bitmap.width=768;bitmap.height=96;const c=bitmap.getContext('2d')!;
   c.fillStyle='#edddb8';c.font='32px Georgia';c.textAlign='center';c.fillText(text,384,61,748);
   const texture=new THREE.CanvasTexture(bitmap);texture.colorSpace=THREE.SRGBColorSpace;
   const mesh=new THREE.Mesh(new THREE.PlaneGeometry(width,.13),new THREE.MeshBasicMaterial({map:texture,transparent:true,depthWrite:false}));mesh.rotation.x=-Math.PI/2;mesh.position.set(x,.929,z);contents.add(mesh);
  };
  const chips=(amount:number,x:number,z:number)=>{
   if(amount<=0)return;
   const count=Math.min(12,Math.max(1,Math.ceil(amount/25)));
   for(let i=0;i<count;i++){const chip=new THREE.Mesh(new THREE.CylinderGeometry(.052,.052,.012,24),new THREE.MeshStandardMaterial({color:i%3===0?'#e3cfa4':'#943b2c',roughness:.62}));chip.position.set(x,.928+i*.013,z);chip.castShadow=true;contents.add(chip);}
  };
  const resize=()=>{const w=el.clientWidth,h=Math.max(1,el.clientHeight);renderer.setSize(w,h);camera.aspect=w/h;camera.updateProjectionMatrix();view.resize();};const observer=new ResizeObserver(resize);observer.observe(el);resize();
  let frame=0,wasDealing=false;
  const tick=(now=performance.now())=>{
   frame=requestAnimationFrame(tick);const p=latest.current,next=JSON.stringify(p);
   if(key!==next){key=next;plan=planPoker(previous,p,initialized&&castProps.current.motion&&!preference.matches);previous=p;initialized=true;started=now;cardObjects=[];scene.remove(contents);disposeCityResources([contents]);contents=new THREE.Group();scene.add(contents);
    label('BLACK LEDGER • BACK ROOM',0,-.35,1.65);
    plan.cards.forEach(move=>addCard(move.card));
    label(p?'YOUR HAND':'TAKE A SEAT',0,1.07,.7);
    if(p){label(`POT $${p.pot}`,0,.35,.9);chips(p.pot,.95,.43);
     p.seats.forEach((seat,i)=>{
      const x=pokerSeatX(i,p.seats.length);
      label(`${seat.name}${seat.folded?' · Folded':''}`,x,-1.07,.95);chips(seat.in,x+.36,-.55);
     });
    }dirty=true;
   }
   const roster=JSON.stringify([p?.seats.map(s=>s.who),castProps.current.player,castProps.current.people.map(w=>[w.id,w.face])]);
   if(models.size===3&&castKey!==roster){
    castKey=roster;cast.clear();costumes.splice(0).forEach(m=>m.dispose());
    const occupants=(p?.seats??[]).map((s,i)=>({id:s.who,face:castProps.current.people.find(w=>w.id===s.who)?.face,x:pokerSeatX(i,p!.seats.length),z:-1.66,yaw:0,isPlayer:false}));
    const self=castProps.current.player;
    if(self?.alive)occupants.push({id:self.name,face:self.face,x:0,z:1.66,yaw:Math.PI,isPlayer:true});
    for(const s of occupants){
     const name=pedestrianModel(s.id,s.face,s.isPlayer),actor=models.get(name)!.clone(true),chair=models.get('gaming-chair')!.clone(true);
     costumes.push(...dressPedestrian(actor,name,wardrobe(s.id,s.face,s.isPlayer)));
     poseInteriorOccupant(actor,{id:s.id,x:s.x-Math.sin(s.yaw)*.18,z:s.z-Math.cos(s.yaw)*.18,yaw:s.yaw,seat:.665});
     chair.position.set(s.x,0,s.z);chair.rotation.y=s.yaw;
     actor.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});cast.add(actor,chair);
    }dirty=true;
   }
   const elapsed=!castProps.current.motion||preference.matches?plan.duration:now-started;
   cardObjects.forEach((object,i)=>{const pose=pokerCardPose(plan.cards[i],elapsed);object.visible=pose.visible;object.position.set(pose.x,pose.y,pose.z);object.rotation.z=pose.turn;});
   const dealing=elapsed<plan.duration;
   if(dealing||wasDealing)dirty=true;wasDealing=dealing;
   if(dirty&&!document.hidden){renderer.render(scene,camera);canvas.dataset.poker=JSON.stringify({board:p?.board??[],mine:p?.mine??[],seats:p?.seats.map(s=>({who:s.who,exposed:s.cards?.length??0,folded:s.folded})),camera:camera.position.toArray(),dealing:elapsed<plan.duration,cast:cast.children.length/2,drawCalls:renderer.info.render.calls});dirty=false;}
  };tick();
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();canvas.removeEventListener('poker-frame',framing);view.dispose();costumes.forEach(m=>m.dispose());disposeCityResources([scene,...models.values()]);renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[]);
 return <div className="poker3d"><div ref={host}/><button className="poker-camera-framing" aria-pressed={wide} onClick={()=>{const next=!wide;setWide(next);host.current?.querySelector('canvas')?.dispatchEvent(new CustomEvent('poker-frame',{detail:next}));}}>{wide?'Read cards':'Whole table'}</button><button className="table-camera-reset" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('table-reset'))}>Reset view</button><small className="table-camera-help">Drag to orbit · Right-drag to pan · Wheel to zoom · Focus table · Hold WASD / arrows to pan · Q/E to orbit</small>{failed&&<p role="status">The 3D table is unavailable. Read the cards below.</p>}</div>;
}
