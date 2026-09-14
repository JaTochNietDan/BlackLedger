import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';
import type {Presence} from './types';
export type InteriorSpot={id:string;x:number;z:number;yaw:number;seat?:number};
// Metres in the glTF frame of saint_agnes_interior in tools/export_city3d.py.
const booths:InteriorSpot[]=[.5,3.2].flatMap((center,row)=>[-1,1].map(side=>({id:`booth-${row}-${side}`,x:-4.65+side*.48,z:-(center+side*.72),yaw:side<0?Math.PI:0,seat:.69})));
const stools:InteriorSpot[]=[-2,0,2,4].map((x,i)=>({id:`stool-${i}`,x,z:1.4,yaw:0,seat:.9}));
const standing:InteriorSpot[]=[{id:'aisle-a',x:1,z:-1,yaw:Math.PI},{id:'aisle-b',x:3,z:-2.8,yaw:Math.PI/2}];
export function interiorPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>(),available=[...booths,...stools,...standing];
 for(const who of [...people].sort((a,b)=>(a.id==='mara'?-1:b.id==='mara'?1:a.id.localeCompare(b.id)))){
  if(result.has(who.id))continue;
  const worker=/\b(bartender|barman|barmaid)\b/i.test(who.role||'');
  if(worker&&![...result.values()].some(s=>s.id==='service'))result.set(who.id,{id:'service',x:0,z:3.45,yaw:Math.PI});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}
export function poseInteriorOccupant(actor:THREE.Group,spot:InteriorSpot) {
 actor.rotation.set(0,spot.yaw,0);
 actor.position.set(spot.x,spot.seat===undefined?.03:spot.seat-.86,spot.z);
 if(spot.seat!==undefined){
  const forward=spot.seat>.8?.08:.18;
  actor.position.x+=Math.sin(spot.yaw)*forward;actor.position.z+=Math.cos(spot.yaw)*forward;
  const bend=1.05;
  for(const side of [-1,1]){
   const hip=actor.getObjectByName(`leg${side}`),knee=actor.getObjectByName(`knee${side}`);
   if(hip)hip.rotation.x=-bend;if(knee)knee.rotation.x=bend;
   aimArm(actor,side,new THREE.Vector3(side*.18,.95,.05));
  }
 }
 actor.updateMatrixWorld(true);
}
