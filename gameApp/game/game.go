package game

import (
	"encoding/json"
	"log"
	"os"

	c "strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/systems"
	"strategy-game/util/ecs"
	"strategy-game/util/ecs/psize"
	"strategy-game/util/gamedata"
	"strategy-game/util/sprite"
	"strategy-game/util/tile"
	"strategy-game/util/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func NewGame() *Game {
	g := &Game{
		world:      ecs.CreateWorld(),
		viewScale:  2,
		frameCount: 0,
		screen:     screen{width: 640, height: 480},
		ui:         ui.CreateGameUI(),
	}

	// VIEW
	g.view = ebiten.NewImage(g.screen.width, g.screen.height)

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
	world      *ecs.World
	view       *ebiten.Image
	viewScale  int
	frameCount int
	screen     screen
	ui         *ui.GameUI
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

func (g *Game) ViewScale() int {
	return g.viewScale
}

func (g *Game) ViewScaleInc() {
	if g.viewScale != 10 {
		g.viewScale++
		g.ui.MinusButton.Active = true
	}
	if g.viewScale == 10 {
		g.ui.PlusButton.Active = false
	}
}

func (g *Game) ViewScaleDec() {
	if g.viewScale != 1 {
		g.viewScale--
		g.ui.PlusButton.Active = true
	}
	if g.viewScale == 1 {
		g.ui.MinusButton.Active = false
	}
}

func (g *Game) RenderWidth() int {
	return g.screen.width
}

func (g *Game) RenderHeight() int {
	return g.screen.height
}

func (g *Game) handleInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		xPos, yPos := g.mousePos()
		if g.ui.PlusButton.InBounds(xPos, yPos) {
			g.ui.PlusButton.Click(g)
		} else if g.ui.MinusButton.InBounds(xPos, yPos) {
			g.ui.MinusButton.Click(g)
		}
	}
}

func (g *Game) mousePos() (int, int) {
	x, y := ebiten.CursorPosition()
	return x / g.ui.CurrentScale, y / g.ui.CurrentScale
}

func (g *Game) Update() error {
	g.frameCount++
	g.handleInput()
	g.world.Update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawWorld(g, screen)
	// print(screen.Bounds().Dx(), screen.Bounds().Dy())
	g.ui.Draw(screen, g)
	// msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	// ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screen.width = outsideWidth
	g.screen.height = outsideHeight
	return g.screen.width, g.screen.height
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
	utilLayer := tilemap.Layers[2]   // 2 layer (util)

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

		// ##UNITS AND UTILITY##
		if utilLayer.Data[i] > 0 {

			tileData := tilesets.Get(utilLayer.Data[i])
			if tileData.IsWall {
				pools.WallFlag.AddExistingEntity(tileEntity)
			} else if tileData.IsUnit {
				class, team := "", ""
				switch tileData.Class {
				case 0: // SHILDER
					class = "shield"
				case 1: // GLAIVER
					class = "glaive"
				case 2: // MAGE
					class = "mage"
				case 3: // ARCHER
					class = "bow"
				default:
					panic("UNEXPECTED CLASS")
				}

				switch tileData.Team {
				case 1: // BLUE
					team = "blue"
				case 2: // RED
					team = "red"
				default:
					panic("UNEXPECTED TEAM")
				}

				img, _, err := ebitenutil.NewImageFromFile("assets/img/" + team + "-" + class + ".png")
				if err != nil {
					panic(err)
				}
				spr := sprite.NewSprite(img, 16, 16)
				// spr.AddAnimation("default", []sprite.Frame{
				// 	{N: 0, Time: 5000},
				// })
				spr.AddAnimation("idle-down", []sprite.Frame{
					{N: 0, Time: 5000},
				})
				spr.AddAnimation("idle-up", []sprite.Frame{
					{N: 1, Time: 5000},
				})
				spr.AddAnimation("idle-left", []sprite.Frame{
					{N: 2, Time: 5000},
				})
				spr.AddAnimation("idle-right", []sprite.Frame{
					{N: 3, Time: 5000},
				})

				spr.AddAnimation("walk-down", []sprite.Frame{
					{N: 0, Time: 200},
					{N: 4, Time: 200},
					{N: 8, Time: 200},
					{N: 12, Time: 200},
				})
				spr.AddAnimation("walk-up", []sprite.Frame{
					{N: 1, Time: 200},
					{N: 5, Time: 200},
					{N: 9, Time: 200},
					{N: 13, Time: 200},
				})
				spr.AddAnimation("walk-left", []sprite.Frame{
					{N: 2, Time: 200},
					{N: 6, Time: 200},
					{N: 10, Time: 200},
					{N: 14, Time: 200},
				})
				spr.AddAnimation("walk-right", []sprite.Frame{
					{N: 3, Time: 200},
					{N: 7, Time: 200},
					{N: 11, Time: 200},
					{N: 15, Time: 200},
				})

				spr.AddAnimation("action-down", []sprite.Frame{
					{N: 16, Time: 500},
				})
				spr.AddAnimation("action-up", []sprite.Frame{
					{N: 17, Time: 500},
				})
				spr.AddAnimation("action-left", []sprite.Frame{
					{N: 18, Time: 500},
				})
				spr.AddAnimation("action-right", []sprite.Frame{
					{N: 19, Time: 500},
				})

				spr.SetAnimation("idle-down")

				spriteComp := c.Sprite{Sprite: spr}
				unitEntity, err := pools.SpritePool.AddNewEntity(spriteComp)
				if err != nil {
					log.Fatal(err)
				}

				opt := ebiten.DrawImageOptions{}
				opt.GeoM.Translate(0, -float64(tileData.Sprite.Height()-16))
				pools.ImageRenderPool.AddExistingEntity(unitEntity, c.ImageRender{Options: opt})

				positionComp := c.Position{X: i % utilLayer.Width, Y: i / utilLayer.Width}
				pools.PositionPool.AddExistingEntity(unitEntity, positionComp)

				pools.GhostFlag.AddExistingEntity(unitEntity)

				occupiedComp.UnitObject = &unitEntity
			}
		}
		// OCCUPIED by any objects
		pools.OccupiedPool.AddExistingEntity(tileEntity, occupiedComp)
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
	pools.WallFlag = ecs.CreateFlagPool(w, psize.Page1024)
	pools.GhostFlag = ecs.CreateFlagPool(w, psize.Page256)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)
}

// INIT SYSTEMS IN ORDER

func InitSystems(w *ecs.World) {
	ecs.AddSystem(w, &systems.TestSystem{})
}

func DrawWorld(g gamedata.GameData, screen *ebiten.Image) {
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

		if occupied.UnitObject != nil {
			unitEntity := *occupied.UnitObject
			img, opt := entityDrawData(unitEntity, frameCount)
			view.DrawImage(img, opt)
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

	}

	for _, ghostEntity := range pools.GhostFlag.Entities() {
		position, err := pools.PositionPool.Component(ghostEntity)
		if err != nil {
			log.Fatal(err)
		}

		sprite, err := pools.SpritePool.Component(ghostEntity)
		if err != nil {
			log.Fatal(err)
		}

		render, err := pools.ImageRenderPool.Component(ghostEntity)
		if err != nil {
			log.Fatal(err)
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
