package ui

import (
	"image/color"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type PauseView struct {
	Active bool
}

func NewPauseView() *PauseView {
	return &PauseView{}
}

func (pv *PauseView) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		pv.Active = false
	}
}

func (pv *PauseView) Draw(screen *ebiten.Image) {
	const (
		boxsize   = logic.ScreenHeight - 256*2
		boxborder = 10
	)
	graphics.DrawRectBorder(
		screen,
		logic.CenterX-boxsize/2, logic.CenterY-boxsize/2,
		boxsize, boxsize, boxborder,
		0, 0, 0, 1,
	)
	vertices, indices := graphics.AppendQuadVerticesIndices(
		nil, nil, 0, &graphics.QuadOpts{
			DstX:      logic.CenterX - boxsize/2 + boxborder,
			DstY:      logic.CenterY - boxsize/2 + boxborder,
			DstWidth:  boxsize - boxborder*2,
			DstHeight: boxsize - boxborder*2,
			R:         0.05,
			G:         0.05,
			B:         0.05,
			A:         0.8,
		},
	)
	screen.DrawTriangles(vertices, indices, graphics.BrushImage, &ebiten.DrawTrianglesOptions{
		Blend: ebiten.BlendSourceOver,
	})
	// Menu title
	rect := text.BoundString(assets.MenuTitleFontFace, "Pause")
	offx := logic.CenterX - float32(rect.Dx())/2
	text.Draw(screen, "Pause", assets.MenuTitleFontFace, int(offx), int(325), color.White)
	// Controls
	controlsStr := "Controls:\n"
	controlsStr += "- WASD: Movement\n"
	controlsStr += "- Arrow keys / Mouse: Camera+zoom\n"
	controlsStr += "- Shift: Run\n"
	controlsStr += "- Space: Jump\n"
	controlsStr += "- Q: Scanner ability\n"
	controlsStr += "- E: Taser ability\n"
	controlsStr += "- Backspace: Restart current level\n"
	controlsStr += "- Tab: Toggle pause menu\n"
	controlsStr += "- Escape: Go back to main menu\n"
	rect = text.BoundString(assets.DescriptionFontFace, controlsStr)
	offx = logic.CenterX - boxsize/2 + 36
	text.Draw(screen, controlsStr, assets.DescriptionFontFace, int(offx), int(400), color.White)
}
