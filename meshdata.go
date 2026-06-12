package r3d

/*
#include <stdlib.h>
#include "raylib.h"
#include "r3d.h"

// Helpers to fill an R3D_MeshData from Go without unsafe pointer arithmetic.
static void r3d_md_set_vertex(R3D_MeshData* d, int i, Vector3 pos, Vector2 uv, Vector3 nrm) {
    d->vertices[i] = R3D_MakeVertex(pos, uv, nrm, (Vector4){1.0f, 0.0f, 0.0f, 1.0f}, WHITE);
}
static void r3d_md_set_index(R3D_MeshData* d, int i, unsigned int idx) {
    d->indices[i] = idx;
}
*/
import "C"

// LoadMeshFromData uploads a triangle mesh built from CPU arrays (positions,
// per-vertex normals, UVs, and triangle indices) to the GPU. The slices must
// satisfy len(normals)==len(positions) and len(uvs)==len(positions); len(indices)
// must be a multiple of 3. Tangents are derived from the normals/UVs. The
// returned Mesh is owned by the caller (Unload when done).
func LoadMeshFromData(positions, normals []Vector3, uvs []Vector2, indices []uint32) Mesh {
	vc, ic := len(positions), len(indices)
	if vc == 0 || ic == 0 {
		return Mesh{}
	}
	data := C.R3D_LoadMeshData(C.int(vc), C.int(ic))
	for i := 0; i < vc; i++ {
		var n Vector3
		if i < len(normals) {
			n = normals[i]
		}
		var uv Vector2
		if i < len(uvs) {
			uv = uvs[i]
		}
		C.r3d_md_set_vertex(&data, C.int(i), positions[i].c(), uv.c(), n.c())
	}
	for i := 0; i < ic; i++ {
		C.r3d_md_set_index(&data, C.int(i), C.uint(indices[i]))
	}
	data.vertexCount = C.int(vc)
	data.indexCount = C.int(ic)

	C.R3D_GenMeshDataTangents(&data, C.R3D_PRIMITIVE_TRIANGLES)
	mesh := Mesh{c: C.R3D_LoadMesh(C.R3D_PRIMITIVE_TRIANGLES, data, nil)}
	C.R3D_UnloadMeshData(data)
	return mesh
}

// UpdateMeshFromData re-uploads new geometry into an existing mesh in place
// (same constraints as LoadMeshFromData). Returns false on failure.
func (m *Mesh) UpdateMeshFromData(positions, normals []Vector3, uvs []Vector2, indices []uint32) bool {
	vc, ic := len(positions), len(indices)
	if vc == 0 || ic == 0 {
		return false
	}
	data := C.R3D_LoadMeshData(C.int(vc), C.int(ic))
	for i := 0; i < vc; i++ {
		var n Vector3
		if i < len(normals) {
			n = normals[i]
		}
		var uv Vector2
		if i < len(uvs) {
			uv = uvs[i]
		}
		C.r3d_md_set_vertex(&data, C.int(i), positions[i].c(), uv.c(), n.c())
	}
	for i := 0; i < ic; i++ {
		C.r3d_md_set_index(&data, C.int(i), C.uint(indices[i]))
	}
	data.vertexCount = C.int(vc)
	data.indexCount = C.int(ic)
	C.R3D_GenMeshDataTangents(&data, C.R3D_PRIMITIVE_TRIANGLES)
	ok := bool(C.R3D_UpdateMesh(&m.c, data, nil))
	C.R3D_UnloadMeshData(data)
	return ok
}
