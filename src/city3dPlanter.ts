import * as THREE from 'three';
import {incendiaryStride} from './city3dIncendiary.js';

export const PLANTER_BLAST=6.2;
const smooth=(n:number)=>{const t=THREE.MathUtils.clamp(n,0,1);return t*t*(3-2*t);};
/** Doorway-local exit: +Z is inside, -Z is the pavement. Caller owns the door. */
export class CityPlanter {
 readonly root=new THREE.Group();
 constructor(readonly actor:THREE.Group){this.root.add(actor);this.update(0);}
 update(seconds:number){
  // Stand behind the door's entire swing until the opening is clear.
  const exit=incendiaryStride(seconds-.5,4.6,3.4,.25);
  const leave=incendiaryStride(seconds-4.3,2,1.5,.25);
  const distance=exit.distance+leave.distance,weight=Math.max(exit.weight,leave.weight);
  const phase=distance/1.15*Math.PI*2;
  this.actor.position.set(-leave.distance,.02+Math.abs(Math.sin(phase))*.025*weight,2.25-exit.distance);
  this.actor.rotation.set(0,Math.PI+Math.PI/2*smooth((seconds-3.9)/.4),0);
  for(const side of [-1,1]){
   const p=phase+(side<0?0:Math.PI);
   for(const [name,angle] of [[`leg${side}`,Math.sin(p)*.35*weight],[`knee${side}`,Math.max(0,Math.sin(p+.7))*.65*weight],[`arm${side}`,-Math.sin(p)*.23*weight],[`elbow${side}`,0]] as const)
    this.actor.getObjectByName(name)?.rotation.set(angle,0,0);
  }
  this.root.updateMatrixWorld(true);
  // Do not close through the person crossing the threshold.
  const door=smooth(seconds/.4)*(1-smooth((seconds-3.3)/.5));
  return {door,blast:seconds>=PLANTER_BLAST,distance,weight};
 }
}
