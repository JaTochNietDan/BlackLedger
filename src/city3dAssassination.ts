import * as THREE from 'three';
import {groundCharacter} from './city3dGround.js';
import {aimArm,poseLongGun,weaponShots} from './city3dWeapons.js';
import {casualtyFall} from './city3dEvents.js';
import type {VisualCue} from './types';

export const ASSASSINATION_SHOT = 3.65;
export const ASSASSINATION_SECONDS = 8.5;
export const ASSASSINATION_VICTIM_X = 5;
export function isExecution(cue: VisualCue) {
  return cue.kind==='gunfight' && cue.strike?.variant==='back-of-head' && cue.attacker?.weapon===1;
}
export function isStagedStrike(cue:VisualCue) {
  return isExecution(cue) || (cue.kind==='gunfight' && ((cue.strike?.variant==='close-shot' && [1,2].includes(cue.attacker?.weapon??0)) || (cue.strike?.variant==='burst' && cue.attacker?.weapon===3))) || (cue.kind==='attack' && cue.strike?.variant==='close-quarters' && cue.attacker?.weapon===0);
}
/** Collapse only an explicitly linked victim cue; unrelated deaths remain visible. */
export function assassinationBatch(cues: VisualCue[]) {
  return cues.filter(cue => !(cue.kind==='killing' && cues.some(gun => isStagedStrike(gun) &&
    gun.target===cue.target && gun.minute===cue.minute &&
    cue.actors?.some(actor=>actor.id===gun.strike!.victim.id))));
}
const ease=(v:number)=>{const t=Math.max(0,Math.min(1,v));return t*t*(3-2*t);};
export function assassinationPose(seconds:number) {
  const approach=Math.min(4.07,Math.max(0,seconds-.35)*1.35);
  const retreat=Math.min(4.07,Math.max(0,seconds-5.3)*1.35);
  const distance=approach-retreat;
  const walking=(seconds>.35&&approach<4.07)||(seconds>5.3&&retreat<4.07);
  const age=seconds-ASSASSINATION_SHOT;
  const aim=ease((seconds-2.85)/.6)*(1-ease((seconds-4.5)/.65));
  const recoil=age>=0?Math.max(0,1-age/.16):0;
  return {distance,walking,phase:(approach+retreat)/1.15*Math.PI*2,aim,recoil,
    yaw:Math.PI/2+Math.PI*ease((seconds-4.85)/.45),
    fall:casualtyFall(Math.max(0,age-.06)/3),age};
}
export const MELEE_IMPACTS = [3.7, 4.15, 4.6] as const;
export function meleePose(seconds:number) {
  const distance=Math.min(4.32,Math.max(0,seconds-.35)*1.35);
  const walking=seconds>.35&&distance<4.32;
  const age=seconds-MELEE_IMPACTS[2];
  const punch=MELEE_IMPACTS.reduce((best,at)=>Math.max(best,
    ease((seconds-at+.18)/.18)*(1-ease((seconds-at)/.24))),0);
  return {distance,walking,phase:distance/1.15*Math.PI*2,punch,age,
    fall:casualtyFall(Math.max(0,age-.06)/3)};
}
export function armedStrikePose(seconds:number,variant:string,model:string) {
  const approach=Math.min(3.2,Math.max(0,seconds-.35)*1.35);
  const retreat=Math.min(3.2,Math.max(0,seconds-4.95)*1.35);
  const walking=(seconds>.35&&approach<3.2)||(seconds>4.95&&retreat<3.2);
  const age=seconds-ASSASSINATION_SHOT;
  const beats=weaponShots(model,variant),last=[...beats].reverse().find(at=>seconds>=at);
  const recoil=last===undefined?0:Math.max(0,1-(seconds-last)/.16);
  const aim=ease((seconds-2.85)/.6)*(1-ease((seconds-4.15)/.4));
  return {distance:approach-retreat,walking,phase:(approach+retreat)/1.15*Math.PI*2,
    aim,recoil,age,yaw:Math.PI/2+Math.PI*ease((seconds-4.5)/.45),
    fall:casualtyFall(Math.max(0,age-.06)/3)};
}
/** One shared world-aligned cast, facing +X. Its full path is reserved before playback. */
export class CityAssassination {
  readonly root=new THREE.Group();
  constructor(readonly attacker:THREE.Group,readonly victim:THREE.Group,readonly weapon?:THREE.Group,readonly variant:string=weapon?'back-of-head':'close-quarters',readonly weaponModel:string='revolver') {
    this.root.add(attacker,victim);if(weapon)attacker.add(weapon);
    attacker.rotation.set(0,Math.PI/2,0);
    this.update(0);
  }
  get duration(){return this.variant==='close-quarters'?6.5:this.variant==='close-shot'||this.variant==='burst'?8:ASSASSINATION_SECONDS;}
  get victimYaw(){return this.variant==='back-of-head'?Math.PI/2:-Math.PI/2;}
  update(seconds:number) {
    const execution=this.variant==='back-of-head';
    const armed=execution?assassinationPose(seconds):armedStrikePose(seconds,this.variant,this.weaponModel),melee=meleePose(seconds);
    const p=this.weapon?armed:melee;
    this.attacker.rotation.y=this.weapon?armed.yaw:Math.PI/2;
    this.attacker.position.set(p.distance,p.walking?Math.abs(Math.sin(p.phase))*.018:0,0);
    for(const name of ['leg1','leg-1','knee1','knee-1','arm-1']) {
      const limb=this.attacker.getObjectByName(name);if(!limb)continue;
      const phase=p.phase+(name.endsWith('-1')?0:Math.PI);
      limb.rotation.x=!p.walking?0:name.startsWith('knee')?Math.max(0,Math.sin(phase+.7))*.65
        :Math.sin(phase+(name.startsWith('arm')?Math.PI:0))*(name.startsWith('arm')?.23:.35);
    }
    if(this.weapon&&this.weaponModel!=='revolver'){
      poseLongGun(this.attacker,this.weapon,-Math.PI/2*armed.aim-.12*armed.recoil,-1.1);
      const pump=this.weapon.getObjectByName('pump-slide'),age=seconds-ASSASSINATION_SHOT;
      if(pump)pump.position.z=age>.15&&age<.70?-.095*Math.sin((age-.15)/.55*Math.PI):0;
    } else if(this.weapon){
      this.weapon.position.set(.3+(.07-.3)*armed.aim,.77+((execution?1.564:1.42)-.77)*armed.aim,.47*armed.aim-.035*armed.recoil);
      this.weapon.rotation.set(Math.PI/2*(1-armed.aim)-.05*armed.recoil,0,0);
      aimArm(this.attacker,1,this.weapon.position);
    } else if(!p.walking) {
      aimArm(this.attacker,1,new THREE.Vector3(.12,1.45,.22+.42*melee.punch));
      aimArm(this.attacker,-1,new THREE.Vector3(-.18,1.38,.22));
    }
    if(this.weapon&&!execution){
      const guard=ease((seconds-2.6)/.5)*(1-ease((seconds-3.75)/.5));
      for(const side of [-1,1])aimArm(this.victim,side,new THREE.Vector3(side*(.3-.1*guard),.85+.55*guard,.18*guard));
    }
    this.victim.position.set(ASSASSINATION_VICTIM_X,p.fall.height-.2,0);
    // Fall away from the attacker, independent of the victim’s initial facing.
    this.victim.quaternion.setFromAxisAngle(new THREE.Vector3(0,1,0),this.victimYaw)
      .premultiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,0,1),p.fall.rotation));
    groundCharacter(this.victim,.005);
    this.root.updateMatrixWorld(true);
    return p;
  }
}
/** Small forward droplets at impact; no invented damage or further shots. */
export function executionSpatter(index:number,seconds:number,head=true) {
  const age=seconds-ASSASSINATION_SHOT;
  const active=age>=0&&age<.65;
  const a=Math.max(0,age),speed=2.2+(index%5)*.2;
  return {x:5.1+speed*a,y:Math.max(.005,(head?1.65:1.2)+(.15+(index%3)*.18)*a-4.9*a*a),
    z:((index%7)-3)*.14*a,size:active?(.06+(index%3)*.015)*(1-a/.65):0};
}
