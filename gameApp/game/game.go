package game

import (
	"encoding/json"
	"os"

	c "strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/game/systems"
	"strategy-game/util/classes"
	"strategy-game/util/ecs"
	"strategy-game/util/ecs/psize"
	"strategy-game/util/sprite"
	"strategy-game/util/teams"
	"strategy-game/util/tile"
	"strategy-game/util/turn"
	"strategy-game/util/turn/turnstate"
	"strategy-game/util/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func InitPools(w *ecs.World) {
	pools.PositionPool = ecs.CreateComponentPool[c.Position](w, psize.Page1024)
	pools.SpritePool = ecs.CreateComponentPool[c.Sprite](w, psize.Page1024)
	pools.MaterialPool = ecs.CreateComponentPool[c.Material](w, psize.Page1024)
	pools.ImageRenderPool = ecs.CreateComponentPool[c.ImageRender](w, psize.Page1024)
	pools.SidePool = ecs.CreateComponentPool[c.Side](w, psize.Page1024)
	pools.OccupiedPool = ecs.CreateComponentPool[c.Occupied](w, psize.Page1024)
	pools.TeamPool = ecs.CreateComponentPool[c.Team](w, psize.Page32)
	pools.ClassPool = ecs.CreateComponentPool[c.Class](w, psize.Page16)
	pools.EnergyPool = ecs.CreateComponentPool[c.Energy](w, psize.Page128)
	pools.TweenPool = ecs.CreateComponentPool[c.Tween](w, psize.Page128)

	// pools.SolidFlag = ecs.CreateFlagPool(w, psize.Page512)
	pools.TileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	pools.WallFlag = ecs.CreateFlagPool(w, psize.Page1024)
	pools.SoftFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.UnitFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.GhostFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.ActiveFlag = ecs.CreateFlagPool(w, psize.Page64)
	pools.TargetUnitFlag = ecs.CreateFlagPool(w, psize.Page8)
	pools.TargetObjectFlag = ecs.CreateFlagPool(w, psize.Page128)
	// isFire := ecs.CreateFlagPool(w, psize.Page32)
	// isIce := ecs.CreateFlagPool(w, psize.Page32)
}

// INIT SYSTEMS IN ORDER

func InitSystems(w *ecs.World) {
	ecs.AddRenderSystem(w, &systems.DrawWorldSystem{})
	ecs.AddRenderSystem(w, &systems.DrawGhostsSystem{})

	ecs.AddSystem(w, &systems.TurnSystem{})
	ecs.AddSystem(w, &systems.MarkActiveUnitsSystem{})
	ecs.AddSystem(w, &systems.MarkActiveTilesSystem{})
	ecs.AddSystem(w, &systems.MoveSystem{})

}

func InitStartData(playerTeam teams.Team) {
	singletons.Turn = turn.Turn{CurrentTurn: teams.Blue, PlayerTeam: playerTeam, State: turnstate.Input}

}

func NewGame(playerTeam teams.Team) *Game {
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
	InitStartData(playerTeam)
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
	// keys := []ebiten.Key{}
	// inpututil.AppendJustPressedKeys(keys)
	// keys[0].
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		xPosUI, yPosUI := g.mousePosUIScale()
		// UI
		if g.ui.PlusButton.InBounds(xPosUI, yPosUI) {
			g.ui.PlusButton.Click(g)
			return
		} else if g.ui.MinusButton.InBounds(xPosUI, yPosUI) {
			g.ui.MinusButton.Click(g)
			return
		}

		// ENT

		if singletons.Turn.State != turnstate.Input {
			return
		}

		// клик на активный (active) либо взятый в цель (target object) объект на экране
		activeEntities := ecs.PoolFilter([]ecs.AnyPool{pools.PositionPool, pools.SpritePool}, []ecs.AnyPool{})
		xPosGame, yPosGame := g.mousePosGameScale()
		for _, entity := range activeEntities {
			// неактивные объекты и объекты не взятые в цель игнорируются
			if !pools.ActiveFlag.HasEntity(entity) && !pools.TargetObjectFlag.HasEntity(entity) {
				continue
			}

			position, err := pools.PositionPool.Component(entity)
			if err != nil {
				panic(err)
			}

			sprite, err := pools.SpritePool.Component(entity)
			if err != nil {
				panic(err)
			}

			if position.X*16 < xPosGame && xPosGame < position.X*16+sprite.Sprite.Width() &&
				position.Y*16 < yPosGame && yPosGame < position.Y*16+sprite.Sprite.Height() {

				// объект, взятый в цель, явл. тайлом (выбрать объект в цель для действия)
				if pools.TargetObjectFlag.HasEntity(entity) && pools.TileFlag.HasEntity(entity) {
					singletons.Turn.State = turnstate.Action
					println("muvin")
					return
				}

				// активный объект не являющийся юнитом (выбрать объект в цель для действия)
				if pools.ActiveFlag.HasEntity(entity) && !pools.UnitFlag.HasEntity(entity) {
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetObjectFlag.AddExistingEntity(entity)
					return
				}

				// компонент team есть у всех юнитов (проверка на юнит выше)
				team, err := pools.TeamPool.Component(entity)
				if err != nil {
					panic(err)
				}

				// активный юнит игрока (выбрать его для управления)
				if pools.ActiveFlag.HasEntity(entity) && team.Team == singletons.Turn.PlayerTeam {
					for _, ent := range pools.TargetUnitFlag.Entities() {
						pools.TargetUnitFlag.RemoveEntity(ent)
					}
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetUnitFlag.AddExistingEntity(entity)
					return
				}

				// активный юнит оппонента (выбрать юнит в цель для действия)
				if pools.ActiveFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
					for _, ent := range pools.TargetObjectFlag.Entities() {
						pools.TargetObjectFlag.RemoveEntity(ent)
					}
					pools.TargetObjectFlag.AddExistingEntity(entity)
					return
				}

				if pools.TargetObjectFlag.HasEntity(entity) && team.Team != singletons.Turn.PlayerTeam {
					// атака...
					return
				}
			}
		}
	}
}

func (g *Game) mousePosUIScale() (int, int) {
	x, y := ebiten.CursorPosition()
	return x / g.ui.CurrentScale, y / g.ui.CurrentScale
}

func (g *Game) mousePosGameScale() (int, int) {
	x, y := ebiten.CursorPosition()
	return x / g.viewScale, y / g.viewScale
}

func (g *Game) Update() error {
	g.frameCount++
	g.handleInput()
	g.world.Update(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// DrawWorld(g, screen)
	g.world.Draw(g, screen)
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

func InitTileEntities(tilesets tile.TilesetArray, tilemapFilepath string) {

	// READ TILEMAP

	contents, err := os.ReadFile(tilemapFilepath)
	if err != nil {
		panic(err)
	}

	var tilemap tile.TilemapJSON
	err = json.Unmarshal(contents, &tilemap)
	if err != nil {
		panic(err)
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
			panic(err)
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
				panic(err)
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

				if tileData.IsSoft {
					pools.SoftFlag.AddExistingEntity(objectEntity)
				}

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
				team := teams.Blue
				class := classes.Shield
				switch tileData.Class {
				case 0: // SHILDER
					class = classes.Shield
				case 1: // GLAIVER
					class = classes.Glaive
				case 2: // ROGUE
					class = classes.Knife
				case 3: // ARCHER
					class = classes.Bow
				default:
					panic("UNEXPECTED CLASS")
				}

				switch tileData.Team {
				case 1: // BLUE
					team = teams.Blue
				case 2: // RED
					team = teams.Red
				default:
					panic("UNEXPECTED TEAM")
				}

				//img, _, err := ebitenutil.NewImageFromFile("assets/img/" + team + "-" + class + ".png")
				// print(class.String())

				// TODO SKINS
				img, _, err := ebitenutil.NewImageFromFile("assets/img/" + team.String() + ".png")
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
					panic(err)
				}

				opt := ebiten.DrawImageOptions{}
				opt.GeoM.Translate(0, -float64(tileData.Sprite.Height()-16))
				pools.ImageRenderPool.AddExistingEntity(unitEntity, c.ImageRender{Options: opt})

				positionComp := c.Position{X: i % utilLayer.Width, Y: i / utilLayer.Width}
				pools.PositionPool.AddExistingEntity(unitEntity, positionComp)

				classComp := c.Class{Class: class}
				pools.ClassPool.AddExistingEntity(unitEntity, classComp)

				teamComp := c.Team{Team: team}
				pools.TeamPool.AddExistingEntity(unitEntity, teamComp)

				pools.EnergyPool.AddExistingEntity(unitEntity, c.Energy{Energy: singletons.ClassStats[class].MaxEnergy})

				pools.GhostFlag.AddExistingEntity(unitEntity)

				pools.UnitFlag.AddExistingEntity(unitEntity)
				occupiedComp.UnitObject = &unitEntity
			}
		}
		// OCCUPIED by any objects
		pools.OccupiedPool.AddExistingEntity(tileEntity, occupiedComp)
	}

}
