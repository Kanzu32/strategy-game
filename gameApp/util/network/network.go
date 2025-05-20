package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strategy-game/game/pools"
	"strategy-game/game/singletons"
	"strategy-game/util/data/gamemode"
	"strategy-game/util/data/teams"
	"strategy-game/util/ecs"
)

type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type GameData struct {
	UnitID ecs.Entity `json:"unitid"`
	TileID ecs.Entity `json:"tileid"`
}

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GameStartData struct {
	Team string `json:"team"`
	Map  string `json:"map"`
}

type Statistics struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	StatName string `json:"statname"`
	Value    int    `json:"value"`
}

// type ServerConnection struct {
// 	conn     net.Conn
// 	m        sync.Mutex
// 	TeamChan chan teams.Team
// }

var conn net.Conn
var TeamChan chan teams.Team

func StartGameRequest() {
	var err error

	conn, err = net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(Packet{"GAMESTART", "GAME START IT'S UNREADED"})
	print(string(b))
	if err != nil {
		panic(err)
	}
	n, err := conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
	println("game start req")
	TeamChan = make(chan teams.Team)
	go gameResponse()
}

func gameResponse() {
	for {
		buf := make([]byte, 40000) // map size
		n, err := conn.Read(buf)
		if n == 0 || err != nil { // TODO закрываться здесь
			return
		}

		var packet Packet
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
		if err != nil {
			panic(err)
		}

		switch packet.Type {
		case "GAMESTART":
			SendChecksum()
			var data GameStartData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				panic(err)
			}
			if data.Team == "BLUE" {
				TeamChan <- teams.Blue
			} else if data.Team == "RED" {
				TeamChan <- teams.Red
			} else {
				print("wrong team")
			}
			singletons.RawMap = data.Map
		case "GAMEDATA":
			SendChecksum()
			var data GameData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				panic(err)
			}

			_, err = pools.TargetUnitFlag.AddExistingEntity(data.UnitID)
			if err != nil {
				println("Target unit already been marked, its OK")
			}
			pools.TargetObjectFlag.AddExistingEntity(data.TileID)
			println("ACTIVE: ", data.UnitID.Id, data.TileID.Id)
		case "SKIP":
			SendChecksum()
			singletons.Turn.IsTurnEnds = true
			println("SKIP")
		default:
			println("!!!DEFAULT!!! :", packet.Type)
		}
	}

}

func EndGame() {
	if singletons.AppState.GameMode == gamemode.Online {
		conn.Close()
	}
}

func SendGameData(unitID ecs.Entity, tileID ecs.Entity) {
	println("start send data")
	SendChecksum()

	b, err := json.Marshal(GameData{UnitID: unitID, TileID: tileID})
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(Packet{Type: "GAMEDATA", Data: string(b)})
	if err != nil {
		panic(err)
	}
	print("Send: ", string(b))
	n, err := conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
	// s.conn.Read(b) // response
}

func SendSkip() {
	println("start send skip")
	SendChecksum()
	b, err := json.Marshal(Packet{Type: "SKIP", Data: ""})
	if err != nil {
		panic(err)
	}
	print("Send: ", string(b))
	n, err := conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
	// s.conn.Read(b) // response
}

func SendChecksum() {
	fmt.Printf("%x\n", pools.CalcHash())

	b, err := json.Marshal(Packet{Type: "CHECKSUM", Data: fmt.Sprintf("%x", pools.CalcHash())})
	if err != nil {
		panic(err)
	}
	print("Send checksum: ", string(b))
	n, err := conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
}

// statName: "total_damage", "total_cells"
func SendStatistics(statName string, value int) {
	if singletons.AppState.GameMode != gamemode.Online || singletons.Turn.CurrentTurn != singletons.Turn.PlayerTeam {
		return
	}
	b, err := json.Marshal(
		Statistics{Email: singletons.UserLogin.Email,
			Password: singletons.UserLogin.Password,
			StatName: statName,
			Value:    value})
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(Packet{Type: "STATISTICS", Data: string(b)})
	if err != nil {
		panic(err)
	}
	print("Send: ", string(b))
	n, err := conn.Write(b)
	if n == 0 || err != nil {
		return
	}
}

func LoginRequest(email string, password string) int {

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(UserData{email, password})
	resp, err := http.Post("http://127.0.0.1:8080/api/login", "application/json", &buf)
	if err != nil {
		println("ошибка при входе")
		return http.StatusBadRequest
	}

	println(resp.Status)
	if resp.StatusCode == http.StatusUnauthorized {
		println("неверное имя пользователя или пароль")
	} else if resp.StatusCode == http.StatusOK {
		println("good log")
	} else {
		println("ошибка при входе")
	}

	return resp.StatusCode
	// do stuff
}

func RegisterRequest(email string, password string) int {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(UserData{email, password})
	resp, err := http.Post("http://127.0.0.1:8080/api/register", "application/json", &buf)
	if err != nil {
		println("ошибка при регистрации")
		return http.StatusBadRequest
	}

	println(resp.Status)
	if resp.StatusCode == http.StatusConflict {
		println("пользователь уже зарегистрирован")
	} else if resp.StatusCode == http.StatusOK {
		println("good reg")
	} else {
		println("ошибка при регистрации")
	}

	return resp.StatusCode

	// do stuff
}

func StatisticsRequest(email string) string {
	var buf bytes.Buffer
	var req struct {
		Email    string `json:"email"`
		Language string `json:"language"`
	}
	req.Email = email
	req.Language = singletons.Settings.Language
	json.NewEncoder(&buf).Encode(req)
	resp, err := http.Post("http://127.0.0.1:8080/api/statistics", "application/json", &buf)
	if err != nil {
		panic(err)
	}
	b := make([]byte, 1024)
	resp.Body.Read(b)
	return string(bytes.Trim(b, "\x00"))
}
