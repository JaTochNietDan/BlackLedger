import {hasInterior,interiorSettings} from './interiorSettings.js';
import {isStagedStrike} from './city3dAssassination.js';
import type {VisualCue} from './types';
/** Only an explicit committed home setting selects private-room playback. */
export function isHomeStrike(cue:VisualCue|null|undefined):boolean {
 return !!cue&&cue.strike?.setting==='home'&&['apartment','mercercourt','riverside','estate','room'].includes(cue.target)&&isStagedStrike(cue)&&!!cue.attacker;
}
export function homeStrikeRoom(target:string){
 return target==='room'
  ? {model:'interior-lodging-room',yaw:Math.PI/4,focusZ:0,origin:{x:-2.45,z:2.45},bounds:{x:3,z:3}}
  : target==='estate'
  ? {model:'interior-cypress',yaw:0,focusZ:2.7,origin:{x:-3.3,z:3.4},bounds:{x:5,z:4.5}}
  : {model:'interior-flat',yaw:0,focusZ:2,origin:{x:-3.3,z:2.95},bounds:{x:4,z:3.5}};
}
/** A selected death headline can own the explicitly linked attack playback. */
export function homeStrikeFor(cue:VisualCue|null|undefined,batch:VisualCue[]):VisualCue|undefined {
 if(cue&&isHomeStrike(cue))return cue;
 if(cue?.kind!=='killing')return;
 return batch.find(hit=>isHomeStrike(hit)&&hit.target===cue.target&&hit.minute===cue.minute&&cue.actors?.some(person=>person.id===hit.strike!.victim.id));
}

export function isInteriorStrike(cue:VisualCue|null|undefined):boolean {
 return !!cue&&cue.strike?.setting==='interior'&&isStagedStrike(cue)&&!!cue.attacker;
}
export function interiorStrikeFor(cue:VisualCue|null|undefined,batch:VisualCue[]):VisualCue|undefined {
 if(cue&&isInteriorStrike(cue))return cue;
 if(cue?.kind!=='killing')return;
 return batch.find(hit=>isInteriorStrike(hit)&&hit.target===cue.target&&hit.minute===cue.minute&&cue.actors?.some(person=>person.id===hit.strike!.victim.id));
}
export function strikeRoom(cue:VisualCue){
 if(cue.strike?.setting!=='interior')return homeStrikeRoom(cue.target);
 if(!hasInterior(cue.target))return undefined;
 const room=interiorSettings[cue.target];
 return {model:room.model,yaw:0,focusZ:0,origin:{x:-3,z:0},bounds:{x:room.span,z:room.span}};
}
