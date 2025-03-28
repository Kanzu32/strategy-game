package ui

import (
	"image/color"
	_ "image/png"
	"regexp"
	"strategy-game/assets"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/turn/turnstate"
	"strategy-game/util/ecs"

	"strategy-game/util/ui/uistate"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	ui               ebitenui.UI
	textFace         *text.GoXFace
	sliceStandard    *image.NineSlice
	sliceIron        *image.NineSlice
	sliceWood        *image.NineSlice
	slicePaper       *image.NineSlice
	sliceIronLight   *image.NineSlice
	backButtonImage  *ebiten.Image
	plusButtonImage  *ebiten.Image
	minusButtonImage *ebiten.Image
}

func CreateUI() UI {
	f, _ := loadFont(36)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3.0, 3.0)

	img, _, err := ebitenutil.NewImageFromFile(assets.NineSliceStandard)
	if err != nil {
		panic(err)
	}
	newImg := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceStandard := image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceIron)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceIron := image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceWood)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceWood := image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSlicePaper)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	slicePaper := image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceIronLight)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceIronLight := image.NewNineSliceSimple(newImg, 3*3, 10*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.BackIcon)
	if err != nil {
		panic(err)
	}
	backButton := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	backButton.DrawImage(img, &opt)

	img, _, err = ebitenutil.NewImageFromFile(assets.PlusIcon)
	if err != nil {
		panic(err)
	}
	plusButton := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	plusButton.DrawImage(img, &opt)

	img, _, err = ebitenutil.NewImageFromFile(assets.MinusIcon)
	if err != nil {
		panic(err)
	}
	minusButton := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	minusButton.DrawImage(img, &opt)

	u := UI{ebitenui.UI{}, text.NewGoXFace(f), sliceStandard, sliceIron, sliceWood, slicePaper, sliceIronLight, backButton, plusButton, minusButton}

	return u
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.ui.Draw(screen)
}

func (u *UI) Update() {
	u.ui.Update()

	if singletons.AppState.UIState == uistate.Game {
		handleGameInput()
	}
}

func mousePosGameScale() (int, int) {
	x, y := ebiten.CursorPosition()
	return x / singletons.View.Scale, y / singletons.View.Scale
}

func handleGameInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		// ENT

		if singletons.Turn.State != turnstate.Input {
			return
		}

		// клик на активный (active) либо взятый в цель (target object) объект на экране
		activeEntities := ecs.PoolFilter([]ecs.AnyPool{pools.PositionPool, pools.SpritePool}, []ecs.AnyPool{})
		xPosGame, yPosGame := mousePosGameScale()
		for _, entity := range activeEntities {
			// неактивные объекты и объекты не взятые в цель игнорируются
			if !pools.ActiveFlag.HasEntity(entity) && !pools.TargetObjectFlag.HasEntity(entity) {
				continue
			}

			position, err := pools.PositionPool.Component(entity)
			if err != nil {
				panic(err)
			}

			sprite, err := pools.SpritePool.Component(entity)
			if err != nil {
				panic(err)
			}

			if position.X*16 < xPosGame && xPosGame < position.X*16+sprite.Sprite.Width() &&
				position.Y*16 < yPosGame && yPosGame < position.Y*16+sprite.Sprite.Height() {

				// объект, взятый в цель, явл. тайлом (выбрать объект в цель для действия)
				if pools.TargetObjectFlag.HasEntity(entity) && pools.TileFlag.HasEntity(entity) {
					singletons.Turn.State = turnstate.Action
					println("muvin")
					return
				}

				// активный объект не являющийся юнитом (выбрать объект в цель для действия)
				if pools.ActiveFlag.HasEntity(entity) && !pools.UnitFlag.HasEntity(entity) {
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetObjectFlag.AddExistingEntity(entity)
					return
				}

				// компонент team есть у всех юнитов (проверка на юнит выше)
				team, err := pools.TeamPool.Component(entity)
				if err != nil {
					panic(err)
				}

				// активный юнит игрока (выбрать его для управления)
				if pools.ActiveFlag.HasEntity(entity) && team.Team == singletons.Turn.PlayerTeam {
					for _, ent := range pools.TargetUnitFlag.Entities() {
						pools.TargetUnitFlag.RemoveEntity(ent)
					}
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetUnitFlag.AddExistingEntity(entity)
					return
				}

				// активный юнит оппонента (выбрать юнит в цель для действия)
				if pools.ActiveFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetObjectFlag.AddExistingEntity(entity)
					return
				}

				if pools.TargetObjectFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
					// атака...
					return
				}
			}
		}
	}
}

func (u *UI) ShowGameControls() {
	u.ui.Container = widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(2),
		widget.GridLayoutOpts.Spacing(0, 0),
		widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
	)))

	u.ui.Container.AddChild(widget.NewContainer())

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceIronLight, Pressed: u.sliceIronLight, Hover: u.sliceIronLight, Disabled: u.sliceIronLight}),
		widget.ButtonOpts.Graphic(u.plusButtonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			widget.WidgetOpts.MinSize(42, 42),
		),
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
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceIronLight, Pressed: u.sliceIronLight, Hover: u.sliceIronLight, Disabled: u.sliceIronLight}),
		widget.ButtonOpts.Graphic(u.minusButtonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			widget.WidgetOpts.MinSize(42, 42),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("dec")
			if singletons.View.Scale != 1 {
				singletons.View.Scale--
			}
		}),
	))

}

// func (u *UI) GetFocused() widget.Focuser {
// 	return u.ui.GetFocusedWidget()
// }

func (u *UI) ShowMainMenu() {

	u.ui.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(100)),
				widget.RowLayoutOpts.Spacing(20),
			),
		),
		// widget.ContainerOpts.BackgroundImage(u.slicePaper),
	)

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Play online", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
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
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Settings", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play online")
			// singletons.UIState = uistate.Game
			// singletons.GameMode = gamemode.Online
		}),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Login", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("login")
			singletons.AppState.UIState = uistate.Login
			singletons.AppState.StateChanged = true
		}),
	))
}

func (u *UI) ShowLogin() {
	// u.ui.Container.RemoveChildren()
	u.ui.Container = widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 0, Left: 100, Right: 100, Bottom: 0}),
				widget.RowLayoutOpts.Spacing(20),
			),
		),
		// widget.ContainerOpts.BackgroundImage(u.slicePaper),
	)

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(54, 54),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				// Stretch: true,
			}),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("main")
			singletons.AppState.UIState = uistate.Main
			singletons.AppState.StateChanged = true
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceIronLight, Pressed: u.sliceIronLight}),
		widget.ButtonOpts.Graphic(u.backButtonImage),
	))
	u.ui.Container.AddChild(widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 10),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				// Position: widget.RowLayoutPositionCenter,
				Stretch: true,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     u.sliceStandard,
			Disabled: u.sliceStandard,
		}),
		widget.TextInputOpts.Face(u.textFace),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(20)),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.Black,
			Disabled:      color.Black,
			Caret:         color.Black,
			DisabledCaret: color.Black,
		}),
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			res, err := regexp.Match(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, []byte(newInputText))
			if err != nil {
				panic(err)
			}
			return res, &newInputText
		}),
		widget.TextInputOpts.Placeholder("Email"),
		widget.TextInputOpts.IgnoreEmptySubmit(true),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			return
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(u.textFace, 5),
		),
	))

	u.ui.Container.AddChild(widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     u.sliceStandard,
			Disabled: u.sliceStandard,
		}),
		widget.TextInputOpts.Face(u.textFace),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(20)),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.Black,
			Disabled:      color.Black,
			Caret:         color.Black,
			DisabledCaret: color.Black,
		}),
		widget.TextInputOpts.Placeholder("Password"),
		widget.TextInputOpts.IgnoreEmptySubmit(true),
		widget.TextInputOpts.Secure(true),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			return
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(u.textFace, 5),
		),
	))

	innerContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(0, 0),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{false}),
			),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{})),
	)

	innerContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(180, 10)),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Login", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("login")
			// singletons.AppState.GameMode = gamemode.Local
			// singletons.AppState.UIState = uistate.Game
			// singletons.AppState.StateChanged = true
		}),
	))

	innerContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: u.sliceStandard, Pressed: u.sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(180, 10)),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Register", u.textFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("register")
			// singletons.AppState.GameMode = gamemode.Local
			// singletons.AppState.UIState = uistate.Game
			// singletons.AppState.StateChanged = true
		}),
	))

	u.ui.Container.AddChild(innerContainer)

}
