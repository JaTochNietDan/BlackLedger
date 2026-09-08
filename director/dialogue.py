"""Finite conversation graphs. The simulation owns every material consequence."""
NPC_CONTRACTS={
            'trade': 'I SELL 16 rounds to the player and RECEIVE 6 scrap. Alternatively I BUY 4 rations FROM the player and PAY them 7 scrap. I do not sell rations or accept food as payment for bullets.',
            'recruit': 'I alone ask to join the colony; the player spends 3 food to welcome me. My faction is not joining.',
            'aid': 'The player gives me 2 medicine; I give them 6 scrap and our goodwill improves.',
            'envoy': 'The player gives me 5 rounds; I give them 8 rations and our goodwill improves.',
            'cache': 'The player gives me 3 rations; they receive 10 scrap and 5 rounds.',
}
COUNTS={'war_envoy':3,'recruit':2,'trade':3,'aid':2,'envoy':2,'cache':2,'mediation':3}
SCHEMA={'type':'array','minItems':2,'maxItems':5,'items':{'type':'object','properties':{
    'line':{'type':'string'},
    'choices':{'type':'array','minItems':2,'maxItems':3,'items':{'type':'object','properties':{
        'label':{'type':'string'},'reaction':{'type':'string'},
        'next':{'type':'integer','minimum':-1,'maximum':4},
        'action':{'type':'integer','minimum':-1,'maximum':2}},
        'required':['label','reaction','next','action'],'additionalProperties':False}}},
    'required':['line','choices'],'additionalProperties':False}}
# Encode terminal-vs-discussion structure in the grammar itself.
_choice=SCHEMA['items']['properties']['choices']['items']
import copy
_discuss=copy.deepcopy(_choice)
_discuss['properties']['action']={'type':'integer','enum':[-1]}
_discuss['properties']['next']={'type':'integer','minimum':1,'maximum':4}
_final=copy.deepcopy(_choice)
_final['properties']['action']={'type':'integer','minimum':0,'maximum':2}
_final['properties']['next']={'type':'integer','enum':[-1]}
SCHEMA['items']['properties']['choices']['items']={'anyOf':[_discuss,_final]}
INSTRUCTIONS='''
Prefer exactly TWO nodes: node 0 may ask a question with next=1/action=-1; node 1 contains ONLY final decisions with next=-1. An array of length 2 has no node 2. No branch may point beyond the array. Write dialogue: a finite array of 2–5 conversation nodes, entry node 0. Each node has line (NPC first-person speech, <=45 words) and 2–3 choices with label, reaction, next, action.
The player can ask about motives, challenge terms, or choose an actual supported outcome. Vary questions and responses around the visitor's dilemma and recorded history. No invented prior deeds or secret knowledge.
For a discussion choice action=-1 and next is a STRICTLY LATER node index. No costs or rewards for discussion. Reaction is a brief NPC spoken reply (<=25 words), without claiming effects.
For a final decision next=-1 and action selects the visitor's EXACT existing outcome index:
recruit: 0 join for 3 food, 1 refuse.
trade: 0 buy 16 ammo for 6 scrap, 1 sell 4 food for 7 scrap, 2 decline.
aid: 0 give 2 medicine for 6 scrap and goodwill, 1 refuse.
envoy: 0 give 5 ammo for 8 food and goodwill, 1 refuse.
cache: 0 give 3 food for 10 scrap and 5 ammo, 1 decline.
Never promise additional rewards, new quests, unsupported bargains, automatic betrayal or deaths. Consequences and costs are rendered by the game. Every node must have a path to the decline outcome, and every node must be reachable from node 0. Final reactions should express the visitor’s feelings in a brief farewell or thanks, without repeating quantities, prices or transfer directions. The game renders exact final action labels. Write actual spoken lines, not descriptions of what the speaker does.
'''


def decline_index(kind):
    return 1 if kind == 'war_envoy' else COUNTS[kind]-1

def validate(raw,kind):
    if not isinstance(raw,list) or not 2<=len(raw)<=5 or kind not in COUNTS:
        raise ValueError('invalid dialogue size')
    nodes=[]
    def text(value,limit):
        if not isinstance(value,str) or not value.strip(): raise ValueError('missing dialogue text')
        return ' '.join(value.split())[:limit]
    for i,node in enumerate(raw):
        if not isinstance(node,dict): raise ValueError('invalid node')
        choices=node.get('choices')
        if not isinstance(choices,list) or not 2<=len(choices)<=3: raise ValueError('invalid choice count')
        clean=[]
        for choice in choices:
            if not isinstance(choice,dict): raise ValueError('invalid choice')
            nxt=choice.get('next'); action=choice.get('action')
            if type(nxt) is not int or type(action) is not int: raise ValueError('invalid reference')
            if action==-1:
                if not i<nxt<len(raw): raise ValueError('cyclic or dangling branch')
            elif not (0<=action<COUNTS[kind] and nxt==-1): raise ValueError('invalid consequence')
            clean.append({'label':text(choice.get('label'),60),'reaction':text(choice.get('reaction'),180),'next':nxt,'action':action})
        nodes.append({'line':text(node.get('line'),280),'choices':clean})
    reachable={0}
    for i,node in enumerate(nodes):
        if i in reachable: reachable.update(c['next'] for c in node['choices'] if c['next']>=0)
    if len(reachable)!=len(nodes): raise ValueError('unreachable node')
    can_decline=set()
    for i in reversed(range(len(nodes))):
        if any(c['action']==decline_index(kind) or c['next'] in can_decline for c in nodes[i]['choices']): can_decline.add(i)
    if len(can_decline)!=len(nodes): raise ValueError('conversation traps player')
    return nodes


def visitor_history(plan, colony):
    """Recorded game outcomes have explicit player perspective; speech remains attributed."""
    matches=[row for row in colony.get('threads',[]) if row.get('name')==plan.get('name') and row.get('faction',plan.get('faction'))==plan.get('faction')]
    if not matches: return []
    row=matches[-1]
    result={key:row[key] for key in ('name','began','status') if key in row}
    result['records']=[]
    recent=row.get('records',[])[-3:]
    earlier=row.get('earlier_record')
    records=([earlier] if isinstance(earlier,dict) and earlier not in recent else [])+recent
    for record in records:
        item={key:record[key] for key in ('day','decision','conversation') if key in record}
        outcome=record.get('outcome',{})
        if outcome:
            item['confirmed_game_outcome']={key:outcome[key] for key in ('player_spent','player_received','visitor_joined','faction_goodwill_change','goodwill_changes','power_changes','stability_changes') if key in outcome}
            item['perspective']='PLAYER means the crossing, not this visitor. Spent is removed from the crossing; received is added to the crossing. Joining food is a welcome cost, not a trade.'
        else:
            item['outcome_detail']='Legacy record: only the original decision label is known. Do not invent quantities or who traded which goods.'
        result['records'].append(item)
    return [result]

def arrival_presentation():
    return {'people':1,'visible_weapon':'none','transport':'on foot','meeting':'present resident; no gate assumed'}

def public_context(plan, state):
    """Only the visitor's current contract and established public facts reach their voice."""
    faction = plan.get('faction', '')
    politics = state.get('politics', {})
    colony = state.get('colony', {})
    enemies = [next(name for name in key.split('|') if name != faction)
               for key, value in politics.get('relations', {}).items()
               if faction and faction in key.split('|') and len(key.split('|')) == 2 and value <= -40]
    day = state.get('current_day')
    expiry = politics.get('ceasefires', {}).get(faction)
    treaty = {'status': 'not recorded'}
    if type(day) is int and day > 0:
        active = type(expiry) is int and expiry > day
        treaty = {'status': 'active' if active else 'no active ceasefire', 'current_day': day}
        if active:
            treaty.update(expires_at_dawn_of_day=expiry,
                          terms='My faction will not raid or plot against the crossing before expiry. Seizing our toll breaks this treaty. This treaty covers only my faction and the independent crossing. It says nothing about whether any regional war exists.')
    contract = ''
    if plan.get('kind')=='mediation':
        contract = ('I represent '+faction+' in a dispute with '+str(plan.get('other',''))+'. Hosting costs 6 food, improves our goodwill by 8 and schedules talks two days after acceptance, without guaranteeing peace. Arming us costs 8 rounds, improves our goodwill by 12, reduces rival goodwill by 15 and worsens our mutual relations by 15. Neutrality costs nothing. Terms are rechecked before payment.')
    if plan.get('return_motive') == 'asylum':
        contract = 'Accepting shelter costs 3 food and loses 12 goodwill with my faction. Refusal does not change faction goodwill.'
    elif plan.get('return_motive') == 'war_supplies':
        contract = ('Exchanging ammunition for food strengthens my faction by 4 power and costs 5 goodwill with each current war enemy: ' + ', '.join(enemies) + '.'
                    if enemies else 'My faction has no current war enemies. This is now an ordinary envoy exchange: no military power gain or enemy goodwill penalty.')
        contract += ' Alternatively, civilian relief costs the crossing 3 medicine, gives no goods back, improves our goodwill by 8 and stability by up to 4 (capped at 100), without military power or rival goodwill changes. Declining remains outcome 1; civilian relief is outcome 2.'
    return {
        'arrival_presentation': arrival_presentation(),
        'speaker_role': 'I am an outside visitor. The crossing belongs to the player’s residents, not my faction; I do not command its defenses.',
        'current_war_enemies': enemies,
        'political_contract': contract,
        'exchange_status': 'Offered only. The player has not accepted this visit’s exchange. Existing ceasefires do not imply acceptance of a trade.',
        'crossing_ceasefire': treaty,
        'npc_contract': contract if plan.get('kind')=='mediation' else NPC_CONTRACTS.get(plan.get('kind'), ''),
        'visitor': {key: plan[key] for key in ('kind', 'name', 'faction', 'return_motive', 'other') if key in plan},
        'current_faction': [{key: row[key] for key in ('name', 'leader', 'agenda', 'type') if key in row}
                            for row in politics.get('factions', []) if row.get('name') == faction][:1],
        'public_faction_profiles': [{key: row[key] for key in ('name', 'leader', 'agenda', 'type') if key in row}
                                    for row in politics.get('factions', []) if row.get('name')][:8],
        'public_relations': {key: value for key, value in politics.get('relations', {}).items()
                             if faction and faction in key.split('|')},
        'public_relationship_meanings': [
            {'factions': key.split('|'), 'status': 'allied' if value>=40 else 'at war' if value<=-40 else 'neither allied nor at war'}
            for key,value in politics.get('relations',{}).items() if faction and faction in key.split('|')],
        'goodwill_with_crossing': state.get('faction_relations', {}).get(faction, 0),
        'visitor_history': visitor_history(plan,colony),
        'faction_delegations': [{key: row[key] for key in ('day','action','status','paid','received') if key in row}
                                for row in politics.get('recent_delegations', [])
                                if isinstance(row, dict) and row.get('faction') == faction
                                and row.get('status') == 'completed player delegation'][-4:],
        'loss_day': day,
        'faction_loss_grievances': [{key: row[key] for key in ('person_id','name','lost_person_id','lost_name','day') if key in row}
                                   for row in colony.get('social', {}).get('political_loss_grievances', [])
                                   if isinstance(row, dict) and faction and row.get('faction') == faction][-3:],
        'loss_recall_scope': 'Recorded personal losses attributed to my faction. They do not establish my personal involvement, present mourning, revenge plans, or acceptance of this offer.',
        'plot_day': day,
        'faction_plot_outcomes': [{key: row[key] for key in ('kind','status','outcome','resolved_day','target_name','target_hp_before','target_hp_after') if key in row}
                                  for row in politics.get('archive', [])
                                  if isinstance(row, dict) and faction and row.get('faction') == faction and row.get('status') == 'resolved'][-3:],
        'plot_recall_scope': 'Historical crossing records about my faction, not proof I took part, witnessed the event or knew about the plot beforehand. A historical injury does not establish current health or survival. No new threat or agreement is implied.',
        'delegation_day': day,
        'delegation_scope': 'Recorded dealings with my faction, not evidence I personally visited or witnessed them. Transfers use the crossing perspective. These do not accept the present offer.',

        'living_residents': [row.get('name') for row in state.get('survivors', [])
                             if row.get('name') and not row.get('away') and row.get('health',row.get('hp',100))>0],
        'away_residents': [row.get('name') for row in state.get('survivors', [])
                           if row.get('name') and row.get('away') and row.get('health',row.get('hp',100))>0],
        'resident_availability_scope': 'living_residents are living residents currently at home. away_residents are absent on expeditions; do not put them in this scene, assume a safe return, or claim they have permanently departed.',
        'known_dead': [row.get('name') for row in state.get('dead', []) if row.get('name')],
        'known_departed': [row.get('name') for row in state.get('departed', []) if row.get('name')],
        'recorded_history': state.get('recent_events', [])[-4:],
        'old_public_records': state.get('recalled_events', [])[:4],
    }
