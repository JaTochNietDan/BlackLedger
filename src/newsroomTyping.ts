import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';

/** Cosmetic work at an occupied Herald desk; it writes no simulation news. */
export class NewsroomTyping {
 seconds=0;
 private hands=new Map<number,THREE.Mesh>();
 private box=new THREE.Box3();
 constructor(private actor:THREE.Group,private phase=0){
  for(const side of [-1,1])actor.getObjectByName(`elbow${side}`)?.traverse(o=>{if(o instanceof THREE.Mesh&&o.name.startsWith('hand'))this.hands.set(side,o);});
  this.pose();
 }
 step(delta:number,motion:boolean,hidden:boolean){
  if(!motion||hidden||!Number.isFinite(delta)||delta<=0)return false;
  this.seconds+=Math.min(.05,delta);this.pose();return true;
 }
 pose(){
  const t=this.seconds+this.phase,rest=t%7>5.4;
  for(const [side,hand] of this.hands){
   // Alternating key presses, interrupted by a short reading pause. The hand
   // stays over the physical keyboard; each full sleeve keeps its own length.
   const lift=rest?.035:.008+.035*(.5+.5*Math.sin(t*15+side*Math.PI/2));
   const target=new THREE.Vector3(this.actor.position.x-.15+side*.18,.97+lift,-2.19);
   this.actor.updateMatrixWorld(true);this.actor.worldToLocal(target);
   for(let i=0;i<6;i++){
    aimArm(this.actor,side,target,new THREE.Vector3(side*.25,-.2,-1));this.actor.updateMatrixWorld(true);this.box.setFromObject(hand,true);
    target.y+=.970+lift-this.box.min.y;
   }
   aimArm(this.actor,side,target,new THREE.Vector3(side*.25,-.2,-1));
  }
  this.actor.updateMatrixWorld(true);
 }
}
