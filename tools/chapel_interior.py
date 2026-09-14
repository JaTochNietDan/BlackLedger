"""Thorne & Sons public chapel of rest; original Blender geometry in metres."""
import math
import bpy


def build(box,cylinder,material):
    plaster=material('Thorne warm plaster',(.53,.48,.39));wood=material('Thorne walnut',(.22,.105,.053))
    panel=material('Thorne panel inset',(.31,.17,.082));velvet=material('Thorne burgundy upholstery',(.23,.065,.065))
    brass=material('Thorne aged brass',(.53,.39,.16),.65);ivory=material('Thorne ivory',(.83,.78,.65))
    carpet=material('Thorne muted carpet',(.23,.25,.20));green=material('Thorne leaves',(.13,.22,.11))
    dark=material('Thorne lettering',(.045,.04,.033));opal=material('Thorne lamps',(.93,.76,.47),0,.7)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,m,rot,vertices);ob.parent=parent;return ob
    def label(name,text,p,size):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(ivory);ob.parent=back
    b('chapel foundation',(0,0,-.09),(10,11,.20),wood)
    b('fitted carpet',(0,0,.008),(9.92,10.92,.016),carpet)
    b('back plaster',(0,5.5,1.9),(10,.16,3.8),plaster,parent=back)
    b('left plaster',(-5,0,1.9),(.16,11,3.8),plaster,parent=left)
    for x in [-4,-2,0,2,4]:
        b('back wainscot',(x,5.39,.62),(1.91,.08,1.18),wood,.012,back)
        b('back panel inset',(x,5.33,.62),(1.60,.025,.87),panel,.012,back)
    for y in [-4.4,-2.2,0,2.2,4.4]:
        b('side wainscot',(-4.89,y,.62),(.08,2.11,1.18),wood,.012,left)
        b('side panel inset',(-4.83,y,.62),(.025,1.80,.87),panel,.012,left)
    for z,h in [(1.24,.09),(3.60,.18)]:
        b('back cornice',(0,5.31,z),(10,.18,h),wood,.02,back)
        b('side cornice',(-4.81,0,z),(.18,11,h),wood,.02,left)
    # Closed display casket on a raised bier; no invented deceased occupant.
    outline=[(-.42,-1.05),(.42,-1.05),(.56,-.42),(.46,1.04),(-.46,1.04),(-.56,-.42)]
    verts=[(x,y+3.65,z) for z in [.70,1.21] for x,y in outline]
    faces=[tuple(range(5,-1,-1)),tuple(range(6,12))]+[(i,(i+1)%6,(i+1)%6+6,i+6) for i in range(6)]
    mesh=bpy.data.meshes.new('casket joinery');mesh.from_pydata(verts,[],faces);mesh.materials.append(wood)
    ob=bpy.data.objects.new('closed display casket',mesh);bpy.context.collection.objects.link(ob)
    bevel=ob.modifiers.new('soft joinery','BEVEL');bevel.width=.035;bevel.segments=3;ob.modifiers.new('normals','WEIGHTED_NORMAL')
    b('casket lid',(0,3.65,1.24),(.83,1.95,.10),panel,.075)
    for y in [2.9,4.35]:
        b('bier cross rail',(0,y,.63),(1.55,.16,.14),wood,.02)
        for x in [-.62,.62]:b('bier leg',(x,y,.31),(.15,.15,.60),wood,.02)
    for x in [-.52,.52]:
        for y in [3.05,3.65,4.25]:
            for yy in [y-.13,y+.13]:c('handle mount',(x,yy,.95),.055,.035,brass,(0,math.pi/2,0))
            b('casket handle',(x*1.10,y,.94),(.04,.32,.04),brass,.015)
    b('chapel plaque',(0,5.35,2.65),(3.4,.10,.80),wood,.025,back)
    label('firm name','THORNE & SONS',(0,5.281,2.70),.27)
    label('chapel sign','CHAPEL OF REST',(0,5.279,2.43),.15)
    # Rows flank the central aisle. Their cushion tops are exact cast anchors.
    for x in [-2.65,-1.65,1.65,2.65]:
        for y in [-.4,1.1]:
            b('visitor cushion',(x,y,.555),(.61,.58,.07),velvet,.04)
            for dx in [-.25,.25]:
                for dy in [-.24,.24]:b('chair leg',(x+dx,y+dy,.275),(.055,.055,.53),wood,.009)
            for dx in [-.26,.26]:b('back upright',(x+dx,y-.26,.82),(.06,.06,.68),wood,.008)
            b('padded chair back',(x,y-.26,1.02),(.60,.09,.34),velvet,.035)
    # Reception desk and register beside, not across, the entrance route.
    b('reception desk',(-3.4,-3.65,.49),(2.35,.82,.94),wood,.025)
    b('reception writing top',(-3.4,-3.65,1.0),(2.47,.96,.10),panel,.035)
    b('appointments ledger',(-3.45,-3.65,1.079),(.65,.44,.055),ivory,.009)
    for i in range(6):b('ledger rule',(-3.45,-3.82+i*.055,1.11),(.54,.009,.003),dark)
    c('pen stand',(-2.65,-3.60,1.1),.075,.1,brass)
    c('desk pen',(-2.65,-3.60,1.26),.012,.23,dark,vertices=12)
    for x in [-2.1,2.1]:
        c('flower stand foot',(x,3.8,.08),.24,.12,brass)
        c('flower stand column',(x,3.8,.65),.055,1.08,brass)
        c('flower vase',(x,3.8,1.35),.17,.38,brass)
        for k in range(7):
            a=k*math.tau/7;xx=x+.19*math.cos(a);yy=3.8+.19*math.sin(a);zz=1.70+(k%3)*.10
            c('flower stem',(xx,yy,zz-.13),.008,.55,green,vertices=8)
            for j in range(5):
                ang=j*math.tau/5
                petal=b('ivory flower petal',(xx+.058*math.cos(ang),yy+.058*math.sin(ang),zz),(.115,.045,.025),ivory,.02);petal.rotation_euler.z=ang
            c('flower centre',(xx,yy,zz+.02),.023,.025,brass,vertices=12)
    for y in [-2,2.3]:
        b('wall sconce plate',(-4.82,y,2.55),(.06,.18,.36),brass,.03,left)
        c('sconce shade',(-4.61,y,2.58),.17,.28,opal,parent=left)
    for x in [-2.1,2.1]:
        c('pendant cord',(x,-.4,3.4),.012,.6,dark,vertices=12)
        c('pendant shade',(x,-.4,3.03),.29,.20,opal)
