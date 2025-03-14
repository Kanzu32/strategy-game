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
var TeamPool *ecs.ComponentPool[c.Team]
var ClassPool *ecs.ComponentPool[c.Class]
var EnergyPool *ecs.ComponentPool[c.Energy]
var TweenPool *ecs.ComponentPool[c.Tween]
var MovePool *ecs.ComponentPool[c.MoveDireaction]

// var StandOnPool *ecs.ComponentPool[c.StandOn]

var TileFlag *ecs.FlagPool
var WallFlag *ecs.FlagPool
var SoftFlag *ecs.FlagPool
var UnitFlag *ecs.FlagPool
var GhostFlag *ecs.FlagPool
var ActiveFlag *ecs.FlagPool
var TargetUnitFlag *ecs.FlagPool
var TargetObjectFlag *ecs.FlagPool
