package agents

import "github.com/Zyko0/Reverse/logic"

type Ability byte

const (
	AbilityLightsOff Ability = iota
	AbilityTasing
	AbilityScouting
)

// Properties
const (
	TasingTicks    = logic.TPS
	TasingCooldown = logic.TPS * 5
	TasingRadius   = 10

	ScoutingTicks    = logic.TPS * 5
	ScoutingCoolDown = logic.TPS * 20
)
