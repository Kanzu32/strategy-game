package uistate

type UIState uint8

//go:generate stringer -type=UIState
const (
	Main UIState = iota + 1
	Statistics
	Settings
	Game
	Login
	Results
)
