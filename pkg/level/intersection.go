package level

import (
	"math"

	"github.com/Zyko0/Reverse/pkg/geom"
)

// Lazy
func (lm *Map) CastRay(src, dst geom.Vec3, max float64) bool {
	// Test the ray from head to head
	src.Y, dst.Y = src.Y+0.49, dst.Y+0.49
	dist := src.DistanceTo(dst)
	if dist > max {
		return false
	}
	max = dist
	// DDA
	dx := dst.X - src.X
	dy := dst.Y - src.Y
	dz := dst.Z - src.Z
	step := 0.
	if math.Abs(dx) >= math.Abs(dy) {
		step = math.Abs(dx)
	} else {
		step = math.Abs(dy)
	}
	if math.Abs(dz) >= step {
		step = math.Abs(dz)
	}
	dx /= step
	dy /= step
	dz /= step
	x := src.X
	y := src.Y
	z := src.Z
	for i := 1.; i <= step; i++ {
		vox := lm.HeightMap[int(z)][int(x)]
		if y <= float64(vox.Height) {
			return false
		}
		x += dx
		y += dy
		z += dz
	}

	return true
}
