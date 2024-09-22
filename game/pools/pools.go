package pools

import (
	c "strategy-game/components"
	"strategy-game/ecs"
)

var PositionPool *ecs.ComponentPool[c.Position]
var SpritePool *ecs.ComponentPool[c.Sprite]
var MaterialPool *ecs.ComponentPool[c.Material]
var ViewPool *ecs.ComponentPool[c.View]
var ImageRenderPool *ecs.ComponentPool[c.ImageRender]

var TileFlag *ecs.FlagPool
