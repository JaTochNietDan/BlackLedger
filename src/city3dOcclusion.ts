import * as THREE from 'three';

/** Orthographic sight lines are parallel, including away from the screen centre. */
export function blockingBuildings(camera: THREE.Camera, target: THREE.Vector3, buildings: Map<string, THREE.Group>) {
  const direction = camera.getWorldDirection(new THREE.Vector3());
  const ray = new THREE.Raycaster(target.clone().addScaledVector(direction, -2000), direction, 0, 1999.95);
  const blocked = new Set<string>();
  const point = new THREE.Vector3();
  for (const [id, building] of buildings) {
    const bounds = building.userData.sightBounds as THREE.Box3 | undefined;
    if (bounds && !ray.ray.intersectBox(bounds, point)) continue;
    if (ray.intersectObject(building, true).length) blocked.add(id);
  }
  return blocked;
}
