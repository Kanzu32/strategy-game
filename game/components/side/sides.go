package sides

type Sides uint8

const (
	Center Sides = iota + 1
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
