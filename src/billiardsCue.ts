import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';

export function createBilliardsCue(){
 const cue=new THREE.Group();cue.name='billiards-cue';
 const materials:THREE.MeshStandardMaterial[]=[];
 const part=(length:number,tip:number,butt:number,centre:number,color:string)=>{
  const material=new THREE.MeshStandardMaterial({color,roughness:.45,transparent:true});materials.push(material);
  const mesh=new THREE.Mesh(new THREE.CylinderGeometry(tip,butt,length,20),material);mesh.position.y=centre;mesh.castShadow=true;cue.add(mesh);
 };
 part(1.08,.005,.010,.15,'#c4a46e');part(.32,.010,.014,-.55,'#453124');part(.016,.005,.005,.698,'#e9dfbb');part(.008,.005,.005,.710,'#668f91');
 return {cue,materials};
}

// Upright waiting pose. The rigid sleeve solver places the palm on the shaft;
// this is separate from the bending/bridge pose required for a live stroke.
export function holdBilliardsCue(actor:THREE.Group,cue:THREE.Group){
 const grip=new THREE.Vector3(.34,1.02,.12);
 cue.rotation.set(0,0,-.08);
 const direction=new THREE.Vector3(0,1,0).applyQuaternion(cue.quaternion);
 cue.position.copy(grip).addScaledVector(direction,-.23);
 actor.add(cue);
 aimArm(actor,1,grip);
 actor.updateMatrixWorld(true);
}
