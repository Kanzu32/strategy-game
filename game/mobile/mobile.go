package mobile

import (
	"strategy-game/game"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func init() {
	// ebiten.SetWindowSize(640, 480)
	// ebiten.SetWindowTitle("Troublemakers!")
	mobile.SetGame(game.NewGame())
}

func Dummy() {}
