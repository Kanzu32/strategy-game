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
	"strategy-game/util/network"

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

var TextFace *text.GoXFace
var sliceStandard *image.NineSlice
var sliceIron *image.NineSlice
var sliceWood *image.NineSlice
var slicePaper *image.NineSlice
var sliceIronLight *image.NineSlice
var backButtonImage *ebiten.Image
var plusButtonImage *ebiten.Image
var minusButtonImage *ebiten.Image
var skipButtonImage *ebiten.Image

type UI struct {
	ui ebitenui.UI
	// textFace         *text.GoXFace
	// sliceStandard    *image.NineSlice
	// sliceIron        *image.NineSlice
	// sliceWood        *image.NineSlice
	// slicePaper       *image.NineSlice
	// sliceIronLight   *image.NineSlice
	// backButtonImage  *ebiten.Image
	// plusButtonImage  *ebiten.Image
	// minusButtonImage *ebiten.Image
}

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

func CreateUI() UI {
	f, _ := loadFont(36)
	TextFace = text.NewGoXFace(f)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Scale(3.0, 3.0)

	img, _, err := ebitenutil.NewImageFromFile(assets.NineSliceStandard)
	if err != nil {
		panic(err)
	}
	newImg := ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceStandard = image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceIron)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceIron = image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceWood)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceWood = image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSlicePaper)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	slicePaper = image.NewNineSliceSimple(newImg, 6*3, 4*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.NineSliceIronLight)
	if err != nil {
		panic(err)
	}
	newImg = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	newImg.DrawImage(img, &opt)
	sliceIronLight = image.NewNineSliceSimple(newImg, 3*3, 10*3)

	img, _, err = ebitenutil.NewImageFromFile(assets.BackIcon)
	if err != nil {
		panic(err)
	}
	backButtonImage = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	backButtonImage.DrawImage(img, &opt)

	img, _, err = ebitenutil.NewImageFromFile(assets.PlusIcon)
	if err != nil {
		panic(err)
	}
	plusButtonImage = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	plusButtonImage.DrawImage(img, &opt)

	img, _, err = ebitenutil.NewImageFromFile(assets.MinusIcon)
	if err != nil {
		panic(err)
	}
	minusButtonImage = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	minusButtonImage.DrawImage(img, &opt)

	img, _, err = ebitenutil.NewImageFromFile(assets.SkipIcon)
	if err != nil {
		panic(err)
	}
	skipButtonImage = ebiten.NewImage(img.Bounds().Dx()*3, img.Bounds().Dy()*3)
	skipButtonImage.DrawImage(img, &opt)

	// u := UI{ebitenui.UI{}, text.NewGoXFace(f), sliceStandard, sliceIron, sliceWood, slicePaper, sliceIronLight, backButton, plusButton, minusButton}
	u := UI{ebitenui.UI{}}

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
	return (x - singletons.View.ShiftX) / singletons.View.Scale, (y - singletons.View.ShiftY) / singletons.View.Scale
}

var isDragging bool = false
var lastPositionX int = 0
var lastPositionY int = 0

func handleGameInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		isDragging = true
		lastPositionX, lastPositionY = ebiten.CursorPosition()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		isDragging = false
	}

	if isDragging {
		x, y := ebiten.CursorPosition()
		singletons.View.ShiftX -= lastPositionX - x
		singletons.View.ShiftY -= lastPositionY - y
		lastPositionX = x
		lastPositionY = y
		println(singletons.View.ShiftX, singletons.View.ShiftY)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && !isDragging {

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
					println("action")
					return
				}

				// активный объект не являющийся юнитом (выбрать объект в цель для действия)
				if pools.ActiveFlag.HasEntity(entity) && !pools.UnitFlag.HasEntity(entity) {
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetObjectFlag.AddExistingEntity(entity)
					println("sht its hapening 1")
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

				// // активный юнит оппонента (выбрать юнит в цель для действия)
				// if pools.ActiveFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
				// 	for _, ent := range pools.TargetObjectFlag.Entities() {
				// 		pools.TargetObjectFlag.RemoveEntity(ent)
				// 	}
				// 	pools.TargetObjectFlag.AddExistingEntity(entity)
				// 	println("sht its hapening 2")
				// 	return
				// }

				// if pools.TargetObjectFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
				// 	// атака...
				// 	println("sht its hapening 3")
				// 	return
				// }
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

	u.ui.Container.AddChild(widget.NewButton( // TODO skip button
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceIronLight, Pressed: sliceIronLight, Hover: sliceIronLight, Disabled: sliceIronLight}),
		widget.ButtonOpts.Graphic(skipButtonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			widget.WidgetOpts.MinSize(42, 42),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("skip")

			if singletons.Turn.State == turnstate.Input {
				if singletons.AppState.GameMode == gamemode.Online {
					network.SendSkip()
				}
				singletons.Turn.IsTurnEnds = true
			}
		}),
	))

	// spacer
	u.ui.Container.AddChild(widget.NewContainer())
	u.ui.Container.AddChild(widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(42, 42),
		),
	))

	u.ui.Container.AddChild(widget.NewContainer())

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceIronLight, Pressed: sliceIronLight, Hover: sliceIronLight, Disabled: sliceIronLight}),
		widget.ButtonOpts.Graphic(plusButtonImage),
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
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceIronLight, Pressed: sliceIronLight, Hover: sliceIronLight, Disabled: sliceIronLight}),
		widget.ButtonOpts.Graphic(minusButtonImage),
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
		widget.ContainerOpts.BackgroundImage(sliceIron),
	)

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Play online", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play online")
			singletons.AppState.GameMode = gamemode.Online
			singletons.AppState.UIState = uistate.Game
			singletons.AppState.StateChanged = true
		}),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Play Offline", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play offline")
			singletons.AppState.GameMode = gamemode.Local
			singletons.AppState.UIState = uistate.Game
			singletons.AppState.StateChanged = true
		}),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Settings", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("play online")
			// singletons.UIState = uistate.Game
			// singletons.GameMode = gamemode.Online
		}),
	))

	u.ui.Container.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Login", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
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
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 20, Left: 100, Right: 100, Bottom: 0}),
				widget.RowLayoutOpts.Spacing(20),
			),
		),
		widget.ContainerOpts.BackgroundImage(sliceIron),
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
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceIronLight, Pressed: sliceIronLight}),
		widget.ButtonOpts.Graphic(backButtonImage),
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
			Idle:     sliceStandard,
			Disabled: sliceStandard,
		}),
		widget.TextInputOpts.Face(TextFace),
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
			widget.CaretOpts.Size(TextFace, 5),
		),
	))

	u.ui.Container.AddChild(widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     sliceStandard,
			Disabled: sliceStandard,
		}),
		widget.TextInputOpts.Face(TextFace),
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
			widget.CaretOpts.Size(TextFace, 5),
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
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(180, 10)),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Login", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("login")

			emailInput := u.ui.Container.Children()[1].(*widget.TextInput)
			passwordInput := u.ui.Container.Children()[2].(*widget.TextInput)

			status := network.LoginRequest(emailInput.GetText(), passwordInput.GetText())
			statusText := u.ui.Container.Children()[len(u.ui.Container.Children())-1].(*widget.Text)

			switch status {
			case 200:
				statusText.Color = color.RGBA{0, 255, 0, 0}
				statusText.Label = "Успешный вход"
			case 401:
				statusText.Color = color.RGBA{255, 0, 0, 0}
				statusText.Label = "Неверная почта или пароль"
			default:
				statusText.Color = color.RGBA{255, 0, 0, 0}
				statusText.Label = "Ошибка при входе"
			}
			// singletons.AppState.GameMode = gamemode.Local
			// singletons.AppState.UIState = uistate.Game
			// singletons.AppState.StateChanged = true
		}),
	))

	innerContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{Idle: sliceStandard, Pressed: sliceStandard}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(180, 10)),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Text("Register", TextFace, &widget.ButtonTextColor{Idle: color.Black, Pressed: color.Black}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(20)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("register")
			emailInput := u.ui.Container.Children()[1].(*widget.TextInput)
			passwordInput := u.ui.Container.Children()[2].(*widget.TextInput)

			status := network.RegisterRequest(emailInput.GetText(), passwordInput.GetText())
			statusText := u.ui.Container.Children()[len(u.ui.Container.Children())-1].(*widget.Text)

			switch status {
			case 200:
				statusText.Color = color.RGBA{0, 255, 0, 0}
				statusText.Label = "Успешная регистрация и вход"
			case 409:
				statusText.Color = color.RGBA{255, 0, 0, 0}
				statusText.Label = "Пользователь уже зарегистрирован"
			default:
				statusText.Color = color.RGBA{255, 0, 0, 0}
				statusText.Label = "Ошибка при регистрации"
			}
			// singletons.AppState.GameMode = gamemode.Local
			// singletons.AppState.UIState = uistate.Game
			// singletons.AppState.StateChanged = true
		}),
	))

	u.ui.Container.AddChild(innerContainer)

	u.ui.Container.AddChild(widget.NewText(widget.TextOpts.Text("", TextFace, color.Black)))
}
