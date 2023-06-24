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
	ticks      uint64
	sheetdrawn bool
	mapdrawn   bool
	drawn      bool

	heightmap *ebiten.Image
	offscreen *ebiten.Image
}

type State struct {
	Level  int
	Map    *level.Map
	Camera *core.Camera
	Player *agents.Player
	Agent  agents.Agent
}

func NewRenderer() *Renderer {
	return &Renderer{
		ticks:      0,
		sheetdrawn: false,
		mapdrawn:   false,
		drawn:      false,

		heightmap: ebiten.NewImage(logic.MapWidth, logic.MapDepth),
		offscreen: ebiten.NewImage(logic.ScreenWidth, logic.ScreenHeight),
	}
}

func (r *Renderer) agentStates(a agents.Agent) (float32, float32, float32, float32) {
	var idle, walk, run, jump float32

	switch a.GetState() {
	case agents.StateIdle:
		idle = 1
	case agents.StateJumping:
		// Must have prio over walk / run
		jump = 1 - float32(a.GetJumpingTicks())/float32(agents.JumpingTicks)
		//pjump *= 2
	case agents.StateWalking:
		walk = 1
	case agents.StateRunning:
		run = 1
	case agents.StateFalling:
		jump = 1
		run = 1
	}

	return idle, walk, run, jump
}

func (r *Renderer) Update() {
	r.drawn = false
	r.ticks++
}

var (
	vertices []ebiten.Vertex
	indices  []uint16
)

func (r *Renderer) Draw(screen *ebiten.Image, state *State) {
	if !r.drawn {
		// Spritesheet update
		if !r.sheetdrawn {
			SheetImage.DrawImage(assets.Sheet0Image, nil)
			r.sheetdrawn = true
		}
		// Map generation
		if !r.mapdrawn {
			r.heightmap.WritePixels(state.Map.CompileBytes())
			r.mapdrawn = true
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
		// Agents states
		pidle, pwalk, prun, pjump := r.agentStates(state.Player)
		aidle, awalk, arun, ajump := r.agentStates(state.Agent)
		agentPosition := state.Agent.GetPosition()
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
				"Zoom":          float32(state.Camera.Zoom * logic.BaseZoom * float64(logic.MapWidth)),
				"AmbientColors": AmbientColorsByLevel[state.Level],
				// Player
				"PlayerIdle":    pidle,
				"PlayerWalking": pwalk,
				"PlayerRunning": prun,
				"PlayerJumping": pjump,
				// Agent
				"AgentPosition": []float32{
					float32(agentPosition.X) - logic.MapWidth/2,
					float32(agentPosition.Y),
					float32(agentPosition.Z) - logic.MapDepth/2,
				},
				"AgentIdle":    aidle,
				"AgentWalking": awalk,
				"AgentRunning": arun,
				"AgentJumping": ajump,
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
