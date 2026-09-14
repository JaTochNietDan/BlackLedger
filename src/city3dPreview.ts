import type {Snapshot, VisualCue} from './types';

const strikes = {
  'Assassination · shot from behind': {variant: 'back-of-head', weapon: 1},
  'Assassination · revolver': {variant: 'close-shot', weapon: 1},
  'Assassination · shotgun': {variant: 'close-shot', weapon: 2},
  'Assassination · Thompson': {variant: 'burst', weapon: 3},
  'Assassination · close quarters': {variant: 'close-quarters', weapon: 0},
} as const;
export const previewScenes = ['Gunfight', ...Object.keys(strikes) as (keyof typeof strikes)[], 'Explosion', 'Arrest', 'Raid'] as const;
export type PreviewScene = typeof previewScenes[number];
/** A private presentation snapshot. No API command or campaign object is changed. */
export function previewScene(state: Snapshot, target: string, scene: PreviewScene, token: string) {
  const strike = scene in strikes ? strikes[scene as keyof typeof strikes] : undefined;
  const kinds = strike ? ['killing', strike.weapon ? 'gunfight' : 'attack'] : [scene.toLowerCase()];
  const cues: VisualCue[] = kinds.map((kind, i) => ({
    id: `preview:${token}:${i}`, kind, target, minute: state.minute,
    caption: `Visual preview: ${scene}`,
    strike:strike?{variant:strike.variant,victim:{id:'preview-victim',name:'Preview character'}}:undefined,
    attacker:strike&&kind!=='killing'?{id:'preview-assassin',name:'Preview assassin',weapon:strike.weapon}:undefined, detainee: kind==='arrest'?{id:'preview-detainee',name:'Preview detainee'}:undefined, actors: kind === 'killing'
      ? [{id: 'preview-victim', name: 'Preview character'}] : [],
  }));
  return {state: {...state,
    building_fires:scene==='Explosion'?[...(state.building_fires||[]).filter(f=>f.target!==target),
      {id:`preview-fire:${token}`,target,minute:state.minute,brigade_at:state.minute+10,extinguished_at:state.minute+180,cleanup_at:state.minute+240}]:state.building_fires,
    id: `${state.id}:preview:${token}`,
    last_result: {from_location: target, to_location: target, elapsed: 0, records: [], cash: 0, respect: 0, heat: 0, health: 0, cues}}, cue: cues[0]};
}
