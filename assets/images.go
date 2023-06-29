package assets

import (
	"bytes"
	_ "embed"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MaterialTextureSize = 128
)

var (
	//go:embed images/sheet0.jpg
	sheet0Bytes []byte
	Sheet0Image *ebiten.Image

	//go:embed images/cursor.png
	cursorBytes []byte
	CursorImage *ebiten.Image

	//go:embed images/cursor_hover.png
	cursorHoverBytes []byte
	CursorHoverImage *ebiten.Image

	//go:embed images/cursor_click.png
	cursorClickBytes []byte
	CursorClickImage *ebiten.Image
)

func init() {
	img, err := jpeg.Decode(bytes.NewReader(sheet0Bytes))
	if err != nil {
		log.Fatal(err)
	}
	Sheet0Image = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(cursorBytes))
	if err != nil {
		log.Fatal(err)
	}
	CursorImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(cursorHoverBytes))
	if err != nil {
		log.Fatal(err)
	}
	CursorHoverImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(cursorClickBytes))
	if err != nil {
		log.Fatal(err)
	}
	CursorClickImage = ebiten.NewImageFromImage(img)
}
