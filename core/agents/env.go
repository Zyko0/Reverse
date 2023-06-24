package agents

import (
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Env struct {
	Map       level.Map
	LastHeard geom.Vec3
}
