package level

import (
	"bytes"
	"encoding/gob"

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

func (lm *Map) CompileBytes() []byte {
	pixels := make([]byte, logic.MapDepth*logic.MapWidth*4)
	for z := 0; z < logic.MapDepth; z++ {
		for x := 0; x < logic.MapWidth; x++ {
			i := z*(logic.MapWidth*4) + x*4
			c := lm.HeightMap[z][x]
			pixels[i+0] = byte(float64(c.Height) / logic.MapHeight * 255)
			pixels[i+1] = c.TopMaterialID
			pixels[i+2] = c.SideMaterialID
			pixels[i+3] = 255
		}
	}

	return pixels
}
