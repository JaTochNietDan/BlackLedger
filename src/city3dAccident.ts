import * as THREE from 'three';
const smooth=(value:number)=>{const t=Math.max(0,Math.min(1,value));return t*t*(3-2*t);};
/** Initial accident fall study. The saved outcome controls whether recovery starts.
 * Ground contact is measured from the posed meshes, including coat and hands. */
export class CityAccident {
 readonly root=new THREE.Group();
 readonly duration=4;
 private meshes:THREE.Mesh[]=[];
 private point=new THREE.Vector3();
 private matrix=new THREE.Matrix4();
 private inverse=new THREE.Matrix4();
 constructor(readonly actor:THREE.Group,readonly fatal:boolean){
  this.root.add(actor);
  actor.traverse(o=>{if(o instanceof THREE.Mesh)this.meshes.push(o);});
  this.update(0);
 }
 update(seconds:number){
  const t=Math.max(0,seconds),fall=smooth((t-.12)/.9);
  const recover=this.fatal?0:smooth((t-2)/1.5);
  this.actor.position.set(0,0,.55*fall);
  this.actor.rotation.set(-Math.PI/2*fall+.8*recover,0,0);
  for(const side of [-1,1]){
   const set=(name:string,angle:number)=>{const joint=this.actor.getObjectByName(name+side);if(joint)joint.rotation.x=angle;};
   set('leg',-.25*fall-.55*recover);
   set('knee',.4*fall-.2*recover);
   set('arm',-1.1*smooth(t/.25)*(1-.8*smooth((t-.65)/.7))-.3*recover);
   set('elbow',.3*fall+.65*recover);
  }
  this.root.updateMatrixWorld(true);this.inverse.copy(this.root.matrixWorld).invert();
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
  return {fatal:this.fatal,recovering:recover>0,done:t>=this.duration};
 }
}
