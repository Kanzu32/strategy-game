package gamemode

type GameMode uint8

const (
	Local GameMode = iota + 1
	Online
)
