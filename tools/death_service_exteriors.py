"""Original period death-service premises, Blender metres/Z-up."""
import bpy
import math


def build(kind, box, cylinder, material, brick, roof_texture):
    masonry=material(kind+' russet masonry',(.32,.21,.16));brick(masonry,91)
    stone=material(kind+' dressed limestone',(.58,.55,.44))
    slate=material(kind+' slate',(.11,.15,.15));roof_texture(slate,True)
    oak=material(kind+' painted oak',(.13,.18,.14))
    glass=material(kind+' muted glazing',(.13,.23,.24),.4)
    iron=material(kind+' iron',(.045,.055,.05),.6)
    lawn=material(kind+' lawn',(.21,.29,.15))
    gravel=material(kind+' gravel',(.47,.43,.33))
    brass=material(kind+' brass',(.58,.43,.18),.7)
    def marker(name,p,scale=None):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob);ob.location=p
        if scale:ob.scale=scale
        return ob
    def window(x,y,z,w=1.2,h=1.5):
        box('window stone surround',(x,y,z),(w+.22,.18,h+.22),stone,.02)
        box('recessed glass',(x,y+.11,z),(w,.045,h),glass)
        for dx in [-w/2,0,w/2]:box('sash stile',(x+dx,y+.15,z),(.055,.065,h+.06),oak)
        box('sash rail',(x,y+.15,z),(w,.065,.065),oak)
        box('projecting sill',(x,y+.16,z-h/2-.09),(w+.35,.32,.12),stone,.02)
        marker('fire-window-'+str(x),(x,y+.2,z))
    def door(x,y):
        marker('entrance-threshold',(x,y,0))
        hinge=marker('entrance-door-hinge',(x-.7,y,0))
        for name,p,d,m in [('oak entry',(x,y,1.3),(1.4,.16,2.6),oak),('entry glass',(x,y+.1,1.85),(1.05,.06,.95),glass),('door brass pull',(x+.48,y+.16,1.1),(.06,.06,.32),brass)]:
            ob=box(name,p,d,m,.02);ob.parent=hinge;ob.location-=hinge.location
        box('door lintel',(x,y,2.77),(1.8,.35,.24),stone,.025)
    box('lot paving',(0,0,-.055),(15.6,15.6,.1),gravel)
    if kind=='cemetery':
        box('burial lawn',(0,-1,.02),(14.8,12.8,.08),lawn)
        box('central gravel walk',(0,0,.08),(2.1,15.4,.08),gravel)
        box('cross path',(0,-1,.086),(14.9,1.3,.08),gravel)
        for side in [-1,1]:
            for y in [-5,-3,1,3]:
                for x in [3.0,5.4]:
                    xx=x*side
                    box('grave stone plinth',(xx,y,.15),(.9,.58,.18),stone,.03)
                    box('rounded headstone',(xx,y,.68),(.68,.22,.95),stone,.10)
                    box('inscribed stone face',(xx,y+.13,.73),(.46,.025,.48),slate,.03)
        # Low caretaker lodge leaves the central entrance and paths unobstructed.
        box('caretaker lodge',(-4.6,5.2,1.65),(4.8,3.6,3.3),masonry)
        box('lodge roof',(-4.6,5.2,3.38),(5.2,4,.22),slate,.035)
        window(-5.5,7.05,1.9,1.1,1.35)
        for x in [-7.3,7.3]:
            box('boundary base',(x,0,.22),(.25,15,.44),stone)
            for y in range(-7,8):
                cylinder('iron boundary rail',(x,y,1.0),.034,1.55,iron)
            box('iron top rail',(x,0,1.64),(.06,15,.06),iron)
        for x in [-1.65,1.65]:
            box('gate pier',(x,7.15,1.15),(.55,.6,2.3),masonry)
            box('gate cap',(x,7.15,2.37),(.72,.76,.18),stone,.025)
        marker('entrance-threshold',(0,7.5,0))
        marker('sign-anchor',(0,7.25,2.6),(.6,1,.25))
    else:
        width=12 if kind=='mortuary' else 10
        # Central recessed doorway is a real void, rather than a door pasted on a solid wall.
        for side in [-1,1]:box('masonry wing',(side*(width/4+.45),1.4,2.3),(width/2-.9,10.2,4.6),masonry)
        box('entry rear',(0,-.1,2.3),(1.8,7.2,4.6),masonry)
        box('entry overdoor',(0,5.1,3.72),(1.8,2.8,1.76),masonry)
        box('limestone roof cornice',(0,1.4,4.64),(width+.35,10.55,.25),stone,.025)
        box('slate flat roof',(0,1.4,4.83),(width+.12,10.3,.16),slate)
        for x in [-width/2,width/2]:box('corner pier',(x,6.57,2.3),(.32,.25,4.6),stone,.02)
        for x in [-3.5,3.5]:window(x,6.55,2.1,2.5,1.9)
        door(0,6.5)
        marker('sign-anchor',(0,6.68,3.5),(.8,1,.28))
        if kind=='mortuary':
            box('receiving canopy',(-3.5,7.0,3.4),(4.2,1.4,.16),oak,.03)
            for x in [-5.4,-1.6]:cylinder('canopy support',(x,7.4,1.65),.045,3.3,iron)
            for x in [-3.5,0,3.5]:box('roof ventilation',(x,-1,5.15),(.9,1.2,.5),slate,.025)
        else:
            box('furnace chimney',(3,-2.5,4.8),(1.3,1.5,9.6),masonry)
            for z in [6.9,9.2]:box('chimney stone collar',(3,-2.5,z),(1.52,1.72,.24),stone)
            box('chimney dark mouth',(3,-2.5,9.62),(.85,1.05,.05),iron)
            for x in [-6.6,6.6]:
                cylinder('memorial planter',(x,4,.34),.55,.68,stone)
                cylinder('clipped evergreen',(x,4,1.1),.46,1.1,lawn)
