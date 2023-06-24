package agents

import (
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	base
}

func NewPlayer() *Player {
	return &Player{
		base{
			Grounded: true,
			Position: level.StartAgentPosition,
		},
	}
}

func (p *Player) Update(env *Env) {
	p.base.update()

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

	p.Running = ebiten.IsKeyPressed(ebiten.KeyShift)
	ms := AgentDefaultMS // TODO: make this better
	if p.Running {
		ms = AgentRunMS
	}
	if !p.Intent.Zero() {
		p.Intent = p.Intent.Normalize()
		p.Intent = p.Intent.MulN(ms)
	}
}
