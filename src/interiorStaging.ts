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

export function garagePlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[1.3,2.8,4.3].flatMap((z,row)=>[-1.1,.5,2.1].map((x,col)=>({id:`garage-customer-${row}-${col}`,x,z,yaw:Math.PI})));
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(/\b(mechanic|garage hand)\b/i.test(who.role||'')&&![...result.values()].some(s=>s.id==='garage-mechanic'))
   result.set(who.id,{id:'garage-mechanic',x:.7,z:-2.5,yaw:Math.PI/2});
  else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

export function restaurantPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const available:InteriorSpot[]=[
  ...[-2.7,2.7].flatMap((x,col)=>[-1.5,1.5].flatMap((z,row)=>[-1,1].map(side=>({id:`dining-${col}-${row}-${side}`,x,z:z+side*.98,yaw:side<0?0:Math.PI,seat:.69})))),
  {id:'dining-wait-a',x:-1.25,z:3.8,yaw:Math.PI},
  {id:'dining-wait-b',x:1.25,z:3.8,yaw:Math.PI},
 ];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const role=who.role||'';
  if(/\b(waiter|waitress|restaurateur|cook|chef|cellarman)\b/i.test(role)){
   const stations:InteriorSpot[]=/\b(cook|chef)\b/i.test(role)
    ? [{id:'dining-kitchen',x:3.05,z:-3.7,yaw:Math.PI}]
    : /\bcellarman\b/i.test(role)
    ? [{id:'dining-wine',x:-3.1,z:-3.65,yaw:Math.PI}]
    : [{id:'dining-service',x:-.65,z:-3.7,yaw:0},{id:'dining-service-2',x:.65,z:-3.7,yaw:0}];
   const spot=stations.find(s=>![...result.values()].some(used=>used.id===s.id));
   if(spot)result.set(who.id,spot);
  } else {const spot=available.shift();if(spot)result.set(who.id,spot);}
 }
 return result;
}

export function exchangePlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[-1.8,.6,3].map((x,i)=>({id:`exchange-clerk-${i}`,x,z:-3.7,yaw:0}));
 const available:InteriorSpot[]=[
  ...[-3.6,-1.4].flatMap((x,col)=>[-.3,1.3].map((z,row)=>({id:`exchange-reader-${col}-${row}`,x,z,yaw:col===0?Math.PI/2:-Math.PI/2,seat:.69}))),
  ...[-.8,.8].flatMap((z,row)=>[.3,1.9,3.5].map((x,col)=>({id:`exchange-floor-${row}-${col}`,x,z,yaw:Math.PI}))),
 ];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const spot=/\b(clerk|broker|teller|exchange attendant)\b/i.test(who.role||'')?staff.shift():available.shift();
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export function precinctPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[{id:'station-sergeant',x:-2.1,z:-3.2,yaw:0},{id:'station-report',x:2.8,z:-2.65,yaw:0}];
 const available:InteriorSpot[]=[
  ...[-.6,.8,2.2].map((z,i)=>({id:`station-bench-${i}`,x:-4.2,z,yaw:Math.PI/2,seat:.69})),
  ...[0,1.6].flatMap((z,row)=>[-2.5,-.8,1,2.8].map((x,col)=>({id:`station-floor-${row}-${col}`,x,z,yaw:Math.PI}))),
 ];
 const order=[...people].sort((a,b)=>(/sergeant/i.test(a.role||'')?0:1)-(/sergeant/i.test(b.role||'')?0:1)||a.id.localeCompare(b.id));
 for(const who of order){
  if(result.has(who.id))continue;
  const spot=/\b(sergeant|police officer|constable|detective|commissioner)\b/i.test(who.role||'')?staff.shift():available.shift();
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export function cabstandPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();let dispatcher=false;
 const seats:InteriorSpot[]=[-.6,.8].map(z=>({id:`cab-bench-${z}`,x:3.85,z,yaw:-Math.PI/2,seat:.59}));
 const spots:InteriorSpot[]=[...seats,...[-1,.8,2.6].flatMap(z=>[-1.3,1.3].map(x=>({id:`cab-waiting-${x}-${z}`,x,z,yaw:x<0?Math.PI/2:-Math.PI/2})))];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(!dispatcher&&/dispatcher/i.test(who.role||'')){result.set(who.id,{id:'cab-dispatch',x:0,z:-3.8,yaw:0});dispatcher=true;continue;}
  const spot=spots.shift();if(spot)result.set(who.id,spot);
 }
 return result;
}

export function docksPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();let clerk=false;
 const spots:InteriorSpot[]=[-3,-1,1,3].flatMap(z=>[-1.3,1.3].map(x=>({id:`dock-aisle-${x}-${z}`,x,z,yaw:x<0?Math.PI/2:-Math.PI/2})));
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(!clerk&&/clerk|foreman/i.test(who.role||'')){result.set(who.id,{id:'dock-dispatch',x:-4,z:4.75,yaw:Math.PI});clerk=true;continue;}
  const spot=spots.shift();if(spot)result.set(who.id,spot);
 }
 return result;
}

export function chapelPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const seats:InteriorSpot[]=[1.1,-.4].flatMap(y=>[-2.65,-1.65,1.65,2.65].map(x=>({id:`chapel-chair-${x}-${y}`,x,z:-y,yaw:Math.PI,seat:.59})));
 let desk=false;
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  if(!desk&&/undertaker|funeral|mortician/i.test(who.role||'')){result.set(who.id,{id:'chapel-reception',x:-3.4,z:4.75,yaw:Math.PI});desk=true;continue;}
  const seat=seats.shift();if(seat)result.set(who.id,seat);
 }
 return result;
}

export function heraldPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[-2.55,2.55].map((x,i)=>({id:`herald-desk-${i}`,x,z:-2.82,yaw:0,seat:.59}));
 const visitors:InteriorSpot[]=[{id:'herald-aisle-0',x:0,z:-1,yaw:Math.PI/2},{id:'herald-aisle-1',x:0,z:.6,yaw:-Math.PI/2},...[-2.55,0,2.55].map((x,i)=>({id:`herald-visitor-${i}`,x,z:2.2,yaw:Math.PI}))];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const spot=/editor|reporter|journalist|copywriter/i.test(who.role||'')?staff.shift()??visitors.shift():visitors.shift();
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export function pawnPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[-2.2,0,2.2].map((x,i)=>({id:`pawn-staff-${i}`,x,z:-2.9,yaw:0}));
 const available:InteriorSpot[]=[0,1.65].flatMap((z,row)=>[-1.8,0,1.8].map((x,col)=>({id:`pawn-customer-${row}-${col}`,x,z,yaw:Math.PI})));
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const role=who.role||'';
  const station=/\bcounter clerk\b/i.test(role)?2:/\bvaluer\b/i.test(role)?1:/\bpawnbroker\b/i.test(role)?0:-1;
  const spot=station<0?available.shift():[...result.values()].some(s=>s.id===staff[station].id)?undefined:staff[station];
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export function poolhallPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[{id:'pool-marker',x:3.3,z:-7.85,yaw:0},{id:'pool-table-hand',x:-1.4,z:-6.4,yaw:0}];
 const available:InteriorSpot[]=[...[-4.25,0,4.25].map(z=>({id:`pool-spectator-${z}`,x:-7.2,z,yaw:Math.PI/2,seat:.69})),...[-4.25,0,4.25].flatMap(z=>[-1.7,1.7].map(x=>({id:`pool-table-side-${x}-${z}`,x,z,yaw:x<0?-Math.PI/2:Math.PI/2})))];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const station=/^marker$/i.test(who.role||'')?0:/^table hand$/i.test(who.role||'')?1:-1;
  const spot=station<0?available.shift():[...result.values()].some(s=>s.id===staff[station].id)?undefined:staff[station];
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export function tailorPlacements(people:Presence[]) {
 const result=new Map<string,InteriorSpot>();
 const staff:InteriorSpot[]=[-1.5,0,1.5].map((x,i)=>({id:`tailor-staff-${i}`,x,z:-3.8,yaw:0}));
 const available:InteriorSpot[]=[{id:'tailor-client-chair',x:-3.8,z:3.2,yaw:Math.PI/2,seat:.69},...[.2,1.9].flatMap((z,row)=>[-2.1,-.5,1.1].map((x,col)=>({id:`tailor-customer-${row}-${col}`,x,z,yaw:Math.PI})))];
 for(const who of [...people].sort((a,b)=>a.id.localeCompare(b.id))){
  if(result.has(who.id))continue;
  const role=who.role||'';
  const station=/\bfitting clerk\b/i.test(role)?2:/\bfinisher\b/i.test(role)?1:/\bcutter\b/i.test(role)?0:-1;
  const spot=station<0?available.shift():[...result.values()].some(s=>s.id===staff[station].id)?undefined:staff[station];
  if(spot)result.set(who.id,spot);
 }
 return result;
}

export type InteriorPlace='casino'|'cemetery'|'crematorium'|'mortuary'|'cabstand'|'docks'|'chapel'|'herald'|'poolhall'|'tailor'|'pawn'|'riverside'|'bar'|'mercercourt'|'room'|'laundry'|'estate'|'apartment'|'flat'|'butcher'|'garage'|'lodging'|'restaurant'|'market'|'precinct';
export function placementsForInterior(place:InteriorPlace,people:Presence[]){
 if(place==='casino'){
  const spots:InteriorSpot[]=[{id:'cashier',x:2.5,z:-5.31,yaw:0},
   ...[-4.8,-3.6,-2.4].map((x,i)=>({id:`slot-seat-${i}`,x,z:-3.3,yaw:Math.PI,seat:.725})),
   ...[[-2.7,-1.8],[2.7,-.5],[-2.7,2.8]].flatMap(([x,z],i)=>[-1,1].map(side=>(i===0&&side<0?{id:'roulette-side',x:-4.6,z:-1.8,yaw:Math.PI/2}:{id:`table-${i}-${side}`,x,z:z+side*1.5,yaw:side>0?Math.PI:0}))),
   {id:'lounge-a',x:5.15,z:3.7,yaw:-Math.PI/2,seat:.65},{id:'lounge-b',x:5.15,z:2.3,yaw:-Math.PI/2,seat:.65}];
  return new Map([...people].sort((a,b)=>Number(/cashier|croupier/i.test(b.role||''))-Number(/cashier|croupier/i.test(a.role||''))||a.id.localeCompare(b.id)).slice(0,spots.length).map((p,i)=>[p.id,spots[i]]));
 }

 if(place==='mortuary'||place==='cemetery'||place==='crematorium'){
  const spots:InteriorSpot[]=[{id:'registry-clerk',x:3.1,z:-1.5,yaw:Math.PI},...[-1,1,2.7].flatMap((z,row)=>[-.8,.9].map((x,col)=>({id:`receiving-${row}-${col}`,x,z,yaw:Math.PI})))];
  if(place!=='mortuary')spots.push(...[-.3,1.3,3].map((z,i)=>({id:`memorial-aisle-${i}`,x:-2.5,z,yaw:Math.PI})),{id:'memorial-front',x:-.8,z:4.2,yaw:Math.PI});
  if(!people.some(p=>/attendant|receptionist|clerk/i.test(p.role||'')))spots.shift();
  return new Map([...people].sort((a,b)=>Number(/attendant|receptionist|clerk/i.test(b.role||''))-Number(/attendant|receptionist|clerk/i.test(a.role||''))||a.id.localeCompare(b.id)).slice(0,spots.length).map((p,i)=>[p.id,spots[i]]));
 }

 if(place==='riverside'){
  const spots:InteriorSpot[]=[...[-.6,.8,2.2].map(z=>({id:'bench',x:-4.2,z,yaw:Math.PI/2,seat:.69})),...[0,1.6].flatMap(z=>[-2.5,-.8,1,2.8].map(x=>({id:'lobby',x,z,yaw:Math.PI})))];
  return new Map([...people].sort((a,b)=>a.id.localeCompare(b.id)).slice(0,spots.length).map((p,i)=>[p.id,spots[i]]));
 }
 if(place==='flat'||place==='lodging')return new Map<string,InteriorSpot>();
 return (place==='cabstand'?cabstandPlacements:place==='docks'?docksPlacements:place==='chapel'?chapelPlacements:place==='herald'?heraldPlacements:place==='poolhall'?poolhallPlacements:place==='tailor'?tailorPlacements:place==='pawn'?pawnPlacements:place==='precinct'?precinctPlacements:place==='market'?exchangePlacements:place==='restaurant'?restaurantPlacements:place==='garage'?garagePlacements:place==='butcher'?butcherPlacements:place==='apartment'?ashburyPlacements:place==='estate'?cypressPlacements:place==='laundry'?laundryPlacements:place==='bar'?interiorPlacements:place==='room'?marinerLobbyPlacements:mercerLobbyPlacements)(people);
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
 if(place==='casino')return {id:"player-entry",x:0,z:4.8,yaw:Math.PI};
 if(place==='mortuary'||place==='cemetery'||place==='crematorium')return {id:"player-entry",x:1.7,z:4,yaw:Math.PI};
 if(place==='cabstand')return {id:'player-entry',x:0,z:4,yaw:Math.PI};
 if(place==='docks')return {id:'player-entry',x:0,z:4.5,yaw:Math.PI};
 if(place==='chapel')return {id:'player-entry',x:0,z:3.6,yaw:Math.PI};
 if(place==='herald')return {id:'player-entry',x:0,z:4.1,yaw:Math.PI};
 if(place==='poolhall')return {id:'player-entry',x:0,z:6.9,yaw:Math.PI};
 if(place==='tailor')return {id:'player-entry',x:0,z:4.1,yaw:Math.PI};
 if(place==='pawn')return {id:'player-entry',x:1.4,z:3.6,yaw:Math.PI};
 if(place==='precinct'||place==='riverside')return {id:'player-entry',x:1.4,z:4.1,yaw:Math.PI};
 if(place==='market')return {id:'player-entry',x:.2,z:3.8,yaw:Math.PI};
 if(place==='restaurant')return {id:'player-entry',x:0,z:4.1,yaw:Math.PI};
 if(place==='lodging')return {id:'player-entry',x:-1.8,z:1.8,yaw:Math.PI/4};
 if(place==='garage')return {id:'player-entry',x:-3,z:4.25,yaw:Math.PI};
 if(place==='butcher')return {id:'player-entry',x:2.7,z:3.5,yaw:Math.PI};
 if(place==='flat')return {id:'player-entry',x:1,z:2.7,yaw:Math.PI};
 if(place==='apartment')return {id:'player-entry',x:1.6,z:4.15,yaw:Math.PI};
 if(place==='estate')return {id:'player-entry',x:1.6,z:3.3,yaw:Math.PI};
 if(place==='room'||place==='laundry')return {id:'player-entry',x:1.6,z:3.15,yaw:Math.PI};
 return place==='mercercourt'
  ? {id:'player-entry',x:2,z:4.15,yaw:Math.PI}
  : {id:'player-entry',x:3.65,z:-.7,yaw:-Math.PI/2};
}
