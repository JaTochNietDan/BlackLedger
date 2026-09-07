import sys,unittest,tempfile,json,copy
from pathlib import Path
sys.path.insert(0,str(Path(__file__).resolve().parents[1]))
import engine,server

class WorldTests(unittest.TestCase):
 def setUp(self):self.w=engine.new_world()
 def act(self,kind,target=None,**kw):return engine.action(self.w,dict(kind=kind,target=target or self.w['player']['location'],**kw))
 def test_reading_never_advances_time(self):
  before=copy.deepcopy(self.w)
  for i in range(10):engine.public(self.w)
  self.assertEqual(before,self.w)
 def test_travel_consumes_minutes(self):
  before=self.w['minute'];self.act('travel','bar');self.assertEqual(self.w['player']['location'],'bar');self.assertEqual(self.w['minute']-before,engine.travel_minutes('room','bar'))
 def test_unavailable_action_rejected(self):
  with self.assertRaises(engine.InvalidAction):self.act('acquire','laundry')
 def test_jobs_require_presence(self):
  with self.assertRaises(engine.InvalidAction):self.act('courier','bar')
 def test_early_progression(self):
  self.act('travel','bar')
  for i in range(3):
   self.act('courier')
   if self.w['event']:self.act('choice',event=self.w['event']['id'],choice='accept')
  self.assertGreaterEqual(self.w['player']['respect'],6)
  self.act('travel','laundry');self.act('acquire');self.assertTrue(engine.own(self.w,'laundry'))
  cash=self.w['player']['cash'];self.act('wait');self.assertGreater(self.w['player']['cash'],cash)
 def test_hidden_hit_not_exposed(self):
  self.act('travel','club');self.act('provoke');self.assertEqual(len(self.w['plots']),1)
  public=json.dumps(engine.public(self.w));self.assertNotIn('"plots"',public);self.assertNotIn(self.w['plots'][0]['id'],public);self.assertNotIn('"rng"',public)
 def test_contacts_warn(self):
  self.w['player']['contacts']=2;engine.retaliation(self.w);engine.advance(self.w,151)
  self.assertTrue(self.w['plots'][0]['known']);self.assertTrue(any(r['kind']=='danger' for r in self.w['history']))
 def test_attack_interrupts_time(self):
  self.w['player']['security']=1;engine.retaliation(self.w);engine.advance(self.w,400)
  self.assertEqual(self.w['minute'],720);self.assertEqual(self.w['event']['kind'],'attack')
 def test_absent_target_not_killed_at_home(self):
  self.w['player']['location']='bar';engine.retaliation(self.w);engine.advance(self.w,300)
  self.assertTrue(self.w['player']['alive']);self.assertEqual(self.w['properties']['room']['condition'],55)
 def test_tribute_cancels_actual_threat(self):
  self.w['player']['cash']=300;self.act('travel','club');self.act('provoke');self.act('audience');self.act('choice',event=self.w['event']['id'],choice='tribute');self.assertEqual(self.w['plots'],[])
 def test_wrong_event_id_rejected(self):
  self.act('travel','club');self.act('audience')
  with self.assertRaises(engine.InvalidAction):self.act('choice',event='stale',choice='leave')
 def test_security_changes_survival(self):
  for strength in [0,3]:
   w=engine.new_world(seed=50);w['player']['security']=strength;w['player']['contacts']=2;engine.retaliation(w);engine.advance(w,240)
   engine.action(w,dict(kind='choice',event=w['event']['id'],choice='defend'))
   if strength==0:weak=w['player']['alive']
   else:strong=w['player']['alive']
  self.assertFalse(weak);self.assertTrue(strong)
 def test_ironman_same_city(self):
  id=self.w['id'];self.w['properties']['laundry']['owner']='player:1';engine.die(self.w,'Test death');self.act('new_life')
  self.assertEqual(self.w['id'],id);self.assertEqual(self.w['life'],2);self.assertFalse(engine.own(self.w,'laundry'));self.assertEqual(self.w['player']['cash'],90);self.assertEqual(len(self.w['dead']),1)
 def test_living_cannot_reroll(self):
  with self.assertRaises(engine.InvalidAction):self.act('new_life')
 def test_delegate_completes_once(self):
  self.w['player']['crew']=[dict(id='leo',name='Leo',loyalty=65)];self.act('delegate');before=self.w['player']['cash'];engine.advance(self.w,130);self.assertEqual(self.w['player']['cash'],before+65);engine.advance(self.w,50);self.assertEqual(self.w['player']['cash'],before+65)
 def test_bills_not_wall_clock(self):
  self.w['minute']=1439;engine.advance(self.w,1);self.assertEqual(self.w['player']['cash'],75)
 def test_ai_validation_rejects_arbitrary_patches(self):
  for raw in [{},dict(speaker='unknown'),dict(speaker='mara',operation='kill_player')]:
   with self.assertRaises(engine.InvalidAction):engine.validate_proposal(self.w,raw)
 def test_ai_rewards_are_rule_owned(self):
  e=engine.validate_proposal(self.w,dict(speaker='mara',operation='courier',title='Sealed envelope',body='Would you deliver this envelope?',outcome='You delivered the envelope.',cash=99999))
  self.assertEqual(e['effect']['reward'],55)
 def test_future_outcome_not_in_browser(self):
  self.w['event']=engine.validate_proposal(self.w,dict(speaker='mara',operation='courier',title='Sealed envelope',body='Would you deliver this envelope?',outcome='Secret future line.'))
  self.assertNotIn('Secret future line',json.dumps(engine.public(self.w)))

class SaveTests(unittest.TestCase):
 def setUp(self):self.temp=tempfile.TemporaryDirectory();self.old=server.DB;server.DB=Path(self.temp.name)/'game.sqlite3';server.initialize()
 def tearDown(self):server.DB=self.old;self.temp.cleanup()
 def test_duplicate_commit_is_exactly_once(self):
  p=dict(kind='travel',target='bar',revision=0,request_id='duplicate-test');a=server.transact(p);b=server.transact(p);self.assertEqual(a,b);self.assertEqual(server.state()['revision'],1)
 def test_stale_revision_does_not_write(self):
  server.transact(dict(kind='travel',target='bar',revision=0,request_id='first-test'))
  with self.assertRaises(engine.InvalidAction):server.transact(dict(kind='courier',target='bar',revision=0,request_id='second-test'))
  self.assertEqual(server.state()['player']['cash'],90)
 def test_refresh_pending_incident(self):
  with server.connection() as db:
   w=server.load(db);w['player']['security']=1;engine.retaliation(w);engine.advance(w,240);server.save(db,w)
  before=server.state();server.initialize();self.assertEqual(server.state(),before)
 def test_failed_payment_rolls_back(self):
  before=server.state()
  with self.assertRaises(engine.InvalidAction):server.transact(dict(kind='security',target='room',revision=0,request_id='broke-test'))
  self.assertEqual(server.state(),before)

if __name__=='__main__':unittest.main()
