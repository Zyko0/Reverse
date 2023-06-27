package agents

import (
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Env struct {
	Map           *level.HMap
	Goal          geom.Vec3
	LastHeard     geom.Vec3
	TimeRemaining uint64
}
