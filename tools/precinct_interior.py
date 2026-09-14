"""Original period police booking room; Blender metres, Z up. No external assets."""
import math
import bpy


def build(box,cylinder,material):
    plaster=material('Ward Station aged plaster',(.61,.61,.51))
    green=material('Ward Station institutional green',(.20,.29,.24))
    tile=material('Ward Station grey terrazzo',(.40,.44,.42))
    border=material('Ward Station dark floor border',(.11,.16,.14))
    oak=material('Ward Station dark oak',(.19,.105,.055))
    trim=material('Ward Station polished oak',(.34,.23,.12))
    steel=material('Ward Station painted steel',(.18,.23,.22),.2)
    metal=material('Ward Station nickel',(.47,.53,.51),.7)
    brass=material('Ward Station brass',(.48,.34,.12),.6)
    black=material('Ward Station black bakelite',(.024,.031,.028))
    paper=material('Ward Station report paper',(.79,.76,.62))
    cork=material('Ward Station cork',(.36,.27,.16))
    light=material('Ward Station opal shade',(.82,.83,.70),0,.5)
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def c(name,p,r,d,m,rotation=(0,0,0),vertices=24,parent=None):
        o=cylinder(name,p,r,d,m,rotation,vertices);o.parent=parent;return o
    def label(name,words,p,size,parent=None):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(black);o.parent=parent
    b('station foundation',(0,0,-.09),(10,11,.20),border)
    for x in range(20):
        for y in range(22):
            material=border if x in [0,19] or y in [0,21] else tile if (x+y)%2 else plaster
            b('terrazzo tile',(-4.75+x*.5,-5.25+y*.5,.011),(.488,.488,.013),material)
    b('rear wall',(0,5.5,2),(10,.16,4),plaster,parent=back)
    b('side wall',(-5,0,2),(.16,11,4),plaster,parent=left)
    b('rear painted dado',(0,5.405,.67),(10,.025,1.31),green,parent=back)
    b('side painted dado',(-4.905,0,.67),(.025,11,1.31),green,parent=left)
    for z,h in [(.09,.15),(1.34,.06),(3.8,.12)]:
        b('rear rail',(0,5.37,z),(10,.08,h),trim,.009,back)
        b('side rail',(-4.87,0,z),(.08,11,h),trim,.009,left)
    # The sergeant's counter has an open rear working aisle.
    b('booking counter',(-2.1,2,.54),(4.8,1.02,1.06),oak,.03)
    b('booking counter ledge',(-2.1,1.98,1.12),(4.98,1.18,.12),trim,.035)
    for x in [-3.8,-2.1,-.4]:b('counter panel',(x,1.47,.57),(1.43,.045,.81),trim,.018)
    b('register lower cover',(-2.4,1.85,1.195),(1.08,.7,.035),black,.015)
    b('open register',(-2.4,1.85,1.222),(1.01,.66,.024),paper,.008)
    for x in [-2.66,-2.15]:
        for i in range(9):b('register ruled line',(x,1.61+i*.057,1.236),(.41,.006,.002),green)
    b('counter nameplate',(-1.25,1.51,1.29),(.9,.09,.20),brass,.015)
    label('sergeant lettering','DESK SERGEANT',(-1.25,1.457,1.25),.075)
    c('desk bell foot',(-.25,1.78,1.218),.13,.048,black)
    c('desk bell dome',(-.25,1.78,1.271),.105,.07,brass,vertices=32)
    c('bell button',(-.25,1.78,1.321),.032,.03,brass)
    # Heavy black telephone with receiver, dial and coiled lead.
    b('telephone base',(-3.73,1.92,1.25),(.49,.38,.14),black,.045)
    c('rotary dial',(-3.73,1.84,1.33),.11,.015,metal,vertices=32)
    for i in range(10):
        a=i*math.tau/10
        c('dial finger hole',(-3.73+math.sin(a)*.078,1.84+math.cos(a)*.078,1.340),.018,.006,black,vertices=12)
    b('telephone receiver',(-3.73,2.07,1.43),(.53,.11,.09),black,.04)
    for x in [-3.95,-3.51]:c('receiver cup',(x,2.07,1.398),.078,.11,black)
    for i in range(16):
        a=i*math.tau/4
        c('telephone cord',(-4.08+.025*math.cos(a),2.03+.025*math.sin(a),1.22+i*.011),.012,.022,black,vertices=8)
    # Report desk and a mechanical typewriter; staff have a distinct rear station.
    b('report desk top',(2.8,1.4,.97),(2.3,1.1,.11),trim,.025)
    for x in [1.85,3.75]:
        b('report desk pedestal',(x,1.4,.46),(.34,.92,.87),oak,.02)
        for z in [.24,.51,.78]:
            b('report desk drawer',(x,.93,z),(.29,.03,.22),trim,.009)
            b('report drawer pull',(x,.896,z),(.15,.025,.02),brass,.005)
    b('typewriter base',(2.45,1.33,1.078),(.66,.54,.15),black,.035)
    b('typewriter body',(2.45,1.50,1.21),(.62,.22,.20),steel,.025)
    c('typewriter platen',(2.45,1.60,1.35),.055,.70,black,(0,math.pi/2,0),32)
    b('typewriter paper',(2.45,1.62,1.54),(.48,.012,.40),paper)
    for row in range(4):
        for col in range(10):c('typewriter key',(2.18+col*.06,1.1+row*.056,1.163),.020,.012,metal,vertices=12)
    b('typewriter space bar',(2.45,1.04,1.161),(.38,.032,.015),metal,.005)
    for x in [3.1,3.52]:b('report folder',(x,1.36,1.057),(.34,.48,.05),paper)
    # Lockers and filing drawers belong to the visible office, not a fake jail roster.
    for y in [-.75,.15]:
        b('steel locker',(4.48,y,1.05),(.70,.82,2.07),steel,.025)
        b('locker door',(4.105,y,1.06),(.035,.72,1.94),green,.01)
        for z in [1.55,1.64,1.73]:b('locker vent',(4.078,y,z),(.013,.43,.027),black)
        b('locker handle',(4.04,y-.25,1.03),(.07,.04,.18),metal,.015)
    # Actual holding area is beyond this secured door; no invented prisoner models.
    b('holding door frame',(3.7,5.30,1.47),(1.83,.20,2.94),steel,.025,back)
    b('holding door',(3.7,5.17,1.43),(1.58,.08,2.77),green,.02,back)
    b('holding observation window',(3.7,5.113,2.12),(.85,.018,.60),black,parent=back)
    for x in [3.37,3.59,3.81,4.03]:b('observation bar',(x,5.08,2.12),(.026,.04,.62),metal,parent=back)
    b('holding lock plate',(4.22,5.105,1.23),(.14,.035,.28),steel,parent=back)
    b('holding lever',(4.13,5.05,1.27),(.27,.06,.035),metal,.01,back)
    b('holding sign',(3.7,5.105,2.66),(1.15,.025,.22),paper,parent=back)
    label('holding lettering','HOLDING',(3.7,5.085,2.61),.13,back)
    b('reports board frame',(-2.15,5.28,2.03),(4.55,.10,1.62),oak,.02,back)
    b('reports board cork',(-2.15,5.21,2.03),(4.37,.03,1.45),cork,parent=back)
    for x in [-3.8,-2.95,-2.1,-1.25,-.4]:
        b('report sheet',(x,5.18,2.02),(.62,.012,.99),paper,parent=back)
        for i in range(9):b('report text line',(x,5.164,2.30-i*.075),(.48,.009,.012),green,parent=back)
    b('station title plaque',(-.5,5.21,3.31),(6.5,.1,.47),paper,.02,back)
    label('station title','WARD STREET STATION',(-.5,5.147,3.20),.28,back)
    # Three-person waiting bench; blank walls leave the public aisle readable.
    b('waiting bench base',(-4.20,-.8,.28),(.71,4.3,.48),oak,.025)
    b('waiting bench seat',(-4.20,-.8,.61),(.75,4.3,.16),trim,.035)
    b('waiting bench back',(-4.63,-.8,.96),(.16,4.4,.76),oak,.02)
    for y in [-2.2,-.8,.6]:
        b('bench back panel',(-4.527,y,.98),(.035,1.22,.52),trim,.016)
    for x,y in [(-2,1.9),(2.6,1.4),(0,-2.1)]:
        c('lamp cable',(x,y,3.43),.012,.62,black,vertices=12)
        c('opal ceiling lamp',(x,y,3.08),.32,.11,light,vertices=32)
        c('lamp rim',(x,y,3.01),.34,.035,metal,vertices=32)
