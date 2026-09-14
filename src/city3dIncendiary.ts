import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';

export const INCENDIARY_RELEASE=3.1;
export const INCENDIARY_IMPACT=3.9;
export const INCENDIARY_SECONDS=6.8;
const smooth=(n:number)=>{const t=THREE.MathUtils.clamp(n,0,1);return t*t*(3-2*t);};
/** Fixed-distance walk with smooth acceleration/braking and unchanged beat times. */
export function incendiaryStride(seconds:number,distance:number,duration:number,ramp:number){
 const t=THREE.MathUtils.clamp(seconds,0,duration),r=Math.min(ramp,duration/2),peak=distance/(duration-r);
 const integral=(u:number)=>u*u*u-u*u*u*u/2;
 let travelled:number,speed:number;
 if(t<r){const u=t/r;travelled=peak*r*integral(u);speed=peak*smooth(u);}
 else if(t>duration-r){const u=(duration-t)/r;travelled=distance-peak*r*integral(u);speed=peak*smooth(u);}
 else {travelled=peak*(t-r/2);speed=peak;}
 return {distance:travelled,weight:speed/peak};
}
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
  const incoming=incendiaryStride(seconds-.2,2,2/1.2,.22);
  const outgoing=incendiaryStride(seconds-4.25,4,4/2.1,.28);
  const approach=incoming.distance,escape=outgoing.distance;
  const weight=Math.max(incoming.weight,outgoing.weight);
  const phase=(approach+escape)/1.15*Math.PI*2;
  this.actor.position.set(-2+approach-escape,.02+Math.abs(Math.sin(phase))*.025*weight,0);
  this.actor.rotation.set(0,seconds<1.87?Math.PI/2*(1-smooth((seconds-1.45)/.42)):-Math.PI/2*smooth((seconds-4)/.25),0);
  for(const side of [-1,1]){
   const swing=phase+(side<0?0:Math.PI);
   for(const [name,angle] of [[`leg${side}`,Math.sin(swing)*.4*weight],[`knee${side}`,Math.max(0,Math.sin(swing+.7))*.65*weight],[`arm${side}`,-Math.sin(swing)*.25*weight],[`elbow${side}`,0]] as const){
    this.actor.getObjectByName(name)?.rotation.set(angle,0,0);
   }
  }
  const wind=smooth((seconds-1.9)/.7),throwing=smooth((seconds-2.75)/.35),lower=smooth((seconds-3.25)/.5);
  const carry=new THREE.Vector3(.43,1.02,.18),back=new THREE.Vector3(.46,1.65,-.28),forward=new THREE.Vector3(.40,1.65,.43);
  const grip=carry.clone().lerp(back,wind).lerp(forward,throwing).lerp(carry,lower);
  if(seconds<4){
   const joints=['arm1','elbow1'].map(name=>this.actor.getObjectByName(name)).filter((joint):joint is THREE.Object3D=>!!joint);
   const relaxed=joints.map(joint=>joint.quaternion.clone());
   aimArm(this.actor,1,grip,new THREE.Vector3(1,0,-.2));
   const settle=smooth((seconds-3.75)/.25);
   joints.forEach((joint,i)=>joint.quaternion.slerp(relaxed[i],settle));
  }
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
 const ray=new THREE.Raycaster(),entry=new THREE.Vector3();
 // Mesh.raycast tests an infinite ray against its box before testing triangles.
 // Our rays are short segments: reject boxes reached only beyond their end.
 const colliders:{mesh:THREE.Mesh;bounds:THREE.Box3}[]=[];
 building.traverse(o=>{if(o instanceof THREE.Mesh){
  o.geometry.computeBoundingBox();
  colliders.push({mesh:o,bounds:o.geometry.boundingBox!.clone().applyMatrix4(o.matrixWorld)});
 }});
 for(const window of windows)for(const loft of [.8,.4,.15,1.4,2,2.8]){
  // The bottle body contacts the facade before its grip reaches the glass.
  const target=window.clone().add(new THREE.Vector3(0,0,-.30));
  let clear=true,previous=start.clone();
  for(let frame=1;frame<=48&&clear;frame++){
   const t=frame/48,point=start.clone().lerp(target,t);point.y+=loft*4*t*(1-t);
   const direction=point.clone().sub(previous),length=direction.length();direction.normalize();
   for(const offset of offsets){ray.set(previous.clone().add(offset),direction);ray.near=0;ray.far=length;
    for(const {mesh,bounds} of colliders){
     if(!bounds.containsPoint(ray.ray.origin)){
      const contact=ray.ray.intersectBox(bounds,entry);
      if(!contact||contact.distanceToSquared(ray.ray.origin)>length*length+1e-12)continue;
     }
     if(ray.intersectObject(mesh,false).length){clear=false;break;}
    }
    if(!clear)break;
   }
   previous=point;
  }
  if(clear)return {target,loft};
 }
 return null;
}

/** Check the complete actor corridor against authored triangles, not whole-building boxes. */
export function incendiaryCorridorClear(root:{x:number;z:number},building:THREE.Object3D){
 const corridor=new THREE.Box3(new THREE.Vector3(root.x-4.6,.20,root.z-.9),new THREE.Vector3(root.x+1.6,2.5,root.z+.9));
 const triangle=new THREE.Triangle(),bounds=new THREE.Box3();let clear=true;
 building.updateWorldMatrix(true,true);
 building.traverse(o=>{
  if(!clear||!(o instanceof THREE.Mesh))return;
  o.geometry.computeBoundingBox();bounds.copy(o.geometry.boundingBox!).applyMatrix4(o.matrixWorld);
  if(!bounds.intersectsBox(corridor))return;
  const position=o.geometry.attributes.position,index=o.geometry.index;
  const count=index?.count??position.count;
  for(let i=0;i<count;i+=3){
   triangle.a.fromBufferAttribute(position,index?index.getX(i):i).applyMatrix4(o.matrixWorld);
   triangle.b.fromBufferAttribute(position,index?index.getX(i+1):i+1).applyMatrix4(o.matrixWorld);
   triangle.c.fromBufferAttribute(position,index?index.getX(i+2):i+2).applyMatrix4(o.matrixWorld);
   if(corridor.intersectsTriangle(triangle)){clear=false;break;}
  }
 });
 return clear;
}

/** Validate a forecourt before reserving it; the release is local to the cast. */
export function incendiaryStagingFlight(root:{x:number;z:number},release:THREE.Vector3,windows:THREE.Vector3[],building:THREE.Object3D){
 if(!incendiaryCorridorClear(root,building))return null;
 const origin=new THREE.Vector3(root.x,.2,root.z);
 const ordered=[...windows].sort((a,b)=>Math.abs(a.x-root.x)-Math.abs(b.x-root.x));
 return incendiaryFlight(release.clone().add(origin),ordered,building);
}

/** Local ballistic glass scatter; begins exactly at the bottle's impact. */
export function incendiaryShard(index:number,age:number,height:number){
 const vx=Math.sin(index*2.399)*(.3+index%3*.12),vz=-(.65+index%4*.22),vy=.3+index%3*.22;
 const landing=(vy+Math.sqrt(vy*vy+19.62*Math.max(0,height)))/9.81;
 const t=THREE.MathUtils.clamp(age,0,landing),fade=1-smooth((age-1.7)/.7);
 return {x:vx*t,y:Math.max(-height,vy*t-4.905*t*t),z:vz*t,
  rx:index+t*5,ry:index*1.7+t*3,rz:index*.8+t*4,
  scale:age<0?0:(.65+index%4*.12)*fade};
}

/** Trace the curved interval, including skipped frames, with chip clearance samples. */
export function incendiaryShardObstructed(index:number,fromAge:number,toAge:number,height:number,impact:THREE.Vector3,building:THREE.Object3D){
 if(toAge<0||toAge<=fromAge)return false;
 const start=Math.max(0,fromAge),end=Math.min(2.4,toAge);
 if(start>=end)return false;
 const meshes:{mesh:THREE.Mesh;bounds:THREE.Box3}[]=[];
 building.updateWorldMatrix(true,true);
 building.traverse(o=>{if(o instanceof THREE.Mesh){if(!o.geometry.boundingBox)o.geometry.computeBoundingBox();meshes.push({mesh:o,bounds:o.geometry.boundingBox!.clone().applyMatrix4(o.matrixWorld)});}});
 const ray=new THREE.Raycaster(),entry=new THREE.Vector3(),from=new THREE.Vector3(),to=new THREE.Vector3(),delta=new THREE.Vector3();
 const offsets=[[0,0,0],[.07,0,0],[-.07,0,0],[0,.07,0],[0,-.07,0],[0,0,.07],[0,0,-.07]];
 const steps=Math.ceil((end-start)*120);
 let previous=incendiaryShard(index,start,height);
 for(let step=1;step<=steps;step++){
  const next=incendiaryShard(index,start+(end-start)*step/steps,height);
  from.set(previous.x,previous.y,previous.z).add(impact);to.set(next.x,next.y,next.z).add(impact);
  delta.copy(to).sub(from);const length=delta.length();
  if(length>1e-7){
   ray.ray.direction.copy(delta).divideScalar(length);ray.near=0;ray.far=length;
   for(const offset of offsets){
    ray.ray.origin.copy(from).add(new THREE.Vector3(...offset));
    for(const {mesh,bounds} of meshes){
     if(!bounds.containsPoint(ray.ray.origin)){
      const contact=ray.ray.intersectBox(bounds,entry);
      if(!contact||entry.distanceToSquared(ray.ray.origin)>length*length+1e-12)continue;
     }
     if(ray.intersectObject(mesh,false).length)return true;
    }
   }
  }
  previous=next;
 }
 return false;
}
