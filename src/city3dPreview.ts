import type {Snapshot, VisualCue} from './types';

const strikes = {
  'Assassination · shot from behind': {variant: 'back-of-head', weapon: 1},
  'Assassination · revolver': {variant: 'close-shot', weapon: 1},
  'Assassination · shotgun': {variant: 'close-shot', weapon: 2},
  'Assassination · Thompson': {variant: 'burst', weapon: 3},
  'Assassination · close quarters': {variant: 'close-quarters', weapon: 0},
} as const;
const gunfights = {'Gunfight':1,'Gunfight · shotgun':2,'Gunfight · Thompson':3} as const;
export const previewScenes = [...Object.keys(gunfights) as (keyof typeof gunfights)[], ...Object.keys(strikes) as (keyof typeof strikes)[], 'Explosion', 'Explosion · casualty', 'Explosion · premature', 'Explosion · fatal accident', 'Incendiary', 'Building drive-by', 'Arrest', 'Raid'] as const;
export type PreviewScene = typeof previewScenes[number];
/** A private presentation snapshot. No API command or campaign object is changed. */
export function previewScene(state: Snapshot, target: string, scene: PreviewScene, token: string) {
  const strike = scene in strikes ? strikes[scene as keyof typeof strikes] : undefined;
  const firearm=scene in gunfights?gunfights[scene as keyof typeof gunfights]:undefined;
  const explosion=scene.startsWith('Explosion');
  const premature=scene==='Explosion · premature'||scene==='Explosion · fatal accident';
  const kinds = scene==='Explosion · casualty'?['killing','explosion']:strike ? ['killing', strike.weapon ? 'gunfight' : 'attack'] : [scene==='Building drive-by'?'driveby-building':firearm?'gunfight':explosion?'explosion':scene.toLowerCase()];
  const cues: VisualCue[] = kinds.map((kind, i) => ({
    id: `preview:${token}:${i}`, kind, target, minute: state.minute,
    caption: `Visual preview: ${scene}`,
    detonation:kind==='explosion'?(premature?'premature':'planted'):undefined,
    accident:premature?{health_lost:scene==='Explosion · fatal accident'?40:30,fatal:scene==='Explosion · fatal accident'}:undefined,
    strike:strike?{variant:strike.variant,victim:{id:'preview-victim',name:'Preview character'}}:undefined,
    drive_by:scene==='Building drive-by'?{driver:{id:'preview-driver',name:'Preview driver'},vehicle:'A Packard',vehicle_tier:3,condition_before:100,condition_after:72}:undefined,
    attacker:scene==='Building drive-by'?{id:'preview-shooter',name:'Preview shooter',weapon:3}:kind==='explosion'?{id:!premature?'player':'preview-planter',name:!premature?state.player?.name||'Preview planter':'Preview planter',weapon:0}:scene==='Incendiary'?{id:'preview-arsonist',name:'Preview arsonist',weapon:0}:strike&&kind!=='killing'?{id:'preview-assassin',name:'Preview assassin',weapon:strike.weapon}:firearm?{id:'preview-attacker',name:'Preview attacker',weapon:firearm??0}:undefined, detainee: kind==='arrest'?{id:'preview-detainee',name:'Preview detainee'}:undefined, actors: kind === 'killing'
      ? [{id: 'preview-victim', name: 'Preview character'}] : [],
  }));
  return {state: {...state,
    building_fires:['Explosion','Explosion · casualty','Incendiary'].includes(scene)?[...(state.building_fires||[]).filter(f=>f.target!==target),
      {id:`preview-fire:${token}`,target,minute:state.minute,brigade_at:state.minute+10,extinguished_at:state.minute+180,cleanup_at:state.minute+240}]:state.building_fires,
    id: `${state.id}:preview:${token}`,
    last_result: {from_location: target, to_location: target, elapsed: 0, records: [], cash: 0, respect: 0, heat: 0, health: 0, cues}}, cue: cues[0]};
}
