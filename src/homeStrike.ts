import {isStagedStrike} from './city3dAssassination.js';
import type {VisualCue} from './types';
/** Only an explicit committed home setting selects private-room playback. */
export function isHomeStrike(cue:VisualCue|null|undefined):boolean {
 return !!cue&&cue.strike?.setting==='home'&&['apartment','mercercourt','estate'].includes(cue.target)&&isStagedStrike(cue)&&!!cue.attacker;
}
export function homeStrikeRoom(target:string){
 return target==='estate'
  ? {model:'interior-cypress',focusZ:2.7,origin:{x:-3.3,z:3.4},bounds:{x:5,z:4.5}}
  : {model:'interior-flat',focusZ:2,origin:{x:-3.3,z:2.95},bounds:{x:4,z:3.5}};
}
/** A selected death headline can own the explicitly linked attack playback. */
export function homeStrikeFor(cue:VisualCue|null|undefined,batch:VisualCue[]):VisualCue|undefined {
 if(cue&&isHomeStrike(cue))return cue;
 if(cue?.kind!=='killing')return;
 return batch.find(hit=>isHomeStrike(hit)&&hit.target===cue.target&&hit.minute===cue.minute&&cue.actors?.some(person=>person.id===hit.strike!.victim.id));
}
