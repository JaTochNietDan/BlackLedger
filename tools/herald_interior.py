"""Original 1950s newspaper city desk, authored in metres, Blender Z up."""
import math
import bpy


def build(box,cylinder,material):
    plaster=material('Herald tobacco plaster',(.55,.51,.40));wood=material('Herald worn oak',(.28,.17,.085))
    trim=material('Herald dark trim',(.15,.105,.065));floor=material('Herald linoleum',(.24,.29,.25))
    paper=material('Herald newsprint',(.82,.78,.63));ink=material('Herald enamel',(.035,.048,.043))
    steel=material('Herald filing steel',(.30,.36,.32),.25);brass=material('Herald brass',(.48,.35,.14),.6)
    glass=material('Herald green lamp',(.07,.22,.12));opal=material('Herald opal light',(.88,.81,.61),0,.5)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),vertices=20,parent=None):
        ob=cylinder(name,p,r,d,m,rot,vertices);ob.parent=parent;return ob
    def label(name,text,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(ink);ob.parent=parent
    b('newsroom foundation',(0,0,-.09),(10,10,.2),trim)
    for x in range(10):
        for y in range(10):b('linoleum tile',(x-4.5,y-4.5,.015),(.988,.988,.016),floor if (x+y)%2 else steel)
    b('back plaster',(0,5,1.9),(10,.16,3.8),plaster,parent=back)
    b('left plaster',(-5,0,1.9),(.16,10,3.8),plaster,parent=left)
    for z,h in [(.13,.23),(1.13,.07),(3.65,.18)]:
        b('back moulding',(0,4.87,z),(10,.12,h),trim,.012,back)
        b('side moulding',(-4.87,0,z),(.12,10,h),trim,.012,left)
    # Four independent writing desks, with a broad central circulation aisle.
    for row,y in enumerate([2,-.4]):
        for x in [-2.55,2.55]:
            # Desk chairs face the platen, with a real cushion for cast seating.
            cy=y+.82
            b('desk chair cushion',(x,cy,.555),(.55,.53,.07),trim,.03)
            for dx in [-.22,.22]:
                for dy in [-.21,.21]:b('chair leg',(x+dx,cy+dy,.275),(.05,.05,.53),wood,.008)
            for dx in [-.23,.23]:b('chair back upright',(x+dx,cy+.23,.82),(.05,.05,.70),wood,.008)
            b('chair back',(x,cy+.23,1.03),(.53,.06,.28),wood,.025)
            b('copy desk top',(x,y,.80),(2.25,1.04,.10),wood,.025)
            for dx in [-.85,.85]:
                b('desk pedestal',(x+dx,y,.40),(.43,.82,.72),trim,.018)
                for z in [.25,.50,.70]:
                    b('desk drawer',(x+dx,y+.425,z),(.36,.04,.15),wood,.008)
                    b('drawer pull',(x+dx,y+.455,z),(.15,.045,.018),brass,.006)
            b('typewriter foot',(x-.18,y,.89),(.60,.47,.07),ink,.025)
            b('typewriter carriage',(x-.18,y-.15,1.06),(.70,.16,.17),steel,.02)
            c('typewriter platen',(x-.18,y-.13,1.15),.046,.58,ink,(0,math.pi/2,0))
            b('paper in platen',(x-.18,y-.16,1.31),(.43,.014,.29),paper)
            for line in range(7):b('typed line',(x-.18,y-.169,1.23+line*.026),(.32-(line%3)*.04,.004,.006),ink)
            for r in range(3):
                for k in range(10):c('typewriter key',(x-.43+k*.055,y+.04+r*.06,.951+r*.008),.019,.018,paper,vertices=12)
            b('space bar',(x-.18,y+.24,.962),(.29,.035,.025),ink,.006)
            b('copy stack',(x+.65,y-.02,.90),(.42,.57,.09),paper,.006)
            for j in range(5):b('copy text',(x+.65,y-.21+j*.08,.948),(.32,.01,.003),ink)
            c('desk lamp base',(x-.86,y-.28,.88),.13,.04,brass)
            c('desk lamp stem',(x-.86,y-.28,1.08),.018,.38,brass)
            b('green desk shade',(x-.86,y-.28,1.30),(.40,.21,.13),glass,.06)
    # Filing wall and pigeonholes remain behind the staff working aisle.
    for x in [-3.7,-2.7,2.7,3.7]:
        b('filing cabinet',(x,4.43,1.0),(.84,.75,1.94),steel,.035,back)
        for z in [.30,.76,1.22,1.68]:
            b('file drawer',(x,4.035,z),(.73,.035,.39),floor,.018,back)
            b('file drawer handle',(x,3.99,z+.04),(.24,.065,.035),brass,.008,back)
            b('file label',(x,4.007,z-.075),(.25,.015,.055),paper,parent=back)
    b('masthead board',(0,4.82,3.03),(5.1,.08,.70),paper,.02,back)
    label('masthead','THE BELLWETHER HERALD',(0,4.764,3.09),.245)
    label('city desk sign','CITY DESK  /  EDITORIAL',(0,4.76,2.84),.15)
    for x in [-.9,0,.9]:
        b('copy pigeonhole',(x,4.42,1.36),(.82,.7,1.4),wood,.012,back)
        for z in [.90,1.25,1.60,1.95]:
            b('pigeonhole opening',(x,4.054,z),(.69,.02,.25),ink,parent=back)
            b('file bundle',(x,4.01,z-.055),(.52,.16,.09),paper,parent=back)
    # Wire machine, paper roll and receiving baskets at the left wall.
    b('wire service bench',(-4.36,-2.2,.78),(.85,2.0,.10),wood,.02)
    for y in [-2.9,-1.5]:b('wire bench leg',(-4.36,y,.39),(.70,.15,.72),trim)
    b('teleprinter housing',(-4.36,-1.95,1.10),(.65,.65,.53),ink,.04)
    c('wire paper roll',(-4.36,-1.6,1.48),.13,.46,paper,(0,math.pi/2,0))
    b('wire paper feed',(-4.36,-1.83,1.39),(.42,.36,.015),paper)
    b('wire copy basket',(-4.36,-2.8,.97),(.60,.40,.25),steel,.02)
    for z in [.88,.95,1.02,1.09]:b('basket slot',(-3.99,-2.8,z),(.02,.31,.025),ink)
    # Editorial pinboard is decorative, with no invented live headlines.
    b('notice frame',(-4.85,.45,2.2),(.10,2.4,1.45),wood,.02,left)
    b('notice cork',(-4.785,.45,2.2),(.022,2.22,1.27),trim,parent=left)
    for y in [-.28,.42,1.1]:
        for z in [1.95,2.5]:
            b('pinned clipping',(-4.765,y,z),(.012,.48,.43),paper,parent=left)
            c('pin',(-4.748,y,z+.17),.019,.012,brass,(0,math.pi/2,0),12,left)
    for x in [-2.5,2.5]:
        c('light cable',(x,0,3.45),.012,.6,ink,vertices=12)
        c('newsroom pendant',(x,0,3.11),.33,.16,opal,vertices=32)
        c('pendant rim',(x,0,3.015),.35,.04,steel,vertices=32)
