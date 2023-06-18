package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/Zyko0/Reverse/assets"
	"github.com/Zyko0/Reverse/graphics"
	"github.com/Zyko0/Reverse/logic"
	"github.com/Zyko0/Reverse/pkg/geom"
	"github.com/Zyko0/Reverse/pkg/level"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	BrushImage = ebiten.NewImage(3, 3)
)

func init() {
	BrushImage.Fill(color.White)
}

const (
	TileSize = 32
)

type App struct {
	err      error
	filename string
	lastsave time.Time

	offset      geom.Vec2
	selected    bool
	selectedSrc geom.Vec2
	selectedDst geom.Vec2
	hovered     geom.Vec2
	zoom        float64

	lvl *level.Map

	txtImg *ebiten.Image
}

func New(lvl *level.Map, filename string) *App {
	return &App{
		err:      nil,
		filename: filename,
		lastsave: time.Time{},

		offset: geom.Vec2{
			X: logic.MapWidth / 2,
			Y: logic.MapDepth / 2,
		},
		selected: false,
		zoom:     1,

		lvl: lvl,

		txtImg: ebiten.NewImage(logic.ScreenWidth, logic.ScreenHeight),
	}
}

func (g *App) Save() {
	g.err = nil

	data, err := g.lvl.Serialize()
	if err != nil {
		g.err = err
		return
	}

	f, err := os.OpenFile(g.filename, os.O_CREATE, 0644)
	if err != nil {
		g.err = err
		return
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		g.err = err
		return
	}

	g.lastsave = time.Now()
}

func (g *App) Update() error {
	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("quit")
	}
	// Save
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Save()
	}
	// Zoom
	if _, yoff := ebiten.Wheel(); yoff != 0 {
		const zoomSens = 0.05

		g.zoom -= yoff * zoomSens
		if g.zoom < zoomSens {
			g.zoom = zoomSens
		}
	}
	// Map traversal
	const moveSens = 0.05
	if d := inpututil.KeyPressDuration(ebiten.KeyRight); d > 0 {
		g.offset.X += float64(d) * moveSens * g.zoom
	}
	if d := inpututil.KeyPressDuration(ebiten.KeyLeft); d > 0 {
		g.offset.X -= float64(d) * moveSens * g.zoom
	}
	if d := inpututil.KeyPressDuration(ebiten.KeyUp); d > 0 {
		g.offset.Y -= float64(d) * moveSens * g.zoom
	}
	if d := inpututil.KeyPressDuration(ebiten.KeyDown); d > 0 {
		g.offset.Y += float64(d) * moveSens * g.zoom
	}
	// Mouse coordinates
	x, y := ebiten.CursorPosition()
	mapw := float64(logic.MapWidth * TileSize / g.zoom)
	mapd := float64(logic.MapDepth * TileSize / g.zoom)
	x -= logic.ScreenWidth / 2
	y -= logic.ScreenHeight / 2
	vx := (float64(x) + g.offset.X*TileSize/g.zoom) / mapw * logic.MapWidth //+ g.offset.X
	vz := (float64(y) + g.offset.Y*TileSize/g.zoom) / mapd * logic.MapDepth //+ g.offset.Y
	blockx, blockz := math.Floor(vx), math.Floor(vz)
	if blockx < 0 {
		blockx = 0
	}
	if blockx >= logic.MapWidth {
		blockx = logic.MapWidth - 1
	}
	if blockz < 0 {
		blockz = 0
	}
	if blockz >= logic.MapDepth {
		blockz = logic.MapDepth - 1
	}
	// Mouse hovering
	g.hovered.X, g.hovered.Y = blockx, blockz
	// Mouse selection
	dst := ebiten.IsKeyPressed(ebiten.KeyShift)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if !dst {
			// Reset selection if clicking the already selected spot
			if g.selected && g.selectedSrc == g.hovered {
				g.selected = false
			} else {
				g.selectedSrc = g.hovered
				g.selectedDst = g.hovered
				g.selected = true
			}
		} else {
			if g.selected {
				if g.selectedSrc == g.hovered {
					g.selected = false
				} else {
					g.selectedDst = g.hovered
				}
			}
		}
	}
	// Reset selection by right clicking also
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.selected = false
	}
	// Height, top and side material modification
	hoff, tmoff, smoff := 0, 0, 0
	if g.selected {
		if inpututil.IsKeyJustPressed(ebiten.KeyNumpadAdd) {
			hoff = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyNumpadSubtract) {
			hoff = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyI) {
			tmoff = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyK) {
			tmoff = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyO) {
			smoff = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			smoff = -1
		}
		if hoff != 0 || tmoff != 0 || smoff != 0 {
			area := image.Rect(
				int(g.selectedSrc.X), int(g.selectedSrc.Y),
				int(g.selectedDst.X), int(g.selectedDst.Y),
			)
			for az := area.Min.Y; az <= area.Max.Y; az++ {
				for ax := area.Min.X; ax <= area.Max.X; ax++ {
					c := g.lvl.HeightMap[az][ax]
					h := c.Height
					tm := c.TopMaterialID
					sm := c.SideMaterialID
					if hoff == -1 && h > 0 {
						h -= 1
					} else if hoff == 1 && h < logic.MapHeight {
						h += 1
					}
					if tmoff == -1 && tm > 0 {
						tm -= 1
					} else if tmoff == 1 && tm < 255 {
						tm += 1
					}
					if smoff == -1 && sm > 0 {
						sm -= 1
					} else if smoff == 1 && sm < 255 {
						sm += 1
					}
					g.lvl.HeightMap[az][ax].Height = h
					g.lvl.HeightMap[az][ax].TopMaterialID = tm
					g.lvl.HeightMap[az][ax].SideMaterialID = sm
				}
			}
		}
	}

	return nil
}

var (
	vertices []ebiten.Vertex
	indices  []uint16
)

func (g *App) Draw(screen *ebiten.Image) {
	g.txtImg.Clear()
	// Display grid
	vertices, indices = vertices[:0], indices[:0]
	index := 0
	size := float32(TileSize / g.zoom)
	mapsize := float64(logic.MapWidth * size)
	off := geom.Vec2{
		X: logic.ScreenWidth / 2,
		Y: logic.ScreenHeight / 2,
	}.Sub(g.offset.Div(geom.Vec2{
		X: logic.MapWidth,
		Y: logic.MapDepth,
	}).Mul(geom.Vec2{
		X: mapsize,
		Y: mapsize,
	}))
	selectArea := image.Rect(
		int(g.selectedSrc.X), int(g.selectedSrc.Y),
		int(g.selectedDst.X), int(g.selectedDst.Y),
	)
	selectArea.Max.X += 1
	selectArea.Max.Y += 1
	// TODO: start z,x values directly at the offsets it matters instead of the lazy 'continue'
	txtColor := ebiten.ColorScale{}
	txtColor.Scale(0, 0, 0, 1)
	vxBuffer := [][]ebiten.Vertex{}
	ixBuffer := [][]uint16{}
	for z := 0.; z < logic.MapDepth; z++ {
		vz := float32(z)*size + float32(off.Y)
		if vz > logic.ScreenHeight || vz+size < 0 {
			continue
		}
		for x := 0.; x < logic.MapWidth; x++ {
			vx := float32(x)*size + float32(off.X)
			if vx > logic.ScreenWidth || vx+size < 0 {
				continue
			}

			var tr, tg, tb float32
			block := g.lvl.HeightMap[int(z)][int(x)]
			// Warmer color based on height
			tr = 1
			tg = 1 - float32(block.Height)/logic.MapHeight
			tb = tg - 1
			// Mark magenta if selected
			pt := image.Point{
				X: int(x),
				Y: int(z),
			}
			if g.selected && pt.In(selectArea) {
				tr, tg, tb = 0.5, 0, 0.5
			}
			// Mark blue if hovered
			if x == g.hovered.X && z == g.hovered.Y {
				tr, tg, tb = 0, 0, 0.5
			}
			// Tile content
			vertices, indices = graphics.AppendQuadVerticesIndices(
				// +1 -2 to account for 1px border
				vertices, indices, index, &graphics.QuadOpts{
					DstX:      vx + 1,
					DstY:      vz + 1,
					DstWidth:  size - 2,
					DstHeight: size - 2,
					SrcWidth:  3,
					SrcHeight: 3,
					R:         tr,
					G:         tg,
					B:         tb,
					A:         1,
				},
			)
			index++
			// If overflowing indices buffer need to flush to a draw call
			if len(indices)+6 > ebiten.MaxIndicesCount {
				//screen.DrawTriangles(vertices, indices, BrushImage, nil)
				vx := make([]ebiten.Vertex, len(vertices))
				ix := make([]uint16, len(indices))
				copy(vx, vertices)
				copy(ix, indices)
				vxBuffer = append(vxBuffer, vx)
				ixBuffer = append(ixBuffer, ix)
				vertices, indices = vertices[:0], indices[:0]
				index = 0
			}
			// Text
			const txtScale = 1. / 3.
			str := fmt.Sprintf("H:%d\nT:%d\nS:%d",
				block.Height, block.TopMaterialID, block.SideMaterialID,
			)
			rect := text.BoundString(assets.MapGenFontFace, str)
			opts := &ebiten.DrawImageOptions{
				ColorScale: txtColor,
			}
			opts.GeoM.Scale(1/g.zoom*txtScale, 1/g.zoom*txtScale)
			opts.GeoM.Translate(
				float64(vx+size/2)-float64(rect.Dx())/2*txtScale/g.zoom,
				float64(vz+size/3),
			)
			text.DrawWithOptions(g.txtImg, str, assets.MapGenFontFace, opts)
		}
	}
	// Draw tiles
	for i := range vxBuffer {
		screen.DrawTriangles(vxBuffer[i], ixBuffer[i], BrushImage, nil)
	}
	screen.DrawTriangles(vertices, indices, BrushImage, nil)
	// Draw text
	screen.DrawImage(g.txtImg, nil)
	// Debug
	debugStr := fmt.Sprintf("Zoom: %.2f Hovered: (%.1f, %.1f)", g.zoom, g.hovered.X, g.hovered.Y)
	if g.selected {
		debugStr += fmt.Sprintf(" Selection(<x0:%.2f; z0:%.2f>, <x1:%.2f, z1:%.2f>)",
			g.selectedSrc.X, g.selectedSrc.Y,
			g.selectedDst.X, g.selectedDst.Y,
		)
	} else {
		debugStr += " Selection(<none>)"
	}
	if !g.lastsave.IsZero() {
		debugStr += fmt.Sprintf(" saved: %v - err: %v", g.lastsave.Format(time.RFC1123), g.err)
	}
	ebitenutil.DebugPrint(screen, debugStr)
	// Commands
	ebitenutil.DebugPrintAt(screen, "- Left click to select (shift click for a region) - Right click to reset", 0, 12)
	ebitenutil.DebugPrintAt(screen, "- Press pnum + or - to change height", 0, 24)
	ebitenutil.DebugPrintAt(screen, "- Press I/K to change top material id", 0, 36)
	ebitenutil.DebugPrintAt(screen, "- Press O/L to change side material id", 0, 48)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logic.ScreenWidth, logic.ScreenHeight
}

func main() {
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")

	if len(os.Args) <= 1 {
		log.Fatal("need to provide a input/output map filename")
	}
	filename := os.Args[1]
	fmt.Println("filename", filename)

	f, err := os.Open(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("couldn't open file: %v", err)
		}
	}

	lvl := &level.Map{}
	data := []byte{}
	if f != nil {
		data, err = io.ReadAll(f)
		if err != nil {
			f.Close()
			log.Fatalf("couldn't read file's content: %v", err)
		}
		if err = lvl.Deserialize(data); err != nil {
			f.Close()
			log.Fatalf("couldn't deserialize level data: %v", err)
		}
	}
	f.Close()
	// Initialize the slice if empty, so that there are default value for the editor
	if len(lvl.HeightMap) == 0 {
		lvl.HeightMap = make([][]level.Column, logic.MapDepth)
		for z := 0; z < logic.MapDepth; z++ {
			lvl.HeightMap[z] = make([]level.Column, logic.MapWidth)
		}
	}

	ebiten.SetFullscreen(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetMaxTPS(logic.TPS)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	editor := New(lvl, filename)
	if err := ebiten.RunGame(editor); err != nil {
		fmt.Println("err:", err)
	}
	// Save on exit
	if editor.Save(); editor.err != nil {
		log.Fatalf("couldn't save level: %v", err)
	}
}
