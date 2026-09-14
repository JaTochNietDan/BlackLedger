"""Render fallback plates from our authored GLBs; no external image assets."""
import bpy
import os
from mathutils import Vector


def render(model, destination, interior=False):
    bpy.ops.object.select_all(action='SELECT');bpy.ops.object.delete(use_global=False)
    bpy.ops.import_scene.gltf(filepath=os.path.abspath('public/art/models/'+model+'.glb'))
    scene=bpy.context.scene
    scene.render.engine='CYCLES';scene.cycles.samples=24
    scene.render.resolution_x=960;scene.render.resolution_y=720;scene.render.resolution_percentage=100
    scene.world.color=(.25,.25,.25)
    # Model is converted back from glTF Y-up to Blender Z-up by the importer.
    focus=Vector((0,0,1.0 if interior else 2.2))
    bpy.ops.object.camera_add(location=(15,-19,16) if interior else (23,28,21))
    camera=bpy.context.object;camera.rotation_euler=(focus-camera.location).to_track_quat('-Z','Y').to_euler()
    camera.data.type='ORTHO';camera.data.ortho_scale=17 if interior else 25;scene.camera=camera
    for position,power,size in [((2,-8,16),1800,10),((-8,5,10),1200,8),((8,8,15),1700,10)]:
        bpy.ops.object.light_add(type='AREA',location=position)
        light=bpy.context.object;light.data.energy=power;light.data.shape='DISK';light.data.size=size
        light.rotation_euler=(focus-light.location).to_track_quat('-Z','Y').to_euler()
    scene.view_settings.view_transform='AgX'
    scene.render.image_settings.file_format='JPEG';scene.render.image_settings.quality=92
    scene.render.filepath=os.path.abspath(destination)
    bpy.ops.render.render(write_still=True)

for kind in ('mortuary','cemetery','crematorium'):
    render(kind,'public/art/fronts/front-'+kind+'-v1.jpg')
for kind in ('mortuary','cemetery','crematorium'):
    render('interior-'+kind,'public/art/rooms/room-'+kind+'-v1.jpg',True)
