"""Local-only HTTP app with transactional ironman saves and bounded AI proposals."""
from pathlib import Path
import copy,json,os,sqlite3,threading,time,urllib.request,urllib.error,uuid,mimetypes
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from urllib.parse import urlparse
import engine
ROOT=Path(__file__).resolve().parent
PORT=int(os.environ.get('BLACK_LEDGER_PORT','8791'))
DB=Path(os.environ.get('BLACK_LEDGER_DB',str(ROOT/'.runtime/campaign.sqlite3')))
AI_LOCK=threading.Lock()
MODEL=os.environ.get('BLACK_LEDGER_MODEL','qwen3:14b')
OLLAMA=os.environ.get('BLACK_LEDGER_OLLAMA','http://127.0.0.1:11435')

def connection():
    DB.parent.mkdir(parents=True,exist_ok=True)
    db=sqlite3.connect(DB,timeout=20);db.execute('PRAGMA journal_mode=WAL');return db

def initialize():
    with connection() as db:
        db.execute('CREATE TABLE IF NOT EXISTS campaign (id INTEGER PRIMARY KEY, state TEXT NOT NULL)')
        db.execute('CREATE TABLE IF NOT EXISTS receipts (id TEXT PRIMARY KEY, result TEXT NOT NULL)')
        db.execute('INSERT OR IGNORE INTO campaign VALUES (1,?)',(json.dumps(engine.new_world()),))
        w=json.loads(db.execute('SELECT state FROM campaign WHERE id=1').fetchone()[0])
        if w['director']['status']=='writing':w['director']['status']='available';w['director']['detail']='Previous preparation interrupted. Ready to retry.';save(db,w)

def load(db):return json.loads(db.execute('SELECT state FROM campaign WHERE id=1').fetchone()[0])
def save(db,w):db.execute('UPDATE campaign SET state=? WHERE id=1',(json.dumps(w),))
def state():
    with connection() as db:return engine.public(load(db))

def transact(payload):
    key=payload.get('request_id')
    if not isinstance(key,str) or not 8<=len(key)<=100:raise engine.InvalidAction('A request ID is required.')
    with connection() as db:
        db.execute('BEGIN IMMEDIATE')
        receipt=db.execute('SELECT result FROM receipts WHERE id=?',(key,)).fetchone()
        if receipt:return json.loads(receipt[0])
        w=load(db)
        if payload.get('revision')!=w['revision']:raise engine.InvalidAction('The city has changed. Refresh your view before deciding.')
        engine.action(w,payload);save(db,w);result=engine.public(w)
        db.execute('INSERT INTO receipts VALUES (?,?)',(key,json.dumps(result)))
        return result

SYSTEM='''You are the director of Black Ledger, a single-player fictional 1930s mafia drama. Propose ONE small opportunity from an established NPC. Return JSON only: title, body, speaker, operation, outcome. Speaker is one of the supplied NPC IDs. operation must be courier, mediation, or collection. body is that character's spoken offer (40–90 words), not a narration of accomplished events. outcome is one sentence describing successful completion of that limited job, not future predictions. Do not create deaths, injuries, property transfers, marriages, police arrests, new people, or debts as accomplished facts. Respect the exact provided state and history. A dead former character is not the current player. Flavor can reference existing characters, rivalries and businesses, but must not assert hidden ownership or unsupported past events. Do not include money amounts or duration in prose: the rules supply them separately. Vary motives, social tensions and dialogue. Do not repeat recent offers. Never mention JSON, AI, simulation or game mechanics in prose.'''

def request_ai():
    if not AI_LOCK.acquire(blocking=False):return False
    with connection() as db:
        db.execute('BEGIN IMMEDIATE');w=load(db)
        if not w['player']['alive'] or len(w['offers'])>=2:AI_LOCK.release();return False
        snapshot=copy.deepcopy(w);w['director']=dict(status='writing',detail='Preparing a new arrangement.',last_request=w['minute']);save(db,w)
    def work():
        try:
            context={k:snapshot[k] for k in ['life','minute','player','factions','npcs','dead']}
            context['recent_history']=snapshot['history'][-12:]
            context['places']=[dict(id=p['id'],name=p['name'],owner=snapshot['properties'][p['id']]['owner']) for p in engine.LOCATIONS]
            body=dict(model=MODEL,stream=False,think=False,format='json',messages=[dict(role='system',content=SYSTEM),dict(role='user',content=json.dumps(context))],options=dict(temperature=.8,num_predict=700))
            req=urllib.request.Request(OLLAMA+'/api/chat',json.dumps(body).encode(),{'Content-Type':'application/json'})
            with urllib.request.urlopen(req,timeout=100) as response:data=json.load(response)
            raw=json.loads(data['message']['content']);event=engine.validate_proposal(snapshot,raw)
            with connection() as db:
                db.execute('BEGIN IMMEDIATE');current=load(db)
                if current['id']!=snapshot['id'] or current['life']!=snapshot['life'] or not current['player']['alive']:return
                current['offers'].append(dict(ready=current['minute']+30,event=event))
                current['director'].update(status='ready',detail='A local AI encounter is ready. It can arrive after your next activity.');save(db,current)
        except Exception as exc:
            with connection() as db:
                db.execute('BEGIN IMMEDIATE');w=load(db);w['director'].update(status='offline',detail='Local AI unavailable or proposal rejected. Authored play remains available.');save(db,w)
            print('Director:',type(exc).__name__,str(exc)[:200],flush=True)
        finally:AI_LOCK.release()
    threading.Thread(target=work,daemon=True).start();return True

class Handler(BaseHTTPRequestHandler):
    def json(self,data,status=200):self.reply(json.dumps(data).encode(),'application/json',status)
    def reply(self,data,content='text/plain',status=200):
        self.send_response(status);self.send_header('Content-Type',content);self.send_header('Content-Length',str(len(data)));self.send_header('Cache-Control','no-store');self.send_header('X-Content-Type-Options','nosniff');self.end_headers()
        try:self.wfile.write(data)
        except (BrokenPipeError,ConnectionResetError):pass
    def do_GET(self):
        if self.path=='/api/state':return self.json(state())
        if self.path=='/api/health':return self.json(dict(ok=True,game='Black Ledger'))
        path=ROOT/'web'/('index.html' if self.path=='/' else self.path.split('?')[0].lstrip('/'))
        if not path.resolve().is_relative_to((ROOT/'web').resolve()) or not path.is_file():return self.json({'error':'Not found'},404)
        self.reply(path.read_bytes(),mimetypes.guess_type(str(path))[0] or 'application/octet-stream')
    def do_POST(self):
        origin=self.headers.get('Origin','')
        if origin and origin not in [f'http://127.0.0.1:{PORT}',f'http://localhost:{PORT}']:return self.json({'error':'Origin rejected'},403)
        if not self.headers.get('Content-Type','').startswith('application/json'):return self.json({'error':'JSON required'},415)
        try:
            length=int(self.headers.get('Content-Length','0'))
            if length<2 or length>12000:raise ValueError('Invalid body length')
            payload=json.loads(self.rfile.read(length))
            if not isinstance(payload,dict):raise ValueError('Object required')
            if self.path=='/api/action':return self.json(transact(payload))
            if self.path=='/api/director':return self.json(dict(started=request_ai()))
            if self.path=='/api/speech':
                w=state();e=w['event']
                # Only synthesize the current saved line, not arbitrary text from a web caller.
                if not e or payload.get('event')!=e['id']:return self.json({'error':'Conversation ended'},409)
                person=next((n for n in w['npcs'] if n['id']==e['speaker']),w['npcs'][0])
                voice=dict(text=e['body'],speaker=person['name'],voice='warm')
                endpoint=os.environ.get('AFTERLIGHT_DIRECTOR_URL','http://127.0.0.1:8787')+'/speech'
                req=urllib.request.Request(endpoint,json.dumps(voice).encode(),{'Content-Type':'application/json'})
                with urllib.request.urlopen(req,timeout=24) as r:audio=r.read(8000001)
                if len(audio)>8000000 or not audio.startswith(b'RIFF'):raise ValueError('Invalid audio')
                return self.reply(audio,'audio/wav')
            self.json({'error':'Not found'},404)
        except engine.InvalidAction as exc:self.json({'error':str(exc)},409)
        except (ValueError,TypeError,KeyError) as exc:self.json({'error':str(exc)},400)
        except (urllib.error.URLError,TimeoutError) as exc:self.json({'error':'Voice service unavailable. You can continue reading.'},503)
        except Exception as exc:
            print('Request error',type(exc).__name__,str(exc),flush=True);self.json({'error':'Unable to complete this request; your previous save is intact.'},500)
    def log_message(self,*args):pass

if __name__=='__main__':
    initialize();print(f'Black Ledger · http://127.0.0.1:{PORT} · {DB}',flush=True)
    ThreadingHTTPServer(('127.0.0.1',PORT),Handler).serve_forever()
