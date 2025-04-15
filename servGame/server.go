package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	database    *Database
	gameManager *GameManager
}

func NewServer(database *Database) *Server {
	log.Println("Создание нового сервера")
	return &Server{
		database:    database,
		gameManager: NewGameManager(),
	}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/api/register", s.handleRegister)
	http.HandleFunc("/api/login", s.handleLogin)
	http.HandleFunc("/api/game/create", s.handleCreateGame)
	http.HandleFunc("/api/game/endturn", s.handleEndTurn)
	http.HandleFunc("/api/game/state", s.handleGameState)

	log.Printf("Сервер запущен на порту %s", port)
	log.Println("Доступные эндпоинты:")
	log.Println("POST /api/register - Регистрация нового аккаунта")
	log.Println("POST /api/login - Авторизация пользователя")
	log.Println("POST /api/game/create - Создание игровой сессии")
	log.Println("POST /api/game/endturn - Завершение хода")
	log.Println("GET /api/game/state - Получение состояния игры")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Обработка запроса на регистрацию")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка декодирования запроса: %v", err)
		respondWithError(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	account := &Account{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if !s.database.Register(account) {
		log.Printf("Пользователь %s уже существует", req.Username)
		respondWithError(w, "Имя пользователя уже занято", http.StatusConflict)
		return
	}

	log.Printf("Успешная регистрация пользователя %s", req.Username)
	respondWithJSON(w, map[string]string{"status": "success"})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if !s.database.Authenticate(req.Username, req.Password) {
		respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, map[string]string{"status": "success"})
}

func (s *Server) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Player1 string `json:"player1"`
		Player2 string `json:"player2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	sessionID, err := s.gameManager.CreateSession(req.Player1, req.Player2)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.gameManager.StartSession(sessionID); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, map[string]string{
		"session_id": sessionID,
		"status":     "created",
	})
}

func (s *Server) handleEndTurn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
		PlayerID  string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := s.gameManager.EndTurn(req.SessionID, req.PlayerID); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, map[string]string{"status": "success"})
}

func (s *Server) handleGameState(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		respondWithError(w, "player_id parameter is required", http.StatusBadRequest)
		return
	}

	session, err := s.gameManager.GetSession(playerID)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, session)
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	log.Printf("Ошибка: %s (код %d)", message, code)
	respondWithJSON(w, map[string]string{"error": message})
	w.WriteHeader(code)
}
