package assets

import (
	_ "embed"
	"log"

	"github.com/Zyko0/Reverse/pkg/level"
)

var (
	//go:embed levels/0.rev
	level0Src []byte
	Level0    = &level.HMap{}
)

func init() {
	err := Level0.Deserialize(level0Src)
	if err != nil {
		log.Fatal(err)
	}
}
