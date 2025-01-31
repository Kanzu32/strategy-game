package turn

import "strategy-game/util/teams"

type Turn struct {
	CurrentTurn teams.Team
	PlayerTeam  teams.Team
}
