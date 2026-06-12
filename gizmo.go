package r3d

/*
#include "raylib.h"
#include "r3d.h"
*/
import "C"

// Matrix support for editor gizmos: camera view/projection matrices and
// matrix-based mesh drawing. Matrices are passed as column-major [16]float32
// (OpenGL / ImGuizmo order).

func matToArr(m C.Matrix) [16]float32 {
	return [16]float32{
		float32(m.m0), float32(m.m1), float32(m.m2), float32(m.m3),
		float32(m.m4), float32(m.m5), float32(m.m6), float32(m.m7),
		float32(m.m8), float32(m.m9), float32(m.m10), float32(m.m11),
		float32(m.m12), float32(m.m13), float32(m.m14), float32(m.m15),
	}
}

func arrToMat(a [16]float32) C.Matrix {
	return C.Matrix{
		m0: C.float(a[0]), m1: C.float(a[1]), m2: C.float(a[2]), m3: C.float(a[3]),
		m4: C.float(a[4]), m5: C.float(a[5]), m6: C.float(a[6]), m7: C.float(a[7]),
		m8: C.float(a[8]), m9: C.float(a[9]), m10: C.float(a[10]), m11: C.float(a[11]),
		m12: C.float(a[12]), m13: C.float(a[13]), m14: C.float(a[14]), m15: C.float(a[15]),
	}
}

// ViewMatrix returns the camera's view matrix (column-major).
func (rc RCamera) ViewMatrix() [16]float32 {
	return matToArr(C.R3D_GetCameraView(rc.c))
}

// ProjMatrix returns the camera's projection matrix for the given aspect ratio.
func (rc RCamera) ProjMatrix(aspect float64) [16]float32 {
	return matToArr(C.R3D_GetCameraProj(rc.c, C.double(aspect)))
}

// DrawMeshPro draws a mesh with a full transform matrix (column-major).
func DrawMeshPro(mesh Mesh, mat Material, model [16]float32) {
	C.R3D_DrawMeshPro(mesh.c, mat.c, arrToMat(model))
}

// DrawModelPro draws a model with a full transform matrix (column-major).
func DrawModelPro(model Model, transform [16]float32) {
	C.R3D_DrawModelPro(model.c, arrToMat(transform))
}
