package level

import (
	"bytes"
	"encoding/gob"
	"math"

	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

type Column struct {
	Height         byte
	TopMaterialID  byte
	SideMaterialID byte
}

type Map struct {
	HeightMap [][]Column
	Start     geom.Vec3
	Goal      geom.Vec3
}

func (lm *Map) Deserialize(data []byte) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(lm)
	if err != nil {
		return err
	}

	return nil
}

func (lm *Map) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(lm)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (lm *Map) At(x, z int) Column {
	return lm.HeightMap[z][x]
}

func (lm *Map) CompileBytes() []byte {
	pixels := make([]byte, logic.MapDepth*logic.MapWidth*4)
	for z := 0; z < logic.MapDepth; z++ {
		for x := 0; x < logic.MapWidth; x++ {
			const rowSize = logic.MapWidth * 4

			i := z*rowSize + x*4
			c := lm.At(x, z)
			pixels[i+0] = byte(float64(c.Height) / logic.MapHeight * 255)
			pixels[i+1] = c.TopMaterialID
			pixels[i+2] = c.SideMaterialID
			pixels[i+3] = 255
		}
	}

	return pixels
}

func (lm *Map) CastRay(src, dst geom.Vec3, max float64) bool {
	// Test the ray from head to head
	src.Y, dst.Y = src.Y+0.49, dst.Y+0.49
	dist := src.DistanceTo(dst)
	if dist > max {
		return false
	}
	max = dist
	// DDA
	dx := dst.X - src.X
	dy := dst.Y - src.Y
	dz := dst.Z - src.Z
	step := 0.
	if math.Abs(dx) >= math.Abs(dy) {
		step = math.Abs(dx)
	} else {
		step = math.Abs(dy)
	}
	if math.Abs(dz) >= step {
		step = math.Abs(dz)
	}
	dx /= step
	dy /= step
	dz /= step
	x := src.X
	y := src.Y
	z := src.Z
	for i := 1.; i <= step; i++ {
		vox := lm.HeightMap[int(z)][int(x)]
		if y <= float64(vox.Height) {
			return false
		}
		x += dx
		y += dy
		z += dz
	}

	return true
}

func (lm *Map) GetReachableNeighbours(x, y, z, agility int, allowRun bool) []geom.Vec3 {
	var neighbours []geom.Vec3

	for xoff := -agility; xoff <= agility; xoff++ {
		nx := x + xoff
		if nx < 0 || nx > logic.MapWidth {
			continue
		}
		/*for yoff := -1; yoff <= 1.; yoff++ {
		ny := y+yoff
		if ny < 0 || ny > logic.MapHeight {
			continue
		}*/
		for zoff := -agility; zoff <= agility; zoff++ {
			nz := z + zoff
			if nz < 0 || nz > logic.MapDepth {
				continue
			}
			neighbours = append(neighbours, geom.Vec3{
				X: float64(x),
				Y: float64(lm.HeightMap[z][x].Height),
				Z: float64(z),
			})
		}
		//}
	}

	return neighbours
}
