package agents

import (
	"math"

	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Agent0 struct {
	base

	CurrentPath   []geom.Vec3
	PathIndex     int
	EscapingTicks uint64
	Solved        bool
}

func NewAgent0() *Agent0 {
	return &Agent0{
		base: base{
			Position: level.StartAgentPosition,
			Grounded: true,
		},
	}
}

var (
	v3HalfXZ = geom.Vec3{
		X: 0.5,
		Y: 0,
		Z: 0.5,
	}
	v3ZeroY = geom.Vec3{
		X: 1,
		Y: 0,
		Z: 1,
	}
)

func (a0 *Agent0) Update(env *Env) {
	a0.base.update()

	const escapeTicks = logic.TPS * 5
	// Goal solved
	runc := 1.
	solved := a0.Position.DistanceTo(env.Goal)+TasingRadius*4 < env.LastKnownAt.DistanceTo(env.Goal)
	if solved {
		a0.EscapingTicks = 0
		a0.CurrentPath = nil
		a0.PathIndex = 0
		a0.Solved = true
	}
	if a0.Solved {
		runc = AgentRunMSMultiplier
		a0.Running = true
	}
	if !a0.Solved {
		// Determine if should escape
		if env.Seen || (env.LastKnownAtUpdate && a0.Position.DistanceTo(env.LastKnownAt) < PlayerVisibilityRadius/2) {
			// Turn in escape mode
			a0.EscapingTicks = escapeTicks
			a0.CurrentPath = nil
			a0.PathIndex = 0
		}
		if a0.EscapingTicks > 0 {
			runc = AgentRunMSMultiplier
			a0.Running = true
			a0.EscapingTicks--
			if a0.EscapingTicks == 0 {
				runc = 1
				/*a0.CurrentPath = nil
				a0.PathIndex = 0*/
			}
		}
		// Pathfinding
		// if player has been seen or heard somewhere else far enough from goal
		if a0.ticks < CheatVisionTicks && env.LastKnownAtUpdate {
			a0.CurrentPath = nil
			a0.PathIndex = 0
		}
	}
	// If no path in mind
	if a0.CurrentPath == nil {
		var goal geom.Vec3
		switch {
		// Solved
		case a0.Solved:
			env.Map.UpdateProximityCosts(env.LastKnownAt, TasingRadius+5)
			goal = env.Goal
		// Try escaping
		case a0.EscapingTicks > 0:
			spot := level.HideSpotsByLevel[env.Level][0]
			dist := env.LastKnownAt.DistanceTo(spot)
			for _, s := range level.HideSpotsByLevel[env.Level] {
				d := env.LastKnownAt.DistanceTo(s)
				if d > dist {
					spot = s
					dist = d
				}
			}
			env.Map.UpdateProximityCosts(env.LastKnownAt, PlayerVisibilityRadius+5)
			goal = spot
		// Try reaching goal safely if "feels" like it is
		case a0.Position.DistanceTo(env.Goal)+TasingRadius+5 < env.LastKnownAt.DistanceTo(env.Goal):
			env.Map.UpdateProximityCosts(env.LastKnownAt, PlayerVisibilityRadius+5)
			goal = env.Goal
		// Go higher to take information
		default:
			spot := level.InfoSpotsByLevel[env.Level][0]
			dist := env.LastKnownAt.DistanceTo(spot)
			for _, s := range level.InfoSpotsByLevel[env.Level] {
				if d := env.LastKnownAt.DistanceTo(s); d > dist {
					spot = s
					dist = d
				}
			}
			env.Map.UpdateProximityCosts(env.LastKnownAt, PlayerVisibilityRadius+5)
			goal = spot
		}
		// Try finding a path
		path, found := env.Map.AStar(
			a0.Position,
			goal,
		)
		if found {
			a0.CurrentPath = path
			a0.PathIndex = 0
		}
	}
	// Follow the path
	current := a0.Position
	if a0.CurrentPath != nil {
		// Block by block
		if a0.PathIndex < len(a0.CurrentPath) {
			x := int(current.X * 10)
			z := int(current.Z * 10)
			if x == int((a0.CurrentPath[a0.PathIndex].X+0.5)*10) &&
				z == int((a0.CurrentPath[a0.PathIndex].Z+0.5)*10) &&
				a0.Grounded {
				a0.PathIndex++
			}
			if a0.PathIndex < len(a0.CurrentPath) {
				next := a0.CurrentPath[a0.PathIndex].Add(v3HalfXZ)
				currHeight := float64(env.Map.At(int(current.X), int(current.Z)).Height)
				nextHeight := float64(env.Map.At(int(next.X), int(next.Z)).Height)
				// Yump
				if nextHeight > currHeight && a0.JumpingTicks == 0 && a0.Grounded {
					a0.JumpingTicks = JumpingTicks
					a0.Grounded = false
				}

				dist := geom.Vec3{current.X, 0, current.Z}.DistanceTo(
					geom.Vec3{next.X, 0, next.Z},
				)
				if dist == 0 {
					// Fine we're just waiting to land
					return
				}
				ms := math.Min(AgentDefaultMS*runc, dist)
				a0.Intent = next.Sub(a0.Position)
				a0.Intent.Y = 0
				a0.Intent = a0.Intent.Normalize().MulN(ms)
			}
		} else {
			// We reached, so let's allow a new path
			a0.CurrentPath = nil
			a0.PathIndex = 0
		}
	}
}

func (a0 *Agent0) HasAbility(ability Ability) bool {
	return false
}

func (a0 *Agent0) Cooldown(ability Ability) float64 {
	return 0
}
