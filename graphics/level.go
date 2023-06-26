package graphics

const (
	EyesModifiersNone float32 = iota
	EyesModifiersJoyful
	EyesModifiersAngry
)

var (
	AmbientColorsByLevel = [][]float32{
		0: {
			1, 1, 1,
			0.5, 0, 1,
			1, 0, 0.5,
			0.5, 0, 1,
			1, 0, 0.5,
			1, 1, 1,
		},
	}
	AgentColorsByLevel = [][]float32{
		0: {1, 0.5, 0.25},
		1: {0.5, 0.25, 1},
	}
	AgentEyesByLevel = []float32{
		0: EyesModifiersJoyful,
		1: EyesModifiersAngry,
	}
)
