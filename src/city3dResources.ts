import * as THREE from 'three';

/** Dispose one mounted city's shared resources exactly once, including prototypes. */
export function disposeCityResources(roots: THREE.Object3D[], extras: {
  geometries?: THREE.BufferGeometry[];
  materials?: THREE.Material[];
  textures?: THREE.Texture[];
} = {}) {
  const geometries = new Set(extras.geometries);
  const materials = new Set(extras.materials);
  const textures = new Set(extras.textures);
  const instances = new Set<THREE.InstancedMesh>();
  const shadows = new Set<THREE.LightShadow>();
  const bitmaps = new Set<ImageBitmap>();
  for (const root of roots) root.traverse(object => {
    if (object instanceof THREE.Mesh) {
      geometries.add(object.geometry);
      for (const material of Array.isArray(object.material) ? object.material : [object.material]) materials.add(material);
      if (object instanceof THREE.InstancedMesh) instances.add(object);
    } else if (object instanceof THREE.Sprite) materials.add(object.material);
    if (object instanceof THREE.DirectionalLight || object instanceof THREE.SpotLight || object instanceof THREE.PointLight)
      shadows.add(object.shadow);
  });
  for (const material of materials)
    for (const value of Object.values(material)) if (value instanceof THREE.Texture) textures.add(value);
  for (const texture of textures) {
    const images = Array.isArray(texture.source.data) ? texture.source.data : [texture.source.data];
    for (const image of images) if (typeof ImageBitmap !== 'undefined' && image instanceof ImageBitmap) bitmaps.add(image);
  }
  instances.forEach(instance => instance.dispose());
  shadows.forEach(shadow => shadow.dispose());
  geometries.forEach(geometry => geometry.dispose());
  materials.forEach(material => material.dispose());
  textures.forEach(texture => texture.dispose());
  bitmaps.forEach(bitmap => bitmap.close());
}
