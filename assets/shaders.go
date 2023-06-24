package assets

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaders/scene.kage
	sceneShaderSrc []byte
	SceneShader    *ebiten.Shader

	//go:embed shaders/map.kage
	mapShaderSrc []byte
	MapShader    *ebiten.Shader

	//go:embed shaders/character.kage
	charShaderSrc []byte
	CharShader    *ebiten.Shader
)

func init() {
	var err error

	SceneShader, err = ebiten.NewShader(sceneShaderSrc)
	if err != nil {
		log.Fatal(err)
	}

	MapShader, err = ebiten.NewShader(mapShaderSrc)
	if err != nil {
		log.Fatal(err)
	}

	CharShader, err = ebiten.NewShader(charShaderSrc)
	if err != nil {
		log.Fatal(err)
	}
}
