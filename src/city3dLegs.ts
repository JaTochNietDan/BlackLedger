import * as THREE from 'three';

/** Rigid two-segment legs with independent ankles for level stair contact. */
export class CityLegs {
 private readonly limbs=new Map<number,{hip:THREE.Object3D;knee:THREE.Object3D;ankle:THREE.Group;upper:number;lower:number}>();
 constructor(readonly actor:THREE.Group){
  actor.updateMatrixWorld(true);
  for(const side of [-1,1]){
   const hip=actor.getObjectByName(`leg${side}`),knee=actor.getObjectByName(`knee${side}`);
   const shoe=knee?.children.find(o=>o instanceof THREE.Mesh&&o.name.startsWith('shoe'));
   if(!hip||!knee||!shoe)continue;
   const ankle=new THREE.Group();ankle.name=`ankle${side}`;ankle.position.set(0,-.37,0);
   knee.add(ankle);actor.updateMatrixWorld(true);ankle.attach(shoe);
   this.limbs.set(side,{hip,knee,ankle,upper:knee.position.length(),lower:.37});
  }
 }
 /** Ankle target in actor-local metres; returns false for an unreachable target. */
 place(side:number,target:THREE.Vector3){
  const limb=this.limbs.get(side);if(!limb)return false;
  const {hip,knee,ankle,upper,lower}=limb;
  const delta=target.clone().sub(hip.position),length=delta.length();
  if(!Number.isFinite(length)||length<Math.abs(upper-lower)+1e-5||length>upper+lower-1e-5)return false;
  const direction=delta.multiplyScalar(1/length),bend=new THREE.Vector3(0,0,1);
  bend.addScaledVector(direction,-bend.dot(direction));
  if(bend.lengthSq()<1e-8)return false;
  bend.normalize();
  const along=(upper*upper-lower*lower+length*length)/(2*length);
  const joint=hip.position.clone().addScaledVector(direction,along).addScaledVector(bend,Math.sqrt(Math.max(0,upper*upper-along*along)));
  const down=new THREE.Vector3(0,-1,0);
  hip.quaternion.setFromUnitVectors(down,joint.clone().sub(hip.position).normalize());
  knee.quaternion.setFromUnitVectors(down,target.clone().sub(joint).applyQuaternion(hip.quaternion.clone().invert()).normalize());
  ankle.quaternion.copy(hip.quaternion).multiply(knee.quaternion).invert();
  this.actor.updateMatrixWorld(true);
  return true;
 }
}
