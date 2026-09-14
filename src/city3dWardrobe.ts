import * as THREE from 'three';
import {castFace} from './city3dCast.js';

// Palette direction from the shipped 24-cell noir portrait sheet, in sheet order.
// These are authored appearance choices, with no gameplay meaning.
const coats = ['#293438','#30342f','#33352d','#32372e','#303438','#303538',
  '#2c302d','#423b34','#30352e','#464a3c','#343732','#3d302c',
  '#32393b','#353936','#3d4235','#494033','#292e2b','#454b37',
  '#40392b','#363831','#344047','#323831','#4b4c40','#343e39'];
const hair = ['#68685d','#5b5040','#29271f','#29281e','#adae9e','#343329',
  '#3b3528','#655039','#3e3525','#51412b','#483724','#503726',
  '#b8b5a4','#2e2b22','#36362d','#c1beb0','#3f3022','#9c6a35',
  '#666252','#634e37','#443127','#333329','#bab9a7','#342a20'];
const skin = ['#be9774','#c4a17f','#d1ae85','#bf9870','#c9ad8b','#b68b64',
  '#c6a37d','#c5a079','#ccb18b','#c3a27e','#b68e68','#cbb190',
  '#c6a47f','#c6ad8c','#b58e69','#c9aa86','#d0b28e','#cfaf86',
  '#a77d55','#bf9672','#d2b99c','#c0a07b','#cab492','#d0b390'];
export function wardrobe(id: string, face = 0, player = false) {
  const selected = Number.isInteger(face) && face >= 1 && face <= 24;
  const index = selected ? face - 1 : castFace(id);
  // Named painted cast retain a restrained palette independent of the generated sheet.
  const named: Record<string, number> = {mara: 2, elena: 8, leo: 7, vittorio: 14, harlow: 5, 'Alex Varga': 3};
  const chosen = (!player || !selected) && Object.hasOwn(named, id) ? named[id] : index;
  const hairStyle = [0,1,12,15].includes(chosen) ? 'receding' : chosen === 9 ? 'bald' : 'full';
  const headwear = chosen === 18 ? 'cap' : ((!player || !selected) && ['leo','vittorio','harlow','Alex Varga'].includes(id)) ? 'fedora' : 'none';
  return {coat: coats[chosen], hair: hair[chosen], skin: skin[chosen], hat: coats[chosen], shirt: '#c8c3ad', hairStyle, headwear};
}

/** Clone once per source material, sharing geometry and embedded texture maps. */
export function dressPedestrian(object: THREE.Group, model: string, look: ReturnType<typeof wardrobe>) {
  for (const style of ['full','receding']) {
    const group = object.getObjectByName('hair-' + style);
    if (group) group.visible = look.hairStyle === style;
  }
  for (const style of ['fedora','cap']) {
    const group = object.getObjectByName('headwear-' + style);
    if (group) group.visible = look.headwear === style;
  }
  const copies = new Map<THREE.Material, THREE.MeshStandardMaterial>();
  const woolBase = new THREE.Color().setRGB(...(model === 'woman' ? [.29,.13,.12] : [.16,.19,.21]) as [number,number,number], THREE.SRGBColorSpace);
  const tint = (source: THREE.Material) => {
    if (!(source instanceof THREE.MeshStandardMaterial)) return source;
    const colour = source.name === 'wool suit' ? look.coat : source.name === 'skin' ? look.skin
      : source.name === 'felt hat' ? look.hat : ['waved chestnut hair','barbered hair'].includes(source.name) ? look.hair
      : source.name === 'ivory shirt' ? look.shirt : undefined;
    if (!colour) return source;
    let own = copies.get(source);
    if (!own) {
      own = source.clone(); own.userData.castSource=source.uuid; own.color.set(colour);
      // Existing wool maps contain the original suit pigment: remove that tint
      // before applying the cast palette, retaining its weave and fibre contrast.
      if (source.name === 'wool suit') own.color.setRGB(own.color.r/woolBase.r,own.color.g/woolBase.g,own.color.b/woolBase.b);
      copies.set(source, own);
    }
    return own;
  };
  object.traverse(part => {if (part instanceof THREE.Mesh) part.material = Array.isArray(part.material) ? part.material.map(tint) : tint(part.material);});
  return [...copies.values()];
}
