import type {Snapshot,Place} from "./types";
const esc=(s:unknown)=>String(s??'').replace(/[&<>"\']/g,c=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","\'":"&#39;"}[c]||c));
const paths:Record<string,string>={city:'M3 20h18M5 20V9h5v11M12 20V4h6v16M18 12h3v8M6 12h2m5-5h3m-3 4h3m-3 4h3',ledger:'M5 3h14v18H5zM8 7h8M8 11h8M8 15h5',crew:'M15 8a3 3 0 1 1-6 0 3 3 0 0 1 6 0M5 21v-3a7 7 0 0 1 14 0v3M19 5a3 3 0 0 1 1 6',families:'M12 3l9 4v6c0 4-9 8-9 8s-9-4-9-8V7zM8 12h8M12 8v8',help:'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20M9 8a3 3 0 1 1 4 3l-1 2M12 17h.01',cash:'M3 6h18v12H3zM14 12a2 2 0 1 1-4 0 2 2 0 0 1 4 0',health:'M12 20S2 14 2 8c0-5 7-6 10-1 3-5 10-4 10 1 0 6-10 12-10 12',respect:'M12 2l3 6 7 1-5 5 1 7-6-3-6 3 1-7-5-5 7-1z',heat:'M12 2c4 6-1 7 3 10 1-3 3-4 3-4s6 9-1 13C2 25 1 9 12 2',voice:'M3 9h4l6-5v16l-6-5H3zM17 7q7 5 0 10',mute:'M3 9h4l6-5v16l-6-5H3zM17 8l5 8m0-8-5 8',clock:'M22 12a10 10 0 1 0-20 0 10 10 0 0 0 20 0M12 6v6l4 3',key:'M10 9a4 4 0 1 0-8 0 4 4 0 0 0 8 0M10 9h12m-4 0v4m-4-4v3',arrow:'M4 12h16m-6-6 6 6-6 6',news:'M3 4h15v16H3zM18 8h4v12H3M6 8h9M6 12h4m2 0h3M6 16h9',settings:'M12 3v3m0 12v3M3 12h3m12 0h3M5 5l2 2m10 10 2 2M5 19l2-2M17 7l2-2M16 12a4 4 0 1 0-8 0 4 4 0 0 0 8 0',held:'M4 3v18M9 3v18M15 3v18M20 3v18M2 8h20M2 16h20'};
export function icon(id:string){return `<svg class="icon" viewBox="0 0 24 24" aria-hidden="true"><path d="${paths[id]||paths.city}"/></svg>`}
// Until the painted faces cover the whole cast, everybody else is drawn from
// their own id: the same person is always the same face. The old version gave
// every one of them the same head, the same haircut and the same collar in one
// of four skin tones, so a room of fifteen read as one man repeated. This
// varies the things a face is actually recognised by — the shape of the head,
// the hairline, the brow, whether there is a hat or a collar or a moustache —
// and keeps the palette the painted portraits use so the two sit together.
const skins = ['#c3a184', '#a67c5b', '#8a6247', '#d0b697', '#6f4c37', '#b98f6a'];
const hairs = ['#26221d', '#3c332a', '#5a4632', '#7a6a52', '#8e8375', '#1c1a17'];
const coats = ['#2b2a26', '#33302a', '#26292b', '#3a332c', '#242524'];

export function portrait(id: string, size = '') {
  let h = 2166136261;
  for (const c of id) { h ^= c.charCodeAt(0); h = Math.imul(h, 16777619) }
  const n = (bit: number, mod: number) => Math.abs(h >> (bit % 24)) % mod;

  const skin = skins[n(0, skins.length)];
  const hair = hairs[n(3, hairs.length)];
  const coat = coats[n(6, coats.length)];
  const wide = 20 + n(9, 5);            // how broad the head is
  const jaw = 34 + n(11, 8);            // where the chin falls
  const brow = n(13, 4);                // heavy brows or not
  const style = n(15, 5);               // hairline
  const hat = n(17, 5) === 0;
  const collar = n(19, 3);
  const tache = n(21, 6) === 0;
  const light = n(2, 2) ? 1 : -1;       // which side the lamp is on

  const cx = 32;
  return `<svg class="portrait ${size}" viewBox="0 0 64 78" aria-hidden="true">` +
    `<defs><radialGradient id="g-${id.replace(/[^a-z0-9]/gi, '')}" cx="${light > 0 ? .32 : .68}" cy=".3" r=".9">` +
    `<stop stop-color="#6a6152"/><stop offset="1" stop-color="#1b1a17"/></radialGradient></defs>` +
    `<rect width="64" height="78" fill="url(#g-${id.replace(/[^a-z0-9]/gi, '')})"/>` +
    // shoulders and coat
    `<path d="M2 78 q6-16 16-19 l14-4 14 4 q10 3 16 19z" fill="${coat}"/>` +
    (collar === 0
      ? `<path d="M26 55 l6 12 6-12 4 2-8 16h-4l-8-16z" fill="#cfc6ae"/><path d="M31 62l2 5 2-5-2-2z" fill="#4a4038"/>`
      : collar === 1
        ? `<path d="M25 56 l7 10 7-10 3 3-10 14-10-14z" fill="#b9b09a"/>`
        : `<path d="M24 57 q8 6 16 0 l2 4 q-10 7-20 0z" fill="#8e8a7a"/>`) +
    // neck and head
    `<path d="M27 44h10v10q-5 4-10 0z" fill="${skin}" opacity=".85"/>` +
    `<path d="M${cx - wide / 2} 26 q0-16 ${wide / 2} -16 t${wide / 2} 16 q0 ${jaw / 2} -${wide / 2} ${jaw / 2} q-${wide / 2} 0 -${wide / 2} -${jaw / 2}z" fill="${skin}"/>` +
    // the lamp on one cheek
    `<path d="M${cx} 10 q${light * wide / 2} 0 ${light * wide / 2} 16 q0 ${jaw / 2} -${light * wide / 2} ${jaw / 2}z" fill="#000" opacity=".22"/>` +
    // hair
    (style === 0 ? `<path d="M${cx - wide / 2 - 1} 26 q1-19 ${wide / 2 + 1} -19 t${wide / 2 + 1} 19 q-3-10 -${wide / 2 + 1} -10 q-${wide / 2 + 1} 0 -${wide / 2 + 1} 10z" fill="${hair}"/>`
      : style === 1 ? `<path d="M${cx - wide / 2 - 1} 24 q2-18 ${wide / 2 + 1} -18 q6 0 8 6 q-9-1 -14 6 q-3 2 -3 6z" fill="${hair}"/>`
      : style === 2 ? `<path d="M${cx - wide / 2 - 1} 27 q0-20 ${wide / 2 + 1} -20 t${wide / 2 + 1} 20 q-2-13 -${wide / 2 + 1} -13 t-${wide / 2 + 1} 13z" fill="${hair}"/><path d="M${cx - wide / 2 - 2} 26 q-2 10 1 16 q-4-6 -3-16z" fill="${hair}"/><path d="M${cx + wide / 2 + 1} 26 q2 10 -1 16 q4-6 3-16z" fill="${hair}"/>`
      : style === 3 ? `<path d="M${cx - wide / 2} 22 q4-15 ${wide / 2} -15 t${wide / 2} 15 q-4-8 -${wide / 2} -8 t-${wide / 2} 8z" fill="${hair}"/>`
      : `<path d="M${cx - wide / 2 - 1} 25 q1-18 ${wide / 2 + 1} -18 t${wide / 2 + 1} 18 q-1-8 -6-11 q-8 5 -16 3 q-3 3 -3 8z" fill="${hair}"/>`) +
    // brow, eyes, mouth
    (brow === 0 ? '' : `<path d="M${cx - 9} 27h6m6 0h6" stroke="#2a231c" stroke-width="${1 + brow * .5}" opacity=".8"/>`) +
    `<ellipse cx="${cx - 6}" cy="31" rx="2.1" ry="1.5" fill="#241e19"/>` +
    `<ellipse cx="${cx + 6}" cy="31" rx="2.1" ry="1.5" fill="#241e19"/>` +
    `<path d="M${cx} 32 v5 l-2 2" stroke="#00000055" fill="none"/>` +
    (tache ? `<path d="M${cx - 5} 40 q5-2 10 0 q-5 2 -10 0z" fill="${hair}"/>` : '') +
    `<path d="M${cx - 4} ${43 + (jaw - 34) / 4} q4 2 8 0" stroke="#5c4234" stroke-width="1.2" fill="none" opacity=".85"/>` +
    (hat ? `<path d="M${cx - wide / 2 - 5} 20 h${wide + 10} l-3-4 h-${wide + 4}z" fill="#1e1c19"/>` +
      `<path d="M${cx - wide / 2 - 1} 16 q${wide / 2 + 1}-11 ${wide + 2} 0z" fill="#25221e"/>` : '') +
    `<rect x=".5" y=".5" width="63" height="77" fill="none" stroke="#000" opacity=".45"/>` +
    `</svg>`;
}

export function building(p:Place,mini=false){let x=mini?100:p.x,y=mini?75:p.y,w=p.type==='casino'?78:p.type==='home'?52:62,h=p.type==='casino'?52:40,z=p.type==='home'?39:25;let c=p.owned?'#7e8255':p.owner==='bellandi'?'#786359':p.type==='home'?'#6d7c71':'#737e69';let windows='';for(let i=0;i<4;i++)windows+=`<rect x="${x-w/2+7+i*(w-13)/4}" y="${y-z+3}" width="7" height="9" fill="${p.locked?'#566358':'#d4b773'}" opacity=".72"/>`;return `<g><path d="M${x-w/2+8} ${y+h/2+6}l${w+16} -8 13 -${h} -${w} -8z" fill="#071312" opacity=".38"/><rect class="footprint" x="${x-w/2}" y="${y-h/2}" width="${w}" height="${h}" rx="2" fill="#293c33" stroke="#56654f"/><path d="M${x-w/2} ${y-z}h${w}v${h/2+z}h-${w}z" fill="${c}"/><path d="M${x+w/2} ${y-z}l14 -10v${h/2+z}l-14 10z" fill="#394c40"/><path d="M${x-w/2} ${y-z}l14 -10h${w}l-14 10z" fill="#8c9177"/><path d="M${x-w/2+5} ${y-z-3}h${w-10}l7 -4h-${w-10}z" fill="#475b4b"/>${windows}<rect x="${x-5}" y="${y+4}" width="10" height="15" fill="#172926"/>${p.type==='casino'?`<rect x="${x-w/2+3}" y="${y-10}" width="${w-6}" height="10" fill="#bb8c5e"/><path d="M${x-w/2-5} ${y+2}h${w+10}l-6 8h-${w-2}z" fill="#6d4540"/>`:''}${p.condition<70?`<path d="M${x-13} ${y-z}l18 22-8 12 13 13" stroke="#17201c" stroke-width="5"/>`:''}</g>`}
export function citySVG(state:Snapshot,selected:string){let roads=`<path d="M55 300H1020M70 515H970M215 150V665M445 135V665M670 100V665M885 120V665" stroke="#34443a" stroke-width="32" fill="none"/><path d="M55 300H1020M70 515H970M215 150V665M445 135V665M670 100V665M885 120V665" stroke="#819076" stroke-width="1" opacity=".45" fill="none" stroke-dasharray="7 10"/>`;
let blocks='',trees='';for(let i=0;i<55;i++){let x=75+(i*137)%920,y=155+(i*83)%470;if(state.locations.some(p=>Math.abs(p.x-x)<80&&Math.abs(p.y-y)<65)||[215,445,670,885].some(v=>Math.abs(v-x)<35)||[300,515].some(v=>Math.abs(v-y)<35))continue;blocks+=`<rect x="${x}" y="${y}" width="${20+i%4*7}" height="${19+i%3*6}" fill="#3e5143" stroke="#64705a" opacity=".7"/><rect x="${x+3}" y="${y-4}" width="${20+i%4*7}" height="${19+i%3*6}" fill="#53614e" opacity=".6"/>`}
for(let i=0;i<38;i++){let x=270+(i*73)%720,y=130+(i*137)%525;trees+=`<circle cx="${x+3}" cy="${y+3}" r="6" fill="#101f1d"/><circle cx="${x}" cy="${y}" r="6" fill="#52644b"/>`}
let current=state.locations.find(l=>l.id===state.player.location)||state.locations[0];return `<svg class="map" viewBox="0 0 1080 740" role="group" aria-label="Map of Bellwether. Select a property to inspect its actions."><defs><pattern id="water" width="38" height="25" patternUnits="userSpaceOnUse"><path d="M3 14h14m5-5h10" class="waterline"/></pattern><radialGradient id="cityglow"><stop stop-color="#4d5d40" stop-opacity=".5"/><stop offset="1" stop-color="#182725" stop-opacity="0"/></radialGradient></defs><rect width="1080" height="740" fill="#20322f"/><path d="M0 0h185L70 130 50 340 100 460 48 570 100 740H0" fill="#14282c"/><path d="M0 0h185L70 130 50 340 100 460 48 570 100 740H0" fill="url(#water)"/><path d="M65 585H1050V680H125zM110 140H1020V275H92zM80 340H1020V480H130z" fill="#2c3e31"/>${roads}${blocks}${trees}<ellipse cx="480" cy="370" rx="390" ry="295" fill="url(#cityglow)"/><text class="district" x="320" y="115">OLD HARBOR</text><text class="district" x="735" y="110">ASHBURY</text><text class="district" x="710" y="700">THE HEIGHTS</text><text x="35" y="620" fill="#74918d" font-size="11" letter-spacing="4" transform="rotate(-90 35 620)">BELLWETHER BAY</text>${state.locations.map(p=>`<g class="property ${p.id===selected?'selected':''} ${p.locked?'locked':''}" data-place="${p.id}" role="button" tabindex="0" aria-label="${esc(p.name)}${p.locked?', district not yet accessible':''}"><rect x="${p.x-65}" y="${p.y-65}" width="130" height="112" fill="transparent"/>${building(p)}<text class="property-label" x="${p.x}" y="${p.y+43}" text-anchor="middle">${esc(p.name)}</text>${p.owned?`<text class="property-sub" x="${p.x}" y="${p.y+57}" text-anchor="middle">YOUR BUSINESS</text>`:p.owner==='bellandi'?`<text class="property-sub" x="${p.x}" y="${p.y+57}" text-anchor="middle">BELLANDI TERRITORY</text>`:''}</g>`).join('')}<g id="player-marker" transform="translate(${current.x},${current.y+24})"><circle r="12" fill="#101b19" stroke="#e1c083" stroke-width="2"/><circle r="4" fill="#f6dfae"/><path d="M0-18v-8" stroke="#f6dfae" stroke-width="2"/></g><text x="1030" y="705" fill="#a0a88f" font-size="12">N ↑</text></svg>`}

// A newspaper picture, drawn rather than photographed: a coarse black and white
// block of whatever the story is about, screened with halftone dots the way a
// 1930s press would have printed it. Deterministic from the story, so the same
// story always carries the same picture.
type PressSubject = {kind:string;id?:string;name:string};

function pressHash(seed:string){let h=2166136261;for(const c of seed){h^=c.charCodeAt(0);h=Math.imul(h,16777619)}return Math.abs(h)}

function halftone(seed:number,density:number){
  let dots='';
  for(let y=0;y<11;y++)for(let x=0;x<16;x++){
    const n=(seed>>((x+y*3)%23))&7;
    if(n/8>density)continue;
    const r=.7+((seed>>((x*2+y)%17))&3)*.35;
    dots+=`<circle cx="${5+x*6}" cy="${5+y*6}" r="${r.toFixed(2)}"/>`;
  }
  return `<g fill="#14100c" opacity=".55">${dots}</g>`;
}

// The shapes are blocky on purpose: a press this old could not hold detail, and
// a silhouette reads at this size where a drawing would not.
function pressShape(kind:string,subject:PressSubject,seed:number){
  const lean=(seed%7)-3;
  switch(subject.kind){
    case 'person':
      return `<g fill="#14100c"><path d="M48 74q0-14 16-18l-3-6q-6-4-6-13 0-13 13-13t13 13q0 9-6 13l-3 6q16 4 16 18z"/></g>`+
        (kind==='killing'?`<g stroke="#14100c" stroke-width="3" fill="none"><path d="M22 70l24-16m0 16L22 54"/></g>`:'');
    case 'place':
      return `<g fill="#14100c"><path d="M20 74V34l25-14 25 14v40z"/><path d="M78 74V46l14-8 14 8v28z"/></g>`+
        `<g fill="#efe7d8"><rect x="30" y="42" width="8" height="9"/><rect x="44" y="42" width="8" height="9"/><rect x="30" y="57" width="8" height="9"/><rect x="44" y="57" width="8" height="9"/><rect x="86" y="54" width="7" height="8"/></g>`+
        (kind==='attack'?`<g stroke="#14100c" stroke-width="3" fill="none"><path d="M45 20l-9 16 12 4-10 14"/></g>`:'')+
        (kind==='police'?`<g fill="#14100c"><rect x="8" y="62" width="30" height="12" rx="3"/><circle cx="16" cy="76" r="4"/><circle cx="31" cy="76" r="4"/></g>`:'');
    default:
      return `<g fill="#14100c"><path d="M6 74V44l12-7v37zM24 74V34l16-9v49zM46 74V22l18-10v62zM70 74V38l14-8v44zM90 74V48l14-7v33z"/></g>`+
        `<g fill="#efe7d8"><rect x="50" y="30" width="5" height="7"/><rect x="58" y="42" width="5" height="7"/><rect x="28" y="46" width="5" height="7"/></g>`;
  }
  void lean;
}

export function pressPlate(kind:string,subject:PressSubject,headline:string){
  const seed=pressHash(headline+subject.kind+(subject.id||''));
  const sky=(seed%3)*.06;
  return `<svg class="press-plate" viewBox="0 0 112 80" role="img" aria-label="${esc(subject.name)}">`+
    `<rect width="112" height="80" fill="#efe7d8"/>`+
    halftone(seed,.34+sky)+
    pressShape(kind,subject,seed)+
    `<rect x="1" y="1" width="110" height="78" fill="none" stroke="#14100c" stroke-width="2"/>`+
    `</svg>`;
}
