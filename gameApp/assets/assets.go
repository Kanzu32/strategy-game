package assets

import (
	_ "embed"
	"strategy-game/util/teams"
)

//go:embed ui/font/monogram-extended.ttf
var MonogramTTF []byte

var GroundTileset = "assets/tiles/tilesets/1_ground-tileset.json"
var DecalsTileset = "assets/tiles/tilesets/2_decals-tileset.json"
var ActiveObject = "assets/tiles/tilesets/3_active-objects-tileset.json"
var ObjectsTileset1 = "assets/tiles/tilesets/4_objects1-tileset.json"
var ObjectsTileset2 = "assets/tiles/tilesets/5_objects2-tileset.json"
var ObjectsTileset3 = "assets/tiles/tilesets/6_objects3-tileset.json"
var ObjectsTileset4 = "assets/tiles/tilesets/7_objects4-tileset.json"
var ObjectsTileset5 = "assets/tiles/tilesets/8_objects5-tileset.json"
var ObjectsTileset6 = "assets/tiles/tilesets/9_objects6-tileset.json"
var UtilTileset = "assets/tiles/tilesets/10_util-tileset.json"

var Tilemap = "assets/tiles/tilemaps/tilemap.json"

var Characters = map[teams.Team]string{
	teams.Blue: "assets/img/blue.png",
	teams.Red:  "assets/img/red.png",
}

var MainMenuBackground = "assets/ui/nine_slice/main_background.png"
var MainUIButton = "assets/ui/nine_slice/nine_slice_ui_1.png"
