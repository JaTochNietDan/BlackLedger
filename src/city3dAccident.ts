import * as THREE from 'three';
import {aimArm} from './city3dWeapons.js';
import {CityLegs} from './city3dLegs.js';
const smooth=(value:number)=>{const t=Math.max(0,Math.min(1,value));return t*t*(3-2*t);};
/** Accident fall and recovery. The saved outcome controls whether recovery starts.
 * Ground contact is measured from the posed meshes, including coat and hands. */
export class CityAccident {
 readonly root=new THREE.Group();
 get duration(){return this.fatal?4:7.5;}
 private meshes:THREE.Mesh[]=[];
 private wrists=new Map<number,THREE.Group>();
 private point=new THREE.Vector3();
 private matrix=new THREE.Matrix4();
 private inverse=new THREE.Matrix4();
 constructor(readonly actor:THREE.Group,readonly fatal:boolean){
  this.root.add(actor);
  new CityLegs(actor);
  for(const side of [-1,1]){
   const elbow=actor.getObjectByName('elbow'+side);
   const hand=elbow?.children.find(o=>o.name.startsWith('hand'));
   if(!elbow||!hand)continue;
   const wrist=new THREE.Group();wrist.name='accident-wrist'+side;wrist.position.y=-.295;elbow.add(wrist);
   actor.updateMatrixWorld(true);wrist.attach(hand);this.wrists.set(side,wrist);
  }
  actor.traverse(o=>{if(o instanceof THREE.Mesh)this.meshes.push(o);});
  this.update(0);
 }
 update(seconds:number){
  for(const wrist of this.wrists.values())wrist.quaternion.identity();
  for(const side of [-1,1])this.actor.getObjectByName('ankle'+side)?.quaternion.identity();
  for(const side of [-1,1])for(const name of ['arm','elbow'])this.actor.getObjectByName(name+side)?.quaternion.identity();
  const t=Math.max(0,seconds),fall=smooth((t-.12)/.9);
  const recover=this.fatal?0:smooth((t-2)/1.5);
  const gather=this.fatal?0:smooth((t-4)/1.3);
  const rise=this.fatal?0:smooth((t-5.3)/2.2);
  const blend=(a:number,b:number)=>a+(b-a)*gather;
  this.actor.position.set(0,0,.55*fall);
  this.actor.rotation.set(blend(-Math.PI/2*fall+.8*recover,.55)*(1-rise),0,0);
  for(const side of [-1,1]){
   const set=(name:string,angle:number)=>{const joint=this.actor.getObjectByName(name+side);if(joint)joint.rotation.x=angle;};
   set('leg',blend(-.25*fall-.55*recover,-1.3)*(1-rise));
   set('knee',blend(.4*fall-.2*recover,2)*(1-rise));
   set('arm',blend(-1.1*smooth(t/.25)*(1-.8*smooth((t-.65)/.7))-.3*recover,-.65)*(1-rise));
   set('elbow',blend(.3*fall+.65*recover,.9)*(1-rise));
   const ankle=this.actor.getObjectByName('ankle'+side),hip=this.actor.getObjectByName('leg'+side),knee=this.actor.getObjectByName('knee'+side);
   if(ankle&&hip&&knee)ankle.quaternion.slerp(this.actor.quaternion.clone().multiply(hip.quaternion).multiply(knee.quaternion).invert(),gather);
  }
  this.root.updateMatrixWorld(true);this.inverse.copy(this.root.matrixWorld).invert();
  // Once the feet are gathered underneath, push up over the same planted soles.
  if(gather>0){
   let near=Infinity,far=-Infinity;
   for(const mesh of this.meshes)if(mesh.name.startsWith('shoe')){
    this.matrix.multiplyMatrices(this.inverse,mesh.matrixWorld);
    const positions=mesh.geometry.attributes.position;
    for(let i=0;i<positions.count;i++){
     this.point.fromBufferAttribute(positions,i).applyMatrix4(this.matrix);
     near=Math.min(near,this.point.z);far=Math.max(far,this.point.z);
    }
   }
   if(Number.isFinite(near))this.actor.position.z+=gather*(.32-(near+far)/2);
   this.root.updateMatrixWorld(true);
  }
  let floor=Infinity;
  for(const mesh of this.meshes){
   this.matrix.multiplyMatrices(this.inverse,mesh.matrixWorld);
   const positions=mesh.geometry.attributes.position;
   for(let i=0;i<positions.count;i++){
    this.point.fromBufferAttribute(positions,i).applyMatrix4(this.matrix);
    floor=Math.min(floor,this.point.y);
   }
  }
  if(Number.isFinite(floor))this.actor.position.y=.005-floor;
  this.root.updateMatrixWorld(true);
  const brace=this.fatal?0:smooth((t-1.35)/.6)*(1-smooth((t-4)/.65));
  if(brace>0)for(const side of [-1,1]){
   const arm=this.actor.getObjectByName('arm'+side),elbow=this.actor.getObjectByName('elbow'+side),wrist=this.wrists.get(side);
   if(!arm||!elbow||!wrist)continue;
   const start=this.root.worldToLocal(wrist.getWorldPosition(new THREE.Vector3()));
   const contact=start.lerp(new THREE.Vector3(side*.43,.084,-.6),brace)
    .add(new THREE.Vector3(0,.12*Math.sin(Math.PI*brace),0));
   const target=this.actor.worldToLocal(this.root.localToWorld(contact));
   aimArm(this.actor,side,target,new THREE.Vector3(side*.2,1,0).applyQuaternion(this.actor.quaternion.clone().invert()));
   const orientation=this.actor.quaternion.clone().multiply(arm.quaternion).multiply(elbow.quaternion).invert()
    .multiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(1,0,0),Math.PI/2));
   wrist.quaternion.slerp(orientation,brace);
   this.root.updateMatrixWorld(true);
  }
  return {fatal:this.fatal,recovering:recover>0,done:t>=this.duration};
 }
}
