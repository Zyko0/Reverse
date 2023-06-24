package graphics

import (
	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/core"
	"github.com/Zyko0/Reverse/core/agents"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	ticks uint64
	sheet bool
	drawn bool

	heightmap *ebiten.Image
	offscreen *ebiten.Image
}

type State struct {
	Level  *level.Map
	Camera *core.Camera
	Player *agents.Player
	Agent  agents.Agent
}

func NewRenderer() *Renderer {
	return &Renderer{
		ticks: 0,
		sheet: false,
		drawn: false,

		heightmap: ebiten.NewImage(logic.MapWidth, logic.MapDepth),
		offscreen: ebiten.NewImage(logic.ScreenWidth, logic.ScreenHeight),
	}
}

func (r *Renderer) Update() {
	r.drawn = false
	r.ticks++
}

var (
	vertices []ebiten.Vertex
	indices  []uint16
)

var mapdrawn bool // TODO: make it better ofc

func (r *Renderer) Draw(screen *ebiten.Image, state *State) {
	if !r.drawn {
		// Spritesheet update
		if !r.sheet {
			SheetImage.DrawImage(assets.Sheet0Image, nil) // TODO: make it dynamic
			r.sheet = true
		}
		// Map generation
		if !mapdrawn {
			r.heightmap.WritePixels(state.Level.CompileBytes())
			mapdrawn = true
		}
		// Scene rendering
		vertices, indices = AppendQuadVerticesIndices(
			vertices[:0], indices[:0], 0, &QuadOpts{
				DstWidth:  logic.ScreenWidth,
				DstHeight: logic.ScreenHeight,
				SrcWidth:  logic.MapWidth,
				SrcHeight: logic.MapDepth,
			},
		)
		r.offscreen.Clear()
		// Player states
		var pidle, pwalk, prun, pjump float32
		switch state.Player.GetState() {
		case agents.StateIdle:
			pidle = 1
		case agents.StateJumping:
			// Must have prio over walk / run
			pjump = 1 - float32(state.Player.JumpingTicks)/float32(agents.JumpingTicks)
			//pjump *= 2
		case agents.StateWalking:
			pwalk = 1
		case agents.StateRunning:
			prun = 1
		case agents.StateFalling:
			pjump = 1
			prun = 1
		}
		// Agent states

		// Render scene
		r.offscreen.DrawTrianglesShader(vertices, indices, assets.SceneShader, &ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"Time": float64(r.ticks) / logic.TPS,
				"Offset": []float32{
					float32(state.Player.Position.X) - logic.MapWidth/2,
					float32(state.Player.Position.Y),
					float32(state.Player.Position.Z) - logic.MapDepth/2,
				},
				"HorizontalAngle": float32(state.Camera.HAngle),
				"MapSize": []float32{
					logic.MapWidth,
					logic.MapHeight,
					logic.MapDepth,
				},
				"Zoom": float32(state.Camera.Zoom * logic.BaseZoom * float64(logic.MapWidth)),
				// Player
				"PlayerIdle":    pidle,
				"PlayerWalking": pwalk,
				"PlayerRunning": prun,
				"PlayerJumping": pjump,
			},
			Images: [4]*ebiten.Image{
				r.heightmap,
				SheetImage,
			},
		})
		// Mark frame as drawn for this tick
		r.drawn = true
	}

	opts := &ebiten.DrawImageOptions{
		//Filter: ebiten.FilterLinear,
	}
	//opts.GeoM.Scale(2, 2) // TODO: Resolution parameter
	//opts.GeoM.Translate(0, logic.ScreenHeight/2-logic.ScreenWidth/2)
	screen.DrawImage(r.offscreen, opts)
	screen.DrawImage(SheetImage, nil)
	//screen.DrawImage(r.heightmap, nil)
	/*opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(16, 16)
	screen.DrawImage(r.heightmap, opts)*/
}
