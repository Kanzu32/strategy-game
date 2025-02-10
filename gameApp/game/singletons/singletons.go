package singletons

import (
	"strategy-game/util/classes"
	"strategy-game/util/stats"
	"strategy-game/util/turn"
	"strategy-game/util/ui/uistate"
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

var UIState uistate.UIState = uistate.Main
