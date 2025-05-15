package main

import (
	"log"
	"strategy-game/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO WALK
// TODO GAME UI
// TODO ATTACK SYSTEM
// TODO упакоука
//
// TODO ...

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Смутьяны!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
