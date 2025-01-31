package main

import (
	"log"
	"strategy-game/game"
	"strategy-game/util/teams"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO SPLIT TARGET UNIT AND TARGET OBJECT

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Смутьяны!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game.NewGame(teams.Blue)); err != nil {
		log.Fatal(err)
	}
}
