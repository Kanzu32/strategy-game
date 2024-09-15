package comp

import (
	"strategy-game/components/material"
	"strategy-game/ecs"
	"strategy-game/sprite"
)

type Position struct {
	X int
	Y int
}

type ScreenRender struct {
	X int
	Y int
}

type Sprite struct {
	Sprite sprite.Sprite
}

type Material struct {
	Material material.Materials
}

type OccupiedBy struct {
	Entity ecs.Entity
}
