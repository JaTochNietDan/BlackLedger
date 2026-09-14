import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';

/** A public bartender wipes the clear rear strip of Saint Agnes's marble bar.
 * Contact follows the articulated hand, not an unrelated floating prop track.
 */
export class CounterWipe {
 readonly clothPosition=new THREE.Vector3();
 private hand?:THREE.Mesh;
 private box=new THREE.Box3();
 private target=new THREE.Vector3();
 constructor(private actor:THREE.Group){
  actor.getObjectByName('elbow1')?.traverse(o=>{if(o instanceof THREE.Mesh&&o.name.startsWith('hand'))this.hand=o;});
 }
 get available(){return !!this.hand;}
 pose(seconds:number){
  if(!this.hand)return;
  this.target.set(.26+.14*Math.sin(seconds*1.4),1.22,.36+.025*Math.sin(seconds*2.8));
  // Align the underside of the real hand with the cloth, accounting for the
  // forearm's changing angle. Do not shorten the rigid arm segments.
  for(let i=0;i<5;i++){
   aimArm(this.actor,1,this.target);this.actor.updateMatrixWorld(true);
   this.box.setFromObject(this.hand,true);
   this.target.y+=1.201-this.box.min.y;
  }
  aimArm(this.actor,1,this.target);this.actor.updateMatrixWorld(true);this.box.setFromObject(this.hand,true);
  this.box.getCenter(this.clothPosition);this.clothPosition.y=1.19;
 }
}

/** A laundry worker smooths linen on the rear edge of Bluebird's counter. */
export class LinenPress {
 seconds=0;
 active=false;
 private hands=new Map<number,THREE.Mesh>();
 private restPose=new Map<THREE.Object3D,THREE.Quaternion>();
 private box=new THREE.Box3();
 private target=new THREE.Vector3();
 constructor(private actor:THREE.Group){
  for(const side of [-1,1]){
   actor.getObjectByName(`elbow${side}`)?.traverse(o=>{if(o instanceof THREE.Mesh&&o.name.startsWith('hand'))this.hands.set(side,o);});
   for(const part of ['arm','elbow']){const joint=actor.getObjectByName(part+side);if(joint)this.restPose.set(joint,joint.quaternion.clone());}
  }
 }
 get available(){return this.hands.size===2;}
 step(delta:number,working:boolean,motion:boolean,hidden:boolean){
  if(!this.available)return false;
  if(!working){
   if(!this.active)return false;
   this.restPose.forEach((q,joint)=>joint.quaternion.copy(q));this.actor.updateMatrixWorld(true);this.active=false;return true;
  }
  const changed=!this.active||(motion&&!hidden&&Number.isFinite(delta)&&delta>0);
  if(!changed)return false;
  if(motion&&!hidden&&Number.isFinite(delta)&&delta>0)this.seconds+=Math.min(.05,delta);
  this.active=true;
  const spread=.16+.05*Math.sin(this.seconds*1.35),forward=.49+.015*Math.sin(this.seconds*2.7);
  for(const [side,hand] of this.hands){
   this.target.set(side*spread,1.17,forward);
   // Contact is solved against the exported hand's underside, while keeping
   // each sleeve segment rigid and the feet at their assigned floor spot.
   for(let i=0;i<6;i++){
    aimArm(this.actor,side,this.target,new THREE.Vector3(side*.8,.5,-.3));this.actor.updateMatrixWorld(true);this.box.setFromObject(hand,true);
    this.target.y+=1.109-this.box.min.y;
   }
   aimArm(this.actor,side,this.target,new THREE.Vector3(side*.8,.5,-.3));
  }
  this.actor.updateMatrixWorld(true);
  return true;
 }
}
