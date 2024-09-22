package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	c "strategy-game/components"
	"strategy-game/components/material"
	"strategy-game/ecs"
	"strategy-game/ecs/psize"
	"strategy-game/pools"
	"strategy-game/sprite"
	"strategy-game/systems"
	tile "strategy-game/tilemap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var RenderWidth int = 320
var RenderHeight int = 240

var frameCounter int = 0

var w *ecs.World

type Game struct{}

func (g *Game) Update() error {
	frameCounter++
	w.Update(frameCounter)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// geom := ebiten.GeoM{}
	// geom.Rotate(0.05 * float64(frameCounter))

	// screen.DrawImage(s.Animate(frameCounter), &ebiten.DrawImageOptions{GeoM: geom, Filter: ebiten.FilterNearest})
	// screen.DrawImage(img, &ebiten.DrawImageOptions{})
	DrawView(screen)
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return RenderWidth, RenderHeight
}

func main() {
	groundTilesetImg, _, err := ebitenutil.NewImageFromFile("assets/ground-tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	objectTilesetImg, _, err := ebitenutil.NewImageFromFile("assets/object-tileset.png")
	if err != nil {
		log.Fatal(err)
	}
	// s = sprite.NewSprite(img, 16, 16)
	// s.AddAnimation("anim", []sprite.Frame{
	// 	{N: 43, Time: 500},
	// 	{N: 48, Time: 500},
	// 	{N: 53, Time: 500},
	// 	{N: 58, Time: 500},
	// })
	// s.SetAnimation("anim")

	w = ecs.CreateWorld()
	pools.PositionPool = ecs.CreateComponentPool[c.Position](w, psize.Page1024)
	pools.SpritePool = ecs.CreateComponentPool[c.Sprite](w, psize.Page1024)
	pools.MaterialPool = ecs.CreateComponentPool[c.Material](w, psize.Page1024)
	pools.ViewPool = ecs.CreateComponentPool[c.View](w, psize.Page1)
	pools.ImageRenderPool = ecs.CreateComponentPool[c.ImageRender](w, psize.Page1024)

	pools.TileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)

	InitView()
	InitTileEntities("assets/tilemap.json", "assets/ground-tileset.json", "assets/object-tileset.json", groundTilesetImg, objectTilesetImg)

	InitSystems()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// INIT ECS POOLS

func InitView() {
	img := ebiten.NewImage(RenderWidth, RenderHeight)
	pools.ViewPool.AddNewEntity(c.View{Img: img})
}

func InitTileEntities(tilemapFilepath string, groundFilepath string, objectsFilepath string, groundTilesetImg *ebiten.Image, objectTilesetImg *ebiten.Image) {
	var GRASS_TILES = []int{22, 85}
	var SAND_TILES = []int{148, 211}

	// READ TILEMAP

	contents, err := os.ReadFile(tilemapFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var tilemap tile.TilemapJSON
	err = json.Unmarshal(contents, &tilemap)
	if err != nil {
		log.Fatal(err)
	}

	// READ TILESET

	contents, err = os.ReadFile(groundFilepath)
	if err != nil {
		log.Fatal(err)
	}
	var groundTileset tile.TilesetJSON
	err = json.Unmarshal(contents, &groundTileset)
	if err != nil {
		log.Fatal(err)
	}

	contents, err = os.ReadFile(objectsFilepath)
	if err != nil {
		log.Fatal(err)
	}
	var objectTileset tile.TilesetJSON
	err = json.Unmarshal(contents, &objectTileset)
	if err != nil {
		log.Fatal(err)
	}

	groundLayer := tilemap.Layers[0] // 1 layer (ground)
	objectLayer := tilemap.Layers[1] // 2 layer (objects)

	for i := 0; i < groundLayer.Height*groundLayer.Width; i++ {

		// ##GROUND##
		if groundLayer.Data[i] > 0 {
			id := groundLayer.Data[i] - 1

			// SPRITE
			s := sprite.NewSprite(groundTilesetImg, groundTileset.Width, groundTileset.Height)
			for _, a := range groundTileset.AnimatedTiles {
				if a.Id == id {
					s.AddAnimation("default", a.Frames)
					break
				}
			}
			if len(s.Animations) == 0 {
				s.AddAnimation("default", []sprite.Frame{{N: id, Time: 5000}})
			}
			s.SetAnimation("default")
			spriteComp := c.Sprite{}
			spriteComp.Sprite = s
			entity, err := pools.SpritePool.AddNewEntity(spriteComp)
			if err != nil {
				log.Fatal(err)
			}

			// SCREEN RENDERER
			opt := ebiten.DrawImageOptions{}
			pools.ImageRenderPool.AddExistingEntity(entity, c.ImageRender{Options: opt})

			// MATERIAL
			materialComp := c.Material{}
			if slices.Contains(GRASS_TILES, id) {
				materialComp.Material = material.Grass
			} else if slices.Contains(SAND_TILES, id) {
				materialComp.Material = material.Sand
			} else {
				materialComp.Material = material.Water
			}
			pools.MaterialPool.AddExistingEntity(entity, materialComp)

			// POSITION
			positionComp := c.Position{X: i % groundLayer.Width, Y: i / groundLayer.Width}
			pools.PositionPool.AddExistingEntity(entity, positionComp)

			// FLAGS
			pools.TileFlag.AddExistingEntity(entity)
		}

		// ##OBJECTS##
		if objectLayer.Data[i]-groundTileset.TileCount > 0 {
			id := objectLayer.Data[i] - groundTileset.TileCount - 1

			// SPRITE
			s := sprite.NewSprite(objectTilesetImg, objectTileset.Width, objectTileset.Height)

			for _, a := range objectTileset.AnimatedTiles {
				if a.Id == id {
					s.AddAnimation("default", a.Frames)
					break
				}
			}
			if len(s.Animations) == 0 {
				s.AddAnimation("default", []sprite.Frame{{N: id, Time: 5000}})
			}
			s.SetAnimation("default")
			spriteComp := c.Sprite{}
			spriteComp.Sprite = s
			entity, err := pools.SpritePool.AddNewEntity(spriteComp)
			if err != nil {
				log.Fatal(err)
			}

			opt := ebiten.DrawImageOptions{}
			opt.GeoM.Translate(0, -float64(s.Height()-16))
			pools.ImageRenderPool.AddExistingEntity(entity, c.ImageRender{Options: opt})

			// POSITION
			positionComp := c.Position{X: i % objectLayer.Width, Y: i / objectLayer.Width}
			pools.PositionPool.AddExistingEntity(entity, positionComp)

			// FLAGS
			pools.TileFlag.AddExistingEntity(entity) // OBJECT, NOT TILE
		}

	}

	// OBJECTS

	// for i := 0; i < objectLayer.Height*objectLayer.Width; i++ {
	// 	id := objectLayer.Data[i] - groundTileset.TileCount - 1

	// 	// SPRITE
	// 	s := sprite.NewSprite(objectTilesetImg, objectTileset.Width, objectTileset.Height)

	// 	for _, a := range objectTileset.AnimatedTiles {
	// 		if a.Id == id {
	// 			s.AddAnimation("default", a.Frames)
	// 			break
	// 		}
	// 	}
	// 	if len(s.Animations) == 0 {
	// 		s.AddAnimation("default", []sprite.Frame{{N: id, Time: 5000}})
	// 	}
	// 	s.SetAnimation("default")
	// 	spriteComp := c.Sprite{}
	// 	spriteComp.Sprite = s
	// 	entity, err := pools.SpritePool.AddNewEntity(spriteComp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	opt := ebiten.DrawImageOptions{}
	// 	opt.GeoM.Translate(0, -float64(s.Height()-16))
	// 	pools.ImageRenderPool.AddExistingEntity(entity, c.ImageRender{Options: opt})

	// 	// POSITION
	// 	positionComp := c.Position{X: i % objectLayer.Width, Y: i / objectLayer.Width}
	// 	pools.PositionPool.AddExistingEntity(entity, positionComp)

	// 	// FLAGS
	// 	pools.TileFlag.AddExistingEntity(entity)
	// }

}

// INIT SYSTEMS IN ORDER

func InitSystems() {
	ecs.AddSystem(w, &systems.DrawTilemapSystem{})
}

// REDRAW SCREEN

func DrawView(screen *ebiten.Image) {
	ent := pools.ViewPool.Entities()[0]
	view, err := pools.ViewPool.Component(ent)
	if err != nil {
		log.Fatal(err)
	}

	screen.DrawImage(view.Img, &ebiten.DrawImageOptions{})
}
