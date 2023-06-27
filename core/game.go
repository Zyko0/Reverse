package core

import (
	"time"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
)

type GameStatus = byte

const (
	GameStatusDefeat GameStatus = iota
	GameStatusVictory
)

type Game struct {
	ticks     uint64
	lvl       int
	status    GameStatus
	lastHeard geom.Vec3

	Level  *level.HMap
	Player *agents.Player
	Agent  agents.Agent
	Camera *Camera
}

func NewGame() *Game {
	return &Game{
		lastHeard: level.StartPlayerPosition,

		Level:  assets.Level0,
		Player: agents.NewPlayer(),
		Agent:  agents.NewAgent0(),
		Camera: newCamera(),
	}
}

func (g *Game) Update() {
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
	// Record last heard if player is noisy
	if g.Player.GetHeard() {
		g.lastHeard = g.Player.GetPosition()
	}
	// Check if tasing and if it hits
	if g.Player.HasAbility(agents.AbilityTasing) {
		tased := g.Level.CastRay(g.Player.Position, g.Agent.GetPosition(), agents.TasingRadius)
		if tased {
			g.status = GameStatusVictory
		}
		// TODO: play success sfx
		// TODO: play failing sfx
	}

	// Agent
	env := &agents.Env{
		Map:           g.Level,
		Goal:          level.GoalPosition,
		LastHeard:     g.lastHeard,
		TimeRemaining: uint64(g.TimeRemaining()),
	}
	g.Agent.Update(env)
	intent := g.Agent.GetIntent()
	g.ResolveCollisions(g.Agent, geom.Vec3{
		X: intent.X,
		Y: 0,
		Z: intent.Z,
	})
	// Camera
	g.Camera.Update()

	g.ticks++
}

func (g *Game) GetLevel() int {
	return g.lvl
}

func (g *Game) Status() GameStatus {
	return g.status
}

func (g *Game) TimeRemaining() time.Duration {
	if g.ticks > level.LevelsTime[g.lvl] {
		return 0
	}

	return (time.Duration(level.LevelsTime[g.lvl]-g.ticks) * time.Second) / logic.TPS
}

func (g *Game) IsOver() bool {
	return g.ticks > level.LevelsTime[g.lvl] || g.status == GameStatusVictory
}
