package ui

import (
	"image/color"
	"strconv"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type LevelView struct {
	ticks        uint64
	active       bool
	started      bool
	hoveredLevel int
	pickedLevel  int
	lx, ly       int
	mx, my       int
}

func NewLevelView() *LevelView {
	x, y := ebiten.CursorPosition()
	return &LevelView{
		ticks:        0,
		started:      false,
		active:       true,
		hoveredLevel: -1,
		pickedLevel:  0,
		lx:           x,
		ly:           y,
		mx:           logic.ScreenWidth / 2,
		my:           logic.ScreenHeight / 2,
	}
}

func (lv *LevelView) Active() bool {
	return lv.active
}

func (lv *LevelView) Activate() {
	lv.ticks = 0
	lv.started = false
	lv.active = true
	lv.hoveredLevel = -1
	lv.pickedLevel = 0
	lv.lx, lv.ly = ebiten.CursorPosition()
	lv.mx = logic.ScreenWidth / 2
	lv.my = logic.ScreenWidth / 2
}

func (lv *LevelView) Deactivate() {
	lv.active = false
}

func (lv *LevelView) LevelStarted() (int, bool) {
	return lv.pickedLevel, lv.started
}

func (lv *LevelView) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		lv.active = false
	}

	cx, cy := ebiten.CursorPosition()
	dx, dy := cx-lv.lx, cy-lv.ly
	lv.mx, lv.my = lv.mx+dx, lv.my+dy
	lv.lx, lv.ly = cx, cy
	if lv.mx < 0 {
		lv.mx = 0
	}
	if lv.mx > logic.ScreenWidth {
		lv.mx = logic.ScreenWidth
	}
	if lv.my < 0 {
		lv.my = 0
	}
	if lv.my > logic.ScreenHeight {
		lv.my = logic.ScreenHeight
	}

	lv.ticks++
}

func (lv *LevelView) Draw(screen *ebiten.Image) {
	const (
		boxSize   = logic.ScreenHeight - 32*2
		boxBorder = 10
	)
	// View
	graphics.DrawRectBorder(
		screen,
		logic.CenterX-boxSize/2, logic.CenterY-boxSize/2,
		boxSize, boxSize, boxBorder,
		0.5, 0.5, 0.5, 1,
	)
	vertices, indices := graphics.AppendQuadVerticesIndices(
		nil, nil, 0, &graphics.QuadOpts{
			DstX:      logic.CenterX - boxSize/2 + boxBorder,
			DstY:      logic.CenterY - boxSize/2 + boxBorder,
			DstWidth:  boxSize - boxBorder*2,
			DstHeight: boxSize - boxBorder*2,
			R:         0.1,
			G:         0.1,
			B:         0.1,
			A:         0.8,
		},
	)
	screen.DrawTriangles(vertices, indices, graphics.BrushImage, &ebiten.DrawTrianglesOptions{
		Blend: ebiten.BlendSourceOver,
	})
	// Agents
	const agentSize = 512
	// Player
	vertices, indices = graphics.AppendQuadVerticesIndices(
		vertices[:0], indices[:0], 0, &graphics.QuadOpts{
			DstX:      logic.ScreenWidth/2 - agentSize + 48,
			DstY:      logic.ScreenHeight/2 - agentSize/2 - 128,
			DstWidth:  agentSize,
			DstHeight: agentSize,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         1,
			G:         1,
		},
	)
	screen.DrawTrianglesShader(vertices, indices, assets.CharShader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Time":       float64(lv.ticks) / logic.TPS,
			"Idle":       0.,
			"Walking":    1.,
			"Running":    0.,
			"Jumping":    0.,
			"AgentColor": []float32{1, 1, 1},
			"Eyes":       0.,
		},
	})
	// Agent
	lvl := lv.pickedLevel
	if lv.hoveredLevel != -1 {
		lvl = lv.hoveredLevel
	}
	vertices, indices = graphics.AppendQuadVerticesIndices(
		vertices[:0], indices[:0], 0, &graphics.QuadOpts{
			DstX:      logic.ScreenWidth/2 - 48,
			DstY:      logic.ScreenHeight/2 - agentSize/2 - 128,
			DstWidth:  agentSize,
			DstHeight: agentSize,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         1,
			G:         -1,
		},
	)
	screen.DrawTrianglesShader(vertices, indices, assets.CharShader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Time":       float64(lv.ticks+logic.TPS*2) / logic.TPS,
			"Idle":       0.,
			"Walking":    0.,
			"Running":    1.,
			"Jumping":    0.,
			"AgentColor": graphics.AgentColorsByLevel[lvl],
			"Eyes":       graphics.AgentEyesByLevel[lvl],
		},
	})
	// Buttons
	const (
		lvlButtonSize     = 96
		borderWidth       = 5
		lvlButtonOffY     = logic.ScreenHeight/2 + 256
		lvlButtonSpacingX = 32
	)
	lv.hoveredLevel = -1
	click := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	vertices, indices = vertices[:0], indices[:0]
	for i := 0; i < 4; i++ {
		x := float32(logic.ScreenWidth/2) + float32((lvlButtonSize+lvlButtonSpacingX)*i)
		x -= (lvlButtonSize*2 + lvlButtonSpacingX*1.5)
		y := float32(lvlButtonOffY)
		vertices, indices = graphics.AppendQuadVerticesIndices(
			vertices[:0], indices[:0], 0, &graphics.QuadOpts{
				DstX:      x,
				DstY:      y,
				DstWidth:  lvlButtonSize,
				DstHeight: lvlButtonSize,
				SrcWidth:  1,
				SrcHeight: 1,
				R:         (float32(i) / 4) * 0.5,
				G:         0,
				B:         1 - float32(i)/4,
				A:         1,
			},
		)
		screen.DrawTriangles(vertices, indices, graphics.BrushImage, nil)
		// Hovering
		if i < 2 {
			if lv.mx >= int(x) && lv.mx <= int(x+lvlButtonSize) && lv.my >= int(y) && lv.my <= int(y+lvlButtonSize) {
				lv.hoveredLevel = i
			}
			// Draw border
			if i == lv.hoveredLevel || i == lv.pickedLevel {
				graphics.DrawRectBorder(
					screen,
					x-borderWidth, y-borderWidth,
					lvlButtonSize+borderWidth*2, lvlButtonSize+borderWidth*2,
					borderWidth, 1, 1, 1, 1,
				)
			}
		}
		// Level text
		str := strconv.Itoa(i + 1)
		if i > 1 {
			str += " :("
		}
		rect := text.BoundString(assets.GameInfoFontFace, str)
		text.Draw(
			screen, str, assets.GameInfoFontFace,
			int(x)+lvlButtonSize/2-rect.Dx()/2, int(y)+lvlButtonSize/2+rect.Dy()/2,
			color.White,
		)
	}
	// Start
	const (
		startButtonWidth  = 512
		startButtonHeight = 96
		startButtonOffY   = logic.CenterY + 376
	)
	x := logic.ScreenWidth/2 - float32(startButtonWidth)/2
	y := float32(startButtonOffY)
	vertices, indices = graphics.AppendQuadVerticesIndices(
		vertices[:0], indices[:0], 0, &graphics.QuadOpts{
			DstX:      x,
			DstY:      y,
			DstWidth:  startButtonWidth,
			DstHeight: startButtonHeight,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         0,
			G:         0.5,
			B:         0.125,
			A:         1,
		},
	)
	screen.DrawTriangles(vertices, indices, graphics.BrushImage, nil)
	startStr := "Start"
	rect := text.BoundString(assets.GameInfoFontFace, startStr)
	text.Draw(
		screen, startStr, assets.GameInfoFontFace,
		int(logic.CenterX-rect.Dx()/2),
		int(startButtonOffY+startButtonHeight/2+rect.Dy()/2),
		color.White,
	)
	startHovered := lv.mx >= int(x) && lv.mx <= int(x+startButtonWidth) && lv.my >= int(y) && lv.my <= int(y+startButtonHeight)
	if startHovered {
		graphics.DrawRectBorder(
			screen,
			x-borderWidth, y-borderWidth,
			startButtonWidth+borderWidth*2, startButtonHeight+borderWidth*2,
			borderWidth, 1, 1, 1, 1,
		)
	}
	// Menu title
	title := "New game"
	rect = text.BoundString(assets.MenuTitleFontFace, title)
	offx := logic.CenterX - float32(rect.Dx())/2
	text.Draw(screen, title, assets.MenuTitleFontFace, int(offx), int(108), color.White)
	// Draw cursor
	cursorImg := assets.CursorImage
	if lv.hoveredLevel != -1 {
		cursorImg = assets.CursorHoverImage
	}
	if click {
		cursorImg = assets.CursorClickImage
		if lv.hoveredLevel != -1 {
			lv.pickedLevel = lv.hoveredLevel
		}
		if startHovered {
			lv.started = true
		}
	}
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(lv.mx)-16, float64(lv.my)-16)
	screen.DrawImage(cursorImg, opts)
}
