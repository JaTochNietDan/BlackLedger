import * as THREE from 'three';
import {trafficSize,type TrafficPose} from './city3dTraffic.js';

/** Union actual staged paths, including support cars and the full escort sweep. */
export function stagedSceneBounds(slots:readonly {pose:TrafficPose;model:string}[]){
 const bounds=new THREE.Box3();
 for(const {pose,model} of slots){
  const size=trafficSize(model),c=Math.cos(pose.heading),s=Math.sin(pose.heading);
  for(const x of [-size.width/2,size.width/2])for(const z of [-size.length/2,size.length/2])
   for(const y of [0,3])bounds.expandByPoint(new THREE.Vector3(pose.x+x*c+z*s,y,pose.z-x*s+z*c));
 }
 return bounds;
}


/** Fit an entire world envelope, including its height, within the clear centre
 * of an orthographic viewport. Independent of the user's previous zoom. */
export function frameScene(camera: THREE.OrthographicCamera, target: THREE.Vector3, bounds: THREE.Box3) {
  if (bounds.isEmpty()) return;
  const centre = bounds.getCenter(new THREE.Vector3());
  const offset = camera.position.clone().sub(target);
  target.copy(centre);
  camera.position.copy(centre).add(offset);
  camera.lookAt(centre);
  camera.updateMatrixWorld(true);
  const view = new THREE.Box3();
  for (const x of [bounds.min.x, bounds.max.x])
    for (const y of [bounds.min.y, bounds.max.y])
      for (const z of [bounds.min.z, bounds.max.z])
        view.expandByPoint(new THREE.Vector3(x, y, z).applyMatrix4(camera.matrixWorldInverse));
  const size = view.getSize(new THREE.Vector3());
  // Leave space for the toolbar, scene caption and soft particle edges.
  camera.zoom = Math.min(32, (camera.right-camera.left)*.58/Math.max(1,size.x),
    (camera.top-camera.bottom)*.58/Math.max(1,size.y));
  camera.updateProjectionMatrix();
}

/** Short screen-space impact, with no persistent camera or control displacement. */
export function impactPulse(age: number, strength: number) {
  if(age<0 || age>=.32)return {x:0,y:0};
  const decay=(1-age/.32)**2;
  return {x:Math.sin(age*83)*strength*.55*decay,y:Math.cos(age*61)*strength*decay};
}
export function renderImpact(camera: THREE.OrthographicCamera, x: number, y: number, width: number, height: number, render: ()=>void) {
  const matrix=camera.projectionMatrix, oldX=matrix.elements[12], oldY=matrix.elements[13];
  if(width<=0 || height<=0 || (!x&&!y)){render();return;}
  matrix.elements[12]+=2*Math.max(-9,Math.min(9,x))/width;
  matrix.elements[13]+=2*Math.max(-12,Math.min(12,y))/height;
  camera.projectionMatrixInverse.copy(matrix).invert();
  try{render();}finally{
    matrix.elements[12]=oldX;matrix.elements[13]=oldY;
    camera.projectionMatrixInverse.copy(matrix).invert();
  }
}
