package core

import (
	"math"

	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

func (g *Game) getNeighbourColumns(x, y, z int) []geom.Vec3 {
	var neighbours []geom.Vec3

	for xoff := -1; xoff <= 1; xoff++ {
		nx := x + xoff
		if nx < 0 || nx > logic.MapWidth {
			continue
		}
		/*for yoff := -1; yoff <= 1.; yoff++ {
		ny := y+yoff
		if ny < 0 || ny > logic.MapHeight {
			continue
		}*/
		for zoff := -1; zoff <= 1.; zoff++ {
			nz := z + zoff
			if nz < 0 || nz > logic.MapDepth {
				continue
			}
			neighbours = append(neighbours, geom.Vec3{
				X: float64(x),
				Y: float64(g.Level.HeightMap[z][x].Height),
				Z: float64(z),
			})
		}
		//}
	}

	return neighbours
}

const (
	JumpVelocityIncr     = 0.02
	FallVelocityTerminal = -0.5
	FallVelocityIncr     = -0.025
)

func (g *Game) ResolveCollisions(agent agents.Agent, intent geom.Vec3) {
	// Agent
	pos := agent.GetPosition()
	fx, fz := pos.X, pos.Z
	ix, iz := int(math.Floor(fx)), int(math.Floor(fz))
	current := g.Level.At(ix, iz)
	// Handle Y movement
	switch {
	case agent.GetJumpingTicks() > 0:
		agent.SetYVelocity(agent.GetYVelocity() + JumpVelocityIncr)
	case pos.Y-0.5 > float64(current.Height):
		v := agent.GetYVelocity()
		v += FallVelocityIncr
		if v < FallVelocityTerminal {
			v = FallVelocityTerminal
		}
		agent.SetYVelocity(v)
		agent.SetGrounded(false)
	}
	pos.Y += agent.GetYVelocity()
	if pos.Y-0.5 < float64(current.Height) && agent.GetJumpingTicks() == 0 {
		pos.Y = float64(current.Height) + 0.5
		agent.SetYVelocity(0)
		agent.SetGrounded(true)
	}
	// Handle X,Z movement
	bz := int(math.Floor(fz + intent.Z))
	if bz != iz {
		c := g.Level.At(ix, bz)
		off := 0.
		if intent.Z > 0 {
			off = 0.99
		}
		if float64(c.Height) > pos.Y-0.5 {
			pos.Z = math.Floor(fz) + off
		} else {
			pos.Z += intent.Z
		}
	} else {
		pos.Z += intent.Z
	}

	bx := int(math.Floor(fx + intent.X))
	if bx != ix {
		c := g.Level.At(bx, iz)
		off := 0.
		if intent.X > 0 {
			off = 0.99
		}
		if float64(c.Height) > pos.Y-0.5 {
			pos.X = math.Floor(fx) + off
		} else {
			pos.X += intent.X
		}
	} else {
		pos.X += intent.X
	}
	// Ensure no oob
	if pos.X > logic.MapWidth-0.5 {
		pos.X = logic.MapWidth - 0.5
	} else if pos.X-0.5 < 0 {
		pos.X = 0.5
	}
	if pos.Z > logic.MapDepth-0.5 {
		pos.Z = logic.MapDepth - 0.5
	} else if pos.Z-0.5 < 0 {
		pos.Z = 0.5
	}

	agent.SetPosition(pos)
}
