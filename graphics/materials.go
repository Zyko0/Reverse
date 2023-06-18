package graphics

import (
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SheetWidth, SheetHeight = logic.MapWidth, logic.MapDepth
)

var (
	SheetImage = ebiten.NewImage(logic.MapWidth, logic.MapDepth)
)
