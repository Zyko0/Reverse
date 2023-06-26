package core

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	ticks       uint64
	lastCursorX int
	zoomValue   float64

	HAngle float64
	Zoom   float64
}

func newCamera() *Camera {
	return &Camera{
		ticks:       0,
		lastCursorX: 0,

		HAngle: 0,
		Zoom:   1,
	}
}

var (
	minZoom = math.Log(0.1) //math.Log(1) // TODO: restore
	maxZoom = math.Log(7)
)

func sign(n float64) float64 {
	if n == 0 {
		return 0
	}
	if n > 0 {
		return 1
	}
	return -1
}

func (c *Camera) Update() {
	const (
		zoomSens      = 0.05
		camRotateSens = 0.001
	)

	// Horizontal angle
	x, _ := ebiten.CursorPosition()
	// Note: hack because lastcursor is to 0 by default so the gap is too huge
	if c.ticks > 1 {
		// Record new horizontal rotation
		if delta := x - c.lastCursorX; delta != 0 {
			c.lastCursorX = x
			c.HAngle = math.Mod(
				c.HAngle-float64(delta)*camRotateSens,
				2*math.Pi,
			)
		}
	} else {
		c.lastCursorX = x
	}
	// Zoom
	_, y := ebiten.Wheel()
	y = sign(y)
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
