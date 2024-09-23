package systems

import (
	"log"
	"strategy-game/pools"

	"github.com/hajimehoshi/ebiten/v2"
)

type DrawTilemapSystem struct{}

func (s *DrawTilemapSystem) Run(frameCount int) {
	ent := pools.ViewPool.Entities()[0]
	view, err := pools.ViewPool.Component(ent)
	if err != nil {
		log.Fatal(err)
	}

	for _, ent := range pools.TileFlag.Entities() {
		position, err := pools.PositionPool.Component(ent)
		if err != nil {
			log.Fatal(err)
		}

		sprite, err := pools.SpritePool.Component(ent)
		if err != nil {
			log.Fatal(err)
		}

		render, err := pools.ImageRenderPool.Component(ent)
		if err != nil {
			log.Fatal(err)
		}

		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16)) // TILE SIZE???
		options.GeoM.Concat(render.Options.GeoM)
		options.Blend = render.Options.Blend
		options.ColorScale = render.Options.ColorScale
		options.Filter = render.Options.Filter
		view.Img.DrawImage(sprite.Sprite.Animate(frameCount), &options)
	}
}
