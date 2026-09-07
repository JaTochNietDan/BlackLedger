"""Black Ledger's deterministic, action-clock world. No model-written state patches."""
import copy
import hashlib
import random
import uuid

LOCATIONS = [
    dict(id='room', name='The Mariner', type='home', district=0, x=220,y=205, cost=0, blurb='A rented room above the waterfront. Thin walls. One stairway. Yours, for now.'),
    dict(id='bar', name='Saint Agnes', type='bar', district=0, x=415,y=210,cost=0,blurb='Coffee by daylight, whispered arrangements after dark. Your contact keeps the corner booth.'),
    dict(id='docks',name='Pier 14',type='work',district=0,x=130,y=420,cost=0,blurb='Cargo arrives every morning. Honest work and less honest favors both pay cash.'),
    dict(id='laundry',name='Bluebird Laundry',type='racket',district=0,x=320,y=425,cost=180,blurb='A neighborhood business seeking protection. A small beginning, if you can honor the arrangement.'),
    dict(id='club',name='The Monarch',type='casino',district=0,x=570,y=380,cost=0,blurb='Bellandi men watch the doors. This casino belongs to a family that does not tolerate outsiders making demands.'),
    dict(id='market',name='Mercer Exchange',type='market',district=0,x=585,y=180,cost=0,blurb='People sell information here. Its value depends on whom you ask, and who sees you asking.'),
    dict(id='apartment',name='Ashbury Court',type='home',district=1,x=790,y=200,cost=180,blurb='A private entrance, a respectable address, and a landlord willing to mind his own business.'),
    dict(id='garage',name='Russo Motor Works',type='racket',district=1,x=800,y=405,cost=350,blurb='A working garage at the edge of Russo territory. An investment that requires political tact.'),
    dict(id='casino',name='The Blue Hour',type='casino',district=1,x=700,y=585,cost=850,blurb='A shuttered gambling house. Restore its reputation and the tables could support an entire crew.'),
    dict(id='estate',name='Cypress House',type='home',district=2,x=930,y=95,cost=3500,blurb='High walls above the city. A statement of success, and a very conspicuous address.'),
]
BY_ID={p['id']:p for p in LOCATIONS}
HOUSING={'room':dict(name='Rented room',rent=15,security=0),'apartment':dict(name='Private apartment',rent=35,security=1),'estate':dict(name='Cypress estate',rent=90,security=2)}

class InvalidAction(ValueError): pass

def record(w,title,text,kind='city'):
    row=dict(id=uuid.uuid4().hex[:12],minute=w['minute'],life=w['life'],title=title,text=text,kind=kind)
    w['history'].append(row)
    w['history']=w['history'][-180:]
    return row

def new_person(life):
    names=['Alex Varga','Nico Ward','Frankie Vale','Sam Costa','Jamie Moretti','Robin Hale']
    return dict(name=names[(life-1)%len(names)],cash=90,health=100,respect=0,heat=0,location='room',home='room',security=0,contacts=0,crew=[],alive=True,earned=0,job_count=0)

def new_world(seed=27):
    w=dict(version=1,id=uuid.uuid4().hex,revision=0,life=1,minute=8*60,rng=seed,player=new_person(1),district=0,
      factions=[dict(id='bellandi',name='Bellandi Family',leader='Vittorio Bellandi',power=90,goodwill=0,cash=8000),dict(id='russo',name='Russo Outfit',leader='Elena Russo',power=58,goodwill=0,cash=4500)],
      npcs=[dict(id='mara',name='Mara Bell',role='Fixer',trust=10,voice='af_heart',color='#a48761'),dict(id='leo',name='Leo Carver',role='Driver',trust=20,voice='am_michael',color='#9ca795'),dict(id='vittorio',name='Vittorio Bellandi',role='Bellandi boss',trust=0,voice='bm_george',color='#ad7970'),dict(id='elena',name='Elena Russo',role='Russo boss',trust=0,voice='bf_emma',color='#83989b')],
      properties={p['id']:dict(owner='bellandi' if p['id']=='club' else 'independent',condition=100,income=14 if p['id']=='laundry' else 24 if p['id']=='garage' else 48 if p['id']=='casino' else 0,carry=0) for p in LOCATIONS},
      plots=[],tasks=[],event=None,history=[],dead=[],director=dict(status='authored',detail='Authored opening. Local AI can prepare additional encounters.',last_request=-9999),offers=[],last_result=None)
    record(w,'A room. A name. No protection.','Mara Bell left word at Saint Agnes: there is work, if you can be discreet. Your room costs $15 each midnight.','personal')
    return w

def roll(w):
    w['rng']=(1664525*w['rng']+1013904223)%2**32
    return w['rng']/2**32

def npc(w,id): return next((p for p in w['npcs'] if p['id']==id),w['npcs'][0])
def own(w,id): return w['properties'][id]['owner']==f"player:{w['life']}"
def guard(w): return w['player']['security']+HOUSING[w['player']['home']]['security']
def pay(w,amount):
    if w['player']['cash']<amount: raise InvalidAction('You cannot afford this commitment.')
    w['player']['cash']-=amount

def earn(w,amount): w['player']['cash']+=amount;w['player']['earned']+=amount

def options(w,place=None):
    p=w['player'];loc=BY_ID[place or p['location']];here=p['location']==loc['id'];out=[]
    def add(id,label,minutes=0,cost=0,reason='',detail=''):
        out.append(dict(id=id,label=label,minutes=minutes,cost=cost,disabled=bool(reason) or p['cash']<cost,reason=reason or ('Not enough cash' if p['cash']<cost else ''),detail=detail,target=loc['id']))
    if not p['alive'] or w['event']:return out
    if loc['district']>w['district']:
        add('expand','Establish contacts across town',90,100,reason='Earn 10 respect first' if p['respect']<10 else '',detail='Opens the next district. Existing businesses and rivals remain.');return out
    if not here:
        minutes=travel_minutes(p['location'],loc['id']);add('travel','Visit '+loc['name'],minutes,detail='Travel advances the city clock. Known threats may interrupt you.');return out
    if loc['id']=='bar':
        add('courier','Carry a discreet envelope',45,detail='Earn $45 and 2 respect. A reliable introduction to the neighborhood.')
        add('contact','Buy Mara a coffee',30,10,detail='Build trust and an information network. Contacts may warn you of trouble.')
        add('recruit','Recruit Leo Carver',30,90,reason='Leo is already in your crew' if p['crew'] else 'Earn 6 respect first' if p['respect']<6 else '',detail='A driver and collector. $12 daily wages; loyalty matters.')
    if loc['id']=='docks':add('dockwork','Work the night cargo',90,detail='Earn $75 and 1 respect. Small chance of a work injury.')
    if loc['id']=='market':
        add('investigate','Ask about threats',45,30,detail='Investigate your exposure. Evidence reveals existing threats; it does not create them.')
        add('lie_low','Keep a low profile',120,15,detail='Lose 10 heat. Time still passes for rivals and businesses.')
    if loc['id']=='club':
        add('audience','Request an audience',45,detail='Discuss your standing with the Bellandi family.')
        add('provoke','Demand protection money',30,detail='EXTREME RISK. Bellandi owns this casino. Challenging him can bring lethal retaliation.')
    if loc['id'] in ['laundry','garage','casino']:
        if own(w,loc['id']):
            add('inspect','Meet the manager',30,detail='Inspect income and repair needs; build neighborhood standing.')
            add('repair','Repair the property',60,50,reason='Already in good condition' if w['properties'][loc['id']]['condition']>=100 else '',detail='Restore 40 condition.')
        else:
            req=6 if loc['id']=='laundry' else 10 if loc['id']=='garage' else 20
            add('acquire','Establish protection' if loc['type']=='racket' else 'Reopen the casino',60,loc['cost'],reason=f'Earn {req} respect first' if p['respect']<req else 'This property belongs to another organization' if w['properties'][loc['id']]['owner']!='independent' else '',detail=f"Earn up to ${w['properties'][loc['id']]['income']}/hour. Income accrues automatically; rivals may take notice.")
    if loc['type']=='home':
        if p['home']!=loc['id']:
            add('move_home','Buy this residence' if loc['id']=='estate' else 'Rent this apartment',60,loc['cost'],detail=f"${HOUSING[loc['id']]['rent']}/day upkeep. Base security {HOUSING[loc['id']]['security']}. Moving resets hired security.")
        else:
            add('rest','Rest for four hours',240,detail='Recover 25 health and reduce heat. You are at home while time passes.')
            add('security','Hire another security detail',30,100,reason='Maximum security hired' if p['security']>=3 else '',detail='Improves detection and survival at home. Each detail adds $10/day upkeep.')
    if p['crew']:
        add('delegate','Send Leo on collections',15,reason='Leo is already on assignment' if w['tasks'] else '',detail='Completes after 120 game minutes: $65. He cannot provide support while away.')
    add('wait','Let an hour pass',60,detail='Income, operations, rent and rival plans continue.')
    return out

def travel_minutes(a,b):
    x,y=BY_ID[a],BY_ID[b];return max(10,round((abs(x['x']-y['x'])+abs(x['y']-y['y']))/12/5)*5)

def scene(w,id,title,body,speaker,choices,kind='conversation',source='authored'):
    w['event']=dict(id=id,title=title,body=body,speaker=speaker,choices=choices,kind=kind,source=source,minute=w['minute'])

def retaliation(w):
    if any(t['life']==w['life'] for t in w['plots']):return
    w['plots'].append(dict(id=uuid.uuid4().hex,kind='hit',life=w['life'],due=w['minute']+240,actor='bellandi',strength=5,known=False))

def die(w,cause):
    p=w['player'];p['alive']=False;p['health']=0
    w['dead'].append(dict(name=p['name'],minute=w['minute'],life=w['life'],cause=cause))
    for prop in w['properties'].values():
        if prop['owner']==f"player:{w['life']}":prop['owner']='former:'+p['name'];prop['income']=max(5,prop['income']-2)
    w['tasks']=[];w['event']=None
    record(w,p['name']+' is dead',cause+' Your life ends here. The city continues.','death')

def attack(w,plot):
    p=w['player'];w['plots'].remove(plot)
    if p['location']!=p['home']:
        prop=w['properties'][p['home']];prop['condition']=max(10,prop['condition']-45)
        record(w,'Someone came looking','You were away. Armed men damaged your residence and left before anyone could identify them. Repairs will take time.','danger');return
    warning=plot['known'] or guard(w)>0 or p['contacts']>=2
    if warning:
        scene(w,'attack-'+plot['id'],'Headlights outside',f"A car stops outside {BY_ID[p['home']]['name']}. "+('Your security raises the alarm.' if guard(w)>0 else 'A contact calls: leave by the back, now.')+' You have moments to act.','mara',[
            dict(id='escape',label='Leave through the rear',detail='A chance to escape. Security and contacts improve your odds.'),
            dict(id='defend',label='Hold the entrance',detail='Rely on your security. A weak defense may be fatal.'),
            dict(id='bargain',label='Offer $180 to stand down',cost=180,detail='Money may settle this incident, but your standing suffers.')],kind='attack')
    elif roll(w)<.78:die(w,'An attack at your residence caught you without warning or protection.')
    else:p['health']=max(1,p['health']-65);record(w,'You survived by inches','The attackers leave you wounded. Nobody warned you. You need medical rest and protection.','danger')

def advance(w,minutes):
    """Advance in minute steps so intervening events stop all subsequent time."""
    p=w['player']
    for _ in range(minutes):
        if not p['alive'] or w['event']:break
        w['minute']+=1
        for prop in w['properties'].values():
            if prop['owner']==f"player:{w['life']}":
                prop['carry']+=prop['income']*prop['condition']/100/60
                amount=int(prop['carry']);prop['carry']-=amount;earn(w,amount)
        for task in list(w['tasks']):
            if task['due']<=w['minute']:
                earn(w,65);record(w,'Leo returns','$65 from collections. He is available again.','business');w['tasks'].remove(task)
        if w['minute']%1440==0:
            bill=HOUSING[p['home']]['rent']+10*p['security']+12*len(p['crew'])
            if p['cash']>=bill:p['cash']-=bill;record(w,'Accounts settled',f'${bill} paid for housing, security and crew.','business')
            else:
                p['cash']=max(0,p['cash']-15);p['home']='room';p['security']=0
                if p['crew']:p['crew'][0]['loyalty']=max(0,p['crew'][0]['loyalty']-20)
                record(w,'Your arrangements unravel','You could not cover the bills. Security leaves; your residence is now a rented room. Unpaid crew lose loyalty.','danger')
        for plot in list(w['plots']):
            if plot['life']!=w['life']:continue
            if not plot['known'] and p['contacts']>=2 and plot['due']-w['minute']<=90:
                plot['known']=True;record(w,'Mara has heard something','Bellandi men have been asking where you sleep. You may have very little time.','danger')
            if plot['due']<=w['minute']:attack(w,plot);break


def action(w,a):
    p=w['player'];kind=a.get('kind');oldtime=w['minute'];oldloc=p['location'];start=len(w['history'])
    if kind=='new_life':
        if p['alive']:raise InvalidAction('This life is still in progress.')
        w['life']+=1;w['player']=new_person(w['life']);w['event']=None;w['offers']=[];w['plots']=[];w['minute']+=480
        for f in w['factions']:f['goodwill']=0
        record(w,'A stranger arrives',w['player']['name']+' rents a room at the Mariner. The previous life left its mark on the city.','personal')
    elif not p['alive']:raise InvalidAction('This life has ended.')
    elif kind=='choice':
        e=w['event']
        if not e or a.get('event')!=e['id']:raise InvalidAction('That conversation is no longer current.')
        choice=next((c for c in e['choices'] if c['id']==a.get('choice')),None)
        if not choice:raise InvalidAction('Unknown decision.')
        pay(w,choice.get('cost',0));w['event']=None;c=choice['id']
        if e['kind']=='attack':
            if c=='bargain':p['respect']=max(0,p['respect']-5);record(w,'A costly reprieve','The men accept your money and withdraw. This does not make Bellandi your friend.','danger')
            else:
                chance=min(.9,(.5 if c=='escape' else .12)+guard(w)*.16+p['contacts']*.04)
                if roll(w)>chance:die(w,'You did not survive the attack at your residence.')
                else:
                    lost=min(p['security'],1+int(roll(w)*2));p['security']-=lost;p['health']=max(1,p['health']-(20 if c=='escape' else 40));p['location']='bar' if c=='escape' else p['home']
                    record(w,'Alive, at a price',f"You survive wounded. {lost} security details are lost. "+('Mara lets you shelter at Saint Agnes.' if c=='escape' else 'The attackers retreat from your home.'),'danger')
        elif e['kind']=='audience':
            if c=='tribute':w['factions'][0]['goodwill']+=8;w['plots']=[t for t in w['plots'] if t['actor']!='bellandi'];record(w,'A temporary understanding','Bellandi accepts the tribute and calls off his current operation against you.','politics')
            else:record(w,'You leave the Monarch','Nothing was agreed. Existing threats remain.','personal')
        elif e['kind']=='proposal':
            if c=='accept':
                effect=e['effect'];earn(w,effect['reward']);p['respect']+=effect['respect'];p['heat']=min(100,p['heat']+effect['heat']);npc(w,e['speaker'])['trust']+=3
                record(w,e['title'],e['outcome']+f" (${effect['reward']:+}, respect +{effect['respect']})",'story')
                advance(w,effect['minutes'])
            else:record(w,'An offer declined','You decline '+npc(w,e['speaker'])['name']+"'s proposal. No payment changes hands.",'story')
    else:
        if w['event']:raise InvalidAction('Resolve the current situation first.')
        target=a.get('target',p['location']);available=options(w,target)
        option=next((o for o in available if o['id']==kind),None)
        if not option or option['disabled']:raise InvalidAction(option['reason'] if option else 'That action is not available here.')
        pay(w,option['cost'])
        # Travel commits destination after time; incidents therefore know the player is in transit.
        if kind=='travel':
            p['location']='transit';advance(w,option['minutes'])
            if p['alive'] and not w['event']:p['location']=target;record(w,'Arrived at '+BY_ID[target]['name'],f"The journey took {w['minute']-oldtime} minutes.",'travel')
        else:
            if kind=='courier':earn(w,45);p['respect']+=2;p['job_count']+=1;record(w,'Envelope delivered','Mara pays $45. A small favor, completed without questions.','work')
            elif kind=='dockwork':
                earn(w,75);p['respect']+=1
                if roll(w)<.15:p['health']=max(1,p['health']-10)
                record(w,'Cargo shifted','$75 for a long shift on Pier 14.','work')
            elif kind=='contact':p['contacts']=min(5,p['contacts']+1);npc(w,'mara')['trust']+=5;record(w,'A useful conversation','Mara will keep an ear open. Your information network improves.','personal')
            elif kind=='recruit':p['crew'].append(dict(id='leo',name='Leo Carver',loyalty=65));p['respect']+=2;record(w,'Your first associate','Leo Carver joins you. He expects $12 a day and a boss who keeps their word.','personal')
            elif kind=='investigate':
                threats=[t for t in w['plots'] if t['life']==w['life']]
                for t in threats:t['known']=True
                record(w,'Word on the street','Bellandi has commissioned an operation against you. Avoid home, negotiate, or arrange protection.' if threats else 'Your sources have no evidence of an active operation against you. This is not a guarantee of safety.','intel')
            elif kind=='lie_low':p['heat']=max(0,p['heat']-10);record(w,'Out of the spotlight','You avoid attention for a while.','personal')
            elif kind=='provoke':
                p['respect']+=1;w['factions'][0]['goodwill']-=35;retaliation(w);record(w,'A demand nobody forgets','The manager refuses. A Bellandi man watches you leave. You have challenged a powerful family on its own ground.','politics')
            elif kind=='audience':
                # Conversation is opened after the appointment time elapses.
                pass
            elif kind=='acquire':
                w['properties'][target]['owner']=f"player:{w['life']}";p['respect']+=4;record(w,'A foothold in the city',BY_ID[target]['name']+' now produces income for you. Earnings accrue as game time passes.','business')
                if target=='garage':w['factions'][1]['goodwill']-=10
            elif kind=='inspect':p['respect']+=1;record(w,'The books are open',f"{BY_ID[target]['name']}: {w['properties'][target]['condition']}% condition. Income depends on its condition.",'business')
            elif kind=='repair':w['properties'][target]['condition']=min(100,w['properties'][target]['condition']+40);record(w,'Repairs arranged','The property is restored by 40 condition.','business')
            elif kind=='move_home':p['home']=target;p['security']=0;p['respect']+=3;record(w,'A different view',BY_ID[target]['name']+' is now your residence. Hired security must be arranged here.','personal')
            elif kind=='rest':p['health']=min(100,p['health']+25);p['heat']=max(0,p['heat']-5);record(w,'Time to recover','You settle in at home. Recovery does not stop other people making plans.','personal')
            elif kind=='security':p['security']+=1;record(w,'Someone at the door','Another security detail is assigned to your residence. It adds $10 a day to your expenses.','personal')
            elif kind=='delegate':w['tasks'].append(dict(id=uuid.uuid4().hex,name='Leo · collections',due=w['minute']+120));record(w,'Leo heads out','Collections should be completed in two hours.','work')
            elif kind=='expand':w['district']=min(2,w['district']+1);record(w,'The city opens up','Your contacts introduce you to another district. More properties are accessible.','city')
            advance(w,option['minutes'])
            if kind=='audience' and p['alive'] and not w['event']:
                scene(w,uuid.uuid4().hex,'A seat across from Bellandi','“People mistake an open door for an invitation. Tell me you understand the difference.”','vittorio',[dict(id='tribute',label='Offer $150 in tribute',cost=150,detail='Improves relations and cancels his current operation against you.'),dict(id='leave',label='Leave without an agreement',detail='No payment. Existing threats remain.')],kind='audience')
        if p['alive'] and not w['event'] and kind not in ['travel','provoke','audience']:
            if w['offers'] and w['minute']>=w['offers'][0]['ready']:
                offer=w['offers'].pop(0);w['event']=offer['event'];w['event']['minute']=w['minute']
            elif p['job_count']==2 and not any(r['title']=='A favor with a price' and r['life']==w['life'] for r in w['history']):
                e=dict(id=uuid.uuid4().hex,title='A favor with a price',body='“A merchant wants a sealed ledger moved before his partners arrive. Fifty dollars. I would understand if you preferred the ordinary work.”',speaker='mara',kind='proposal',source='authored',minute=w['minute'],effect=dict(reward=50,respect=2,heat=3,minutes=45),outcome='You moved the ledger. Mara now knows you can handle sensitive work.',choices=[dict(id='accept',label='Take the sealed ledger',detail='$50 · 45 minutes · +2 respect · +3 heat'),dict(id='decline',label='Keep your hands clear',detail='No cost or time.')]);w['event']=e
                record(w,'A favor with a price','Mara offers more sensitive work.','story')
    if p['alive'] and p['location']=='transit':p['location']=oldloc
    w['revision']+=1
    w['last_result']=dict(from_location=oldloc,to_location=w['player']['location'],elapsed=w['minute']-oldtime,records=w['history'][start:])
    return w


def public(w):
    """Explicit projection: hidden plans, random seed and unresolved outcomes never reach the browser."""
    p=copy.deepcopy(w['player']);e=copy.deepcopy(w['event'])
    if e:
        for key in ['effect','outcome']:e.pop(key,None)
        e['choices']=[dict(c,disabled=p['cash']<c.get('cost',0)) for c in e['choices']]
    return dict(id=w['id'],revision=w['revision'],life=w['life'],minute=w['minute'],player=p,district=w['district'],factions=copy.deepcopy(w['factions']),npcs=copy.deepcopy(w['npcs']),
      locations=[dict(loc,**w['properties'][loc['id']],owned=own(w,loc['id']),locked=loc['district']>w['district'],actions=options(w,loc['id'])) for loc in LOCATIONS],
      event=e,history=w['history'][-60:],dead=w['dead'],tasks=w['tasks'],director=w['director'],last_result=w['last_result'],daily_cost=HOUSING[p['home']]['rent']+p['security']*10+len(p['crew'])*12,
      income=sum(v['income']*v['condition']/100 for k,v in w['properties'].items() if own(w,k)),security=guard(w))


def validate_proposal(w,raw):
    """Allow bounded offers, not arbitrary model-authorized changes to the world."""
    if not isinstance(raw,dict):raise InvalidAction('Director response must be an object.')
    if raw.get('speaker') not in {n['id'] for n in w['npcs']}:raise InvalidAction('Unknown speaker.')
    catalog={'courier':dict(reward=55,respect=2,heat=3,minutes=45),'mediation':dict(reward=35,respect=3,heat=0,minutes=60),'collection':dict(reward=70,respect=1,heat=5,minutes=90)}
    if raw.get('operation') not in catalog:raise InvalidAction('Unsupported operation.')
    for k in ['title','body','outcome']:
        if not isinstance(raw.get(k),str) or not 3<=len(raw[k])<=(70 if k=='title' else 600):raise InvalidAction('Invalid director text.')
    fx=catalog[raw['operation']]
    return dict(id=uuid.uuid4().hex,title=raw['title'],body=raw['body'],outcome=raw['outcome'],speaker=raw['speaker'],kind='proposal',source='local-ai',effect=fx,choices=[dict(id='accept',label='Accept the arrangement',detail=f"${fx['reward']} · {fx['minutes']} minutes · +{fx['respect']} respect · +{fx['heat']} heat"),dict(id='decline',label='Decline the arrangement',detail='No cost or time.')])
