import * as THREE from 'three';
import type {Presence} from './types';
import {poseInteriorOccupant} from './interiorStaging.js';
import {aimArm} from './city3dWeapons.js';
/** Presence alone does not make a patron a participant. Only a public dealer role qualifies. */
export function blackjackDealer(people:Presence[]){
 return [...people].filter(p=>/\b(croupier|dealer)\b/i.test(p.role??'')).sort((a,b)=>a.id.localeCompare(b.id))[0];
}
export function poseBlackjackDealer(actor:THREE.Group){
 actor.position.set(0,.03,-1.56);actor.rotation.set(0,0,0);
 for(const side of [-1,1])aimArm(actor,side,new THREE.Vector3(side*.22,1.05,.42));
 actor.updateMatrixWorld(true);
}

export const blackjackPlayerSeat={x:-.85,z:1.6,yaw:2.65,height:.665};
export function poseBlackjackPlayer(actor:THREE.Group){
 const s=blackjackPlayerSeat;
 poseInteriorOccupant(actor,{id:'blackjack-player',x:s.x-Math.sin(s.yaw)*.18,z:s.z-Math.cos(s.yaw)*.18,yaw:s.yaw,seat:s.height});
}
