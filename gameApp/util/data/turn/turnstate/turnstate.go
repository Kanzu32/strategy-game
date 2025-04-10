package turnstate

type TurnState uint8

//go:generate stringer -type=TurnState
const (
	Input TurnState = iota + 1
	Action
	Wait
)
