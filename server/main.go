package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")

	var conn1 net.Conn
	var conn2 net.Conn

	for {
		conn1, err = listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn1.Close()
			continue
		}
		conn2, err = listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn1.Close()
			continue
		}
		go handleGame(conn1, conn2) // запускаем горутину для обработки запроса
	}
}

type ConnectionData struct {
	A string `json:"a"`
	B string `json:"b"`
}

// обработка подключения
// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	// считываем полученные в запросе данные
// 	dec := gob.NewDecoder(conn)
// 	var source ConnectionData
// 	err := dec.Decode(&source)

// 	if err != nil {
// 		fmt.Println("Decode error:", err)
// 		return
// 	}

// 	fmt.Println(source)

// 	answer := "OK"
// 	conn.Write([]byte(answer))
// }

func handleGame(conn1 net.Conn, conn2 net.Conn) {
	defer conn1.Close()
	defer conn2.Close()

	// dec := gob.NewDecoder(conn1)
	// var source ConnectionData
	// err := dec.Decode(&source)

	// if err != nil {
	// 	fmt.Println("Decode error:", err)
	// 	return
	// }
	buf := make([]byte, 1024)
	conn1.Read(buf)
	conn2.Write(buf)

	fmt.Println("RESEND")

	answer := "OK"
	conn1.Write([]byte(answer))

}
