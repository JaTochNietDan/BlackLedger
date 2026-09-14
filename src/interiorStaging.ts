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
// Mercer Court uses its own authored bench and clear lobby floor, not bar seats.
export function mercerLobbyPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[
  ...[.15,1.85].map((z,i)=>({id:`lobby-bench-${i}`,x:-4.85,z,yaw:Math.PI/2,seat:.77})),
  ...[-3.5,-1.5,.5,2.5].flatMap((z,row)=>[-2.6,-.7,1.2].map((x,col)=>({id:`lobby-floor-${row}-${col}`,x,z,yaw:row%2?Math.PI:0})))
 ];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\bsuperintendent\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='lobby-register'))
   result.set(who.id,{id:'lobby-register',x:-2,z:-5.5,yaw:0});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

// Mariner reception and its own bench, with a clear approach to the stairs.
export function marinerLobbyPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[
  ...[.05,1.75].map((z,i)=>({id:`boarding-bench-${i}`,x:-4.2,z,yaw:Math.PI/2,seat:.77})),
  ...[-2.5,-.5,1.5].flatMap((z,row)=>[-2.3,-.5,1.3].map((x,col)=>({id:`boarding-floor-${row}-${col}`,x,z,yaw:row%2?Math.PI:0})))
 ];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\b(landlady|landlord|receptionist)\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='boarding-reception'))
   result.set(who.id,{id:'boarding-reception',x:-2.5,z:-5.23,yaw:0});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

// The washers occupy the rear strip; the collection counter has its own clerk aisle.
export function laundryPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[
  ...[-.4,1.1].map((z,i)=>({id:`laundry-bench-${i}`,x:-4.2,z,yaw:Math.PI/2,seat:.77})),
  ...[-2,-.2,1.6].flatMap((z,row)=>[-2.3,-.5,1.1].map((x,col)=>({id:`laundry-floor-${row}-${col}`,x,z,yaw:row%2?Math.PI:0})))
 ];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\b(launderer|laundress|laundry worker|clerk)\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='laundry-counter'))
   result.set(who.id,{id:'laundry-counter',x:3.35,z:-2.0,yaw:0});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

// Cypress drawing-room seats match the authored sofa and facing armchairs.
export function cypressPlacements(people:Presence[]) {
 const available:InteriorSpot[]=[
  ...[-1.25,-.3,.65].map((z,i)=>({id:`cypress-sofa-${i}`,x:-3.58,z,yaw:Math.PI/2,seat:.69})),
  ...[-1.35,.8].map((z,i)=>({id:`cypress-chair-${i}`,x:.86,z,yaw:-Math.PI/2,seat:.69})),
  {id:'cypress-study',x:3.35,z:-.4,yaw:Math.PI},
  {id:'cypress-visitor',x:-1.2,z:2.8,yaw:Math.PI},
 ];
 const result=new Map<string,InteriorSpot>();
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const spot=available.shift();if(spot)result.set(who.id,spot);
 }
 return result;
}

export function ashburyPlacements(people:Presence[]) {
 const available:InteriorSpot[]=[
  ...[2.8,1.4,0].map((z,i)=>({id:`ashbury-bench-${i}`,x:-4.13,z,yaw:Math.PI/2,seat:.77})),
  ...[-.4,1.2,2.8].flatMap((z,row)=>[-2.2,-.6,1].map((x,col)=>({id:`ashbury-hall-${row}-${col}`,x,z,yaw:Math.PI}))),
 ];
 const result=new Map<string,InteriorSpot>();
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\b(concierge|superintendent|caretaker)\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='ashbury-register'))
   result.set(who.id,{id:'ashbury-register',x:-2.65,z:-3.35,yaw:0});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

export function butcherPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[-.6,.7,2].flatMap((z,row)=>[-2.5,-.7,1.1].map((x,col)=>({id:`butcher-customer-${row}-${col}`,x,z:z+1.7,yaw:Math.PI})));
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\b(butcher|shopkeeper|clerk)\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='butcher-service'))
   result.set(who.id,{id:'butcher-service',x:-.6,z:-1.8,yaw:0});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

export type InteriorPlace='bar'|'mercercourt'|'room'|'laundry'|'estate'|'apartment'|'flat'|'butcher';
export function placementsForInterior(place:InteriorPlace,people:Presence[]){
 if(place==='flat')return new Map<string,InteriorSpot>();
 return (place==='butcher'?butcherPlacements:place==='apartment'?ashburyPlacements:place==='estate'?cypressPlacements:place==='laundry'?laundryPlacements:place==='bar'?interiorPlacements:place==='room'?marinerLobbyPlacements:mercerLobbyPlacements)(people);
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

// Reserved clear floor positions; these never displace a public occupant.
export function interiorPlayerSpot(place:InteriorPlace):InteriorSpot {
 if(place==='butcher')return {id:'player-entry',x:2.7,z:3.5,yaw:Math.PI};
 if(place==='flat')return {id:'player-entry',x:1,z:2.7,yaw:Math.PI};
 if(place==='apartment')return {id:'player-entry',x:1.6,z:4.15,yaw:Math.PI};
 if(place==='estate')return {id:'player-entry',x:1.6,z:3.3,yaw:Math.PI};
 if(place==='room'||place==='laundry')return {id:'player-entry',x:1.6,z:3.15,yaw:Math.PI};
 return place==='mercercourt'
  ? {id:'player-entry',x:2,z:4.15,yaw:Math.PI}
  : {id:'player-entry',x:3.65,z:-.7,yaw:-Math.PI/2};
}
