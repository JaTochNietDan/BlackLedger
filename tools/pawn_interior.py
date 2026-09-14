"""Original Ackerman pawnshop, Blender metres with Z up; no external assets."""
import math
import bpy


def build(box,cylinder,material):
    plaster=material('Ackerman aged plaster',(.56,.50,.38));oak=material('Ackerman dark oak',(.19,.105,.045))
    trim=material('Ackerman polished oak',(.34,.21,.095));floor=material('Ackerman floorboards',(.33,.23,.13))
    brass=material('Ackerman brass',(.53,.38,.13),.65);silver=material('Ackerman silver',(.61,.64,.60),.8)
    paper=material('Ackerman pledge paper',(.82,.75,.58));black=material('Ackerman black enamel',(.03,.04,.032))
    velvet=material('Ackerman faded velvet',(.16,.22,.17));cloth=material('Ackerman radio cloth',(.37,.32,.20))
    opal=material('Ackerman opal lamp',(.90,.78,.54),0,.65)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,mat,bevel=0,parent=None):
        ob=box(name,p,d,mat,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,mat,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,mat,rot,vertices);ob.parent=parent;return ob
    def label(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(black);ob.parent=parent
    b('shop foundation',(0,0,-.09),(8,9,.20),oak)
    for col in range(20):
        for row in range(6):
            b('floorboard',(-3.8+col*.4,-3.75+row*1.5,.012),(.39,1.486,.012),floor if (row+col)%3 else trim)
    b('back wall',(0,4.5,1.9),(8,.16,3.8),plaster,parent=back)
    b('left wall',(-4,0,1.9),(.16,9,3.8),plaster,parent=left)
    for z,h in [(.15,.25),(1.03,.07),(3.63,.15)]:
        b('back rail',(0,4.36,z),(8,.11,h),oak,.015,back)
        b('side rail',(-3.86,0,z),(.11,9,h),oak,.015,left)
    # A wide clear staff aisle behind the valuation counter.
    b('counter body',(0,1.6,.54),(6.4,1.0,1.04),oak,.03)
    b('counter top',(0,1.57,1.115),(6.6,1.14,.11),trim,.03)
    for x in [-2.2,0,2.2]:
        b('counter raised panel',(x,1.077,.55),(1.98,.035,.80),trim,.02)
        b('velvet valuation pad',(x,1.56,1.178),(1.35,.73,.016),velvet,.015)
    # Watchmaker's tray with rings, chains and pocket watches. These are scenery,
    # never a substitute for the authoritative possessions offered by the shop.
    b('jewellery tray',(-2.2,1.55,1.215),(1.03,.56,.065),oak,.02)
    b('tray velvet',(-2.2,1.55,1.254),(.93,.46,.015),velvet)
    for x in [-2.49,-2.20,-1.91]:
        c('pocket watch',(x,1.55,1.281),.095,.035,brass,vertices=32)
        c('watch dial',(x,1.55,1.302),.077,.006,paper,vertices=32)
        b('watch minute hand',(x,1.57,1.307),(.01,.083,.004),black)
        b('watch hour hand',(x+.018,1.55,1.309),(.05,.009,.004),black)
        c('watch crown',(x,1.665,1.282),.024,.035,brass)
    # Small mechanical balance for valuation, its pans hanging clear of the pad.
    b('scale foot',(0,1.65,1.23),(.62,.34,.1),black,.04)
    c('scale upright',(0,1.65,1.63),.035,.72,brass)
    b('scale beam',(0,1.65,1.96),(.95,.04,.045),brass,.01)
    for x in [-.39,.39]:
        c('pan suspension',(x,1.65,1.78),.012,.33,brass,vertices=12)
        c('scale pan',(x,1.65,1.61),.18,.027,brass,vertices=32)
    # Register, pledge tickets and an ink stamp at the clerk's end.
    b('cash register base',(2.2,1.64,1.26),(.88,.64,.15),black,.03)
    b('cash register case',(2.2,1.80,1.48),(.80,.30,.35),brass,.035)
    b('register number window',(2.2,1.64,1.56),(.63,.035,.11),black,.01)
    for row in range(3):
        for col in range(7):c('register key',(1.90+col*.10,1.41+row*.07,1.355),.032,.025,silver,vertices=16)
    b('pledge ticket book',(1.42,1.51,1.206),(.42,.52,.045),paper,.012)
    for j in range(5):b('pledge rule',(1.42,1.32+j*.08,1.231),(.34,.008,.003),black)
    c('stamp grip',(2.91,1.55,1.29),.06,.20,oak)
    b('stamp foot',(2.91,1.55,1.203),(.17,.13,.05),black,.01)
    # Rear shelves: tabletop radios, concertina cases and camera boxes.
    for z in [.55,1.5,2.45]:b('rear stock shelf',(0,4.11,z),(7.35,.58,.09),trim,.018,back)
    for x in [-3.55,-1.2,1.2,3.55]:b('stock shelf upright',(x,4.20,1.49),(.10,.36,2.02),oak,.012,back)
    for x in [-2.8,-1.75,-.65,.45,1.55,2.65]:
        b('radio case',(x,4.12,1.89),(.81,.41,.67),oak,.05,back)
        b('radio speaker cloth',(x,3.90,1.93),(.62,.025,.43),cloth,parent=back)
        for xx in [-.23,-.115,0,.115,.23]:b('radio grille',(x+xx,3.88,1.97),(.025,.028,.36),trim,parent=back)
        for xx in [-.20,.20]:c('radio tuning knob',(x+xx,3.862,1.65),.048,.05,black,(math.pi/2,0,0),16,back)
        b('camera box',(x,4.10,.81),(.59,.38,.41),black,.03,back)
        c('camera lens',(x,3.87,.82),.115,.13,silver,(math.pi/2,0,0),24,back)
        c('camera lens glass',(x,3.795,.82),.083,.025,black,(math.pi/2,0,0),24,back)
    for i,x in enumerate([-2.7,0,2.7]):
        b('clock case',(x,4.25,3.02),(1.05,.35,.85),oak,.08,back)
        c('clock face',(x,4.04,3.05),.32,.025,paper,(math.pi/2,0,0),48,back)
        for k in range(12):
            a=k*math.tau/12
            c('clock hour mark',(x+math.sin(a)*.266,4.02,3.05+math.cos(a)*.266),.012,.008,black,(math.pi/2,0,0),12,back)
        b('clock long hand',(x,4.008,3.16),(.018,.018,.23),black,parent=back)
        hand=b('clock short hand',(x+.07,3.999,3.07),(.17,.018,.025),black,parent=back);hand.rotation_euler.y=(i-1)*.6
    # A side display of framed pictures; no portrait claims to be an NPC.
    for y in [-2.3,-.3,1.7]:
        b('picture frame',(-3.83,y,2.0),(.14,1.55,1.42),trim,.03,left)
        b('picture canvas',(-3.747,y,2.0),(.022,1.35,1.22),velvet,parent=left)
        b('picture horizon',(-3.73,y,1.90),(.012,1.34,.13),plaster,parent=left)
    # Cases wait along the right wall, leaving the public floor clear.
    for y in [-1.8,-.6,.6]:
        b('pledged trunk',(3.45,y,.39),(.73,.90,.74),oak,.045)
        for yy in [y-.29,y+.29]:b('trunk strap',(3.45,yy,.40),(.75,.075,.76),brass,.015)
        b('trunk handle',(3.054,y,.42),(.05,.24,.05),black,.015)
    b('shop name plaque',(0,4.23,3.59),(3.75,.06,.24),paper,.01,back)
    label('shop lettering','ACKERMAN & SON',(0,4.186,3.53),.18)
    for x,y in [(-2,0),(1.9,0)]:
        c('pendant cable',(x,y,3.28),.014,.72,black,vertices=12)
        c('opal pendant',(x,y,2.91),.28,.13,opal,vertices=32)
        c('pendant rim',(x,y,2.835),.29,.032,brass,vertices=32)
