package geom

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	v.X += v2.X
	v.Y += v2.Y

	return v
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	v.X -= v2.X
	v.Y -= v2.Y

	return v
}

func (v Vec2) Mul(v2 Vec2) Vec2 {
	v.X *= v2.X
	v.Y *= v2.Y

	return v
}

func (v Vec2) MulN(n float64) Vec2 {
	v.X *= n
	v.Y *= n

	return v
}

func (v Vec2) Div(v2 Vec2) Vec2 {
	v.X /= v2.X
	v.Y /= v2.Y

	return v
}

func (v Vec2) Rotate(angle float64) Vec2 {
	s, c := math.Sincos(angle)
	x := v.X*c - v.Y*s
	y := v.X*s + v.Y*c
	v.X, v.Y = x, y

	return v
}

func (v Vec2) RotateAroundCenter(center Vec2, angle float64) Vec2 {
	s, c := math.Sincos(angle)
	// Sub
	v.X -= center.X
	v.Y -= center.Y
	x := v.X*c - v.Y*s
	y := v.X*s + v.Y*c
	v.X, v.Y = x, y
	// Add back
	v.X += center.X
	v.Y += center.Y

	return v
}

func (v Vec2) Normalize() Vec2 {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y)
	v.X /= mag
	v.Y /= mag

	return v
}

func (v Vec2) String() string {
	return fmt.Sprintf("%.2f,%.2f", v.X, v.Y)
}
