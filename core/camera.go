package core

import (
	"math"

	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	cameraOffset = geom.Vec3{
		X: 0,
		Y: -128,
		Z: -128,
	}
	camDir = geom.Vec3{
		X: -1,
		Y: 0,
		Z: -1,
	}.Normalize()
)

type Camera struct {
	ticks       uint64
	lastCursorX int
	zoomValue   float64

	HAngle    float64
	Position  geom.Vec3
	Direction geom.Vec3
	Zoom      float64
}

func newCamera(position geom.Vec3) *Camera {
	p := position.Add(cameraOffset)

	return &Camera{
		ticks:       0,
		lastCursorX: 0,

		HAngle:    0.,
		Position:  p,
		Direction: camDir,
		Zoom:      1.,
	}
}

func (c *Camera) UpdateDirection(playerPosition geom.Vec3) {
	const (
		camRotateSens = 0.001
	)

	x, _ := ebiten.CursorPosition()
	// Note: hack because lastcursor is to 0 by default so the gap is too huge
	if c.ticks > 1 {
		// Record new horizontal rotation
		if delta := x - c.lastCursorX; delta != 0 {
			c.lastCursorX = x
			c.HAngle = math.Mod(
				c.HAngle+float64(delta)*camRotateSens,
				2*math.Pi,
			)
		}
	} else {
		c.lastCursorX = x
	}

	// Update camera position
	pxz := geom.Vec2{
		X: camDir.X,
		Y: camDir.Z,
	}.Rotate(
		c.HAngle,
	)
	c.Direction.X = pxz.X
	c.Direction.Z = pxz.Y
	c.Direction = c.Direction.Normalize()

	/*c.Position.X = playerPosition.X + pxz.X
	c.Position.Z = playerPosition.Z + pxz.Y*/
	// Update camera direction
	//c.Direction = c.Position.Sub(playerPosition).Normalize()
}

func (c *Camera) UpdatePosition(playerPosition geom.Vec3) {
	//c.Position = playerPosition.Add(c.Direction.MulN(cameraOffsetMag))
}

var (
	minZoom = math.Log(0.25)
	maxZoom = math.Log(4)
)

func (c *Camera) Update() {
	const (
		zoomSens = 0.05
	)
	// TODO: update zoom
	_, y := ebiten.Wheel()
	c.zoomValue += y * zoomSens
	if c.zoomValue < 0 {
		c.zoomValue = 0
	}
	if c.zoomValue > 1 {
		c.zoomValue = 1
	}
	// Actual lerp zoom
	c.Zoom = math.Exp(minZoom + (maxZoom-minZoom)*c.zoomValue)

	c.ticks++
}

// AsUniforms returns two vec3 for position and direction respectively
func (c *Camera) AsUniforms() []float32 {
	/*fmt.Printf("pos: %.2f %.2f %.2f\n",
		float32(c.Position.X/MapWidth*2-1),
		float32(c.Position.Y/MapHeight),
		float32(c.Position.Z/MapDepth*2-1),
	)*/

	return []float32{
		// Position
		float32(c.Position.X/logic.MapWidth*2 - 1),
		float32(c.Position.Y / logic.MapHeight),
		float32(c.Position.Z/logic.MapDepth*2 - 1),
		// Direction
		float32(c.Direction.X),
		float32(c.Direction.Y),
		float32(c.Direction.Z),
	}
}
