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


def roof_texture(mat, slate=False):
    """Physical 2.5m tiles: overlapping slate courses or mineral surfaced felt."""
    n=256 if slate else 128;rng=random.Random(1957 if slate else 1958)
    base=mat.diffuse_color[:3];pixels=[];heights=[];rough=[]
    tones=[rng.uniform(.76,1.18) for _ in range(64)]
    for y in range(n):
        for x in range(n):
            grain=rng.uniform(-1,1)
            if slate:
                # X runs down the roof slope; Y runs parallel to its ridge.
                row=x//32;offset=(y+(row%2)*16)%n;u=x%32;v=offset%32
                seam=u<2 or v<1
                cleft=math.sin(v*.8+u*.17)*.018+math.sin(v*2.1-u*.4)*.012
                tone=(.38 if seam else tones[row*8+offset//32]+grain*.065+cleft)
                height=.05 if seam else .35+u/32*.28+cleft
                r=.83+grain*.045
            else:
                broad=math.sin(x*math.tau/n)*math.cos(y*math.tau/n)*.045
                mineral=.15 if grain>.58 else (-.08 if grain<-.7 else 0)
                tone=.90+broad+mineral+grain*.08
                height=.3+grain*.065;r=.92+grain*.035
            # Encode linear material colour for an sRGB image.
            pixels.extend((*[1.055*(c*tone)**(1/2.4)-.055 for c in base],1));heights.append(height)
            rough.extend((r,r,r,1))
    normals=[]
    for y in range(n):
        for x in range(n):
            dx=heights[y*n+(x-1)%n]-heights[y*n+(x+1)%n]
            dy=heights[((y-1)%n)*n+x]-heights[((y+1)%n)*n+x]
            v=Vector((dx,dy,1)).normalized();normals.extend((v.x*.5+.5,v.y*.5+.5,v.z*.5+.5,1))
    tree=mat.node_tree;shader=tree.nodes['Principled BSDF']
    for suffix,data,target in [('surface',pixels,'Base Color'),('relief',normals,'Normal'),('roughness',rough,'Roughness')]:
        image=bpy.data.images.new(mat.name+' '+suffix,width=n,height=n)
        if target!='Base Color':image.colorspace_settings.name='Non-Color'
        image.pixels=data;image.pack();tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        if target=='Normal':
            normal=tree.nodes.new('ShaderNodeNormalMap');normal.inputs['Strength'].default_value=.65 if slate else .4
            tree.links.new(tex.outputs['Color'],normal.inputs['Color']);tree.links.new(normal.outputs['Normal'],shader.inputs['Normal'])
        else:tree.links.new(tex.outputs['Color'],shader.inputs[target])


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
    dark = material('tar roof', (.13,.14,.13));roof_texture(dark)
    glass = material('smoked glass', (.12,.19,.21), .25)
    warm = material('occupied windows', (.78,.46,.17), 0, .35)
    iron = material('painted iron', (.12,.15,.14), .5)
    height = floors * 3.15
    glazing={}
    for state in ('intact','broken'):
        group=bpy.data.objects.new('window-'+state,None);bpy.context.collection.objects.link(group);glazing[state]=group
    recess=material('unlit window recess',(.018,.023,.024))
    shards=material('broken glass edges',(.19,.31,.32),.45)
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor);anchor.location=(0,depth/2+.26,2.85)
    # A real front vestibule, not a painted doorway on a solid block.
    # The 2m clear opening and 3m recess admit an articulated person.
    for side in (-1,1):
        box('masonry wing',(side*(width/4+.5),0,height/2),((width-2)/2,depth,height),wall)
        box('foundation wing',(side*(width/4+.5625),0,.25),((width-2)/2+.125,depth+.25,.5),stone)
    box('vestibule lintel',(0,0,(height+2.6)/2),(2,depth,height-2.6),wall)
    box('vestibule rear',(0,-1.5,1.3),(2,depth-3,2.6),wall)
    box('vestibule paving',(0,depth/2-1.5,-.04),(2,3,.08),stone)
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
                    if axis==0 and side==1 and floor==0 and abs(x)<2.2:
                        if abs(x)<.01: continue
                        x=2.2 if x>0 else -2.2
                    facade('window surround', x,z,(1.85 if floor==0 else 1.36,.16,1.95),stone)
                    if axis==0 and side==1:
                        vent=bpy.data.objects.new('fire-window-'+str(floor)+'-'+str(col),None)
                        bpy.context.collection.objects.link(vent);vent.location=(x,d/2+.22,z)
                    pane_width=1.58 if floor==0 else 1.09
                    pane=facade('window pane', x,z,(pane_width,.24,1.64),warm if (floor+col+side+seed)%5==0 else glass)
                    if axis==0 and side==1 and (floor>0 or floors==1):
                        pane.parent=glazing['intact']
                        backing=box('broken window recess',(x,d/2+.155,z),(pane_width,.025,1.64),recess,0)
                        backing.parent=glazing['broken']
                        # Four separated jagged remnants cling to the edges; the centre is missing.
                        outline=[(-1,-1),(-.18,-1),(-.35,-.62),(-.56,-.74),(-.70,-.30),(-1,-.18)]
                        for sx,sz in ((1,1),(-1,1),(1,-1),(-1,-1)):
                            verts=[(x+px*sx*pane_width/2,d/2+.195,z+pz*sz*.82) for px,pz in outline]
                            mesh=bpy.data.meshes.new('jagged glass');mesh.from_pydata(verts,[],[tuple(range(len(verts)))]);mesh.update()
                            ob=bpy.data.objects.new('window glass remnant',mesh);bpy.context.collection.objects.link(ob);ob.data.materials.append(shards);ob.parent=glazing['broken']
                            # Two-sided thin glazing exports visible edges from either orbit direction.
                            solid=ob.modifiers.new('glass thickness','SOLIDIFY');solid.thickness=.008

                    facade('window mullion', x,z,(.065,.26,1.64),iron)
                    facade('window sill',x,z-.91,(1.5,.42,.14),stone)
            facade('shop fascia',0,2.85,(w+.2,.35,.36),iron)
            if not(axis==0 and side==1):
                facade('entrance',0,1.15,(1.35,.32,2.25),iron)
                facade('door glass',0,1.45,(1.02,.36,1.4),glass)
            for col in range(count+1):
                x=(col-count/2)*w/count
                if not(axis==0 and side==1 and abs(x)<1):
                    facade('ground pier',x,1.4,(.21,.27,2.65),stone)
            if kind in ('tavern','shop','casino') and axis==0 and side==1:
                cloth=material('burgundy canvas',(.24,.055,.045))
                for col in range(count):
                    x=(col-(count-1)/2)*w/count
                    awning=facade('canvas awning',x,2.62,(w/count-.18,1.35,.12),cloth)
                    awning.location.y += .62
                    awning.rotation_euler.x=math.radians(-12)
                    facade('awning valance',x,2.4,(w/count-.18,1.5,.26),cloth)

    oak=material('entrance polished oak',(.16,.075,.029))
    brass=material('entrance aged brass',(.48,.34,.12),.65)
    pivot=bpy.data.objects.new('entrance-door-hinge',None)
    bpy.context.collection.objects.link(pivot);pivot.location=(-.92,depth/2+.08,.02)
    def doorpart(name,xyz,dims,mat,bevel=.01):
        ob=box(name,xyz,dims,mat,bevel)
        ob.parent=pivot;ob.location-=pivot.location
        return ob
    # Thin independent leaf with framed glazing and recessed lower panels.
    doorpart('door bottom rail',(0,depth/2+.08,.14),(1.84,.12,.24),oak)
    doorpart('door middle rail',(0,depth/2+.08,1.0),(1.84,.12,.16),oak)
    doorpart('door top rail',(0,depth/2+.08,2.28),(1.84,.12,.16),oak)
    for x in (-.84,.84):doorpart('door stile',(x,depth/2+.08,1.19),(.16,.12,2.34),oak)
    doorpart('door lower panel',(0,depth/2+.07,.58),(1.52,.07,.7),oak)
    doorpart('door glazing',(0,depth/2+.08,1.65),(1.52,.045,1.14),glass)
    doorpart('door brass kickplate',(0,depth/2+.15,.22),(1.5,.025,.22),brass)
    doorpart('door brass pull',(.63,depth/2+.20,1.14),(.045,.075,.3),brass)
    for x in (-1.08,1.08):box('entry stone jamb',(x,depth/2+.04,1.28),(.16,.32,2.56),stone)
    box('entry stone head',(0,depth/2+.04,2.49),(2.32,.32,.14),stone)
    entry=bpy.data.objects.new('entrance-threshold',None)
    bpy.context.collection.objects.link(entry);entry.location=(0,depth/2+.08,0)

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


def car(kind, articulated=True, police=False):
    car_objects_before=set(bpy.context.scene.objects)
    length = {'ford':4.5,'hudson':4.9,'packard':5.6}[kind]
    color = {'ford':(.23,.31,.27),'hudson':(.30,.12,.09),'packard':(.08,.1,.12)}[kind]
    paint=material('enamel',(.065,.075,.08) if police else color,.5)
    chrome=material('chrome',(.65,.67,.63),.85)
    glass=material('car glass',(.15,.24,.27),.45)
    rubber=material('rubber',(.025,.027,.025))
    cream=material('whitewalls',(.72,.69,.59))
    lamp=material('headlamps',(.95,.81,.48),0,.7)
    red=material('tail lamps',(.55,.035,.015),0,.3)
    # Nose along -Y; exporter maps this to browser +Z.
    if not articulated:
        box('chassis',(0,0,.56),(1.8,length,.55),paint,.19)
        box('cabin glass',(0,.14,1.15),(1.50,length*.4,.63),glass,.20)
    else:
        # Open cabin above the floor: entry no longer passes through a solid box.
        leather=material('oxblood seat leather',(.12,.065,.043))
        lining=material('door lining',(.22,.17,.105))
        door_paint=material('police cream',(.7,.69,.60)) if police else paint
        box('floor pan',(0,0,.09),(1.70,length-.25,.12),rubber,.07)
        box('front body',(0,-length*.37,.64),(1.8,length*.26,.44),paint,.17)
        box('rear body',(0,length*.365,.67),(1.8,length*.27,.48),paint,.17)
        # Separate front/rear glazing, pillars and door leaves expose real seats.
        front=-length*.17;rear=length*.22
        box('split windscreen',(0,front,1.19),(1.40,.045,.53),glass,.055)
        box('rear windscreen',(0,rear,1.20),(1.39,.045,.49),glass,.06)
        box('windscreen divider',(0,front-.035,1.20),(.025,.028,.55),chrome,.005)
        for x in (-.72,.72):
            for y in (front,rear):box('window pillar',(x,y,1.19),(.07,.09,.59),paint,.025)
        for y in (.04,length*.175):
            box('bench cushion',(0,y,.32),(1.34,.43,.16),leather,.07)
            box('bench backrest',(0,y+.18,.63),(1.34,.13,.48),leather,.065)
            for x in range(9):box('upholstery piping',(-.56+x*.14,y-.015,.406),(.014,.31,.008),lining,.004)
        box('dashboard',(0,front+.1,1.02),(1.4,.22,.19),paint,.05)
        for x in (-.43,-.25,-.07):
            cylinder('dashboard dial',(x,front+.22,1.035),.055,.015,chrome,(math.pi/2,0,0),24)
            cylinder('dial face',(x,front+.232,1.035),.044,.017,rubber,(math.pi/2,0,0),24)
        bpy.ops.mesh.primitive_torus_add(major_radius=.16,minor_radius=.017,major_segments=32,minor_segments=8,location=(-.42,front+.40,1.08),rotation=(math.radians(60),0,0))
        bpy.context.object.name='steering wheel';bpy.context.object.data.materials.append(rubber)
        for side in (-1,1):
            grip=bpy.data.objects.new('seat-driver-grip-'+('left' if side<0 else 'right'),None)
            bpy.context.collection.objects.link(grip);grip.location=((- .42+side*.16)*.984,(front+.40)+(.15-(front+.40))*.035,1.08)
        for side in (-1,1):
            box('cabin sill',(side*.79,.08,.62),(.12,length*.41,.13),paint,.035)
            for index,(a,b) in enumerate(((front+.05,.09),(.19,rear-.04))):
                before=set(bpy.context.scene.objects);middle=(a+b)/2;span=b-a
                box('door skin',(side*.80,middle,.855),(.095,span,.38),door_paint,.035)
                box('door upholstery',(side*.744,middle,.865),(.025,span-.05,.30),lining,.015)
                box('door window',(side*.755,middle,1.225),(.032,span-.07,.36),glass,.025)
                for y in (a+.025,b-.025):box('door window surround',(side*.78,y,1.23),(.04,.04,.45),chrome,.01)
                for z in (1.04,1.445):box('window belt trim',(side*.785,middle,z),(.045,span,.03),chrome,.01)
                box('outside door handle',(side*.866,b-.13,.99),(.035,.17,.028),chrome,.013)
                box('inside door pull',(side*.712,middle,.9),(.03,.16,.025),chrome,.01)
                meshes=set(bpy.context.scene.objects)-before
                hinge=bpy.data.objects.new('car-door-'+('front' if index==0 else 'rear')+('-left' if side<0 else '-right'),None)
                bpy.context.collection.objects.link(hinge);hinge.location=(side*.8,a,.68)
                for ob in meshes:ob.parent=hinge;ob.location-=hinge.location
                seat=bpy.data.objects.new('seat-'+('front' if index==0 else 'rear')+('-left' if side<0 else '-right'),None)
                bpy.context.collection.objects.link(seat);seat.location=(side*.36,.04 if index==0 else length*.175,.495)
            box('mirror stem',(side*.86,front+.08,1.07),(.16,.025,.025),chrome,.008)
            box('wing mirror',(side*.94,front+.08,1.11),(.065,.16,.10),chrome,.025)
    box('bonnet',(0,-length*.3,.87),(1.7,length*.34,.32),paint,.14)
    box('roof',(0,.2,1.51),(1.48,length*(.42 if articulated else .31),.16),paint,.12)
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
        cylinder('headlamp bezel',(side*.61,-length/2-.025,.77),.16,.07,chrome,(math.pi/2,0,0),32)
        cylinder('headlight',(side*.61,-length/2-.065,.77),.135,.025,lamp,(math.pi/2,0,0),32)
        box('tail lamp',(side*.65,length/2,.75),(.18,.07,.16),red,.04)
    for y in (-length/2,length/2):
        box('bumper',(0,y,.46),(1.9,.14,.16),chrome,.05)
    for x in range(9):
        box('grille',(-.48+x*.12,-length/2-.025,.69),(.04,.07,.23),chrome)
    if articulated:
        # Taper the greenhouse into a raked sedan roof; transform door geometry
        # in world space so each separately hinged window keeps the same seam.
        bpy.context.view_layer.update()
        for ob in set(bpy.context.scene.objects)-car_objects_before:
            if ob.type!='MESH':continue
            world=ob.matrix_world.copy();inverse=world.inverted()
            for vertex in ob.data.vertices:
                co=world @ vertex.co
                t=max(0,min(1,(co.z-1.03)/.50))
                if t:
                    co.x*=1-.16*t
                    co.y+=(.15-co.y)*(.35 if co.y<.15 else .22)*t
                    vertex.co=inverse @ co
            for polygon in ob.data.polygons:polygon.use_smooth=True
            bpy.context.view_layer.objects.active=ob
            normal=ob.modifiers.new('Manufactured surface normals','WEIGHTED_NORMAL');normal.keep_sharp=True
            bpy.ops.object.modifier_apply(modifier=normal.name)


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


def wool_texture(mat):
    n=256;pixels=[];heights=[];rough=[];base=mat.diffuse_color[:3]
    for y in range(n):
        for x in range(n):
            warp=math.sin(x*math.pi/2);weft=math.sin(y*math.pi/2)
            twill=1 if ((x//2+y//2)%4)<2 else -1
            stripe=.13 if x%16<1 else 0
            tone=.91+.035*warp+.035*weft+.025*twill+stripe
            pixels.extend((*[c*tone for c in base],1))
            heights.append(.06*warp+.06*weft+.035*twill)
            value=.78+.045*warp*weft;rough.extend((value,value,value,1))
    normals=[]
    for y in range(n):
        for x in range(n):
            dx=heights[y*n+(x-1)%n]-heights[y*n+(x+1)%n]
            dy=heights[((y-1)%n)*n+x]-heights[((y+1)%n)*n+x]
            v=Vector((dx,dy,1)).normalized();normals.extend((v.x*.5+.5,v.y*.5+.5,v.z*.5+.5,1))
    tree=mat.node_tree;shader=tree.nodes['Principled BSDF']
    for name,data,target in [('wool weave colour',pixels,'Base Color'),('wool weave normal',normals,'Normal'),('wool fibre roughness',rough,'Roughness')]:
        image=bpy.data.images.new(name,width=n,height=n)
        if target!='Base Color':image.colorspace_settings.name='Non-Color'
        image.pixels=data;image.pack();tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        if target=='Normal':
            normal=tree.nodes.new('ShaderNodeNormalMap');normal.inputs['Strength'].default_value=.35
            tree.links.new(tex.outputs['Color'],normal.inputs['Color']);tree.links.new(normal.outputs['Normal'],shader.inputs['Normal'])
        else:tree.links.new(tex.outputs['Color'],shader.inputs[target])


def person(waved=False):
    coat=material('wool suit',(.29,.13,.12) if waved else (.16,.19,.21));wool_texture(coat)
    skin=material('skin',(.58,.38,.24))
    hat=material('felt hat',(.12,.1,.08))
    shoe=material('leather',(.05,.04,.035))
    shirt=material('ivory shirt',(.76,.72,.6))
    tie=material('wine silk tie',(.27,.045,.035))
    hair=material('waved chestnut hair' if waved else 'barbered hair',(.095,.043,.023))
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
    def tailored(name,center,rings,mat):
        # Rounded cross sections retain a pressed front while removing box corners.
        # Radius changes form shoulder caps, waist suppression and cloth folds.
        segments=24;vertices=[];faces=[]
        for z,w,d in rings:
            for i in range(segments):
                angle=i*math.tau/segments
                x,y=math.cos(angle),math.sin(angle)
                vertices.append((center[0]+w*math.copysign(abs(x)**.72,x),
                                 center[1]+d*math.copysign(abs(y)**.72,y),z))
        faces.append(tuple(reversed(range(segments))))
        for row in range(len(rings)-1):
            for i in range(segments):
                a=row*segments+i;b=row*segments+(i+1)%segments
                faces.append((a,b,b+segments,a+segments))
        faces.append(tuple((len(rings)-1)*segments+i for i in range(segments)))
        mesh=bpy.data.meshes.new(name);mesh.from_pydata(vertices,[],faces);mesh.update();mesh.materials.append(mat)
        ob=bpy.data.objects.new(name,mesh);bpy.context.collection.objects.link(ob)
        for polygon in mesh.polygons:polygon.use_smooth=len(polygon.vertices)==4
        return ob
    waist=.17 if waved else .20;shoulder=.23 if waved else .255
    tailored('jacket',(0,0),[
        (.88 if waved else .83,.215,.142),(.92,.23,.15),
        (1.06,waist,.135),(1.19,shoulder*.95,.15),
        (1.32,shoulder,.15),(1.38,shoulder*.91,.138),(1.4,.18,.12)],coat)
    oval('neck',(0,0,1.46),(.075,.08,.13),skin)
    oval('head',(0,-.005,1.64),(.125,.115,.17),skin)
    oval('nose',(0,-.119,1.63),(.035,.045,.045),skin)
    for side in (-1,1):oval('ear',(side*.125,0,1.64),(.025,.035,.052),skin)
    eye=material('muted eye whites',(.42,.39,.33));eye.node_tree.nodes['Principled BSDF'].inputs['Roughness'].default_value=.3
    detail=material('facial detail',(.10,.06,.04))
    for side in (-1,1):
        oval('eye white',(side*.052,-.110,1.675),(.027,.014,.012),eye)
        oval('iris',(side*.052,-.125,1.675),(.010,.005,.010),detail)
        oval('eyebrow',(side*.052,-.103,1.704),(.032,.013,.007),detail)
    oval('mouth seam',(0,-.114,1.57),(.034,.008,.006),detail)
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
        # Separate root groups let the renderer choose one portrait silhouette.
        for style in ('full','receding'):
            root=joint('hair-'+style,(0,0,0));vertices=[];faces=[]
            segments=32;rows=10
            for row in range(rows+1):
                for segment in range(segments+1):
                    phi=segment*math.tau/segments
                    start=.72 if style=='receding' else .02
                    edge=(1.23+.58*math.cos(phi)) if style=='receding' else (1.35+.45*math.cos(phi))
                    theta=start+(max(start+.015,edge)-start)*row/rows
                    radius=1+.014*math.sin(phi*15+theta*4)
                    vertices.append((.132*math.sin(theta)*math.sin(phi)*radius,
                                     -.005+.122*math.sin(theta)*math.cos(phi)*radius,
                                     1.64+.181*math.cos(theta)))
            for row in range(rows):
                for segment in range(segments):
                    i=row*(segments+1)+segment;faces.append((i,i+1,i+segments+2,i+segments+1))
            mesh=bpy.data.meshes.new('barbered '+style);mesh.from_pydata(vertices,[],faces);mesh.materials.append(hair)
            ob=bpy.data.objects.new('barbered '+style,mesh);bpy.context.collection.objects.link(ob)
            for polygon in mesh.polygons:polygon.use_smooth=True
            attach(ob,root)
            if style=='full':
                lock=oval('side parted forelock',(-.038,-.065,1.785),(.087,.07,.043),hair)
                lock.rotation_euler.y=math.radians(-12);attach(lock,root)
        fedora=joint('headwear-fedora',(0,0,0))
        attach(oval('fedora brim',(0,0,1.79),(.235,.205,.022),hat),fedora)
        attach(oval('fedora crown',(0,.015,1.865),(.155,.145,.095),hat),fedora)
        attach(cylinder('hat ribbon',(0,.015,1.827),.154,.028,shoe),fedora)
        cap=joint('headwear-cap',(0,0,0))
        attach(oval('cloth cap crown',(0,-.01,1.835),(.20,.17,.07),hat),cap)
        attach(oval('cap peak',(0,-.095,1.795),(.15,.11,.015),hat),cap)
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
        attach(tailored('trouser upper',(side*.12,0),[(.47,.083,.094),(.51,.092,.103),(.65,.095,.11),(.80,.091,.106),(.87,.077,.085)],coat),hip)
        knee=joint('knee'+str(side),(side*.12,0,.47),hip)
        attach(tailored('trouser lower',(side*.12,0),[(.095,.075,.086),(.125,.08,.0925),(.19,.072,.081),(.32,.073,.085),(.44,.08,.09),(.485,.074,.08)],coat),knee)
        attach(box('shoe',(side*.12,-.07,.075),(.18,.32,.14),shoe,.05),knee)
        arm=joint('arm'+str(side),(side*.3,0,1.35))
        attach(tailored('jacket upper sleeve',(side*.31,0),[(1.06,.057,.074),(1.10,.064,.082),(1.23,.07,.095),(1.32,.067,.088),(1.35,.044,.061)],coat),arm)
        elbow=joint('elbow'+str(side),(side*.3,0,1.065),arm)
        attach(tailored('jacket forearm',(side*.31,0),[(.8305,.052,.071),(.86,.060,.078),(.92,.061,.083),(1.00,.0675,.09),(1.045,.060,.081),(1.0755,.053,.068)],coat),elbow)
        attach(box('shirt cuff',(side*.31,0,.842),(.125,.175,.045),shirt,.015),elbow)
        attach(oval('hand',(side*.31,-.005,.77),(.066,.065,.095),skin),elbow)
    # Bake the weave at one physical scale before joints animate; torso meshes
    # and manufactured limb meshes share the same 20cm texture tile.
    bpy.context.view_layer.update()
    for ob in bpy.context.scene.objects:
        if ob.type!='MESH' or coat not in list(ob.data.materials):continue
        uv=ob.data.uv_layers.active or ob.data.uv_layers.new(name='wool physical scale')
        for polygon in ob.data.polygons:
            normal=ob.matrix_world.to_3x3() @ polygon.normal
            axis=max(range(3),key=lambda i:abs(normal[i]));axes=[i for i in range(3) if i!=axis]
            for index in polygon.loop_indices:
                co=ob.matrix_world @ ob.data.vertices[ob.data.loops[index].vertex_index].co
                uv.data[index].uv=(co[axes[0]]/.2,co[axes[1]]/.2)


def handcuffs():
    steel=material('brushed restraint steel',(.42,.45,.47),.8)
    for side in (-1,1):
        bpy.ops.mesh.primitive_torus_add(major_radius=.065,minor_radius=.008,major_segments=32,minor_segments=8,location=(side*.085,0,0),rotation=(0,math.pi/2,0))
        ring=bpy.context.object;ring.name='wrist cuff';ring.data.materials.append(steel)
        for polygon in ring.data.polygons:polygon.use_smooth=True
    for i in range(3):
        bpy.ops.mesh.primitive_torus_add(major_radius=.022,minor_radius=.006,major_segments=16,minor_segments=6,location=((i-1)*.035,0,-.06),rotation=(math.pi/2 if i%2 else 0,0,0))
        link=bpy.context.object;link.name='cuff chain link';link.data.materials.append(steel)


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


def long_gun(kind):
    steel=material('blued gun steel',(.065,.078,.084),.82)
    wood=material('oiled walnut stock',(.24,.105,.035))
    rubber=material('stock butt plate',(.024,.024,.022))
    bore=material('dark barrel bore',(.006,.008,.009))
    # Packed grain follows the stock length; bevels keep small highlights legible.
    n=128;image=bpy.data.images.new('walnut gun grain',width=n,height=n);pixels=[]
    for y in range(n):
        for x in range(n):
            grain=.83+.12*math.sin(x*.47+math.sin(y*.05)*1.5)+.04*math.sin(x*2.9+y*.1)
            pixels.extend((.32*grain,.15*grain,.058*grain,1))
    image.pixels=pixels;image.pack();node=wood.node_tree.nodes.new('ShaderNodeTexImage');node.image=image
    wood.node_tree.links.new(node.outputs['Color'],wood.node_tree.nodes['Principled BSDF'].inputs['Base Color'])
    box('walnut shoulder stock',(0,.22,.03),(.068,.38,.135),wood,.032)
    box('ribbed butt plate',(0,.414,.03),(.078,.018,.15),rubber,.009)
    box('stock wrist',(0,.03,.025),(.05,.13,.09),wood,.018)
    box('receiver',(0,-.085,.075),(.065,.20,.095),steel,.01)
    for side in (-1,1):
        box('trigger guard side',(side*.023,-.025,-.025),(.009,.09,.06),steel,.004)
        cylinder('receiver pin',(side*.034,-.08,.073),.008,.006,steel,(0,math.pi/2,0),16)
    box('trigger guard bow',(0,-.025,-.054),(.045,.09,.01),steel,.004)
    box('trigger',(0,-.034,-.023),(.008,.018,.035),steel,.003)
    box('ejection port',(.034,-.112,.082),(.004,.08,.032),bore,.004)
    length=.51 if kind=='shotgun' else .32
    end=-.18-length
    cylinder('barrel',(0,-.18-length/2,.105),.023 if kind=='shotgun' else .021,length,steel,(math.pi/2,0,0),32)
    cylinder('muzzle bore',(0,end-.002,.105),.014 if kind=='shotgun' else .011,.004,bore,(math.pi/2,0,0),24)
    box('front sight',(0,end+.03,.133),(.01,.022,.018),steel,.003)
    box('rear sight',(0,-.07,.13),(.035,.035,.025),steel,.003)
    if kind=='shotgun':
        cylinder('magazine tube',(0,-.345,.055),.02,.31,steel,(math.pi/2,0,0),24)
        pump=bpy.data.objects.new('pump-slide',None);bpy.context.collection.objects.link(pump)
        ob=box('pump forearm',(0,-.31,.045),(.074,.19,.075),wood,.025);ob.parent=pump
        for i in range(10):
            ob=box('pump groove',(0,-.23-i*.016,.078),(.066,.004,.007),rubber,.001);ob.parent=pump
    else:
        grip=box('pistol grip',(0,-.018,-.06),(.06,.075,.15),wood,.016);grip.rotation_euler.x=-.2
        box('stick magazine',(0,-.165,-.07),(.038,.067,.27),steel,.008)
        box('foregrip',(0,-.31,.055),(.07,.18,.065),wood,.014)
        for i in range(11):cylinder('barrel cooling fin',(0,-.22-i*.019,.105),.029,.008,steel,(math.pi/2,0,0),24)
        box('charging handle',(.055,-.07,.105),(.04,.021,.024),steel,.005)
    for name,pos in [('muzzle',(0,end-.009,.105)),('support-grip',(0,-.27,.005))]:
        anchor=bpy.data.objects.new(name,None);bpy.context.collection.objects.link(anchor);anchor.location=pos
        if kind=='shotgun' and name=='support-grip':anchor.parent=pump


def export(name):
    # Join by material except animated limbs: a building becomes ~6 draws.
    groups={}
    for ob in list(bpy.context.scene.objects):
        if ob.type=='MESH' and (ob.parent is None or ob.parent.name.startswith(('wheel-roll-','interior-wall-','entrance-door-','window-','pump-','car-door-','laundry-drum-','slot-coins'))) and not ob.name.startswith(('leg','arm','shoe','clock-hand')):
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


def harbour_pier():
    timber=material('salt worn dock timber',(.32,.27,.19))
    # Long grain follows each board; packed colour and normal detail survive glTF.
    n=128;colour=[];normal=[]
    for y in range(n):
        for x in range(n):
            g=math.sin(x*.83+math.sin(y*.06))*.08+math.sin(x*2.4+y*.03)*.025
            colour.extend((.55+g,.49+g,.38+g,1))
            v=Vector((math.cos(x*.83+math.sin(y*.06))*.18,0,1)).normalized()
            normal.extend((v.x*.5+.5,.5,v.z*.5+.5,1))
    tree=timber.node_tree;shader=tree.nodes['Principled BSDF']
    for suffix,data in [('grain',colour),('normal',normal)]:
        image=bpy.data.images.new('dock timber '+suffix,width=n,height=n)
        if suffix=='normal':image.colorspace_settings.name='Non-Color'
        image.pixels=data;image.pack();tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        if suffix=='normal':
            bump=tree.nodes.new('ShaderNodeNormalMap');bump.inputs['Strength'].default_value=.5
            tree.links.new(tex.outputs['Color'],bump.inputs['Color']);tree.links.new(bump.outputs['Normal'],shader.inputs['Normal'])
        else:tree.links.new(tex.outputs['Color'],shader.inputs['Base Color'])
    iron=material('dock black iron',(.065,.075,.07),.55)
    for i in range(32):box('landing deck board',(-7.75+i*.5,0,.08),(.48,7.6,.24),timber,.018)
    for y in (-3.35,0,3.35):box('landing bearer',(0,y,-.20),(16,.24,.32),timber)
    for x in (-7,0,7):
        for y in (-3.35,3.35):
            cylinder('timber pile',(x,y,-1.15),.23,2.7,timber,vertices=16)
            cylinder('pile iron collar',(x,y,-.2),.25,.14,iron,vertices=16)
    for x in (-6.5,5.5):
        for y in (-3.25,3.25):
            cylinder('mooring bollard base',(x,y,.24),.3,.08,iron,vertices=16)
            cylinder('mooring bitt',(x,y,.53),.12,.54,iron,vertices=16)
            cylinder('bollard head',(x,y,.80),.22,.10,iron,vertices=16)
    # An open working quay: low edge timbers rather than a fence through the loading space.
    for y in (-3.7,3.7):box('edge rubbing timber',(0,y,.29),(16,.18,.18),timber,.025)
    rubber=material('fender rubber',(.035,.035,.03))
    for x in (-7,0,7):
        for y in (-3.72,3.72):
            bpy.ops.mesh.primitive_torus_add(major_segments=16,minor_segments=8,major_radius=.3,minor_radius=.10,location=(x,y,-.25),rotation=(math.pi/2,0,0))
            ob=bpy.context.object;ob.name='rubber dock fender';ob.data.materials.append(rubber)


def quay_section():
    stone=material('harbour retaining masonry',(.36,.35,.30));brick(stone,1959)
    cap=material('worn quay coping',(.44,.43,.37))
    box('quay retaining wall',(0,0,-1.35),(.5,8,3),stone)
    box('quay coping',(0,0,.19),(.75,8,.12),cap,.025)


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
    slate=material('funeral slate',(.105,.135,.145));roof_texture(slate,True)
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
        for uv in roof.data.uv_layers.active.data:uv.uv.x*=side
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


def mariner():
    wall=material('mariner salt worn brick',(.39,.24,.18));brick(wall,61)
    stone=material('mariner limestone',(.61,.58,.48))
    timber=material('mariner painted timber',(.12,.20,.18))
    iron=material('mariner blackened iron',(.095,.115,.11),.45)
    slate=material('mariner weathered slate',(.19,.23,.25));roof_texture(slate,True)
    glass=material('mariner sash glass',(.13,.21,.23),.2)
    warm=material('mariner occupied room',(.68,.44,.20),0,.25)
    brass=material('mariner door brass',(.48,.36,.16),.6)
    box('lodging house brick',(0,0,4.725),(12,10,9.45),wall)
    box('lodging house foundation',(0,0,.22),(12.25,10.25,.44),stone)
    for z in (3.12,6.27,9.42):box('lodging house string course',(0,0,z),(12.3,10.3,.15),stone)
    for axis in (0,1):
        for side in (-1,1):
            depth=5 if axis==0 else 6
            positions=(-4.3,-2.0,2.0,4.3) if axis==0 else (-3.4,0,3.4)
            def facade(name,along,z,dims,mat):
                xyz=(along,side*(depth+.06),z) if axis==0 else (side*(depth+.06),along,z)
                ds=dims if axis==0 else (dims[1],dims[0],dims[2])
                return box(name,xyz,ds,mat)
            for floor in range(3):
                for i,x in enumerate(positions):
                    z=1.75+floor*3.15
                    facade('sash stone lintel',x,z+1,(1.60,.3,.2),stone)
                    facade('sash painted frame',x,z,(1.35,.22,1.85),timber)
                    facade('sash pane',x,z,(1.10,.27,1.59),warm if (i+floor+axis+side)%6==0 else glass)
                    facade('sash central rail',x,z,(1.12,.31,.065),stone)
                    facade('sash upright',x,z,(.055,.32,1.6),stone)
                    facade('sash sill',x,z-.96,(1.62,.42,.14),stone)
    # One public entrance and a shallow shelter facing the frontage.
    box('mariner oak door',(0,5.2,1.25),(1.6,.18,2.5),timber,.02)
    for side in (-1,1):
        box('door jamb',(side*.94,5.15,1.42),(.22,.32,2.84),stone)
        box('door raised panel',(side*.38,5.305,.7),(.57,.035,.88),timber,.025)
        box('door glazing',(side*.38,5.31,1.72),(.56,.04,.84),glass)
        box('shelter post',(side*1.38,6.4,1.40),(.13,.13,2.8),timber)
        beam('shelter bracket',(side*1.38,6.4,2.2),(side*.90,6.4,2.75),.09,timber)
    box('door pull',(.20,5.355,1.24),(.045,.055,.28),brass,.012)
    box('entry transom',(0,5.17,2.73),(1.65,.26,.22),stone)
    box('entry shelter',(0,5.78,2.88),(3.2,1.65,.19),slate)
    box('lodging sign board',(0,5.21,3.50),(5.5,.22,.67),timber)
    anchor=bpy.data.objects.new('sign-anchor',None);bpy.context.collection.objects.link(anchor);anchor.location=(0,5.34,3.5);anchor.scale=(.71,1,.52)
    # Closed brick gables and two sloping slate planes, with a proper ridge.
    for side in (-1,1):
        mesh=bpy.data.meshes.new('mariner gable')
        mesh.from_pydata([(-6,side*5.02,9.45),(6,side*5.02,9.45),(0,side*5.02,11.63)],[],[(0,1,2) if side<0 else (2,1,0)])
        mesh.materials.append(wall);uv=mesh.uv_layers.new(name='masonry scale')
        for i,loop in enumerate(mesh.loops):
            co=mesh.vertices[loop.vertex_index].co;uv.data[i].uv=(co.x/2.5,co.z/2.5)
        ob=bpy.data.objects.new('mariner brick gable',mesh);bpy.context.collection.objects.link(ob)
        roof=box('mariner pitched slate',(side*3.1,0,10.60),(6.7,10.7,.15),slate)
        roof.rotation_euler.y=side*math.radians(20)
        for uv in roof.data.uv_layers.active.data:uv.uv.x*=side
        cylinder('gable vent surround',(0,side*5.08,10.35),.43,.12,stone,(math.pi/2,0,0),24)
        cylinder('gable vent dark',(0,side*5.16,10.35),.32,.06,iron,(math.pi/2,0,0),24)
        for i in range(-2,3):box('gable vent louvre',(0,side*5.20,10.35+i*.10),(.48,.06,.035),timber)
    box('slate ridge cap',(0,0,11.78),(.25,10.8,.18),slate,.06)
    box('lodging chimney',(2,-2,11.13),(1.1,1.05,2.6),wall)
    box('chimney cap',(2,-2,12.47),(1.35,1.3,.16),stone)
    for x in (1.75,2.25):cylinder('chimney clay pot',(x,-2,12.70),.16,.34,wall,vertices=16)
    for level in (3.15,6.3):
        box('rear escape landing',(0,-5.75,level),(3.4,1.4,.12),iron)
        for x in (-1.6,1.6):box('rear escape upright',(x,-6.4,level+.5),(.055,.055,1),iron)
        box('rear escape handrail',(0,-6.4,level+1),(3.3,.055,.055),iron)
        for x in range(12):box('rear escape baluster',(-1.5+x*.27,-6.4,level+.5),(.025,.025,1),iron)
        for step in range(9):box('rear escape stair',(-1.2+step*.3,-5.75,level+step*.35),(.37,.9,.08),iron)


def vacant_lot():
    gravel=material('vacant yard earth and gravel',(.30,.28,.22));gravel_texture(gravel)
    timber=material('vacant yard silvered timber',(.39,.37,.30))
    iron=material('vacant yard rusted fixings',(.22,.16,.10),.4)
    grass=material('vacant yard dry weeds',(.30,.32,.17))
    # The footprint is strictly within the same 17m reserve as occupied lots.
    box('vacant yard ground',(0,0,.015),(16.4,16.4,.03),gravel)
    n=128;pixels=[]
    for y in range(n):
        for x in range(n):
            grain=math.sin(x*.72+math.sin(y*.05)*1.7)*.09+math.sin(x*2.2+y*.02)*.04
            tone=.89+grain
            pixels.extend((.39*tone,.37*tone,.30*tone,1))
    image=bpy.data.images.new('silvered fence grain',width=n,height=n);image.pixels=pixels;image.pack()
    tex=timber.node_tree.nodes.new('ShaderNodeTexImage');tex.image=image
    timber.node_tree.links.new(tex.outputs['Color'],timber.node_tree.nodes['Principled BSDF'].inputs['Base Color'])
    rng=random.Random(1958)
    for axis in (0,1):
        for side in (-1,1):
            def fence(name,along,z,dims,mat):
                xyz=(along,side*8,z) if axis==0 else (side*8,along,z)
                shape=dims if axis==0 else (dims[1],dims[0],dims[2])
                return box(name,xyz,shape,mat)
            for along in (-7.8,-4,0,4,7.8):fence('fence post',along,.88,(.16,.16,1.76),timber)
            for z in (.42,1.20):fence('fence rail',0,z,(15.7,.09,.12),timber)
            for i in range(52):
                along=-7.65+i*.3
                if (i+axis*7+(side+1)*3)%19==0:continue
                height=1.48+rng.uniform(-.12,.12)
                fence('weathered fence board',along,height/2,(.27,.07,height),timber)
                for z in (.42,1.20):fence('fence nail',along,z,(.025,.09,.025),iron)
    # Sparse clusters of folded grass blades, not a solid green carpet.
    vertices=[];faces=[]
    for i in range(240):
        x=rng.uniform(-7.5,7.5);y=rng.uniform(-7.5,7.5)
        if abs(x)<3 and abs(y)<3:continue
        for j in range(5):
            angle=rng.random()*math.tau;h=rng.uniform(.16,.48);w=.025
            start=len(vertices);dx=math.cos(angle);dy=math.sin(angle)
            vertices.extend([(x-w*dy,y+w*dx,.035),(x+w*dy,y-w*dx,.035),(x+.12*dx,y+.12*dy,h)])
            faces.append((start,start+1,start+2))
    mesh=bpy.data.meshes.new('vacant yard weed blades');mesh.from_pydata(vertices,[],faces);mesh.update()
    ob=bpy.data.objects.new('vacant yard weed blades',mesh);bpy.context.collection.objects.link(ob);mesh.materials.append(grass)


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


def street_tree(x, y, iron):
    bark=material('street tree bark',(.23,.18,.12))
    foliage=[material('street tree foliage '+str(i),color) for i,color in enumerate([(.18,.24,.12),(.25,.31,.16),(.32,.36,.20)])]
    for mat in foliage:
        n=128;pixels=[];base=mat.diffuse_color[:3]
        for row in range(n):
            v=row/(n-1)-.5
            for col in range(n):
                u=col/(n-1)
                midrib=math.exp(-abs(v)*100)
                vein=math.exp(-abs(math.sin((u-abs(v)*.65)*math.pi*9))*35)
                tone=.85+.1*math.sin(u*math.pi)+.12*midrib+.06*vein
                pixels.extend((*[channel*tone for channel in base],1))
        image=bpy.data.images.new(mat.name+' veins',width=n,height=n);image.pixels=pixels;image.pack()
        tree=mat.node_tree;tex=tree.nodes.new('ShaderNodeTexImage');tex.image=image
        tree.links.new(tex.outputs['Color'],tree.nodes['Principled BSDF'].inputs['Base Color'])
    # A small maintained tree: the entire crown fits a 1.9m pavement envelope.
    cylinder('tree soil',(x,y,.025),.39,.05,material('tree soil',(.12,.10,.065)),vertices=24)
    bpy.ops.mesh.primitive_torus_add(major_segments=24,minor_segments=6,location=(x,y,.06),major_radius=.42,minor_radius=.035)
    bpy.context.object.name='tree grate rim';bpy.context.object.data.materials.append(iron)
    for i in range(12):
        angle=i*math.tau/12
        beam('tree grate spoke',(x+.13*math.cos(angle),y+.13*math.sin(angle),.055),(x+.4*math.cos(angle),y+.4*math.sin(angle),.055),.035,iron)
    bpy.ops.mesh.primitive_cone_add(vertices=12,radius1=.12,radius2=.055,depth=3.0,location=(x,y,1.52))
    bpy.context.object.name='tree tapered trunk';bpy.context.object.data.materials.append(bark)
    for side in (-1,1):
        box('tree iron guard',(x+side*.24,y,.62),(.04,.04,1.2),iron)
        beam('tree guard brace',(x+side*.24,y,.9),(x+side*.08,y,1.08),.025,iron)
    rng=random.Random(1957)
    for i in range(10):
        angle=i*2.4;z=2.15+i*.15
        beam('tree branch',(x,y,z),(x+.53*math.cos(angle),y+.53*math.sin(angle),z+.65),.045,bark)
    vertices=[];faces=[];tones=[]
    for i in range(540):
        angle=rng.random()*math.tau;v=rng.uniform(-1,1)
        radius=.74*math.sqrt(1-v*v)*math.sqrt(rng.random())
        centre=Vector((x+radius*math.cos(angle),y+radius*math.sin(angle),3.55+v*1.12))
        azimuth=rng.random()*math.tau;tilt=rng.uniform(-.7,.7)
        along=Vector((math.cos(azimuth)*math.cos(tilt),math.sin(azimuth)*math.cos(tilt),math.sin(tilt)))
        across=Vector((-math.sin(azimuth),math.cos(azimuth),0))
        length=rng.uniform(.12,.19);width=length*.42
        # Folded six-point blades keep foliage irregular and catch soft light.
        points=[-along*length, -along*length*.35+across*width, along*length*.45+across*width, along*length, along*length*.45-across*width, -along*length*.35-across*width, Vector((0,0,.025))]
        start=len(vertices);vertices.extend(centre+point for point in points)
        for j in range(6):faces.append((start+j,start+(j+1)%6,start+6));tones.append(i%3)
    mesh=bpy.data.meshes.new('street tree leaf blades');mesh.from_pydata(vertices,[],faces);mesh.update()
    ob=bpy.data.objects.new('street tree leaf blades',mesh);bpy.context.collection.objects.link(ob)
    for mat in foliage:mesh.materials.append(mat)
    uv=mesh.uv_layers.new(name='leaf veins')
    coordinates=[(0,.5),(.325,1),(.725,1),(1,.5),(.725,0),(.325,0),(.5,.5)]
    for polygon,tone in zip(mesh.polygons,tones):
        polygon.material_index=tone
        for index in polygon.loop_indices:uv.data[index].uv=coordinates[mesh.loops[index].vertex_index%7]


def streetside():
    iron=material('street furniture iron',(.12,.16,.14),.6)
    wood=material('weathered bench timber',(.29,.19,.10))
    red=material('hydrant enamel',(.39,.075,.035),.35)
    zinc=material('galvanized bin',(.34,.36,.31),.65)
    # Rear pavement furnishing; tree crown remains inside the reserved band.
    street_tree(-6.8,-.65,iron)
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


def saint_agnes_interior():
    wood=material('cafe polished walnut',(.105,.047,.024))
    plaster=material('cafe aged ivory plaster',(.53,.48,.35))
    cream=material('cafe ivory mosaic',(.58,.55,.43))
    charcoal=material('cafe charcoal mosaic',(.036,.048,.045))
    brass=material('cafe aged brass',(.46,.29,.085),.75)
    leather=material('cafe oxblood upholstery',(.16,.026,.025))
    marble=material('cafe cream marble',(.66,.62,.48))
    porcelain=material('cafe porcelain',(.79,.76,.62))
    iron=material('cafe black enamel',(.025,.032,.031),.35)
    mirror=material('cafe silvered mirror',(.24,.32,.29),.8)
    glow=material('cafe opal lamps',(.94,.68,.30),0,1.6)
    # Packed walnut colour and normal grain, shared across panelling/furniture.
    n=128;image=bpy.data.images.new('cafe walnut grain',width=n,height=n);pixels=[]
    for y in range(n):
        for x in range(n):
            grain=.88+.12*math.sin(x*.7+math.sin(y*.045)*2)+.04*math.sin(x*2.8)
            pixels.extend((*[1.055*(c*grain)**(1/2.4)-.055 for c in wood.diffuse_color[:3]],1))
    image.pixels=pixels;image.pack();tex=wood.node_tree.nodes.new('ShaderNodeTexImage');tex.image=image
    wood.node_tree.links.new(tex.outputs['Color'],wood.node_tree.nodes['Principled BSDF'].inputs['Base Color'])
    wood.node_tree.nodes['Principled BSDF'].inputs['Roughness'].default_value=.34
    leather.node_tree.nodes['Principled BSDF'].inputs['Roughness'].default_value=.43
    box('cafe floor slab',(0,0,-.13),(12,10,.26),charcoal)
    for ix in range(24):
        for iy in range(20):
            box('mosaic tile',(-5.75+ix*.5,-4.75+iy*.5,.014),(.485,.485,.028),cream if (ix+iy)%2 else charcoal)
    box('back plaster wall',(0,-5,1.75),(12,.2,3.5),plaster)
    box('left plaster wall',(-6,0,1.75),(.2,10,3.5),plaster)
    for x in [-5.5+i for i in range(12)]:
        box('back walnut panel',(x,-4.85,.65),(.91,.12,1.18),wood,.022)
        for z in (.18,1.12):box('panel inset moulding',(x,-4.77,z),(.77,.035,.035),brass,.01)
    for y in [-4.5+i for i in range(10)]:
        box('left walnut panel',(-5.85,y,.65),(.12,.91,1.18),wood,.022)
    for z in (.12,1.3,3.36):
        box('back dado and cornice',(0,-4.75,z),(12,.18,.12),wood,.018)
        box('left dado and cornice',(-5.75,0,z),(.18,10,.12),wood,.018)
    # Back bar and mirrored display, with layered frames and bottle shelves.
    box('back bar cupboard',(1,-4.25,.55),(7,1.1,1.1),wood,.045)
    for x in (-2,-.5,1,2.5,4):
        box('cupboard raised door',(x,-3.67,.55),(1.35,.05,.87),wood,.018)
        cylinder('cupboard brass pull',(x+.44,-3.61,.6),.035,.055,brass,(math.pi/2,0,0),16)
    box('silvered back bar mirror',(1,-4.82,2.25),(6.8,.05,1.65),mirror)
    for x in (-2.5,4.5):box('mirror carved stile',(x,-4.7,2.25),(.13,.16,1.9),wood,.025)
    for z in (1.33,3.18):box('mirror cornice',(1,-4.7,z),(7.15,.2,.16),wood,.025)
    bottles=[material('bottle green',(.022,.12,.052),.15),material('bottle amber',(.25,.105,.021),.12)]
    label=material('bottle paper labels',(.59,.52,.34))
    for z in (1.4,2.18):
        box('bottle shelf',(1,-4.3,z),(6.8,.7,.065),wood,.012)
        for i in range(22):
            x=-2.15+i*.30;m=bottles[i%2];h=.32+(i%3)*.04
            cylinder('bottle body',(x,-4.22,z+h/2+.05),.075,h,m,vertices=16)
            cylinder('bottle neck',(x,-4.22,z+h+.11),.033,.14,m,vertices=16)
            cylinder('bottle label',(x,-4.22,z+.18),.076,.12,label,vertices=16)
    box('front bar panel',(1,-2.65,.53),(7, .7,1.06),wood,.055)
    box('rounded marble counter',(1,-2.65,1.12),(7.25,1.15,.13),marble,.055)
    for x in (-2,-.5,1,2.5,4):
        box('bar inset field',(x,-2.27,.57),(1.3,.05,.67),wood,.018)
    cylinder('brass foot rail',(1,-1.95,.23),.037,6.8,brass,(0,math.pi/2,0),24)
    for x in (-2,0,2,4):
        cylinder('foot rail post',(x,-2.05,.14),.025,.28,brass,vertices=16)
        cylinder('stool pedestal',(x,-1.4,.35),.065,.65,brass,vertices=24)
        cylinder('stool foot',(x,-1.4,.055),.24,.08,iron,vertices=32)
        cylinder('round padded stool',(x,-1.4,.72),.26,.16,leather,vertices=32)
    # Two deeply upholstered booths beside the wall, with clear central aisle.
    for y in (.5,3.2):
        for side in (-1,1):
            yy=y+side*.72
            box('booth walnut plinth',(-4.65,yy,.24),(1.8,.6,.48),wood,.04)
            box('booth seat cushion',(-4.65,yy,.53),(1.8,.65,.19),leather,.07)
            box('booth padded back',(-4.65,yy+side*.25,.96),(1.8,.19,.85),leather,.065)
            for x in (-5.3,-4.85,-4.4,-3.95):
                cylinder('upholstery button',(x,yy+side*.14,1.08),.025,.018,brass,(math.pi/2,0,0),12)
        box('booth marble table',(-4.65,y,.83),(1.6,.8,.09),marble,.04)
        cylinder('table cast pedestal',(-4.65,y,.4),.07,.78,iron,vertices=24)
        cylinder('table base',(-4.65,y,.07),.3,.10,iron,vertices=24)
        cylinder('sugar bowl',(-4.65,y,.93),.09,.12,porcelain,vertices=24)
        for x in (-5.15,-4.15):
            cylinder('coffee saucer',(x,y,.893),.115,.025,porcelain,vertices=32)
            cylinder('coffee cup',(x,y,.97),.065,.13,porcelain,vertices=32)
    # Period espresso boiler, register and service ware.
    copper=material('cafe copper boiler',(.43,.19,.073),.8)
    cylinder('espresso boiler',(3,-4,1.54),.25,.76,copper,vertices=32)
    cylinder('boiler lid',(3,-4,1.94),.28,.06,brass,vertices=32)
    for x in (2.82,3.18):cylinder('espresso tap',(x,-3.72,1.5),.035,.2,brass,(math.pi/2,0,0),16)
    box('cash register base',(-1.7,-2.6,1.29),(.65,.55,.23),brass,.04)
    box('cash register head',(-1.7,-2.75,1.57),(.58,.25,.36),brass,.04)
    for row in range(3):
        for col in range(7):cylinder('register ivory key',(-1.94+col*.08,-2.43-row*.08,1.44),.025,.045,porcelain,vertices=12)
    for x in (-.5,1,2):
        cylinder('counter saucer',(x,-2.35,1.2),.12,.024,porcelain,vertices=24)
        cylinder('counter cup',(x,-2.35,1.28),.067,.13,porcelain,vertices=24)
    for x in (-3,1,4):
        cylinder('pendant suspension',(x,-2.7,3.2),.016,.55,brass,vertices=12)
        cylinder('opal pendant',(x,-2.7,2.85),.24,.27,glow,vertices=32)
        cylinder('pendant brass shade',(x,-2.7,3),.30,.06,brass,vertices=32)
    # Framed local prints on the booth wall.
    printmat=material('cafe sepia print',(.26,.22,.14))
    for y in (.5,3.2):
        box('picture frame',(-5.76,y,2.25),(.10,1.5,1.1),wood,.03)
        box('picture mount',(-5.69,y,2.25),(.03,1.3,.91),porcelain)
        box('sepia picture',(-5.66,y,2.25),(.02,1.06,.68),printmat)
    # Retain wall assemblies independently for browser camera cutaways.
    walls={}
    for side in ('left','back'):
        group=bpy.data.objects.new('interior-wall-'+side,None)
        bpy.context.collection.objects.link(group);walls[side]=group
    for ob in list(bpy.context.scene.objects):
        if ob.type!='MESH':continue
        if ob.name.startswith(('left plaster','left walnut','left dado','picture frame','picture mount','sepia picture')):
            ob.parent=walls['left']
        elif ob.name.startswith(('back plaster','back walnut','back dado','panel inset','silvered back','mirror carved','mirror cornice')):
            ob.parent=walls['back']


def fire_engine():
    paint=material('fire engine red enamel',(.39,.035,.018),.45)
    chrome=material('fire engine bright metal',(.58,.60,.57),.82)
    rubber=material('fire engine rubber',(.024,.026,.025))
    glass=material('fire engine glazing',(.11,.19,.21),.3)
    brass=material('pump brass',(.51,.34,.11),.7)
    canvas=material('woven hose canvas',(.46,.40,.27))
    wood=material('ladder varnished ash',(.39,.22,.085))
    lamp=material('engine headlamps',(.94,.83,.56),0,.4)
    box('truck frame',(0,0,.53),(1.8,5.4,.25),rubber,.05)
    box('rounded bonnet',(0,-1.86,1.13),(1.65,1.6,.72),paint,.22)
    box('cab lower',(0,-.54,1.0),(2.02,1.38,.78),paint,.14)
    box('split windscreen',(0,-.72,1.73),(1.8,1.13,.76),glass,.07)
    box('cab roof',(0,-.58,2.17),(2.05,1.42,.14),paint,.09)
    box('windscreen divider',(0,-1.30,1.79),(.055,.045,.7),chrome)
    for side in (-1,1):
        box('cab rear pillar',(side*.93,.01,1.77),(.10,.10,.7),paint)
        box('door handle',(side*1.025,-.25,1.26),(.04,.24,.035),chrome,.01)
        box('running board',(side*.99,.05,.54),(.28,3.7,.12),chrome,.03)
        box('rear equipment locker',(side*.70,1.40,1.0),(.62,2.32,.75),paint,.06)
        for y in (.65,1.45,2.15):
            box('locker inset',(side*1.02,y,1.04),(.025,.66,.50),chrome,.018)
            box('locker latch',(side*1.044,y,1.06),(.025,.13,.035),brass)
        for axle,y in enumerate((-1.83,1.87)):
            cylinder('truck tyre',(side*.94,y,.52),.5,.24,rubber,(0,math.pi/2,0),48)
            cylinder('red wheel',(side*1.067,y,.52),.34,.025,paint,(0,math.pi/2,0),40)
            cylinder('chrome wheel hub',(side*1.086,y,.52),.16,.045,chrome,(0,math.pi/2,0),32)
            for bolt in range(8):
                angle=bolt*math.tau/8
                cylinder('wheel lug',(side*1.115,y+math.sin(angle)*.245,.52+math.cos(angle)*.245),.025,.02,chrome,(0,math.pi/2,0),10)
        for rail in (-1,1):box('ladder rail',(side*.82+rail*.14,.70,1.82),(.05,3.75,.07),wood,.008)
        for rung in range(14):box('ladder rung',(side*.82,-1.05+rung*.26,1.82),(.32,.035,.04),wood,.008)
        cylinder('round headlamp',(side*.70,-2.70,1.10),.17,.13,chrome,(math.pi/2,0,0),32)
        cylinder('headlamp lens',(side*.70,-2.78,1.10),.135,.02,lamp,(math.pi/2,0,0),32)
        # Side pump controls ahead of the hose bed.
        cylinder('pressure gauge',(side*1.04,.27,1.17),.09,.055,chrome,(0,math.pi/2,0),24)
        for y in (.27,.52):cylinder('hose outlet',(side*1.065,y,.85),.105,.12,brass,(0,math.pi/2,0),24)
    for y in (-2.77,2.77):box('truck bumper',(0,y,.57),(2.12,.14,.17),chrome,.045)
    box('radiator',(0,-2.68,1.15),(1.10,.08,.69),rubber,.04)
    for rib in range(13):box('radiator chrome rib',(-.49+rib*.082,-2.735,1.15),(.025,.025,.61),chrome)
    box('hose bed',(0,1.4,1.11),(.76,2.2,.20),rubber)
    for row in range(5):
        for layer in range(3):box('folded hose',(-.28+row*.14,1.4,1.28+layer*.07),(.115,2.05,.06),canvas,.024)
    cylinder('beacon base',(0,-.65,2.29),.19,.08,chrome,vertices=32)
    cylinder('red rotating beacon',(0,-.65,2.43),.16,.22,material('engine beacon',(.8,.025,.01),0,1),vertices=32)


def fire_nozzle():
    brass=material('nozzle brass',(.48,.31,.10),.72)
    rubber=material('nozzle grip',(.035,.038,.031))
    cylinder('hose coupling',(0,.12,0),.065,.10,brass,(math.pi/2,0,0),32)
    cylinder('nozzle barrel',(0,-.08,0),.043,.32,brass,(math.pi/2,0,0),32)
    cylinder('nozzle hand grip',(0,.015,0),.049,.14,rubber,(math.pi/2,0,0),32)
    cylinder('nozzle tip',(0,-.255,0),.031,.065,brass,(math.pi/2,0,0),24)
    box('shutoff lever',(0,-.035,.065),(.12,.025,.025),brass,.008)
    tip=bpy.data.objects.new('water-outlet',None);bpy.context.collection.objects.link(tip);tip.location=(0,-.29,0)


def firefighter():
    person(False)
    for name in ('headwear-fedora','headwear-cap','hair-receding'):
        group=bpy.data.objects.get(name)
        if group:
            for child in list(group.children_recursive):bpy.data.objects.remove(child,do_unlink=True)
            bpy.data.objects.remove(group,do_unlink=True)
    canvas=bpy.data.materials['wool suit'];canvas.name='fire brigade dark woven coat'
    image=canvas.node_tree.nodes['Principled BSDF'].inputs['Base Color'].links[0].from_node.image
    pixels=list(image.pixels)
    for i in range(0,len(pixels),4):
        pixels[i]*=.72;pixels[i+1]*=.67;pixels[i+2]*=.52
    image.pixels=pixels;image.pack()
    leather=material('fire helmet leather',(.055,.041,.025))
    brass=material('fire coat brass',(.47,.32,.10),.65)
    cylinder('helmet swept brim',(0,.035,1.78),.27,.04,leather,vertices=40)
    bpy.ops.mesh.primitive_uv_sphere_add(segments=32,ring_count=16,location=(0,0,1.86))
    dome=bpy.context.object;dome.name='leather helmet crown';dome.scale=(.195,.205,.125)
    bpy.ops.object.transform_apply(location=False,rotation=False,scale=True);dome.data.materials.append(leather)
    for polygon in dome.data.polygons:polygon.use_smooth=True
    box('helmet crest',(0,0,1.975),(.045,.32,.035),leather,.014)
    box('helmet front shield',(0,-.21,1.855),(.12,.025,.14),brass,.02)
    box('coat storm flap',(0,-.17,1.14),(.10,.045,.40),canvas,.01)
    for z in (1.02,1.14,1.26):box('coat brass clasp',(.035,-.201,z),(.1,.02,.025),brass,.004)
    for side in (-1,1):box('coat pocket',(side*.16,-.173,.94),(.15,.035,.13),canvas,.01)
    for name in ('ivory shirt','wine silk tie'):
        mat=bpy.data.materials.get(name)
        if mat:
            mat.diffuse_color=(.09,.085,.065,1)
            mat.node_tree.nodes['Principled BSDF'].inputs['Base Color'].default_value=mat.diffuse_color


def police_officer():
    person(False)
    # A distinct uniform and peaked cap, retaining the articulated cast rig.
    for name in ('headwear-fedora','headwear-cap','hair-receding'):
        group=bpy.data.objects.get(name)
        if group:
            for child in list(group.children_recursive):bpy.data.objects.remove(child,do_unlink=True)
            bpy.data.objects.remove(group,do_unlink=True)
    serge=bpy.data.materials['wool suit'];serge.name='police woven navy serge'
    # Bake the uniform pigment into the packed wool image for glTF fidelity.
    shader=serge.node_tree.nodes['Principled BSDF']
    image=shader.inputs['Base Color'].links[0].from_node.image
    pixels=list(image.pixels)
    for i in range(0,len(pixels),4):
        pixels[i]*=.85;pixels[i+1]*=.95;pixels[i+2]*=1.25
    image.pixels=pixels;image.pack()
    navy=material('police cap navy',(.025,.04,.075))
    brass=material('police badge brass',(.62,.43,.12),.8)
    leather=material('police duty leather',(.018,.02,.018))
    cylinder('uniform cap band',(0,0,1.8),.15,.07,navy,vertices=32)
    cylinder('uniform cap crown',(0,0,1.86),.18,.07,navy,vertices=32)
    peak=box('polished cap visor',(0,-.16,1.78),(.3,.19,.035),leather,.03)
    cylinder('cap badge',(0,-.154,1.83),.034,.014,brass,(math.pi/2,0,0),16)
    box('duty belt',(0,0,.96),(.435,.30,.075),leather,.015)
    box('belt buckle',(0,-.159,.96),(.065,.018,.065),brass,.009)
    box('closed leather holster',(.255,.02,.89),(.08,.14,.23),leather,.025)
    box('utility pouch',(-.24,.03,.93),(.075,.15,.14),leather,.018)
    for x in (-.13,.13):
        box('uniform breast pocket',(x,-.18,1.26),(.11,.028,.13),navy,.008)
        box('pocket flap',(x,-.197,1.31),(.12,.015,.04),navy,.006)
    for z in (1.12,1.04):cylinder('brass tunic button',(0,-.162,z),.016,.014,brass,(math.pi/2,0,0),16)
    cylinder('shield badge',(-.13,-.214,1.31),.038,.012,brass,(math.pi/2,0,0),6)
    for side in (-1,1):
        box('shoulder epaulette',(side*.23,0,1.405),(.10,.20,.026),navy,.01)


def bar_cloth():
    linen=material('washed bar linen',(.65,.61,.49))
    stripe=material('woven blue border',(.20,.29,.32))
    box('folded wiping cloth',(0,0,.005),(.19,.16,.01),linen,.004)
    for x in (-.073,.073):box('cloth border',(x,0,.0105),(.009,.145,.002),stripe,.001)
    for i in range(4):box('soft linen fold',(-.05+i*.035,0,.011),(.017,.14,.002),linen,.001)


def slot_cabinet():
    """1950s electromechanical cabinet; curved reel paper receives public symbols."""
    red=material('slot oxblood enamel',(.25,.035,.025),.22)
    cream=material('slot ivory enamel',(.72,.65,.47),.18)
    brass=material('slot satin brass',(.55,.35,.10),.8)
    chrome=material('slot polished steel',(.55,.59,.56),.9)
    dark=material('slot black bakelite',(.018,.023,.018),.15)
    glass=material('slot opal crown',(.80,.64,.31),0,.3)
    box('cabinet foot',(0,0,.07),(1.02,.72,.14),dark,.04)
    box('cabinet body',(0,.08,.84),(.95,.62,1.50),red,.10)
    box('lower face',(0,-.28,.49),(.91,.12,.66),cream,.07)
    box('crown face',(0,-.25,1.39),(.91,.15,.35),cream,.09)
    box('top crown',(0,.01,1.62),(.80,.58,.16),red,.08)
    box('crown name plaque',(0,-.338,1.43),(.73,.023,.17),glass,.025)
    def label(name,text,xyz,size,mat=dark):
        curve=bpy.data.curves.new(name,'FONT');curve.body=text;curve.size=size;curve.align_x='CENTER';curve.extrude=.0005
        ob=bpy.data.objects.new(name,curve);bpy.context.collection.objects.link(ob);ob.location=xyz;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(mat)
        bpy.ops.object.select_all(action='DESELECT');ob.select_set(True);bpy.context.view_layer.objects.active=ob;bpy.ops.object.convert(target='MESH');ob.select_set(False)
    label('machine title','LUCKY BELL',(0,-.355,1.40),.075)
    label('payout label','PAY OUT',(0,-.350,.69),.035)
    label('coin label','INSERT COIN',(.17,-.356,.49),.026)
    for x in (-.455,.455):box('front bright trim',(x,-.35,.96),(.035,.035,1.14),chrome,.014)
    # Curved paper strips, each its own material so the browser can print the
    # backend's strip and animate its travel without inventing a symbol.
    for i,x in enumerate((-.28,0,.28)):
        paper=material(f'slot-reel-{i}',(.90,.87,.72))
        verts=[];faces=[];steps=24
        for row in range(steps+1):
            a=-.65+row/steps*1.30
            for side in (-1,1):verts.append((x+side*.115,-.29-.14*math.cos(a),1.055+.28*math.sin(a)))
        for row in range(steps):faces.append((row*2,row*2+1,row*2+3,row*2+2))
        mesh=bpy.data.meshes.new(f'reel paper {i}');mesh.from_pydata(verts,[],faces);mesh.materials.append(paper)
        uv=mesh.uv_layers.new(name='Reel print')
        for poly in mesh.polygons:
            for index in poly.loop_indices:
                v=mesh.loops[index].vertex_index;uv.data[index].uv=(v%2,v//2/steps)
        ob=bpy.data.objects.new(f'reel-paper-{i}',mesh);bpy.context.collection.objects.link(ob)
        for poly in mesh.polygons:poly.use_smooth=True
        for edge in (-1,1):box('window side trim',(x+edge*.128,-.40,1.055),(.025,.055,.37),brass,.008)
    for z in (.875,1.235):box('window rail',(0,-.415,z),(.87,.055,.04),chrome,.012)
    for x in (-.445,.445):
        ob=box('payline arrow',(x,-.442,1.055),(.065,.008,.018),red,.004)
    # Recessed tray has a floor and raised rim, rather than a solid rectangle.
    box('tray dark recess',(-.1,-.351,.31),(.67,.035,.24),dark,.025)
    box('tray bottom',(-.1,-.46,.21),(.70,.27,.035),chrome,.016)
    for x in (-.455,.255):box('tray side',(x,-.46,.255),(.03,.27,.11),chrome,.012)
    box('tray lip',(-.1,-.60,.255),(.72,.035,.10),chrome,.013)
    box('coin surround',(.23,-.35,.56),(.16,.025,.09),chrome,.01)
    box('coin slot',(.23,-.369,.56),(.09,.008,.018),dark,.005)
    for x in (-.4,.4):
        for z in (.40,1.32):cylinder('face screw',(x,-.37,z),.012,.008,chrome,(math.pi/2,0,0),12)
    coins=bpy.data.objects.new('slot-coins',None);bpy.context.collection.objects.link(coins)
    for i in range(9):
        coin=cylinder('payout coin',(-.31+(i%3)*.17,-.53+(i//3)*.073,.239+(i%2)*.006),.038,.009,brass,vertices=32);coin.parent=coins
    # A proper pivot is retained for the handle's pull animation.
    lever=bpy.data.objects.new('slot-lever',None);bpy.context.collection.objects.link(lever);lever.location=(.56,.05,1.0)
    cylinder('lever axle',(.53,.05,1.0),.075,.18,chrome,(0,math.pi/2,0),32)
    parts=[beam('lever stem',(.63,.05,1.0),(.63,.05,1.47),.025,chrome),cylinder('lever hub',(.63,.05,1.0),.055,.06,brass,(0,math.pi/2,0),24)]
    bpy.ops.mesh.primitive_uv_sphere_add(segments=24,ring_count=12,radius=.07,location=(.63,.05,1.48));knob=bpy.context.object;knob.name='lever knob';knob.data.materials.append(dark);parts.append(knob)
    bpy.context.view_layer.update()
    for ob in parts:
        world=ob.matrix_world.copy();ob.parent=lever;ob.matrix_world=world


def laundry_interior():
    """Bluebird's working floor: belt-era drum washers, sorting and collection."""
    cream=material('Bluebird ivory enamel',(.66,.65,.54),.18)
    teal=material('Bluebird machine green',(.10,.24,.22),.25)
    iron=material('Bluebird dark cast iron',(.035,.045,.044),.55)
    steel=material('Bluebird polished rims',(.42,.46,.43),.82)
    brass=material('Bluebird valves',(.45,.29,.09),.7)
    plaster=material('Bluebird limewash',(.52,.49,.38))
    tile=material('Bluebird worn floor',(.29,.31,.27))
    wood=material('Bluebird scrubbed beech',(.38,.25,.12))
    linen=material('Bluebird cotton weave',(.68,.65,.52))
    ink=material('Bluebird printed ink',(.035,.075,.065))
    # Packed fine weave survives the browser export, including at close zoom.
    n=128;pixels=[];rng=random.Random(1953)
    for y in range(n):
        for x in range(n):
            tone=.95+(.04 if (x+y)%2 else -.04)+rng.uniform(-.018,.018)
            pixels.extend((*[1.055*(c*tone)**(1/2.4)-.055 for c in linen.diffuse_color[:3]],1))
    img=bpy.data.images.new('Bluebird cotton weave',width=n,height=n);img.pixels=pixels;img.pack()
    tex=linen.node_tree.nodes.new('ShaderNodeTexImage');tex.image=img
    linen.node_tree.links.new(tex.outputs['Color'],linen.node_tree.nodes['Principled BSDF'].inputs['Base Color'])
    box('floor foundation',(0,1,-.13),(10,10,.25),iron)
    for x in range(20):
        for y in range(20):box('quarry tile',(-4.75+x*.5,-3.75+y*.5,0),(.488,.488,.035),tile,.003)
    box('rear plaster',(0,6,2.1),(10,.2,4.2),plaster)
    box('west plaster',(-5,1,2.1),(.2,10,4.2),plaster)
    for z in range(6):
        for x in range(20):box('rear glazed tile',(-4.75+x*.5,5.88,.15+z*.3),(.485,.045,.285),cream,.004)
        for y in range(20):box('west glazed tile',(-4.88,-3.75+y*.5,.15+z*.3),(.045,.485,.285),cream,.004)
    for z in (.12,1.87):
        box('rear green border',(0,5.83,z),(10,.06,.08),teal)
        box('west green border',(-4.83,1,z),(.06,10,.08),teal)
    def ring(name,xyz,major,minor,mat):
        bpy.ops.mesh.primitive_torus_add(major_radius=major,minor_radius=minor,major_segments=48,minor_segments=10,location=xyz,rotation=(math.pi/2,0,0))
        ob=bpy.context.object;ob.name=name;ob.data.materials.append(mat)
        for p in ob.data.polygons:p.use_smooth=True
    def label(name,text,xyz,size):
        curve=bpy.data.curves.new(name,'FONT');curve.body=text;curve.size=size;curve.align_x='CENTER';curve.extrude=.001
        ob=bpy.data.objects.new(name,curve);bpy.context.collection.objects.link(ob);ob.location=xyz;ob.rotation_euler=(math.pi/2,0,0);ob.data.materials.append(ink)
        bpy.ops.object.select_all(action='DESELECT');ob.select_set(True);bpy.context.view_layer.objects.active=ob;bpy.ops.object.convert(target='MESH');ob.select_set(False)
    label('rear bluebird sign','BLUEBIRD LAUNDRY',(0,5.85,3.3),.43)
    label('rear service sign','WASH  /  PRESS  /  FOLD',(0,5.84,2.91),.19)
    # Three front-loading industrial washers face a reserved service aisle.
    for i,x in enumerate((-3.35,-1.35,.65)):
        box('washer plinth',(x,4.55,.12),(1.65,1.5,.24),iron,.06)
        box('washer enamel cabinet',(x,4.55,.99),(1.55,1.35,1.5),cream,.10)
        box('washer green crown',(x,4.55,1.78),(1.58,1.39,.18),teal,.06)
        cylinder('washer dark drum',(x,3.86,.94),.49,.04,iron,(math.pi/2,0,0),48)
        ring('washer polished door rim',(x,3.81,.94),.49,.047,steel)
        ring('washer rubber seal',(x,3.805,.94),.421,.017,iron)
        drum=bpy.data.objects.new(f'laundry-drum-{i}',None);bpy.context.collection.objects.link(drum);drum.location=(x,3.815,.94)
        moving=[]
        for j in range(16):
            a=j*math.tau/16
            moving.append(cylinder('drum perforation',(x+math.sin(a)*.33,3.827,.94+math.cos(a)*.33),.026,.015,steel,(math.pi/2,0,0),8))
        for j in range(3):
            ob=box('linen inside drum',(x-.16+j*.14,3.815,.80+j*.10),(.20,.028,.15),linen,.04);ob.rotation_euler.y=j*.45;moving.append(ob)
        bpy.context.view_layer.update()
        for ob in moving:
            world=ob.matrix_world.copy();ob.parent=drum;ob.matrix_world=world
        box('washer door hinge',(x-.54,3.79,.94),(.10,.13,.28),steel,.025)
        beam('washer locking handle',(x+.52,3.74,.82),(x+.52,3.74,1.07),.045,brass)
        cylinder('washer selector dial',(x-.44,3.83,1.53),.075,.05,iron,(math.pi/2,0,0),24)
        label('washer number',f'No. {i+1}',(x+.19,3.855,1.51),.11)
        beam('washer supply riser',(x,5.35,.3),(x,5.35,2.3),.045,steel)
        beam('washer fill elbow',(x,5.35,2.3),(x,4.8,2.3),.045,steel)
        cylinder('washer valve wheel',(x,5.25,2.17),.09,.025,brass,(math.pi/2,0,0),16)
    beam('rear steam main',(-4.6,5.6,2.5),(4.6,5.6,2.5),.085,steel)
    for x in (-4.6,4.6):beam('rear steam riser',(x,5.6,.15),(x,5.6,3.9),.085,steel)
    # Collection counter on the right, leaving the centre and entrance open.
    box('folding counter',(3.35,1.1,.51),(2.6,1,.98),teal,.035)
    box('scrubbed folding top',(3.35,1.1,1.04),(2.75,1.14,.10),wood,.035)
    for x in (2.55,3.35,4.15):box('counter recessed panel',(x,.585,.52),(.65,.025,.67),cream,.02)
    for i in range(5):
        box('folded sheets',(2.7,1.05,1.115+i*.055),(.65,.50,.05),linen,.025)
        box('linen blue binding',(2.7,1.05,1.143+i*.055),(.052,.51,.008),teal)
    box('linen under hands',(3.35,1.49,1.100),(.64,.32,.014),linen,.006)
    for x in (3.06,3.64):box('linen stitched hem',(x,1.49,1.108),(.012,.30,.002),teal)
    box('collection ticket book',(3.85,.85,1.115),(.35,.27,.05),linen,.01)
    label('counter collection sign','COLLECTION',(3.35,.557,.65),.13)
    # Rear shelving beside the washers, with individually wrapped bundles.
    for z in (.3,1.05,1.8,2.55):box('linen shelf',(3.35,5.12,z),(2.3,.8,.07),wood,.018)
    for x in (2.17,4.53):box('shelf upright',(x,5.12,1.4),(.09,.85,2.8),teal,.014)
    for row in range(3):
        for col in range(3):
            x=2.62+col*.73;z=.48+row*.75
            box('wrapped laundry bundle',(x,5.10,z),(.61,.62,.27),linen,.055)
            box('bundle tie',(x,4.782,z),(.025,.012,.28),wood)
    # A low waiting bench along the left wall, facing the public floor.
    for y in (-1.8,.8):
        for x in (-4.5,-3.85):box('waiting bench leg',(x,y,.35),(.10,.10,.70),iron,.015)
    for x in (-4.52,-4.30,-4.08,-3.86):box('waiting bench slat',(x,-.5,.72),(.19,3.1,.10),wood,.02)
    for z in (1,1.22,1.44):box('waiting bench back',(-4.65,-.5,z),(.12,3.15,.16),teal,.02)
    # Ceiling-free room keeps suspended opal work lights visible in cutaway.
    glow=material('Bluebird opal shades',(.78,.76,.59),0,.7)
    for x in (-2.5,2.5):
        cylinder('work lamp stem',(x,1.5,3.45),.018,.55,iron)
        cylinder('work lamp shade',(x,1.5,3.14),.33,.09,teal,vertices=40)
        cylinder('work lamp opal',(x,1.5,3.08),.24,.055,glow,vertices=32)
    walls={}
    for side in ('left','back'):
        group=bpy.data.objects.new('interior-wall-'+side,None);bpy.context.collection.objects.link(group);walls[side]=group
    for ob in list(bpy.context.scene.objects):
        if ob.type!='MESH':continue
        if ob.name.startswith('west '):ob.parent=walls['left']
        elif ob.name.startswith('rear '):ob.parent=walls['back']


def mariner_lobby():
    """Boarding-house reception: 24 keys, rent book, waiting bench and stairs."""
    plaster=material('Mariner tobacco cream plaster',(.52,.46,.34))
    wood=material('Mariner worn oak',(.22,.115,.055))
    green=material('Mariner painted dado',(.12,.21,.18))
    brass=material('Mariner reception brass',(.50,.35,.13),.65)
    iron=material('Mariner black iron',(.035,.045,.04),.45)
    paper=material('Mariner ledger paper',(.70,.65,.50))
    linen=material('Mariner folded linen',(.60,.57,.47))
    amber=material('Mariner opal shades',(.85,.68,.38),0,.8)
    # Deterministic fine oak grain, carried by exported UVs rather than nodes
    # that only render inside Blender.
    rng=random.Random(1963);n=128;pixels=[]
    for y in range(n):
        for x in range(n):
            grain=.9+.09*math.sin(y*.62+math.sin(x*.06)*1.4)+rng.uniform(-.04,.04)
            pixels.extend((*[1.055*(c*grain)**(1/2.4)-.055 for c in wood.diffuse_color[:3]],1))
    img=bpy.data.images.new('Mariner oak grain',width=n,height=n);img.pixels=pixels;img.pack()
    tex=wood.node_tree.nodes.new('ShaderNodeTexImage');tex.image=img
    wood.node_tree.links.new(tex.outputs['Color'],wood.node_tree.nodes['Principled BSDF'].inputs['Base Color'])
    box('foundation',(0,1,-.14),(10,10,.25),wood)
    for row in range(40):
        for col in range(5):
            box('oak floorboard',(-4+col*2,-3.875+row*.25,0),(1.985,.242,.035),wood,.004)
    box('rear plaster wall',(0,6,2.15),(10,.22,4.3),plaster)
    box('west plaster wall',(-5,1,2.15),(.22,10,4.3),plaster)
    for x in range(-19,20):box('rear dado board',(x*.25,5.85,.65),(.242,.08,1.3),green)
    for y in range(-15,24):box('west dado board',(-4.85,y*.25,.65),(.08,.242,1.3),green)
    for z in (.1,1.35,4.1):
        box('rear moulding',(0,5.78,z),(10,.15,.12),wood,.025)
        box('west moulding',(-4.78,1,z),(.15,10,.12),wood,.025)
    # The desk faces the public floor, leaving an actual service aisle behind.
    box('reception desk',(-2.5,4.35,.54),(3.3,.82,1.08),green,.025)
    box('reception oak top',(-2.5,4.35,1.12),(3.5,1,.10),wood,.04)
    for x in (-3.65,-2.5,-1.35):
        box('reception recessed panel',(x,3.925,.56),(.9,.035,.72),wood,.025)
        box('reception panel inset',(x,3.90,.56),(.72,.025,.54),green,.016)
    cylinder('service bell base',(-1.4,4.17,1.20),.13,.05,iron,vertices=24)
    cylinder('service bell dome',(-1.4,4.17,1.255),.10,.08,brass,vertices=32)
    cylinder('service bell button',(-1.4,4.17,1.31),.035,.04,brass,vertices=16)
    box('rent ledger cover',(-2.8,4.1,1.19),(.66,.43,.035),iron,.008)
    box('open rent ledger',(-2.8,4.1,1.214),(.62,.40,.016),paper,.005)
    for i in range(9):box('ledger ruled line',(-2.8,3.94+i*.039,1.224),(.56,.002,.002),green)
    box('ledger spine',(-2.8,4.1,1.226),(.009,.40,.003),wood)
    beam('ledger pencil',(-2.39,3.94,1.185),(-2.36,4.24,1.185),.014,brass)
    box('rear key board',(-2.5,5.7,2.3),(3.7,.12,1.62),wood,.03)
    for row in range(4):
        for col in range(6):
            x=-4+col*.6;z=1.76+row*.36
            cylinder('rear key hook',(x,5.58,z),.022,.13,brass,(math.pi/2,0,0),12)
            box('rear numbered key tag',(x,5.50,z-.09),(.14,.026,.18),paper,.01)
            cylinder('rear hanging key stem',(x+.09,5.49,z-.11),.012,.14,brass,vertices=8)
    # Wall labels remain physical lettering and are part of the cutaway.
    def label(name,text,x,y,z,size):
        curve=bpy.data.curves.new(name,'FONT');curve.body=text;curve.size=size;curve.align_x='CENTER';curve.extrude=.001
        obj=bpy.data.objects.new(name,curve);bpy.context.collection.objects.link(obj);obj.location=(x,y,z);obj.rotation_euler=(math.pi/2,0,0);obj.data.materials.append(iron)
        bpy.ops.object.select_all(action='DESELECT');obj.select_set(True);bpy.context.view_layer.objects.active=obj;bpy.ops.object.convert(target='MESH');obj.select_set(False)
    label('rear house name','THE MARINER',-2.5,5.72,3.46,.36)
    for row in range(4):
        for col in range(6):label('rear room number',str(row*6+col+1),-4+col*.6,5.477,1.63+row*.36,.065)
    # Corridor and linen cupboard, separate from the staircase clear space.
    box('rear corridor door',(.45,5.83,1.35),(1.55,.12,2.7),wood,.03)
    for x in (-.4,1.3):box('rear door jamb',(x,5.73,1.42),(.13,.19,2.84),green,.02)
    box('rear door lintel',(.45,5.73,2.85),(1.84,.19,.13),green,.02)
    box('rear door glass',(.45,5.745,1.93),(1.13,.035,.72),paper,.02)
    label('rear corridor sign','ROOMS',.45,5.72,1.89,.16)
    cylinder('rear door handle',(1.01,5.69,1.12),.045,.08,brass,(math.pi/2,0,0),16)
    for step in range(14):
        y=.2+step*.4;h=(step+1)*.2
        box('stair tread',(3.6,y,h/2),(2.25,.4,h),wood,.016)
        box('stair runner',(3.6,y-.01,h+.013),(1.05,.36,.022),green)
        if step%2==0:
            for x in (2.42,4.77):cylinder('stair spindle',(x,y,h+.46),.025,.92,iron,vertices=12)
    for x in (2.42,4.77):beam('stair rail',(x,.2,1.03),(x,5.4,3.63),.07,wood)
    # Two seated tenants face into the room from an oak bench.
    for y in (-2.35,.65):
        for x in (-4.5,-3.85):box('waiting bench leg',(x,y,.35),(.12,.12,.70),wood,.02)
    for x in (-4.52,-4.30,-4.08,-3.86):box('waiting bench slat',(x,-.85,.72),(.19,3.25,.10),wood,.022)
    for z in (1,1.22,1.44):box('waiting bench back',(-4.65,-.85,z),(.12,3.3,.16),green,.02)
    for y in (-2.52,.82):box('waiting bench arm',(-4.23,y,1.06),(.94,.12,.10),wood,.03)
    # Steam radiator, valve and piping under the stair; no floor occupants here.
    for i in range(9):box('radiator fin',(4.28,4.6+i*.11,.58),(.60,.055,.90),iron,.025)
    beam('radiator feed',(4.7,4.6,.18),(4.7,5.7,.18),.055,iron)
    cylinder('radiator valve',(4.7,4.57,.4),.08,.035,brass,vertices=16)
    # Spare folded sheets wait behind reception, clear of the clerk.
    for i in range(3):box('folded boarding linen',(-3.8,4.4,1.2+i*.045),(.42,.52,.04),linen,.015)
    for x in (-4.1,.5):
        box('rear sconce back',(x,5.7,3.22),(.18,.1,.3),brass,.02)
        cylinder('rear opal sconce',(x,5.46,3.27),.13,.30,amber,vertices=24)
    box('entrance coir mat',(0,-2.9,.038),(2.2,1.35,.035),material('Mariner coir',(.20,.14,.07)))
    walls={}
    for side in ('left','back'):
        group=bpy.data.objects.new('interior-wall-'+side,None);bpy.context.collection.objects.link(group);walls[side]=group
    for ob in list(bpy.context.scene.objects):
        if ob.type!='MESH':continue
        if ob.name.startswith('west '):ob.parent=walls['left']
        elif ob.name.startswith('rear '):ob.parent=walls['back']


def mercer_lobby():
    plaster=material('aged cream plaster',(.56,.51,.40))
    wood=material('varnished walnut',(.18,.09,.045))
    stone=material('worn stone steps',(.38,.38,.33))
    brass=material('tenant brass',(.46,.32,.12),.6)
    iron=material('black stair rail',(.035,.045,.04),.5)
    black=material('black terrazzo',(.09,.10,.09))
    cream=material('cream terrazzo',(.55,.51,.40))
    glass=material('frosted glass',(.33,.42,.39),.2)
    box('foundation',(0,1,-.15),(12,12,.25),stone)
    for x in range(-12,12):
        for y in range(-10,14):
            box('terrazzo tile',(x*.5+.25,y*.5+.25,0),(.492,.492,.035),black if (x+y)%2 else cream)
    box('rear plaster wall',(0,7,2.35),(12,.25,4.7),plaster)
    box('west plaster wall',(-6,1,2.35),(.25,12,4.7),plaster)
    for x in range(-23,24):
        box('rear tongue groove board',(x*.25,6.83,.67),(.242,.08,1.34),wood)
    for y in range(-19,28):
        box('west tongue groove board',(-5.83,y*.25,.67),(.08,.242,1.34),wood)
    box('rear dado rail',(0,6.76,1.37),(12,.12,.12),wood,.025)
    box('west dado rail',(-5.76,1,1.37),(.12,12,.12),wood,.025)
    for z in (.12,4.45):
        box('rear moulding',(0,6.78,z),(12,.19,.18),cream,.025)
        box('west moulding',(-5.78,1,z),(.19,12,.18),cream,.025)
    # Forty-eight numbered compartments, matching the accommodation register.
    box('letterbox backboard',(-2,6.68,2.55),(4.75,.16,2.15),wood,.04)
    for row in range(6):
        for col in range(8):
            x=-4.0+col*.57;z=1.7+row*.34
            box('individual letterbox',(x,6.55,z),(.53,.15,.30),brass,.012)
            box('mail slot',(x,6.465,z+.06),(.33,.025,.025),iron)
            cylinder('letterbox keyhole',(x+.16,6.46,z-.07),.022,.025,iron,(math.pi/2,0,0),12)
            # Small engraved number plaques, actual type geometry.
            curve=bpy.data.curves.new('mailbox number','FONT');curve.body=str(row*8+col+1);curve.size=.095;curve.align_x='CENTER';curve.extrude=.001
            obj=bpy.data.objects.new('mailbox number',curve);bpy.context.collection.objects.link(obj);obj.location=(x-.08,6.457,z-.085);obj.rotation_euler=(math.pi/2,0,0);obj.data.materials.append(iron)
            bpy.context.view_layer.objects.active=obj;obj.select_set(True);bpy.ops.object.convert(target='MESH');obj.select_set(False)
    # Apartment corridor door with transom glazing and a worn brass push plate.
    box('corridor door',(1.15,6.76,1.35),(1.8,.16,2.7),wood,.04)
    for x in (.16,2.14):box('door casing',(x,6.61,1.5),(.17,.23,3),cream,.02)
    box('door casing head',(1.15,6.61,2.96),(2.15,.23,.17),cream,.02)
    box('door frosted pane',(1.15,6.655,1.93),(1.38,.03,.85),glass)
    for x in (.75,1.55):box('recessed door panel',(x,6.66,.75),(.58,.035,.8),wood,.035)
    box('door push plate',(1.85,6.64,1.35),(.14,.025,.43),brass,.01)
    # A stone stair with a continuous handrail, open toward the lobby.
    for step in range(15):
        y=.2+step*.42;h=(step+1)*.20
        box('stair tread',(4.25,y,h/2),(2.7,.42,h),stone,.018)
        box('stair nosing',(4.25,y-.205,h),(2.76,.06,.06),cream,.012)
        if step%2==0:
            for x in (2.9,5.6):cylinder('stair baluster',(x,y,h+.47),.032,.94,iron)
    for x in (2.9,5.6):beam('stair handrail',(x,.2,1.1),(x,6.1,3.9),.09,wood)
    # Waiting bench: slatted seat, back, curved-looking arm supports and legs.
    for y in (-2.5,.5):
        for x in (-5.15,-4.48):box('bench leg',(x,y,.35),(.12,.12,.7),wood,.02)
    for x in (-5.18,-4.96,-4.74,-4.52):box('seat slat',(x,-1,.72),(.19,3.3,.10),wood,.025)
    for z in (1.0,1.22,1.44):box('bench back slat',(-5.3,-1,z),(.12,3.4,.16),wood,.025)
    for y in (-2.7,.7):box('bench arm',(-4.88,y,1.06),(.97,.12,.10),wood,.035)
    for x in (-4.7,.1):
        box('sconce back',(x,6.65,3.72),(.22,.15,.46),brass,.03)
        cylinder('sconce opal globe',(x,6.37,3.78),.17,.35,material('opal lamp '+str(x),(.9,.72,.42),0,1.5),vertices=24)
    box('entrance mat',(0,-3.4,.04),(2.8,1.5,.035),material('coir doormat',(.22,.16,.08)))
    # Keep walls and their mounted fixtures independently removable when the
    # browser camera orbits behind them. Furniture remains on the floor.
    walls={}
    for side in ('left','back'):
        group=bpy.data.objects.new('interior-wall-'+side,None);bpy.context.collection.objects.link(group);walls[side]=group
    for ob in list(bpy.context.scene.objects):
        if ob.type!='MESH':continue
        if ob.name.startswith('west '):ob.parent=walls['left']
        elif ob.name.startswith(('rear ','letterbox ','individual letterbox','mail slot','mailbox number','corridor door','door ','recessed door','sconce ')):ob.parent=walls['back']



def mercer_plate(path, camera_at, target, scale):
    scene=bpy.context.scene
    broken=bpy.data.objects.get('window-broken')
    if broken:
        for child in broken.children_recursive:child.hide_render=True
    scene.render.engine='CYCLES';scene.cycles.samples=32;scene.cycles.use_denoising=True
    scene.render.resolution_x=1280;scene.render.resolution_y=900;scene.render.resolution_percentage=100
    scene.world=bpy.data.worlds.new('Mercer ambient sky');scene.world.color=(.22,.22,.20)
    for name,position,power,size in [('soft daylight',(-8,-9,22),2600,12),('warm reflected light',(6,2,15),1600,10)]:
        data=bpy.data.lights.new(name,'AREA');data.energy=power;data.shape='DISK';data.size=size
        obj=bpy.data.objects.new(name,data);scene.collection.objects.link(obj);obj.location=position;obj.rotation_euler=(Vector(target)-obj.location).to_track_quat('-Z','Y').to_euler()
    data=bpy.data.cameras.new('plate camera');obj=bpy.data.objects.new('plate camera',data);scene.collection.objects.link(obj);scene.camera=obj
    obj.location=camera_at;obj.rotation_euler=(Vector(target)-obj.location).to_track_quat('-Z','Y').to_euler();data.type='ORTHO';data.ortho_scale=scale
    scene.render.image_settings.file_format='JPEG';scene.render.image_settings.quality=92;scene.render.filepath=path
    bpy.ops.render.render(write_still=True)


def mercer_court():
    """Five-storey workers' flats: limestone lintels, iron escapes and roof laundry."""
    building('tenement',5,14,12,57,palette=(.40,.22,.15))
    stone=bpy.data.materials['limestone']
    iron=bpy.data.materials['painted iron']
    brass=material('Mercer oxidised brass',(.38,.29,.13),.5)
    linen=material('Mercer washed linen',(.69,.65,.52))
    slate=material('Mercer blue laundry',(.23,.29,.33))
    # A raised central name tablet and paired entrance sconces identify the
    # address without turning the residential facade into a shop frontage.
    box('Mercer entry tablet',(0,6.28,3.05),(3.35,.20,.65),stone,.045)
    anchor=bpy.data.objects['sign-anchor'];anchor.location=(0,6.40,3.05);anchor.scale=(.42,1,.5)
    for x in (-1.6,1.6):
        box('entry lamp bracket',(x,6.36,2.28),(.16,.22,.33),iron,.025)
        box('entry lamp glass',(x,6.49,2.34),(.20,.18,.26),bpy.data.materials['occupied windows'],.02)
    for floor in range(1,5):
        for x in (-5.6,-2.8,0,2.8,5.6):
            box('window lintel cap',(x,6.19,floor*3.15+2.83),(1.63,.29,.15),stone,.018)
    # Cast dentils below the projecting roof cornice.
    for x in range(-13,14):
        box('cornice dentil',(x*.49,6.20,15.42),(.22,.27,.24),stone,.01)
    # Clotheslines live entirely within the roof parapet; no walking route
    # crosses this roof. Deterministic hanging cloth panels include folds.
    for x in (-5.4,-.4):
        for y in (-3,0):
            cylinder('laundry upright',(x,y,17.0),.045,2.15,iron)
    for y in (-3,0):
        beam('laundry line',(-5.4,y,17.7),(-.4,y,17.7),.025,iron)
        for i in range(4):
            x=-4.8+i*1.05
            for fold in range(7):
                box('linen folded panel',(x+fold*.10,y+(.04 if fold%2 else -.04),17.20),(.105,.035,.95),linen if i%2 else slate)
            for pin in (0,.6):
                box('wood clothes peg',(x+pin,y,17.72),(.035,.06,.12),brass)
    for x in (-6.2,6.2):
        cylinder('rainwater downpipe',(x,6.28,7.65),.075,15.2,iron)
        for z in (1,4,7,10,13):
            box('pipe wall collar',(x,6.22,z),(.21,.20,.08),iron)
    # Individual brass letter boxes flank the working entrance door.
    for side in (-1,1):
        for row in range(3):
            for col in range(2):
                x=side*(1.4+col*.22)
                box('tenant letter box',(x,6.28,.98+row*.19),(.20,.08,.16),brass,.009)
                box('letter slot',(x,6.325,1.02+row*.19),(.13,.01,.014),iron)

if __name__ == '__main__' and '--only=slot-cabinet' in __import__('sys').argv:
    clear();slot_cabinet()
    manifest_path=os.path.join(OUT,'manifest.json')
    with open(manifest_path) as f: selected_manifest=json.load(f)
    selected_manifest['slot-cabinet']=export('slot-cabinet')
    with open(manifest_path,'w') as f:json.dump(selected_manifest,f,indent=2)
    raise SystemExit(0)

if __name__ == '__main__' and '--only=interior-laundry' in __import__('sys').argv:
    clear();laundry_interior()
    manifest_path=os.path.join(OUT,'manifest.json')
    with open(manifest_path) as f: selected_manifest=json.load(f)
    selected_manifest['interior-laundry']=export('interior-laundry')
    with open(manifest_path,'w') as f:json.dump(selected_manifest,f,indent=2)
    raise SystemExit(0)

if __name__ == '__main__' and '--only=interior-mariner' in __import__('sys').argv:
    clear();mariner_lobby()
    manifest_path=os.path.join(OUT,'manifest.json')
    with open(manifest_path) as f: selected_manifest=json.load(f)
    selected_manifest['interior-mariner']=export('interior-mariner')
    with open(manifest_path,'w') as f:json.dump(selected_manifest,f,indent=2)
    raise SystemExit(0)

if __name__ == '__main__' and '--only=bar-cloth' in __import__('sys').argv:
    clear();bar_cloth()
    manifest_path=os.path.join(OUT,'manifest.json')
    with open(manifest_path) as f: selected_manifest=json.load(f)
    selected_manifest['bar-cloth']=export('bar-cloth')
    with open(manifest_path,'w') as f:json.dump(selected_manifest,f,indent=2)
    raise SystemExit(0)

# Export just this new address without rewriting reviewed assets.
if __name__ == '__main__' and '--only=mercer-court' in __import__('sys').argv:
    clear();mercer_court()
    manifest_path=os.path.join(OUT,'manifest.json')
    with open(manifest_path) as f: selected_manifest=json.load(f)
    selected_manifest['mercer-court']=export('mercer-court')
    mercer_plate(os.path.abspath('public/art/fronts/front-mercercourt-v1.jpg'),(28,35,27),(0,0,9),36)
    clear();mercer_lobby();selected_manifest['interior-mercer-court']=export('interior-mercer-court')
    mercer_plate(os.path.abspath('public/art/rooms/room-mercercourt-v1.jpg'),(11,-15,10),(0,2,1.6),15)
    with open(manifest_path,'w') as f:json.dump(selected_manifest,f,indent=2)
    raise SystemExit(0)


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
clear();car('ford',police=True)
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
clear();fire_nozzle();manifest['fire-nozzle']=export('fire-nozzle')
clear();firefighter();manifest['firefighter']=export('firefighter')
clear();fire_engine();manifest['fire-engine']=export('fire-engine')
clear();police_officer();manifest['police-officer']=export('police-officer')
clear();undertaker();manifest['undertaker']=export('undertaker')
clear();handcuffs();manifest['handcuffs']=export('handcuffs')
clear();revolver();manifest['revolver']=export('revolver')
for name in ('shotgun','thompson'):
    clear();long_gun(name);manifest[name]=export(name)
clear();blast_fragment();manifest['blast-fragment']=export('blast-fragment')
clear();bar_cloth();manifest['bar-cloth']=export('bar-cloth')
clear();mariner();manifest['mariner']=export('mariner')
clear();slot_cabinet();manifest['slot-cabinet']=export('slot-cabinet')
clear();laundry_interior();manifest['interior-laundry']=export('interior-laundry')
clear();mariner_lobby();manifest['interior-mariner']=export('interior-mariner')
clear();mercer_court();manifest['mercer-court']=export('mercer-court')
clear();mercer_lobby();manifest['interior-mercer-court']=export('interior-mercer-court')
clear();saint_agnes_interior();manifest['interior-saint-agnes']=export('interior-saint-agnes')
clear();harbour_pier();manifest['harbour-pier']=export('harbour-pier')
clear();quay_section();manifest['quay-section']=export('quay-section')
clear();vacant_lot();manifest['vacant-lot']=export('vacant-lot')
clear();street_bed();manifest['street-bed']=export('street-bed')
clear();streetside();manifest['streetside']=export('streetside')
with open(os.path.join(OUT,'manifest.json'),'w') as f: json.dump(manifest,f,indent=2)
print('Exported',len(manifest),'models')
