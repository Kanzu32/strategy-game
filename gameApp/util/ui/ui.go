package ui

import (
	"image/color"
	_ "image/png"
	"strategy-game/assets"
	"strategy-game/game/singletons"
	"strategy-game/util/gamemode"
	"strategy-game/util/ui/uistate"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(assets.MonogramTTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

type UI struct {
	ui              ebitenui.UI
	textFace        *text.GoXFace
	uiSlice         *image.NineSlice
	backgroundSlice *image.NineSlice
}

func CreateUI() UI {
	img, _, err := ebitenutil.NewImageFromFile(assets.MainUIButton)
	if err != nil {
		panic(err)
	}

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3.0, 3.0)
	newImg := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)

	uiSlice := image.NewNineSliceSimple(newImg, 6*3, 4*3)

	backgroundImage, _, err := ebitenutil.NewImageFromFile(assets.MainMenuBackground)
	if err != nil {
		panic(err)
	}

	backgroundSlice := image.NewNineSliceSimple(backgroundImage, 10, 10)

	f, _ := loadFont(36)
	u := UI{ebitenui.UI{}, text.NewGoXFace(f), uiSlice, backgroundSlice}

	return u
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.ui.Draw(screen)
}

func (u *UI) Update() {
	u.ui.Update()
}

func (u *UI) ShowGameControls() {
	u.ui.Container = widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(2),
		widget.GridLayoutOpts.Spacing(0, 0),
		widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
	)))

	u.ui.Container.AddChild(widget.NewContainer())

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.uiSlice, Pressed: u.uiSlice, Hover: u.uiSlice, Disabled: u.uiSlice}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		// widget.ButtonOpts.ClickedHandler(g.ViewScaleInc),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("inc")
			if singletons.View.Scale != 10 {
				singletons.View.Scale++
			}
		}),
	))

	u.ui.Container.AddChild(widget.NewContainer())

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.uiSlice, Pressed: u.uiSlice}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("dec")
			if singletons.View.Scale != 1 {
				singletons.View.Scale--
			}
		}),
	))

}

func (u *UI) GetFocused() widget.Focuser {
	return u.ui.GetFocusedWidget()
}

func (u *UI) ShowMainMenu() {

	u.ui.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(100)),
				widget.RowLayoutOpts.Spacing(20),
			),
		),
		// widget.ContainerOpts.BackgroundImage(u.backgroundSlice),
	)

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.uiSlice, Pressed: u.uiSlice}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Play Online", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.uiSlice, Pressed: u.uiSlice}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Play Offline", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play offline")
			singletons.AppState.GameMode = gamemode.Local
			singletons.AppState.UIState = uistate.Game
			singletons.AppState.StateChanged = true
		}),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.uiSlice, Pressed: u.uiSlice}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Settings", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play online")
			// singletons.UIState = uistate.Game
			// singletons.GameMode = gamemode.Online
		}),
	))
}
