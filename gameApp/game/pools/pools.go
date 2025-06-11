package pools

import (
	"crypto/sha256"
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
var HealthPool *ecs.ComponentPool[c.Health]
var TweenPool *ecs.ComponentPool[c.Tween]
var MovePool *ecs.ComponentPool[c.MoveDirection]
var DirectionPool *ecs.ComponentPool[c.Direction]
var AttackPool *ecs.ComponentPool[c.Attack]
var DamagePool *ecs.ComponentPool[c.Damage]

var TileFlag *ecs.FlagPool
var WallFlag *ecs.FlagPool
var SoftFlag *ecs.FlagPool
var UnitFlag *ecs.FlagPool
var GhostFlag *ecs.FlagPool
var ActiveFlag *ecs.FlagPool
var TargetUnitFlag *ecs.FlagPool
var TargetObjectFlag *ecs.FlagPool
var DeadFlag *ecs.FlagPool

func CalcHash() [32]byte {
	data := TileFlag.String() + WallFlag.String() + SoftFlag.String() + UnitFlag.String() + DeadFlag.String() + PositionPool.String() + SidePool.String() +
		OccupiedPool.String() + TeamPool.String() + ClassPool.String() + EnergyPool.String() + HealthPool.String() + DirectionPool.String()

	return sha256.Sum256([]byte(data))
}
