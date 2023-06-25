package agents

import (
	"github.com/Zyko0/Reverse/pkg/level"
)

type Agent0 struct {
	base
}

func NewAgent0() *Agent0 {
	return &Agent0{
		base{
			Position: level.StartAgentPosition,
			Grounded: true,
		},
	}
}

func (a0 *Agent0) Update(env *Env) {

}

func (a0 *Agent0) HasAbility(ability Ability) bool {
	return false
}
