package game

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	c "strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/systems"
	"strategy-game/util/ecs"
	"strategy-game/util/ecs/psize"
	"strategy-game/util/tile"
	"strategy-game/util/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/furex/v2"
)

func NewGame() *Game {
	g := &Game{
		world:       ecs.CreateWorld(),
		viewOptions: &ebiten.DrawImageOptions{},
		frameCount:  0,
		screen:      screen{width: 640, height: 480},
	}

	// VIEW
	g.view = ebiten.NewImage(g.screen.width, g.screen.height)
	g.viewOptions.GeoM.Scale(2, 2)

	// ECS THINGS
	InitPools(g.world)

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
	InitSystems(g.world)

	return g
}

type Game struct {
	world       *ecs.World
	view        *ebiten.Image
	viewOptions *ebiten.DrawImageOptions
	frameCount  int
	initOnce    sync.Once
	screen      screen
	ui          *furex.View
	// renderWidth  int
	// renderHeight int
}

type screen struct {
	width  int
	height int
}

func (g *Game) FrameCount() int {
	return g.frameCount
}

func (g *Game) View() *ebiten.Image {
	return g.view
}

func (g *Game) ViewOptions() *ebiten.DrawImageOptions {
	return g.viewOptions
}

func (g *Game) RenderWidth() int {
	return g.screen.width
}

func (g *Game) RenderHeight() int {
	return g.screen.height
}

func (g *Game) Update() error {
	g.initOnce.Do(func() {
		g.setupUI()
	})
	g.ui.UpdateWithSize(g.screen.width, g.screen.height)

	g.frameCount++
	g.world.Update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawWorld(g, screen)
	g.ui.Draw(screen)
	// msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	// ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screen.width, g.screen.height
}

func (g *Game) setupUI() {
	g.ui = ui.CreateGameUi(g.screen.width, g.screen.height)
}

// INIT ECS POOLS

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

func InitPools(w *ecs.World) {
	pools.PositionPool = ecs.CreateComponentPool[c.Position](w, psize.Page1024)
	pools.SpritePool = ecs.CreateComponentPool[c.Sprite](w, psize.Page1024)
	pools.MaterialPool = ecs.CreateComponentPool[c.Material](w, psize.Page1024)
	pools.ImageRenderPool = ecs.CreateComponentPool[c.ImageRender](w, psize.Page1024)
	pools.SidePool = ecs.CreateComponentPool[c.Side](w, psize.Page1024)
	pools.OccupiedPool = ecs.CreateComponentPool[c.Occupied](w, psize.Page1024)

	pools.SolidFlag = ecs.CreateFlagPool(w, psize.Page512)
	pools.TileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)
}

// INIT SYSTEMS IN ORDER

func InitSystems(w *ecs.World) {
	ecs.AddSystem(w, &systems.TestSystem{})
}

func DrawWorld(g ecs.GameData, screen *ebiten.Image) {
	frameCount := g.FrameCount()
	view := g.View()
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

		options.Blend = render.Options.Blend
		options.ColorScale = render.Options.ColorScale
		options.Filter = render.Options.Filter
		view.DrawImage(sprite.Sprite.Animate(frameCount), &options)

		// OBJECT
		occupied, err := pools.OccupiedPool.Component(tileEntity)
		if err != nil {
			log.Fatal(err)
		}

		if occupied.ActiveObject != nil {
			objectEntity := *occupied.ActiveObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.DrawImage(img, opt)
		}

		if occupied.StaticObject != nil {
			objectEntity := *occupied.StaticObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.DrawImage(img, opt)
		}

		if occupied.UnitObject != nil { // TODO: UNIT
			objectEntity := *occupied.UnitObject
			img, opt := entityDrawData(objectEntity, frameCount)
			view.DrawImage(img, opt)
		}
	}

	screen.DrawImage(view, g.ViewOptions())
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
