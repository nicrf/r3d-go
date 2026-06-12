// Package r3d provides a cgo binding for the R3D 3D-rendering library
// (https://github.com/Bigfoot71/r3d), an extension of raylib.
//
// The native static libraries (libr3d.a, libraylib.a, libassimp.a) and the
// public headers are vendored under lib/ and include/ so the package is
// self-contained: `go build` links them directly. The C include paths are
// common to all platforms; the per-OS library search paths and link flags
// live in the cgo_<goos>.go files. Supported targets: macOS/arm64 and
// Windows/amd64, each with its own vendored archives under lib/<goos_goarch>/
// (see cgo_darwin.go / cgo_windows.go).
package r3d

/*
#cgo CFLAGS: -I${SRCDIR}/include -I${SRCDIR}/include/r3d

#include <stdlib.h>
#include "raylib.h"
#include "raymath.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// ---------------------------------------------------------------------------
// Value types (memory-compatible with the matching raylib structs)
// ---------------------------------------------------------------------------

// Vector2 maps to raylib's Vector2.
type Vector2 struct{ X, Y float32 }

// Vector3 maps to raylib's Vector3.
type Vector3 struct{ X, Y, Z float32 }

// Vector4 maps to raylib's Vector4 / Quaternion.
type Vector4 struct{ X, Y, Z, W float32 }

// Quaternion is an alias of Vector4, matching raylib semantics.
type Quaternion = Vector4

// Color maps to raylib's Color (RGBA, one byte per channel).
type Color struct{ R, G, B, A uint8 }

// QuaternionIdentity returns the identity quaternion (0,0,0,1).
func QuaternionIdentity() Quaternion { return Quaternion{0, 0, 0, 1} }

func (v Vector2) c() C.Vector2 { return C.Vector2{x: C.float(v.X), y: C.float(v.Y)} }
func (v Vector3) c() C.Vector3 {
	return C.Vector3{x: C.float(v.X), y: C.float(v.Y), z: C.float(v.Z)}
}
func (v Vector4) c() C.Vector4 {
	return C.Vector4{x: C.float(v.X), y: C.float(v.Y), z: C.float(v.Z), w: C.float(v.W)}
}
func (c Color) c() C.Color {
	return C.Color{r: C.uchar(c.R), g: C.uchar(c.G), b: C.uchar(c.B), a: C.uchar(c.A)}
}

func goVec3(v C.Vector3) Vector3 { return Vector3{float32(v.x), float32(v.y), float32(v.z)} }

// ---------------------------------------------------------------------------
// Camera
// ---------------------------------------------------------------------------

// Camera projection modes (raylib CameraProjection).
const (
	CameraPerspective  = int(C.CAMERA_PERSPECTIVE)
	CameraOrthographic = int(C.CAMERA_ORTHOGRAPHIC)
)

// Camera movement modes for UpdateCamera (raylib CameraMode).
const (
	CameraCustom   = int(C.CAMERA_CUSTOM)
	CameraFree     = int(C.CAMERA_FREE)
	CameraOrbital  = int(C.CAMERA_ORBITAL)
	CameraFirstPer = int(C.CAMERA_FIRST_PERSON)
	CameraThirdPer = int(C.CAMERA_THIRD_PERSON)
)

// Camera3D maps to raylib's Camera3D.
type Camera3D struct {
	Position   Vector3
	Target     Vector3
	Up         Vector3
	Fovy       float32
	Projection int32
}

func (cam Camera3D) c() C.Camera3D {
	return C.Camera3D{
		position:   cam.Position.c(),
		target:     cam.Target.c(),
		up:         cam.Up.c(),
		fovy:       C.float(cam.Fovy),
		projection: C.int(cam.Projection),
	}
}

func goCamera(c C.Camera3D) Camera3D {
	return Camera3D{
		Position:   goVec3(c.position),
		Target:     goVec3(c.target),
		Up:         goVec3(c.up),
		Fovy:       float32(c.fovy),
		Projection: int32(c.projection),
	}
}

// ---------------------------------------------------------------------------
// R3D core
// ---------------------------------------------------------------------------

// Init initializes the R3D renderer at the given internal resolution.
// Call after a window exists (see InitWindow). Returns false on failure.
func Init(resWidth, resHeight int) bool {
	return bool(C.R3D_Init(C.int(resWidth), C.int(resHeight)))
}

// Close releases all R3D resources.
func Close() { C.R3D_Close() }

// Begin starts an R3D render pass with the given camera.
func Begin(cam Camera3D) { C.R3D_Begin(cam.c()) }

// End finishes the current R3D render pass and composites the frame.
func End() { C.R3D_End() }

// ---------------------------------------------------------------------------
// Lighting
// ---------------------------------------------------------------------------

// LightType enumerates the supported light kinds.
type LightType int32

const (
	LightDir  LightType = C.R3D_LIGHT_DIR
	LightSpot LightType = C.R3D_LIGHT_SPOT
	LightOmni LightType = C.R3D_LIGHT_OMNI
)

// Light is a handle to a renderer-managed light.
type Light int32

// CreateLight creates a new light of the given type.
func CreateLight(t LightType) Light { return Light(C.R3D_CreateLight(C.R3D_LightType(t))) }

// Destroy removes the light.
func (l Light) Destroy() { C.R3D_DestroyLight(C.R3D_Light(l)) }

// SetActive enables or disables the light.
func (l Light) SetActive(active bool) {
	C.R3D_SetLightActive(C.R3D_Light(l), C.bool(active))
}

// SetDirection sets the light direction (directional / spot lights).
func (l Light) SetDirection(dir Vector3) {
	C.R3D_SetLightDirection(C.R3D_Light(l), dir.c())
}

// SetPosition sets the light position (spot / omni lights).
func (l Light) SetPosition(pos Vector3) {
	C.R3D_SetLightPosition(C.R3D_Light(l), pos.c())
}

// SetColor sets the light color.
func (l Light) SetColor(col Color) {
	C.R3D_SetLightColor(C.R3D_Light(l), col.c())
}

// SetEnergy sets the light energy (intensity multiplier).
func (l Light) SetEnergy(energy float32) {
	C.R3D_SetLightEnergy(C.R3D_Light(l), C.float(energy))
}

// EnableShadow turns on shadow casting for this light.
func (l Light) EnableShadow() {
	C.R3D_EnableShadow(C.R3D_Light(l))
}

// ---------------------------------------------------------------------------
// Meshes
// ---------------------------------------------------------------------------

// Mesh wraps an R3D_Mesh.
type Mesh struct{ c C.R3D_Mesh }

// GenMeshSphere generates a UV sphere mesh.
func GenMeshSphere(radius float32, rings, slices int) Mesh {
	return Mesh{c: C.R3D_GenMeshSphere(C.float(radius), C.int(rings), C.int(slices))}
}

// GenMeshCube generates a box mesh.
func GenMeshCube(width, height, length float32) Mesh {
	return Mesh{c: C.R3D_GenMeshCube(C.float(width), C.float(height), C.float(length))}
}

// GenMeshPlane generates a subdivided plane mesh.
func GenMeshPlane(width, length float32, resX, resZ int) Mesh {
	return Mesh{c: C.R3D_GenMeshPlane(C.float(width), C.float(length), C.int(resX), C.int(resZ))}
}

// GenMeshCylinder generates a cylinder mesh.
func GenMeshCylinder(radius, height float32, slices int) Mesh {
	return Mesh{c: C.R3D_GenMeshCylinder(C.float(radius), C.float(height), C.int(slices))}
}

// GenMeshTorus generates a torus mesh.
func GenMeshTorus(radius, size float32, radSeg, sides int) Mesh {
	return Mesh{c: C.R3D_GenMeshTorus(C.float(radius), C.float(size), C.int(radSeg), C.int(sides))}
}

// IsValid reports whether the mesh holds valid GPU data.
func (m Mesh) IsValid() bool { return bool(C.R3D_IsMeshValid(m.c)) }

// Unload frees the mesh GPU resources.
func (m Mesh) Unload() { C.R3D_UnloadMesh(m.c) }

// ---------------------------------------------------------------------------
// Materials
// ---------------------------------------------------------------------------

// Material wraps an R3D_Material.
type Material struct{ c C.R3D_Material }

// GetDefaultMaterial returns the renderer's default PBR material.
func GetDefaultMaterial() Material { return Material{c: C.R3D_GetDefaultMaterial()} }

// SetAlbedoColor sets the material albedo (base) color multiplier.
func (m *Material) SetAlbedoColor(col Color) { m.c.albedo.color = col.c() }

// SetEmissionColor sets the material emission color.
func (m *Material) SetEmissionColor(col Color) { m.c.emission.color = col.c() }

// Unload frees material GPU resources.
func (m Material) Unload() { C.R3D_UnloadMaterial(m.c) }

// ---------------------------------------------------------------------------
// Drawing
// ---------------------------------------------------------------------------

// DrawMesh draws a single mesh at a position with uniform scale.
func DrawMesh(m Mesh, mat Material, pos Vector3, scale float32) {
	C.R3D_DrawMesh(m.c, mat.c, pos.c(), C.float(scale))
}

// DrawMeshEx draws a single mesh with full position/rotation/scale.
func DrawMeshEx(m Mesh, mat Material, pos Vector3, rot Quaternion, scale Vector3) {
	C.R3D_DrawMeshEx(m.c, mat.c, pos.c(), rot.c(), scale.c())
}

// DrawMeshInstanced draws `count` instances of a mesh from an instance buffer.
func DrawMeshInstanced(m Mesh, mat Material, buf InstanceBuffer, count int) {
	C.R3D_DrawMeshInstanced(m.c, mat.c, buf.c, C.int(count))
}

// ---------------------------------------------------------------------------
// Instance buffers
// ---------------------------------------------------------------------------

// InstanceFlags is a bitmask of per-instance attributes.
type InstanceFlags uint32

const (
	InstancePosition InstanceFlags = 1 << 0
	InstanceRotation InstanceFlags = 1 << 1
	InstanceScale    InstanceFlags = 1 << 2
	InstanceColor    InstanceFlags = 1 << 3
	InstanceCustom   InstanceFlags = 1 << 4
)

// InstanceBuffer wraps an R3D_InstanceBuffer plus its capacity for safe mapping.
type InstanceBuffer struct {
	c        C.R3D_InstanceBuffer
	capacity int
}

// LoadInstanceBuffer allocates an instance buffer for the given attribute set.
// With the default formats, mapped attributes have these native layouts:
// position/scale -> Vector3, rotation -> Quaternion, color -> Color.
func LoadInstanceBuffer(capacity int, flags InstanceFlags) InstanceBuffer {
	return InstanceBuffer{
		c:        C.R3D_LoadInstanceBuffer(C.int(capacity), C.R3D_InstanceFlags(flags)),
		capacity: capacity,
	}
}

// Capacity returns the buffer's instance capacity.
func (b InstanceBuffer) Capacity() int { return b.capacity }

func (b InstanceBuffer) mapPtr(flag InstanceFlags, discard bool) unsafe.Pointer {
	return unsafe.Pointer(C.R3D_MapInstances(b.c, C.R3D_InstanceFlags(flag), C.bool(discard)))
}

// MapPositions returns the position attribute as a writable Go slice backed by
// the mapped GPU buffer. Call UnmapInstances when done.
func (b InstanceBuffer) MapPositions(discard bool) []Vector3 {
	return unsafe.Slice((*Vector3)(b.mapPtr(InstancePosition, discard)), b.capacity)
}

// MapRotations returns the rotation attribute as a writable Go slice.
func (b InstanceBuffer) MapRotations(discard bool) []Quaternion {
	return unsafe.Slice((*Quaternion)(b.mapPtr(InstanceRotation, discard)), b.capacity)
}

// MapScales returns the scale attribute as a writable Go slice.
func (b InstanceBuffer) MapScales(discard bool) []Vector3 {
	return unsafe.Slice((*Vector3)(b.mapPtr(InstanceScale, discard)), b.capacity)
}

// MapColors returns the color attribute as a writable Go slice.
func (b InstanceBuffer) MapColors(discard bool) []Color {
	return unsafe.Slice((*Color)(b.mapPtr(InstanceColor, discard)), b.capacity)
}

// Unmap flushes the given attributes back to the GPU after mapping.
func (b InstanceBuffer) Unmap(flags InstanceFlags) {
	C.R3D_UnmapInstances(b.c, C.R3D_InstanceFlags(flags))
}

// Unload frees the instance buffer.
func (b InstanceBuffer) Unload() { C.R3D_UnloadInstanceBuffer(b.c) }

// ---------------------------------------------------------------------------
// Environment
// ---------------------------------------------------------------------------

// SetBackgroundColor sets the clear/background color used when no skybox is set.
func SetBackgroundColor(col Color) {
	C.R3D_GetEnvironment().background.color = col.c()
}

// SetAmbientColor sets the ambient light color used when no IBL map is set.
func SetAmbientColor(col Color) {
	C.R3D_GetEnvironment().ambient.color = col.c()
}
