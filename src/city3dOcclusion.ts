import * as THREE from 'three';

/** Sample the pedestrian outline as well as the torso. A narrow roof edge can
 * cover a head or foot while leaving the centre ray unobstructed. */
export function characterSightPoints(camera: THREE.Camera, centre: THREE.Vector3) {
  const right = new THREE.Vector3().setFromMatrixColumn(camera.matrixWorld, 0).normalize();
  return [centre, ...[-.8, .8].flatMap(y => [-.4, .4].map(x =>
    centre.clone().addScaledVector(right, x).add(new THREE.Vector3(0, y, 0))))];
}

/** Orthographic sight lines are parallel, including away from the screen centre. */
export function blockingBuildings(camera: THREE.Camera, target: THREE.Vector3, buildings: Map<string, THREE.Group>) {
  const direction = camera.getWorldDirection(new THREE.Vector3());
  const ray = new THREE.Raycaster(target.clone().addScaledVector(direction, -2000), direction, 0, 1999.95);
  const blocked = new Set<string>();
  const point = new THREE.Vector3();
  for (const [id, building] of buildings) {
    const bounds = building.userData.sightBounds as THREE.Box3 | undefined;
    if (bounds && !ray.ray.intersectBox(bounds, point)) continue;
    const visibleHit = ray.intersectObject(building, true).some(hit => {
      for (let part: THREE.Object3D | null = hit.object; part; part = part.parent)
        if (!part.visible) return false;
      if (hit.object instanceof THREE.Mesh) {
        const material = Array.isArray(hit.object.material)
          ? hit.object.material[hit.face?.materialIndex ?? 0] : hit.object.material;
        if (material && !material.visible) return false;
      }
      return true;
    });
    if (visibleHit) blocked.add(id);
  }
  return blocked;
}
