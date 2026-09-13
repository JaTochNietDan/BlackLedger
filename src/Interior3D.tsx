import {cameraCommand, screenPan} from './city3dControls';
import {useEffect, useRef, useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';
import type {Presence} from './types';
import {wardrobe, dressPedestrian} from './city3dWardrobe';
import {pedestrianModel} from './city3dCast';
import {disposeCityResources} from './city3dResources';
import './interior3d.css';

export function Interior3D(props:{people:Presence[];picked:string;onPick:(id:string)=>void;minute:number}) {
 const host=useRef<HTMLDivElement>(null), latest=useRef(props);latest.current=props;
 const [status,setStatus]=useState('Opening Saint Agnes…');
 useEffect(()=>{
  const el=host.current!;let dead=false,frame=0,dirty=true,renderedFrames=0;
  const scene=new THREE.Scene();scene.background=new THREE.Color(0x171b18);
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;renderer.toneMappingExposure=1.1;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label','Saint Agnes 3D interior. Drag or Q/E to orbit, scroll or +/- to zoom, Home resets, WASD or arrows pan; click a person to select their actions.');canvas.tabIndex=0;el.append(canvas);
  const camera=new THREE.OrthographicCamera(-9,9,7,-7,.1,100);camera.position.set(13,14,-17);
  const controls=new OrbitControls(camera,canvas);controls.target.set(0,1,0);controls.minZoom=.7;controls.maxZoom=3;
  controls.minPolarAngle=.35;controls.maxPolarAngle=1.15;controls.enablePan=false;controls.update();
  const changed=()=>{dirty=true;};controls.addEventListener('change',changed);
  const resize=()=>{const w=el.clientWidth,h=el.clientHeight;renderer.setSize(w,h);camera.left=-7*w/h;camera.right=7*w/h;camera.top=7;camera.bottom=-7;camera.updateProjectionMatrix();dirty=true;};
  const observer=new ResizeObserver(resize);observer.observe(el);resize();
  const ambient=new THREE.HemisphereLight(0xffe7bc,0x443e32,2);scene.add(ambient);
  const sun=new THREE.DirectionalLight(0xffe3b0,3);sun.position.set(2,10,-8);sun.castShadow=true;sun.shadow.mapSize.set(1024,1024);
  Object.assign(sun.shadow.camera,{left:-9,right:9,top:9,bottom:-9,near:.1,far:35});sun.shadow.bias=-.0003;scene.add(sun);
  for(const x of [-3,1,4]){const lamp=new THREE.PointLight(0xffba68,12,7,2);lamp.position.set(x,2.65,2.7);scene.add(lamp);}
  const models=new Map<string,THREE.Group>();const actors=new Map<string,THREE.Group>();let costumes:THREE.Material[]=[];
  const selected=new THREE.Mesh(new THREE.RingGeometry(.45,.5,40),new THREE.MeshBasicMaterial({color:0xcba85c,side:THREE.DoubleSide}));selected.rotation.x=-Math.PI/2;selected.position.y=.04;scene.add(selected);
  const loader=new GLTFLoader();
  Promise.all(['interior-saint-agnes','person','woman'].map(async name=>{
   const gltf=await loader.loadAsync(`/art/models/${name}.glb`);
   if(dead){disposeCityResources([gltf.scene]);return;}models.set(name,gltf.scene);
  })).then(()=>{if(dead)return;const room=models.get('interior-saint-agnes')!;
   room.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});scene.add(room);dirty=true;setStatus('');
  }).catch(()=>{if(!dead)setStatus('The 3D room could not load. The people and actions below remain available.');});
  let roster='',presentation='';
  const pick=new THREE.Raycaster();const pointer=new THREE.Vector2();let down={x:0,y:0};
  const press=(event:PointerEvent)=>{down={x:event.clientX,y:event.clientY};};
  const release=(event:PointerEvent)=>{
   if(Math.hypot(event.clientX-down.x,event.clientY-down.y)>5)return;
   const rect=canvas.getBoundingClientRect();pointer.set((event.clientX-rect.left)/rect.width*2-1,1-(event.clientY-rect.top)/rect.height*2);
   pick.setFromCamera(pointer,camera);const hit=pick.intersectObjects([...actors.values()],true)[0];
   let object:THREE.Object3D|undefined=hit?.object;while(object&&!object.userData.person)object=object.parent||undefined;
   if(object)latest.current.onPick(object.userData.person);
  };
  canvas.addEventListener('pointerdown',press);canvas.addEventListener('pointerup',release);
  const keys=(event:KeyboardEvent)=>{
   const command=cameraCommand(event);if(!command)return;
   if(command.startsWith('pan-')){
    const delta=screenPan({x:camera.position.x,z:camera.position.z},{x:controls.target.x,z:controls.target.z},command,.6/camera.zoom);
    const before=controls.target.clone();controls.target.x=THREE.MathUtils.clamp(controls.target.x+delta.x,-6,6);controls.target.z=THREE.MathUtils.clamp(controls.target.z+delta.z,-5,5);
    camera.position.add(controls.target.clone().sub(before));
   }else if(['rotate-left','rotate-right'].includes(command)){
    const offset=camera.position.clone().sub(controls.target);
    offset.applyAxisAngle(new THREE.Vector3(0,1,0),command==='rotate-left'?-.12:.12);
    camera.position.copy(controls.target).add(offset);
   }else if(['+','=','-'].includes(event.key)){
    camera.zoom=THREE.MathUtils.clamp(camera.zoom*(event.key==='-'?1/1.12:1.12),controls.minZoom,controls.maxZoom);camera.updateProjectionMatrix();
   }else if(event.key==='Home'){
    camera.position.set(13,14,-17);controls.target.set(0,1,0);camera.zoom=1;camera.updateProjectionMatrix();
   }else return;
   event.preventDefault();controls.update();dirty=true;
  };canvas.addEventListener('keydown',keys);
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current;
   const key=JSON.stringify(p.people.slice(0,9).map(w=>[w.id,w.face]));
   if(models.size===3&&key!==roster){roster=key;dirty=true;actors.forEach(a=>scene.remove(a));actors.clear();costumes.forEach(m=>m.dispose());costumes=[];
    p.people.slice(0,9).forEach((who,i)=>{
     const model=pedestrianModel(who.id,who.face),object=models.get(model)!.clone(true);
     costumes.push(...dressPedestrian(object,model,wardrobe(who.id,who.face)));
     object.position.set(-1+(i%3)*2,.03,-Math.floor(i/3)*2);object.rotation.y=Math.PI;
     object.userData.person=who.id;object.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});
     actors.set(who.id,object);scene.add(object);
    });
   }
   const chosen=actors.get(p.picked);selected.visible=!!chosen;if(chosen){selected.position.x=chosen.position.x;selected.position.z=chosen.position.z;}
   const hour=((p.minute/60)%24+24)%24;sun.intensity=hour>=6&&hour<20?3:.5;
   const stateKey=`${p.picked}:${p.minute}`;if(stateKey!==presentation){presentation=stateKey;dirty=true;}
   controls.update();
   if(dirty&&!document.hidden){
    const room=models.get('interior-saint-agnes');
    const left=room?.getObjectByName('interior-wall-left'),back=room?.getObjectByName('interior-wall-back');
    if(left)left.visible=camera.position.x>=-5.8;
    if(back)back.visible=camera.position.z<=4.8;
    renderer.render(scene,camera);renderedFrames++;dirty=false;
    if(models.size===3)canvas.dataset.interior=JSON.stringify({people:[...actors.keys()],picked:p.picked,
      drawCalls:renderer.info.render.calls,triangles:renderer.info.render.triangles,renderedFrames,
      zoom:camera.zoom,cutawayWalls:[...(!left?.visible?['left']:[]),...(!back?.visible?['back']:[])]});
   }
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();controls.removeEventListener('change',changed);controls.dispose();canvas.removeEventListener('keydown',keys);canvas.removeEventListener('pointerdown',press);canvas.removeEventListener('pointerup',release);disposeCityResources([scene,...models.values()]);renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[]);
 return <div className="interior3d"><div ref={host} className="interior3d-canvas"/><span className="interior3d-caption">SAINT AGNES · Drag / Q/E: orbit · Scroll / +/−: zoom · WASD / arrows: pan · Home: reset · Select a person</span>{status&&<p role="status">{status}</p>}</div>;
}
