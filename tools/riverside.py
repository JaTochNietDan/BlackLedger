"""Original Riverside Courts: four residential wings and their communal entrance."""
import math
import bpy


def lettering(name,words,p,size,mat,parent=None):
    cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
    o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(mat);o.parent=parent


def exterior(box,cylinder,material):
    brick=material('Riverside red brick',(.39,.20,.13));stone=material('Riverside pale stone',(.60,.55,.43))
    roof=material('Riverside roof tar',(.105,.12,.115));iron=material('Riverside painted iron',(.04,.065,.055),.35)
    glass=material('Riverside window glass',(.15,.24,.25),.25);lit=material('Riverside occupied glass',(.64,.40,.16),0,.35)
    grass=material('Riverside courtyard green',(.18,.25,.14));wood=material('Riverside courtyard wood',(.29,.18,.08))
    box('courtyard plinth',(0,0,.1),(16.8,16.8,.2),stone,.04)
    # Four nine-storey wings: eight flats on eight floors, four on the top floor.
    # The map compresses apartment depth, but entrances and floor heights fit the cast.
    for wing,(x,y) in enumerate([(-5,-5),(5,-5),(-5,5),(5,5)]):
        box('residential wing',(x,y,13.55),(5.6,5.6,26.9),brick,.055)
        box('wing foundation',(x,y,.38),(5.8,5.8,.6),stone,.03)
        box('wing cornice',(x,y,26.9),(5.95,5.95,.28),stone,.035)
        box('wing roof',(x,y,27.12),(5.85,5.85,.22),roof,.025)
        for floor in range(9):
            z=1.75+floor*3
            for side in [-1,1]:
                for offset in [-1.5,0,1.5]:
                    pane=lit if (floor*7+wing*3+int(offset*2))%5==0 else glass
                    box('front wing window',(x+offset,y+side*2.815,z),( .85,.045,1.35),pane)
                    box('front stone sill',(x+offset,y+side*2.85,z-.75),(1.03,.16,.12),stone)
                    box('side wing window',(x+side*2.815,y+offset,z),(.045,.85,1.35),pane)
            box('wing floor course',(x,y,z+1.05),(5.66,5.66,.035),stone)
        box('wing entry',(x,y+2.84,1.25),(1.3,.08,2.30),iron,.015)
        box('entry canopy',(x,y+3.04,2.58),(1.75,.65,.16),stone)
        lettering('wing letter',chr(65+wing),(x,y+2.901,2.90),.35,stone)
    box('central garden',(0,0,.24),(3,8,.25),grass,.12)
    for y in [-3,0,3]:
        box('garden planter',(0,y,.5),(1.5,1.4,.4),stone,.05)
        for x in [-.4,.1,.5]:
            cylinder('courtyard shrub',(x,y,.93),.40,.7,grass,vertices=12)
    for x in [-2.15,2.15]:
        box('courtyard bench',(x,0,.65),(.6,2.4,.14),wood,.03)
        for y in [-.9,.9]:box('bench support',(x,y,.37),(.4,.12,.5),iron)
    box('entry directory',(0,7.6,1.1),(2.7,.20,1.7),stone,.04)
    lettering('address','RIVERSIDE COURTS',(0,7.715,1.56),.20,iron)
    lettering('wing register','A  1-68     B  69-136\nC  137-204  D  205-272',(0,7.72,1.16),.15,iron)
    for o in bpy.context.scene.objects:
        if o.type=='FONT':o.rotation_euler.z=math.pi


def lobby(box,cylinder,material):
    cream=material('Riverside lobby cream',(.65,.61,.49));green=material('Riverside lobby green',(.23,.33,.27))
    tile=material('Riverside terrazzo',(.46,.47,.39));wood=material('Riverside oak',(.28,.16,.07));brass=material('Riverside brass',(.49,.35,.14),.5)
    ink=material('Riverside lettering',(.04,.075,.065));glass=material('Riverside frosted glass',(.39,.51,.46),.1)
    light=material('Riverside opal shade',(.85,.78,.57),0,.6)
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    b('lobby foundation',(0,0,-.09),(10,11,.20),green)
    for x in range(20):
        for y in range(22):b('terrazzo tile',(-4.75+x*.5,-5.25+y*.5,.011),(.488,.488,.013),green if x in [0,19] or y in [0,21] else tile if (x+y)%2 else cream)
    b('rear wall',(0,5.5,2),(10,.16,4),cream,parent=back);b('side wall',(-5,0,2),(.16,11,4),cream,parent=left)
    b('rear dado',(0,5.405,.68),(10,.025,1.32),green,parent=back);b('side dado',(-4.905,0,.68),(.025,11,1.32),green,parent=left)
    for z in [.12,1.35,3.82]:
        b('rear trim',(0,5.37,z),(10,.07,.08),wood,parent=back);b('side trim',(-4.87,0,z),(.07,11,.08),wood,parent=left)
    for i,x in enumerate([-3.6,-1.2,1.2,3.6]):
        b('wing door frame',(x,5.30,1.40),(1.97,.20,2.8),wood,.025,back)
        b('wing door',(x,5.17,1.34),(1.75,.07,2.59),green,.02,back)
        b('wing door glazing',(x,5.125,1.95),(1.25,.02,.75),glass,parent=back)
        b('wing push plate',(x+.62,5.09,1.18),(.15,.035,.40),brass,.01,back)
        lettering('wing directory',f'{chr(65+i)}   {i*68+1}-{(i+1)*68}',(x,5.075,2.55),.17,cream,back)
    lettering('lobby title','RIVERSIDE COURTS',(0,5.27,3.27),.38,ink,back)
    b('waiting bench base',(-4.20,-.8,.28),(.71,4.3,.48),wood,.025)
    b('waiting bench seat',(-4.20,-.8,.61),(.75,4.3,.16),wood,.035)
    b('waiting bench back',(-4.63,-.8,.96),(.16,4.4,.76),wood,.02)
    for x,y in [(-2,1.9),(2.6,1.4),(0,-2.1)]:
        cylinder('lamp cable',(x,y,3.43),.012,.62,ink,vertices=12)
        cylinder('opal ceiling lamp',(x,y,3.08),.32,.11,light,vertices=32)
        cylinder('lamp rim',(x,y,3.01),.34,.035,brass,vertices=32)
    b('resident noticeboard',(4.44,0,1.4),(.16,2.4,2.0),wood,.03)
    for y in [-.9,.9]:b('noticeboard leg',(4.44,y,.3),(.14,.14,.6),wood,.015)
    for y in [-.75,0,.75]:
        b('tenant notice',(4.343,y,1.5),(.02,.54,.8),cream)
        for z in [1.65,1.5,1.35]:b('notice line',(4.325,y,z),(.015,.38,.015),ink)
