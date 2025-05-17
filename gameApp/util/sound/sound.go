package sound

import (
	"bytes"
	"strategy-game/assets"
	"strategy-game/game/singletons"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var musicPlayer *audio.Player
var hitPlayer *audio.Player
var killPlayer *audio.Player
var context *audio.Context

func Init() {
	context := audio.NewContext(48000)
	wavStream, err := wav.DecodeF32(bytes.NewReader(assets.Hit2Sound))
	if err != nil {
		panic(err)
	}
	hitPlayer, err = context.NewPlayerF32(wavStream)
	if err != nil {
		panic(err)
	}

	wavStream, err = wav.DecodeF32(bytes.NewReader(assets.KillSound))
	if err != nil {
		panic(err)
	}
	killPlayer, err = context.NewPlayerF32(wavStream)
	if err != nil {
		panic(err)
	}

	wavStream, err = wav.DecodeF32(bytes.NewReader(assets.Music1))
	if err != nil {
		panic(err)
	}
	musicPlayer, err = context.NewPlayerF32(wavStream)
	if err != nil {
		panic(err)
	}

}

func PlayHit() {
	hitPlayer.SetVolume(float64(singletons.Settings.Sound) / 50)
	hitPlayer.Rewind()
	hitPlayer.Play()
}

func PlayKill() {
	killPlayer.SetVolume(float64(singletons.Settings.Sound) / 50)
	killPlayer.Rewind()
	killPlayer.Play()
}

func RestartMusicIfNeeds() {
	musicPlayer.SetVolume(float64(singletons.Settings.Sound) / 50)
	if !musicPlayer.IsPlaying() {
		musicPlayer.Play()
	}
}

// func PlayMusic() {
// 	musicPlayer.
// }
