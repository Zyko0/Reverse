package level

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"time"

	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

type HColumn struct {
	Height   byte
	RowIndex byte
	Rule     byte
}

type HMap struct {
	HeightMap      [][]HColumn
	Neighbours1    [][]uint16
	Neighbours2    [][]uint16
	Neighbours3    [][]uint16
	Neighbours3Run [][]uint16
}

func (hm *HMap) At(x, z int) HColumn {
	return hm.HeightMap[z][x]
}

func (hm *HMap) CompileBytes() []byte {
	pixels := make([]byte, logic.MapDepth*logic.MapWidth*4)
	for z := 0; z < logic.MapDepth; z++ {
		for x := 0; x < logic.MapWidth; x++ {
			const rowSize = logic.MapWidth * 4

			i := z*rowSize + x*4
			c := hm.At(x, z)
			pixels[i+0] = byte(float64(c.Height) / logic.MapHeight * 255)
			pixels[i+1] = c.RowIndex
			pixels[i+2] = c.Rule
			pixels[i+3] = 255
		}
	}

	return pixels
}

func (hm *HMap) GetReachableNeighbours(x, y, z, agility int, allowRun bool) []geom.Vec3 {
	var ns []geom.Vec3

	runc := 1
	if allowRun {
		runc = 3 // Agent running ms multiplier
	}
	for xoff := -agility * runc; xoff <= agility*runc; xoff++ {
		nx := x + xoff
		if nx < 0 || nx >= logic.MapWidth {
			continue
		}
		for zoff := -agility * runc; zoff <= agility*runc; zoff++ {
			nz := z + zoff
			if nz < 0 || nz >= logic.MapDepth {
				continue
			}
			// Avoid current block // TODO: maybe allow it
			/*if xoff == 0 && zoff == 0 {
				continue
			}*/
			h := int(hm.HeightMap[nz][nx].Height)
			// Too high can't reach (thcr)
			if h > y && h-y > 1 {
				continue
			}
			ns = append(ns, geom.Vec3{
				X: float64(nx),
				Y: float64(h),
				Z: float64(nz),
			})
		}
	}

	return ns
}

func (hm *HMap) Deserialize(data []byte) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(hm)
	if err != nil {
		return err
	}

	return nil
}

func (hm *HMap) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(hm)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (hm *HMap) getNeighboursHashes(agility int, run bool) [][]uint16 {
	t := time.Now()
	neighbours := make([][]uint16, 65536)
	for x := 0; x < logic.MapWidth; x++ {
		for z := 0; z < logic.MapDepth; z++ {
			c := hm.HeightMap[z][x]
			pos := geom.Vec3{
				X: float64(x),
				Y: float64(c.Height),
				Z: float64(z),
			}
			ns := hm.GetReachableNeighbours(x, int(c.Height), z, agility, run)
			hashes := make([]uint16, len(ns))
			for i, n := range ns {
				hashes[i] = n.AsUHashXZ()
			}
			neighbours[pos.AsUHashXZ()] = hashes
		}
	}
	fmt.Println("neighbourg init - agility", agility, "run", run, "time", time.Since(t))

	return neighbours
}

func (hm *HMap) BuildStaticNeighbours() {
	hm.Neighbours1 = hm.getNeighboursHashes(1, false)
	//hm.Neighbours2 = hm.getNeighboursHashes(2, false)
	//hm.Neighbours3 = hm.getNeighboursHashes(3, false)
	//hm.Neighbours3Run = hm.getNeighboursHashes(3, true)
}

func (hm *HMap) CastRay(src, dst geom.Vec3, max float64) bool {
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
		vox := hm.HeightMap[int(z)][int(x)]
		if y <= float64(vox.Height) {
			return false
		}
		x += dx
		y += dy
		z += dz
	}

	return true
}

var (
	qStore        [1024]uint16
	nullNodeValue = uint16(65535)
	nullNodeSlice = make([]uint16, 65536)
	from          = make([]uint16, 65536)
)

func init() {
	for i := range nullNodeSlice {
		nullNodeSlice[i] = nullNodeValue
	}
}

func (hm *HMap) BFS(start, goal geom.Vec3, agility int, allowRun bool) ([]geom.Vec3, bool) {
	startHash, goalHash := start.AsUHashXZ(), goal.AsUHashXZ()
	queue := qStore[:1]
	queue[0] = startHash

	copy(from, nullNodeSlice)

	found := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == goalHash {
			found = true
			break
		}
		/*for _, next := range neighbours[current] { //lm.GetReachableNeighbours(int(pos.X), int(pos.Y), int(pos.Z), agility, allowRun) {
			if c := from[next]; c == nullNodeValue {
				queue = append(queue, next)
				from[next] = current
			}
		}*/
	}
	if !found {
		return nil, false
	}

	current := goalHash
	path := make([]geom.Vec3, 0, 1024)
	for current != startHash {
		v := geom.Vec3{}.FromUHashXZ(current)
		v.Y = float64(hm.HeightMap[int(v.Z)][int(v.X)].Height)
		path = append(path, v)
		current = from[current]
	}
	// Reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, true
}
