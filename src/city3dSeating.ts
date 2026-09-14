import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';
/** Position the pelvis on an authored seat while keeping original limb lengths. */
export function seatDriver(actor:THREE.Group,seat:THREE.Vector3, grips:readonly THREE.Vector3[]){
 const lean=.38;
 for(const name of ['headwear-fedora','headwear-cap']){const hat=actor.getObjectByName(name);if(hat)hat.visible=false;}
 actor.rotation.set(lean,0,0);
 actor.position.copy(seat).sub(new THREE.Vector3(0,.86,0).applyEuler(actor.rotation));
 for(const side of [-1,1]){
  const hip=actor.getObjectByName(`leg${side}`),knee=actor.getObjectByName(`knee${side}`);
  if(hip)hip.rotation.x=-Math.PI/2-lean;
  if(knee)knee.rotation.x=Math.PI/2+1.3;
  const grip=grips[side<0?0:1];
  const target=grip.clone().sub(actor.position).applyQuaternion(actor.quaternion.clone().invert());
  aimArm(actor,side,target,new THREE.Vector3(0,-1,0));
 }
 actor.updateMatrixWorld(true);
}
