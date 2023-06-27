package geom

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Zero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z

	return v
}

func (v Vec3) Sub(v2 Vec3) Vec3 {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z

	return v
}

func (v Vec3) Mul(v2 Vec3) Vec3 {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z

	return v
}

func (v Vec3) MulN(n float64) Vec3 {
	v.X *= n
	v.Y *= n
	v.Z *= n

	return v
}

func (v Vec3) Floor() Vec3 {
	v.X = math.Floor(v.X)
	v.Y = math.Floor(v.Y)
	v.Z = math.Floor(v.Z)

	return v
}

func (v Vec3) Normalize() Vec3 {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	v.X /= mag
	v.Y /= mag
	v.Z /= mag

	return v
}

func (v Vec3) DistanceTo(v2 Vec3) float64 {
	d := (v2.X-v.X)*(v2.X-v.X) + (v2.Y-v.Y)*(v2.Y-v.Y) + (v2.Z-v.Z)*(v2.Z-v.Z)

	return math.Sqrt(d)
}

func (v Vec3) FromUHash(hash uint32) Vec3 {
	v.X = float64(hash & 255)
	v.Z = float64((hash >> 8) & 255)
	v.Y = float64((hash >> 16) & 255)
	return v
}

func (v Vec3) AsUHash() uint32 {
	return uint32(v.X) | (uint32(v.Z) << 8) | (uint32(v.Y) << 16)
}

func (v Vec3) FromUHashXZ(hash uint16) Vec3 {
	v.X = float64(hash & 255)
	v.Z = float64((hash >> 8) & 255)
	return v
}

func (v Vec3) AsUHashXZ() uint16 {
	return uint16(v.X) | (uint16(v.Z) << 8)
}

func (v Vec3) String() string {
	return fmt.Sprintf("%.2f,%.2f,%.2f", v.X, v.Y, v.Z)
}
