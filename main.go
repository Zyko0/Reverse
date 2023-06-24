package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Zyko0/Reverse/core"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	renderer *graphics.Renderer
	game     *core.Game
}

func New() *Game {
	return &Game{
		renderer: graphics.NewRenderer(),
		game:     core.NewGame(),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("quit")
	}

	// Update game
	g.game.Update()
	// Update renderer
	g.renderer.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render game
	g.renderer.Draw(screen, &graphics.State{
		Level:  g.game.Level,
		Camera: g.game.Camera,
		Player: g.game.Player,
	})
	// Debug
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f - FPS %.02f - PPos (%v) CamPos(%v) CamDir(%v) Intent(%v) Hangle %.4f - Block(%d,%d)",
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			g.game.Player.Position,
			g.game.Camera.Position, g.game.Camera.Direction, g.game.Player.Intent,
			g.game.Camera.HAngle,
			int(g.game.Player.Position.X), int(g.game.Player.Position.Z),
		),
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	/*f, err := os.Create("beat.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println("couldn't profile:", err)
		return
	}
	defer pprof.StopCPUProfile()*/

	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(false) // TODO: remove
	ebiten.SetWindowSize(logic.ScreenWidth, logic.ScreenHeight)
	ebiten.SetMaxTPS(logic.TPS)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	// (broken) go get github.com/hajimehoshi/ebiten/v2@1c09ec5e44727a0c38b605552d93e4d470a128ab
	// (stable) v2.5.0-alpha.12.0.20230228174701-7c0fbce0cfd8
	if err := ebiten.RunGame(New()); err != nil {
		// TODO: gracefull
	}
}
