package directions

type Direction uint8

//go:generate stringer -type=Direction
const (
	Down Direction = iota + 1
	Right
	Left
	Up
)
