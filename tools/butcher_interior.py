"""Original period butcher shop; Blender metres, Z up. No external assets."""
import math
import bpy


def build(box,cylinder,material):
    tile=material('Butcher glazed ivory tile',(.73,.73,.61))
    green=material('Butcher bottle green enamel',(.085,.19,.15),.15)
    dark=material('Butcher charcoal grout',(.15,.17,.14))
    nickel=material('Butcher brushed nickel',(.55,.59,.55),.75)
    wood=material('Butcher end grain maple',(.47,.30,.14))
    paper=material('Butcher wrapping paper',(.71,.61,.43))
    red=material('Butcher fresh cuts',(.36,.065,.055))
    fat=material('Butcher ivory fat',(.78,.64,.48))
    black=material('Butcher scale numerals',(.025,.035,.025))
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def label(name,text,p,size,parent=None):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(black);o.parent=parent
        return o
    b('shop foundation',(0,0,-.08),(9,9,.18),dark)
    for row in range(18):
        for col in range(18):
            b('quarry tile',(-4.25+col*.5,-4.25+row*.5,.01),(.487,.487,.015),tile if (row+col)%2 else green)
    b('left wall',(-4.5,0,1.8),(.14,9,3.6),dark,parent=left)
    b('rear wall',(0,4.5,1.8),(9,.14,3.6),dark,parent=back)
    for row in range(11):
        z=.16+row*.30
        for col in range(18):
            b('rear glazed tile',(-4.25+col*.5,4.40,z),(.487,.035,.287),tile,parent=back)
            b('left glazed tile',(-4.40,-4.25+col*.5,z),(.035,.487,.287),tile,parent=left)
    for z in [.12,2.42,3.35]:
        b('rear tile border',(0,4.36,z),(9,.055,.085),green,parent=back)
        b('left tile border',(-4.36,0,z),(.055,9,.085),green,parent=left)
    # Low open display, with a glazed customer face and a readable row of trays.
    b('display cabinet',(-.65,.55,.46),(5.2,1.25,.87),green,.05)
    for x in [-2.7,-1.2,.3,1.4]:b('cabinet ventilation',(x,-.087,.37),(.8,.025,.28),dark)
    b('display bed',(-.65,.55,.94),(5.15,1.2,.09),nickel,.02)
    for x in [-2.65,-1.35,-.05,1.25]:
        b('display tray',(x,.55,1.015),(1.12,.91,.035),tile,.04)
        for j in range(3):
            o=cylinder('cut of beef',(x+(j-1)*.28,.56,1.07),.19,.065,red,vertices=18);o.scale.y=1.45
            cylinder('marrow centre',(x+(j-1)*.28,.55,1.108),.042,.009,fat,vertices=14)
        b('tray ticket',(x,.05,1.08),(.35,.06,.13),paper,.01)
    glass=material('Butcher display glass',(.58,.74,.68))
    shader=glass.node_tree.nodes.get('Principled BSDF');shader.inputs['Alpha'].default_value=.16
    glass.diffuse_color=(.58,.74,.68,.16);glass.surface_render_method='DITHERED'
    b('display front glass',(-.65,-.12,1.2),(5.16,.018,.50),glass)
    for x in [-3.26,1.96]:b('display end glass',(x,.53,1.2),(.018,1.3,.50),glass)
    b('display rail',(-.65,-.14,1.47),(5.28,.045,.045),nickel,.015)
    for x in [-3.28,1.98]:b('display upright',(x,-.14,1.2),(.045,.045,.57),nickel,.01)
    # Scale on the right service island, with physical dial and tick marks.
    b('service island',(3.1,.7,.50),(1.5,1.45,.97),green,.025)
    b('marble service top',(3.1,.7,1.02),(1.65,1.56,.08),tile,.025)
    b('scale foot',(3.1,.63,1.14),(.62,.65,.16),green,.04)
    cylinder('scale dial housing',(3.1,.60,1.48),.28,.18,nickel,(math.pi/2,0,0),48)
    cylinder('scale dial face',(3.1,.495,1.48),.246,.012,tile,(math.pi/2,0,0),48)
    for i in range(17):
        a=math.pi*.18+i*math.pi*1.64/16
        o=b('scale dial mark',(3.1+math.sin(a)*.21,.481,1.48+math.cos(a)*.21),(.009,.009,.028),black);o.rotation_euler.y=a
    o=b('scale needle',(3.13,.47,1.56),(.012,.012,.17),black);o.rotation_euler.y=.32
    b('scale tray',(3.1,.72,1.80),(.77,.59,.045),nickel,.05)
    # Work zone: block, knife rack and paper dispenser, outside customer aisles.
    b('endgrain chopping block',(-1.8,3.12,.99),(2.4,1.1,.32),wood,.03)
    for x in [-2.72,-.88]:
        for y in [2.74,3.5]:b('block leg',(x,y,.44),(.16,.16,.84),green,.015)
    for i in range(15):b('block endgrain seam',(-2.9+i*.157,3.12,1.156),(.005,1.03,.002),paper)
    cylinder('paper roll',(1.25,3.38,1.19),.19,1.1,paper,(0,math.pi/2,0),32)
    for x in [.62,1.88]:b('paper roll stand',(x,3.38,1.1),(.065,.5,.5),nickel)
    b('packing table',(1.25,3.18,.89),(1.6,1.08,.10),wood,.015)
    for x in [.65,1.85]:
        for y in [2.8,3.56]:b('packing leg',(x,y,.44),(.065,.065,.84),nickel)
    b('knife rack',(-2,4.26,1.86),(1.8,.12,.14),wood,.015,back)
    for x in [-2.6,-2.2,-1.8,-1.4]:
        b('knife blade',(x,4.14,1.66),(.10,.025,.28),nickel,.015,back)
        b('knife handle',(x,4.14,1.96),(.06,.05,.24),black,.015,back)
    b('cold room surround',(3.4,4.22,1.32),(1.62,.24,2.6),nickel,.04,back)
    b('cold room insulated door',(3.4,4.06,1.32),(1.40,.13,2.39),green,.04,back)
    b('cold room lever',(2.91,3.95,1.21),(.10,.12,.49),nickel,.025,back)
    label('cold room sign','COLD ROOM',(3.4,3.978,2.1),.12,back)
    b('shop sign board',(-.55,4.27,2.95),(4.8,.1,.57),paper,.02,back)
    label('shop name','FASSANO MEATS',(-.55,4.20,2.78),.32,back)
    for x in [-2.4,1.2]:
        cylinder('opal pendant',(x,1.5,2.95),.29,.13,tile,vertices=32)
        cylinder('pendant suspension',(x,1.5,3.35),.018,.67,nickel,vertices=10)
