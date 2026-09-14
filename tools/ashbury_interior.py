"""Ashbury Court's locally authored entrance hall. Blender metres, Z up."""
import math
import bpy


def build(box, cylinder, material):
    cream=material('Ashbury travertine',(.66,.61,.49))
    pale=material('Ashbury ivory plaster',(.76,.72,.61))
    black=material('Ashbury dark marble',(.075,.085,.076))
    stone=material('Ashbury stair stone',(.45,.43,.36))
    walnut=material('Ashbury walnut',(.20,.095,.045))
    green=material('Ashbury green upholstery',(.095,.20,.14))
    brass=material('Ashbury brushed brass',(.52,.36,.13),.65)
    iron=material('Ashbury lift iron',(.045,.055,.048),.4)
    glass=material('Ashbury opal glass',(.78,.72,.53),0,.3)
    paper=material('Ashbury stationery',(.84,.77,.60))
    ink=material('Ashbury lettering',(.07,.085,.065))
    terracotta=material('Ashbury glazed planter',(.29,.15,.075))
    leaf=material('Ashbury foliage',(.07,.19,.10))
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,xyz,dims,mat,bevel=0,parent=None):
        ob=box(name,xyz,dims,mat,bevel);ob.parent=parent;return ob
    def text(name,value,xyz,size,mat,parent=None):
        curve=bpy.data.curves.new(name,'FONT');curve.body=value;curve.align_x='CENTER';curve.size=size;curve.extrude=.001
        ob=bpy.data.objects.new(name,curve);bpy.context.collection.objects.link(ob);ob.location=xyz;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(mat);ob.parent=parent
        bpy.context.view_layer.objects.active=ob;ob.select_set(True);bpy.ops.object.convert(target='MESH');ob.select_set(False)
    b('floor substrate',(0,0,-.09),(10,10,.20),black)
    for ix in range(20):
        for iy in range(20):
            x,y=-4.75+ix*.5,-4.75+iy*.5
            border=ix in [1,18] or iy in [1,18]
            b('stone floor tile',(x,y,.0135),(.495,.495,.008),black if border else cream)
    # Two retained walls, with a clear cutaway foreground and east side.
    b('west plaster',(-5,0,2.1),(.16,10,4.2),pale,parent=left)
    # Upper stair opening continues beyond the cutaway, rather than ending at a wall.
    b('rear plaster',(-1.275,5,2.1),(7.45,.16,4.2),pale,parent=back)
    b('rear right pier',(4.8,5,2.1),(.4,.16,4.2),pale,parent=back)
    b('wall below landing',(3.525,5,1.3),(2.15,.16,2.6),pale,parent=back)
    b('upper landing',(3.53,4.71,2.57),(2,.58,.14),stone)
    for z,height in [(.18,.3),(1.13,.09),(3.95,.24)]:
        b('west stone molding',(-4.86,0,z),(.14,10,height),cream,parent=left)
        if z<2.6:b('rear stone molding',(0,4.86,z),(10,.14,height),cream,parent=back)
        else:b('rear stone molding',(-1.275,4.86,z),(7.45,.14,height),cream,parent=back)
    for y in [-4,-2,0,2,4]:
        b('west pilaster',(-4.79,y,2),(.25,.27,3.8),cream,parent=left)
        for z in [.35,3.75]:b('pilaster capital',(-4.75,y,z),(.32,.4,.16),brass,.02,left)
    # Reception counter on west rear; a real aisle behind it for the clerk.
    b('reception front',(-2.65,2.1,.51),(2.5,.65,.95),walnut,.035)
    b('reception marble top',(-2.65,2.1,1.02),(2.66,.8,.10),black,.045)
    for x in [-3.45,-2.65,-1.85]:
        b('counter inset',(x,1.76,.55),(.65,.025,.66),cream,.02)
        b('counter brass trim',(x,1.74,.89),(.65,.025,.025),brass)
    b('register',(-2.7,2.1,1.10),(.65,.42,.05),paper)
    b('register spine',(-2.7,2.1,1.13),(.025,.42,.015),walnut)
    cylinder('service bell',(-1.72,2.02,1.12),.09,.11,brass,vertices=24)
    # Sixty-four numbered boxes make the residential capacity visible.
    b('mail wall',(-2.6,4.77,2.1),(3.65,.18,1.92),walnut,parent=back)
    for row in range(8):
        for col in range(8):
            x=-4.16+col*.445;z=1.31+row*.225
            b('mailbox door',(x,4.64,z),(.42,.065,.205),brass,.01,back)
            b('mail slot',(x,4.60,z+.045),(.27,.012,.02),iron,parent=back)
            text('mail number',str(row*8+col+1),(x,4.598,z-.045),.066,ink,back)
    text('hall name','ASHBURY COURT',(-2.6,4.78,3.28),.28,brass,back)
    # Lift recess with brass framed folding gate. This is presentation, not a
    # selectable new destination or a claim that every flat has been rendered.
    b('lift recess',(.15,4.79,1.55),(1.5,.18,2.95),iron,parent=back)
    for x in [-.7,1]:b('lift jamb',(x,4.55,1.6),(.15,.3,3.18),cream,.02,back)
    b('lift lintel',(.15,4.55,3.17),(1.85,.3,.2),cream,.02,back)
    for i in range(11):b('gate vertical',(-.53+i*.135,4.48,1.48),(.018,.018,2.75),brass,parent=back)
    for z in [.4,.9,1.4,1.9,2.4]:
        for side in [-1,1]:
            ob=b('gate cross brace',(.15,4.46,z),(.023,.018,1.65),brass,parent=back);ob.rotation_euler.y=side*.95
    text('lift sign','LIFT',(.15,4.36,3.48),.18,brass,back)
    cylinder('lift button',(1.26,4.5,1.25),.045,.035,brass,(math.pi/2,0,0),16).parent=back
    # Stone stair along east side, away from public staging bays.
    for step in range(12):
        y=.05+step*.38;h=(step+1)*.22
        b('stair riser',(3.53,y,h/2),(2,.38,h),stone)
        b('stair nosing',(3.53,y-.17,h+.018),(2.06,.045,.035),cream)
        for x in [2.47,4.58]:
            cylinder('stair baluster',(x,y,h+.48),.025,.96,brass,vertices=10)
    def beam(name,a,c,width,mat):
        from mathutils import Vector
        a,c=Vector(a),Vector(c);ob=b(name,(a+c)/2,(width,width,(c-a).length),mat,.01);ob.rotation_euler=(c-a).to_track_quat('Z','Y').to_euler()
    for x in [2.47,4.58]:beam('stair handrail',(x,-.1,1.1),(x,4.4,3.64),.075,walnut)
    text('stairs sign','FLATS 1–64',(3.5,4.82,2.30),.16,brass,back)
    # Upholstered west bench, facing into the central hall.
    for y in [-2.8,-1.4,0]:
        b('bench cushion',(-4.13,y,.725),(.76,1.27,.09),green,.04)
    b('bench frame',(-4.15,-1.4,.61),(.88,4.3,.12),walnut,.025)
    b('bench back',(-4.55,-1.4,1.05),(.14,4.4,.80),walnut,.03)
    b('bench upholstery',(-4.45,-1.4,1.09),(.1,4.12,.55),green,.045)
    for y in [-3.4,.6]:
        for x in [-4.45,-3.85]:b('bench foot',(x,y,.29),(.09,.09,.55),brass,.01)
    # Notice board and period lamps are attached to the cutaway wall.
    b('notice frame',(-4.72,-3.9,2.15),(.14,1.30,1.15),walnut,.035,left)
    b('notice backing',(-4.63,-3.9,2.15),(.03,1.12,.98),green,parent=left)
    for y,z in [(-4.18,2.3),(-3.65,2.28),(-3.85,1.99)]:b('tenant notice',(-4.606,y,z),(.01,.32,.37),paper,parent=left)
    for y in [-2.3,1.1]:
        b('sconce bracket',(-4.65,y,2.95),(.45,.055,.055),brass,parent=left)
        ob=cylinder('opal sconce',(-4.43,y,3.1),.15,.36,glass,vertices=24);ob.parent=left
    for x,y in [(1.75,3.5),(3.7,-3.75)]:
        cylinder('plant pot',(x,y,.27),.27,.5,terracotta,vertices=24)
        cylinder('plant soil',(x,y,.525),.24,.015,walnut,vertices=24)
        for i in range(9):
            angle=i*math.tau/9
            beam('leaf stem',(x,y,.54),(x+math.cos(angle)*.28,y+math.sin(angle)*.28,1.15+(i%3)*.18),.02,leaf)
            ob=b('rubber plant leaf',(x+math.cos(angle)*.30,y+math.sin(angle)*.30,1.14+(i%3)*.18),(.16,.38,.035),leaf,.04);ob.rotation_euler=(.3,0,angle)
