package agents

import (
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Agent0 struct {
	base

	CurrentPath []geom.Vec3
	PathIndex   int
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
	a0.base.update()
	// Pathfinding
	if a0.CurrentPath == nil {
		path, found := env.Map.BFS(
			a0.Position,
			env.Goal,
			jumpDistanceAwarenessByAgent[0],
			false,
		)
		if found {
			a0.CurrentPath = path
			a0.PathIndex = 0
		}
	}
	// Get closer to target
	if a0.CurrentPath != nil && a0.PathIndex < len(a0.CurrentPath) {
		current := a0.Position.AsUHashXZ()
		for ; a0.PathIndex < len(a0.CurrentPath) && current == a0.CurrentPath[a0.PathIndex].AsUHashXZ(); a0.PathIndex++ {
		}
		if a0.PathIndex < len(a0.CurrentPath) {
			next := a0.CurrentPath[a0.PathIndex].AsUHashXZ()
			a0.Intent = geom.Vec3{}.FromUHashXZ(next)
			height := float64(env.Map.At(int(a0.Intent.X), int(a0.Intent.Z)).Height)
			// Yump
			if height > a0.Position.Y && a0.JumpingTicks == 0 {
				a0.JumpingTicks = JumpingTicks
				a0.Grounded = false
			}
			a0.Intent = a0.Intent.Sub(a0.Position.Floor())
			a0.Intent.Y = 0
			a0.Intent = a0.Intent.Normalize().MulN(AgentDefaultMS)
		}
	}
}

func (a0 *Agent0) HasAbility(ability Ability) bool {
	return false
}
