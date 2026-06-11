package r3d

/*
#include "raylib.h"
*/
import "C"

// RenderTexture wraps a raylib RenderTexture2D (an offscreen framebuffer).
// Used as a target for R3D multi-view rendering (e.g. editor viewports shown
// inside Dear ImGui windows).
type RenderTexture struct{ c C.RenderTexture2D }

// LoadRenderTexture creates an offscreen render target of the given size.
func LoadRenderTexture(width, height int) RenderTexture {
	return RenderTexture{c: C.LoadRenderTexture(C.int(width), C.int(height))}
}

// Unload frees the render target.
func (rt RenderTexture) Unload() { C.UnloadRenderTexture(rt.c) }

// TextureID returns the GL texture id of the color attachment (for ImGui::Image).
func (rt RenderTexture) TextureID() uint32 { return uint32(rt.c.texture.id) }

// Width returns the color texture width.
func (rt RenderTexture) Width() int { return int(rt.c.texture.width) }

// Height returns the color texture height.
func (rt RenderTexture) Height() int { return int(rt.c.texture.height) }

// IsValid reports whether the render target was created successfully.
func (rt RenderTexture) IsValid() bool { return rt.c.id != 0 }

// DrawToScreen blits the render target onto the current framebuffer at the
// given rectangle (handles the vertical flip of GL render textures).
func (rt RenderTexture) DrawToScreen(x, y, width, height float32) {
	src := C.Rectangle{x: 0, y: 0, width: C.float(rt.Width()), height: C.float(-rt.Height())}
	dst := C.Rectangle{x: C.float(x), y: C.float(y), width: C.float(width), height: C.float(height)}
	C.DrawTexturePro(rt.c.texture, src, dst, C.Vector2{}, 0, C.Color{r: 255, g: 255, b: 255, a: 255})
}
