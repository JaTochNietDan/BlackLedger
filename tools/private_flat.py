"""Locally authored 1950s bedsit; the resident's private rooms, not the lobby."""
import math
import bpy


def build(box,cylinder,material):
    wood=material('Flat varnished oak',(.31,.18,.085))
    edge=material('Flat dark wood',(.13,.075,.035))
    plaster=material('Flat ivory plaster',(.72,.66,.52))
    blue=material('Flat blue upholstery',(.16,.26,.28))
    linen=material('Flat bed linen',(.79,.74,.59))
    quilt=material('Flat woven blanket',(.34,.37,.23))
    cream=material('Flat enamel',(.73,.72,.61))
    chrome=material('Flat nickel',(.48,.49,.43),.75)
    black=material('Flat iron',(.035,.04,.034),.3)
    glass=material('Flat daylight glass',(.38,.49,.48),0,.15)
    curtain=material('Flat striped curtains',(.50,.31,.14))
    paper=material('Flat pages',(.8,.74,.58))
    rug=material('Flat woven rug',(.37,.15,.105))
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    b('floor base',(0,0,-.08),(8,7,.18),edge)
    for row in range(28):
        for col in range(5):
            low=max(-4,-4.8+col*1.9+(row%2)*.9);high=min(4,-4.8+(col+1)*1.9+(row%2)*.9)
            if high>low:b('oak floorboard',((low+high)/2,-3.375+row*.25,.0135),(high-low-.01,.244,.008),wood)
    b('left plaster',(-4,0,1.65),(.14,7,3.3),plaster,parent=left)
    b('back plaster',(0,3.5,1.65),(8,.14,3.3),plaster,parent=back)
    for z,h in [(.13,.22),(3.18,.14)]:
        b('left molding',(-3.88,0,z),(.13,7,h),linen,parent=left)
        b('back molding',(0,3.38,z),(8,.13,h),linen,parent=back)
    # Sleeping alcove with a full mattress, folded blanket, pillows and bedside radio.
    b('bed frame',(-2.55,1.6,.3),(1.68,2.45,.32),edge,.035)
    b('mattress',(-2.55,1.58,.56),(1.59,2.3,.25),linen,.1)
    b('headboard',(-2.55,2.88,.65),(1.8,.12,1.25),wood,.035)
    b('blanket',(-2.55,1.1,.72),(1.60,1.26,.06),quilt,.025)
    for x in [-3,-2.15]:b('bed pillow',(x,2.36,.73),(.65,.45,.12),linen,.09)
    for x in [-3.22,-1.88]:
        for y in [.58,2.58]:b('bed leg',(x,y,.12),(.08,.08,.22),chrome)
    from residential_storage import bedside
    bedside(box,-1.25,2.65,.68,.64,.78,wood,edge,chrome)
    b('radio',(-1.25,2.65,.99),(.47,.25,.31),edge,.035)
    for x in [-1.4,-1.32,-1.24]:b('radio grille',(x,2.515,1),(.018,.015,.18),chrome)
    # Window and folded curtains beside the bed.
    b('window surround',(-2.5,3.35,2.15),(2.3,.15,1.65),linen,.02,back)
    b('window panes',(-2.5,3.24,2.15),(2.1,.025,1.45),glass,parent=back)
    for x in [-3.5,-2.5,-1.5]:b('sash bar',(x,3.20,2.15),(.035,.04,1.47),wood,parent=back)
    b('sash crossbar',(-2.5,3.19,2.15),(2.1,.04,.035),wood,parent=back)
    b('window sill',(-2.5,3.13,1.3),(2.4,.3,.08),linen,.02,back)
    for side in [-1,1]:
        for fold in range(5):
            o=cylinder('curtain fold',(-2.5+side*(1.13+fold*.055),3.10,2.13),.055,1.83,curtain,vertices=10);o.parent=back
    # Kitchen with enamel stove, burners, oven, sink and refrigerator.
    b('kitchen linoleum',(2.02,2.52,.022),(3.8,1.65,.014),quilt)
    for x in [.75,1.75]:
        b('kitchen cupboard',(x,2.9,.49),(.95,.87,.94),cream,.025)
        b('cabinet door',(x,2.44,.46),(.84,.035,.77),linen,.012)
        b('cabinet handle',(x+.30,2.4,.71),(.025,.04,.18),chrome,.01)
    b('countertop',(1.25,2.85,1.0),(2.02,1.04,.09),edge,.035)
    b('sink rim',(.75,2.79,1.056),(.68,.55,.025),chrome,.055)
    b('sink dark basin',(.75,2.79,1.072),(.52,.39,.008),black,.055)
    cylinder('tap stem',(.75,3.17,1.21),.025,.32,chrome,vertices=16)
    b('tap spout',(.75,3.07,1.37),(.055,.23,.05),chrome,.025)
    b('stove',(2.7,2.92,.48),(.81,.9,.93),cream,.06)
    b('oven glass',(2.7,2.451,.44),(.60,.025,.43),black,.025)
    b('oven rail',(2.7,2.4,.70),(.53,.045,.035),chrome,.015)
    for x in [2.48,2.92]:
        for y in [2.7,3.1]:
            cylinder('stove burner',(x,y,.96),.13,.025,black,vertices=24)
        cylinder('stove knob',(x,2.43,.82),.045,.04,black,(math.pi/2,0,0),16)
    b('refrigerator',(3.55,2.9,.91),(.76,.91,1.78),cream,.10)
    b('fridge door',(3.55,2.41,.94),(.67,.09,1.6),linen,.075)
    b('fridge handle',(3.27,2.33,1.19),(.045,.06,.36),chrome,.02)
    # Crockery shelves and a kettle, no invented contraband or fitted safe.
    for z in [1.82,2.42]:b('kitchen shelf',(1.25,3.21,z),(2.0,.40,.07),wood,.015,back)
    for x in [.55,.9,1.25]:
        cylinder('cup',(x,3.2,1.94),.065,.18,linen,vertices=16).parent=back
    cylinder('kettle',(1.72,2.86,1.17),.115,.22,chrome,vertices=24)
    # Living corner separated from the clear front entry aisle.
    b('sofa base',(-3.1,-1.42,.37),(1.1,1.9,.42),blue,.10)
    b('sofa back',(-3.58,-1.42,.88),(.22,1.95,.94),blue,.09)
    for y in [-1.88,-.96]:b('sofa cushion',(-3.04,y,.64),(.86,.86,.14),blue,.05)
    for y in [-2.42,-.42]:b('sofa arm',(-3.07,y,.79),(1.13,.17,.36),blue,.075)
    b('rug',(-1.49,-1.4,.031),(1.78,2.0,.027),rug)
    b('coffee table',(-1.5,-1.4,.55),(.75,1.22,.08),wood,.025)
    for x in [-1.78,-1.22]:
        for y in [-1.88,-.92]:b('coffee leg',(x,y,.28),(.06,.06,.52),edge)
    b('book',(-1.5,-1.4,.615),(.35,.43,.045),paper,.01)
    b('dining table',(1.5,.28,.80),(1.38,.96,.075),wood,.03)
    for x in [1,2]:
        for y in [-.06,.62]:b('table leg',(x,y,.4),(.06,.06,.77),chrome)
    for x in [1.1,1.9]:
        cylinder('dinner plate',(x,.28,.85),.17,.02,linen,vertices=24)
    for x in [1.08,1.92]:
        b('dining chair seat',(x,-.68,.58),(.58,.54,.07),wood,.025)
        b('dining chair back',(x,-.93,.94),(.56,.07,.60),wood,.03)
        for dx in [-.22,.22]:
            for y in [-.9,-.46]:b('chair foot',(x+dx,y,.28),(.045,.045,.54),chrome)
    # Bathroom door indicates the rest of the dwelling beyond the cutaway.
    b('bathroom architrave',(-3.87,-2.95,1.1),(.13,.88,2.15),wood,parent=left)
    b('bathroom door',(-3.77,-2.95,1.09),(.06,.75,2.02),linen,parent=left)
    b('bathroom latch',(-3.72,-2.7,1),(.05,.035,.13),chrome,parent=left)
