import type {Person,Presence} from './types';
import {pedestrianModel} from './city3dCast';
import {wardrobe,dressPedestrian} from './city3dWardrobe';
import {poseBlackjackDealer,poseBlackjackPlayer,blackjackPlayerSeat} from './blackjackDealer';
import {cardPose,type planCards} from './blackjackPresentation';
import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {type Card,knownCard,pipOf,isRedSuit} from './cards';
import {disposeCityResources} from './city3dResources';
import './blackjackTable3d.css';

function cardTexture(card?:Card){
 const canvas=document.createElement('canvas');canvas.width=256;canvas.height=368;
 const c=canvas.getContext('2d')!;
 c.fillStyle='#eee6cf';c.fillRect(0,0,256,368);
 if(!card||!knownCard(card)){
  c.fillStyle='#702920';c.fillRect(12,12,232,344);c.strokeStyle='#d1ad77';c.lineWidth=1;
  for(let y=-240;y<600;y+=14){c.beginPath();c.moveTo(18,y);c.lineTo(238,y+220);c.stroke();c.beginPath();c.moveTo(238,y);c.lineTo(18,y+220);c.stroke();}
  c.strokeStyle='#eee6cf';c.lineWidth=5;c.strokeRect(20,20,216,328);
 }else{
  c.fillStyle=isRedSuit(card.suit)?'#a32922':'#1f231f';
  for(let i=0;i<2;i++){c.save();if(i){c.translate(256,368);c.rotate(Math.PI);}c.font='bold 42px Georgia';c.fillText(card.rank,18,48);c.font='40px Georgia';c.fillText(pipOf(card.suit),18,91);c.restore();}
  c.textAlign='center';
  const n=Number(card.rank);
  const rows:Record<number,number[][]>={2:[[0,-1],[0,1]],3:[[0,-1],[0,0],[0,1]],4:[[-1,-1],[1,-1],[-1,1],[1,1]],5:[[-1,-1],[1,-1],[0,0],[-1,1],[1,1]],6:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1]],7:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1],[0,-.5]],8:[[-1,-1],[1,-1],[-1,0],[1,0],[-1,1],[1,1],[0,-.5],[0,.5]],9:[[-1,-1],[1,-1],[-1,-.33],[1,-.33],[-1,.33],[1,.33],[-1,1],[1,1],[0,0]],10:[[-1,-1],[1,-1],[-1,-.33],[1,-.33],[-1,.33],[1,.33],[-1,1],[1,1],[0,-.67],[0,.67]]};
  if(rows[n]){c.font='43px Georgia';for(const [x,y] of rows[n])c.fillText(pipOf(card.suit),128+x*39,195+y*87);}
  else {c.font='bold 76px Georgia';c.fillText(card.rank,128,176);c.font='64px Georgia';c.fillText(pipOf(card.suit),128,246);}
 }
 const texture=new THREE.CanvasTexture(canvas);texture.colorSpace=THREE.SRGBColorSpace;texture.flipY=false;texture.anisotropy=4;return texture;
}

export function BlackjackTable3D({mine,theirs,hidden,presentation,dealer,player}:{player?:Pick<Person,"name"|"face"|"alive">;dealer?:Presence;mine:Card[];theirs:Card[];hidden:number;presentation:{plan:ReturnType<typeof planCards>;start:number;active:boolean}}){
 const host=useRef<HTMLDivElement>(null),latest=useRef({mine,theirs,hidden,presentation});latest.current={mine,theirs,hidden,presentation};
 const [status,setStatus]=useState('Opening the card table…');
 useEffect(()=>{
  const el=host.current!;let dead=false,frame=0,dirty=true,key='',rendered=0;
  const scene=new THREE.Scene();scene.background=new THREE.Color('#173128');
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','3D blackjack table with your cards nearest you and the dealer opposite.');el.append(canvas);
  const camera=new THREE.PerspectiveCamera(38,1,.01,20);camera.position.set(0,player?.alive?4.8:3.8,player?.alive?5.3:3.5);camera.lookAt(0,player?.alive?.8:dealer?1.02:.72,player?.alive?.25:dealer?-.18:0);
  scene.add(new THREE.HemisphereLight(0xffebcb,0x17271f,1.6));
  const light=new THREE.DirectionalLight(0xffe0b5,2.4);light.position.set(-2,5,2);light.castShadow=true;light.shadow.mapSize.set(2048,2048);Object.assign(light.shadow.camera,{left:-3.5,right:3.5,top:3.5,bottom:-3.5,near:.1,far:12});light.shadow.bias=-.0001;light.shadow.normalBias=.008;scene.add(light);
  const resize=()=>{const w=el.clientWidth,h=el.clientHeight;renderer.setSize(w,h);camera.aspect=w/h;camera.position.set(0,player?.alive?4.8:3.8,player?.alive?5.3:3.5).multiplyScalar(Math.max(1,1.35/camera.aspect));camera.lookAt(0,player?.alive?.8:dealer?1.02:.72,player?.alive?.25:dealer?-.18:0);camera.updateProjectionMatrix();dirty=true;};
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
  const tick=()=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current,k=JSON.stringify(p);
   if(prototype&&k!==key){
    key=k;cards.clear();textures.splice(0).forEach(t=>t.dispose());materials.splice(0).forEach(m=>m.dispose());
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
   const elapsed=p.presentation.active?performance.now()-p.presentation.start:Infinity;
   cards.children.forEach((object,i)=>{const move=p.presentation.plan.cards[i];if(!move)return;const pose=cardPose(move,elapsed);object.visible=pose.visible;object.position.fromArray(pose.position);object.rotation.z=pose.rotation;});
   if(p.presentation.active)dirty=true;
   if(!dirty||document.hidden)return;renderer.render(scene,camera);rendered++;canvas.dataset.blackjack=JSON.stringify({mine:p.mine,theirs:p.theirs,hidden:p.hidden,dealing:p.presentation.active,flips:p.presentation.plan.cards.filter(c=>c.flip).length,dealer:dealer?.id,player:player?.alive?player.name:undefined,rendered,drawCalls:renderer.info.render.calls});dirty=false;
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();disposeCityResources([scene,...models],{textures,materials:[...materials,...costumes]});renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[dealer?.id,dealer?.face,player?.name,player?.face,player?.alive]);
 return <div className="blackjack3d"><div ref={host}/>{status&&<p role="status">{status}</p>}</div>;
}
