package systems

import (
	"log"
	"strategy-game/ecs"
	"strategy-game/pools"

	"github.com/hajimehoshi/ebiten/v2"
)

type DrawTilesSystem struct{}

func (s *DrawTilesSystem) Run(frameCount int) {
	viewEntity := pools.ViewPool.Entities()[0]
	view, err := pools.ViewPool.Component(viewEntity)
	if err != nil {
		log.Fatal(err)
	}

	for _, tileEntity := range pools.TileFlag.Entities() {
		position, err := pools.PositionPool.Component(tileEntity)
		if err != nil {
			log.Fatal(err)
		}

		sprite, err := pools.SpritePool.Component(tileEntity)
		if err != nil {
			log.Fatal(err)
		}

		render, err := pools.ImageRenderPool.Component(tileEntity)
		if err != nil {
			log.Fatal(err)
		}

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
		options.GeoM.Concat(render.Options.GeoM)

		// options.GeoM.Scale(4.0, 4.0)
		// options.GeoM.Rotate(float64(frameCount) * 0.01)

		options.Blend = render.Options.Blend
		options.ColorScale = render.Options.ColorScale
		options.Filter = render.Options.Filter
		view.Img.DrawImage(sprite.Sprite.Animate(frameCount), &options)

		// OBJECT
		occupied, err := pools.OccupiedPool.Component(tileEntity)
		if err != nil {
			log.Fatal(err)
		}

		if occupied.ActiveObject != nil {
			objectEntity := *occupied.ActiveObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.Img.DrawImage(img, opt)
		}

		if occupied.StaticObject != nil {
			objectEntity := *occupied.StaticObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.Img.DrawImage(img, opt)
		}

		if occupied.UnitObject != nil { // TODO: UNIT
			objectEntity := *occupied.UnitObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.Img.DrawImage(img, opt)
		}
	}
}

func entityDrawData(objectEntity ecs.Entity, frameCount int) (*ebiten.Image, *ebiten.DrawImageOptions) {
	position, err := pools.PositionPool.Component(objectEntity)
	if err != nil {
		log.Fatal(err)
	}

	sprite, err := pools.SpritePool.Component(objectEntity)
	if err != nil {
		log.Fatal(err)
	}

	render, err := pools.ImageRenderPool.Component(objectEntity)
	if err != nil {
		log.Fatal(err)
	}

	options := ebiten.DrawImageOptions{}
	options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
	options.GeoM.Concat(render.Options.GeoM)

	options.Blend = render.Options.Blend
	options.ColorScale = render.Options.ColorScale
	options.Filter = render.Options.Filter
	return sprite.Sprite.Animate(frameCount), &options
}
