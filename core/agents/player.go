package agents

import (
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	base

	tasingTicks   uint64
	tasingCD      uint64
	scoutingTicks uint64
	scoutingCD    uint64
}

func NewPlayer() *Player {
	return &Player{
		base: base{
			Grounded: true,
			Position: level.StartPlayerPosition,
		},

		scoutingCD: ScoutingCoolDown,
	}
}

func (p *Player) Update(env *Env) {
	p.base.update()

	if p.tasingCD > 0 {
		p.tasingCD--
	}
	if p.tasingTicks > 0 {
		p.tasingTicks--
		p.tasingCD = TasingCooldown
	}
	if p.scoutingCD > 0 {
		p.scoutingCD--
	}
	if p.scoutingTicks > 0 {
		p.scoutingTicks--
		p.scoutingCD = ScoutingCoolDown
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Intent.Z += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Intent.Z -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Intent.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Intent.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.JumpingTicks == 0 && p.YVelocity == 0 {
		p.JumpingTicks = JumpingTicks
		p.Grounded = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) && p.tasingTicks == 0 && p.tasingCD == 0 {
		p.tasingTicks = TasingTicks
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) && p.scoutingTicks == 0 && p.scoutingCD == 0 {
		p.scoutingCD = ScoutingTicks
	}

	p.Running = ebiten.IsKeyPressed(ebiten.KeyShift)
	ms := AgentDefaultMS // TODO: make this better
	if p.Running {
		ms = AgentRunMS
	}
	if !p.Intent.Zero() {
		p.Intent = p.Intent.Normalize()
		p.Intent.X = -p.Intent.X
		p.Intent = p.Intent.MulN(ms)
	}
}

func (p *Player) HasAbility(ability Ability) bool {
	switch ability {
	case AbilityTasing:
		return p.tasingTicks > 0
	case AbilityScouting:
		return p.scoutingTicks > 0
	default:
		return false
	}
}

func (p *Player) Cooldown(ability Ability) float64 {
	switch ability {
	case AbilityTasing:
		return float64(p.tasingCD) / float64(TasingCooldown)
	case AbilityScouting:
		return float64(p.scoutingCD) / float64(ScoutingCoolDown)
	default:
		return 0
	}
}
