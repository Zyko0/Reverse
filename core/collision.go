package core

import (
	"math"

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

func (g *Game) ResolveCollisions(intent geom.Vec3) {
	// Player
	fx, fy, fz := g.Player.Position.X, g.Player.Position.Y, g.Player.Position.Z
	ix, _, iz := int(fx), int(fy), int(fz)
	current := g.Level.HeightMap[iz][ix]
	// Handle Y movement
	switch {
	case g.Player.JumpingTicks > 0:
		g.Player.YVelocity += JumpVelocityIncr
	case g.Player.Position.Y-0.5 > float64(current.Height):
		g.Player.YVelocity += FallVelocityIncr
		if g.Player.YVelocity < FallVelocityTerminal {
			g.Player.YVelocity = FallVelocityTerminal
		}
	}
	g.Player.Position.Y += g.Player.YVelocity
	if g.Player.Position.Y-0.5 < float64(current.Height) && g.Player.JumpingTicks == 0 {
		g.Player.Position.Y = float64(current.Height) + 0.5
		g.Player.YVelocity = 0
	}
	// Handle X,Z movement
	pos := g.Player.Position
	bx := int(fx + intent.X)
	if bx != ix {
		c := g.Level.HeightMap[iz][bx]
		off := 0.
		if intent.X > 0 {
			off = 0.99
		}
		if float64(c.Height) > g.Player.Position.Y-0.5 {
			pos.X = math.Floor(fx) + off
		} else {
			pos.X += intent.X
		}
	} else {
		pos.X += intent.X
	}

	bz := int(fz + intent.Z)
	if bz != iz {
		c := g.Level.HeightMap[bz][ix]
		off := 0.
		if intent.Z > 0 {
			off = 0.99
		}
		if float64(c.Height) > g.Player.Position.Y-0.5 {
			pos.Z = math.Floor(fz) + off
		} else {
			pos.Z += intent.Z
		}
	} else {
		pos.Z += intent.Z
	}

	g.Player.Position = pos
}
