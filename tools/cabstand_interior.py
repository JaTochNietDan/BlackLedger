"""Original Vance Cab Company dispatch office, Blender metres/Z-up."""
import math
import bpy


def build(box,cylinder,material):
    plaster=material('Vance cream plaster',(.60,.56,.43))
    green=material('Vance enamel green',(.14,.23,.18),.25)
    wood=material('Vance dark oak',(.25,.14,.065))
    floor=material('Vance linoleum',(.28,.31,.24))
    paper=material('Vance ivory paper',(.85,.79,.62))
    black=material('Vance bakelite',(.035,.038,.032))
    brass=material('Vance brass',(.57,.40,.13),.65)
    glass=material('Vance milk glass',(.91,.84,.64),0,.7)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,m,rot,vertices);ob.parent=parent;return ob
    def text(name,body,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=body;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(paper);ob.parent=parent
    b('office foundation',(0,0,-.10),(10,10,.23),floor)
    for x in range(-5,5):
        for y in range(-5,5):
            if (x+y)%2==0:b('worn floor tile',(x+.5,y+.5,.016),(.98,.98,.003),green)
    b('rear plaster',(0,5,1.8),(10,.14,3.6),plaster,parent=back)
    b('side plaster',(-5,0,1.8),(.14,10,3.6),plaster,parent=left)
    b('rear dado',(0,4.89,.60),(10,.08,1.20),green,parent=back)
    b('side dado',(-4.89,0,.60),(.08,10,1.20),green,parent=left)
    b('dispatch sign',(0,4.79,3.04),(5.8,.10,.65),green,.015,back)
    text('company name','VANCE CAB COMPANY',(0,4.71,3.03),.30)
    text('company subtitle','DISPATCH OFFICE',(0,4.71,2.80),.15)
    # Rear dispatch counter leaves a working aisle between it and the wall.
    b('dispatch counter',(0,2.6,.53),(4.4,.90,1.04),wood,.025)
    b('counter brass footrail',(0,1.98,.20),(4.25,.055,.055),brass,.012)
    b('counter worktop',(0,2.6,1.09),(4.6,1.04,.12),green,.02)
    for x in [-1.35,1.25]:
        b('telephone base',(x,2.6,1.19),(.42,.31,.14),black,.035)
        c('telephone dial',(x,2.54,1.273),.105,.018,brass)
        for a in range(10):
            t=a*math.tau/10;c('dial finger hole',(x+math.sin(t)*.074,2.54+math.cos(t)*.074,1.286),.012,.009,black,vertices=12)
        b('telephone receiver',(x,2.69,1.35),(.47,.09,.08),black,.025)
        for side in [-1,1]:c('receiver earpiece',(x+side*.18,2.69,1.31),.07,.10,black)
    b('fare ledger',(-.22,2.45,1.17),(.66,.46,.04),paper)
    b('ledger spine',(-.55,2.45,1.18),(.035,.47,.05),wood)
    b('pencil',(.05,2.42,1.20),(.015,.28,.015),brass)
    # Route board is deliberately schematic decor, not live jobs or cab positions.
    b('route board',(-2.75,4.77,2.06),(3.1,.10,1.13),black,.02,back)
    for x in [-3.8,-3.2,-2.6,-2]:b('route north-south',(x,4.70,2.06),(.018,.015,.78),paper,parent=back)
    for z in [1.76,2.02,2.30]:b('route east-west',(-2.85,4.69,z),(2.4,.018,.018),paper,parent=back)
    text('routes label','BELLWETHER ROUTES',(-2.75,4.68,2.48),.14)
    b('fare pigeonholes',(3.0,4.72,1.78),(2.1,.30,1.62),wood,.02,back)
    for x in [2.05,2.68,3.31,3.94]:b('pigeonhole divider',(x,4.51,1.78),(.045,.28,1.56),green,parent=back)
    for z in [1.05,1.52,1.99,2.46]:b('pigeonhole shelf',(3,4.51,z),(2.05,.30,.045),green,parent=back)
    for x in [2.35,2.98,3.61]:
        for z in [1.10,1.57,2.04]:b('folded fare sheets',(x,4.48,z),(.40,.20,.055),paper,parent=back)
    # Driver lockers and coffee shelf stay at the sides, outside public lanes.
    for y in [-2.9,-2,-1.1]:
        b('driver locker',(-4.45,y,1.08),(.70,.84,2.12),green,.025)
        b('locker door',(-4.08,y,1.08),(.035,.75,1.98),wood,.012)
        b('locker handle',(-4.04,y-.24,1.02),(.05,.03,.20),brass)
        for z in [1.65,1.72,1.79]:b('locker vent',(-4.05,y,z),(.014,.49,.014),black)
    b('drivers bench cushion',(3.85,.1,.53),(.76,2.75,.12),wood,.035)
    b('drivers bench back',(4.24,.1,.87),(.12,2.8,.82),green,.025)
    for y in [-1.05,1.25]:
        for x in [3.57,4.13]:b('bench leg',(x,y,.25),(.09,.09,.50),wood)
    b('coffee shelf',(3.9,-2.4,.82),(1.5,.66,.12),wood,.02)
    for x in [3.25,4.55]:b('shelf leg',(x,-2.4,.40),(.08,.50,.80),green)
    c('coffee urn',(3.65,-2.4,1.16),.17,.54,brass)
    c('urn lid',(3.65,-2.4,1.44),.19,.04,black)
    for x in [4.0,4.25]:c('coffee cup',(x,-2.45,.96),.06,.13,paper)
    c('office clock',(.25,4.75,2.15),.30,.12,wood,(math.pi/2,0,0),48,back)
    c('clock face',(.25,4.67,2.15),.26,.018,paper,(math.pi/2,0,0),48,back)
    b('clock hour',(.20,4.65,2.20),(.12,.015,.023),black,parent=back)
    b('clock minute',(.25,4.64,2.25),(.018,.015,.22),black,parent=back)
    for x in [-2,2]:
        c('pendant stem',(x,0,3.31),.022,.55,black)
        c('pendant shade',(x,0,3.02),.37,.10,green)
        c('pendant diffuser',(x,0,2.95),.28,.08,glass)
