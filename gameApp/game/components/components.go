package components

import (
	"strategy-game/util/data/classes"
	"strategy-game/util/data/material"
	"strategy-game/util/data/side"
	"strategy-game/util/data/sprite"
	"strategy-game/util/data/teams"
	"strategy-game/util/data/tween"
	"strategy-game/util/ecs"

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

type Team struct {
	Team teams.Team
}

type Class struct {
	Class classes.Class
}

type Energy struct {
	Energy uint8
}

type Tween struct {
	Animation tween.TweenAnimation
}

type MoveDirection struct {
	X int8
	Y int8
}

// type StandOn struct {
// 	Tile *ecs.Entity
// }
