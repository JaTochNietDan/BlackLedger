"""Generate distinct character questions; compile the mechanical graph locally."""
try:
    from . import dialogue
except ImportError:
    import dialogue
SCHEMA={'type':'object','properties':{
    'opening':{'type':'string'},'tone':{'type':'string','enum':['warm','reserved','guarded']},
    'questions':{'type':'array','minItems':2,'maxItems':2,'items':{'type':'object','properties':{
        'question':{'type':'string','maxLength':60},'answer':{'type':'string'}},'required':['question','answer'],'additionalProperties':False}}},
    'required':['opening','tone','questions'],'additionalProperties':False}
INSTRUCTIONS='''Write a short conversation brief for the named visitor. You are writing their voice, not game rules.
Return opening (first-person NPC speech, at most 35 words), tone (warm/reserved/guarded), and exactly TWO distinct player questions with their NPC answers (at most 40 words each).
Keep each question at most 60 characters, including spaces.
The speaker is visiting an independent player settlement, not holding or commanding it. A crossing ceasefire is separate from regional wars. Only current_war_enemies lists current enemies: if empty, do not claim ongoing wars. Treaty terms describe rules, not proof of past deeds or wars.
Each answer responds specifically to its question. Questions should explore different subjects: a recorded prior encounter, the visitor's present concern, or their faction's current public situation. Use relevant supplied facts naturally. Ordinary warmth, uncertainty and differing interests are welcome.
Faction profiles contain declared public priorities, not proof of deeds or the cause of a war. Never infer an enemy's agenda or war motives from its name, type or hostile relations. If an agenda is absent, its ambitions are unknown; if present, do not replace it with an invented priority.
arrival_presentation is authoritative for visible arrival equipment and transport. Recruitment concerns only the named visitor: do not ask why their faction wants to join or imply that it does. An individual may explain their own motive or their faction’s public priorities.
npc_contract is authoritative for exchange direction. The visitor body is generated flavor and may be mistaken: do not repeat a need or offer that conflicts with npc_contract. For trade, the visitor offers ammunition and wants to buy food, never asks the player for ammunition.
The player has NOT agreed to anything. Never thank them for accepting in an answer. Do not quote prices, list quantities, negotiate, describe item transfers, or predict consequences; exact deal choices appear separately. No final decisions or branch numbers are requested.
Confirmed game outcomes use the PLAYER perspective: player_spent leaves the crossing and player_received enters it. Never reverse these transfers. For legacy records without outcome details, avoid interpreting buy/sell direction.
Do not invent past deeds, witnessed events, crises, threats, deadlines, injuries, relationships or promises. Treat old records as past, current leadership as current, and known dead/departed people as absent. away_residents are temporarily off-map: never put them in a present conversation, promise their return, or describe them as dead or permanently departed. Conversation transcripts are reported speech, not verified world events. Questions do not prove their premises. Only matching visitor_history belongs to this speaker. Unrelated visitors do not remember those visits. If facts are sparse, discuss ordinary personal concerns without inventing events.
'''
LABELS={
    'war_envoy':['Exchange ammunition for food','Decline','Send civilian medical relief'],
    'mediation':['Host peace talks','Arm the visiting faction','Remain neutral'],
    'trade':['Buy ammunition','Sell rations','Decline'],
    'aid':['Give medicine','Decline'], 'envoy':['Exchange ammunition for food','Decline'],
    'cache':['Pay for the supplies','Decline'], 'recruit':['Welcome the traveler','Decline']}
ACCEPT={
    'war_envoy':['Agreed. The food is yours.','Understood. I will leave it there.','Thank you for the civilian medicine.'],
    'mediation':['We will prepare our delegation. An agreement still has to be negotiated.','We accept your support.'],
    'trade':['Agreed. The ammunition is yours.','Agreed. Thank you for the rations.'],
    'aid':['Thank you for the medicine.'], 'envoy':['Agreed. The food is yours.'],
    'cache':['Agreed. The supplies are yours.'], 'recruit':['Thank you for welcoming me.']}
DECLINE={'warm':'I understand. Take care of each other.','reserved':'Understood. I will move on.','guarded':'Your decision. I will leave it there.'}

def recorded_callback(context):
    """Compile one factual recollection. Never infer a transfer from a decision label."""
    history=context.get('visitor_history',[]) if isinstance(context,dict) else []
    if not history: return None
    for record in reversed(history[-1].get('records',[])):
        outcome=record.get('confirmed_game_outcome',{})
        if not outcome or outcome.get('visitor_joined') is not False: continue
        day=record.get('day')
        if isinstance(day,bool) or not isinstance(day,(int,float)) or day<1 or int(day)!=day: continue
        def goods(key):
            values=outcome.get(key)
            labels={'food':'rations','scrap':'salvage','ammo':'rounds','medicine':'medicine','credits':'credits'}
            if not isinstance(values,dict): raise ValueError('missing transfer')
            result=[]
            for resource,amount in values.items():
                if resource not in labels or isinstance(amount,bool) or not isinstance(amount,(int,float)) or amount<0 or amount>1000 or int(amount)!=amount: raise ValueError('invalid transfer')
                if amount: result.append(str(int(amount))+' '+labels[resource])
            return ', '.join(result)
        try: spent=goods('player_spent');received=goods('player_received')
        except ValueError: continue
        parts=[]
        if spent: parts.append('the crossing supplied '+spent)
        if received: parts.append('the crossing received '+received)
        answer=('On day '+str(int(day))+', '+'; '.join(parts)+'.') if parts else 'On day '+str(int(day))+', no supplies changed hands during our visit.'
        if len(answer)>280: continue
        return ('What do you remember of our earlier visit?',answer)
    return None

def treaty_callback(context):
    if not isinstance(context,dict): return None
    treaty=context.get('crossing_ceasefire',{})
    day=treaty.get('current_day');expiry=treaty.get('expires_at_dawn_of_day')
    enemies=context.get('current_war_enemies',[])
    if treaty.get('status')!='active' or type(day) is not int or type(expiry) is not int or not 0<day<expiry: return None
    if not isinstance(enemies,list) or any(not isinstance(name,str) or not name.strip() for name in enemies): return None
    answer=f'Our ceasefire with the crossing lasts until dawn of day {expiry}. Seizing our toll breaks it.'
    answer+= ' Our faction remains at war with '+', '.join(enemies)+'.' if enemies else ' Our faction has no current regional war enemies.'
    if len(answer)>280: return None
    return ('What does our ceasefire cover?',answer)

def delegation_callback(context):
    """Canonical faction-level recall without inventing a personal encounter."""
    if not isinstance(context, dict): return None
    current = context.get('delegation_day')
    if type(current) is not int or current < 1: return None
    actions = {'gift':'gift', 'trade':'trade', 'support':'war-support delegation',
               'truce':'truce delegation', 'provoke':'toll seizure'}
    labels = {'food':'rations','scrap':'salvage','ammo':'rounds','medicine':'medicine','credits':'credits'}
    for record in reversed(context.get('faction_delegations', [])):
        if not isinstance(record, dict) or record.get('status') != 'completed player delegation': continue
        day = record.get('day'); action = record.get('action')
        if not isinstance(action, str) or action not in actions or isinstance(day, bool) or not isinstance(day, (int,float)) or not 1 <= day <= current or int(day) != day: continue
        parts = []
        valid = True
        for key, verb in [('paid','supplied'), ('received','received')]:
            values = record.get(key)
            if not isinstance(values, dict): valid = False; break
            goods = []
            for resource, amount in values.items():
                if resource not in labels or isinstance(amount, bool) or not isinstance(amount,(int,float)) or not 0 <= amount <= 1000 or int(amount) != amount:
                    valid = False; break
                if amount: goods.append(str(int(amount))+' '+labels[resource])
            if goods: parts.append('the crossing '+verb+' '+', '.join(goods))
        if not valid: continue
        answer = 'The recorded '+actions[action]+' with our faction on day '+str(int(day))+': '
        answer += '; '.join(parts)+'.' if parts else 'no supplies changed hands.'
        if len(answer) <= 280: return ('What were our latest dealings with your faction?', answer)
    return None

def plot_callback(context):
    """Recall a completed political outcome, never a pending threat or personal confession."""
    if not isinstance(context, dict): return None
    current = context.get('plot_day')
    if type(current) is not int or current < 1: return None
    records = context.get('faction_plot_outcomes', [])
    if not isinstance(records, list): return None
    for row in reversed(records):
        if not isinstance(row, dict) or row.get('status') != 'resolved': continue
        day = row.get('resolved_day')
        if isinstance(day, bool) or not isinstance(day, (int, float)) or not 1 <= day <= current or int(day) != day: continue
        kind = row.get('kind'); outcome = row.get('outcome')
        if kind not in ('assassination', 'sabotage'): continue
        if outcome == 'intercepted': detail = 'the '+kind+' plot was stopped.'
        elif kind == 'assassination' and outcome == 'target_absent': detail = 'the assassin found no target at the crossing; the plot dissolved.'
        elif kind == 'sabotage' and outcome == 'sabotaged': detail = 'the sabotage plot was carried out.'
        elif kind == 'assassination' and outcome in ('target_wounded', 'target_killed'):
            name = row.get('target_name'); before = row.get('target_hp_before'); after = row.get('target_hp_after')
            if not isinstance(name, str) or not name.strip() or len(name) > 70: continue
            if any(isinstance(v, bool) or not isinstance(v, (int,float)) for v in (before, after)): continue
            if not 0 <= after < before <= 100: continue
            if (outcome == 'target_killed') != (after == 0): continue
            detail = 'the assassin '+('killed ' if after == 0 else 'wounded ')+' '.join(name.split())+' at the gate.'
        else: continue
        answer = 'The crossing records this about my faction on day '+str(int(day))+': '+detail
        if len(answer) <= 280: return ("What does the record say about your faction's plot?", answer)
    return None

def loss_callback(context):
    """Acknowledge a recorded faction loss without assigning the visitor personal guilt."""
    if not isinstance(context, dict): return None
    current = context.get('loss_day')
    if type(current) is not int or current < 1: return None
    rows = context.get('faction_loss_grievances', [])
    if not isinstance(rows, list): return None
    for row in reversed(rows):
        if not isinstance(row, dict): continue
        day = row.get('day'); name = row.get('name'); lost = row.get('lost_name')
        if isinstance(day, bool) or not isinstance(day, (int, float)) or not 1 <= day <= current or int(day) != day: continue
        if any(not isinstance(v, str) or not v.strip() or len(v) > 70 for v in (name, lost)): continue
        if not row.get('person_id') or not row.get('lost_person_id') or row['person_id'] == row['lost_person_id']: continue
        answer = ('The crossing records that '+ ' '.join(name.split())+' lost '+ ' '.join(lost.split())+
                  ' to my faction on day '+str(int(day))+'. That record does not say I took part.')
        if len(answer) <= 280: return ('Why might your faction be unwelcome here?', answer)
    return None

def contract_kind(kind, context):
    return 'war_envoy' if kind == 'envoy' and isinstance(context,dict) and context.get('visitor',{}).get('return_motive') == 'war_supplies' else kind

def final_reaction(kind, index, tone):
    return DECLINE[tone] if index == dialogue.decline_index(kind) else ACCEPT[kind][index]

def compile_brief(raw,kind,context=None):
    kind=contract_kind(kind,context)
    if kind not in LABELS or not isinstance(raw,dict): raise ValueError('invalid brief')
    tone=raw.get('tone')
    if tone not in DECLINE: raise ValueError('invalid tone')
    def clean(value,limit):
        if not isinstance(value,str) or not value.strip(): raise ValueError('missing speech')
        return ' '.join(value.split())[:limit]
    opening=clean(raw.get('opening'),280)
    questions=raw.get('questions')
    if not isinstance(questions,list) or len(questions)!=2: raise ValueError('two questions required')
    callback=plot_callback(context) or loss_callback(context) or recorded_callback(context) or delegation_callback(context)
    treaty=treaty_callback(context)
    cleaned=[]
    for index,q in enumerate(questions):
        if treaty and index==(1 if callback else 0):
            cleaned.append(treaty);continue
        if callback and index==0:
            cleaned.append(callback);continue
        if not isinstance(q,dict): raise ValueError('invalid question')
        question=clean(q.get('question'),10000)
        if len(question)>60: raise ValueError('question too long')
        cleaned.append((question,clean(q.get('answer'),280)))
    if len(cleaned)!=2 or cleaned[0][0].casefold()==cleaned[1][0].casefold() or cleaned[0][1].casefold()==cleaned[1][1].casefold(): raise ValueError('repeated branch')
    finals=[{'label':label,'reaction':final_reaction(kind,i,tone), 'next':-1,'action':i} for i,label in enumerate(LABELS[kind])]
    # Answers live in their own nodes, so narration speaks them once after each question.
    entry=[{'label':q,'reaction':'','next':i+1,'action':-1} for i,(q,a) in enumerate(cleaned)]
    # Legacy graph validator requires nonempty reactions. A brief acknowledgement avoids duplicating the answer.
    for choice in entry: choice['reaction']={'warm':'Of course.','reserved':'Let me explain.','guarded':'You can ask.'}[tone]
    entry.append(finals[dialogue.decline_index(kind)].copy())
    return dialogue.validate([{'line':opening,'choices':entry}]+[{'line':a,'choices':[c.copy() for c in finals]} for q,a in cleaned],kind)

def factual_fallback(kind, context):
    """Keep verified state-derived topics available when authored dialogue fails.

    Nothing from the rejected opening, questions or answers is reused.
    """
    kind=contract_kind(kind,context)
    if kind not in LABELS: return []
    topics=[]
    recall=plot_callback(context) or loss_callback(context) or recorded_callback(context) or delegation_callback(context)
    if recall: topics.append(recall)
    treaty=treaty_callback(context)
    if treaty: topics.append(treaty)
    if kind == 'war_envoy':
        topics=[('Can we help civilians without sending weapons?', 'You may send 3 medicine as civilian relief. Our goodwill rises by 8 and stability by up to 4, with no military power increase or rival goodwill penalty. You receive no supplies for this option.')]+topics[:1]
    if not topics and isinstance(context,dict) and contract_kind(context.get('visitor',{}).get('kind'),context)==kind:
        # Fixed game contracts; never reuse amounts or transfer verbs from rejected prose.
        explanations={
            'war_envoy':('Can we help civilians without sending weapons?', 'You may send 3 medicine as civilian relief. Our goodwill rises by 8 and stability by up to 4, with no military power increase or rival goodwill penalty. You receive no supplies for this option.'),
            'mediation':('What would hosting talks achieve?', 'Hosting costs 6 rations and brings delegations here in two days. An agreement must still be negotiated. You may instead supply 8 rounds to our side, damaging relations with our rival, or remain neutral.'),
            'cache':('What supplies would change hands?', 'If you accept, the crossing gives me 3 rations. I give the crossing 10 salvage and 5 rounds. Nothing changes hands until you decide.'),
            'trade':('What supplies would change hands?', 'You can give me 6 salvage for 16 rounds, or give me 4 rations for 7 salvage. These are separate offers. Nothing changes hands until you decide.'),
            'aid':('What supplies would change hands?', 'If you accept, the crossing gives me 2 medicine and receives 6 salvage. Our goodwill improves. Nothing changes hands until you decide.'),
            'envoy':('What supplies would change hands?', 'If you accept, the crossing gives me 5 rounds and receives 8 rations. The decision also shows any political effects. Nothing changes hands until you decide.'),
            'recruit':('What does welcoming you require?', 'Three rations will get me through the first stretch. If you agree, I will settle here with your people.')}
        topics.append(explanations[kind])
    if not topics: return []
    finals=[{'label':label,'reaction':final_reaction(kind,i,'reserved'),
             'next':-1,'action':i} for i,label in enumerate(LABELS[kind])]
    entry=[{'label':question,'reaction':'Let me explain.','next':i+1,'action':-1}
           for i,(question,answer) in enumerate(topics)]
    entry.append(finals[dialogue.decline_index(kind)].copy())
    return dialogue.validate([{'line':'We can talk before you decide.','choices':entry}]+
                             [{'line':answer,'choices':[choice.copy() for choice in finals]}
                              for question,answer in topics],kind)
