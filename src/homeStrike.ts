import {isStagedStrike} from './city3dAssassination.js';
import type {VisualCue} from './types';
/** Only an explicit committed home setting selects private-room playback. */
export function isFlatHomeStrike(cue:VisualCue|null|undefined):boolean {
 return !!cue&&cue.strike?.setting==='home'&&['apartment','mercercourt'].includes(cue.target)&&isStagedStrike(cue)&&!!cue.attacker;
}
export const HOME_STRIKE_ORIGIN={x:-3.3,z:2.95};
/** A selected death headline can own the explicitly linked attack playback. */
export function flatHomeStrikeFor(cue:VisualCue|null|undefined,batch:VisualCue[]):VisualCue|undefined {
 if(cue&&isFlatHomeStrike(cue))return cue;
 if(cue?.kind!=='killing')return;
 return batch.find(hit=>isFlatHomeStrike(hit)&&hit.target===cue.target&&hit.minute===cue.minute&&cue.actors?.some(person=>person.id===hit.strike!.victim.id));
}
