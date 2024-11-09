package gamedata

import "github.com/hajimehoshi/ebiten/v2"

type GameData interface {
	FrameCount() int
	RenderHeight() int
	RenderWidth() int
	View() *ebiten.Image
	ViewScale() int
	ViewScaleInc()
	ViewScaleDec()
}
