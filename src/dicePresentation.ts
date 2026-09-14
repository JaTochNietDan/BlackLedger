import * as THREE from 'three';
export const DICE_ROLL_MS=1100;
const normals=[new THREE.Vector3(),new THREE.Vector3(0,1,0),new THREE.Vector3(1,0,0),new THREE.Vector3(0,0,1),new THREE.Vector3(0,0,-1),new THREE.Vector3(-1,0,0),new THREE.Vector3(0,-1,0)];
export class PresentedDie {
 private points:number[]=[];
 private rotation=new THREE.Matrix4();
 private tumble=new THREE.Quaternion();
 private end=new THREE.Quaternion();
 constructor(readonly object:THREE.Group,private index:number){
  object.updateMatrixWorld(true);const inverse=object.matrixWorld.clone().invert(),point=new THREE.Vector3();
  object.traverse(o=>{if(!(o instanceof THREE.Mesh))return;const positions=o.geometry.attributes.position;for(let i=0;i<positions.count;i++){point.fromBufferAttribute(positions,i).applyMatrix4(o.matrixWorld).applyMatrix4(inverse);this.points.push(point.x,point.y,point.z);}});
 }
 pose(face:number,progress:number){
  this.object.visible=Number.isInteger(face)&&face>=1&&face<=6;if(!this.object.visible)return;
  const t=THREE.MathUtils.clamp(progress,0,1),ease=1-Math.pow(1-t,3),side=this.index===0?-1:1;
  this.end.setFromUnitVectors(normals[face],new THREE.Vector3(0,1,0));
  this.end.premultiply(new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0,1,0),side*.32));
  this.tumble.setFromEuler(new THREE.Euler((1-t)**2*Math.PI*4,0,(1-t)**2*Math.PI*3*side));
  this.object.quaternion.copy(this.end).multiply(this.tumble);
  this.rotation.makeRotationFromQuaternion(this.object.quaternion);const e=this.rotation.elements;let low=Infinity;
  for(let i=0;i<this.points.length;i+=3)low=Math.min(low,e[1]*this.points[i]+e[5]*this.points[i+1]+e[9]*this.points[i+2]);
  const bounce=.22*Math.abs(Math.sin(t*Math.PI*3))*(1-t);
  this.object.position.set(side*THREE.MathUtils.lerp(.75,.23,ease),.026-low+bounce,THREE.MathUtils.lerp(.37,side*.045,ease));
  this.object.updateMatrixWorld(true);
 }
}
