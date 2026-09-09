/**
 * Where an action goes in the room panel.
 *
 * The core owns the list of groups and says so — "Groups is the ordered list,
 * for anything that renders them" — and for a long time nothing rendered them.
 * The panel kept its own copy, and the copy was missing "people", so any action
 * in that group whose subject was not standing in the room disappeared: no
 * button, no reason, nothing. Paying the crew a bonus while they are out on
 * collections is the ordinary case of that.
 *
 * So the order comes from the core, and anything carrying a group this build
 * has not heard of joins the first section rather than falling out of the
 * panel. That is the same fallback the core applies to an action nobody has
 * classified.
 */
export interface Placed<A> {id: string; title: string; blurb: string; mine: A[]}

export function placeActions<A extends {group: string}>(
  order: readonly {id: string; title: string; blurb: string}[],
  actions: readonly A[],
): Placed<A>[] {
  if (order.length === 0) return [];
  const known = new Set(order.map(g => g.id));
  const first = order[0].id;
  return order.map(({id, title, blurb}) => ({
    id, title, blurb,
    mine: actions.filter(a => a.group === id || (id === first && !known.has(a.group))),
  }));
}
