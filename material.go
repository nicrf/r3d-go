package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// Material modes -------------------------------------------------------------

// BillboardMode controls how a material orients toward the camera.
type BillboardMode int32

const (
	BillboardDisabled BillboardMode = C.R3D_BILLBOARD_DISABLED
	BillboardFront    BillboardMode = C.R3D_BILLBOARD_FRONT
	BillboardYAxis    BillboardMode = C.R3D_BILLBOARD_Y_AXIS
)

// BlendMode controls color blending for a material.
type BlendMode int32

const (
	BlendMix                BlendMode = C.R3D_BLEND_MIX
	BlendAdditive           BlendMode = C.R3D_BLEND_ADDITIVE
	BlendMultiply           BlendMode = C.R3D_BLEND_MULTIPLY
	BlendPremultipliedAlpha BlendMode = C.R3D_BLEND_PREMULTIPLIED_ALPHA
)

// TransparencyMode controls how transparent surfaces are rendered.
type TransparencyMode int32

const (
	TransparencyDisabled TransparencyMode = C.R3D_TRANSPARENCY_DISABLED
	TransparencyPrepass  TransparencyMode = C.R3D_TRANSPARENCY_PREPASS
	TransparencyAlpha    TransparencyMode = C.R3D_TRANSPARENCY_ALPHA
)

// CullMode controls face culling for a material.
type CullMode int32

const (
	CullNone  CullMode = C.R3D_CULL_NONE
	CullBack  CullMode = C.R3D_CULL_BACK
	CullFront CullMode = C.R3D_CULL_FRONT
)

// Texture maps ---------------------------------------------------------------

// AlbedoMap wraps an R3D_AlbedoMap (base color texture + tint).
type AlbedoMap struct{ c C.R3D_AlbedoMap }

// NormalMap wraps an R3D_NormalMap.
type NormalMap struct{ c C.R3D_NormalMap }

// EmissionMap wraps an R3D_EmissionMap.
type EmissionMap struct{ c C.R3D_EmissionMap }

// LoadAlbedoMap loads a base-color texture, multiplied by color.
func LoadAlbedoMap(fileName string, color Color) AlbedoMap {
	c := C.CString(fileName)
	defer C.free(unsafe.Pointer(c))
	return AlbedoMap{c: C.R3D_LoadAlbedoMap(c, color.c())}
}

// LoadNormalMap loads a tangent-space normal map with the given strength.
func LoadNormalMap(fileName string, scale float32) NormalMap {
	c := C.CString(fileName)
	defer C.free(unsafe.Pointer(c))
	return NormalMap{c: C.R3D_LoadNormalMap(c, C.float(scale))}
}

// LoadEmissionMap loads an emission texture tinted by color at the given energy.
func LoadEmissionMap(fileName string, color Color, energy float32) EmissionMap {
	c := C.CString(fileName)
	defer C.free(unsafe.Pointer(c))
	return EmissionMap{c: C.R3D_LoadEmissionMap(c, color.c(), C.float(energy))}
}

// Material setters (extends the type declared in r3d.go) ---------------------

// SetAlbedoMap assigns a loaded albedo map.
func (m *Material) SetAlbedoMap(a AlbedoMap) { m.c.albedo = a.c }

// SetNormalMap assigns a loaded normal map.
func (m *Material) SetNormalMap(n NormalMap) { m.c.normal = n.c }

// SetEmissionMap assigns a loaded emission map.
func (m *Material) SetEmissionMap(e EmissionMap) { m.c.emission = e.c }

// SetEmission sets the emission color and energy.
func (m *Material) SetEmission(color Color, energy float32) {
	m.c.emission.color = color.c()
	m.c.emission.energy = C.float(energy)
}

// SetMetalness sets the PBR metalness factor [0,1].
func (m *Material) SetMetalness(v float32) { m.c.orm.metalness = C.float(v) }

// SetRoughness sets the PBR roughness factor [0,1].
func (m *Material) SetRoughness(v float32) { m.c.orm.roughness = C.float(v) }

// SetOcclusion sets the ambient-occlusion factor [0,1].
func (m *Material) SetOcclusion(v float32) { m.c.orm.occlusion = C.float(v) }

// SetSpecular sets the specular factor.
func (m *Material) SetSpecular(v float32) { m.c.orm.specular = C.float(v) }

// SetNormalScale sets the normal-map strength.
func (m *Material) SetNormalScale(v float32) { m.c.normal.scale = C.float(v) }

// SetUVScale sets the UV tiling factor.
func (m *Material) SetUVScale(x, y float32) {
	m.c.uvScale = C.Vector2{x: C.float(x), y: C.float(y)}
}

// SetUVOffset sets the UV offset.
func (m *Material) SetUVOffset(x, y float32) {
	m.c.uvOffset = C.Vector2{x: C.float(x), y: C.float(y)}
}

// SetBillboardMode sets the material billboard mode.
func (m *Material) SetBillboardMode(mode BillboardMode) {
	m.c.billboardMode = C.R3D_BillboardMode(mode)
}

// SetBlendMode sets the material blend mode.
func (m *Material) SetBlendMode(mode BlendMode) { m.c.blendMode = C.R3D_BlendMode(mode) }

// SetTransparencyMode sets the material transparency mode.
func (m *Material) SetTransparencyMode(mode TransparencyMode) {
	m.c.transparencyMode = C.R3D_TransparencyMode(mode)
}

// SetCullMode sets the material face-culling mode.
func (m *Material) SetCullMode(mode CullMode) { m.c.cullMode = C.R3D_CullMode(mode) }

// SetUnlit toggles whether the material ignores lighting.
func (m *Material) SetUnlit(unlit bool) { m.c.unlit = C.bool(unlit) }

// CompareMode is a depth/stencil comparison function.
type CompareMode int32

const (
	CompareLess     CompareMode = C.R3D_COMPARE_LESS
	CompareLEqual   CompareMode = C.R3D_COMPARE_LEQUAL
	CompareEqual    CompareMode = C.R3D_COMPARE_EQUAL
	CompareGreater  CompareMode = C.R3D_COMPARE_GREATER
	CompareGEqual   CompareMode = C.R3D_COMPARE_GEQUAL
	CompareNotEqual CompareMode = C.R3D_COMPARE_NOTEQUAL
	CompareAlways   CompareMode = C.R3D_COMPARE_ALWAYS
	CompareNever    CompareMode = C.R3D_COMPARE_NEVER
)

// StencilOp is a stencil-buffer update operation.
type StencilOp int32

const (
	StencilKeep    StencilOp = C.R3D_STENCIL_KEEP
	StencilZero    StencilOp = C.R3D_STENCIL_ZERO
	StencilReplace StencilOp = C.R3D_STENCIL_REPLACE
	StencilIncr    StencilOp = C.R3D_STENCIL_INCR
	StencilDecr    StencilOp = C.R3D_STENCIL_DECR
)

// SetDepthMode sets the depth comparison function (e.g. CompareAlways to
// ignore depth for x-ray passes).
func (m *Material) SetDepthMode(mode CompareMode) {
	m.c.depth.mode = C.R3D_CompareMode(mode)
}

// SetStencil configures the full stencil state for this material
// (used for outline / x-ray / masking effects).
func (m *Material) SetStencil(mode CompareMode, ref, mask uint8, opFail, opZFail, opPass StencilOp) {
	m.c.stencil.mode = C.R3D_CompareMode(mode)
	m.c.stencil.ref = C.uint8_t(ref)
	m.c.stencil.mask = C.uint8_t(mask)
	m.c.stencil.opFail = C.R3D_StencilOp(opFail)
	m.c.stencil.opZFail = C.R3D_StencilOp(opZFail)
	m.c.stencil.opPass = C.R3D_StencilOp(opPass)
}
