package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"
*/
import "C"

import "unsafe"

// Texture wraps a raylib Texture (used for custom surface-shader samplers).
type Texture struct{ c C.Texture }

// LoadTexture loads a GPU texture from an image file.
func LoadTexture(fileName string) Texture {
	c := C.CString(fileName)
	defer C.free(unsafe.Pointer(c))
	return Texture{c: C.LoadTexture(c)}
}

// Unload frees the texture.
func (t Texture) Unload() { C.UnloadTexture(t.c) }

// ID returns the OpenGL texture id (e.g. to use as an ImGui image handle).
func (t Texture) ID() uint32 { return uint32(t.c.id) }

// Width returns the texture width in pixels.
func (t Texture) Width() int { return int(t.c.width) }

// Height returns the texture height in pixels.
func (t Texture) Height() int { return int(t.c.height) }

// IsValid reports whether the texture handle is a non-zero GL id.
func (t Texture) IsValid() bool { return t.c.id != 0 }

// SurfaceShader wraps an R3D_SurfaceShader (custom material/decal shader).
type SurfaceShader struct{ c *C.R3D_SurfaceShader }

// LoadSurfaceShader loads a surface shader from a .glsl file.
func LoadSurfaceShader(filePath string) SurfaceShader {
	c := C.CString(filePath)
	defer C.free(unsafe.Pointer(c))
	return SurfaceShader{c: C.R3D_LoadSurfaceShader(c)}
}

// LoadSurfaceShaderFromMemory compiles a surface shader from source code.
func LoadSurfaceShaderFromMemory(code string) SurfaceShader {
	c := C.CString(code)
	defer C.free(unsafe.Pointer(c))
	return SurfaceShader{c: C.R3D_LoadSurfaceShaderFromMemory(c)}
}

// Unload frees the surface shader.
func (s SurfaceShader) Unload() { C.R3D_UnloadSurfaceShader(s.c) }

// IsValid reports whether the shader handle is non-nil.
func (s SurfaceShader) IsValid() bool { return s.c != nil }

func (s SurfaceShader) setUniform(name string, ptr unsafe.Pointer) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.R3D_SetSurfaceShaderUniform(s.c, cn, ptr)
}

// SetUniformFloat sets a float uniform.
func (s SurfaceShader) SetUniformFloat(name string, v float32) {
	cv := C.float(v)
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetUniformInt sets an int (or bool, as 0/1) uniform.
func (s SurfaceShader) SetUniformInt(name string, v int) {
	cv := C.int(v)
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetUniformVec3 sets a vec3 uniform.
func (s SurfaceShader) SetUniformVec3(name string, v Vector3) {
	cv := v.c()
	s.setUniform(name, unsafe.Pointer(&cv))
}

// SetSampler binds a texture to a sampler uniform.
func (s SurfaceShader) SetSampler(name string, t Texture) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.R3D_SetSurfaceShaderSampler(s.c, cn, t.c)
}

// SetShader applies a custom surface shader to this material.
func (m *Material) SetShader(s SurfaceShader) { m.c.shader = s.c }

// ClearShader removes any custom surface shader from this material.
func (m *Material) ClearShader() { m.c.shader = nil }

// TextureWrap modes (raylib TextureWrap).
type TextureWrap int32

const (
	WrapRepeat TextureWrap = C.TEXTURE_WRAP_REPEAT
	WrapClamp  TextureWrap = C.TEXTURE_WRAP_CLAMP
	WrapMirror TextureWrap = C.TEXTURE_WRAP_MIRROR_REPEAT
)

// SetTextureWrap sets the default texture wrap mode used by R3D loaders.
func SetTextureWrap(w TextureWrap) { C.R3D_SetTextureWrap(C.TextureWrap(w)) }
