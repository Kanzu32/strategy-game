package gameemu

import (
	"testing"
)

func TestGameSessionManagement(t *testing.T) {
	t.Run("Создание и базовая работа сессии", func(t *testing.T) {
		gm := NewGameManager()

		// Тест создания сессии
		sessionID, err := gm.CreateSession("player1", "player2")
		if err != nil {
			t.Fatalf("Ошибка создания сессии: %v", err)
		}

		// Проверка состояния новой сессии
		session, err := gm.GetSession("player1")
		if err != nil {
			t.Fatalf("Ошибка получения сессии: %v", err)
		}

		if session.State != "ожидание" {
			t.Errorf("Новая сессия должна быть в состоянии 'ожидание', получили: %s", session.State)
		}

		// Тест старта сессии
		err = gm.StartSession(sessionID)
		if err != nil {
			t.Fatalf("Ошибка старта сессии: %v", err)
		}

		// Проверка состояния после старта
		session, _ = gm.GetSession("player1")
		if session.State != "активна" {
			t.Errorf("После старта сессия должна быть 'активна', получили: %s", session.State)
		}
		if session.Current != "player1" {
			t.Errorf("После старта текущий игрок должен быть player1, получили: %s", session.Current)
		}
	})

	t.Run("Очередность ходов", func(t *testing.T) {
		gm := NewGameManager()
		sessionID, _ := gm.CreateSession("player1", "player2")
		gm.StartSession(sessionID)

		// Первый ход player1
		err := gm.EndTurn(sessionID, "player1")
		if err != nil {
			t.Fatalf("Ошибка при первом ходе: %v", err)
		}

		session, _ := gm.GetSession("player1")
		if session.Current != "player2" {
			t.Errorf("После хода player1 текущий должен быть player2, получили: %s", session.Current)
		}

		// Ход player2
		err = gm.EndTurn(sessionID, "player2")
		if err != nil {
			t.Fatalf("Ошибка при ходе player2: %v", err)
		}

		session, _ = gm.GetSession("player1")
		if session.Current != "player1" {
			t.Errorf("После хода player2 текущий должен быть player1, получили: %s", session.Current)
		}
	})
}

func TestGameFlow(t *testing.T) {
	gm := NewGameManager()
	player1 := "user1"
	player2 := "user2"

	// Создаем и запускаем сессию
	sessionID, err := gm.CreateSession(player1, player2)
	if err != nil {
		t.Fatal("Ошибка создания сессии:", err)
	}

	if err := gm.StartSession(sessionID); err != nil {
		t.Fatal("Ошибка старта сессии:", err)
	}

	t.Run("Корректная последовательность ходов", func(t *testing.T) {
		// Первый ход player1 -> player2
		if err := gm.EndTurn(sessionID, player1); err != nil {
			t.Error("Ошибка при первом ходе player1:", err)
		}

		session, _ := gm.GetSession(player1)
		if session.Current != player2 {
			t.Errorf("После хода player1 текущий должен быть %s, получили %s", player2, session.Current)
		}

		// Ход player2 -> player1
		if err := gm.EndTurn(sessionID, player2); err != nil {
			t.Error("Ошибка при ходе player2:", err)
		}

		session, _ = gm.GetSession(player1)
		if session.Current != player1 {
			t.Errorf("После хода player2 текущий должен быть %s, получили %s", player1, session.Current)
		}
	})

	t.Run("Некорректные ходы", func(t *testing.T) {
		// Сейчас текущий player1
		// Попытка player2 ходить вне очереди
		err := gm.EndTurn(sessionID, player2)
		if err == nil {
			t.Error("Ожидалась ошибка при ходе player2 вне очереди")
		}

		// Проверяем, что текущий игрок не изменился
		session, _ := gm.GetSession(player1)
		if session.Current != player1 {
			t.Error("Текущий игрок не должен меняться при ошибочном ходе")
		}

		// Корректный ход player1
		if err := gm.EndTurn(sessionID, player1); err != nil {
			t.Error("Ошибка при корректном ходе player1:", err)
		}

		// Попытка player1 ходить дважды подряд
		err = gm.EndTurn(sessionID, player1)
		if err == nil {
			t.Error("Ожидалась ошибка при повторном ходе player1")
		}
	})

	t.Run("Попытка хода в неактивной сессии", func(t *testing.T) {
		// Создаем новую неактивную сессию
		gm := NewGameManager()
		sessionID, _ := gm.CreateSession("p1", "p2")

		err := gm.EndTurn(sessionID, "p1")
		if err == nil {
			t.Error("Ожидалась ошибка при ходе в неактивной сессии")
		}
	})
}
