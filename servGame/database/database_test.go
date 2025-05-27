package database

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *Database {
	// Используем тестовую MongoDB (можно заменить на in-memory для изолированных тестов)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Используем тестовую базу данных
	db := client.Database("test_game_server")

	// Очищаем коллекции перед тестами
	_, _ = db.Collection("test_users").DeleteMany(context.TODO(), map[string]interface{}{})
	_, _ = db.Collection("test_stats").DeleteMany(context.TODO(), map[string]interface{}{})
	_, _ = db.Collection("test_maps").DeleteMany(context.TODO(), map[string]interface{}{})

	return &Database{
		client:   client,
		accounts: db.Collection("test_users"),
		stats:    db.Collection("test_stats"),
		maps:     db.Collection("test_maps"),
	}
}

func TestRegisterAndAuthenticate(t *testing.T) {
	db := setupTestDB(t)
	email := "test@example.com"
	password := "securepassword"

	// Тест регистрации
	err := db.Register(email, password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Попытка повторной регистрации
	err = db.Register(email, password)
	if err == nil {
		t.Error("Expected error for duplicate registration, got nil")
	}

	// Тест аутентификации
	err = db.Authenticate(email, password)
	if err != nil {
		t.Errorf("Authenticate failed: %v", err)
	}

	// Неверный пароль
	err = db.Authenticate(email, "wrongpassword")
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}

	// Несуществующий пользователь
	err = db.Authenticate("nonexistent@example.com", password)
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
}

func TestUpdateUserStats(t *testing.T) {
	db := setupTestDB(t)
	email := "stats@example.com"
	password := "password123"

	// Сначала регистрируем пользователя
	err := db.Register(email, password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Обновляем статистику
	err = db.UpdateUserStats(email, password, "total_damage", 10)
	if err != nil {
		t.Fatalf("UpdateUserStats failed: %v", err)
	}

	// Проверяем статистику
	stats, err := db.GetStatistics(email)
	if err != nil {
		t.Fatalf("GetStatistics failed: %v", err)
	}

	if stats.TotalDamage != 10 {
		t.Errorf("Expected TotalDamage=10, got %d", stats.TotalDamage)
	}

	// Обновляем еще раз
	err = db.UpdateUserStats(email, password, "total_damage", 5)
	if err != nil {
		t.Fatalf("Second UpdateUserStats failed: %v", err)
	}

	stats, err = db.GetStatistics(email)
	if err != nil {
		t.Fatalf("GetStatistics failed: %v", err)
	}

	if stats.TotalDamage != 15 {
		t.Errorf("Expected TotalDamage=15 after second update, got %d", stats.TotalDamage)
	}
}

func TestGetStatistics_NoStats(t *testing.T) {
	db := setupTestDB(t)
	email := "nostats@example.com"

	// Для пользователя без статистики должен возвращаться нулевой объект
	stats, err := db.GetStatistics(email)
	if err != nil {
		t.Fatalf("GetStatistics failed: %v", err)
	}

	if stats.TotalDamage != 0 || stats.TotalCells != 0 || stats.WinCount != 0 {
		t.Errorf("Expected zero stats, got %+v", stats)
	}
}
