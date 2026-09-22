// Goods and their units come from the public state. Commands still send the
// original integer amount; this only distinguishes quantities from money.
export function actionUnit(
  id: string,
  goods: readonly {id: string; unit: string}[] = [],
): string | undefined {
  const match = /^(?:buy|sell):(.+)$/.exec(id);
  if (!match) return undefined;
  return goods.find(g => g.id === match[1])?.unit || 'unit';
}

export function formatActionAmount(
  amount: number,
  money: (value: number) => string,
  unit?: string,
): string {
  return unit ? `${amount.toLocaleString()} ${unit}${amount === 1 ? '' : 's'}` : money(amount);
}
