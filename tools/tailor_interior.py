"""Original Ruttledge & Vance fitting salon; Blender metres, Z up."""
import math
import bpy


def build(box,cylinder,material):
    cream=material('Ruttledge cream plaster',(.68,.64,.52));oak=material('Ruttledge walnut',(.23,.13,.065))
    trim=material('Ruttledge polished walnut',(.39,.26,.13));floor=material('Ruttledge parquet',(.39,.28,.17))
    brass=material('Ruttledge brass',(.52,.37,.13),.7);iron=material('Ruttledge black iron',(.04,.055,.045),.35)
    paper=material('Ruttledge pattern paper',(.85,.79,.64));cloth=material('Ruttledge green wool',(.20,.28,.22))
    navy=material('Ruttledge navy wool',(.10,.15,.19));grey=material('Ruttledge grey wool',(.30,.32,.30));wine=material('Ruttledge curtain',(.31,.13,.10))
    mirror=material('Ruttledge silver mirror',(.55,.63,.61),.85);glass=material('Ruttledge opal lamp',(.90,.82,.64),0,.6)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,mat,bevel=0,parent=None):
        ob=box(name,p,d,mat,bevel);ob.parent=parent;return ob
    def c(name,p,r,d,mat,rot=(0,0,0),vertices=24,parent=None):
        ob=cylinder(name,p,r,d,mat,rot,vertices);ob.parent=parent;return ob
    def label(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(iron);ob.parent=parent
    b('salon foundation',(0,0,-.09),(10,10,.20),oak)
    for x in range(20):
        for y in range(10):b('parquet plank',(-4.75+x*.5,-4.5+y,.012),(.488,.988,.012),floor if (x+y)%3 else trim)
    b('rear plaster',(0,5,2),(10,.16,4),cream,parent=back);b('side plaster',(-5,0,2),(.16,10,4),cream,parent=left)
    for z,h in [(.13,.22),(1.07,.07),(3.74,.14)]:
        b('rear moulding',(0,4.85,z),(10,.13,h),oak,.014,back);b('side moulding',(-4.85,0,z),(.13,10,h),oak,.014,left)
    # Fitting counter, measurement book and swatches.
    b('fitting counter',(0,2.6,.53),(4.9,.90,1.03),oak,.025)
    b('counter polished edge',(0,2.6,1.10),(5.05,1.04,.11),trim,.03)
    for x in [-1.6,0,1.6]:b('counter front panel',(x,2.13,.54),(1.38,.045,.78),trim,.02)
    b('open measurement book',(1.48,2.52,1.19),(.82,.60,.06),paper,.012)
    for x in [1.26,1.69]:
        for i in range(6):b('measurement ruled line',(x,2.31+i*.075,1.224),(.32,.009,.003),iron)
    for i,mat in enumerate([cloth,navy,grey,wine]):b('cloth swatch',(1.97+i*.08,2.65,1.19+i*.012),(.36,.42,.018),mat)
    # Small mechanical finishing machine; the main cutting room is upstairs.
    b('machine bed',(0,2.64,1.215),(.97,.53,.11),iron,.035)
    b('machine pillar',(.29,2.72,1.45),(.17,.20,.44),iron,.045)
    b('machine arm',(.01,2.72,1.68),(.71,.21,.15),iron,.045)
    b('needle housing',(-.31,2.72,1.56),(.15,.16,.27),iron,.045)
    c('needle',(-.31,2.72,1.365),.008,.19,brass,vertices=12)
    c('machine wheel',(.47,2.73,1.59),.16,.055,iron,(0,math.pi/2,0),32)
    c('wheel hub',(.505,2.73,1.59),.055,.025,brass,(0,math.pi/2,0),24)
    c('thread spool',(.16,2.72,1.84),.045,.13,paper)
    b('work under needle',(-.30,2.55,1.28),(.44,.55,.018),cream)
    b('pattern paper',(-1.51,2.56,1.17),(1.0,.66,.018),paper)
    for x in [-1.72,-1.34]:
        blade=b('tailor shear blade',(x,2.58,1.189),(.045,.38,.016),brass,.008);blade.rotation_euler.z=.24 if x<-1.5 else -.24
        c('shear handle',(x,2.36,1.192),.067,.025,iron,vertices=24)
    b('measuring tape',(-1.51,2.80,1.188),(.86,.055,.015),paper)
    for i in range(17):b('tape measure mark',(-1.91+i*.05,2.796,1.198),(.008,.038,.003),iron)
    # Folded bolts on rear shelves, set behind the staff aisle.
    for z in [.72,1.30,1.88,2.46]:b('fabric shelf',(-.2,4.57,z),(5.30,.55,.08),trim,.015,back)
    for x in [-2.82,2.42]:b('shelf upright',(x,4.61,1.60),(.09,.42,2.00),oak,parent=back)
    for row in range(3):
        for col in range(7):
            x=-2.39+col*.73
            b('folded wool',(x,4.57,.90+row*.58),(.62,.42,.24),[cloth,navy,grey,wine][(row+col)%4],.035,back)
            b('bolt label',(x,4.346,.89+row*.58),(.17,.018,.07),paper,parent=back)
    b('house name board',(-.2,4.73,3.25),(5.4,.08,.62),oak,.025,back)
    b('house name inset',(-.2,4.679,3.25),(5.18,.025,.45),paper,parent=back)
    label('house lettering','RUTTLEDGE & VANCE',(-.2,4.657,3.16),.27)
    # Full-length fitting mirror and curtain screen, visible from the public floor.
    b('mirror surround',(-3.83,4.74,1.65),(1.64,.18,2.75),trim,.055,back)
    b('mirror glass',(-3.83,4.63,1.65),(1.38,.035,2.49),mirror,.025,back)
    for y in [2.4,4.35]:c('screen upright',(-2.92,y,1.33),.038,2.60,brass)
    c('screen curtain rail',(-2.92,3.375,2.61),.035,2.03,brass,(math.pi/2,0,0))
    for j in range(16):b('curtain fold',(-2.92+(.045 if j%2 else -.045),2.45+j*.12,1.55),(.08,.13,2.02),wine,.012)
    # Garment rail against the side wall.
    for y in [-1.6,1.1]:c('garment rail post',(-4.35,y,1.20),.035,2.35,brass)
    c('garment rail',(-4.35,-.25,2.38),.035,2.8,brass,(math.pi/2,0,0))
    for j in range(8):
        y=-1.4+j*.34
        b('hanging folded cloth',(-4.32,y,1.64),(.54,.16,1.25),[navy,grey,cloth][j%3],.05)
    # Three headless display forms distinguish the salon from a crowd of NPCs.
    for i,y in enumerate([-3.2,-1.8,-.4]):
        x=4.10;mat=[navy,grey,cloth][i]
        c('display form base',(x,y,.07),.34,.11,iron,vertices=32)
        c('display form stand',(x,y,.66),.033,1.15,brass)
        b('jacket body',(x,y,1.36),(.55,.30,.80),mat,.10)
        for side in [-1,1]:
            sleeve=b('jacket sleeve',(x+side*.33,y,1.36),(.17,.27,.69),mat,.075);sleeve.rotation_euler.y=side*.17
            lapel=b('jacket lapel',(x+side*.13,y-.167,1.54),(.13,.022,.40),trim,.015);lapel.rotation_euler.y=-side*.27
        c('display neck finial',(x,y,1.81),.07,.10,brass)
        for z in [1.24,1.38]:c('jacket button',(x+.05,y-.17,z),.018,.012,brass,(math.pi/2,0,0),12)
    # A client chair in the front corner.
    b('client chair base',(-3.8,-3.2,.29),(.72,.72,.49),oak,.025)
    b('client chair cushion',(-3.8,-3.2,.61),(.76,.76,.16),cloth,.055)
    b('client chair back',(-4.23,-3.2,.98),(.14,.79,.79),oak,.035)
    b('client chair back cushion',(-4.137,-3.2,1.0),(.06,.65,.58),cloth,.03)
    # Stair to the cutting room remains outside the visitor and counter aisles.
    for step in range(11):
        y=.92+step*.35;h=(step+1)*.25
        b('cutting room stair',(3.68,y,h/2),(1.72,.35,h),trim,.012)
        if step%2==0:
            for x in [2.87,4.49]:c('stair baluster',(x,y,h+.40),.025,.80,brass,vertices=16)
    for x in [2.87,4.49]:
        # A joined handrail with its angle matching the rise.
        rail=b('stair rail',(x,2.66,2.25),(.065,4.4,.065),oak,.02);rail.rotation_euler.x=math.atan2(2.5,3.5)
    b('upper landing',(3.68,4.72,2.68),(1.72,.50,.14),trim,.012)
    b('upper doorway wall',(3.68,5,4.0),(1.92,.16,2.6),cream,parent=back)
    b('upstairs door',(3.68,4.89,3.85),(1.52,.10,2.20),oak,.025,back)
    b('upstairs door lintel',(3.68,4.79,5.01),(1.78,.16,.12),trim,.012,back)
    label('cutting room sign','CUTTING ROOM',(3.68,4.828,4.57),.12)
    for x,y in [(-2,.6),(1.3,.6)]:
        c('pendant stem',(x,y,3.30),.016,.66,brass)
        c('opal pendant',(x,y,2.95),.27,.13,glass,vertices=32)
        c('pendant lip',(x,y,2.88),.28,.03,brass,vertices=32)
