"""Reusable original bedside case with a physically sliding upper drawer."""
import bpy


def bedside(box,x,y,width,depth,height,wood,trim,metal):
    def b(name,p,d,m,bevel=.008,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    t=.045
    for xx in [x-width/2+t/2,x+width/2-t/2]:
        b('cabinet side',(xx,y,height/2),(t,depth,height),wood)
    for z in [.025,height,.43*height]:
        b('cabinet horizontal',(x,y,z),(width,depth,t),wood)
    b('cabinet rear',(x,y+depth/2-t/2,height/2),(width,t,height),wood)
    b('lower drawer face',(x,y-depth/2-.012,height*.22),(width-.065,.045,height*.36),trim)
    b('lower drawer pull',(x,y-depth/2-.05,height*.22),(.17,.03,.025),metal)
    drawer=bpy.data.objects.new('burglary-drawer',None);bpy.context.collection.objects.link(drawer)
    # Group origin remains zero, so +glTF Z translates the full tray outward.
    z=height*.72;tray_depth=depth-.09;front=y-depth/2-.012
    b('upper drawer face',(x,front,z),(width-.065,.045,height*.45),trim,parent=drawer)
    b('upper drawer pull',(x,front-.045,z),(.17,.03,.025),metal,parent=drawer)
    b('drawer bottom',(x,y-.015,height*.49),(width-.105,tray_depth,.025),wood,parent=drawer)
    for xx in [x-(width-.13)/2,x+(width-.13)/2]:
        b('drawer tray side',(xx,y-.015,height*.66),(.025,tray_depth,height*.32),wood,parent=drawer)
    b('drawer tray rear',(x,y+tray_depth/2-.015,height*.66),(width-.13,.025,height*.32),wood,parent=drawer)
    return drawer
