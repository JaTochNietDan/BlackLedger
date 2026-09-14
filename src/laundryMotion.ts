import * as THREE from 'three';
import type {Place} from './types';
export type LaundryOperation=Pick<Place,'trading'|'staff'|'supply'|'condition'|'trouble'>;

// The public trading figure controls visual workload, not production or time.
// No machinery runs when the public premises cannot support a wash cycle.
export function laundryRunningMachines(p?:LaundryOperation){
 if(!p||!Number.isFinite(p.trading)||(p.trading??0)<=0||(p.staff??0)<=0||(p.supply??0)<=0||p.condition<40||p.trouble)return 0;
 return Math.min(3,Math.ceil(p.trading!*3));
}
export class LaundryMotion {
 seconds=0;
 private readonly drums:{object:THREE.Object3D;base:THREE.Quaternion;axis:THREE.Vector3;seconds:number}[]=[];
 private readonly rotation=new THREE.Quaternion();
 constructor(room:THREE.Group){
  room.updateMatrixWorld(true);
  for(let i=0;i<3;i++){
   const object=room.getObjectByName(`laundry-drum-${i}`);if(!object)continue;
   const world=object.getWorldQuaternion(new THREE.Quaternion());
   this.drums.push({object,base:object.quaternion.clone(),axis:new THREE.Vector3(0,0,1).applyQuaternion(world.invert()),seconds:0});
  }
 }
 get count(){return this.drums.length;}
 step(delta:number,running:number,motion:boolean,hidden:boolean){
  if(!motion||hidden||running<=0||!Number.isFinite(delta)||delta<=0)return false;
  const elapsed=Math.min(.05,delta);let changed=false;
  this.drums.forEach((drum,i)=>{
   if(i>=running)return;
   drum.seconds+=elapsed;
   drum.object.quaternion.copy(drum.base).multiply(this.rotation.setFromAxisAngle(drum.axis,drum.seconds*(2.5+i*.17)));
   changed=true;
  });
  if(changed)this.seconds+=elapsed;
  return changed;
 }
}
