package singletons

import (
	"strategy-game/util/data/classes"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/stats"
	"strategy-game/util/data/turn"
	"strategy-game/util/network"
	"strategy-game/util/ui/uistate"

	"github.com/hajimehoshi/ebiten/v2"
)

var Turn turn.Turn

var ClassStats = map[classes.Class]stats.Stats{
	classes.Shield: {
		MaxEnergy:     5,
		EnergyPerTurn: 1,
		MoveCost:      1,
		AttackCost:    1,
		ActionCost:    1,
	},

	classes.Glaive: {
		MaxEnergy:     6,
		EnergyPerTurn: 1,
		MoveCost:      1,
		AttackCost:    1,
		ActionCost:    1,
	},

	classes.Bow: {
		MaxEnergy:     4,
		EnergyPerTurn: 1,
		MoveCost:      1,
		AttackCost:    1,
		ActionCost:    1,
	},

	classes.Knife: {
		MaxEnergy:     8,
		EnergyPerTurn: 1,
		MoveCost:      1,
		AttackCost:    1,
		ActionCost:    1,
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
	Image *ebiten.Image
	Scale int
}

var Settings struct {
	DefaultGameScale int
	UIScale          int
	Sound            int
}

var Connection network.ServerConnection

// var World *ecs.World

// var UI ui.UI
