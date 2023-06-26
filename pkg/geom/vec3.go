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

func (v Vec3) String() string {
	return fmt.Sprintf("%.2f,%.2f,%.2f", v.X, v.Y, v.Z)
}
