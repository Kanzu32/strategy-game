package network

import (
	"strategy-game/util/ecs"
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

type GameData struct {
	GameID         uint16
	TargetedUnit   ecs.Entity
	TargetedObject ecs.Entity
}

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

// func CreateServerConnection() ServerConnection {
// 	conn, err := net.Dial("tcp", "127.0.0.1:4545")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return ServerConnection{conn: conn}
// }

// type ServerConnection struct {
// 	conn net.Conn
// }

// func (c *ServerConnection) Close() {
// 	c.conn.Close()
// }
