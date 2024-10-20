package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	c "strategy-game/components"
	"strategy-game/ecs"
	"strategy-game/ecs/psize"
	"strategy-game/pools"
	"strategy-game/systems"
	tile "strategy-game/tilemap"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/furex/v2"
)

var RenderWidth int = 320
var RenderHeight int = 240

var frameCounter int = 0

var w *ecs.World

type Game struct {
	initOnce sync.Once
	screen   screen
	gameUI   *furex.View
}

type screen struct {
	Width  int
	Height int
}

func (g *Game) Update() error {
	g.initOnce.Do(func() {
		g.setupUI()
	})
	g.gameUI.UpdateWithSize(ebiten.WindowSize())

	frameCounter++
	w.Update(frameCounter)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawView(screen)
	g.gameUI.Draw(screen)

	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screen.Width = outsideWidth
	g.screen.Height = outsideHeight
	return g.screen.Width, g.screen.Height
}

func (g *Game) setupUI() {
	g.gameUI = &furex.View{
		Width:        g.screen.Width,
		Height:       g.screen.Height,
		Direction:    furex.Column,
		Justify:      furex.JustifyEnd,
		AlignItems:   furex.AlignItemStretch,
		AlignContent: furex.AlignContentStretch,
		Wrap:         furex.NoWrap,
	}

	// UI HERE

	bottomMargin := 0
	BottomMenu := furex.View{
		Height:       100,
		Handler:      &BottomMenu{},
		Bottom:       &bottomMargin,
		Display:      furex.DisplayFlex,
		Direction:    furex.Row,
		Justify:      furex.JustifySpaceAround,
		AlignItems:   furex.AlignItemCenter,
		AlignContent: furex.AlignContentCenter,
		Wrap:         furex.NoWrap,
	}

	g.gameUI.AddChild(&BottomMenu)

	img, _, err := ebitenutil.NewImageFromFile("assets/ui/test-icon.png")
	if err != nil {
		log.Fatal(err)
	}

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})

	BottomMenu.AddChild(&furex.View{
		Width:   80,
		Height:  80,
		Handler: &SpellBottom{Img: img},
	})
}

type BottomMenu struct{}

func (b *BottomMenu) Draw(screen *ebiten.Image, frame image.Rectangle, view *furex.View) {
	vector.DrawFilledRect(
		screen,
		float32(frame.Min.X),
		float32(frame.Min.Y),
		float32(frame.Size().X),
		float32(frame.Size().Y),
		color.White,
		false,
	)
}

type SpellBottom struct {
	Img *ebiten.Image
}

func (b *SpellBottom) Draw(screen *ebiten.Image, frame image.Rectangle, view *furex.View) {
	opt := ebiten.DrawImageOptions{}
	opt.Filter = ebiten.FilterNearest
	opt.GeoM.Scale(2, 2) // UI SCALE
	opt.GeoM.Translate(float64(frame.Min.X), float64(frame.Min.Y))
	screen.DrawImage(b.Img, &opt)
	// vector.DrawFilledRect(
	// 	screen,
	// 	float32(frame.Min.X),
	// 	float32(frame.Min.Y),
	// 	float32(frame.Size().X),
	// 	float32(frame.Size().Y),
	// 	color.Black,
	// 	false,
	// )
}

func main() {
	// s = sprite.NewSprite(img, 16, 16)
	// s.AddAnimation("anim", []sprite.Frame{
	// 	{N: 43, Time: 500},
	// 	{N: 48, Time: 500},
	// 	{N: 53, Time: 500},
	// 	{N: 58, Time: 500},
	// })
	// s.SetAnimation("anim")

	// INIT
	w = ecs.CreateWorld()
	InitPools()

	InitView()

	tilesets := tile.CreateTilesetArray([]string{
		"assets/tiles/tilesets/1_ground-tileset.json",
		"assets/tiles/tilesets/2_decals-tileset.json",
		"assets/tiles/tilesets/3_active-objects-tileset.json",
		"assets/tiles/tilesets/4_objects1-tileset.json",
		"assets/tiles/tilesets/5_objects2-tileset.json",
		"assets/tiles/tilesets/6_objects3-tileset.json",
		"assets/tiles/tilesets/7_objects4-tileset.json",
		"assets/tiles/tilesets/8_objects5-tileset.json",
		"assets/tiles/tilesets/9_objects6-tileset.json",
		"assets/tiles/tilesets/10_util-tileset.json",
	})
	InitTileEntities(tilesets, "assets/tiles/tilemaps/tilemap.json")

	InitSystems()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Troublemakers!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

// INIT ECS POOLS

func InitView() {
	img := ebiten.NewImage(RenderWidth, RenderHeight)
	ent, err := pools.ViewPool.AddNewEntity(c.View{Img: img})
	if err != nil {
		log.Fatal(err)
	}
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2, 2)
	pools.ImageRenderPool.AddExistingEntity(ent, c.ImageRender{Options: opt})
}

func InitTileEntities(tilesets tile.TilesetArray, tilemapFilepath string) {

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

	groundLayer := tilemap.Layers[0] // 1 layer (ground)
	objectLayer := tilemap.Layers[1] // 2 layer (objects)
	// utilLayer := tilemap.Layers[2]   // 2 layer (objects)

	for i := 0; i < groundLayer.Height*groundLayer.Width; i++ {

		// ##GROUND##
		tileData := tilesets.Get(groundLayer.Data[i])

		// SPRITE

		spriteComp := c.Sprite{Sprite: tileData.Sprite}
		tileEntity, err := pools.SpritePool.AddNewEntity(spriteComp)
		if err != nil {
			log.Fatal(err)
		}

		// IMAGE RENDERER
		opt := ebiten.DrawImageOptions{}
		pools.ImageRenderPool.AddExistingEntity(tileEntity, c.ImageRender{Options: opt})

		// MATERIAL
		materialComp := c.Material{Material: tileData.Material}
		pools.MaterialPool.AddExistingEntity(tileEntity, materialComp)

		// SIDE
		sideComp := c.Side{Side: tileData.Side}
		pools.SidePool.AddExistingEntity(tileEntity, sideComp)

		// POSITION
		positionComp := c.Position{X: i % groundLayer.Width, Y: i / groundLayer.Width}
		pools.PositionPool.AddExistingEntity(tileEntity, positionComp)

		// FLAGS
		pools.TileFlag.AddExistingEntity(tileEntity)

		// OCCUPIED

		occupiedComp := c.Occupied{StaticObject: nil, ActiveObject: nil, UnitObject: nil}

		// ##OBJECTS##
		if objectLayer.Data[i] > 0 {
			tileData := tilesets.Get(objectLayer.Data[i])

			// SPRITE
			spriteComp := c.Sprite{Sprite: tileData.Sprite}
			objectEntity, err := pools.SpritePool.AddNewEntity(spriteComp)
			if err != nil {
				log.Fatal(err)
			}

			// IMAGE RENDERER
			opt := ebiten.DrawImageOptions{}
			opt.GeoM.Translate(0, -float64(tileData.Sprite.Height()-16))
			pools.ImageRenderPool.AddExistingEntity(objectEntity, c.ImageRender{Options: opt})

			// POSITION
			positionComp := c.Position{X: i % objectLayer.Width, Y: i / objectLayer.Width}
			pools.PositionPool.AddExistingEntity(objectEntity, positionComp)

			// OCCUPIED (tile) by object
			if tileData.IsActive {
				occupiedComp.ActiveObject = &objectEntity
			} else {
				occupiedComp.StaticObject = &objectEntity
			}

			// FLAGS
			// pools.TileFlag.AddExistingEntity(entity) // OBJECT, NOT TILE
		}

		// OCCUPIED
		pools.OccupiedPool.AddExistingEntity(tileEntity, occupiedComp) // Doesn't include units yet!!!
	}

}

func InitPools() {
	pools.PositionPool = ecs.CreateComponentPool[c.Position](w, psize.Page1024)
	pools.SpritePool = ecs.CreateComponentPool[c.Sprite](w, psize.Page1024)
	pools.MaterialPool = ecs.CreateComponentPool[c.Material](w, psize.Page1024)
	pools.ViewPool = ecs.CreateComponentPool[c.View](w, psize.Page1)
	pools.ImageRenderPool = ecs.CreateComponentPool[c.ImageRender](w, psize.Page1024)
	pools.SidePool = ecs.CreateComponentPool[c.Side](w, psize.Page1024)
	pools.OccupiedPool = ecs.CreateComponentPool[c.Occupied](w, psize.Page1024)

	pools.SolidFlag = ecs.CreateFlagPool(w, psize.Page512)
	pools.TileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)
}

// INIT SYSTEMS IN ORDER

func InitSystems() {
	ecs.AddSystem(w, &systems.DrawTilesSystem{})
}

// REDRAW SCREEN

func DrawView(screen *ebiten.Image) {
	ent := pools.ViewPool.Entities()[0]
	view, err := pools.ViewPool.Component(ent)
	if err != nil {
		log.Fatal(err)
	}

	imgRender, err := pools.ImageRenderPool.Component(ent)
	if err != nil {
		log.Fatal(err)
	}

	screen.DrawImage(view.Img, &imgRender.Options)
}
