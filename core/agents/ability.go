package agents

import "github.com/Zyko0/Reverse/logic"

type Ability byte

const (
	AbilityLightsOff Ability = iota
	AbilityTasing
	AbilityScanning
)

// Properties
const (
	TasingTicks    = logic.TPS
	TasingCooldown = logic.TPS * 5
	TasingRadius   = 10

	ScanningTicks    = logic.TPS * 5
	ScanningCoolDown = logic.TPS * 20
)
