export interface VoiceHandle {play():Promise<void>;pause():void;dispose():void;onEnded(fn:()=>void):void}
export interface VoiceDependencies {load(event:string,signal:AbortSignal):Promise<Blob>;audio(blob:Blob):VoiceHandle;status(text:string):void;unavailable():void}
/** Owns cancellation and media lifetime, independently of React or model latency. */
export class VoicePlayer {
 private epoch=0;private request:AbortController|null=null;private playing:VoiceHandle|null=null;
 constructor(private deps:VoiceDependencies){}
 stop(){this.epoch++;this.request?.abort();this.request=null;if(this.playing){this.playing.pause();this.playing.dispose();this.playing=null}this.deps.status('Read aloud')}
 async speak(event:string){this.stop();const token=this.epoch,abort=new AbortController();this.request=abort;this.deps.status('Preparing voice…');try{
  const blob=await this.deps.load(event,abort.signal);if(token!==this.epoch)return;
  const player=this.deps.audio(blob);this.playing=player;player.onEnded(()=>{if(token===this.epoch){player.dispose();this.playing=null;this.deps.status('Replay voice')}});
  await player.play();if(token!==this.epoch){player.pause();return}this.deps.status('Speaking…');
 }catch(err){if(token!==this.epoch)return;if(this.playing){this.playing.pause();this.playing.dispose();this.playing=null}if((err as Error).name!=='AbortError'){this.deps.status('Retry voice');this.deps.unavailable()}}}
}
