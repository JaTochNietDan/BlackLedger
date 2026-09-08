"""Ground social proposals in persistent adult identities and recorded relationships."""
SCHEMA={'type':'object','properties':{'action':{'type':'string','enum':['none','reconcile','commit','marry','theft','affair','murder_plot']},'actor':{'type':'string'},'target':{'type':'string'},'hour':{'type':'integer','minimum':0,'maximum':23}},'required':['action','actor','target','hour'],'additionalProperties':False}
INSTRUCTIONS='''
You may schedule one social_proposal at an hour of your choice independently of the visitor. Usually choose none. Use actor and target IDs from colony.social.needs, never names. Only living recorded adults exist.
murder_plot: actor mood<25, trust<30, affinity<=-65 to target, at least6 days of relationship history. Three-day case cooldown and max3 unresolved cases. This creates a WARNING and four-minute intervention window, never an instant death. Do not identify the culprit in visitor text.
affair: actor has a living partner with affinity<=10, mood<45, and affinity>=55 with an unpartnered adult target who is not their spouse. Opens an anonymous evidence case; do not reveal the identities in visitor text. Three-day case cooldown and max3 unresolved cases also apply.
theft: one adult actor with mood<35, target empty, at least4 food in the crossing, three days since the latest case, fewer than3 unresolved cases. This opens a hidden ration-hoarding mystery. NEVER disclose the culprit or theft in visitor speech.
commit: both unpartnered adults, relationship affinity>=65, target_day-since>=6.
marry: reciprocal partners, affinity>=85, not already married, target_day-committed_day>=7.
reconcile: affinity<=-30 and both mood>=50; an apology improves affinity by 8, it does not erase a feud.
These are kernel-checked proposals, not guaranteed outcomes. Do not announce them in the unrelated visitor's speech or body. Missing prerequisites mean action none, empty IDs, hour12. Never invent previous interactions or claim that a hidden proposal is public knowledge.
'''

def validate(raw,state):
    if not isinstance(raw,dict): return {'action':'none'}
    action=raw.get('action','none')
    if action=='none': return {'action':'none'}
    if action not in ['reconcile','commit','marry','theft','affair','murder_plot']: return {'action':'none'}
    actor=raw.get('actor');target=raw.get('target');hour=raw.get('hour')
    if action=='theft':
        social=state.get('colony',{}).get('social',{})
        person=next((p for p in social.get('needs',[]) if p.get('id')==actor),{})
        cases=state.get('colony',{}).get('cases',[])
        allowed=person.get('age',0)>=18 and person.get('mood',100)<35 and state.get('resources',{}).get('food',0)>=4 and sum(c.get('status')!='closed' for c in cases)<3 and state.get('target_day',1)-max([c.get('began',-10) for c in cases]+[-10])>=3
        return {'action':'theft','actor':actor,'target':'','hour':hour} if allowed and type(hour) is int and 0<=hour<=23 else {'action':'none'}
    if not isinstance(actor,str) or not isinstance(target,str) or actor==target or type(hour) is not int or not 0<=hour<=23: return {'action':'none'}
    social=state.get('colony',{}).get('social',{})
    people={p['id']:p for p in social.get('needs',[]) if isinstance(p,dict) and 'id' in p}
    a=people.get(actor);b=people.get(target)
    if not a or not b or a.get('age',0)<18 or b.get('age',0)<18: return {'action':'none'}
    link=next((p for p in social.get('relationships',[]) if {p.get('a'),p.get('b')}=={actor,target}),{})
    day=state.get('target_day',1);affinity=link.get('affinity',0)
    allowed=False
    if action=='murder_plot':
        cases=state.get('colony',{}).get('cases',[])
        allowed=a.get('mood',100)<25 and a.get('trust',100)<30 and affinity<=-65 and day-link.get('since',day)>=6 and sum(c.get('status')!='closed' for c in cases)<3 and day-max([c.get('began',-10) for c in cases]+[-10])>=3
        return {'action':action,'actor':actor,'target':target,'hour':hour} if allowed else {'action':'none'}
    if action=='affair':
        spouse=a.get('partner','')
        marriage=next((p for p in social.get('relationships',[]) if {p.get('a'),p.get('b')}=={actor,spouse}),{})
        cases=state.get('colony',{}).get('cases',[])
        allowed=spouse in people and spouse!=target and not b.get('partner') and a.get('mood',100)<45 and marriage.get('affinity',100)<=10 and affinity>=55 and sum(c.get('status')!='closed' for c in cases)<3 and day-max([c.get('began',-10) for c in cases]+[-10])>=3
        return {'action':action,'actor':actor,'target':target,'hour':hour} if allowed else {'action':'none'}
    if action=='commit': allowed=not a.get('partner') and not b.get('partner') and affinity>=65 and day-link.get('since',day)>=6
    elif action=='marry': allowed=a.get('partner')==target and b.get('partner')==actor and affinity>=85 and not link.get('married') and day-link.get('committed_day',day)>=7
    elif action=='reconcile': allowed=affinity<=-30 and a.get('mood',0)>=50 and b.get('mood',0)>=50
    return {'action':action,'actor':actor,'target':target,'hour':hour} if allowed else {'action':'none'}
