package systems

import (
	"math"
	"strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/util/ecs"
	"strategy-game/util/gamedata"

	"github.com/hajimehoshi/ebiten/v2"
)

// ###
// LOGIC SYSTEMS
// ###

type ActiveUnitsSystem struct{}

func (s *ActiveUnitsSystem) Run(g gamedata.GameData) {
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

type ActionSystem struct{}

func (s *ActionSystem) Run(g gamedata.GameData) {
	// get targeted unit
	units := ecs.PoolFilter([]ecs.AnyPool{pools.TargetUnitFlag, pools.EnergyPool, pools.ClassPool, pools.PositionPool}, []ecs.AnyPool{})
	if len(units) > 1 {
		panic("More than one targeted units")
	}

	// all tiles
	tiles := ecs.PoolFilter([]ecs.AnyPool{pools.TileFlag, pools.PositionPool, pools.OccupiedPool}, []ecs.AnyPool{})

	for _, unit := range units {
		// energy, err := pools.EnergyPool.Component(unit)
		// if err != nil {
		// 	panic(err)
		// }

		// class, err := pools.ClassPool.Component(unit)
		// if err != nil {
		// 	panic(err)
		// }

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

// ###
// RENDER SYSTEMS
// ###

type DrawWorldSystem struct{}

func (s *DrawWorldSystem) Run(g gamedata.GameData, screen *ebiten.Image) {
	frameCount := g.FrameCount()
	view := g.View()
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

		// OBJECT
		occupied, err := pools.OccupiedPool.Component(tileEntity)
		if err != nil {
			panic(err)
		}

		if occupied.UnitObject != nil {
			unitEntity := *occupied.UnitObject
			img, opt := entityImage(unitEntity, frameCount)

			// units highlight
			if pools.TargetUnitFlag.HasEntity(unitEntity) {
				opt.ColorScale.Scale(2, 1, 1, 1)
			} else if pools.ActiveFlag.HasEntity(unitEntity) {
				opt.ColorScale.Scale(1, 2, 1, 1)
			}

			view.DrawImage(img, opt)
		}

		if occupied.ActiveObject != nil {
			objectEntity := *occupied.ActiveObject
			img, opt := entityImage(objectEntity, frameCount)

			// objects highlight
			if pools.TargetObjectFlag.HasEntity(objectEntity) {
				println("1")
				opt.ColorScale.Scale(2, 1, 1, 1)
			} else if pools.ActiveFlag.HasEntity(objectEntity) {
				println("2")
				opt.ColorScale.Scale(1, 2, 1, 1)
			}

			view.DrawImage(img, opt)
		}

		if occupied.StaticObject != nil {
			objectEntity := *occupied.StaticObject
			img, opt := entityImage(objectEntity, frameCount)
			view.DrawImage(img, opt)
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(g.ViewScale()), float64(g.ViewScale()))
	screen.DrawImage(view, opt)
}

type DrawGhostsSystem struct{}

func (s *DrawGhostsSystem) Run(g gamedata.GameData, screen *ebiten.Image) {
	frameCount := g.FrameCount()
	view := g.View()
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

		options.ColorScale = render.Options.ColorScale
		options.ColorScale.ScaleAlpha(0.6)
		options.Filter = render.Options.Filter
		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(g.ViewScale()), float64(g.ViewScale()))
	screen.DrawImage(view, opt)
}

// type DrawActiveSystem struct{}

// func (s *DrawActiveSystem) Run(g gamedata.GameData, screen *ebiten.Image) {
// 	frameCount := g.FrameCount()
// 	view := g.View()
// 	for _, HighlightedEntity := range pools.ActiveFlag.Entities() {
// 		position, err := pools.PositionPool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		sprite, err := pools.SpritePool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		render, err := pools.ImageRenderPool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		options := ebiten.DrawImageOptions{}
// 		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
// 		options.GeoM.Concat(render.Options.GeoM)

// 		options.ColorScale = render.Options.ColorScale
// 		options.ColorScale.Scale(1.5, 1.5, 1.5, 1)
// 		options.Filter = render.Options.Filter
// 		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)
// 	}

// 	opt := &ebiten.DrawImageOptions{}
// 	opt.GeoM.Scale(float64(g.ViewScale()), float64(g.ViewScale()))
// 	screen.DrawImage(view, opt)
// }

// type DrawTargetedSystem struct{}

// func (s *DrawTargetedSystem) Run(g gamedata.GameData, screen *ebiten.Image) {
// 	frameCount := g.FrameCount()
// 	view := g.View()
// 	for _, HighlightedEntity := range pools.TargetFlag.Entities() {
// 		position, err := pools.PositionPool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		sprite, err := pools.SpritePool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		render, err := pools.ImageRenderPool.Component(HighlightedEntity)
// 		if err != nil {
// 			panic(err)
// 		}

// 		options := ebiten.DrawImageOptions{}
// 		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
// 		options.GeoM.Concat(render.Options.GeoM)

// 		options.ColorScale = render.Options.ColorScale
// 		options.ColorScale.Scale(2, 1, 1, 1)
// 		options.Filter = render.Options.Filter
// 		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)
// 	}

// 	opt := &ebiten.DrawImageOptions{}
// 	opt.GeoM.Scale(float64(g.ViewScale()), float64(g.ViewScale()))
// 	screen.DrawImage(view, opt)
// }

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
