import * as THREE from 'three';
import {CityLegs} from './city3dLegs.js';
const smooth=(n:number)=>{const t=THREE.MathUtils.clamp(n,0,1);return t*t*(3-2*t);};
/** Three 20cm risers, 40cm treads. Root is the first ankle's tread position. */
export class CityStairDescent {
 readonly root=new THREE.Group();
 readonly duration=4.65;
 private readonly legs:CityLegs;
 constructor(readonly actor:THREE.Group){this.root.add(actor);this.legs=new CityLegs(actor);this.update(0);}
 update(seconds:number){
  const beat=THREE.MathUtils.clamp((seconds-.35)/.65,0,6),index=Math.min(5,Math.floor(beat)),u=beat-index;
  const feet=[new THREE.Vector3(.12,.1,0),new THREE.Vector3(-.12,.1,0)];
  for(let i=0;i<index;i++){feet[i%2].y-=.2;feet[i%2].z-=.4;}
  const moving=index%2,foot=feet[moving];
  // Lift before advancing past the nosing, then lower onto the new tread.
  foot.z-=.4*smooth((u-.2)/.55);
  const lift=.12*smooth(u/.2)*(1-smooth((u-.65)/.35));
  foot.y+=lift-.2*smooth((u-.65)/.35);
  // Prepare with both feet planted, then recover standing height on the pavement.
  const standing=.145*(1-smooth(seconds/.35)+smooth((seconds-4.25)/.4));
  this.actor.position.set(0,(feet[0].y+feet[1].y-lift)/2-.25+standing,(feet[0].z+feet[1].z)/2);
  this.actor.rotation.set(0,Math.PI,0);this.root.updateMatrixWorld(true);
  const reached=feet.map((p,i)=>this.legs.place(i===0?-1:1,p.clone().sub(this.actor.position).applyQuaternion(this.actor.quaternion.clone().invert())));
  return {feet,moving,reached,done:seconds>=this.duration};
 }
}
