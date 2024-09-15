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
	"strategy-game/sprite"
	tile "strategy-game/tilemap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var frameCounter int = 0

var w *ecs.World

var positionPool *ecs.ComponentPool[c.Position]
var spritePool *ecs.ComponentPool[c.Sprite]
var materialPool *ecs.ComponentPool[c.Material]

var tileFlag *ecs.FlagPool

type Game struct{}

func (g *Game) Update() error {
	frameCounter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// geom := ebiten.GeoM{}
	// geom.Rotate(0.05 * float64(frameCounter))

	// screen.DrawImage(s.Animate(frameCounter), &ebiten.DrawImageOptions{GeoM: geom, Filter: ebiten.FilterNearest})
	// screen.DrawImage(img, &ebiten.DrawImageOptions{})
	DrawTiles(screen)
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	tilesetImg, _, err := ebitenutil.NewImageFromFile("assets/ground-tileset.png")
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
	positionPool = ecs.CreateComponentPool[c.Position](w, psize.Page1024)
	spritePool = ecs.CreateComponentPool[c.Sprite](w, psize.Page1024)
	materialPool = ecs.CreateComponentPool[c.Material](w, psize.Page1024)

	tileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)

	InitTileEntities("assets/ground-tilemap.json", "assets/ground-tileset.json", tilesetImg)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// INIT ECS POOLS

func InitTileEntities(tilemapFilepath string, tilesetFilepath string, image *ebiten.Image) {
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

	contents, err = os.ReadFile(tilesetFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var tileset tile.TilesetJSON
	err = json.Unmarshal(contents, &tileset)
	if err != nil {
		log.Fatal(err)
	}

	groundLayer := tilemap.Layers[0] // 1 layer (ground)

	for i := 0; i < groundLayer.Height*groundLayer.Width; i++ {
		id := groundLayer.Data[i] - 1

		// SPRITE
		s := sprite.NewSprite(image, tileset.Width, tileset.Height)
		for _, a := range tileset.AnimatedTiles {
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
		entity, err := spritePool.AddNewEntity(spriteComp)
		if err != nil {
			log.Fatal(err)
		}

		// MATERIAL
		materialComp := c.Material{}
		if slices.Contains(GRASS_TILES, id) {
			materialComp.Material = material.Grass
		} else if slices.Contains(SAND_TILES, id) {
			materialComp.Material = material.Sand
		} else {
			materialComp.Material = material.Water
		}
		materialPool.AddExistingEntity(entity, materialComp)

		// POSITION
		positionComp := c.Position{X: i % groundLayer.Width, Y: i / groundLayer.Width}
		positionPool.AddExistingEntity(entity, positionComp)

		// FLAGS
		tileFlag.AddExistingEntity(entity)
	}
}

// DRAW FUNCTIONS

// func AnimateSprites(frameCounter int) {
// 	for _, ent := range spritePool.Entities() {
// 		sprite, err := spritePool.Component(ent)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		sprite.Sprite.Animate(frameCounter)
// 	}
// }

func DrawTiles(screen *ebiten.Image) {
	for _, ent := range tileFlag.Entities() {
		position, err := positionPool.Component(ent)
		if err != nil {
			log.Fatal(err)
		}
		sprite, err := spritePool.Component(ent)
		if err != nil {
			log.Fatal(err)
		}
		options := ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(position.X*16), float64(position.Y*16))
		screen.DrawImage(sprite.Sprite.Animate(frameCounter), &options)
	}
}
