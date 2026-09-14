import {useEffect,useRef,useState} from 'react';
import type {Snapshot} from './types';
import {MapMenu} from './MapMenu';
import {PressCut} from './Herald';
import {playMoment} from './sound';
export function SceneNewspaper({article,visible,voice,onClose}:{article:NonNullable<Snapshot['newspaper']>[number];visible:boolean;voice:boolean;onClose:()=>void}){
 const [clip,setClip]=useState<Blob|null>(null),[status,setStatus]=useState('Preparing narration…');
 const [retry,setRetry]=useState(0);const player=useRef<HTMLAudioElement|null>(null);
 useEffect(()=>{
  setClip(null);if(!voice)return;
  const abort=new AbortController();setStatus('Preparing narration…');
  fetch('/api/newspaper/speech',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({story:article.id}),signal:abort.signal})
   .then(async response=>{if(!response.ok)throw Error('unavailable');return response.blob();})
   .then(blob=>{if(!abort.signal.aborted){setClip(blob);setStatus('Narration ready');}})
   .catch(()=>{if(!abort.signal.aborted)setStatus('Narration unavailable');});
  return()=>abort.abort();
 },[article.id,article.headline,article.body,voice,retry]);
 useEffect(()=>{if(visible)return playMoment('newspaper');},[visible]);
 useEffect(()=>{
  if(!visible||!voice||!clip)return;
  let disposed=false;const url=URL.createObjectURL(clip),audio=new Audio(url);player.current=audio;
  audio.onended=()=>{if(!disposed)setStatus('Narration finished');};
  audio.play().then(()=>{if(!disposed)setStatus('Narrating…');}).catch(()=>{if(!disposed)setStatus('Play narration');});
  return()=>{disposed=true;audio.pause();audio.onended=null;URL.revokeObjectURL(url);if(player.current===audio)player.current=null;};
 },[visible,voice,clip]);
 if(!visible)return null;
 return <MapMenu title="The Bellwether Herald" onClose={onClose}>
  <div className="scene-newspaper">
   <div className="paper-sheet"><div className="paper">
    <div className="masthead"><h1>The Bellwether Herald</h1><div className="rule"><span>{article.dateline||'Bellwether'}</span><span>Late city edition</span><span>Day {article.day} · Five cents</span></div></div>
    <article className="scene-article">
     {article.subject&&<PressCut kind={article.kind} subject={article.subject} headline={article.headline}/>}
     <h2>{article.headline}</h2>{article.standfirst&&<p className="standfirst">{article.standfirst}</p>}
     <p className="byline">{article.byline||'Herald city desk'} · {article.time}</p><p className="body">{article.body}</p>
    </article>
   </div></div>
   {voice&&<div className="newspaper-voice" role="status"><span>{status}</span><button onClick={()=>{if(player.current){void player.current.play().catch(()=>setStatus('Play narration'));}else setRetry(n=>n+1);}}>Read aloud</button><button onClick={()=>{player.current?.pause();setStatus('Narration paused');}}>Pause</button></div>}
   <button className="plain" onClick={onClose}>Return to the city →</button>
  </div>
 </MapMenu>;
}
