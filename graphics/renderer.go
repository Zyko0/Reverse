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
	ticks    uint64
	mapdrawn bool
	drawn    bool

	heightmap *ebiten.Image
	offscreen *ebiten.Image
}

type State struct {
	Level      int
	Map        *level.HMap
	Camera     *core.Camera
	Player     *agents.Player
	Agent      agents.Agent
	GameStatus core.GameStatus
	PlayerSeen bool
	AgentSeen  bool
}

func NewRenderer() *Renderer {
	return &Renderer{
		ticks:    0,
		mapdrawn: false,
		drawn:    false,

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
		ptasing := float32(0)
		if state.Player.HasAbility(agents.AbilityTasing) {
			ptasing = 1
		}
		tasingClr := []float32{1, 0, 0}
		if state.GameStatus == core.GameStatusVictory {
			ptasing = 1
			tasingClr = []float32{0, 0, 1}
		}
		aidle, awalk, arun, ajump := r.agentStates(state.Agent)
		agentPosition := state.Agent.GetPosition()
		// Abilities
		lightsOff := float32(0)
		if state.Agent.HasAbility(agents.AbilityLightsOff) {
			lightsOff = 1
		}
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
				"TaseColor":     tasingClr,
				"TaseRadius":    float32(agents.TasingRadius),
				"AmbientColors": AmbientColorsByLevel[state.Level],
				// Player
				"PlayerIdle":    pidle,
				"PlayerWalking": pwalk,
				"PlayerRunning": prun,
				"PlayerJumping": pjump,
				"PlayerTasing":  ptasing,
				// Agent
				"AgentPosition": []float32{
					float32(agentPosition.X) - logic.MapWidth/2,
					float32(agentPosition.Y),
					float32(agentPosition.Z) - logic.MapDepth/2,
				},
				"AgentEyes":    AgentEyesByLevel[state.Level],
				"AgentColor":   AgentColorsByLevel[state.Level],
				"AgentIdle":    aidle,
				"AgentWalking": awalk,
				"AgentRunning": arun,
				"AgentJumping": ajump,
				// Abilities
				"LightsOff": lightsOff,
			},
			Images: [4]*ebiten.Image{
				r.heightmap,
				assets.Sheet0Image,
			},
		})
		// Render minimap
		const minimapSize = 256
		vertices, indices = AppendQuadVerticesIndices(
			vertices[:0], indices[:0], 0, &QuadOpts{
				DstX:      logic.ScreenWidth - minimapSize*1.5,
				DstY:      0,
				DstWidth:  minimapSize * 1.5,
				DstHeight: minimapSize * 1.5,
				SrcWidth:  minimapSize,
				SrcHeight: minimapSize,
			},
		)
		// Agent properties
		pseen := float32(0)
		psig := float32(0)
		aseen := float32(0)
		asig := float32(0)
		if state.PlayerSeen {
			pseen = 1
		}
		if state.Player.GetHeard() {
			psig = 1
		}
		if state.AgentSeen {
			aseen = 1
		}
		if state.Agent.GetHeard() || state.Player.HasAbility(agents.AbilityScanning) {
			asig = 1
			aseen = 1
		}
		r.offscreen.DrawTrianglesShader(vertices, indices, assets.MinimapShader, &ebiten.DrawTrianglesShaderOptions{
			Images: [4]*ebiten.Image{
				r.heightmap,
			},
			Uniforms: map[string]any{
				"Time": float64(r.ticks) / logic.TPS,
				"MapSize": []float32{
					logic.MapWidth, logic.MapDepth,
				},
				"GoalPosition": []float32{
					float32(level.GoalPosition.X),
					float32(level.GoalPosition.Z),
				},
				"PlayerPosition": []float32{
					float32(state.Player.Position.X),
					float32(state.Player.Position.Z),
				},
				"PlayerSignal": float32(psig),
				"PlayerTasing": ptasing,
				"PlayerSeen":   pseen,
				"TaseColor":    tasingClr,
				"AgentPosition": []float32{
					float32(agentPosition.X),
					float32(agentPosition.Z),
				},
				"AgentSignal": float32(asig),
				"AgentSeen":   aseen,
			},
		})
		// Render abilities cooldowns
		r.renderAbilities(state)
		// Mark frame as drawn for this tick
		r.drawn = true
	}

	opts := &ebiten.DrawImageOptions{}
	//opts.GeoM.Scale(2, 2) // TODO: Resolution parameter
	//opts.GeoM.Translate(0, logic.ScreenHeight/2-logic.ScreenWidth/2)
	screen.DrawImage(r.offscreen, opts)
}
