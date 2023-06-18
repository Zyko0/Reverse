package core

import (
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PlayerRadius         = 0.5
	PlayerDefaultMS      = 0.2
	FallVelocityTerminal = -0.5
	FallVelocityIncr     = -0.025
	JumpingTicks         = logic.TPS * 0.2
	JumpVelocityIncr     = 0.05
)

type Player struct {
	Position geom.Vec3
	Intent   geom.Vec3

	YVelocity    float64
	JumpingTicks uint64
}

func newPlayer(position geom.Vec3) *Player {
	return &Player{
		Position: position,
	}
}

func (p *Player) Update() {
	if p.JumpingTicks > 0 {
		p.JumpingTicks--
	}

	p.Intent.X, p.Intent.Y, p.Intent.Z = 0, 0, 0
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

	}
	if !p.Intent.Zero() {
		p.Intent = p.Intent.Normalize()
		p.Intent = p.Intent.MulN(PlayerDefaultMS)
	}

	// TODO: handle jump / fall
}
