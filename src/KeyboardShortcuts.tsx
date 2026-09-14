import {useEffect,useRef,useState} from 'react';
import {cameraCommand} from './city3dControls';
import './keyboardShortcuts.css';

type Entry={element:HTMLElement;label:string};
function available(element:HTMLElement){for(let parent=element.parentElement;parent;parent=parent.parentElement){if(parent.matches('details:not([open])')&&!parent.querySelector(':scope > summary')?.contains(element))return false;}return !element.closest('[inert],[hidden],.shortcut-shade')&&!element.matches(':disabled')&&element.getClientRects().length>0&&getComputedStyle(element).visibility!=='hidden';}
function scope(){return [...document.querySelectorAll<HTMLElement>('[role="dialog"][aria-modal="true"]')].filter(available).at(-1)||document.body;}
function actions():Entry[]{return [...scope().querySelectorAll<HTMLElement>('button,summary,select,input:not([type="hidden"]),textarea,a[href]')].filter(available).map(element=>{let label=(element.getAttribute('aria-label')||element.textContent||element.getAttribute('placeholder')||element.getAttribute('title')||'Input').replace(/\s+/g,' ').trim();const title=element.getAttribute('title');if(title&&!label.includes(title))label+=' — '+title;return{element,label};});}
function activate(element:HTMLElement){if(!element.isConnected||!available(element))return;element.focus({preventScroll:true});element.scrollIntoView({block:'nearest'});if(element.matches('select,input,textarea'))return;element.click();}
const hints=[['1–6','People · Families · Market · Ledger · Herald · Guide'],['7 / 8','Settings / daily accounts'],['V','Toggle voices'],['T','Toggle travel playback speed'],['G','Travel to selected address / step inside'],['B','Leave the building'],['F / O / Z','Follow character / whole city / selected address'],['J','Focus the address directory (type a name, then Enter)'],['WASD / arrows','Pan the visible table, room or city camera'],['Q / E','Rotate camera'],['+ / − / wheel','Zoom; stop following'],['Home','Reset the visible camera'],['/','Search available actions; arrows to select, Enter to use'],['Tab / Shift Tab','Move between controls; Enter or Space to use'],['Esc','Close the current paper or menu'],['?','This keyboard reference']];
export function KeyboardShortcuts(){
 const [mode,setMode]=useState<'actions'|'help'|null>(null),[entries,setEntries]=useState<Entry[]>([]),[query,setQuery]=useState(''),[index,setIndex]=useState(0);
 const current=useRef(mode);current.current=mode;
 const previous=useRef<HTMLElement|null>(null),panel=useRef<HTMLElement>(null),input=useRef<HTMLInputElement>(null);
 const close=()=>{setMode(null);previous.current?.isConnected&&previous.current.focus({preventScroll:true});};
 useEffect(()=>{
  const key=(event:KeyboardEvent)=>{
   if(event.defaultPrevented||event.repeat||event.isComposing||event.ctrlKey||event.metaKey||event.altKey||current.current)return;
   const target=event.target as HTMLElement;
   if(target?.isContentEditable||target?.closest('input,textarea,select,[role="textbox"]'))return;
   if(event.key==='/'||event.key==='?'){
    event.preventDefault();previous.current=document.activeElement as HTMLElement;setEntries(actions());setQuery('');setIndex(0);setMode(event.key==='/'?'actions':'help');return;
   }
   const activeScope=scope();
   const shortcut=activeScope===document.body?[...document.querySelectorAll<HTMLElement>('[data-shortcut]')].find(el=>el.dataset.shortcut===event.key.toLowerCase()&&available(el)):undefined;
   if(shortcut){event.preventDefault();activate(shortcut);return;}
   if(cameraCommand(event)){
    const canvas=['.poker3d canvas','.blackjack3d canvas','.pool-render canvas','.interior3d canvas','.city3d canvas'].flatMap(selector=>[...activeScope.querySelectorAll<HTMLElement>(selector)]).find(available);
    if(canvas&&available(canvas)&&target!==canvas){event.preventDefault();canvas.focus({preventScroll:true});canvas.dispatchEvent(new KeyboardEvent('keydown',{key:event.key,code:event.code,shiftKey:event.shiftKey}));}
   }
  };
  document.addEventListener('keydown',key);return()=>document.removeEventListener('keydown',key);
 },[]);
 useEffect(()=>{if(mode==='actions')input.current?.focus();else if(mode)panel.current?.focus();},[mode]);
 const filtered=entries.filter(e=>query.toLowerCase().split(/\s+/).every(word=>e.label.toLowerCase().includes(word)));
 const chosen=Math.min(index,Math.max(0,filtered.length-1));
 return <>{!mode&&<button className="shortcut-reference" aria-label="Keyboard shortcuts" onClick={()=>{previous.current=document.activeElement as HTMLElement;setMode('help');}}><kbd>?</kbd> Keys</button>}{mode&&<div className="shortcut-shade" onClick={e=>{if(e.target===e.currentTarget)close();}}><section ref={panel} tabIndex={-1} className="shortcut-paper" role="dialog" aria-modal="true" aria-label={mode==='help'?'Keyboard shortcuts':'Available actions'} onKeyDown={e=>{
  if(e.key==='Escape'){e.preventDefault();e.stopPropagation();close();}
  if(e.key==='Tab'){const items=[...e.currentTarget.querySelectorAll<HTMLElement>('button,input')];const at=items.indexOf(document.activeElement as HTMLElement);if(e.shiftKey&&at<=0){e.preventDefault();items.at(-1)?.focus();}else if(!e.shiftKey&&(at<0||at===items.length-1)){e.preventDefault();items[0]?.focus();}}
 }}><header><h2>{mode==='help'?'A pocket guide to the keys':'What would you like to do?'}</h2><button onClick={close} aria-label="Close keyboard panel">Put away ×</button></header>{mode==='help'?<dl>{hints.map(([keys,what])=><div key={keys}><dt><kbd>{keys}</kbd></dt><dd>{what}</dd></div>)}</dl>:<><label htmlFor="shortcut-search">Search actions available here</label><input id="shortcut-search" ref={input} value={query} autoComplete="off" aria-controls="shortcut-results" aria-activedescendant={filtered.length?`shortcut-${chosen}`:undefined} role="combobox" aria-expanded="true" onChange={e=>{setQuery(e.target.value);setIndex(0);}} onKeyDown={e=>{
 if(e.key==='ArrowDown'||e.key==='ArrowUp'){e.preventDefault();setIndex(Math.max(0,Math.min(filtered.length-1,chosen+(e.key==='ArrowDown'?1:-1))));}
 if(e.key==='Enter'&&!e.nativeEvent.isComposing){e.preventDefault();if(filtered[chosen]){const element=filtered[chosen].element;close();activate(element);}}
 }}/><ul id="shortcut-results" role="listbox">{filtered.map((entry,i)=><li id={`shortcut-${i}`} key={i} role="option" aria-selected={i===chosen}><button tabIndex={-1} onMouseEnter={()=>setIndex(i)} onClick={()=>{close();activate(entry.element);}} ref={el=>{if(i===chosen)el?.scrollIntoView({block:'nearest'});}}>{entry.label}</button></li>)}</ul>{!filtered.length&&<p>No available actions match that search.</p>}<p className="shortcut-note">↑ ↓ Choose · Enter Use · Esc Put away. Actions keep their usual costs and consequences.</p></>}</section></div>}</>;
}
