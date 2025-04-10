package classes

type Class uint8

//go:generate stringer -type=Class
const (
	Shield Class = iota + 1
	Glaive
	Knife
	Bow
)
