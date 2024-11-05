package components

import (
	"strategy-game/util/ecs"
	"strategy-game/util/material"
	"strategy-game/util/side"
	"strategy-game/util/sprite"

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
	Material material.Material
}

type Side struct {
	Side side.Side
}

type Occupied struct {
	ActiveObject *ecs.Entity
	UnitObject   *ecs.Entity
	StaticObject *ecs.Entity
}
