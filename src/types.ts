import type {PoolInput,PoolState,PoolOpponent,PoolTournamentState,PoolTournamentNoticeState} from './billiards';
export interface Group {
  id: string;
  title: string;
  blurb: string;
}
export interface VisualCue {
  burglary?: {intruder:{id:string;name:string};resident:{id:string;name:string};resident_present:boolean;success:boolean;taken:number;health_lost:number;identified:boolean;fatal:boolean};
  detonation?: 'planted' | 'premature';
  accident?: {health_lost: number; fatal: boolean};
  drive_by?: {driver: {id: string; name: string}; vehicle: string; vehicle_tier: number; condition_before: number; condition_after: number};
  id: string;
  kind: string;
  target: string;
  caption: string;
  headline?: string;
  detainee?: {id: string; name: string};
  strike?: {setting?: 'home'; variant: 'back-of-head' | 'close-shot' | 'burst' | 'close-quarters'; victim: {id: string; name: string}};
  attacker?: {id: string; name: string; weapon: number};
  actors?: {id: string; name: string}[];
  gravity?: number;
  minute?: number;
}
export interface Sum {
  least: number;
  most: number;
  preset: number;
  label: string;
}
export interface Action {
  anywhere?: boolean;
  group: string;
  subject?: string;
  id: string;
  label: string;
  minutes: number;
  away?: number;
  cost: number;
  asks?: number;
  tier?: number;
  choice?: string;
  sum?: Sum;
  disabled: boolean;
  reason: string;
  detail: string;
  target: string;
}
export interface Presence {
  home_id?: string;
  home_name?: string;
  accommodation?: string;
  says?: string;
  face?: number;
  id: string;
  name: string;
  role?: string;
  faction?: string;
  standing: string;
  doing?: string;
  where?: string;
  where_id?: string;
  lost?: string;
  because?: string;
  trust?: number;
  sore?: number;
  owes?: number;
  overdue?: boolean;
  yours?: boolean;
  carrying?: string;
  driving?: string;
  restless?: boolean;
  known?: boolean;
  temperament?: string;
  walking?: boolean;
  minutes?: number;
}
export interface Place {
	 rent_register?: {capacity:number; occupied:number; daily:number; tenants:{id:string;name:string;daily:number;accommodation:string;account?:{day:number;due:number;paid:number;arrears:number;collected:number}}[]} | null;
  crossing?: {
    to: string;
    to_id: string;
    minutes: number;
    driving: boolean;
    plate: number;
    plate_max: number;
    warned: boolean;
    note: string;
  };
  trading?: number;
  people?: Presence[];
  note?: string;
  note_warn?: boolean;
  room?: string;
  away?: number;
  travel_note?: string;
  trade?: {custom: number; multiplier: number; order: boolean; order_pays: number} | null;
  bankroll?: number;
  handle?: number;
  still?: boolean;
  staff?: number;
  wage?: number;
  runs?: string;
  hands?: {id: string; name: string; role: string; here: boolean}[];
  positions?: number;
  run_as?: string;
  skimmed?: boolean;
  back_room?: boolean;
  supply?: number;
  unpaid?: number;
  trouble?: boolean;
  capacity?: number;
  holder?: string;
  id: string;
  name: string;
  type: string;
  district: number;
  x: number;
  y: number;
  cost: number;
  blurb: string;
  owner: string;
  condition: number;
  income: number;
  owned: boolean;
  locked: boolean;
  actions: Action[];
}
export interface Person {
  face?: number;
  stock?: {[good: string]: number};
  earned?: number;
  name: string;
  cash: number;
  health: number;
  respect: number;
  heat: number;
  location: string;
  home: string;
  security: number;
  contacts: number;
  crew: {id: string; name: string; loyalty: number}[];
  alive: boolean;
  job_count: number;
}
export interface NPC {
  id: string;
  name: string;
  role: string;
  trust: number;
  voice: string;
  color: string;
}
export interface Event {
  conditions?: string;
  connection?: {id: string; title: string; result: string} | null;
  id: string;
  title: string;
  body: string;
  speaker: string;
  kind: string;
  source: string;
  minute: number;
  choices: {
    id: string;
    label: string;
    detail: string;
    cost: number;
    disabled: boolean;
    pay?: number;
    minutes?: number;
    respect?: number;
    heat?: number;
    reason?: string;
  }[];
}
export interface Record {
  id: string;
  minute: number;
  life: number;
  title: string;
  text: string;
  kind: string;
  count?: number;
}
export interface Coming {
  id: string;
  name: string;
  where: string;
  to?: string;
  leaving?: boolean;
  note: string;
}
export interface Journeying {
  vehicle?: string;
  id: string;
  name: string;
  from_id: string;
  from: string;
  to_id: string;
  to: string;
  because: string;
  minutes: number;
  progress: number;
  yours?: boolean;
}
export interface StreetSegment extends Journeying {from_minute:number;to_minute:number;end_progress:number;}
export interface Snapshot {
  pool?:PoolState|null;
  pool_tournament?:PoolTournamentState|null;
  pool_tournament_notice?:PoolTournamentNoticeState|null;
  pool_opponents?:PoolOpponent[];
	apartment_market?: {id:string;building:string;number:number;address:string;owned:boolean;home:boolean;available:boolean;locked?:boolean;owner:string;resident:string;asking:number;offer:number;daily_rent:number;rent_paid_today?:number;rent_arrears?:number;vacant?:boolean;neighborhood_index?:number}[];
	property_market?: {id:string;name:string;owned:boolean;available:boolean;holder:string;asking:number;offer:number;condition:number;neighborhood_index?:number;residents:number;home:boolean;locked:boolean}[];
  housing_shortage?: number;
  house?: {
    games: boolean;
    place?: string;
    limit?: number;
    machine?: number;
    least?: number;
    usual?: number;
    pull?: number;
    floor?: number;
    ceiling?: number;
    yours?: boolean;
    high?: number;
  };
  street?: Journeying[];
  street_note?: string;
  city?: {
    scrutiny: number;
    state: string;
    crackdown: boolean;
    premium: number;
    raids_sooner: number;
  };
  service?: {
    serving: boolean;
    name?: string;
    title?: string;
    work?: number;
    pay?: number;
    next?: number;
  };
  pacts?: {id: string; name: string; tribute: number; strength: boolean; power: number}[];
  own_people?: {
    id: string;
    name: string;
    trust: number;
    temperament: string;
    wage: number;
    restless: boolean;
  }[];
  organization?: {named: boolean; name?: string; power?: number; peak?: number; needs?: string[]};
  roles?: {role: string; title: string; name: string; id?: string}[];
  machine?: {
    pulled: boolean;
    place?: string;
    stake?: number;
    line?: string[];
    pays?: number;
    won?: boolean;
    stops: number;
    edge: number;
    two_cherries: number;
    one_cherry: number;
    strip: {id: string; face: string; stops: number; pays: number}[];
  };
  seated?: string;
  seated_to?: string;
  dice?: {
    playing: boolean;
    settled: boolean;
    place?: string;
    where?: string;
    bet?: string;
    bet_label?: string;
    stake?: number;
    point?: number;
    dice?: number[];
    total?: number;
    rolls?: number;
    won?: boolean;
    outcome?: string;
    bets?: {id: string; label: string; detail: string}[];
  };
  cards?: {
    place: string;
    ante: number;
    pot: number;
    mine: {rank: string; suit: string; value: number}[];
    board: {rank: string; suit: string; value: number}[];
    street: string;
    street_name: string;
    hand: string;
    seats: {
      who: string;
      name: string;
      threw: number;
      in: number;
      folded: boolean;
      said: string;
      sore?: number;
      moved?: number;
      stack: number;
      cards?: {rank: string; suit: string; value: number}[];
      hand?: string;
    }[];
    bet: number;
    my_bet: number;
    facing: boolean;
    folded: boolean;
    done: boolean;
    outcome: string;
    won: number;
    stack: number;
    buy_in: number;
    hands: number;
    over: boolean;
    ended: string;
    up: number;
  } | null;
  hand?: {
    playing: boolean;
    place?: string;
    stake?: number;
    player?: number;
    dealer?: number;
    cards?: number;
    mine?: {rank: string; suit: string; value: number}[];
    theirs?: {rank: string; suit: string; value: number}[];
  };
  wheel?: {
    spun: boolean;
    place?: string;
    stake?: number;
    bet?: string;
    pocket?: number;
    colour?: string;
    won?: boolean;
    pays?: number;
    back?: number;
    chips?: {bet: string; label: string; amount: number; won: boolean; back: number}[];
  };
  population?: {living: number; organized: number; street: number; jobs: number; known: number};
  armoury?: {
    held: boolean;
    place?: string;
    crates?: number;
    capacity?: number;
    price?: number;
    buyers?: number;
    attention?: number;
  };
  retainers?: {
    id: string;
    name: string;
    role: string;
    retainer: number;
    outbid: boolean;
    detail: string;
  }[];
  arms?: {weapon: string; armour: string; charges: number};
  everyone?: Presence[];
  cast?: {
    id: string;
    name: string;
    role: string;
    temperament: string;
    manner: string;
    organization: string;
    known_for: string[];
  }[];
  grudges?: {holder: string; against: string; because: string}[];
  commissions?: {
    id: string;
    brief: string;
    giver: string;
    patron: string;
    progress: string;
    met: boolean;
    pay: number;
    respect: number;
    due: number;
    minutes_left: number;
  }[];
  residence?: {
    home: string;
    comforts: {id: string; label: string; upkeep: number}[];
    upkeep: number;
    sheltered: number;
    reach: number;
  };
  vehicle?: {
    car: string;
    condition: number;
    fuel?: number;
    tank?: number;
    pace: number;
    concealed: number;
    exposed: number;
    upkeep: number;
    running: boolean;
    plate?: number;
    plate_max?: number;
  };
  appearance?: {attire: string; condition: number; standing: number; presence: number};
  offshore?: {balance: number; reachable: boolean};
  arrangements?: {target: string; hired: string; paid: number; status: string}[];
  editions?: {
    day: number;
    life: number;
    dateline: string;
    count: number;
    current: boolean;
    mine: boolean;
    stories: {
      id: string;
      headline: string;
      body: string;
      kind: string;
      minute: number;
      day: number;
      time: string;
      weight?: number;
      standfirst?: string;
      byline?: string;
      dateline?: string;
      subject?: {kind: string; id?: string; name: string};
    }[];
  }[];
  newspaper?: {
    id: string;
    headline: string;
    body: string;
    kind: string;
    minute: number;
    day: number;
    time: string;
    weight?: number;
    standfirst?: string;
    byline?: string;
    dateline?: string;
    subject?: {kind: string; id?: string; name: string};
  }[];
  goods?: {id: string; name: string; unit: string; base: number; price: number; heat: number}[];
  conflicts?: {between: string[]; state: string; since: number}[];
  business_truces?: {[faction: string]: number};
  known_threats?: string[];
  opportunity?: {title: string; detail: string; target: string} | null;
  id: string;
  revision: number;
  life: number;
  minute: number;
  sky?: {kind: string; wet: number};
  player: Person;
  district: number;
  factions: {
    id: string;
    name: string;
    leader: string;
    power: number;
    goodwill: number;
    cash: number;
    strength: string;
    money: string;
    knowledge: number;
    leader_id?: string;
    holdings?: string[];
    people?: number;
    hands: string;
    fighting?: string[];
    standing?: string;
    yours?: boolean;
  }[];
  npcs: NPC[];
  locations: Place[];
  event: Event | null;
  history: Record[];
  dead: {name: string; minute: number; life: number; cause: string; estate?: string}[];
  tasks: {id: string; name: string; due: number}[];
  director: {status: string; detail: string; last_request: number};
  building_fires?: {id: string; target: string; minute: number; brigade_at: number; extinguished_at: number; cleanup_at: number}[];
  police_presence?: {id: string; target: string; minute: number; cleanup_at: number}[];
  aftermath?: {id: string; target: string; victim: {id: string; name: string}; minute: number; police_at: number; cleanup_at: number; cause?: 'charge-accident'; face?: number}[];
  last_result: {
    street_travel?: StreetSegment[] | null;
    comings?: Coming[];
    cues?: VisualCue[];
    action?: string;
    kind?: string;
    from_location: string;
    to_location: string;
    elapsed: number;
    records: Record[];
    cash: number;
    respect: number;
    heat: number;
    health: number;
  } | null;
  epitaph?: {
    name: string;
    cause: string;
    day: number;
    life: number;
    respect: number;
    earned: number;
    estate: string;
    became: string;
    standing: string[];
    inherits: string;
    headlines: string[];
  } | null;
  dashboard?: {
    id: string;
    label: string;
    value: string;
    note?: string;
    meaning: string;
    warn?: boolean;
  }[];
  guide?: {title: string; what: string; open: boolean; reason?: string; done?: boolean}[];
  rules?: string[];
  groups?: Group[];
  books?: {
    income: number;
    costs: number;
    net: number;
    lines: {label: string; amount: number; detail?: string}[];
    behind: {
      id: string;
      place: string;
      nights: number;
      hands: number;
      positions: number;
      shut: boolean;
    }[];
    cash: number;
    sheltered: number;
    lent: number;
    owed: number;
    offshore: number;
    holdings: number;
    earned: number;
  };
  daily_cost: number;
  income: number;
  security: number;
}
export interface Command {
  pool?:PoolInput;
  pool_game?:number;
  pool_host?:{fee:number;cut_percent:number;enter:boolean};
  kind: string;
  target?: string;
  event?: string;
  choice?: string;
  amount?: number;
  chips?: {bet: string; amount: number}[];
  request_id?: string;
  revision?: number;
}
