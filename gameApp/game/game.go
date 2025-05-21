package game

import (
	"encoding/json"
	"os"

	"strategy-game/assets"
	c "strategy-game/game/components"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/game/systems"

	"strategy-game/util/data/classes"
	"strategy-game/util/data/directions"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/sprite"
	"strategy-game/util/data/teams"
	"strategy-game/util/data/turn"
	"strategy-game/util/data/turn/turnstate"
	"strategy-game/util/data/userstatus"
	"strategy-game/util/ecs"
	"strategy-game/util/ecs/psize"
	"strategy-game/util/network"
	"strategy-game/util/sound"
	"strategy-game/util/tile"
	"strategy-game/util/ui"
	"strategy-game/util/ui/uistate"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	pools.HealthPool = ecs.CreateComponentPool[c.Health](w, psize.Page128)
	pools.TweenPool = ecs.CreateComponentPool[c.Tween](w, psize.Page128)
	pools.MovePool = ecs.CreateComponentPool[c.MoveDirection](w, psize.Page128)
	pools.DirectionPool = ecs.CreateComponentPool[c.Direction](w, psize.Page64)
	pools.AttackPool = ecs.CreateComponentPool[c.Attack](w, psize.Page32)
	pools.DamagePool = ecs.CreateComponentPool[c.Damage](w, psize.Page16)
	// pools.StandOnPool = ecs.CreateComponentPool[c.StandOn](w, psize.Page64)

	pools.TileFlag = ecs.CreateFlagPool(w, psize.Page1024)
	pools.WallFlag = ecs.CreateFlagPool(w, psize.Page1024)
	pools.SoftFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.UnitFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.GhostFlag = ecs.CreateFlagPool(w, psize.Page16)
	pools.ActiveFlag = ecs.CreateFlagPool(w, psize.Page64)
	pools.TargetUnitFlag = ecs.CreateFlagPool(w, psize.Page8)
	pools.TargetObjectFlag = ecs.CreateFlagPool(w, psize.Page128)
	pools.DeadFlag = ecs.CreateFlagPool(w, psize.Page16)
}

// INIT SYSTEMS IN ORDER

func InitSystems(w *ecs.World) {
	ecs.AddRenderSystem(w, &systems.DrawWorldSystem{})
	ecs.AddRenderSystem(w, &systems.DrawGhostsSystem{})
	ecs.AddRenderSystem(w, &systems.DrawStatsSystem{})

	ecs.AddSystem(w, &systems.TurnSystem{})
	ecs.AddSystem(w, &systems.MarkActiveUnitsSystem{})
	ecs.AddSystem(w, &systems.MarkActiveTilesSystem{})
	ecs.AddSystem(w, &systems.NetworkSystem{})

	// ecs.AddSystem(w, &systems.TweenMoveSystem{})
	// ecs.AddSystem(w, &systems.UnitMoveSystem{})
	ecs.AddSystem(w, &systems.MoveSystem{})
	ecs.AddSystem(w, &systems.AttackSystem{})
	ecs.AddSystem(w, &systems.DamageSystem{})
	ecs.AddSystem(w, &systems.EnergySystem{})

}

func InitStartData() {
	if singletons.AppState.GameMode == gamemode.Local {
		singletons.Turn = turn.Turn{CurrentTurn: teams.Blue, PlayerTeam: teams.Blue, State: turnstate.Input, IsAttackAllowed: true, IsTurnEnds: false, IsGameEnds: false}
	} else {
		team := <-network.TeamChan
		// println("TEAM:", team)
		if team == teams.Blue {
			singletons.Turn = turn.Turn{CurrentTurn: teams.Blue, PlayerTeam: team, State: turnstate.Input, IsAttackAllowed: true, IsTurnEnds: false, IsGameEnds: false}
		} else {
			singletons.Turn = turn.Turn{CurrentTurn: teams.Blue, PlayerTeam: team, State: turnstate.Wait, IsAttackAllowed: true, IsTurnEnds: false, IsGameEnds: false}
		}
	}
	singletons.View.Scale = singletons.Settings.DefaultGameScale
	singletons.View.ShiftX = 0
	singletons.View.ShiftY = 0
}

func NewGame() *Game {
	g := &Game{
		world: ecs.CreateWorld(),
		ui:    ui.CreateUI(),
	}
	singletons.Render.Height = 640
	singletons.Render.Width = 480
	// VIEW
	singletons.View.Image = ebiten.NewImage(singletons.Render.Width, singletons.Render.Height)
	singletons.View.Scale = 2

	singletons.UserLogin.Status = userstatus.Offline

	sound.Init()
	g.ui.ShowMainMenu()
	// g.ui.ShowLogin()
	return g
}

type Game struct {
	world *ecs.World
	ui    ui.UI
}

func (g *Game) StartGame() {
	// ECS THINGS
	if singletons.AppState.GameMode == gamemode.Online {
		network.StartGameRequest()
	}
	g.world = ecs.CreateWorld()
	InitPools(g.world)

	tilesets := tile.CreateTilesetArray([]string{
		assets.GroundTileset,
		assets.DecalsTileset,
		assets.ActiveObject,
		assets.ObjectsTileset1,
		assets.ObjectsTileset2,
		assets.ObjectsTileset3,
		assets.ObjectsTileset4,
		assets.ObjectsTileset5,
		assets.ObjectsTileset6,
		assets.UtilTileset,
	})

	InitStartData()
	if singletons.AppState.GameMode == gamemode.Local {
		b, err := os.ReadFile(assets.Tilemap)
		if err != nil {
			panic(err)
		}
		singletons.RawMap = string(b)
		InitTileEntities(tilesets)
	} else if singletons.AppState.GameMode == gamemode.Online {
		InitTileEntities(tilesets)
	}

	InitSystems(g.world)
}

// func (g *Game) mousePosGameScale() (int, int) {
// 	x, y := ebiten.CursorPosition()
// 	return x / singletons.View.Scale, y / singletons.View.Scale
// }

func (g *Game) Update() error {
	if singletons.AppState.StateChanged {
		switch singletons.AppState.UIState {
		case uistate.Game:
			g.StartGame()
			g.ui.ShowGameControls()
			// println("Game start!!!!!!")
		case uistate.Main:
			g.ui.ShowMainMenu()
		case uistate.Login:
			g.ui.ShowLogin()
		case uistate.Settings:
			g.ui.ShowSettings()
		case uistate.Statistics:
			g.ui.ShowStatistics()
		case uistate.Results:
			g.ui.ShowGameResult()
		}
		singletons.AppState.StateChanged = false
	}

	g.ui.Update()

	if singletons.AppState.UIState == uistate.Game {
		singletons.FrameCount++

		if singletons.Turn.IsGameEnds == true {
			singletons.AppState.UIState = uistate.Results
			singletons.AppState.StateChanged = true
			network.EndGame()
			singletons.Turn.IsGameEnds = false
		} else {
			g.world.Update()
		}
	}

	sound.RestartMusicIfNeeds()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if singletons.AppState.UIState == uistate.Game {
		g.world.Draw(screen)
	}
	g.ui.Draw(screen)
	// msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	// ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	singletons.Render.Width = outsideWidth
	singletons.Render.Height = outsideHeight
	return singletons.Render.Width, singletons.Render.Height
}

func InitTileEntities(tilesets tile.TilesetArray) {

	// READ TILEMAP

	var tilemap tile.TilemapJSON
	singletons.MapMutex.Lock()

	err := json.Unmarshal([]byte(singletons.RawMap), &tilemap)
	singletons.MapMutex.Unlock()
	if err != nil {
		panic(err)
	}

	groundLayer := tilemap.Layers[0] // 1 layer (ground)
	objectLayer := tilemap.Layers[1] // 2 layer (objects)
	utilLayer := tilemap.Layers[2]   // 2 layer (util)

	singletons.MapSize.Height = 30
	singletons.MapSize.Width = 30

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

				// #1# ИНИЦИАЛИЗАЦИЯ ЮНИТА И ДОБАВЛЕНИЕ НА ТАЙЛ

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
					// class = classes.Bow // TODO add bows
					class = classes.Knife
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

				img, _, err := ebitenutil.NewImageFromFile("assets/img/" + team.String() + "-" + class.String() + ".png")
				// print(class.String())

				// img, _, err := ebitenutil.NewImageFromFile(assets.Characters[team])
				if err != nil {
					panic(err)
				}
				spr := sprite.NewSprite(img, 48, 48)
				// spr.AddAnimation("default", []sprite.Frame{
				// 	{N: 0, Time: 5000},
				// })
				spr.AddAnimation("idle-down", []sprite.Frame{
					{N: 0, Time: 5000},
				})
				spr.AddAnimation("idle-right", []sprite.Frame{
					{N: 11, Time: 5000},
				})
				spr.AddAnimation("idle-up", []sprite.Frame{
					{N: 22, Time: 5000},
				})

				spr.AddAnimation("idle-left", []sprite.Frame{
					{N: 33, Time: 5000},
				})
				spr.AddAnimation("walk-down", []sprite.Frame{
					{N: 0, Time: 100},
					{N: 1, Time: 100},
					{N: 2, Time: 100},
				})
				spr.AddAnimation("walk-right", []sprite.Frame{
					{N: 11, Time: 100},
					{N: 12, Time: 100},
					{N: 13, Time: 100},
				})
				spr.AddAnimation("walk-up", []sprite.Frame{
					{N: 22, Time: 100},
					{N: 23, Time: 100},
					{N: 24, Time: 100},
				})
				spr.AddAnimation("walk-left", []sprite.Frame{
					{N: 33, Time: 100},
					{N: 34, Time: 100},
					{N: 35, Time: 100},
				})

				spr.AddAnimation("hit-down", []sprite.Frame{
					{N: 7, Time: 250},
					{N: 8, Time: 550},
					{N: 9, Time: 5000},
				})
				spr.AddAnimation("hit-right", []sprite.Frame{
					{N: 18, Time: 250},
					{N: 19, Time: 550},
					{N: 20, Time: 5000},
				})
				spr.AddAnimation("hit-up", []sprite.Frame{
					{N: 29, Time: 250},
					{N: 30, Time: 550},
					{N: 31, Time: 5000},
				})
				spr.AddAnimation("hit-left", []sprite.Frame{
					{N: 40, Time: 250},
					{N: 41, Time: 550},
					{N: 42, Time: 5000},
				})

				spr.AddAnimation("dead-down", []sprite.Frame{
					{N: 10, Time: 5000},
				})
				spr.AddAnimation("dead-right", []sprite.Frame{
					{N: 21, Time: 5000},
				})
				spr.AddAnimation("dead-up", []sprite.Frame{
					{N: 32, Time: 5000},
				})
				spr.AddAnimation("dead-left", []sprite.Frame{
					{N: 43, Time: 5000},
				})

				// атака идёт 1000мс не зависимо от анимации спрайта
				if class == classes.Glaive {
					spr.AddAnimation("attack-down", []sprite.Frame{
						{N: 0, Time: 750},
						// {N: 3, Time: 5000},
						{N: 44, Time: 50},
						{N: 45, Time: 150},
						{N: 46, Time: 500},
						// {N: 3, Time: 5000},
					}) // 44 45 46
					spr.AddAnimation("attack-right", []sprite.Frame{
						{N: 11, Time: 750},
						{N: 55, Time: 50},
						{N: 56, Time: 150},
						{N: 57, Time: 500},
					}) // 55 56 57
					spr.AddAnimation("attack-up", []sprite.Frame{
						{N: 22, Time: 750},
						{N: 66, Time: 50},
						{N: 67, Time: 150},
						{N: 68, Time: 500},
					}) // 66 67 68
					spr.AddAnimation("attack-left", []sprite.Frame{
						{N: 33, Time: 750},
						{N: 77, Time: 50},
						{N: 78, Time: 150},
						{N: 79, Time: 500},
					}) // 77 78 79
				} else {
					spr.AddAnimation("attack-down", []sprite.Frame{
						{N: 0, Time: 750},
						{N: 3, Time: 5000},
					})
					spr.AddAnimation("attack-right", []sprite.Frame{
						{N: 11, Time: 750},
						{N: 14, Time: 5000},
					})
					spr.AddAnimation("attack-up", []sprite.Frame{
						{N: 22, Time: 750},
						{N: 25, Time: 5000},
					})
					spr.AddAnimation("attack-left", []sprite.Frame{
						{N: 33, Time: 750},
						{N: 36, Time: 5000},
					})
				}

				spr.SetAnimation("idle-down")

				spriteComp := c.Sprite{Sprite: spr}
				unitEntity, err := pools.SpritePool.AddNewEntity(spriteComp)
				if err != nil {
					panic(err)
				}

				opt := ebiten.DrawImageOptions{}
				opt.GeoM.Translate(-float64(16), -float64(16))
				pools.ImageRenderPool.AddExistingEntity(unitEntity, c.ImageRender{Options: opt})

				positionComp := c.Position{X: i % utilLayer.Width, Y: i / utilLayer.Width}
				pools.PositionPool.AddExistingEntity(unitEntity, positionComp)

				classComp := c.Class{Class: class}
				pools.ClassPool.AddExistingEntity(unitEntity, classComp)

				teamComp := c.Team{Team: team}
				pools.TeamPool.AddExistingEntity(unitEntity, teamComp)

				pools.EnergyPool.AddExistingEntity(unitEntity, c.Energy{Energy: singletons.ClassStats[class].MaxEnergy})

				pools.HealthPool.AddExistingEntity(unitEntity, c.Health{Health: singletons.ClassStats[class].MaxHealth})

				pools.DirectionPool.AddExistingEntity(unitEntity, c.Direction{Direction: directions.Down})

				pools.GhostFlag.AddExistingEntity(unitEntity)

				// Добавление юнита на тайл
				pools.UnitFlag.AddExistingEntity(unitEntity)
				occupiedComp.UnitObject = &unitEntity

				// pools.StandOnPool.AddExistingEntity(unitEntity, c.StandOn{Tile: &tileEntity})
			}
		}
		// OCCUPIED by any objects
		pools.OccupiedPool.AddExistingEntity(tileEntity, occupiedComp)
	}

}
