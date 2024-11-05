package pools

import (
	c "strategy-game/game/components"
	"strategy-game/util/ecs"
)

var PositionPool *ecs.ComponentPool[c.Position]
var SpritePool *ecs.ComponentPool[c.Sprite]
var MaterialPool *ecs.ComponentPool[c.Material]
var SidePool *ecs.ComponentPool[c.Side]
var ImageRenderPool *ecs.ComponentPool[c.ImageRender]
var OccupiedPool *ecs.ComponentPool[c.Occupied]

var SolidFlag *ecs.FlagPool
var TileFlag *ecs.FlagPool
