import * as THREE from 'three';
import {aimArm,poseCustody} from './city3dWeapons.js';
const smooth=(v:number)=>{const t=Math.max(0,Math.min(1,v));return t*t*(3-2*t);};
/** The officer approaches from behind; both cast members share one swept reserve. */
export class CityCustody {
 readonly root=new THREE.Group();
 constructor(readonly detainee:THREE.Group,readonly officer:THREE.Group){
  this.root.add(detainee,officer);detainee.rotation.y=officer.rotation.y=Math.PI/2;this.update(0);
 }
 update(seconds:number){
  const distance=Math.min(2.42,Math.max(0,seconds-.35)*1.25),walking=seconds>.35&&distance<2.42;
  const phase=distance/1.15*Math.PI*2;
  this.detainee.position.set(3,0,0);poseCustody(this.detainee,seconds-1.3);
  this.officer.position.set(distance,walking?Math.abs(Math.sin(phase))*.018:0,0);
  for(const name of ['leg1','leg-1','knee1','knee-1']){
   const limb=this.officer.getObjectByName(name);if(!limb)continue;
   const p=phase+(name.endsWith('-1')?0:Math.PI);
   limb.rotation.x=!walking?0:name.startsWith('knee')?Math.max(0,Math.sin(p+.7))*.65:Math.sin(p)*.35;
  }
  const reach=smooth((seconds-2.25)/.55);
  for(const side of [-1,1])aimArm(this.officer,side,new THREE.Vector3(side*(.31+(.085-.31)*reach),.77+.21*reach,.34*reach));
  this.root.updateMatrixWorld(true);
  return {distance,walking,restrained:seconds>=2.8};
 }
}
