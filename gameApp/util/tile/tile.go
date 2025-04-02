package tile

import (
	"encoding/json"
	"image"
	_ "image/png"
	"log"
	"os"
	"strategy-game/util/data/material"
	"strategy-game/util/data/side"
	"strategy-game/util/data/sprite"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileInfoJSON struct {
	Id         int            `json:"id"`
	Frames     []sprite.Frame `json:"animation"`
	Properties []PropertyJSON `json:"properties"`
}

type PropertyJSON struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

type tilesetJSON struct {
	Width     int            `json:"tilewidth"`
	Height    int            `json:"tileheight"`
	TilesInfo []TileInfoJSON `json:"tiles"`
	TileCount int            `json:"tilecount"`
	ImagePath string         `json:"image"`
	Image     *ebiten.Image
}

type TilesetArray struct {
	Data []tilesetJSON
}

// ID тайлсета с объектами для взаимодействия. (бочки, кактусы, трава...)
const ACTIVE_OBJECTS_TILESET = 2

// ID вспомогательного тайлсета
const UTILITY_TILESET = 9

// ID вспомогательных тайлов
const WALL_ID = 1

// SHIELDER, GLAVER, MAGE, ARCHER
// var BLUE_TEAM = []int{1, 6, 11, 16}
// var RED_TEAM = []int{2, 7, 12, 17}

var BLUE_TEAM = []int{2, 7, 12, 17}
var RED_TEAM = []int{3, 8, 13, 18}

// мягкие объекты (трава)
// var SOFT_OBJECTS = []int{17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 34, 36, 37, 38, 39, 40, 41, 42, 43}
var SOFT_OBJECTS = []int{18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44}

type tileData struct {
	Sprite   sprite.Sprite
	Material material.Material
	Side     side.Side
	IsActive bool
	IsSoft   bool
	IsWall   bool
	IsUnit   bool
	Team     int
	Class    int
}

func (arr *TilesetArray) Get(n int) tileData {
	var target tilesetJSON
	isActive := false
	isSoft := false
	isWall := false
	isUnit := false
	team := 0
	class := 0
	tilesetN := n

	for i, tileset := range arr.Data {
		if tilesetN > tileset.TileCount {
			tilesetN = tilesetN - tileset.TileCount
		} else {
			target = tileset
			if i == ACTIVE_OBJECTS_TILESET {
				isActive = true
				for _, id := range SOFT_OBJECTS {
					if tilesetN == id {
						isSoft = true
						break
					}
				}
			}
			if i == UTILITY_TILESET {
				if tilesetN == WALL_ID {
					isWall = true
					break
				}
				for cl, id := range BLUE_TEAM {
					if tilesetN == id {
						isUnit = true
						team = 1
						class = cl
						break
					}
				}
				for cl, id := range RED_TEAM {
					if tilesetN == id {
						isUnit = true
						team = 2
						class = cl
						break
					}
				}
			}
			break
		}
	}
	id := tilesetN - 1

	spr := sprite.NewSprite(target.Image, target.Width, target.Height)
	var m material.Material
	var s side.Side

	hasDefaultAnimation := false

	for _, info := range target.TilesInfo {
		if info.Id == id {
			if info.Frames != nil {
				spr.AddAnimation("default", info.Frames)
				hasDefaultAnimation = true
			}
			if info.Properties != nil {
				for _, property := range info.Properties {

					if property.Name == "material" {
						switch property.Value {
						case "sand":
							m = material.Sand
							spr.AddAnimation("wet", []sprite.Frame{{N: id + 9, Time: 5000}})
						case "grass":
							m = material.Grass
							spr.AddAnimation("wet", []sprite.Frame{{N: id + 9, Time: 5000}})
						case "water":
							m = material.Water
						}
					}

					if property.Name == "side" {
						switch property.Value {
						case "Up":
							s = side.Up
						case "Down":
							s = side.Down
						case "Right":
							s = side.Right
						case "Left":
							s = side.Left
						case "LeftUp":
							s = side.LeftUp
						case "LeftDown":
							s = side.LeftDown
						case "RightUp":
							s = side.RightUp
						case "RightDown":
							s = side.RightDown
						case "Center":
							s = side.Center
						case "RightCorner":
							s = side.RightCorner
						case "LeftCorner":
							s = side.LeftCorner
						}
					}
				}
			}
		}
	}
	if !hasDefaultAnimation {
		spr.AddAnimation("default", []sprite.Frame{{N: id, Time: 5000}})
	}
	spr.SetAnimation("default")

	return tileData{
		Sprite:   spr,
		Material: m,
		Side:     s,
		IsActive: isActive,
		IsSoft:   isSoft,
		IsWall:   isWall,
		IsUnit:   isUnit,
		Team:     team,
		Class:    class,
	}
}

func CreateTilesetArray(paths []string) TilesetArray {
	tilesetArray := TilesetArray{}
	for _, path := range paths {
		// fmt.Println(path)
		var tileset tilesetJSON
		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(contents, &tileset)

		f, err := os.Open("assets/tiles/tilesets/" + tileset.ImagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		tileset.Image = ebiten.NewImageFromImage(img)

		tilesetArray.Data = append(tilesetArray.Data, tileset)
	}
	return tilesetArray
}
