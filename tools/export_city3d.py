"""Author Bellwether's browser models in Blender. No external asset packs.
Run: .venv-blender/bin/python tools/export_city3d.py
Coordinates are metres, Blender Z up, glTF Y up. Export bounds include trim.
Materials use locally generated, deterministic brick UV textures, shared by
all facades. Sources and GLBs are reproducible; see models/manifest.json.
"""
import bpy
import json
import math
import os
import random
from mathutils import Vector

OUT = os.path.abspath('public/art/models')
os.makedirs(OUT, exist_ok=True)
random.seed(1950)


def material(name, color, metal=0, emission=0):
    mat = bpy.data.materials.new(name)
    mat.diffuse_color = (*color, 1)
    mat.use_nodes = True
    shader = mat.node_tree.nodes.get('Principled BSDF')
    shader.inputs['Base Color'].default_value = (*color, 1)
    shader.inputs['Roughness'].default_value = .72 if not metal else .28
    shader.inputs['Metallic'].default_value = metal
    if emission:
        shader.inputs['Emission Color'].default_value = (*color, 1)
        shader.inputs['Emission Strength'].default_value = emission
    return mat


def brick(mat, seed):
    rng = random.Random(seed)
    n = 256
    img = bpy.data.images.new(mat.name + '-brick', width=n, height=n)
    pixels = []
    heights = []
    base = mat.diffuse_color[:3]
    tones = [rng.uniform(.75, 1.15) for _ in range(64)]
    for y in range(n):
        row = y // 16
        for x in range(n):
            sx = (x + (row % 2) * 16) % n
            mortar = y % 16 < 2 or sx % 32 < 2
            tone = tones[(row * 8 + sx // 32) % 64] + rng.uniform(-.09, .09)
            c = (.23, .22, .19) if mortar else [v * tone for v in base]
            pixels.extend((*c, 1))
            heights.append(.1 if mortar else .8 + rng.uniform(-.04,.04))
    img.pixels = pixels
    img.pack()
    tree = mat.node_tree
    tex = tree.nodes.new('ShaderNodeTexImage')
    tex.image = img
    tree.links.new(tex.outputs['Color'], tree.nodes['Principled BSDF'].inputs['Base Color'])
    normal = bpy.data.images.new(mat.name+'-normal',width=n,height=n)
    normal.colorspace_settings.name='Non-Color'
    normal_pixels=[]
    for y in range(n):
        for x in range(n):
            dx=heights[y*n+(x-1)%n]-heights[y*n+(x+1)%n]
            dy=heights[((y-1)%n)*n+x]-heights[((y+1)%n)*n+x]
            v=Vector((dx,dy,1.8)).normalized()
            normal_pixels.extend((v.x*.5+.5,v.y*.5+.5,v.z*.5+.5,1))
    normal.pixels=normal_pixels
    normal.pack()
    normal_tex=tree.nodes.new('ShaderNodeTexImage');normal_tex.image=normal
    normal_map=tree.nodes.new('ShaderNodeNormalMap');normal_map.inputs['Strength'].default_value=.65
    tree.links.new(normal_tex.outputs['Color'],normal_map.inputs['Color'])
    tree.links.new(normal_map.outputs['Normal'],tree.nodes['Principled BSDF'].inputs['Normal'])


def box(name, xyz, dims, mat, bevel=0):
    bpy.ops.mesh.primitive_cube_add(size=1, location=xyz)
    ob = bpy.context.object
    ob.name = name
    ob.scale = dims
    bpy.ops.object.transform_apply(location=False, rotation=False, scale=True)
    ob.data.materials.append(mat)
    # Physical scale, matching mortar courses around every face.
    for poly in ob.data.polygons:
        axis = max(range(3), key=lambda i: abs(poly.normal[i]))
        axes = [i for i in range(3) if i != axis]
        for index in poly.loop_indices:
            co = ob.data.vertices[ob.data.loops[index].vertex_index].co
            ob.data.uv_layers.active.data[index].uv = (co[axes[0]] / 2.5, co[axes[1]] / 2.5)
    if bevel:
        mod = ob.modifiers.new('Rounded manufactured edges', 'BEVEL')
        mod.width = bevel
        mod.segments = 2
        bpy.context.view_layer.objects.active = ob
        bpy.ops.object.modifier_apply(modifier=mod.name)
    return ob


def cylinder(name, xyz, radius, depth, mat, rotation=(0, 0, 0), vertices=12):
    bpy.ops.mesh.primitive_cylinder_add(vertices=vertices, radius=radius, depth=depth, location=xyz, rotation=rotation)
    ob = bpy.context.object
    ob.name = name
    ob.data.materials.append(mat)
    return ob


def clear():
    bpy.ops.wm.read_factory_settings(use_empty=True)


def blast_fragment():
    clay=material('broken fired clay',(.46,.25,.14));brick(clay,1954)
    ob=box('chipped masonry fragment',(0,0,0),(.18,.09,.08),clay,.012)
    # A broken end, rather than a perfect miniature brick.
    for vertex in ob.data.vertices:
        if vertex.co.x>.065:
            vertex.co.x-=.009*(1+math.sin(vertex.co.y*91+vertex.co.z*73))


def rooftop_tank(height, iron):
    # Coopered timber cistern, strapped steel hoops and a braced rooftop stand.
    wood=material('weathered cistern cedar',(.58,.49,.35))
    n=256;pixels=[];normals=[]
    for y in range(n):
        for x in range(n):
            grain=math.sin(x*.69+math.sin(y*.024)*1.8)*.09+math.sin(x*2.17+y*.013)*.035
            tone=.91+grain+math.sin(x*.12)*.10
            pixels.extend((.58*tone,.49*tone,.35*tone,1))
            vec=Vector((math.cos(x*.69+math.sin(y*.024)*1.8)*.16,0,1)).normalized()
            normals.extend((vec.x*.5+.5,.5,vec.z*.5+.5,1))
    tree=wood.node_tree
    for name,data,is_normal in [('cistern cedar grain',pixels,False),('cistern cedar normal',normals,True)]:
        image=bpy.data.images.new(name,width=n,height=n);image.pixels=data
        if is_normal:image.colorspace_settings.name='Non-Color'
        image.pack();tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        if is_normal:
            normal=tree.nodes.new('ShaderNodeNormalMap');normal.inputs['Strength'].default_value=.4
            tree.links.new(tex.outputs['Color'],normal.inputs['Color'])
            tree.links.new(normal.outputs['Normal'],tree.nodes['Principled BSDF'].inputs['Normal'])
        else:tree.links.new(tex.outputs['Color'],tree.nodes['Principled BSDF'].inputs['Base Color'])
    for x in (1,3):
        for y in (-3,-1):
            box('tank stand leg',(x,y,height+.65),(.13,.13,1.3),iron)
        beam('tank stand cross brace',(x,-3,height+.12),(x,-1,height+1.25),.065,iron)
        beam('tank stand cross brace',(x,-1,height+.12),(x,-3,height+1.25),.065,iron)
    for y in (-3,-1):
        beam('tank stand cross brace',(1,y,height+.12),(3,y,height+1.25),.065,iron)
        beam('tank stand cross brace',(3,y,height+.12),(1,y,height+1.25),.065,iron)
        box('tank support girder',(2,y,height+1.25),(2.8,.16,.20),iron)
    cylinder('tank floor',(2,-2,height+1.34),1.22,.12,wood)
    for stave in range(32):
        angle=stave*math.tau/32
        ob=box('tank cedar stave',(2+1.2*math.cos(angle),-2+1.2*math.sin(angle),height+2.62),(.232,.1,2.5),wood,.009)
        ob.rotation_euler.z=angle+math.pi/2
    for z in (1.50,2.22,2.96,3.70):
        bpy.ops.mesh.primitive_torus_add(major_segments=32,minor_segments=6,location=(2,-2,height+z),major_radius=1.255,minor_radius=.038)
        ob=bpy.context.object;ob.name='tank iron hoop';ob.data.materials.append(iron)
    bpy.ops.mesh.primitive_cone_add(vertices=32,radius1=1.37,radius2=.08,depth=.62,location=(2,-2,height+4.14))
    ob=bpy.context.object;ob.name='tank conical cap';ob.data.materials.append(iron)
    cylinder('tank cap vent',(2,-2,height+4.53),.085,.2,iron)
    for x in (1.72,2.28):
        box('tank ladder upright',(x,-3.48,height+2.1),(.065,.065,4.2),iron)
    for rung in range(14):
        box('tank ladder rung',(2,-3.48,height+.2+rung*.28),(.56,.065,.065),iron)
    for z in (.45,1.15):
        beam('tank ladder bracket',(2,-3.48,height+z),(2,-3,height+z),.065,iron)


def building(kind, floors, width=12, depth=12, seed=0, palette=None, accent=None):
    palettes = [(0.43,.23,.15), (.51,.46,.35), (.31,.33,.26), (.37,.21,.18)]
    wall = material('weathered masonry', palette or palettes[seed % 4]); brick(wall, seed)
    stone = material('limestone', (.59,.55,.45))
    dark = material('tar roof', (.13,.14,.13))
    glass = material('smoked glass', (.12,.19,.21), .25)
    warm = material('occupied windows', (.78,.46,.17), 0, .35)
    iron = material('painted iron', (.12,.15,.14), .5)
    height = floors * 3.15
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor);anchor.location=(0,depth/2+.26,2.85)
    box('masonry', (0,0,height/2), (width,depth,height), wall)
    box('foundation', (0,0,.25), (width+.25,depth+.25,.5), stone)
    box('cornice lower', (0,0,height-.18), (width+.45,depth+.45,.18), stone)
    box('cornice cap', (0,0,height), (width+.65,depth+.65,.22), stone)
    box('roof', (0,0,height+.13), (width,depth,.08), dark)
    for side in (-1,1):
        box('parapet front', (0,side*(depth/2-.15),height+.4), (width,.3,.65),wall)
        box('parapet coping', (0,side*(depth/2-.15),height+.77), (width+.2,.45,.12),stone)
        box('parapet side', (side*(width/2-.15),0,height+.4), (.3,depth,.65),wall)
        box('side coping', (side*(width/2-.15),0,height+.77), (.45,depth+.2,.12),stone)
    for seam in range(int(width)):
        box('roof tar seam',(-width/2+seam+.5,0,height+.19),(.045,depth-.5,.02),iron if 'iron' in locals() else dark)

    for axis in (0, 1):
        for side in (-1, 1):
            w = width if axis == 0 else depth
            d = depth if axis == 0 else width
            def facade(name, along, z, dims, mat):
                xyz = (along, side*(d/2+.06), z) if axis == 0 else (side*(d/2+.06), along, z)
                ds = dims if axis == 0 else (dims[1],dims[0],dims[2])
                return box(name, xyz, ds, mat)
            count = max(3, int(w/2.6))
            for floor in range(floors):
                for col in range(count):
                    x = (col-(count-1)/2)*w/count
                    z = floor*3.15 + 1.8
                    facade('window surround', x,z,(1.85 if floor==0 else 1.36,.16,1.95),stone)
                    facade('window pane', x,z,(1.58 if floor==0 else 1.09,.24,1.64),warm if (floor+col+side+seed)%5==0 else glass)
                    facade('window mullion', x,z,(.065,.26,1.64),iron)
                    facade('window sill',x,z-.91,(1.5,.42,.14),stone)
            facade('shop fascia',0,2.85,(w+.2,.35,.36),iron)
            facade('entrance',0,1.15,(1.35,.32,2.25),iron)
            facade('door glass',0,1.45,(1.02,.36,1.4),glass)
            for col in range(count+1):
                x=(col-count/2)*w/count
                facade('ground pier',x,1.4,(.21,.27,2.65),stone)
            if kind in ('tavern','shop','casino') and axis==0 and side==1:
                cloth=material('burgundy canvas',(.24,.055,.045))
                for col in range(count):
                    x=(col-(count-1)/2)*w/count
                    awning=facade('canvas awning',x,2.62,(w/count-.18,1.35,.12),cloth)
                    awning.location.y += .62
                    awning.rotation_euler.x=math.radians(-12)
                    facade('awning valance',x,2.4,(w/count-.18,1.5,.26),cloth)

    # Rooftop skylight, chimney, water tank, fire escape: silhouette at any angle.
    box('skylight',(-2,1,height+.45),(2.2,2.8,.6),iron,.08)
    box('skylight glazing',(-2,1,height+.79),(2,2.6,.08),glass)
    box('chimney',(width*.3,depth*.25,height+1.2),(1.15,1.1,2.4),wall)
    box('chimney coping',(width*.3,depth*.25,height+2.4),(1.4,1.35,.25),stone)
    if floors >= 3:
        for z in range(1,floors):
            h=z*3.15
            box('escape platform',(0,-depth/2-.8,h),(3.5,1.4,.13),iron)
            for x in (-1.6,1.6):
                box('escape rail',(x,-depth/2-1.4,h+.5),(.07,.07,1),iron)
            box('escape handrail',(0,-depth/2-1.4,h+1),(3.4,.07,.07),iron)
            for step in range(9):
                box('escape stair',(-1.2+step*.3,-depth/2-.7,h+step*.35),(.4,.9,.08),iron)
        if kind not in ('casino','civic'):
            rooftop_tank(height,iron)
    if kind == 'casino':
        for ob in list(bpy.context.scene.objects):
            if ob.name.startswith(('skylight','tank','water tank')):
                bpy.data.objects.remove(ob,do_unlink=True)
        anchor.location=(0,depth/2+2.14,3.3)
        anchor.scale=(1,1,.45)
        box('marquee',(0,depth/2+1,3.2),(9,2,.55),stone,.1)
        for i in range(16):
            box('marquee bulb',(-4+i*.53,depth/2+2.05,2.88),(.16,.12,.16),warm)
        brass=material('aged marquee brass',(.52,.37,.15),.65)
        neon=material('neon tubing',accent or (.8,.025,.018),0,1.8)
        for side in (-1,1):
            for fin in range(3):
                x=side*(width/2-.25-fin*.32)
                box('deco vertical fin',(x,depth/2+.22,height*.6),(.17,.5,height+1.25-fin*.7),stone)
            box('neon blade casing',(side*(width/2-.95),depth/2+.6,height*.78),(.65,.6,4),iron,.1)
            for z in range(6):
                box('neon blade tube',(side*(width/2-.95),depth/2+.93,height*.78-1.55+z*.61),(.43,.045,.08),neon)
        for tier in range(3):
            box('stepped deco crown',(0,0,height+.55+tier*.52),(width-1.5-tier*2,depth-1.5-tier*2,.35),stone)
        for y in (depth/2+.16,depth/2+1.85):
            box('marquee brass edge',(0,y,3.44),(8.8,.065,.09),brass)
    if kind == 'civic':
        for ob in list(bpy.context.scene.objects):
            if ob.name.startswith(('skylight','tank','water tank')):
                bpy.data.objects.remove(ob,do_unlink=True)
        for tier,(w,d,h) in enumerate([(7,6,1.4),(5,4.8,1.2),(3.4,3.4,2.8)]):
            base=height+.8+sum([1.4,1.2,2.8][:tier])
            box('municipal tower',(0,0,base+h/2),(w,d,h),stone)
            box('tower ledge',(0,0,base+h),(w+.3,d+.3,.18),iron)
        clock_z=height+5.2
        face=material('ivory clock face',(.81,.76,.60),0,.2)
        cylinder('clock bezel',(0,1.77,clock_z),.82,.14,iron,(math.pi/2,0,0))
        cylinder('clock face',(0,1.86,clock_z),.71,.05,face,(math.pi/2,0,0))
        for name,length in [('minute',.56),('hour',.38)]:
            hand=box('clock-hand-'+name,(0,1.91,clock_z+length/2),(.045,.03,length),iron)
            bpy.context.scene.cursor.location=(0,1.91,clock_z)
            bpy.ops.object.origin_set(type='ORIGIN_CURSOR')
        for x in (-4,-2,2,4):
            box('civic pilaster',(x,depth/2+.25,3.1),(.42,.48,5.7),stone)
            box('column capital',(x,depth/2+.27,5.95),(.7,.6,.25),stone)
    if kind == 'warehouse':
        for x in (-3,3):
            box('loading door',(x,-depth/2-.1,1.6),(3.8,.22,3),iron)
            for z in range(8):
                box('door rib',(x,-depth/2-.25,.3+z*.35),(3.8,.08,.035),stone)


def car(kind, articulated=True):
    length = {'ford':4.5,'hudson':4.9,'packard':5.6}[kind]
    color = {'ford':(.23,.31,.27),'hudson':(.30,.12,.09),'packard':(.08,.1,.12)}[kind]
    paint=material('enamel',color,.5)
    chrome=material('chrome',(.65,.67,.63),.85)
    glass=material('car glass',(.15,.24,.27),.45)
    rubber=material('rubber',(.025,.027,.025))
    cream=material('whitewalls',(.72,.69,.59))
    lamp=material('headlamps',(.95,.81,.48),0,.7)
    red=material('tail lamps',(.55,.035,.015),0,.3)
    # Nose along -Y; exporter maps this to browser +Z.
    box('chassis',(0,0,.56),(1.8,length,.55),paint,.19)
    box('bonnet',(0,-length*.3,.87),(1.7,length*.34,.32),paint,.14)
    box('cabin glass',(0,.14,1.15),(1.50,length*.4,.63),glass,.20)
    box('roof',(0,.2,1.51),(1.48,length*.31,.16),paint,.12)
    for side in (-1,1):
        box('centre pillar',(side*.78,.15,1.22),(.065,.1,.5),chrome)
        box('chrome sill',(side*.91,0,.7),(.045,length*.85,.05),chrome)
        for axle,y in enumerate((-length*.29,length*.30)):
            before=set(bpy.context.scene.objects)
            cylinder('wheel',(side*.87,y,.39),.37,.22,rubber,(0,math.pi/2,0),32 if articulated else 12)
            cylinder('whitewall',(side*1.0,y,.39),.29,.02,cream,(0,math.pi/2,0),32 if articulated else 12)
            cylinder('hubcap',(side*1.015,y,.39),.18,.025,chrome,(0,math.pi/2,0),32 if articulated else 12)
            if articulated:
                # Valve and wheel bolts provide visible rotation without exaggerated treads.
                cylinder('tyre valve',(side*1.025,y+.22,.49),.014,.03,chrome,(0,math.pi/2,0))
                for bolt in range(5):
                    angle=bolt*math.tau/5
                    cylinder('hub bolt',(side*1.031,y+.10*math.cos(angle),.39+.10*math.sin(angle)),.016,.014,chrome,(0,math.pi/2,0))
                meshes=set(bpy.context.scene.objects)-before
                pivot=bpy.data.objects.new('wheel-roll-'+('front' if axle==0 else 'rear')+('-left' if side<0 else '-right'),None)
                bpy.context.collection.objects.link(pivot);pivot.location=(side*.87,y,.39)
                for ob in meshes:
                    ob.parent=pivot;ob.location-=pivot.location
        box('headlight',(side*.61,-length/2-.015,.77),(.32,.06,.22),lamp,.08)
        box('tail lamp',(side*.65,length/2,.75),(.18,.07,.16),red,.04)
    for y in (-length/2,length/2):
        box('bumper',(0,y,.46),(1.9,.14,.16),chrome,.05)
    for x in range(9):
        box('grille',(-.48+x*.12,-length/2-.025,.69),(.04,.07,.23),chrome)


def villa():
    wall=material('estate pale stucco',(.64,.59,.45))
    stone=material('estate limestone',(.72,.67,.54))
    tile=material('terracotta roof',(.34,.12,.07))
    dark=material('estate ironwork',(.09,.12,.1),.45)
    glass=material('estate glazing',(.13,.22,.23),.2)
    timber=material('shutters',(.13,.22,.17))
    box('estate foundation',(0,0,.3),(11.5,10.5,.6),stone)
    box('stucco residence',(0,0,3.25),(11,10,5.9),wall)
    for side in (-1,1):
        vertices=[(x,side*5+dy,z) for dy in (-.08,.08) for x,z in [(-5.5,6.2),(5.5,6.2),(0,8.6)]]
        mesh=bpy.data.meshes.new('gable masonry')
        mesh.from_pydata(vertices,[],[(0,2,1),(3,4,5),(0,1,4,3),(1,2,5,4),(2,0,3,5)])
        mesh.materials.append(wall)
        ob=bpy.data.objects.new('estate gable',mesh);bpy.context.collection.objects.link(ob)
    # Two true roof planes and repeated tile rolls give the estate a distinct
    # residential silhouette at street level and when orbiting above it.
    for side in (-1,1):
        roof=box('pitched tiled roof',(side*2.9,0,7.3),(6.8,11,.18),tile)
        roof.rotation_euler.y=side*math.radians(24)
        for y in range(23):
            roll=cylinder('clay tile roll',(side*2.9,-5.35+y*.48,7.43),.07,6.85,tile)
            roll.rotation_euler.y=math.pi/2+side*math.radians(24)
    cylinder('ridge tiles',(0,0,8.64),.16,11.2,tile,(math.pi/2,0,0))
    box('estate chimney',(-3,-2.7,7.8),(.9,1,3),wall)
    box('chimney cap',(-3,-2.7,9.35),(1.2,1.3,.2),stone)
    for side in (-1,1):
        for x in (-3.4,0,3.4):
            for z in (1.8,4.7):
                box('window frame',(x,side*5.05,z),(1.5,.16,1.85),stone)
                box('estate window',(x,side*5.16,z),(1.22,.12,1.58),glass)
                box('window crossbar',(x,side*5.24,z),(1.25,.055,.07),stone)
                for edge in (-1,1):
                    box('louvered shutter',(x+edge*.94,side*5.1,z),(.44,.18,1.8),timber)
                    for slat in range(8):
                        box('shutter slat',(x+edge*.94,side*5.22,z-.72+slat*.2),(.4,.055,.045),dark)
    box('front door',(0,5.26,1.5),(1.4,.13,2.4),timber)
    warm=material('estate occupied windows',(.73,.43,.16),0,.4)
    for side in (-1,1):
        for y in (-3,0,3):
            for z in (1.8,4.7):
                box('side window frame',(side*5.56,y,z),(.16,1.5,1.85),stone)
                box('side glazing',(side*5.66,y,z),(.12,1.22,1.58),warm if y==0 else glass)
                box('side window mullion',(side*5.74,y,z),(.045,.06,1.58),stone)
    for step in range(3):
        box('entry steps',(0,6.7-step*.38,.1+step*.13),(4,1.2,.2+step*.26),stone)
    box('porch canopy',(0,5.9,3.15),(5.8,2.7,.28),stone)
    for x in (-2.5,2.5):
        cylinder('porch column',(x,6.8,1.6),.18,3,stone)
        box('porch capital',(x,6.8,3),(.55,.55,.2),stone)
    for side in (-1,1):
        for y in range(12):
            box('estate fence picket',(side*7.25,-5.5+y, .85),(.055,.055,1.6),dark)
        for z in (.45,1.35):box('estate fence rail',(side*7.25,0,z),(.065,11.5,.065),dark)
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor)
    anchor.location=(0,7.27,3.2)


def person(waved=False):
    coat=material('wool suit',(.29,.13,.12) if waved else (.16,.19,.21))
    skin=material('skin',(.58,.38,.24))
    hat=material('felt hat',(.12,.1,.08))
    shoe=material('leather',(.05,.04,.035))
    shirt=material('ivory shirt',(.76,.72,.6))
    tie=material('wine silk tie',(.27,.045,.035))
    hair=material('waved chestnut hair',(.095,.043,.023)) if waved else None
    def oval(name,at,scale,mat):
        bpy.ops.mesh.primitive_uv_sphere_add(segments=16,ring_count=8,location=at)
        ob=bpy.context.object;ob.name=name;ob.scale=scale
        bpy.ops.object.transform_apply(location=False,rotation=False,scale=True)
        ob.data.materials.append(mat)
        for polygon in ob.data.polygons:polygon.use_smooth=True
        return ob
    def joint(name,at,parent=None):
        ob=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(ob)
        ob.location=at
        bpy.context.view_layer.update()
        if parent:
            matrix=ob.matrix_world.copy();ob.parent=parent;ob.matrix_world=matrix
        return ob
    def attach(ob,parent):
        matrix=ob.matrix_world.copy();ob.parent=parent;ob.matrix_world=matrix
        return ob
    # Three jacket rings produce shoulders, a fitted waist and a wider hem.
    vertices=[]
    for z,w,d in ([(.88,.23,.15),(1.09,.17,.13),(1.4,.23,.15)] if waved else [(.83,.23,.15),(1.06,.20,.135),(1.4,.255,.15)]):
        vertices.extend([(-w,-d,z),(w,-d,z),(w,d,z),(-w,d,z)])
    faces=[(3,2,1,0),(8,9,10,11)]
    for ring in range(2):
        for i in range(4):faces.append((ring*4+i,ring*4+(i+1)%4,(ring+1)*4+(i+1)%4,(ring+1)*4+i))
    mesh=bpy.data.meshes.new('tailored jacket');mesh.from_pydata(vertices,[],faces);mesh.materials.append(coat)
    torso=bpy.data.objects.new('jacket',mesh);bpy.context.collection.objects.link(torso)
    oval('neck',(0,0,1.46),(.075,.08,.13),skin)
    oval('head',(0,-.005,1.64),(.125,.115,.17),skin)
    oval('nose',(0,-.119,1.63),(.035,.045,.045),skin)
    for side in (-1,1):oval('ear',(side*.125,0,1.64),(.025,.035,.052),skin)
    if waved:
        oval('swept hair crown',(0,.025,1.76),(.145,.125,.08),hair)
        oval('pinned hair back',(0,.087,1.64),(.145,.07,.15),hair)
        for side in (-1,1):
            for wave in range(4):
                oval('sculpted hair wave',(side*(.119+wave*.003),.022,1.74-wave*.052),(.04,.105,.039),hair)
        for wave in range(3):
            lock=oval('front finger wave',(-.075+wave*.06,-.09,1.765-wave*.01),(.07,.035,.035),hair)
            lock.rotation_euler.y=math.radians(-20)
    else:
        oval('fedora brim',(0,0,1.79),(.235,.205,.022),hat)
        oval('fedora crown',(0,.015,1.865),(.155,.145,.095),hat)
        cylinder('hat ribbon',(0,.015,1.827),.154,.028,shoe)
    box('shirt front',(0,-.151,1.285),(.19,.015,.23),shirt,.012)
    for side in (-1,1):
        lapel=box('notched lapel',(side*.108,-.164,1.265),(.09,.025,.27),coat,.012)
        lapel.rotation_euler.y=side*math.radians(22)
    if waved:
        for side in (-1,1):
            collar=box('blouse collar',(side*.052,-.174,1.39),(.065,.022,.09),shirt,.012)
            collar.rotation_euler.y=side*math.radians(25)
        oval('gold lapel pin',(.14,-.185,1.31),(.023,.014,.023),material('brass pin',(.62,.43,.12),.65))
    else:
        box('tie',(0,-.171,1.28),(.04,.018,.21),tie,.008)
    for z in (1.09,.99):oval('jacket button',(.045,-.151,z),(.013,.012,.013),shoe)
    for side in (-1,1):
        box('welt pocket',(side*.14,-.153,.995),(.115,.02,.018),shoe)
        hip=joint('leg'+str(side),(side*.12,0,.86))
        attach(box('trouser upper',(side*.12,0,.67),(.19,.22,.4),coat,.045),hip)
        knee=joint('knee'+str(side),(side*.12,0,.47),hip)
        attach(box('trouser lower',(side*.12,0,.29),(.16,.185,.39),coat,.035),knee)
        attach(box('shoe',(side*.12,-.07,.075),(.18,.32,.14),shoe,.05),knee)
        arm=joint('arm'+str(side),(side*.3,0,1.35))
        attach(box('jacket sleeve',(side*.31,0,1.095),(.14,.19,.53),coat,.045),arm)
        attach(box('shirt cuff',(side*.31,0,.842),(.125,.175,.045),shirt,.015),arm)
        attach(oval('hand',(side*.31,-.005,.77),(.066,.065,.095),skin),arm)


def revolver():
    steel=material('blued revolver steel',(.055,.065,.075),.8)
    wood=material('walnut grip',(.17,.07,.025))
    dark=material('revolver bore',(.012,.014,.016))
    grip=box('walnut grip',(0,.025,-.018),(.055,.075,.13),wood,.014)
    grip.rotation_euler.x=math.radians(-12)
    box('steel frame',(0,-.025,.062),(.055,.14,.065),steel,.009)
    cylinder('six shot cylinder',(0,-.055,.083),.041,.073,steel,(math.pi/2,0,0))
    cylinder('barrel',(0,-.191,.086),.021,.22,steel,(math.pi/2,0,0))
    cylinder('muzzle bore',(0,-.302,.086),.010,.003,dark,(math.pi/2,0,0))
    box('front sight',(0,-.265,.11),(.012,.018,.014),steel,.003)
    box('hammer',(0,.027,.113),(.014,.036,.029),steel,.004)
    # A small open trigger guard, rather than a solid block beneath the frame.
    for side in (-1,1):box('guard side',(0,-.061+side*.023,.006),(.012,.011,.048),steel,.003)
    box('guard bottom',(0,-.061,-.018),(.012,.055,.01),steel,.003)
    box('trigger',(0,-.055,.012),(.01,.01,.029),steel,.003)
    anchor=bpy.data.objects.new('muzzle',None);bpy.context.collection.objects.link(anchor)
    anchor.location=(0,-.307,.086)


def export(name):
    # Join by material except animated limbs: a building becomes ~6 draws.
    groups={}
    for ob in list(bpy.context.scene.objects):
        if ob.type=='MESH' and (ob.parent is None or ob.parent.name.startswith('wheel-roll-')) and not ob.name.startswith(('leg','arm','shoe','clock-hand')):
            key=(ob.parent.name if ob.parent else '',ob.data.materials[0].name)
            groups.setdefault(key,[]).append(ob)
    for obs in groups.values():
        bpy.ops.object.select_all(action='DESELECT')
        for ob in obs: ob.select_set(True)
        bpy.context.view_layer.objects.active=obs[0]
        bpy.ops.object.join()
    bpy.context.view_layer.update()
    points=[ob.matrix_world @ vertex.co for ob in bpy.context.scene.objects if ob.type=='MESH' for vertex in ob.data.vertices]
    bounds=[[round(min(p[i] for p in points),4) for i in range(3)], [round(max(p[i] for p in points),4) for i in range(3)]]
    bpy.ops.export_scene.gltf(filepath=os.path.join(OUT,name+'.glb'),export_format='GLB',export_yup=True,export_cameras=False,export_lights=False)
    return {'file':name+'.glb','bounds_blender':bounds,'bytes':os.path.getsize(os.path.join(OUT,name+'.glb'))}


def beam(name, a, b, width, mat):
    a,b=Vector(a),Vector(b)
    ob=box(name,(a+b)/2,(width,width,(b-a).length),mat)
    ob.rotation_euler=(b-a).to_track_quat('Z','Y').to_euler()
    return ob


def corrugated_roof(name, xyz, width, depth, mat):
    # Closed zinc sheet profile with real ridges, physical-scale UVs and no
    # independent mesh per rib. Joined with other zinc parts at export.
    count=round(width/.25)*4
    vertices=[]
    for i in range(count+1):
        x=-width/2+width*i/count
        z=.035*(1-math.cos(i*math.pi/2))
        vertices.extend([(x,-depth/2,z),(x,depth/2,z),(x,-depth/2,-.25),(x,depth/2,-.25)])
    faces=[]
    for i in range(count):
        a=i*4;b=a+4
        faces.extend([(a,b,b+1,a+1),(a+2,a+3,b+3,b+2),(a+2,b+2,b,a),(a+1,b+1,b+3,a+3)])
    faces.extend([(0,1,3,2),(count*4+2,count*4+3,count*4+1,count*4)])
    mesh=bpy.data.meshes.new(name);mesh.from_pydata(vertices,[],faces);mesh.update()
    ob=bpy.data.objects.new(name,mesh);bpy.context.collection.objects.link(ob);ob.location=xyz
    mesh.materials.append(mat);uv=mesh.uv_layers.new(name='UVMap')
    for polygon in mesh.polygons:
        for index in polygon.loop_indices:
            co=mesh.vertices[mesh.loops[index].vertex_index].co
            uv.data[index].uv=(co.x/2.5,co.y/2.5)
    return ob


def zinc_texture(mat):
    # Tileable mottled galvanising and restrained oxidation; authored pixels
    # survive glTF export without requiring procedural runtime shaders.
    n=256;rng=random.Random(1956);pixels=[];normals=[]
    for y in range(n):
        for x in range(n):
            mottling=math.sin(x*math.tau/32+math.sin(y*math.tau/64))*math.sin(y*math.tau/16+x*math.tau/64)
            grain=rng.uniform(-.035,.035)
            oxide=max(0,mottling-.50)*.35
            tone=.86+.12*mottling+grain
            pixels.extend((.34*tone+oxide*.20,.37*tone-oxide*.10,.36*tone-oxide*.19,1))
            normals.extend((.5+grain,.5+grain*.5,1,1))
    tree=mat.node_tree;shader=tree.nodes['Principled BSDF'];shader.inputs['Roughness'].default_value=.68
    for suffix,data,normal in [('galvanised',pixels,False),('micro-normal',normals,True)]:
        image=bpy.data.images.new('zinc-'+suffix,width=n,height=n)
        if normal:image.colorspace_settings.name='Non-Color'
        image.pixels=data;image.pack();tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        if normal:
            node=tree.nodes.new('ShaderNodeNormalMap');node.inputs['Strength'].default_value=.35
            tree.links.new(tex.outputs['Color'],node.inputs['Color']);tree.links.new(node.outputs['Normal'],shader.inputs['Normal'])
        else:tree.links.new(tex.outputs['Color'],shader.inputs['Base Color'])


def industrial(kind):
    stone=material('aged concrete',(.45,.43,.37))
    wall=material('industrial brick',(.41,.22,.14));brick(wall,8)
    roof=material('corrugated zinc',(.27,.30,.29),.45)
    if kind in ('garage','dealer','docks','haulage'):zinc_texture(roof)
    iron=material('rusted steel',(.24,.18,.11),.55)
    glass=material('industrial glazing',(.18,.27,.28),.25)
    cream=material('cream enamel',(.76,.69,.52),.2)
    red=material('pump red',(.42,.07,.04),.3)
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor)
    anchor.location=(0,-3.28,4.05)
    if kind=='filling':
        anchor.location=(0,-5.06,3.52)
        box('station office',(0,4,1.65),(12,5,3.3),wall)
        box('station roof',(0,4,3.4),(12.4,5.4,.25),cream)
        box('office glazing',(0,1.44,1.85),(8,.1,1.6),glass)
        box('pump canopy',(0,-2,3.4),(13,6,.25),cream,.1)
        box('canopy stripe',(0,-2,3.5),(13.1,6.1,.16),red,.05)
        for x in (-4,4):
            cylinder('canopy column',(x,-2,1.65),.15,3.3,cream)
            for y in (-3,-1):
                box('pump base',(x,y,.8),(.75,.65,1.6),red,.12)
                box('pump dial',(x,y-.35,1.25),(.55,.06,.4),cream)
                cylinder('pump globe',(x,y,1.8),.27,.18,cream)
        return
    if kind in ('garage','dealer'):
        for side in (-1,1):
            for y in (-1,2,5):
                box('workshop side window',(side*7.07,y,2.8),(.16,2,1.2),cream)
                box('workshop glazing',(side*7.17,y,2.8),(.06,1.75,.95),glass)
        for x in (-4.5,0,4.5):
            box('back clerestory',(x,7.1,3.2),(3.6,.15,.85),glass)
            cylinder('roof vent',(x,2,5),.35,1,iron,vertices=24)
            cylinder('vent flashing',(x,2,4.72),.52,.08,roof,vertices=24)
            cylinder('vent rain cap',(x,2,5.50),.50,.10,roof,vertices=24)

        box('workshop',(0,2,2.25),(14,10,4.5),wall)
        corrugated_roof('workshop corrugated roof',(0,2,4.675),14.4,10.4,roof)
        for x in (-4.5,0,4.5):
            box('garage bay',(x,-3.1,1.9),(3.6,.2,3.6),iron)
            for z in range(12):
                box('door seam',(x,-3.24,.3+z*.29),(3.5,.035,.035),roof)
            box('bay windows',(x,-3.27,2.8),(3.2,.05,.6),glass)
        if kind=='dealer':
            before=set(bpy.context.scene.objects)
            car('ford',articulated=False)
            for ob in set(bpy.context.scene.objects)-before: ob.location += Vector((-3,-6,0))
        return
    if kind=='chapel':
        anchor.location=(0,-5.57,4.0)
        box('chapel nave',(0,1,2.8),(9,12,5.6),wall)
        for x in (-2.3,2.3):
            ob=box('pitched roof',(x,1,6.5),(5.5,12.7,.24),roof)
            ob.rotation_euler.y=math.radians(25 if x>0 else -25)
        timber=material('chapel oak door',(.17,.085,.033))
        brass=material('chapel door brass',(.52,.36,.12),.65)
        box('chapel door',(0,-5.61,1.65),(1.95,.16,3.1),timber,.025)
        for side in (-1,1):
            box('entry stone jamb',(side*1.12,-5.59,1.75),(.24,.22,3.5),stone,.02)
            for height in (.8,1.7,2.6):
                box('oak raised panel',(side*.48,-5.704,height),(.72,.035,.64),timber,.025)
            box('entry pull',(side*.12,-5.747,1.55),(.035,.035,.25),brass,.01)
        box('entry lintel',(0,-5.59,3.52),(2.48,.22,.24),stone,.02)
        box('entry threshold',(0,-5.76,.08),(2.48,.75,.16),stone,.02)
        box('tower',(0,-4,6),(3,3,12),wall)
        for z in (8,10):box('tower window',(0,-5.55,z),(.9,.12,1.4),glass)
        beam('cross', (0,-4,12),(0,-4,13.5),.15,iron)
        beam('cross arms',(-.5,-4,13),( .5,-4,13),.15,iron)
        return
    anchor.location=(-3,-3.6,4.25)
    # Dock and haulage yard: cargo shed plus a lattice derrick, fully inside lot.
    box('cargo shed',(-3,1.5,2.6),(7,10,5.2),wall)
    corrugated_roof('cargo corrugated roof',(-3,1.5,5.35),7.4,10.4,roof)
    box('loading gate',(-3,-3.58,1.9),(4.5,.16,3.7),iron)
    for x,y in [(3,4),(5,4),(3,2),(5,2)]:
        box('shipping crate',(x,y,.65),(1.6,1.6,1.3),iron)
        for z in (.25,1.05):box('crate strapping',(x,y,z),(1.64,1.64,.07),roof)
    if kind=='docks':
        for x,y in [(3,-3),(5,-3),(3,-1),(5,-1)]:beam('derrick upright',(x,y,0),(4,-2,10),.13,iron)
        beam('crane boom',(4,-2,9),(6,5,13),.3,iron)
        beam('hoist cable',(6,5,13),(6,5,2),.035,iron)
        for z in range(2,9,2):
            beam('derrick brace',(3,-3,z),(5,-1,z+1.6),.09,iron)


def gravel_texture(mat):
    # Periodic aggregate texture in Blender; box UVs keep its scale in metres.
    n=256;cells=64;rand=random.Random(1953)
    stones=[(rand.random(),rand.random(),rand.random()) for _ in range(cells*cells)]
    pixels=[];heights=[];base=mat.diffuse_color[:3]
    for y in range(n):
        for x in range(n):
            u=x*cells/n;v=y*cells/n;cx=int(u);cy=int(v);nearest=100;shade=0
            for dy in (-1,0,1):
                for dx in (-1,0,1):
                    sx,sy,colour=stones[((cy+dy)%cells)*cells+(cx+dx)%cells]
                    dist=(u-cx-dx-sx)**2+(v-cy-dy-sy)**2
                    if dist<nearest:nearest=dist;shade=colour
            height=max(0,1-math.sqrt(nearest)/.88);heights.append(height)
            light=.65+shade*.5+height*.18
            pixels.extend((*[c*light for c in base],1))
    image=bpy.data.images.new('yard aggregate colour',width=n,height=n);image.pixels=pixels;image.pack()
    tree=mat.node_tree;tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
    tree.links.new(tex.outputs['Color'],tree.nodes['Principled BSDF'].inputs['Base Color'])
    pixels=[]
    for y in range(n):
        for x in range(n):
            dx=heights[y*n+(x-1)%n]-heights[y*n+(x+1)%n]
            dy=heights[((y-1)%n)*n+x]-heights[((y+1)%n)*n+x]
            vec=Vector((dx,dy,1.6)).normalized();pixels.extend((vec.x*.5+.5,vec.y*.5+.5,vec.z*.5+.5,1))
    image=bpy.data.images.new('yard aggregate normal',width=n,height=n);image.colorspace_settings.name='Non-Color';image.pixels=pixels;image.pack()
    tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
    normal=tree.nodes.new('ShaderNodeNormalMap');normal.inputs['Strength'].default_value=.35
    tree.links.new(tex.outputs['Color'],normal.inputs['Color'])
    tree.links.new(normal.outputs['Normal'],tree.nodes['Principled BSDF'].inputs['Normal'])


def undertaker():
    wall=material('funeral red brick',(.28,.15,.115));brick(wall,47)
    stone=material('funeral sandstone',(.53,.49,.40))
    slate=material('funeral slate',(.105,.135,.145))
    oak=material('funeral oak',(.085,.055,.031))
    glass=material('funeral glazing',(.075,.115,.12),.3)
    brass=material('funeral brass',(.57,.42,.16),.7)
    iron=material('funeral iron',(.055,.07,.065),.6)
    gravel=material('coach yard gravel',(.29,.285,.24));gravel_texture(gravel)
    box('coach yard',(0,-3.25,-.01),(15.4,8.9,.04),gravel)
    box('funeral foundation',(0,4.4,.18),(13.2,6.2,.36),stone)
    box('funeral premises',(0,4.4,3.3),(13,6,6.6),wall)
    for z in (.42,3.45,6.55):box('front stone course',(0,7.47,z),(13.2,.2,.18),stone)
    for x in (-6.2,2.8,6.2):box('front pilaster',(x,7.5,1.8),(.28,.28,3.25),stone,.02)
    box('long window surround',(-1.65,7.52,1.85),(7.4,.18,2.42),oak,.025)
    box('empty long window',(-1.65,7.64,1.85),(7.08,.08,2.12),glass)
    for x in (-4,-1.65,.7):box('display window mullion',(x,7.70,1.85),(.045,.06,2.15),brass)
    box('window sill',(-1.65,7.68,.7),(7.55,.40,.15),stone,.025)
    box('funeral entry',(4.5,7.55,1.57),(1.75,.2,2.94),oak,.025)
    for x in (4.08,4.92):
        for z in (.72,1.55,2.38):box('door panels',(x,7.675,z),(.65,.065,.61),oak,.02)
    box('door pull',(4.35,7.73,1.42),(.04,.04,.26),brass,.008)
    box('brass business plate',(3.2,7.69,1.77),(.45,.045,.63),brass,.02)
    box('entry lintel',(4.5,7.59,3.17),(2.15,.28,.23),stone,.02)
    box('entry threshold',(4.5,7.7,.075),(2.1,.62,.15),stone,.025)
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor)
    anchor.location=(-1.2,7.61,3.13);anchor.scale=(.82,1,.40)
    def front_window(x):
        box('upper surround',(x,7.46,4.9),(1.45,.18,1.88),stone)
        box('upper window',(x,7.58,4.9),(1.2,.08,1.63),glass)
        box('sash horizontal',(x,7.64,4.9),(1.23,.04,.055),oak)
        box('upper sill',(x,7.62,3.99),(1.57,.34,.13),stone)
    for x in (-4.75,-1.65,1.45,4.55):front_window(x)
    for side in (-1,1):
        for y in (2.65,5.9):
            for z in (1.85,4.9):
                box('side surround',(side*6.56,y,z),(.18,1.3,1.85),stone)
                box('side sash',(side*6.67,y,z),(.07,1.06,1.60),glass)
                box('side crossbar',(side*6.72,y,z),(.04,1.08,.055),oak)
    box('rear coach door',(-3,1.32,1.6),(3.3,.16,3.05),oak)
    for x in (-3.85,-2.15):box('rear door glazing',(x,1.215,2.15),(1.30,.045,.85),glass)
    for x in (1,4.5):box('rear sash',(x,1.30,4.8),(1.35,.17,1.75),glass)
    # Close the roof ends with brick gables and physical-scale UVs.
    for y in (1.4,7.4):
        vertices=[(x,y+dy,z) for dy in (-.07,.07) for x,z in [(-6.5,6.6),(6.5,6.6),(0,8.72)]]
        mesh=bpy.data.meshes.new('funeral gable');mesh.from_pydata(vertices,[],[(0,2,1),(3,4,5),(0,1,4,3),(1,2,5,4),(2,0,3,5)])
        mesh.materials.append(wall);uv=mesh.uv_layers.new()
        for poly in mesh.polygons:
            for loop in poly.loop_indices:
                co=mesh.vertices[mesh.loops[loop].vertex_index].co;uv.data[loop].uv=(co.x/2.5,co.z/2.5)
        ob=bpy.data.objects.new('brick gable',mesh);bpy.context.collection.objects.link(ob)
    # A low gabled slate roof, with courses and separate chimney caps.
    for side in (-1,1):
        roof=box('slate roof',(side*3.4,4.4,7.65),(7.15,6.65,.17),slate)
        roof.rotation_euler.y=side*math.radians(18)
        for course in range(8):
            x=side*(.42+course*.86);z=8.72-abs(x)*math.tan(math.radians(18))
            box('slate course',(x,4.4,z+.045),(.025,6.65,.018),iron)
    cylinder('slate ridge',(0,4.4,8.77),.1,6.7,slate,(math.pi/2,0,0))
    for x in (-4.7,4.7):
        box('brick chimney',(x,5.55,7.82),(.7,.9,1.7),wall)
        box('chimney cap',(x,5.55,8.72),(.88,1.05,.18),stone)
    # Rear yard has a wide opening; the gates swing inward within the lot.
    for side in (-1,1):
        for y in range(9):box('side fence picket',(side*7.65,-7.55+y,.8),(.055,.055,1.6),iron)
        for z in (.45,1.3):box('side fence rail',(side*7.65,-3.55,z),(.07,8.1,.055),iron)
        for x in (3,4,5,6,7):box('back fence picket',(side*x,-7.65,.8),(.055,.055,1.6),iron)
        for z in (.45,1.3):box('back fence rail',(side*5,-7.65,z),(5.2,.07,.055),iron)
        box('gate post',(side*2.45,-7.65,.88),(.13,.13,1.76),iron)
        # Each gate leaf is authored as a rotated group of iron members.
        before=set(bpy.context.scene.objects)
        for x in range(6):box('gate upright',(side*(.15+x*.45),-7.65,.8),(.045,.06,1.6),iron)
        for z in (.45,1.3):box('gate rail',(side*1.3,-7.65,z),(2.3,.06,.055),iron)
        hinge=Vector((side*2.45,-7.65,0));angle=-side*math.radians(85)
        from mathutils import Matrix
        rot=Matrix.Rotation(angle,4,'Z')
        for ob in set(bpy.context.scene.objects)-before:
            ob.location=hinge+rot.to_3x3()@(ob.location-hinge);ob.rotation_euler.z=angle
    before=set(bpy.context.scene.objects)
    car('packard',articulated=False)
    black=bpy.data.materials['enamel']
    box('hearse coach body',(0,.72,1.36),(1.75,3.35,.75),black,.14)
    box('hearse coach roof',(0,.72,1.78),(1.82,3.45,.16),black,.10)
    for side in (-1,1):
        box('hearse long glass',(side*.90,.72,1.46),(.06,2.75,.43),glass)
        for y in (-.35,.55,1.45):box('hearse pillars',(side*.94,y,1.46),(.045,.055,.46),brass)
    for ob in set(bpy.context.scene.objects)-before:
        x,y,z=ob.location;ob.location=(-2-y,-3.8+x,z+.015);ob.rotation_euler.z+=math.pi/2


def street_bed():
    stone=material('weathered kerbstone',(.47,.46,.40))
    pale=material('replacement kerbstone',(.56,.53,.45))
    iron=material('street cast iron',(.12,.135,.13),.65)
    groove=material('recessed drain',(.035,.045,.04),.4)
    # Four runs of individually jointed stone. Top is just above the existing
    # pavement, while the exposed outside face bridges its raised edge.
    for side in (-1,1):
        for i in range(16):
            at=-11.25+i*1.5
            mat=pale if i%7==3 else stone
            box('kerb frontage',(at,side*11.88,.045),(1.47,.24,.27),mat,.025)
            flank_at=at+(.1175 if i==0 else -.1175 if i==15 else 0)
            flank_length=1.235 if i in (0,15) else 1.47
            box('kerb flank',(side*11.88,flank_at,.045),(.24,flank_length,.27),mat,.025)
    # One cover in the road bordering the frontage. Model coordinates are
    # Blender Z-up; +Y is the road north of this parcel after glTF conversion.
    cylinder('manhole rim',(0,16,-.091),.49,.016,iron)
    cylinder('manhole recess',(0,16,-.080),.447,.007,groove)
    for i in range(-5,6):
        x=i*.072
        length=2*math.sqrt(max(0,.425**2-x*x))
        box('cover ribs',(x,16,-.073),(.025,length,.01),iron,.003)
    for side in (-1,1):
        box('lifting recess',(side*.30,16,-.064),(.07,.14,.006),groove,.006)
    for x in (-8.5,8.5):
        box('drain surround',(x,12.4,-.09),(.82,.42,.018),iron,.018)
        box('drain cavity',(x,12.4,-.078),(.72,.32,.006),groove)
        for i in range(9):
            box('drain bars',(x-.32+i*.08,12.4,-.068),(.027,.32,.012),iron,.003)


def streetside():
    iron=material('street furniture iron',(.12,.16,.14),.6)
    wood=material('weathered bench timber',(.29,.19,.10))
    red=material('hydrant enamel',(.39,.075,.035),.35)
    zinc=material('galvanized bin',(.34,.36,.31),.65)
    # Compact band behind each parcel: no object projects more than 0.4m
    # in depth, so this authored set fits between the facade and kerb.
    for x in (-3.8,-2.2):
        for y in (-.23,.23):box('bench leg',(x,y,.24),(.09,.09,.48),iron)
        box('bench back support',(x,.27,.68),(.08,.08,.95),iron)
    for y in (-.24,-.08,.08,.24):box('seat slat',(-3,y,.5),(2.2,.12,.065),wood,.015)
    for z in (.76,.96,1.16):box('back slat',(-3,.28,z),(2.2,.075,.14),wood,.015)
    for x in (-4,-2):box('bench armrest',(x,0,.76),(.075,.65,.08),iron,.025)
    cylinder('litter bin',(0,0,.43),.28,.86,zinc)
    cylinder('bin rolled rim',(0,0,.88),.3,.065,iron)
    cylinder('bin lid',(0,0,.93),.3,.06,zinc)
    for i in range(12):
        angle=i*math.tau/12
        box('bin vertical rib',(.285*math.cos(angle),.285*math.sin(angle),.43),(.028,.028,.75),iron)
    cylinder('hydrant pedestal',(3,0,.1),.23,.2,iron)
    cylinder('hydrant body',(3,0,.51),.17,.8,red)
    cylinder('hydrant cap',(3,0,.94),.22,.14,red)
    cylinder('hydrant crown',(3,0,1.05),.1,.1,iron)
    for side in (-1,1):
        cylinder('hydrant outlet',(3+side*.22,0,.63),.105,.2,red,(0,math.pi/2,0))
        cylinder('outlet nut',(3+side*.33,0,.63),.075,.05,iron,(0,math.pi/2,0))


manifest={}
for i,(name,floors,w,d) in enumerate([('tenement',4,12,11),('tavern',2,12,12),('casino',2,13,11),('warehouse',1,13,12),('civic',3,13,12),('shop',3,12,11),('villa',2,11,11)]):
    clear()
    if name=='villa':villa()
    else:building(name,floors,w,d,i)
    manifest[name]=export(name)
for name in ('filling','garage','dealer','chapel','docks','haulage'):
    clear();industrial(name);manifest[name]=export(name)
for name,floors,width,depth,seed,palette,accent in [
    ('monarch',3,13,11,31,(.52,.44,.30),(.95,.40,.06)),
    ('bluehour',2,12,12,32,(.24,.30,.34),(.08,.5,.8)),
    ('goldenlily',3,10,12,33,(.26,.34,.25),(.55,.75,.14)),
    ('papermoon',1,14,10,34,(.42,.25,.27),(.85,.18,.35)),
]:
    clear();building('casino',floors,width,depth,seed,palette,accent)
    manifest[name]=export(name)
for name in ('ford','hudson','packard'):
    clear();car(name);manifest[name]=export(name)
clear();car('ford')
paint=bpy.data.materials['enamel'];paint.diffuse_color=(.065,.075,.08,1);paint.node_tree.nodes['Principled BSDF'].inputs['Base Color'].default_value=paint.diffuse_color
box('police door panel',(0,0,.82),(1.83,1.5,.32),material('police cream',(.7,.69,.60)))
cylinder('red beacon',(0,0,1.76),.18,.28,material('beacon',(.8,.02,.01),0,2))
manifest['police']=export('police')
for name in ('person','woman'):
    clear();person(name == 'woman');manifest[name]=export(name)
    motion_points=[]
    for sample in range(48):
        for ob in bpy.context.scene.objects:
            if ob.name.startswith(('leg','arm','knee')):
                phase=sample*math.tau/48+(0 if ob.name.endswith('-1') else math.pi)
                ob.rotation_euler.x=(max(0,math.sin(phase+.7))*.65 if ob.name.startswith('knee') else
                    math.sin(phase+(math.pi if ob.name.startswith('arm') else 0))*(.23 if ob.name.startswith('arm') else .35))
        bpy.context.view_layer.update()
        motion_points.extend(ob.matrix_world @ Vector(c) for ob in bpy.context.scene.objects if ob.type=='MESH' for c in ob.bound_box)
    manifest[name]['motion_bounds_blender']=[[round(min(p[i] for p in motion_points),4) for i in range(3)],[round(max(p[i] for p in motion_points),4) for i in range(3)]]
clear();undertaker();manifest['undertaker']=export('undertaker')
clear();revolver();manifest['revolver']=export('revolver')
clear();blast_fragment();manifest['blast-fragment']=export('blast-fragment')
clear();street_bed();manifest['street-bed']=export('street-bed')
clear();streetside();manifest['streetside']=export('streetside')
with open(os.path.join(OUT,'manifest.json'),'w') as f: json.dump(manifest,f,indent=2)
print('Exported',len(manifest),'models')
