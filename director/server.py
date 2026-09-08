"""Afterlight local director. Standard library only; never executes model output."""
import copy
import threading
from contextlib import contextmanager
import hashlib
from pathlib import Path
import subprocess
try:
    from .speech import render as render_speech
except ImportError:
    try:
        from director.speech import render as render_speech
    except ImportError:
        from speech import render as render_speech
import json
import os
import time
import urllib.request
import urllib.error
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

try:
    from director import dialogue, social, dialogue_brief, dialogue_review, planning, summit_dialogue
except ImportError:
    import dialogue, social, dialogue_brief, dialogue_review, planning, summit_dialogue

def source_fingerprint():
    digest=hashlib.sha256()
    for path in sorted(Path(__file__).resolve().parent.glob('*.py')):
        digest.update(path.name.encode());digest.update(b'\0');digest.update(path.read_bytes());digest.update(b'\0')
    return digest.hexdigest()[:16]

# Captured at startup: subsequent edits must not masquerade as loaded code.
SOURCE_FINGERPRINT=source_fingerprint()

MODEL_LOCK = threading.Lock()

@contextmanager
def model_slot():
    if not MODEL_LOCK.acquire(blocking=False):
        raise BlockingIOError('Local model busy; retry shortly')
    try:
        yield
    finally:
        MODEL_LOCK.release()

MODEL = os.environ.get('AFTERLIGHT_MODEL', 'qwen3:14b')
OLLAMA_URL = os.environ.get('AFTERLIGHT_OLLAMA_URL', 'http://127.0.0.1:11434')
PORT = int(os.environ.get('AFTERLIGHT_DIRECTOR_PORT', '8787'))
# Unity prepares chapters asynchronously, including slower cold starts. Legacy
# launchers retain their shorter request budget unless explicitly configured.
PLAN_TIMEOUT = max(10.0, min(180.0, float(os.environ.get('AFTERLIGHT_PLAN_TIMEOUT', '48'))))
CUES = ['story','arrival','politics','danger','loss','discovery','bond','radio']
BASIC_KINDS = ['recruit', 'trade', 'aid', 'envoy', 'cache']
KINDS = BASIC_KINDS + ['mediation']
FACTIONS = ['Orchard Union', 'Roadkeepers', 'Iron Company']
SCHEMA = {'type': 'object', 'properties': {
    'loadout': {'type':'string','enum':['scavengers','patrol','breachers','hunters']},
    'arrival_hour': {'type':'integer','minimum':0,'maximum':23},
    'political_hour': {'type':'integer','minimum':0,'maximum':23},
    'raid_hour': {'type':'integer','minimum':0,'maximum':23},
    'raid': {'type':'boolean'},
    'cue': {'type':'string', 'enum':CUES},
    'kind': {'type':'string', 'enum':KINDS},
    'other': {'type':'string'},
    'name': {'type':'string'},
    'faction': {'type':'string', 'enum':FACTIONS},
    'title': {'type':'string'},
    'body': {'type':'string'},
    'rumor': {'type':'string'},
    'lane': {'type':'string', 'enum':['north','east','south']},
    'threat': {'type':'string','enum':['undead','raiders']},
    'pressure': {'type':'integer', 'minimum':-1, 'maximum':1},
}, 'required':['kind','name','faction','title','body','rumor','lane','threat','pressure'], 'additionalProperties':False}
SYSTEM = '''You are the director of Afterlight, a grounded zombie-survival settlement game.
Return one JSON plan for the requested next day. Choose loadout scavengers (mostly pistols, little armor), patrol (carbines and vests), breachers (fewer armored shotgun attackers), or hunters (fewer unarmored hunting-rifle shooters with a shotgun escort). Hunters favor range and precision but carry small magazines; they are not automatically allies, a new faction or proof of past attacks. Vary equipment to fit the encounter; avoid constant heavy assaults. Choose arrival_hour, political_hour and raid_hour independently on a 24-hour clock (0–23). The political proposal occurs at political_hour even if the gate is occupied or combat is happening; its validity is checked at execution. Do not write visitor speech as if a later political proposal has already happened. Set raid true or false. Quiet days are welcome. Attacks and visitors can occur in daylight or darkness; they are not forced night waves. Threat/lane describe the planned attack if raid is true. Use only recorded facts. Dead people stay dead.
arrival_presentation describes the rendered visitor: use that visible equipment, transport and meeting behavior. Do not invent visible firearms or wagons, or assume a built gate. Recruit scenes concern one individual applying to join, never their whole faction.
encounter_options lists the supported dilemmas available for this planning turn. Choose its kind and write a present-day scene fitting that transaction. Evidence explains colony pressure to you, not automatic visitor knowledge. No transaction has occurred yet.
Choose a NEW visitor name, distinct from every current survivor. Never turn a resident into an outside visitor. Create a visitor with a plausible motive, and a specific dilemma with an unresolved thread or a callback to a previous choice.
No scripted player outcomes, impossible knowledge, deaths, inventory changes, or events that already happened unless in history.
Visitor kinds have EXACT supported choices; your prose must fit them:
recruit: visitor asks to join, admitting costs 3 food; can refuse. Capacity 6.
trade: caravan sells 16 ammunition for 6 scrap; can sell 4 food for 7 scrap; can decline.
aid: visitor requests 2 medicine; helping grants 6 scrap and faction goodwill; can refuse.
envoy: faction requests 5 ammunition, promises 8 food and goodwill; can refuse.
cache: visitor reveals salvage; pay 3 food to gain 10 scrap and 5 ammunition, or decline.
mediation: choose faction and other from an available politics.summit_options pair, in either order. The visitor asks you to host talks for 6 food (+8 goodwill with their faction), arm their side for 8 rounds (+12 goodwill with their faction, -15 with the rival and -15 relations between them), or remain neutral. Talks arrive two days after acceptance; peace is not guaranteed. Set other to the rival faction name; otherwise other is empty. Never invent a war cause.
Never promise additional effects. Reference one actual recent decision if appropriate, but do not repeat every history entry.
Body max 42 words, title max 5 words, name max 3 words. Body is polished in-world prose only. Never print lane, pressure, JSON fields, game instructions, or a list of choices in body. Never invent past actions by named survivors. Write only the visitor’s present situation. Choose threat undead or raiders. Avoid making every day violent. Raiders are unaffiliated armed scavengers, not a war with a named faction. Use undead before day 3. After day 2, mix occasional human raids with undead attacks. Rumor max 22 words: a clear, concrete warning of the chosen threat and approach lane. Do not assign an unsupported political motive. Shooting, building, medicine and melee levels describe earned survivor skills; do not invent extra perks. away=true residents are currently off-map and cannot participate in a present settlement scene. Their expedition remaining value is simulation seconds, not world-clock minutes or a guaranteed safe return. colony.training describes an unfinished lesson: it grants no XP until completed, and reserved materials are already deducted from supplies. Progress needs 60 seconds together; waiting is not productive training. Do not announce the student has learned the skill or spent additional supplies before completion. Do not turn medical expertise into unimplemented diagnoses, surgery or disease cures.
Choose lane north/east/south. Pressure -1,0,1 bounded variety, lower if survivors hurt; do not punish success automatically.
Example of appropriate prose for an aid visitor (do not copy): "A nurse waits beside a handcart. Two travelers were hurt on the road, and her last bandage is already soaked through. The farming communities can offer salvage if you spare medicine."
Choose a cue for the visitor notification: arrival, story, politics, danger, loss, discovery, bond, or radio. Match the situation; reserve danger for a real supported warning. Avoid constant betrayals. Quiet help and ordinary trade belong in this world. No markdown.'''

PROPOSAL_SCHEMA = {'type':'object','properties':{
    'action':{'type':'string','enum':['none','dispute','splinter','plot','summit']},
    'actor':{'type':'string'},'target':{'type':'string'},'name':{'type':'string'},
    'leader':{'type':'string'},'plot_kind':{'type':'string','enum':['sabotage','assassination']}
},'required':['action','actor','target','name','leader','plot_kind'],'additionalProperties':False}
SCHEMA['properties']['proposal']=PROPOSAL_SCHEMA
SCHEMA['required'] += ['other','proposal','cue','arrival_hour','political_hour','raid_hour','raid']
SYSTEM = SYSTEM.replace('Raiders are unaffiliated armed scavengers, not a war with a named faction.',
    'Raiders may belong to a known faction only if its goodwill is -30 or lower; otherwise they are unaffiliated.')
SYSTEM = SYSTEM.replace('Do not assign an unsupported political motive.',
    'Political motives must be supported by recorded faction goodwill, relations or pending plots.')
SYSTEM += """
reserved_cast_names are established people. Choose a NEW name for this arrival; scheduled returning characters are handled by the game.
The politics object is the persistent regional simulation. Names in politics.factions are the ONLY existing factions.
You may propose ONE bounded political development. Prefer a relevant development when an unresolved feud or instability supports it; otherwise action none is appropriate. Proposals are checked again when applied.
- dispute: actor and target must be two different existing factions; relations worsen by 8. Do not claim war has begun.
- summit: choose actor and target from politics.summit_options. This schedules delegations to arrive two days after the proposal executes. The player later chooses guarded talks, open talks or cancellation. Do not claim peace, betrayal or casualties in advance.
- splinter: actor is an existing parent with stability <=40, target empty. Name and leader describe a new independent group. Only allowed if fewer than 8 factions and target_day-last_birth >=4. Explain a plausible political disagreement in the visitor body without claiming it already happened.
- plot: actor must have player goodwill <=-25. plot_kind is sabotage or assassination. Sabotage creates a warning; assassination is a hidden plan that later sends a visitor seeking admission. Never reveal a proposed assassination, its target or its cover in the visitor story or speech. Never kill anyone or claim success. Use only when no plot by that actor is pending.
Use empty strings for unused proposal fields and plot_kind sabotage when action none. Do not propose novelty for its own sake. Political hostility must come from recorded state.
politics.summit_history records completed talks, cancellations or a betrayal that started combat. Leader names are historical identities; do not invent deaths, victory or peace. Use relation_before/after and goodwill_before/after as the actual changes. A talks_completed outcome may still leave factions at war.
politics.resolved_outcomes are completed historical events, not pending threats. Use their exact recorded losses or recovered supplies; never replay them as new attacks. leader_at_warning is a historical identity, not necessarily the current leader. colony.settlement.damaged counts structures that currently need repair; do not claim an old damaged structure remains damaged if current state does not support it.
colony.recent_encounters lists completed basic visits in resolution order, including refusals; it excludes queued visitors and unaccepted pending offers. Vary new encounters using this history: after two consecutive visits of one kind, prefer a different kind with a useful affordable choice. Preserve justified returning-character stories; do not invent earlier meetings, infer resource transfers from decision labels, or assume a newcomer personally knows these visits.
unaffordable_paid_encounters_now lists basic encounters whose paid choices cannot currently be funded. Declining remains free, and resources may change before arrival. Prefer an encounter with a useful affordable choice when the colony is struggling; unaffordable offers are occasional preparation dilemmas, not completed promises or guaranteed aid. Do not invent discounts or alternate payments.
colony.recent_trades are completed supply-exchange transactions from the settlement perspective: bought means items received and credits paid; sold means items given and credits received. quantity and credits are transaction totals, already applied. They are not promises, current inventory, or future deliveries. No counterparty or faction is recorded: do not attribute these transactions to a named visitor, faction, smuggler, or political deal. Do not infer player motives, replay trades, or grant rewards again; use current supplies for current needs.
colony.recent_expeditions are completed historical journeys, not current absences. received supplies have already been added; do not grant them again or infer they remain unspent. returned_wounded is a past injury, not proof of current health. fatal_injury means the traveler delivered supplies but died; traveler_lost means no returning traveler or delivered supplies and does not establish a perpetrator. Do not resurrect lost travelers, invent attackers or political motives, or replay these outcomes. Use survivor state for current health and availability.
colony.recent_attacks are completed encounters with historical dates. Named losses and net_health_loss compare residents at home at the start with the end; they do not prove each injury's perpetrator or an assassination motive. structures_lost_or_spent can include consumed traps. structures_damaged records surviving structures and their net health lost during that raid, grouped by kind; they were not destroyed. structural_damage_recorded=false means this detail was not recorded. Later repairs do not change historical damage: use current settlement damage to decide whether repairs are still needed. weapons_dropped are historical battlefield drops, not recovered or owned equipment. Current inventory and health take precedence over these old observations. details_recorded=false means details are unknown, not zero casualties. Never replay a completed attack or reward; a newcomer does not automatically know private losses. You may use recorded losses to motivate an appropriate offered aid or trading dilemma without inventing new mechanics or completed aid.
Colony structures, equipment, memories and unresolved story threads are context, not permission to invent actions or new mechanics. Refer to old recorded decisions when relevant. Ordinary trade and human warmth should remain common.
"""

SCHEMA['properties']['dialogue']=dialogue.SCHEMA
SCHEMA['required'].append('dialogue')
SYSTEM+=dialogue.INSTRUCTIONS+social.INSTRUCTIONS
SYSTEM+="\nConversation records and reported_speech document questions and statements, not proof that their contents happened. A question about a murder does not establish a murder. NPC speech can be mistaken. Only simulation event records establish outcomes; never use an unverified spoken claim alone as grounds for a death, feud, betrayal or reward.\n"
SYSTEM+="\nrecalled_events are exact OLD journal facts with dates and generations. They may inspire a relevant callback, debt or changed relationship, but are not new events or permission to repeat rewards. Preserve their original time and participants. Do not make an unrelated newcomer claim personal knowledge of private events. Present-day mechanics and identities take precedence over old circumstances.\n"
SYSTEM+="\ncolony.settlement.completed counts finished placed structures, including unavailable ones. unavailable is the union of destroyed and marked_for_demolition structures; the two reasons may overlap, so never add them together. Available placed counts are completed minus unavailable; this does not prove reachability or shelter. destroyed means zero health, not a repairable damaged facility; marked_for_demolition means designated for removal, not already removed. damaged counts surviving structures needing repair that are not designated for removal. under_construction is unfinished and grants no facility benefits. Legacy buildings are separate older facilities, not a substitute for placed structures. colony.shelter is the simulation's actual enclosed, roofed bed capacity; visible beds alone do not prove shelter. colony.land.remaining is finite gatherable yield, not owned inventory. arms_market stock is individual weapons and field equipment offered by a trader, not owned gear. vest is a protective vest, machete is a melee weapon, and kit is a single-use field medical kit; buying places these in spare storage, not on a resident. Its dated prices are purchase quotes; resale is half the quote rounded down, using the same trader wallet in supply_market. stored_equipment contains spare owned items, excluding equipped weapons on survivors. Do not invent weapon purchases, ownership, or availability on a future day from a current shipment. supply_market stock is lots offered by a trader, not colony supplies; its dated disruption and prices describe that shipment only. Resident requests are offered or promised obligations until their status records fulfillment. You may acknowledge these facts, but do not claim construction, harvesting, purchases or promises have completed unless recorded.\n"
SCHEMA['properties']['social_proposal']=social.SCHEMA
SCHEMA['required'] += ['social_proposal','loadout']

def reserved_cast_names(state):
    people=state.get('survivors',[])+state.get('dead',[])+state.get('departed',[])+state.get('colony',{}).get('threads',[])
    names=[p.get('name','') for p in people if isinstance(p,dict)]
    names+=state.get('reserved_cast_names',[])
    names+=[f.get('leader','') for f in state.get('politics',{}).get('factions',[]) if isinstance(f,dict)]
    return {' '.join(n.split()).casefold() for n in names if isinstance(n,str) and n.strip()}

def faction_names(state):
    entries=state.get('politics',{}).get('factions',[]) if isinstance(state.get('politics',{}),dict) else []
    names=[f['name'] for f in entries if isinstance(f,dict) and isinstance(f.get('name'),str)]
    legacy=state.get('faction_relations',{})
    return names or ([name for name in legacy if isinstance(name,str) and name.strip()] if isinstance(legacy,dict) else []) or FACTIONS

def validate(raw, state=None):
    if not isinstance(raw, dict):
        raise ValueError('plan must be an object')
    out = {}
    loadout=raw.get('loadout','scavengers')
    out['loadout']=loadout if isinstance(loadout,str) and loadout in ['scavengers','patrol','breachers','hunters'] else 'scavengers'
    for key, default in [('arrival_hour',9),('political_hour',12),('raid_hour',21)]:
        value=raw.get(key,default)
        if type(value) is not int or not 0<=value<=23: raise ValueError('invalid schedule hour')
        out[key]=value
    if type(raw.get('raid',True)) is not bool: raise ValueError('invalid raid switch')
    out['raid']=raw.get('raid',True)
    cue=raw.get('cue','arrival')
    out['cue']=cue if isinstance(cue,str) and cue in CUES else 'arrival'
    for key, limit in [('name',48),('title',72),('body',550),('rumor',180)]:
        value = raw.get(key)
        if not isinstance(value,str) or not value.strip():
            raise ValueError('missing text: '+key)
        out[key] = ' '.join(value.split())[:limit]
    for key, allowed in [('kind',KINDS),('faction',faction_names(state or {})),('lane',['north','east','south']),('threat',['undead','raiders'])]:
        if raw.get(key) not in allowed:
            raise ValueError('invalid '+key)
        out[key] = raw[key]
    if type(raw.get('pressure')) is not int or not -1 <= raw['pressure'] <= 1:
        raise ValueError('invalid pressure')
    out['pressure'] = raw['pressure']
    out['other'] = raw.get('other', '') if out['kind']=='mediation' else ''
    if not isinstance(out['other'],str): raise ValueError('invalid mediation rival')
    if out['kind']=='mediation' and not any({out['faction'],out['other']}=={p['actor'],p['target']} for p in planning.summit_options(state or {})):
        raise ValueError('mediation requires an available hostile faction pair')
    proposal=raw.get('proposal',{'action':'none'})
    if not isinstance(proposal,dict) or proposal.get('action') not in ['none','dispute','splinter','plot','summit']:
        raise ValueError('invalid proposal')
    out['proposal']={'action':proposal['action']}
    for key in ['actor','target','name','leader','plot_kind']:
        value=proposal.get(key,'')
        if not isinstance(value,str): raise ValueError('invalid proposal text')
        out['proposal'][key]=' '.join(value.split())[:32 if key in ['name','leader'] else 48]
    out['social_proposal']=social.validate(raw.get('social_proposal',{}),state or {})
    if raw.get('dialogue'):
        try: out['dialogue']=dialogue.validate(raw['dialogue'],out['kind'])
        except ValueError: out['dialogue_error']='Invalid conversation omitted; standard choices retained'
    return out

def validate_context(plan, state):
    if plan['kind'] not in [option['kind'] for option in planning.encounter_options(state)]:
        raise ValueError('visitor does not match active survival dilemmas')
    resident_names = {str(p.get('name','')).split()[0].casefold() for p in state.get('survivors',[])+state.get('dead',[])+state.get('departed',[]) if isinstance(p,dict) and p.get('name')}
    if plan['name'].split()[0].casefold() in resident_names:
        raise ValueError('visitor is already a resident')
    if ' '.join(plan['name'].split()).casefold() in reserved_cast_names(state):
        raise ValueError('new visitor reuses an established cast name')
    if state.get('target_day',2)<3 and plan['threat']=='raiders':
        raise ValueError('raiders unavailable before day 3')
    proposal=plan.get('proposal',{'action':'none'})
    action=proposal['action']; actor=proposal.get('actor',''); target=proposal.get('target','')
    world=state.get('politics',{})
    known={f['name']:f for f in world.get('factions',[]) if isinstance(f,dict) and 'name' in f}
    if action!='none' and actor not in known: raise ValueError('proposal actor does not exist')
    if action=='dispute' and (target not in known or actor==target): raise ValueError('invalid dispute target')
    if action=='summit' and not any({option['actor'],option['target']}=={actor,target} for option in planning.summit_options(state)):
        raise ValueError('summit requires an existing war and an available meeting slot')
    if action=='splinter':
        if known[actor].get('stability',100)>40 or len(known)>=8 or state.get('target_day',2)-world.get('last_birth',0)<4:
            raise ValueError('splinter has no political foundation')
        if len(proposal.get('name',''))<3 or proposal['name'] in known or not proposal.get('leader'):
            raise ValueError('invalid new faction identity')
    if action=='plot':
        if state.get('faction_relations',{}).get(actor,0)>-25 or proposal.get('plot_kind') not in ['sabotage','assassination']:
            raise ValueError('plot has no hostile motive')
        if any(p.get('faction')==actor for p in world.get('plots',[])): raise ValueError('plot already exists')
    plan['rumor']=('Armed raiders' if plan['threat']=='raiders' else 'Undead')+' approach from the '+plan['lane']+'. Scouts expect movement around %02d:00.'%plan.get('raid_hour',21)
    if not plan.get('raid',True): plan['rumor']='Scouts report no organized attack planned for this day.'
    return plan

def fallback(state, reason='Model unavailable'):
    day = int(state.get('target_day',2))
    kind = BASIC_KINDS[(day-2)%len(BASIC_KINDS)]
    world=str(state.get('world_id',''))
    seed=int(hashlib.sha256((world+'|'+str(day)).encode()).hexdigest()[:12],16)
    if world: kind=BASIC_KINDS[(day+int(hashlib.sha256(world.encode()).hexdigest()[:8],16))%len(BASIC_KINDS)]
    options=[option for option in planning.encounter_options(state) if option['kind'] in BASIC_KINDS]
    if kind not in [option['kind'] for option in options]: kind=options[seed%len(options)]['kind']
    names=['Rowan','Tess','Edda','Finch','Wren','Bram','Soren','Maren','Cal','Jory','Ansel','Petra','Lena','Hale','Nico','Dara','Orin','Vera','Kit','Leif','Rhea','Seth','Mae','Alden','Tilda','Perrin','Esme','Bren','Luca','Greta','Ivo','Nessa','Ari','Dorian','Mira','Ewan','Sasha','Laurel','Enid','Reed','Cora','Asa','Faye','Ellis','Blythe','Kellan','Della','Morgan','Rory','Galen','Lark','Davin','Nolan','Astrid','Tarin','Jessa','Gus','Thea','Vaughn','Bea','Lyle','Zara','Hollis','Quinn']
    if world:
        offset=seed%len(names);names=names[offset:]+names[:offset]
    residents={name.split()[0] for name in reserved_cast_names(state)}
    name=next((n for n in names if n.casefold() not in residents),'Wayfarer-'+format(seed,'x'))
    suffix=2
    while name.split()[0].casefold() in residents:
        name='Wayfarer-'+format(seed,'x')+'-'+str(suffix);suffix+=1
    faction={'recruit':'Orchard Union','trade':'Roadkeepers','aid':'Orchard Union','envoy':'Iron Company','cache':'Roadkeepers'}[kind]
    if faction not in faction_names(state): faction=faction_names(state)[seed%len(faction_names(state))]
    bodies = {
        'recruit': name+' arrives with a patched coat and steady hands. They have traveled from '+faction+' and ask for a bed, a share of your food, and work at the crossing.',
        'trade': name+' of '+faction+' stops a cart at the gate, with ammunition to sell and empty food crates to fill.',
        'aid': name+' of '+faction+' asks for medicine for an injured traveling companion. They offer salvage in return for your help.',
        'envoy': name+' brings an offer from '+faction+': ammunition from your stores in exchange for a food shipment.',
        'cache': name+' of '+faction+' has located supplies in an abandoned service station. Rations will pay for the supplies already found.'}
    return dict(arrival_hour=8+day%10,raid_hour=12+day%11,raid=day%4!=0,proposal={'action':'none'},kind=kind,name=name,faction=faction,title={'recruit':'Another place at the table','trade':'The road still provides','aid':'Someone else’s wounded','envoy':'The price of a crossing','cache':'Beneath the old service station'}[kind],body=bodies[kind],rumor=('Scouts report no organized attack planned for this day.' if day%4==0 else ('Armed raiders' if day%3==0 else 'Undead')+' approach from the east around %02d:00.'%(12+day%11)),lane='east',threat='raiders' if day%3==0 else 'undead',pressure=0,source='Prepared fallback',detail=reason)

def generate_dialogue(plan,state,diagnostics=None,review_arrival=False):
    deadline=time.monotonic()+24
    stage='draft'
    arrival_approved=False
    def report(source,reason):
        if diagnostics is not None:
            diagnostics.update(source=source,stage=stage,reason=reason)
        chosen=graph if source!='fallback' else fallback
        if review_arrival:
            plan['arrival_text_source']='reviewed' if arrival_approved else 'spoken_fallback'
            if not arrival_approved: plan['body']=plan.get('name','The visitor')+' says, “'+chosen[0]['line']+'”'
        return chosen
    system=dialogue_brief.INSTRUCTIONS
    context=dialogue.public_context(plan,state)
    fallback=dialogue_brief.factual_fallback(plan.get('kind'),context)
    payload={'model':MODEL,'stream':False,'think':False,'format':dialogue_brief.SCHEMA,
             'messages':[{'role':'system','content':system},{'role':'user','content':json.dumps(context)}],
             'options':{'temperature':0.55,'num_ctx':8192,'num_predict':900},'keep_alive':'15m'}
    req=urllib.request.Request(OLLAMA_URL+'/api/chat',data=json.dumps(payload).encode(),headers={'Content-Type':'application/json'})
    try:
        with model_slot(), urllib.request.urlopen(req,timeout=16) as response: result=json.load(response)
        graph=dialogue_brief.compile_brief(json.loads(result['message']['content']),plan['kind'],context)
        if dialogue_review.premature_acceptance(graph): return report('fallback','premature_acceptance')
        remaining=deadline-time.monotonic()
        if remaining<=0: return report('fallback','deadline_before_review')
        stage='review'
        review=dialogue_review.payload(graph,context,MODEL,plan.get("body","") if review_arrival else None)
        req=urllib.request.Request(OLLAMA_URL+'/api/chat',data=json.dumps(review).encode(),headers={'Content-Type':'application/json'})
        with model_slot(), urllib.request.urlopen(req,timeout=remaining) as response: result=json.load(response)
        verdict=json.loads(result['message']['content'])
        if time.monotonic()>deadline: return report('fallback','review_deadline')
        if review_arrival: verdict,arrival_approved=dialogue_review.split_arrival_review(verdict)
        if not dialogue_review.approved(verdict):
            print(json.dumps({'dialogue_review':'not approved','visitor':plan.get('name',''),'issues':str(verdict.get('checks',[]))[:1600] if isinstance(verdict,dict) else 'Malformed verdict'}),flush=True)
            retained=dialogue_review.retain_reviewed(graph,fallback,verdict)
            if retained:
                graph=dialogue.validate(retained,dialogue_brief.contract_kind(plan['kind'],context))
                return report('reviewed_partial','unsupported_passages_removed')
            return report('fallback','review_rejected')
        return report('generated','approved')
    except (urllib.error.URLError,TimeoutError,ValueError,KeyError,OSError) as exc:
        return report('fallback',type(exc).__name__)

def generate_summit(visitor,state):
    context=summit_dialogue.context(visitor,state)
    payload={'model':MODEL,'stream':False,'think':True,'format':summit_dialogue.BRIEF_SCHEMA,
             'messages':[{'role':'system','content':summit_dialogue.BRIEF_INSTRUCTIONS},
                         {'role':'user','content':json.dumps(context)}],
             'options':{'temperature':0.35,'num_ctx':8192,'num_predict':6144},'keep_alive':'15m'}
    deadline=time.monotonic()+85
    try:
        for attempt in range(2):
            request=urllib.request.Request(OLLAMA_URL+'/api/chat',data=json.dumps(payload).encode(),headers={'Content-Type':'application/json'})
            with model_slot(), urllib.request.urlopen(request,timeout=max(1,deadline-time.monotonic())) as response: result=json.load(response)
            content=result['message']['content']
            try:
                graph,prepared_lines=summit_dialogue.ground_prose(summit_dialogue.compile_brief(json.loads(content),context))
                review=summit_dialogue.review_payload(graph,context,MODEL)
                review_request=urllib.request.Request(OLLAMA_URL+'/api/chat',data=json.dumps(review).encode(),headers={'Content-Type':'application/json'})
                with model_slot(), urllib.request.urlopen(review_request,timeout=max(1,deadline-time.monotonic())) as response: reviewed=json.load(response)
                verdict=json.loads(reviewed['message']['content'])
                if not summit_dialogue.review_approved(verdict):
                    raise ValueError('Content review rejected script: '+json.dumps(verdict)+'. Correct speaker identities, remove invented incidents or treaty terms, and keep merged branches neutral.')
                return {'summit_dialogue':graph,'source':'Local AI summit'+(f' · prepared lines: {prepared_lines}' if prepared_lines else ''),'prepared_lines':prepared_lines}
            except ValueError as error:
                if attempt==1 or time.monotonic()>deadline: raise
                payload['messages'] += [{'role':'assistant','content':content},
                    {'role':'user','content':'Repair this entire summit brief. Validation error: '+str(error)+'. Return the complete corrected object with every named scene, valid moves, concise lines and labels.'}]
        raise ValueError('summit unavailable')
    except (urllib.error.URLError,TimeoutError,ValueError,KeyError,OSError) as error:
        return {'source':'Prepared summit fallback','detail':str(error)[:160]}

def arrival_dialogue(payload):
    visitor=payload.get('visitor')
    state=payload.get('state')
    if not isinstance(visitor,dict) or not isinstance(state,dict): raise ValueError('invalid conversation context')
    if visitor.get('kind')=='summit': return generate_summit(visitor,state)
    if visitor.get('kind') not in KINDS: raise ValueError('unsupported visitor')
    for field in ('name','faction','title','body'):
        if not isinstance(visitor.get(field),str) or len(visitor[field])>550: raise ValueError('invalid visitor text')
    motive=visitor.get('return_motive','')
    if not isinstance(motive,str) or motive not in ['', 'asylum','war_supplies','commerce','relief','salvage','settling']: raise ValueError('invalid return motive')
    if motive=='asylum' and visitor['kind']!='recruit' or motive=='war_supplies' and visitor['kind']!='envoy': raise ValueError('mismatched political contract')
    if visitor['faction'] not in faction_names(state): raise ValueError('inactive faction')
    if visitor['kind']=='mediation' and (visitor.get('other') not in faction_names(state) or visitor['other']==visitor['faction']): raise ValueError('invalid mediation rival')
    # Reuse the public-fact whitelist and compiled mechanical contracts.
    return {'dialogue':generate_dialogue(visitor,state)}

def generate(state):
    started = time.monotonic()
    schema=copy.deepcopy(SCHEMA)
    schema['properties'].pop('dialogue',None)
    schema['required'].remove('dialogue')
    schema['properties']['faction']['enum']=faction_names(state)
    schema['properties']['kind']['enum']=[option['kind'] for option in planning.encounter_options(state)]
    schema['properties']['proposal']['properties']['action']['enum']=planning.political_actions(state)
    if state.get('target_day',2)<3: schema['properties']['threat']['enum']=['undead']
    payload = {'model':MODEL,'stream':False,'think':False,'format':schema,'messages':[{'role':'system','content':SYSTEM.replace(dialogue.INSTRUCTIONS,'')},{'role':'user','content':json.dumps(planning.context(state),separators=(',',':'))}], 'options':{'temperature':0.8,'num_ctx':8192,'num_predict':900},'keep_alive':'15m'}
    req = urllib.request.Request(OLLAMA_URL+'/api/chat',data=json.dumps(payload).encode(),headers={'Content-Type':'application/json'})
    try:
        with model_slot(), urllib.request.urlopen(req,timeout=PLAN_TIMEOUT) as response:
            result=json.load(response)
        plan=validate_context(validate(json.loads(result['message']['content']),state),state)
        dialogue_diagnostics={}
        conversation=generate_dialogue(plan,state,dialogue_diagnostics,review_arrival=True)
        plan["dialogue_diagnostics"]=dialogue_diagnostics
        if conversation: plan['dialogue']=conversation
        else: plan['dialogue_error']='Conversation unavailable; standard choices retained'
        plan.update(source='Local AI · '+MODEL,detail='Generated in %.1fs'%(time.monotonic()-started),telemetry={'model':MODEL,'prompt_tokens':result.get('prompt_eval_count',0),'output_tokens':result.get('eval_count',0),'generation_seconds':round(time.monotonic()-started,2),'proposal':plan.get('proposal',{}).get('action','none')})
        return plan
    except (urllib.error.URLError, TimeoutError, ValueError, KeyError, OSError) as exc:
        return fallback(state,type(exc).__name__+': '+str(exc)[:160])

class Handler(BaseHTTPRequestHandler):
    def send(self, status, payload):
        body=json.dumps(payload).encode()
        self.send_response(status); self.send_header('Content-Type','application/json'); self.send_header('Content-Length',str(len(body))); self.end_headers(); self.wfile.write(body)
    def do_GET(self):
        if self.path=='/health': self.send(200,{'service':'afterlight-director','model':MODEL,'version':'0.15.5','source_fingerprint':SOURCE_FINGERPRINT})
        else: self.send(404,{'error':'not found'})
    def do_POST(self):
        if self.path not in ['/plan','/speech','/dialogue']: return self.send(404,{'error':'not found'})
        try:
            length=int(self.headers.get('Content-Length','0'))
            if length<=0 or length>64000: return self.send(413,{'error':'invalid request size'})
            state=json.loads(self.rfile.read(length))
            if not isinstance(state,dict): raise ValueError('expected object')
            if self.path=='/speech':
                try:
                    audio=render_speech(state)
                except BlockingIOError:
                    return self.send(429,{'error':'speech busy'})
                except (OSError, subprocess.SubprocessError):
                    return self.send(503,{'error':'local voice unavailable'})
                self.send_response(200)
                self.send_header('Content-Type','audio/wav')
                self.send_header('Content-Length',str(len(audio)))
                self.end_headers(); self.wfile.write(audio)
                return
            if self.path=='/dialogue': return self.send(200,arrival_dialogue(state))
            plan=generate(state)
            self.send(200,plan)
            print(json.dumps({'day':state.get('target_day'),'source':plan['source'],'detail':plan['detail']}),flush=True)
        except (ValueError,TypeError): self.send(400,{'error':'invalid JSON request'})
        except (BrokenPipeError,ConnectionResetError): pass

if __name__=='__main__':
    print(f'Afterlight director at http://127.0.0.1:{PORT}; model={MODEL}',flush=True)
    ThreadingHTTPServer(('127.0.0.1',PORT),Handler).serve_forever()
