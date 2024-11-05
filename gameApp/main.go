package main

import (
	"log"
	"strategy-game/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Troublemakers!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}

}
