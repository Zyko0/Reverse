package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	boxIndices       = [6]uint16{0, 1, 2, 1, 2, 3}
	borderBoxIndices = []uint16{0, 2, 4, 2, 4, 6, 1, 3, 5, 3, 5, 7}

	BrushImage = ebiten.NewImage(1, 1)
)

func init() {
	BrushImage.Fill(color.White)
}

type QuadOpts struct {
	DstX, DstY          float32
	SrcX, SrcY          float32
	DstWidth, DstHeight float32
	SrcWidth, SrcHeight float32
	R, G, B, A          float32
}

func AppendQuadVerticesIndices(vertices []ebiten.Vertex, indices []uint16, index int, opts *QuadOpts) ([]ebiten.Vertex, []uint16) {
	sx, sy, dx, dy := opts.SrcX, opts.SrcY, opts.DstX, opts.DstY
	sw, sh, dw, dh := opts.SrcWidth, opts.SrcHeight, opts.DstWidth, opts.DstHeight
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices, []ebiten.Vertex{
		{
			DstX:   dx,
			DstY:   dy,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx,
			DstY:   dy + dh,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy + dh,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}...)

	indiceCursor := uint16(index * 4)
	indices = append(indices, []uint16{
		boxIndices[0] + indiceCursor,
		boxIndices[1] + indiceCursor,
		boxIndices[2] + indiceCursor,
		boxIndices[3] + indiceCursor,
		boxIndices[4] + indiceCursor,
		boxIndices[5] + indiceCursor,
	}...)

	return vertices, indices
}

func DrawRectBorder(dst *ebiten.Image, x, y, width, height, borderWidth, r, g, b, a float32) {
	dst.DrawTriangles([]ebiten.Vertex{
		{
			DstX:   x,
			DstY:   y,
			SrcX:   0,
			SrcY:   0,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX: x + borderWidth,
			DstY: y + borderWidth,
			SrcX: 0,
			SrcY: 0,
		},
		{
			DstX:   x + width,
			DstY:   y,
			SrcX:   1,
			SrcY:   0,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX: x + width - borderWidth,
			DstY: y + borderWidth,
			SrcX: 1,
			SrcY: 0,
		},
		{
			DstX:   x,
			DstY:   y + height,
			SrcX:   0,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX: x + borderWidth,
			DstY: y + height - borderWidth,
			SrcX: 0,
			SrcY: 1,
		},
		{
			DstX:   x + width,
			DstY:   y + height,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX: x + width - borderWidth,
			DstY: y + height - borderWidth,
			SrcX: 1,
			SrcY: 1,
		},
	}, borderBoxIndices, BrushImage, &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	})
}
