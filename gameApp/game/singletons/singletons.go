package singletons

import (
	"strategy-game/util/data/classes"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/stats"
	"strategy-game/util/data/turn"
	"strategy-game/util/data/userstatus"
	"strategy-game/util/ui/uistate"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var Turn turn.Turn

var ClassStats = map[classes.Class]stats.Stats{
	classes.Shield: {
		MaxEnergy:           5,
		EnergyPerTurn:       1,
		MoveCost:            1,
		AttackCost:          1,
		ActionCost:          1,
		MaxHealth:           25,
		Attack:              100, //2
		AttackDistanceStart: 1,
		AttackDistanceEnd:   1,
	},

	classes.Glaive: {
		MaxEnergy:           6,
		EnergyPerTurn:       1,
		MoveCost:            1,
		AttackCost:          1,
		ActionCost:          1,
		MaxHealth:           20,
		Attack:              100, //4
		AttackDistanceStart: 2,
		AttackDistanceEnd:   2.25,
	},

	classes.Bow: {
		MaxEnergy:           4,
		EnergyPerTurn:       1,
		MoveCost:            1,
		AttackCost:          1,
		ActionCost:          1,
		MaxHealth:           10,
		Attack:              4,
		AttackDistanceStart: 2,
		AttackDistanceEnd:   2.25,
	},

	classes.Knife: {
		MaxEnergy:           8,
		EnergyPerTurn:       1,
		MoveCost:            1,
		AttackCost:          1,
		ActionCost:          1,
		MaxHealth:           15,
		Attack:              100, //2
		AttackDistanceStart: 1,
		AttackDistanceEnd:   1,
	},
}

var AppState struct {
	UIState      uistate.UIState
	GameMode     gamemode.GameMode
	StateChanged bool
}

var FrameCount int = 0

var Render struct {
	Width  int
	Height int
}

var View struct {
	Image  *ebiten.Image
	Scale  int
	ShiftX int
	ShiftY int
}

var UserLogin struct {
	Email    string
	Password string
	Status   userstatus.UserStatus
}

var Settings struct {
	DefaultGameScale int    `json:"DefaultGameScale"`
	Sound            int    `json:"Sound"`
	Fullscreen       bool   `json:"Fullscreen"`
	Language         string `json:"Language"`
}

var MapSize struct {
	Width  int
	Height int
}

type textLines struct {
	On               string
	Off              string
	WinOnline        string
	LoseOnline       string
	WinBlue          string
	WinRed           string
	ToMainMenu       string
	PlayOnline       string
	PlayOffline      string
	Settings         string
	Login            string
	Register         string
	Statistics       string
	Email            string
	Password         string
	Language         string
	GameDefaultScale string
	Sound            string
	Fullscreen       string
	OnlineStatus     string
	OfflineStatus    string
	LoginError       string
	ConnectionError  string
	RegisterError    string
	EmailError       string
}

var LanguageText = map[string]textLines{
	"Eng": {
		On:               "On",
		Off:              "Off",
		WinOnline:        "You win",
		LoseOnline:       "You lose",
		WinBlue:          "Blue team win",
		WinRed:           "Red team win",
		ToMainMenu:       "Go to main menu",
		PlayOnline:       "Play Online",
		PlayOffline:      "Play Offline",
		Settings:         "Settings",
		Login:            "Login",
		Register:         "Register",
		Statistics:       "Statistics",
		Email:            "Email",
		Password:         "Password",
		Language:         "Language",
		GameDefaultScale: "Default Game Scale",
		Sound:            "Sound",
		Fullscreen:       "Fullscreen",
		OnlineStatus:     "Online: ",
		OfflineStatus:    "Offline",
		LoginError:       "Wrong email or password",
		ConnectionError:  "Error connecting to server",
		RegisterError:    "User is already registered",
		EmailError:       "Wrong email format",
	},
	"Rus": {
		On:               "Вкл",
		Off:              "Выкл",
		WinOnline:        "Победа",
		LoseOnline:       "Поражение",
		WinBlue:          "Синяя команда победила",
		WinRed:           "Красная команда победила",
		ToMainMenu:       "В главное меню",
		PlayOnline:       "Играть по сети",
		PlayOffline:      "Играть на одном компьютере",
		Settings:         "Настройки",
		Login:            "Вход",
		Register:         "Регистрация",
		Statistics:       "Статистика",
		Email:            "Электронная почта",
		Password:         "Пароль",
		Language:         "Язык",
		GameDefaultScale: "Масштабирование игрового поля",
		Sound:            "Громкость",
		Fullscreen:       "Полноэкранный режим",
		OnlineStatus:     "Вход произведён: ",
		OfflineStatus:    "Вход не произведён",
		LoginError:       "Неверная почта или пароль",
		ConnectionError:  "Ошибка при соединении с сервером",
		RegisterError:    "Пользователь уже зарегистрирован",
		EmailError:       "Неверный формат элекетронной почты",
	},
}

var RawMap string
var MapMutex sync.Mutex

// var World *ecs.World

// var UI ui.UI
