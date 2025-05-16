package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type UserStats struct {
	Email       string `bson:"email"`
	TotalDamage int    `bson:"total_damage"`
	TotalCells  int    `bson:"total_cells"`
}

type GameMap struct {
	SessionID string      `bson:"session_id"`
	MapData   interface{} `bson:"map_data"`
}

type Database struct {
	client   *mongo.Client
	accounts *mongo.Collection
	stats    *mongo.Collection
	maps     *mongo.Collection
}

func NewDatabase(uri, dbName, usersColl, statsColl, mapsColl string) *Database {
	log.Printf("Подключение к MongoDB по URI: %s", uri)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Не удалось проверить подключение к MongoDB: %v", err)
	}

	log.Printf("Успешное подключение к БД %s", dbName)
	db := client.Database(dbName)

	return &Database{
		client:   client,
		accounts: db.Collection(usersColl),
		stats:    db.Collection(statsColl),
		maps:     db.Collection(mapsColl),
	}
}

func (db *Database) Authenticate(email, password string) error {
	// log.Printf("Аутентификация пользователя: %s", email)
	filter := bson.M{"email": email}
	var result Account
	err := db.accounts.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка поиска пользователя в базе данных: %v", err))
	}

	if result.Password != password {
		// log.Printf("Неверный пароль для пользователя %s", email)
		return errors.New(fmt.Sprintf("Неверный пароль для пользователя %s", email))
	}

	// log.Printf("Успешная аутентификация пользователя %s", email)
	return nil
}

func (db *Database) Register(account *Account) error {
	// log.Printf("Регистрация нового пользователя: %s", account.Email)
	filter := bson.M{"email": account.Email}
	var existing Account
	err := db.accounts.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		// log.Printf("Пользователь %s уже существует в базе данных", account.Email)

		return errors.New(fmt.Sprintf("Пользователь %s уже существует в базе данных", account.Email))
	}

	_, err = db.accounts.InsertOne(context.TODO(), account)
	if err != nil {
		// log.Printf("Ошибка при добавлении пользователя в базу данных: %v", err)
		return errors.New(fmt.Sprintf("Ошибка при добавлении пользователя в базу данных: %v", err))
	}

	// log.Printf("Пользователь %s успешно добавлен в базу данных", account.Email)
	return nil
}

func (db *Database) UpdateUserStats(email, password string, damage, cells int) error {
	var account Account
	err := db.accounts.FindOne(context.TODO(), bson.M{
		"email":    email,
		"password": password,
	}).Decode(&account)

	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	filter := bson.M{"email": email}
	update := bson.M{
		"$inc": bson.M{
			"total_damage": damage,
			"total_cells":  cells,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err = db.stats.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}

	return nil
}

func (db *Database) GetGameMap(sessionID string) (*GameMap, error) {
	var result GameMap
	err := db.maps.FindOne(context.TODO(), bson.M{"session_id": sessionID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
