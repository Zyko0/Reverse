package agents

import "github.com/Zyko0/Reverse/pkg/geom"

type Agent0 struct {
	base
}

func NewAgent0(position geom.Vec3) *Agent0 {
	return &Agent0{
		base{
			Position: position,
			Grounded: true,
		},
	}
}

func (a0 *Agent0) Update() {
	
}