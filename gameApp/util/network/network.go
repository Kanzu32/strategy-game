package network

import (
	"bytes"
	"encoding/json"
	"net"
	"strategy-game/game/pools"
	"strategy-game/util/data/teams"
	"strategy-game/util/ecs"
	"sync"
)

// type ConnectionData struct {
// 	A string `json:"a"`
// 	B string `json:"b"`
// }

// func Test() {
// 	conn, err := net.Dial("tcp", "127.0.0.1:4545")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	enc := gob.NewEncoder(conn)
// 	enc.Encode(ConnectionData{"Test message a", "Test message b"})

// 	fmt.Print("Ответ:")
// 	buff := make([]byte, 1024)
// 	n, err := conn.Read(buff)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Print(string(buff[0:n]))
// 	fmt.Println()
// }

// func SendGameData(data GameData) {
// 	conn, err := net.Dial("tcp", "127.0.0.1:4545")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	enc := gob.NewEncoder(conn)
// 	enc.Encode(data)

// 	fmt.Print("Ответ:")
// 	buff := make([]byte, 1024)
// 	n, err := conn.Read(buff)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Print(string(buff[0:n]))
// 	fmt.Println()
// }

// func SendGameRequest() {
// 	conn, err := net.Dial("tcp", "127.0.0.1:4545")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	conn

// 	fmt.Print("Ответ:")
// 	buff := make([]byte, 1024)
// 	n, err := conn.Read(buff)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Print(string(buff[0:n]))
// 	fmt.Println()
// }

// func StartSession() {
// 	conn, err := net.Dial("tcp", "127.0.0.1:4545")
// 	if err != nil {
// 		panic(err)
// 	}
// }

type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// type MessageData struct {
// 	Message string
// }

type GameData struct {
	UserTeam string     `json:"userteam"`
	UnitID   ecs.Entity `json:"unitid"`
	TileID   ecs.Entity `json:"tileid"`
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

func (s *ServerConnection) SendGameData(userTeam teams.Team, unitID ecs.Entity, tileID ecs.Entity) {
	s.m.Lock()
	defer s.m.Unlock()
	println("start send data")
	teamString := ""
	if userTeam == teams.Blue {
		teamString = "BLUE"
	} else {
		teamString = "RED"
	}

	b, err := json.Marshal(GameData{UserTeam: teamString, UnitID: unitID, TileID: tileID})
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

func (s *ServerConnection) LoginRequest() {
	conn, err := net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// do stuff
}

func (s *ServerConnection) RegisterRequest() {
	conn, err := net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// do stuff
}

func (s *ServerConnection) StatisticsRequest() {
	conn, err := net.Dial("tcp", "127.0.0.1:4545")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// do stuff
}

// type ServerConnection struct {
// 	conn net.Conn
// }

// func (c *ServerConnection) Close() {
// 	c.conn.Close()
// }
