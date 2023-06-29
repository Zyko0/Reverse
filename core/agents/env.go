package agents

import (
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Env struct {
	Level             int
	Map               *level.HMap
	Goal              geom.Vec3
	LastKnownAt       geom.Vec3
	LastKnownAtUpdate bool
	TimeRemaining     uint64
	Seen              bool
	CanSeePlayer      bool
}
