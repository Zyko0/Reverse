package assets

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	defaultMusicVolume = 0.25
	defaultSFXVolume   = 0.75
)

var (
	ctx = audio.NewContext(44100)

	//go:embed sfx/fall.wav
	fallBytes  []byte
	fallPlayer *audio.Player

	//go:embed sfx/footsteps.wav
	footstepsBytes  []byte
	footstepsPlayer *audio.Player

	//go:embed sfx/scan.wav
	scanBytes  []byte
	scanPlayer *audio.Player
)

func init() {
	wavReader, err := wav.Decode(ctx, bytes.NewReader(fallBytes))
	if err != nil {
		log.Fatal(err)
	}
	fallPlayer, err = ctx.NewPlayer(wavReader)
	if err != nil {
		log.Fatal(err)
	}
	fallPlayer.SetVolume(defaultSFXVolume)

	wavReader, err = wav.Decode(ctx, bytes.NewReader(footstepsBytes))
	if err != nil {
		log.Fatal(err)
	}
	infiniteReader := audio.NewInfiniteLoop(wavReader, wavReader.Length())
	footstepsPlayer, err = ctx.NewPlayer(infiniteReader)
	if err != nil {
		log.Fatal(err)
	}
	footstepsPlayer.SetVolume(defaultSFXVolume)

	wavReader, err = wav.Decode(ctx, bytes.NewReader(scanBytes))
	if err != nil {
		log.Fatal(err)
	}
	scanPlayer, err = ctx.NewPlayer(wavReader)
	if err != nil {
		log.Fatal(err)
	}
	scanPlayer.SetVolume(defaultSFXVolume)
}

func PlayFall() {
	fallPlayer.Rewind()
	fallPlayer.Play()
}

func PlayFootsteps() {
	if !footstepsPlayer.IsPlaying() {
		footstepsPlayer.Rewind()
		footstepsPlayer.Play()
	}
}

func StopFootsteps() {
	footstepsPlayer.Pause()
}

func PlayScan() {
	scanPlayer.Rewind()
	scanPlayer.Play()
}
