"""The Green Baize: six original billiard tables, Blender Z-up metres."""
import math
import bpy
import bmesh


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
            width,length,radius,height=1.27,2.54,.028575,.78
            for dx in [-.48,.48]:
                for dy in [-.97,.97]:
                    b('table turned foot',(x+dx,y+dy,.10),(.25,.25,.18),wood,.035)
                    c('table leg',(x+dx,y+dy,.36),.10,.50,edge,vertices=16)
                    c('leg brass collar',(x+dx,y+dy,.19),.115,.055,brass,vertices=16)
            for dx in [-width/2-.10,width/2+.10]:
                for dy in [-length/4,length/4]:
                    b('split table apron',(x+dx,y+dy,height-.19),(.18,length/2-.18,.28),wood,.012)
                    b('walnut side rail',(x+dx+( .02 if dx>0 else -.02),y+dy,height+.041),(.10,length/2-.19,.045),edge,.008)
            for dy in [-length/2-.10,length/2+.10]:
                b('table end apron',(x,y+dy,height-.19),(width-.18,.18,.28),wood,.012)
                b('walnut end rail',(x,y+dy+(.02 if dy>0 else -.02),height+.041),(width-.17,.10,.045),edge,.008)
            cloth=b('playing cloth',(x,y,height-.009),(width+.26,length+.26,.018),felt)
            pockets=[(-width/2-.026,-length/2-.026),(-width/2-.026,length/2+.026),(width/2+.026,-length/2-.026),(width/2+.026,length/2+.026),(-width/2-.045,0),(width/2+.045,0)]
            for dx,dy in pockets:
                hole=c('pocket cutter',(x+dx,y+dy,height),.076,.35,black,vertices=32)
                bpy.context.view_layer.objects.active=cloth
                mod=cloth.modifiers.new('open pocket','BOOLEAN');mod.operation='DIFFERENCE';mod.object=hole
                bpy.ops.object.modifier_apply(modifier=mod.name);bpy.data.objects.remove(hole,do_unlink=True)
                # Open leather cup and recessed floor; no solid cylinder cap at cloth height.
                verts=[(x+dx+math.cos(i*math.tau/32)*r,y+dy+math.sin(i*math.tau/32)*r,z) for r,z in [(.076,height),(.063,height-.15)] for i in range(32)]
                faces=[(i,(i+1)%32,(i+1)%32+32,i+32) for i in range(32)]
                faces.append(tuple(range(63,31,-1)))
                mesh=bpy.data.meshes.new('leather pocket basket');mesh.from_pydata(verts,[],faces);mesh.materials.append(black)
                cup=bpy.data.objects.new('leather pocket basket',mesh);bpy.context.collection.objects.link(cup)
                bpy.ops.mesh.primitive_torus_add(major_segments=32,minor_segments=8,location=(x+dx,y+dy,height+.002),major_radius=.077,minor_radius=.006)
                bpy.context.object.name='leather pocket lip';bpy.context.object.data.materials.append(leather)
            rails=[(.085,0,width-.085,0),(.085,length,width-.085,length)]
            for xx in [0,width]:
                outside=-.045 if xx==0 else width+.045
                rails.extend([(xx,.085,xx,length/2-.075),(xx,length/2+.075,xx,length-.085),(xx,length/2-.075,outside,length/2-.065),(xx,length/2+.075,outside,length/2+.065)])
            for xx in [0,width]:
                for yy in [0,length]:
                    sx=1 if xx==0 else -1;sy=1 if yy==0 else -1
                    rails.extend([(xx+sx*.085,yy,xx+sx*.049,yy-sy*.035),(xx,yy+sy*.085,xx-sx*.035,yy+sy*.049)])
            for ax,ay,bx,by in rails:
                distance=math.hypot(bx-ax,by-ay);nx,ny=-(by-ay)/distance,(bx-ax)/distance
                if nx*(width/2-(ax+bx)/2)+ny*(length/2-(ay+by)/2)>0:nx,ny=-nx,-ny
                profile=[(0,radius),(.012,.052),(.059,.052),(.065,0),(.024,0)]
                verts=[(x+xx+nx*out-width/2,y+yy+ny*out-length/2,height+z) for xx,yy in [(ax,ay),(bx,by)] for out,z in profile]
                faces=[(0,4,3,2,1),(5,6,7,8,9)]+[(i,(i+1)%5,(i+1)%5+5,i+5) for i in range(5)]
                mesh=bpy.data.meshes.new('physical cushion profile');mesh.from_pydata(verts,[],faces);mesh.materials.append(felt)
                ob=bpy.data.objects.new('physical cushion profile',mesh);bpy.context.collection.objects.link(ob)
                bm=bmesh.new();bm.from_mesh(mesh);bmesh.ops.recalc_face_normals(bm,faces=list(bm.faces));bm.to_mesh(mesh);bm.free()
            for dx in [-width/2-.12,width/2+.12]:
                for i in [1,2,3,5,6,7]:c('rail sight',(x+dx,y-length/2+length*i/8,height+.065),.008,.003,ivory,vertices=4)
            for i,(dx,dy) in enumerate([(-.48,-.68),(.39,.82),(.16,-.33),(-.18,.60),(.57,.19)]):
                bpy.ops.mesh.primitive_uv_sphere_add(segments=16,ring_count=8,radius=radius,location=(x+dx,y+dy,height+radius))
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
