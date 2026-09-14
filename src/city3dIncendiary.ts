import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';

export const INCENDIARY_RELEASE=3.1;
export const INCENDIARY_IMPACT=3.9;
export const INCENDIARY_SECONDS=6.8;
const smooth=(n:number)=>{const t=THREE.MathUtils.clamp(n,0,1);return t*t*(3-2*t);};
/** A presentation cast in local metres. The caller reserves and places its path. */
export class CityIncendiary {
 readonly root=new THREE.Group();
 readonly release=new THREE.Vector3();
 loft=.8;
 private readonly hand=new THREE.Vector3();
 private readonly releaseRotation=new THREE.Quaternion();
 constructor(readonly actor:THREE.Group,readonly bottle:THREE.Group,readonly target:THREE.Vector3){
  this.root.add(actor,bottle);
  this.poseActor(INCENDIARY_RELEASE);
  this.hold();
  this.release.copy(this.bottle.position);
  this.releaseRotation.copy(this.bottle.quaternion);
  this.update(0);
 }
 get duration(){return INCENDIARY_SECONDS;}
 private grip(){
  // The authored hand centre is 29.5cm below the elbow, with its natural offset.
  const elbow=this.actor.getObjectByName('elbow1');
  this.hand.set(.01,-.295,.005);
  if(elbow){elbow.localToWorld(this.hand);this.root.worldToLocal(this.hand);}
  return this.hand;
 }
 private hold(){
  this.bottle.position.copy(this.grip());
  const elbow=this.actor.getObjectByName('elbow1');
  const along=elbow?this.root.worldToLocal(elbow.getWorldPosition(new THREE.Vector3())).sub(this.bottle.position).normalize():new THREE.Vector3(0,1,0);
  // Keep both the body and protruding cloth clear of the sleeve.
  const forward=new THREE.Vector3(0,0,1).applyQuaternion(this.actor.quaternion);
  const across=along.cross(forward).normalize();
  const outside=new THREE.Vector3(1,0,0).applyQuaternion(this.actor.quaternion);
  if(across.dot(outside)>0)across.negate();
  this.bottle.quaternion.setFromUnitVectors(new THREE.Vector3(0,1,0),across);
 }
 private poseActor(seconds:number){
  const approach=THREE.MathUtils.clamp((seconds-.2)*1.2,0,2);
  const fleeing=Math.max(0,seconds-4.25),escape=Math.min(4,fleeing*2.1);
  const walking=approach>0&&approach<2||escape>0&&escape<4;
  const phase=(approach+escape)/1.15*Math.PI*2;
  this.actor.position.set(-2+approach-escape,.02+(walking?Math.abs(Math.sin(phase))*.025:0),0);
  this.actor.rotation.set(0,seconds<1.87?Math.PI/2*(1-smooth((seconds-1.45)/.42)):-Math.PI/2*smooth((seconds-4)/.25),0);
  for(const side of [-1,1]){
   const swing=phase+(side<0?0:Math.PI);
   for(const [name,angle] of [[`leg${side}`,walking?Math.sin(swing)*.4:0],[`knee${side}`,walking?Math.max(0,Math.sin(swing+.7))*.65:0],[`arm${side}`,walking?-Math.sin(swing)*.25:0],[`elbow${side}`,0]] as const){
    this.actor.getObjectByName(name)?.rotation.set(angle,0,0);
   }
  }
  const wind=smooth((seconds-1.9)/.7),throwing=smooth((seconds-2.75)/.35),lower=smooth((seconds-3.25)/.5);
  const carry=new THREE.Vector3(.43,1.02,.18),back=new THREE.Vector3(.46,1.65,-.28),forward=new THREE.Vector3(.40,1.65,.43);
  const grip=carry.clone().lerp(back,wind).lerp(forward,throwing).lerp(carry,lower);
  if(seconds<4)aimArm(this.actor,1,grip,new THREE.Vector3(1,0,-.2));
  this.actor.updateMatrixWorld(true);
 }
 update(seconds:number){
  this.poseActor(seconds);
  const released=seconds>=INCENDIARY_RELEASE;
  this.bottle.visible=seconds<INCENDIARY_IMPACT;
  if(!released){
   this.hold();
  }else{
   const t=THREE.MathUtils.clamp((seconds-INCENDIARY_RELEASE)/(INCENDIARY_IMPACT-INCENDIARY_RELEASE),0,1);
   this.bottle.position.copy(this.release).lerp(this.target,t);
   this.bottle.position.y+=this.loft*4*t*(1-t);
   this.bottle.quaternion.copy(this.releaseRotation).multiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(1,0,0),t*Math.PI*2));
  }
  this.root.updateMatrixWorld(true);
  return {released,impact:seconds>=INCENDIARY_IMPACT,escaping:seconds>4.25};
 }
}

type Flight={target:THREE.Vector3;loft:number};
const flightCache=new WeakMap<THREE.Object3D,{geometry:string;paths:Map<string,Flight|null>}>();
const bufferIDs=new WeakMap<object,number>();let nextBufferID=1;
const bufferID=(buffer:object)=>{let id=bufferIDs.get(buffer);if(id===undefined){id=nextBufferID++;bufferIDs.set(buffer,id);}return id;};
const copyFlight=(flight:Flight|null)=>flight?{target:flight.target.clone(),loft:flight.loft}:null;
/** Cached only for static building geometry; matrices and buffer versions invalidate it. */
export function incendiaryFlight(start:THREE.Vector3,windows:THREE.Vector3[],building:THREE.Object3D){
 building.updateWorldMatrix(true,true);
 const geometry:string[]=[];
 building.traverse(o=>{if(o instanceof THREE.Mesh){
  const attributes=(Object.values(o.geometry.attributes) as (THREE.BufferAttribute|THREE.InterleavedBufferAttribute)[]).map(a=>a instanceof THREE.InterleavedBufferAttribute?`${bufferID(a)}:${bufferID(a.data)}:${a.data.version}`:`${bufferID(a)}:${a.version}`);
  const materials=Array.isArray(o.material)?o.material:[o.material];
  geometry.push(o.geometry.uuid+':'+attributes.join(',')+':'+(o.geometry.index?`${bufferID(o.geometry.index)}:${o.geometry.index.version}`:'none')+':'+o.matrixWorld.elements.join(',')+':'+materials.map(m=>m.side).join(','));
 }});
 const signature=geometry.join(';');
 let cache=flightCache.get(building);
 if(!cache||cache.geometry!==signature){cache={geometry:signature,paths:new Map()};flightCache.set(building,cache);}
 const key=JSON.stringify([start.toArray(),windows.map(w=>w.toArray())]);
 if(cache.paths.has(key))return copyFlight(cache.paths.get(key)!);
 const result=findIncendiaryFlight(start,windows,building);
 if(cache.paths.size>=12)cache.paths.delete(cache.paths.keys().next().value!);
 cache.paths.set(key,copyFlight(result));return result;
}

/** Test a conservative sampled bottle envelope once when selecting a facade target. */
function findIncendiaryFlight(start:THREE.Vector3,windows:THREE.Vector3[],building:THREE.Object3D):Flight|null{
 const offsets=[new THREE.Vector3(),...[-1,1].flatMap(side=>[new THREE.Vector3(side*.27,0,0),new THREE.Vector3(0,side*.27,0),new THREE.Vector3(0,0,side*.27)])];
 const ray=new THREE.Raycaster();
 for(const window of windows)for(const loft of [.8,.4,.15,1.4,2,2.8]){
  // The bottle body contacts the facade before its grip reaches the glass.
  const target=window.clone().add(new THREE.Vector3(0,0,-.30));
  let clear=true,previous=start.clone();
  for(let frame=1;frame<=48&&clear;frame++){
   const t=frame/48,point=start.clone().lerp(target,t);point.y+=loft*4*t*(1-t);
   const direction=point.clone().sub(previous),length=direction.length();direction.normalize();
   for(const offset of offsets){ray.set(previous.clone().add(offset),direction);ray.near=0;ray.far=length;
    if(ray.intersectObject(building,true).length){clear=false;break;}
   }
   previous=point;
  }
  if(clear)return {target,loft};
 }
 return null;
}
