package agents

import (
	"fmt"
	"time"

	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Agent0 struct {
	base

	CurrentPath []geom.Vec3
}

func NewAgent0() *Agent0 {
	return &Agent0{
		base: base{
			Position: level.StartAgentPosition,
			Grounded: true,
		},
	}
}

func (a0 *Agent0) Update(env *Env) {
	if a0.CurrentPath == nil {
		t := time.Now()
		path, found := env.Map.BFS(
			a0.Position,
			env.Goal,
			jumpDistanceAwarenessByAgent[0],
			false,
		)
		fmt.Println("start", a0.Position, "goal", env.Goal, "awareness", jumpDistanceAwarenessByAgent[0])
		fmt.Println("time spent", time.Since(t), "found", found, "len", len(path))
		if found {
			a0.CurrentPath = path
		}
		t = time.Now()
		//fmt.Println("found", path)
	}
}

func (a0 *Agent0) HasAbility(ability Ability) bool {
	return false
}
