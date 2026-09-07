import buildings from '../public/art/buildings.json';
export const paintedLocations=buildings.map(b=>b.id);
export const paintedAsset=(id:string)=>{const b=buildings.find(b=>b.id===id);return b?'/art/'+b.file:null};
export const paintedBuildings=buildings;
