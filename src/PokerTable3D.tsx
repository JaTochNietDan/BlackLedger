import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {cardTexture} from './cardTexture';
import {TableCamera} from './tableCamera';
import {disposeCityResources} from './city3dResources';
import type {CardsState} from './Tables';
import type {Card} from './cards';
import './cardRoom.css';

// The projection is the only source of exposed cards, names and chip amounts.
export function PokerTable3D({cards}:{cards:CardsState|null}){
 const host=useRef<HTMLDivElement>(null),latest=useRef(cards);latest.current=cards;
 const [failed,setFailed]=useState(false);
 useEffect(()=>{
  const el=host.current!;let renderer:THREE.WebGLRenderer;
  try{renderer=new THREE.WebGLRenderer({antialias:true});}catch{setFailed(true);return;}
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFSoftShadowMap;
  renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Poker table. Your cards are nearest you; shared cards are in the center.');el.append(canvas);
  const scene=new THREE.Scene();scene.background=new THREE.Color('#10251e');
  scene.add(new THREE.HemisphereLight(0xffeed0,0x182821,2.5));
  const light=new THREE.DirectionalLight(0xffddaa,3);light.position.set(-2,6,3);light.castShadow=true;light.shadow.mapSize.set(1024,1024);Object.assign(light.shadow.camera,{left:-3,right:3,top:3,bottom:-3});light.shadow.bias=-.0002;scene.add(light);
  const camera=new THREE.PerspectiveCamera(38,1,.01,30);let dirty=true;
  const view=new TableCamera(camera,canvas,()=>{dirty=true;},new THREE.Vector3(0,.83,.12),4.2);
  function oval(radius:number,height:number,y:number,color:string,scaleX:number){
   const mesh=new THREE.Mesh(new THREE.CylinderGeometry(radius,radius,height,96),new THREE.MeshStandardMaterial({color,roughness:.75}));mesh.scale.x=scaleX;mesh.position.y=y;mesh.castShadow=true;mesh.receiveShadow=true;scene.add(mesh);return mesh;
  }
  oval(1.42,.17,.7,'#342217',1.48);oval(1.39,.07,.81,'#815b32',1.48);oval(1.32,.085,.85,'#291e18',1.48);oval(1.22,.025,.889,'#285744',1.52);
  const floor=new THREE.Mesh(new THREE.PlaneGeometry(25,25),new THREE.MeshStandardMaterial({color:'#232921',roughness:1}));floor.rotation.x=-Math.PI/2;scene.add(floor);
  for(const x of [-1.1,1.1]){const leg=new THREE.Mesh(new THREE.CylinderGeometry(.14,.24,.7,16),new THREE.MeshStandardMaterial({color:'#392719'}));leg.position.set(x,.35,0);scene.add(leg);}
  let contents=new THREE.Group();scene.add(contents);let key='';
  const addCard=(card:Card|undefined,x:number,z:number)=>{
   const stock=new THREE.Mesh(new THREE.BoxGeometry(.265,.009,.38),new THREE.MeshStandardMaterial({color:'#d5c9ae'}));stock.position.set(x,.918,z);stock.castShadow=true;contents.add(stock);
   const texture=cardTexture(card);texture.flipY=true;
   const face=new THREE.Mesh(new THREE.PlaneGeometry(.257,.372),new THREE.MeshStandardMaterial({map:texture,roughness:.86}));face.rotation.x=-Math.PI/2;face.position.set(x,.924,z);contents.add(face);
  };
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
  let frame=0;
  const tick=()=>{
   frame=requestAnimationFrame(tick);const p=latest.current,next=JSON.stringify(p);
   if(key!==next){key=next;scene.remove(contents);disposeCityResources([contents]);contents=new THREE.Group();scene.add(contents);
    label('BLACK LEDGER • BACK ROOM',0,-.35,1.65);
    (p?.board??[]).forEach((card,i)=>addCard(card,(i-2)*.31,0));
    (p?.mine??[]).forEach((card,i)=>addCard(card,(i-.5)*.30,.78));
    label(p?'YOUR HAND':'TAKE A SEAT',0,1.07,.7);
    if(p){label(`POT $${p.pot}`,0,.35,.9);chips(p.pot,.95,.43);
     p.seats.forEach((seat,i)=>{
      const x=(i-(p.seats.length-1)/2)*Math.min(1.05,2.65/Math.max(1,p.seats.length-1));
      if(!seat.folded){const exposed=seat.cards;for(let j=0;j<2;j++)addCard(exposed?.[j],x+(j-.5)*.29,-.78);}
      label(`${seat.name}${seat.folded?' · Folded':''}`,x,-1.07,.95);chips(seat.in,x+.36,-.55);
     });
    }dirty=true;
   }
   if(dirty&&!document.hidden){renderer.render(scene,camera);canvas.dataset.poker=JSON.stringify({board:p?.board??[],mine:p?.mine??[],seats:p?.seats.map(s=>({who:s.who,exposed:s.cards?.length??0,folded:s.folded})),camera:camera.position.toArray(),drawCalls:renderer.info.render.calls});dirty=false;}
  };tick();
  return()=>{cancelAnimationFrame(frame);observer.disconnect();view.dispose();disposeCityResources([scene]);renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[]);
 return <div className="poker3d"><div ref={host}/><button className="table-camera-reset" onClick={()=>host.current?.querySelector('canvas')?.dispatchEvent(new Event('table-reset'))}>Reset view</button><small className="table-camera-help">Drag to orbit · Right-drag to pan · Wheel to zoom · Focus table for WASD / arrows</small>{failed&&<p role="status">The 3D table is unavailable. Read the cards below.</p>}</div>;
}
