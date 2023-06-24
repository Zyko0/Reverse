package assets

import (
	"bytes"
	_ "embed"
	"image/jpeg"
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
)

func init() {
	img, err := jpeg.Decode(bytes.NewReader(sheet0Bytes))
	if err != nil {
		log.Fatal(err)
	}
	Sheet0Image = ebiten.NewImageFromImage(img)
}
