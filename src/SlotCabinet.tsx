import {useEffect,useRef,useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {drumFaces,drumRun,drumRunIDs,type Reel} from './cards';
import {reelArt} from './reels';
import {disposeCityResources} from './city3dResources';
import './slot3d.css';
type Props={strip:Reel[];line:string[];pulled:boolean;rolling:boolean[];turn:number;paid:number;disabled:boolean;onPull:()=>void};

export function SlotCabinet(props:Props){
 const host=useRef<HTMLDivElement>(null),latest=useRef(props);latest.current=props;
 const [status,setStatus]=useState('Opening the machine…');
 useEffect(()=>{
  const el=host.current!;let dead=false,frame=0,dirty=true,model:THREE.Group|undefined;
  const scene=new THREE.Scene();scene.background=new THREE.Color(0x173128);
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));
  renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;renderer.toneMappingExposure=1.25;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Lucky Bell 3D slot machine. Click the side lever to pull; the Pull button below also works.');el.append(canvas);
  const camera=new THREE.PerspectiveCamera(34,1,.01,20);camera.position.set(1.65,1.6,3.7);camera.lookAt(0,.85,0);
  scene.add(new THREE.HemisphereLight(0xffecc8,0x203027,3));
  const key=new THREE.DirectionalLight(0xffdeb0,4);key.position.set(-2,4,4);scene.add(key);
  const fill=new THREE.DirectionalLight(0xd5e6ff,2);fill.position.set(3,2,-2);scene.add(fill);
  const resize=()=>{renderer.setSize(el.clientWidth,el.clientHeight);camera.aspect=el.clientWidth/el.clientHeight;camera.updateProjectionMatrix();dirty=true;};
  const observer=new ResizeObserver(resize);observer.observe(el);resize();
  const papers=[0,1,2].map(()=>{const surface=document.createElement('canvas');surface.width=192;surface.height=384;const texture=new THREE.CanvasTexture(surface);texture.colorSpace=THREE.SRGBColorSpace;texture.flipY=false;return {surface,texture,context:surface.getContext('2d')!};});
  const symbols=new Map<string,HTMLImageElement>();let symbolKey='';
  const loadSymbols=(strip:Reel[])=>{
   const key=strip.map(s=>s.id).join(':');if(key===symbolKey)return;symbolKey=key;
   for(const s of strip){if(symbols.has(s.id))continue;const art=reelArt(s.id);if(!art)continue;
    const img=new Image();img.onload=()=>{if(!dead){symbols.set(s.id,img);dirty=true;}};img.src='data:image/svg+xml;charset=utf-8,'+encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50">${art}</svg>`);
   }
  };
  new GLTFLoader().loadAsync('/art/models/slot-cabinet.glb').then(g=>{
   if(dead){disposeCityResources([g.scene]);return;}model=g.scene;
   model.traverse(o=>{if(!(o instanceof THREE.Mesh))return;const materials=Array.isArray(o.material)?o.material:[o.material];for(const mat of materials){const index=/^slot-reel-(\d)$/.exec(mat.name);if(index&&mat instanceof THREE.MeshStandardMaterial){mat.map=papers[Number(index[1])].texture;mat.color.set(0xffffff);mat.roughness=.8;mat.needsUpdate=true;}}});
   scene.add(model);dirty=true;setStatus('');
  }).catch(()=>{if(!dead)setStatus('The cabinet could not load. Use the pull control below; the result is still reported in words.');});
  const ray=new THREE.Raycaster(),point=new THREE.Vector2();
  const click=(event:MouseEvent)=>{if(latest.current.disabled||!model)return;const rect=canvas.getBoundingClientRect();point.set((event.clientX-rect.left)/rect.width*2-1,1-(event.clientY-rect.top)/rect.height*2);ray.setFromCamera(point,camera);const lever=model.getObjectByName('slot-lever');if(lever&&ray.intersectObject(lever,true).length)latest.current.onPull();};canvas.addEventListener('click',click);
  let seen='',start=0,rendered=0;let wasRolling=false;
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current;loadSymbols(p.strip);
   const rolling=p.rolling.some(Boolean),key=JSON.stringify([p.line,p.pulled,p.turn,p.rolling,p.strip,p.paid]);
   if(rolling&&!wasRolling)start=now;wasRolling=rolling;
   if(key!==seen){seen=key;dirty=true;}
   if(rolling)dirty=true;
   if(!dirty||document.hidden)return;
   const windows=drumFaces(p.strip,p.line,p.pulled);
   for(let i=0;i<3;i++){
    const {context:c,surface,texture}=papers[i],faces=drumRun(p.strip,windows[i],i,14),ids=drumRunIDs(p.strip,windows[i],i,14);
    const progress=p.rolling[i]?Math.min(1,(now-start)/(700+i*450)):1;
    const offset=14*Math.pow(1-progress,3),cell=surface.height/3;
    c.fillStyle='#f4e6bf';c.fillRect(0,0,surface.width,surface.height);
    for(let at=Math.floor(offset);at<=Math.ceil(offset)+3;at++){
     const y=(at-offset)*cell,img=symbols.get(ids[at]);
     if(img)c.drawImage(img,32,y+8,128,112);
     else {c.fillStyle='#342313';c.font='bold 28px Georgia';c.textAlign='center';c.fillText(faces[at]??'—',96,y+75);}
    }
    texture.needsUpdate=true;
   }
   const lever=model?.getObjectByName('slot-lever');if(lever){const age=(now-start)/1000;lever.rotation.x=rolling?.75*Math.sin(Math.min(1,age/.45)*Math.PI):0;}
   const coins=model?.getObjectByName('slot-coins');if(coins)coins.visible=p.paid>0&&!rolling;
   renderer.render(scene,camera);rendered++;canvas.dataset.slot=JSON.stringify({line:p.line,rolling:p.rolling,paid:p.paid,rendered,drawCalls:renderer.info.render.calls,triangles:renderer.info.render.triangles});dirty=false;
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();canvas.removeEventListener('click',click);disposeCityResources([scene]);papers.forEach(p=>p.texture.dispose());renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[]);
 return <div className="slot3d"><div ref={host}/>{status&&<p role="status">{status}</p>}</div>;
}
