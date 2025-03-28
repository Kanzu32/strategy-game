package systems

import (
	"math"
	"strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/util/data/turn/turnstate"
	"strategy-game/util/data/tween"
	"strategy-game/util/data/tween/tweentype"
	"strategy-game/util/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// ###
// LOGIC SYSTEMS
// ###

type TurnSystem struct{}

func (s *TurnSystem) Run() { // highlight active units
	if singletons.Turn.State == turnstate.Action {
		for _, ent := range pools.ActiveFlag.Entities() {
			pools.ActiveFlag.RemoveEntity(ent)
		}
	}
}

type MarkActiveUnitsSystem struct{}

func (s *MarkActiveUnitsSystem) Run() { // highlight active units
	if singletons.Turn.State != turnstate.Input {
		return
	}

	if singletons.Turn.CurrentTurn == singletons.Turn.PlayerTeam {
		entities := ecs.PoolFilter([]ecs.AnyPool{pools.TeamPool, pools.EnergyPool}, []ecs.AnyPool{pools.ActiveFlag}) // all inactive units
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

	// get targeted unit
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag, pools.EnergyPool, pools.ClassPool, pools.PositionPool}, []ecs.AnyPool{})
	if len(units) > 1 {
		panic("More than one targeted units")
	}

	// all tiles
	tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.PositionPool, pools.OccupiedPool}, []ecs.AnyPool{})

	for _, unit := range units {

		unitPosition, err := pools.PositionPool.Component(unit)
		if err != nil {
			panic(err)
		}

		for _, tile := range tiles {
			tilePostion, err := pools.PositionPool.Component(tile)
			if err != nil {
				panic(err)
			}

			occupied, err := pools.OccupiedPool.Component(tile)
			if err != nil {
				panic(err)
			}

			if positionsDistance(unitPosition, tilePostion) == 1 &&
				!pools.ActiveFlag.HasEntity(tile) &&
				!pools.WallFlag.HasEntity(tile) &&
				occupied.UnitObject == nil &&
				(occupied.ActiveObject == nil || pools.SoftFlag.HasEntity(*occupied.ActiveObject)) {

				pools.ActiveFlag.AddExistingEntity(tile)
			} else if positionsDistance(unitPosition, tilePostion) != 1 && pools.ActiveFlag.HasEntity(tile) {
				pools.ActiveFlag.RemoveEntity(tile)
			}
		}
	}
}

type TweenMoveSystem struct{}

func (s *TweenMoveSystem) Run() {

	if singletons.Turn.State != turnstate.Action {
		return
	}

	// get targeted unit
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

	println("addin")
	unitPos, err := pools.PositionPool.Component(unit)
	if err != nil {
		panic(err)
	}

	tilePos, err := pools.PositionPool.Component(tile)
	if err != nil {
		panic(err)
	}

	pools.TweenPool.AddExistingEntity(unit, components.Tween{Animation: tween.CreateTween(tweentype.Linear, 1000, (tilePos.X-unitPos.X)*16, (tilePos.Y-unitPos.Y)*16, 0)})

	for _, ent := range tiles {
		pools.TargetObjectFlag.RemoveEntity(ent)
	}

	pools.MovePool.AddExistingEntity(unit, components.MoveDirection{X: int8(tilePos.X - unitPos.X), Y: int8(tilePos.Y - unitPos.Y)})
}

type UnitMoveSystem struct{}

func (s *UnitMoveSystem) Run() {

	if singletons.Turn.State != turnstate.Action {
		return
	}

	// get targeted unit
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag}, []ecs.AnyPool{})
	if len(units) > 1 {
		panic("More than one targeted units")
	} else if len(units) == 0 {
		panic("Zero targeted units")
	}
	unit := units[0]

	if !pools.TweenPool.HasEntity(unit) {
		return
	}

	t, err := pools.TweenPool.Component(unit)
	if err != nil {
		panic(err)
	}

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

			if occupied.UnitObject == &unit && pos.X == unitPos.X && pos.Y == unitPos.Y {
				occupied.UnitObject = nil
			}

			if occupied.UnitObject == nil && pos.X == unitPos.X+int(move.X) && pos.Y == unitPos.Y+int(move.Y) {
				occupied.UnitObject = &unit
			}
		}

		unitPos.X += int(move.X)
		unitPos.Y += int(move.Y)

		pools.TweenPool.RemoveEntity(unit)
		pools.MovePool.RemoveEntity(unit)

		singletons.Turn.State = turnstate.Input
	}

}

// ###
// RENDER SYSTEMS
// ###

type DrawWorldSystem struct{}

func (s *DrawWorldSystem) Run(screen *ebiten.Image) {
	unitRenderQueue := []ecs.Entity{}
	objectRenderQueue := []ecs.Entity{}
	frameCount := singletons.FrameCount
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
		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)

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
		img, opt := entityImage(unitEntity, frameCount)

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
			println("animatin")
			val := tweenComp.Animation.Animate()
			println(val.X, val.Y, val.Angle)
			opt.GeoM.Translate(float64(val.X), float64(val.Y))
			opt.GeoM.Rotate(float64(val.Angle))
		}

		view.DrawImage(img, opt)
	}

	// object queue
	for _, objectEntity := range objectRenderQueue {
		img, opt := entityImage(objectEntity, frameCount)

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
	screen.DrawImage(view, opt)
}

type DrawGhostsSystem struct{}

func (s *DrawGhostsSystem) Run(screen *ebiten.Image) {
	frameCount := singletons.FrameCount
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
		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(singletons.View.Scale), float64(singletons.View.Scale))
	screen.DrawImage(view, opt)
}

func entityImage(objectEntity ecs.Entity, frameCount int) (*ebiten.Image, *ebiten.DrawImageOptions) {
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
	return sprite.Sprite.Animate(frameCount), &options
}

func positionsDistance(a *components.Position, b *components.Position) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}
