package level

import (
	"bytes"
	"encoding/gob"
	"math"

	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
)

type HColumn struct {
	Height   byte
	RowIndex byte
	Rule     byte
}

type HMap struct {
	HeightMap   [][]HColumn
	Neighbours1 [][]uint16
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

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
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
			v := geom.Vec3{
				X: float64(nx),
				Y: float64(h),
				Z: float64(nz),
			}
			// Check if diagonal is accessible
			if xoff != 0 && zoff != 0 {
				h0, h1 := y-1, y-1
				if tz := z + sign(zoff); tz >= 0 && tz < logic.MapDepth {
					h0 = int(hm.HeightMap[tz][x].Height)
				}
				if tx := x + sign(xoff); tx >= 0 && tx < logic.MapWidth {
					h1 = int(hm.HeightMap[z][tx].Height)
				}
				if (h0 < y || h0-y < 2) && (h1 < y || h1-y < 2) {
					ns = append(ns, v)
				}
			} else {
				ns = append(ns, v)
			}
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

	return neighbours
}

func (hm *HMap) BuildStaticNeighbours() {
	hm.Neighbours1 = hm.getNeighboursHashes(1, false)
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
	qStore         [1024]uint16
	nullNodeValue  = uint16(65535)
	nullNodeSlice  = make([]uint16, 65536)
	zeroSlice      = make([]uint16, 65536)
	from           = make([]uint16, 65536)
	costs          = make([]uint16, 65536)
	proximityCosts = make([]uint16, 65536)
)

func init() {
	for i := range nullNodeSlice {
		nullNodeSlice[i] = nullNodeValue
	}
}

func (hm *HMap) cost(current, next uint16) uint16 {
	if (current&0xFF00)^(next&0xFF00) > 0 && (current&0x00FF)^(next&0x00FF) > 0 {
		return 14 // sqrt2 for diagonal
	}
	return 10
}

func (hm *HMap) AStar(start, goal geom.Vec3) ([]geom.Vec3, bool) {
	startHash, goalHash := start.AsUHashXZ(), goal.AsUHashXZ()
	queue := qStore[:1]
	queue[0] = startHash

	copy(from, nullNodeSlice)
	copy(costs, nullNodeSlice)
	costs[startHash] = proximityCosts[startHash]
	found := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == goalHash {
			found = true
			break
		}
		for _, next := range hm.Neighbours1[current] {
			nc := hm.cost(current, next)
			// If jump involved make the cost higher
			h0 := hm.HeightMap[current>>8][current&255].Height
			h1 := hm.HeightMap[next>>8][next&255].Height
			if h1 > h0 {
				nc *= 10
			}
			nc += costs[current]
			if c := costs[next]; nc < c && proximityCosts[next] <= proximityCosts[current] {
				costs[next] = nc
				from[next] = current
				queue = append(queue, next)
			}
		}
	}
	if !found {
		return nil, false
	}

	current := goalHash
	path := make([]geom.Vec3, 0, 256)
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

var (
	proximities = make([]uint16, 512)
)

func (hm *HMap) UpdateProximityCosts(pos geom.Vec3, radius float64) {
	r := uint16(radius * 10)
	hash := pos.AsUHashXZ()
	queue := qStore[:1]
	queue[0] = hash

	copy(proximityCosts, zeroSlice)
	copy(costs, nullNodeSlice)
	costs[hash] = 0
	proximities = proximities[:0]
	for len(queue) > 0 {
		current := queue[0]
		proximities = append(proximities, current)
		queue = queue[1:]
		for _, next := range hm.Neighbours1[current] {
			nc := costs[current] + 10 //hm.cost(current, next)
			if nc < r && nc < costs[next] {
				costs[next] = nc
				queue = append(queue, next)
			}
		}
	}
	// Update player proximity costs
	for _, h := range proximities {
		proximityCosts[h] = r - costs[h]
	}
}
