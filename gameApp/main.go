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

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Смутьяны!")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
