import * as THREE from 'three';
import {CityPlanter} from './city3dPlanter.js';
import {CityStairDescent} from './city3dStairs.js';
/** Villa landing-local approach followed by its authored stair descent. */
export class CityVillaExit {
 readonly root=new THREE.Group();
 readonly duration=12.9;
 private readonly approach:CityPlanter;
 private readonly stairs:CityStairDescent;
 private readonly departure:CityPlanter;
 private readonly clearance:CityStairDescent;
 constructor(readonly actor:THREE.Group){
  this.approach=new CityPlanter(actor,3.36,0);
  this.stairs=new CityStairDescent(actor);
  this.stairs.root.position.z=-1.11;
  this.clearance=new CityStairDescent(actor,0,1);
  this.clearance.root.position.set(0,-.6,-2.31);
  this.departure=new CityPlanter(actor,0,2);
  this.departure.root.position.set(0,-.6,-4.96);
  this.root.add(this.approach.root,this.stairs.root,this.clearance.root,this.departure.root);this.update(0);
 }
 update(seconds:number){
  let reached=true;
  if(seconds<3.9){this.approach.root.add(this.actor);this.approach.update(seconds);}
  else if(seconds<8.55){this.stairs.root.add(this.actor);reached=this.stairs.update(seconds-3.9).reached.every(Boolean);}
  else if(seconds<10.6){this.clearance.root.add(this.actor);reached=this.clearance.update(seconds-8.55).reached.every(Boolean);}
  else {this.departure.root.add(this.actor);this.departure.update(seconds-10.6+3.9);}
  const t=THREE.MathUtils.clamp(seconds/.4,0,1),close=THREE.MathUtils.clamp((seconds-3.3)/.5,0,1);
  return {door:t*t*(3-2*t)*(1-close*close*(3-2*close)),reached,done:seconds>=this.duration};
 }
}
