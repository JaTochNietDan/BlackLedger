"""Original Bellwether Mortuary receiving room, Blender metres/Z-up."""
import math
import bpy


def build(box, cylinder, material):
    ivory=material('Mortuary glazed ivory',(.72,.75,.68))
    jade=material('Mortuary jade tile',(.24,.38,.33))
    stone=material('Mortuary terrazzo',(.39,.42,.38))
    steel=material('Mortuary brushed steel',(.47,.53,.53),.7)
    enamel=material('Mortuary cream enamel',(.68,.69,.56),.3)
    wood=material('Mortuary oak reception',(.25,.14,.065))
    black=material('Mortuary bakelite and rubber',(.028,.036,.033))
    paper=material('Mortuary registry paper',(.83,.78,.62))
    glow=material('Mortuary milk glass',(.91,.87,.73),0,.8)
    brass=material('Mortuary nameplates',(.57,.40,.16),.6)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),parent=None):
        ob=cylinder(name,p,r,d,m,rot,24);ob.parent=parent;return ob
    def label(name,text,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(black);ob.parent=parent
    b('receiving foundation',(0,0,-.1),(10,10,.2),stone)
    for x in range(20):
        for y in range(20):
            b('terrazzo floor',(-4.75+x*.5,-4.75+y*.5,.012),(.486,.486,.025),ivory if (x+y)%2 else stone)
    b('back tiled wall',(0,5,1.85),(10,.16,3.7),ivory,parent=back)
    b('left tiled wall',(-5,0,1.85),(.16,10,3.7),ivory,parent=left)
    for z in [.25,.75,1.25]:
        for x in range(20):b('back jade tile',(-4.75+x*.5,4.9,z),(.485,.045,.48),jade,parent=back)
        for y in range(20):b('side jade tile',(-4.9,-4.75+y*.5,z),(.045,.485,.48),jade,parent=left)
    # Six individual cold cabinets with physical hinges, latch handles and plates.
    for x in [-3.5,-1.9,-.3]:
        b('insulated cabinet',(x,4.1,1.35),(1.5,1.55,2.7),enamel,.025,back)
        for z in [.72,1.95]:
            b('rubber door gasket',(x,3.30,z),(1.35,.05,1.12),black,.035,back)
            b('cold cabinet door',(x,3.25,z),(1.27,.10,1.04),steel,.04,back)
            b('door latch',(x+.43,3.14,z),(.10,.12,.36),black,.025,back)
            b('cabinet name plate',(x,3.18,z+.25),(.34,.025,.12),brass,.01,back)
            for dz in [-.3,.3]:b('door hinge',(x-.59,3.15,z+dz),(.10,.08,.15),steel,.012,back)
    # Empty preparation trolley and stainless basin; no invented bodies.
    for x in [-3.65,-2.65]:
        for y in [-.65,1.25]:
            c('trolley leg',(x,y,.49),.038,.9,steel)
            c('rubber caster',(x,y,.12),.10,.065,black,(math.pi/2,0,0))
    b('trolley shelf',(-3.15,.3,.3),(1.15,2.05,.065),steel,.02)
    b('preparation top',(-3.15,.3,1.0),(1.32,2.30,.1),steel,.045)
    for x in [-3.79,-2.51]:b('raised tray edge',(x,.3,1.07),(.03,2.30,.1),steel,.01)
    b('sink cabinet',(-4.35,2.4,.55),(1,1.10,1.1),enamel,.03)
    b('sink rim',(-4.35,2.4,1.15),(1.10,1.2,.08),steel,.02)
    b('basin hollow',(-4.35,2.4,1.196),(.72,.75,.02),black,.06)
    c('tap stem',(-4.35,2.9,1.36),.035,.37,steel)
    c('tap spout',(-4.35,2.80,1.53),.035,.22,steel,(math.pi/2,0,0))
    # Visitor registry counter, partition screen, filing drawers and telephone.
    b('reception counter',(3.2,.4,.52),(2.5,1.0,1.04),wood,.035)
    b('counter ledge',(3.2,.4,1.08),(2.65,1.15,.10),wood,.025)
    b('open register',(2.8,.3,1.15),(.6,.40,.035),paper,.01)
    b('telephone base',(3.85,.4,1.20),(.34,.30,.15),black,.055)
    b('telephone receiver',(3.85,.4,1.33),(.46,.12,.10),black,.04)
    b('filing cabinet',(4,3.9,.85),(1.15,1.2,1.7),enamel,.03)
    for z in [.27,.67,1.07,1.47]:
        b('filing drawer',(4,3.27,z),(1.03,.07,.33),steel,.01)
        b('drawer pull',(4,3.19,z),(.27,.08,.045),black,.015)
    for y in [-3.9,-2.7]:
        b('visitor chair seat',(-4.1,y,.55),(.7,.7,.12),wood,.035)
        b('visitor chair back',(-4.43,y,.98),(.1,.72,.78),wood,.025)
        for x in [-4.36,-3.85]:
            for dy in [-.27,.27]:b('chair leg',(x,y+dy,.26),(.065,.065,.52),wood)
    label('receiving sign','BELLWETHER MORTUARY',(1.2,4.88,3.13),.25)
    label('care sign','RECEIVING & CARE',(2.9,4.88,2.70),.16)
    for x in [-2,2]:
        c('lamp stem',(x,0,3.22),.028,.5,steel)
        b('opal task light',(x,0,2.94),(1.3,.40,.16),glow,.045)
