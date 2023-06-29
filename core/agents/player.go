package agents

import (
	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	base

	tasingTicks   uint64
	tasingCD      uint64
	scanningTicks uint64
	scanningCD    uint64
}

func NewPlayer() *Player {
	return &Player{
		base: base{
			Grounded: true,
			Position: level.StartPlayerPosition,
		},

		scanningCD: ScanningCoolDown,
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
	if p.scanningCD > 0 {
		p.scanningCD--
	}
	if p.scanningTicks > 0 {
		p.scanningTicks--
		p.scanningCD = ScanningCoolDown
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
	if ebiten.IsKeyPressed(ebiten.KeyQ) && p.scanningTicks == 0 && p.scanningCD == 0 {
		p.scanningTicks = ScanningTicks
		// SFX
		assets.PlayScan()
	}

	p.Running = ebiten.IsKeyPressed(ebiten.KeyShift)
	ms := AgentDefaultMS
	if p.Running {
		ms = AgentRunMS
		// SFX
		assets.PlayFootsteps()
	} else {
		assets.StopFootsteps()
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
	case AbilityScanning:
		return p.scanningTicks > 0
	default:
		return false
	}
}

func (p *Player) Cooldown(ability Ability) float64 {
	switch ability {
	case AbilityTasing:
		return float64(p.tasingCD) / float64(TasingCooldown)
	case AbilityScanning:
		return float64(p.scanningCD) / float64(ScanningCoolDown)
	default:
		return 0
	}
}
