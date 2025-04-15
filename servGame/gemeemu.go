package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

type GameSession struct {
	ID        string    `json:"id"`
	Player1   string    `json:"player1"`
	Player2   string    `json:"player2"`
	Current   string    `json:"current"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	mu        sync.Mutex
}

type GameManager struct {
	sessions map[string]*GameSession
	players  map[string]string
	mu       sync.RWMutex
}

func NewGameManager() *GameManager {
	log.Println("Инициализация менеджера игровых сессий")
	return &GameManager{
		sessions: make(map[string]*GameSession),
		players:  make(map[string]string),
	}
}

func (gm *GameManager) CreateSession(player1, player2 string) (string, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	log.Printf("Попытка создания сессии для игроков %s и %s", player1, player2)

	if _, exists := gm.players[player1]; exists {
		return "", errors.New("игрок 1 уже в другой сессии")
	}
	if _, exists := gm.players[player2]; exists {
		return "", errors.New("игрок 2 уже в другой сессии")
	}

	sessionID := "sess_" + time.Now().Format("20060102_150405")
	session := &GameSession{
		ID:        sessionID,
		Player1:   player1,
		Player2:   player2,
		Current:   player1,
		State:     "ожидание",
		CreatedAt: time.Now(),
	}

	gm.sessions[sessionID] = session
	gm.players[player1] = sessionID
	gm.players[player2] = sessionID

	log.Printf("Создана новая игровая сессия %s", sessionID)
	return sessionID, nil
}

func (gm *GameManager) StartSession(sessionID string) error {
	gm.mu.RLock()
	session, exists := gm.sessions[sessionID]
	gm.mu.RUnlock()

	if !exists {
		return errors.New("сессия не найдена")
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	if session.State != "ожидание" {
		return errors.New("сессия уже начата")
	}

	session.State = "активна"
	log.Printf("Сессия %s начата", sessionID)
	return nil
}

func (gm *GameManager) EndTurn(sessionID, playerID string) error {
	gm.mu.RLock()
	session, exists := gm.sessions[sessionID]
	gm.mu.RUnlock()

	if !exists {
		return errors.New("сессия не найдена")
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	// Проверяем, что игра активна
	if session.State != "активна" {
		return errors.New("игра не активна")
	}

	// Проверяем, что ходит правильный игрок
	if session.Current != playerID {
		return errors.New("не ваш ход")
	}

	// Меняем текущего игрока
	if session.Current == session.Player1 {
		session.Current = session.Player2
	} else {
		session.Current = session.Player1
	}

	log.Printf("Ход передан в сессии %s. Теперь ходит %s", sessionID, session.Current)
	return nil
}

func (gm *GameManager) GetSession(playerID string) (*GameSession, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	sessionID, exists := gm.players[playerID]
	if !exists {
		return nil, errors.New("игрок не в сессии")
	}

	session, exists := gm.sessions[sessionID]
	if !exists {
		return nil, errors.New("сессия не найдена")
	}

	return session, nil
}
