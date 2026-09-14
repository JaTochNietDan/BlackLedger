"""Original modest Mariner lodging room; metres, Blender Z up."""
import math
import bpy


def build(box,cylinder,material):
    wood=material('Lodging worn pine',(.38,.25,.13))
    dark=material('Lodging dark varnish',(.16,.10,.055))
    plaster=material('Lodging aged plaster',(.62,.58,.44))
    linen=material('Lodging cotton linen',(.74,.7,.56))
    green=material('Lodging wool blanket',(.24,.31,.23))
    metal=material('Lodging enamel iron',(.08,.10,.085),.4)
    nickel=material('Lodging nickel',(.5,.52,.47),.7)
    glass=material('Lodging window',(.35,.47,.46),0,.1)
    paper=material('Lodging paper',(.70,.65,.49))
    curtain=material('Lodging faded curtain',(.43,.32,.19))
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def c(name,p,r,d,m,rot=(0,0,0),parent=None):
        o=cylinder(name,p,r,d,m,rot,20);o.parent=parent;return o
    b('room base',(0,0,-.08),(6,6,.18),dark)
    for row in range(24):
        for col in range(4):
            lo=max(-3,-3.6+col*1.8+(row%2)*.8);hi=min(3,-3.6+(col+1)*1.8+(row%2)*.8)
            if hi>lo:b('pine floorboard',((lo+hi)/2,-2.875+row*.25,.014),(hi-lo-.01,.242,.007),wood)
    b('west plaster',(-3,0,1.55),(.14,6,3.1),plaster,parent=left)
    b('back plaster',(0,3,1.55),(6,.14,3.1),plaster,parent=back)
    for z,h in [(.13,.21),(2.98,.13)]:
        b('west molding',(-2.9,0,z),(.09,6,h),dark,parent=left)
        b('back molding',(0,2.9,z),(6,.09,h),dark,parent=back)
    # Narrow iron bed occupies the back-left corner, clear of the diagonal aisle.
    b('iron bed frame',(-1.95,1.55,.34),(1.35,2.18,.14),metal,.025)
    b('mattress',(-1.95,1.55,.53),(1.27,2.07,.25),linen,.07)
    b('folded blanket',(-1.95,1.08,.69),(1.28,1.15,.065),green,.02)
    b('bed pillow',(-1.95,2.29,.71),(.88,.42,.14),linen,.09)
    for y in [.45,2.65]:
        for x in [-2.62,-1.28]:c('bedpost',(x,y,.50),.035,.96,metal)
        c('bed end rail',(-1.95,y,.97),.034,1.40,metal,(0,math.pi/2,0))
        for x in [-2.45,-2.2,-1.95,-1.7,-1.45]:c('bed spindle',(x,y,.69),.018,.52,metal)
    b('bedside drawer',(-.67,2.48,.42),(.61,.65,.77),wood,.025)
    for z in [.24,.57]:
        b('drawer face',(-.67,2.137,z),(.53,.03,.25),dark,.012)
        c('drawer knob',(-.67,2.1,z),.035,.045,nickel,(math.pi/2,0,0))
    c('bedside lamp base',(-.67,2.47,.85),.13,.055,nickel)
    c('bedside lamp stem',(-.67,2.47,1.05),.018,.36,nickel)
    bpy.ops.mesh.primitive_cone_add(vertices=24,radius1=.22,radius2=.12,depth=.28,location=(-.67,2.47,1.3));bpy.context.object.name='linen lamp shade';bpy.context.object.data.materials.append(linen)
    b('window frame',(-1.68,2.84,2.15),(1.85,.16,1.24),linen,parent=back)
    b('window glazing',(-1.68,2.73,2.15),(1.67,.025,1.06),glass,parent=back)
    for x in [-2.51,-1.68,-.85]:b('window sash',(x,2.7,2.15),(.035,.045,1.08),dark,parent=back)
    b('window crossbar',(-1.68,2.7,2.15),(1.7,.045,.035),dark,parent=back)
    for x in [-2.73,-.63]:
        for i in range(4):c('curtain fold',(x+i*.043,2.66,2.14),.043,1.45,curtain,parent=back)
    # Front-right writing desk, washstand, case and coat hooks.
    b('writing desk',(1.8,-1.85,.79),(1.67,.7,.095),dark,.025)
    for x in [1.1,2.5]:
        for y in [-2.1,-1.6]:b('desk leg',(x,y,.4),(.07,.07,.77),wood)
    b('writing paper',(1.65,-1.85,.85),(.34,.4,.018),paper)
    b('pencil',(1.94,-1.8,.866),(.025,.3,.02),metal)
    b('chair seat',(1.8,-2.55,.53),(.57,.51,.08),wood,.02)
    b('chair back',(1.8,-2.79,.91),(.58,.06,.68),dark,.02)
    for x in [1.58,2.02]:
        for y in [-2.74,-2.36]:b('chair leg',(x,y,.26),(.045,.045,.49),dark)
    b('washstand',(2.55,.1,.44),(.61,.78,.83),wood,.025)
    b('washstand marble',(2.55,.1,.89),(.72,.84,.08),linen,.025)
    c('wash basin',(2.55,.1,.96),.24,.08,linen)
    c('basin hollow',(2.55,.1,1.005),.185,.008,glass)
    c('water jug',(2.55,.38,1.13),.09,.34,linen)
    b('travel trunk',(2.0,-.95,.19),(1.15,.48,.34),dark,.05)
    for x in [1.7,2.3]:b('trunk leather strap',(x,-.95,.37),(.06,.49,.018),metal)
    b('coat hook board',(-2.83,-1.2,1.8),(.09,1.2,.12),dark,parent=left)
    for y in [-1.6,-1.2,-.8]:
        b('coat hook',(-2.75,y,1.78),(.16,.025,.035),nickel,parent=left)
    # Room door on the removed foreground wall is represented by its clear aisle.
