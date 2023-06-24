package core

import (
	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Game struct {
	Level  *level.Map
	Player *agents.Player
	Agent  agents.Agent
	Camera *Camera
}

func NewGame() *Game {
	position := geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 4.5,   // TODO: default height
		Z: 244.5, //MapDepth / 2,
	}

	return &Game{
		Level:  assets.Level0,
		Player: agents.NewPlayer(position),
		Camera: newCamera(), // TODO: set initial direction
	}
}

func (g *Game) Update() {
	g.Player.Update(nil)
	// Camera
	// TODO: process new player position with intents and new cam direction just here
	//g.Player.Position = g.Player.Position.Add(g.Camera.Direction.Mul(g.Player.Intent))
	pxz := geom.Vec2{
		X: g.Player.Intent.X,
		Y: g.Player.Intent.Z,
	}.Rotate(g.Camera.HAngle)
	g.ResolveCollisions(g.Player, geom.Vec3{
		X: pxz.X,
		Y: g.Player.Intent.Y,
		Z: pxz.Y,
	})
	// TODO: ^
	g.Camera.Update() // Only ticks for now
}
