package comp

import (
	"strategy-game/components/material"
	"strategy-game/ecs"
	"strategy-game/sprite"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X int
	Y int
}

type ImageRender struct {
	Options ebiten.DrawImageOptions
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

type View struct {
	Img *ebiten.Image
}
