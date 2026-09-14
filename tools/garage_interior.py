"""Original 1950s motor workshop. Blender metres, Z up; no external assets."""
import math
import bpy


def build(box,cylinder,material):
    concrete=material('Garage worn concrete',(.30,.31,.27))
    brick=material('Garage dusty brick',(.38,.23,.14))
    grout=material('Garage mortar',(.46,.42,.32))
    green=material('Garage enamel green',(.10,.22,.18),.15)
    red=material('Garage tool red',(.40,.09,.055),.15)
    iron=material('Garage iron',(.055,.065,.059),.45)
    metal=material('Garage machined steel',(.48,.52,.48),.8)
    rubber=material('Garage tyre rubber',(.022,.025,.021))
    wood=material('Garage worn worktop',(.38,.24,.11))
    paper=material('Garage aged sign',(.79,.72,.54))
    lamp=material('Garage opal tubes',(.77,.83,.73),0,.5)
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def c(name,p,r,d,m,rotation=(0,0,0),vertices=24,parent=None):
        o=cylinder(name,p,r,d,m,rotation,vertices);o.parent=parent;return o
    def label(name,text,p,size,parent=None):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(iron);o.parent=parent
    b('workshop foundation',(0,0,-.10),(12,11,.21),iron)
    for x in range(6):
        for y in range(6):
            lo=-5.5+y*2;hi=min(lo+2,5.5)
            if hi>lo:b('concrete slab',(-5+x*2,(lo+hi)/2,.011),(1.985,hi-lo-.015,.012),concrete)
    b('left masonry',(-6,0,2.1),(.18,11,4.2),grout,parent=left)
    b('rear masonry',(0,5.5,2.1),(12,.18,4.2),grout,parent=back)
    for row in range(20):
        for col in range(25):
            x=-6.2+col*.5+(row%2)*.25
            lo=max(-5.94,x);hi=min(5.94,x+.475)
            if hi>lo:b('rear brick',((lo+hi)/2,5.395,.105+row*.2),(hi-lo,.035,.183),brick,parent=back)
        for col in range(23):
            y=-5.7+col*.5+(row%2)*.25
            lo=max(-5.44,y);hi=min(5.44,y+.475)
            if hi>lo:b('side brick',(-5.895,(lo+hi)/2,.105+row*.2),(.035,hi-lo,.183),brick,parent=left)
    for x in [-5.65,0,5.65]:b('rear steel column',(x,5.23,2.12),(.16,.19,4.22),green,parent=back)
    b('rear steel beam',(0,5.23,4.03),(11.6,.24,.23),green,parent=back)
    # Empty inspection lift; no decorative car impersonates a public vehicle.
    for x in [-4.55,-1.65]:
        b('service bay stripe',(x,-.65,.02),(.09,6.4,.005),paper)
    for x in [-4.0,-2.2]:
        b('lift runway',(x,.25,.23),(.43,4.55,.21),iron,.025)
        for y in [-1.5,1.8]:b('lift support',(x,y,.10),(.24,.30,.18),metal)
        for y in [-1.65,-.8,.05,.9,1.75]:b('runway grip',(x,y,.344),(.4,.03,.012),metal)
    c('hydraulic ram',(-3.1,.35,.14),.26,.27,metal)
    b('lift crossmember',(-3.1,.35,.19),(2.4,.32,.23),green,.025)
    # Bench and hung tools behind the lift.
    b('workbench cabinet',(-3.15,4.66,.49),(4.4,.98,.94),green,.025)
    b('workbench timber top',(-3.15,4.66,1.03),(4.6,1.12,.13),wood,.025)
    for x in [-4.76,-3.69,-2.62,-1.55]:
        for z in [.22,.5,.78]:
            b('tool drawer',(x,4.15,z),(.99,.045,.22),green,.008)
            b('drawer handle',(x,4.1,z),(.41,.04,.035),metal,.015)
    b('tool board',(-3.2,5.20,2.03),(4.7,.075,1.35),wood,parent=back)
    for i in range(12):
        x=-5.25+i*.37
        b('hanging spanner',(x,5.12,2.05),(.055,.025,.45),metal,.012,back)
        c('spanner ring',(x,5.12,2.32),.07,.025,metal,(math.pi/2,0,0),16,back)
        c('spanner socket',(x,5.096,2.32),.038,.006,iron,(math.pi/2,0,0),16,back)
    b('vice base',(-4.55,4.60,1.15),(.44,.43,.11),iron,.02)
    b('vice jaws',(-4.55,4.60,1.34),(.46,.26,.22),metal,.015)
    c('vice screw',(-4.55,4.28,1.25),.035,.48,metal,(math.pi/2,0,0),16)
    # Engine on its stand in the right bay; cylinder head and manifold ribs.
    for x in [1.55,2.65]:b('engine stand foot',(x,2.55,.11),(.12,1.28,.12),green)
    b('engine stand crossbar',(2.1,2.55,.18),(1.35,.14,.14),green)
    c('engine stand upright',(2.1,2.85,.64),.09,.95,green)
    b('engine block',(2.1,2.55,1.10),(.8,1.17,.65),iron,.06)
    b('cylinder head',(2.1,2.55,1.47),(.89,1.23,.16),metal,.03)
    for y in [2.08,2.32,2.56,2.8,3.04]:
        c('spark plug',(2.1,y,1.61),.042,.14,paper,vertices=12)
        c('exhaust branch',(1.62,y,1.2),.045,.42,metal,(0,math.pi/2,0),16)
    c('engine pulley',(2.1,1.91,1.1),.23,.09,metal,(math.pi/2,0,0),32)
    c('pulley hub',(2.1,1.845,1.1),.075,.06,iron,(math.pi/2,0,0),16)
    # Tyre rack, barrels with pumps and period air compressor.
    for x in [.1,4.7]:
        b('rack upright',(x,4.98,1.24),(.09,.09,2.45),green)
    for z in [.35,1.65]:
        for y in [4.48,5.1]:b('rack rail',(2.4,y,z),(4.7,.08,.08),green)
        for x in [.65,1.55,2.45,3.35,4.25]:
            c('spare tyre',(x,4.76,z+.38),.37,.25,rubber,(math.pi/2,0,0),36)
            c('tyre hub opening',(x,4.617,z+.38),.21,.012,iron,(math.pi/2,0,0),32)
            for a in range(20):
                t=a*math.tau/20
                o=b('tyre tread',(x+math.sin(t)*.366,4.76,z+.38+math.cos(t)*.366),(.06,.27,.03),iron,.005);o.rotation_euler.y=t
    for x in [4.65,5.4]:
        c('oil drum',(x,2.1,.47),.30,.9,red,vertices=32)
        for z in [.15,.43,.78]:c('drum reinforcing hoop',(x,2.1,z),.31,.035,iron,vertices=32)
        c('drum pump',(x,2.1,1.12),.027,.46,metal,vertices=12)
        b('pump handle',(x+.08,2.1,1.38),(.25,.04,.04),iron,.01)
    c('compressor reservoir',(5,.1,.46),.30,1.1,green,(math.pi/2,0,0),32)
    for y in [-.32,.5]:
        for x in [4.73,5.27]:c('compressor wheel',(x,y,.16),.13,.07,rubber,(0,math.pi/2,0),20)
    b('compressor motor',(5,.15,.85),(.48,.51,.24),iron,.035)
    # Foreground service desk and paperwork; leaves the middle entrance clear.
    b('service desk',(4.45,-3.05,.58),(2.15,1.05,1.12),green,.025)
    b('desk top',(4.45,-3.05,1.17),(2.29,1.17,.09),wood,.02)
    b('job ledger',(4.1,-3.05,1.235),(.43,.52,.035),paper,.01)
    for y in [-3.2,-3.12,-3.04,-2.96]:b('ledger line',(4.1,y,1.255),(.32,.003,.002),iron)
    b('counter sign',(4.93,-3.24,1.39),(.70,.035,.30),paper,.01)
    label('service sign','SERVICE',(4.93,-3.265,1.34),.105)
    b('garage sign',(0,5.23,3.4),(6.1,.1,.55),paper,.02,back)
    label('garage lettering','RUSSO MOTOR WORKS',(0,5.16,3.23),.32,back)
    for x in [-3,2.5]:
        for dx in [-.82,.82]:c('light suspension cable',(x+dx,.5,3.86),.012,.56,iron,vertices=8)
        b('fluorescent housing',(x,.5,3.55),(2.2,.29,.10),green,.02)
        for y in [.41,.59]:c('fluorescent tube',(x,y,3.47),.034,2.06,lamp,(0,math.pi/2,0),16)
