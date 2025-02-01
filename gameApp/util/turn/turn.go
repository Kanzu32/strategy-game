package turn

import (
	"strategy-game/util/teams"
	"strategy-game/util/turnstate"
)

type Turn struct {
	CurrentTurn teams.Team
	PlayerTeam  teams.Team
	State       turnstate.TurnState
}
