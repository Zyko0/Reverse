package utils

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var inputGiven bool

func EnsureCursorCaptured() bool {
	if ebiten.CursorMode() == ebiten.CursorModeCaptured {
		return true
	}

	inputGiven = inputGiven || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if runtime.GOOS != "js" || inputGiven {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		inputGiven = false
	}

	return false
}
