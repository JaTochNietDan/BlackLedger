"""Original Oak Ridge office and Stillwater remembrance/furnace room."""
import bpy
import math


def build(kind, box, cylinder, material):
    wood=material(kind+' aged oak',(.25,.14,.065))
    plaster=material(kind+' lime plaster',(.65,.61,.49))
    green=material(kind+' dark green',(.17,.26,.20))
    floor=material(kind+' stone floor',(.42,.42,.35))
    paper=material(kind+' registry paper',(.82,.77,.62))
    black=material(kind+' cast iron',(.045,.05,.043),.5)
    brass=material(kind+' aged brass',(.57,.40,.15),.6)
    glow=material(kind+' opal lights',(.93,.85,.66),0,.7)
    red=material(kind+' furnace firebrick',(.37,.20,.12))
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),parent=None):
        ob=cylinder(name,p,r,d,m,rot,24);ob.parent=parent;return ob
    def text(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(paper);ob.parent=parent
    b('foundation',(0,0,-.1),(10,10,.2),wood)
    for x in range(10):
        for y in range(10):b('stone floor slab',(-4.5+x,-4.5+y,.015),(.984,.984,.03),floor)
    b('back wall',(0,5,1.85),(10,.16,3.7),plaster,parent=back)
    b('west wall',(-5,0,1.85),(.16,10,3.7),plaster,parent=left)
    for z in [.18,1.2]:
        b('back wood rail',(0,4.88,z),(10,.08,.13),wood,parent=back)
        b('side wood rail',(-4.88,0,z),(.08,10,.13),wood,parent=left)
    for y in [-3,-1]:
        b('visitor bench',(-4,y,.57),(1.1,1.7,.15),wood,.04)
        b('bench back',(-4.5,y,1.02),(.14,1.7,.8),wood,.025)
        for x in [-4.4,-3.6]:
            for dy in [-.6,.6]:b('bench leg',(x,y+dy,.28),(.09,.09,.56),wood)
    b('registry desk',(3,.2,.55),(2.6,1.05,1.1),wood,.035)
    b('desk top',(3,.2,1.14),(2.8,1.18,.12),wood,.025)
    b('open register',(2.6,.15,1.225),(.62,.42,.04),paper,.01)
    c('inkwell',(3.2,.35,1.25),.07,.10,black)
    b('telephone',(3.8,.2,1.29),(.34,.28,.16),black,.06)
    b('receiver',(3.8,.2,1.42),(.45,.12,.10),black,.035)
    if kind=='cemetery':
        b('plot map frame',(-1.4,4.80,2.15),(4.6,.16,1.9),wood,.025,back)
        b('plot map paper',(-1.4,4.69,2.15),(4.35,.045,1.67),paper,parent=back)
        for x in range(8):
            for z in range(3):b('plot boundary',(-3.20+x*.50,4.65,1.89+z*.43),(.35,.018,.28),green,parent=back)
        b('archive cabinet',(3.6,4.1,1.05),(2.05,1.35,2.1),wood,.025,back)
        for z in [.32,.82,1.32,1.82]:
            for x in [3.12,4.08]:
                b('archive drawer',(x,3.38,z),(.88,.08,.43),green,.015,back)
                b('drawer label',(x,3.32,z+.08),(.28,.02,.08),paper,parent=back)
                b('drawer handle',(x,3.29,z-.09),(.2,.07,.035),brass,.01,back)
        b('grounds tool rack',(-4.55,3.3,1.65),(.16,2.3,.12),wood,parent=left)
        for y in [2.5,3.2,3.9]:
            c('shovel handle',(-4.4,y,1.0),.035,1.8,wood,parent=left)
            b('shovel blade',(-4.4,y,.23),(.07,.28,.43),black,.06,left)
        b('office sign',(.4,4.80,3.35),(5.8,.15,.42),green,parent=back)
        text('Oak Ridge lettering','OAK RIDGE  /  PLOT OFFICE',(.4,4.69,3.25),.20)
    else:
        # Closed furnace doors behind the public remembrance area.
        for x in [-3,-.6]:
            b('firebrick furnace',(x,3.75,1.15),(2.05,2.1,2.3),red,.025,back)
            b('furnace iron surround',(x,2.64,1.05),(1.50,.15,1.55),black,.04,back)
            b('closed furnace door',(x,2.52,1.05),(1.26,.16,1.29),green,.08,back)
            b('furnace door handle',(x+.42,2.39,1.05),(.07,.12,.4),brass,.02,back)
            c('flue',(x,4.1,2.9),.3,1.3,black,parent=back)
            c('temperature gauge',(x,2.45,2.04),.14,.07,paper,(math.pi/2,0,0),back)
        b('urn shelving',(3.5,4.3,1.75),(2.0,.7,2.5),wood,.02,back)
        for z in [.65,1.35,2.05,2.75]:b('urn shelf',(3.5,3.92,z),(2.0,.8,.08),wood,parent=back)
        for x in [2.85,3.5,4.15]:
            for z in [.86,1.56,2.26]:
                c('memorial urn',(x,3.90,z),.17,.32,brass,parent=back)
                c('urn lid',(x,3.90,z+.18),.19,.06,brass,parent=back)
        b('remembrance sign',(-1.8,4.85,3.45),(5.8,.1,.38),green,parent=back)
        text('Stillwater lettering','STILLWATER  /  REMEMBRANCE',(-1.8,4.77,3.35),.19)
    for x in [-2,2]:
        c('lamp suspension',(x,0,3.2),.025,.5,brass)
        c('opal pendant',(x,0,2.94),.30,.15,glow)
