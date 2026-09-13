import * as THREE from 'three';

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
