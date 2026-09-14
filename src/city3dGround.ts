import * as THREE from 'three';

/** Seat a posed rig on its parent's horizontal surface, using actual vertices. */
export function groundCharacter(actor:THREE.Group, height:number) {
  actor.updateWorldMatrix(true,true);
  const inverse=actor.parent?.matrixWorld.clone().invert()||new THREE.Matrix4();
  const transform=new THREE.Matrix4(),point=new THREE.Vector3();
  let lowest=Infinity;
  actor.traverseVisible(part=>{
    if(!(part instanceof THREE.Mesh))return;
    transform.multiplyMatrices(inverse,part.matrixWorld);
    const vertices=part.geometry.attributes.position;
    for(let i=0;i<vertices.count;i++)lowest=Math.min(lowest,point.fromBufferAttribute(vertices,i).applyMatrix4(transform).y);
  });
  if(Number.isFinite(lowest))actor.position.y+=height-lowest;
  actor.updateMatrixWorld(true);
}
