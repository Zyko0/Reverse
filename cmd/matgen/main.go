package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	fragSrc = []byte(
		`
package main

const (
	Pi        = 3.1415925359
	TwoPi = Pi*2
)

var Time float

func colorize(t, seed float) vec3 {
	return vec3(.6 + .6*cos(TwoPi * t + seed*vec3(0,20,10)))
}

func hash(p vec2, seed float) float {
	return fract(sin(dot(p, vec2(12.9898, 4.1414)*seed)) * 43758.5453)
}

func rotate2D(v vec2, a float) vec2 {
	s := sin(a)
	c := cos(a)
	m := mat2(c, -s, s, c)
	return m * v
}

const (
	NoiseHashFloored = 0.5
	NoiseDot = 1.5
	NoiseSinCos = 2.5
)

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	uv := texCoord*2-1
	seed := color.r
	scale := color.g
	noise := color.b

	// Scale
	uv /= scale

	// Value
	var v float
	uv = sin(uv*Pi)+cos(uv*Pi)
	v = hash(floor(uv), 1)
	if noise < NoiseHashFloored {
	} else if noise < NoiseDot {
		v *= dot(uv, uv)
	} else if noise < NoiseSinCos {
		//vv := dot(uv, uv)
		vv := length(uv)-0.25*scale
		vv *= vv
		v = vv
	}

	clr := colorize(v, seed)

	return vec4(clr, 1)
}
`)

	s *ebiten.Shader
)

func init() {
	var err error

	s, err = ebiten.NewShader(fragSrc)
	if err != nil {
		log.Fatal(err)
	}
}

type App struct {
	ticks uint64

	seed float64
	scale float64
	noise int
}

func New() *App {
	return &App{}
}

func (g *App) Save() {

}

func (g *App) Update() error {
	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("quit")
	}
	// Save
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Save()
	}
	// Parameters
	const (
		wheelSens = 0.005
		angleSens = 0.01
	)
	// Angle
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.seed += angleSens
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.seed -= angleSens
	}
	// Scale
	_, yoff := ebiten.Wheel()
	g.scale += yoff * wheelSens
	if g.scale < wheelSens {
		g.scale = 0.0001
	}
	// Noise
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.noise++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.noise--
	}

	g.ticks++

	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	const size = 512

	// Draw material
	vertices, indices := graphics.AppendQuadVerticesIndices(
		nil, nil, 0, &graphics.QuadOpts{
			DstX:      logic.ScreenWidth/2 - size/2,
			DstY:      logic.ScreenHeight/2 - size/2,
			DstWidth:  size,
			DstHeight: size,
			SrcWidth:  1,
			SrcHeight: 1,
			R:         float32(g.seed),
			G:         float32(g.scale),
			B:         float32(g.noise),
			A:         1,
		},
	)
	screen.DrawTrianglesShader(vertices, indices, s, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Time": float64(g.ticks) / logic.TPS,
		},
	})
	// Debug
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Angle: %6f - Scale %.6f - Noise %d",
		g.seed,
		g.scale,
		g.noise,
	))
	// Commands
	ebitenutil.DebugPrintAt(screen, "- Mouse wheel to change the scale", 0, 12)
	ebitenutil.DebugPrintAt(screen, "- Press I/K to change noise method", 0, 24)
	ebitenutil.DebugPrintAt(screen, "- Press O/L to change side material id", 0, 36)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetMaxTPS(logic.TPS)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	if err := ebiten.RunGame(New()); err != nil {
		fmt.Println("err:", err)
	}

}
