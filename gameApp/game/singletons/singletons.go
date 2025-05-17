package singletons

import (
	"strategy-game/util/data/classes"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/stats"
	"strategy-game/util/data/turn"
	"strategy-game/util/data/userstatus"
	"strategy-game/util/ui/uistate"

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
		Attack:              2,
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
		Attack:              4,
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
		Attack:              2,
		AttackDistanceStart: 1,
		AttackDistanceEnd:   1,
	},
}

// var UIState uistate.UIState = uistate.Main

// var GameMode gamemode.GameMode = gamemode.Local

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
	DefaultGameScale int  `json:"DefaultGameScale"`
	Sound            int  `json:"Sound"`
	Fullscreen       bool `json:"Fullscreen"`
}

// var World *ecs.World

// var UI ui.UI
