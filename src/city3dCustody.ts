import * as THREE from 'three';
import {aimArm,poseCustody} from './city3dWeapons.js';
const smooth=(v:number)=>{const t=Math.max(0,Math.min(1,v));return t*t*(3-2*t);};
function gait(actor:THREE.Group,distance:number,walking:boolean){
 const phase=distance/1.15*Math.PI*2;
 actor.position.y=walking?Math.abs(Math.sin(phase))*.018:0;
 for(const name of ['leg1','leg-1','knee1','knee-1']){
  const limb=actor.getObjectByName(name);if(!limb)continue;
  const p=phase+(name.endsWith('-1')?0:Math.PI);
  limb.rotation.x=!walking?0:name.startsWith('knee')?Math.max(0,Math.sin(p+.7))*.65:Math.sin(p)*.35;
 }
}
/** Approach, cuff, move alongside, then escort within one reserved forecourt path. */
export class CityCustody {
 readonly root=new THREE.Group();
 constructor(readonly detainee:THREE.Group,readonly officer:THREE.Group){
  this.root.add(detainee,officer);this.update(0);
 }
 update(seconds:number){
  const distance=Math.min(2.42,Math.max(0,seconds-.35)*1.25),walking=seconds>.35&&distance<2.42;
  const alongside=smooth((seconds-3.4)/1.25),arc=alongside*Math.PI/2;
  const turn=smooth((seconds-4.75)/.65);
  const escort=Math.min(2.3,Math.max(0,seconds-5.5)*1.1),escorting=seconds>5.5&&escort<2.3;
  this.detainee.position.set(3-escort,0,0);
  this.detainee.rotation.y=Math.PI/2-turn*Math.PI;
  poseCustody(this.detainee,seconds-1.3);
  gait(this.detainee,escort,escorting);
  this.officer.position.set(seconds<3.4?distance:3-.58*Math.cos(arc)-escort,0,.75*Math.sin(arc));
  const tangent=Math.atan2(.58*Math.sin(arc),.75*Math.cos(arc));
  this.officer.rotation.y=seconds<3.4?Math.PI/2*(1-smooth((seconds-3.05)/.35)):tangent-turn*Math.PI;
  gait(this.officer,distance+alongside*1.05+escort,walking||(seconds>3.4&&seconds<4.65)||escorting);
  const reach=smooth((seconds-2.25)/.55)*(1-smooth((seconds-2.95)/.35));
  for(const side of [-1,1])aimArm(this.officer,side,new THREE.Vector3(side*(.31+(.085-.31)*reach),.77+.21*reach,.34*reach));
  // The escorting hand rests on the prisoner's upper arm, following the gait.
  if(seconds>=5.15){
   this.root.updateMatrixWorld(true);
   const sleeve=this.detainee.getObjectByName('arm1');
   if(sleeve){
    const contact=this.officer.worldToLocal(sleeve.localToWorld(new THREE.Vector3(0,-.08,0)));
    const hand=new THREE.Vector3(-.31,.77,0).lerp(contact,smooth((seconds-5.15)/.35));
    aimArm(this.officer,-1,hand);
   }
  }
  this.root.updateMatrixWorld(true);
  return {distance,walking,restrained:seconds>=2.8,alongside,escort,escorting};
 }
}
