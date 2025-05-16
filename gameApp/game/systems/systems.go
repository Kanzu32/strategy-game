package systems

import (
	"log"
	"math"
	"strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/util/data/classes"
	"strategy-game/util/data/damagetype"
	"strategy-game/util/data/directions"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/teams"
	"strategy-game/util/data/turn/turnstate"
	"strategy-game/util/data/tween"
	"strategy-game/util/data/tween/tweentype"
	"strategy-game/util/ecs"
	"strategy-game/util/network"
	"strategy-game/util/ui"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ###
// LOGIC SYSTEMS
// ###

type TurnSystem struct{}

func (s *TurnSystem) Run() { // highlight active units
	if singletons.Turn.State == turnstate.Action {
		for _, ent := range pools.ActiveFlag.Entities() {
			pools.ActiveFlag.RemoveEntity(ent)
			println("FUCK 1")
		}
	}
}

type MarkActiveUnitsSystem struct{}

func (s *MarkActiveUnitsSystem) Run() { // highlight active units
	if singletons.Turn.State != turnstate.Input || pools.TargetUnitFlag.EntityCount() > 0 { // если игрок уже походил юнитом, то другого выбирать не надо
		return
	}

	if singletons.Turn.CurrentTurn == singletons.Turn.PlayerTeam {
		entities := ecs.PoolFilter([]ecs.AnyPool{pools.TeamPool, pools.EnergyPool}, []ecs.AnyPool{pools.ActiveFlag, pools.DeadFlag}) // all inactive units
		for _, entity := range entities {
			energyComp, err := pools.EnergyPool.Component(entity)
			if err != nil {
				panic(err)
			}

			teamComp, err := pools.TeamPool.Component(entity)
			if err != nil {
				panic(err)
			}

			if energyComp.Energy > 0 && teamComp.Team == singletons.Turn.PlayerTeam {
				println("active unit")
				pools.ActiveFlag.AddExistingEntity(entity) // highlight units
			}
		}
	}
}

type MarkActiveTilesSystem struct{}

func (s *MarkActiveTilesSystem) Run() {
	if singletons.Turn.State != turnstate.Input {
		return
	}

	// Смотрим какой юнит взят в таргет
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag, pools.EnergyPool, pools.ClassPool, pools.PositionPool}, []ecs.AnyPool{})
	if len(units) > 1 {
		panic("More than one targeted units")
	}

	// Ищим тайлы доступлые для клика
	tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.PositionPool, pools.OccupiedPool}, []ecs.AnyPool{})

	for _, unit := range units {

		unitPosition, err := pools.PositionPool.Component(unit)
		if err != nil {
			panic(err)
		}

		unitEnergy, err := pools.EnergyPool.Component(unit)
		if err != nil {
			panic(err)
		}

		class, err := pools.ClassPool.Component(unit)
		if err != nil {
			panic(err)
		}

		for _, tile := range tiles {

			if pools.ActiveFlag.HasEntity(tile) {
				pools.ActiveFlag.RemoveEntity(tile)
			}

			tilePostion, err := pools.PositionPool.Component(tile)
			if err != nil {
				panic(err)
			}

			occupied, err := pools.OccupiedPool.Component(tile)
			if err != nil {
				panic(err)
			}

			distance := positionsDistance(unitPosition, tilePostion)

			if distance == 1 && // Проверка на передвижение
				!pools.WallFlag.HasEntity(tile) &&
				occupied.UnitObject == nil &&
				(occupied.ActiveObject == nil || pools.SoftFlag.HasEntity(*occupied.ActiveObject)) &&
				unitEnergy.Energy >= singletons.ClassStats[class.Class].MoveCost {

				pools.ActiveFlag.AddExistingEntity(tile)
			} else if singletons.Turn.IsAttackAllowed &&
				distance >= singletons.ClassStats[class.Class].AttackDistanceStart && // проверка на атаку
				distance <= singletons.ClassStats[class.Class].AttackDistanceEnd &&
				!pools.WallFlag.HasEntity(tile) &&
				occupied.UnitObject != nil &&
				!pools.DeadFlag.HasEntity(*occupied.UnitObject) &&
				(occupied.ActiveObject == nil || pools.SoftFlag.HasEntity(*occupied.ActiveObject)) &&
				unitEnergy.Energy >= singletons.ClassStats[class.Class].AttackCost {

				// Проверка на то, что в цель взят не союзник
				team, err := pools.TeamPool.Component(*occupied.UnitObject)
				if err != nil {
					panic(err)
				}
				if team.Team != singletons.Turn.PlayerTeam {
					pools.ActiveFlag.AddExistingEntity(tile)
				}
			}
			// if distance == 1 && // Проверка на передвижение
			// 	!pools.ActiveFlag.HasEntity(tile) &&
			// 	!pools.WallFlag.HasEntity(tile) &&
			// 	occupied.UnitObject == nil &&
			// 	(occupied.ActiveObject == nil || pools.SoftFlag.HasEntity(*occupied.ActiveObject)) &&
			// 	unitEnergy.Energy >= singletons.ClassStats[class.Class].MoveCost {

			// 	pools.ActiveFlag.AddExistingEntity(tile)
			// } else if singletons.Turn.IsAttackAllowed &&
			// 	distance >= singletons.ClassStats[class.Class].AttackDistanceStart && // проверка на атаку
			// 	distance <= singletons.ClassStats[class.Class].AttackDistanceEnd &&
			// 	!pools.ActiveFlag.HasEntity(tile) &&
			// 	!pools.WallFlag.HasEntity(tile) &&
			// 	occupied.UnitObject != nil &&
			// 	!pools.DeadFlag.HasEntity(*occupied.UnitObject) &&
			// 	(occupied.ActiveObject == nil || pools.SoftFlag.HasEntity(*occupied.ActiveObject)) &&
			// 	unitEnergy.Energy >= singletons.ClassStats[class.Class].AttackCost {

			// 	// Проверка на то, что в цель взят не союзник
			// 	team, err := pools.TeamPool.Component(*occupied.UnitObject)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	if team.Team != singletons.Turn.PlayerTeam {
			// 		pools.ActiveFlag.AddExistingEntity(tile)
			// 	}
			// } else if class.Class != classes.Glaive && distance != 1 && pools.ActiveFlag.HasEntity(tile) {

			// 	pools.ActiveFlag.RemoveEntity(tile)
			// 	println("FUCK 2")

			// } else if class.Class == classes.Glaive && (distance != 1 ||
			// 	distance >= singletons.ClassStats[class.Class].AttackDistanceStart &&
			// 		distance <= singletons.ClassStats[class.Class].AttackDistanceEnd &&
			// 		occupied.UnitObject == nil) &&
			// 		pools.ActiveFlag.HasEntity(tile) {

			// 	pools.ActiveFlag.RemoveEntity(tile)
			// 	println("FUCK 3")
			// }

			//  else if distance != 1 &&
			// 	(distance < singletons.ClassStats[class.Class].AttackDistanceStart || distance > singletons.ClassStats[class.Class].AttackDistanceEnd) &&
			// 	pools.ActiveFlag.HasEntity(tile) {

			// 	pools.ActiveFlag.RemoveEntity(tile)
			// }
			// else if positionsDistance(unitPosition, tilePostion) != 1 && pools.ActiveFlag.HasEntity(tile) {
			// 	pools.ActiveFlag.RemoveEntity(tile)
			// }
		}
	}
}

type NetworkSystem struct{}

func (s *NetworkSystem) Run() {
	if singletons.Turn.State != turnstate.Action {
		return
	}

	if singletons.AppState.GameMode == gamemode.Online {
		units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag}, []ecs.AnyPool{})
		if len(units) > 1 {
			panic("More than one targeted units")
		} else if len(units) == 0 {
			panic("Zero targeted units")
		}
		unit := units[0]

		// get targeted object
		tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.TargetObjectFlag}, []ecs.AnyPool{})
		if len(tiles) > 1 {
			panic("More than one targeted objects")
		} else if len(tiles) == 0 {
			return
		}
		tile := tiles[0]

		network.SendGameData(unit, tile)
	}

}

type MoveSystem struct{}

func (s *MoveSystem) Run() {
	// Если пользователь не смотрит анимацию (не Action или Wait)
	if singletons.Turn.State == turnstate.Input {
		return
	}

	// #1# Проверка на необходимость начать новую анимацию передвижения

	// если есть юнит взятый в цель, который не двигается...
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag}, []ecs.AnyPool{pools.MovePool, pools.TweenPool})

	// ... и тайл взятый в цель
	tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.TargetObjectFlag}, []ecs.AnyPool{})

	if len(units) > 1 { // если их больше чем 1 всё плохо
		panic("More than one targeted units")
	} else if len(tiles) > 1 {
		panic("More than one targeted objects")
	} else if len(units) == 1 && len(tiles) == 1 { // если их по одному, то стартуем новую анимацию ИНАЧЕ идём проверять #2#
		unit := units[0]
		tile := tiles[0]

		occupied, err := pools.OccupiedPool.Component(tile)
		if err != nil {
			panic(err)
		}

		// Пропускаем если тайл занят (это атака или взаимодействие)
		if occupied.UnitObject != nil || occupied.StaticObject != nil || (occupied.ActiveObject != nil && !pools.SoftFlag.HasEntity(*occupied.ActiveObject)) {
			return
		}

		println("now we muvin")
		unitPos, err := pools.PositionPool.Component(unit)
		if err != nil {
			panic(err)
		}

		tilePos, err := pools.PositionPool.Component(tile)
		if err != nil {
			panic(err)
		}

		println("tween")
		pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Linear, 500, (tilePos.X-unitPos.X)*16, (tilePos.Y-unitPos.Y)*16, 0)})

		sprite, err := pools.SpritePool.Component(unit)
		if err != nil {
			panic(err)
		}

		for _, ent := range tiles {
			pools.TargetObjectFlag.RemoveEntity(ent)
		}
		println("move")
		c := components.MoveDirection{X: int8(tilePos.X - unitPos.X), Y: int8(tilePos.Y - unitPos.Y)}
		pools.MovePool.AddExistingEntity(unit, c)

		dir, err := pools.DirectionPool.Component(unit)
		if err != nil {
			panic(err)
		}
		if c.X > 0 && c.Y == 0 {
			dir.Direction = directions.Right
			sprite.Sprite.SetAnimation("walk-right")
		} else if c.X < 0 && c.Y == 0 {
			dir.Direction = directions.Left
			sprite.Sprite.SetAnimation("walk-left")
		} else if c.X == 0 && c.Y < 0 {
			dir.Direction = directions.Up
			sprite.Sprite.SetAnimation("walk-up")
		} else {
			dir.Direction = directions.Down
			sprite.Sprite.SetAnimation("walk-down")
		}

		return // завершаем т.к. нет необходимости проверять завершение ходьбы на этом кадре
	}

	// #2# Проверка условий завершения анимации передвижения

	// Проверим есть ли юниты передвигающиеся в данный момент
	units = ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag, pools.TweenPool, pools.MovePool}, []ecs.AnyPool{})

	if len(units) > 1 { // Если взято в цель больше одного юнита -> всё плохо
		panic("More than one targeted units")
	} else if len(units) == 1 { // Нормальное состояние
		unit := units[0]

		if !pools.TweenPool.HasEntity(unit) { //useless?
			return
		}

		t, err := pools.TweenPool.Component(unit)
		if err != nil {
			panic(err)
		}

		// Если закончилась анимация ходьбы анимацию нужно убрать и изменить позицию юнита
		if t.Animation.IsEnded() {
			unitPos, err := pools.PositionPool.Component(unit)
			if err != nil {
				panic(err)
			}

			move, err := pools.MovePool.Component(unit)
			if err != nil {
				panic(err)
			}

			for _, entity := range pools.OccupiedPool.Entities() {
				occupied, err := pools.OccupiedPool.Component(entity)
				if err != nil {
					panic(err)
				}

				pos, err := pools.PositionPool.Component(entity)
				if err != nil {
					panic(err)
				}

				if occupied.UnitObject != nil && occupied.UnitObject.Equals(unit) && pos.X == unitPos.X && pos.Y == unitPos.Y {
					println("remove from ", pos.X, " ", pos.Y)
					occupied.UnitObject = nil
				}

				if occupied.UnitObject == nil && pos.X == unitPos.X+int(move.X) && pos.Y == unitPos.Y+int(move.Y) {
					println("add to ", pos.X, " ", pos.Y)
					occupied.UnitObject = &unit
				}
			}

			unitPos.X += int(move.X)
			unitPos.Y += int(move.Y)

			pools.TweenPool.RemoveEntity(unit)
			pools.MovePool.RemoveEntity(unit)

			sprite, err := pools.SpritePool.Component(unit)
			if err != nil {
				panic(err)
			}

			dir, err := pools.DirectionPool.Component(unit)
			if err != nil {
				panic(err)
			}

			switch dir.Direction {
			case directions.Down:
				sprite.Sprite.SetAnimation("idle-down")
			case directions.Left:
				sprite.Sprite.SetAnimation("idle-left")
			case directions.Up:
				sprite.Sprite.SetAnimation("idle-up")
			case directions.Right:
				sprite.Sprite.SetAnimation("idle-right")
			}

			if singletons.Turn.State == turnstate.Action {
				singletons.Turn.State = turnstate.Input
			}

			energy, err := pools.EnergyPool.Component(unit)
			if err != nil {
				panic(err)
			}

			class, err := pools.ClassPool.Component(unit)
			if err != nil {
				panic(err)
			}

			energy.Energy -= singletons.ClassStats[class.Class].MoveCost

			// в конце удаляем флаги таргета для наблюдателя
			// if singletons.Turn.State == turnstate.Wait {
			// 	for _, ent := range pools.TargetUnitFlag.Entities() {
			// 		pools.TargetUnitFlag.RemoveEntity(ent)
			// 	}

			// 	for _, ent := range pools.TargetObjectFlag.Entities() {
			// 		pools.TargetObjectFlag.RemoveEntity(ent)
			// 	}
			// }
		}
	}
}

type AttackSystem struct{}

func (s *AttackSystem) Run() {
	// Если пользователь смотрит анимацию (Action или Wait)
	if singletons.Turn.State == turnstate.Input {
		return
	}

	// если есть юнит взятый в цель, который не атакует...
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag}, []ecs.AnyPool{pools.TweenPool, pools.AttackPool})

	// ... и тайл взятый в цель
	tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.TargetObjectFlag}, []ecs.AnyPool{})

	if len(units) > 1 { // если их больше чем 1 всё плохо
		panic("More than one targeted units")
	} else if len(tiles) > 1 {
		panic("More than one targeted objects")
	} else if len(units) == 1 && len(tiles) == 1 { // если их по одному, то стартуем новую анимацию ИНАЧЕ идём проверять #2#

		println("u can win dis figt")

		unit := units[0]
		tile := tiles[0]

		occupied, err := pools.OccupiedPool.Component(tile)
		if err != nil {
			panic(err)
		}

		// Пропускаем если таргет не на тайле с юнитом (это ходьба или взаимодействие)
		if occupied.UnitObject == nil || occupied.StaticObject != nil || (occupied.ActiveObject != nil && !pools.SoftFlag.HasEntity(*occupied.ActiveObject)) {
			return
		}
		// иначе начинаем анимацию атаки

		println("now we figtin")

		position, err := pools.PositionPool.Component(*occupied.UnitObject)
		if err != nil {
			panic(err)
		}

		attackerDirection, err := pools.DirectionPool.Component(unit)
		if err != nil {
			panic(err)
		}

		attackerSprite, err := pools.SpritePool.Component(unit)
		if err != nil {
			panic(err)
		}

		attackerPosition, err := pools.PositionPool.Component(unit)
		if err != nil {
			panic(err)
		}

		// TODO WE ARE HERE
		deltaX := position.X - attackerPosition.X
		deltaY := position.Y - attackerPosition.Y
		if math.Abs(float64(deltaX)) > math.Abs(float64(deltaY)) {
			if deltaX > 0 {
				attackerDirection.Direction = directions.Right
			} else {
				attackerDirection.Direction = directions.Left
			}
		} else {
			if deltaY > 0 {
				attackerDirection.Direction = directions.Down
			} else {
				attackerDirection.Direction = directions.Up
			}
		}

		switch attackerDirection.Direction {
		case directions.Down:
			pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Back75Forward25, 1000, 0, 8, 0)})
			attackerSprite.Sprite.SetAnimation("attack-down")
		case directions.Up:
			pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Back75Forward25, 1000, 0, -8, 0)})
			attackerSprite.Sprite.SetAnimation("attack-up")
		case directions.Right:
			pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Back75Forward25, 1000, 8, 0, 0)})
			attackerSprite.Sprite.SetAnimation("attack-right")
		case directions.Left:
			pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Back75Forward25, 1000, -8, 0, 0)})
			attackerSprite.Sprite.SetAnimation("attack-left")
		}

		pools.AttackPool.AddExistingEntity(unit, components.Attack{Target: occupied.UnitObject})

		for _, ent := range tiles {
			pools.TargetObjectFlag.RemoveEntity(ent)
		}

		return
	}

	// #2# Проверка условий завершения анимации атаки

	// Проверим есть ли юниты атакующие в данный момент
	units = ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag, pools.TweenPool, pools.AttackPool}, []ecs.AnyPool{})

	if len(units) > 1 { // Если взято в цель больше одного юнита -> всё плохо
		panic("More than one targeted units")
	} else if len(units) == 1 { // Нормальное состояние
		unit := units[0]

		t, err := pools.TweenPool.Component(unit)
		if err != nil {
			panic(err)
		}

		// Если закончилась анимация атаки, то анимацию нужно сменить и пометить юнит, которому нанесли урон
		if t.Animation.IsEnded() {
			attack, err := pools.AttackPool.Component(unit)
			if err != nil {
				panic(err)
			}

			attackerSprite, err := pools.SpritePool.Component(unit)
			if err != nil {
				panic(err)
			}

			attackerClass, err := pools.ClassPool.Component(unit)
			if err != nil {
				panic(err)
			}

			attackerEnergy, err := pools.EnergyPool.Component(unit)
			if err != nil {
				panic(err)
			}

			attackerDir, err := pools.DirectionPool.Component(unit)
			if err != nil {
				panic(err)
			}

			attackerEnergy.Energy -= singletons.ClassStats[attackerClass.Class].AttackCost
			damageMult := 1

			// ножь в спину дамажит больше
			if attackerClass.Class == classes.Knife {
				targetDir, err := pools.DirectionPool.Component(*attack.Target)
				if err != nil {
					panic(err)
				}
				if attackerDir.Direction == targetDir.Direction {
					damageMult = 2
				}
			}

			// глефа бъёт по площади
			if attackerClass.Class == classes.Glaive {
				attackerPos, err := pools.PositionPool.Component(unit)
				if err != nil {
					panic(err)
				}

				offestX := 0
				offestY := 0

				switch attackerDir.Direction {
				case directions.Down:
					offestY = 2
				case directions.Up:
					offestY = -2
				case directions.Right:
					offestX = 2
				case directions.Left:
					offestX = -2
				}

				for _, tile := range pools.TileFlag.Entities() {
					tilePos, err := pools.PositionPool.Component(tile)
					if err != nil {
						panic(err)
					}

					occupied, err := pools.OccupiedPool.Component(tile)
					if err != nil {
						panic(err)
					}
					if occupied.UnitObject != nil && !pools.DeadFlag.HasEntity(*occupied.UnitObject) &&
						(offestY == 0 && tilePos.Y >= attackerPos.Y-1 && tilePos.Y <= attackerPos.Y+1 && tilePos.X == attackerPos.X+offestX ||
							offestX == 0 && tilePos.X >= attackerPos.X-1 && tilePos.X <= attackerPos.X+1 && tilePos.Y == attackerPos.Y+offestY) {

						pools.DamagePool.AddExistingEntity(*occupied.UnitObject, components.Damage{Value: singletons.ClassStats[attackerClass.Class].Attack, Type: damagetype.Hit})
					}
				}
			} else {
				pools.DamagePool.AddExistingEntity(*attack.Target, components.Damage{Value: singletons.ClassStats[attackerClass.Class].Attack * uint8(damageMult), Type: damagetype.Hit})
			}

			pools.TweenPool.RemoveEntity(unit)
			pools.AttackPool.RemoveEntity(unit)
			singletons.Turn.IsAttackAllowed = false

			switch attackerDir.Direction {
			case directions.Down:
				attackerSprite.Sprite.SetAnimation("idle-down")
			case directions.Up:
				attackerSprite.Sprite.SetAnimation("idle-up")
			case directions.Right:
				attackerSprite.Sprite.SetAnimation("idle-right")
			case directions.Left:
				attackerSprite.Sprite.SetAnimation("idle-left")
			}

			// TODO ОТКЛЮЧИТЬ
			// if singletons.Turn.State == turnstate.Action {
			// 	singletons.Turn.State = turnstate.Input
			// }
		}
	}
}

type DamageSystem struct{}

func (s *DamageSystem) Run() {
	// Если пользователь смотрит анимацию (Action или Wait)
	// if singletons.Turn.State == turnstate.Input {
	// 	return
	// }

	// #1# если есть юнит получивший урон и не анимируемый в данный момент
	units := ecs.PoolFilter([]ecs.AnyPool{pools.DamagePool}, []ecs.AnyPool{pools.TweenPool})

	for _, unit := range units {
		damage, err := pools.DamagePool.Component(unit)
		if err != nil {
			panic(err)
		}

		health, err := pools.HealthPool.Component(unit)
		if err != nil {
			panic(err)
		}

		switch damage.Type {
		case damagetype.Hit:
			println("START DAMAGE")
			spr, err := pools.SpritePool.Component(unit)
			if err != nil {
				panic(err)
			}

			dir, err := pools.DirectionPool.Component(unit)
			if err != nil {
				panic(err)
			}

			pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.XSin, 800, 0, 0, 0)})

			switch dir.Direction {
			case directions.Down:
				spr.Sprite.SetAnimation("hit-down")
			case directions.Right:
				spr.Sprite.SetAnimation("hit-right")
			case directions.Left:
				spr.Sprite.SetAnimation("hit-left")
			case directions.Up:
				spr.Sprite.SetAnimation("hit-up")
			}

			if health.Health < damage.Value {
				health.Health = 0
			} else {
				health.Health -= damage.Value
			}

		default:
			panic("UNEXPECTED DAMAGE TYPE")
		}
	}

	// #2# если есть юнит получивший урон и анимируемый в данный момент
	units = ecs.PoolFilter([]ecs.AnyPool{pools.DamagePool, pools.TweenPool}, []ecs.AnyPool{})

	for _, unit := range units {
		tween, err := pools.TweenPool.Component(unit)
		if err != nil {
			panic(err)
		}

		if tween.Animation.IsEnded() {
			println("STOP DAMAGE")
			spr, err := pools.SpritePool.Component(unit)
			if err != nil {
				panic(err)
			}

			dir, err := pools.DirectionPool.Component(unit)
			if err != nil {
				panic(err)
			}

			health, err := pools.HealthPool.Component(unit)
			if err != nil {
				panic(err)
			}

			if health.Health > 0 {
				switch dir.Direction {
				case directions.Down:
					spr.Sprite.SetAnimation("idle-down")
				case directions.Left:
					spr.Sprite.SetAnimation("idle-left")
				case directions.Right:
					spr.Sprite.SetAnimation("idle-right")
				case directions.Up:
					spr.Sprite.SetAnimation("idle-up")
				}
			} else {
				switch dir.Direction {
				case directions.Down:
					spr.Sprite.SetAnimation("dead-down")
				case directions.Left:
					spr.Sprite.SetAnimation("dead-left")
				case directions.Right:
					spr.Sprite.SetAnimation("dead-right")
				case directions.Up:
					spr.Sprite.SetAnimation("dead-up")
				}

				pools.DeadFlag.AddExistingEntity(unit)
			}

			pools.TweenPool.RemoveEntity(unit)
			pools.DamagePool.RemoveEntity(unit)

			if singletons.Turn.State == turnstate.Action {
				singletons.Turn.State = turnstate.Input
			}
		}

	}

}

type EnergySystem struct{}

func (s *EnergySystem) Run() {
	// TODO if any need
	if singletons.Turn.State == turnstate.Action {
		return
	}
	// Проверим не закончилась ли энергия у юнита в таргете
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag}, []ecs.AnyPool{})
	if len(units) > 1 { // если их больше чем 1 всё плохо
		panic("More than one targeted units")
	} else if len(units) == 1 {
		unit := units[0]
		energy, err := pools.EnergyPool.Component(unit)
		if err != nil {
			panic(err)
		}

		if energy.Energy == 0 {
			println("CUZ U EXAUSTED")
			singletons.Turn.IsTurnEnds = true
		}
	}

	if singletons.Turn.IsTurnEnds {
		log.Println("SKIPED!!!")
		// передаём ход другому игроку
		if singletons.Turn.CurrentTurn == teams.Blue {
			singletons.Turn.CurrentTurn = teams.Red
		} else {
			singletons.Turn.CurrentTurn = teams.Blue
		}

		for _, ent := range pools.TargetUnitFlag.Entities() {
			pools.TargetUnitFlag.RemoveEntity(ent)
		}

		for _, ent := range pools.TargetObjectFlag.Entities() {
			pools.TargetObjectFlag.RemoveEntity(ent)
		}

		for _, ent := range pools.ActiveFlag.Entities() { // ? useless
			pools.ActiveFlag.RemoveEntity(ent)
		}

		if singletons.AppState.GameMode == gamemode.Online {
			println("I AM HERE WITH TURNSTATE: ", singletons.Turn.State.String())
			if singletons.Turn.State == turnstate.Wait {
				singletons.Turn.State = turnstate.Input
			} else {
				// network.SendSkip()
				singletons.Turn.State = turnstate.Wait
			}
		} else {
			if singletons.Turn.PlayerTeam == teams.Blue {
				singletons.Turn.PlayerTeam = teams.Red
			} else {
				singletons.Turn.PlayerTeam = teams.Blue
			}
			singletons.Turn.State = turnstate.Input
		}

		units = ecs.PoolFilter([]ecs.AnyPool{pools.UnitFlag}, []ecs.AnyPool{})
		for _, unit := range units {
			energy, err := pools.EnergyPool.Component(unit)
			if err != nil {
				panic(err)
			}

			class, err := pools.ClassPool.Component(unit)
			if err != nil {
				panic(err)
			}
			if energy.Energy < singletons.ClassStats[class.Class].MaxEnergy {
				energy.Energy += singletons.ClassStats[class.Class].EnergyPerTurn
			}
		}
		singletons.Turn.IsTurnEnds = false
		singletons.Turn.IsAttackAllowed = true
		println("WEARE ACTUALY ENDING")
	}
}

// ###
// RENDER SYSTEMS
// ###

type DrawWorldSystem struct{}

func (s *DrawWorldSystem) Run(screen *ebiten.Image) {
	unitRenderQueue := []ecs.Entity{}
	objectRenderQueue := []ecs.Entity{}
	view := singletons.View.Image

	// tile render
	for _, tileEntity := range pools.TileFlag.Entities() {
		position, err := pools.PositionPool.Component(tileEntity)
		if err != nil {
			panic(err)
		}

		sprite, err := pools.SpritePool.Component(tileEntity)
		if err != nil {
			panic(err)
		}

		render, err := pools.ImageRenderPool.Component(tileEntity)
		if err != nil {
			panic(err)
		}

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
		options.GeoM.Concat(render.Options.GeoM)

		options.Blend = render.Options.Blend

		options.ColorScale = render.Options.ColorScale

		if pools.TargetObjectFlag.HasEntity(tileEntity) { // draw active tile
			options.ColorScale.Scale(2, 1, 1, 1)
		} else if pools.ActiveFlag.HasEntity(tileEntity) {
			options.ColorScale.Scale(1, 2, 1, 1)
		}

		options.Filter = render.Options.Filter
		view.DrawImage(sprite.Sprite.Animate(), &options)

		// add objects to queue
		occupied, err := pools.OccupiedPool.Component(tileEntity)
		if err != nil {
			panic(err)
		}

		if occupied.UnitObject != nil {
			unitRenderQueue = append(unitRenderQueue, *occupied.UnitObject)
		}

		if occupied.ActiveObject != nil {
			objectRenderQueue = append(objectRenderQueue, *occupied.ActiveObject)
		}

		if occupied.StaticObject != nil {
			objectRenderQueue = append(objectRenderQueue, *occupied.StaticObject)
		}
	}

	//unit render
	for _, unitEntity := range unitRenderQueue {
		img, opt := entityImage(unitEntity)

		// units highlight
		if pools.TargetUnitFlag.HasEntity(unitEntity) {
			opt.ColorScale.Scale(2, 1, 1, 1)
		} else if pools.ActiveFlag.HasEntity(unitEntity) {
			opt.ColorScale.Scale(1, 2, 1, 1)
		}

		if pools.TweenPool.HasEntity(unitEntity) {
			tweenComp, err := pools.TweenPool.Component(unitEntity)
			if err != nil {
				panic(err)
			}

			val := tweenComp.Animation.Animate()

			opt.GeoM.Translate(float64(val.X), float64(val.Y))
			opt.GeoM.Rotate(float64(val.Angle))
		}

		view.DrawImage(img, opt)
	}

	// object queue
	for _, objectEntity := range objectRenderQueue {
		img, opt := entityImage(objectEntity)

		// objects highlight
		if pools.TargetObjectFlag.HasEntity(objectEntity) {
			println("Targeted object")
			opt.ColorScale.Scale(2, 1, 1, 1)
		} else if pools.ActiveFlag.HasEntity(objectEntity) {
			println("Active object")
			opt.ColorScale.Scale(1, 2, 1, 1)
		}

		view.DrawImage(img, opt)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(singletons.View.Scale), float64(singletons.View.Scale))
	opt.GeoM.Translate(float64(singletons.View.ShiftX), float64(singletons.View.ShiftY))
	screen.DrawImage(view, opt)
}

type DrawGhostsSystem struct{}

func (s *DrawGhostsSystem) Run(screen *ebiten.Image) {
	view := singletons.View.Image
	for _, ghostEntity := range pools.GhostFlag.Entities() {
		position, err := pools.PositionPool.Component(ghostEntity)
		if err != nil {
			panic(err)
		}

		sprite, err := pools.SpritePool.Component(ghostEntity)
		if err != nil {
			panic(err)
		}

		render, err := pools.ImageRenderPool.Component(ghostEntity)
		if err != nil {
			panic(err)
		}

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
		options.GeoM.Concat(render.Options.GeoM)

		if pools.TweenPool.HasEntity(ghostEntity) {
			tweenComp, err := pools.TweenPool.Component(ghostEntity)
			if err != nil {
				panic(err)
			}
			val := tweenComp.Animation.GetValue()
			options.GeoM.Translate(float64(val.X), float64(val.Y))
			options.GeoM.Rotate(float64(val.Angle))
		}

		options.ColorScale = render.Options.ColorScale
		options.ColorScale.ScaleAlpha(0.6)
		options.Filter = render.Options.Filter
		view.DrawImage(sprite.Sprite.Animate(), &options)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(singletons.View.Scale), float64(singletons.View.Scale))
	opt.GeoM.Translate(float64(singletons.View.ShiftX), float64(singletons.View.ShiftY))
	screen.DrawImage(view, opt)
}

// TODO move with tween
type DrawStatsSystem struct{}

func (s *DrawStatsSystem) Run(screen *ebiten.Image) {
	// Текст необходимо рисовать сразу на screen избегая view т.к. масштабировать текст плохо
	// view := singletons.View.Image
	for _, unitEntity := range pools.UnitFlag.Entities() {
		if pools.TweenPool.HasEntity(unitEntity) || pools.DeadFlag.HasEntity(unitEntity) {
			continue
		}

		position, err := pools.PositionPool.Component(unitEntity)
		if err != nil {
			panic(err)
		}

		energy, err := pools.EnergyPool.Component(unitEntity)
		if err != nil {
			panic(err)
		}

		health, err := pools.HealthPool.Component(unitEntity)
		if err != nil {
			panic(err)
		}

		op := &text.DrawOptions{}
		op.ColorScale.SetR(255)
		op.ColorScale.SetG(255)
		op.ColorScale.SetB(0)
		op.ColorScale.SetA(255)

		op.GeoM.Scale(float64(singletons.View.Scale)*0.25, float64(singletons.View.Scale)*0.25)
		op.GeoM.Translate(float64((position.X*16+8)*singletons.View.Scale), float64((position.Y*16-4)*singletons.View.Scale))
		op.GeoM.Translate(float64(singletons.View.ShiftX), float64(singletons.View.ShiftY))

		text.Draw(screen, strconv.FormatUint(uint64(energy.Energy), 10), ui.TextFace, op)

		op = &text.DrawOptions{}
		// op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 0})
		op.ColorScale.SetR(0)
		op.ColorScale.SetG(225)
		op.ColorScale.SetB(0)
		op.ColorScale.SetA(255)

		op.GeoM.Scale(float64(singletons.View.Scale)*0.25, float64(singletons.View.Scale)*0.25)
		op.GeoM.Translate(float64((position.X*16)*singletons.View.Scale), float64((position.Y*16-4)*singletons.View.Scale))
		op.GeoM.Translate(float64(singletons.View.ShiftX), float64(singletons.View.ShiftY))

		text.Draw(screen, strconv.FormatUint(uint64(health.Health), 10), ui.TextFace, op)

		// opt := &ebiten.DrawImageOptions{}
		// opt.GeoM.Scale(float64(singletons.View.Scale), float64(singletons.View.Scale))
		// opt.GeoM.Translate(float64(singletons.View.ShiftX), float64(singletons.View.ShiftY))
		// screen.DrawImage(view, opt)
	}
}

func entityImage(objectEntity ecs.Entity) (*ebiten.Image, *ebiten.DrawImageOptions) {
	position, err := pools.PositionPool.Component(objectEntity)
	if err != nil {
		panic(err)
	}

	sprite, err := pools.SpritePool.Component(objectEntity)
	if err != nil {
		panic(err)
	}

	render, err := pools.ImageRenderPool.Component(objectEntity)
	if err != nil {
		panic(err)
	}

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
	options.GeoM.Concat(render.Options.GeoM)

	options.Blend = render.Options.Blend
	options.ColorScale = render.Options.ColorScale
	options.Filter = render.Options.Filter
	return sprite.Sprite.Animate(), &options
}

func positionsDistance(a *components.Position, b *components.Position) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}
