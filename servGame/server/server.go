package server

import (
	"SERV/database"
	"SERV/queue"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

type Server struct {
	database  *database.Database
	userQueue *queue.Queue
}

type Packet struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Entity struct {
	State   uint8  `json:"state"`
	Id      uint16 `json:"id"`
	Version uint8  `json:"version"`
}

type GameData struct {
	UnitID Entity `json:"unitid"`
	TileID Entity `json:"tileid"`
	Skip   bool   // только на сервере
}

type GameStartData struct {
	Team string `json:"team"`
}

func NewServer(database *database.Database) *Server {
	log.Println("Создание нового сервера")
	return &Server{
		database:  database,
		userQueue: queue.NewQueue(),
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/api/register", s.handleRegister)
	http.HandleFunc("/api/login", s.handleLogin)
	// http.HandleFunc("/api/game/create", s.handleCreateGame)
	// http.HandleFunc("/api/game/endturn", s.handleEndTurn)
	// http.HandleFunc("/api/game/state", s.handleGameState)

	log.Printf("Сервер запущен на порте %s", port)
	log.Println("Доступные эндпоинты:")
	log.Println("POST /api/register - Регистрация нового аккаунта")
	log.Println("POST /api/login - Авторизация пользователя")
	log.Println("POST /api/game/create - Создание игровой сессии")
	log.Println("POST /api/game/endturn - Завершение хода")
	log.Println("GET /api/game/state - Получение состояния игры")

	go s.waitForGameConnections()

	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func (s *Server) waitForGameConnections() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		log.Println("Ошибка при инициализации tcp соединения, ", err)
		return
	}
	defer listener.Close()
	log.Println("Сервер ожидает подключений для создания игровых сессий...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка при попытке подключения пользователя, ", err)
			conn.Close()
			continue
		}
		go s.handleGameConnection(conn) // запускаем горутину для обработки запроса
	}
}

func (s *Server) handleGameConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if n == 0 || err != nil {
		log.Println("Ошибка при попытке получения данных от пользователя, ", err)
		return
	}

	var packet Packet

	err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
	if err != nil {
		log.Println("Ошибка при попытке десериализации данных от пользователя, ", err)
		return
	}

	log.Println("Создано новое соединение с игроком, ", packet.Type, packet.Data)
	switch packet.Type {
	case "GAMESTART":
		s.userQueue.Add(conn)
		log.Println("Добавление игрока в очередь, длина очереди: ", s.userQueue.Count())
		if s.userQueue.Count() >= 2 {
			log.Println("Созданно новое лобби")
			conn1, err := s.userQueue.Remove()
			if err != nil {
				log.Println("Ошибка при удалении пользователя из очереди")
				return
			}
			conn2, err := s.userQueue.Remove()
			if err != nil {
				log.Println("Ошибка при удалении пользователя из очереди")
				return
			}
			go s.handleGame(conn1, conn2)
		}
	default:
		println("DEFAULT: ", packet.Type, " ", packet.Data)
	}
}

func (s *Server) handleGame(connBlue net.Conn, connRed net.Conn) {
	defer connBlue.Close()
	defer connRed.Close()

	// blue ready
	b, err := json.Marshal(GameStartData{"BLUE"})
	if err != nil {
		log.Println("Ошибка при сериализации GameStartData")
		return
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		log.Println("Ошибка при сериализации пакета GameStartData")
		return
	}
	n, err := connBlue.Write(b)
	if n == 0 || err != nil {
		log.Println("Ошибка при передаче пакета GameStartData СИНЕЙ команде")
		return
	}

	// red ready
	b, err = json.Marshal(GameStartData{"RED"})
	if err != nil {
		log.Println("Ошибка при сериализации GameStartData")
		return
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		log.Println("Ошибка при сериализации пакета GameStartData")
		return
	}
	n, err = connRed.Write(b)
	if n == 0 || err != nil {
		log.Println("Ошибка при передаче пакета GameStartData КРАСНОЙ команде")
		return
	}

	log.Println("Игровая сессия успешно создана")
	blueInputChan := make(chan GameData)
	redInputChan := make(chan GameData)

	go s.handleClientInput(connBlue, blueInputChan) // BLUE
	go s.handleClientInput(connRed, redInputChan)   // RED

	for {
		select {
		case gameData, ok := <-blueInputChan:

			if ok == false {
				log.Println("СИНЯЯ команда закрыла соединение, завершение сессии")
				return
			}

			if gameData.Skip == true {
				b, err = json.Marshal(Packet{"SKIP", ""})
				if err != nil {
					log.Println("Ошибка при сериализации пакета gameData skip для КРАСНОЙ команды")
					return
				}
				connRed.Write(b)
				log.Println("С -> К Пакет: SKIP")
			} else {
				b, err := json.Marshal(gameData)
				if err != nil {
					log.Println("Ошибка при сериализации gameData для КРАСНОЙ команды")
					return
				}
				b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
				if err != nil {
					log.Println("Ошибка при сериализации пакета gameData для КРАСНОЙ команды")
					return
				}
				connRed.Write(b)
				log.Println("С -> К Пакет: ", gameData.UnitID.Id, gameData.TileID.Id)
			}

		case gameData, ok := <-redInputChan:

			if ok == false {
				log.Println("КРАСНАЯ команда закрыла соединение, завершение сессии")
				return
			}

			if gameData.Skip == true {
				b, err = json.Marshal(Packet{"SKIP", ""})
				if err != nil {
					log.Println("Ошибка при сериализации пакета gameData skip для СИНЕЙ команды")
					return
				}
				connBlue.Write(b)
				log.Println("К -> С Пакет: SKIP")
			} else {
				b, err := json.Marshal(gameData)
				if err != nil {
					log.Println("Ошибка при сериализации gameData для СИНЕЙ команды")
					return
				}
				b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
				if err != nil {
					log.Println("Ошибка при сериализации пакета gameData для СИНЕЙ команды")
					return
				}
				connBlue.Write(b)
				log.Println("К -> C Пакет: ", gameData.UnitID.Id, gameData.TileID.Id)
			}
		}
	}
}

func (s *Server) handleClientInput(conn net.Conn, inputChan chan<- GameData) {
	var packet Packet
	for {
		println("ждём ввода")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			log.Println("Ошибка при получении данных от пользователя:", err)
			close(inputChan)
			log.Println("Канал закрыт")
			return
		}
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
		if err != nil {
			log.Println("Ошибка при десериализации пакета от пользователя:", err)
			return
		}

		if packet.Type == "GAMEDATA" {
			var data GameData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				log.Println("Ошибка при десериализации данных GameData от пользователя:", err)
				return
			}
			inputChan <- data
		} else if packet.Type == "SKIP" {
			var data GameData
			data.Skip = true
			inputChan <- data
		} else {
			log.Println("Неизвестный тип пакета от пользователя")
		}

	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Обработка запроса на регистрацию")
	var req struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка декодирования запроса: %v", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	account := &database.Account{
		Password: req.Password,
		Email:    req.Email,
	}

	err := s.database.Register(account)

	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusConflict)
		return
	}

	log.Printf("Успешная регистрация пользователя %s", req.Email)
	respondWithJSON(w, map[string]string{"status": "success"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Обработка запроса на авторизацию")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка декодирования запроса: %v", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	err := s.database.Authenticate(req.Email, req.Password)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusUnauthorized)
		return
	}

	log.Printf("Успешная авторизация пользователя %s", req.Email)
	respondWithJSON(w, map[string]string{"status": "success"})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, code int) {
	log.Printf("Возвращаемый код ошибки: %d", code)
	// respondWithJSON(w, map[string]string{"error": message})
	w.WriteHeader(code)
}
