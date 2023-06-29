package level

import (
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

// Positions
var (
	GoalPosition = geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 1,
		Z: logic.MapDepth * 0.05,
	}
	StartPlayerPosition = geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 1,
		Z: logic.MapDepth * 0.1,
	}
	StartAgentPosition = geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 1,
		Z: logic.MapDepth * 0.9,
	}
)

var (
	LevelsTime = []uint64{
		0: logic.TPS * 60 * 5,
		1: logic.TPS * 60 * 5,
		2: logic.TPS * 60 * 5,
		3: logic.TPS * 60 * 5,
	}

	HideSpotsByLevel = [][]geom.Vec3{
		0: {
			// Labyrinths
			{
				X: 108,
				Y: 2,
				Z: 176,
			},
			{
				X: 147,
				Y: 2,
				Z: 81,
			},
			// First floor
			{
				X: 103,
				Y: 10,
				Z: 156,
			},
			{
				X: 103,
				Y: 10,
				Z: 149,
			},
			{
				X: 152,
				Y: 10,
				Z: 156,
			},
			{
				X: 152,
				Y: 10,
				Z: 149,
			},
			{
				X: 152,
				Y: 10,
				Z: 108,
			},
			{
				X: 152,
				Y: 10,
				Z: 101,
			},
			{
				X: 103,
				Y: 10,
				Z: 101,
			},
			{
				X: 103,
				Y: 10,
				Z: 108,
			},
		},
	}
	InfoSpotsByLevel = [][]geom.Vec3{
		0: {
			// Middle tower
			{
				X: 128,
				Y: 38,
				Z: 128,
			},
			{
				X: 127,
				Y: 38,
				Z: 129,
			},
			// Ebitens
			{
				X: 134,
				Y: 11,
				Z: 157,
			},
			{
				X: 121,
				Y: 11,
				Z: 157,
			},
			{
				X: 121,
				Y: 11,
				Z: 100,
			},
			{
				X: 134,
				Y: 11,
				Z: 100,
			},
		},
		1: {
			// Middle tower
			{
				X: 128,
				Y: 38,
				Z: 128,
			},
			{
				X: 127,
				Y: 38,
				Z: 129,
			},
			// Columns
			{
				X: 83,
				Y: 33,
				Z: 173,
			},
			{
				X: 173,
				Y: 7,
				Z: 173,
			},
			{
				X: 173,
				Y: 33,
				Z: 83,
			},
			{
				X: 83,
				Y: 7,
				Z: 83,
			},
		},
		2: nil,
		3: nil,
	}
)
