import buildings from '../public/art/buildings.json';
import previews from '../public/art/previews.json';
export const paintedLocations=buildings.map(b=>b.id);
export const paintedAsset=(id:string,condition=100)=>{const b=buildings.find(b=>b.id===id);return b?'/art/'+(b.damage&&condition<b.damage.below?b.damage.file:b.file):previews.some(p=>p.id===id)?'/art/'+previews.find(p=>p.id===id)!.file:null};
export const paintedBuildings=buildings;

export const paintedMask=(id:string)=>{const p=previews.find(p=>p.id===id);return p?'/art/'+p.mask:paintedAsset(id)};
