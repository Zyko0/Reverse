package assets

import (
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
)

var (
	MapGenFontFace font.Face
)

func init() {
	tfont, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		log.Fatal(err)
	}
	MapGenFontFace = truetype.NewFace(tfont, &truetype.Options{
		Size: 27,
	})
}
