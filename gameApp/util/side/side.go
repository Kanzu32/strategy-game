package side

type Side uint8

const (
	Center Side = iota + 1
	LeftUp
	Left
	RightUp
	Right
	Up
	RightDown
	LeftDown
	Down
	LeftCorner
	RightCorner
)
