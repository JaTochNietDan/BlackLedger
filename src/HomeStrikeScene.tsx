import {useEffect,useRef,useState,type ReactNode} from 'react';
import * as THREE from 'three';
import {GLTFLoader} from 'three/addons/loaders/GLTFLoader.js';
import {CityAssassination,executionSpatter,MELEE_IMPACTS,ASSASSINATION_SHOT} from './city3dAssassination';
import {homeStrikeRoom} from './homeStrike';
import {pedestrianModel} from './city3dCast';
import {wardrobe,dressPedestrian} from './city3dWardrobe';
import {GunfireAudio,BlastAudio,gunfightPose} from './city3dEvents';
import {weaponShots} from './city3dWeapons';
import {playCityGunshot,playMoment,playRecordedEffect,preloadCityGunshots,preloadCityEffects,soundOn} from './sound';
import {disposeCityResources} from './city3dResources';
import type {VisualCue,Snapshot} from './types';
import './homeStrike.css';

/** A recorded result, never a second combat simulation. No actions are sent here. */
export function HomeStrikeScene({cue,world,motion,overlay,onDone}:{cue:VisualCue;world:Snapshot;motion:boolean;overlay:ReactNode;onDone:(id:string)=>void}){
 const host=useRef<HTMLDivElement>(null),latest=useRef({motion,onDone});latest.current={motion,onDone};
 const [status,setStatus]=useState('Opening the home…');
 useEffect(()=>{
  const roomSettings=homeStrikeRoom(cue.target),origin=roomSettings.origin;
  const el=host.current!;let dead=false,frame=0,cast:CityAssassination|undefined,start=0,finished=false;
  const scene=new THREE.Scene();scene.background=new THREE.Color(0x171b18);
  const renderer=new THREE.WebGLRenderer({antialias:true});renderer.setPixelRatio(Math.min(devicePixelRatio,1.5));renderer.shadowMap.enabled=true;
  renderer.outputColorSpace=THREE.SRGBColorSpace;renderer.toneMapping=THREE.ACESFilmicToneMapping;renderer.toneMappingExposure=1.05;
  el.append(renderer.domElement);renderer.domElement.setAttribute('aria-label','Recorded attack inside the resident’s home');
  const camera=new THREE.PerspectiveCamera(43,1,.1,60);camera.position.set(5.7,6.4,10.8);camera.lookAt(0,.7,roomSettings.focusZ);
  const resize=()=>{renderer.setSize(el.clientWidth,el.clientHeight);camera.aspect=el.clientWidth/el.clientHeight;camera.updateProjectionMatrix();if(cast)renderer.render(scene,camera);};
  const observer=new ResizeObserver(resize);observer.observe(el);resize();
  scene.add(new THREE.HemisphereLight(0xffe7bc,0x443e32,2));const sun=new THREE.DirectionalLight(0xffe3b0,3);sun.position.set(2,9,5);sun.castShadow=true;sun.shadow.mapSize.set(1024,1024);Object.assign(sun.shadow.camera,{left:-6,right:6,top:6,bottom:-6});scene.add(sun);
  const tier=cue.attacker!.weapon,weaponName=tier===3?'thompson':tier===2?'shotgun':tier===1?'revolver':undefined;
  const beats=weaponName?weaponShots(weaponName,cue.strike!.variant):MELEE_IMPACTS;
  const audio=new GunfireAudio(()=>weaponName?playCityGunshot(weaponName):playMoment('body-hit'),beats,true);
  const reaction=new BlastAudio(()=>playRecordedEffect('pain'));
  void preloadCityGunshots();void preloadCityEffects();
  const flash=new THREE.PointLight(0xffd08a,0,5);scene.add(flash);
  const muzzleGlow=new THREE.Mesh(new THREE.SphereGeometry(.065,8,6),new THREE.MeshBasicMaterial({color:0xffdfa0}));muzzleGlow.visible=false;scene.add(muzzleGlow);
  const droplets=new THREE.InstancedMesh(new THREE.SphereGeometry(1,6,4),new THREE.MeshBasicMaterial({color:0x720c12}),20);droplets.frustumCulled=false;scene.add(droplets);
  const temp=new THREE.Object3D(),models:THREE.Group[]=[],costumes:THREE.Material[]=[];
  const loader=new GLTFLoader();
  const actor=async(id:string)=>{
   const player=id==='player',face=player?world.player.face:world.everyone?.find(p=>p.id===id)?.face;
   const identity=player?world.player.name:id,model=pedestrianModel(identity,face,player);
   const loaded=await loader.loadAsync(`/art/models/${model}.glb`);models.push(loaded.scene);
   if(dead){disposeCityResources([loaded.scene]);return loaded.scene;}
   costumes.push(...dressPedestrian(loaded.scene,model,wardrobe(identity,face,player)));return loaded.scene;
  };
  Promise.all([loader.loadAsync(`/art/models/${roomSettings.model}.glb`),actor(cue.attacker!.id),actor(cue.strike!.victim.id),weaponName?loader.loadAsync(`/art/models/${weaponName}.glb`):Promise.resolve(undefined)]).then(([room,attacker,victim,weapon])=>{
   models.push(room.scene);if(weapon)models.push(weapon.scene);
   if(dead){disposeCityResources(models);return;}
   scene.add(room.scene);cast=new CityAssassination(attacker,victim,weapon?.scene,cue.strike!.variant,weaponName);
   cast.root.position.set(origin.x,0,origin.z);scene.add(cast.root);
   scene.traverse(o=>{if(o instanceof THREE.Mesh){o.castShadow=true;o.receiveShadow=true;}});
   start=performance.now();setStatus('');
  }).catch(()=>{if(!dead){setStatus('The home scene could not load. The recorded result is available below.');finished=true;latest.current.onDone(cue.id);}});
  const reduced=matchMedia('(prefers-reduced-motion: reduce)');
  const tick=(now:number)=>{
   if(dead)return;frame=requestAnimationFrame(tick);if(!cast||finished)return;
   const enabled=latest.current.motion&&!reduced.matches,seconds=enabled?(now-start)/1000:cast.duration;
   cast.update(Math.min(seconds,cast.duration));
   const pulse=weaponName?gunfightPose(seconds,beats).flash:false;
   const muzzle=cast.weapon?.getObjectByName('muzzle');if(muzzle)muzzle.getWorldPosition(flash.position);flash.intensity=pulse?12:0;muzzleGlow.position.copy(flash.position);muzzleGlow.visible=!!muzzle&&pulse;
   for(let i=0;i<20;i++){
    const drop=executionSpatter(i,seconds,cue.strike!.variant==='back-of-head');
    temp.position.set(origin.x+drop.x,drop.y,origin.z+drop.z);temp.scale.setScalar(weaponName?drop.size:0);temp.updateMatrix();droplets.setMatrixAt(i,temp.matrix);
   }
   droplets.instanceMatrix.needsUpdate=true;
   audio.update(seconds,enabled&&soundOn());reaction.update(seconds-(weaponName?ASSASSINATION_SHOT:MELEE_IMPACTS[2]),enabled&&soundOn());
   renderer.render(scene,camera);renderer.domElement.dataset.homeStrike=JSON.stringify({cue:cue.id,seconds,ready:true,finished:seconds>=cast.duration,attacker:cue.attacker!.id,victim:cue.strike!.victim.id});
   if(seconds>=cast.duration){finished=true;latest.current.onDone(cue.id);}
  };frame=requestAnimationFrame(tick);
  return()=>{dead=true;cancelAnimationFrame(frame);observer.disconnect();audio.dispose();reaction.dispose();disposeCityResources([scene,...models]);costumes.forEach(m=>m.dispose());renderer.dispose();renderer.forceContextLoss();renderer.domElement.remove();};
 },[cue.id]);
 return <div className="home-strike-scene"><div ref={host} className="home-strike-canvas"/>{status&&<p role="status">{status}</p>}<div className="city3d-story">{overlay}</div></div>;
}
