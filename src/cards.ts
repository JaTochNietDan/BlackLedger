// What a playing card looks like, and what a roulette cloth is laid out as.
// Nothing here decides anything: the core deals the cards, spins the wheel and
// settles the money. This is only how those facts are arranged on a table.

export interface Card {rank:string; suit:string; value:number}

// The four suits, as the core names them, with the character to print and
// whether it is a red suit. A card the core sends with a suit not in here is
// drawn as a plain rectangle rather than as a guess.
const suits:{[id:string]:{pip:string; red:boolean}} = {
  spades: {pip: '♠', red: false},
  hearts: {pip: '♥', red: true},
  diamonds: {pip: '♦', red: true},
  clubs: {pip: '♣', red: false},
};

export const pipOf = (suit:string) => suits[suit]?.pip ?? '';
export const isRedSuit = (suit:string) => suits[suit]?.red ?? false;

// A card the view knows how to draw. Anything else is shown face down rather
// than invented.
export const knownCard = (c:Card) => !!c && !!suits[c.suit] && !!c.rank;

// The cloth. A roulette table is three columns of twelve with the nought
// across the top, and it has to be laid out in that order or a player who has
// stood at one will not recognise it.
export const clothRows = ():number[][] => {
  const rows:number[][] = [];
  for (let row = 0; row < 12; row++) {
    rows.push([row * 3 + 3, row * 3 + 2, row * 3 + 1]);
  }
  return rows;
};

// The red pockets, which are not every other number. This is the same layout
// the core holds; the view needs it only to colour the cloth, and the core
// remains the authority on what actually won.
const reds = new Set([1,3,5,7,9,12,14,16,18,19,21,23,25,27,30,32,34,36]);
export const clothColour = (n:number) => n === 0 ? 'green' : reds.has(n) ? 'red' : 'black';

// The outside bets along the bottom, in the order they sit on a real cloth,
// paired with the bet ids the core accepts.
export const outsideBets = () => [
  {id: 'dozen1', label: '1st 12', wide: true},
  {id: 'dozen2', label: '2nd 12', wide: true},
  {id: 'dozen3', label: '3rd 12', wide: true},
  {id: 'low', label: '1-18', wide: false},
  {id: 'even', label: 'Even', wide: false},
  {id: 'red', label: 'Red', wide: false},
  {id: 'black', label: 'Black', wide: false},
  {id: 'odd', label: 'Odd', wide: false},
  {id: 'high', label: '19-36', wide: false},
];

// Where a pocket sits on the wheel face, in degrees clockwise from the top.
// The order is the wheel's own, which is not the cloth's and not numerical:
// it alternates colours and spreads the low numbers around the rim.
export const wheelOrder = [
  0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10,
  5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26,
];

export const wheelAngle = (pocket:number) => {
  const at = wheelOrder.indexOf(pocket);
  return at < 0 ? 0 : at * (360 / wheelOrder.length);
};
