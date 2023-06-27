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
		Y: 5,
		Z: logic.MapDepth * 0.1,
	}
	StartAgentPosition = geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 5,
		Z: logic.MapDepth * 0.9,
	}
)

// Time level
var (
	LevelsTime = []uint64{
		0: logic.TPS * 60 * 2,
		1: logic.TPS * 60 * 2,
		2: logic.TPS * 60 * 2,
		3: logic.TPS * 60 * 5,
	}
)
