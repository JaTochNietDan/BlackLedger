import {useEffect,useRef,type ReactNode} from 'react';
export function MapMenu({title,onClose,children}:{title:string;onClose:()=>void;children:ReactNode}){
 const panel=useRef<HTMLElement>(null),close=useRef(onClose);close.current=onClose;
 useEffect(()=>{
  const previous=document.activeElement as HTMLElement|null;
  panel.current?.focus({preventScroll:true});
  return()=>{if(previous?.isConnected)previous.focus({preventScroll:true});};
 },[]);
 return <div className="map-menu-shade" onClick={e=>{if(e.target===e.currentTarget)close.current();}}>
  <section ref={panel} className="map-menu" role="dialog" aria-modal="true" aria-label={title} tabIndex={-1} onKeyDown={e=>{
   if(e.key==='Escape'){e.stopPropagation();close.current();}
   if(e.key==='Tab'){
    const items=[...e.currentTarget.querySelectorAll<HTMLElement>('button:not(:disabled),a[href],input,select,textarea,[tabindex="0"]')].filter(el=>el.getClientRects().length);
    const index=items.indexOf(document.activeElement as HTMLElement);
    if(e.shiftKey&&index<=0){e.preventDefault();items.at(-1)?.focus();}
    else if(!e.shiftKey&&(index<0||index===items.length-1)){e.preventDefault();items[0]?.focus();}
   }
  }}>
   <header className="map-menu-title"><span>{title}</span><button onClick={onClose} aria-label={`Close ${title}`}>Close ×</button></header>
   <div className="map-menu-body">{children}</div>
  </section>
 </div>;
}
