package agents

import (
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	base

	tasingTicks uint64
	tasingCD    uint64
}

func NewPlayer() *Player {
	return &Player{
		base: base{
			Grounded: true,
			Position: level.StartPlayerPosition,
		},
	}
}

func (p *Player) Update(env *Env) {
	p.base.update()

	if p.tasingCD > 0 {
		p.tasingCD--
	}
	if p.tasingTicks > 0 {
		p.tasingTicks--
		if p.tasingTicks == 0 {
			p.tasingCD = TasingCooldown
		}
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
	return ability == AbilityTasing && p.tasingTicks > 0
}
