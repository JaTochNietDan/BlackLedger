import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';
import {groundCharacter} from './city3dGround.js';
import type {VisualCue} from './types';
export function isIndoorSearch(cue:VisualCue|null|undefined){
 return !!cue&&cue.kind==='robbery'&&!!cue.burglary&&cue.burglary.success&&!cue.burglary.resident_present&&!cue.burglary.fatal&&['room','apartment','mercercourt'].includes(cue.target);
}
export function burglaryRoute(target:string):[number,number][]{
 return target==='room'?[[-2.45,2.45],[-.2,.8],[-.1,-.65],[-.67,-1.63]]:[[0,3.15],[0,1.1],[-.4,-.3],[-1.25,-1.81]];
}
const smooth=(v:number)=>{const t=THREE.MathUtils.clamp(v,0,1);return t*t*(3-2*t);};
/** Successful unattended search only; failure/confrontation needs its own cast. */
export class BurglarySearch {
 readonly root=new THREE.Group();readonly distance:number;readonly arrival:number;readonly duration:number;
 readonly money=new THREE.Mesh(new THREE.BoxGeometry(.13,.018,.07),new THREE.MeshStandardMaterial({color:0x788264}));
 private route:THREE.Vector3[];private hand?:THREE.Mesh;private box=new THREE.Box3();
 constructor(readonly attacker:THREE.Group,private drawer:THREE.Object3D,target:string,readonly taken:number){
  this.root.add(attacker,this.money);this.route=burglaryRoute(target).map(([x,z])=>new THREE.Vector3(x,0,z));
  this.distance=this.route.slice(1).reduce((sum,p,i)=>sum+p.distanceTo(this.route[i]),0);
  this.arrival=this.distance/.95;this.duration=this.arrival+4.7+this.distance/1.25+.4;
  attacker.getObjectByName('elbow1')?.traverse(o=>{if(o instanceof THREE.Mesh&&o.name.startsWith('hand'))this.hand=o;});
  this.update(0);
 }
 update(seconds:number){
  const t=Math.max(0,seconds),search=t-this.arrival,leaving=search>4.7;
  const along=leaving?Math.max(0,this.distance-(search-4.7)*1.25):Math.min(this.distance,t*.95);
  let left=along,point=this.route[0].clone(),heading=0;
  for(let i=1;i<this.route.length;i++){
   const a=this.route[i-1],b=this.route[i],length=a.distanceTo(b);heading=Math.atan2(b.x-a.x,b.z-a.z);
   if(left<=length||i===this.route.length-1){point.lerpVectors(a,b,Math.min(1,left/length));break;}left-=length;
  }
  const walking=(t<this.arrival)||(leaving&&along>0),phase=(leaving?this.distance+(this.distance-along):along)/1.15*Math.PI*2;
  const bend=search>=0&&!leaving?smooth(search/.4)*(1-smooth((search-3.8)/.6)):.0;
  this.attacker.visible=true;
  this.attacker.position.copy(point);this.attacker.rotation.set(.30*bend,walking?heading+(leaving?Math.PI:0):Math.PI,0,'YXZ');
  for(const side of [-1,1]){
   const leg=this.attacker.getObjectByName(`leg${side}`),knee=this.attacker.getObjectByName(`knee${side}`),arm=this.attacker.getObjectByName(`arm${side}`),elbow=this.attacker.getObjectByName(`elbow${side}`);
   const swing=phase+(side===1?Math.PI:0);
   if(leg)leg.rotation.set(walking?Math.sin(swing)*.35:-.30*bend,0,0);
   if(knee)knee.rotation.set(walking?Math.max(0,Math.sin(swing+.7))*.65:0,0,0);
   if(arm)arm.rotation.set(walking?-Math.sin(swing)*.23:0,0,0);if(elbow)elbow.rotation.set(0,0,0);
  }
  this.drawer.position.z=.34*smooth(search/.6)*(1-smooth((search-3.8)/.6));
  if(bend>0){
   const rummage=search>.6&&search<2.7?Math.sin(search*7)*.025:0;
   const pocket=this.taken>0?smooth((search-2.7)/.9):0;
   aimArm(this.attacker,1,new THREE.Vector3(.12+.08*pocket,.88+rummage,.30*bend-.26*pocket));
   aimArm(this.attacker,-1,new THREE.Vector3(-.17,.87,.24*bend));
  }
  groundCharacter(this.attacker,.019);
  this.attacker.visible=t<this.duration-.2;
  this.money.visible=this.taken>0&&search>=2.7&&search<3.6&&!!this.hand;
  if(this.money.visible&&this.hand){this.box.setFromObject(this.hand,true).getCenter(this.money.position);this.root.worldToLocal(this.money.position);this.money.rotation.copy(this.attacker.rotation);}
  this.root.updateMatrixWorld(true);
  return {search,walking,drawer:this.drawer.position.z};
 }
}
