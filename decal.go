package r3d

/*
#include "raylib.h"
#include "r3d.h"

// R3D_DECAL_BASE is a compound-literal macro; cgo can't use it directly.
static R3D_Decal r3d_decal_base(void) { return R3D_DECAL_BASE; }
*/
import "C"

// Decal wraps an R3D_Decal projected onto scene geometry.
type Decal struct{ c C.R3D_Decal }

// NewDecal returns a decal initialized with sensible defaults (R3D_DECAL_BASE).
func NewDecal() Decal { return Decal{c: C.r3d_decal_base()} }

// SetAlbedoMap assigns the decal base-color map.
func (d *Decal) SetAlbedoMap(a AlbedoMap) { d.c.albedo = a.c }

// SetNormalMap assigns the decal normal map.
func (d *Decal) SetNormalMap(n NormalMap) { d.c.normal = n.c }

// SetNormalThreshold sets the max angle (degrees) between surface and decal
// before the decal fades out.
func (d *Decal) SetNormalThreshold(deg float32) { d.c.normalThreshold = C.float(deg) }

// SetFadeWidth sets the edge fade width (degrees).
func (d *Decal) SetFadeWidth(deg float32) { d.c.fadeWidth = C.float(deg) }

// DrawDecal projects the decal at a position with uniform scale (box volume).
func DrawDecal(d Decal, pos Vector3, scale float32) {
	C.R3D_DrawDecal(d.c, pos.c(), C.float(scale))
}

// DrawDecalEx projects the decal with full position/rotation/scale.
func DrawDecalEx(d Decal, pos Vector3, rot Quaternion, scale Vector3) {
	C.R3D_DrawDecalEx(d.c, pos.c(), rot.c(), scale.c())
}
