package r3d

/*
#include "raylib.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// ---------------------------------------------------------------------------
// Extra light controls (extends Light declared in r3d.go)
// ---------------------------------------------------------------------------

// SetColor sets the light color. (declared here to keep r3d.go minimal)
// SetRange sets the light effective range (omni/spot falloff distance).
func (l Light) SetRange(r float32) { C.R3D_SetLightRange(C.R3D_Light(l), C.float(r)) }

// SetSpecular sets the light specular contribution.
func (l Light) SetSpecular(s float32) { C.R3D_SetLightSpecular(C.R3D_Light(l), C.float(s)) }

// LookAt orients the light from position toward target (spot/dir lights).
func (l Light) LookAt(position, target Vector3) {
	C.R3D_LightLookAt(C.R3D_Light(l), position.c(), target.c())
}

// SetInnerCutOff sets the spot light inner cone angle (degrees).
func (l Light) SetInnerCutOff(deg float32) {
	C.R3D_SetLightInnerCutOff(C.R3D_Light(l), C.float(deg))
}

// SetOuterCutOff sets the spot light outer cone angle (degrees).
func (l Light) SetOuterCutOff(deg float32) {
	C.R3D_SetLightOuterCutOff(C.R3D_Light(l), C.float(deg))
}

// SetShadowSoftness sets the soft-shadow blur amount.
func (l Light) SetShadowSoftness(s float32) {
	C.R3D_SetShadowSoftness(C.R3D_Light(l), C.float(s))
}

// SetShadowBias sets the shadow depth bias (reduces acne).
func (l Light) SetShadowBias(v float32) {
	C.R3D_SetShadowDepthBias(C.R3D_Light(l), C.float(v))
}

// ---------------------------------------------------------------------------
// Instance buffer streaming upload (for particle systems, etc.)
// ---------------------------------------------------------------------------

// UploadPositions streams `count` positions into the buffer's POSITION attribute.
func (b InstanceBuffer) UploadPositions(positions []Vector3, count int) {
	if count == 0 {
		return
	}
	C.R3D_UploadInstances(b.c, C.R3D_InstanceFlags(InstancePosition),
		0, C.int(count), unsafe.Pointer(&positions[0]), C.bool(true))
}

// ---------------------------------------------------------------------------
// Texture filtering
// ---------------------------------------------------------------------------

// Texture filter modes (raylib TextureFilter).
type TextureFilter int32

const (
	FilterPoint         TextureFilter = C.TEXTURE_FILTER_POINT
	FilterBilinear      TextureFilter = C.TEXTURE_FILTER_BILINEAR
	FilterTrilinear     TextureFilter = C.TEXTURE_FILTER_TRILINEAR
	FilterAnisotropic4x TextureFilter = C.TEXTURE_FILTER_ANISOTROPIC_4X
	FilterAnisotropic8x TextureFilter = C.TEXTURE_FILTER_ANISOTROPIC_8X
)

// SetTextureFilter sets the default texture filtering used by R3D loaders.
func SetTextureFilter(f TextureFilter) { C.R3D_SetTextureFilter(C.TextureFilter(f)) }
