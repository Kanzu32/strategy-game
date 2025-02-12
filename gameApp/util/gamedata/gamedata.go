package gamedata

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameData interface {
	FrameCount() int
	RenderHeight() int
	RenderWidth() int
	View() *ebiten.Image
	ViewScale() int
	ViewScaleInc(args *widget.ButtonClickedEventArgs)
	ViewScaleDec(args *widget.ButtonClickedEventArgs)
}
