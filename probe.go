package r3d

/*
#include "raylib.h"
#include "r3d.h"
*/
import "C"

// ProbeFlags selects which IBL products a reflection probe captures.
type ProbeFlags uint32

const (
	ProbeIllumination ProbeFlags = 1 << 0 // diffuse irradiance
	ProbeReflection   ProbeFlags = 1 << 1 // specular prefiltered map
)

// ProbeUpdateMode controls how often a probe re-bakes its environment.
type ProbeUpdateMode int32

const (
	ProbeUpdateOnce   ProbeUpdateMode = C.R3D_PROBE_UPDATE_ONCE
	ProbeUpdateAlways ProbeUpdateMode = C.R3D_PROBE_UPDATE_ALWAYS
)

// Probe is a handle to a local reflection/illumination probe.
type Probe int32

// CreateProbe creates a reflection probe capturing the requested products.
func CreateProbe(flags ProbeFlags) Probe {
	return Probe(C.R3D_CreateProbe(C.R3D_ProbeFlags(flags)))
}

// Destroy removes the probe.
func (p Probe) Destroy() { C.R3D_DestroyProbe(C.R3D_Probe(p)) }

// SetActive enables or disables the probe.
func (p Probe) SetActive(active bool) { C.R3D_SetProbeActive(C.R3D_Probe(p), C.bool(active)) }

// SetUpdateMode sets the probe refresh policy.
func (p Probe) SetUpdateMode(mode ProbeUpdateMode) {
	C.R3D_SetProbeUpdateMode(C.R3D_Probe(p), C.R3D_ProbeUpdateMode(mode))
}

// SetPosition sets the probe world position.
func (p Probe) SetPosition(pos Vector3) { C.R3D_SetProbePosition(C.R3D_Probe(p), pos.c()) }

// SetRange sets the probe influence radius.
func (p Probe) SetRange(r float32) { C.R3D_SetProbeRange(C.R3D_Probe(p), C.float(r)) }

// SetFalloff sets the probe influence falloff.
func (p Probe) SetFalloff(f float32) { C.R3D_SetProbeFalloff(C.R3D_Probe(p), C.float(f)) }

// SetInterior marks the probe as an interior (bounded) probe.
func (p Probe) SetInterior(interior bool) {
	C.R3D_SetProbeInterior(C.R3D_Probe(p), C.bool(interior))
}

// SetShadows toggles whether the probe bakes shadows.
func (p Probe) SetShadows(enabled bool) {
	C.R3D_SetProbeShadows(C.R3D_Probe(p), C.bool(enabled))
}
