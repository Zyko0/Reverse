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
	ticks       uint64
	lvl         int
	reached     bool
	status      GameStatus
	lastKnownAt geom.Vec3

	Level  *level.HMap
	Player *agents.Player
	Agent  agents.Agent
	Camera *Camera
}

func NewGame(lvl int) *Game {
	return &Game{
		ticks:       0,
		lvl:         lvl,
		reached:     false,
		status:      GameStatusDefeat,
		lastKnownAt: level.StartPlayerPosition,

		Level:  assets.Level0,
		Player: agents.NewPlayer(),
		Agent:  agents.NewAgentByLevel(lvl), // TODO: change depending on the level
		Camera: newCamera(),
	}
}

func (g *Game) Update() {
	// Camera
	g.Camera.Update()
	// Gameover
	if g.IsOver() {
		return
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
	// Record last known player position
	heardSeen := g.Player.HasAbility(agents.AbilityScanning)
	canSeePlayer := false
	if g.Player.GetHeard() { // if noisy
		g.lastKnownAt = g.Player.GetPosition()
		heardSeen = true
	}
	// Allow player to be seen first 15seconds to ease it for AI
	if g.PlayerSeen() || g.ticks < agents.CheatVisionTicks { // if seen
		g.lastKnownAt = g.Player.GetPosition()
		heardSeen = true
		canSeePlayer = true
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
		Level:             g.lvl,
		Map:               g.Level,
		Goal:              level.GoalPosition,
		LastKnownAt:       g.lastKnownAt,
		LastKnownAtUpdate: heardSeen,
		TimeRemaining:     uint64(g.TimeRemaining()),
		Seen:              g.AgentSeen(),
		CanSeePlayer:      canSeePlayer,
	}
	g.Agent.Update(env)
	intent := g.Agent.GetIntent()
	g.ResolveCollisions(g.Agent, geom.Vec3{
		X: intent.X,
		Y: 0,
		Z: intent.Z,
	})
	// If Agent took the goal point
	if g.Agent.GetPosition().Floor() == level.GoalPosition.Floor() {
		g.reached = true
		g.status = GameStatusDefeat
	}

	g.ticks++
}

func (g *Game) GetLevel() int {
	return g.lvl
}

func (g *Game) PlayerSeen() bool {
	a := g.Agent.GetPosition()
	p := g.Player.GetPosition()
	dist := a.DistanceTo(p)

	return g.Level.CastRay(a, p, dist)
}

func (g *Game) AgentSeen() bool {
	return g.PlayerSeen() && g.Player.GetPosition().DistanceTo(g.Agent.GetPosition()) <= agents.PlayerVisibilityRadius
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
	return g.ticks > level.LevelsTime[g.lvl] || g.reached || g.status == GameStatusVictory
}
