import * as THREE from 'three';

export function sceneWeapon(tier?:number):'revolver'|'shotgun'|'thompson'|null {
 if(tier===undefined)return 'revolver'; // Unattributed legacy gunfire.
 return tier===1?'revolver':tier===2?'shotgun':tier===3?'thompson':null;
}

/** Two rigid sleeve segments reach the grip; arm lengths remain unchanged. */
export function aimArm(actor:THREE.Group,side:number,target:THREE.Vector3){
 const arm=actor.getObjectByName(`arm${side}`),elbow=actor.getObjectByName(`elbow${side}`);if(!arm||!elbow)return;
 const origin=arm.position,delta=target.clone().sub(origin),distance=delta.length();
 const upper=.285,lower=.295,d=Math.max(.011,Math.min(upper+lower-.001,distance));
 const direction=delta.normalize();
 const along=(upper*upper-lower*lower+d*d)/(2*d);
 const bend=new THREE.Vector3(side*.5,-1,-.4);bend.addScaledVector(direction,-bend.dot(direction)).normalize();
 const joint=origin.clone().addScaledVector(direction,along).addScaledVector(bend,Math.sqrt(Math.max(0,upper*upper-along*along)));
 arm.quaternion.setFromUnitVectors(new THREE.Vector3(0,-1,0),joint.clone().sub(origin).normalize());
 const forearm=origin.clone().addScaledVector(direction,d).sub(joint).applyQuaternion(arm.quaternion.clone().invert()).normalize();
 elbow.quaternion.setFromUnitVectors(new THREE.Vector3(0,-1,0),forearm);
}

export function poseLongGun(actor:THREE.Group,weapon:THREE.Group,armAngle:number){
 const aim=Math.max(0,Math.min(1,-armAngle/(Math.PI/2))),recoil=Math.max(0,-armAngle-Math.PI/2);
 weapon.position.set(0,1.12+.08*aim,.10+.04*aim-recoil*.06);
 weapon.rotation.set(.45*(1-aim)-recoil*.3,0,0);
 actor.updateMatrixWorld(true);
 const grip=actor.worldToLocal(weapon.getWorldPosition(new THREE.Vector3()));
 const support=weapon.getObjectByName('support-grip');
 if(support){const target=actor.worldToLocal(support.getWorldPosition(new THREE.Vector3()));aimArm(actor,-1,target);}
 aimArm(actor,1,grip);actor.updateMatrixWorld(true);
}

const revolverShots=[.7,1.05,1.5,1.9] as const;
const shotgunShots=[.7,1.7] as const;
const thompsonShots=[.7,.79,.88,1.5,1.59,1.68] as const;
export function weaponShots(model?:string):readonly number[]{return model==='shotgun'?shotgunShots:model==='thompson'?thompsonShots:revolverShots;}
export function pumpOffset(seconds:number){
 const age=seconds-(seconds>=1.7?1.7:.7);
 return age>.15&&age<.70?-.095*Math.sin((age-.15)/.55*Math.PI):0;
}
