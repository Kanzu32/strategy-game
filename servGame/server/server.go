package server

import (
	"SERV/database"
	"SERV/queue"
	"SERV/terminal"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
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
	Map  string `json:"map"`
}

type Statistics struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	StatName string `json:"statname"`
	Value    int    `json:"value"`
}

func NewServer(database *database.Database) *Server {
	terminal.Log("Создание нового сервера")
	return &Server{
		database:  database,
		userQueue: queue.NewQueue(),
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/api/register", s.handleRegister)
	http.HandleFunc("/api/login", s.handleLogin)
	// http.HandleFunc("/api/map", s.handleGetMap)
	http.HandleFunc("/api/statistics", s.handleGetStatistics)
	// http.HandleFunc("/api/game/create", s.handleCreateGame)
	// http.HandleFunc("/api/game/endturn", s.handleEndTurn)
	// http.HandleFunc("/api/game/state", s.handleGameState)

	terminal.Log("Сервер запущен на порте", port)
	terminal.Log("Доступные эндпоинты:")
	terminal.Log("POST /api/register - Регистрация нового аккаунта")
	terminal.Log("POST /api/login - Авторизация пользователя")
	terminal.Log("GET /api/statistics - Получение статистики")
	// terminal.Log("POST /api/game/create - Создание игровой сессии")
	// terminal.Log("POST /api/game/endturn - Завершение хода")
	// terminal.Log("GET /api/game/state - Получение состояния игры")

	go s.waitForGameConnections()

	terminal.LogFatal(http.ListenAndServe(":"+port, nil))

}

func (s *Server) waitForGameConnections() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		terminal.Log("Ошибка при инициализации tcp соединения, ", err)
		return
	}
	defer listener.Close()
	terminal.Log("Сервер ожидает подключений для создания игровых сессий...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			terminal.Log("Ошибка при попытке подключения пользователя, ", err)
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
		terminal.Log("Ошибка при попытке получения данных от пользователя, ", err)
		return
	}

	var packet Packet

	err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
	if err != nil {
		terminal.Log("Ошибка при попытке десериализации данных от пользователя, ", err)
		return
	}

	terminal.Log("Создано новое соединение с игроком, ", packet.Type, packet.Data)
	switch packet.Type {
	case "GAMESTART":
		s.userQueue.Add(conn)
		terminal.Log("Добавление игрока в очередь, длина очереди: ", s.userQueue.Count())
		if s.userQueue.Count() >= 2 {
			terminal.Log("Созданно новое лобби")
			conn1, err := s.userQueue.Remove()
			if err != nil {
				terminal.Log("Ошибка при удалении пользователя из очереди")
				return
			}
			conn2, err := s.userQueue.Remove()
			if err != nil {
				terminal.Log("Ошибка при удалении пользователя из очереди")
				return
			}
			go s.handleGame(conn1, conn2)
		}
	default:
		terminal.Log("Неизвестный игровой пакет пакет:", packet.Type, packet.Data)
	}
}

func (s *Server) handleGame(connBlue net.Conn, connRed net.Conn) {
	defer connBlue.Close()
	defer connRed.Close()

	// взять в базе данных случайную карту
	mapStr, err := s.database.GetGameMap()
	if err != nil {
		terminal.Log("Ошибка при получении карты из БД")
		return
	}

	// blue ready
	b, err := json.Marshal(GameStartData{"BLUE", *mapStr})
	if err != nil {
		terminal.Log("Ошибка при сериализации GameStartData")
		return
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		terminal.Log("Ошибка при сериализации пакета GameStartData")
		return
	}
	n, err := connBlue.Write(b)
	if n == 0 || err != nil {
		terminal.Log("Ошибка при передаче пакета GameStartData СИНЕЙ команде")
		return
	}

	// red ready
	b, err = json.Marshal(GameStartData{"RED", *mapStr})
	if err != nil {
		terminal.Log("Ошибка при сериализации GameStartData")
		return
	}
	b, err = json.Marshal(Packet{"GAMESTART", string(b)})
	if err != nil {
		terminal.Log("Ошибка при сериализации пакета GameStartData")
		return
	}
	n, err = connRed.Write(b)
	if n == 0 || err != nil {
		terminal.Log("Ошибка при передаче пакета GameStartData КРАСНОЙ команде")
		return
	}

	terminal.Log("Игровая сессия успешно создана")
	blueInputChan := make(chan GameData)
	redInputChan := make(chan GameData)

	blueCheckSumChan := make(chan string)
	blueLastCheckSum := ""

	redCheckSumChan := make(chan string)
	redLastCheckSum := ""

	blueStatisticsChan := make(chan Statistics)
	redStatisticsChan := make(chan Statistics)

	go s.handleClientInput(connBlue, blueInputChan, blueCheckSumChan, blueStatisticsChan) // BLUE
	go s.handleClientInput(connRed, redInputChan, redCheckSumChan, redStatisticsChan)     // RED

	for {
		select {
		case gameData, ok := <-blueInputChan:

			if ok == false {
				terminal.Log("СИНЯЯ команда закрыла соединение, завершение сессии")
				return
			}

			if gameData.Skip == true {
				b, err = json.Marshal(Packet{"SKIP", ""})
				if err != nil {
					terminal.Log("Ошибка при сериализации пакета gameData skip для КРАСНОЙ команды")
					return
				}
				connRed.Write(b)
				terminal.Log("С -> К Пакет: SKIP")
			} else {
				b, err := json.Marshal(gameData)
				if err != nil {
					terminal.Log("Ошибка при сериализации gameData для КРАСНОЙ команды")
					return
				}
				b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
				if err != nil {
					terminal.Log("Ошибка при сериализации пакета gameData для КРАСНОЙ команды")
					return
				}
				connRed.Write(b)
				terminal.Log("С -> К Пакет: ", gameData.UnitID.Id, gameData.TileID.Id)
			}

		case gameData, ok := <-redInputChan:

			if ok == false {
				terminal.Log("КРАСНАЯ команда закрыла соединение, завершение сессии")
				return
			}

			if gameData.Skip == true {
				b, err = json.Marshal(Packet{"SKIP", ""})
				if err != nil {
					terminal.Log("Ошибка при сериализации пакета gameData skip для СИНЕЙ команды")
					return
				}
				connBlue.Write(b)
				terminal.Log("К -> С Пакет: SKIP")
			} else {
				b, err := json.Marshal(gameData)
				if err != nil {
					terminal.Log("Ошибка при сериализации gameData для СИНЕЙ команды")
					return
				}
				b, err = json.Marshal(Packet{"GAMEDATA", string(b)})
				if err != nil {
					terminal.Log("Ошибка при сериализации пакета gameData для СИНЕЙ команды")
					return
				}
				connBlue.Write(b)
				terminal.Log("К -> C Пакет: ", gameData.UnitID.Id, gameData.TileID.Id)
			}

		case checksum, ok := <-blueCheckSumChan:
			if ok == false {
				terminal.Log("СИНЯЯ команда закрыла соединение, завершение сессии")
				return
			}

			if redLastCheckSum != "" && redLastCheckSum == checksum {
				redLastCheckSum = ""
				// КОНТРОЛЬНЫЕ СУММЫ СОВПАЛИ
			} else if redLastCheckSum != "" && redLastCheckSum != checksum {
				terminal.Log("КОНТРОЛЬНЫЕ СУММЫ НЕ СОВПАЛИ!")
			} else {
				blueLastCheckSum = checksum
			}
		case checksum, ok := <-redCheckSumChan:
			if ok == false {
				terminal.Log("КРАСНАЯ команда закрыла соединение, завершение сессии")
				return
			}

			if blueLastCheckSum != "" && blueLastCheckSum == checksum {
				blueLastCheckSum = ""
				// КОНТРОЛЬНЫЕ СУММЫ СОВПАЛИ
				terminal.Log("Контрольные суммы совпали")
			} else if blueLastCheckSum != "" && blueLastCheckSum != checksum {
				terminal.Log("КОНТРОЛЬНЫЕ СУММЫ НЕ СОВПАЛИ!")
			} else {
				redLastCheckSum = checksum
			}

		case statistics, ok := <-blueStatisticsChan:
			if ok == false {
				terminal.Log("СИНЯЯ команда закрыла соединение, завершение сессии")
				return
			}

			err := s.database.UpdateUserStats(statistics.Email, statistics.Password, statistics.StatName, statistics.Value)
			if err != nil {
				terminal.Log(err)
			}
		case statistics, ok := <-redStatisticsChan:
			if ok == false {
				terminal.Log("КРАСНАЯ команда закрыла соединение, завершение сессии")
				return
			}

			err := s.database.UpdateUserStats(statistics.Email, statistics.Password, statistics.StatName, statistics.Value)
			if err != nil {
				terminal.Log(err)
			}
		}
	}
}

func (s *Server) handleClientInput(conn net.Conn, inputChan chan<- GameData, checksumChan chan<- string, statisticsChan chan<- Statistics) {
	var packet Packet
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			terminal.Log("Ошибка при получении данных от пользователя:", err)
			close(inputChan)
			terminal.Log("Канал закрыт")
			return
		}
		err = json.Unmarshal(bytes.Trim(buf, "\x00"), &packet)
		if err != nil {
			terminal.Log("Ошибка при десериализации пакета от пользователя:", err)
			return
		}

		if packet.Type == "GAMEDATA" {
			var data GameData
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				terminal.Log("Ошибка при десериализации данных GameData от пользователя:", err)
				return
			}
			inputChan <- data
		} else if packet.Type == "SKIP" {
			var data GameData
			data.Skip = true
			inputChan <- data
		} else if packet.Type == "CHECKSUM" {
			checksumChan <- packet.Data
		} else if packet.Type == "STATISTICS" {
			var data Statistics
			err := json.Unmarshal([]byte(packet.Data), &data)
			if err != nil {
				terminal.Log("Ошибка при десериализации данных Statistics от пользователя:", err)
				return
			}
			statisticsChan <- data
		} else {
			terminal.Log("Неизвестный тип пакета от пользователя")
		}

	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	terminal.Log("Обработка запроса на регистрацию")
	var req struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		terminal.Log("Ошибка декодирования запроса:", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	err := s.database.Register(req.Email, req.Password)

	if err != nil {
		terminal.Log(err)
		respondWithError(w, http.StatusConflict)
		return
	}

	terminal.Log("Успешная регистрация пользователя", req.Email)
	respondWithJSON(w, map[string]string{"status": "success"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	terminal.Log("Обработка запроса на авторизацию")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		terminal.Log("Ошибка декодирования запроса:", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	err := s.database.Authenticate(req.Email, req.Password)
	if err != nil {
		terminal.Log(err)
		respondWithError(w, http.StatusUnauthorized)
		return
	}

	terminal.Log("Успешная авторизация пользователя", req.Email)
	respondWithJSON(w, map[string]string{"status": "success"})
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, code int) {
	terminal.Log("Возвращаемый код ошибки:", code)
	// respondWithJSON(w, map[string]string{"error": message})
	w.WriteHeader(code)
}

// type StatsRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// 	Damage   int    `json:"damage"`
// 	Cells    int    `json:"cells"`
// }

// TODO delete

// func (s *Server) handleUpdateStats(w http.ResponseWriter, r *http.Request) {
// 	var req StatsRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	err := s.database.UpdateUserStats(req.Email, req.Password, req.Damage, req.Cells)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Stats updated"))
// }

// func (s *Server) handleGetMap(w http.ResponseWriter, r *http.Request) {
// 	sessionID := r.URL.Query().Get("session_id")
// 	if sessionID == "" {
// 		http.Error(w, "session_id is required", http.StatusBadRequest)
// 		return
// 	}

// 	gameMap, err := s.database.GetGameMap(sessionID)
// 	if err != nil {
// 		http.Error(w, "Map not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"map_data": gameMap.MapData,
// 	})
// }

func (s *Server) handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Language string `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		terminal.Log("Ошибка декодирования запроса:", err)
		respondWithError(w, http.StatusBadRequest)
		return
	}

	statistics, err := s.database.GetStatistics(req.Email)
	terminal.Log("Выдача статистики для пользователя", statistics.Email)
	data := ""
	if req.Language == "Rus" {
		data = "Всего нанесено урона: " + strconv.Itoa(statistics.TotalDamage) +
			"\r\nВсего сделано шагов: " + strconv.Itoa(statistics.TotalCells) +
			"\r\nКоличество побед: " + strconv.Itoa(statistics.WinCount)
	} else {
		data = "Total damage: " + strconv.Itoa(statistics.TotalDamage) +
			"\r\nTotal steps: " + strconv.Itoa(statistics.TotalCells) +
			"\r\nWin count: " + strconv.Itoa(statistics.WinCount)
	}

	if err != nil {
		http.Error(w, "Statistics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	// json.NewEncoder(w).Encode(data)
	w.Write([]byte(data))
}
