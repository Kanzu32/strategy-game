package comp

import (
	materials "strategy-game/components/material"
	sides "strategy-game/components/side"
	"strategy-game/ecs"
	sprites "strategy-game/sprite"

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
	Sprite sprites.Sprite
}

type Material struct {
	Material materials.Materials
}

type Side struct {
	Side sides.Sides
}

type Occupied struct {
	ActiveObject *ecs.Entity
	UnitObject   *ecs.Entity
	StaticObject *ecs.Entity
}

type View struct {
	Img *ebiten.Image
}
