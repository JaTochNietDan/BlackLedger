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
