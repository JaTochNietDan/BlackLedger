import {interiorSettings} from './interiorSettings';
import {LaundryMotion,laundryRunningMachines,type LaundryOperation} from './laundryMotion';
import {InteriorArrival} from './interiorArrival';
import {CounterWipe,LinenPress} from './interiorService';
import {InteriorCastBatch} from './interiorCastBatch';
import {placementsForInterior,interiorPlayerSpot,poseInteriorOccupant,type InteriorPlace} from './interiorStaging';
import {cameraCommand, KeyboardPan, bindKeyboardPan} from './city3dControls';
import {useEffect, useRef, useState} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {OrbitControls} from 'three/addons/controls/OrbitControls.js';
import {Reflector} from 'three/addons/objects/Reflector.js';
import type {Presence, Person} from './types';
import {wardrobe, dressPedestrian} from './city3dWardrobe';
import {pedestrianModel} from './city3dCast';
import {disposeCityResources} from './city3dResources';
import './interior3d.css';

export function Interior3D(props:{place:InteriorPlace;operation?:LaundryOperation;motion:boolean;player:Pick<Person,'name'|'face'|'alive'>;people:Presence[];picked:string;onPick:(id:string)=>void;minute:number}) {
 const host=useRef<HTMLDivElement>(null), latest=useRef(props);latest.current=props;
 const settings=interiorSettings[props.place],roomName=settings.name;
 const [enlarged,setEnlarged]=useState(false);
 const [drawerOpen,setDrawerOpen]=useState(false),drawerState=useRef(false);drawerState.current=drawerOpen;
 const expanded=useRef(false);expanded.current=enlarged;
 const extraPeople=props.people.length-placementsForInterior(props.place,props.people).size;
 const [status,setStatus]=useState(`Opening ${roomName}…`);
 useEffect(()=>{
  const flat=props.place==='flat'||props.place==='lodging',laundry=props.place==='laundry',lobby=props.place!=='bar',roomModel=settings.model;
  const origin=new THREE.Vector3(13,14,lobby?17:-17), centre=new THREE.Vector3(0,1,lobby?-1:0),span=settings.span;
  setStatus(`Opening ${roomName}…`);
  const el=host.current!;let dead=false,frame=0,dirty=true,renderedFrames=0;
  const scene=new THREE.Scene();scene.background=new THREE.Color(0x171b18);
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));
  renderer.shadowMap.enabled=true;renderer.shadowMap.type=THREE.PCFShadowMap;
  renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;renderer.toneMappingExposure=1.1;
  const canvas=renderer.domElement;canvas.setAttribute('aria-label',`${roomName} 3D interior. Drag or Q/E to orbit, scroll or +/- to zoom, Home resets, WASD or arrows pan${flat?'.':'; click a person to select their actions.'}`);canvas.tabIndex=0;el.append(canvas);
  const camera=new THREE.OrthographicCamera(-9,9,7,-7,.1,100);camera.position.copy(origin);
  const controls=new OrbitControls(camera,canvas);controls.target.copy(centre);controls.minZoom=.7;controls.maxZoom=3;
  controls.minPolarAngle=.35;controls.maxPolarAngle=1.15;controls.enablePan=false;controls.update();
  const changed=()=>{dirty=true;};controls.addEventListener('change',changed);
  const resize=()=>{
   const w=Math.max(1,el.clientWidth),h=Math.max(1,el.clientHeight),aspect=w/h;
   // Enlarging a room in a narrow window must not crop its side walls.
   // Keep the authored vertical framing on wide screens, and preserve a
   // minimum horizontal field of view as the viewport becomes taller.
   const viewSpan=span*Math.max(1,1.35/aspect);
   renderer.setSize(w,h);camera.left=-viewSpan*aspect;camera.right=viewSpan*aspect;
   camera.top=viewSpan;camera.bottom=-viewSpan;camera.updateProjectionMatrix();dirty=true;
  };
  const observer=new ResizeObserver(resize);observer.observe(el);resize();
  const ambient=new THREE.HemisphereLight(0xffe7bc,0x443e32,2);scene.add(ambient);
  const sun=new THREE.DirectionalLight(0xffe3b0,3);sun.position.set(2,10,lobby?8:-8);sun.castShadow=true;sun.shadow.mapSize.set(1024,1024);
  Object.assign(sun.shadow.camera,{left:-9,right:9,top:9,bottom:-9,near:.1,far:35});sun.shadow.bias=-.0003;scene.add(sun);
  for(const position of settings.lamps){const lamp=new THREE.PointLight(0xffba68,settings.lampIntensity??12,7,2);lamp.position.set(...position);scene.add(lamp);}
  const reduce=matchMedia('(prefers-reduced-motion: reduce)');let reduced=reduce.matches;
  const reduction=()=>{reduced=reduce.matches;dirty=true;};reduce.addEventListener('change',reduction);
  const models=new Map<string,THREE.Group>();const actors=new Map<string,THREE.Group>();let costumes:THREE.Material[]=[];
  const selected=new THREE.Mesh(new THREE.RingGeometry(.45,.5,40),new THREE.MeshBasicMaterial({color:0xcba85c,side:THREE.DoubleSide}));selected.rotation.x=-Math.PI/2;selected.position.y=.04;scene.add(selected);
  const playerMarker=new THREE.Mesh(new THREE.RingGeometry(.36,.41,40),new THREE.MeshBasicMaterial({color:0xede2bd,side:THREE.DoubleSide}));
  playerMarker.rotation.x=-Math.PI/2;playerMarker.position.y=.065;playerMarker.visible=false;scene.add(playerMarker);
  let machines:LaundryMotion|undefined;
  let fittingMirror:Reflector|undefined;
  renderer.info.autoReset=false;
  const loader=new GLTFLoader();
  const modelNames=[roomModel,'person','woman',...(!lobby?['bar-cloth']:[])];
  Promise.all(modelNames.map(async name=>{
   const gltf=await loader.loadAsync(`/art/models/${name}.glb`);
   if(dead){disposeCityResources([gltf.scene]);return;}models.set(name,gltf.scene);
  })).then(()=>{if(dead)return;const room=models.get(roomModel)!;
   room.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});scene.add(room);
   if(props.place==='tailor'){
    // Replace the authored silver surface so the reflected camera cannot see
    // its opaque back face through the oblique clipping plane.
    room.traverse(o=>{if(o instanceof THREE.Mesh&&(Array.isArray(o.material)?o.material:[o.material]).some(m=>m.name==='Ruttledge silver mirror'))o.visible=false;});
    fittingMirror=new Reflector(new THREE.PlaneGeometry(1.36,2.47),{color:0xd8d8cd,textureWidth:512,textureHeight:1024,multisample:0,clipBias:0});
    fittingMirror.name='fitting-mirror';fittingMirror.position.set(-3.83,1.65,-4.60);scene.add(fittingMirror);
    room.updateMatrixWorld(true);room.getObjectByName('interior-wall-back')!.attach(fittingMirror);
   }
   if(laundry)machines=new LaundryMotion(room);dirty=true;setStatus('');
  }).catch(()=>{if(!dead)setStatus('The 3D room could not load. The people and actions below remain available.');});
  let linenService:LinenPress|undefined;
  let service:CounterWipe|undefined,cloth:THREE.Group|undefined,serviceSeconds=0;
  let castBatch:InteriorCastBatch|undefined;
  let arrival:InteriorArrival|undefined,arrivalSeconds=0,arriving=false;
  let roster='',presentation='',drawerTravel=0;
  const pick=new THREE.Raycaster();const pointer=new THREE.Vector2();let down={x:0,y:0};
  const press=(event:PointerEvent)=>{down={x:event.clientX,y:event.clientY};};
  const release=(event:PointerEvent)=>{
   if(Math.hypot(event.clientX-down.x,event.clientY-down.y)>5)return;
   const rect=canvas.getBoundingClientRect();pointer.set((event.clientX-rect.left)/rect.width*2-1,1-(event.clientY-rect.top)/rect.height*2);
   pick.setFromCamera(pointer,camera);const hit=pick.intersectObjects(castBatch?[castBatch.root,...castBatch.unbatched]:[...actors.values()],true)[0];
   if(hit?.instanceId!==undefined){const id=hit.object.userData.people[hit.instanceId];latest.current.onPick(id==='player'?'':id);return;}
   let object:THREE.Object3D|undefined=hit?.object;while(object&&!object.userData.person)object=object.parent||undefined;
   if(object)latest.current.onPick(object.userData.person==='player'?'':object.userData.person);
  };
  canvas.addEventListener('pointerdown',press);canvas.addEventListener('pointerup',release);
  const keyboardPan=new KeyboardPan(), unbindPan=bindKeyboardPan(canvas,keyboardPan);
  let panTime=performance.now();
  const keys=(event:KeyboardEvent)=>{
   if(event.key==='Escape'&&expanded.current){event.preventDefault();event.stopPropagation();setEnlarged(false);return;}
   const command=cameraCommand(event);if(!command)return;
   if(command.startsWith('pan-')||command.startsWith('rotate-')){
    keyboardPan.press(event);
   }else if(['+','=','-'].includes(event.key)){
    camera.zoom=THREE.MathUtils.clamp(camera.zoom*(event.key==='-'?1/1.12:1.12),controls.minZoom,controls.maxZoom);camera.updateProjectionMatrix();
   }else if(event.key==='Home'){
    camera.position.copy(origin);controls.target.copy(centre);camera.zoom=1;camera.updateProjectionMatrix();
   }else return;
   event.preventDefault();controls.update();dirty=true;
  };canvas.addEventListener('keydown',keys);
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);const p=latest.current;
   const seconds=(now-panTime)/1000;
   const delta=keyboardPan.step(camera.position,controls.target,seconds,8/camera.zoom);panTime=now;
   const turn=keyboardPan.rotation(seconds);
   if(turn){const offset=camera.position.clone().sub(controls.target).applyAxisAngle(new THREE.Vector3(0,1,0),turn);camera.position.copy(controls.target).add(offset);dirty=true;}
   if(delta.x||delta.z){
    const before=controls.target.clone();controls.target.x=THREE.MathUtils.clamp(controls.target.x+delta.x,-6,6);controls.target.z=THREE.MathUtils.clamp(controls.target.z+delta.z,-5,5);
    camera.position.add(controls.target.clone().sub(before));dirty=true;
   }
   const key=JSON.stringify([p.people.map(w=>[w.id,w.face,w.role]),[p.player.name,p.player.face,p.player.alive]]);
   if(models.size===modelNames.length&&key!==roster){roster=key;dirty=true;arrival=undefined;service=undefined;linenService=undefined;cloth?.removeFromParent();cloth=undefined;castBatch?.dispose();actors.forEach(a=>scene.remove(a));actors.clear();costumes.forEach(m=>m.dispose());costumes=[];
    const placements=placementsForInterior(p.place,p.people);
    p.people.forEach(who=>{
     const spot=placements.get(who.id);if(!spot)return;
     const model=pedestrianModel(who.id,who.face),object=models.get(model)!.clone(true);
     costumes.push(...dressPedestrian(object,model,wardrobe(who.id,who.face)));
     poseInteriorOccupant(object,spot);object.userData.spot=spot.id;
     object.userData.person=who.id;object.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});
     actors.set(who.id,object);scene.add(object);
    });
    if(p.player.alive){
     const model=pedestrianModel(p.player.name,p.player.face,true),object=models.get(model)!.clone(true);
     costumes.push(...dressPedestrian(object,model,wardrobe(p.player.name,p.player.face,true)));
     const spot=interiorPlayerSpot(p.place);poseInteriorOccupant(object,spot);
     object.userData.spot=spot.id;object.userData.person='player';
     object.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});
     actors.set('player',object);scene.add(object);arrival=new InteriorArrival(object,p.place);
     if(!p.motion||reduced)arrivalSeconds=arrival.duration;
     arriving=arrival.pose(arrivalSeconds);
    }
    if(!lobby){const bartender=[...actors.values()].find(a=>a.userData.spot==='service');if(bartender){
     service=new CounterWipe(bartender);if(service.available){cloth=models.get('bar-cloth')!.clone(true);cloth.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});scene.add(cloth);service.pose(serviceSeconds);cloth.position.copy(service.clothPosition);}
    }}
    if(laundry){const worker=[...actors.values()].find(a=>a.userData.spot==='laundry-counter');if(worker)linenService=new LinenPress(worker);}
    castBatch=new InteriorCastBatch(actors);scene.add(castBatch.root);
   }
   const running=laundryRunningMachines(p.operation);
   if(machines?.step(seconds,running,p.motion&&!reduced,document.hidden))dirty=true;
   let poseChanged=linenService?.step(seconds,running>0,p.motion&&!reduced,document.hidden)??false;
   if(arrival&&arriving&&!document.hidden){
    arrivalSeconds=!p.motion||reduced?arrival.duration:Math.min(arrival.duration,arrivalSeconds+Math.min(.05,Math.max(0,seconds)));
    arriving=arrival.pose(arrivalSeconds);poseChanged=true;
   }
   if(service?.available&&cloth&&p.motion&&!reduced&&!document.hidden){serviceSeconds+=Math.min(.05,Math.max(0,seconds));service.pose(serviceSeconds);cloth.position.copy(service.clothPosition);poseChanged=true;}
   if(poseChanged){castBatch?.update();dirty=true;}
   const playerActor=actors.get('player');playerMarker.visible=!!playerActor;
   if(playerActor){playerMarker.position.x=playerActor.position.x;playerMarker.position.z=playerActor.position.z;}
   const chosen=actors.get(p.picked);selected.visible=!!chosen;if(chosen){selected.position.x=chosen.position.x;selected.position.z=chosen.position.z;}
   const hour=((p.minute/60)%24+24)%24;sun.intensity=hour>=6&&hour<20?3:.5;
   const stateKey=`${p.picked}:${p.minute}`;if(stateKey!==presentation){presentation=stateKey;dirty=true;}
   const drawer=models.get(roomModel)?.getObjectByName('burglary-drawer');
   if(drawer){
    const goal=drawerState.current ? .34 : 0;
    const next=(!p.motion||reduced)?goal:THREE.MathUtils.damp(drawerTravel,goal,12,Math.min(seconds,.1));
    if(Math.abs(next-drawerTravel)>.00001){drawerTravel=next;drawer.position.z=next;dirty=true;}
   }
   controls.update();
   if(dirty&&!document.hidden){
    const room=models.get(roomModel);
    const left=room?.getObjectByName('interior-wall-left'),back=room?.getObjectByName('interior-wall-back');
    if(left)left.visible=camera.position.x>=settings.leftWall;
    if(back)back.visible=lobby?camera.position.z>=settings.backWall:camera.position.z<=settings.backWall;
    // Moving service poses refresh the cast above; orbiting a static room
    // does not re-upload its unchanged instance buffers.
    renderer.info.reset();renderer.render(scene,camera);renderedFrames++;dirty=false;
    if(models.size===modelNames.length)canvas.dataset.interior=JSON.stringify({place:p.place,people:[...actors.keys()],occupants:[...actors].map(([id,a])=>({id,spot:a.userData.spot,x:a.position.x,y:a.position.y,z:a.position.z})),picked:p.picked,
      mirror:fittingMirror?{visible:!!back?.visible,resolution:[512,1024]}:undefined,drawCalls:renderer.info.render.calls,triangles:renderer.info.render.triangles,renderedFrames,
      zoom:camera.zoom,linen:linenService?{seconds:linenService.seconds,active:linenService.active}:undefined,machines:machines?{count:machines.count,running,seconds:machines.seconds}:undefined,arrival:arrival?{seconds:arrivalSeconds,duration:arrival.duration,moving:arriving}:undefined,service:service?.available?{seconds:serviceSeconds,cloth:service.clothPosition.toArray()}:undefined,omitted:Math.max(0,p.people.length-actors.size+(playerActor?1:0)),cutawayWalls:[...(!left?.visible?['left']:[]),...(!back?.visible?['back']:[])]});
   }
  };frame=requestAnimationFrame(tick);
  return()=>{reduce.removeEventListener('change',reduction);unbindPan();dead=true;cancelAnimationFrame(frame);observer.disconnect();controls.removeEventListener('change',changed);controls.dispose();canvas.removeEventListener('keydown',keys);canvas.removeEventListener('pointerdown',press);canvas.removeEventListener('pointerup',release);castBatch?.dispose();if(fittingMirror){fittingMirror.removeFromParent();fittingMirror.geometry.dispose();fittingMirror.dispose();}disposeCityResources([scene,...models.values()]);renderer.dispose();renderer.forceContextLoss();canvas.remove();};
 },[props.place]);
 return <div className={`interior3d${enlarged?' enlarged':''}`}>{(props.place==='flat'||props.place==='lodging')&&<button className="interior3d-drawer" aria-pressed={drawerOpen} onClick={()=>setDrawerOpen(!drawerOpen)}>{drawerOpen?'Close bedside drawer':'Open bedside drawer'}</button>}<button className="interior3d-expand" aria-pressed={enlarged} onClick={()=>{setEnlarged(!enlarged);host.current?.querySelector('canvas')?.focus();}}>{enlarged?'Standard room view':'Enlarge room'}</button><div ref={host} className="interior3d-canvas"/><span className="interior3d-caption">{roomName} · {props.player.name}: pale ring · Drag / Q/E: orbit · Scroll / +/−: zoom · WASD / arrows: pan · Home: reset{props.place!=='flat'&&props.place!=='lodging'&&' · Select a person'}{extraPeople>0&&` · ${extraPeople} more in the people list`}</span>{status&&<p role="status">{status}</p>}</div>;
}
