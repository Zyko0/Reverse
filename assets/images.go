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
	//go:embed images/material.jpg
	materialBytes []byte
	MaterialImage *ebiten.Image
)

func init() {
	img, err := jpeg.Decode(bytes.NewReader(materialBytes))
	if err != nil {
		log.Fatal(err)
	}
	MaterialImage = ebiten.NewImageFromImage(img)
}
