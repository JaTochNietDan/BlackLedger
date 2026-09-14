/**
 * How much of the paper the player has not read.
 *
 * The badge on the Herald used to count every story of this life. On a save at
 * day 32 that was seventeen, spread over three days — and the Herald opens on
 * today's edition, which held six. So the button promised seventeen, showed
 * six, and the act of looking marked all seventeen read: eleven stories from
 * earlier days were struck off without ever being on screen. Left alone it
 * grows without bound; by day four hundred it would be a number nobody could
 * act on.
 *
 * The count is what the paper will actually show when it is opened: the
 * stories in the latest issue that the player has not seen yet. The archive is
 * still there behind "All issues" — it is the prompt that has to be honest.
 */
export interface Filed {
  id: string;
  day: number;
}

export function unreadInLatest(paper: readonly Filed[], seen: string | null): number {
  if (paper.length === 0) return 0;
  const latest = paper[0].day;
  const today = paper.filter(s => s.day === latest);
  if (!seen) return today.length;
  const at = today.findIndex(s => s.id === seen);
  // Seen a story from an older issue, or none of these: the whole issue is new.
  return at < 0 ? today.length : at;
}

/** Exact public headline/minute linkage avoids showing an unrelated new story. */
export function articleForScene<T extends {headline:string;minute:number}>(cue:{headline?:string;minute?:number},stories:readonly T[]):T|undefined {
 return cue.headline&&cue.minute!==undefined?stories.find(story=>story.headline===cue.headline&&story.minute===cue.minute):undefined;
}
