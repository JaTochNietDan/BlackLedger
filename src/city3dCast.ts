/** Appearance follows the existing portrait cast, never a guess from a name.
 * Face numbers are the public one-based portrait selection; zero means default.
 * Keep the sheet mapping aligned with core/voices.go and tools/portraits.py. */
export const castLooks = 'mmfmfmfmfmmfmfmmfmmmfmff';
const painted: Record<string, string> = {
  mara: 'f', elena: 'f', leo: 'm', vittorio: 'm', harlow: 'm', 'Alex Varga': 'm',
};
export function castFace(id: string) {
  let hash = 2166136261;
  for (const byte of new TextEncoder().encode(id)) hash = Math.imul(hash ^ byte, 16777619) >>> 0;
  return hash % 24;
}
export function pedestrianModel(id: string, face = 0, player = false): 'person' | 'woman' {
  const chosen = Number.isInteger(face) && face >= 1 && face <= 24;
  const look = player && chosen ? castLooks[face - 1]
    : painted[id] || castLooks[chosen ? face - 1 : castFace(id)];
  return look === 'f' ? 'woman' : 'person';
}
export function isPedestrian(model: string) {
  return model === 'person' || model === 'woman';
}
