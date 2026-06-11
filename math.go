package r3d

import "math"

// Vector3 helpers ------------------------------------------------------------

// Add returns v + o.
func (v Vector3) Add(o Vector3) Vector3 { return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z} }

// Scale returns v * s.
func (v Vector3) Scale(s float32) Vector3 { return Vector3{v.X * s, v.Y * s, v.Z * s} }

// Length returns the Euclidean length of v.
func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize returns v scaled to unit length (or v if it is zero-length).
func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if l == 0 {
		return v
	}
	return v.Scale(1 / l)
}

// Quaternion helpers ---------------------------------------------------------

// QuaternionFromAxisAngle builds a rotation quaternion (angle in radians).
func QuaternionFromAxisAngle(axis Vector3, angle float32) Quaternion {
	axis = axis.Normalize()
	half := float64(angle) * 0.5
	s := float32(math.Sin(half))
	return Quaternion{axis.X * s, axis.Y * s, axis.Z * s, float32(math.Cos(half))}
}

// QuaternionFromEuler builds a rotation quaternion from pitch/yaw/roll (radians),
// matching raylib's convention.
func QuaternionFromEuler(pitch, yaw, roll float32) Quaternion {
	x0 := float32(math.Cos(float64(pitch) * 0.5))
	x1 := float32(math.Sin(float64(pitch) * 0.5))
	y0 := float32(math.Cos(float64(yaw) * 0.5))
	y1 := float32(math.Sin(float64(yaw) * 0.5))
	z0 := float32(math.Cos(float64(roll) * 0.5))
	z1 := float32(math.Sin(float64(roll) * 0.5))
	return Quaternion{
		X: x1*y0*z0 - x0*y1*z1,
		Y: x0*y1*z0 + x1*y0*z1,
		Z: x0*y0*z1 - x1*y1*z0,
		W: x0*y0*z0 + x1*y1*z1,
	}
}

// Mul returns the Hamilton product q*o (apply o then q).
func (q Quaternion) Mul(o Quaternion) Quaternion {
	return Quaternion{
		X: q.X*o.W + q.W*o.X + q.Y*o.Z - q.Z*o.Y,
		Y: q.Y*o.W + q.W*o.Y + q.Z*o.X - q.X*o.Z,
		Z: q.Z*o.W + q.W*o.Z + q.X*o.Y - q.Y*o.X,
		W: q.W*o.W - q.X*o.X - q.Y*o.Y - q.Z*o.Z,
	}
}
