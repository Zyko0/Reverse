package graphics

import (
	"image/color"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	abilityOnCDColor      = color.RGBA{128, 128, 128, 255}
	abilityAvailableColor = color.RGBA{255, 255, 255, 255}
)

func (r *Renderer) renderAbilities(s *State) {
	const (
		abilityWidth  = 256 * 0.75
		abilityHeight = abilityWidth / 2
		offx          = logic.ScreenWidth - abilityWidth*2
		offy          = abilityWidth * 2
	)

	// Q - Scan ability
	vertices, indices := AppendQuadVerticesIndices(
		nil, nil, 0, &QuadOpts{
			DstX:      offx,
			DstY:      offy,
			DstWidth:  abilityWidth,
			DstHeight: abilityHeight,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0.15,
			G:         0.15,
			B:         0.5,
			A:         0.75,
		},
	)
	qcd := float32(s.Player.Cooldown(agents.AbilityScanning))
	vertices, indices = AppendQuadVerticesIndices(
		vertices, indices, 1, &QuadOpts{
			DstX:      offx,
			DstY:      offy,
			DstWidth:  abilityWidth * qcd,
			DstHeight: abilityHeight,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0.1,
			G:         0.1,
			B:         0.1,
			A:         0.5,
		},
	)
	// E - Taser ability
	vertices, indices = AppendQuadVerticesIndices(
		vertices, indices, 2, &QuadOpts{
			DstX:      offx + abilityWidth,
			DstY:      offy,
			DstWidth:  abilityWidth,
			DstHeight: abilityHeight,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0.5,
			G:         0.05,
			B:         0.0,
			A:         0.75,
		},
	)
	ecd := float32(s.Player.Cooldown(agents.AbilityTasing))
	vertices, indices = AppendQuadVerticesIndices(
		vertices, indices, 3, &QuadOpts{
			DstX:      offx + abilityWidth,
			DstY:      offy,
			DstWidth:  abilityWidth * ecd,
			DstHeight: abilityHeight,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0.1,
			G:         0.1,
			B:         0.1,
			A:         0.5,
		},
	)
	r.offscreen.DrawTriangles(vertices, indices, BrushImage, nil)
	// Text scanner
	clr := abilityOnCDColor
	if qcd == 0 {
		clr = abilityAvailableColor
	}
	qstr := "Scan (Q)"
	rect := text.BoundString(assets.DescriptionFontFace, qstr)
	x := offx + abilityWidth/2 - float32(rect.Dx())/2
	y := offy + abilityHeight/2 + float32(rect.Dy())/2
	text.Draw(r.offscreen, qstr, assets.DescriptionFontFace, int(x), int(y), clr)
	// Text taser
	clr = abilityOnCDColor
	if ecd == 0 {
		clr = abilityAvailableColor
	}
	estr := "Tase (E)"
	rect = text.BoundString(assets.DescriptionFontFace, estr)
	x = offx + abilityWidth + abilityWidth/2 - float32(rect.Dx())/2
	y = offy + abilityHeight/2 + float32(rect.Dy())/2
	text.Draw(r.offscreen, estr, assets.DescriptionFontFace, int(x), int(y), clr)
}
