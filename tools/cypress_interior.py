"""Locally authored Cypress House drawing room; Blender metres, Z up.

Invoked by export_city3d.py. No downloaded models or textures.
"""
import math
import random
import bpy


def build(box, cylinder, material):
    rng = random.Random(1957)
    oak = material('Cypress smoked walnut', (.19, .105, .055))
    trim = material('Cypress polished edges', (.31, .19, .09))
    plaster = material('Cypress warm plaster', (.66, .61, .47))
    green = material('Cypress sage wallpaper', (.27, .34, .25))
    brass = material('Cypress aged brass', (.54, .37, .13), .7)
    stone = material('Cypress limestone', (.63, .60, .51))
    dark = material('Cypress fireplace iron', (.035, .031, .028), .4)
    leather = material('Cypress oxblood leather', (.23, .044, .035))
    velvet = material('Cypress ochre velvet', (.45, .29, .105))
    cream = material('Cypress linen', (.8, .71, .51))
    glass = material('Cypress dusk window', (.27, .39, .39), .1, .25)
    rug = material('Cypress wine rug', (.26, .075, .054))
    ink = material('Cypress green rug motif', (.15, .23, .18))
    paper = material('Cypress book pages', (.71, .65, .48))
    book_colors = [material('Cypress book ' + str(i), c) for i, c in enumerate([
        (.25,.08,.05), (.09,.19,.14), (.13,.14,.2), (.34,.24,.1)])]
    # Thin, individually staggered boards give the floor direction and age.
    boards = [material('Cypress parquet ' + str(i), (.24+i*.018,.135+i*.011,.07+i*.006)) for i in range(6)]
    box('floor foundation',(0,0,-.09),(10,9,.20),oak)
    for row in range(30):
        y = -4.35 + row*.3
        for col in range(7):
            lo=max(-5,-5.8+col*1.7+(row%2)*.85)
            hi=min(5,-5.8+(col+1)*1.7+(row%2)*.85)
            if hi>lo: box('parquet',((lo+hi)/2,y,.0075),(hi-lo-.008,.294,.02),rng.choice(boards))
    # Walls are independent groups so orbiting outside cuts them away.
    def wall_group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left=wall_group('interior-wall-left');back=wall_group('interior-wall-back')
    def b(name,xyz,dims,mat,bevel=0,parent=None):
        o=box(name,xyz,dims,mat,bevel);o.parent=parent;return o
    b('west wallpaper',(-5,0,2),( .16,9,4),green,parent=left)
    b('north wallpaper',(0,4.5,2),(10,.16,4),green,parent=back)
    for z,h in [(.14,.23),(1.03,.09),(3.86,.18),(3.65,.06)]:
        b('west molding',(-4.87,0,z),(.13,9,h),trim,parent=left)
        b('north molding',(0,4.37,z),(10,.13,h),trim,parent=back)
    for i in range(10):
        x=-4.5+i
        b('north wainscot',(x,4.38,.57),(.92,.1,.77),oak,parent=back)
        b('panel inset',(x,4.31,.57),(.70,.025,.52),trim,.012,back)
    for i in range(9):
        y=-4+i
        b('west wainscot',(-4.87,y,.57),(.1,.92,.77),oak,parent=left)
        b('panel inset',(-4.80,y,.57),(.025,.70,.52),trim,.012,left)
    # Two sash windows and pleated curtains flank the central fireplace.
    for x in [-3.25,3.25]:
        b('window frame',(x,4.28,2.4),(2.05,.16,2.35),cream,.025,back)
        b('window glass',(x,4.17,2.4),(1.83,.025,2.13),glass,parent=back)
        for dx in [-.9,0,.9]: b('window mullion',(x+dx,4.12,2.4),(.045,.05,2.16),trim,parent=back)
        for z in [1.35,2.4,3.45]:b('window sash',(x,4.1,z),(1.83,.06,.045),trim,parent=back)
        b('window sill',(x,4.02,1.3),(2.18,.36,.10),cream,.02,back)
        for side in [-1,1]:
            for fold in range(6):
                o=cylinder('curtain pleat',(x+side*(.96+fold*.07),3.99,2.30),.07,2.5,velvet,vertices=10);o.parent=back
        b('curtain pelmet',(x,3.94,3.61),(2.8,.25,.2),velvet,.06,back)
    # Hearth, soot-dark opening, and a quiet grate (no unsimulated fire).
    b('hearth slab',(0,3.75,.085),(2.6,1.4,.14),stone,.025)
    b('fireplace cavity',(0,4.22,.67),(1.45,.16,1.18),dark,parent=back)
    for x in [-1,1]:b('mantel pier',(x,4.03,.77),(.36,.58,1.5),stone,.025,back)
    b('mantel lintel',(0,4.01,1.47),(2.35,.6,.3),stone,.025,back)
    b('mantel shelf',(0,3.98,1.68),(2.65,.72,.13),trim,.025,back)
    for x in [-.6,-.4,-.2,0,.2,.4,.6]:b('grate bar',(x,3.98,.31),(.025,.26,.4),dark,parent=back)
    b('mirror frame',(0,4.32,2.62),(1.72,.10,1.42),brass,.035,back)
    b('aged mirror',(0,4.25,2.62),(1.53,.02,1.23),glass,parent=back)
    b('mantel clock case',(0,3.96,1.99),(.55,.21,.49),oak,.07,back)
    o=cylinder('clock face',(0,3.835,2.02),.18,.025,cream,(math.pi/2,0,0),32);o.parent=back
    b('clock hand',(0,3.815,2.08),(.015,.01,.13),dark,parent=back)
    b('clock hand',(.065,3.815,2.02),(.13,.01,.015),dark,parent=back)
    # Rug and modest geometric woven border, away from the entrance aisle.
    b('Persian rug',(-1.3,.3,.03),(5,4.4,.025),rug)
    for x in [-3.67,1.07]:b('rug binding',(x,.3,.045),(.1,4.25,.006),cream)
    for y in [-1.77,2.37]:b('rug binding',(-1.3,y,.045),(4.8,.1,.006),cream)
    for y in [-1.5,2.1]:
        for x in [-3.2,-2.5,-1.8,-1.1,-.4,.3]:
            o=b('woven diamond',(x,y,.05),(.20,.20,.006),ink);o.rotation_euler.z=math.pi/4
    # West sofa faces east: seat height is .69m for existing rig poses.
    b('sofa base',(-3.65,.3,.36),(1.12,3,.42),leather,.12)
    b('sofa back',(-4.14,.3,.88),(.25,3,1.02),leather,.1)
    for y in [-.65,.3,1.25]:
        b('sofa cushion',(-3.58,y,.61),(.91,.91,.16),leather,.07)
        for yy in [-.23,.23]: b('tuft brass',(-3.995,y+yy,1.04),(.02,.025,.025),brass,.006)
    for y in [-1.22,1.82]:b('rolled sofa arm',(-3.65,y,.75),(1.17,.22,.36),leather,.095)
    # East armchairs face west, each with a clear seat and exposed wooden feet.
    for y in [-.8,1.35]:
        for x in [.55,1.25]:
            for dy in [-.37,.37]:b('chair leg',(x,y+dy,.18),(.085,.085,.32),trim,.015)
        b('chair cushion',(.86,y,.61),(.87,.86,.16),velvet,.07)
        b('chair back',(1.34,y,.99),(.20,1.05,.90),velvet,.10)
        for dy in [-.5,.5]:b('chair arm',(.92,y+dy,.81),(1.08,.17,.3),velvet,.065)
    b('coffee tabletop',(-1.45,.3,.54),(1.45,1.95,.09),trim,.055)
    for x in [-2,-.9]:
        for y in [-.45,1.05]:b('coffee leg',(x,y,.27),(.075,.075,.5),oak,.01)
    b('folded newspaper',(-1.45,.3,.6),(.52,.36,.028),paper)
    for y in [.21,.26,.31,.36]:b('news lines',(-1.45,y,.616),(.36,.012,.002),dark)
    # Reading lamps, bookcase, and working study desk in the east alcove.
    for x,y in [(-3.7,2.7),(1,2.7)]:
        cylinder('lamp foot',(x,y,.05),.26,.075,brass,vertices=24)
        cylinder('lamp stem',(x,y,.94),.025,1.8,brass,vertices=12)
        bpy.ops.mesh.primitive_cone_add(vertices=32,radius1=.35,radius2=.19,depth=.43,location=(x,y,1.98))
        bpy.context.object.name='pleated reading shade';bpy.context.object.data.materials.append(cream)
    b('bookcase back',(-4.75,2.94,1.46),(.22,2.4,2.84),oak,parent=left)
    for z in [.18,.76,1.36,1.96,2.56,2.89]:b('bookcase shelf',(-4.51,2.94,z),(.52,2.4,.075),trim,parent=left)
    for z in [.78,1.38,1.98,2.58]:
        for i in range(18):
            y=1.83+i*.123;h=rng.uniform(.22,.43)
            b('leather volume',(-4.46,y,z+h/2),(.34,.095,h),rng.choice(book_colors),.005,left)
            b('book spine band',(-4.275,y,z+h*.78),(.012,.095,.012),brass,parent=left)
    b('study desk top',(3.35,1.7,.82),(2.2,1,.12),trim,.035)
    from residential_storage import bedside
    bedside(box,2.5,1.7,.45,.86,.74,oak,trim,brass)
    for x in [4.2]:
        b('desk pedestal',(x,1.7,.41),(.45,.86,.74),oak,.02)
        for z in [.2,.43,.66]:
            b('drawer front',(x,1.24,z),(.40,.045,.20),trim,.012)
            b('drawer pull',(x,1.2,z),(.14,.025,.025),brass,.009)
    b('blotter',(3.35,1.61,.888),(.9,.6,.01),ink)
    b('letter',(3.3,1.6,.902),(.30,.39,.012),paper)
    b('radio cabinet',(4.06,1.87,1.10),(.52,.35,.43),oak,.055)
    for x in [3.86,3.92,3.98,4.04,4.10]:b('radio grille',(x,1.68,1.12),(.018,.014,.23),brass)
    # A standing desk visitor occupies x3.35, Blender y.4 (glTF z-.4).
    # Front sideboard leaves the east entry approach unobstructed.
    b('sideboard',(3.7,-2.15,.55),(1.85,.55,1.05),oak,.025)
    b('sideboard top',(3.7,-2.15,1.1),(1.98,.65,.09),trim,.02)
    for x in [3.25,4.12]:
        b('sideboard door',(x,-2.45,.59),(.78,.045,.77),trim,.015)
        b('cabinet knob',(x,-2.49,.77),(.05,.04,.05),brass,.018)
    cylinder('silver tray',(3.7,-2.15,1.17),.3,.025,brass,vertices=32)
    cylinder('decanter',(3.6,-2.15,1.33),.08,.29,glass,vertices=12)
    cylinder('decanter neck',(3.6,-2.15,1.51),.03,.09,glass,vertices=12)
    for x in [3.85,4.04]:cylinder('tumbler',(x,-2.1,1.23),.045,.12,glass,vertices=16)
