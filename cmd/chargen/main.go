package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type App struct {
	ticks uint64

	side    float32
	front   float32
	jumping float32
}

func New() *App {
	return &App{
		side:    1,
		front:   1,
		jumping: 0,
	}
}

func (g *App) Update() error {
	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("quit")
	}

	// Side
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.side = -1
	} else {
		g.side = 1
	}
	// Front / Back
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) {
		g.front = -1
	} else {
		g.front = 1
	}
	// Jumping
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		g.jumping = 1
	} else {
		g.jumping = 0
	}

	g.ticks++

	return nil
}

const (
	EyesModifiersNone float32 = iota
	EyesModifiersJoyful
	EyesModifiersAngry
)

func (g *App) Draw(screen *ebiten.Image) {
	const size = 512

	// Draw material
	screen.Fill(color.RGBA{128, 0, 255, 255})
	vertices, indices := graphics.AppendQuadVerticesIndices(
		nil, nil, 0, &graphics.QuadOpts{
			DstX:      logic.ScreenWidth/2 - size/2,
			DstY:      logic.ScreenHeight/2 - size/2,
			DstWidth:  size,
			DstHeight: size,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         g.front,
			G:         1,
		},
	)
	screen.DrawTrianglesShader(vertices, indices, assets.CharShader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Time":       float64(g.ticks) / logic.TPS,
			"Idle":       0.,
			"Walking":    1.,
			"Running":    0.,
			"Jumping":    g.jumping,
			"AgentColor": []float32{0.5, 0.25, 1},
			"Eyes":       EyesModifiersAngry,
		},
	})
	// Debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"FPS: %.2f", ebiten.ActualFPS(),
	))
	// Commands
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(false)
	ebiten.SetMaxTPS(logic.TPS)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	if err := ebiten.RunGame(New()); err != nil {
		fmt.Println("err:", err)
	}
}
