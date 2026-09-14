import {useEffect,useRef,useState} from 'react';
import type {Snapshot} from './types';
import {PressCut} from './Herald';
import {playMoment} from './sound';
export function SceneNewspaper({article,visible,voice,onClose}:{article:NonNullable<Snapshot['newspaper']>[number];visible:boolean;voice:boolean;onClose:()=>void}){
 const [clip,setClip]=useState<Blob|null>(null),[status,setStatus]=useState('Preparing narration…');
 const [timing,setTiming]=useState<{preparedMs:number|null;startDelayMs:number|null}>({preparedMs:null,startDelayMs:null});
 const revealed=useRef(0),panel=useRef<HTMLElement>(null);
 useEffect(()=>{
  if(!visible)return;
  const previous=document.activeElement as HTMLElement|null;
  panel.current?.focus({preventScroll:true});
  return()=>{if(previous?.isConnected)previous.focus({preventScroll:true});};
 },[visible]);
 const [retry,setRetry]=useState(0);const player=useRef<HTMLAudioElement|null>(null);
 useEffect(()=>{
  setClip(null);if(!voice)return;
  const began=performance.now();setTiming({preparedMs:null,startDelayMs:null});
  const abort=new AbortController();setStatus('Preparing narration…');
  fetch('/api/newspaper/speech',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({story:article.id}),signal:abort.signal})
   .then(async response=>{if(!response.ok)throw Error('unavailable');return response.blob();})
   .then(blob=>{if(!abort.signal.aborted){setClip(blob);setTiming(t=>({...t,preparedMs:Math.round(performance.now()-began)}));setStatus('Narration ready');}})
   .catch(()=>{if(!abort.signal.aborted)setStatus('Narration unavailable');});
  return()=>abort.abort();
 },[article.id,article.headline,article.body,voice,retry]);
 useEffect(()=>{if(visible){revealed.current=performance.now();return playMoment('newspaper');}},[visible]);
 useEffect(()=>{
  if(!visible||!voice||!clip)return;
  let disposed=false;const url=URL.createObjectURL(clip),audio=new Audio(url);player.current=audio;
  audio.onended=()=>{if(!disposed)setStatus('Narration finished');};
  audio.play().then(()=>{if(!disposed){setStatus('Narrating…');setTiming(t=>({...t,startDelayMs:Math.round(performance.now()-revealed.current)}));}}).catch(()=>{if(!disposed)setStatus('Play narration');});
  return()=>{disposed=true;audio.pause();audio.onended=null;URL.revokeObjectURL(url);if(player.current===audio)player.current=null;};
 },[visible,voice,clip]);
 const diagnostic=JSON.stringify({article:article.id,visible,status,...timing});
 if(!visible)return <span hidden data-news-audio={diagnostic}/>;
 return <div className="newspaper-backdrop" onClick={e=>{if(e.target===e.currentTarget)onClose();}}>
  <section ref={panel} className="scene-newspaper" role="dialog" aria-modal="true" aria-label="The Bellwether Herald" tabIndex={-1} data-news-audio={diagnostic} onKeyDown={e=>{
   if(e.key==='Escape'){e.stopPropagation();onClose();}
   if(e.key==='Tab'){
    const items=[...e.currentTarget.querySelectorAll<HTMLButtonElement>('button:not(:disabled)')];
    const index=items.indexOf(document.activeElement as HTMLButtonElement);
    if(e.shiftKey&&index<=0){e.preventDefault();items.at(-1)?.focus();}
    else if(!e.shiftKey&&(index<0||index===items.length-1)){e.preventDefault();items[0]?.focus();}
   }
  }}>
   <div className="paper-sheet"><div className="paper">
    <div className="newspaper-tools"><span>From the streets of Bellwether</span><button onClick={onClose} aria-label="Close The Bellwether Herald">Fold away <span aria-hidden="true">×</span></button></div>
    <div className="masthead"><h1>The Bellwether Herald</h1><div className="rule"><span>{article.dateline||'Bellwether'}</span><span>Late city edition</span><span>Day {article.day} · Five cents</span></div></div>
    <article className="scene-article">
     {article.subject&&<PressCut kind={article.kind} subject={article.subject} headline={article.headline}/>}
     <h2>{article.headline}</h2>{article.standfirst&&<p className="standfirst">{article.standfirst}</p>}
     <p className="byline">{article.byline||'Herald city desk'} · {article.time}</p><p className="body">{article.body}</p>
    </article>
   <footer className="newspaper-footer">
   {voice&&<div className="newspaper-voice"><span role="status">{status}</span><button disabled={status==='Preparing narration…'||status==='Narrating…'} onClick={()=>{if(player.current){const audio=player.current;void audio.play().then(()=>{if(player.current===audio)setStatus('Narrating…');}).catch(()=>{if(player.current===audio)setStatus('Play narration');});}else setRetry(n=>n+1);}}>{status==='Narration paused'?'Resume reading':'Read aloud'}</button><button disabled={status!=='Narrating…'} onClick={()=>{player.current?.pause();setStatus('Narration paused');}}>Pause</button></div>}
   <span className="newspaper-imprint">The Bellwether Herald · Independent city press</span>
   </footer>
   </div></div>
  </section>
 </div>;
}
