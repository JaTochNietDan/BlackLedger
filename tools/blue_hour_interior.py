"""The Blue Hour casino floor. Original Blender geometry, Z-up metres."""
import math
import bpy


def build(box, cylinder, material):
    navy=material('Blue Hour midnight velvet',(.035,.09,.15))
    carpet=material('Blue Hour patterned blue carpet',(.055,.16,.21))
    walnut=material('Blue Hour polished walnut',(.20,.105,.05))
    brass=material('Blue Hour brass fittings',(.64,.43,.16),.7)
    cream=material('Blue Hour ivory plaster',(.77,.71,.55))
    felt=material('Blue Hour teal baize',(.025,.25,.22))
    black=material('Blue Hour bakelite',(.025,.027,.032))
    red=material('Blue Hour burgundy leather',(.30,.035,.05))
    white=material('Blue Hour ivory chips',(.9,.84,.65))
    glow=material('Blue Hour frosted glass',(.91,.77,.49),0,.6)
    def group(name):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);return ob
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        ob=box(name,p,d,m,bevel);ob.parent=parent;return ob
    def c(name,p,r,h,m,parent=None):
        ob=cylinder(name,p,r,h,m,(0,0,0),32);ob.parent=parent;return ob
    def label(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.002
        ob=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(ob);ob.location=p;ob.rotation_euler=(math.pi/2,0,0);cu.materials.append(brass);ob.parent=parent
    b('casino foundation',(0,0,-.10),(12,12,.20),walnut)
    b('casino carpet',(0,0,.0075),(11.8,11.8,.015),carpet)
    for x in range(-5,6):
        for y in range(-5,6):
            ob=b('carpet diamond',(x,y,.016),(.13,.13,.001),brass);ob.rotation_euler.z=math.pi/4
    b('rear wall',(0,6,2.1),(12,.16,4.2),navy,parent=back)
    b('left wall',(-6,0,2.1),(.16,12,4.2),navy,parent=left)
    for z,h in [(.18,.28),(1.14,.09),(3.97,.18)]:
        b('rear moulding',(0,5.87,z),(12,.12,h),walnut,.015,back)
        b('left moulding',(-5.87,0,z),(.12,12,h),walnut,.015,left)
    for x in [-5,-3,-1,1,3,5]:
        b('deco rear pier',(x,5.85,2.4),(.14,.15,2.8),brass,.015,back)
    for y in [-4,-1,2,5]:
        b('deco side pier',(-5.85,y,2.4),(.15,.14,2.8),brass,.015,left)
        b('wall light',(-5.65,y,2.9),(.22,.30,.60),glow,.04,left)
    label('casino title','THE BLUE HOUR',(0,5.78,3.25),.42)
    # Cashier cage, individual bars, grilled cash window and a closed safe.
    b('cashier cabinet',(3.9,4.25,.58),(3.0,1.1,1.16),walnut,.04)
    b('cashier marble ledge',(3.9,4.25,1.20),(3.2,1.10,.12),cream,.04)
    for x in [2.45+i*.24 for i in range(13)]:
        c('cashier grille',(x,4.2,1.95),.018,1.4,brass)
    b('cage header',(3.9,4.2,2.7),(3.2,.10,.12),brass,.01)
    b('cash tray',(3.9,3.95,1.29),(.65,.38,.035),black,.02)
    b('vault cabinet',(4.3,5.48,1.1),(1.4,.7,2.2),black,.05,back)
    c('vault dial',(4.3,5.05,1.25),.12,.06,brass,back).rotation_euler.x=math.pi/2
    label('cashier sign','CASHIER',(3.9,4.18,2.85),.21,None)
    # Three gaming surfaces leave a full-width central circulation aisle.
    for x,y,kind in [(-2.7,1.8,'roulette'),(2.7,.5,'cards'),(-2.7,-2.8,'cards')]:
        c(kind+' pedestal',(x,y,.46),.27,.9,walnut)
        top=c(kind+' padded rail',(x,y,.94),1.28,.18,walnut);top.scale.y=.70
        top=c(kind+' felt',(x,y,1.04),1.17,.035,felt);top.scale.y=.70
        if kind=='roulette':
            c('roulette bowl',(x-.48,y,1.12),.47,.13,walnut)
            c('roulette brass track',(x-.48,y,1.20),.40,.04,brass)
            c('roulette wheel',(x-.48,y,1.23),.32,.03,black)
            for i in range(24):
                a=i*math.pi/12
                ob=b('roulette pocket',(x-.48+.275*math.cos(a),y+.275*math.sin(a),1.255),(.055,.06,.015),red if i%2 else white);ob.rotation_euler.z=a
            c('roulette spindle',(x-.48,y,1.33),.05,.20,brass)
            for row in range(3):
                for col in range(4):b('roulette betting box',(x+.20+col*.18,y-.2+row*.18,1.065),(.16,.16,.005),red if (row+col)%2 else black)
        else:
            for dx in [-.6,0,.6]:
                b('empty betting mark',(x+dx,y-.23,1.066),(.27,.23,.004),brass,.015)
            b('dealer chip rack',(x,y+.36,1.11),(.90,.20,.10),walnut,.015)
            for i in range(6):
                for j in range(3):c('rack chip',(x-.35+i*.14,y+.36,1.17+j*.015),.055,.012,red if i%2 else white)
    # Slot cabinets and stools along the rear-left wall, away from the entry.
    for x in [-4.8,-3.6,-2.4]:
        b('slot walnut case',(x,4.65,1.12),(.86,.8,2.24),walnut,.05)
        b('slot brass face',(x,4.20,1.65),(.75,.07,.88),brass,.035)
        b('slot reel window',(x,4.14,1.66),(.65,.03,.31),black,.02)
        for dx in [-.21,0,.21]:
            b('slot ivory reel',(x+dx,4.10,1.66),(.17,.035,.25),white,.015)
            b('slot reel symbol',(x+dx,4.074,1.66),(.07,.009,.09),red,.02)
        b('slot payout tray',(x,4.08,.97),(.56,.25,.11),brass,.02)
        c('slot stool stem',(x,3.3,.30),.055,.6,brass)
        c('slot stool base',(x,3.3,.055),.29,.08,black)
        c('slot stool cushion',(x,3.3,.66),.30,.13,red)
    # Entry hostess stand and corner lounge, without fictitious occupants.
    b('hostess stand',(4.6,-4.4,.64),(1.05,.72,1.28),walnut,.04)
    b('hostess book',(4.6,-4.4,1.31),(.65,.40,.045),cream,.02)
    for y in [-3.7,-2.3]:
        b('lounge seat',(5.15,y,.54),(.85,1.1,.22),red,.10)
        b('lounge back',(5.53,y,.97),(.20,1.1,.85),red,.07)
    for x,y in [(-2.7,1.8),(2.7,.5),(-2.7,-2.8)]:
        c('pendant stem',(x,y,3.60),.025,.7,brass)
        c('pendant canopy',(x,y,3.23),.52,.10,brass)
        c('pendant milk glass',(x,y,3.16),.45,.10,glow)
