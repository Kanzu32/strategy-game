package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strategy-game/game"
	"strategy-game/game/singletons"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	file, err := os.Open("conf.json")
	if err != nil {
		println("configuration missing, standard settings")
		singletons.Settings.DefaultGameScale = 4
		singletons.Settings.Fullscreen = false
		singletons.Settings.Sound = 3
		singletons.Settings.Language = "Eng"

		file, err = os.Create("conf.json")
		if err != nil {
			println("Unable to create new config file:", err.Error())
		}

		b, err := json.Marshal(singletons.Settings)
		if err != nil {
			println("Unable to marshal settings:", err.Error())
		}
		_, err = file.Write(b)
		if err != nil {
			println("Unable to write new config:", err.Error())
		}
	} else {
		buf := make([]byte, 1024)
		_, err = file.Read(buf)
		if err != nil {
			println("Unable to read config file:", err.Error())
		}
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &singletons.Settings)
		// err = json.Unmarshal(buf, &singletons.Settings)
		if err != nil {
			println("Unable to unmarshal config file:", err.Error())
		}
	}
	file.Close()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Смутьяны!")
	ebiten.SetFullscreen(singletons.Settings.Fullscreen)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
