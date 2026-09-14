"""Original 1950s neighbourhood dining room. Blender metres, Z up."""
import math
import bpy


def build(box,cylinder,material):
    cream=material('Vittoria warm plaster',(.66,.57,.40))
    wood=material('Vittoria dark walnut',(.15,.075,.035))
    edge=material('Vittoria walnut highlights',(.29,.16,.065))
    red=material('Vittoria oxblood upholstery',(.30,.035,.026))
    brass=material('Vittoria aged brass',(.50,.34,.12),.65)
    linen=material('Vittoria ivory linen',(.85,.79,.63))
    china=material('Vittoria glazed china',(.84,.86,.75))
    steel=material('Vittoria silverware',(.56,.61,.59),.8)
    green=material('Vittoria bottle glass',(.055,.13,.055),.12)
    dark=material('Vittoria charcoal',(.032,.027,.02))
    tile=material('Vittoria terrazzo',(.48,.43,.32))
    lamp=material('Vittoria frosted lamp',(.90,.72,.41),0,.6)
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def c(name,p,r,d,m,rotation=(0,0,0),vertices=24,parent=None):
        o=cylinder(name,p,r,d,m,rotation,vertices);o.parent=parent;return o
    def text(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(brass);o.parent=parent
    b('dining foundation',(0,0,-.08),(10,11,.18),dark)
    for ix in range(20):
        for iy in range(22):
            b('terrazzo tile',(-4.75+ix*.5,-5.25+iy*.5,.011),(.487,.487,.013),tile if (ix+iy)%2 else cream)
    b('left plaster',(-5,0,1.8),(.15,11,3.6),cream,parent=left)
    b('rear plaster',(0,5.5,1.8),(10,.15,3.6),cream,parent=back)
    for z,h in [(.10,.18),(1.18,.10),(3.35,.12),(3.47,.07)]:
        b('rear moulding',(0,5.36,z),(10,.15,h),wood,.012,back)
        b('left moulding',(-4.86,0,z),(.15,11,h),wood,.012,left)
    b('rear wainscot',(0,5.40,.62),(10,.09,1.05),wood,parent=back)
    b('left wainscot',(-4.90,0,.62),(.09,11,1.05),wood,parent=left)
    for x in [-4.5,-3.5,-2.5,-1.5,-.5,.5,1.5,2.5,3.5,4.5]:
        b('rear panel inset',(x,5.33,.63),(.82,.035,.76),edge,.025,back)
    for y in [-4.75,-3.75,-2.75,-1.75,-.75,.25,1.25,2.25,3.25,4.25]:
        b('side panel inset',(-4.83,y,.63),(.035,.82,.76),edge,.025,left)
    # Four facing pairs; seat centres map exactly to restaurantPlacements.
    for x in [-2.7,2.7]:
        for y in [-1.5,1.5]:
            c('table pedestal',(x,y,.49),.10,.95,wood)
            c('table foot',(x,y,.08),.37,.10,wood,vertices=32)
            b('table walnut rim',(x,y,1.02),(1.48,1.08,.10),wood,.05)
            b('tablecloth',(x,y,1.079),(1.45,1.05,.018),linen,.025)
            for side in [-1,1]:
                cy=y+side*.98
                b('banquette plinth',(x,cy,.26),(1.35,.72,.43),wood,.03)
                b('banquette cushion',(x,cy,.61),(1.36,.73,.16),red,.055)
                b('banquette back',(x,cy+side*.39,.95),(1.4,.17,.73),red,.055)
                for dx in [-.42,0,.42]:
                    b('upholstery piping',(x+dx,cy+side*.29,.98),(.013,.014,.51),brass,.005)
                py=y+side*.27
                c('charger',(x,py,1.097),.21,.013,brass,vertices=40)
                c('dinner plate',(x,py,1.109),.19,.014,china,vertices=40)
                c('plate well',(x,py,1.119),.145,.009,linen,vertices=40)
                b('folded napkin',(x+.39,py,1.103),(.14,.24,.02),linen,.005)
                b('table knife',(x+.27,py,1.101),(.025,.25,.016),steel,.006)
                b('table fork',(x-.27,py,1.101),(.018,.20,.014),steel,.004)
                for j in range(3):b('fork tine',(x-.285+j*.015,py+.10,1.102),(.007,.055,.011),steel)
                c('glass foot',(x-.47,py,1.108),.065,.011,steel)
                c('glass stem',(x-.47,py,1.16),.009,.10,steel,vertices=12)
                c('wine glass bowl',(x-.47,py,1.235),.058,.08,china,vertices=24)
            c('bud vase',(x,y,1.16),.055,.15,green)
            c('flower stem',(x,y,1.31),.007,.18,green,vertices=8)
            for dx,dy in [(-.025,0),(.025,0),(0,.025),(0,-.025)]:c('red carnation',(x+dx,y+dy,1.41),.035,.035,red,vertices=12)
    # Wine sideboard and open racks are fixed to the rear cutaway group.
    b('wine sideboard',(-3.1,4.83,.56),(2.5,.87,1.08),wood,.025,back)
    b('sideboard stone',(-3.1,4.83,1.14),(2.64,.95,.09),cream,.025,back)
    for x in [-3.9,-3.1,-2.3]:
        b('cabinet door',(x,4.37,.59),(.73,.05,.89),edge,.018,back)
        c('cabinet knob',(x+.25,4.30,.64),.03,.06,brass,(math.pi/2,0,0),16,back)
    for z in [1.54,2.15]:
        b('wine shelf',(-3.1,5.04,z),(2.7,.61,.065),wood,parent=back)
        for i in range(8):
            x=-4.23+i*.32
            c('wine bottle',(x,5.03,z+.22),.078,.35,green,parent=back)
            c('bottle neck',(x,5.03,z+.44),.032,.14,green,parent=back)
            b('wine label',(x,4.949,z+.23),(.105,.009,.13),linen,parent=back)
    b('private door surround',(0,5.32,1.35),(1.65,.15,2.7),wood,.025,back)
    b('private door',(0,5.21,1.29),(1.38,.08,2.5),edge,.02,back)
    for z in [.61,1.76]:b('private door panel',(0,5.15,z),(1.1,.04,.91),wood,.025,back)
    c('private brass knob',(.49,5.07,1.17),.045,.08,brass,(math.pi/2,0,0),24,back)
    text('private sign','PRIVATE',(0,5.10,2.18),.15)
    b('service hatch dark recess',(3.05,5.38,2.02),(2.62,.035,1.40),dark,parent=back)
    for x in [1.68,4.42]:b('hatch side',(x,5.23,2.01),(.12,.27,1.55),wood,parent=back)
    b('hatch lintel',(3.05,5.23,2.78),(2.87,.27,.12),wood,parent=back)
    b('service counter',(3.05,4.95,1.25),(2.93,.90,.11),cream,.035,back)
    for x in [2.15,2.65]:
        for z in [1.32,1.35,1.38,1.41]:c('stacked plates',(x,4.92,z),.18,.024,china,vertices=32,parent=back)
    c('coffee urn',(3.62,4.96,1.65),.22,.67,steel,vertices=32,parent=back)
    c('urn lid',(3.62,4.96,2.01),.25,.055,steel,parent=back)
    b('urn tap',(3.62,4.68,1.48),(.07,.18,.07),brass,parent=back)
    text('restaurant name',"VITTORIA'S",(0,5.29,3.02),.28)
    for x in [-2.7,2.7]:
        for y in [-1.5,1.5]:
            c('pendant cable',(x,y,3.15),.012,.65,dark,vertices=12)
            c('pendant shade',(x,y,2.79),.34,.15,lamp,vertices=40)
            c('shade brass rim',(x,y,2.73),.35,.027,brass,vertices=40)
    # Framed wall decoration: layered geometric still life, not an external image.
    for y in [-2.7,.2,3.1]:
        b('picture frame',(-4.81,y,2.28),(.09,1.28,.90),brass,.015,left)
        b('picture canvas',(-4.75,y,2.28),(.025,1.13,.75),linen,parent=left)
        b('painted table',(-4.73,y,2.05),(.012,.94,.06),wood,parent=left)
        b('painted bottle',(-4.715,y-.21,2.30),(.012,.19,.44),green,parent=left)
        c('painted fruit',(-4.70,y+.18,2.17),.14,.012,red,(0,math.pi/2,0),32,left)
