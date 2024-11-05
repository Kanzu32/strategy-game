package ui

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/furex/v2"
)

func CreateGameUi(width int, height int) *furex.View {
	view := &furex.View{
		Width:        width,
		Height:       height,
		Direction:    furex.Column,
		Justify:      furex.JustifyEnd,
		AlignItems:   furex.AlignItemStretch,
		AlignContent: furex.AlignContentStretch,
		Wrap:         furex.NoWrap,
	}

	bottomMargin := 0
	BottomMenu := furex.View{
		Height:       100,
		Handler:      &BottomMenu{},
		Bottom:       &bottomMargin,
		Display:      furex.DisplayFlex,
		Direction:    furex.Row,
		Justify:      furex.JustifySpaceAround,
		AlignItems:   furex.AlignItemCenter,
		AlignContent: furex.AlignContentCenter,
		Wrap:         furex.NoWrap,
	}

	view.AddChild(&BottomMenu)

	f, err := os.Open("assets/ui/test-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	fImg, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	img := ebiten.NewImageFromImage(fImg)

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})
	return view
}

type BottomMenu struct{}

func (b *BottomMenu) Draw(screen *ebiten.Image, frame image.Rectangle, view *furex.View) {
	vector.DrawFilledRect(
		screen,
		float32(frame.Min.X),
		float32(frame.Min.Y),
		float32(frame.Size().X),
		float32(frame.Size().Y),
		color.White,
		false,
	)
}

type SpellBottom struct {
	Img *ebiten.Image
}

func (b *SpellBottom) Draw(screen *ebiten.Image, frame image.Rectangle, view *furex.View) {
	opt := ebiten.DrawImageOptions{}
	opt.Filter = ebiten.FilterNearest
	opt.GeoM.Scale(2, 2) // UI SCALE
	opt.GeoM.Translate(float64(frame.Min.X), float64(frame.Min.Y))
	screen.DrawImage(b.Img, &opt)
	// vector.DrawFilledRect(
	// 	screen,
	// 	float32(frame.Min.X),
	// 	float32(frame.Min.Y),
	// 	float32(frame.Size().X),
	// 	float32(frame.Size().Y),
	// 	color.Black,
	// 	false,
	// )
}
