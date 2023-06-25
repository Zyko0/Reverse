package core

import (
	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Game struct {
	lvl int

	Level  *level.Map
	Player *agents.Player
	Agent  agents.Agent
	Camera *Camera
}

func NewGame() *Game {
	return &Game{
		Level:  assets.Level0,
		Player: agents.NewPlayer(),
		Agent:  agents.NewAgent0(),
		Camera: newCamera(), // TODO: set initial direction
	}
}

func (g *Game) Update() {
	env := &agents.Env{
		Map:       g.Level,
		LastHeard: g.Player.Position,
	}
	// Player
	g.Player.Update(nil)
	pxz := geom.Vec2{
		X: g.Player.Intent.X,
		Y: g.Player.Intent.Z,
	}.Rotate(g.Camera.HAngle)
	g.ResolveCollisions(g.Player, geom.Vec3{
		X: pxz.X,
		Y: 0,
		Z: pxz.Y,
	})
	// Agent
	g.Agent.Update(env)
	intent := g.Agent.GetIntent()
	g.ResolveCollisions(g.Agent, geom.Vec3{
		X: intent.X,
		Y: 0,
		Z: intent.Z,
	})
	// Camera
	g.Camera.Update()
}

func (g *Game) GetLevel() int {
	return g.lvl
}
