package ui

import (
	_ "image/png"
	"strategy-game/util/gamedata"
	"strategy-game/util/ui/mainuistate"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GAME UI
type GameUI struct {
	MenuBackground *ebiten.Image
	WalkButton     Button
	AttackButton   Button
	SkillButton    Button
	PlusButton     Button
	MinusButton    Button
	Portraits      []*ebiten.Image
	Skills         []*ebiten.Image
	CurrentScale   int
}

func CreateGameUI() *GameUI {
	ui := GameUI{}

	i := icon{}
	img, _, err := ebitenutil.NewImageFromFile("assets/ui/plus.png")
	if err != nil {
		panic(err)
	}
	i.Active = img

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/scale-inactive.png")
	if err != nil {
		panic(err)
	}
	i.Inactive = img
	ui.PlusButton = Button{Active: true, icon: i, Handler: plusHandler}

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/minus.png")
	if err != nil {
		panic(err)
	}
	i.Active = img

	ui.MinusButton = Button{Active: true, icon: i, Handler: minusHandler}

	// BACKGROUND
	img, _, err = ebitenutil.NewImageFromFile("assets/ui/unit-background.png")
	if err != nil {
		panic(err)
	}
	ui.MenuBackground = img

	// LOAD ICONS

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/walk-active-icon.png")
	if err != nil {
		panic(err)
	}
	i.Active = img

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/walk-inactive-icon.png")
	if err != nil {
		panic(err)
	}
	i.Inactive = img

	ui.WalkButton = Button{Active: true, icon: i}

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/attack-active-icon.png")
	if err != nil {
		panic(err)
	}
	i.Active = img

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/attack-inactive-icon.png")
	if err != nil {
		panic(err)
	}
	i.Inactive = img

	ui.AttackButton = Button{Active: true, icon: i}

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/skill-test-icon.png")
	if err != nil {
		panic(err)
	}
	i.Active = img

	// img, _, err = ebitenutil.NewImageFromFile("assets/ui/skill-test-icon.png")
	// if err != nil {
	// 	panic(err)
	// }
	// i.Inactive = img

	ui.SkillButton = Button{Active: true, icon: i}

	//PORTRAITS
	img, _, err = ebitenutil.NewImageFromFile("assets/ui/portraits/man-portrait.png")
	if err != nil {
		panic(err)
	}
	ui.Portraits = append(ui.Portraits, img)

	img, _, err = ebitenutil.NewImageFromFile("assets/ui/portraits/knight-portrait.png")
	if err != nil {
		panic(err)
	}
	ui.Portraits = append(ui.Portraits, img)

	return &ui
}

func (ui *GameUI) Draw(screen *ebiten.Image, g gamedata.GameData) {
	screenSize := screen.Bounds().Dx()
	menuSize := ui.MenuBackground.Bounds().Dx()
	scale := screenSize / (menuSize * 2)
	if scale < 2 {
		scale = 2
	}
	ui.CurrentScale = scale

	//unit menu
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(scale), float64(scale))
	opt.GeoM.Translate(float64(screenSize/2-menuSize*scale/2), float64(screen.Bounds().Dy()-ui.MenuBackground.Bounds().Dy()*scale))
	screen.DrawImage(ui.MenuBackground, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(scale), float64(scale))
	opt.GeoM.Translate(
		float64(screenSize/2-menuSize*scale/2),
		float64(screen.Bounds().Dy()-ui.MenuBackground.Bounds().Dy()*scale))
	screen.DrawImage(ui.Portraits[0], opt) //portrait

	// CONTROLL BUTTONS
	ui.DrawButton(
		&ui.WalkButton,
		screenSize/2-menuSize*scale/2+60*scale,
		screen.Bounds().Dy()-ui.MenuBackground.Bounds().Dy()*scale+9*scale,
		scale,
		screen,
	)

	ui.DrawButton(
		&ui.AttackButton,
		screenSize/2-menuSize*scale/2+92*scale,
		screen.Bounds().Dy()-ui.MenuBackground.Bounds().Dy()*scale+9*scale,
		scale,
		screen,
	)

	ui.DrawButton(
		&ui.SkillButton,
		screenSize/2-menuSize*scale/2+124*scale,
		screen.Bounds().Dy()-ui.MenuBackground.Bounds().Dy()*scale+9*scale,
		scale,
		screen,
	)

	// SCALE BUTTONS

	ui.DrawButton(
		&ui.PlusButton,
		screenSize-ui.PlusButton.Image().Bounds().Dx()*scale-3*scale,
		screen.Bounds().Dy()/2-ui.PlusButton.Image().Bounds().Dx()*scale-10*scale,
		scale,
		screen,
	)

	ui.DrawButton(
		&ui.MinusButton,
		screenSize-ui.MinusButton.Image().Bounds().Dx()*scale-3*scale,
		screen.Bounds().Dy()/2-ui.MinusButton.Image().Bounds().Dx()*scale+10*scale,
		scale,
		screen,
	)
}

func (ui *GameUI) DrawButton(button *Button, x int, y int, scale int, screen *ebiten.Image) {
	button.UnscaledX = x / scale
	button.UnscaledY = y / scale
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(scale), float64(scale))
	opt.GeoM.Translate(
		float64(x),
		float64(y))
	screen.DrawImage(button.Image(), opt)
}

// MAIN UI
type MainUI struct {
	State mainuistate.UIState
}

func (ui *MainUI) DrawButton(button *Button, x int, y int, scale int, screen *ebiten.Image) {
	button.UnscaledX = x / scale
	button.UnscaledY = y / scale
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(scale), float64(scale))
	opt.GeoM.Translate(
		float64(x),
		float64(y))
	screen.DrawImage(button.Image(), opt)
}

func (ui *MainUI) DrawTextInput(textInput *TextInput, x int, y int, scale int, screen *ebiten.Image) {
	textInput.UnscaledX = x / scale
	textInput.UnscaledY = y / scale
	text.Draw(screen, textInput.Text, textInput.Face, &textInput.Options)
}

// INPUT HANDLERS

type handler func(g gamedata.GameData)

func plusHandler(g gamedata.GameData) {
	g.ViewScaleInc()
}

func minusHandler(g gamedata.GameData) {
	g.ViewScaleDec()
}

// UI UTILS
type icon struct {
	Active   *ebiten.Image
	Inactive *ebiten.Image
}

type Button struct {
	Active    bool
	Handler   handler
	icon      icon
	UnscaledX int
	UnscaledY int
}

func (b *Button) Click(g gamedata.GameData) {
	b.Handler(g)
}

func (b *Button) Image() *ebiten.Image {
	if b.Active {
		return b.icon.Active
	}
	return b.icon.Inactive
}

func (b *Button) InBounds(x int, y int) bool {
	if b.UnscaledX <= x && x <= b.UnscaledX+b.Image().Bounds().Dx() &&
		b.UnscaledY <= y && y <= b.UnscaledY+b.Image().Bounds().Dy() {

		return true
	}
	return false
}

type TextInput struct {
	Text      string
	Face      text.Face
	Options   text.DrawOptions
	Focus     bool
	UnscaledX int
	UnscaledY int
	Height    int
	Width     int
}

func (b *TextInput) InBounds(x int, y int) bool {
	if b.UnscaledX <= x && x <= b.UnscaledX+b.Width &&
		b.UnscaledY <= y && y <= b.UnscaledY+b.Height {

		return true
	}
	return false
}
