package r3d

/*
#include "raylib.h"
#include "rlgl.h"
*/
import "C"

import "unsafe"

// Low-level rlgl + raylib 2D state, exposed so external code (e.g. a Dear ImGui
// backend) can render immediate-mode geometry into the current frame.

// rlgl draw modes.
const (
	RlLines     = C.RL_LINES
	RlTriangles = C.RL_TRIANGLES
	RlQuads     = C.RL_QUADS
)

// Pixel formats (subset).
const PixelFormatR8G8B8A8 = C.RL_PIXELFORMAT_UNCOMPRESSED_R8G8B8A8

// Blend modes (raylib BlendMode).
const (
	BlendAlpha            = C.BLEND_ALPHA
	BlendModeAdditive     = C.BLEND_ADDITIVE
	BlendModeMultiplied   = C.BLEND_MULTIPLIED
	BlendModeAddColors    = C.BLEND_ADD_COLORS
	BlendModeSubtractCols = C.BLEND_SUBTRACT_COLORS
)

// RlBegin starts immediate-mode vertex submission in the given mode.
func RlBegin(mode int) { C.rlBegin(C.int(mode)) }

// RlEnd finishes immediate-mode vertex submission.
func RlEnd() { C.rlEnd() }

// RlVertex2f submits a 2D vertex position.
func RlVertex2f(x, y float32) { C.rlVertex2f(C.float(x), C.float(y)) }

// RlVertex3f submits a 3D vertex position.
func RlVertex3f(x, y, z float32) { C.rlVertex3f(C.float(x), C.float(y), C.float(z)) }

// RlTexCoord2f submits a texture coordinate for the next vertex.
func RlTexCoord2f(u, v float32) { C.rlTexCoord2f(C.float(u), C.float(v)) }

// RlColor4ub submits a color for the next vertex.
func RlColor4ub(r, g, b, a uint8) {
	C.rlColor4ub(C.uchar(r), C.uchar(g), C.uchar(b), C.uchar(a))
}

// RlSetTexture binds a GL texture id for subsequent immediate-mode geometry (0 = none).
func RlSetTexture(id uint32) { C.rlSetTexture(C.uint(id)) }

// RlDrawRenderBatchActive flushes the current render batch.
func RlDrawRenderBatchActive() { C.rlDrawRenderBatchActive() }

// RlEnableDepthTest / RlDisableDepthTest toggle the depth test.
func RlEnableDepthTest()  { C.rlEnableDepthTest() }
func RlDisableDepthTest() { C.rlDisableDepthTest() }

// RlEnableBackfaceCulling / RlDisableBackfaceCulling toggle face culling.
func RlEnableBackfaceCulling()  { C.rlEnableBackfaceCulling() }
func RlDisableBackfaceCulling() { C.rlDisableBackfaceCulling() }

// RlLoadTexture uploads pixel data to the GPU and returns its GL texture id.
func RlLoadTexture(data []byte, width, height, format, mipmaps int) uint32 {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	return uint32(C.rlLoadTexture(p, C.int(width), C.int(height), C.int(format), C.int(mipmaps)))
}

// RlUpdateTexture updates a region of an existing GPU texture.
func RlUpdateTexture(id uint32, offsetX, offsetY, width, height, format int, data []byte) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	C.rlUpdateTexture(C.uint(id), C.int(offsetX), C.int(offsetY), C.int(width), C.int(height), C.int(format), p)
}

// RlUnloadTexture frees a GL texture by id.
func RlUnloadTexture(id uint32) { C.rlUnloadTexture(C.uint(id)) }

// DrawGLTexture draws a raw GL texture id (RGBA8) onto the current framebuffer,
// scaled. Useful for debugging texture contents.
func DrawGLTexture(id uint32, srcW, srcH int, x, y, scale float32) {
	tex := C.Texture2D{
		id: C.uint(id), width: C.int(srcW), height: C.int(srcH),
		mipmaps: 1, format: C.PIXELFORMAT_UNCOMPRESSED_R8G8B8A8,
	}
	src := C.Rectangle{x: 0, y: 0, width: C.float(srcW), height: C.float(srcH)}
	dst := C.Rectangle{x: C.float(x), y: C.float(y), width: C.float(srcW) * C.float(scale), height: C.float(srcH) * C.float(scale)}
	C.DrawTexturePro(tex, src, dst, C.Vector2{}, 0, C.Color{r: 255, g: 255, b: 255, a: 255})
}

// BeginBlendMode / EndBlendMode set the active blend mode (raylib).
func BeginBlendMode(mode int) { C.BeginBlendMode(C.int(mode)) }
func EndBlendMode()           { C.EndBlendMode() }

// BeginScissorMode / EndScissorMode clip subsequent drawing to a rectangle (raylib, screen pixels).
func BeginScissorMode(x, y, width, height int) {
	C.BeginScissorMode(C.int(x), C.int(y), C.int(width), C.int(height))
}
func EndScissorMode() { C.EndScissorMode() }
