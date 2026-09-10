// The arrangement of the city, as opposed to the facts about it.
//
// Where an address stands comes from the core, and nothing here may argue with
// it. What is here is everything the core has no opinion about: which picture
// stands in which slot of which block, and how far it has been nudged off the
// exact centre of its ground. That is a matter of taste, so it is arranged by
// hand and kept in a file two people can both edit — one by dragging in the
// browser, the other in a text editor, both arriving as ordinary diffs.

export type Placement = {sprite: string; dx?: number; dy?: number};
export type Layout = {slots: Record<string, Placement>; editable?: boolean};

export const EMPTY: Layout = {slots: {}};

// A slot's name. Block column, block row, and which slot of the terrace — the
// same key the renderer already builds when it lays the blocks out, so the file
// can be read by a person and matched to a place on the map.
export const slotKey = (col: number, row: number, index: number) => `${col},${row},${index}`;

// A sprite deliberately left empty, so a slot can be cleared and stay cleared
// rather than falling back to whatever the automatic choice would have been.
export const NOTHING = '';

export async function loadLayout(): Promise<Layout> {
  try {
    const reply = await fetch('/api/layout');
    if (!reply.ok) return EMPTY;
    const body = await reply.json();
    return {slots: body.slots || {}, editable: !!body.editable};
  } catch {
    return EMPTY; // a city with no arrangement is the automatic one
  }
}

export async function saveLayout(layout: Layout): Promise<string> {
  const reply = await fetch('/api/layout', {
    method: 'POST',
    headers: {'content-type': 'application/json'},
    body: JSON.stringify({slots: layout.slots}),
  });
  if (reply.ok) return '';
  try {
    return (await reply.json()).error || `save failed (${reply.status})`;
  } catch {
    return `save failed (${reply.status})`;
  }
}
