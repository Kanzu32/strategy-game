package tile

import (
	"strategy-game/sprite"
)

type TilesetJSON struct {
	Width         int              `json:"tilewidth"`
	Height        int              `json:"tileheight"`
	AnimatedTiles []AnimationsJSON `json:"tiles"`
}

type AnimationsJSON struct {
	Id     int            `json:"id"`
	Frames []sprite.Frame `json:"animation"`
}

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

// func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
// 	contents, err := os.ReadFile(filepath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var tilemapJSON TilemapJSON
// 	err = json.Unmarshal(contents, &tilemapJSON)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &tilemapJSON, nil
// }

// func NewTilesetJSON(filepath string) (*TilesetJSON, error) {
// 	contents, err := os.ReadFile(filepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var tilesetJSON TilesetJSON
// 	err = json.Unmarshal(contents, &tilesetJSON)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &tilesetJSON, nil
// }

// func TileEntities(tilemapFilepath string, tilesetFilepath string, image *ebiten.Image, w *ecs.World, ) {

// 	// READ TILEMAP

// 	contents, err := os.ReadFile(tilemapFilepath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var tilemap TilemapJSON
// 	err = json.Unmarshal(contents, &tilemap)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// READ TILESET

// 	contents, err = os.ReadFile(tilesetFilepath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var tileset TilesetJSON
// 	err = json.Unmarshal(contents, &tileset)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	groundLayer := tilemap.Layers[0] // 1 layer (ground)

// 	for i := 0; i <= groundLayer.Height*groundLayer.Width; i++ {

// 	}
// }
