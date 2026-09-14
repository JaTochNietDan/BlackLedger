import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';
import {casualtyFall} from './city3dEvents.js';
import type {VisualCue} from './types';

export const ASSASSINATION_SHOT = 3.65;
export const ASSASSINATION_SECONDS = 6.5;
export const ASSASSINATION_VICTIM_X = 5;
export function isExecution(cue: VisualCue) {
  return cue.kind==='gunfight' && cue.strike?.variant==='back-of-head' && cue.attacker?.weapon===1;
}
/** Collapse only an explicitly linked victim cue; unrelated deaths remain visible. */
export function assassinationBatch(cues: VisualCue[]) {
  return cues.filter(cue => !(cue.kind==='killing' && cues.some(gun => isExecution(gun) &&
    gun.target===cue.target && gun.minute===cue.minute &&
    cue.actors?.some(actor=>actor.id===gun.strike!.victim.id))));
}
const ease=(v:number)=>{const t=Math.max(0,Math.min(1,v));return t*t*(3-2*t);};
export function assassinationPose(seconds:number) {
  const distance=Math.min(4.07,Math.max(0,seconds-.35)*1.35);
  const walking=seconds>.35&&distance<4.07;
  const age=seconds-ASSASSINATION_SHOT;
  const aim=ease((seconds-2.85)/.6)*(1-ease((seconds-4.5)/.65));
  const recoil=age>=0?Math.max(0,1-age/.16):0;
  return {distance,walking,phase:distance/1.15*Math.PI*2,aim,recoil,
    fall:casualtyFall(Math.max(0,age-.06)/3),age};
}
/** One shared world-aligned cast, facing +X. Its full path is reserved before playback. */
export class CityAssassination {
  readonly root=new THREE.Group();
  constructor(readonly attacker:THREE.Group,readonly victim:THREE.Group,readonly weapon:THREE.Group) {
    this.root.add(attacker,victim);attacker.add(weapon);
    attacker.rotation.set(0,Math.PI/2,0);
    this.update(0);
  }
  update(seconds:number) {
    const p=assassinationPose(seconds);
    this.attacker.position.set(p.distance,p.walking?Math.abs(Math.sin(p.phase))*.018:0,0);
    for(const name of ['leg1','leg-1','knee1','knee-1','arm-1']) {
      const limb=this.attacker.getObjectByName(name);if(!limb)continue;
      const phase=p.phase+(name.endsWith('-1')?0:Math.PI);
      limb.rotation.x=!p.walking?0:name.startsWith('knee')?Math.max(0,Math.sin(phase+.7))*.65
        :Math.sin(phase+(name.startsWith('arm')?Math.PI:0))*(name.startsWith('arm')?.23:.35);
    }
    this.weapon.position.set(.3+(.07-.3)*p.aim,.77+(1.564-.77)*p.aim,.47*p.aim-.035*p.recoil);
    this.weapon.rotation.set(Math.PI/2*(1-p.aim)-.05*p.recoil,0,0);
    aimArm(this.attacker,1,this.weapon.position);
    this.victim.position.set(ASSASSINATION_VICTIM_X,p.fall.height-.2,0);
    // World-space fall preserves the forward direction after starting with back to shooter.
    this.victim.quaternion.setFromAxisAngle(new THREE.Vector3(0,1,0),Math.PI/2)
      .premultiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,0,1),p.fall.rotation));
    this.root.updateMatrixWorld(true);
    return p;
  }
}
/** Small forward droplets at impact; no invented damage or further shots. */
export function executionSpatter(index:number,seconds:number) {
  const age=seconds-ASSASSINATION_SHOT;
  const active=age>=0&&age<.65;
  const a=Math.max(0,age),speed=2.2+(index%5)*.2;
  return {x:5.1+speed*a,y:Math.max(.005,1.65+(.15+(index%3)*.18)*a-4.9*a*a),
    z:((index%7)-3)*.14*a,size:active?(.06+(index%3)*.015)*(1-a/.65):0};
}
