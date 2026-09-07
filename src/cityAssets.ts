import buildings from '../public/art/buildings.json';
export const paintedLocations=buildings.map(b=>b.id);
export const paintedAsset=(id:string,condition=100)=>{const b=buildings.find(b=>b.id===id);return b?'/art/'+(b.damage&&condition<b.damage.below?b.damage.file:b.file):null};
export const paintedBuildings=buildings;
