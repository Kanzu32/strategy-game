package teams

type Team uint8

//go:generate stringer -type=Team
const (
	Blue Team = iota + 1
	Red
)
