"""Original Mercer Exchange message hall and reading room; Blender metres, Z up."""
import math
import bpy


def build(box,cylinder,material):
    stone=material('Mercer cream limestone',(.61,.57,.46))
    plaster=material('Mercer muted plaster',(.58,.59,.48))
    oak=material('Mercer dark oak',(.19,.105,.045))
    edge=material('Mercer polished oak',(.34,.22,.105))
    brass=material('Mercer brass fittings',(.51,.37,.15),.6)
    green=material('Mercer leather writing surface',(.07,.16,.12))
    paper=material('Mercer envelopes and newsprint',(.80,.74,.59))
    ink=material('Mercer printer ink',(.035,.043,.035))
    cork=material('Mercer noticeboard',(.40,.28,.14))
    globe=material('Mercer opal globes',(.88,.79,.58),0,.55)
    floor=material('Mercer stone floor',(.43,.46,.40))
    def group(name):
        o=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(o);return o
    left,back=group('interior-wall-left'),group('interior-wall-back')
    def b(name,p,d,m,bevel=0,parent=None):
        o=box(name,p,d,m,bevel);o.parent=parent;return o
    def c(name,p,r,d,m,rotation=(0,0,0),vertices=24,parent=None):
        o=cylinder(name,p,r,d,m,rotation,vertices);o.parent=parent;return o
    def text(name,words,p,size,parent=back):
        cu=bpy.data.curves.new(name,'FONT');cu.body=words;cu.size=size;cu.align_x='CENTER';cu.extrude=.001
        o=bpy.data.objects.new(name,cu);bpy.context.collection.objects.link(o);o.location=p;o.rotation_euler=(math.pi/2,0,0);o.data.materials.append(ink);o.parent=parent
    b('hall foundation',(0,0,-.09),(10,10,.20),oak)
    for x in range(20):
        for y in range(20):b('floor tile',(-4.75+x*.5,-4.75+y*.5,.011),(.488,.488,.013),floor if (x+y)%2 else stone)
    b('rear plaster',(0,5,2),(10,.16,4),plaster,parent=back)
    b('side plaster',(-5,0,2),(.16,10,4),plaster,parent=left)
    for z,h in [(.12,.2),(1.1,.08),(3.75,.14),(3.9,.08)]:
        b('rear cornice',(0,4.85,z),(10,.18,h),stone,.012,back)
        b('side cornice',(-4.85,0,z),(.18,10,h),stone,.012,left)
    for x in [-4.7,-.7,4.7]:
        b('rear pilaster',(x,4.77,2),(.32,.24,3.70),stone,parent=back)
        b('pilaster capital',(x,4.72,3.6),(.53,.34,.18),stone,.02,back)
    # Three message windows; no decorative amounts pretend to be live prices.
    b('counter case',(.6,2.6,.52),(7.2,.95,1.02),oak,.025)
    b('counter polished ledge',(.6,2.55,1.075),(7.45,1.16,.09),edge,.035)
    for x in [-1.8,.6,3.0]:
        b('counter inset',(x,2.11,.54),(2.1,.035,.77),edge,.02)
        b('writing pad',(x,2.05,1.128),(1.03,.38,.012),green,.012)
        for xx in [x-1.1,x+1.1]:
            b('window upright',(xx,2.7,1.94),(.10,.15,1.72),brass,.014)
        b('window lintel',(x,2.7,2.80),(2.3,.17,.11),oak,.015)
        for j in range(13):
            xx=x-1.0+j/6
            b('window grille',(xx,2.7,2.12),(.017,.024,1.24),brass)
        for z in [1.51,2.05,2.72]:b('grille crossrail',(x,2.7,z),(2.1,.034,.022),brass)
        b('letter tray',(x,2.80,1.14),(.48,.38,.045),brass,.02)
        for k in range(4):b('waiting envelope',(x,2.80,1.173+k*.012),(.32,.22,.008),paper)
    for x,words in [(-1.8,'MESSAGES'),(.6,'ENQUIRIES'),(3.0,'NOTICES')]:
        b('window sign',(x,2.67,2.97),(1.74,.07,.26),paper,.01)
        text('window lettering',words,(x,2.625,2.9),.13,None)
    # Pigeonholes, folders and wire baskets behind the clerks.
    for cx in [-2.6,2.3]:
        b('archive back',(cx,4.86,1.82),(3.55,.08,2.0),oak,parent=back)
        for z in [.85,1.31,1.77,2.23,2.69]:b('archive shelf',(cx,4.56,z),(3.6,.60,.06),edge,parent=back)
        for i in range(7):b('archive divider',(cx-1.75+i*.583,4.56,1.77),(.035,.6,1.88),edge,parent=back)
        for row in range(4):
            for col in range(6):
                x=cx-1.46+col*.583
                b('filed papers',(x,4.59,.93+row*.46),(.45,.39,.12),paper,parent=back)
                b('folder label',(x,4.383,1.05+row*.46),(.17,.014,.06),paper,parent=back)
    b('exchange plaque frame',(.1,4.53,3.38),(5.05,.10,.62),oak,.025,back)
    b('exchange plaque',(.1,4.466,3.38),(4.91,.025,.48),paper,parent=back)
    text('exchange lettering','MERCER EXCHANGE',(.1,4.44,3.26),.34)
    # Long reading table and facing upholstered benches.
    b('reading tabletop',(-2.5,-.5,1.04),(1.5,3.0,.12),edge,.035)
    b('reading leather',(-2.5,-.5,1.107),(1.24,2.72,.014),green,.02)
    for x in [-3.02,-1.98]:
        for y in [-1.68,.68]:b('reading table leg',(x,y,.51),(.11,.11,1.00),oak,.01)
    for x,side in [(-3.6,-1),(-1.4,1)]:
        b('bench plinth',(x,-.5,.3),(.72,3.1,.48),oak,.025)
        b('bench seat',(x,-.5,.61),(.75,3.1,.16),green,.045)
        b('bench back',(x+side*.40,-.5,.94),(.15,3.15,.72),oak,.025)
        for y in [-1.65,-.9,-.1,.65]:b('bench back inset',(x+side*.30,y,.99),(.025,.58,.43),green,.018)
    for y in [-1.3,.3]:
        b('open newspaper',(-2.5,y,1.123),(1.0,.61,.014),paper,.005)
        for column in range(4):
            for row in range(6):b('newspaper column',(-2.87+column*.25,y-.22+row*.073,1.132),(.19,.015,.003),ink)
    # Side noticeboards retain blank forms: actual offers remain in the UI.
    for y in [-2.8,.2,3.2]:
        b('noticeboard frame',(-4.80,y,2.24),(.10,2.25,1.38),oak,.02,left)
        b('notice cork',(-4.735,y,2.24),(.035,2.08,1.20),cork,parent=left)
        for dy,z in [(-.58,2.34),(0,2.20),(.58,2.37)]:
            b('pinned notice',(-4.708,y+dy,z),(.012,.43,.61),paper,parent=left)
            c('notice pin',(-4.692,y+dy,z+.25),.018,.012,brass,(0,math.pi/2,0),12,left)
            for j in range(5):b('notice ruled line',(-4.694,y+dy,z+.12-j*.075),(.005,.31,.01),ink,parent=left)
    # A public writing stand in the east corner leaves the entrance aisle free.
    b('forms desk',(3.35,-2.2,1.04),(1.65,.70,.10),edge,.02)
    for x in [2.68,4.02]:b('forms desk leg',(x,-2.2,.51),(.11,.48,1.02),oak,.01)
    for x in [3.0,3.65]:b('stack of forms',(x,-2.2,1.13),(.40,.29,.10),paper)
    for x,y in [(-2.5,-.5),(1.4,-.5),(1.4,2.6)]:
        c('ceiling pendant',(x,y,3.4),.018,.6,brass,vertices=12)
        bpy.ops.mesh.primitive_uv_sphere_add(segments=24,ring_count=12,radius=.22,location=(x,y,3.03))
        bpy.context.object.name='opal globe';bpy.context.object.data.materials.append(globe)
        c('globe cap',(x,y,3.24),.095,.045,brass)
