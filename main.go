package main

import (
	"image/color"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/xfmt"
	"github.com/Zyko0/Reverse/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	levelview *ui.LevelView
	pauseview *ui.PauseView

	game     *core.Game
	renderer *graphics.Renderer
}

func New() *Game {
	return &Game{
		levelview: ui.NewLevelView(),
		pauseview: ui.NewPauseView(),

		game:     core.NewGame(0),
		renderer: graphics.NewRenderer(),
	}
}

func (g *Game) Update() error {
	// Level view
	if g.levelview.Active() {
		if level, started := g.levelview.LevelStarted(); started {
			g.game = core.NewGame(level)
			g.levelview.Deactivate()
		} else {
			g.levelview.Update()
			return nil
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || g.game.IsOver() && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.levelview.Activate()
		return nil
	}
	// Restart level
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.game = core.NewGame(g.game.GetLevel())
	}
	// Pause
	if g.pauseview.Active {
		g.pauseview.Update()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.pauseview.Active = true
		return nil
	}

	// Update game
	g.game.Update()
	// Update renderer
	g.renderer.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Level view
	if g.levelview.Active() {
		g.levelview.Draw(screen)
		return
	}
	// Render game
	g.renderer.Draw(screen, &graphics.State{
		Level:      g.game.GetLevel(),
		Map:        g.game.Level,
		Camera:     g.game.Camera,
		Player:     g.game.Player,
		Agent:      g.game.Agent,
		GameStatus: g.game.Status(),
		PlayerSeen: g.game.PlayerSeen(),
		AgentSeen:  g.game.AgentSeen(),
	})
	// Remaining time
	// TODO: clean below
	timeTxt := xfmt.Duration(g.game.TimeRemaining())
	rect := text.BoundString(assets.GameInfoFontFace, timeTxt)
	text.Draw(screen, timeTxt, assets.GameInfoFontFace,
		logic.ScreenWidth/2-rect.Dx()/2, 36,
		color.White,
	)
	playerStatusTxt := ""
	if g.game.Player.GetHeard() {
		playerStatusTxt = "HEARD"
	}
	if g.game.PlayerSeen() {
		playerStatusTxt = "SEEN"
	}
	if g.game.IsOver() {
		playerStatusTxt = "GAME OVER (VICTORY)"
		if g.game.Status() == core.GameStatusDefeat {
			playerStatusTxt = "PLAYER WINS (DEFEAT)"
		}
	}
	if playerStatusTxt != "" {
		rect := text.BoundString(assets.GameInfoFontFace, playerStatusTxt)
		text.Draw(screen, playerStatusTxt, assets.GameInfoFontFace,
			logic.ScreenWidth/2-rect.Dx()/2, 72,
			color.White,
		)
	}
	// UI
	if g.pauseview.Active {
		g.pauseview.Draw(screen)
	}
	// Debug
	// TODO: remove below
	/*ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f - FPS %.02f - PPos (%v) Intent(%v) Hangle %.4f - Block(%d,%d) - Seen %v",
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			g.game.Player.Position, g.game.Player.Intent,
			g.game.Camera.HAngle,
			int(g.game.Player.Position.X), int(g.game.Player.Position.Z),
			g.game.AgentSeen(),
		),
	)*/
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(logic.ScreenWidth, logic.ScreenHeight)
	ebiten.SetMaxTPS(logic.TPS)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	// (broken) go get github.com/hajimehoshi/ebiten/v2@1c09ec5e44727a0c38b605552d93e4d470a128ab
	// (stable) v2.5.0-alpha.12.0.20230228174701-7c0fbce0cfd8
	if err := ebiten.RunGameWithOptions(New(), &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil {
		// TODO: gracefull
	}
}
