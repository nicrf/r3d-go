package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"

// `map` is a Go keyword, so cgo can't reach env->ambient.map directly.
static void r3d_set_ambient_map(R3D_AmbientMap m) { R3D_GetEnvironment()->ambient.map = m; }

// R3D_PROCEDURAL_SKY_BASE is a compound-literal macro; cgo can't use it directly.
static R3D_ProceduralSky r3d_procedural_sky_base(void) { return R3D_PROCEDURAL_SKY_BASE; }
*/
import "C"

import "unsafe"

// ---------------------------------------------------------------------------
// Skybox, cubemaps and image-based lighting
// ---------------------------------------------------------------------------

// CubemapLayout describes how a source image maps onto a cubemap.
type CubemapLayout int32

const (
	CubemapAutoDetect     CubemapLayout = C.R3D_CUBEMAP_LAYOUT_AUTO_DETECT
	CubemapLineVertical   CubemapLayout = C.R3D_CUBEMAP_LAYOUT_LINE_VERTICAL
	CubemapLineHorizontal CubemapLayout = C.R3D_CUBEMAP_LAYOUT_LINE_HORIZONTAL
	CubemapCross3x4       CubemapLayout = C.R3D_CUBEMAP_LAYOUT_CROSS_THREE_BY_FOUR
	CubemapCross4x3       CubemapLayout = C.R3D_CUBEMAP_LAYOUT_CROSS_FOUR_BY_THREE
	CubemapPanorama       CubemapLayout = C.R3D_CUBEMAP_LAYOUT_PANORAMA
)

// AmbientFlags selects which IBL products to generate from a cubemap.
type AmbientFlags uint32

const (
	AmbientIllumination AmbientFlags = 1 << 0 // diffuse irradiance
	AmbientReflection   AmbientFlags = 1 << 1 // specular prefiltered map
)

// Cubemap wraps an R3D_Cubemap (used as a skybox source).
type Cubemap struct{ c C.R3D_Cubemap }

// AmbientMap wraps an R3D_AmbientMap (irradiance + prefiltered reflection).
type AmbientMap struct{ c C.R3D_AmbientMap }

// LoadCubemap loads a cubemap from an image/HDR panorama on disk.
func LoadCubemap(fileName string, layout CubemapLayout) Cubemap {
	c := C.CString(fileName)
	defer C.free(unsafe.Pointer(c))
	return Cubemap{c: C.R3D_LoadCubemap(c, C.R3D_CubemapLayout(layout))}
}

// Unload frees the cubemap.
func (c Cubemap) Unload() { C.R3D_UnloadCubemap(c.c) }

// GenAmbientMap generates IBL maps (irradiance/reflection) from a cubemap.
func GenAmbientMap(cube Cubemap, flags AmbientFlags) AmbientMap {
	return AmbientMap{c: C.R3D_GenAmbientMap(cube.c, C.R3D_AmbientFlags(flags))}
}

// Unload frees the ambient map.
func (a AmbientMap) Unload() { C.R3D_UnloadAmbientMap(a.c) }

// SetSky sets the skybox cubemap drawn as the scene background.
func SetSky(cube Cubemap) { C.R3D_GetEnvironment().background.sky = cube.c }

// SetSkyBlur sets the skybox blur factor [0,1].
func SetSkyBlur(blur float32) { C.R3D_GetEnvironment().background.skyBlur = C.float(blur) }

// SetBackgroundRotation rotates the skybox/background by a quaternion.
func SetBackgroundRotation(rot Quaternion) {
	C.R3D_GetEnvironment().background.rotation = rot.c()
}

// SetBackgroundEnergy scales the background (skybox or color) intensity.
func SetBackgroundEnergy(energy float32) {
	C.R3D_GetEnvironment().background.energy = C.float(energy)
}

// SetAmbientMap sets the IBL ambient map used for diffuse + reflections.
func SetAmbientMap(a AmbientMap) { C.r3d_set_ambient_map(a.c) }

// SetAmbientEnergy scales ambient light intensity.
func SetAmbientEnergy(energy float32) { C.R3D_GetEnvironment().ambient.energy = C.float(energy) }

// ---------------------------------------------------------------------------
// Post-processing
// ---------------------------------------------------------------------------

// BloomMode selects the bloom blend operator.
type BloomMode int32

const (
	BloomDisabled BloomMode = C.R3D_BLOOM_DISABLED
	BloomMix      BloomMode = C.R3D_BLOOM_MIX
	BloomAdditive BloomMode = C.R3D_BLOOM_ADDITIVE
	BloomScreen   BloomMode = C.R3D_BLOOM_SCREEN
)

// FogMode selects the fog distribution.
type FogMode int32

const (
	FogDisabled FogMode = C.R3D_FOG_DISABLED
	FogLinear   FogMode = C.R3D_FOG_LINEAR
	FogExp2     FogMode = C.R3D_FOG_EXP2
	FogExp      FogMode = C.R3D_FOG_EXP
)

// TonemapMode selects the HDR tone-mapping operator.
type TonemapMode int32

const (
	TonemapLinear   TonemapMode = C.R3D_TONEMAP_LINEAR
	TonemapReinhard TonemapMode = C.R3D_TONEMAP_REINHARD
	TonemapFilmic   TonemapMode = C.R3D_TONEMAP_FILMIC
	TonemapACES     TonemapMode = C.R3D_TONEMAP_ACES
	TonemapAGX      TonemapMode = C.R3D_TONEMAP_AGX
)

// TonemapModeCount is the number of tonemap operators (for cycling).
const TonemapModeCount = int(C.R3D_TONEMAP_COUNT)

// SetBloom configures the bloom post effect.
func SetBloom(mode BloomMode, intensity, filterRadius float32) {
	env := C.R3D_GetEnvironment()
	env.bloom.mode = C.R3D_Bloom(mode)
	env.bloom.intensity = C.float(intensity)
	env.bloom.filterRadius = C.float(filterRadius)
}

// SetSSAO configures screen-space ambient occlusion.
func SetSSAO(enabled bool, intensity, radius, power float32) {
	env := C.R3D_GetEnvironment()
	env.ssao.enabled = C.bool(enabled)
	env.ssao.intensity = C.float(intensity)
	env.ssao.radius = C.float(radius)
	env.ssao.power = C.float(power)
}

// SetFog configures atmospheric fog.
func SetFog(mode FogMode, color Color, start, end, density float32) {
	env := C.R3D_GetEnvironment()
	env.fog.mode = C.R3D_Fog(mode)
	env.fog.color = color.c()
	env.fog.start = C.float(start)
	env.fog.end = C.float(end)
	env.fog.density = C.float(density)
}

// SetDoF configures the depth-of-field effect.
func SetDoF(enabled bool, focusPoint, focusScale, maxBlurSize float32) {
	env := C.R3D_GetEnvironment()
	if enabled {
		env.dof.mode = C.R3D_DOF_ENABLED
	} else {
		env.dof.mode = C.R3D_DOF_DISABLED
	}
	env.dof.focusPoint = C.float(focusPoint)
	env.dof.focusScale = C.float(focusScale)
	env.dof.maxBlurSize = C.float(maxBlurSize)
}

// SetTonemap sets the tone-mapping operator and exposure.
func SetTonemap(mode TonemapMode, exposure float32) {
	env := C.R3D_GetEnvironment()
	env.tonemap.mode = C.R3D_Tonemap(mode)
	env.tonemap.exposure = C.float(exposure)
}

// SetColorAdjustment sets post color grading (brightness/contrast/saturation).
func SetColorAdjustment(brightness, contrast, saturation float32) {
	env := C.R3D_GetEnvironment()
	env.color.brightness = C.float(brightness)
	env.color.contrast = C.float(contrast)
	env.color.saturation = C.float(saturation)
}

// SetSSR toggles and tunes screen-space reflections.
func SetSSR(enabled bool, maxDistance, thickness float32) {
	env := C.R3D_GetEnvironment()
	env.ssr.enabled = C.bool(enabled)
	env.ssr.maxDistance = C.float(maxDistance)
	env.ssr.thickness = C.float(thickness)
}

// SetSSGI toggles and tunes screen-space global illumination.
func SetSSGI(enabled bool, intensity float32) {
	env := C.R3D_GetEnvironment()
	env.ssgi.enabled = C.bool(enabled)
	env.ssgi.intensity = C.float(intensity)
}

// SetSSIL toggles and tunes screen-space indirect lighting.
func SetSSIL(enabled bool, giIntensity, aoIntensity float32) {
	env := C.R3D_GetEnvironment()
	env.ssil.enabled = C.bool(enabled)
	env.ssil.giIntensity = C.float(giIntensity)
	env.ssil.aoIntensity = C.float(aoIntensity)
}

// ---------------------------------------------------------------------------
// Procedural sky
// ---------------------------------------------------------------------------

// ProceduralSky describes a procedurally generated sky (gradient + sun disk).
type ProceduralSky struct{ c C.R3D_ProceduralSky }

// DefaultProceduralSky returns a procedural sky with sensible defaults.
func DefaultProceduralSky() ProceduralSky {
	return ProceduralSky{c: C.r3d_procedural_sky_base()}
}

// SetSun configures the sun disk direction, color and brightness.
func (s *ProceduralSky) SetSun(direction Vector3, color Color, energy float32) {
	s.c.sunDirection = direction.c()
	s.c.sunColor = color.c()
	s.c.sunEnergy = C.float(energy)
}

// SetSkyColors sets the sky gradient colors (zenith / horizon) and brightness.
func (s *ProceduralSky) SetSkyColors(top, horizon Color, energy float32) {
	s.c.skyTopColor = top.c()
	s.c.skyHorizonColor = horizon.c()
	s.c.skyEnergy = C.float(energy)
}

// GenProceduralSky bakes the procedural sky into a cubemap of the given size.
func GenProceduralSky(size int, sky ProceduralSky) Cubemap {
	return Cubemap{c: C.R3D_GenProceduralSky(C.int(size), sky.c)}
}
