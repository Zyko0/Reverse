package agents

import "github.com/Zyko0/Reverse/logic"

type Ability byte

const (
	AbilityLightsOff Ability = iota
	AbilityTasing
)

// Properties
const (
	TasingTicks    = logic.TPS
	TasingCooldown = logic.TPS * 5
)
