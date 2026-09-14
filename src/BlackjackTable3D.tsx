import {blackjackCamera} from './blackjackPresentation';
import type {Person,Presence} from './types';
import {pedestrianModel} from './city3dCast';
import {wardrobe,dressPedestrian} from './city3dWardrobe';
import {poseBlackjackDealer,poseBlackjackPlayer,blackjackPlayerSeat} from './blackjackDealer';
import {cardPose,type planCards} from './blackjackPresentation';
import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {type Card} from './cards';
import {cardTexture} from './cardTexture';
import {TableCamera} from './tableCamera';
import {disposeCityResources} from './city3dResources';
import './blackjackTable3d.css';


export function BlackjackTable3D({mine,theirs,hidden,presentation,dealer,player}:{player?:Pick<Person,"name"|"face"|"alive">;dealer?:Presence;mine:Card[];theirs:Card[];hidden:number;presentation:{plan:ReturnType<typeof planCards>;start:number;active:boolean}}){
 const host=useRef<HTMLDivElement>(null),latest=useRef({mine,theirs,hidden,presentation});latest.current={mine,theirs,hidden,presentation};
 const [status,setStatus]=useState('Opening the card table…'),[wide,setWide]=useState(false);
 useEffect(()=>{
  const el=host.current!;let dead=false,frame=0,dirty=true,key='',stateKey='',rendered=0,cardBuilds=0,previousFrame=0,wasPresenting=false;
  let intervals:number[]=[],playback:{frames:number;fps:number;p95FrameMs:number;worstFrameMs:number}|undefined;
  const scene=new THREE.Scene();scene.background=new THREE.Color('#173128');
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','3D blackjack table with your cards nearest you and the dealer opposite.');el.append(canvas);
  const camera=new THREE.PerspectiveCamera(38,1,.01,20);
  const view=new TableCamera(camera,canvas,()=>{dirty=true;},new THREE.Vector3(0,.87,0),wide?blackjackCamera.wide:blackjackCamera.close,0,blackjackCamera.fitAspect);
  const framing=(event:Event)=>view.frameView(new THREE.Vector3(0,.87,0),(event as CustomEvent<boolean>).detail?blackjackCamera.wide:blackjackCamera.close);
  canvas.addEventListener('blackjack-frame',framing);
  scene.add(new THREE.HemisphereLight(0xffebcb,0x17271f,1.6));
  const light=new THREE.DirectionalLight(0xffe0b5,2.4);light.position.set(-2,5,2);light.castShadow=true;light.shadow.mapSize.set(2048,2048);Object.assign(light.shadow.camera,{left:-3.5,right:3.5,top:3.5,bottom:-3.5,near:.1,far:12});light.shadow.bias=-.0001;light.shadow.normalBias=.008;scene.add(light);
  const resize=()=>{const w=el.clientWidth,h=Math.max(1,el.clientHeight);renderer.setSize(w,h);camera.aspect=w/h;camera.updateProjectionMatrix();view.resize();dirty=true;};
  const observer=new ResizeObserver(resize);observer.observe(el);resize();
  const models:THREE.Group[]=[];let prototype:THREE.Group|undefined;
  const cards=new THREE.Group();scene.add(cards);
  const textures:THREE.Texture[]=[],materials:THREE.Material[]=[],costumes:THREE.Material[]=[];
  const shadow=(root:THREE.Object3D,cast=true)=>root.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=cast;o.receiveShadow=true;}});
  const loader=new GLTFLoader();
  loader.loadAsync('/art/models/gaming-floor.glb').then(g=>{if(dead){disposeCityResources([g.scene]);return;}models.push(g.scene);shadow(g.scene,false);scene.add(g.scene);dirty=true;}).catch(()=>{});
  if(dealer){const name=pedestrianModel(dealer.id,dealer.face);loader.loadAsync(`/art/models/${name}.glb`).then(g=>{
    if(dead){disposeCityResources([g.scene]);return;}models.push(g.scene);
    const actor=g.scene.clone(true);costumes.push(...dressPedestrian(actor,name,wardrobe(dealer.id,dealer.face)));
    poseBlackjackDealer(actor);shadow(actor);scene.add(actor);dirty=true;
  }).catch(()=>{ /* The named dealer and authoritative hand remain readable if the model is unavailable. */ });}
  if(player?.alive){const name=pedestrianModel(player.name,player.face,true);
   Promise.all([name,'gaming-chair'].map(async n=>{const g=await loader.loadAsync(`/art/models/${n}.glb`);if(dead){disposeCityResources([g.scene]);return;}models.push(g.scene);return g.scene;})).then(([rig,chair])=>{
    if(dead||!rig||!chair)return;const actor=rig.clone(true);costumes.push(...dressPedestrian(actor,name,wardrobe(player.name,player.face,true)));poseBlackjackPlayer(actor);
    chair.position.set(blackjackPlayerSeat.x,0,blackjackPlayerSeat.z);chair.rotation.y=blackjackPlayerSeat.yaw;shadow(actor);shadow(chair);scene.add(actor,chair);dirty=true;
   }).catch(()=>{});
  }
  Promise.all(['blackjack-table','playing-card'].map(async name=>{const g=await loader.loadAsync(`/art/models/${name}.glb`);if(dead){disposeCityResources([g.scene]);return;}models.push(g.scene);return g.scene;})).then(([table,card])=>{if(dead||!table||!card)return;shadow(table);scene.add(table);prototype=card;key='';dirty=true;setStatus('');}).catch(()=>{if(!dead)setStatus('The table could not load. Your cards are listed below.');});
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current;
   const state=JSON.stringify([p.presentation.start,p.presentation.active]);if(state!==stateKey){stateKey=state;dirty=true;}
   if(p.presentation.active&&!wasPresenting){intervals=[];previousFrame=0;playback=undefined;}
   if(p.presentation.active&&!document.hidden){if(previousFrame)intervals.push(now-previousFrame);previousFrame=now;}else previousFrame=0;
   if(wasPresenting&&!p.presentation.active&&intervals.length){const sorted=[...intervals].sort((a,b)=>a-b);playback={frames:intervals.length+1,fps:1000*intervals.length/intervals.reduce((sum,n)=>sum+n,0),p95FrameMs:sorted[Math.ceil(sorted.length*.95)-1],worstFrameMs:sorted[sorted.length-1]};}
   wasPresenting=p.presentation.active;
   // Geometry and textures depend on card identity, not animation timing or its completion.
   const k=JSON.stringify(p.presentation.plan.cards.map(c=>c.card??null));
   if(prototype&&k!==key){
    key=k;cardBuilds++;cards.clear();textures.splice(0).forEach(t=>t.dispose());materials.splice(0).forEach(m=>m.dispose());
    p.presentation.plan.cards.forEach(move=>{
      const object=prototype!.clone(true),texture=cardTexture(move.card);textures.push(texture);
      object.traverse(o=>{if(o instanceof THREE.Mesh&&o.material.name==='card printed face'){const m=o.material.clone();m.map=texture;m.color.set('#ffffff');m.needsUpdate=true;o.material=m;materials.push(m);}});
      // A physical reverse lets the card turn over without exposing its face through the stock.
      let face:THREE.Mesh|undefined;
      object.traverse(o=>{if(o instanceof THREE.Mesh&&o.material.name==='card printed face')face=o;});
      if(face){const reverse=face.clone();const back=cardTexture();textures.push(back);const material=(face.material as THREE.MeshStandardMaterial).clone();material.map=back;materials.push(material);reverse.material=material;reverse.rotation.z=Math.PI;object.add(reverse);}
      object.scale.set(1.4,1,1.4);shadow(object);cards.add(object);
    });dirty=true;
   }
   if(p.presentation.active)dirty=true;
   if(!dirty||document.hidden)return;
   const elapsed=p.presentation.active?performance.now()-p.presentation.start:Infinity;
   cards.children.forEach((object,i)=>{const move=p.presentation.plan.cards[i];if(!move)return;const pose=cardPose(move,elapsed);object.visible=pose.visible;object.position.fromArray(pose.position);object.rotation.z=pose.rotation;});
   renderer.render(scene,camera);rendered++;canvas.dataset.blackjack=JSON.stringify({mine:p.mine,theirs:p.theirs,hidden:p.hidden,dealing:p.presentation.active,flips:p.presentation.plan.cards.filter(c=>c.flip).length,dealer:dealer?.id,player:player?.alive?player.name:undefined,rendered,cardBuilds,playback,camera:camera.position.toArray(),drawCalls:renderer.info.render.calls});dirty=false;
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();canvas.removeEventListener('blackjack-frame',framing);view.dispose();disposeCityResources([scene,...models],{textures,materials:[...materials,...costumes]});renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[dealer?.id,dealer?.face,player?.name,player?.face,player?.alive]);
 return <div className="blackjack3d"><div ref={host}/><button className="table-camera-framing" aria-pressed={wide} onClick={()=>{const next=!wide;setWide(next);host.current?.querySelector('canvas')?.dispatchEvent(new CustomEvent('blackjack-frame',{detail:next}));}}>{wide?'Read cards':'Whole table'}</button><button className="table-camera-reset" onClick={()=>host.current?.querySelector("canvas")?.dispatchEvent(new Event("table-reset"))}>Reset view</button><small className="table-camera-help">Drag to orbit · Right-drag to pan · Wheel to zoom · Focus table · Hold WASD / arrows to pan · Q/E to orbit</small>{status&&<p role="status">{status}</p>}</div>;
}
