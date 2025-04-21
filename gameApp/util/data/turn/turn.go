package turn

import (
	"strategy-game/util/data/teams"
	"strategy-game/util/data/turn/turnstate"
)

type Turn struct {
	CurrentTurn teams.Team
	PlayerTeam  teams.Team
	State       turnstate.TurnState
	IsTurnEnds  bool
}
