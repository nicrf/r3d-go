package r3d

/*
#include "raylib.h"
#include "r3d.h"
*/
import "C"

// Projection selects a camera projection mode.
type Projection int32

const (
	ProjectionPerspective  Projection = C.R3D_PROJECTION_PERSPECTIVE
	ProjectionOrthographic Projection = C.R3D_PROJECTION_ORTHOGRAPHIC
)

// Rectangle mirrors raylib's Rectangle (used for viewports).
type Rectangle struct{ X, Y, Width, Height float32 }

func (r Rectangle) c() C.Rectangle {
	return C.Rectangle{x: C.float(r.X), y: C.float(r.Y), width: C.float(r.Width), height: C.float(r.Height)}
}

// RCamera wraps R3D's native camera (position + quaternion orientation),
// used by multi-view rendering (BeginPro).
type RCamera struct{ c C.R3D_Camera }

// CameraFromRL converts a raylib Camera3D into an R3D camera.
func CameraFromRL(cam Camera3D) RCamera { return RCamera{c: C.R3D_CameraFromRL(cam.c())} }

// SetPosition sets the camera world position.
func (rc *RCamera) SetPosition(p Vector3) { rc.c.position = p.c() }

// LookAt orients the camera toward target with the given up vector.
func (rc *RCamera) LookAt(target, up Vector3) {
	C.R3D_CameraLookAt(&rc.c, target.c(), up.c())
}

// View describes a render pass: a camera, a viewport, and an optional offscreen
// target. If Target is nil the backbuffer is used; otherwise R3D renders into
// the given RenderTexture (e.g. an editor viewport shown in a Dear ImGui window).
type View struct {
	Camera   RCamera
	Viewport Rectangle
	Target   *RenderTexture
}

// BeginPro starts an R3D render pass for the view (split-screen, minimaps,
// render-to-texture). Pair with End.
func BeginPro(v View) {
	cv := C.R3D_View{
		camera:   v.Camera.c,
		viewport: v.Viewport.c(),
	}
	if v.Target != nil {
		cv.target = v.Target.c
	}
	C.R3D_BeginPro(cv)
}
