package core

import (
	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core/agent"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type Game struct {
	Level  *level.Map
	Player *Player
	Agent  agent.Agent
	Camera *Camera
}

func NewGame() *Game {
	position := geom.Vec3{
		X: logic.MapWidth / 2,
		Y: 4.5, // TODO: default height
		Z: 244.5, //MapDepth / 2,
	}

	return &Game{
		Level:  assets.Level0,
		Player: newPlayer(position),
		Camera: newCamera(position), // TODO: set initial direction
	}
}

func (g *Game) Update() {
	g.Player.Update()
	// Camera
	g.Camera.UpdateDirection(g.Player.Position)
	// TODO: process new player position with intents and new cam direction just here
	//g.Player.Position = g.Player.Position.Add(g.Camera.Direction.Mul(g.Player.Intent))
	pxz := geom.Vec2{
		X: g.Player.Intent.X,
		Y: g.Player.Intent.Z,
	}.Rotate(g.Camera.HAngle)
	g.ResolveCollisions(geom.Vec3{
		X: pxz.X,
		Y: g.Player.Intent.Y,
		Z: pxz.Y,
	})
	// TODO: ^
	g.Camera.UpdatePosition(g.Player.Position)
	g.Camera.Update() // Only ticks for now
}
