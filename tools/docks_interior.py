"""Pier 14 transit shed. Original Blender geometry; no external asset packs."""
import math
import bpy


def build(box,cylinder,material):
    concrete=material('Pier worn concrete',(.36,.38,.34));steel=material('Pier painted steel',(.18,.23,.21),.35)
    timber=material('Pier pine boards',(.40,.28,.14));lightwood=material('Pier crate slats',(.53,.39,.22))
    rust=material('Pier aged metal',(.29,.17,.09),.3);black=material('Pier rubber and stencil',(.035,.045,.039))
    paper=material('Pier docket paper',(.78,.73,.59));brass=material('Pier scale brass',(.51,.37,.13),.6)
    rope=material('Pier hemp',(.42,.36,.23));opal=material('Pier industrial lamps',(.89,.81,.62),0,.6)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,m,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,m,rot,vertices);ob.parent=parent;return ob
    def label(name,text,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=text;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(black);ob.parent=parent
    b('shed foundation',(0,0,-.10),(12,12,.23),concrete)
    for x in [-4,-2,0,2,4]:b('floor expansion seam',(x,0,.016),(.012,11.9,.003),steel)
    for y in [-4,-2,0,2,4]:b('floor expansion seam',(0,y,.016),(11.9,.012,.003),steel)
    b('rear cladding',(0,6,2.1),(12,.12,4.2),steel,parent=back)
    b('side cladding',(-6,0,2.1),(.12,12,4.2),steel,parent=left)
    for x in range(-24,25):b('rear corrugation',(x*.245,5.90,2.1),(.055,.075,4.12),rust if x%12==0 else steel,parent=back)
    for y in range(-24,25):b('side corrugation',(-5.90,y*.245,2.1),(.075,.055,4.12),steel,parent=left)
    for x in [-5.75,-2.85,0,2.85,5.75]:
        b('rear structural column',(x,5.74,2.15),(.18,.28,4.3),rust,parent=back)
        b('column foot',(x,5.70,.12),(.36,.40,.20),steel,parent=back)
    b('rear lintel',(0,5.70,4.12),(11.8,.32,.25),rust,parent=back)
    b('pier sign',(0,5.48,3.55),(3.2,.06,.66),paper,.015,back)
    label('pier lettering','PIER 14',(0,5.44,3.56),.32)
    label('shed lettering','RECEIVING & DISPATCH',(0,5.44,3.34),.14)
    def crate(x,y,z,w=1.25,d=1.1,h=1.0):
        b('freight crate',(x,y,z+h/2),(w,d,h),timber,.012)
        for xx in [-w/2+.07,w/2-.07]:
            for yy in [-d/2-.014,d/2+.014]:b('crate upright',(x+xx,y+yy,z+h/2),(.105,.055,h),lightwood,.007)
        for zz in [.08,h-.08]:
            for yy in [-d/2-.04,d/2+.04]:b('crate cross slat',(x,y+yy,z+zz),(w,.06,.10),lightwood,.006)
        for i in range(5):b('crate lid plank',(x-w/2+(i+.5)*w/5,y,z+h+.015),(w/5-.012,d,.035),lightwood)
        brace=b('crate diagonal brace',(x,y-d/2-.055,z+h/2),(w*.95,.045,.095),lightwood,.006);brace.rotation_euler.y=-math.atan2(h*.70,w*.9)
        b('freight docket',(x+.18,y-d/2-.08,z+h*.65),(.26,.012,.20),paper)
    # Cargo remains out of the public centre aisle. These are neutral scenery,
    # not a claim about the player's stock or the business inventory ledger.
    for x in [-4.25,4.25]:
        for y in [1.6,4.15]:
            for xx in [-.48,.48]:b('pallet bearer',(x+xx,y,.095),(.13,1.50,.16),timber)
            for j in range(6):b('pallet deck',(x,y-.65+j*.26,.20),(1.65,.17,.08),lightwood)
            crate(x,y,.24)
            crate(x-.08,y+.04,1.24,1.05,.98,.83)
    # Platform scale with dial and a clear approach at the front-right.
    b('weighing platform',(3.6,-1.6,.14),(1.45,1.35,.24),steel,.035)
    for x in range(8):b('platform rib',(3.0+x*.17,-1.6,.269),(.022,1.18,.018),rust)
    c('scale upright',(3.6,-.93,.91),.055,1.60,steel)
    c('scale dial case',(3.6,-.93,1.73),.31,.13,steel,(math.pi/2,0,0),48)
    c('scale dial',(3.6,-1.01,1.73),.27,.016,paper,(math.pi/2,0,0),48)
    for k in range(16):
        a=k*math.tau/16;b('scale graduation',(3.6+math.sin(a)*.22,-1.025,1.73+math.cos(a)*.22),(.014,.01,.036),black)
    b('scale needle',(3.6,-1.04,1.82),(.015,.015,.21),black)
    # Clerk's desk at the entrance, clipboard and telephone.
    b('dispatch desk',(-4,-3.6,.48),(2.1,.9,.92),timber,.02)
    b('desk top',(-4,-3.6,.99),(2.24,1.02,.10),lightwood,.025)
    b('shipping ledger',(-4.2,-3.6,1.06),(.65,.5,.04),paper)
    for j in range(6):b('ledger rule',(-4.2,-3.79+j*.065,1.083),(.54,.009,.003),black)
    b('telephone base',(-3.35,-3.65,1.10),(.35,.28,.12),black,.045)
    c('telephone dial',(-3.35,-3.70,1.175),.085,.022,brass)
    b('telephone receiver',(-3.35,-3.56,1.24),(.43,.075,.07),black,.03)
    # Hand truck and rope coils beside the left wall.
    for x in [-5.33,-4.91]:c('hand truck wheel',(x,-.45,.18),.16,.065,black,(0,math.pi/2,0))
    b('truck toe',(-5.12,-.25,.08),(.57,.48,.055),steel)
    for x in [-5.34,-4.90]:b('truck upright',(x,-.49,.79),(.055,.055,1.42),steel,.014)
    for z in [.38,.75,1.12,1.46]:b('truck cross rail',(-5.12,-.49,z),(.47,.055,.055),steel,.014)
    for level in range(4):
        bpy.ops.mesh.primitive_torus_add(major_radius=.32,minor_radius=.026,major_segments=32,minor_segments=8,location=(-4.6,-1.95,.07+level*.05));bpy.context.object.name='coiled hemp rope';bpy.context.object.data.materials.append(rope)
    for x in [-2.4,2.4]:
        c('lamp cable',(x,0,3.87),.013,.50,black,vertices=12)
        c('enamel shade',(x,0,3.58),.33,.16,steel)
        c('lamp diffuser',(x,0,3.486),.26,.022,opal)
