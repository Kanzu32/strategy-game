package network

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strategy-game/game/pools"
	"strategy-game/util/data/teams"
	"strategy-game/util/ecs"
	"sync"
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
}

type ServerConnection struct {
	conn     net.Conn
	m        sync.Mutex
	TeamChan chan teams.Team
}

func (s *ServerConnection) StartGameRequest() {
	var err error
	s.m.Lock()
	defer s.m.Unlock()
	s.conn, err = net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(Packet{"GAMESTART", "GAME START IT'S UNREADED"})
	print(string(b))
	if err != nil {
		panic(err)
	}
	n, err := s.conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
	println("game start req")
	s.TeamChan = make(chan teams.Team)
	go s.gameResponse()
}

func (s *ServerConnection) gameResponse() {
	for {
		buf := make([]byte, 1024)
		n, err := s.conn.Read(buf)
		if n == 0 || err != nil {
			panic(err)
		}

		var packet Packet
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
		if err != nil {
			panic(err)
		}

		switch packet.Type {
		case "GAMESTART":
			var data GameStartData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				panic(err)
			}
			if data.Team == "BLUE" {
				s.TeamChan <- teams.Blue
			} else if data.Team == "RED" {
				s.TeamChan <- teams.Red
			} else {
				print("wrong team")
			}
		case "GAMEDATA":
			var data GameData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				panic(err)
			}

			pools.TargetUnitFlag.AddExistingEntity(data.UnitID)
			pools.TargetObjectFlag.AddExistingEntity(data.TileID)
			println("ACTIVE: ", data.UnitID.Id, data.TileID.Id)
		default:
			print(packet.Type)
		}
	}

}

func (s *ServerConnection) SendGameData(unitID ecs.Entity, tileID ecs.Entity) {
	s.m.Lock()
	defer s.m.Unlock()
	println("start send data")

	b, err := json.Marshal(GameData{UnitID: unitID, TileID: tileID})
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(Packet{Type: "GAMEDATA", Data: string(b)})
	if err != nil {
		panic(err)
	}
	print("Send: ", string(b))
	n, err := s.conn.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}
	// s.conn.Read(b) // response
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

func StatisticsRequest(email string, password string) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(UserData{email, password})
	resp, err := http.Post("http://127.0.0.1:8080/api/register", "application/json", &buf)
	if err != nil {
		panic(err)
	}

	println(resp)

	// do stuff
}
