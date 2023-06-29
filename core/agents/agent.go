package agents

import (
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

const (
	PlayerVisibilityRadius = 50

	AgentDefaultMS       = 0.1
	AgentRunMSMultiplier = 3
	AgentRunMS           = AgentDefaultMS * AgentRunMSMultiplier
	JumpingTicks         = logic.TPS / 6
	HeardForTicks        = logic.TPS * 2
	CheatVisionTicks     = logic.TPS * 15
)

var (
	jumpDistanceAwarenessByAgent = []int{
		0: 1,
		1: 2,
		2: 3,
		3: 3,
	}
)

type State = byte

const (
	StateIdle State = iota
	StateWalking
	StateRunning
	StateJumping
	StateWalkJumping
	StateFalling
)

type Agent interface {
	Update(env *Env)

	GetState() State
	GetAngle() float64
	GetPosition() geom.Vec3
	SetPosition(position geom.Vec3)
	GetIntent() geom.Vec3
	GetGrounded() bool
	SetGrounded(grounded bool)
	GetYVelocity() float64
	SetYVelocity(v float64)
	GetJumpingTicks() uint64
	GetHeard() bool
	SetHeard()

	HasAbility(ability Ability) bool
	Cooldown(ability Ability) float64
}

func NewAgentByLevel(lvl int) Agent {
	switch lvl {
	case 0:
		return NewAgent0()
	default:
		return nil
	}
}

type base struct {
	ticks uint64

	Angle    float64
	Position geom.Vec3
	Intent   geom.Vec3

	Grounded      bool
	YVelocity     float64
	JumpingTicks  uint64
	Running       bool
	HeardForTicks uint64
}

func (b *base) update() {
	b.Running = false
	if b.JumpingTicks > 0 {
		b.JumpingTicks--
	}
	if b.HeardForTicks > 0 {
		b.HeardForTicks--
	}
	b.Intent.X, b.Intent.Y, b.Intent.Z = 0, 0, 0
	b.ticks++
}

func (b *base) GetState() State {
	switch {
	case b.JumpingTicks > 0:
		return StateJumping
	case !b.Grounded:
		return StateFalling
	case !b.Intent.Zero():
		if b.Running {
			return StateRunning
		}
		return StateWalking
	default:
		return StateIdle
	}
}

func (b *base) GetAngle() float64 {
	return b.Angle
}

func (b *base) GetPosition() geom.Vec3 {
	return b.Position
}

func (b *base) SetPosition(position geom.Vec3) {
	b.Position = position
}

func (b *base) GetIntent() geom.Vec3 {
	return b.Intent
}

func (b *base) GetGrounded() bool {
	return b.Grounded
}

func (b *base) SetGrounded(grounded bool) {
	b.Grounded = grounded
}

func (b *base) GetYVelocity() float64 {
	return b.YVelocity
}

func (b *base) SetYVelocity(v float64) {
	b.YVelocity = v
}

func (b *base) GetJumpingTicks() uint64 {
	return b.JumpingTicks
}

func (b *base) GetHeard() bool {
	return b.HeardForTicks > 0 || b.Running
}

func (b *base) SetHeard() {
	b.HeardForTicks = HeardForTicks
}
