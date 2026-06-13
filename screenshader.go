package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// ScreenShaderStage selects where a fullscreen post-process pass runs.
type ScreenShaderStage int32

const (
	ScreenStageScene  ScreenShaderStage = C.R3D_SCREEN_SHADER_STAGE_SCENE  // before built-in post (HDR scene)
	ScreenStagePost   ScreenShaderStage = C.R3D_SCREEN_SHADER_STAGE_POST   // after HDR post, before output conversion
	ScreenStageOutput ScreenShaderStage = C.R3D_SCREEN_SHADER_STAGE_OUTPUT // after tonemap (LDR), before AA
	ScreenStageFinal  ScreenShaderStage = C.R3D_SCREEN_SHADER_STAGE_FINAL  // after AA, before final blit
)

// ScreenShader wraps an R3D_ScreenShader (a custom fullscreen post-process pass).
// Its GLSL defines `void fragment()` and writes the output to COLOR, reading the
// input via SampleColor(TEXCOORD) / FetchColor(PIXCOORD).
type ScreenShader struct{ c *C.R3D_ScreenShader }

// LoadScreenShaderFromMemory compiles a screen shader from source code.
func LoadScreenShaderFromMemory(code string) ScreenShader {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	return ScreenShader{c: C.R3D_LoadScreenShaderFromMemory(c)}
}

// LoadScreenShader loads a screen shader from a .glsl file.
func LoadScreenShader(filePath string) ScreenShader {
	c := C.CString(filePath)
	defer C.free(unsafe.Pointer(c))
	return ScreenShader{c: C.R3D_LoadScreenShader(c)}
}

// Unload frees the screen shader.
func (s ScreenShader) Unload() { C.R3D_UnloadScreenShader(s.c) }

// IsValid reports whether the shader handle is non-nil.
func (s ScreenShader) IsValid() bool { return s.c != nil }

func (s ScreenShader) setUniform(name string, ptr unsafe.Pointer) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.R3D_SetScreenShaderUniform(s.c, cn, ptr)
}

// SetUniformFloat sets a float uniform.
func (s ScreenShader) SetUniformFloat(name string, v float32) {
	cv := C.float(v)
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetUniformInt sets an int (or bool, as 0/1) uniform.
func (s ScreenShader) SetUniformInt(name string, v int) {
	cv := C.int(v)
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetUniformVec3 sets a vec3 uniform.
func (s ScreenShader) SetUniformVec3(name string, v Vector3) {
	cv := v.c()
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetSampler binds a texture to a sampler uniform.
func (s ScreenShader) SetSampler(name string, t Texture) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.R3D_SetScreenShaderSampler(s.c, cn, t.c)
}

// SetScreenShaderChain installs an ordered chain of screen shaders at a stage
// (replacing the current chain). Passing an empty slice clears the stage.
func SetScreenShaderChain(stage ScreenShaderStage, shaders []ScreenShader) {
	if len(shaders) == 0 {
		C.R3D_SetScreenShaderChain(C.R3D_ScreenShaderStage(stage), nil, 0)
		return
	}
	arr := make([]*C.R3D_ScreenShader, len(shaders))
	for i, s := range shaders {
		arr[i] = s.c
	}
	C.R3D_SetScreenShaderChain(C.R3D_ScreenShaderStage(stage), &arr[0], C.int(len(arr)))
}

// SetScreenShader is a convenience for a single-shader chain at a stage (nil
// shader clears the stage).
func SetScreenShader(stage ScreenShaderStage, s ScreenShader) {
	if !s.IsValid() {
		SetScreenShaderChain(stage, nil)
		return
	}
	SetScreenShaderChain(stage, []ScreenShader{s})
}
