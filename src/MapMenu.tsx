import {useEffect,useRef,type ReactNode} from 'react';
import './documentMenu.css';
const editions:Record<string,{label:string;mark:string}>={
 crew:{label:'Personal address book',mark:'Contacts'},families:{label:'Private family dossiers',mark:'Confidential'},
 market:{label:'Mercer Exchange · Trade circular',mark:'Market quotations'},ledger:{label:'Private accounts & memoranda',mark:'Black Ledger'},
 news:{label:'The Bellwether Herald · Reading room',mark:'Press archive'},settings:{label:'Black Ledger · Control desk',mark:'Preferences'},
 help:{label:'A newcomer’s handbook',mark:'Bellwether'},
};
export function MapMenu({title,onClose,children,edition}:{title:string;onClose:()=>void;children:ReactNode;edition?:string}){
 const design=editions[edition||'']||{label:'Private papers',mark:'Bellwether'};
 const panel=useRef<HTMLElement>(null),close=useRef(onClose);close.current=onClose;
 useEffect(()=>{
  const previous=document.activeElement as HTMLElement|null;
  panel.current?.focus({preventScroll:true});
  return()=>{if(previous?.isConnected)previous.focus({preventScroll:true});};
 },[]);
 return <div className="map-menu-shade" onClick={e=>{if(e.target===e.currentTarget)close.current();}}>
  <section ref={panel} className={`map-menu document-menu document-${edition||'papers'}`} role="dialog" aria-modal="true" aria-label={title} tabIndex={-1} onKeyDown={e=>{
   if(e.key==='Escape'){e.stopPropagation();close.current();}
   if(e.key==='Tab'){
    const items=[...e.currentTarget.querySelectorAll<HTMLElement>('button:not(:disabled),a[href],input,select,textarea,[tabindex="0"]')].filter(el=>el.getClientRects().length);
    const index=items.indexOf(document.activeElement as HTMLElement);
    if(e.shiftKey&&index<=0){e.preventDefault();items.at(-1)?.focus();}
    else if(!e.shiftKey&&(index<0||index===items.length-1)){e.preventDefault();items[0]?.focus();}
   }
  }}>
   <header className="map-menu-title"><div className="document-heading"><span className="document-tab">{title}</span><span className="document-label">{design.label}</span></div><button onClick={onClose} aria-label={`Close ${title}`}>Put away <span aria-hidden="true">×</span></button></header>
   <div className="map-menu-body">{children}</div>
   <footer className="document-footer"><span>{design.mark}</span><span>Bellwether · 1953</span></footer>
  </section>
 </div>;
}
