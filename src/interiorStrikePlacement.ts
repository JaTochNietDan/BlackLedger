import * as THREE from 'three';
import type {CityAssassination} from './city3dAssassination';
/** Test the full recorded approach and fall against the venue's authored room.
 * Fail visibly rather than put actors through furniture or move the hit outside. */
export function interiorStrikePlacement(room:THREE.Object3D,cast:CityAssassination,bounds:{x:number;z:number}){
 room.updateMatrixWorld(true);
 const ray=new THREE.Raycaster(),at=new THREE.Vector3(),directions=[new THREE.Vector3(1,0,0),new THREE.Vector3(-1,0,0),new THREE.Vector3(0,0,1),new THREE.Vector3(0,0,-1)];
 const clear=(x:number,z:number,yaw:number)=>{
  cast.root.position.set(x,0,z);cast.root.rotation.y=yaw;
  for(let step=0;step<=20;step++){
   cast.update(cast.duration*step/20);cast.root.updateMatrixWorld(true);
   for(const actor of [cast.attacker,cast.victim]){
    actor.getWorldPosition(at);
    ray.set(new THREE.Vector3(at.x,.18,at.z),new THREE.Vector3(0,-1,0));ray.near=0;ray.far=.22;
    if(!ray.intersectObject(room,true).length)return false;
    const box=new THREE.Box3().setFromObject(actor,true);
    for(const px of [box.min.x,box.max.x])for(const pz of [box.min.z,box.max.z]){
     ray.set(new THREE.Vector3(px,.18,pz),new THREE.Vector3(0,-1,0));ray.far=.22;
     if(!ray.intersectObject(room,true).length)return false;
    }
    for(const y of [.4,1,1.6])for(const direction of directions){
     ray.set(new THREE.Vector3(at.x,y,at.z),direction);ray.far=.38;
     if(ray.intersectObject(room,true).length)return false;
    }
   }
  }
  return true;
 };
 for(const yaw of [0,Math.PI/2,Math.PI,-Math.PI/2])for(let z=-bounds.z+1;z<bounds.z;z++)for(let x=-bounds.x+1;x<bounds.x;x++){
  if(clear(x,z,yaw)){cast.update(0);return{x,z,yaw};}
 }
 cast.update(0);return undefined;
}
