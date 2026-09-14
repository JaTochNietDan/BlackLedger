import * as THREE from 'three';
import {interiorPlayerSpot,type InteriorPlace} from './interiorStaging.js';
const smooth=(v:number)=>{const t=THREE.MathUtils.clamp(v,0,1);return t*t*(3-2*t);};

/** Cosmetic entrance along the reserved foreground aisle. It never advances
 * game time, delays actions or changes the player's authoritative location. */
export class InteriorArrival {
 readonly duration:number;
 private readonly end;
 private readonly start:THREE.Vector3;
 private readonly distance:number;
 private readonly heading:number;
 private readonly limbs=new Map<string,THREE.Object3D>();
 private readonly shoes:THREE.Mesh[]=[];
 private readonly box=new THREE.Box3();
 private readonly shoeBox=new THREE.Box3();
 constructor(private actor:THREE.Group,private place:InteriorPlace){
  this.end=interiorPlayerSpot(place);
  this.start=place==='poolhall'?new THREE.Vector3(0,.03,8.5):place==='pawn'?new THREE.Vector3(0,.03,4.15):place==='lodging'?new THREE.Vector3(-2.45,.03,2.45):place==='garage'?new THREE.Vector3(-3,.03,5.02):place==='butcher'?new THREE.Vector3(2.7,.03,4.12):place==='flat'?new THREE.Vector3(0,.03,3.3):place==='estate'?new THREE.Vector3(0,.03,4.3):place==='bar'?new THREE.Vector3(5.25,.03,-.7):(place==='room'||place==='laundry')?new THREE.Vector3(0,.03,3.6):new THREE.Vector3(0,.03,4.6);
  this.distance=Math.hypot(this.end.x-this.start.x,this.end.z-this.start.z);
  this.duration=this.distance/1.15+.4;
  this.heading=Math.atan2(this.end.x-this.start.x,this.end.z-this.start.z);
  actor.traverse(o=>{if(/^(leg|knee|arm|elbow)-?1$/.test(o.name))this.limbs.set(o.name,o);if(o instanceof THREE.Mesh&&o.name.startsWith('shoe'))this.shoes.push(o);});
 }
 pose(seconds:number){
  const progress=smooth(seconds/this.duration),travel=this.distance*progress;
  const strength=smooth(travel/.22)*smooth((this.distance-travel)/.22);
  this.actor.position.set(THREE.MathUtils.lerp(this.start.x,this.end.x,progress),.03,THREE.MathUtils.lerp(this.start.z,this.end.z,progress));
  const turn=Math.atan2(Math.sin(this.end.yaw-this.heading),Math.cos(this.end.yaw-this.heading));
  this.actor.rotation.y=this.heading+turn*smooth((progress-.75)/.25);
  for(const [name,limb] of this.limbs){
   const phase=travel/1.15*Math.PI*2+(name.endsWith('-1')?0:Math.PI);
   limb.rotation.x=strength*(name.startsWith('elbow')?0:name.startsWith('knee')?Math.max(0,Math.sin(phase+.7))*.65:Math.sin(phase+(name.startsWith('arm')?Math.PI:0))*(name.startsWith('arm')?.23:.35));
  }
  // Fit the real soles to the floor or entrance mat. The rig's rigid leg
  // segments retain their length, and neither shoe can fall through the floor.
  this.actor.updateMatrixWorld(true);this.box.makeEmpty();
  for(const shoe of this.shoes)this.box.union(this.shoeBox.setFromObject(shoe,true));
  if(!this.box.isEmpty()){
   const mat=this.place==='room'?{x:1.1,z0:2.225,z1:3.575,top:.0555}:this.place==='mercercourt'?{x:1.4,z0:2.65,z1:4.15,top:.0575}:undefined;
   const onMat=mat&&this.box.min.x<mat.x&&this.box.max.x> -mat.x&&this.box.min.z<mat.z1&&this.box.max.z>mat.z0;
   const ground=onMat?mat.top:.0175;
   this.actor.position.y+=ground+.004-this.box.min.y;
  }
  this.actor.updateMatrixWorld(true);
  return progress<1;
 }
}
