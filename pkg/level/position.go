package level

import (
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

var (
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
