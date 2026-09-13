import type {Snapshot, VisualCue} from './types';

export const previewScenes = ['Gunfight', 'Assassination', 'Explosion', 'Arrest', 'Raid'] as const;
export type PreviewScene = typeof previewScenes[number];
/** A private presentation snapshot. No API command or campaign object is changed. */
export function previewScene(state: Snapshot, target: string, scene: PreviewScene, token: string) {
  const kinds = scene === 'Assassination' ? ['killing', 'gunfight'] : [scene.toLowerCase()];
  const cues: VisualCue[] = kinds.map((kind, i) => ({
    id: `preview:${token}:${i}`, kind, target, minute: state.minute,
    caption: `Visual preview: ${scene}`, detainee: kind==='arrest'?{id:'preview-detainee',name:'Preview detainee'}:undefined, actors: kind === 'killing'
      ? [{id: 'preview-victim', name: 'Preview character'}] : [],
  }));
  return {state: {...state, id: `${state.id}:preview:${token}`,
    last_result: {from_location: target, to_location: target, elapsed: 0, records: [], cash: 0, respect: 0, heat: 0, health: 0, ...state.last_result, cues}}, cue: cues[0]};
}
