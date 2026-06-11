package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// GenMeshQuad generates a flat quad facing frontDir (handy for billboards/sprites).
func GenMeshQuad(width, length float32, resX, resZ int, frontDir Vector3) Mesh {
	return Mesh{c: C.R3D_GenMeshQuad(C.float(width), C.float(length), C.int(resX), C.int(resZ), frontDir.c())}
}

// Model wraps an R3D_Model (meshes + materials + optional skeleton).
type Model struct{ c C.R3D_Model }

// LoadModel loads a model (glTF/GLB, OBJ, FBX, … via Assimp).
func LoadModel(filePath string) Model {
	c := C.CString(filePath)
	defer C.free(unsafe.Pointer(c))
	return Model{c: C.R3D_LoadModel(c)}
}

// MeshCount returns how many meshes compose the model.
func (m Model) MeshCount() int { return int(m.c.meshCount) }

// AABB returns the model's axis-aligned bounding box (min, max corners).
func (m Model) AABB() (min, max Vector3) {
	return goVec3(m.c.aabb.min), goVec3(m.c.aabb.max)
}

// Unload frees the model. If unloadMaterials is true, its materials/textures are freed too.
func (m Model) Unload(unloadMaterials bool) {
	C.R3D_UnloadModel(m.c, C.bool(unloadMaterials))
}

// DrawModel draws the model at a position with uniform scale.
func DrawModel(m Model, pos Vector3, scale float32) {
	C.R3D_DrawModel(m.c, pos.c(), C.float(scale))
}

// DrawModelEx draws the model with full position/rotation/scale.
func DrawModelEx(m Model, pos Vector3, rot Quaternion, scale Vector3) {
	C.R3D_DrawModelEx(m.c, pos.c(), rot.c(), scale.c())
}

// ---------------------------------------------------------------------------
// Skeletal animation
// ---------------------------------------------------------------------------

// AnimationLib wraps an R3D_AnimationLib (a set of named animation clips).
type AnimationLib struct{ c C.R3D_AnimationLib }

// LoadAnimationLib loads all animation clips from a model file.
func LoadAnimationLib(filePath string) AnimationLib {
	c := C.CString(filePath)
	defer C.free(unsafe.Pointer(c))
	return AnimationLib{c: C.R3D_LoadAnimationLib(c)}
}

// Count returns the number of animation clips in the library.
func (a AnimationLib) Count() int { return int(a.c.count) }

// Unload frees the animation library.
func (a AnimationLib) Unload() { C.R3D_UnloadAnimationLib(a.c) }

// AnimationPlayer wraps an R3D_AnimationPlayer driving a model's skeleton.
type AnimationPlayer struct{ c C.R3D_AnimationPlayer }

// LoadAnimationPlayer binds a model's skeleton to an animation library.
func LoadAnimationPlayer(m Model, lib AnimationLib) AnimationPlayer {
	return AnimationPlayer{c: C.R3D_LoadAnimationPlayer(m.c.skeleton, lib.c)}
}

// Play starts the clip at the given index.
func (p *AnimationPlayer) Play(animIndex int) {
	C.R3D_PlayAnimation(&p.c, C.int(animIndex))
}

// SetLoop sets whether the clip at animIndex loops.
func (p *AnimationPlayer) SetLoop(animIndex int, loop bool) {
	C.R3D_SetAnimationLoop(&p.c, C.int(animIndex), C.bool(loop))
}

// SetSpeed sets the playback speed of the clip at animIndex.
func (p *AnimationPlayer) SetSpeed(animIndex int, speed float32) {
	C.R3D_SetAnimationSpeed(&p.c, C.int(animIndex), C.float(speed))
}

// Update advances the animation by dt seconds and recomputes the pose.
func (p *AnimationPlayer) Update(dt float32) {
	C.R3D_UpdateAnimationPlayer(&p.c, C.float(dt))
}

// Unload frees the animation player.
func (p AnimationPlayer) Unload() { C.R3D_UnloadAnimationPlayer(p.c) }

// DrawAnimatedModel draws a model posed by the player at a position and uniform scale.
func DrawAnimatedModel(m Model, p AnimationPlayer, pos Vector3, scale float32) {
	C.R3D_DrawAnimatedModel(m.c, p.c, pos.c(), C.float(scale))
}

// DrawAnimatedModelEx draws a posed model with full position/rotation/scale.
func DrawAnimatedModelEx(m Model, p AnimationPlayer, pos Vector3, rot Quaternion, scale Vector3) {
	C.R3D_DrawAnimatedModelEx(m.c, p.c, pos.c(), rot.c(), scale.c())
}
