"""The Green Baize: six original billiard tables, Blender Z-up metres."""
import math
import bpy


def build(box,cylinder,material):
    wall=material('Baize tobacco plaster',(.46,.43,.32));wood=material('Baize dark walnut',(.16,.09,.043))
    edge=material('Baize polished walnut',(.30,.18,.075));floor=material('Baize worn boards',(.30,.24,.15))
    felt=material('Baize green cloth',(.055,.23,.135));leather=material('Baize pocket leather',(.08,.045,.019))
    black=material('Baize black enamel',(.018,.027,.022));brass=material('Baize brass',(.50,.35,.11),.7)
    ivory=material('Baize ivory',(.88,.84,.66));lamp=material('Baize lit shade',(.79,.76,.51),0,.6)
    colors=[ivory,material('Baize yellow ball',(.92,.62,.06)),material('Baize blue ball',(.035,.12,.42)),material('Baize red ball',(.55,.035,.025)),black]
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,mat,bevel=0,parent=None):
        ob=box(name,p,d,mat,bevel);ob.parent=parent
        if parent==back and name not in ['rear wall','rear rail']:ob.location.y+=1
        return ob
    def c(name,p,r,d,mat,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,mat,rot,vertices);ob.parent=parent
        if parent==back:ob.location.y+=1
        return ob
    def label(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);cu.materials.append(ivory);ob.parent=parent
        if parent==back:ob.location.y+=1
    b('hall foundation',(0,0,-.09),(16,18,.20),wood)
    for x in range(32):
        for y in range(9):b('floorboard',(-7.75+x*.5,-8+y*2,.012),(.487,1.98,.012),floor if (x+y)%4 else edge)
    b('rear wall',(0,9,2),(16,.16,4),wall,parent=back);b('side wall',(-8,0,2),(.16,18,4),wall,parent=left)
    for z,h in [(.16,.26),(1.10,.08),(3.83,.14)]:
        b('rear rail',(0,8.84,z),(16,.12,h),wood,.012,back);b('side rail',(-7.84,0,z),(.12,18,h),wood,.012,left)
    for x in [-3.7,3.7]:
        for row,y in enumerate([-4.25,0,4.25]):
            for dx in [-.83,.83]:
                for dy in [-1.03,1.03]:
                    b('table turned foot',(x+dx,y+dy,.14),(.34,.34,.25),wood,.06)
                    c('table leg',(x+dx,y+dy,.53),.14,.60,edge,vertices=16)
                    c('leg brass collar',(x+dx,y+dy,.24),.16,.09,brass,vertices=16)
            b('table apron',(x,y,.83),(2.42,2.94,.35),wood,.07)
            b('table slate',(x,y,1.03),(2.50,3.04,.17),black,.035)
            # Cut the cloth at every pocket so the mouth has real depth.
            cloth=b('playing cloth',(x,y,1.122),(2.18,2.70,.028),felt)
            for dx in [-1.09,1.09]:
                for dy in [-1.35,0,1.35]:
                    hole=c('pocket cutter',(x+dx,y+dy,1.13),.135,.30,black,vertices=24)
                    bpy.context.view_layer.objects.active=cloth
                    mod=cloth.modifiers.new('open pocket','BOOLEAN');mod.operation='DIFFERENCE';mod.object=hole
                    bpy.ops.object.modifier_apply(modifier=mod.name);bpy.data.objects.remove(hole,do_unlink=True)
                    c('pocket mouth',(x+dx,y+dy,1.097),.14,.018,black,vertices=24)
                    c('leather pocket cup',(x+dx,y+dy,.95),.155,.24,leather,vertices=24)
            for dx in [-1.19,1.19]:
                for dy in [-.68,.68]:b('side cushion',(x+dx,y+dy,1.15),(.19,1.08,.14),edge,.025)
            for dy in [-1.45,1.45]:b('end cushion',(x,y+dy,1.15),(1.91,.20,.14),edge,.025)
            for dx in [-1.19,1.19]:
                for dy in [-.97,-.39,.39,.97]:c('rail sight',(x+dx,y+dy,1.225),.019,.006,ivory,vertices=12)
            for i,(dx,dy) in enumerate([(-.48,-.68),(.39,.82),(.16,-.33),(-.18,.60),(.57,.19)]):
                bpy.ops.mesh.primitive_uv_sphere_add(segments=16,ring_count=8,radius=.053,location=(x+dx,y+dy,1.188))
                ob=bpy.context.object;ob.name='billiard ball';ob.data.materials.append(colors[(i+row)%len(colors)])
                for poly in ob.data.polygons:poly.use_smooth=True
            # Individual green-shaded table lights leave the centre aisle open.
            for dy in [-.60,.60]:
                c('lamp suspension',(x,y+dy,3.54),.012,.76,black,vertices=12)
                bpy.ops.mesh.primitive_cone_add(vertices=32,radius1=.43,radius2=.16,depth=.23,location=(x,y+dy,3.08))
                shade=bpy.context.object;shade.name='green table shade';shade.data.materials.append(felt)
                c('table shade glow',(x,y+dy,2.962),.37,.012,lamp,vertices=32)
                c('shade brass rim',(x,y+dy,2.968),.435,.028,brass,vertices=32)
    # Marker counter, ledger and ball storage along the rear right wall.
    b('marker counter',(3.7,6.70,.54),(3.9,.76,1.04),wood,.025)
    b('marker counter top',(3.7,6.70,1.10),(4.04,.89,.09),edge,.03)
    b('marker ledger',(3.3,6.66,1.17),(.76,.54,.055),ivory,.012)
    for i in range(6):b('ledger rule',(3.3,6.48+i*.067,1.201),(.60,.007,.003),black)
    for z in [.62,1.20,1.78]:b('ball tray shelf',(4.0,7.62,z),(3.4,.38,.07),edge,.012,back)
    for x in [2.6,3.3,4,4.7,5.4]:b('ball return box',(x,7.63,.86),(.58,.34,.40),felt,.02,back)
    b('name board',(-2.4,7.72,3.25),(5.6,.13,.65),wood,.025,back)
    label('hall lettering','THE GREEN BAIZE',(-2.4,7.64,3.13),.38)
    b('back room door',(0,7.85,1.18),(1.25,.11,2.34),wood,.025,back)
    label('back room lettering','PRIVATE',(0,7.78,1.78),.17)
    c('back door knob',(.43,7.76,1.13),.05,.07,brass,(math.pi/2,0,0),16,back)
    # Payphone and a small score board are distinct silhouettes at the rear.
    b('payphone case',(-6.35,7.59,1.78),(.57,.40,.94),black,.04,back)
    b('coin plate',(-6.35,7.373,2.02),(.40,.03,.28),brass,.015,back)
    for xx in [-6.44,-6.27]:b('coin slot',(xx,7.351,2.06),(.065,.018,.013),black,parent=back)
    c('phone dial',(-6.35,7.357,1.70),.145,.035,brass,(math.pi/2,0,0),32,back)
    b('phone receiver',(-6.73,7.37,1.79),(.14,.17,.55),black,.055,back)
    for zz in [1.55,2.03]:c('receiver end',(-6.73,7.37,zz),.105,.15,black,vertices=24,parent=back)
    for i in range(12):c('phone curled cord',(-6.72,7.4,1.43-i*.03),.03,.013,black,vertices=12,parent=back)
    b('scoreboard',(6.62,7.74,2.45),(1.65,.12,1.45),wood,.025,back)
    for zz in [2.04,2.37,2.70]:
        b('score slider rail',(6.62,7.657,zz),(1.4,.025,.025),brass,parent=back)
        for i in range(10):c('score bead',(6.13+i*.10,7.63,zz),.039,.06,ivory,(0,math.pi/2,0),12,back)
    # Cue racks sit between spectator chairs, outside every occupant footprint.
    for y in [-2.1,2.1]:
        for z in [.20,1.67]:b('cue rack',(-7.64,y,z),(.30,1.45,.10),edge,.012,left)
        for j in range(6):
            c('rack cue',(-7.56,y-.55+j*.22,.99),.019,1.63,edge,vertices=12,parent=left)
            c('cue tip',(-7.56,y-.55+j*.22,1.817),.018,.03,ivory,vertices=12,parent=left)
    for y in [-4.25,0,4.25]:
        b('spectator chair base',(-7.2,y,.29),(.72,.72,.49),wood,.025)
        b('spectator cushion',(-7.2,y,.61),(.76,.76,.16),felt,.045)
        b('spectator chair back',(-7.62,y,.99),(.13,.78,.78),edge,.025)
