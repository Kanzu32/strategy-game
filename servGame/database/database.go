package database

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	Email string `bson:"email"`
	Hash  string `bson:"hash"`
}

type UserStats struct {
	Email       string `bson:"email"`
	TotalDamage int    `bson:"total_damage"`
	TotalCells  int    `bson:"total_cells"`
	WinCount    int    `bson:"win_count"`
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
	hash := hashPassword(password)
	filter := bson.M{"email": email}
	var result Account
	err := db.accounts.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка поиска пользователя в базе данных: %v", err))
	}

	if result.Hash != hash {
		// log.Printf("Неверный пароль для пользователя %s", email)
		return errors.New(fmt.Sprintf("Неверный пароль для пользователя %s", email))
	}

	// log.Printf("Успешная аутентификация пользователя %s", email)
	return nil
}

func (db *Database) Register(email, password string) error {
	// log.Printf("Регистрация нового пользователя: %s", account.Email)
	filter := bson.M{"email": email}
	var existing Account
	err := db.accounts.FindOne(context.TODO(), filter).Decode(&existing)
	if err == nil {
		// log.Printf("Пользователь %s уже существует в базе данных", account.Email)

		return errors.New(fmt.Sprintf("Пользователь %s уже существует в базе данных", email))
	}

	_, err = db.accounts.InsertOne(context.TODO(), Account{Hash: hashPassword(password), Email: email})
	if err != nil {
		// log.Printf("Ошибка при добавлении пользователя в базу данных: %v", err)
		return errors.New(fmt.Sprintf("Ошибка при добавлении пользователя в базу данных: %v", err))
	}

	// log.Printf("Пользователь %s успешно добавлен в базу данных", account.Email)
	return nil
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func (db *Database) UpdateUserStats(email, password, statName string, value int) error {
	var account Account
	err := db.accounts.FindOne(context.TODO(), bson.M{
		"email": email,
		"hash":  hashPassword(password),
	}).Decode(&account)

	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	filter := bson.M{"email": email}

	update := bson.M{
		"$inc": bson.M{
			statName: value,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err = db.stats.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}

	return nil
}

func (db *Database) GetGameMap() (*string, error) {
	var result bson.M

	num, err := db.maps.CountDocuments(context.TODO(), bson.D{})
	println(num)
	if err != nil {
		return nil, err
	}

	opts := options.FindOne().SetSkip(rand.Int64N(num))
	err = db.maps.FindOne(context.TODO(), bson.D{}, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	b, err := bson.MarshalExtJSON(result, false, false)
	if err != nil {
		return nil, err
	}

	s := string(b)
	return &s, nil
}

func (db *Database) GetStatistics(email string) (*UserStats, error) {
	var result UserStats
	err := db.stats.FindOne(context.TODO(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		result = UserStats{"", 0, 0, 0}
	}
	return &result, nil
}
