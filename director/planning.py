"""Small planning context; full state remains authoritative for validation."""
import copy
from collections import Counter
try:
    from .dialogue import NPC_CONTRACTS, arrival_presentation
except ImportError:
    from dialogue import NPC_CONTRACTS, arrival_presentation

def fields(value, keys):
    return {k:copy.deepcopy(value[k]) for k in keys if k in value}

def settlement(state):
    """Summarize placed structures separately from unfinished construction and legacy totals."""
    colony = state.get('colony', {})
    built, planned, damaged = Counter(), Counter(), Counter()
    unavailable, destroyed, demolition = Counter(), Counter(), Counter()
    for structure in colony.get('structures', []):
        kind = structure.get('kind')
        if isinstance(kind, str) and kind:
            (built if structure.get('done') is True else planned)[kind] += 1
            hp, maximum = structure.get('hp'), structure.get('max_hp')
            finished = structure.get('done') is True
            dead = isinstance(hp, (int, float)) and hp <= 0
            removing = structure.get('deconstruct') is True
            if finished and (dead or removing):
                unavailable[kind] += 1
            if finished and dead:
                destroyed[kind] += 1
            if finished and removing:
                demolition[kind] += 1
            if structure.get('done') is True and isinstance(hp, (int, float)) and isinstance(maximum, (int, float)) and 0 < hp < maximum and not removing:
                damaged[kind] += 1
    result = {'completed': dict(sorted(built.items())), 'under_construction': dict(sorted(planned.items()))}
    if damaged:
        result['damaged'] = dict(sorted(damaged.items()))
    for key, counts in (('unavailable', unavailable), ('destroyed', destroyed), ('marked_for_demolition', demolition)):
        if counts:
            result[key] = dict(sorted(counts.items()))
    return result


def encounter_options(state):
    """Bind urgent survival scenes to existing transactions, without awarding outcomes."""
    colony=state.get('colony', {})
    meal=colony.get('last_rations', {})
    day=state.get('current_day', 0)
    residents=len(state.get('survivors', []))
    recent=0 <= day-meal.get('day', -99) <= 1
    if recent and meal.get('shortage', 0)>0 and state.get('resources', {}).get('food', 0)<max(1,residents*2):
        return [
            {'kind':'envoy','pressure':'food_shortage','dilemma':'Trade defense supplies for food, or retain ammunition.',
             'cost':{'ammo':5},'gain':{'food':8},'evidence':{'day':meal['day'],'shortage':meal['shortage']}},
            {'kind':'cache','pressure':'food_shortage','dilemma':'Spend scarce rations on salvage and ammunition, or keep food for residents.',
             'cost':{'food':3},'gain':{'scrap':10,'ammo':5},'evidence':{'day':meal['day'],'shortage':meal['shortage']}}
        ]
    options=[{'kind':kind} for kind in ('recruit','trade','aid','envoy','cache')]
    history=colony.get('recent_encounters',[])
    if len(history)>=2:
        previous,last=history[-2:]
        if previous.get('kind')==last.get('kind') and 0<=day-last.get('day',-99)<=2:
            options=[option for option in options if option['kind']!=last.get('kind')]
    return options + ([{'kind':'mediation','pairs':summit_options(state)}] if summit_options(state) else [])


def context(state):
    out=fields(state,('current_day','target_day','generation','resources','wall_health','buildings','faction_relations','last_visitor','reserved_cast_names'))
    out['arrival_presentation']=arrival_presentation()
    out['survivors']=[fields(p,('name','health','role','away','shooting','building','medicine','melee','trait','perks','weapon','magazine','stance','armor','blade')) for p in state.get('survivors',[])[:6]]
    # All known dead/departed names remain excluded, even when their detailed biography is omitted.
    out['unavailable_people']=[fields(p,('name','day','cause')) for p in state.get('dead',[])+state.get('departed',[])]
    out['recent_events']=[fields(r,('day','generation','record_type','cue','text')) for r in state.get('recent_events',[])[-6:]]
    out['recalled_events']=[fields(r,('day','generation','record_type','text')) for r in state.get('recalled_events',[])[-3:]]
    politics=state.get('politics',{})
    out['politics']=fields(politics,('factions','relations','plots','ceasefires','last_birth','summits'))
    out['politics']['plots']=[p for p in out['politics'].get('plots',[]) if not p.get('covert')]
    out['politics']['summit_options'] = summit_options(state)
    out['politics']['summit_history'] = [fields(r,('a','b','resolved_day','security','outcome','attacker','relation_before','relation_after','goodwill_before','goodwill_after','trust','player_spent','moves','power_before','power_after')) | {'leaders':[fields(leader,('name',)) for leader in r.get('leaders',[])]} for r in politics.get('summit_history',[])[-6:]]
    out['politics']['resolved_outcomes'] = []
    for record in politics.get('archive', [])[-6:]:
        if record.get('covert'): continue
        item = fields(record, ('kind','faction','resolved_day','outcome','target_name','structure_name','structure_destroyed','structure_missing','actual_losses','actual_received'))
        item['leader_at_warning'] = fields(record.get('instigator', {}), ('name',))
        out['politics']['resolved_outcomes'].append(item)
    out['politics']['recent_delegations']=politics.get('recent_delegations',[])[-2:]
    colony=state.get('colony',{})
    out['colony']={'threads':[], 'people':colony.get('people',[])[:6], 'cases':colony.get('cases',[])[:3], 'social':colony.get('social',{})}
    out['encounter_options']=[dict(option,npc_contract=NPC_CONTRACTS.get(option['kind'],'Political mediation; use the listed summit options.')) for option in encounter_options(state)]
    if 'resources' in state:
        costs={'recruit':[{'food':3}], 'trade':[{'scrap':6},{'food':4}],
               'aid':[{'medicine':2}], 'envoy':[{'ammo':5}], 'cache':[{'food':3}]}
        resources=state['resources']
        unavailable=[option['kind'] for option in out['encounter_options']
                     if option['kind'] in costs and not any(
                         all(resources.get(item,0)>=quantity for item,quantity in cost.items())
                         for cost in costs[option['kind']])]
        if unavailable: out['unaffordable_paid_encounters_now']=unavailable
    out['colony']['recent_attacks']=[]
    for attack in colony.get('recent_attacks', [])[-4:]:
        item=fields(attack,('id','day','faction','threat','outcome','confirmed_kills','details_recorded','structural_damage_recorded'))
        if attack.get('details_recorded') is True:
            for key, allowed, limit in (
                ('losses', ('name',), 6),
                ('net_health_loss', ('name','health_lost'), 6),
                ('structures_lost_or_spent', ('kind','count'), 32),
                ('structures_damaged', ('kind','count','health_lost'), 32),
                ('weapons_dropped', ('weapon','count'), 12)):
                item[key]=[fields(row,allowed) for row in attack.get(key,[])[:limit]]
        out['colony']['recent_attacks'].append(item)
    production=colony.get('production',{})
    out['colony']['production']={
        'stations':[fields(s,('kind','operational','under_construction','unavailable')) for s in production.get('stations',[])[:2]],
        'orders':[fields(o,('recipe','station','held','progress_percent','reserved_materials','planned_output','assigned_crafter')) for o in production.get('orders',[])[:5]],
        'meaning':'Queued outputs are not owned inventory or completed trades. Reserved materials are already deducted from current resources; do not deduct them again. Operational stations may still be unreachable. Assigned crafters may be walking or paused; holding an order stops its production.'}
    order_count=production.get('order_count',len(production.get('orders',[])))
    shown=len(out['colony']['production']['orders'])
    if isinstance(order_count,int) and order_count>shown:
        out['colony']['production']['order_count']=order_count
        out['colony']['production']['orders_not_shown']=order_count-shown
    out['colony']['settlement']=settlement(state)
    out['colony']['last_rations']=fields(colony.get('last_rations',{}),('day','required','served','shortage','remaining'))
    out['colony']['shelter']=fields(colony.get('shelter',{}),('residents','sheltered_beds','shortfall'))
    out['colony']['land']=fields(colony.get('land',{}),('scenario','remaining','cleared_deposits'))
    growth=colony.get('land',{}).get('woodland_regrowth',{})
    if growth:
        out['colony']['land']['woodland_regrowth']=fields(growth,('scheduled_sites','suppressed_sites','awaiting_schedule','earliest_maturity_day','potential_wood'))
        out['colony']['land']['woodland_regrowth']['meaning']='Future harvestable timber, not stored resources or currently harvestable land. Cleared woodland takes six days to regrow. Construction suppresses growth; occupied trunk sites delay maturity. Earliest day is conditional, not a guaranteed delivery. Colonists must harvest mature trees.'
    out['colony']['supply_market']=fields(colony.get('commodity_market',{}),('day','disruption','stock','prices','credits'))
    out['colony']['arms_market']=fields(colony.get('arms_market',{}),('day','disruption','stock','prices'))
    out['colony']['stored_equipment']=fields(colony.get('inventory',{}),('weapons','machete','vest','kit'))
    out['colony']['expeditions']=[fields(e,('name','kind','remaining')) for e in colony.get('expeditions',[])[:2]]
    if colony.get('recent_encounters'):
        out['colony']['recent_encounters']=[fields(e,('day','kind','name','faction','decision')) for e in colony['recent_encounters'][-6:]]
    if colony.get('recent_trades'):
        out['colony']['recent_trades']=[fields(t,('day','direction','category','item','name','quantity','credits')) for t in colony['recent_trades'][-6:]]
    out['colony']['recent_expeditions']=[fields(e,('name','kind','day','outcome','received')) for e in colony.get('recent_expeditions',[])[-4:]]
    out['colony']['training']=fields(colony.get('training',{}),('teacher','student','skill','progress_seconds','waiting_seconds','reserved_materials'))
    out['colony']['requests']=[fields(r,('name','kind','title','requirement','status','target','deadline','outcome')) for r in colony.get('requests',[])[-8:]]
    for thread in colony.get('threads',[])[-4:]:
        item=fields(thread,('name','faction','status','closed_reason','last_kind','last_motive','personal_trust','helped','next_day'))
        if thread.get('earlier_record'): item['earlier_record']=fields(thread['earlier_record'],('day','decision','outcome','kind'))
        item['records']=[fields(r,('day','decision','outcome','kind')) for r in thread.get('records',[])[-2:]]
        out['colony']['threads'].append(item)
    return out


def summit_options(state):
    world=state.get('politics',{})
    if len(world.get('summits',[]))>=3:return []
    names=[f['name'] for f in world.get('factions',[]) if 'name' in f]
    pending={frozenset((s.get('a'),s.get('b'))) for s in world.get('summits',[])}
    return [{'actor':a,'target':b,'relation':world.get('relations',{}).get('|'.join(sorted((a,b))),0)}
            for i,a in enumerate(names) for b in names[i+1:]
            if world.get('relations',{}).get('|'.join(sorted((a,b))),0)<=-40 and frozenset((a,b)) not in pending]


def political_actions(state):
    world=state.get('politics',{})
    factions=world.get('factions',[])
    actions=['none']
    if len(factions)>1:actions.append('dispute')
    if summit_options(state):actions.append('summit')
    if len(factions)<8 and state.get('target_day',2)-world.get('last_birth',0)>=4 and any(f.get('stability',100)<=40 for f in factions):actions.append('splinter')
    day=state.get('target_day',2)
    if len(world.get('plots',[]))<2 and state.get('survivors') and any(state.get('faction_relations',{}).get(f['name'],0)<=-25 and world.get('ceasefires',{}).get(f['name'],0)<=day and not any(p.get('faction')==f['name'] for p in world.get('plots',[])) for f in factions):actions.append('plot')
    return actions
