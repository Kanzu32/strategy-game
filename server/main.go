package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"sync"
)

type queue struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []net.Conn
}

func NewQueue() *queue {
	return &queue{sync.Mutex{}, make([]net.Conn, 0)}
}

func (s *queue) Add(v net.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *queue) Remove() (net.Conn, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil, errors.New("empty stack")
	}

	res := s.s[0]
	s.s = s.s[1:l]
	return res, nil
}

func (s *queue) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)
}

var userQueue *queue

func main() {
	listener, err := net.Listen("tcp", ":4545")
	userQueue = NewQueue()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn) // запускаем горутину для обработки запроса
	}
}

// type ConnectionData struct {
// 	A string `json:"a"`
// 	B string `json:"b"`
// }

type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// type MessageData struct {
// 	Message string `json:"message"`
// }

type Entity struct {
	State   uint8  `json:"state"`
	Id      uint16 `json:"id"`
	Version uint8  `json:"version"`
}

type GameData struct {
	UserTeam string `json:"userteam"`
	UnitID   Entity `json:"unitid"`
	TileID   Entity `json:"tileid"`
}

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GameStartData struct {
	Team string `json:"team"`
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if n == 0 || err != nil {
		panic(err)
	}

	var packet Packet

	err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
	if err != nil {
		panic(err)
	}

	println(packet.Type, packet.Data)
	println("readed")
	switch packet.Type {
	case "GAMESTART":
		print("ADD ")
		userQueue.Add(conn)
		println(userQueue.Count())
		if userQueue.Count() >= 2 {
			println("CREATED LOBBY")
			conn1, err := userQueue.Remove()
			if err != nil {
				panic(err)
			}
			conn2, err := userQueue.Remove()
			if err != nil {
				panic(err)
			}
			go handleGame(conn1, conn2)
		}
	default:
		println("DEFAULT: ", packet.Type, " ", packet.Data)
	}
}

func handleGame(connBlue net.Conn, connRed net.Conn) {
	defer connBlue.Close()
	defer connRed.Close()

	// blue ready
	b, err := json.Marshal(GameStartData{"BLUE"})
	if err != nil {
		panic(err)
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		panic(err)
	}
	n, err := connBlue.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}

	// red ready
	b, err = json.Marshal(GameStartData{"RED"})
	if err != nil {
		panic(err)
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		panic(err)
	}
	n, err = connRed.Write(b)
	if n == 0 || err != nil {
		panic(err)
	}

	print("GAME STARTED")
	inputChan := make(chan GameData)
	go handleClientInput(connBlue, inputChan) // BLUE
	go handleClientInput(connRed, inputChan)  // RED

	for {
		gameData := <-inputChan

		if gameData.UserTeam == "BLUE" {
			b, err := json.Marshal(gameData)
			if err != nil {
				panic(err)
			}
			b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
			if err != nil {
				panic(err)
			}
			connRed.Write(b)
			println(gameData.UserTeam, gameData.UnitID.Id, gameData.TileID.Id)

			// b, err = json.Marshal(Packet{"OK", ""})
			// if err != nil {
			// 	panic(err)
			// }
			// connBlue.Write(b)
		} else {
			b, err := json.Marshal(gameData)
			if err != nil {
				panic(err)
			}
			b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
			if err != nil {
				panic(err)
			}
			connBlue.Write(b)
			println(gameData.UserTeam, gameData.UnitID.Id, gameData.TileID.Id)

			// b, err = json.Marshal(Packet{"OK", ""})
			// if err != nil {
			// 	panic(err)
			// }
			// connRed.Write(b)
		}
	}

}

func handleClientInput(conn net.Conn, inputChan chan<- GameData) {
	var packet Packet
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
		if err != nil {
			fmt.Println("Unmarshal input error:", err)
			continue
		}

		if packet.Type == "GAMEDATA" {
			var data GameData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				panic(err)
			}
			inputChan <- data
		} else {
			println("unexpected client input")
		}

	}
}
